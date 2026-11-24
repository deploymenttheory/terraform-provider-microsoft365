package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Kinds of non-standard function. Only helpers are subject to the helper rule.
const (
	kindHelper    = "helper"
	kindKiota     = "kiota"
	kindFramework = "framework"
	// kindStateMapper is the package's required Graph-to-state mapper. Its name and file are checked; the caller and line rules don't apply.
	kindStateMapper = "state-mapper"
)

// Standard names of the package's Graph-to-state mapper, for resource and data source packages.
const (
	stateMapperName           = "MapRemoteStateToTerraform"
	dataSourceStateMapperName = "MapRemoteStateToDataSource"
)

// dataSourcesPath marks a data source package: any package under internal/services/datasources/.
const dataSourcesPath = "/internal/services/datasources/"

// stateMapperFor returns the standard state mapper name for the package in dir.
func stateMapperFor(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	if strings.Contains(filepath.ToSlash(abs)+"/", dataSourcesPath) {
		return dataSourceStateMapperName
	}
	return stateMapperName
}

// crudMethods are the resource methods whose direct calls are counted.
var crudMethods = map[string]bool{"Create": true, "Read": true, "Update": true, "Delete": true, "ImportState": true, "ModifyPlan": true}

// standardFuncs are left out of the report entirely: CRUD, schema and metadata methods and constructResource.
var standardFuncs = map[string]bool{
	"Create": true, "Read": true, "Update": true, "Delete": true, "ImportState": true, "ModifyPlan": true,
	"Metadata": true, "Schema": true, "Configure": true, "ValidateConfig": true, "ConfigValidators": true, "UpgradeState": true, "IdentitySchema": true, "Exists": true,
	"constructResource": true,
}

// kiotaMethods implement the Kiota Parsable, AdditionalDataHolder and BackedModel interfaces.
var kiotaMethods = map[string]bool{
	"Serialize": true, "GetFieldDeserializers": true,
	"GetAdditionalData": true, "SetAdditionalData": true,
	"GetBackingStore": true, "SetBackingStore": true,
}

// kiotaFactory matches the parse-node factories Kiota models register for discriminator lookups.
var kiotaFactory = regexp.MustCompile(`^[cC]reate\w*FromDiscriminatorValue$`)

// frameworkMethod matches the methods that schema validators and plan modifiers must implement.
var frameworkMethod = regexp.MustCompile(`^(Description|MarkdownDescription|Validate[A-Z]\w*|PlanModify[A-Z]\w*)$`)

// Func is one non-standard function found in a package.
type Func struct {
	Package     string
	File        string
	Name        string
	Receiver    string
	Kind        string
	Lines       int
	Refs        int
	CRUDCallers []string
	NameCheck   string
}

// Rule returns "pass" when a helper meets all three parts of the helper rule, or the parts it misses.
func (f Func) Rule() string {
	if f.Kind == kindStateMapper {
		if f.NameCheck != nameOK {
			return "fails: name"
		}
		return "exempt"
	}
	if f.Kind != kindHelper {
		return "excluded"
	}
	var misses []string
	if len(f.CRUDCallers) < minCRUDCallers {
		misses = append(misses, "callers")
	}
	if f.Lines <= maxShortLines {
		misses = append(misses, "lines")
	}
	if f.NameCheck != nameOK {
		misses = append(misses, "name")
	}
	if len(misses) == 0 {
		return "pass"
	}
	return "fails: " + strings.Join(misses, ",")
}

const (
	// minCRUDCallers is the number of distinct CRUD methods that must call a helper directly.
	minCRUDCallers = 2
	// maxShortLines is the declaration span at or below which a helper counts as short.
	maxShortLines = 10
)

// scanDir parses the non-test Go files in dir and returns its non-standard functions in file and name order.
func scanDir(dir string) ([]Func, error) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	byPackage := map[string][]*ast.File{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", filepath.Join(dir, name), err)
		}
		byPackage[file.Name.Name] = append(byPackage[file.Name.Name], file)
	}
	var funcs []Func
	mapper := stateMapperFor(dir)
	for _, files := range byPackage {
		funcs = append(funcs, scanPackage(fset, filepath.Base(dir), files, mapper)...)
	}
	sort.Slice(funcs, func(i, j int) bool {
		if funcs[i].File != funcs[j].File {
			return funcs[i].File < funcs[j].File
		}
		if funcs[i].Name != funcs[j].Name {
			return funcs[i].Name < funcs[j].Name
		}
		return funcs[i].Receiver < funcs[j].Receiver
	})
	return funcs, nil
}

