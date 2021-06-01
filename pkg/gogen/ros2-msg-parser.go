/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

var goKeywords = map[string]struct{}{
	"break": {}, "default": {}, "func": {}, "interface": {}, "select": {},
	"case": {}, "defer": {}, "go": {}, "map": {}, "struct": {},
	"chan": {}, "else": {}, "goto": {}, "package": {}, "switch": {},
	"const": {}, "fallthrough": {}, "if": {}, "range": {}, "type": {},
	"continue": {}, "for": {}, "import": {}, "return": {}, "var": {},
}

func CName(rosName string) string {
	if _, ok := goKeywords[rosName]; ok {
		return "_" + rosName
	}
	return rosName
}

func parseService(service *ROS2Service, source string) error {
	currentMsg := service.Request
	for i, line := range strings.Split(source, "\n") {
		if line == "---" {
			if currentMsg == service.Response {
				return errors.New("too many '---' delimeters")
			}
			currentMsg = service.Response
			continue
		}
		if err := parseLine(currentMsg, line); err != nil {
			return lineErr(i+1, err)
		}
	}
	return nil
}

func parseLine(msg *ROS2Message, line string) error {
	obj, err := parseMessageLine(line, msg)
	if err != nil {
		return err
	}

	switch obj := obj.(type) {
	case *ROS2Constant:
		msg.Constants = append(msg.Constants, obj)
	case *ROS2Field:
		msg.Fields = append(msg.Fields, obj)
		switch obj.PkgName {
		case "":
		case ".":
		case "time":
			msg.GoImports["time"] = ""
		case "rosidl_runtime_c":
			msg.GoImports["github.com/tiiuae/rclgo/pkg/ros2/"+obj.PkgName] = obj.GoPkgName
		default:
			msg.GoImports["github.com/tiiuae/rclgo/pkg/ros2/msgs/"+obj.PkgName+"/msg"] = obj.GoPkgName
			msg.CImports[obj.PkgName] = true
		}
	case nil:
	default:
		return fmt.Errorf("Couldn't parse the input row '%s'", line)
	}
	return nil
}

func lineErr(line int, err error) error {
	return fmt.Errorf("error on line %d: %w", line, err)
}

// ParseROS2Message parses a message definition.
func ParseROS2Message(res *ROS2Message, content string) error {
	for i, line := range strings.Split(content, "\n") {
		if err := parseLine(res, line); err != nil {
			return lineErr(i+1, err)
		}
	}
	return nil
}

var ros2messagesCommentsBuffer = strings.Builder{} // Collect pre-field comments here to be included in the comments. Flushed on empty lines.

func parseMessageLine(testRow string, ros2msg *ROS2Message) (interface{}, error) {
	testRow = strings.TrimSpace(testRow)

	re.R(&testRow, `m!^#\s*(.*)$!`) // Extract comments from comment-only lines to be included in the pre-field comments
	if re.R0.Matches > 0 {
		if re.R0.S[1] != "" {
			ros2messagesCommentsBuffer.WriteString(re.R0.S[1])
		}
		return nil, nil
	}
	if testRow == "" { // do not process empty lines or comment lines
		ros2messagesCommentsBuffer.Reset()
		return nil, nil
	}

	typeChar, capture := isRowConstantOrField(testRow, ros2msg)
	switch typeChar {
	case 'c':
		con, err := ParseROS2MessageConstant(capture, ros2msg)
		if err == nil {
			return con, err
		}
	case 'f':
		f, err := ParseROS2MessageField(capture, ros2msg)
		if err == nil {
			return f, err
		}
	}
	return nil, fmt.Errorf("Couldn't parse the input row as either ROS2 Field or Constant? input '%s'", testRow)
}

func isRowConstantOrField(textRow string, ros2msg *ROS2Message) (byte, map[string]string) {
	re.R(&textRow, `m!
	# This regex might be overly complex for constant parsing but it works. It was originally a copy-paste of the ROS2 field-parser
	^
	(?:(?P<package>\w+)/)?
	(?P<type>\w+)
	(?P<array>\[(?P<bounded><=)?(?P<size>\d*)\])?
	\s+
	(?P<field>\w+)
	\s*=\s*
	(?P<default>[^#]+)?
	(?:\s*\#\s*(?P<comment>.*))?
	$
	!x`)
	if re.R0.Matches > 0 {
		return 'c', re.R0.Z
	}

	re.R(&textRow, `m!
	^(?:(?P<package>\w+)/)?
	(?P<type>\w+)
	(?P<array>\[(?P<bounded><=)?(?P<size>\d*)\])?
	(?P<bounded><=)?                                    # Special case for bounded strings
	(?P<size>\d*)?
	\s+
	(?P<field>\w+)
	\s*
	(?P<default>[^#]+)?
	(?:\s+\#\s*(?P<comment>.*))?
	$
	!x`)
	if re.R0.Matches > 0 {
		return 'f', re.R0.Z
	}

	return 'e', nil
}

