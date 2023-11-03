package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	std_msgs_msg "github.com/tiiuae/rclgo/examples/publisher_subscriber/msgs/std_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

func run() error {
	rclArgs, restArgs, err := rclgo.ParseArgs(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse ROS args: %v", err)
	}

	if err := rclgo.Init(rclArgs); err != nil {
		return fmt.Errorf("failed to initialize rclgo: %v", err)
	}
	defer rclgo.Uninit()

	node, err := rclgo.NewNode("publisher", "")
	if err != nil {
		return fmt.Errorf("failed to create node: %v", err)
	}
	defer node.Close()

	pub, err := std_msgs_msg.NewStringPublisher(node, "hello", nil)
	if err != nil {
		return fmt.Errorf("failed to create publisher: %v", err)
	}
	defer pub.Close()

	msg := std_msgs_msg.NewString()
	if len(restArgs) > 0 {
		msg.Data = restArgs[0]
	} else {
		msg.Data = "gopher"
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		node.Logger().Infof("Publishing: %#v", msg)
		if err := pub.Publish(msg); err != nil {
			return fmt.Errorf("failed to publish: %v", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
