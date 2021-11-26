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
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
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

var (
	defaultContext   *Context
	initNotCalledErr = errors.New("Init has not been called")
)

// DefaultContext returns the global default context or nil if Init has not yet
// been called.
func DefaultContext() *Context { return defaultContext }

// Init initializes the global default context and logging system if they have
// not been initalized yet. Calling Init multiple times after a successful
// (returning nil) call is a no-op.
func Init(args *Args) (err error) {
	if defaultContext == nil {
		defaultContext, err = NewContext(0, args)
	}
	return
}

// Uninit uninitializes the default context if it has been initialized. Calling
// Uninit multiple times without calling Init in between the calls is a no-op.
// Uninit should be called before program termination if Init has been called
// successfully.
func Uninit() (err error) {
	if defaultContext != nil {
		err = defaultContext.Close()
		defaultContext = nil
	}
	return
}

// Spin starts and waits for all ROS resources in DefaultContext() that need
// waiting such as nodes and subscriptions. Spin returns when an error occurs or
// ctx is canceled.
func Spin(ctx context.Context) error {
	if defaultContext == nil {
		return initNotCalledErr
	}
	return defaultContext.Spin(ctx)
}

// Context manages resources for a set of RCL entities.
type Context struct {
	rcl_allocator_t *C.rcutils_allocator_t
	rcl_context_t   *C.rcl_context_t
	defaultClock    *Clock
	clock           *Clock

	rosResourceStore
}

/*
NewContext initializes a new RCL context.

clockType is the type of the default clock created for the Context. If 0,
ClockTypeROSTime is used.

A nil rclArgs is treated as en empty argument list.

If logging has not yet been initialized, NewContext will initialize it
automatically using rclArgs for logging configuration.
*/
func NewContext(clockType ClockType, rclArgs *Args) (ctx *Context, err error) {
	ctx = &Context{}
	defer onErr(&err, ctx.Close)

	if rclArgs == nil {
		rclArgs, _, err = ParseArgs(nil)
		if err != nil {
			return nil, err
		}
	}

	ctx.rcl_allocator_t = (*C.rcl_allocator_t)(C.malloc(C.sizeof_rcl_allocator_t))
	*ctx.rcl_allocator_t = C.rcl_get_default_allocator()
	ctx.rcl_context_t = (*C.rcl_context_t)(C.malloc(C.sizeof_rcl_context_t))
	*ctx.rcl_context_t = C.rcl_get_zero_initialized_context()

	if err := rclInitLogging(rclArgs, false); err != nil {
		return nil, err
	}

	rcl_init_options_t := C.rcl_get_zero_initialized_init_options()
	rc := C.rcl_init_options_init(&rcl_init_options_t, *ctx.rcl_allocator_t)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}

	rc = C.rcl_init(rclArgs.argc(), rclArgs.argv(), &rcl_init_options_t, ctx.rcl_context_t)
	runtime.KeepAlive(rclArgs)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}

	if clockType == 0 {
		clockType = ClockTypeROSTime
	}
	ctx.defaultClock, err = ctx.NewClock(clockType)
	if err != nil {
		return nil, err
	}
	ctx.clock = ctx.defaultClock

	return ctx, nil
}

func (c *Context) Close() error {
	if c.rcl_context_t == nil && c.rcl_allocator_t == nil {
		return closeErr("context")
	}
	errs := multierror.Append(c.rosResourceStore.Close())
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
	return errs.ErrorOrNil()
}

func (c *Context) Clock() *Clock {
	return c.clock
}

func (c *Context) SetClock(newClock *Clock) {
	if newClock == nil {
		c.clock = c.defaultClock
	} else {
		c.clock = newClock
	}
}

// Spin starts and waits for all ROS resources in the context that need waiting
// such as nodes and subscriptions. Spin returns when an error occurs or ctx is
// canceled.
func (c *Context) Spin(ctx context.Context) error {
	ws, err := c.NewWaitSet()
	if err != nil {
		return spinErr("context", err)
	}
	defer ws.Close()
	ws.addResources(&c.rosResourceStore)
	return spinErr("context", ws.Run(ctx))
}
