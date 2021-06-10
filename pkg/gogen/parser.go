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

func cName(rosName string) string {
	if _, ok := goKeywords[rosName]; ok {
		return "_" + rosName
	}
	return rosName
}

type parser struct {
	config *Config
}

func (p *parser) parseService(service *ROS2Service, source string) error {
	currentMsg := service.Request
	for i, line := range strings.Split(source, "\n") {
		if line == "---" {
			if currentMsg == service.Response {
				return errors.New("too many '---' delimeters")
			}
			currentMsg = service.Response
			continue
		}
		if err := p.parseLine(currentMsg, line); err != nil {
			return lineErr(i+1, err)
		}
	}
	return nil
}

func (p *parser) parseLine(msg *ROS2Message, line string) error {
	obj, err := p.parseMessageLine(line, msg)
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
		case "primitives":
			msg.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/"+obj.PkgName] = obj.GoPkgName
		default:
			msg.GoImports[p.config.MessageModulePrefix+"/"+obj.PkgName+"/msg"] = obj.GoPkgName
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
func (p *parser) ParseROS2Message(res *ROS2Message, content string) error {
	for i, line := range strings.Split(content, "\n") {
		if err := p.parseLine(res, line); err != nil {
			return lineErr(i+1, err)
		}
	}
	return nil
}

var ros2messagesCommentsBuffer = strings.Builder{} // Collect pre-field comments here to be included in the comments. Flushed on empty lines.

func (p *parser) parseMessageLine(testRow string, ros2msg *ROS2Message) (interface{}, error) {
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
		f, err := p.ParseROS2MessageField(capture, ros2msg)
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

	t, ok := primitiveTypeMappings[d.RosType]
	if !ok {
		d.GoType = fmt.Sprintf("<MISSING translation from ROS2 Constant type '%s'>", d.RosType)
		return d, fmt.Errorf("Unknown ROS2 Constant type '%s'\n", d.RosType)
	}
	d.GoType = t.GoType
	d.PkgName = t.PackageName
	return d, nil
}

func (p *parser) ParseROS2MessageField(capture map[string]string, ros2msg *ROS2Message) (*ROS2Field, error) {
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
		GoName:       snakeToCamel(capture["field"]),
		RosName:      capture["field"],
		CName:        cName(capture["field"]),
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
	case "", "time", "primitives":
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
	p.cSerializationCode(f, ros2msg)
	p.goSerializationCode(f, ros2msg)

	return f, nil
}

func translateROS2Type(f *ROS2Field, m *ROS2Message) (pkgName string, cType string, goType string) {
	t, ok := primitiveTypeMappings[f.RosType]
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

func (p *parser) cSerializationCode(f *ROS2Field, m *ROS2Message) string {
	if f.PkgName == "" {
	}
	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Array local package reference
		return ucFirst(f.RosType) + `__Array_to_C(mem.` + f.CName + `[:], m.` + f.GoName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Array remote package reference
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + f.GoPkgReference() + ucFirst(f.RosType) + `__Array_to_C(*(*[]` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)), m.` + f.GoName + `[:])`
	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Slice local package reference
		return ucFirst(f.RosType) + `__Sequence_to_C(&mem.` + f.CName + `, m.` + f.GoName + `)`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Slice remote package reference
		return f.GoPkgReference() + ucFirst(f.RosType) + `__Sequence_to_C((*` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)), m.` + f.GoName + `)`

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return f.GoPkgReference() + f.GoType + "TypeSupport.AsCStruct(unsafe.Pointer(&mem." + f.CName + "), &m." + f.GoName + ")"

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + `primitives.` + ucFirst(f.RosType) + `__Array_to_C(*(*[]primitives.C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)), m.` + f.GoName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
		return `primitives.` + ucFirst(f.RosType) + `__Sequence_to_C((*primitives.C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)), m.` + f.GoName + `)`

	} else if f.TypeArray == "" && f.PkgName == "" {
		// Primitive value single

		// string and U16String are special cases because they have custom
		// serialization implementations but still use a non-generated type in
		// generated message fields.
		if f.RosType == "string" {
			m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
			return "primitives.StringAsCStruct(unsafe.Pointer(&mem." + f.CName + "), m." + f.GoName + ")"
		} else if f.RosType == "U16String" {
			m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
			return "primitives.U16StringAsCStruct(unsafe.Pointer(&mem." + f.CName + "), m." + f.GoName + ")"
		}
		return `mem.` + f.CName + ` = C.` + f.CType + `(m.` + f.GoName + `)`
	}
	return "//<MISSING cSerializationCode!!>"
}

