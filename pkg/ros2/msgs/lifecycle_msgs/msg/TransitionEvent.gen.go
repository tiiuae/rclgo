/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package lifecycle_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -llifecycle_msgs__rosidl_typesupport_c -llifecycle_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <lifecycle_msgs/msg/transition_event.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("lifecycle_msgs/TransitionEvent", &TransitionEvent{})
}

// Do not create instances of this type directly. Always use NewTransitionEvent
// function instead.
type TransitionEvent struct {
	Timestamp uint64 `yaml:"timestamp"`// The time point at which this event occurred.
	Transition Transition `yaml:"transition"`// The id and label of this transition event.
	StartState State `yaml:"start_state"`// The starting state from which this event transitioned.
	GoalState State `yaml:"goal_state"`// The end state of this transition event.
}

// NewTransitionEvent creates a new TransitionEvent with default values.
func NewTransitionEvent() *TransitionEvent {
	self := TransitionEvent{}
	self.SetDefaults(nil)
	return &self
}

func (t *TransitionEvent) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.Transition.SetDefaults(nil)
	t.StartState.SetDefaults(nil)
	t.GoalState.SetDefaults(nil)
	
	return t
}

func (t *TransitionEvent) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__lifecycle_msgs__msg__TransitionEvent())
}
func (t *TransitionEvent) PrepareMemory() unsafe.Pointer { //returns *C.lifecycle_msgs__msg__TransitionEvent
	return (unsafe.Pointer)(C.lifecycle_msgs__msg__TransitionEvent__create())
}
func (t *TransitionEvent) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.lifecycle_msgs__msg__TransitionEvent__destroy((*C.lifecycle_msgs__msg__TransitionEvent)(pointer_to_free))
}
func (t *TransitionEvent) AsCStruct() unsafe.Pointer {
	mem := (*C.lifecycle_msgs__msg__TransitionEvent)(t.PrepareMemory())
	mem.timestamp = C.uint64_t(t.Timestamp)
	mem.transition = *(*C.lifecycle_msgs__msg__Transition)(t.Transition.AsCStruct())
	mem.start_state = *(*C.lifecycle_msgs__msg__State)(t.StartState.AsCStruct())
	mem.goal_state = *(*C.lifecycle_msgs__msg__State)(t.GoalState.AsCStruct())
	return unsafe.Pointer(mem)
}
func (t *TransitionEvent) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.lifecycle_msgs__msg__TransitionEvent)(ros2_message_buffer)
	t.Timestamp = uint64(mem.timestamp)
	t.Transition.AsGoStruct(unsafe.Pointer(&mem.transition))
	t.StartState.AsGoStruct(unsafe.Pointer(&mem.start_state))
	t.GoalState.AsGoStruct(unsafe.Pointer(&mem.goal_state))
}
func (t *TransitionEvent) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CTransitionEvent = C.lifecycle_msgs__msg__TransitionEvent
type CTransitionEvent__Sequence = C.lifecycle_msgs__msg__TransitionEvent__Sequence

func TransitionEvent__Sequence_to_Go(goSlice *[]TransitionEvent, cSlice CTransitionEvent__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]TransitionEvent, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.lifecycle_msgs__msg__TransitionEvent__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_lifecycle_msgs__msg__TransitionEvent * uintptr(i)),
		))
		(*goSlice)[i] = TransitionEvent{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func TransitionEvent__Sequence_to_C(cSlice *CTransitionEvent__Sequence, goSlice []TransitionEvent) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.lifecycle_msgs__msg__TransitionEvent)(C.malloc((C.size_t)(C.sizeof_struct_lifecycle_msgs__msg__TransitionEvent * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.lifecycle_msgs__msg__TransitionEvent)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_lifecycle_msgs__msg__TransitionEvent * uintptr(i)),
		))
		*cIdx = *(*C.lifecycle_msgs__msg__TransitionEvent)(v.AsCStruct())
	}
}
func TransitionEvent__Array_to_Go(goSlice []TransitionEvent, cSlice []CTransitionEvent) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func TransitionEvent__Array_to_C(cSlice []CTransitionEvent, goSlice []TransitionEvent) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.lifecycle_msgs__msg__TransitionEvent)(goSlice[i].AsCStruct())
	}
}

