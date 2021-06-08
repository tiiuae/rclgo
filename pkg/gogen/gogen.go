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
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kivilahtio/go-re/v0"
)

type Config struct {
	RclgoImportPath     string
	MessageModulePrefix string
}

var DefaultConfig = Config{
	RclgoImportPath:     "github.com/tiiuae/rclgo",
	MessageModulePrefix: "github.com/tiiuae/rclgo/pkg/rclgo/msgs",
}

// RclgoRepoRootPath returns the path to the root of the rclgo repository.
// Panics if the path can't be determined.
func RclgoRepoRootPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not determine rclgo repo root path")
	}
	return filepath.Join(file, "../../..")
}

func GeneratePrimitives(c *Config, destFilePath string) error {
	destFilePath = filepath.Join(destFilePath, "pkg/rclgo/primitives/primitives.gen.go")
	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	fmt.Printf("Generating primitive types: %s\n", destFilePath)
	return primitiveTypes.Execute(destFile, map[string]interface{}{
		"PMap":   &primitiveTypeMappings,
		"Config": c,
	})
}

func GenerateROS2AllMessagesImporter(c *Config, destPathPkgRoot string) error {
	msgDirs, err := filepath.Glob(filepath.Join(destPathPkgRoot, "*/msg"))
	if err != nil {
		return err
	}
	srvDirs, err := filepath.Glob(filepath.Join(destPathPkgRoot, "*/srv"))
	if err != nil {
		return err
	}
	pkgs := map[string]struct{}{}
	for _, d := range msgDirs {
		pkgs[path.Join(path.Base(path.Dir(d)), path.Base(d))] = struct{}{}
	}
	for _, d := range srvDirs {
		pkgs[path.Join(path.Base(path.Dir(d)), path.Base(d))] = struct{}{}
	}

	destFilePath := filepath.Join(destPathPkgRoot, "msgs.gen.go")

	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	fmt.Printf("Generating all importer: %s\n", destFilePath)
	return ros2MsgImportAllPackage.Execute(destFile, map[string]interface{}{
		"Packages": pkgs,
		"Config":   c,
	})
}

func GenerateGolangMessageTypes(c *Config, rootPath string, destPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		skip, blacklistEntry := blacklisted(path)
		if skip {
			fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
			return nil
		}

		if re.M(filepath.ToSlash(path), `m!/msg/.+?\.msg$!`) {
			fmt.Printf("Generating: %s\n", path)
			_, err := generateMessage(c, path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", path, destPath, err)
			}
		} else if re.M(filepath.ToSlash(path), `m!/srv/.+?\.srv$!`) {
			fmt.Printf("Generating: %s\n", path)
			_, err := generateService(c, path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Service '%s' to '%s', error: %v\n", path, destPath, err)
				return nil
			}
		}
		return nil
	})
}

func generateMessage(c *Config, sourcePath string, destPathPkgRoot string) (*ROS2Message, error) {
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

	parser := parser{config: c}
	err = parser.ParseROS2Message(md, string(content))
	if err != nil {
		return nil, err
	}

	destFile, err := createTargetGolangTypeFile(destPathPkgRoot, md.Metadata)
	if err != nil {
		return nil, err
	}
	defer destFile.Close()

	err = ros2MsgToGolangTypeTemplate.Execute(destFile, map[string]interface{}{
		"Message":             md,
		"Config":              c,
		"cSerializationCode":  parser.cSerializationCode,
		"goSerializationCode": parser.goSerializationCode,
	})
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

func createTargetGolangTypeFile(destPathPkgRoot string, m *Metadata) (*os.File, error) {
	destFilePath := filepath.Join(destPathPkgRoot, m.ImportPath(), m.Name+".gen.go")
	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return nil, err
	}
	return destFile, err
}

func generateService(c *Config, srcPath, dstPathPkgRoot string) (*ROS2Service, error) {
	m, err := parseMetadataFromPath(srcPath)
	if err != nil {
		return nil, err
	}
	service := NewROS2Service(m.Package, m.Name)
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	parser := parser{config: c}
	if err = parser.parseService(service, string(srcFile)); err != nil {
		return nil, err
	}

	srvFile, err := createTargetGolangTypeFile(dstPathPkgRoot, service.Metadata)
	if err != nil {
		return nil, err
	}
	defer srvFile.Close()
	err = ros2ServiceToGolangTypeTemplate.Execute(srvFile, map[string]interface{}{
		"Service": service,
		"Config":  c,
	})
	if err != nil {
		return nil, err
	}

	reqFile, err := createTargetGolangTypeFile(dstPathPkgRoot, service.Request.Metadata)
	if err != nil {
		return nil, err
	}
	defer reqFile.Close()
	err = ros2MsgToGolangTypeTemplate.Execute(reqFile, map[string]interface{}{
		"Message":             service.Request,
		"Config":              c,
		"cSerializationCode":  parser.cSerializationCode,
		"goSerializationCode": parser.goSerializationCode,
	})
	if err != nil {
		return nil, err
	}

	respFile, err := createTargetGolangTypeFile(dstPathPkgRoot, service.Response.Metadata)
	if err != nil {
		return nil, err
	}
	defer respFile.Close()
	err = ros2MsgToGolangTypeTemplate.Execute(respFile, map[string]interface{}{
		"Message":             service.Response,
		"Config":              c,
		"cSerializationCode":  parser.cSerializationCode,
		"goSerializationCode": parser.goSerializationCode,
	})
	if err != nil {
		return nil, err
	}

	return service, nil
}
