package rclgo

import (
	"fmt"
)

func ExampleParseArgs() {
	rosArgs, restArgs, err := ParseArgs([]string{"--extra0", "args0", "--ros-args", "--log-level", "DEBUG", "--", "--extra1", "args1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = ParseArgs([]string{"--ros-args", "--log-level", "INFO"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = ParseArgs([]string{"--extra0", "args0", "--extra1", "args1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: [%v]\n", rosArgs)
	fmt.Printf("restArgs: %+v\n\n", restArgs)

	rosArgs, restArgs, err = ParseArgs(nil)
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
