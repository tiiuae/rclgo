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

### Developing with custom interface types

By default `rclgo-gen generate` looks for interface definitions in
`$AMENT_PREFIX_PATH` and generates Go bindings for all interfaces it finds.
`$AMENT_PREFIX_PATH` contains a list of paths to the ROS 2 underlay and overlays
you have sourced. If you want to use interfaces which are not available in
`$AMENT_PREFIX_PATH`, you can either add the directories to `$AMENT_PREFIX_PATH`
(use the package as an overlay) or pass paths using the `--root-path` option. If
multiple `--root-path` options are passed, the paths are searched in the order
the options are passed. If multiple identically named package and interface name
combinations are found the first one is used. When `--root-path` is used,
`rclgo-gen generate` won't look for interface definitions from
`$AMENT_PREFIX_PATH`. If you want to generate bindings for files using both
`$AMENT_PREFIX_PATH` and custom paths, you can pass
`--root-path="$AMENT_PREFIX_PATH"` in addition to the other paths.

In addition to generating the Go bindings, you must also generate the C
bindings, compile them and make them available to the Go tool. Generated Go
bindings include the `include` and `lib` subdirectories of the root paths used
when generating the bindings in the header and library search paths,
respectively. When building a package using the Go bindings the environment
variables `CGO_CFLAGS` and `CGO_LDFLAGS` can be used to pass additional `-I` and
`-L` options, respectively, if needed.

An example is available in
[examples/custom_message_package](examples/custom_message_package).

[docs]: https://pkg.go.dev/github.com/tiiuae/rclgo/pkg/rclgo
