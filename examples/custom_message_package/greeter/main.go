package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	greeting_msgs "github.com/tiiuae/rclgo/examples/custom_message_package/greeter/msgs/greeting_msgs/msg"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

//go:generate go run github.com/tiiuae/rclgo/cmd/rclgo-gen generate -d msgs --include-go-package-deps ./...

func run() error {
	rclArgs, restArgs, err := rclgo.ParseArgs(os.Args[1:])
	if err != nil {
		return fmt.Errorf("failed to parse ROS args: %v", err)
	}
	if err := rclgo.Init(rclArgs); err != nil {
		return fmt.Errorf("failed to initialize rclgo: %v", err)
	}
	defer rclgo.Uninit()
	node, err := rclgo.NewNode("greeter", "")
	if err != nil {
		return fmt.Errorf("failed to create node: %v", err)
	}
	pub, err := greeting_msgs.NewHelloPublisher(node, "~/hello", nil)
	if err != nil {
		return fmt.Errorf("failed to create publisher: %v", err)
	}
	greeting := greeting_msgs.NewHello()
	if len(restArgs) > 0 {
		greeting.Name = restArgs[0]
	} else {
		greeting.Name = "gopher"
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	for {
		node.Logger().Infof("Publishing greeting: %#v", greeting)
		if err := pub.Publish(greeting); err != nil {
			return fmt.Errorf("failed to publish: %v", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
