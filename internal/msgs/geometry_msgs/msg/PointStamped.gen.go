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

#include <geometry_msgs/msg/point_stamped.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("geometry_msgs/PointStamped", PointStampedTypeSupport)
}

// Do not create instances of this type directly. Always use NewPointStamped
// function instead.
type PointStamped struct {
	Header std_msgs_msg.Header `yaml:"header"`
	Point Point `yaml:"point"`
}

// NewPointStamped creates a new PointStamped with default values.
func NewPointStamped() *PointStamped {
	self := PointStamped{}
	self.SetDefaults()
	return &self
}

func (t *PointStamped) Clone() *PointStamped {
	c := &PointStamped{}
	c.Header = *t.Header.Clone()
	c.Point = *t.Point.Clone()
	return c
}

func (t *PointStamped) CloneMsg() types.Message {
	return t.Clone()
}

func (t *PointStamped) SetDefaults() {
	t.Header.SetDefaults()
	t.Point.SetDefaults()
	
}

// ClonePointStampedSlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func ClonePointStampedSlice(dst, src []PointStamped) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var PointStampedTypeSupport types.MessageTypeSupport = _PointStampedTypeSupport{}

type _PointStampedTypeSupport struct{}

func (t _PointStampedTypeSupport) New() types.Message {
	return NewPointStamped()
}

func (t _PointStampedTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.geometry_msgs__msg__PointStamped
	return (unsafe.Pointer)(C.geometry_msgs__msg__PointStamped__create())
}

func (t _PointStampedTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.geometry_msgs__msg__PointStamped__destroy((*C.geometry_msgs__msg__PointStamped)(pointer_to_free))
}

func (t _PointStampedTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*PointStamped)
	mem := (*C.geometry_msgs__msg__PointStamped)(dst)
	std_msgs_msg.HeaderTypeSupport.AsCStruct(unsafe.Pointer(&mem.header), &m.Header)
	PointTypeSupport.AsCStruct(unsafe.Pointer(&mem.point), &m.Point)
}

func (t _PointStampedTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*PointStamped)
	mem := (*C.geometry_msgs__msg__PointStamped)(ros2_message_buffer)
	std_msgs_msg.HeaderTypeSupport.AsGoStruct(&m.Header, unsafe.Pointer(&mem.header))
	PointTypeSupport.AsGoStruct(&m.Point, unsafe.Pointer(&mem.point))
}

func (t _PointStampedTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__geometry_msgs__msg__PointStamped())
}

type CPointStamped = C.geometry_msgs__msg__PointStamped
type CPointStamped__Sequence = C.geometry_msgs__msg__PointStamped__Sequence

func PointStamped__Sequence_to_Go(goSlice *[]PointStamped, cSlice CPointStamped__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]PointStamped, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.geometry_msgs__msg__PointStamped__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__PointStamped * uintptr(i)),
		))
		PointStampedTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func PointStamped__Sequence_to_C(cSlice *CPointStamped__Sequence, goSlice []PointStamped) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.geometry_msgs__msg__PointStamped)(C.malloc((C.size_t)(C.sizeof_struct_geometry_msgs__msg__PointStamped * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.geometry_msgs__msg__PointStamped)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__PointStamped * uintptr(i)),
		))
		PointStampedTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func PointStamped__Array_to_Go(goSlice []PointStamped, cSlice []CPointStamped) {
	for i := 0; i < len(cSlice); i++ {
		PointStampedTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func PointStamped__Array_to_C(cSlice []CPointStamped, goSlice []PointStamped) {
	for i := 0; i < len(goSlice); i++ {
		PointStampedTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}