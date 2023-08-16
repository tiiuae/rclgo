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
func (w *WaitSet) Context() *Context {
	return w.context
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
func (w *WaitSet) Run(ctx context.Context) (err error) {
	for _, subscription := range w.Subscriptions {
		if subscription.waitable.reserve() {
			defer subscription.waitable.release()
		}
	}
	for _, timer := range w.Timers {
		if timer.waitable.reserve() {
			defer timer.waitable.release()
		}
	}
	for _, service := range w.Services {
		if service.waitable.reserve() {
			defer service.waitable.release()
		}
	}
	for _, client := range w.Clients {
		if client.waitable.reserve() {
			defer client.waitable.release()
		}
	}
	for _, actionClient := range w.ActionClients {
		if actionClient.waitable.reserve() {
			defer actionClient.waitable.release()
		}
	}
	for _, actionServer := range w.ActionServers {
		if actionServer.waitable.reserve() {
			defer actionServer.waitable.release()
		}
	}
	for _, guardCondition := range w.guardConditions {
		if guardCondition.waitable.reserve() {
			defer guardCondition.waitable.release()
		}
	}
	if ctx == nil {
		return errors.New("context must not be nil")
	}
	errs := make(chan error, 1)
	defer func() {
		err = errors.Join(err, <-errs)
	}()
	errctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer close(errs)
		<-errctx.Done()
		errs <- w.cancelWait.Trigger()
	}()
	for {
		if err := w.initEntities(); err != nil {
			return err
		}
		if rc := C.rcl_wait(&w.rcl_wait_set_t, -1); rc != C.RCL_RET_OK {
			return errorsCast(rc)
		}
		guardConditions := unsafe.Slice(w.rcl_wait_set_t.guard_conditions, len(w.guardConditions))
		for i := range w.guardConditions {
			if guardConditions[i] == w.cancelWait.rclGuardCondition {
				return ctx.Err()
			}
		}
		timers := unsafe.Slice(w.rcl_wait_set_t.timers, len(w.Timers))
		for i, t := range w.Timers {
			if timers[i] != nil {
				t.Reset() //nolint:errcheck
				t.Callback(t)
			}
		}
		subs := unsafe.Slice(w.rcl_wait_set_t.subscriptions, len(w.Subscriptions))
		for i, s := range w.Subscriptions {
			if subs[i] != nil {
				s.Callback(s)
			}
		}
		svcs := unsafe.Slice(w.rcl_wait_set_t.services, len(w.Services))
		for i, s := range w.Services {
			if svcs[i] != nil {
				s.handleRequest()
			}
		}
		clients := unsafe.Slice(w.rcl_wait_set_t.clients, len(w.Clients))
		for i, c := range w.Clients {
			if clients[i] != nil {
				c.sender.HandleResponse()
			}
		}
		for _, s := range w.ActionServers {
			s.handleReadyEntities(ctx, w)
		}
		for _, c := range w.ActionClients {
			c.handleReadyEntities(w)
		}
	}
}

func (w *WaitSet) initEntities() error {
	if !C.rcl_wait_set_is_valid(&w.rcl_wait_set_t) {
		return errorsCastC(C.RCL_RET_WAIT_SET_INVALID, fmt.Sprintf("rcl_wait_set_is_valid() failed for wait_set='%v'", w))
	}
	var rc C.rcl_ret_t = C.rcl_wait_set_clear(&w.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_clear() failed for wait_set='%v'", w))
	}
	rc = C.rcl_wait_set_resize(
		&w.rcl_wait_set_t,
		C.size_t(len(w.Subscriptions)+2*len(w.ActionClients)),
		C.size_t(len(w.guardConditions)),
		C.size_t(len(w.Timers)+len(w.ActionServers)),
		C.size_t(len(w.Clients)+3*len(w.ActionClients)),
		C.size_t(len(w.Services)+3*len(w.ActionServers)),
		w.rcl_wait_set_t.size_of_events,
	)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_resize() failed for wait_set='%v'", w))
	}
	for _, sub := range w.Subscriptions {
		rc = C.rcl_wait_set_add_subscription(&w.rcl_wait_set_t, sub.rcl_subscription_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_subscription() failed for wait_set='%v'", w))
		}
	}
	for _, timer := range w.Timers {
		rc = C.rcl_wait_set_add_timer(&w.rcl_wait_set_t, timer.rcl_timer_t, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_timer() failed for wait_set='%v'", w))
		}
	}
	for _, service := range w.Services {
		rc = C.rcl_wait_set_add_service(&w.rcl_wait_set_t, service.rclService, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_service() failed for wait_set='%v'", w))
		}
	}
	for _, client := range w.Clients {
		rc = C.rcl_wait_set_add_client(&w.rcl_wait_set_t, client.rclClient, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_client() failed for wait_set='%v'", w))
		}
	}
	for _, guardCondition := range w.guardConditions {
		rc = C.rcl_wait_set_add_guard_condition(&w.rcl_wait_set_t, guardCondition.rclGuardCondition, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_guard_condition() failed for wait_set='%v'", w))
		}
	}
	for _, server := range w.ActionServers {
		rc = C.rcl_action_wait_set_add_action_server(&w.rcl_wait_set_t, &server.rclServer, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_action_server() failed for wait_set='%v'", w))
		}
	}
	for _, client := range w.ActionClients {
		rc = C.rcl_action_wait_set_add_action_client(&w.rcl_wait_set_t, &client.rclClient, nil, nil)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, fmt.Sprintf("rcl_wait_set_add_action_client() failed for wait_set='%v'", w))
		}
	}
	return nil
}

/*
Close frees the allocated memory
*/
func (w *WaitSet) Close() (err error) {
	if w.context == nil {
		return closeErr("wait set")
	}
	w.context.removeResource(w)
	w.context = nil
	rc := C.rcl_wait_set_fini(&w.rcl_wait_set_t)
	if rc != C.RCL_RET_OK {
		err = errors.Join(err, errorsCast(rc))
	}
	var closeError closeError
	cancelWaitErr := w.cancelWait.Close()
	if cancelWaitErr != nil && !errors.As(cancelWaitErr, &closeError) {
		err = errors.Join(err, cancelWaitErr)
	}
	return err
}
