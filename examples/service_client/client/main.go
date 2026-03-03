package main

import (
	"context"
	"fmt"

	example_srvs_srv "github.com/tiiuae/rclgo/examples/service_client/client/msgs/example_srvs/srv"
	"github.com/tiiuae/rclgo/pkg/rclgo"
)

//go:generate go run github.com/tiiuae/rclgo/cmd/rclgo-gen generate -d msgs --include-go-package-deps ./... --cgo-flags-path ""

func run() {
	spinCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	qosProfile := rclgo.NewDefaultServiceQosProfile()

	serviceCtx, err := rclgo.NewContext(0, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := serviceCtx.NewNode("client_node", "/test")
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := node.NewClient("add", example_srvs_srv.AddTwoIntsTypeSupport, &rclgo.ClientOptions{Qos: qosProfile})
	if err != nil {
		fmt.Println(err)
		return
	}

	go serviceCtx.Spin(spinCtx)

	req := example_srvs_srv.NewAddTwoInts_Request()
	req.A = 42
	req.B = 27

	res, _, err := client.Send(spinCtx, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)

}

func main() {
	run()
}