// scanPackage type-checks one package without its imports, so that references to its own functions and methods resolve exactly.
func scanPackage(fset *token.FileSet, pkgName string, files []*ast.File, mapper string) []Func {
	info := &types.Info{Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}}
	conf := types.Config{Importer: stubImporter{}, Error: func(error) {}}
	// Imports are stubbed, so type errors are expected; the Error handler keeps checking going and they are ignored.
	_, _ = conf.Check(pkgName, fset, files, info)

	type decl struct {
		fd   *ast.FuncDecl
		file string
		obj  types.Object
	}
	var decls []decl
	index := map[types.Object]int{}
	for _, file := range files {
		fileName := filepath.Base(fset.Position(file.Pos()).Filename)
		for _, d := range file.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			decls = append(decls, decl{fd: fd, file: fileName, obj: info.Defs[fd.Name]})
			if obj := info.Defs[fd.Name]; obj != nil {
				index[obj] = len(decls) - 1
			}
		}
	}

	refs := make([]int, len(decls))
	callers := make([]map[string]bool, len(decls))
	for i, d := range decls {
		// Each CRUD method is identified by receiver type and name, so two resources in one package stay distinct.
		caller := ""
		if d.fd.Recv != nil && crudMethods[d.fd.Name.Name] {
			caller = receiverName(d.fd) + "." + d.fd.Name.Name
		}
		ast.Inspect(d.fd.Body, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			obj := info.Uses[ident]
			if fn, isFunc := obj.(*types.Func); isFunc {
				obj = fn.Origin()
			}
			target, found := index[obj]
			if !found || target == i {
				return true
			}
			refs[target]++
			if caller != "" {
				if callers[target] == nil {
					callers[target] = map[string]bool{}
				}
				callers[target][caller] = true
			}
			return true
		})
	}

	var funcs []Func
	for i, d := range decls {
		kind, report := classify(d.fd, mapper)
		if !report {
			continue
		}
		f := Func{Package: pkgName, File: d.file, Name: d.fd.Name.Name, Receiver: receiverName(d.fd), Kind: kind, Refs: refs[i]}
		f.Lines = fset.Position(d.fd.End()).Line - fset.Position(d.fd.Pos()).Line + 1
		for caller := range callers[i] {
			f.CRUDCallers = append(f.CRUDCallers, caller)
		}
		sort.Strings(f.CRUDCallers)
		switch kind {
		case kindHelper:
			f.NameCheck = checkName(f.Name, f.File, mapper)
		case kindStateMapper:
			f.NameCheck = checkStateMapperFile(f.Name, f.File)
		}
		funcs = append(funcs, f)
	}
	return funcs
}

// classify reports whether a function belongs in the report and, if so, its kind. mapper is the package's standard state mapper name.
func classify(fd *ast.FuncDecl, mapper string) (string, bool) {
	name := fd.Name.Name
	if standardFuncs[name] || (fd.Recv == nil && strings.HasPrefix(name, "New")) || name == "init" || name == "main" {
		return "", false
	}
	if fd.Recv == nil && name == mapper {
		return kindStateMapper, true
	}
	if (fd.Recv != nil && kiotaMethods[name]) || (fd.Recv == nil && kiotaFactory.MatchString(name)) {
		return kindKiota, true
	}
	if fd.Recv != nil && frameworkMethod.MatchString(name) {
		return kindFramework, true
	}
	return kindHelper, true
}

// receiverName returns the receiver's type name, without pointer or type parameters, or "" for a plain function.
func receiverName(fd *ast.FuncDecl) string {
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return ""
	}
	expr := fd.Recv.List[0].Type
	for {
		switch e := expr.(type) {
		case *ast.StarExpr:
			expr = e.X
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		case *ast.ParenExpr:
			expr = e.X
		case *ast.Ident:
			return e.Name
		default:
			return ""
		}
	}
}

// stubImporter returns an empty package for every import, so the scanner needs no module downloads or builds.
type stubImporter struct{}

// Import returns an empty, complete package named after the last path element.
func (stubImporter) Import(path string) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}
	pkg := types.NewPackage(path, filepath.Base(path))
	pkg.MarkComplete()
	return pkg, nil
}
