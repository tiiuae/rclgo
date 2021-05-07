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
#include <px4_msgs/msg/position_setpoint.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("px4_msgs/PositionSetpoint", &PositionSetpoint{})
}
const (
	PositionSetpoint_SETPOINT_TYPE_POSITION uint8 = 0// position setpoint
	PositionSetpoint_SETPOINT_TYPE_VELOCITY uint8 = 1// velocity setpoint
	PositionSetpoint_SETPOINT_TYPE_LOITER uint8 = 2// loiter setpoint
	PositionSetpoint_SETPOINT_TYPE_TAKEOFF uint8 = 3// takeoff setpoint
	PositionSetpoint_SETPOINT_TYPE_LAND uint8 = 4// land setpoint, altitude must be ignored, descend until landing
	PositionSetpoint_SETPOINT_TYPE_IDLE uint8 = 5// do nothing, switch off motors or keep at idle speed (MC)
	PositionSetpoint_SETPOINT_TYPE_FOLLOW_TARGET uint8 = 6// setpoint in NED frame (x, y, z, vx, vy, vz) set by follow target
	PositionSetpoint_VELOCITY_FRAME_LOCAL_NED uint8 = 1// MAV_FRAME_LOCAL_NED
	PositionSetpoint_VELOCITY_FRAME_BODY_NED uint8 = 8// MAV_FRAME_BODY_NED
)

// Do not create instances of this type directly. Always use NewPositionSetpoint
// function instead.
type PositionSetpoint struct {
	Timestamp uint64 `yaml:"timestamp"`// time since system start (microseconds)
	Valid bool `yaml:"valid"`// true if setpoint is valid
	Type uint8 `yaml:"type"`// setpoint type to adjust behavior of position controller
	Vx float32 `yaml:"vx"`// local velocity setpoint in m/s in NED
	Vy float32 `yaml:"vy"`// local velocity setpoint in m/s in NED
	Vz float32 `yaml:"vz"`// local velocity setpoint in m/s in NED
	VelocityValid bool `yaml:"velocity_valid"`// true if local velocity setpoint valid
	VelocityFrame uint8 `yaml:"velocity_frame"`// to set velocity setpoints in NED or body
	AltValid bool `yaml:"alt_valid"`// do not set for 3D position control. Set to true if you want z-position control while doing vx,vy velocity control.
	Lat float64 `yaml:"lat"`// latitude, in deg
	Lon float64 `yaml:"lon"`// longitude, in deg
	Alt float32 `yaml:"alt"`// altitude AMSL, in m
	Yaw float32 `yaml:"yaw"`// yaw (only for multirotors), in rad [-PI..PI), NaN = hold current yaw
	YawValid bool `yaml:"yaw_valid"`// true if yaw setpoint valid
	Yawspeed float32 `yaml:"yawspeed"`// yawspeed (only for multirotors, in rad/s)
	YawspeedValid bool `yaml:"yawspeed_valid"`// true if yawspeed setpoint valid
	LandingGear int8 `yaml:"landing_gear"`// landing gear: see definition of the states in landing_gear.msg
	LoiterRadius float32 `yaml:"loiter_radius"`// loiter radius (only for fixed wing), in m
	LoiterDirection int8 `yaml:"loiter_direction"`// loiter direction: 1 = CW, -1 = CCW
	AcceptanceRadius float32 `yaml:"acceptance_radius"`// navigation acceptance_radius if we're doing waypoint navigation
	CruisingSpeed float32 `yaml:"cruising_speed"`// the generally desired cruising speed (not a hard constraint)
	CruisingThrottle float32 `yaml:"cruising_throttle"`// the generally desired cruising throttle (not a hard constraint)
	DisableWeatherVane bool `yaml:"disable_weather_vane"`// VTOL: disable (in auto mode) the weather vane feature that turns the nose into the wind
}

// NewPositionSetpoint creates a new PositionSetpoint with default values.
func NewPositionSetpoint() *PositionSetpoint {
	self := PositionSetpoint{}
	self.SetDefaults(nil)
	return &self
}

