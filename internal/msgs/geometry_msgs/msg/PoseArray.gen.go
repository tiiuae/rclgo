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

#include <geometry_msgs/msg/pose_array.h>

*/
import "C"

func init() {
	typemap.RegisterMessage("geometry_msgs/PoseArray", PoseArrayTypeSupport)
}

// Do not create instances of this type directly. Always use NewPoseArray
// function instead.
type PoseArray struct {
	Header std_msgs_msg.Header `yaml:"header"`
	Poses []Pose `yaml:"poses"`
}

// NewPoseArray creates a new PoseArray with default values.
func NewPoseArray() *PoseArray {
	self := PoseArray{}
	self.SetDefaults()
	return &self
}

func (t *PoseArray) Clone() *PoseArray {
	c := &PoseArray{}
	c.Header = *t.Header.Clone()
	if t.Poses != nil {
		c.Poses = make([]Pose, len(t.Poses))
		ClonePoseSlice(c.Poses, t.Poses)
	}
	return c
}

func (t *PoseArray) CloneMsg() types.Message {
	return t.Clone()
}

func (t *PoseArray) SetDefaults() {
	t.Header.SetDefaults()
	t.Poses = nil
}

// ClonePoseArraySlice clones src to dst by calling Clone for each element in
// src. Panics if len(dst) < len(src).
func ClonePoseArraySlice(dst, src []PoseArray) {
	for i := range src {
		dst[i] = *src[i].Clone()
	}
}

// Modifying this variable is undefined behavior.
var PoseArrayTypeSupport types.MessageTypeSupport = _PoseArrayTypeSupport{}

type _PoseArrayTypeSupport struct{}

func (t _PoseArrayTypeSupport) New() types.Message {
	return NewPoseArray()
}

func (t _PoseArrayTypeSupport) PrepareMemory() unsafe.Pointer { //returns *C.geometry_msgs__msg__PoseArray
	return (unsafe.Pointer)(C.geometry_msgs__msg__PoseArray__create())
}

func (t _PoseArrayTypeSupport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.geometry_msgs__msg__PoseArray__destroy((*C.geometry_msgs__msg__PoseArray)(pointer_to_free))
}

func (t _PoseArrayTypeSupport) AsCStruct(dst unsafe.Pointer, msg types.Message) {
	m := msg.(*PoseArray)
	mem := (*C.geometry_msgs__msg__PoseArray)(dst)
	std_msgs_msg.HeaderTypeSupport.AsCStruct(unsafe.Pointer(&mem.header), &m.Header)
	Pose__Sequence_to_C(&mem.poses, m.Poses)
}

func (t _PoseArrayTypeSupport) AsGoStruct(msg types.Message, ros2_message_buffer unsafe.Pointer) {
	m := msg.(*PoseArray)
	mem := (*C.geometry_msgs__msg__PoseArray)(ros2_message_buffer)
	std_msgs_msg.HeaderTypeSupport.AsGoStruct(&m.Header, unsafe.Pointer(&mem.header))
	Pose__Sequence_to_Go(&m.Poses, mem.poses)
}

func (t _PoseArrayTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__geometry_msgs__msg__PoseArray())
}

type CPoseArray = C.geometry_msgs__msg__PoseArray
type CPoseArray__Sequence = C.geometry_msgs__msg__PoseArray__Sequence

func PoseArray__Sequence_to_Go(goSlice *[]PoseArray, cSlice CPoseArray__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]PoseArray, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.geometry_msgs__msg__PoseArray__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__PoseArray * uintptr(i)),
		))
		PoseArrayTypeSupport.AsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}
func PoseArray__Sequence_to_C(cSlice *CPoseArray__Sequence, goSlice []PoseArray) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.geometry_msgs__msg__PoseArray)(C.malloc((C.size_t)(C.sizeof_struct_geometry_msgs__msg__PoseArray * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.geometry_msgs__msg__PoseArray)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_geometry_msgs__msg__PoseArray * uintptr(i)),
		))
		PoseArrayTypeSupport.AsCStruct(unsafe.Pointer(cIdx), &v)
	}
}
func PoseArray__Array_to_Go(goSlice []PoseArray, cSlice []CPoseArray) {
	for i := 0; i < len(cSlice); i++ {
		PoseArrayTypeSupport.AsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}
func PoseArray__Array_to_C(cSlice []CPoseArray, goSlice []PoseArray) {
	for i := 0; i < len(goSlice); i++ {
		PoseArrayTypeSupport.AsCStruct(unsafe.Pointer(&cSlice[i]), &goSlice[i])
	}
}
