package rclgo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func requireTopicNamesAndTypes(t *testing.T, node *rclgo.Node, expected map[string][]string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	for {
		actual, err := node.GetTopicNamesAndTypes(true)
		require.NoError(t, err)
		if assert.ObjectsAreEqualValues(expected, actual) {
			return
		}
		select {
		case <-ctx.Done():
			require.EqualValues(t, expected, actual)
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func TestNodeGetTopicNamesAndTypes(t *testing.T) {
	setNewDomainID()

	rclctx1, err := newDefaultRCLContext()
	require.NoError(t, err)
	defer rclctx1.Close()
	node1, err := rclctx1.NewNode("node1", "topic_names_and_types_test")
	require.NoError(t, err)

	rclctx2, err := newDefaultRCLContext()
	require.NoError(t, err)
	defer rclctx2.Close()
	node2, err := rclctx2.NewNode("node2", "topic_names_and_types_test")
	require.NoError(t, err)

	t.Log("node1 in empty network")
	requireTopicNamesAndTypes(t, node1, map[string][]string{
		"/rosout": {"rcl_interfaces/msg/Log"},
	})

	t.Log("node2 in empty network")
	requireTopicNamesAndTypes(t, node2, map[string][]string{
		"/rosout": {"rcl_interfaces/msg/Log"},
	})

	t.Log("new publisher")
	_, err = std_msgs_msg.NewBoolPublisher(node1, "test_topic", nil)
	require.NoError(t, err)

	t.Log("node1 after publisher")
	requireTopicNamesAndTypes(t, node1, map[string][]string{
		"/rosout":                                {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic": {"std_msgs/msg/Bool"},
	})

	t.Log("node2 after publisher")
	requireTopicNamesAndTypes(t, node2, map[string][]string{
		"/rosout":                                {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic": {"std_msgs/msg/Bool"},
	})

	t.Log("new int publisher")
	intpub, err := std_msgs_msg.NewInt64Publisher(node1, "test_topic2", nil)
	require.NoError(t, err)

	t.Log("node1 after creating int publisher")
	requireTopicNamesAndTypes(t, node1, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})

	t.Log("node2 after creating int publisher")
	requireTopicNamesAndTypes(t, node2, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})

	t.Log("publish int")
	err = intpub.Publish(std_msgs_msg.NewInt64())
	require.NoError(t, err)

	t.Log("node1 after publishing int")
	requireTopicNamesAndTypes(t, node1, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})

	t.Log("node2 after publishing int")
	requireTopicNamesAndTypes(t, node2, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})

	t.Log("new string publisher")
	_, err = std_msgs_msg.NewStringPublisher(node2, "test_topic", nil)
	require.NoError(t, err)

	t.Log("node1 after second publisher")
	requireTopicNamesAndTypes(t, node1, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool", "std_msgs/msg/String"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})

	t.Log("node2 after second publisher")
	requireTopicNamesAndTypes(t, node2, map[string][]string{
		"/rosout":                                 {"rcl_interfaces/msg/Log"},
		"/topic_names_and_types_test/test_topic":  {"std_msgs/msg/Bool", "std_msgs/msg/String"},
		"/topic_names_and_types_test/test_topic2": {"std_msgs/msg/Int64"},
	})
}
