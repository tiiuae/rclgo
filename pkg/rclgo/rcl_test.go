//#nosec G404

package rclgo_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func getNextDomainID() int {
	id, _ := strconv.ParseUint(os.Getenv("ROS_DOMAIN_ID"), 10, 8)
	newID := rand.Intn(101)
	for ; newID == int(id); newID = rand.Intn(101) {
	}
	return newID
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	os.Setenv("ROS_DOMAIN_ID", fmt.Sprint(getNextDomainID()))
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
