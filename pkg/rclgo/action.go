package rclgo

/*
#include <stdlib.h>

#include <rcl_action/rcl_action.h>
*/
import "C"
import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/hashicorp/go-multierror"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

func NewRmwQosProfileStatusDefault() RmwQosProfile {
	return RmwQosProfile{
		History:                 RmwQosHistoryPolicyKeepLast,
		Depth:                   1,
		Reliability:             RmwQosReliabilityPolicyReliable,
		Durability:              RmwQosDurabilityPolicyTransientLocal,
		Deadline:                RmwQosDeadlineDefault,
		Lifespan:                RmwQosLifespanDefault,
		Liveliness:              RmwQosLivelinessPolicySystemDefault,
		LivelinessLeaseDuration: RmwQosLivelinessLeaseDurationDefault,
	}
}

type goalIDMessage interface {
	types.Message
	GetGoalID() *types.GoalID
	SetGoalID(*types.GoalID)
}

type goalRequestMessage interface {
	goalIDMessage
	GetGoalDescription() types.Message
	SetGoalDescription(types.Message)
}

type goalResponseMessage interface {
	GetGoalAccepted() bool
}

type forEach interface {
	CallForEach(func(interface{}))
}

type GoalStatus int8

//go:generate go run golang.org/x/tools/cmd/stringer -type GoalStatus -linecomment

const (
	GoalUnknown GoalStatus = iota // Unknown

	// Active states
	GoalAccepted  // Accepted
	GoalExecuting // Executing
	GoalCanceling // Canceling

	// Terminal states
	GoalSucceeded // Succeeded
	GoalCanceled  // Canceled
	GoalAborted   // Aborted
)

// FeedbackSender is used to send feedback about a goal.
type FeedbackSender struct {
	goal *GoalHandle
}

// Send sends msg to clients listening for feedback messages.
//
// The type support of msg must be types.ActionTypeSupport.Feedback().
func (s *FeedbackSender) Send(msg types.Message) error {
	return s.goal.server.sendFeedback(s.goal, msg)
}

// GoalHandle is used to keep track of the status of a goal sent to an
// ActionServer.
type GoalHandle struct {
	// The ID of the goal. Modifying this is undefined behavior.
	ID types.GoalID
	// Description is a message whose type support is
	// types.ActionTypeSupport.Goal().
	Description types.Message

	server       *ActionServer
	cancelContex context.CancelFunc
	requestID    C.rmw_request_id_t
	handle       *C.rcl_action_goal_handle_t
	info         C.rcl_action_goal_info_t

	result     types.Message
	resultCond *sync.Cond
}

func newEmptyGoal(s *ActionServer, cancel context.CancelFunc) *GoalHandle {
	goal := &GoalHandle{
		server:       s,
		cancelContex: cancel,
		info:         C.rcl_action_get_zero_initialized_goal_info(),
		resultCond:   sync.NewCond(&sync.Mutex{}),
	}
	runtime.SetFinalizer(goal, func(g *GoalHandle) {
		if g.handle != nil {
			C.rcl_action_goal_handle_fini(g.handle)
			g.handle = nil
		}
	})
	return goal
}

// Server returns the ActionServer that is handling g.
func (g *GoalHandle) Server() *ActionServer {
	return g.server
}

// Logger is a shorthand for g.Server().Node().Logger().
func (g *GoalHandle) Logger() *Logger {
	return g.server.node.logger
}

func (g *GoalHandle) status() GoalStatus {
	if g.handle == nil {
		return GoalUnknown
	}
	var status C.schar
	rc := C.rcl_action_goal_handle_get_status(g.handle, &status)
	if rc != C.RCL_RET_OK {
		panic(errorsCastC(rc, "failed to get goal status"))
	}
	return GoalStatus(status)
}

func (g *GoalHandle) setState(event C.rcl_action_goal_event_t) {
	rc := C.rcl_action_update_goal_state(g.handle, event)
	if rc == C.RCL_RET_OK {
		g.server.sendStatuses(true)
	} else {
		g.server.node.Logger().Debug(errorsCastC(rc, "failed to update goal state"))
	}
}

// Accept accepts g and returns a FeedbackSender that can be used to send
// feedback about the goal to action clients. Calls after the first successful
// (returning a nil error) do not change the state of the goal and only return
// valid feedback senders. Accept should be called as soon as the goal is
// decided to be accepted. If Accept returns a non-nil error the returned
// FeedbackSender is nil and g is left in an unspecified but valid state. In
// that case it is usually appropriate to stop executing the goal and return the
// error returned by Accept.
func (g *GoalHandle) Accept() (s *FeedbackSender, err error) {
	defer wrapErr("failed to accept goal: %w", &err)
	g.resultCond.L.Lock()
	defer g.resultCond.L.Unlock()
	if g.status() == GoalUnknown {
		g.info.goal_id.uuid = *(*[types.GoalIDLen]C.uchar)(unsafe.Pointer(&g.ID))
		if err := g.server.acceptGoal(g); err != nil {
			return nil, err
		}
	}
	if g.status() == GoalAccepted {
		stamp := time.Duration(g.info.stamp.sec) * time.Second
		stamp += time.Duration(g.info.stamp.nanosec)
		if err := g.server.sendGoalResponseWithStamp(g, stamp); err != nil {
			return nil, err
		}
		g.setState(C.GOAL_EVENT_EXECUTE)
	}
	return &FeedbackSender{goal: g}, nil
}

