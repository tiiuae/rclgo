/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
	http://www.apache.org/licenses/LICENSE-2.0
*/

/*
THIS FILE IS AUTOGENERATED BY 'rclgo-gen generate'
*/

package visualization_msgs
import (
	"unsafe"

	"github.com/tiiuae/rclgo/pkg/ros2/ros2types"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
	geometry_msgs "github.com/tiiuae/rclgo/pkg/ros2/msgs/geometry_msgs/msg"
	rosidl_runtime_c "github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"
	
)
/*
#cgo LDFLAGS: -L/opt/ros/foxy/lib -Wl,-rpath=/opt/ros/foxy/lib -lrcl -lrosidl_runtime_c -lrosidl_typesupport_c -lrcutils -lrmw_implementation
#cgo LDFLAGS: -lvisualization_msgs__rosidl_typesupport_c -lvisualization_msgs__rosidl_generator_c
#cgo LDFLAGS: -lgeometry_msgs__rosidl_typesupport_c -lgeometry_msgs__rosidl_generator_c

#cgo CFLAGS: -I/opt/ros/foxy/include

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <visualization_msgs/msg/interactive_marker_control.h>
*/
import "C"

func init() {
	ros2_type_dispatcher.RegisterROS2MsgTypeNameAlias("visualization_msgs/InteractiveMarkerControl", &InteractiveMarkerControl{})
}
const (
	InteractiveMarkerControl_INHERIT uint8 = 0// Orientation mode: controls how orientation changes.INHERIT: Follow orientation of interactive markerFIXED: Keep orientation fixed at initial stateVIEW_FACING: Align y-z plane with screen (x: forward, y:left, z:up).
	InteractiveMarkerControl_FIXED uint8 = 1// Orientation mode: controls how orientation changes.INHERIT: Follow orientation of interactive markerFIXED: Keep orientation fixed at initial stateVIEW_FACING: Align y-z plane with screen (x: forward, y:left, z:up).
	InteractiveMarkerControl_VIEW_FACING uint8 = 2// Orientation mode: controls how orientation changes.INHERIT: Follow orientation of interactive markerFIXED: Keep orientation fixed at initial stateVIEW_FACING: Align y-z plane with screen (x: forward, y:left, z:up).
	InteractiveMarkerControl_NONE uint8 = 0// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_MENU uint8 = 1// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_BUTTON uint8 = 2// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_MOVE_AXIS uint8 = 3// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_MOVE_PLANE uint8 = 4// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_ROTATE_AXIS uint8 = 5// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_MOVE_ROTATE uint8 = 6// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS.
	InteractiveMarkerControl_MOVE_3D uint8 = 7// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS."3D" interaction modes work with the mouse+SHIFT+CTRL or with 3D cursors.MOVE_3D: Translate freely in 3D space.ROTATE_3D: Rotate freely in 3D space about the origin of parent frame.MOVE_ROTATE_3D: Full 6-DOF freedom of translation and rotation about the cursor origin.
	InteractiveMarkerControl_ROTATE_3D uint8 = 8// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS."3D" interaction modes work with the mouse+SHIFT+CTRL or with 3D cursors.MOVE_3D: Translate freely in 3D space.ROTATE_3D: Rotate freely in 3D space about the origin of parent frame.MOVE_ROTATE_3D: Full 6-DOF freedom of translation and rotation about the cursor origin.
	InteractiveMarkerControl_MOVE_ROTATE_3D uint8 = 9// Interaction mode for this controlNONE: This control is only meant for visualization; no context menu.MENU: Like NONE, but right-click menu is active.BUTTON: Element can be left-clicked.MOVE_AXIS: Translate along local x-axis.MOVE_PLANE: Translate in local y-z plane.ROTATE_AXIS: Rotate around local x-axis.MOVE_ROTATE: Combines MOVE_PLANE and ROTATE_AXIS."3D" interaction modes work with the mouse+SHIFT+CTRL or with 3D cursors.MOVE_3D: Translate freely in 3D space.ROTATE_3D: Rotate freely in 3D space about the origin of parent frame.MOVE_ROTATE_3D: Full 6-DOF freedom of translation and rotation about the cursor origin.
)

// Do not create instances of this type directly. Always use NewInteractiveMarkerControl
// function instead.
type InteractiveMarkerControl struct {
	Name rosidl_runtime_c.String `yaml:"name"`// Identifying string for this control.You need to assign a unique value to this to receive feedback from the GUIon what actions the user performs on this control (e.g. a button click).
	Orientation geometry_msgs.Quaternion `yaml:"orientation"`// Defines the local coordinate frame (relative to the pose of the parentinteractive marker) in which is being rotated and translated.Default: Identity
	OrientationMode uint8 `yaml:"orientation_mode"`
	InteractionMode uint8 `yaml:"interaction_mode"`
	AlwaysVisible bool `yaml:"always_visible"`// If true, the contained markers will also be visiblewhen the gui is not in interactive mode.
	Markers []Marker `yaml:"markers"`// Markers to be displayed as custom visual representation.Leave this empty to use the default control handles.Note:- The markers can be defined in an arbitrary coordinate frame,but will be transformed into the local frame of the interactive marker.- If the header of a marker is empty, its pose will be interpreted asrelative to the pose of the parent interactive marker.
	IndependentMarkerOrientation bool `yaml:"independent_marker_orientation"`// In VIEW_FACING mode, set this to true if you don't want the markersto be aligned with the camera view point. The markers will show upas in INHERIT mode.
	Description rosidl_runtime_c.String `yaml:"description"`// Short description (< 40 characters) of what this control does,e.g. "Move the robot".Default: A generic description based on the interaction mode
}

