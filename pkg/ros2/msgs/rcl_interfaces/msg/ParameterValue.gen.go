/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package rcl_interfaces
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	rosidl_runtime_c "github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lrcl_interfaces__rosidl_typesupport_c -lrcl_interfaces__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rcl_interfaces/msg/parameter_value.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("rcl_interfaces/ParameterValue", &ParameterValue{})
}

// Do not create instances of this type directly. Always use NewParameterValue
// function instead.
type ParameterValue struct {
	Type uint8 `yaml:"type"`// The type of this parameter, which corresponds to the appropriate field below.
	BoolValue bool `yaml:"bool_value"`// Boolean value, can be either true or false.
	IntegerValue int64 `yaml:"integer_value"`// Integer value ranging from -9,223,372,036,854,775,808 to9,223,372,036,854,775,807.
	DoubleValue float64 `yaml:"double_value"`// A double precision floating point value following IEEE 754.
	StringValue rosidl_runtime_c.String `yaml:"string_value"`// A textual value with no practical length limit.
	ByteArrayValue []byte `yaml:"byte_array_value"`// An array of bytes, used for non-textual information.
	BoolArrayValue []bool `yaml:"bool_array_value"`// An array of boolean values.
	IntegerArrayValue []int64 `yaml:"integer_array_value"`// An array of 64-bit integer values.
	DoubleArrayValue []float64 `yaml:"double_array_value"`// An array of 64-bit floating point values.
	StringArrayValue []rosidl_runtime_c.String `yaml:"string_array_value"`// An array of string values.
}

// NewParameterValue creates a new ParameterValue with default values.
func NewParameterValue() *ParameterValue {
	self := ParameterValue{}
	self.SetDefaults(nil)
	return &self
}

func (t *ParameterValue) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.StringValue.SetDefaults("")
	
	return t
}

func (t *ParameterValue) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__rcl_interfaces__msg__ParameterValue())
}
func (t *ParameterValue) PrepareMemory() unsafe.Pointer { //returns *C.rcl_interfaces__msg__ParameterValue
	return (unsafe.Pointer)(C.rcl_interfaces__msg__ParameterValue__create())
}
func (t *ParameterValue) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.rcl_interfaces__msg__ParameterValue__destroy((*C.rcl_interfaces__msg__ParameterValue)(pointer_to_free))
}
func (t *ParameterValue) AsCStruct() unsafe.Pointer {
	mem := (*C.rcl_interfaces__msg__ParameterValue)(t.PrepareMemory())
	mem._type = C.uint8_t(t.Type)
	mem.bool_value = C.bool(t.BoolValue)
	mem.integer_value = C.int64_t(t.IntegerValue)
	mem.double_value = C.double(t.DoubleValue)
	mem.string_value = *(*C.rosidl_runtime_c__String)(t.StringValue.AsCStruct())
	rosidl_runtime_c.Byte__Sequence_to_C((*rosidl_runtime_c.CByte__Sequence)(unsafe.Pointer(&mem.byte_array_value)), t.ByteArrayValue)
	rosidl_runtime_c.Bool__Sequence_to_C((*rosidl_runtime_c.CBool__Sequence)(unsafe.Pointer(&mem.bool_array_value)), t.BoolArrayValue)
	rosidl_runtime_c.Int64__Sequence_to_C((*rosidl_runtime_c.CInt64__Sequence)(unsafe.Pointer(&mem.integer_array_value)), t.IntegerArrayValue)
	rosidl_runtime_c.Float64__Sequence_to_C((*rosidl_runtime_c.CFloat64__Sequence)(unsafe.Pointer(&mem.double_array_value)), t.DoubleArrayValue)
	rosidl_runtime_c.String__Sequence_to_C((*rosidl_runtime_c.CString__Sequence)(unsafe.Pointer(&mem.string_array_value)), t.StringArrayValue)
	return unsafe.Pointer(mem)
}
func (t *ParameterValue) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.rcl_interfaces__msg__ParameterValue)(ros2_message_buffer)
	t.Type = uint8(mem._type)
	t.BoolValue = bool(mem.bool_value)
	t.IntegerValue = int64(mem.integer_value)
	t.DoubleValue = float64(mem.double_value)
	t.StringValue.AsGoStruct(unsafe.Pointer(&mem.string_value))
	rosidl_runtime_c.Byte__Sequence_to_Go(&t.ByteArrayValue, *(*rosidl_runtime_c.CByte__Sequence)(unsafe.Pointer(&mem.byte_array_value)))
	rosidl_runtime_c.Bool__Sequence_to_Go(&t.BoolArrayValue, *(*rosidl_runtime_c.CBool__Sequence)(unsafe.Pointer(&mem.bool_array_value)))
	rosidl_runtime_c.Int64__Sequence_to_Go(&t.IntegerArrayValue, *(*rosidl_runtime_c.CInt64__Sequence)(unsafe.Pointer(&mem.integer_array_value)))
	rosidl_runtime_c.Float64__Sequence_to_Go(&t.DoubleArrayValue, *(*rosidl_runtime_c.CFloat64__Sequence)(unsafe.Pointer(&mem.double_array_value)))
	rosidl_runtime_c.String__Sequence_to_Go(&t.StringArrayValue, *(*rosidl_runtime_c.CString__Sequence)(unsafe.Pointer(&mem.string_array_value)))
}
func (t *ParameterValue) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CParameterValue = C.rcl_interfaces__msg__ParameterValue
type CParameterValue__Sequence = C.rcl_interfaces__msg__ParameterValue__Sequence

func ParameterValue__Sequence_to_Go(goSlice *[]ParameterValue, cSlice CParameterValue__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]ParameterValue, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.rcl_interfaces__msg__ParameterValue__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rcl_interfaces__msg__ParameterValue * uintptr(i)),
		))
		(*goSlice)[i] = ParameterValue{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func ParameterValue__Sequence_to_C(cSlice *CParameterValue__Sequence, goSlice []ParameterValue) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.rcl_interfaces__msg__ParameterValue)(C.malloc((C.size_t)(C.sizeof_struct_rcl_interfaces__msg__ParameterValue * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.rcl_interfaces__msg__ParameterValue)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rcl_interfaces__msg__ParameterValue * uintptr(i)),
		))
		*cIdx = *(*C.rcl_interfaces__msg__ParameterValue)(v.AsCStruct())
	}
}
func ParameterValue__Array_to_Go(goSlice []ParameterValue, cSlice []CParameterValue) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func ParameterValue__Array_to_C(cSlice []CParameterValue, goSlice []ParameterValue) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.rcl_interfaces__msg__ParameterValue)(goSlice[i].AsCStruct())
	}
}

