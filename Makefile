GO111MODULE = on

all: configure build #default make target

configure:
	echo "" # make configure unimplemented. Unified deployment API between all software components.

build: vet
	cd cmd/rclgo     && go build
	cd cmd/rclgo-gen && go build

install:
	cd cmd/rclgo     && go install
	cd cmd/rclgo-gen && go install

.PHONY: test
test:
	go test -v `go list ./... | sed -e '\@^github\.com/tiiuae/rclgo/pkg/ros2/msgs@D' -e '\@^github\.com/tiiuae/rclgo/cmd@D'`

generate:
	go run cmd/rclgo-gen/main.go generate

vet:
	go vet ./...