func (g *GoalHandle) abort() {
	g.resultCond.L.Lock()
	defer g.resultCond.L.Unlock()
	g.setState(C.GOAL_EVENT_ABORT)
	g.resultCond.Broadcast()
}

func (g *GoalHandle) setResult(msg types.Message) {
	g.resultCond.L.Lock()
	defer g.resultCond.L.Unlock()
	g.setState(C.GOAL_EVENT_SUCCEED)
	g.result = msg
	g.server.rclServerMu.Lock()
	defer g.server.rclServerMu.Unlock()
	C.rcl_action_notify_goal_done(&g.server.rclServer)
	g.resultCond.Broadcast()
}

func (g *GoalHandle) startCancel() {
	g.resultCond.L.Lock()
	defer g.resultCond.L.Unlock()
	g.cancelContex()
	g.setState(C.GOAL_EVENT_CANCEL_GOAL)
}

func (g *GoalHandle) finishCancel() {
	g.resultCond.L.Lock()
	defer g.resultCond.L.Unlock()
	g.setState(C.GOAL_EVENT_CANCELED)
	g.resultCond.Broadcast()
}

func (g *GoalHandle) waitResult() (GoalStatus, types.Message) {
	g.resultCond.L.Lock()
	for {
		status := g.status()
		switch status {
		case GoalSucceeded, GoalAborted, GoalCanceled:
			defer g.resultCond.L.Unlock()
			return status, g.result
		default:
			g.resultCond.Wait()
		}
	}
}

type ExecuteGoalFunc = func(context.Context, *GoalHandle) (types.Message, error)

// Action can execute goals.
type Action interface {
	// ExecuteGoal executes a goal.
	//
	// The description of the goal is passed in the GoalHandle.
	//
	// First ExecuteGoal must decide whether to accept the goal or not. The goal
	// can be accepted by calling GoalHandle.Accept. GoalHandle.Accept should be
	// called as soon as the decision to accept the goal is made, before
	// starting to execute the goal.
	//
	// ExecuteGoal returns a pair of (result, error). If ExecuteGoal returns a
	// nil error, the goal is assumed to be executed successfully to completion.
	// In this case the result must be non-nil, and its type support must be
	// TypeSupport().Result(). If ExecuteGoal returns a non-nil error, the
	// result is ignored. If the error is returned before accepting the goal,
	// the goal is considered to have been rejected. If the error is returned
	// after accepting the goal, the goal is considered to have been aborted.
	//
	// The context is used to notify cancellation of the goal. If the context is
	// canceled, ExecuteGoal should stop all processing as soon as possible. In
	// this case the return values of ExecuteGoal are ignored.
	//
	// ExecuteGoal may be called multiple times in parallel by the ActionServer.
	// Each call will receive a different GoalHandle.
	ExecuteGoal(context.Context, *GoalHandle) (types.Message, error)

	// TypeSupport returns the type support for the action. The same value
	// must be returned on every invocation.
	TypeSupport() types.ActionTypeSupport
}

// NewAction returns an Action implementation that uses typeSupport and executes
// goals using executeGoal.
func NewAction(
	typeSupport types.ActionTypeSupport,
	executeGoal ExecuteGoalFunc,
) Action {
	return &action{
		typeSupport: typeSupport,
		executeGoal: executeGoal,
	}
}

type action struct {
	typeSupport types.ActionTypeSupport
	executeGoal ExecuteGoalFunc
}

func (a *action) ExecuteGoal(
	ctx context.Context, goal *GoalHandle,
) (types.Message, error) {
	return a.executeGoal(ctx, goal)
}

func (a *action) TypeSupport() types.ActionTypeSupport {
	return a.typeSupport
}

type ActionServerOptions struct {
	GoalServiceQos   RmwQosProfile
	CancelServiceQos RmwQosProfile
	ResultServiceQos RmwQosProfile
	FeedbackTopicQos RmwQosProfile
	StatusTopicQos   RmwQosProfile
	ResultTimeout    time.Duration
	Clock            *Clock
}

func NewDefaultActionServerOptions() *ActionServerOptions {
	return &ActionServerOptions{
		GoalServiceQos:   NewRmwQosProfileServicesDefault(),
		CancelServiceQos: NewRmwQosProfileServicesDefault(),
		ResultServiceQos: NewRmwQosProfileServicesDefault(),
		FeedbackTopicQos: NewRmwQosProfileDefault(),
		StatusTopicQos:   NewRmwQosProfileStatusDefault(),
		ResultTimeout:    15 * time.Minute,
	}
}

// ActionServer listens for and executes goals sent by action clients.
type ActionServer struct {
	rosID
	node          *Node
	action        Action
	typeSupport   types.ActionTypeSupport
	resultTimeout time.Duration
	clock         *Clock

	rclServer    C.rcl_action_server_t
	rclServerMu  sync.Mutex
	timerCancels []func()

	goals   map[types.GoalID]*GoalHandle
	goalsMu sync.RWMutex
}

