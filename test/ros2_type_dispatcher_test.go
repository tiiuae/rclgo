// This test is not in the package ros2_type_dispatcher to avoid an import
// cycle.

package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	std_srvs_srv "github.com/tiiuae/rclgo/pkg/ros2/msgs/std_srvs/srv"
	"github.com/tiiuae/rclgo/pkg/ros2/ros2_type_dispatcher"
)

func TestTranslateROS2ServiceNameToType(t *testing.T) {
	Convey("Scenario: test service name to type translation", t, func() {
		Convey("Translating the name of an imported message should work", func() {
			srv, ok := ros2_type_dispatcher.TranslateROS2ServiceTypeNameToType("std_srvs/Empty")
			So(ok, ShouldBeTrue)
			So(srv, ShouldNotBeNil)
			So(srv, ShouldHaveSameTypeAs, std_srvs_srv.Empty)
		})
		Convey("Translating the name of a non-imported message should not work", func() {
			srv, ok := ros2_type_dispatcher.TranslateROS2ServiceTypeNameToType("std_srvs/ThisServiceDoesNotExist")
			So(ok, ShouldBeFalse)
			So(srv, ShouldBeNil)
		})
	})
}
