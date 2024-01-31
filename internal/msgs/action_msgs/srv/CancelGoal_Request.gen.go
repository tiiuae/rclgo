/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package action_msgs_srv
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	action_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/action_msgs/msg"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <action_msgs/srv/cancel_goal.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("action_msgs/CancelGoal_Request", CancelGoal_RequestTypeSupport)
	typemap.RegisterMessage("action_msgs/srv/CancelGoal_Request", CancelGoal_RequestTypeSupport)
}

type CancelGoal_Request struct {
	GoalInfo action_msgs_msg.GoalInfo `yaml:"goal_info"`// Goal info describing the goals to cancel, see above.
}

// NewCancelGoal_Request creates a new CancelGoal_Request with default values.
func NewCancelGoal_Request() *CancelGoal_Request {
	self := CancelGoal_Request{}
	self.SetDefaults()
	return &self
}

func (t *CancelGoal_Request) Clone() *CancelGoal_Request {
	c := &CancelGoal_Request{}
	c.GoalInfo = *t.GoalInfo.Clone()
	return c
}

func (t *CancelGoal_Request) CloneMsg() types.Message {
	return t.Clone()
}

func (t *CancelGoal_Request) SetDefaults() {
	t.GoalInfo.SetDefaults()
}

func (t *CancelGoal_Request) GetTypeSupport() types.MessageTypeSupport {
	return CancelGoal_RequestTypeSupport
}
func (t *CancelGoal_Request) GetGoalID() *types.GoalID {
	return (*types.GoalID)(&t.GoalInfo.GoalId.Uuid)
}

func (t *CancelGoal_Request) SetGoalID(id *types.GoalID) {
	t.GoalInfo.GoalId.Uuid = *id
}

// CancelGoal_RequestPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type CancelGoal_RequestPublisher struct {
	*rclgo.Publisher
}

// NewCancelGoal_RequestPublisher creates and returns a new publisher for the
// CancelGoal_Request
func NewCancelGoal_RequestPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*CancelGoal_RequestPublisher, error) {
	pub, err := node.NewPublisher(topic_name, CancelGoal_RequestTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &CancelGoal_RequestPublisher{pub}, nil
}

func (p *CancelGoal_RequestPublisher) Publish(msg *CancelGoal_Request) error {
	return p.Publisher.Publish(msg)
}

// CancelGoal_RequestSubscription wraps rclgo.Subscription to provide type safe helper
// functions
type CancelGoal_RequestSubscription struct {
	*rclgo.Subscription
}

// CancelGoal_RequestSubscriptionCallback type is used to provide a subscription
// handler function for a CancelGoal_RequestSubscription.
type CancelGoal_RequestSubscriptionCallback func(msg *CancelGoal_Request, info *rclgo.MessageInfo, err error)

// NewCancelGoal_RequestSubscription creates and returns a new subscription for the
// CancelGoal_Request
func NewCancelGoal_RequestSubscription(node *rclgo.Node, topic_name string, opts *rclgo.SubscriptionOptions, subscriptionCallback CancelGoal_RequestSubscriptionCallback) (*CancelGoal_RequestSubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg CancelGoal_Request
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, CancelGoal_RequestTypeSupport, opts, callback)
	if err != nil {
		return nil, err
	}
	return &CancelGoal_RequestSubscription{sub}, nil
}

func (s *CancelGoal_RequestSubscription) TakeMessage(out *CancelGoal_Request) (*rclgo.MessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneCancelGoal_RequestSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneCancelGoal_RequestSlice(dst, src []CancelGoal_Request) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var CancelGoal_RequestTypeSupport types.MessageTypeSupport = _CancelGoal_RequestTypeSupport{}

type _CancelGoal_RequestTypeSupport struct{}

func (t _CancelGoal_RequestTypeSupport) New() types.Message {
	return NewCancelGoal_Request()
}

func (t _CancelGoal_RequestTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.action_msgs__srv__CancelGoal_Request
	return (unsafe.Pointer)(C.action_msgs__srv__CancelGoal_Request__create())
}

func (t _CancelGoal_RequestTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.action_msgs__srv__CancelGoal_Request__destroy((*C.action_msgs__srv__CancelGoal_Request)(pointer_to_free))
}

func (t _CancelGoal_RequestTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*CancelGoal_Request)
	mem := (*C.action_msgs__srv__CancelGoal_Request)(dst)
	action_msgs_msg.GoalInfoTypeSupport.AsCStruct(unsafe.Pointer(&mem.goal_info), &m.GoalInfo)
}

func (t _CancelGoal_RequestTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*CancelGoal_Request)
	mem := (*C.action_msgs__srv__CancelGoal_Request)(ros2_message_buffer)
	action_msgs_msg.GoalInfoTypeSupport.AsGoStruct(&m.GoalInfo, unsafe.Pointer(&mem.goal_info))
}

func (t _CancelGoal_RequestTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__action_msgs__srv__CancelGoal_Request())
}

type CCancelGoal_Request = C.action_msgs__srv__CancelGoal_Request
type CCancelGoal_Request__Sequence = C.action_msgs__srv__CancelGoal_Request__Sequence

func CancelGoal_Request__Sequence_to_Go(goSlice *[]CancelGoal_Request, cSlice CCancelGoal_Request__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]CancelGoal_Request, cSlice.size)
	src := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range src {
		CancelGoal_RequestTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(&src[i]))
	}
}
func CancelGoal_Request__Sequence_to_C(cSlice *CCancelGoal_Request__Sequence, goSlice []CancelGoal_Request) {
	if len(goSlice) == 0 {
		cSlice.data = nil
		cSlice.capacity = 0
		cSlice.size = 0
		return
	}
	cSlice.data = (*C.action_msgs__srv__CancelGoal_Request)(C.malloc(C.sizeof_struct_action_msgs__srv__CancelGoal_Request * C.size_t(len(goSlice))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity
	dst := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range goSlice {
		CancelGoal_RequestTypeSupport.AsCStruct(unsafe.Pointer(&dst[i]), &goSlice[i])
	}
}
func CancelGoal_Request__Array_to_Go(goSlice []CancelGoal_Request, cSlice []CCancelGoal_Request) {
	for i := 0; i < len(cSlice); i++ {
		CancelGoal_RequestTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func CancelGoal_Request__Array_to_C(cSlice []CCancelGoal_Request, goSlice []CancelGoal_Request) {
	for i := 0; i < len(goSlice); i++ {
		CancelGoal_RequestTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