// NewActionServer creates a new action server.
//
// opts must not be modified after passing it to this function. If opts is nil,
// default options are used.
func (n *Node) NewActionServer(
	name string,
	action Action,
	opts *ActionServerOptions,
) (*ActionServer, error) {
	if opts == nil {
		opts = NewDefaultActionServerOptions()
	}
	s := &ActionServer{
		node:          n,
		action:        action,
		typeSupport:   action.TypeSupport(),
		resultTimeout: opts.ResultTimeout,
		clock:         opts.Clock,

		rclServer: C.rcl_action_get_zero_initialized_server(),

		goals: make(map[types.GoalID]*GoalHandle),
	}
	if s.clock == nil {
		s.clock = n.context.Clock()
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rclOpts := C.rcl_action_server_options_t{
		allocator: *n.context.rcl_allocator_t,
		result_timeout: C.rcl_duration_t{
			nanoseconds: C.long(opts.ResultTimeout),
		},
	}
	opts.GoalServiceQos.asCStruct(&rclOpts.goal_service_qos)
	opts.CancelServiceQos.asCStruct(&rclOpts.cancel_service_qos)
	opts.ResultServiceQos.asCStruct(&rclOpts.result_service_qos)
	opts.FeedbackTopicQos.asCStruct(&rclOpts.feedback_topic_qos)
	opts.StatusTopicQos.asCStruct(&rclOpts.status_topic_qos)
	rc := C.rcl_action_server_init(
		&s.rclServer,
		n.rcl_node_t,
		s.clock.rcl_clock_t,
		(*C.rosidl_action_type_support_t)(s.typeSupport.TypeSupport()),
		cname,
		&rclOpts,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to create action server")
	}
	n.addResource(s)
	return s, nil
}

func (s *ActionServer) Close() (err error) {
	if s.rclServer == C.rcl_action_get_zero_initialized_server() {
		return closeErr("action")
	}
	s.node.removeResource(s)
	s.rclServerMu.Lock()
	for _, f := range s.timerCancels {
		f()
	}
	s.rclServerMu.Unlock()
	rc := C.rcl_action_server_fini(&s.rclServer, s.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = errorsCastC(rc, "failed to finalize action server")
	}
	s.rclServer = C.rcl_action_get_zero_initialized_server()
	return err
}

// Node returns the node s was created with.
func (s *ActionServer) Node() *Node {
	return s.node
}

func (s *ActionServer) takeGoalRequest(goal *GoalHandle) error {
	ts := s.typeSupport.SendGoal().Request()
	reqBuf := ts.PrepareMemory()
	defer ts.ReleaseMemory(reqBuf)
	rc := C.rcl_action_take_goal_request(&s.rclServer, &goal.requestID, reqBuf)
	if rc != C.RCL_RET_OK {
		return fmt.Errorf("%v", errorsCastC(rc, "failed to take goal request"))
	}
	req := ts.New().(goalRequestMessage)
	ts.AsGoStruct(req, reqBuf)
	goal.ID = *req.GetGoalID()
	goal.Description = req.GetGoalDescription()
	return nil
}

func (s *ActionServer) acceptGoal(goal *GoalHandle) error {
	s.rclServerMu.Lock()
	defer s.rclServerMu.Unlock()
	goal.handle = C.rcl_action_accept_new_goal(&s.rclServer, &goal.info)
	if goal.handle == nil {
		return errors.New("accepting a goal handle failed: " + errorString())
	}
	s.sendStatuses(false)
	s.goalsMu.Lock()
	defer s.goalsMu.Unlock()
	if s.goals[goal.ID] == nil {
		s.goals[goal.ID] = goal
		return nil
	}
	return fmt.Errorf("goal with ID %v has already been accepted", goal.ID)
}

func (s *ActionServer) sendGoalResponseWithStamp(goal *GoalHandle, stamp time.Duration) error {
	resp := s.typeSupport.NewSendGoalResponse(
		goal.status() != GoalUnknown,
		stamp,
	)
	ts := s.typeSupport.SendGoal().Response()
	buf := ts.PrepareMemory()
	defer ts.ReleaseMemory(buf)
	ts.AsCStruct(buf, resp)
	rc := C.rcl_action_send_goal_response(&s.rclServer, &goal.requestID, buf)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "failed to send goal response")
	}
	return nil
}

func (s *ActionServer) sendGoalResponse(goal *GoalHandle) {
	now, err := s.node.context.Clock().now()
	if err != nil {
		s.logGoalError(goal, err)
		return
	}
	if err = s.sendGoalResponseWithStamp(goal, now); err != nil {
		s.logGoalError(goal, err)
	}
}

func (s *ActionServer) handleGoalRequest(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	goal := newEmptyGoal(s, cancel)
	err := s.takeGoalRequest(goal)
	go func() {
		defer cancel()
		if err != nil {
			s.node.Logger().Error(err)
			return
		}
		var result types.Message
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("rclgo: panic when executing goal: %v", r)
				}
			}()
			result, err = s.action.ExecuteGoal(ctx, goal)
		}()
		if ctx.Err() == nil {
			if err == nil {
				if result == nil {
					s.logGoalError(goal, "a nil result was returned even though the goal was executed successfully")
					goal.abort()
				} else {
					goal.setResult(result)
				}
			} else if goal.status() == GoalUnknown {
				s.sendGoalResponse(goal)
			} else {
				goal.abort()
			}
		} else {
			goal.finishCancel()
		}
		if s.resultTimeout >= 0 {
			s.scheduleRemoval()
		}
	}()
}

