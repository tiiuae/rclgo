/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package std_msgs_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	builtin_interfaces_msg "github.com/tiiuae/rclgo/internal/msgs/builtin_interfaces/msg"
	primitives "github.com/tiiuae/rclgo/pkg/rclgo/primitives"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lstd_msgs__rosidl_typesupport_c -lstd_msgs__rosidl_generator_c
#cgo LDFLAGS: -lbuiltin_interfaces__rosidl_typesupport_c -lbuiltin_interfaces__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>

#include <std_msgs/msg/header.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("std_msgs/Header", HeaderTypeSupport)
}

// Do not create instances of this type directly. Always use NewHeader
// function instead.
type Header struct {
	Stamp builtin_interfaces_msg.Time `yaml:"stamp"`// Two-integer timestamp that is expressed as seconds and nanoseconds.
	FrameId string `yaml:"frame_id"`// Transform frame with which this data is associated.
}

// NewHeader creates a new Header with default values.
func NewHeader() *Header {
	self := Header{}
	self.SetDefaults()
	return &self
}

func (t *Header) Clone() *Header {
	c := &Header{}
	c.Stamp = *t.Stamp.Clone()
	c.FrameId = t.FrameId
	return c
}

func (t *Header) CloneMsg() types.Message {
	return t.Clone()
}

func (t *Header) SetDefaults() {
	t.Stamp.SetDefaults()
	t.FrameId = ""
}

// CloneHeaderSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneHeaderSlice(dst, src []Header) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var HeaderTypeSupport types.MessageTypeSupport = _HeaderTypeSupport{}

type _HeaderTypeSupport struct{}

func (t _HeaderTypeSupport) New() types.Message {
	return NewHeader()
}

func (t _HeaderTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.std_msgs__msg__Header
	return (unsafe.Pointer)(C.std_msgs__msg__Header__create())
}

func (t _HeaderTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.std_msgs__msg__Header__destroy((*C.std_msgs__msg__Header)(pointer_to_free))
}

func (t _HeaderTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*Header)
	mem := (*C.std_msgs__msg__Header)(dst)
	builtin_interfaces_msg.TimeTypeSupport.AsCStruct(unsafe.Pointer(&mem.stamp), &m.Stamp)
	primitives.StringAsCStruct(unsafe.Pointer(&mem.frame_id), m.FrameId)
}

func (t _HeaderTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*Header)
	mem := (*C.std_msgs__msg__Header)(ros2_message_buffer)
	builtin_interfaces_msg.TimeTypeSupport.AsGoStruct(&m.Stamp, unsafe.Pointer(&mem.stamp))
	primitives.StringAsGoStruct(&m.FrameId, unsafe.Pointer(&mem.frame_id))
}

func (t _HeaderTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__std_msgs__msg__Header())
}

type CHeader = C.std_msgs__msg__Header
type CHeader__Sequence = C.std_msgs__msg__Header__Sequence

func Header__Sequence_to_Go(goSlice *[]Header, cSlice CHeader__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Header, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.std_msgs__msg__Header__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_std_msgs__msg__Header * uintptr(i)),
		))
		HeaderTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func Header__Sequence_to_C(cSlice *CHeader__Sequence, goSlice []Header) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.std_msgs__msg__Header)(C.malloc((C.size_t)(C.sizeof_struct_std_msgs__msg__Header * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.std_msgs__msg__Header)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_std_msgs__msg__Header * uintptr(i)),
		))
		HeaderTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func Header__Array_to_Go(goSlice []Header, cSlice []CHeader) {
	for i := 0; i < len(cSlice); i++ {
		HeaderTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func Header__Array_to_C(cSlice []CHeader, goSlice []Header) {
	for i := 0; i < len(goSlice); i++ {
		HeaderTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
