package rclgo_test

import (
	"context"
	"errors"

	example_interfaces_action "github.com/tiiuae/rclgo/internal/msgs/example_interfaces/action"
	"github.com/tiiuae/rclgo/pkg/rclgo"
	"github.com/tiiuae/rclgo/pkg/rclgo/types"
)

var fibonacci = rclgo.NewAction(
	example_interfaces_action.FibonacciTypeSupport,
	func(ctx context.Context, goal *rclgo.GoalHandle) (types.Message, error) {
		description := goal.Description.(*example_interfaces_action.Fibonacci_Goal)
		if description.Order < 0 {
			return nil, errors.New("order must be non-negative")
		}
		sender, err := goal.Accept()
		if err != nil {
			return nil, err
		}
		result := example_interfaces_action.NewFibonacci_Result()
		fb := example_interfaces_action.NewFibonacci_Feedback()
		var x, y, i int32
		for y = 1; i < description.Order; x, y, i = y, x+y, i+1 {
			result.Sequence = append(result.Sequence, x)
			fb.Sequence = result.Sequence
			if err = sender.Send(fb); err != nil {
				goal.Logger().Error("failed to send feedback: ", err)
			}
		}
		return result, nil
	},
)

func ExampleActionServer() {
	err := rclgo.Init(nil)
	if err != nil {
		// handle error
	}
	defer rclgo.Uninit()
	node, err := rclgo.NewNode("my_node", "my_namespace")
	if err != nil {
		// handle error
	}
	_, err = node.NewActionServer("fibonacci", fibonacci, nil)
	if err != nil {
		// handle error
	}
	ctx := context.Background()
	if err = rclgo.Spin(ctx); err != nil {
		// handle error
	}
}
