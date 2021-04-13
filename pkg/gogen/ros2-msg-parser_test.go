/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var nilAry []string

func TestParseROS2Field(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	testFunc := func(description string, line string, ros2msg *ROS2Message, f interface{}) {
		testName := ros2msg.RosPackage + "." + ros2msg.RosMsgName + " " + rosName(f)
		if description != "" {
			testName += " : " + description
		}
		Convey(testName, func() {
			m, err := ParseROS2MessageRow(line, ros2msg)
			So(err, ShouldBeNil)
			So(m, ShouldResemble, f)
		})
	}

	Convey("Parse ROS2 Fields", t, func() {
		testFunc("",
			"unique_identifier_msgs/UUID goal_id",
			ROS2MessageNew("action_msgs", "GoalInfo"),
			&ROS2Field{
				RosType:      "UUID",
				GoType:       "UUID",
				CType:        "UUID",
				TypeArray:    "",
				ArrayBounded: "",
				ArraySize:    0,
				DefaultValue: "",
				PkgName:      "unique_identifier_msgs",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "goal_id",
				GoName:       "GoalId",
				CName:        "goal_id",
			},
		)
		testFunc("",
			"string[] full_node_names",
			ROS2MessageNew("composition_interfaces", "ListNodes_Response"),
			&ROS2Field{
				RosType:      "string",
				GoType:       "String",
				CType:        "String",
				TypeArray:    "[]",
				ArrayBounded: "",
				ArraySize:    0,
				DefaultValue: "",
				PkgName:      "rosidl_runtime_c",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "full_node_names",
				GoName:       "FullNodeNames",
				CName:        "full_node_names",
			},
		)
		testFunc("",
			"float64[3] float64_values_default [3.1415, 0.0, -3.1415]",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
			&ROS2Field{
				RosType:      "float64",
				GoType:       "float64",
				CType:        "double",
				TypeArray:    "[3]",
				ArrayBounded: "",
				ArraySize:    3,
				DefaultValue: "[3.1415, 0.0, -3.1415]",
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "float64_values_default",
				GoName:       "Float64ValuesDefault",
				CName:        "float64_values_default",
			},
		)
		testFunc("",
			"BasicTypes[3] basic_types_values",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
			&ROS2Field{
				RosType:      "BasicTypes",
				GoType:       "BasicTypes",
				CType:        "BasicTypes",
				TypeArray:    "[3]",
				ArrayBounded: "",
				ArraySize:    3,
				DefaultValue: "",
				PkgName:      ".",
				PkgIsLocal:   true,
				Comment:      "",
				RosName:      "basic_types_values",
				GoName:       "BasicTypesValues",
				CName:        "basic_types_values",
			},
		)
		testFunc("",
			"bool[3] bool_values_default [false, true, false]",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
			&ROS2Field{
				RosType:      "bool",
				GoType:       "bool",
				CType:        "bool",
				TypeArray:    "[3]",
				ArrayBounded: "",
				ArraySize:    3,
				DefaultValue: "[false, true, false]",
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "bool_values_default",
				GoName:       "BoolValuesDefault",
				CName:        "bool_values_default",
			},
		)
		testFunc("",
			`string[3] string_values_default ["", "max value", "min value"]`,
			ROS2MessageNew("test_msgs", "Arrays_Response"),
			&ROS2Field{
				RosType:      "string",
				GoType:       "String",
				CType:        "String",
				TypeArray:    "[3]",
				ArrayBounded: "",
				ArraySize:    3,
				DefaultValue: `["", "max value", "min value"]`,
				PkgName:      "rosidl_runtime_c",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "string_values_default",
				GoName:       "StringValuesDefault",
				CName:        "string_values_default",
			},
		)
		testFunc("Fields with comments containing '=' do not get identified as Constants",
			`float32 v_ref                      # ADC channel voltage reference, use to calculate LSB voltage(lsb=scale/resolution)`,
			ROS2MessageNew("px4_msgs", "AdcReport"),
			&ROS2Field{
				RosType:      "float32",
				GoType:       "float32",
				CType:        "float",
				TypeArray:    "",
				ArrayBounded: "",
				ArraySize:    0,
				DefaultValue: ``,
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "ADC channel voltage reference, use to calculate LSB voltage(lsb=scale/resolution)",
				RosName:      "v_ref",
				GoName:       "VRef",
				CName:        "v_ref",
			},
		)
		testFunc("Array size int is big enough to store the C array size.",
			`uint8[512] junk`,
			ROS2MessageNew("px4_msgs", "OrbTestLarge"),
			&ROS2Field{
				RosType:      "uint8",
				GoType:       "uint8",
				CType:        "uint8_t",
				TypeArray:    "[512]",
				ArrayBounded: "",
				ArraySize:    512,
				DefaultValue: ``,
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "junk",
				GoName:       "Junk",
				CName:        "junk",
			},
		)
		testFunc(`C struct has reserved Go keyword .test. Accessed with ._test instead.`,
			`uint8 type`,
			ROS2MessageNew("sensor_msgs", "JoyFeedback"),
			&ROS2Field{
				RosType:      "uint8",
				GoType:       "uint8",
				CType:        "uint8_t",
				TypeArray:    "",
				ArrayBounded: "",
				ArraySize:    0,
				DefaultValue: ``,
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "type",
				GoName:       "Type",
				CName:        "_type",
			},
		)
		testFunc(`Bounded sequence.`,
			`BasicTypes[<=3] basic_types_values`,
			ROS2MessageNew("test_msgs", "BoundedSequences"),
			&ROS2Field{
				RosType:      "BasicTypes",
				GoType:       "BasicTypes",
				CType:        "BasicTypes",
				TypeArray:    "[]",
				ArrayBounded: "<=3",
				ArraySize:    0,
				DefaultValue: ``,
				PkgName:      ".",
				PkgIsLocal:   true,
				Comment:      "",
				RosName:      "basic_types_values",
				GoName:       "BasicTypesValues",
				CName:        "basic_types_values",
			},
		)
		testFunc(`Bounded sequence with defaults.`,
			`int8[<=3] int8_values_default [0, 127, -128]`,
			ROS2MessageNew("test_msgs", "BoundedSequences"),
			&ROS2Field{
				RosType:      "int8",
				GoType:       "int8",
				CType:        "int8_t",
				TypeArray:    "[]",
				ArrayBounded: "<=3",
				ArraySize:    0,
				DefaultValue: `[0, 127, -128]`,
				PkgName:      "",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "int8_values_default",
				GoName:       "Int8ValuesDefault",
				CName:        "int8_values_default",
			},
		)
		testFunc(`Bounded string with defaults.`,
			`string<=22 bounded_string_value "this is yet another"`,
			ROS2MessageNew("test_msgs", "Strings"),
			&ROS2Field{
				RosType:      "string",
				GoType:       "String",
				CType:        "String",
				TypeArray:    "",
				ArrayBounded: "<=22",
				ArraySize:    0,
				DefaultValue: `"this is yet another"`,
				PkgName:      "rosidl_runtime_c",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "bounded_string_value",
				GoName:       "BoundedStringValue",
				CName:        "bounded_string_value",
			},
		)
		testFunc(`Bounded string.`,
			`string<=22 bounded_string_value`,
			ROS2MessageNew("test_msgs", "Strings"),
			&ROS2Field{
				RosType:      "string",
				GoType:       "String",
				CType:        "String",
				TypeArray:    "",
				ArrayBounded: "<=22",
				ArraySize:    0,
				DefaultValue: ``,
				PkgName:      "rosidl_runtime_c",
				PkgIsLocal:   false,
				Comment:      "",
				RosName:      "bounded_string_value",
				GoName:       "BoundedStringValue",
				CName:        "bounded_string_value",
			},
		)
	})

	Convey("Parse ROS2 Constants", t, func() {
		testFunc("",
			"uint8 NUM_ACTUATOR_CONTROLS = 8",
			ROS2MessageNew("px4_msgs", "ActuatorControls0"),
			&ROS2Constant{
				Value:   "8",
				RosType: "uint8",
				GoType:  "uint8",
				Comment: "",
				RosName: "NUM_ACTUATOR_CONTROLS",
			},
		)
		testFunc("",
			"uint8 TYPE_LED    = 0",
			ROS2MessageNew("sensor_msgs", "JoyFeedback"),
			&ROS2Constant{
				Value:   "0",
				RosType: "uint8",
				GoType:  "uint8",
				Comment: "",
				RosName: "TYPE_LED",
			},
		)
		testFunc("",
			"uint8 BATTERY_WARNING_CRITICAL = 2        # critical voltage, return / abort immediately",
			ROS2MessageNew("px4_msgs", "BatteryStatus"),
			&ROS2Constant{
				Value:   "2",
				RosType: "uint8",
				GoType:  "uint8",
				Comment: "critical voltage, return / abort immediately",
				RosName: "BATTERY_WARNING_CRITICAL",
			},
		)
		testFunc("",
			"byte BYTE_CONST=50",
			ROS2MessageNew("test_msgs", "Constants"),
			&ROS2Constant{
				Value:   "50",
				RosType: "byte",
				GoType:  "byte",
				Comment: "",
				RosName: "BYTE_CONST",
			},
		)
	})

	/*
		ROS2 snake-case-camel-case -conversions use different rules than
		"github.com/stoewer/go-strcase"
		or
		"github.com/iancoleman/strcase"
	*/
	Convey("Case transformations", t, func() {
		So(camelToSnake("ColorRGBA"), ShouldEqual, "color_rgba")
		So(camelToSnake("MultiDOFJointTrajectoryPoint"), ShouldEqual, "multi_dof_joint_trajectory_point")
		So(camelToSnake("TFMessage"), ShouldEqual, "tf_message")
		So(camelToSnake("WStrings"), ShouldEqual, "w_strings")
		So(camelToSnake("Float32MultiArray"), ShouldEqual, "float32_multi_array")
		So(camelToSnake("PointCloud2"), ShouldEqual, "point_cloud2")
		So(camelToSnake("GoalID"), ShouldEqual, "goal_id")
		So(camelToSnake("WString"), ShouldEqual, "w_string")
		So(camelToSnake("TF2Error"), ShouldEqual, "tf2_error")
	})

	Convey("Defaults string parser", t, func() {
		So(splitMsgDefaultArrayValues("int", ``), ShouldResemble, nilAry)
		So(splitMsgDefaultArrayValues("int", `[]`), ShouldResemble, nilAry)
		So(splitMsgDefaultArrayValues("int", `[1,2,3]`), ShouldResemble, []string{`1`, `2`, `3`})
		So(splitMsgDefaultArrayValues("string", `["", "this is a", "test msg"]`), ShouldResemble, []string{`""`, `"this is a"`, `"test msg"`})
		So(splitMsgDefaultArrayValues("string", `[1  ,  2 ,   "3"]`), ShouldResemble, []string{`"1  "`, `"2 "`, `"3"`})
		So(defaultValueSanitizer("string", `"Hello world!"`), ShouldEqual, `"Hello world!"`)
		So(defaultValueSanitizer("string", `"Hello\"world!"`), ShouldEqual, `"Hello\"world!"`)
	})

	Convey("defaultCode() generator", t, func() {
		So(defaultCode(&ROS2Field{
			TypeArray:    "[3]",
			ArraySize:    3,
			DefaultValue: "",
			PkgName:      "StringValues",
			PkgIsLocal:   false,
			RosType:      "string",
			CType:        "String",
			GoType:       "String",
			RosName:      "string_values",
			GoName:       "StringValues",
			CName:        "string_values",
			Comment:      "",
		}), ShouldEqual, `t.StringValues[0].SetDefaults(nil)
	t.StringValues[1].SetDefaults(nil)
	t.StringValues[2].SetDefaults(nil)
	`)
	})
}

