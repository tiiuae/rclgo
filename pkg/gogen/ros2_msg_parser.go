package gogen

import (
	"fmt"
	"regexp"
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

			for unicode.IsUpper(tmp[i] + 1) {
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
		// remove comments
		commentsRegex := regexp.MustCompile(`#\s*(.*)$`)
		reSubmatch := commentsRegex.FindStringSubmatch(line)
		var comment string = ""
		if len(reSubmatch) > 0 {
			comment = reSubmatch[1]
		}
		line = commentsRegex.ReplaceAllString(line, "")

		// remove leading and trailing spaces
		line = strings.TrimSpace(line)

		// do not process empty lines
		if line == "" {
			continue
		}

		// definition
		if strings.Contains(line, "=") {
			matches := regexp.MustCompile(`^([a-z0-9]+)(\s|\t)+([A-Z0-9_]+)(\s|\t)*=(\s|\t)*(.+?)$`).FindStringSubmatch(line)
			if matches == nil {
				return fmt.Errorf("unable to parse definition (%s)\n", line)
			}

			d := ROS2Constant{
				RosType: matches[1],
				RosName: matches[3],
				Value:   matches[6],
				Comment: comment,
			}

			_, _, d.GoType = Ros2PrimitiveTypeToGoC(d.RosType)

			res.Constants = append(res.Constants, d)

			// field
		} else {
			// remove multiple spaces between type and name
			line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")

			parts := strings.Split(line, " ")
			if len(parts) != 2 {
				return fmt.Errorf("unable to parse field (%s)\n", line)
			}

			f := ROS2Field{Comment: comment}

			// use NameOverride if a bidirectional conversion between snake and
			// camel is not possible
			f.GoName = SnakeToCamel(parts[1])
			f.RosName = parts[1]
			f.CName = CName(f.RosName)

			f.RosType = parts[0]

			// split TypeArray and Type
			ma := regexp.MustCompile(`^(.+?)(\[.*?\])$`).FindStringSubmatch(f.RosType)
			if ma != nil {
				f.TypeArray = ma[2]
				f.RosType = ma[1]
			}

			f.PkgName, f.CType, f.GoType = func() (string, string, string) {
				// explicit package
				parts := strings.Split(f.RosType, "/")
				if len(parts) == 2 {
					// type of same package
					if parts[0] == res.RosPackage {
						return "", f.RosType, parts[1]
					}

					// type of other package
					return parts[0], f.RosType, parts[1]
				}

				// implicit package, type of std_msgs
				if res.RosPackage != "std_msgs" {
					switch f.RosType {
					case "Bool", "ColorRGBA",
						"Duration", "Empty", "Float32MultiArray", "Float32",
						"Float64MultiArray", "Float64", "Header", "Int8MultiArray",
						"Int8", "Int16MultiArray", "Int16", "Int32MultiArray", "Int32",
						"Int64MultiArray", "Int64", "MultiArrayDimension", "MultiarrayLayout",
						"String", "Time", "UInt8MultiArray", "UInt8", "UInt16MultiArray", "UInt16",
						"UInt32MultiArray", "UInt32", "UInt64MultiArray", "UInt64":
						return "std_msgs", f.RosType, parts[0]
					}
				}

				return Ros2PrimitiveTypeToGoC(f.RosType)
			}()
			if f.PkgName == "." {
				f.PkgIsLocal = true
			}

			res.Fields = append(res.Fields, f)
		}
	}

	res.GoImports = map[string]struct{}{}
	res.CImports = map[string]bool{}
	for _, f := range res.Fields {
		switch f.PkgName {
		case "":
		case ".":
		case "time":
			res.GoImports["time"] = struct{}{}
		default:
			res.GoImports["github.com/tiiuae/rclgo/pkg/ros2/msgs/"+f.PkgName] = struct{}{}
			res.CImports[f.PkgName] = true
		}
	}

	return nil
}

func Ros2PrimitiveTypeToGoC(ros2PrimitiveType string) (rosPackage string, cType string, goType string) {
	switch ros2PrimitiveType {
	case "bool", "byte",
		"string":
		return "", ros2PrimitiveType, ros2PrimitiveType
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
