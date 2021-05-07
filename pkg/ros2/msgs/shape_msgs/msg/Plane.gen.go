/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package shape_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	rosidl_runtime_c "github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lshape_msgs__rosidl_typesupport_c -lshape_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <shape_msgs/msg/plane.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("shape_msgs/Plane", &Plane{})
}

// Do not create instances of this type directly. Always use NewPlane
// function instead.
type Plane struct {
	Coef [4]float64 `yaml:"coef"`// Representation of a plane, using the plane equation ax + by + cz + d = 0.a := coef[0]b := coef[1]c := coef[2]d := coef[3]
}

// NewPlane creates a new Plane with default values.
func NewPlane() *Plane {
	self := Plane{}
	self.SetDefaults(nil)
	return &self
}

func (t *Plane) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *Plane) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__shape_msgs__msg__Plane())
}
func (t *Plane) PrepareMemory() unsafe.Pointer { //returns *C.shape_msgs__msg__Plane
	return (unsafe.Pointer)(C.shape_msgs__msg__Plane__create())
}
func (t *Plane) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.shape_msgs__msg__Plane__destroy((*C.shape_msgs__msg__Plane)(pointer_to_free))
}
func (t *Plane) AsCStruct() unsafe.Pointer {
	mem := (*C.shape_msgs__msg__Plane)(t.PrepareMemory())
	cSlice_coef := mem.coef[:]
	rosidl_runtime_c.Float64__Array_to_C(*(*[]rosidl_runtime_c.CFloat64)(unsafe.Pointer(&cSlice_coef)), t.Coef[:])
	return unsafe.Pointer(mem)
}
func (t *Plane) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.shape_msgs__msg__Plane)(ros2_message_buffer)
	cSlice_coef := mem.coef[:]
	rosidl_runtime_c.Float64__Array_to_Go(t.Coef[:], *(*[]rosidl_runtime_c.CFloat64)(unsafe.Pointer(&cSlice_coef)))
}
func (t *Plane) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CPlane = C.shape_msgs__msg__Plane
type CPlane__Sequence = C.shape_msgs__msg__Plane__Sequence

func Plane__Sequence_to_Go(goSlice *[]Plane, cSlice CPlane__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Plane, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.shape_msgs__msg__Plane__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_shape_msgs__msg__Plane * uintptr(i)),
		))
		(*goSlice)[i] = Plane{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func Plane__Sequence_to_C(cSlice *CPlane__Sequence, goSlice []Plane) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.shape_msgs__msg__Plane)(C.malloc((C.size_t)(C.sizeof_struct_shape_msgs__msg__Plane * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.shape_msgs__msg__Plane)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_shape_msgs__msg__Plane * uintptr(i)),
		))
		*cIdx = *(*C.shape_msgs__msg__Plane)(v.AsCStruct())
	}
}
func Plane__Array_to_Go(goSlice []Plane, cSlice []CPlane) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func Plane__Array_to_C(cSlice []CPlane, goSlice []Plane) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.shape_msgs__msg__Plane)(goSlice[i].AsCStruct())
	}
}

