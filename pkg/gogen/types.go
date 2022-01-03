/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"path/filepath"
)

type Metadata struct {
	Name, Package, Type string
}

func (m *Metadata) ImportPath() string {
	return filepath.Join(m.Package, m.Type)
}

func (m *Metadata) GoPackage() string {
	return m.Package + "_" + m.Type
}

/*
ROS2Message is a message definition. https://design.ros2.org/articles/legacy_interface_definition.html
Use ROS2MessageNew() to initialize the struct
*/
type ROS2Message struct {
	*Metadata
	Fields    []*ROS2Field
	Constants []*ROS2Constant
	GoImports map[string]string
	CImports  stringSet
}

func ROS2MessageNew(pkg, name string) *ROS2Message {
	return newMessageWithType(pkg, name, "msg")
}

func newMessageWithType(pkg, name, typ string) *ROS2Message {
	return &ROS2Message{
		Metadata: &Metadata{
			Name:    name,
			Package: pkg,
			Type:    typ,
		},
		GoImports: map[string]string{},
		CImports:  stringSet{},
	}
}

// ROS2Constant is a message definition.
type ROS2Constant struct {
	RosType    string
	GoType     string
	RosName    string
	Value      string
	Comment    string
	PkgName    string
	PkgIsLocal bool
}

func (t *ROS2Constant) GoPkgReference() string {
	if t.PkgName == "" || t.PkgIsLocal {
		return ""
	}
	return t.PkgName + "."
}

// Field is a message field.
type ROS2Field struct {
	TypeArray    string
	ArrayBounded string
	ArraySize    int
	DefaultValue string
	PkgName      string
	GoPkgName    string
	PkgIsLocal   bool
	RosType      string
	CType        string
	GoType       string
	RosName      string
	CName        string
	GoName       string
	Comment      string
}

func (t *ROS2Field) GoPkgReference() string {
	if t.PkgName == "" || t.PkgIsLocal {
		return ""
	}
	return t.GoPkgName + "."
}

func (t *ROS2Field) IsSingleComplex() bool {
	return t.TypeArray == "" && t.PkgName != ""
}

type ROS2Service struct {
	*Metadata
	Request  *ROS2Message
	Response *ROS2Message
}

func NewROS2Service(pkg, name string) *ROS2Service {
	s := &ROS2Service{
		Metadata: &Metadata{
			Name:    name,
			Package: pkg,
			Type:    "srv",
		},
		Request:  newMessageWithType(pkg, name+"_Request", "srv"),
		Response: newMessageWithType(pkg, name+"_Response", "srv"),
	}
	return s
}

/*
ROS2ErrorType must have fields exported otherwise they cannot be used by the test/template -package
*/
type ROS2ErrorType struct {
	Name      string
	Rcl_ret_t string // The function call return value the error is mapped to
	Reference string // This is a reference to another type, so we just redefine the same type with another name
	Comment   string // Any found comments before or over the #definition
}

type rosidl_runtime_c_type_mapping struct {
	RosType     string
	GoType      string
	CType       string
	CStructName string
	PackageName string
	SkipAutogen bool
}

var primitiveTypeMappings = map[string]rosidl_runtime_c_type_mapping{
	"string":   {RosType: "string", GoType: "string", CStructName: "String", CType: "String", SkipAutogen: true},
	"time":     {RosType: "time", GoType: "Time", CStructName: "Time", CType: "time", SkipAutogen: true},
	"duration": {RosType: "duration", GoType: "Duration", CStructName: "Duration", CType: "duration", SkipAutogen: true},
	"float32":  {RosType: "float32", GoType: "float32", CStructName: "float", CType: "float"},
	"float64":  {RosType: "float64", GoType: "float64", CStructName: "double", CType: "double"},
	"bool":     {RosType: "bool", GoType: "bool", CStructName: "boolean", CType: "bool"},
	"byte":     {RosType: "byte", GoType: "byte", CStructName: "octet", CType: "uint8_t"},
	"char":     {RosType: "char", GoType: "byte", CStructName: "char", CType: "uchar", SkipAutogen: true}, // Autogen sequence/array handlers have C-type schar, but char everywhere else is uchar?
	"int8":     {RosType: "int8", GoType: "int8", CStructName: "int8", CType: "int8_t"},
	"int16":    {RosType: "int16", GoType: "int16", CStructName: "int16", CType: "int16_t"},
	"int32":    {RosType: "int32", GoType: "int32", CStructName: "int32", CType: "int32_t"},
	"int64":    {RosType: "int64", GoType: "int64", CStructName: "int64", CType: "int64_t"},
	"uint8":    {RosType: "uint8", GoType: "uint8", CStructName: "uint8", CType: "uint8_t"},
	"uint16":   {RosType: "uint16", GoType: "uint16", CStructName: "uint16", CType: "uint16_t"},
	"uint32":   {RosType: "uint32", GoType: "uint32", CStructName: "uint32", CType: "uint32_t"},
	"uint64":   {RosType: "uint64", GoType: "uint64", CStructName: "uint64", CType: "uint64_t"},
	// for wstring the RosType is actually "wstring", but the way the generator is implemented, this is a reasonable hack to make it work with this fringe-case without extensive refactoring
	"wstring": {RosType: "U16String", GoType: "string", CStructName: "U16String", CType: "U16String", SkipAutogen: true},
}

/*
blacklistedMessages is matched against the paths gogen inspects if it is a ROS2 Message file and needs to be turned into a Go type.
If the path matches the blacklist, it is ignored and a notification is logged.
*/
var blacklistedMessages = []string{
	"libstatistics_collector/msg/DummyMessage",
	"this-is-a-test-blacklist-entry-do-not-remove-used-for-internal-testing",
}

/*
cErrorTypeFiles are looked for #definitions and parsed as Golang ros2 error types
*/
var cErrorTypeFiles = []string{
	"rcl/types.h",
	"rmw/ret_types.h",
}
