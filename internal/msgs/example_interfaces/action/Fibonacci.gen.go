/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by rclgo-gen. DO NOT EDIT.

package example_interfaces_action

/*
#cgo LDFLAGS: -L/opt/ros/galactic/lib -Wl,-rpath=/opt/ros/galactic/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lexample_interfaces__rosidl_typesupport_c -lexample_interfaces__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/galactic/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <example_interfaces/action/fibonacci.h>
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/rclgo/typemap"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"

	action_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/action_msgs/msg"
	action_msgs_srv "github.com/tiiuae/rclgo/internal/msgs/action_msgs/srv"
)

func init() {
	typemap.RegisterAction("example_interfaces/Fibonacci", FibonacciTypeSupport)
}

type _FibonacciTypeSupport struct {}

func (s _FibonacciTypeSupport) Goal() types.MessageTypeSupport {
	return Fibonacci_GoalTypeSupport
}

func (s _FibonacciTypeSupport) SendGoal() types.ServiceTypeSupport {
	return Fibonacci_SendGoalTypeSupport
}

func (s _FibonacciTypeSupport) NewSendGoalResponse(accepted bool, stamp time.Duration) types.Message {
	msg := NewFibonacci_SendGoal_Response()
	msg.Accepted = accepted
	secs := stamp.Truncate(time.Second)
	msg.Stamp.Sec = int32(secs)
	msg.Stamp.Nanosec = uint32(stamp - secs)
	return msg
}

func (s _FibonacciTypeSupport) Result() types.MessageTypeSupport {
	return Fibonacci_ResultTypeSupport
}

func (s _FibonacciTypeSupport) GetResult() types.ServiceTypeSupport {
	return Fibonacci_GetResultTypeSupport
}

func (s _FibonacciTypeSupport) NewGetResultResponse(status int8, result types.Message) types.Message {
	msg := NewFibonacci_GetResult_Response()
	msg.Status = status
	if result == nil {
		msg.Result = *NewFibonacci_Result()
	} else {
		msg.Result = *result.(*Fibonacci_Result)
	}
	return msg
}

func (s _FibonacciTypeSupport) CancelGoal() types.ServiceTypeSupport {
	return action_msgs_srv.CancelGoalTypeSupport
}

func (s _FibonacciTypeSupport) Feedback() types.MessageTypeSupport {
	return Fibonacci_FeedbackTypeSupport
}

func (s _FibonacciTypeSupport) FeedbackMessage() types.MessageTypeSupport {
	return Fibonacci_FeedbackMessageTypeSupport
}

func (s _FibonacciTypeSupport) NewFeedbackMessage(goalID *types.GoalID, feedback types.Message) types.Message {
	msg := NewFibonacci_FeedbackMessage()
	msg.GoalID.Uuid = *goalID
	msg.Feedback = *feedback.(*Fibonacci_Feedback)
	return msg
}

func (s _FibonacciTypeSupport) GoalStatusArray() types.MessageTypeSupport {
	return action_msgs_msg.GoalStatusArrayTypeSupport
}

func (s _FibonacciTypeSupport) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_action_type_support_handle__example_interfaces__action__Fibonacci())
}

// Modifying this variable is undefined behavior.
var FibonacciTypeSupport types.ActionTypeSupport = _FibonacciTypeSupport{}
