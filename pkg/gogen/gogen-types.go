package gogen

// ROS2Message is a message definition.
type ROS2Message struct {
	RosMsgName string
	RosPackage string
	Url        string
	Fields     []ROS2Field
	Constants  []ROS2Constant
	GoImports  map[string]struct{}
	CImports   map[string]bool
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
	TypeArray  string
	PkgName    string
	PkgIsLocal bool
	RosType    string
	CType      string
	GoType     string
	RosName    string
	CName      string
	GoName     string
	Comment    string
}
