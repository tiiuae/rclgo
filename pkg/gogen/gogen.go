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
	DestPath            string

	RegexIncludes  RuleSet
	ROSPkgIncludes []string
	GoPkgIncludes  []string

	LicenseHeader string
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

type Generator struct {
	config               *Config
	cImportsByPkgAndType map[string]stringSet
	allPkgs              map[string]*rosPkgRef
	actionMsgsNeeded     bool
}

func New(config *Config) *Generator {
	return &Generator{
		config:               config,
		cImportsByPkgAndType: make(map[string]stringSet),
	}
}

func (g *Generator) GeneratePrimitives() error {
	return g.generateRclgoFile(
		"primitive types",
		filepath.Join(g.config.DestPath, "pkg/rclgo/primitives/primitives.gen.go"),
		primitiveTypes,
		templateData{"PMap": &primitiveTypeMappings},
	)
}

func (g *Generator) GenerateRclgoFlags() error {
	return g.generateRclgoFile(
		"rclgo flags",
		filepath.Join(g.config.DestPath, "pkg/rclgo/flags.gen.go"),
		rclgoFlags,
		nil,
	)
}

func (g *Generator) GenerateTestGogenFlags() error {
	return g.generateRclgoFile(
		"gogen flags",
		filepath.Join(g.config.DestPath, "test/gogen/flags.gen.go"),
		gogenTestFlags,
		nil,
	)
}

func (g *Generator) generateRclgoFile(fileType, destFilePath string, tmpl *template.Template, data templateData) error {
	PrintErrf("Generating %s: %s\n", fileType, destFilePath)
	return g.generateGoFile(destFilePath, tmpl, data)
}

