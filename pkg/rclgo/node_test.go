package rclgo_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	. "github.com/smartystreets/goconvey/convey"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

type obj map[string]interface{}

type resultSet []interface{}

func (s *resultSet) Add(description string, val ...interface{}) {
	if description == "" {
		*s = append(*s, val...)
	} else {
		o := obj{"description": description}
		switch len(val) {
		case 0:
		case 1:
			o["value"] = val[0]
		default:
			o["values"] = val
		}
		*s = append(*s, o)
	}
}

func TestNodeGetTopicNamesAndTypes(t *testing.T) {
	var (
		rclctx1, rclctx2 *rclgo.Context
		node1, node2     *rclgo.Node

		results resultSet
		err     error
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
			namesAndTypes, err := node1.GetTopicNamesAndTypes()
			results.Add("node1 in empty network", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})
			namesAndTypes, err = node2.GetTopicNamesAndTypes()
			results.Add("node2 in empty network", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})

			_, err = std_msgs_msg.NewBoolPublisher(node1, "test_topic", nil)
			results.Add("new publisher error", err)
			namesAndTypes, err = node1.GetTopicNamesAndTypes()
			results.Add("node1 after publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})
			namesAndTypes, err = node2.GetTopicNamesAndTypes()
			results.Add("node2 after publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})

			intpub, err := std_msgs_msg.NewInt64Publisher(node1, "test_topic2", nil)
			results.Add("new int publisher error", err)
			namesAndTypes, err = node1.GetTopicNamesAndTypes()
			results.Add("node1 after creating int publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})
			namesAndTypes, err = node2.GetTopicNamesAndTypes()
			results.Add("node2 after creating int publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})

			err = intpub.Publish(std_msgs_msg.NewInt64())
			results.Add("publish int error from node1", err)
			namesAndTypes, err = node1.GetTopicNamesAndTypes()
			results.Add("node1 after publishing int", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})
			namesAndTypes, err = node2.GetTopicNamesAndTypes()
			results.Add("node2 after publishing int", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})

			_, err = std_msgs_msg.NewStringPublisher(node2, "test_topic", nil)
			results.Add("second publisher error", err)
			namesAndTypes, err = node1.GetTopicNamesAndTypes()
			results.Add("node1 after second publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})
			namesAndTypes, err = node2.GetTopicNamesAndTypes()
			results.Add("node2 after second publisher", obj{
				"namesAndTypes": namesAndTypes,
				"err":           err,
			})

			So(cupaloy.Snapshot(results...), ShouldBeNil)
		})
	})
}
