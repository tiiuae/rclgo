/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package std_msgs_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	primitives "github.com/tiiuae/rclgo/pkg/rclgo/primitives"
	
)
/*
#include <rosidl_runtime_c/message_type_support_struct.h>

#include <std_msgs/msg/int64_multi_array.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("std_msgs/Int64MultiArray", Int64MultiArrayTypeSupport)
}

// Do not create instances of this type directly. Always use NewInt64MultiArray
// function instead.
type Int64MultiArray struct {
	Layout MultiArrayLayout `yaml:"layout"`// specification of data layout
	Data []int64 `yaml:"data"`// array of data
}

// NewInt64MultiArray creates a new Int64MultiArray with default values.
func NewInt64MultiArray() *Int64MultiArray {
	self := Int64MultiArray{}
	self.SetDefaults()
	return &self
}

func (t *Int64MultiArray) Clone() *Int64MultiArray {
	c := &Int64MultiArray{}
	c.Layout = *t.Layout.Clone()
	if t.Data != nil {
		c.Data = make([]int64, len(t.Data))
		copy(c.Data, t.Data)
	}
	return c
}

func (t *Int64MultiArray) CloneMsg() types.Message {
	return t.Clone()
}

func (t *Int64MultiArray) SetDefaults() {
	t.Layout.SetDefaults()
	t.Data = nil
}

// Int64MultiArrayPublisher wraps rclgo.Publisher to provide type safe helper
// functions
type Int64MultiArrayPublisher struct {
	*rclgo.Publisher
}

// NewInt64MultiArrayPublisher creates and returns a new publisher for the
// Int64MultiArray
func NewInt64MultiArrayPublisher(node *rclgo.Node, topic_name string, options *rclgo.PublisherOptions) (*Int64MultiArrayPublisher, error) {
	pub, err := node.NewPublisher(topic_name, Int64MultiArrayTypeSupport, options)
	if err != nil {
		return nil, err
	}
	return &Int64MultiArrayPublisher{pub}, nil
}

func (p *Int64MultiArrayPublisher) Publish(msg *Int64MultiArray) error {
	return p.Publisher.Publish(msg)
}

// Int64MultiArraySubscription wraps rclgo.Subscription to provide type safe helper
// functions
type Int64MultiArraySubscription struct {
	*rclgo.Subscription
}

// Int64MultiArraySubscriptionCallback type is used to provide a subscription
// handler function for a Int64MultiArraySubscription.
type Int64MultiArraySubscriptionCallback func(msg *Int64MultiArray, info *rclgo.RmwMessageInfo, err error)

// NewInt64MultiArraySubscription creates and returns a new subscription for the
// Int64MultiArray
func NewInt64MultiArraySubscription(node *rclgo.Node, topic_name string, subscriptionCallback Int64MultiArraySubscriptionCallback) (*Int64MultiArraySubscription, error) {
	callback := func(s *rclgo.Subscription) {
		var msg Int64MultiArray
		info, err := s.TakeMessage(&msg)
		subscriptionCallback(&msg, info, err)
	}
	sub, err := node.NewSubscription(topic_name, Int64MultiArrayTypeSupport, callback)
	if err != nil {
		return nil, err
	}
	return &Int64MultiArraySubscription{sub}, nil
}

func (s *Int64MultiArraySubscription) TakeMessage(out *Int64MultiArray) (*rclgo.RmwMessageInfo, error) {
	return s.Subscription.TakeMessage(out)
}

// CloneInt64MultiArraySlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneInt64MultiArraySlice(dst, src []Int64MultiArray) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var Int64MultiArrayTypeSupport types.MessageTypeSupport = _Int64MultiArrayTypeSupport{}

type _Int64MultiArrayTypeSupport struct{}

func (t _Int64MultiArrayTypeSupport) New() types.Message {
	return NewInt64MultiArray()
}

func (t _Int64MultiArrayTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.std_msgs__msg__Int64MultiArray
	return (unsafe.Pointer)(C.std_msgs__msg__Int64MultiArray__create())
}

func (t _Int64MultiArrayTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.std_msgs__msg__Int64MultiArray__destroy((*C.std_msgs__msg__Int64MultiArray)(pointer_to_free))
}

func (t _Int64MultiArrayTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*Int64MultiArray)
	mem := (*C.std_msgs__msg__Int64MultiArray)(dst)
	MultiArrayLayoutTypeSupport.AsCStruct(unsafe.Pointer(&mem.layout), &m.Layout)
	primitives.Int64__Sequence_to_C((*primitives.CInt64__Sequence)(unsafe.Pointer(&mem.data)), m.Data)
}

func (t _Int64MultiArrayTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*Int64MultiArray)
	mem := (*C.std_msgs__msg__Int64MultiArray)(ros2_message_buffer)
	MultiArrayLayoutTypeSupport.AsGoStruct(&m.Layout, unsafe.Pointer(&mem.layout))
	primitives.Int64__Sequence_to_Go(&m.Data, *(*primitives.CInt64__Sequence)(unsafe.Pointer(&mem.data)))
}

func (t _Int64MultiArrayTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__Int64MultiArray())
}

type CInt64MultiArray = C.std_msgs__msg__Int64MultiArray
type CInt64MultiArray__Sequence = C.std_msgs__msg__Int64MultiArray__Sequence

func Int64MultiArray__Sequence_to_Go(goSlice *[]Int64MultiArray, cSlice CInt64MultiArray__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Int64MultiArray, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.std_msgs__msg__Int64MultiArray__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_std_msgs__msg__Int64MultiArray * uintptr(i)),
		))
		Int64MultiArrayTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func Int64MultiArray__Sequence_to_C(cSlice *CInt64MultiArray__Sequence, goSlice []Int64MultiArray) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.std_msgs__msg__Int64MultiArray)(C.malloc((C.size_t)(C.sizeof_struct_std_msgs__msg__Int64MultiArray * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.std_msgs__msg__Int64MultiArray)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_std_msgs__msg__Int64MultiArray * uintptr(i)),
		))
		Int64MultiArrayTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func Int64MultiArray__Array_to_Go(goSlice []Int64MultiArray, cSlice []CInt64MultiArray) {
	for i := 0; i < len(cSlice); i++ {
		Int64MultiArrayTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func Int64MultiArray__Array_to_C(cSlice []CInt64MultiArray, goSlice []Int64MultiArray) {
	for i := 0; i < len(goSlice); i++ {
		Int64MultiArrayTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
