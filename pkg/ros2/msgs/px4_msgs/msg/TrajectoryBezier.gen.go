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
	rosidl_runtime_c "github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lpx4_msgs__rosidl_typesupport_c -lpx4_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <px4_msgs/msg/trajectory_bezier.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("px4_msgs/TrajectoryBezier", &TrajectoryBezier{})
}

// Do not create instances of this type directly. Always use NewTrajectoryBezier
// function instead.
type TrajectoryBezier struct {
	Timestamp uint64 `yaml:"timestamp"`// time since system start (microseconds)
	Position [3]float32 `yaml:"position"`// local position x,y,z (metres)
	Yaw float32 `yaml:"yaw"`// yaw angle (rad)
	Delta float32 `yaml:"delta"`// time it should take to get to this waypoint, if this is the final waypoint (seconds)
}

// NewTrajectoryBezier creates a new TrajectoryBezier with default values.
func NewTrajectoryBezier() *TrajectoryBezier {
	self := TrajectoryBezier{}
	self.SetDefaults(nil)
	return &self
}

func (t *TrajectoryBezier) SetDefaults(d interface{}) ros2types.ROS2Msg {
	
	return t
}

func (t *TrajectoryBezier) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__px4_msgs__msg__TrajectoryBezier())
}
func (t *TrajectoryBezier) PrepareMemory() unsafe.Pointer { //returns *C.px4_msgs__msg__TrajectoryBezier
	return (unsafe.Pointer)(C.px4_msgs__msg__TrajectoryBezier__create())
}
func (t *TrajectoryBezier) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.px4_msgs__msg__TrajectoryBezier__destroy((*C.px4_msgs__msg__TrajectoryBezier)(pointer_to_free))
}
func (t *TrajectoryBezier) AsCStruct() unsafe.Pointer {
	mem := (*C.px4_msgs__msg__TrajectoryBezier)(t.PrepareMemory())
	mem.timestamp = C.uint64_t(t.Timestamp)
	cSlice_position := mem.position[:]
	rosidl_runtime_c.Float32__Array_to_C(*(*[]rosidl_runtime_c.CFloat32)(unsafe.Pointer(&cSlice_position)), t.Position[:])
	mem.yaw = C.float(t.Yaw)
	mem.delta = C.float(t.Delta)
	return unsafe.Pointer(mem)
}
func (t *TrajectoryBezier) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.px4_msgs__msg__TrajectoryBezier)(ros2_message_buffer)
	t.Timestamp = uint64(mem.timestamp)
	cSlice_position := mem.position[:]
	rosidl_runtime_c.Float32__Array_to_Go(t.Position[:], *(*[]rosidl_runtime_c.CFloat32)(unsafe.Pointer(&cSlice_position)))
	t.Yaw = float32(mem.yaw)
	t.Delta = float32(mem.delta)
}
func (t *TrajectoryBezier) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CTrajectoryBezier = C.px4_msgs__msg__TrajectoryBezier
type CTrajectoryBezier__Sequence = C.px4_msgs__msg__TrajectoryBezier__Sequence

func TrajectoryBezier__Sequence_to_Go(goSlice *[]TrajectoryBezier, cSlice CTrajectoryBezier__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]TrajectoryBezier, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.px4_msgs__msg__TrajectoryBezier__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__TrajectoryBezier * uintptr(i)),
		))
		(*goSlice)[i] = TrajectoryBezier{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func TrajectoryBezier__Sequence_to_C(cSlice *CTrajectoryBezier__Sequence, goSlice []TrajectoryBezier) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.px4_msgs__msg__TrajectoryBezier)(C.malloc((C.size_t)(C.sizeof_struct_px4_msgs__msg__TrajectoryBezier * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.px4_msgs__msg__TrajectoryBezier)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_px4_msgs__msg__TrajectoryBezier * uintptr(i)),
		))
		*cIdx = *(*C.px4_msgs__msg__TrajectoryBezier)(v.AsCStruct())
	}
}
func TrajectoryBezier__Array_to_Go(goSlice []TrajectoryBezier, cSlice []CTrajectoryBezier) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func TrajectoryBezier__Array_to_C(cSlice []CTrajectoryBezier, goSlice []TrajectoryBezier) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.px4_msgs__msg__TrajectoryBezier)(goSlice[i].AsCStruct())
	}
}

