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
	"io/fs"
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
	return primitiveTypes.Execute(destFile, templateData{
		"PMap":   &primitiveTypeMappings,
		"Config": c,
	})
}

func GenerateROS2AllMessagesImporter(c *Config, destPathPkgRoot string) error {
	pkgs := map[string]struct{}{}
	for _, glob := range []string{"*/msg", "*/srv", "*/action"} {
		dirs, err := filepath.Glob(filepath.Join(destPathPkgRoot, glob))
		if err != nil {
			return err
		}
		for _, d := range dirs {
			pkgs[path.Join(
				filepath.Base(filepath.Dir(d)),
				filepath.Base(d),
			)] = struct{}{}
		}
	}

	destFilePath := filepath.Join(destPathPkgRoot, "msgs.gen.go")

	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	fmt.Printf("Generating all importer: %s\n", destFilePath)
	return ros2MsgImportAllPackage.Execute(destFile, templateData{
		"Packages": pkgs,
		"Config":   c,
	})
}

func GenerateGolangMessageTypes(c *Config, rootPaths []string, destPath string) error {
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
	for meta, path := range findInterfaceFiles(rootPaths) {
		fmt.Printf("Generating: %s\n", path)
		switch meta.Type {
		case "msg":
			result, err := g.generateMessage(&meta, path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", path, destPath, err)
				break
			}
			set := getOrDefault(result.GoPackage())
			set.AddFrom(result.CImports)
		case "srv":
			result, err := g.generateService(&meta, path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Service '%s' to '%s', error: %v\n", path, destPath, err)
				break
			}
			set := getOrDefault(result.GoPackage())
			set.AddFrom(result.Request.CImports)
			set.AddFrom(result.Response.CImports)
		case "action":
			result, err := g.generateAction(path, destPath)
			if err != nil {
				fmt.Printf("Error converting ROS2 Action '%s' to '%s', error: %v\n", path, destPath, err)
			}
			set := getOrDefault(result.GoPackage())
			set.AddFrom(result.Goal.CImports)
			set.AddFrom(result.SendGoal.Request.CImports)
			set.AddFrom(result.SendGoal.Response.CImports)
			set.AddFrom(result.Result.CImports)
			set.AddFrom(result.GetResult.Request.CImports)
			set.AddFrom(result.GetResult.Response.CImports)
			set.AddFrom(result.Feedback.CImports)
			set.AddFrom(result.FeedbackMessage.CImports)
		default:
			fmt.Printf("Interface file %s has an invalid type: %s\n", path, meta.Type)
		}
	}
	for pkg, imports := range cImports {
		err := g.generateCommonPackageGoFile(pkg, imports, destPath)
		if err != nil {
			fmt.Printf("Failed to generate common package file for package %s: %v", pkg, err)
		}
	}
	return nil
}

func findInterfaceFiles(rootPaths []string) map[Metadata]string {
	files := make(map[Metadata]string)
	for i := len(rootPaths) - 1; i >= 0; i-- {
		filepath.Walk(rootPaths[i], func(path string, info fs.FileInfo, err error) error {
			skip, blacklistEntry := blacklisted(path)
			if skip {
				fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
				return nil
			}
			if re.M(filepath.ToSlash(path), `m!/(msg/.+\.msg)|(srv/.+\.srv)|(action/.+\.action)$!`) {
				md, err := parseMetadataFromPath(path)
				if err == nil {
					files[*md] = path
				} else {
					fmt.Printf("Failed to parse metadata from path %s: %v\n", path, err)
				}
			}
			return nil
		})
	}
	return files
}

func (g *generator) generateMessage(md *Metadata, sourcePath string, destPathPkgRoot string) (*ROS2Message, error) {
	msg := ROS2MessageNew("", "")
	var err error

	msg.Metadata = md

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
	dirs := strings.Split(p, string(filepath.Separator))

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

func (g *generator) generateService(m *Metadata, srcPath, dstPathPkgRoot string) (*ROS2Service, error) {
	service := NewROS2Service(m.Package, m.Name)
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	parser := parser{config: g.config}
	if err = parser.ParseService(service, string(srcFile)); err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, dstPathPkgRoot, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (g *generator) generateAction(srcPath, dstPathPkgRoot string) (*ROS2Action, error) {
	m, err := parseMetadataFromPath(srcPath)
	if err != nil {
		return nil, err
	}
	action := NewROS2Action(m.Package, m.Name)
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	parser := parser{config: g.config}
	if err = parser.ParseAction(action, string(srcFile)); err != nil {
		return nil, err
	}
	err = g.generateGoFile(
		dstPathPkgRoot,
		action.Metadata,
		ros2ActionToGolangTypeTemplate,
		templateData{"Action": action},
	)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, dstPathPkgRoot, action.Goal)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, dstPathPkgRoot, action.Result)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, dstPathPkgRoot, action.Feedback)
	if err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, dstPathPkgRoot, action.SendGoal)
	if err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, dstPathPkgRoot, action.GetResult)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, dstPathPkgRoot, action.FeedbackMessage)
	if err != nil {
		return nil, err
	}
	return action, nil
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
