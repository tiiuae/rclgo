/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

import (
	"sync"
)

// Context manages resources for a set of RCL entities.
type Context struct {
	WG       *sync.WaitGroup
	entities *rclEntityWrapper
}

/*
NewRCLContext initializes a new RCL context.

parent can be nil, a new context.Background is created
clockType can be nil, then no clock is initialized, you can later initialize it with NewClock()
rclArgs can be nil
*/
func NewContext(wg *sync.WaitGroup, clockType Rcl_clock_type_t, rclArgs *RCLArgs) (*Context, RCLError) {
	rclEntities, rclError := rclInit(rclArgs)
	if rclError != nil {
		return nil, rclError
	}

	ctx := &Context{
		WG:       wg,
		entities: rclEntities,
	}
	if wg == nil {
		ctx.WG = &sync.WaitGroup{}
	}

	if clockType != 0 {
		_, err := ctx.NewClock(clockType)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (c *Context) Close() *RCLErrors {
	var errs *RCLErrors
	c.WG.Wait() // Wait for gothreads to quit, before GC:ing. Otherwise a ton of null-pointers await.

	for o := c.entities.WaitSets.Front(); o != nil; o = o.Next() {
		err := o.Value.(*WaitSet).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			c.entities.WaitSets.Remove(o)
		}
	}
	for o := c.entities.Publishers.Front(); o != nil; o = o.Next() {
		err := o.Value.(*Publisher).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			c.entities.Publishers.Remove(o)
		}
	}
	for o := c.entities.Subscriptions.Front(); o != nil; o = o.Next() {
		err := o.Value.(*Subscription).Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			c.entities.Subscriptions.Remove(o)
		}
	}
	if c.entities.Clock != nil {
		err := c.entities.Clock.Fini()
		if err != nil {
			errs = RCLErrorsPut(errs, err)
		} else {
			c.entities.Clock = nil
		}
	}

	return errs
}
