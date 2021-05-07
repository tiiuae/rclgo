/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package px4_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lpx4_msgs__rosidl_typesupport_c -lpx4_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <px4_msgs/msg/differential_pressure.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("px4_msgs/DifferentialPressure", &DifferentialPressure{})
}

// Do not create instances of this type directly. Always use NewDifferentialPressure
// function instead.
type DifferentialPressure struct {
	Timestamp uint64 `yaml:"timestamp"`// time since system start (microseconds)
	ErrorCount uint64 `yaml:"error_count"`// Number of errors detected by driver
	DifferentialPressureRawPa float32 `yaml:"differential_pressure_raw_pa"`// Raw differential pressure reading (may be negative)
	DifferentialPressureFilteredPa float32 `yaml:"differential_pressure_filtered_pa"`// Low pass filtered differential pressure reading
	Temperature float32 `yaml:"temperature"`// Temperature provided by sensor, -1000.0f if unknown
	DeviceId uint32 `yaml:"device_id"`// unique device ID for the sensor that does not change between power cycles
}

// NewDifferentialPressure creates a new DifferentialPressure with default values.
func NewDifferentialPressure() *DifferentialPressure {
	self := DifferentialPressure{}
	self.SetDefaults(nil)
	return &self
}

func (t *DifferentialPressure) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *DifferentialPressure) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__px4_msgs__msg__DifferentialPressure())
}
func (t *DifferentialPressure) PrepareMemory() unsafe.Pointer { //returns *C.px4_msgs__msg__DifferentialPressure
	return (unsafe.Pointer)(C.px4_msgs__msg__DifferentialPressure__create())
}
func (t *DifferentialPressure) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.px4_msgs__msg__DifferentialPressure__destroy((*C.px4_msgs__msg__DifferentialPressure)(pointer_to_free))
}
func (t *DifferentialPressure) AsCStruct() unsafe.Pointer {
	mem := (*C.px4_msgs__msg__DifferentialPressure)(t.PrepareMemory())
	mem.timestamp = C.uint64_t(t.Timestamp)
	mem.error_count = C.uint64_t(t.ErrorCount)
	mem.differential_pressure_raw_pa = C.float(t.DifferentialPressureRawPa)
	mem.differential_pressure_filtered_pa = C.float(t.DifferentialPressureFilteredPa)
	mem.temperature = C.float(t.Temperature)
	mem.device_id = C.uint32_t(t.DeviceId)
	return unsafe.Pointer(mem)
}
func (t *DifferentialPressure) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.px4_msgs__msg__DifferentialPressure)(ros2_message_buffer)
	t.Timestamp = uint64(mem.timestamp)
	t.ErrorCount = uint64(mem.error_count)
	t.DifferentialPressureRawPa = float32(mem.differential_pressure_raw_pa)
	t.DifferentialPressureFilteredPa = float32(mem.differential_pressure_filtered_pa)
	t.Temperature = float32(mem.temperature)
	t.DeviceId = uint32(mem.device_id)
}
func (t *DifferentialPressure) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CDifferentialPressure = C.px4_msgs__msg__DifferentialPressure
type CDifferentialPressure__Sequence = C.px4_msgs__msg__DifferentialPressure__Sequence

func DifferentialPressure__Sequence_to_Go(goSlice *[]DifferentialPressure, cSlice CDifferentialPressure__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]DifferentialPressure, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.px4_msgs__msg__DifferentialPressure__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__DifferentialPressure * uintptr(i)),
		))
		(*goSlice)[i] = DifferentialPressure{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func DifferentialPressure__Sequence_to_C(cSlice *CDifferentialPressure__Sequence, goSlice []DifferentialPressure) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.px4_msgs__msg__DifferentialPressure)(C.malloc((C.size_t)(C.sizeof_struct_px4_msgs__msg__DifferentialPressure * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.px4_msgs__msg__DifferentialPressure)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__DifferentialPressure * uintptr(i)),
		))
		*cIdx = *(*C.px4_msgs__msg__DifferentialPressure)(v.AsCStruct())
	}
}
func DifferentialPressure__Array_to_Go(goSlice []DifferentialPressure, cSlice []CDifferentialPressure) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func DifferentialPressure__Array_to_C(cSlice []CDifferentialPressure, goSlice []DifferentialPressure) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.px4_msgs__msg__DifferentialPressure)(goSlice[i].AsCStruct())
	}
}

