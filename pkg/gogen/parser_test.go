/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	. "github.com/smartystreets/goconvey/convey"
)

var nilAry []string

func TestParseROS2Field(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	parser := parser{config: &DefaultConfig}

	testFunc := func(description string, line string, ros2msg *ROS2Message) {
		testName := ros2msg.Package + "." + ros2msg.Name + " " + line
		if description != "" {
			testName += " : " + description
		}
		Convey(testName, func() {
			m, err := parser.parseMessageLine(line, ros2msg)
			So(err, ShouldBeNil)
			sum := md5.Sum([]byte(line))
			So(cupaloy.SnapshotMulti(hex.EncodeToString(sum[:]), m), ShouldBeNil)
		})
	}

	Convey("Parse ROS2 Fields", t, func() {
		testFunc("",
			"unique_identifier_msgs/UUID goal_id",
			ROS2MessageNew("action_msgs", "GoalInfo"),
		)
		testFunc("",
			"string[] full_node_names",
			ROS2MessageNew("composition_interfaces", "ListNodes_Response"),
		)
		testFunc("",
			"float64[3] float64_values_default [3.1415, 0.0, -3.1415]",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
		)
		testFunc("",
			"BasicTypes[3] basic_types_values",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
		)
		testFunc("",
			"bool[3] bool_values_default [false, true, false]",
			ROS2MessageNew("test_msgs", "Arrays_Response"),
		)
		testFunc("",
			`string[3] string_values_default ["", "max value", "min value"]`,
			ROS2MessageNew("test_msgs", "Arrays_Response"),
		)
		testFunc("Fields with comments containing '=' do not get identified as Constants",
			`float32 v_ref                      # ADC channel voltage reference, use to calculate LSB voltage(lsb=scale/resolution)`,
			ROS2MessageNew("px4_msgs", "AdcReport"),
		)
		testFunc("Array size int is big enough to store the C array size.",
			`uint8[512] junk`,
			ROS2MessageNew("px4_msgs", "OrbTestLarge"),
		)
		testFunc(`C struct has reserved Go keyword .test. Accessed with ._test instead.`,
			`uint8 type`,
			ROS2MessageNew("sensor_msgs", "JoyFeedback"),
		)
		testFunc(`Bounded sequence.`,
			`BasicTypes[<=3] basic_types_values`,
			ROS2MessageNew("test_msgs", "BoundedSequences"),
		)
		testFunc(`Bounded sequence with defaults.`,
			`int8[<=3] int8_values_default [0, 127, -128]`,
			ROS2MessageNew("test_msgs", "BoundedSequences"),
		)
		testFunc(`Bounded string with defaults.`,
			`string<=22 bounded_string_value "this is yet another"`,
			ROS2MessageNew("test_msgs", "Strings"),
		)
		testFunc(`Bounded string.`,
			`string<=22 bounded_string_value`,
			ROS2MessageNew("test_msgs", "Strings"),
		)
	})

	testParseService := func(pkg, name, source string) {
		s := NewROS2Service(pkg, name)
		So(parser.parseService(s, source), ShouldBeNil)
		So(cupaloy.SnapshotMulti("service-"+name, s), ShouldBeNil)
	}

	Convey("Parse ROS2 services", t, func() {
		testParseService("action_msgs", "CancelGoal", `
# Cancel one or more goals with the following policy:
#
# - If the goal ID is zero and timestamp is zero, cancel all goals.
# - If the goal ID is zero and timestamp is not zero, cancel all goals accepted
#   at or before the timestamp.
# - If the goal ID is not zero and timestamp is zero, cancel the goal with the
#   given ID regardless of the time it was accepted.
# - If the goal ID is not zero and timestamp is not zero, cancel the goal with
#   the given ID and all goals accepted at or before the timestamp.

# Goal info describing the goals to cancel, see above.
GoalInfo goal_info
---
##
## Return codes.
##

# Indicates the request was accepted without any errors.
#
# One or more goals have transitioned to the CANCELING state. The
# goals_canceling list is not empty.
int8 ERROR_NONE=0

# Indicates the request was rejected.
#
# No goals have transitioned to the CANCELING state. The goals_canceling list is
# empty.
int8 ERROR_REJECTED=1

# Indicates the requested goal ID does not exist.
#
# No goals have transitioned to the CANCELING state. The goals_canceling list is
# empty.
int8 ERROR_UNKNOWN_GOAL_ID=2

# Indicates the goal is not cancelable because it is already in a terminal state.
#
# No goals have transitioned to the CANCELING state. The goals_canceling list is
# empty.
int8 ERROR_GOAL_TERMINATED=3

# Return code, see above definitions.
int8 return_code

# Goals that accepted the cancel request.
GoalInfo[] goals_canceling		
`)
		testParseService("tf2_msgs", "FrameGraph", `
---
string frame_yaml
`)
		testParseService("", "NoResponse", `
string input
---
`)
	})

	Convey("Parse ROS2 Constants", t, func() {
		testFunc("",
			"uint8 NUM_ACTUATOR_CONTROLS = 8",
			ROS2MessageNew("px4_msgs", "ActuatorControls0"),
		)
		testFunc("",
			"uint8 TYPE_LED    = 0",
			ROS2MessageNew("sensor_msgs", "JoyFeedback"),
		)
		testFunc("",
			"uint8 BATTERY_WARNING_CRITICAL = 2        # critical voltage, return / abort immediately",
			ROS2MessageNew("px4_msgs", "BatteryStatus"),
		)
		testFunc("",
			"byte BYTE_CONST=50",
			ROS2MessageNew("test_msgs", "Constants"),
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
		}), ShouldEqual, `for i := range t.StringValues {
		t.StringValues[i].SetDefaults()
	}`)
	})
}

func TestCErrorTypeParser(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	Convey("", t, func() {
		et, err := parseROS2ErrorType("/// Success return code.")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.String(), ShouldEqual, "Success return code.")

		et, err = parseROS2ErrorType("#define RCL_RET_OK RMW_RET_OK")
		So(et, ShouldResemble, &ROS2ErrorType{
			Name:      "RCL_RET_OK",
			Rcl_ret_t: "",
			Reference: "RMW_RET_OK",
			Comment:   "Success return code.",
		})
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.Len(), ShouldEqual, 0)

		et, err = parseROS2ErrorType("/// This comment is flushed because it is not part of a continuous stream.")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.Len(), ShouldBeGreaterThan, 0)

		et, err = parseROS2ErrorType("")
		So(et, ShouldBeNil)
		So(err, ShouldBeNil)
		So(ros2errorTypesCommentsBuffer.String(), ShouldEqual, "")

		et, err = parseROS2ErrorType("#define RCL_RET_NOT_INIT 101")
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

func TestBlacklist(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)

	Convey("Blacklist", t, func() {
		skip, blacklistEntry := blacklisted("/opt/ros/galactic/this-is-a-test-blacklist-entry-do-not-remove-used-for-internal-testing/msgs/Lol.msg")
		So(skip, ShouldBeTrue)
		So(blacklistEntry, ShouldEqual, "this-is-a-test-blacklist-entry-do-not-remove-used-for-internal-testing")
	})
}
