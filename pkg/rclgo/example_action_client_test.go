//nolint:revive
package rclgo_test

import (
	"context"
	"fmt"

	example_interfaces_action "github.com/tiiuae/rclgo/internal/msgs/example_interfaces/action"
	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

func ExampleActionClient() {
	err := rclgo.Init(nil)
	if err != nil {
		// handle error
	}
	defer rclgo.Uninit()
	node, err := rclgo.NewNode("my_node", "my_namespace")
	if err != nil {
		// handle error
	}
	client, err := node.NewActionClient(
		"fibonacci",
		example_interfaces_action.FibonacciTypeSupport,
		nil,
	)
	if err != nil {
		// handle error
	}
	ctx := context.Background()
	goal := example_interfaces_action.NewFibonacci_Goal()
	goal.Order = 10
	result, _, err := client.WatchGoal(ctx, goal, func(ctx context.Context, feedback types.Message) {
		fmt.Println("Got feedback:", feedback)
	})
	if err != nil {
		// handle error
	}
	fmt.Println("Got result:", result)
}

func ExampleActionClient_type_safe_wrapper() {
	err := rclgo.Init(nil)
	if err != nil {
		// handle error
	}
	defer rclgo.Uninit()
	node, err := rclgo.NewNode("my_node", "my_namespace")
	if err != nil {
		// handle error
	}
	client, err := example_interfaces_action.NewFibonacciClient(
		node,
		"fibonacci",
		nil,
	)
	if err != nil {
		// handle error
	}
	ctx := context.Background()
	goal := example_interfaces_action.NewFibonacci_Goal()
	goal.Order = 10
	result, _, err := client.WatchGoal(ctx, goal, func(ctx context.Context, feedback *example_interfaces_action.Fibonacci_FeedbackMessage) {
		fmt.Println("Got feedback:", feedback)
	})
	if err != nil {
		// handle error
	}
	fmt.Println("Got result:", result)
}
