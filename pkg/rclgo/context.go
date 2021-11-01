/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#include <rcl/context.h>
#include <rcl/init_options.h>
#include <rcl/init.h>
#include <rcl/time.h>
#include <rcutils/allocator.h>
*/
import "C"

import (
	"fmt"
	"io"
	"sync"
	"unsafe"

	"github.com/hashicorp/go-multierror"
)

type rosID uint64

func (r *rosID) getID() uint64 {
	return uint64(*r)
}

func (r *rosID) setID(x uint64) {
	*r = rosID(x)
}

type rosResource interface {
	io.Closer
	getID() uint64
	setID(x uint64)
}

// rosResourceStore manages ROS resources. When Close is called, all resources in
// the store are Closed. The zero value is ready for use.
type rosResourceStore struct {
	mutex     sync.Mutex
	resources map[uint64]rosResource
	idCounter uint64
}

func (s *rosResourceStore) addResource(r rosResource) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.resources == nil {
		s.resources = make(map[uint64]rosResource)
		// The counter starts at one to allow removing zero-initialized
		// resources.
		s.idCounter = 1
	}
	r.setID(s.idCounter)
	s.resources[s.idCounter] = r
	s.idCounter++
}

func (s *rosResourceStore) removeResource(r rosResource) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.resources, r.getID())
}

func (s *rosResourceStore) Close() error {
	var err *multierror.Error
	for _, r := range s.resources {
		err = multierror.Append(err, r.Close())
	}
	return err.ErrorOrNil()
}

// Context manages resources for a set of RCL entities.
type Context struct {
	WG *sync.WaitGroup

	rcl_allocator_t *C.rcutils_allocator_t
	rcl_context_t   *C.rcl_context_t
	Clock           *Clock

	rosResourceStore
}

/*
NewContext initializes a new RCL context.

If clockType == 0, no clock is created. You can always create a clock by calling
NewClock.

A nil rclArgs is treated as en empty argument list.

If logging has not yet been initialized, NewContext will initialize it
automatically using rclArgs for logging configuration.
*/
func NewContext(wg *sync.WaitGroup, clockType ClockType, rclArgs *Args) (ctx *Context, err error) {
	if wg == nil {
		wg = &sync.WaitGroup{}
	}
	ctx = &Context{WG: wg}
	defer onErr(&err, ctx.Close)

	if err = rclInit(rclArgs, ctx); err != nil {
		return nil, err
	}

	if clockType != 0 {
		ctx.Clock, err = ctx.NewClock(clockType)
		if err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func (c *Context) Close() error {
	if c.WG == nil {
		return closeErr("context")
	}
	c.WG.Wait() // Wait for gothreads to quit, before GC:ing. Otherwise a ton of null-pointers await.

	var errs *multierror.Error
	errs = multierror.Append(errs, c.rosResourceStore.Close())

	if c.rcl_context_t != nil {
		if rc := C.rcl_shutdown(c.rcl_context_t); rc != C.RCL_RET_OK {
			errs = multierror.Append(errs, errorsCastC(rc, fmt.Sprintf("C.rcl_shutdown(%+v)", c.rcl_context_t)))
		} else if rc := C.rcl_context_fini(c.rcl_context_t); rc != C.RCL_RET_OK {
			errs = multierror.Append(errs, errorsCastC(rc, "rcl_context_fini failed"))
		}
		C.free(unsafe.Pointer(c.rcl_context_t))
		c.rcl_context_t = nil
	}
	if c.rcl_allocator_t != nil {
		C.free(unsafe.Pointer(c.rcl_allocator_t))
		c.rcl_allocator_t = nil
	}

	c.WG = nil
	return errs.ErrorOrNil()
}
