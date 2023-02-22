/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package primitives

//#include "rosidl_runtime_c/u16string.h"
import "C"
import (
	"unicode/utf16"
	"unsafe"
)

func U16StringAsCStruct(dst unsafe.Pointer, m string) { // rosidl_runtime_c__U16String__assignn() does something like this, but to call it we still need to make a C string and free it.
	mem := (*C.rosidl_runtime_c__U16String)(dst)
	runescape := utf16.Encode([]rune(m))

	mem.data = (*C.ushort)(C.malloc(C.sizeof_ushort * C.size_t(len(runescape)+1)))
	mem.size = C.size_t(len(runescape))
	mem.capacity = C.size_t(len(runescape) + 1)
	memData := unsafe.Slice((*uint16)(mem.data), mem.capacity)
	copy(memData, runescape)
	memData[len(memData)-1] = 0
}

func U16StringAsGoStruct(msg *string, ros2_message_buffer unsafe.Pointer) {
	mem := (*C.rosidl_runtime_c__U16String)(ros2_message_buffer)

	*msg = string(utf16.Decode(unsafe.Slice((*uint16)(mem.data), mem.size)))
}

type CU16String = C.rosidl_runtime_c__U16String
type CU16String__Sequence = C.rosidl_runtime_c__U16String__Sequence

func U16String__Sequence_to_Go(goSlice *[]string, cSlice CU16String__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]string, int64(cSlice.size))
	src := unsafe.Slice(cSlice.data, cSlice.size)
	for i := 0; i < int(cSlice.size); i++ {
		U16StringAsGoStruct(&(*goSlice)[i], unsafe.Pointer(&src[i]))
	}
}

func U16String__Sequence_to_C(cSlice *CU16String__Sequence, goSlice []string) {
	if len(goSlice) == 0 {
		cSlice.data = nil
		cSlice.capacity = 0
		cSlice.size = 0
		return
	}
	cSlice.data = (*C.rosidl_runtime_c__U16String)(C.malloc((C.size_t)(C.sizeof_struct_rosidl_runtime_c__U16String * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity
	dst := unsafe.Slice(cSlice.data, cSlice.size)
	for i := range goSlice {
		U16StringAsCStruct(unsafe.Pointer(&dst[i]), goSlice[i])
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
