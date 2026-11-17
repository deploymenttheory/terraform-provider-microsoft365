// Command declarations identifies Go files containing only imports, types, and constants.
package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
)

func declarationOnly(source string) bool {
	file, err := parser.ParseFile(token.NewFileSet(), "coverage.go", source, parser.SkipObjectResolution)
	if err != nil {
		return false
	}
	for _, spec := range file.Imports {
		if spec.Path.Value == `"C"` || (spec.Name != nil && spec.Name.Name == "_") {
			return false
		}
	}
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || (general.Tok != token.IMPORT && general.Tok != token.TYPE && general.Tok != token.CONST) {
			return false
		}
	}
	return true
}

func main() {
	var sources map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&sources); err != nil {
		os.Exit(1)
	}
	verified := make([]string, 0)
	for name, source := range sources {
		if declarationOnly(source) {
			verified = append(verified, name)
		}
	}
	sort.Strings(verified)
	if err := json.NewEncoder(os.Stdout).Encode(verified); err != nil {
		os.Exit(1)
	}
}
