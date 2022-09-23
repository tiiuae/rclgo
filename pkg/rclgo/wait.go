package rclgo

/*
#include <rcl/wait.h>
#include <rcl_action/wait.h>
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/hashicorp/go-multierror"
)

type singleUse atomic.Bool

func (s *singleUse) reserve() bool {
	return (*atomic.Bool)(s).CompareAndSwap(false, true)
}

func (s *singleUse) release() {
	(*atomic.Bool)(s).Store(false)
}

type WaitSet struct {
	rosID
	Subscriptions   []*Subscription
	Timers          []*Timer
	Services        []*Service
	Clients         []*Client
	ActionClients   []*ActionClient
	ActionServers   []*ActionServer
	guardConditions []*guardCondition
	rcl_wait_set_t  C.rcl_wait_set_t
	cancelWait      *guardCondition
	context         *Context
}

func NewWaitSet() (*WaitSet, error) {
	if defaultContext == nil {
		return nil, errInitNotCalled
	}
	return defaultContext.NewWaitSet()
}

func (c *Context) NewWaitSet() (ws *WaitSet, err error) {
	const (
		subscriptionsCount   = 0
		guardConditionsCount = 0
		timersCount          = 0
		clientsCount         = 0
		servicesCount        = 0
		eventsCount          = 0
	)
	ws = &WaitSet{
		context:        c,
		Subscriptions:  []*Subscription{},
		Timers:         []*Timer{},
		Services:       []*Service{},
		Clients:        []*Client{},
		rcl_wait_set_t: C.rcl_get_zero_initialized_wait_set(),
	}
	defer onErr(&err, ws.Close)
	var rc C.rcl_ret_t = C.rcl_wait_set_init(
		&ws.rcl_wait_set_t,
		subscriptionsCount,
		guardConditionsCount,
		timersCount,
		clientsCount,
		servicesCount,
		eventsCount,
		c.rcl_context_t,
		*c.rcl_allocator_t,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCast(rc)
	}
	ws.cancelWait, err = c.newGuardCondition()
	if err != nil {
		return nil, err
	}
	ws.addGuardConditions(ws.cancelWait)
	c.addResource(ws)
	return ws, nil
}

// Context returns the context s belongs to.
func (s *WaitSet) Context() *Context {
	return s.context
}

func (w *WaitSet) AddSubscriptions(subs ...*Subscription) {
	w.Subscriptions = append(w.Subscriptions, subs...)
}

func (w *WaitSet) AddTimers(timers ...*Timer) {
	w.Timers = append(w.Timers, timers...)
}

func (w *WaitSet) AddServices(services ...*Service) {
	w.Services = append(w.Services, services...)
}

func (w *WaitSet) AddClients(clients ...*Client) {
	w.Clients = append(w.Clients, clients...)
}

func (w *WaitSet) AddActionServers(servers ...*ActionServer) {
	w.ActionServers = append(w.ActionServers, servers...)
}

func (w *WaitSet) AddActionClients(clients ...*ActionClient) {
	w.ActionClients = append(w.ActionClients, clients...)
}

func (w *WaitSet) addGuardConditions(guardConditions ...*guardCondition) {
	w.guardConditions = append(w.guardConditions, guardConditions...)
}

func (w *WaitSet) addResources(res *rosResourceStore) {
	for _, res := range res.resources {
		switch res := res.(type) {
		case *Subscription:
			w.AddSubscriptions(res)
		case *Timer:
			w.AddTimers(res)
		case *Service:
			w.AddServices(res)
		case *Client:
			w.AddClients(res)
		case *ActionServer:
			w.AddActionServers(res)
		case *ActionClient:
			w.AddActionClients(res)
		case *guardCondition: // Guard conditions are handled specially
		case *Node:
			w.addResources(&res.rosResourceStore)
		}
	}
}

/*
Run causes the current goroutine to block on this given WaitSet.
WaitSet executes the given timers and subscriptions and calls their callbacks on new events.
*/
func (self *WaitSet) Run(ctx context.Context) (err error) {
	for _, subscription := range self.Subscriptions {
		if subscription.waitable.reserve() {
			defer subscription.waitable.release()
		}
	}
	for _, timer := range self.Timers {
		if timer.waitable.reserve() {
			defer timer.waitable.release()
		}
	}
	for _, service := range self.Services {
		if service.waitable.reserve() {
			defer service.waitable.release()
		}
	}
	for _, client := range self.Clients {
		if client.waitable.reserve() {
			defer client.waitable.release()
		}
	}
	for _, actionClient := range self.ActionClients {
		if actionClient.waitable.reserve() {
			defer actionClient.waitable.release()
		}
	}
	for _, actionServer := range self.ActionServers {
		if actionServer.waitable.reserve() {
			defer actionServer.waitable.release()
		}
	}
	for _, guardCondition := range self.guardConditions {
		if guardCondition.waitable.reserve() {
			defer guardCondition.waitable.release()
		}
	}
	if ctx == nil {
		return errors.New("context must not be nil")
	}
	errs := make(chan error, 1)
	defer func() {
		err = multierror.Append(err, <-errs).ErrorOrNil()
	}()
	errctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer close(errs)
		<-errctx.Done()
		errs <- self.cancelWait.Trigger()
	}()
	for {
		if err := self.initEntities(); err != nil {
			return err
		}
		if rc := C.rcl_wait(&self.rcl_wait_set_t, -1); rc != C.RCL_RET_OK {
			return errorsCast(rc)
		}
		guardConditions := unsafe.Slice(self.rcl_wait_set_t.guard_conditions, len(self.guardConditions))
		for i := range self.guardConditions {
			if guardConditions[i] == self.cancelWait.rclGuardCondition {
				return ctx.Err()
			}
		}
		timers := unsafe.Slice(self.rcl_wait_set_t.timers, len(self.Timers))
		for i, t := range self.Timers {
			if timers[i] != nil {
				t.Reset() //nolint:errcheck
				t.Callback(t)
			}
		}
		subs := unsafe.Slice(self.rcl_wait_set_t.subscriptions, len(self.Subscriptions))
		for i, s := range self.Subscriptions {
			if subs[i] != nil {
				s.Callback(s)
			}
		}
		svcs := unsafe.Slice(self.rcl_wait_set_t.services, len(self.Services))
		for i, s := range self.Services {
			if svcs[i] != nil {
				s.handleRequest()
			}
		}
		clients := unsafe.Slice(self.rcl_wait_set_t.clients, len(self.Clients))
		for i, c := range self.Clients {
			if clients[i] != nil {
				c.sender.HandleResponse()
			}
		}
		for _, s := range self.ActionServers {
			s.handleReadyEntities(ctx, self)
		}
		for _, c := range self.ActionClients {
			c.handleReadyEntities(self)
		}
	}
}

