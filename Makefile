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
	go test -v ./...

generate:
	@pkgs=(
	    builtin_interfaces
	    example_interfaces
	    geometry_msgs
	    sensor_msgs
	    std_msgs
	    std_srvs
	    test_msgs
	)

	dest_path=internal/msgs

	rm -rf "$$dest_path/"*
	for pkg in $${pkgs[@]}; do
	    go run ./cmd/rclgo-gen generate \
	        --message-module-prefix "github.com/tiiuae/rclgo/$$dest_path" \
	        -r "/opt/ros/galactic/share/$$pkg" \
	        -d "$$dest_path" \
			|| exit 1
	done
	rm "$$dest_path/msgs.gen.go" || exit 1
	go run cmd/rclgo-gen/main.go generate-rclgo

vet:
	go vet ./...