func (s *ActionServer) scheduleRemoval() {
	s.rclServerMu.Lock()
	defer s.rclServerMu.Unlock()
	timer := time.AfterFunc(s.resultTimeout, func() {
		if err := s.expireGoals(); err != nil {
			s.node.Logger().Error(err)
		}
	})
	s.timerCancels = append(s.timerCancels, func() {
		if !timer.Stop() {
			<-timer.C
		}
	})
}

func (s *ActionServer) expireGoals() error {
	var numExpired C.size_t
	expiredGoals := make([]C.rcl_action_goal_info_t, 10)
	expiredGoalsHeader := (*reflect.SliceHeader)(unsafe.Pointer(&expiredGoals))
	s.rclServerMu.Lock()
	defer s.rclServerMu.Unlock()
	s.goalsMu.Lock()
	defer s.goalsMu.Unlock()
	for {
		rc := C.rcl_action_expire_goals(
			&s.rclServer,
			(*C.rcl_action_goal_info_t)(unsafe.Pointer(expiredGoalsHeader.Data)),
			C.ulong(expiredGoalsHeader.Len),
			&numExpired,
		)
		if rc != C.RCL_RET_OK {
			return errorsCastC(rc, "failed to expire goals")
		}
		if numExpired == 0 {
			return nil
		}
		for _, goal := range expiredGoals[:int(numExpired)] {
			delete(s.goals, *(*types.GoalID)(unsafe.Pointer(&goal.goal_id.uuid)))
		}
	}
}

func (s *ActionServer) takeResultRequest(header *C.rmw_request_id_t) (goalIDMessage, error) {
	ts := s.typeSupport.GetResult().Request()
	reqBuf := ts.PrepareMemory()
	defer ts.ReleaseMemory(reqBuf)
	rc := C.rcl_action_take_result_request(&s.rclServer, header, reqBuf)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to take result request")
	}
	req := ts.New().(goalIDMessage)
	ts.AsGoStruct(req, reqBuf)
	return req, nil
}

func (s *ActionServer) getGoal(id *types.GoalID, takeLock bool) *GoalHandle {
	if takeLock {
		s.goalsMu.RLock()
		defer s.goalsMu.RUnlock()
	}
	return s.goals[*id]
}

func (s *ActionServer) sendResultResponse(status GoalStatus, result types.Message, header *C.rmw_request_id_t) error {
	ts := s.typeSupport.GetResult().Response()
	resp := s.typeSupport.NewGetResultResponse(int8(status), result)
	buf := ts.PrepareMemory()
	defer ts.ReleaseMemory(buf)
	ts.AsCStruct(buf, resp)
	rc := C.rcl_action_send_result_response(&s.rclServer, header, buf)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "failed to send result response")
	}
	return nil
}

func (s *ActionServer) handleResultRequest() {
	var header C.rmw_request_id_t
	req, err := s.takeResultRequest(&header)
	go func() {
		if err != nil {
			s.node.Logger().Error(err)
			return
		}
		goal := s.getGoal(req.GetGoalID(), true)
		if goal == nil {
			err = s.sendResultResponse(GoalUnknown, nil, &header)
		} else {
			status, result := goal.waitResult()
			err = s.sendResultResponse(status, result, &header)
		}
		if err != nil {
			s.logGoalError(goal, err)
		}
		if err = s.expireGoals(); err != nil {
			s.logGoalError(goal, err)
		}
	}()
}

func (s *ActionServer) processCancelRequest(
	header *C.rmw_request_id_t,
	req unsafe.Pointer,
	resp *C.rcl_action_cancel_response_t,
) error {
	s.rclServerMu.Lock()
	defer s.rclServerMu.Unlock()
	rc := C.rcl_action_process_cancel_request(
		&s.rclServer,
		(*C.rcl_action_cancel_request_t)(req),
		resp,
	)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "failed to process cancel request")
	}
	return nil
}

func (s *ActionServer) handleCancelRequest() {
	ts := s.typeSupport.CancelGoal()
	reqBuf := ts.Request().PrepareMemory()
	respBuf := C.rcl_action_get_zero_initialized_cancel_response()
	var header C.rmw_request_id_t
	rc := C.rcl_action_take_cancel_request(&s.rclServer, &header, reqBuf)
	go func() {
		defer ts.Request().ReleaseMemory(reqBuf)
		if rc != C.RCL_RET_OK {
			s.node.Logger().Error(errorsCastC(rc, "failed to take cancel request"))
			return
		}
		err := s.processCancelRequest(&header, reqBuf, &respBuf)
		if err != nil {
			s.node.Logger().Error(err)
			return
		}
		defer C.rcl_action_cancel_response_fini(&respBuf)
		resp := ts.Response().New()
		ts.Response().AsGoStruct(resp, unsafe.Pointer(&respBuf.msg))
		func() {
			s.goalsMu.RLock()
			defer s.goalsMu.RUnlock()
			resp.(forEach).CallForEach(func(toCancel interface{}) {
				goal := s.getGoal(toCancel.(*types.GoalID), false)
				if goal != nil {
					goal.startCancel()
				}
			})
		}()
		rc = C.rcl_action_send_cancel_response(
			&s.rclServer,
			&header,
			unsafe.Pointer(&respBuf.msg),
		)
		if rc != C.RCL_RET_OK {
			s.node.Logger().Error(errorsCastC(rc, "failed to send cancel response"))
		}
	}()
}

