/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package action_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -laction_msgs__rosidl_typesupport_c -laction_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <action_msgs/msg/goal_status.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("action_msgs/GoalStatus", &GoalStatus{})
}
const (
	GoalStatus_STATUS_UNKNOWN int8 = 0// Indicates status has not been properly set.
	GoalStatus_STATUS_ACCEPTED int8 = 1// The goal has been accepted and is awaiting execution.
	GoalStatus_STATUS_EXECUTING int8 = 2// The goal is currently being executed by the action server.
	GoalStatus_STATUS_CANCELING int8 = 3// The client has requested that the goal be canceled and the action server hasaccepted the cancel request.
	GoalStatus_STATUS_SUCCEEDED int8 = 4// The goal was achieved successfully by the action server.
	GoalStatus_STATUS_CANCELED int8 = 5// The goal was canceled after an external request from an action client.
	GoalStatus_STATUS_ABORTED int8 = 6// The goal was terminated by the action server without an external request.
)

// Do not create instances of this type directly. Always use NewGoalStatus
// function instead.
type GoalStatus struct {
	GoalInfo GoalInfo `yaml:"goal_info"`// Goal info (contains ID and timestamp).
	Status int8 `yaml:"status"`// Action goal state-machine status.
}

// NewGoalStatus creates a new GoalStatus with default values.
func NewGoalStatus() *GoalStatus {
	self := GoalStatus{}
	self.SetDefaults(nil)
	return &self
}

func (t *GoalStatus) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.GoalInfo.SetDefaults(nil)
	
	return t
}

func (t *GoalStatus) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__action_msgs__msg__GoalStatus())
}
func (t *GoalStatus) PrepareMemory() unsafe.Pointer { //returns *C.action_msgs__msg__GoalStatus
	return (unsafe.Pointer)(C.action_msgs__msg__GoalStatus__create())
}
func (t *GoalStatus) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.action_msgs__msg__GoalStatus__destroy((*C.action_msgs__msg__GoalStatus)(pointer_to_free))
}
func (t *GoalStatus) AsCStruct() unsafe.Pointer {
	mem := (*C.action_msgs__msg__GoalStatus)(t.PrepareMemory())
	mem.goal_info = *(*C.action_msgs__msg__GoalInfo)(t.GoalInfo.AsCStruct())
	mem.status = C.int8_t(t.Status)
	return unsafe.Pointer(mem)
}
func (t *GoalStatus) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.action_msgs__msg__GoalStatus)(ros2_message_buffer)
	t.GoalInfo.AsGoStruct(unsafe.Pointer(&mem.goal_info))
	t.Status = int8(mem.status)
}
func (t *GoalStatus) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CGoalStatus = C.action_msgs__msg__GoalStatus
type CGoalStatus__Sequence = C.action_msgs__msg__GoalStatus__Sequence

func GoalStatus__Sequence_to_Go(goSlice *[]GoalStatus, cSlice CGoalStatus__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]GoalStatus, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.action_msgs__msg__GoalStatus__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_action_msgs__msg__GoalStatus * uintptr(i)),
		))
		(*goSlice)[i] = GoalStatus{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func GoalStatus__Sequence_to_C(cSlice *CGoalStatus__Sequence, goSlice []GoalStatus) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.action_msgs__msg__GoalStatus)(C.malloc((C.size_t)(C.sizeof_struct_action_msgs__msg__GoalStatus * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.action_msgs__msg__GoalStatus)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_action_msgs__msg__GoalStatus * uintptr(i)),
		))
		*cIdx = *(*C.action_msgs__msg__GoalStatus)(v.AsCStruct())
	}
}
func GoalStatus__Array_to_Go(goSlice []GoalStatus, cSlice []CGoalStatus) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func GoalStatus__Array_to_C(cSlice []CGoalStatus, goSlice []GoalStatus) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.action_msgs__msg__GoalStatus)(goSlice[i].AsCStruct())
	}
}