func ParseROS2MessageConstant(capture map[string]string, ros2msg *ROS2Message) (*ROS2Constant, error) {
	d := &ROS2Constant{
		RosType: capture["type"],
		RosName: capture["field"],
		Value:   strings.TrimSpace(capture["default"]),
		Comment: commentSerializer(capture["comment"], &ros2messagesCommentsBuffer),
	}

	t, ok := ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING[d.RosType]
	if !ok {
		d.GoType = fmt.Sprintf("<MISSING translation from ROS2 Constant type '%s'>", d.RosType)
		return d, fmt.Errorf("Unknown ROS2 Constant type '%s'\n", d.RosType)
	}
	d.GoType = t.GoType
	d.PkgName = t.PackageName
	return d, nil
}

func ParseROS2MessageField(capture map[string]string, ros2msg *ROS2Message) (*ROS2Field, error) {
	size, err := strconv.ParseInt(capture["size"], 10, 32)
	if err != nil && capture["size"] != "" {
		return nil, err
	}
	if capture["bounded"] != "" {
		capture["array"] = strings.Replace(capture["array"], capture["bounded"]+capture["size"], "", 1)
		capture["bounded"] = capture["bounded"] + capture["size"]
		size = 0
	}
	f := &ROS2Field{
		Comment:      commentSerializer(capture["comment"], &ros2messagesCommentsBuffer),
		GoName:       SnakeToCamel(capture["field"]),
		RosName:      capture["field"],
		CName:        CName(capture["field"]),
		RosType:      capture["type"],
		TypeArray:    capture["array"],
		ArrayBounded: capture["bounded"],
		ArraySize:    int(size),
		DefaultValue: capture["default"],
		PkgName:      capture["package"],
	}

	f.PkgName, f.CType, f.GoType = translateROS2Type(f, ros2msg)
	f.GoPkgName = f.PkgName
	switch f.PkgName {
	case "", "time", "rosidl_runtime_c":
	case ".":
		if ros2msg.Type == "msg" {
			f.PkgIsLocal = true
		} else {
			f.PkgName = ros2msg.Package
			f.GoPkgName = ros2msg.Package + "_msg"
		}
	default:
		f.GoPkgName = f.PkgName + "_msg"
	}
	// Prepopulate extra Go imports
	cSerializationCode(f, ros2msg)
	goSerializationCode(f, ros2msg)

	return f, nil
}

func translateROS2Type(f *ROS2Field, m *ROS2Message) (pkgName string, cType string, goType string) {
	t, ok := ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING[f.RosType]
	if ok {
		f.RosType = t.RosType
		return t.PackageName, t.CType, t.GoType
	}

	if f.PkgName == "" && m.Type != "msg" {
		return m.Package, f.RosType, f.RosType
	}

	// explicit package
	if f.PkgName != "" {
		// type of same package
		if f.PkgName == m.Package {
			return ".", f.RosType, f.RosType
		}

		// type of other package
		return f.PkgName, f.RosType, f.RosType
	}

	// implicit package, type of std_msgs
	if m.Package != "std_msgs" {
		switch f.RosType {
		case "Bool", "ColorRGBA",
			"Duration", "Empty", "Float32MultiArray", "Float32",
			"Float64MultiArray", "Float64", "Header", "Int8MultiArray",
			"Int8", "Int16MultiArray", "Int16", "Int32MultiArray", "Int32",
			"Int64MultiArray", "Int64", "MultiArrayDimension", "MultiarrayLayout",
			"String", "Time", "UInt8MultiArray", "UInt8", "UInt16MultiArray", "UInt16",
			"UInt32MultiArray", "UInt32", "UInt64MultiArray", "UInt64":
			return "std_msgs", f.RosType, f.RosType
		}
	}

	// These are not actually primitive types, but same-package complex types.
	return ".", f.RosType, f.RosType
}

