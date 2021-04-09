GO111MODULE = on

all: configure build #default make target

configure:
	echo "" # make configure unimplemented. Unified deployment API between all software components.

build:
	cd cmd/rclgo     && go build
	cd cmd/rclgo-gen && go build

install:
	cd cmd/rclgo     && go install
	cd cmd/rclgo-gen && go install

.PHONY: test
test:
	go test -v ./...

generate:
	go run cmd/rclgo-gen/main.go generate