func (self *WaitSet) initEntities() error {
	if !C.rcl_wait_set_is_valid(&self.rcl_wait_set_t) {
		return errorsCastC(C.RCL_RET_WAIT_SET_INVALID, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", self))
	}
	var rc C.rcl_ret_t = C.rcl_wait_set_clear(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", self))
	}
	rc = C.rcl_wait_set_resize(
		&self.rcl_wait_set_t,
		C.size_t(len(self.Subscriptions)+2*len(self.ActionClients)),
		C.size_t(len(self.guardConditions)),
		C.size_t(len(self.Timers)+len(self.ActionServers)),
		C.size_t(len(self.Clients)+3*len(self.ActionClients)),
		C.size_t(len(self.Services)+3*len(self.ActionServers)),
		self.rcl_wait_set_t.size_of_events,
	)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_resize() failed for wait_set='%v'", self))
	}
	for _, sub := range self.Subscriptions {
		rc = C.rcl_wait_set_add_subscription(&self.rcl_wait_set_t, sub.rcl_subscription_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", self))
		}
	}
	for _, timer := range self.Timers {
		rc = C.rcl_wait_set_add_timer(&self.rcl_wait_set_t, timer.rcl_timer_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_timer() failed for wait_set='%v'", self))
		}
	}
	for _, service := range self.Services {
		rc = C.rcl_wait_set_add_service(&self.rcl_wait_set_t, service.rclService, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_service() failed for wait_set='%v'", self))
		}
	}
	for _, client := range self.Clients {
		rc = C.rcl_wait_set_add_client(&self.rcl_wait_set_t, client.rclClient, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_client() failed for wait_set='%v'", self))
		}
	}
	for _, guardCondition := range self.guardConditions {
		rc = C.rcl_wait_set_add_guard_condition(&self.rcl_wait_set_t, guardCondition.rclGuardCondition, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_guard_condition() failed for wait_set='%v'", self))
		}
	}
	for _, server := range self.ActionServers {
		rc = C.rcl_action_wait_set_add_action_server(&self.rcl_wait_set_t, &server.rclServer, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_action_server() failed for wait_set='%v'", self))
		}
	}
	for _, client := range self.ActionClients {
		rc = C.rcl_action_wait_set_add_action_client(&self.rcl_wait_set_t, &client.rclClient, nil, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_action_client() failed for wait_set='%v'", self))
		}
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (self *WaitSet) Close() error {
	if self.context == nil {
		return closeErr("wait set")
	}
	var errs *multierror.Error
	self.context.removeResource(self)
	self.context = nil
	rc := C.rcl_wait_set_fini(&self.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		errs = multierror.Append(errs, errorsCast(rc))
	}
	var closeError closeError
	err := self.cancelWait.Close()
	if err != nil && !errors.As(err, &closeError) {
		errs = multierror.Append(errs, err)
	}
	return errs.ErrorOrNil()
}
