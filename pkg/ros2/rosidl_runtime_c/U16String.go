/*
Due to the way the rcl string representation differs from Go representation, for serdes purposes treat the U16String as ros2types.ROS2Msg
so no special string-specific exceptions need to e made to the already complex ROS2 Msg serdes templating.
*/
package rosidl_runtime_c

/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo CFLAGS: -I/opt/ros/foxy/include

#include "rosidl_runtime_c/u16string.h"

*/
import "C"
import (
	"fmt"
	"unicode/utf16"
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
)

type U16String string

func NewU16String() *U16String {
	self := U16String("")
	self.SetDefaults(nil)
	return &self
}

func (t U16String) Equal2(cmp U16String) bool {
	return t == cmp
}
func (t *U16String) Equal(cmp *U16String) bool {
	return *t == *cmp
}

func (t *U16String) SetDefaults(d interface{}) ros2types.ROS2Msg {
	switch d.(type) {
	case string:
		*t = U16String(d.(string))
	case U16String:
		*t = d.(U16String)
	case nil:
		// *t is already ""
	default:
		panic(fmt.Sprintf("interface conversion: interface {} is %#v, not rosidl_runtime_c.U16String\n", d))
	}
	return t
}

func (t *U16String) TypeSupport() unsafe.Pointer {
	fmt.Printf("rosidl_runtime_c.TypeSupport() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
	return unsafe.Pointer(t)
}
func (t *U16String) PrepareMemory() unsafe.Pointer {
	fmt.Printf("rosidl_runtime_c.PrepareMemory() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
	return unsafe.Pointer(t)
}
func (t *U16String) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	fmt.Printf("rosidl_runtime_c.ReleaseMemory() called. This is never meant to be directly addressed as a stand-alone data object in the ROS2 messaging bus.")
}
func (t *U16String) AsCStruct() unsafe.Pointer { // rosidl_runtime_c__U16String__assignn() does something like this, but to call it we still need to make a C string and free it.
	mem := (*C.rosidl_runtime_c__U16String)(C.malloc(C.sizeof_struct_rosidl_runtime_c__U16String)) //TODO add this to template generator
	runescape := utf16.Encode([]rune(*t))

	mem.data = (*C.uint_least16_t)(C.malloc((C.size_t)(C.sizeof_uint_least16_t * uintptr(len(runescape)+1))))

	for i := 0; i < len(runescape); i++ {
		t.setDataCArrayIndex(mem, i, runescape[i])
	}
	t.setDataCArrayIndex(mem, len(runescape), '\x00')
	mem.size = C.size_t(len(runescape))
	mem.capacity = C.size_t(len(runescape) + 1)
	return unsafe.Pointer(mem)
}
func (t *U16String) setDataCArrayIndex(mem *C.rosidl_runtime_c__U16String, i int, v uint16) {
	cIdx := (*C.uint_least16_t)(unsafe.Pointer(
		uintptr(unsafe.Pointer(mem.data)) + (C.sizeof_uint_least16_t * uintptr(i)),
	))
	*cIdx = (C.uint_least16_t)(v)
}
func (t *U16String) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.rosidl_runtime_c__U16String)(ros2_message_buffer)
	sb := make([]uint16, int(mem.size))
	for i := 0; i < int(mem.size); i++ {
		cIdx := (*C.uint_least16_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(mem.data)) + (C.sizeof_uint_least16_t * uintptr(i)),
		))
		sb[i] = (uint16(*cIdx))
	}
	*t = U16String(utf16.Decode(sb))
}
func (t *U16String) Clone() ros2types.ROS2Msg {
	c := *t
	return &c
}

type CU16String = C.rosidl_runtime_c__U16String
type CU16String__Sequence = C.rosidl_runtime_c__U16String__Sequence

func U16String__Sequence_to_Go(goSlice *[]U16String, cSlice CU16String__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]U16String, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.rosidl_runtime_c__U16String)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_rosidl_runtime_c__U16String * uintptr(i)),
		))
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func U16String__Sequence_to_C(cSlice *CU16String__Sequence, goSlice []U16String) {
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
		*cIdx = *(*C.rosidl_runtime_c__U16String)(v.AsCStruct())
	}
}
func U16String__Array_to_Go(goSlice []U16String, cSlice []CU16String) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func U16String__Array_to_C(cSlice []CU16String, goSlice []U16String) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.rosidl_runtime_c__U16String)(goSlice[i].AsCStruct())
	}
}
