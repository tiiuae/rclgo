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
#include <shape_msgs/msg/solid_primitive.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("shape_msgs/SolidPrimitive", &SolidPrimitive{})
}
const (
	SolidPrimitive_BOX uint8 = 1
	SolidPrimitive_SPHERE uint8 = 2
	SolidPrimitive_CYLINDER uint8 = 3
	SolidPrimitive_CONE uint8 = 4
	SolidPrimitive_BOX_X uint8 = 0// For type BOX, the X, Y, and Z dimensions are the length of the corresponding sides of the box.
	SolidPrimitive_BOX_Y uint8 = 1// For type BOX, the X, Y, and Z dimensions are the length of the corresponding sides of the box.
	SolidPrimitive_BOX_Z uint8 = 2// For type BOX, the X, Y, and Z dimensions are the length of the corresponding sides of the box.
	SolidPrimitive_SPHERE_RADIUS uint8 = 0// For the SPHERE type, only one component is used, and it gives the radius of the sphere.
	SolidPrimitive_CYLINDER_HEIGHT uint8 = 0
	SolidPrimitive_CYLINDER_RADIUS uint8 = 1
	SolidPrimitive_CONE_HEIGHT uint8 = 0
	SolidPrimitive_CONE_RADIUS uint8 = 1
)

// Do not create instances of this type directly. Always use NewSolidPrimitive
// function instead.
type SolidPrimitive struct {
	Type uint8 `yaml:"type"`// The type of the shape
	Dimensions []float64 `yaml:"dimensions"`// At no point will dimensions have a length > 3.. The dimensions of the shape
}

// NewSolidPrimitive creates a new SolidPrimitive with default values.
func NewSolidPrimitive() *SolidPrimitive {
	self := SolidPrimitive{}
	self.SetDefaults(nil)
	return &self
}

func (t *SolidPrimitive) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *SolidPrimitive) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__shape_msgs__msg__SolidPrimitive())
}
func (t *SolidPrimitive) PrepareMemory() unsafe.Pointer { //returns *C.shape_msgs__msg__SolidPrimitive
	return (unsafe.Pointer)(C.shape_msgs__msg__SolidPrimitive__create())
}
func (t *SolidPrimitive) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.shape_msgs__msg__SolidPrimitive__destroy((*C.shape_msgs__msg__SolidPrimitive)(pointer_to_free))
}
func (t *SolidPrimitive) AsCStruct() unsafe.Pointer {
	mem := (*C.shape_msgs__msg__SolidPrimitive)(t.PrepareMemory())
	mem._type = C.uint8_t(t.Type)
	rosidl_runtime_c.Float64__Sequence_to_C((*rosidl_runtime_c.CFloat64__Sequence)(unsafe.Pointer(&mem.dimensions)), t.Dimensions)
	return unsafe.Pointer(mem)
}
func (t *SolidPrimitive) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.shape_msgs__msg__SolidPrimitive)(ros2_message_buffer)
	t.Type = uint8(mem._type)
	rosidl_runtime_c.Float64__Sequence_to_Go(&t.Dimensions, *(*rosidl_runtime_c.CFloat64__Sequence)(unsafe.Pointer(&mem.dimensions)))
}
func (t *SolidPrimitive) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CSolidPrimitive = C.shape_msgs__msg__SolidPrimitive
type CSolidPrimitive__Sequence = C.shape_msgs__msg__SolidPrimitive__Sequence

func SolidPrimitive__Sequence_to_Go(goSlice *[]SolidPrimitive, cSlice CSolidPrimitive__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]SolidPrimitive, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.shape_msgs__msg__SolidPrimitive__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_shape_msgs__msg__SolidPrimitive * uintptr(i)),
		))
		(*goSlice)[i] = SolidPrimitive{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func SolidPrimitive__Sequence_to_C(cSlice *CSolidPrimitive__Sequence, goSlice []SolidPrimitive) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.shape_msgs__msg__SolidPrimitive)(C.malloc((C.size_t)(C.sizeof_struct_shape_msgs__msg__SolidPrimitive * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.shape_msgs__msg__SolidPrimitive)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_shape_msgs__msg__SolidPrimitive * uintptr(i)),
		))
		*cIdx = *(*C.shape_msgs__msg__SolidPrimitive)(v.AsCStruct())
	}
}
func SolidPrimitive__Array_to_Go(goSlice []SolidPrimitive, cSlice []CSolidPrimitive) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func SolidPrimitive__Array_to_C(cSlice []CSolidPrimitive, goSlice []SolidPrimitive) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.shape_msgs__msg__SolidPrimitive)(goSlice[i].AsCStruct())
	}
}

