//#nosec G404

package rclgo_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	std_msgs_msg "github.com/tiiuae/rclgo/internal/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

var (
	usedDomainIDs      = map[int]bool{}
	usedDomainIDsMutex = sync.Mutex{}
)

func getUnusedDomainID() int {
	usedDomainIDsMutex.Lock()
	defer usedDomainIDsMutex.Unlock()
	if len(usedDomainIDs) == 102 {
		panic("ran out of unused domain IDs")
	}
	oldID, _ := strconv.Atoi(os.Getenv("ROS_DOMAIN_ID"))
	usedDomainIDs[oldID] = true
	newID := rand.Intn(102)
	for usedDomainIDs[newID] {
		newID = rand.Intn(102)
	}
	usedDomainIDs[newID] = true
	return newID
}

func setNewDomainID() {
	os.Setenv("ROS_DOMAIN_ID", strconv.Itoa(getUnusedDomainID()))
}

func TestMain(m *testing.M) {
	setNewDomainID()
	os.Exit(m.Run())
}

func ExampleParseArgs() {
	rosArgs, restArgs, err := rclgo.ParseArgs(
		[]string{
			"--extra0",
			"args0",
			"--ros-args",
			"--log-level",
			"DEBUG",
			"--",
			"--extra1",
			"args1",
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = rclgo.ParseArgs(
		[]string{"--ros-args", "--log-level", "INFO"},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = rclgo.ParseArgs(
		[]string{"--extra0", "args0", "--extra1", "args1"},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = rclgo.ParseArgs(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n", restArgs)

	// Output: rosArgs: [--log-level DEBUG]
	// restArgs: [--extra0 args0 --extra1 args1]
	//
	// rosArgs: [--log-level INFO]
	// restArgs: []
	//
	// rosArgs: []
	// restArgs: [--extra0 args0 --extra1 args1]
	//
	// rosArgs: []
	// restArgs: []
}

func requireValue[T any](t *testing.T, expected T) func(actual T, err error) {
	t.Helper()
	return func(actual T, err error) {
		t.Helper()
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}

func noErr[T any](v T, err error) func(*testing.T) T {
	return func(t *testing.T) T {
		t.Helper()
		require.NoError(t, err)
		return v
	}
}

func noopSub(*rclgo.Subscription) {}

func TestPubSubCount(t *testing.T) {
	setNewDomainID()

	rclctx := noErr(newDefaultRCLContext())(t)
	defer rclctx.Close()

	pubNode := noErr(rclctx.NewNode("node1", "test_pubsub_counts"))(t)
	subNode := noErr(rclctx.NewNode("node2", "test_pubsub_counts"))(t)

	pub := noErr(pubNode.NewPublisher("string", std_msgs_msg.StringTypeSupport, nil))(t)

	requireValue(t, 0)(pub.GetSubscriptionCount())

	sub1 := noErr(subNode.NewSubscription("string", std_msgs_msg.StringTypeSupport, nil, noopSub))(t)

	requireValue(t, 1)(pub.GetSubscriptionCount())
	requireValue(t, 1)(sub1.GetPublisherCount())

	sub2 := noErr(subNode.NewSubscription("string", std_msgs_msg.StringTypeSupport, nil, noopSub))(t)

	requireValue(t, 2)(pub.GetSubscriptionCount())
	requireValue(t, 1)(sub1.GetPublisherCount())
	requireValue(t, 1)(sub2.GetPublisherCount())

	subOptsUnrealiable := rclgo.NewDefaultSubscriptionOptions()
	subOptsUnrealiable.Qos.Reliability = rclgo.ReliabilityBestEffort
	subOptsUnrealiable.Qos.Durability = rclgo.DurabilityTransientLocal
	subUnreliable := noErr(subNode.NewSubscription("string", std_msgs_msg.StringTypeSupport, subOptsUnrealiable, noopSub))(t)

	requireValue(t, 2)(pub.GetSubscriptionCount())
	requireValue(t, 1)(sub1.GetPublisherCount())
	requireValue(t, 1)(sub2.GetPublisherCount())
	requireValue(t, 0)(subUnreliable.GetPublisherCount())

	pubOptsUnrealiable := rclgo.NewDefaultPublisherOptions()
	pubOptsUnrealiable.Qos.Reliability = rclgo.ReliabilityBestEffort
	pubOptsUnrealiable.Qos.Durability = rclgo.DurabilityTransientLocal
	pubUnreliable := noErr(pubNode.NewPublisher("string", std_msgs_msg.StringTypeSupport, pubOptsUnrealiable))(t)

	requireValue(t, 2)(pub.GetSubscriptionCount())
	requireValue(t, 1)(pubUnreliable.GetSubscriptionCount())
	requireValue(t, 1)(sub1.GetPublisherCount())
	requireValue(t, 1)(sub2.GetPublisherCount())
	requireValue(t, 1)(subUnreliable.GetPublisherCount())

	subOther := noErr(subNode.NewSubscription("other_string", std_msgs_msg.StringTypeSupport, nil, noopSub))(t)

	requireValue(t, 2)(pub.GetSubscriptionCount())
	requireValue(t, 1)(pubUnreliable.GetSubscriptionCount())
	requireValue(t, 1)(sub1.GetPublisherCount())
	requireValue(t, 1)(sub2.GetPublisherCount())
	requireValue(t, 1)(subUnreliable.GetPublisherCount())
	requireValue(t, 0)(subOther.GetPublisherCount())
}