func (s *ActionServer) sendFeedback(goal *GoalHandle, fb types.Message) error {
	msg := s.typeSupport.NewFeedbackMessage(&goal.ID, fb)
	ts := s.typeSupport.FeedbackMessage()
	buf := ts.PrepareMemory()
	defer ts.ReleaseMemory(buf)
	ts.AsCStruct(buf, msg)
	rc := C.rcl_action_publish_feedback(&s.rclServer, buf)
	if rc != C.RCL_RET_OK {
		return errorsCastC(rc, "failed to send feedback")
	}
	return nil
}

func (s *ActionServer) sendStatuses(lock bool) {
	if lock {
		s.rclServerMu.Lock()
		defer s.rclServerMu.Unlock()
	}
	statuses := C.rcl_action_get_zero_initialized_goal_status_array()
	defer C.rcl_action_goal_status_array_fini(&statuses)
	rc := C.rcl_action_get_goal_status_array(&s.rclServer, &statuses)
	if rc != C.RCL_RET_OK {
		s.node.Logger().Error("failed to get status array: ", errorsCast(rc))
		return
	}
	rc = C.rcl_action_publish_status(&s.rclServer, unsafe.Pointer(&statuses.msg))
	if rc != C.RCL_RET_OK {
		s.node.Logger().Error("failed to publish status array: ", errorsCast(rc))
	}
}

func (s *ActionServer) handleReadyEntities(ctx context.Context, ws *WaitSet) {
	s.rclServerMu.Lock()
	var goalReq, cancelReq, resultReq, expired C.bool
	rc := C.rcl_action_server_wait_set_get_entities_ready(
		&ws.rcl_wait_set_t,
		&s.rclServer,
		&goalReq,
		&cancelReq,
		&resultReq,
		&expired,
	)
	s.rclServerMu.Unlock()
	if rc != C.RCL_RET_OK {
		s.node.Logger().Error(errorsCastC(rc, "failed to get ready entities"))
		return
	}
	if goalReq {
		s.handleGoalRequest(ctx)
	}
	if cancelReq {
		s.handleCancelRequest()
	}
	if resultReq {
		s.handleResultRequest()
	}
	if expired {
		s.expireGoals()
	}
}

func (s *ActionServer) logGoalError(goal *GoalHandle, a ...interface{}) {
	var b strings.Builder
	fmt.Fprint(&b, "goal: ", goal.ID.String(), ": ")
	fmt.Fprint(&b, a...)
	s.node.Logger().Error(b.String())
}

type FeedbackHandler func(context.Context, types.Message)

type StatusHandler func(context.Context, types.Message)

type ActionClientOptions struct {
	GoalServiceQos   RmwQosProfile
	CancelServiceQos RmwQosProfile
	ResultServiceQos RmwQosProfile
	FeedbackTopicQos RmwQosProfile
	StatusTopicQos   RmwQosProfile
}

func NewDefaultActionClientOptions() *ActionClientOptions {
	return &ActionClientOptions{
		GoalServiceQos:   NewRmwQosProfileServicesDefault(),
		CancelServiceQos: NewRmwQosProfileServicesDefault(),
		ResultServiceQos: NewRmwQosProfileServicesDefault(),
		FeedbackTopicQos: NewRmwQosProfileDefault(),
		StatusTopicQos:   NewRmwQosProfileDefault(),
	}
}

type actionClientHandler = func(context.Context, types.Message)

type actionClientHandlerMapEntry struct {
	ctx     context.Context
	handler actionClientHandler
}

type actionClientHandlerMap map[uint64]actionClientHandlerMapEntry

func (m actionClientHandlerMap) call(msg types.Message) {
	for _, entry := range m {
		entry := entry
		go func() { entry.handler(entry.ctx, msg.CloneMsg()) }()
	}
}

type actionClientSubs struct {
	perGoal  map[types.GoalID]actionClientHandlerMap
	allGoals actionClientHandlerMap
}

func newActionClientHandlers() actionClientSubs {
	return actionClientSubs{
		perGoal:  make(map[types.GoalID]actionClientHandlerMap),
		allGoals: make(actionClientHandlerMap),
	}
}

// ActionClient communicates with an ActionServer to initiate and monitor the
// progress of goals.
//
// All methods except Close are safe for concurrent use.
type ActionClient struct {
	rosID
	node *Node

	typeSupport types.ActionTypeSupport
	rclClient   C.rcl_action_client_t
	rclClientMu sync.Mutex

	goalSender   requestSender
	resultSender requestSender
	cancelSender requestSender

	nextSubscriberID uint64
	feedbackSubs     actionClientSubs
	statusSubs       actionClientSubs
}

