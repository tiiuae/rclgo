/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

package rclgo

/*
#include <rcl/types.h>
#include <rcutils/error_handling.h>
*/
import "C"

import (
	"fmt"
)

type rclError struct {
	rclRetCode int
	context    string
	trace      string
}

func (e *rclError) Error() string {
	return e.context
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
func errorString() string {
	rcutils_error_string_str := C.rcutils_get_error_string().str // TODO: Do I need to free this or not?

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
	// upt := (*[1024]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&rcutils_error_string_str))))
	// return string((*upt)[:])
}

func errorsBuildContext(e error, ctx string, stackTrace string) string {
	return fmt.Sprintf("[%T]", e) + " " + ctx + " " + errorString() + "\n" + stackTrace + "\n"
}

func errorsCast(rcl_ret_t C.rcl_ret_t) error {
	return errorsCastC(rcl_ret_t, "")
}

func onErr(err *error, f func() error) {
	if *err != nil {
		f() //nolint:errcheck
	}
}

type closeError string

func (e closeError) Error() string {
	return string(e)
}

func closeErr(s string) error {
	return closeError("tried to close a closed " + s)
}