// NewInteractiveMarkerControl creates a new InteractiveMarkerControl with default values.
func NewInteractiveMarkerControl() *InteractiveMarkerControl {
	self := InteractiveMarkerControl{}
	self.SetDefaults(nil)
	return &self
}

func (t *InteractiveMarkerControl) SetDefaults(d interface{}) ros2types.ROS2Msg {
	t.Name.SetDefaults("")
	t.Orientation.SetDefaults(nil)
	t.Description.SetDefaults("")
	
	return t
}

func (t *InteractiveMarkerControl) TypeSupport() unsafe.Pointer {
	return unsafe.Pointer(C.rosidl_typesupport_c__get_message_type_support_handle__visualization_msgs__msg__InteractiveMarkerControl())
}
func (t *InteractiveMarkerControl) PrepareMemory() unsafe.Pointer { //returns *C.visualization_msgs__msg__InteractiveMarkerControl
	return (unsafe.Pointer)(C.visualization_msgs__msg__InteractiveMarkerControl__create())
}
func (t *InteractiveMarkerControl) ReleaseMemory(pointer_to_free unsafe.Pointer) {
	C.visualization_msgs__msg__InteractiveMarkerControl__destroy((*C.visualization_msgs__msg__InteractiveMarkerControl)(pointer_to_free))
}
func (t *InteractiveMarkerControl) AsCStruct() unsafe.Pointer {
	mem := (*C.visualization_msgs__msg__InteractiveMarkerControl)(t.PrepareMemory())
	mem.name = *(*C.rosidl_runtime_c__String)(t.Name.AsCStruct())
	mem.orientation = *(*C.geometry_msgs__msg__Quaternion)(t.Orientation.AsCStruct())
	mem.orientation_mode = C.uint8_t(t.OrientationMode)
	mem.interaction_mode = C.uint8_t(t.InteractionMode)
	mem.always_visible = C.bool(t.AlwaysVisible)
	Marker__Sequence_to_C(&mem.markers, t.Markers)
	mem.independent_marker_orientation = C.bool(t.IndependentMarkerOrientation)
	mem.description = *(*C.rosidl_runtime_c__String)(t.Description.AsCStruct())
	return unsafe.Pointer(mem)
}
func (t *InteractiveMarkerControl) AsGoStruct(ros2_message_buffer unsafe.Pointer) {
	mem := (*C.visualization_msgs__msg__InteractiveMarkerControl)(ros2_message_buffer)
	t.Name.AsGoStruct(unsafe.Pointer(&mem.name))
	t.Orientation.AsGoStruct(unsafe.Pointer(&mem.orientation))
	t.OrientationMode = uint8(mem.orientation_mode)
	t.InteractionMode = uint8(mem.interaction_mode)
	t.AlwaysVisible = bool(mem.always_visible)
	Marker__Sequence_to_Go(&t.Markers, mem.markers)
	t.IndependentMarkerOrientation = bool(mem.independent_marker_orientation)
	t.Description.AsGoStruct(unsafe.Pointer(&mem.description))
}
func (t *InteractiveMarkerControl) Clone() ros2types.ROS2Msg {
	clone := *t
	return &clone
}

type CInteractiveMarkerControl = C.visualization_msgs__msg__InteractiveMarkerControl
type CInteractiveMarkerControl__Sequence = C.visualization_msgs__msg__InteractiveMarkerControl__Sequence

func InteractiveMarkerControl__Sequence_to_Go(goSlice *[]InteractiveMarkerControl, cSlice CInteractiveMarkerControl__Sequence) {
	if cSlice.size == 0 {
		return
	}
	*goSlice = make([]InteractiveMarkerControl, int64(cSlice.size))
	for i := 0; i < int(cSlice.size); i++ {
		cIdx := (*C.visualization_msgs__msg__InteractiveMarkerControl__Sequence)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_visualization_msgs__msg__InteractiveMarkerControl * uintptr(i)),
		))
		(*goSlice)[i] = InteractiveMarkerControl{}
		(*goSlice)[i].AsGoStruct(unsafe.Pointer(cIdx))
	}
}
func InteractiveMarkerControl__Sequence_to_C(cSlice *CInteractiveMarkerControl__Sequence, goSlice []InteractiveMarkerControl) {
	if len(goSlice) == 0 {
		return
	}
	cSlice.data = (*C.visualization_msgs__msg__InteractiveMarkerControl)(C.malloc((C.size_t)(C.sizeof_struct_visualization_msgs__msg__InteractiveMarkerControl * uintptr(len(goSlice)))))
	cSlice.capacity = C.size_t(len(goSlice))
	cSlice.size = cSlice.capacity

	for i, v := range goSlice {
		cIdx := (*C.visualization_msgs__msg__InteractiveMarkerControl)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cSlice.data)) + (C.sizeof_struct_visualization_msgs__msg__InteractiveMarkerControl * uintptr(i)),
		))
		*cIdx = *(*C.visualization_msgs__msg__InteractiveMarkerControl)(v.AsCStruct())
	}
}
func InteractiveMarkerControl__Array_to_Go(goSlice []InteractiveMarkerControl, cSlice []CInteractiveMarkerControl) {
	for i := 0; i < len(cSlice); i++ {
		goSlice[i].AsGoStruct(unsafe.Pointer(&cSlice[i]))
	}
}
func InteractiveMarkerControl__Array_to_C(cSlice []CInteractiveMarkerControl, goSlice []InteractiveMarkerControl) {
	for i := 0; i < len(goSlice); i++ {
		cSlice[i] = *(*C.visualization_msgs__msg__InteractiveMarkerControl)(goSlice[i].AsCStruct())
	}
}

