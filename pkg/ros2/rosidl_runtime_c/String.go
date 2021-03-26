package rosidl_runtime_c

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo CFLAGS: -I/opt/ros/foxy/include

#include "rosidl_runtime_c/string.h"

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type String string

func (t *String) TypeSupport() unsafe.Pointer {
	fmt.Printf("rosidl_runtime_c.TypeSupport() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
	return unsafe.Pointer(t)
}
func (t *String) PrepareMemory() unsafe.Pointer {
	fmt.Printf("rosidl_runtime_c.PrepareMemory() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
	return unsafe.Pointer(t)
}
func (t *String) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	fmt.Printf("rosidl_runtime_c.ReleaseMemory() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
}
