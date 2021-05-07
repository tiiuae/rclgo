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
#include <px4_msgs/msg/esc_report.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("px4_msgs/EscReport", &EscReport{})
}
const (
	EscReport_FAILURE_NONE uint8 = 0
	EscReport_FAILURE_OVER_CURRENT_MASK uint8 = 1// (1 << 0)
	EscReport_FAILURE_OVER_VOLTAGE_MASK uint8 = 2// (1 << 1)
	EscReport_FAILURE_OVER_TEMPERATURE_MASK uint8 = 4// (1 << 2)
	EscReport_FAILURE_OVER_RPM_MASK uint8 = 8// (1 << 3)
	EscReport_FAILURE_INCONSISTENT_CMD_MASK uint8 = 16// (1 << 4)  Set if ESC received an inconsistent command (i.e out of boundaries)
	EscReport_FAILURE_MOTOR_STUCK_MASK uint8 = 32// (1 << 5)
	EscReport_FAILURE_GENERIC_MASK uint8 = 64// (1 << 6)
)

// Do not create instances of this type directly. Always use NewEscReport
// function instead.
type EscReport struct {
	Timestamp uint64 `yaml:"timestamp"`// time since system start (microseconds)
	EscErrorcount uint32 `yaml:"esc_errorcount"`// Number of reported errors by ESC - if supported
	EscRpm int32 `yaml:"esc_rpm"`// Motor RPM, negative for reverse rotation [RPM] - if supported
	EscVoltage float32 `yaml:"esc_voltage"`// Voltage measured from current ESC [V] - if supported
	EscCurrent float32 `yaml:"esc_current"`// Current measured from current ESC [A] - if supported
	EscTemperature uint8 `yaml:"esc_temperature"`// Temperature measured from current ESC [degC] - if supported
	EscAddress uint8 `yaml:"esc_address"`// Address of current ESC (in most cases 1-8 / must be set by driver)
	EscState uint8 `yaml:"esc_state"`// State of ESC - depend on Vendor
	Failures uint8 `yaml:"failures"`// Bitmask to indicate the internal ESC faults
}

// NewEscReport creates a new EscReport with default values.
func NewEscReport() *EscReport {
	self := EscReport{}
	self.SetDefaults(nil)
	return &self
}

func (t *EscReport) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *EscReport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__px4_msgs__msg__EscReport())
}
func (t *EscReport) PrepareMemory() unsafe.Pointer { //returns *C.px4_msgs__msg__EscReport
	return (unsafe.Pointer)(C.px4_msgs__msg__EscReport__create())
}
func (t *EscReport) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.px4_msgs__msg__EscReport__destroy((*C.px4_msgs__msg__EscReport)(pointer_to_free))
}
func (t *EscReport) AsCStruct() unsafe.Pointer {
	mem := (*C.px4_msgs__msg__EscReport)(t.PrepareMemory())
	mem.timestamp = C.uint64_t(t.Timestamp)
	mem.esc_errorcount = C.uint32_t(t.EscErrorcount)
	mem.esc_rpm = C.int32_t(t.EscRpm)
	mem.esc_voltage = C.float(t.EscVoltage)
	mem.esc_current = C.float(t.EscCurrent)
	mem.esc_temperature = C.uint8_t(t.EscTemperature)
	mem.esc_address = C.uint8_t(t.EscAddress)
	mem.esc_state = C.uint8_t(t.EscState)
	mem.failures = C.uint8_t(t.Failures)
	return unsafe.Pointer(mem)
}
func (t *EscReport) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.px4_msgs__msg__EscReport)(ros2_message_buffer)
	t.Timestamp = uint64(mem.timestamp)
	t.EscErrorcount = uint32(mem.esc_errorcount)
	t.EscRpm = int32(mem.esc_rpm)
	t.EscVoltage = float32(mem.esc_voltage)
	t.EscCurrent = float32(mem.esc_current)
	t.EscTemperature = uint8(mem.esc_temperature)
	t.EscAddress = uint8(mem.esc_address)
	t.EscState = uint8(mem.esc_state)
	t.Failures = uint8(mem.failures)
}
func (t *EscReport) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CEscReport = C.px4_msgs__msg__EscReport
type CEscReport__Sequence = C.px4_msgs__msg__EscReport__Sequence

func EscReport__Sequence_to_Go(goSlice *[]EscReport, cSlice CEscReport__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]EscReport, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.px4_msgs__msg__EscReport__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__EscReport * uintptr(i)),
		))
		(*goSlice)[i] = EscReport{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func EscReport__Sequence_to_C(cSlice *CEscReport__Sequence, goSlice []EscReport) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.px4_msgs__msg__EscReport)(C.malloc((C.size_t)(C.sizeof_struct_px4_msgs__msg__EscReport * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.px4_msgs__msg__EscReport)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__EscReport * uintptr(i)),
		))
		*cIdx = *(*C.px4_msgs__msg__EscReport)(v.AsCStruct())
	}
}
func EscReport__Array_to_Go(goSlice []EscReport, cSlice []CEscReport) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func EscReport__Array_to_C(cSlice []CEscReport, goSlice []EscReport) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.px4_msgs__msg__EscReport)(goSlice[i].AsCStruct())
	}
}

