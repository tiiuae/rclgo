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
)

func GenerateROS2AllMessagesImporter(destPathPkgRoot string, ros2Messages map[string]struct{}) error {

	destFilePath := filepath.Join(destPathPkgRoot, "..", "msgs", "ros2msgs.gen.go")

	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Generating all importer: %s\n", destFilePath)
	return ros2MsgImportAllPackage.Execute(destFile, ros2Messages)
}
