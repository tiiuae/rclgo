package rclgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func TestNodeGetTopicNamesAndTypes(t *testing.T) {
	setNewDomainID()
	var (
		rclctx1, rclctx2 *rclgo.Context
		node1, node2     *rclgo.Node
		intpub           *std_msgs_msg.Int64Publisher
		err              error
	)
	defer func() {
		if rclctx1 != nil {
			rclctx1.Close()
		}
		if rclctx2 != nil {
			rclctx2.Close()
		}
	}()
	Convey("Scenario: Node.GetTopicNamesAndTypes works correctly", t, func() {
		Convey("Create a rcl context and node", func() {
			rclctx1, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			node1, err = rclctx1.NewNode("node1", "topic_names_and_types_test")
			So(err, ShouldBeNil)
			rclctx2, err = newDefaultRCLContext()
			So(err, ShouldBeNil)
			node2, err = rclctx2.NewNode("node2", "topic_names_and_types_test")
			So(err, ShouldBeNil)
		})
		Convey("Check that topic names and types are correct", func() {
			Convey("node1 in empty network", func() {
				namesAndTypes, err := node1.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout": {"rcl_interfaces/msg/Log"},
				})
			})
			Convey("node2 in empty network", func() {
				namesAndTypes, err := node2.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout": {"rcl_interfaces/msg/Log"},
				})
			})
			Convey("new publisher", func() {
				_, err = std_msgs_msg.NewBoolPublisher(node1, "test_topic", nil)
				So(err, ShouldBeNil)
			})
			Convey("node1 after publisher", func() {
				namesAndTypes, err := node1.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout":                                {"rcl_interfaces/msg/Log"},
					"/topic_names_and_types_test/test_topic": {"std_msgs/msg/Bool"},
				})
			})
			Convey("node2 after publisher", func() {
				namesAndTypes, err := node2.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout": {"rcl_interfaces/msg/Log"},
				})
			})
			Convey("new int publisher", func() {
				intpub, err = std_msgs_msg.NewInt64Publisher(node1, "test_topic2", nil)
				So(err, ShouldBeNil)
			})
			Convey("node1 after creating int publisher", func() {
				namesAndTypes, err := node1.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout":                                 {"rcl_interfaces/msg/Log"},
					"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
					"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
				})
			})
			Convey("node2 after creating int publisher", func() {
				namesAndTypes, err := node2.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout": {"rcl_interfaces/msg/Log"},
				})
			})
			Convey("publish int", func() {
				err = intpub.Publish(std_msgs_msg.NewInt64())
				So(err, ShouldBeNil)
			})
			Convey("node1 after publishing int", func() {
				namesAndTypes, err := node1.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout":                                 {"rcl_interfaces/msg/Log"},
					"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
					"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
				})
			})
			Convey("node2 after publishing int", func() {
				namesAndTypes, err := node2.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout": {"rcl_interfaces/msg/Log"},
				})
			})
			Convey("new string publisher", func() {
				_, err = std_msgs_msg.NewStringPublisher(node2, "test_topic", nil)
				So(err, ShouldBeNil)
			})
			Convey("node1 after second publisher", func() {
				namesAndTypes, err := node1.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout":                                 {"rcl_interfaces/msg/Log"},
					"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool", "std_msgs/msg/String"},
					"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
				})
			})
			Convey("node2 after second publisher", func() {
				namesAndTypes, err := node2.GetTopicNamesAndTypes()
				So(err, ShouldBeNil)
				So(namesAndTypes, ShouldResemble, map[string][]string{
					"/rosout":                                {"rcl_interfaces/msg/Log"},
					"/topic_names_and_types_test/test_topic": {"std_msgs/msg/String"},
				})
			})
		})
	})
}
