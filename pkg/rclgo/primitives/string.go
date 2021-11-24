/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package primitives

/*
#cgo LDFLAGS: -L/opt/ros/galactic/lib -Wl,-rpath=/opt/ros/galactic/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo CFLAGS: -I/opt/ros/galactic/include

#include "rosidl_runtime_c/string.h"

*/
import "C"
import (
	"strings"
	"unsafe"
)

func StringAsCStruct(dst unsafe.Pointer, m string) {
	mem := (*C.rosidl_runtime_c__String)(dst) //TODO add this to template generator

	mem.data = (*C.char)(C.malloc((C.size_t)(C.sizeof_char * uintptr(len(m)+1))))

	for i := 0; i < len(m); i++ {
		stringSetDataCArrayIndex(mem, i, m[i])
	}
	stringSetDataCArrayIndex(mem, len(m), '\x00')
	mem.size = C.size_t(len(m))
	mem.capacity = C.size_t(len(m) + 1)
}

func StringAsGoStruct(m *string, ros2_message_buffer unsafe.Pointer) {
	mem := (*C.rosidl_runtime_c__String)(ros2_message_buffer)
	sb := strings.Builder{}
	sb.Grow(int(mem.size))
	sb.Write((*[1 << 30]byte)(unsafe.Pointer(mem.data))[:mem.size])
	*m = sb.String()
}

func stringSetDataCArrayIndex(mem *C.rosidl_runtime_c__String, i int, v byte) {
	cIdx := (*C.uint8_t)(unsafe.Pointer(
		uintptr(unsafe.Pointer(mem.data)) + (C.sizeof_uint8_t * uintptr(i)),
	))
	*cIdx = (C.uint8_t)(v)
}

type CString = C.rosidl_runtime_c__String
type CString__Sequence = C.rosidl_runtime_c__String__Sequence

func String__Sequence_to_Go(goSlice *[]string, cSlice CString__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]string, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.rosidl_runtime_c__String)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rosidl_runtime_c__String * uintptr(i)),
		))
		StringAsGoStruct(&(*goSlice)[i], unsafe.Pointer(cIdx))
	}
}

func String__Sequence_to_C(cSlice *CString__Sequence, goSlice []string) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.rosidl_runtime_c__String)(C.malloc((C.size_t)(C.sizeof_struct_rosidl_runtime_c__String * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.rosidl_runtime_c__String)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rosidl_runtime_c__String * uintptr(i)),
		))
		StringAsCStruct(unsafe.Pointer(cIdx), v)
	}
}

func String__Array_to_Go(goSlice []string, cSlice []CString) {
	for i := 0; i < len(cSlice); i++ {
		StringAsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}

func String__Array_to_C(cSlice []CString, goSlice []string) {
	for i := 0; i < len(goSlice); i++ {
		StringAsCStruct(unsafe.Pointer(&cSlice[i]), goSlice[i])
	}
}
