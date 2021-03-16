package ros2

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo LDFLAGS: /home/kivilahtio/install/rclc/lib/librclc.so
#cgo CFLAGS: -I/opt/ros/foxy/include -I/home/kivilahtio/install/rclc/include/

#include <rmw/ret_types.h>
#include <rcl/types.h>

*/
import "C"
import (
	"fmt"
)

type RCLError interface {
	Error() string
	rcl_ret() int
	context() string
}

type RCL_RET_GOLANG_UNKNOWN_RET_TYPE struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_GOLANG_UNKNOWN_RET_TYPE) Error() string {
	return errStr(fmt.Sprintf(
		"Unknown 'rcl_ret_t' type '%d'! RCL error definitions might have changed and"+
			"the error mapping with Golang bindings needs to be updated!.",
		e.rcl_ret_t), e.ctx)
}
func (e *RCL_RET_GOLANG_UNKNOWN_RET_TYPE) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_GOLANG_UNKNOWN_RET_TYPE) context() string {
	return e.ctx
}

type RCL_RET_ERROR struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_ERROR) Error() string {
	return errStr("RCL_RET_ERROR.", e.ctx)
}
func (e *RCL_RET_ERROR) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_ERROR) context() string {
	return e.ctx
}

type RCL_RET_ALREADY_INIT struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_ALREADY_INIT) Error() string {
	return errStr("RCL_RET_ALREADY_INIT.", e.ctx)
}
func (e *RCL_RET_ALREADY_INIT) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_ALREADY_INIT) context() string {
	return e.ctx
}

type RCL_RET_INVALID_ARGUMENT struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_INVALID_ARGUMENT) Error() string {
	return errStr("RCL_RET_INVALID_ARGUMENT.", e.ctx)
}
func (e *RCL_RET_INVALID_ARGUMENT) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_INVALID_ARGUMENT) context() string {
	return e.ctx
}

type RCL_RET_TOPIC_NAME_INVALID struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_TOPIC_NAME_INVALID) Error() string {
	return errStr("RCL_RET_TOPIC_NAME_INVALID.", e.ctx)
}
func (e *RCL_RET_TOPIC_NAME_INVALID) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_TOPIC_NAME_INVALID) context() string {
	return e.ctx
}

type RCL_RET_NODE_INVALID_NAME struct {
	rcl_ret_t int
	ctx       string
}

func (e *RCL_RET_NODE_INVALID_NAME) Error() string {
	return errStr("RCL_RET_NODE_INVALID_NAME.", e.ctx)
}
func (e *RCL_RET_NODE_INVALID_NAME) rcl_ret() int {
	return e.rcl_ret_t
}
func (e *RCL_RET_NODE_INVALID_NAME) context() string {
	return e.ctx
}

func errStr(strs ...string) string {
	var msg string
	for _, v := range strs {
		if v != "" {
			msg = fmt.Sprintf("%v: %v", msg, v)
		}
	}
	return msg
}

func ErrorsCast(rcl_ret_t C.rcl_ret_t) RCLError {
	return ErrorsCastC(rcl_ret_t, "")
}
func ErrorsCastC(rcl_ret_t C.rcl_ret_t, context string) RCLError {
	// https://stackoverflow.com/questions/9928221/table-of-functions-vs-switch-in-golang
	// switch-case is faster thanks to compiler optimization than a dispatcher?
	switch rcl_ret_t {
	case C.RCL_RET_ERROR:
		return &RCL_RET_ERROR{(int)(rcl_ret_t), context}
	case C.RCL_RET_INVALID_ARGUMENT:
		return &RCL_RET_INVALID_ARGUMENT{(int)(rcl_ret_t), context}
	case C.RCL_RET_ALREADY_INIT:
		return &RCL_RET_ALREADY_INIT{(int)(rcl_ret_t), context}
	case C.RCL_RET_TOPIC_NAME_INVALID:
		return &RCL_RET_TOPIC_NAME_INVALID{(int)(rcl_ret_t), context}
	case C.RCL_RET_NODE_INVALID_NAME:
		return &RCL_RET_NODE_INVALID_NAME{(int)(rcl_ret_t), context}
	default:
		return &RCL_RET_GOLANG_UNKNOWN_RET_TYPE{(int)(rcl_ret_t), ""}
	}
}

