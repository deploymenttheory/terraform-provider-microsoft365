package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// UsageMap tracks SDK usage across the codebase
type UsageMap struct {
	Packages map[string]int            `json:"packages"`
	Imports  map[string][]string       `json:"imports"`
	Types    map[string]map[string]int `json:"types"`
	Methods  map[string]int            `json:"methods"`
	Fields   map[string]map[string]int `json:"fields"`
	Enums    map[string][]string       `json:"enums"`
}

// initializeUsageMap creates a new UsageMap with initialized maps.
func initializeUsageMap() *UsageMap {
	return &UsageMap{
		Packages: make(map[string]int),
		Imports:  make(map[string][]string),
		Types:    make(map[string]map[string]int),
		Methods:  make(map[string]int),
		Fields:   make(map[string]map[string]int),
		Enums:    make(map[string][]string),
	}
}

// main orchestrates SDK usage extraction from repository.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <repo-path>\n", os.Args[0])
		os.Exit(1)
	}

	usage := initializeUsageMap()

	if err := analyzeRepository(os.Args[1], usage); err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing files: %v\n", err)
		os.Exit(1)
	}

	if err := outputResults(usage); err != nil {
		fmt.Fprintf(os.Stderr, "Error outputting results: %v\n", err)
		os.Exit(1)
	}
}

// analyzeRepository walks the repository directory tree and analyzes all Go files.
func analyzeRepository(repoPath string, usage *UsageMap) error {
	servicesPath := filepath.Join(repoPath, "internal", "services")
	return filepath.Walk(servicesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipFile(path) {
			return nil
		}

		return analyzeFile(path, usage)
	})
}

// shouldSkipFile returns true if the file should be excluded from analysis.
func shouldSkipFile(path string) bool {
	return !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")
}

// outputResults serializes usage data to JSON and writes to stdout.
func outputResults(usage *UsageMap) error {
	output, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// analyzeFile parses a Go file and extracts SDK usage patterns.
func analyzeFile(path string, usage *UsageMap) error {
	node, err := parseGoFile(path)
	if err != nil {
		return nil // Skip files with parse errors
	}

	fileImports := collectImports(node, path, usage)
	varTypes := make(map[string]string)

	analyzeAST(node, path, fileImports, varTypes, usage)
	return nil
}

// parseGoFile parses a Go source file and returns its AST.
func parseGoFile(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, path, nil, parser.ParseComments)
}

// collectImports extracts SDK imports from an AST and tracks them in usage map.
// Returns a map of import aliases to full import paths for the file.
func collectImports(node *ast.File, path string, usage *UsageMap) map[string]string {
	fileImports := make(map[string]string)

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)

		if !isSDKImport(importPath) {
			continue
		}

		alias := getImportAlias(imp, importPath)
		fileImports[alias] = importPath

		trackImport(importPath, path, usage)
	}

	return fileImports
}

// isSDKImport checks if an import path is from Microsoft Graph SDK or Kiota.
func isSDKImport(importPath string) bool {
	return strings.Contains(importPath, "microsoftgraph") || strings.Contains(importPath, "kiota")
}

// getImportAlias returns the alias for an import (explicit or inferred from path).
func getImportAlias(imp *ast.ImportSpec, importPath string) string {
	if imp.Name != nil {
		return imp.Name.Name
	}
	parts := strings.Split(importPath, "/")
	return parts[len(parts)-1]
}

// trackImport records an SDK import and the file that imports it.
func trackImport(importPath, filePath string, usage *UsageMap) {
	usage.Packages[importPath]++

	if usage.Imports[importPath] == nil {
		usage.Imports[importPath] = []string{}
	}
	usage.Imports[importPath] = append(usage.Imports[importPath], filePath)
}

// analyzeAST walks the AST and processes different node types to extract SDK usage.
func analyzeAST(node *ast.File, path string, fileImports map[string]string, varTypes map[string]string, usage *UsageMap) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			processAssignments(x, fileImports, varTypes)
		case *ast.SelectorExpr:
			processSelectorExpr(x, fileImports, varTypes, usage)
		case *ast.CallExpr:
			processCallExpr(x, path, fileImports, varTypes, usage)
		case *ast.CompositeLit:
			processCompositeLit(x, fileImports, usage)
		}
		return true
	})
}

// processAssignments tracks variable assignments to infer SDK types.
// Example: user := models.NewUser()
func processAssignments(stmt *ast.AssignStmt, fileImports map[string]string, varTypes map[string]string) {
	for i, rhs := range stmt.Rhs {
		if i >= len(stmt.Lhs) {
			break
		}

		varName := extractVarName(stmt.Lhs[i])
		if varName == "" {
			continue
		}

		if sdkType := extractSDKType(rhs, fileImports); sdkType != "" {
			varTypes[varName] = sdkType
		}
	}
}

