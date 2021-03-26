package gogen

/*
ROS2Message is a message definition. https://design.ros2.org/articles/legacy_interface_definition.html
Use ROS2MessageNew() to initialize the struct
*/
type ROS2Message struct {
	RosMsgName        string
	RosPackage        string
	DataStructureType string // The type, either msg, srv, ...
	Url               string
	Fields            []*ROS2Field
	Constants         []*ROS2Constant
	GoImports         map[string]string
	CImports          map[string]bool
}

func ROS2MessageNew(RosPackage string, RosMsgName string) *ROS2Message {
	return &ROS2Message{
		RosPackage: RosPackage,
		RosMsgName: RosMsgName,
		GoImports:  map[string]string{},
		CImports:   map[string]bool{},
	}
}

// ROS2Constant is a message definition.
type ROS2Constant struct {
	RosType string
	GoType  string
	RosName string
	Value   string
	Comment string
}

// Field is a message field.
type ROS2Field struct {
	TypeArray    string
	ArraySize    int
	DefaultValue string
	PkgName      string
	PkgIsLocal   bool
	RosType      string
	CType        string
	GoType       string
	RosName      string
	CName        string
	GoName       string
	Comment      string
}

type rosidl_runtime_c_type_mapping struct {
	RosType     string
	GoType      string
	CType       string
	CStructName string
	PackageName string
	SkipAutogen bool
}

var ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING = map[string]rosidl_runtime_c_type_mapping{
	"string":   {RosType: "string", GoType: "String", CStructName: "String", CType: "String", PackageName: "rosidl_runtime_c", SkipAutogen: true},
	"time":     {RosType: "time", GoType: "Time", CStructName: "Time", CType: "time", SkipAutogen: true},
	"duration": {RosType: "duration", GoType: "Duration", CStructName: "Duration", CType: "duration", SkipAutogen: true},
	"float32":  {RosType: "float32", GoType: "float32", CStructName: "float", CType: "float"},
	"float64":  {RosType: "float64", GoType: "float64", CStructName: "double", CType: "double"},
	"bool":     {RosType: "bool", GoType: "bool", CStructName: "boolean", CType: "bool"},
	"byte":     {RosType: "byte", GoType: "byte", CStructName: "octet", CType: "uint8_t"},
	"char":     {RosType: "char", GoType: "byte", CStructName: "char", CType: "char"},
	"int8":     {RosType: "int8", GoType: "int8", CStructName: "int8", CType: "int8_t"},
	"int16":    {RosType: "int16", GoType: "int16", CStructName: "int16", CType: "int16_t"},
	"int32":    {RosType: "int32", GoType: "int32", CStructName: "int32", CType: "int32_t"},
	"int64":    {RosType: "int64", GoType: "int64", CStructName: "int64", CType: "int64_t"},
	"uint8":    {RosType: "uint8", GoType: "uint8", CStructName: "uint8", CType: "uint8_t"},
	"uint16":   {RosType: "uint16", GoType: "uint16", CStructName: "uint16", CType: "uint16_t"},
	"uint32":   {RosType: "uint32", GoType: "uint32", CStructName: "uint32", CType: "uint32_t"},
	"uint64":   {RosType: "uint64", GoType: "uint64", CStructName: "uint64", CType: "uint64_t"},
}
