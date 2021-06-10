/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package geometry_msgs_msg
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo/types"
	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lgeometry_msgs__rosidl_typesupport_c -lgeometry_msgs__rosidl_generator_c
#cgo LDFLAGS: -lstd_msgs__rosidl_typesupport_c -lstd_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>

#include <geometry_msgs/msg/accel_with_covariance_stamped.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("geometry_msgs/AccelWithCovarianceStamped", AccelWithCovarianceStampedTypeSupport)
}

// Do not create instances of this type directly. Always use NewAccelWithCovarianceStamped
// function instead.
type AccelWithCovarianceStamped struct {
	Header std_msgs_msg.Header `yaml:"header"`// This represents an estimated accel with reference coordinate frame and timestamp.
	Accel AccelWithCovariance `yaml:"accel"`// This represents an estimated accel with reference coordinate frame and timestamp.
}

// NewAccelWithCovarianceStamped creates a new AccelWithCovarianceStamped with default values.
func NewAccelWithCovarianceStamped() *AccelWithCovarianceStamped {
	self := AccelWithCovarianceStamped{}
	self.SetDefaults()
	return &self
}

func (t *AccelWithCovarianceStamped) Clone() *AccelWithCovarianceStamped {
	c := &AccelWithCovarianceStamped{}
	c.Header = *t.Header.Clone()
	c.Accel = *t.Accel.Clone()
	return c
}

func (t *AccelWithCovarianceStamped) CloneMsg() types.Message {
	return t.Clone()
}

func (t *AccelWithCovarianceStamped) SetDefaults() {
	t.Header.SetDefaults()
	t.Accel.SetDefaults()
}

// CloneAccelWithCovarianceStampedSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func CloneAccelWithCovarianceStampedSlice(dst, src []AccelWithCovarianceStamped) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var AccelWithCovarianceStampedTypeSupport types.MessageTypeSupport = _AccelWithCovarianceStampedTypeSupport{}

type _AccelWithCovarianceStampedTypeSupport struct{}

func (t _AccelWithCovarianceStampedTypeSupport) New() types.Message {
	return NewAccelWithCovarianceStamped()
}

func (t _AccelWithCovarianceStampedTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.geometry_msgs__msg__AccelWithCovarianceStamped
	return (unsafe.Pointer)(C.geometry_msgs__msg__AccelWithCovarianceStamped__create())
}

func (t _AccelWithCovarianceStampedTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.geometry_msgs__msg__AccelWithCovarianceStamped__destroy((*C.geometry_msgs__msg__AccelWithCovarianceStamped)(pointer_to_free))
}

func (t _AccelWithCovarianceStampedTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*AccelWithCovarianceStamped)
	mem := (*C.geometry_msgs__msg__AccelWithCovarianceStamped)(dst)
	std_msgs_msg.HeaderTypeSupport.AsCStruct(unsafe.Pointer(&mem.header), &m.Header)
	AccelWithCovarianceTypeSupport.AsCStruct(unsafe.Pointer(&mem.accel), &m.Accel)
}

func (t _AccelWithCovarianceStampedTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*AccelWithCovarianceStamped)
	mem := (*C.geometry_msgs__msg__AccelWithCovarianceStamped)(ros2_message_buffer)
	std_msgs_msg.HeaderTypeSupport.AsGoStruct(&m.Header, unsafe.Pointer(&mem.header))
	AccelWithCovarianceTypeSupport.AsGoStruct(&m.Accel, unsafe.Pointer(&mem.accel))
}

func (t _AccelWithCovarianceStampedTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__geometry_msgs__msg__AccelWithCovarianceStamped())
}

type CAccelWithCovarianceStamped = C.geometry_msgs__msg__AccelWithCovarianceStamped
type CAccelWithCovarianceStamped__Sequence = C.geometry_msgs__msg__AccelWithCovarianceStamped__Sequence

func AccelWithCovarianceStamped__Sequence_to_Go(goSlice *[]AccelWithCovarianceStamped, cSlice CAccelWithCovarianceStamped__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]AccelWithCovarianceStamped, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.geometry_msgs__msg__AccelWithCovarianceStamped__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__AccelWithCovarianceStamped * uintptr(i)),
		))
		AccelWithCovarianceStampedTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func AccelWithCovarianceStamped__Sequence_to_C(cSlice *CAccelWithCovarianceStamped__Sequence, goSlice []AccelWithCovarianceStamped) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.geometry_msgs__msg__AccelWithCovarianceStamped)(C.malloc((C.size_t)(C.sizeof_struct_geometry_msgs__msg__AccelWithCovarianceStamped * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.geometry_msgs__msg__AccelWithCovarianceStamped)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__AccelWithCovarianceStamped * uintptr(i)),
		))
		AccelWithCovarianceStampedTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func AccelWithCovarianceStamped__Array_to_Go(goSlice []AccelWithCovarianceStamped, cSlice []CAccelWithCovarianceStamped) {
	for i := 0; i < len(cSlice); i++ {
		AccelWithCovarianceStampedTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func AccelWithCovarianceStamped__Array_to_C(cSlice []CAccelWithCovarianceStamped, goSlice []AccelWithCovarianceStamped) {
	for i := 0; i < len(goSlice); i++ {
		AccelWithCovarianceStampedTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
