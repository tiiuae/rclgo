package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tiiuae/rclgo/pkg/ros2"
)

func TestRCLArgs(t *testing.T) {

	SetDefaultFailureMode(FailureContinues)
	Convey("RCLArgs parsing", t, func() {
		args, err := ros2.NewRCLArgs("--ros-args --log-level DEBUG --enclave /enclave")
		So(err, ShouldBeNil)
		So(args.GoArgs, ShouldResemble, []string{"--ros-args", "--log-level", "DEBUG", "--enclave", "/enclave"})
	})
}
