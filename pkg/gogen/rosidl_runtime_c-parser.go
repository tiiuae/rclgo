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

func Generate_rosidl_runtime_c_sequence_handlers(destPathPkgRoot string) error {

	destFilePath := filepath.Join(destPathPkgRoot, "..", "rosidl_runtime_c", "Primitives.go")

	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Generating roside_runtime_c: %s\n", destFilePath)
	return ros2rosidl_runtime_c_handlers.Execute(destFile, map[string]interface{}{
		"PMap": &ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING,
	})
}
