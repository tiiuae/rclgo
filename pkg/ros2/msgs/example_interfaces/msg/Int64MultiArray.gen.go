/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package example_interfaces
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	rosidl_runtime_c "github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lexample_interfaces__rosidl_typesupport_c -lexample_interfaces__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <example_interfaces/msg/int64_multi_array.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("example_interfaces/Int64MultiArray", &Int64MultiArray{})
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
	self.SetDefaults(nil)
	return &self
}

func (t *Int64MultiArray) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.Layout.SetDefaults(nil)
	
	return t
}

func (t *Int64MultiArray) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__example_interfaces__msg__Int64MultiArray())
}
func (t *Int64MultiArray) PrepareMemory() unsafe.Pointer { //returns *C.example_interfaces__msg__Int64MultiArray
	return (unsafe.Pointer)(C.example_interfaces__msg__Int64MultiArray__create())
}
func (t *Int64MultiArray) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.example_interfaces__msg__Int64MultiArray__destroy((*C.example_interfaces__msg__Int64MultiArray)(pointer_to_free))
}
func (t *Int64MultiArray) AsCStruct() unsafe.Pointer {
	mem := (*C.example_interfaces__msg__Int64MultiArray)(t.PrepareMemory())
	mem.layout = *(*C.example_interfaces__msg__MultiArrayLayout)(t.Layout.AsCStruct())
	rosidl_runtime_c.Int64__Sequence_to_C((*rosidl_runtime_c.CInt64__Sequence)(unsafe.Pointer(&mem.data)), t.Data)
	return unsafe.Pointer(mem)
}
func (t *Int64MultiArray) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.example_interfaces__msg__Int64MultiArray)(ros2_message_buffer)
	t.Layout.AsGoStruct(unsafe.Pointer(&mem.layout))
	rosidl_runtime_c.Int64__Sequence_to_Go(&t.Data, *(*rosidl_runtime_c.CInt64__Sequence)(unsafe.Pointer(&mem.data)))
}
func (t *Int64MultiArray) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CInt64MultiArray = C.example_interfaces__msg__Int64MultiArray
type CInt64MultiArray__Sequence = C.example_interfaces__msg__Int64MultiArray__Sequence

func Int64MultiArray__Sequence_to_Go(goSlice *[]Int64MultiArray, cSlice CInt64MultiArray__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]Int64MultiArray, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.example_interfaces__msg__Int64MultiArray__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_example_interfaces__msg__Int64MultiArray * uintptr(i)),
		))
		(*goSlice)[i] = Int64MultiArray{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func Int64MultiArray__Sequence_to_C(cSlice *CInt64MultiArray__Sequence, goSlice []Int64MultiArray) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.example_interfaces__msg__Int64MultiArray)(C.malloc((C.size_t)(C.sizeof_struct_example_interfaces__msg__Int64MultiArray * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.example_interfaces__msg__Int64MultiArray)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_example_interfaces__msg__Int64MultiArray * uintptr(i)),
		))
		*cIdx = *(*C.example_interfaces__msg__Int64MultiArray)(v.AsCStruct())
	}
}
func Int64MultiArray__Array_to_Go(goSlice []Int64MultiArray, cSlice []CInt64MultiArray) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func Int64MultiArray__Array_to_C(cSlice []CInt64MultiArray, goSlice []Int64MultiArray) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.example_interfaces__msg__Int64MultiArray)(goSlice[i].AsCStruct())
	}
}

