package ros2

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRCLArgs(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	Convey("RCLArgs parsing", t, func() {
		args, err := NewRCLArgs("--ros-args --log-level DEBUG --enclave /enclave")
		So(err, ShouldBeNil)
		So(args.GoArgs, ShouldResemble, []string{"--ros-args", "--log-level", "DEBUG", "--enclave", "/enclave"})
	})
}
