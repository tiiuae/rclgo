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

#include "rosidl_runtime_c/u16string.h"

*/
import "C"
import (
	"unicode/utf16"
	"unsafe"
)

func U16StringAsCStruct(dst unsafe.Pointer, m string) { // rosidl_runtime_c__U16String__assignn() does something like this, but to call it we still need to make a C string and free it.
	mem := (*C.rosidl_runtime_c__U16String)(dst)
	runescape := utf16.Encode([]rune(m))

	mem.data = (*C.uint_least16_t)(C.malloc((C.size_t)(C.sizeof_uint_least16_t * uintptr(len(runescape)+1))))

	for i := 0; i < len(runescape); i++ {
		u16StringSetDataCArrayIndex(mem, i, runescape[i])
	}
	u16StringSetDataCArrayIndex(mem, len(runescape), '\x00')
	mem.size = C.size_t(len(runescape))
	mem.capacity = C.size_t(len(runescape) + 1)
}

func U16StringAsGoStruct(msg *string, ros2_message_buffer unsafe.Pointer) {
	mem := (*C.rosidl_runtime_c__U16String)(ros2_message_buffer)

	*msg = string(utf16.Decode((*[1 << 30]uint16)(unsafe.Pointer(mem.data))[:mem.size]))
}

func u16StringSetDataCArrayIndex(mem *C.rosidl_runtime_c__U16String, i int, v uint16) {
	cIdx := (*C.uint_least16_t)(unsafe.Pointer(
		uintptr(unsafe.Pointer(mem.data)) + (C.sizeof_uint_least16_t * uintptr(i)),
	))
	*cIdx = (C.uint_least16_t)(v)
}

type CU16String = C.rosidl_runtime_c__U16String
type CU16String__Sequence = C.rosidl_runtime_c__U16String__Sequence

func U16String__Sequence_to_Go(goSlice *[]string, cSlice CU16String__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]string, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.rosidl_runtime_c__U16String)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rosidl_runtime_c__U16String * uintptr(i)),
		))
		U16StringAsGoStruct((&(*goSlice)[i]), unsafe.Pointer(cIdx))
	}
}

func U16String__Sequence_to_C(cSlice *CU16String__Sequence, goSlice []string) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.rosidl_runtime_c__U16String)(C.malloc((C.size_t)(C.sizeof_struct_rosidl_runtime_c__U16String * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.rosidl_runtime_c__U16String)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rosidl_runtime_c__U16String * uintptr(i)),
		))
		U16StringAsCStruct(unsafe.Pointer(cIdx), v)
	}
}

func U16String__Array_to_Go(goSlice []string, cSlice []CU16String) {
	for i := 0; i < len(cSlice); i++ {
		U16StringAsGoStruct(&goSlice[i], unsafe.Pointer(&cSlice[i]))
	}
}

func U16String__Array_to_C(cSlice []CU16String, goSlice []string) {
	for i := 0; i < len(goSlice); i++ {
		U16StringAsCStruct(unsafe.Pointer(&cSlice[i]), goSlice[i])
	}
}
