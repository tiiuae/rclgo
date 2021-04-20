/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

import (
	"context"
	"sync"
)

/*
RCLContext has a key rclContextImpl which points to the RCLContextImpl
*/
type RCLContext context.Context

type rclContextKey string // context.WithValue best practice?
var rclContextImplKey rclContextKey = "rclContextImpl"

func GetRCLEntities(ctx RCLContext) *rclEntityWrapper {
	return ctx.Value(rclContextImplKey).(*RCLContextImpl).rclEntityWrapper
}
func GetRCLContextImpl(ctx RCLContext) *RCLContextImpl {
	return ctx.Value(rclContextImplKey).(*RCLContextImpl)
}

/*
RCLContextImpl Contains the lifecycle handling for a given portion of RCL entities.
It is wrapped inside the RCLContext (context.Context) because there is no way to fill a single context with parallel values.
*/
type RCLContextImpl struct {
	Quit             func() // call this to make the context stop.
	WG               *sync.WaitGroup
	rclEntityWrapper *rclEntityWrapper
}

/*
NewRCLContext initializes a new RCL context.

parent can be nil, a new context.Background is created
clockType can be nil, then no clock is initialized, you can later initialize it with NewClock()
rclArgs can be nil
*/
func NewRCLContext(parent context.Context, wg *sync.WaitGroup, clockType Rcl_clock_type_t, rclArgs *RCLArgs) (RCLContext, RCLError) {
	rclEntities, rclError := rclInit(rclArgs)
	if rclError != nil {
		return nil, rclError
	}

	if parent == nil {
		parent = context.Background()
	}

	var quitFunc context.CancelFunc
	parent, quitFunc = context.WithCancel(parent)
	rclContextImpl := RCLContextImpl{
		Quit:             quitFunc,
		WG:               wg,
		rclEntityWrapper: rclEntities,
	}
	if wg == nil {
		rclContextImpl.WG = &sync.WaitGroup{}
	}

	newCtx := (RCLContext)(context.WithValue(parent, rclContextImplKey, &rclContextImpl))

	if clockType != 0 {
		_, err := NewClock(newCtx, clockType)
		if err != nil {
			return newCtx, err
		}
	}

	return newCtx, nil
}

/*
NewRCLContextChild TODO:
- Example usage of nested contexts to init ROS2 and then create nodes etc for a nested context.
- Cleanup partially one context at a time.
*/
func NewRCLContextChild(parent context.Context) (*RCLContext, RCLError) {
	return nil, nil
}

func RCLContextFini(rclContext RCLContext) *RCLErrors {
	var errs *RCLErrors
	GetRCLContextImpl(rclContext).Quit()
	GetRCLContextImpl(rclContext).WG.Wait() // Wait for gothreads to quit, before GC:ing. Otherwise a ton of null-pointers await.
	ent := GetRCLEntities(rclContext)

	for o := ent.WaitSets.Front(); o != nil; o = o.Next() {
		err := o.Value.(*WaitSet).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			ent.WaitSets.Remove(o)
		}
	}
	for o := ent.Publishers.Front(); o != nil; o = o.Next() {
		err := o.Value.(*Publisher).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			ent.Publishers.Remove(o)
		}
	}
	for o := ent.Subscriptions.Front(); o != nil; o = o.Next() {
		err := o.Value.(*Subscription).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			ent.Subscriptions.Remove(o)
		}
	}
	if ent.Clock != nil {
		err := ent.Clock.Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			ent.Clock = nil
		}
	}

	return errs
}
