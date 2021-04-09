/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Generate_rosidl_runtime_c_sequence_handlers(destPathPkgRoot string) error {

	destFilePath := filepath.Join(destPathPkgRoot, "..", "rosidl_runtime_c", "Primitives.go")

	_, err := os.Stat(destFilePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' is missing? It should exist relative the the given destination path '%s'", destFilePath, destPathPkgRoot)
		err = os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}

	return ros2rosidl_runtime_c_handlers.Execute(destFile, map[string]interface{}{
		"PMap": &ROSIDL_RUNTIME_C_PRIMITIVE_TYPES_MAPPING,
	})
}

func GenerateGolangTypeFromROS2MessagePath(sourcePath string, destPathPkgRoot string) (*ROS2Message, error) {
	md := ROS2MessageNew("", "")

	err := ParseMessageMetadataFromPath(sourcePath, md)
	if err != nil {
		return md, err
	}

	content, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return md, err
	}

	err = ParseROS2Message(md, string(content))
	if err != nil {
		return md, err
	}

	destFile, err := CreateTargetGolangTypeFile(destPathPkgRoot, md)
	if err != nil {
		return md, err
	}

	err = ros2MsgToGolangTypeTemplate.Execute(destFile, md)
	if err != nil {
		return md, err
	}
	return md, nil
}

func ParseMessageMetadataFromPath(p string, md *ROS2Message) error {
	var dirs []string
	var ros2msgName string
	var ros2pkgName string
	var ros2dataStructureName string

	dirs = strings.Split(p, "/")
	ros2msgName = strings.TrimSuffix(filepath.Base(p), ".msg")

	if len(dirs) >= 2 {
		ros2pkgName = dirs[len(dirs)-3]
		ros2dataStructureName = dirs[len(dirs)-2]
	} else {
		return fmt.Errorf("Path '%s' cannot be parsed for ROS2 package name!", p)
	}

	md.RosMsgName = ros2msgName
	md.DataStructureType = ros2dataStructureName
	md.RosPackage = ros2pkgName
	md.Url = p

	return nil
}

func CreateTargetGolangTypeFile(destPathPkgRoot string, md *ROS2Message) (*os.File, error) {
	destFilePath := filepath.Join(destPathPkgRoot, md.RosPackage, "msg", md.RosMsgName+".go")
	destFileDir := filepath.Join(destPathPkgRoot, md.RosPackage, "msg")
	_, err := os.Stat(destFileDir)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating directory '%s'", destFileDir)
		err = os.MkdirAll(destFileDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return destFile, err
	}
	return destFile, err
}

func ROS2MessageListToArray(l *list.List) []*ROS2Message {
	mdsArray := make([]*ROS2Message, l.Len())
	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		md, ok := e.Value.(*ROS2Message)
		if !ok {
			panic(fmt.Sprintf("ROS2MessageListToArray():> One of the ROS2Messages at index '%d' value '%+v' is not a ROS2Message!", i, e.Value))
		}
		mdsArray[i] = md
		i++
	}
	return mdsArray
}
