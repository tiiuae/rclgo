/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo LDFLAGS: /home/kivilahtio/install/rclc/lib/librclc.so
#cgo CFLAGS: -I/opt/ros/foxy/include -I/home/kivilahtio/install/rclc/include/

#include <rcl/types.h>
#include <rcutils/error_handling.h>

*/
import "C"
import (
	"container/list"
	"fmt"
	"strings"
)

/*
RCLErrors is a list of errors for functions which could return multiple different errors, wrapped in a tight package, easy-to-code.
*/
type RCLErrors struct {
	list.List
	i *list.Element
}

func (self *RCLErrors) Next() RCLError {
	if self.i == nil {
		self.i = self.Front()
	}
	n := self.i.Next()
	if n != nil {
		e := n.Value.(RCLError)
		return e
	}
	return nil
}
func (self *RCLErrors) Put(e RCLError) *RCLErrors {
	self.PushBack(e)
	return self
}
func (self *RCLErrors) String() string {
	sb := strings.Builder{}
	sb.WriteString("RCLErrors happened:\n")
	for e := self.List.Front(); e != nil; e.Next() {
		err := e.Value.(RCLError)
		sb.WriteString(err.context() + "\n")
	}
	return sb.String()
}

/*
RCLErrorsPut has initialization, incrementation, the jizz, jazz and brass all in one! Incredible! Amazing!
*/
func RCLErrorsPut(rclErrors *RCLErrors, e RCLError) *RCLErrors {
	if rclErrors == nil {
		rclErrors = &RCLErrors{}
	}
	return rclErrors.Put(e)
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

type RCLError interface {
	Error() string // Error implements the Golang Error-interface
	rcl_ret() int
	context() string
}

type RCL_RET_struct struct {
	rcl_ret_t int
	ctx       string
	trace     string
}

func (e *RCL_RET_struct) Error() string {
	return e.ctx
}
func (e *RCL_RET_struct) Trace() string {
	return e.trace
}
func (e *RCL_RET_struct) context() string {
	return e.ctx
}
func (e *RCL_RET_struct) rcl_ret() int {
	return e.rcl_ret_t
}

/// Return the error message followed by `, at <file>:<line>` if set, else "error not set".
/**
 * This function is "safe" because it returns a copy of the current error
 * string or one containing the string "error not set" if no error was set.
 * This ensures that the copy is owned by the calling thread and is therefore
 * never invalidated by other error handling calls, and that the C string
 * inside is always valid and null terminated.
 *
 * \return The current error string, with file and line number, or "error not set" if not set.
 */
func ErrorString() string {
	var rcutils_error_string_str = C.rcutils_get_error_string().str // TODO: Do I need to free this or not?

	// Because the C string is null-terminated, we need to find the NULL-character to know where the string ends.
	// Otherwise we create a Go string of length 1024 of NULLs and gibberish
	bytes := make([]byte, len(rcutils_error_string_str))
	for i := 0; i < len(rcutils_error_string_str); i++ {
		if byte(rcutils_error_string_str[i]) == 0x00 {
			return string(bytes[:i]) // ending slice offset is exclusive
		}
		bytes[i] = byte(rcutils_error_string_str[i])
	}
	return string(bytes)

	// This would be much faster I guess.
	//upt := (*[1024]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&rcutils_error_string_str))))
	//return string((*upt)[:])
}

/// Reset the error state by clearing any previously set error state.
func ErrorReset() {
	C.rcutils_reset_error()
}

func errorsBuildContext(e RCLError, ctx string, stackTrace string) string {
	return fmt.Sprintf("[%T]", e) + " " + ctx + " " + ErrorString() + "\n" + stackTrace + "\n"
}
func ErrorsCast(rcl_ret_t C.rcl_ret_t) RCLError {
	return ErrorsCastC(rcl_ret_t, "")
}