// NewActionClient creates an action client that communicates with an action
// server.
func (n *Node) NewActionClient(
	name string,
	ts types.ActionTypeSupport,
	opts *ActionClientOptions,
) (*ActionClient, error) {
	if opts == nil {
		opts = NewDefaultActionClientOptions()
	}
	c := &ActionClient{
		node: n,

		typeSupport: ts,
		rclClient:   C.rcl_action_get_zero_initialized_client(),

		feedbackSubs: newActionClientHandlers(),
		statusSubs:   newActionClientHandlers(),
	}
	c.goalSender = newRequestSender(requestSenderTransport{
		SendRequest:  c.sendGoalRequest,
		TakeResponse: c.takeGoalResponse,
		TypeSupport:  ts.SendGoal(),
		Logger:       n.Logger(),
	})
	c.resultSender = newRequestSender(requestSenderTransport{
		SendRequest:  c.sendResultRequest,
		TakeResponse: c.takeResultResponse,
		TypeSupport:  ts.GetResult(),
		Logger:       n.Logger(),
	})
	c.cancelSender = newRequestSender(requestSenderTransport{
		SendRequest:  c.sendCancelRequest,
		TakeResponse: c.takeCancelResponse,
		TypeSupport:  ts.CancelGoal(),
		Logger:       n.Logger(),
	})
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rclOpts := C.rcl_action_client_options_t{
		allocator: *n.context.rcl_allocator_t,
	}
	opts.GoalServiceQos.asCStruct(&rclOpts.goal_service_qos)
	opts.CancelServiceQos.asCStruct(&rclOpts.cancel_service_qos)
	opts.ResultServiceQos.asCStruct(&rclOpts.result_service_qos)
	opts.FeedbackTopicQos.asCStruct(&rclOpts.feedback_topic_qos)
	opts.StatusTopicQos.asCStruct(&rclOpts.status_topic_qos)
	rc := C.rcl_action_client_init(
		&c.rclClient,
		n.rcl_node_t,
		(*C.rosidl_action_type_support_t)(ts.TypeSupport()),
		cname,
		&rclOpts,
	)
	if rc != C.RCL_RET_OK {
		return nil, errorsCastC(rc, "failed to create action client")
	}
	n.addResource(c)
	return c, nil
}

// Close frees resources used by the ActionClient. A closed ActionClient must
// not be used.
func (c *ActionClient) Close() error {
	if c.typeSupport == nil {
		return closeErr("action client")
	}
	c.node.removeResource(c)
	err := multierror.Append(
		c.cancelSender.Close(),
		c.goalSender.Close(),
		c.resultSender.Close(),
	)
	rc := C.rcl_action_client_fini(&c.rclClient, c.node.rcl_node_t)
	if rc != C.RCL_RET_OK {
		err = multierror.Append(
			err,
			errorsCastC(rc, "failed to finalize action client"),
		)
	}
	c.typeSupport = nil
	return err.ErrorOrNil()
}

// Node returns the node c was created with.
func (c *ActionClient) Node() *Node {
	return c.node
}

// WatchGoal combines functionality of SendGoal and WatchFeedback. It sends a
// goal to the server. If the goal is accepted, feedback for the goal is watched
// until the goal reaches a terminal state or ctx is canceled. If the goal is
// accepted and completes succesfully, its result is returned. Otherwise a
// non-nil error is returned.
//
// onFeedback may be nil, in which case feedback for the goal is not watched.
//
// The type support of goal must be types.ActionTypeSupport.Goal().
//
// The type support of the returned message is types.ActionTypeSupport.Result().
//
// The type support of the message passed to onFeedback is
// types.ActionTypeSupport.FeedbackMessage().
func (c *ActionClient) WatchGoal(ctx context.Context, goal types.Message, onFeedback FeedbackHandler) (types.Message, error) {
	req, err := c.newSendGoalRequest(goal)
	if err != nil {
		return nil, err
	}
	if onFeedback != nil {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		unsub := c.subscribe(ctx, &c.feedbackSubs, req.GetGoalID(), onFeedback)
		defer unsub()
	}
	resp, err := c.SendGoalRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if !resp.(goalResponseMessage).GetGoalAccepted() {
		return nil, errors.New("goal was rejected")
	}
	return c.GetResult(ctx, req.GetGoalID())
}

// SendGoal sends a new goal to the server and returns the status message of the
// goal. The ID for the goal is generated using a cryptographically secure
// random number generator.
//
// A non-nil error is returned only if the processing of the request itself
// failed. SendGoal returns normally even if the goal is rejected, and the
// status can be read from the returned response message.
//
// The type support of goal must be types.ActionTypeSupport.Goal().
//
// The type support of the returned message is
// types.ActionTypeSupport.SendGoal().Response().
func (c *ActionClient) SendGoal(ctx context.Context, goal types.Message) (types.Message, *types.GoalID, error) {
	msg, err := c.newSendGoalRequest(goal)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.SendGoalRequest(ctx, msg)
	return resp, msg.GetGoalID(), err
}

