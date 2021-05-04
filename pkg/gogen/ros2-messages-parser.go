/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

func GenerateGolangMessageTypes(rootPath string, destPath string) map[string]*ROS2Message {
	ros2MessagesList := list.New()
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		skip, blacklistEntry := blacklisted(path)
		if skip {
			fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
			return nil
		}

		if re.M(path, `m!/msg/.+?\.msg$!`) {
			fmt.Printf("Generating: %s\n", path)
			md, err := GenerateGolangTypeFromROS2MessagePath(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", path, destPath, err)
			}
			ros2MessagesList.PushBack(md)
		}
		return nil
	})
	ros2MessagesAry := make(map[string]*ROS2Message, ros2MessagesList.Len())
	e := ros2MessagesList.Front()
	for e != nil {
		m := e.Value.(*ROS2Message)
		ros2MessagesAry[m.RosPackage] = m
		e = e.Next()
	}
	return ros2MessagesAry
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
	destFilePath := filepath.Join(destPathPkgRoot, md.RosPackage, "msg", md.RosMsgName+".gen.go")
	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return nil, err
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
