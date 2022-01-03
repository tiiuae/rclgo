/*
This file is part of rclgo

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Apache License, Version 2.0 (the "License");
    http://www.apache.org/licenses/LICENSE-2.0
*/

package gogen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/kivilahtio/go-re/v0"
)

type Config struct {
	RclgoImportPath     string
	MessageModulePrefix string
}

var DefaultConfig = Config{
	RclgoImportPath:     "github.com/tiiuae/rclgo",
	MessageModulePrefix: "github.com/tiiuae/rclgo-msgs",
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

type generator struct {
	config *Config
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
	cImports := make(map[string]stringSet)
	getOrDefault := func(s string) stringSet {
		set := cImports[s]
		if set == nil {
			set = stringSet{}
			cImports[s] = set
		}
		return set
	}
	g := generator{config: c}
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		skip, blacklistEntry := blacklisted(path)
		if skip {
			fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
			return nil
		}
		slashPath := filepath.ToSlash(path)
		if re.M(slashPath, `m!/msg/.+\.msg$!`) {
			fmt.Printf("Generating: %s\n", path)
			result, err := g.generateMessage(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", path, destPath, err)
			}
			set := getOrDefault(result.GoPackage())
			set.AddFrom(result.CImports)
		} else if re.M(slashPath, `m!/srv/.+\.srv$!`) {
			fmt.Printf("Generating: %s\n", path)
			result, err := g.generateService(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Service '%s' to '%s', error: %v\n", path, destPath, err)
			}
			set := getOrDefault(result.GoPackage())
			set.AddFrom(result.Request.CImports)
			set.AddFrom(result.Response.CImports)
		}
		return nil
	})
	for pkg, imports := range cImports {
		err := g.generateCommonPackageGoFile(pkg, imports, destPath)
		if err != nil {
			fmt.Printf("Failed to generate common package file for package %s: %v", pkg, err)
		}
	}
	return nil
}

func (g *generator) generateMessage(sourcePath string, destPathPkgRoot string) (*ROS2Message, error) {
	msg := ROS2MessageNew("", "")
	var err error

	msg.Metadata, err = parseMetadataFromPath(sourcePath)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	parser := parser{config: g.config}
	err = parser.ParseROS2Message(msg, string(content))
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, destPathPkgRoot, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
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

func (g *generator) generateService(srcPath, dstPathPkgRoot string) (*ROS2Service, error) {
	m, err := parseMetadataFromPath(srcPath)
	if err != nil {
		return nil, err
	}
	service := NewROS2Service(m.Package, m.Name)
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	parser := parser{config: g.config}
	if err = parser.parseService(service, string(srcFile)); err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, dstPathPkgRoot, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

type templateData = map[string]interface{}

func (g *generator) generateGoFile(dstPathPkgRoot string, meta *Metadata, tmpl *template.Template, data templateData) error {
	f, err := createTargetGolangTypeFile(dstPathPkgRoot, meta)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, ok := data["Config"]; !ok {
		data["Config"] = g.config
	}
	return tmpl.Execute(f, data)
}

func (g *generator) generateMessageGoFile(parser *parser, dstPathPkgRoot string, msg *ROS2Message) error {
	return g.generateGoFile(
		dstPathPkgRoot,
		msg.Metadata,
		ros2MsgToGolangTypeTemplate,
		templateData{
			"Message":             msg,
			"cSerializationCode":  parser.cSerializationCode,
			"goSerializationCode": parser.goSerializationCode,
		},
	)
}

func (g *generator) generateServiceGoFiles(parser *parser, dstPathPkgRoot string, srv *ROS2Service) error {
	err := g.generateGoFile(
		dstPathPkgRoot,
		srv.Metadata,
		ros2ServiceToGolangTypeTemplate,
		templateData{"Service": srv},
	)
	if err != nil {
		return err
	}
	err = g.generateMessageGoFile(parser, dstPathPkgRoot, srv.Request)
	if err != nil {
		return err
	}
	return g.generateMessageGoFile(parser, dstPathPkgRoot, srv.Response)
}

func (g *generator) generateCommonPackageGoFile(goPkg string, cImports stringSet, destRoot string) error {
	i := strings.LastIndex(goPkg, "_")
	if i < 0 || i > len(goPkg)-1 {
		return errors.New("package type suffix is missing or incorrect")
	}
	cPkg := goPkg[:i]
	pkgType := goPkg[i+1:]
	f, err := mkdir_p(filepath.Join(destRoot, cPkg, pkgType, "common.gen.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	return ros2PackageCommonTemplate.Execute(f, templateData{
		"GoPackage": goPkg,
		"CPackage":  cPkg,
		"CImports":  cImports,
	})
}
