/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

func GenerateGolangMessageTypes(rootPath string, destPath string) map[string]struct{} {
	generatedMessages := map[string]struct{}{}
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		skip, blacklistEntry := blacklisted(path)
		if skip {
			fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
			return nil
		}

		if re.M(path, `m!/msg/.+?\.msg$!`) {
			fmt.Printf("Generating: %s\n", path)
			md, err := generateMessage(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", path, destPath, err)
			}
			generatedMessages[md.Package+"/"+md.Type] = struct{}{}
		} else if re.M(path, `m!/srv/.+?\.srv$!`) {
			fmt.Printf("Generating: %s\n", path)
			srv, err := generateService(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Service '%s' to '%s', error: %v\n", path, destPath, err)
				return nil
			}
			generatedMessages[srv.Request.Package+"/"+srv.Request.Type] = struct{}{}
			generatedMessages[srv.Request.Package+"/"+srv.Request.Type] = struct{}{}
		}
		return nil
	})
	return generatedMessages
}

func generateMessage(sourcePath string, destPathPkgRoot string) (*ROS2Message, error) {
	md := ROS2MessageNew("", "")
	var err error

	md.Metadata, err = parseMetadataFromPath(sourcePath)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	err = ParseROS2Message(md, string(content))
	if err != nil {
		return nil, err
	}

	destFile, err := CreateTargetGolangTypeFile(destPathPkgRoot, md.Metadata)
	if err != nil {
		return nil, err
	}
	defer destFile.Close()

	err = ros2MsgToGolangTypeTemplate.Execute(destFile, md)
	if err != nil {
		return nil, err
	}
	return md, nil
}

func parseMetadataFromPath(p string) (*Metadata, error) {
	base := filepath.Base(p)
	ext := filepath.Ext(base)
	m := &Metadata{
		Name: strings.TrimSuffix(base, ext),
		Type: ext[1:],
	}
	dirs := strings.Split(p, "/")

	if len(dirs) >= 2 {
		m.Package = dirs[len(dirs)-3]
	} else {
		return nil, fmt.Errorf("Path '%s' cannot be parsed for ROS2 package name!", p)
	}

	return m, nil
}

func CreateTargetGolangTypeFile(destPathPkgRoot string, m *Metadata) (*os.File, error) {
	destFilePath := filepath.Join(destPathPkgRoot, m.ImportPath(), m.Name+".gen.go")
	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return nil, err
	}
	return destFile, err
}

func generateService(srcPath, dstPathPkgRoot string) (*ROS2Service, error) {
	m, err := parseMetadataFromPath(srcPath)
	if err != nil {
		return nil, err
	}
	service := NewROS2Service(m.Package, m.Name)
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	if err = parseService(service, string(srcFile)); err != nil {
		return nil, err
	}

	srvFile, err := CreateTargetGolangTypeFile(dstPathPkgRoot, service.Metadata)
	if err != nil {
		return nil, err
	}
	defer srvFile.Close()
	err = ros2ServiceToGolangTypeTemplate.Execute(srvFile, service)
	if err != nil {
		return nil, err
	}

	reqFile, err := CreateTargetGolangTypeFile(dstPathPkgRoot, service.Request.Metadata)
	if err != nil {
		return nil, err
	}
	defer reqFile.Close()
	err = ros2MsgToGolangTypeTemplate.Execute(reqFile, service.Request)
	if err != nil {
		return nil, err
	}

	respFile, err := CreateTargetGolangTypeFile(dstPathPkgRoot, service.Response.Metadata)
	if err != nil {
		return nil, err
	}
	defer respFile.Close()
	err = ros2MsgToGolangTypeTemplate.Execute(respFile, service.Response)
	if err != nil {
		return nil, err
	}

	return service, nil
}
