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
#include "rosidl_runtime_c/primitives_sequence.h"

*/
import "C"
import (
	"unsafe"
)

/*
Char has some strange naming conventions under the ROS2 IDL hood, so it is easier to define the Char type manually, than refactor the whole generator templating.
*/
type CChar = C.schar
type CChar__Sequence = C.rosidl_runtime_c__char__Sequence

func Char__Sequence_to_Go(goSlice *[]byte, cSlice CChar__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]byte, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.schar)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_schar * uintptr(i)),
		))
		(*goSlice)[i] = byte(*cIdx)
	}
}
func Char__Sequence_to_C(cSlice *CChar__Sequence, goSlice []byte) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.schar)(C.malloc((C.size_t)(C.sizeof_schar * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.schar)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_schar * uintptr(i)),
		))
		*cIdx = (C.schar)(v)
	}
}
func Char__Array_to_Go(goSlice []byte, cSlice []CChar) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i] = byte(cSlice[i])
	}
}
func Char__Array_to_C(cSlice []CChar, goSlice []byte) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = C.schar(goSlice[i])
	}
}