func cSerializationCode(f *ROS2Field, m *ROS2Message) string {

	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Array local package reference
		return ucFirst(f.RosType) + `__Array_to_C(mem.` + f.CName + `[:], t.` + f.GoName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Array remote package reference
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + f.GoPkgReference() + ucFirst(f.RosType) + `__Array_to_C(*(*[]` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)), t.` + f.GoName + `[:])`
	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Slice local package reference
		return ucFirst(f.RosType) + `__Sequence_to_C(&mem.` + f.CName + `, t.` + f.GoName + `)`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Slice remote package reference
		return f.GoPkgReference() + ucFirst(f.RosType) + `__Sequence_to_C((*` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)), t.` + f.GoName + `)`

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return `mem.` + f.CName + ` = *(*C.` + cStructName(f, m) + `)(t.` + f.GoName + `.AsCStruct())`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + `rosidl_runtime_c.` + ucFirst(f.RosType) + `__Array_to_C(*(*[]rosidl_runtime_c.C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)), t.` + f.GoName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + ucFirst(f.RosType) + `__Sequence_to_C((*rosidl_runtime_c.C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)), t.` + f.GoName + `)`

	} else if f.TypeArray == "" && f.PkgName == "" {
		// Primitive value single
		return `mem.` + f.CName + ` = C.` + f.CType + `(t.` + f.GoName + `)`

	}
	return "//<MISSING cSerializationCode!!>"
}

func cStructName(f *ROS2Field, m *ROS2Message) string {
	if f.PkgName == "rosidl_runtime_c" {
		return "rosidl_runtime_c__" + f.CType
	} else if f.PkgName != "" {
		if f.PkgIsLocal {
			return m.Package + "__msg__" + f.CType
		} else {
			return f.PkgName + "__msg__" + f.CType
		}
	}
	return "<MISSING cStructName!!>"
}

func goSerializationCode(f *ROS2Field, m *ROS2Message) string {

	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Array local package reference
		return ucFirst(f.RosType) + `__Array_to_Go(t.` + f.GoName + `[:], mem.` + f.CName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" {
		// Complex value Array remote package reference
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + f.GoPkgReference() + ucFirst(f.RosType) + `__Array_to_Go(t.` + f.GoName + `[:], *(*[]` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)))`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Slice local package reference
		return ucFirst(f.RosType) + `__Sequence_to_Go(&t.` + f.GoName + `, mem.` + f.CName + `)`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Slice remote package reference
		return f.GoPkgReference() + ucFirst(f.RosType) + `__Sequence_to_Go(&t.` + f.GoName + `, *(*` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return `t.` + f.GoName + `.AsGoStruct(unsafe.Pointer(&mem.` + f.CName + `))`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + `rosidl_runtime_c.` + ucFirst(f.RosType) + `__Array_to_Go(t.` + f.GoName + `[:], *(*[]rosidl_runtime_c.C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)))`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + ucFirst(f.RosType) + `__Sequence_to_Go(&t.` + f.GoName + `, *(*rosidl_runtime_c.C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray == "" && f.PkgName == "" {
		// Primitive value single
		return `t.` + f.GoName + ` = ` + f.GoType + `(mem.` + f.CName + `)`

	}
	return "//<MISSING goSerializationCode!!>"
}

func defaultCode(f *ROS2Field) string {
	if f.PkgName != "" && f.TypeArray != "" {
		defaultValues := splitMsgDefaultArrayValues(f.RosType, f.DefaultValue)
		// Complex value array and slice
		sb := strings.Builder{}
		var indexesCount int
		if f.ArraySize > 0 {
			indexesCount = f.ArraySize
			sb.Grow(indexesCount)
		} else if len(defaultValues) > 0 { // Init a slice
			indexesCount = len(defaultValues)
			sb.Grow(indexesCount + 1)

			sb.WriteString(`t.` + f.GoName + ` = make(` + f.TypeArray + f.GoPkgReference() + f.GoType + `, ` + strconv.Itoa(indexesCount) + `)` + "\n\t")
		}

		for i := 0; i < indexesCount; i++ {
			defaultValue := "nil"
			if i < len(defaultValues) {
				defaultValue = defaultValues[i]
			}
			sb.WriteString(`t.` + f.GoName + `[` + strconv.Itoa(i) + `].SetDefaults(` + ValOrNil(defaultValue) + ")\n\t")
		}
		return sb.String()

	} else if f.PkgName != "" && f.TypeArray == "" {
		// Complex value single
		return `t.` + f.GoName + `.SetDefaults(` + ValOrNil(defaultValueSanitizer(f.RosType, f.DefaultValue)) + `)` + "\n\t"

	} else if f.DefaultValue != "" && f.TypeArray != "" {
		// Primitive value array
		return `t.` + f.GoName + ` = ` + f.TypeArray + f.GoPkgReference() + f.GoType + `{` + normalizeMsgDefaultArrayValue(f.DefaultValue) + `}` + "\n\t"

	} else if f.DefaultValue != "" {
		// Primitive value single
		return `t.` + f.GoName + ` = ` + f.DefaultValue + "\n\t"

	} else if f.DefaultValue == "" {
		return ""
	}
	return "//<MISSING defaultCode!!>"
}
