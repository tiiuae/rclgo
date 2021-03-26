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

// Float
type CFloat = C.float
type Crosidl_runtime_c__float__Sequence = C.rosidl_runtime_c__float__Sequence

func Float__Sequence_to_Go(goSlice *[]float32, cSlice Crosidl_runtime_c__float__Sequence) {
	*goSlice = make([]float32, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.float)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_float * uintptr(i)),
		))
		(*goSlice)[i] = float32(*cIdx)
	}
}
func Float__Sequence_to_C(cSlice *Crosidl_runtime_c__float__Sequence, goSlice []float32) {
	cSlice.data = (*C.float)(C.malloc((C.size_t)(C.sizeof_float * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.float)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_float * uintptr(i)),
		))
		*cIdx = (C.float)(v)
	}
}
func Float__Array_to_Go(goSlice []float32, cSlice []CFloat) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i] = float32(cSlice[i])
	}
}
func Float__Array_to_C(cSlice []CFloat, goSlice []float32) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = C.float(goSlice[i])
	}
}

// Int64
type CInt64 = C.int64_t
type Crosidl_runtime_c__int64__Sequence = C.rosidl_runtime_c__int64__Sequence

func Int64__Sequence_to_Go(goSlice *[]int64, cSlice Crosidl_runtime_c__int64__Sequence) {
	*goSlice = make([]int64, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.int64_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_int64_t * uintptr(i)),
		))
		(*goSlice)[i] = int64(*cIdx)
	}
}
func Int64__Sequence_to_C(cSlice *Crosidl_runtime_c__int64__Sequence, goSlice []int64) {
	cSlice.data = (*C.int64_t)(C.malloc((C.size_t)(C.sizeof_int64_t * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.int64_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_int64_t * uintptr(i)),
		))
		*cIdx = (C.int64_t)(v)
	}
}
func Int64__Array_to_Go(goSlice []int64, cSlice []CInt64) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i] = int64(cSlice[i])
	}
}
func Int64__Array_to_C(cSlice []CInt64, goSlice []int64) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = C.int64_t(goSlice[i])
	}
}
