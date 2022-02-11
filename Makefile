GO111MODULE = on

.ONESHELL:
SHELL = /bin/bash

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
	go test -v -count=1 ./...

generate:
	dest_path=internal/msgs

	rm -rf "$$dest_path/"*
	go run ./cmd/rclgo-gen generate \
	    --root-path /opt/ros/galactic \
	    --dest-path "$$dest_path" \
		--message-module-prefix "github.com/tiiuae/rclgo/$$dest_path" \
		--include-package action_msgs \
		--include-package builtin_interfaces \
		--include-package example_interfaces \
		--include-package geometry_msgs \
		--include-package sensor_msgs \
		--include-package std_msgs \
		--include-package std_srvs \
		--include-package test_msgs \
		--include-package unique_identifier_msgs \
		|| exit 1
	rm "$$dest_path/msgs.gen.go" || exit 1
	go run ./cmd/rclgo-gen generate-rclgo --root-path /opt/ros/galactic

vet:
	go vet ./...
