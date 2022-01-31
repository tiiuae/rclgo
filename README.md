rclgo the ROS2 client library Golang wrapper
============================================

[![Go Reference](https://pkg.go.dev/badge/github.com/tiiuae/rclgo.svg)][docs]

## Getting started

rclgo is used with the Go module system like most Go libraries. It also requires
a ROS 2 installation as well as C bindings to any ROS interface type used with
the library. The ROS core components can be installed by installing the Debian
package `ros-galactic-ros-core`. C bindings for ROS interfaces are also usually
distributed as Debian packages.

API documentation is available at [pkg.go.dev][docs].

### ROS 2 interface bindings

rclgo requires Go bindings of all the ROS 2 interfaces to exist. rclgo-gen is
used to generate Go bindings to existing interface types installed. The
generated Go bindings depend on the corresponding C bindings to be installed.

rclgo-gen can be installed globally by running

    go install github.com/tiiuae/rclgo/cmd/rclgo-gen@latest

but it is recommended to add rclgo-gen as a dependency to the project to ensure
the version matches that of rclgo. This can be done by adding a file `tools.go`
to the main package of the project containing something similar to the
following:
```go
//go:build tools

package main

import _ "github.com/tiiuae/rclgo/cmd/rclgo-gen"
```
Then run `go mod tidy`. This version of rclgo-gen can be used by running

    go run github.com/tiiuae/rclgo/cmd/rclgo-gen generate -d msgs

in the project directory. The command can be added as a `go generate` comment to
one of the source files in the project, such as `main.go`, as follows:
```go
//go:generate go run github.com/tiiuae/rclgo/cmd/rclgo-gen generate -d msgs
```

[docs]: https://pkg.go.dev/github.com/tiiuae/rclgo/pkg/rclgo
