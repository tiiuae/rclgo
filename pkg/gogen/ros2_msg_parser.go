package gogen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func SnakeToCamel(in string) string {
	tmp := []rune(in)
	tmp[0] = unicode.ToUpper(tmp[0])
	for i := 0; i < len(tmp); i++ {
		if tmp[i] == '_' {
			tmp[i+1] = unicode.ToUpper(tmp[i+1])
			tmp = append(tmp[:i], tmp[i+1:]...)
			i--
		}
	}
	return string(tmp)
}

func CamelToSnake(in string) string {
	tmp := []rune(in)
	tmp[0] = unicode.ToLower(tmp[0])
	for i := 0; i < len(tmp); i++ {
		if unicode.IsUpper(tmp[i]) {
			tmp[i] = unicode.ToLower(tmp[i])
			tmp = append(tmp[:i], append([]rune{'_'}, tmp[i:]...)...)
			i++ // adjust for the extra '_'-character
			i++ // Look if the next character is UC
			for i < len(tmp) && unicode.IsUpper(tmp[i]) {
				tmp[i] = unicode.ToLower(tmp[i])
				i++
			}
		}
	}
	return string(tmp)
}

func CName(rosName string) string {
	switch rosName {
	case "type":
		return "_type"
	default:
		return rosName
	}
}

// ParseROS2Message parses a message definition.
func ParseROS2Message(res *ROS2Message, content string) error {
	for _, line := range strings.Split(content, "\n") {
		obj, err := ParseROS2MessageRow(line, res)
		if err != nil {
			return err
		}

		switch obj.(type) {
		case *ROS2Constant:
			res.Constants = append(res.Constants, obj.(*ROS2Constant))
		case *ROS2Field:
			res.Fields = append(res.Fields, obj.(*ROS2Field))
		case nil:
			continue
		default:
			return fmt.Errorf("Couldn't parse the input row '%s'", line)
		}
	}

	for _, f := range res.Fields {
		switch f.PkgName {
		case "":
		case ".":
		case "time":
			res.GoImports["time"] = ""
		case "rosidl_runtime_c":
			res.GoImports["github.com/tiiuae/rclgo/pkg/ros2/"+f.PkgName] = f.PkgName
		default:
			res.GoImports["github.com/tiiuae/rclgo/pkg/ros2/msgs/"+f.PkgName+"/msg"] = f.PkgName
			res.CImports[f.PkgName] = true
		}
	}

	return nil
}

