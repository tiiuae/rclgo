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
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/kivilahtio/go-re/v0"
)

type Rule struct {
	Pattern *regexp.Regexp
}

func NewRule(pattern string) (*Rule, error) {
	re, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}
	return &Rule{Pattern: re}, nil
}

type RuleSet []*Rule

func (s RuleSet) Includes(str string) bool {
	for _, r := range s {
		if r.Pattern.MatchString(str) {
			return true
		}
	}
	return false
}

type Config struct {
	RclgoImportPath     string
	MessageModulePrefix string
	RootPaths           []string

	RegexIncludes RuleSet
	PkgIncludes   []string
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
	config        *Config
	destRootPath  string
	cImportsByPkg map[string]stringSet
	allPkgs       map[string]*rosPkgRef
}

func GeneratePrimitives(c *Config, destFilePath string) error {
	return generateRclgoFile(
		"primitive types",
		filepath.Join(destFilePath, "pkg/rclgo/primitives/primitives.gen.go"),
		primitiveTypes,
		templateData{
			"PMap":   &primitiveTypeMappings,
			"Config": c,
		},
	)
}

func GenerateRclgoFlags(c *Config, destFilePath string) error {
	return generateRclgoFile(
		"rclgo flags",
		filepath.Join(destFilePath, "pkg/rclgo/flags.gen.go"),
		rclgoFlags,
		templateData{"Config": c},
	)
}

func GenerateTestGogenFlags(c *Config, destFilePath string) error {
	return generateRclgoFile(
		"gogen flags",
		filepath.Join(destFilePath, "test/gogen/flags.gen.go"),
		gogenTestFlags,
		templateData{"Config": c},
	)
}

func generateRclgoFile(fileType, destFilePath string, tmpl *template.Template, data templateData) error {
	destFile, err := mkdir_p(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	fmt.Printf("Generating %s: %s\n", fileType, destFilePath)
	return tmpl.Execute(destFile, data)
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

type rosPkgRef struct {
	Interfaces map[Metadata]string
	Generated  bool
}

func GenerateGolangMessageTypes(config *Config, destPath string) error {
	gen := generator{
		config:        config,
		destRootPath:  destPath,
		cImportsByPkg: make(map[string]stringSet),
	}
	gen.findPackages(config.RootPaths)
	for _, pkg := range config.PkgIncludes {
		gen.generatePkg(gen.allPkgs[pkg], true)
	}
	for pkg, ref := range gen.allPkgs {
		if config.RegexIncludes.Includes(pkg) {
			gen.generatePkg(ref, false)
		}
	}
	for pkg, imports := range gen.cImportsByPkg {
		err := gen.generateCommonPackageGoFile(pkg, imports)
		if err != nil {
			fmt.Printf("Failed to generate common package file for package %s: %v", pkg, err)
		}
	}
	return nil
}

func (g *generator) generatePkg(ref *rosPkgRef, genDeps bool) {
	if ref != nil && !ref.Generated {
		ref.Generated = true
		for meta, path := range ref.Interfaces {
			cImports := g.generateInterface(meta, path)
			if genDeps {
				for imp := range cImports {
					g.generatePkg(g.allPkgs[imp], genDeps)
				}
			}
		}
	}
}

func (g *generator) getCImportsForPkg(pkg string) stringSet {
	set := g.cImportsByPkg[pkg]
	if set == nil {
		set = stringSet{}
		g.cImportsByPkg[pkg] = set
	}
	return set
}

func (g *generator) generateInterface(meta Metadata, ifacePath string) stringSet {
	fmt.Printf("Generating: %s\n", ifacePath)
	switch meta.Type {
	case "msg":
		result, err := g.generateMessage(&meta, ifacePath)
		if err != nil {
			fmt.Printf("Error converting ROS2 Message '%s' to '%s', error: %v\n", ifacePath, g.destRootPath, err)
			return nil
		}
		set := g.getCImportsForPkg(result.GoPackage())
		set.AddFrom(result.CImports)
		return set
	case "srv":
		result, err := g.generateService(&meta, ifacePath)
		if err != nil {
			fmt.Printf("Error converting ROS2 Service '%s' to '%s', error: %v\n", ifacePath, g.destRootPath, err)
			return nil
		}
		set := g.getCImportsForPkg(result.GoPackage())
		set.AddFrom(result.Request.CImports)
		set.AddFrom(result.Response.CImports)
		return set
	case "action":
		result, err := g.generateAction(ifacePath)
		if err != nil {
			fmt.Printf("Error converting ROS2 Action '%s' to '%s', error: %v\n", ifacePath, g.destRootPath, err)
			return nil
		}
		set := g.getCImportsForPkg(result.GoPackage())
		set.AddFrom(result.Goal.CImports)
		set.AddFrom(result.SendGoal.Request.CImports)
		set.AddFrom(result.SendGoal.Response.CImports)
		set.AddFrom(result.Result.CImports)
		set.AddFrom(result.GetResult.Request.CImports)
		set.AddFrom(result.GetResult.Response.CImports)
		set.AddFrom(result.Feedback.CImports)
		set.AddFrom(result.FeedbackMessage.CImports)
		return set
	default:
		fmt.Printf("Interface file %s has an invalid type: %s\n", ifacePath, meta.Type)
		return nil
	}
}

func (g *generator) findPackages(rootPaths []string) {
	g.allPkgs = map[string]*rosPkgRef{}
	for i := len(rootPaths) - 1; i >= 0; i-- {
		filepath.Walk(rootPaths[i], func(path string, info fs.FileInfo, err error) error { //nolint:errcheck
			skip, blacklistEntry := blacklisted(path)
			if skip {
				fmt.Printf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
				return nil
			}
			if re.M(filepath.ToSlash(path), `m!/(msg/.+\.msg)|(srv/.+\.srv)|(action/.+\.action)$!`) {
				meta, err := parseMetadataFromPath(path)
				if err != nil {
					fmt.Printf("Failed to parse metadata from path %s: %v\n", path, err)
				} else {
					ref := g.allPkgs[meta.Package]
					if ref == nil {
						ref = &rosPkgRef{Interfaces: map[Metadata]string{}}
						g.allPkgs[meta.Package] = ref
					}
					ref.Interfaces[*meta] = path
				}
			}
			return nil
		})
	}
}

func (g *generator) generateMessage(md *Metadata, sourcePath string) (*ROS2Message, error) {
	msg := ROS2MessageNew("", "")
	var err error

	msg.Metadata = md

	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	parser := parser{config: g.config}
	err = parser.ParseROS2Message(msg, string(content))
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, msg)
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

func (g *generator) generateService(m *Metadata, srcPath string) (*ROS2Service, error) {
	service := NewROS2Service(m.Package, m.Name)
	srcFile, err := os.ReadFile(filepath.Clean(srcPath))
	if err != nil {
		return nil, err
	}
	parser := parser{config: g.config}
	if err = parser.ParseService(service, string(srcFile)); err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, service)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (g *generator) generateAction(srcPath string) (*ROS2Action, error) {
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
		action.Metadata,
		ros2ActionToGolangTypeTemplate,
		templateData{"Action": action},
	)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, action.Goal)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, action.Result)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, action.Feedback)
	if err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, action.SendGoal)
	if err != nil {
		return nil, err
	}
	err = g.generateServiceGoFiles(&parser, action.GetResult)
	if err != nil {
		return nil, err
	}
	err = g.generateMessageGoFile(&parser, action.FeedbackMessage)
	if err != nil {
		return nil, err
	}
	return action, nil
}