func cStructName(f *ROS2Field, m *ROS2Message) string {
	if f.PkgName == "primitives" {
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

func (p *parser) goSerializationCode(f *ROS2Field, m *ROS2Message) string {

	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Array local package reference
		return ucFirst(f.RosType) + `__Array_to_Go(m.` + f.GoName + `[:], mem.` + f.CName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" {
		// Complex value Array remote package reference
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + f.GoPkgReference() + ucFirst(f.RosType) + `__Array_to_Go(m.` + f.GoName + `[:], *(*[]` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)))`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && f.PkgIsLocal {
		// Complex value Slice local package reference
		return ucFirst(f.RosType) + `__Sequence_to_Go(&m.` + f.GoName + `, mem.` + f.CName + `)`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" && !f.PkgIsLocal {
		// Complex value Slice remote package reference
		return f.GoPkgReference() + ucFirst(f.RosType) + `__Sequence_to_Go(&m.` + f.GoName + `, *(*` + f.GoPkgReference() + `C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return f.GoPkgReference() + f.GoType + "TypeSupport.AsGoStruct(&m." + f.GoName + ", unsafe.Pointer(&mem." + f.CName + "))"

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
		return `cSlice_` + f.RosName + ` := mem.` + f.CName + `[:]
	` + `primitives.` + ucFirst(f.RosType) + `__Array_to_Go(m.` + f.GoName + `[:], *(*[]primitives.C` + ucFirst(f.RosType) + `)(unsafe.Pointer(&cSlice_` + f.RosName + `)))`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
		return `primitives.` + ucFirst(f.RosType) + `__Sequence_to_Go(&m.` + f.GoName + `, *(*primitives.C` + ucFirst(f.RosType) + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray == "" && f.PkgName == "" {
		// Primitive value single

		// string and U16String are special cases because they have custom
		// serialization implementations but still use a non-generated type in
		// generated message fields.
		if f.RosType == "string" {
			m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
			return "primitives.StringAsGoStruct(&m." + f.GoName + ", unsafe.Pointer(&mem." + f.CName + "))"
		} else if f.RosType == "U16String" {
			m.GoImports[p.config.RclgoImportPath+"/pkg/rclgo/primitives"] = "primitives"
			return "primitives.U16StringAsGoStruct(&m." + f.GoName + ", unsafe.Pointer(&mem." + f.CName + "))"
		}
		return `m.` + f.GoName + ` = ` + f.GoType + `(mem.` + f.CName + `)`

	}
	return "//<MISSING goSerializationCode!!>"
}

func defaultCode(f *ROS2Field) string {
	if f.PkgName != "" && f.TypeArray != "" && f.DefaultValue == "" {
		// Complex value array and slice common default
		if f.ArraySize == 0 {
			return "t." + f.GoName + " = nil"
		}
		return "for i := range t." + f.GoName + " {\n" +
			"\t\tt." + f.GoName + "[i].SetDefaults()\n" +
			"\t}"

	} else if f.PkgName != "" && f.TypeArray != "" {
		defaultValues := splitMsgDefaultArrayValues(f.RosType, f.DefaultValue)
		// Complex value array and slice
		sb := strings.Builder{}
		var indexesCount int
		if f.ArraySize > 0 {
			indexesCount = f.ArraySize
		} else if len(defaultValues) > 0 { // Init a slice
			indexesCount = len(defaultValues)
			fmt.Fprint(&sb, `t.`, f.GoName, ` = make(`, f.TypeArray, f.GoPkgReference(), f.GoType, `, `, indexesCount, ")\n\t")
		}

		fmt.Fprint(&sb,
			"for i := range t.", f.GoName, " {\n",
			"\t\tt.", f.GoName, "[i].SetDefaults()\n",
			"\t}",
		)
		return sb.String()

	} else if f.PkgName != "" && f.TypeArray == "" {
		// Complex value single
		return `t.` + f.GoName + ".SetDefaults()"

	} else if f.DefaultValue != "" && f.TypeArray != "" {
		// Primitive value array and slice
		defaultValues := splitMsgDefaultArrayValues(f.RosType, f.DefaultValue)
		for i := range defaultValues {
			defaultValues[i] = defaultValueSanitizer(f.RosType, defaultValues[i])
		}
		return `t.` + f.GoName + ` = ` + f.TypeArray + f.GoPkgReference() + f.GoType + `{` + strings.Join(defaultValues, ",") + `}`

	} else if f.DefaultValue == "" && f.TypeArray != "" {
		// Primitive value array and slice common default
		if f.ArraySize == 0 {
			return "t." + f.GoName + " = nil"
		}
		return fmt.Sprint("t.", f.GoName, " = [", f.ArraySize, "]", f.GoType, "{}")

	} else if f.DefaultValue != "" {
		// Primitive value single
		return `t.` + f.GoName + ` = ` + defaultValueSanitizer(f.RosType, f.DefaultValue)

	} else if f.DefaultValue == "" {
		// Primitive value single common default
		return "t." + f.GoName + " = " + primitiveCommonDefault(f)
	}
	return "//<MISSING defaultCode!!>"
}

func primitiveCommonDefault(f *ROS2Field) string {
	switch f.RosType {
	case "string", "wstring", "U16String":
		return `""`
	case "bool":
		return "false"
	case "float32", "float64", "byte", "char", "int8", "int16",
		"int32", "int64", "uint8", "uint16", "uint32", "uint64":
		return "0"
	default:
		panic("common default value for ROS type " + f.RosType + " is not defined")
	}
}

func cloneCode(f *ROS2Field) string {
	if f.PkgName != "" && f.TypeArray != "" && f.ArraySize == 0 {
		// complex value slice
		return "if t." + f.GoName + " != nil {\n" +
			"\t\tc." + f.GoName + " = make([]" + f.GoPkgReference() + f.GoType + ", len(t." + f.GoName + "))\n" +
			"\t\t" + f.GoPkgReference() + "Clone" + f.GoType + "Slice(c." + f.GoName + ", t." + f.GoName + ")\n" +
			"\t}"
	} else if f.PkgName != "" && f.TypeArray != "" && f.ArraySize > 0 {
		// complex value array
		return f.GoPkgReference() + "Clone" + f.GoType + "Slice(c." + f.GoName + "[:], t." + f.GoName + "[:])"
	} else if f.PkgName != "" && f.TypeArray == "" {
		// complex value single
		return "c." + f.GoName + " = *t." + f.GoName + ".Clone()"
	} else if f.PkgName == "" && f.TypeArray != "" && f.ArraySize == 0 {
		// primitive value slice
		return "if t." + f.GoName + " != nil {\n" +
			"\t\tc." + f.GoName + " = make([]" + f.GoType + ", len(t." + f.GoName + "))\n" +
			"\t\tcopy(c." + f.GoName + ", t." + f.GoName + ")\n" +
			"\t}"
	} else if f.PkgName == "" {
		// primitive value single and array
		return "c." + f.GoName + " = t." + f.GoName
	}
	return "//<MISSING cloneCode!!>"
}