func ParseROS2MessageRow(testRow string, ros2msg *ROS2Message) (interface{}, error) {
	// remove leading and trailing spaces
	testRow = strings.TrimSpace(testRow)
	// remove multiple spaces
	testRow = regexp.MustCompile(`\s+`).ReplaceAllString(testRow, " ")

	// do not process empty lines or comment lines
	if testRow == "" || regexp.MustCompile(`^#`).MatchString(testRow) {
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
	capture, err := parseNamedCaptureGroupsRegex(textRow, parseROS2Const_regexp)
	if err == nil {
		return 'c', capture
	}
	capture, err = parseNamedCaptureGroupsRegex(textRow, parseROS2Field_regexp)
	if err == nil {
		return 'f', capture
	}
	return 'e', nil
}

var parseROS2Const_regexp = regexp.MustCompile(`^(:?(?P<package>\w+)/)?(?P<type>\w+)(?P<array>\[(?P<arraySize>\d*)\])? (?P<field>\w+) = ?(?P<default>[^#]+)?(:?\s+#\s*(?P<comment>.*))?$`)

func ParseROS2MessageConstant(capture map[string]string, ros2msg *ROS2Message) (*ROS2Constant, error) {
	d := ROS2Constant{
		RosType: capture["type"],
		RosName: capture["field"],
		Value:   capture["default"],
		Comment: capture["comment"],
	}

	_, _, d.GoType = Ros2PrimitiveTypeToGoC(d.RosType)
	return &d, nil
}

var parseROS2Field_regexp *regexp.Regexp = regexp.MustCompile(`^(:?(?P<package>\w+)/)?(?P<type>\w+)(?P<array>\[(?P<arraySize>\d*)\])? (?P<field>\w+) ?(?P<default>[^#]+)?(:?\s+#\s*(?P<comment>.*))?$`)

func ParseROS2MessageField(capture map[string]string, ros2msg *ROS2Message) (*ROS2Field, error) {
	arraySize, err := strconv.ParseInt(capture["arraySize"], 10, 32)
	if err != nil && capture["arraySize"] != "" {
		return nil, err
	}
	f := &ROS2Field{
		Comment:      capture["comment"],
		GoName:       SnakeToCamel(capture["field"]),
		RosName:      capture["field"],
		CName:        CName(capture["field"]),
		RosType:      capture["type"],
		TypeArray:    capture["array"],
		ArraySize:    int(arraySize),
		DefaultValue: capture["default"],
		PkgName:      capture["package"],
	}

	f.PkgName, f.CType, f.GoType = translateROS2Type(f, ros2msg)
	if f.PkgName == "." {
		f.PkgIsLocal = true
	}
	// Prepopulate extra Go imports
	cSerializationCode(f, ros2msg)
	goSerializationCode(f, ros2msg)

	return f, nil
}

func translateROS2Type(f *ROS2Field, m *ROS2Message) (pkgName string, cType string, goType string) {
	// explicit package
	if f.PkgName != "" {
		// type of same package
		if f.PkgName == m.RosPackage {
			return "", f.RosType, f.RosType
		}

		// type of other package
		return f.PkgName, f.RosType, f.RosType
	}

	// implicit package, type of std_msgs
	if m.RosPackage != "std_msgs" {
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

	t, ok := ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING[f.RosType]
	if !ok {
		// These are not actually primitive types, but same-package complex types.
		return ".", f.RosType, f.RosType
	}
	return t.PackageName, t.CType, t.GoType
}

func cSerializationCode(f *ROS2Field, m *ROS2Message) string {
	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" {
		// Complex value Array

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" {
		// Complex value Slice

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return `mem.` + f.CName + ` = *(*C.` + cStructName(f, m) + `)(t.` + f.GoName + `.AsCStruct())`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + strings.Title(f.RosType) + `__Array_to_C(*(*[]rosidl_runtime_c.C` + f.CType + `)(unsafe.Pointer(&mem.` + f.CName + `)), t.` + f.GoName + `[:])`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + strings.Title(f.RosType) + `__Sequence_to_C((*rosidl_runtime_c.Crosidl_runtime_c__` + f.CType + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)), t.` + f.GoName + `)`

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
			return m.RosPackage + "__msg__" + f.CType
		} else {
			return f.PkgName + "__msg__" + f.CType
		}
	}
	return "<MISSING cStructName!!>"
}

func goSerializationCode(f *ROS2Field, m *ROS2Message) string {
	if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName != "" {
		// Complex value Array

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName != "" {
		// Complex value Slice

	} else if f.TypeArray == "" && f.PkgName != "" {
		// Complex value single
		return `t.` + f.GoName + `.AsGoStruct(unsafe.Pointer(&mem.` + f.CName + `))`

	} else if f.TypeArray != "" && f.ArraySize > 0 && f.PkgName == "" {
		// Primitive value Array
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + strings.Title(f.RosType) + `__Array_to_Go(t.` + f.GoName + `[:], *(*[]rosidl_runtime_c.C` + f.CType + `)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray != "" && f.ArraySize == 0 && f.PkgName == "" {
		// Primitive value Slice
		m.GoImports["github.com/tiiuae/rclgo/pkg/ros2/rosidl_runtime_c"] = "rosidl_runtime_c"
		return `rosidl_runtime_c.` + strings.Title(f.RosType) + `__Sequence_to_Go(&t.` + f.GoName + `, *(*rosidl_runtime_c.Crosidl_runtime_c__` + f.CType + `__Sequence)(unsafe.Pointer(&mem.` + f.CName + `)))`

	} else if f.TypeArray == "" && f.PkgName == "" {
		// Primitive value single
		return `t.` + f.GoName + ` = ` + f.GoType + `(mem.` + f.CName + `)`

	}
	return "//<MISSING goSerializationCode!!>"
}

//DEPRECATED
func Ros2PrimitiveTypeToGoC(ros2PrimitiveType string) (rosPackage string, cType string, goType string) {
	switch ros2PrimitiveType {
	case "bool":
		return "", ros2PrimitiveType, ros2PrimitiveType
	case "byte":
		return "", "uint8_t", "byte"
	case "string":
		return "rosidl_runtime_c", "String", "String"
	case "float32":
		return "", "float", "float32"
	case "float64":
		return "", "double", "float64"
	case "int8":
		return "", "int8_t", "int8"
	case "int16":
		return "", "int16_t", "int16"
	case "int32":
		return "", "int32_t", "int32"
	case "int64":
		return "", "int64_t", "int64"
	case "uint8":
		return "", "uint8_t", "uint8"
	case "uint16":
		return "", "uint16_t", "uint16"
	case "uint32":
		return "", "uint32_t", "uint32"
	case "uint64":
		return "", "uint64_t", "uint64"
	case "char":
		return "", ros2PrimitiveType, "byte" // In Golang []byte is converted to string
	case "time", "duration":
		return "time", ros2PrimitiveType, strings.Title(ros2PrimitiveType)
	default:
		// These are not actually primitive types, but same-package complex types.
		return ".", ros2PrimitiveType, ros2PrimitiveType
	}
}

func parseNamedCaptureGroupsRegex(textRow string, regexpStr *regexp.Regexp) (map[string]string, error) {
	subexpNames := regexpStr.SubexpNames()
	namedCaptureGroups := make(map[string]string, len(subexpNames))
	matches := regexpStr.FindAllStringSubmatch(textRow, -1)
	if matches == nil {
		return namedCaptureGroups, fmt.Errorf("Unable to parse text '%s' using regexp '%s'\n", textRow, regexpStr)
	}
	for _, match := range matches {
		for groupIdx, group := range match {
			namedCaptureGroups[subexpNames[groupIdx]] = group
		}
	}
	return namedCaptureGroups, nil
}