func (g *Generator) GenerateROS2AllMessagesImporter() error {
	pkgs := map[string]struct{}{}
	for _, glob := range []string{"*/msg", "*/srv", "*/action"} {
		dirs, err := filepath.Glob(filepath.Join(g.config.DestPath, glob))
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
	return g.generateGoFile(
		filepath.Join(g.config.DestPath, "msgs.gen.go"),
		ros2MsgImportAllPackage,
		templateData{"Packages": pkgs},
	)
}

type rosPkgRef struct {
	Interfaces map[Metadata]string
	Generated  bool
}

func (g *Generator) GenerateGolangMessageTypes() error {
	g.findPackages()
	if len(g.config.RegexIncludes) == 0 && len(g.config.ROSPkgIncludes) == 0 && len(g.config.GoPkgIncludes) == 0 {
		for pkg := range g.allPkgs {
			g.generatePkg(pkg, false)
		}
	} else {
		for _, pkg := range g.config.ROSPkgIncludes {
			g.generatePkg(pkg, true)
		}
		goDeps, err := loadGoPkgDeps(g.config.GoPkgIncludes...)
		if err != nil {
			return fmt.Errorf("failed to load Go deps: %w", err)
		}
		prefix := g.config.MessageModulePrefix + "/"
		for goDep := range goDeps {
			pkgWithType := strings.TrimPrefix(goDep, prefix)
			if pkgWithType != goDep {
				g.generatePkg(path.Dir(pkgWithType), true)
			}
		}
		if g.actionMsgsNeeded {
			g.generatePkg("action_msgs", true)
		}
		for pkg := range g.allPkgs {
			if g.config.RegexIncludes.Includes(pkg) {
				g.generatePkg(pkg, false)
			}
		}
	}
	for pkgAndType, imports := range g.cImportsByPkgAndType {
		err := g.generateCommonPackageGoFile(pkgAndType, imports)
		if err != nil {
			PrintErrf("Failed to generate common package file for package %s: %v\n", pkgAndType, err)
		}
	}
	return nil
}

func (g *Generator) generatePkg(pkg string, genDeps bool) {
	ref := g.allPkgs[pkg]
	if ref == nil {
		PrintErrf("Failed to generate package %s: package not found\n", pkg)
	} else if !ref.Generated {
		ref.Generated = true
		for meta, path := range ref.Interfaces {
			g.generateInterface(meta, path)
		}
		if genDeps {
			for imp := range g.cImportsByPkgAndType[pkg+"_msg"] {
				g.generatePkg(imp, genDeps)
			}
			for imp := range g.cImportsByPkgAndType[pkg+"_srv"] {
				g.generatePkg(imp, genDeps)
			}
			for imp := range g.cImportsByPkgAndType[pkg+"_action"] {
				g.generatePkg(imp, genDeps)
			}
		}
	}
}

func (g *Generator) getCImportsForPkgAndType(pkgAndType string) stringSet {
	set := g.cImportsByPkgAndType[pkgAndType]
	if set == nil {
		set = stringSet{}
		g.cImportsByPkgAndType[pkgAndType] = set
	}
	return set
}

func (g *Generator) generateInterface(meta Metadata, ifacePath string) {
	PrintErrf("Generating: %s\n", ifacePath)
	switch meta.Type {
	case "msg":
		result, err := g.generateMessage(&meta, ifacePath)
		if err != nil {
			PrintErrf("Error converting ROS2 Message '%s' to '%s', error: %v\n", ifacePath, g.config.DestPath, err)
		}
		set := g.getCImportsForPkgAndType(result.GoPackage())
		set.AddFrom(result.CImports)
	case "srv":
		result, err := g.generateService(&meta, ifacePath)
		if err != nil {
			PrintErrf("Error converting ROS2 Service '%s' to '%s', error: %v\n", ifacePath, g.config.DestPath, err)
		}
		set := g.getCImportsForPkgAndType(result.GoPackage())
		set.AddFrom(result.Request.CImports)
		set.AddFrom(result.Response.CImports)
	case "action":
		result, err := g.generateAction(ifacePath)
		if err != nil {
			PrintErrf("Error converting ROS2 Action '%s' to '%s', error: %v\n", ifacePath, g.config.DestPath, err)
		}
		g.actionMsgsNeeded = true
		set := g.getCImportsForPkgAndType(result.GoPackage())
		set.AddFrom(result.Goal.CImports)
		set.AddFrom(result.SendGoal.Request.CImports)
		set.AddFrom(result.SendGoal.Response.CImports)
		set.AddFrom(result.Result.CImports)
		set.AddFrom(result.GetResult.Request.CImports)
		set.AddFrom(result.GetResult.Response.CImports)
		set.AddFrom(result.Feedback.CImports)
		set.AddFrom(result.FeedbackMessage.CImports)
	default:
		PrintErrf("Interface file %s has an invalid type: %s\n", ifacePath, meta.Type)
	}
}

func (g *Generator) findPackages() {
	g.allPkgs = map[string]*rosPkgRef{}
	for i := len(g.config.RootPaths) - 1; i >= 0; i-- {
		filepath.Walk(g.config.RootPaths[i], func(path string, info fs.FileInfo, err error) error { //nolint:errcheck
			skip, blacklistEntry := blacklisted(path)
			if skip {
				PrintErrf("Blacklisted: %s, matched regex '%s'\n", path, blacklistEntry)
				return nil
			}
			if re.M(filepath.ToSlash(path), `m!/(msg/.+\.msg)|(srv/.+\.srv)|(action/.+\.action)$!`) {
				meta, err := parseMetadataFromPath(path)
				if err != nil {
					PrintErrf("Failed to parse metadata from path %s: %v\n", path, err)
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

func (g *Generator) generateMessage(md *Metadata, sourcePath string) (*ROS2Message, error) {
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

func ifaceFilePath(destPathPkgRoot string, m *Metadata) string {
	return filepath.Join(destPathPkgRoot, m.ImportPath(), m.Name+".gen.go")
}

func (g *Generator) generateService(m *Metadata, srcPath string) (*ROS2Service, error) {
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

func (g *Generator) generateAction(srcPath string) (*ROS2Action, error) {
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
	err = g.generateIfaceGoFile(
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

func (g *Generator) generateGoFile(destPath string, tmpl *template.Template, data templateData) error {
	f, err := mkdir_p(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if data == nil {
		data = templateData{"Config": g.config}
	} else if _, ok := data["Config"]; !ok {
		data["Config"] = g.config
	}
	return tmpl.Execute(f, data)
}

func (g *Generator) generateIfaceGoFile(meta *Metadata, tmpl *template.Template, data templateData) error {
	return g.generateGoFile(ifaceFilePath(g.config.DestPath, meta), tmpl, data)
}

func (g *Generator) generateMessageGoFile(parser *parser, msg *ROS2Message) error {
	return g.generateIfaceGoFile(
		msg.Metadata,
		ros2MsgToGolangTypeTemplate,
		templateData{
			"Message":             msg,
			"cSerializationCode":  parser.cSerializationCode,
			"goSerializationCode": parser.goSerializationCode,
		},
	)
}

func (g *Generator) generateServiceGoFiles(parser *parser, srv *ROS2Service) error {
	err := g.generateIfaceGoFile(
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

func (g *Generator) generateCommonPackageGoFile(pkgAndType string, cImports stringSet) error {
	i := strings.LastIndex(pkgAndType, "_")
	if i < 0 || i > len(pkgAndType)-1 {
		return errors.New("package type suffix is missing or incorrect")
	}
	cPkg := pkgAndType[:i]
	pkgType := pkgAndType[i+1:]
	return g.generateGoFile(
		filepath.Join(g.config.DestPath, cPkg, pkgType, "common.gen.go"),
		ros2PackageCommonTemplate,
		templateData{
			"GoPackage": pkgAndType,
			"CPackage":  cPkg,
			"CImports":  cImports,
		},
	)
}
