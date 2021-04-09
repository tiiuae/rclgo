/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package ros2types

import (
	"unsafe"
)

type ROS2Msg interface {
	TypeSupport() unsafe.Pointer //*C.rosidl_message_type_support_t
	PrepareMemory() unsafe.Pointer
	ReleaseMemory(unsafe.Pointer)
	AsCStruct() unsafe.Pointer
	AsGoStruct(unsafe.Pointer)
	Clone() ROS2Msg
	SetDefaults(interface{}) ROS2Msg // func parameter should always be nil
}