// extractVarName extracts variable name from an expression (if it's an identifier).
func extractVarName(expr ast.Expr) string {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// extractSDKType determines if an expression is an SDK type and returns its full name.
func extractSDKType(rhs ast.Expr, fileImports map[string]string) string {
	switch rhsType := rhs.(type) {
	case *ast.CallExpr:
		return extractTypeFromCallExpr(rhsType, fileImports)
	case *ast.UnaryExpr:
		return extractTypeFromUnaryExpr(rhsType, fileImports)
	case *ast.CompositeLit:
		return extractTypeFromCompositeLit(rhsType, fileImports)
	}
	return ""
}

// extractTypeFromCallExpr extracts SDK type from function call expression.
// Example: models.NewUser() -> "github.com/.../models.User"
func extractTypeFromCallExpr(call *ast.CallExpr, fileImports map[string]string) string {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if importPath, exists := fileImports[ident.Name]; exists {
				return fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
			}
		}
	}
	return ""
}

// extractTypeFromUnaryExpr extracts SDK type from unary expression (typically pointers).
// Example: &models.User{} -> "github.com/.../models.User"
func extractTypeFromUnaryExpr(unary *ast.UnaryExpr, fileImports map[string]string) string {
	if unary.Op == token.AND {
		if comp, ok := unary.X.(*ast.CompositeLit); ok {
			return extractTypeFromCompositeLit(comp, fileImports)
		}
	}
	return ""
}

// extractTypeFromCompositeLit extracts SDK type from composite literal.
// Example: models.User{} -> "github.com/.../models.User"
func extractTypeFromCompositeLit(comp *ast.CompositeLit, fileImports map[string]string) string {
	if sel, ok := comp.Type.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if importPath, exists := fileImports[ident.Name]; exists {
				return fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
			}
		}
	}
	return ""
}

// processSelectorExpr handles selector expressions (package.Type or object.Field).
func processSelectorExpr(sel *ast.SelectorExpr, fileImports map[string]string, varTypes map[string]string, usage *UsageMap) {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if importPath, exists := fileImports[ident.Name]; exists {
			// Package-level reference
			fullName := fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
			usage.Methods[fullName]++
		} else if typeName, exists := varTypes[ident.Name]; exists {
			// Field access on typed object
			trackFieldAccess(typeName, sel.Sel.Name, usage)
		}
	}
}

// trackFieldAccess records a field access on an SDK type.
func trackFieldAccess(typeName, fieldName string, usage *UsageMap) {
	if usage.Fields[typeName] == nil {
		usage.Fields[typeName] = make(map[string]int)
	}
	usage.Fields[typeName][fieldName]++
}

// processCallExpr handles function/method calls and tracks enum parser usage.
func processCallExpr(call *ast.CallExpr, path string, fileImports map[string]string, varTypes map[string]string, usage *UsageMap) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	if importPath, exists := fileImports[ident.Name]; exists {
		trackPackageMethod(importPath, sel.Sel.Name, path, usage)
	} else if typeName, exists := varTypes[ident.Name]; exists {
		trackObjectMethod(typeName, sel.Sel.Name, usage)
	}
}

// trackPackageMethod records a package-level method call and checks for enum parsers.
func trackPackageMethod(importPath, methodName, path string, usage *UsageMap) {
	fullMethodName := fmt.Sprintf("%s.%s", importPath, methodName)
	usage.Methods[fullMethodName]++

	// Track enum parser usage
	if strings.HasPrefix(methodName, "Parse") && strings.Contains(importPath, "models") {
		trackEnumUsage(importPath, methodName, path, usage)
	}
}

// trackObjectMethod records a method call on a typed object.
func trackObjectMethod(typeName, methodName string, usage *UsageMap) {
	fullMethodName := fmt.Sprintf("%s.%s", typeName, methodName)
	usage.Methods[fullMethodName]++
}

// trackEnumUsage detects and records enum parser calls.
// Example: ParseRunAsAccountType -> tracks RunAsAccountType enum
func trackEnumUsage(importPath, methodName, path string, usage *UsageMap) {
	enumType := strings.TrimPrefix(methodName, "Parse")
	if enumType == "" {
		return
	}

	fullEnumName := fmt.Sprintf("%s.%s", importPath, enumType)
	if !slices.Contains(usage.Enums[fullEnumName], path) {
		usage.Enums[fullEnumName] = append(usage.Enums[fullEnumName], path)
	}
}

// processCompositeLit handles struct literal instantiation and field usage.
// Example: models.User{DisplayName: "foo"}
func processCompositeLit(comp *ast.CompositeLit, fileImports map[string]string, usage *UsageMap) {
	sel, ok := comp.Type.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	importPath, exists := fileImports[ident.Name]
	if !exists {
		return
	}

	typeName := fmt.Sprintf("%s.%s", importPath, sel.Sel.Name)
	trackTypeInstantiation(typeName, usage)
	trackStructFields(typeName, comp.Elts, usage)
}

// trackTypeInstantiation records that an SDK type was instantiated.
func trackTypeInstantiation(typeName string, usage *UsageMap) {
	if usage.Types[typeName] == nil {
		usage.Types[typeName] = make(map[string]int)
	}
	usage.Types[typeName]["_instantiated"]++
}

// trackStructFields extracts and records fields used in struct literals.
func trackStructFields(typeName string, elts []ast.Expr, usage *UsageMap) {
	if usage.Fields[typeName] == nil {
		usage.Fields[typeName] = make(map[string]int)
	}

	for _, elt := range elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if fieldIdent, ok := kv.Key.(*ast.Ident); ok {
				usage.Fields[typeName][fieldIdent.Name]++
			}
		}
	}
}
