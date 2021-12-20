/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package test_msgs_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <test_msgs/msg/empty.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("test_msgs/Empty", EmptyTypeSupport)
}

// Do not create instances of this type directly. Always use NewEmpty
// function instead.
type Empty struct {
}

// NewEmpty creates a new Empty with default values.
func NewEmpty() *Empty {
	self := Empty{}
	self.SetDefaults()
	return &self
}

func (t *Empty) Clone() *Empty {
	c := &Empty{}
	return c
}

func (t *Empty) CloneMsg() types.Message {
	return t.Clone()
}

func (t *Empty) SetDefaults() {
}

// EmptyPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type EmptyPublisher struct {
	*rclgo.Publisher
}

// NewEmptyPublisher creates and returns a new publisher for the
// Empty
func NewEmptyPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*EmptyPublisher, error) {
	pub, err := node.NewPublisher(topic_name, EmptyTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &EmptyPublisher{pub}, nil
}

func (p *EmptyPublisher) Publish(msg *Empty) error {
	return p.Publisher.Publish(msg)
}

// EmptySubscription wraps rclgo.Subscription to provide type safe helper
// functions
type EmptySubscription struct {
	*rclgo.Subscription
}

// EmptySubscriptionCallback type is used to provide a subscription
// handler function for a EmptySubscription.
type EmptySubscriptionCallback func(msg *Empty, info *rclgo.RmwMessageInfo, err error)

// NewEmptySubscription creates and returns a new subscription for the
// Empty
func NewEmptySubscription(node *rclgo.Node, topic_name string, subscriptionCallback EmptySubscriptionCallback) (*EmptySubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg Empty
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, EmptyTypeSupport, callback)
	if err != nil {
		return nil, err
	}
	return &EmptySubscription{sub}, nil
}

func (s *EmptySubscription) TakeMessage(out *Empty) (*rclgo.RmwMessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneEmptySlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneEmptySlice(dst, src []Empty) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var EmptyTypeSupport types.MessageTypeSupport = _EmptyTypeSupport{}

type _EmptyTypeSupport struct{}

func (t _EmptyTypeSupport) New() types.Message {
	return NewEmpty()
}

func (t _EmptyTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.test_msgs__msg__Empty
	return (unsafe.Pointer)(C.test_msgs__msg__Empty__create())
}

func (t _EmptyTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.test_msgs__msg__Empty__destroy((*C.test_msgs__msg__Empty)(pointer_to_free))
}

func (t _EmptyTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	
}

func (t _EmptyTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	
}

func (t _EmptyTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__test_msgs__msg__Empty())
}

type CEmpty = C.test_msgs__msg__Empty
type CEmpty__Sequence = C.test_msgs__msg__Empty__Sequence

func Empty__Sequence_to_Go(goSlice *[]Empty, cSlice CEmpty__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Empty, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.test_msgs__msg__Empty__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_test_msgs__msg__Empty * uintptr(i)),
		))
		EmptyTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func Empty__Sequence_to_C(cSlice *CEmpty__Sequence, goSlice []Empty) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.test_msgs__msg__Empty)(C.malloc((C.size_t)(C.sizeof_struct_test_msgs__msg__Empty * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.test_msgs__msg__Empty)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_test_msgs__msg__Empty * uintptr(i)),
		))
		EmptyTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func Empty__Array_to_Go(goSlice []Empty, cSlice []CEmpty) {
	for i := 0; i < len(cSlice); i++ {
		EmptyTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func Empty__Array_to_C(cSlice []CEmpty, goSlice []Empty) {
	for i := 0; i < len(goSlice); i++ {
		EmptyTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
