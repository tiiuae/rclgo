/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lstd_msgs__rosidl_generator_c -lstd_msgs__rosidl_typesupport_c -lrcutils -lrmw_implementation -lpx4_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_typesupport_c -lnav_msgs__rosidl_generator_c
#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rcl/validate_topic_name.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func ValidateTopicName(topic_name string) RCLError {
	var validation_result C.int
	var invalid_index C.size_t
	C.rcl_validate_topic_name_with_size(C.CString(topic_name), (C.ulong)(len(topic_name)), &validation_result, &invalid_index)
	if validation_result != 0 {
		var error_description *C.char = C.rcl_topic_name_validation_result_string(validation_result)
		defer C.free(unsafe.Pointer(error_description))
		return ErrorsCastC(C.RCL_RET_TOPIC_NAME_INVALID, fmt.Sprintf("rcl_validate_topic_name_with_size() failed for topic_name='%s', in index='%d', because: '%s'", topic_name, invalid_index, C.GoString(error_description)))
	}
	return nil
}