/*
// Copyright 2014-2018 Open Source Robotics Foundation, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef RMW__RET_TYPES_H_
#define RMW__RET_TYPES_H_

#ifdef __cplusplus
extern "C"
{
#endif

#include <stdint.h>

/// Return code for rmw functions
typedef int32_t rmw_ret_t;
/// The operation ran as expected
#define RMW_RET_OK 0
/// Generic error to indicate operation could not complete successfully
#define RMW_RET_ERROR 1
/// The operation was halted early because it exceeded its timeout critera
#define RMW_RET_TIMEOUT 2
/// The operation or event handling is not supported.
#define RMW_RET_UNSUPPORTED 3

/// Failed to allocate memory
#define RMW_RET_BAD_ALLOC 10
/// Argument to function was invalid
#define RMW_RET_INVALID_ARGUMENT 11
/// Incorrect rmw implementation.
#define RMW_RET_INCORRECT_RMW_IMPLEMENTATION 12

// rmw node specific ret codes in 2XX
/// Failed to find node name
// Using same return code than in rcl
#define RMW_RET_NODE_NAME_NON_EXISTENT 203

#ifdef __cplusplus
}
#endif

#endif  // RMW__RET_TYPES_H_




// Copyright 2014 Open Source Robotics Foundation, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef RCL__TYPES_H_
#define RCL__TYPES_H_

#include <rmw/types.h>

typedef rmw_ret_t rcl_ret_t;
/// Success return code.
#define RCL_RET_OK RMW_RET_OK
/// Unspecified error return code.
#define RCL_RET_ERROR RMW_RET_ERROR
/// Timeout occurred return code.
#define RCL_RET_TIMEOUT RMW_RET_TIMEOUT
/// Failed to allocate memory return code.
#define RCL_RET_BAD_ALLOC RMW_RET_BAD_ALLOC
/// Invalid argument return code.
#define RCL_RET_INVALID_ARGUMENT RMW_RET_INVALID_ARGUMENT
/// Unsupported return code.
#define RCL_RET_UNSUPPORTED RMW_RET_UNSUPPORTED

// rcl specific ret codes start at 100
/// rcl_init() already called return code.
#define RCL_RET_ALREADY_INIT 100
/// rcl_init() not yet called return code.
#define RCL_RET_NOT_INIT 101
/// Mismatched rmw identifier return code.
#define RCL_RET_MISMATCHED_RMW_ID 102
/// Topic name does not pass validation.
#define RCL_RET_TOPIC_NAME_INVALID 103
/// Service name (same as topic name) does not pass validation.
#define RCL_RET_SERVICE_NAME_INVALID 104
/// Topic name substitution is unknown.
#define RCL_RET_UNKNOWN_SUBSTITUTION 105
/// rcl_shutdown() already called return code.
#define RCL_RET_ALREADY_SHUTDOWN 106

// rcl node specific ret codes in 2XX
/// Invalid rcl_node_t given return code.
#define RCL_RET_NODE_INVALID 200
#define RCL_RET_NODE_INVALID_NAME 201
#define RCL_RET_NODE_INVALID_NAMESPACE 202
/// Failed to find node name
#define RCL_RET_NODE_NAME_NON_EXISTENT 203

// rcl publisher specific ret codes in 3XX
/// Invalid rcl_publisher_t given return code.
#define RCL_RET_PUBLISHER_INVALID 300

// rcl subscription specific ret codes in 4XX
/// Invalid rcl_subscription_t given return code.
#define RCL_RET_SUBSCRIPTION_INVALID 400
/// Failed to take a message from the subscription return code.
#define RCL_RET_SUBSCRIPTION_TAKE_FAILED 401

// rcl service client specific ret codes in 5XX
/// Invalid rcl_client_t given return code.
#define RCL_RET_CLIENT_INVALID 500
/// Failed to take a response from the client return code.
#define RCL_RET_CLIENT_TAKE_FAILED 501

// rcl service server specific ret codes in 6XX
/// Invalid rcl_service_t given return code.
#define RCL_RET_SERVICE_INVALID 600
/// Failed to take a request from the service return code.
#define RCL_RET_SERVICE_TAKE_FAILED 601

// rcl guard condition specific ret codes in 7XX

// rcl timer specific ret codes in 8XX
/// Invalid rcl_timer_t given return code.
#define RCL_RET_TIMER_INVALID 800
/// Given timer was canceled return code.
#define RCL_RET_TIMER_CANCELED 801

// rcl wait and wait set specific ret codes in 9XX
/// Invalid rcl_wait_set_t given return code.
#define RCL_RET_WAIT_SET_INVALID 900
/// Given rcl_wait_set_t is empty return code.
#define RCL_RET_WAIT_SET_EMPTY 901
/// Given rcl_wait_set_t is full return code.
#define RCL_RET_WAIT_SET_FULL 902

// rcl argument parsing specific ret codes in 1XXX
/// Argument is not a valid remap rule
#define RCL_RET_INVALID_REMAP_RULE 1001
/// Expected one type of lexeme but got another
#define RCL_RET_WRONG_LEXEME 1002
/// Found invalid ros argument while parsing
#define RCL_RET_INVALID_ROS_ARGS 1003
/// Argument is not a valid parameter rule
#define RCL_RET_INVALID_PARAM_RULE 1010
/// Argument is not a valid log level rule
#define RCL_RET_INVALID_LOG_LEVEL_RULE 1020

// rcl event specific ret codes in 20XX
/// Invalid rcl_event_t given return code.
#define RCL_RET_EVENT_INVALID 2000
/// Failed to take an event from the event handle
#define RCL_RET_EVENT_TAKE_FAILED 2001

/// typedef for rmw_serialized_message_t;
typedef rmw_serialized_message_t rcl_serialized_message_t;

#endif  // RCL__TYPES_H_
*/