func (t *PositionSetpoint) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *PositionSetpoint) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__px4_msgs__msg__PositionSetpoint())
}
func (t *PositionSetpoint) PrepareMemory() unsafe.Pointer { //returns *C.px4_msgs__msg__PositionSetpoint
	return (unsafe.Pointer)(C.px4_msgs__msg__PositionSetpoint__create())
}
func (t *PositionSetpoint) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.px4_msgs__msg__PositionSetpoint__destroy((*C.px4_msgs__msg__PositionSetpoint)(pointer_to_free))
}
func (t *PositionSetpoint) AsCStruct() unsafe.Pointer {
	mem := (*C.px4_msgs__msg__PositionSetpoint)(t.PrepareMemory())
	mem.timestamp = C.uint64_t(t.Timestamp)
	mem.valid = C.bool(t.Valid)
	mem._type = C.uint8_t(t.Type)
	mem.vx = C.float(t.Vx)
	mem.vy = C.float(t.Vy)
	mem.vz = C.float(t.Vz)
	mem.velocity_valid = C.bool(t.VelocityValid)
	mem.velocity_frame = C.uint8_t(t.VelocityFrame)
	mem.alt_valid = C.bool(t.AltValid)
	mem.lat = C.double(t.Lat)
	mem.lon = C.double(t.Lon)
	mem.alt = C.float(t.Alt)
	mem.yaw = C.float(t.Yaw)
	mem.yaw_valid = C.bool(t.YawValid)
	mem.yawspeed = C.float(t.Yawspeed)
	mem.yawspeed_valid = C.bool(t.YawspeedValid)
	mem.landing_gear = C.int8_t(t.LandingGear)
	mem.loiter_radius = C.float(t.LoiterRadius)
	mem.loiter_direction = C.int8_t(t.LoiterDirection)
	mem.acceptance_radius = C.float(t.AcceptanceRadius)
	mem.cruising_speed = C.float(t.CruisingSpeed)
	mem.cruising_throttle = C.float(t.CruisingThrottle)
	mem.disable_weather_vane = C.bool(t.DisableWeatherVane)
	return unsafe.Pointer(mem)
}
func (t *PositionSetpoint) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.px4_msgs__msg__PositionSetpoint)(ros2_message_buffer)
	t.Timestamp = uint64(mem.timestamp)
	t.Valid = bool(mem.valid)
	t.Type = uint8(mem._type)
	t.Vx = float32(mem.vx)
	t.Vy = float32(mem.vy)
	t.Vz = float32(mem.vz)
	t.VelocityValid = bool(mem.velocity_valid)
	t.VelocityFrame = uint8(mem.velocity_frame)
	t.AltValid = bool(mem.alt_valid)
	t.Lat = float64(mem.lat)
	t.Lon = float64(mem.lon)
	t.Alt = float32(mem.alt)
	t.Yaw = float32(mem.yaw)
	t.YawValid = bool(mem.yaw_valid)
	t.Yawspeed = float32(mem.yawspeed)
	t.YawspeedValid = bool(mem.yawspeed_valid)
	t.LandingGear = int8(mem.landing_gear)
	t.LoiterRadius = float32(mem.loiter_radius)
	t.LoiterDirection = int8(mem.loiter_direction)
	t.AcceptanceRadius = float32(mem.acceptance_radius)
	t.CruisingSpeed = float32(mem.cruising_speed)
	t.CruisingThrottle = float32(mem.cruising_throttle)
	t.DisableWeatherVane = bool(mem.disable_weather_vane)
}
func (t *PositionSetpoint) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CPositionSetpoint = C.px4_msgs__msg__PositionSetpoint
type CPositionSetpoint__Sequence = C.px4_msgs__msg__PositionSetpoint__Sequence

func PositionSetpoint__Sequence_to_Go(goSlice *[]PositionSetpoint, cSlice CPositionSetpoint__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]PositionSetpoint, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.px4_msgs__msg__PositionSetpoint__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__PositionSetpoint * uintptr(i)),
		))
		(*goSlice)[i] = PositionSetpoint{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func PositionSetpoint__Sequence_to_C(cSlice *CPositionSetpoint__Sequence, goSlice []PositionSetpoint) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.px4_msgs__msg__PositionSetpoint)(C.malloc((C.size_t)(C.sizeof_struct_px4_msgs__msg__PositionSetpoint * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.px4_msgs__msg__PositionSetpoint)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__PositionSetpoint * uintptr(i)),
		))
		*cIdx = *(*C.px4_msgs__msg__PositionSetpoint)(v.AsCStruct())
	}
}
func PositionSetpoint__Array_to_Go(goSlice []PositionSetpoint, cSlice []CPositionSetpoint) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func PositionSetpoint__Array_to_C(cSlice []CPositionSetpoint, goSlice []PositionSetpoint) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.px4_msgs__msg__PositionSetpoint)(goSlice[i].AsCStruct())
	}
}

