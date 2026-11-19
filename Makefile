GO111MODULE = on

.ONESHELL:
SHELL = /bin/bash

# default make target
all: build

build:
	cd cmd/rclgo-gen && go build

install:
	cd cmd/rclgo-gen && go install

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: test-verbose
test-verbose:
	go test -v -count=1 ./...

generate:
	dest_path=internal/msgs

	rm -rf "$$dest_path/"*
	go run ./cmd/rclgo-gen generate \
		--root-path /usr \
		--root-path /opt/ros/${ROS_DISTRO} \
		--dest-path "$$dest_path" \
		--message-module-prefix "github.com/ATIinc/rclgo/$$dest_path" \
		--license-header-path ./license-header.txt \
		--include-go-package-deps ./... \
		--cgo-flags-path "" \
		|| exit 1
	rm "$$dest_path/msgs.gen.go" || exit 1
	go run ./cmd/rclgo-gen generate-rclgo \
		--root-path /usr \
		--root-path /opt/ros/${ROS_DISTRO} \
		--license-header-path ./license-header.txt

lint:
	golangci-lint run ./...