func (c *ActionClient) newSendGoalRequest(goal types.Message) (goalRequestMessage, error) {
	msg := c.typeSupport.SendGoal().Request().New().(goalRequestMessage)
	goalID := msg.GetGoalID()
	if _, err := rand.Read(goalID[:]); err != nil {
		return nil, fmt.Errorf("failed to generate goal ID: %v", err)
	}
	msg.SetGoalDescription(goal)
	return msg, nil
}

// SendGoalRequest sends a goal to the server and returns the status message of
// the goal.
//
// The type support of request must be types.ActionTypeSupport.SendGoal().Request().
//
// The type support of the returned message is
// types.ActionTypeSupport.SendGoal().Response().
//
// A non-nil error is returned only if the processing of the request itself
// failed. SendGoalRequest returns normally even if the goal is rejected, and
// the status can be read from the returned response message.
func (c *ActionClient) SendGoalRequest(ctx context.Context, request types.Message) (types.Message, error) {
	resp, _, err := c.goalSender.Send(ctx, request)
	return resp, err
}

func (c *ActionClient) sendGoalRequest(req unsafe.Pointer) (C.long, error) {
	var seqNum C.long
	rc := C.rcl_action_send_goal_request(&c.rclClient, req, &seqNum)
	if rc != C.RCL_RET_OK {
		return 0, errorsCastC(rc, "failed to send goal request")
	}
	return seqNum, nil
}

func (c *ActionClient) takeGoalResponse(resp unsafe.Pointer) (C.long, interface{}, error) {
	var header C.rmw_request_id_t
	rc := C.rcl_action_take_goal_response(&c.rclClient, &header, resp)
	if rc != C.RCL_RET_OK {
		return 0, nil, errorsCastC(rc, "failed to send goal request")
	}
	return header.sequence_number, nil, nil
}

// GetResult returns the result of the goal with goalID or an error if getting
// the result fails. If the goal has not yet reached a terminal state, GetResult
// waits for that to happen before returning.
//
// The type support of the returned message is
// types.ActionTypeSupport.GetResult().Response().
func (c *ActionClient) GetResult(ctx context.Context, goalID *types.GoalID) (types.Message, error) {
	msg := c.typeSupport.GetResult().Request().New()
	msg.(goalIDMessage).SetGoalID(goalID)
	resp, _, err := c.resultSender.Send(ctx, msg)
	return resp, err
}

func (c *ActionClient) sendResultRequest(req unsafe.Pointer) (C.long, error) {
	var seqNum C.long
	rc := C.rcl_action_send_result_request(&c.rclClient, req, &seqNum)
	if rc != C.RCL_RET_OK {
		return 0, errorsCastC(rc, "failed to send result request")
	}
	return seqNum, nil
}

func (c *ActionClient) takeResultResponse(resp unsafe.Pointer) (C.long, interface{}, error) {
	var header C.rmw_request_id_t
	rc := C.rcl_action_take_result_response(&c.rclClient, &header, resp)
	if rc != C.RCL_RET_OK {
		return 0, nil, errorsCastC(rc, "failed to send result request")
	}
	return header.sequence_number, nil, nil
}

// CancelGoal cancels goals.
//
// A non-nil error is returned only if the processing of the request itself
// failed. CancelGoal returns normally even if the canceling fails. The status
// can be read from the returned response message.
//
// The request includes a goal ID and a timestamp. If both the ID and the
// timestamp have zero values, all goals are canceled. If the ID is zero but the
// timestamp is not, all goals accepted at or before the timestamp are canceled.
// If the ID is not zero and the timestamp is zero, the goal with the specified
// ID is canceled. If both the ID and the timestamp are non-zero, the goal with
// the specified ID as well as all goals accepted at or before the timestamp are
// canceled.
//
// The type of request is action_msgs/srv/CancelGoal_Request.
//
// The type of the returned message is action_msgs/srv/CancelGoal_Response.
func (c *ActionClient) CancelGoal(ctx context.Context, request types.Message) (types.Message, error) {
	resp, _, err := c.cancelSender.Send(ctx, request)
	return resp, err
}

func (c *ActionClient) sendCancelRequest(req unsafe.Pointer) (C.long, error) {
	var seqNum C.long
	rc := C.rcl_action_send_cancel_request(&c.rclClient, req, &seqNum)
	if rc != C.RCL_RET_OK {
		return 0, errorsCastC(rc, "failed to send cancel request")
	}
	return seqNum, nil
}

func (c *ActionClient) takeCancelResponse(resp unsafe.Pointer) (C.long, interface{}, error) {
	var header C.rmw_request_id_t
	rc := C.rcl_action_take_cancel_response(&c.rclClient, &header, resp)
	if rc != C.RCL_RET_OK {
		return 0, nil, errorsCastC(rc, "failed to send cancel request")
	}
	return header.sequence_number, nil, nil
}

// WatchFeedback calls handler for every feedback message for the goal with id
// goalID. If goalID is nil, handler is called for all feedback messages
// regardless of which goal they belong to.
//
// WatchFeedback returns after the handler has been registered. The returned
// channel will receive exactly one error value, which may be nil, and then the
// channel is closed. Reading the value from the channel is not required.
// Watching can be stopped by canceling ctx.
//
// The type support of the message passed to handler is
// types.ActionTypeSupport.FeedbackMessage().
func (c *ActionClient) WatchFeedback(ctx context.Context, goalID *types.GoalID, handler FeedbackHandler) <-chan error {
	unsub := c.subscribe(ctx, &c.statusSubs, goalID, handler)
	errc := make(chan error, 1)
	go func() {
		defer unsub()
		<-ctx.Done()
		errc <- ctx.Err()
	}()
	return errc
}

