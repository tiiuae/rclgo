package rosidl_runtime_c

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo CFLAGS: -I/opt/ros/foxy/include

#include "rosidl_runtime_c/string.h"
#include "rosidl_runtime_c/primitives_sequence.h"

*/
import "C"
import (
	"unsafe"
)

// Char
type CChar = C.schar
type CChar__Sequence = C.rosidl_runtime_c__char__Sequence

func Char__Sequence_to_Go(goSlice *[]byte, cSlice CChar__Sequence) {
	*goSlice = make([]byte, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.schar)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_schar * uintptr(i)),
		))
		(*goSlice)[i] = byte(*cIdx)
	}
}
func Char__Sequence_to_C(cSlice *CChar__Sequence, goSlice []byte) {
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
