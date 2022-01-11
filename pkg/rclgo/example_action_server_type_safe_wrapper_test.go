package rclgo_test

import (
	"context"
	"errors"

	example_interfaces_action "github.com/tiiuae/rclgo/internal/msgs/example_interfaces/action"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

var typeSafeFibonacci = example_interfaces_action.NewFibonacciAction(
	func(
		ctx context.Context, goal *example_interfaces_action.FibonacciGoalHandle,
	) (*example_interfaces_action.Fibonacci_Result, error) {
		if goal.Description.Order < 0 {
			return nil, errors.New("order must be non-negative")
		}
		sender, err := goal.Accept()
		if err != nil {
			return nil, err
		}
		logger := goal.Server().Node().Logger()
		result := example_interfaces_action.NewFibonacci_Result()
		fb := example_interfaces_action.NewFibonacci_Feedback()
		var x, y, i int32
		for y = 1; i < goal.Description.Order; x, y, i = y, x+y, i+1 {
			result.Sequence = append(result.Sequence, x)
			fb.Sequence = result.Sequence
			if err = sender.Send(fb); err != nil {
				logger.Error("failed to send feedback: ", err)
			}
		}
		return result, nil
	},
)

func ExampleActionServer_type_safe_wrapper() {
	err := rclgo.Init(nil)
	if err != nil {
		// handle error
	}
	defer rclgo.Uninit()
	node, err := rclgo.NewNode("my_node", "my_namespace")
	if err != nil {
		// handle error
	}
	_, err = example_interfaces_action.NewFibonacciServer(node, "fibonacci", typeSafeFibonacci, nil)
	if err != nil {
		// handle error
	}
	ctx := context.Background()
	if err = rclgo.Spin(ctx); err != nil {
		// handle error
	}
}
