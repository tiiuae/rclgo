//#nosec G404

package rclgo_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

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
	os.Setenv("ROS_DOMAIN_ID", fmt.Sprint(getUnusedDomainID()))
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
