package main

import (
	"context"
	"fmt"
	"os"

	example_srvs_srv "github.com/tiiuae/rclgo/examples/service_client/service/msgs/example_srvs/srv"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

//go:generate go run github.com/tiiuae/rclgo/cmd/rclgo-gen generate -d msgs --include-go-package-deps ./... --cgo-flags-path ""

func run() error {
	spinCtx, cancelSpin := context.WithCancel(context.Background())
	defer cancelSpin()

	ctx, err := rclgo.NewContext(0, nil)
	if err != nil {
		return fmt.Errorf("error while creating context: %v", err)
	}

	node, err := ctx.NewNode("service_node", "/test")
	if err != nil {
		return fmt.Errorf("error while creating node: %v", err)
	}

	_, err = example_srvs_srv.NewAddTwoIntsService(node, "add", nil, func(info *rclgo.ServiceInfo, req *example_srvs_srv.AddTwoInts_Request, sender example_srvs_srv.AddTwoIntsServiceResponseSender) {
		res := example_srvs_srv.NewAddTwoInts_Response()
		res.Sum = req.A + req.B
		sender.SendResponse(res)
	})

	if err != nil {
		return fmt.Errorf("error occured starting service: %v", err)
	}

	go ctx.Spin(spinCtx)

	select {}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