type templateData = map[string]interface{}

func (g *generator) generateGoFile(meta *Metadata, tmpl *template.Template, data templateData) error {
	f, err := createTargetGolangTypeFile(g.destRootPath, meta)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, ok := data["Config"]; !ok {
		data["Config"] = g.config
	}
	return tmpl.Execute(f, data)
}

func (g *generator) generateMessageGoFile(parser *parser, msg *ROS2Message) error {
	return g.generateGoFile(
		msg.Metadata,
		ros2MsgToGolangTypeTemplate,
		templateData{
			"Message":             msg,
			"cSerializationCode":  parser.cSerializationCode,
			"goSerializationCode": parser.goSerializationCode,
		},
	)
}

func (g *generator) generateServiceGoFiles(parser *parser, srv *ROS2Service) error {
	err := g.generateGoFile(
		srv.Metadata,
		ros2ServiceToGolangTypeTemplate,
		templateData{"Service": srv},
	)
	if err != nil {
		return err
	}
	err = g.generateMessageGoFile(parser, srv.Request)
	if err != nil {
		return err
	}
	return g.generateMessageGoFile(parser, srv.Response)
}

func (g *generator) generateCommonPackageGoFile(goPkg string, cImports stringSet) error {
	i := strings.LastIndex(goPkg, "_")
	if i < 0 || i > len(goPkg)-1 {
		return errors.New("package type suffix is missing or incorrect")
	}
	cPkg := goPkg[:i]
	pkgType := goPkg[i+1:]
	f, err := mkdir_p(filepath.Join(g.destRootPath, cPkg, pkgType, "common.gen.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	return ros2PackageCommonTemplate.Execute(f, templateData{
		"GoPackage": goPkg,
		"CPackage":  cPkg,
		"CImports":  cImports,
		"RootPaths": g.config.RootPaths,
	})
}
