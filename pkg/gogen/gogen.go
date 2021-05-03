/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"fmt"
	"path/filepath"
	"runtime"
)

/*
Looks up where the rclgo-module is installed and returns a path where to write Golang bindings for ROS2 messages.
*/
func GetGoConvertedROS2MsgPackagesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(file, "../..", "ros2/msgs")
}

func Generate(rootPath, destPath, rclcPath string) {
	err := Generate_rosidl_runtime_c_sequence_handlers(destPath)
	if err != nil {
		fmt.Printf("Error: '%+v'\n", err)
	}
	err = GenerateROS2ErrorTypes(rootPath, destPath, rclcPath)
	if err != nil {
		fmt.Printf("Error: '%+v'\n", err)
	}
	ros2MessagesAry := GenerateGolangMessageTypes(rootPath, destPath)
	err = GenerateROS2AllMessagesImporter(destPath, ros2MessagesAry)
	if err != nil {
		fmt.Printf("Error: '%+v'\n", err)
	}
}
