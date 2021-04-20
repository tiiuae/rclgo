package ros2

import (
	"fmt"
	"os"
)

func ExampleNewRCLArgs() {
	oldOSArgs := os.Args
	defer func() { os.Args = oldOSArgs }()

	os.Args = []string{"--extra0", "args0", "--ros-args", "--log-level", "DEBUG", "--", "--extra1", "args1"}
	rosArgs, err := NewRCLArgs("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> [--extra0 args0 --ros-args --log-level DEBUG -- --extra1 args1]

	rosArgs, err = NewRCLArgs("--log-level INFO")
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> [--ros-args --log-level INFO]

	os.Args = []string{"--extra0", "args0", "--extra1", "args1"}
	rosArgs, err = NewRCLArgs("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("rosArgs: %+v\n", rosArgs.GoArgs) // -> []

	// Output: rosArgs: [--extra0 args0 --ros-args --log-level DEBUG -- --extra1 args1]
	// rosArgs: [--ros-args --log-level INFO]
	// rosArgs: []
}
