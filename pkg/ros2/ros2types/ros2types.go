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
