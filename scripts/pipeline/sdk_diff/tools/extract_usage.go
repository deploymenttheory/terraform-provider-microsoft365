package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// UsageMap tracks SDK usage across the codebase
type UsageMap struct {
	Packages map[string]int              `json:"packages"`
	Imports  map[string][]string         `json:"imports"`
	Types    map[string]map[string]int   `json:"types"`
	Methods  map[string]int              `json:"methods"`
	Fields   map[string]map[string]int   `json:"fields"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <repo-path>\n", os.Args[0])
		os.Exit(1)
	}

	repoPath := os.Args[1]
	usage := &UsageMap{
		Packages: make(map[string]int),
		Imports:  make(map[string][]string),
		Types:    make(map[string]map[string]int),
		Methods:  make(map[string]int),
		Fields:   make(map[string]map[string]int),
	}

	// Walk all Go files in internal/services
	servicesPath := filepath.Join(repoPath, "internal", "services")
	err := filepath.Walk(servicesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip test files and vendor
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return analyzeFile(path, usage)
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing files: %v\n", err)
		os.Exit(1)
	}

	// Output JSON
	output, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func analyzeFile(path string, usage *UsageMap) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		// Skip files with parse errors (e.g., templates)
		return nil
	}

	// Track current file's imports
	fileImports := make(map[string]string) // alias -> full path

	// First pass: collect imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		
		// Only track msgraph SDK imports
		if !strings.Contains(importPath, "microsoftgraph") && !strings.Contains(importPath, "kiota") {
			continue
		}

		// Get alias (or last part of path)
		alias := ""
		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(importPath, "/")
			alias = parts[len(parts)-1]
		}

		fileImports[alias] = importPath
		usage.Packages[importPath]++
		
		// Track which files import this package
		if usage.Imports[importPath] == nil {
			usage.Imports[importPath] = []string{}
		}
		usage.Imports[importPath] = append(usage.Imports[importPath], path)
	}

	// Track type assignments to infer field accesses
	// Maps variable names to their SDK types
	varTypes := make(map[string]string)

	// Second pass: analyze usage
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			// Track variable assignments to infer types
			// e.g., user := models.NewUser() or user := &models.User{}
			for i, rhs := range x.Rhs {
				if i >= len(x.Lhs) {
					break
				}
				
				// Get variable name
				var varName string
				if ident, ok := x.Lhs[i].(*ast.Ident); ok {
					varName = ident.Name
				}
				
				if varName == "" {
					continue
				}
				
				// Check if RHS is SDK type
				switch rhsType := rhs.(type) {
				case *ast.CallExpr:
					// Function call: models.NewUser()
					if sel, ok := rhsType.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := sel.X.(*ast.Ident); ok {
							if importPath, exists := fileImports[ident.Name]; exists {
								varTypes[varName] = fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
							}
						}
					}
				case *ast.UnaryExpr:
					// Pointer: &models.User{}
					if rhsType.Op == token.AND {
						if comp, ok := rhsType.X.(*ast.CompositeLit); ok {
							if sel, ok := comp.Type.(*ast.SelectorExpr); ok {
								if ident, ok := sel.X.(*ast.Ident); ok {
									if importPath, exists := fileImports[ident.Name]; exists {
										varTypes[varName] = fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
									}
								}
							}
						}
					}
				case *ast.CompositeLit:
					// Struct literal: models.User{}
					if sel, ok := rhsType.Type.(*ast.SelectorExpr); ok {
						if ident, ok := sel.X.(*ast.Ident); ok {
							if importPath, exists := fileImports[ident.Name]; exists {
								varTypes[varName] = fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
							}
						}
					}
				}
			}

		case *ast.SelectorExpr:
			// Could be: package.Type, object.Field, or object.Method
			if ident, ok := x.X.(*ast.Ident); ok {
				// Check if this is a package reference
				if importPath, exists := fileImports[ident.Name]; exists {
					// This is package.Something (Type or function)
					fullName := fmt.Sprintf("%s.%s", importPath, x.Sel.Name)
					usage.Methods[fullName]++
				} else {
					// This is object.Field - check if we know the object's type
					if typeName, exists := varTypes[ident.Name]; exists {
						// Track field access
						if usage.Fields[typeName] == nil {
							usage.Fields[typeName] = make(map[string]int)
						}
						usage.Fields[typeName][x.Sel.Name]++
					}
				}
			}

		case *ast.CallExpr:
			// Method calls
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if importPath, exists := fileImports[ident.Name]; exists {
						// Package-level function call
						methodName := fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
						usage.Methods[methodName]++
					} else if typeName, exists := varTypes[ident.Name]; exists {
						// Method call on typed object
						methodName := fmt.Sprintf("%s.%s", typeName, sel.Sel.Name)
						usage.Methods[methodName]++
					}
				}
			}

		case *ast.CompositeLit:
			// Struct literal instantiation: models.User{DisplayName: "foo"}
			if sel, ok := x.Type.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if importPath, exists := fileImports[ident.Name]; exists {
						typeName := fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
						
						// Track the type itself
						if usage.Types[typeName] == nil {
							usage.Types[typeName] = make(map[string]int)
						}
						usage.Types[typeName]["_instantiated"]++

						// Track fields used in struct literal
						if usage.Fields[typeName] == nil {
							usage.Fields[typeName] = make(map[string]int)
						}
						
						for _, elt := range x.Elts {
							if kv, ok := elt.(*ast.KeyValueExpr); ok {
								if fieldIdent, ok := kv.Key.(*ast.Ident); ok {
									usage.Fields[typeName][fieldIdent.Name]++
								}
							}
						}
					}
				}
			}
		}

		return true
	})

	return nil
}
