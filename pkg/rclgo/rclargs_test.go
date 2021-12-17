package rclgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func TestRCLArgs(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	Convey("RCLArgs parsing", t, func() {
		args, err := rclgo.NewRCLArgs(
			"--ros-args --log-level DEBUG --enclave /enclave",
		)
		So(err, ShouldBeNil)
		So(
			args.GoArgs,
			ShouldResemble,
			[]string{
				"--ros-args",
				"--log-level",
				"DEBUG",
				"--enclave",
				"/enclave",
			},
		)
	})
}
