/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package example_interfaces_action
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	unique_identifier_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/unique_identifier_msgs/msg"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <example_interfaces/action/fibonacci.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("example_interfaces/Fibonacci_GetResult_Request", Fibonacci_GetResult_RequestTypeSupport)
	typemap.RegisterMessage("example_interfaces/action/Fibonacci_GetResult_Request", Fibonacci_GetResult_RequestTypeSupport)
}

// Do not create instances of this type directly. Always use NewFibonacci_GetResult_Request
// function instead.
type Fibonacci_GetResult_Request struct {
	GoalID unique_identifier_msgs_msg.UUID `yaml:"goal_id"`
}

// NewFibonacci_GetResult_Request creates a new Fibonacci_GetResult_Request with default values.
func NewFibonacci_GetResult_Request() *Fibonacci_GetResult_Request {
	self := Fibonacci_GetResult_Request{}
	self.SetDefaults()
	return &self
}

func (t *Fibonacci_GetResult_Request) Clone() *Fibonacci_GetResult_Request {
	c := &Fibonacci_GetResult_Request{}
	c.GoalID = *t.GoalID.Clone()
	return c
}

func (t *Fibonacci_GetResult_Request) CloneMsg() types.Message {
	return t.Clone()
}

func (t *Fibonacci_GetResult_Request) SetDefaults() {
	t.GoalID.SetDefaults()
}

func (t *Fibonacci_GetResult_Request) GetTypeSupport() types.MessageTypeSupport {
	return Fibonacci_GetResult_RequestTypeSupport
}
func (t *Fibonacci_GetResult_Request) GetGoalID() *types.GoalID {
	return (*types.GoalID)(&t.GoalID.Uuid)
}

func (t *Fibonacci_GetResult_Request) SetGoalID(id *types.GoalID) {
	t.GoalID.Uuid = *id
}

// Fibonacci_GetResult_RequestPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type Fibonacci_GetResult_RequestPublisher struct {
	*rclgo.Publisher
}

// NewFibonacci_GetResult_RequestPublisher creates and returns a new publisher for the
// Fibonacci_GetResult_Request
func NewFibonacci_GetResult_RequestPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*Fibonacci_GetResult_RequestPublisher, error) {
	pub, err := node.NewPublisher(topic_name, Fibonacci_GetResult_RequestTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &Fibonacci_GetResult_RequestPublisher{pub}, nil
}

func (p *Fibonacci_GetResult_RequestPublisher) Publish(msg *Fibonacci_GetResult_Request) error {
	return p.Publisher.Publish(msg)
}

// Fibonacci_GetResult_RequestSubscription wraps rclgo.Subscription to provide type safe helper
// functions
type Fibonacci_GetResult_RequestSubscription struct {
	*rclgo.Subscription
}

// Fibonacci_GetResult_RequestSubscriptionCallback type is used to provide a subscription
// handler function for a Fibonacci_GetResult_RequestSubscription.
type Fibonacci_GetResult_RequestSubscriptionCallback func(msg *Fibonacci_GetResult_Request, info *rclgo.RmwMessageInfo, err error)

// NewFibonacci_GetResult_RequestSubscription creates and returns a new subscription for the
// Fibonacci_GetResult_Request
func NewFibonacci_GetResult_RequestSubscription(node *rclgo.Node, topic_name string, subscriptionCallback Fibonacci_GetResult_RequestSubscriptionCallback) (*Fibonacci_GetResult_RequestSubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg Fibonacci_GetResult_Request
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, Fibonacci_GetResult_RequestTypeSupport, callback)
	if err != nil {
		return nil, err
	}
	return &Fibonacci_GetResult_RequestSubscription{sub}, nil
}

func (s *Fibonacci_GetResult_RequestSubscription) TakeMessage(out *Fibonacci_GetResult_Request) (*rclgo.RmwMessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneFibonacci_GetResult_RequestSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneFibonacci_GetResult_RequestSlice(dst, src []Fibonacci_GetResult_Request) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var Fibonacci_GetResult_RequestTypeSupport types.MessageTypeSupport = _Fibonacci_GetResult_RequestTypeSupport{}

type _Fibonacci_GetResult_RequestTypeSupport struct{}

func (t _Fibonacci_GetResult_RequestTypeSupport) New() types.Message {
	return NewFibonacci_GetResult_Request()
}

func (t _Fibonacci_GetResult_RequestTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.example_interfaces__action__Fibonacci_GetResult_Request
	return (unsafe.Pointer)(C.example_interfaces__action__Fibonacci_GetResult_Request__create())
}

func (t _Fibonacci_GetResult_RequestTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.example_interfaces__action__Fibonacci_GetResult_Request__destroy((*C.example_interfaces__action__Fibonacci_GetResult_Request)(pointer_to_free))
}

func (t _Fibonacci_GetResult_RequestTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*Fibonacci_GetResult_Request)
	mem := (*C.example_interfaces__action__Fibonacci_GetResult_Request)(dst)
	unique_identifier_msgs_msg.UUIDTypeSupport.AsCStruct(unsafe.Pointer(&mem.goal_id), &m.GoalID)
}

func (t _Fibonacci_GetResult_RequestTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*Fibonacci_GetResult_Request)
	mem := (*C.example_interfaces__action__Fibonacci_GetResult_Request)(ros2_message_buffer)
	unique_identifier_msgs_msg.UUIDTypeSupport.AsGoStruct(&m.GoalID, unsafe.Pointer(&mem.goal_id))
}

func (t _Fibonacci_GetResult_RequestTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__example_interfaces__action__Fibonacci_GetResult_Request())
}

type CFibonacci_GetResult_Request = C.example_interfaces__action__Fibonacci_GetResult_Request
type CFibonacci_GetResult_Request__Sequence = C.example_interfaces__action__Fibonacci_GetResult_Request__Sequence

func Fibonacci_GetResult_Request__Sequence_to_Go(goSlice *[]Fibonacci_GetResult_Request, cSlice CFibonacci_GetResult_Request__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Fibonacci_GetResult_Request, cSlice.size)
	src := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range src {
		Fibonacci_GetResult_RequestTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(&src[i]))
	}
}
func Fibonacci_GetResult_Request__Sequence_to_C(cSlice *CFibonacci_GetResult_Request__Sequence, goSlice []Fibonacci_GetResult_Request) {
	if len(goSlice) == 0 {
		cSlice.data = nil
		cSlice.capacity = 0
		cSlice.size = 0
		return
	}
	cSlice.data = (*C.example_interfaces__action__Fibonacci_GetResult_Request)(C.malloc(C.sizeof_struct_example_interfaces__action__Fibonacci_GetResult_Request * C.size_t(len(goSlice))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity
	dst := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range goSlice {
		Fibonacci_GetResult_RequestTypeSupport.AsCStruct(unsafe.Pointer(&dst[i]), &goSlice[i])
	}
}
func Fibonacci_GetResult_Request__Array_to_Go(goSlice []Fibonacci_GetResult_Request, cSlice []CFibonacci_GetResult_Request) {
	for i := 0; i < len(cSlice); i++ {
		Fibonacci_GetResult_RequestTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func Fibonacci_GetResult_Request__Array_to_C(cSlice []CFibonacci_GetResult_Request, goSlice []Fibonacci_GetResult_Request) {
	for i := 0; i < len(goSlice); i++ {
		Fibonacci_GetResult_RequestTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