func (c *ActionClient) subscribe(
	ctx context.Context,
	subs *actionClientSubs,
	id *types.GoalID,
	handler actionClientHandler,
) (unsubscribe func()) {
	c.rclClientMu.Lock()
	defer c.rclClientMu.Unlock()
	if id == nil {
		subID := c.nextSubscriberID
		c.nextSubscriberID++
		subs.allGoals[subID] = actionClientHandlerMapEntry{
			ctx:     ctx,
			handler: handler,
		}
		return func() {
			c.rclClientMu.Lock()
			defer c.rclClientMu.Unlock()
			if _, ok := subs.allGoals[subID]; ok {
				delete(subs.allGoals, subID)
			}
		}
	}
	goalID := *id
	handlers := subs.perGoal[goalID]
	if handlers == nil {
		handlers = make(actionClientHandlerMap)
		subs.perGoal[goalID] = handlers
	}
	subID := c.nextSubscriberID
	c.nextSubscriberID++
	handlers[subID] = actionClientHandlerMapEntry{
		ctx:     ctx,
		handler: handler,
	}
	return func() {
		c.rclClientMu.Lock()
		defer c.rclClientMu.Unlock()
		if handlers := subs.perGoal[goalID]; handlers != nil {
			if _, ok := handlers[subID]; ok {
				delete(handlers, subID)
				if len(handlers) == 0 {
					delete(subs.perGoal, goalID)
				}
			}
		}
	}
}

func (c *ActionClient) handleFeedback() {
	ts := c.typeSupport.FeedbackMessage()
	buf := ts.PrepareMemory()
	defer ts.ReleaseMemory(buf)
	rc := C.rcl_action_take_feedback(&c.rclClient, buf)
	if rc != C.RCL_RET_OK {
		c.node.Logger().Error(errorsCastC(rc, "failed to take feedback"))
		return
	}
	msg := ts.New().(goalIDMessage)
	ts.AsGoStruct(msg, buf)
	c.feedbackSubs.allGoals.call(msg)
	c.feedbackSubs.perGoal[*msg.GetGoalID()].call(msg)
}

// WatchStatus calls handler for every status message regarding the goal with id
// goalID. If goalID is nil, handler is called for all status messages
// regardless of which goal they belong to.
//
// WatchStatus returns after the handler has been registered. The returned
// channel will receive exactly one error value, which may be nil, and then the
// channel is closed. Reading the value from the channel is not required.
// Watching can be stopped by canceling ctx.
//
// The type of the message passed to handler will be action_msgs/msg/GoalStatus.
func (c *ActionClient) WatchStatus(ctx context.Context, goalID *types.GoalID, handler StatusHandler) <-chan error {
	unsub := c.subscribe(ctx, &c.statusSubs, goalID, handler)
	errc := make(chan error, 1)
	go func() {
		defer unsub()
		<-ctx.Done()
		errc <- ctx.Err()
	}()
	return errc
}

func (c *ActionClient) handleStatus() {
	ts := c.typeSupport.GoalStatusArray()
	buf := ts.PrepareMemory()
	defer ts.ReleaseMemory(buf)
	rc := C.rcl_action_take_status(&c.rclClient, buf)
	if rc != C.RCL_RET_OK {
		c.node.Logger().Error(errorsCastC(rc, "failed to take status"))
		return
	}
	msg := ts.New()
	ts.AsGoStruct(msg, buf)
	msg.(forEach).CallForEach(func(info interface{}) {
		msg := info.(goalIDMessage)
		c.statusSubs.perGoal[*msg.GetGoalID()].call(msg)
		c.statusSubs.allGoals.call(msg)
	})
}

func (c *ActionClient) handleReadyEntities(ws *WaitSet) {
	var feedback, status, goalResp, cancelResp, resultResp C.bool
	c.rclClientMu.Lock()
	defer c.rclClientMu.Unlock()
	rc := C.rcl_action_client_wait_set_get_entities_ready(
		&ws.rcl_wait_set_t,
		&c.rclClient,
		&feedback,
		&status,
		&goalResp,
		&cancelResp,
		&resultResp,
	)
	if rc != C.RCL_RET_OK {
		c.node.Logger().Error(errorsCastC(rc, "failed to get ready entities"))
		return
	}
	if feedback {
		c.handleFeedback()
	}
	if status {
		c.handleStatus()
	}
	if goalResp {
		c.goalSender.HandleResponse()
	}
	if cancelResp {
		c.cancelSender.HandleResponse()
	}
	if resultResp {
		c.resultSender.HandleResponse()
	}
}

func wrapErr(format string, err *error, a ...interface{}) {
	if *err != nil {
		args := make([]interface{}, 0, 1+len(a))
		args = append(args, *err)
		args = append(args, a...)
		*err = fmt.Errorf(format, args...)
	}
}
