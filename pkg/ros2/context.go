/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

/*
#include <rcl/context.h>
#include <rcl/init_options.h>
#include <rcl/init.h>
#include <rcl/time.h>
#include <rcutils/allocator.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"unsafe"

	"github.com/hashicorp/go-multierror"
)

type rosID uint64

func (r *rosID) GetID() uint64 {
	return uint64(*r)
}

func (r *rosID) SetID(x uint64) {
	*r = rosID(x)
}

type rosResource interface {
	io.Closer
	GetID() uint64
	SetID(x uint64)
}

// rosResourceStore manages ROS resources. When Close is called, all resources in
// the store are Closed. The zero value is ready for use.
type rosResourceStore struct {
	resources map[uint64]rosResource
	idCounter uint64
}

func (s *rosResourceStore) addResource(r rosResource) {
	if s.resources == nil {
		s.resources = make(map[uint64]rosResource)
		// The counter starts at one to allow removing zero-initialized
		// resources.
		s.idCounter = 1
	}
	r.SetID(s.idCounter)
	s.resources[s.idCounter] = r
	s.idCounter++
}

func (s *rosResourceStore) removeResource(r rosResource) {
	delete(s.resources, r.GetID())
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

	rcl_allocator_t    *C.rcutils_allocator_t
	rcl_context_t      *C.rcl_context_t
	Clock              *Clock
	rcl_init_options_t *C.rcl_init_options_t

	rosResourceStore
}

/*
NewRCLContext initializes a new RCL context.

parent can be nil, a new context.Background is created
clockType can be nil, then no clock is initialized, you can later initialize it with NewClock()
rclArgs can be nil
*/
func NewContext(wg *sync.WaitGroup, clockType Rcl_clock_type_t, rclArgs *RCLArgs) (ctx *Context, err error) {
	ctx = &Context{WG: wg}
	defer onErr(&err, func() { ctx.Close() })

	if err = rclInit(rclArgs, ctx); err != nil {
		return nil, err
	}

	if wg == nil {
		ctx.WG = &sync.WaitGroup{}
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
		return errors.New("tried to close a closed Context")
	}
	c.WG.Wait() // Wait for gothreads to quit, before GC:ing. Otherwise a ton of null-pointers await.

	var errs *multierror.Error
	errs = multierror.Append(errs, c.rosResourceStore.Close())

	var rc C.rcl_ret_t
	rc = C.rcl_init_options_fini(c.rcl_init_options_t)
	if rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, fmt.Sprintf("C.rcl_init_options_fini(%+v)", c.rcl_init_options_t)))
	} else {
		c.rcl_init_options_t = nil
	}
	if rc = C.rcl_shutdown(c.rcl_context_t); rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, fmt.Sprintf("C.rcl_shutdown(%+v)", c.rcl_context_t)))
	} else if rc = C.rcl_context_fini(c.rcl_context_t); rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, ErrorsCastC(rc, "rcl_context_fini failed"))
	}
	C.free(unsafe.Pointer(c.rcl_allocator_t))
	c.rcl_allocator_t = nil

	c.WG = nil
	return errs.ErrorOrNil()
}