func TestCErrorTypeParser(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	Convey("", t, func() {
		et, err := ParseROS2ErrorType("/// Success return code.")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.String(), ShouldEqual, "Success return code.")

		et, err = ParseROS2ErrorType("#define RCL_RET_OK RMW_RET_OK")
		So(et, ShouldResemble, &ROS2ErrorType{
			Name:      "RCL_RET_OK",
			Rcl_ret_t: "",
			Reference: "RMW_RET_OK",
			Comment:   "Success return code.",
		})
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.Len(), ShouldEqual, 0)

		et, err = ParseROS2ErrorType("/// This comment is flushed because it is not part of a continuous stream.")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.Len(), ShouldBeGreaterThan, 0)

		et, err = ParseROS2ErrorType("")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.String(), ShouldEqual, "")

		et, err = ParseROS2ErrorType("#define RCL_RET_NOT_INIT 101")
		So(et, ShouldResemble, &ROS2ErrorType{
			Name:      "RCL_RET_NOT_INIT",
			Rcl_ret_t: "101",
			Reference: "",
			Comment:   "",
		})
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.Len(), ShouldEqual, 0)
	})
}

/*
func TestSerDesSimple(t *testing.T) {

	Convey("example_interfaces/Float32MultiArray.msg", t, func() {
		ob := test_msgs.Arrays{}
	})
}*/

func rosName(obj interface{}) string {
	switch obj.(type) {
	case *ROS2Constant:
		return obj.(*ROS2Constant).RosName
	case *ROS2Field:
		return obj.(*ROS2Field).RosName
	default:
		panic(fmt.Sprintf("Unable to get the ROS2 Message row-object (Field|Constant) name!%+v\n", obj))
	}
}

func TestBlacklist(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("Blacklist", t, func() {
		skip, blacklistEntry := blacklisted("/opt/ros/foxy/this-is-a-test-blacklist-entry-do-not-remove-used-for-internal-testing/msgs/Lol.msg")
		So(skip, ShouldBeTrue)
		So(blacklistEntry, ShouldEqual, "this-is-a-test-blacklist-entry-do-not-remove-used-for-internal-testing")
	})
}
