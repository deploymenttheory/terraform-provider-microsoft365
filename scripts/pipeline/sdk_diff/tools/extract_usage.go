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

// ResourceInfo tracks SDK usage for a single Terraform entity
type ResourceInfo struct {
	ResourcePath    string          `json:"resource_path"`
	SDKDependencies SDKDependencies `json:"sdk_dependencies"`
	Files           []string        `json:"files"`
}

// SDKDependencies tracks what SDK components an entity uses
type SDKDependencies struct {
	Types         []string            `json:"types"`
	FieldsUsed    map[string][]string `json:"fields_used"`
	MethodsCalled []string            `json:"methods_called"`
	EnumsUsed     []EnumUsage         `json:"enums_used"`
}

// EnumUsage tracks enum usage
type EnumUsage struct {
	Enum         string   `json:"enum"`
	ValuesInUse  []string `json:"values_in_use,omitempty"`
}

// Statistics provides summary statistics
type Statistics struct {
	TotalResources      int `json:"total_resources"`
	TotalActions        int `json:"total_actions"`
	TotalListActions    int `json:"total_list_actions"`
	TotalEphemerals     int `json:"total_ephemerals"`
	TotalDataSources    int `json:"total_data_sources"`
	TotalSDKTypesUsed   int `json:"total_sdk_types_used"`
	TotalSDKMethodsUsed int `json:"total_sdk_methods_used"`
	TotalEnumsTracked   int `json:"total_enums_tracked"`
}

// UsageMapV2 is the new resource-centric structure
type UsageMapV2 struct {
	TerraformResources   map[string]*ResourceInfo `json:"terraform_resources"`
	TerraformActions     map[string]*ResourceInfo `json:"terraform_actions"`
	TerraformListActions map[string]*ResourceInfo `json:"terraform_list_actions"`
	TerraformEphemerals  map[string]*ResourceInfo `json:"terraform_ephemerals"`
	TerraformDataSources map[string]*ResourceInfo `json:"terraform_data_sources"`
	SDKToResourceIndex   map[string][]string      `json:"sdk_to_resource_index"`
	Statistics           Statistics               `json:"statistics"`
}

// main orchestrates SDK usage extraction from repository.
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <repo-path>\n", os.Args[0])
		os.Exit(1)
	}

	usage := initializeUsageMapV2()

	if err := analyzeRepository(os.Args[1], usage); err != nil {
		fmt.Fprintf(os.Stderr, "Error analyzing files: %v\n", err)
		os.Exit(1)
	}

	if err := outputResults(usage); err != nil {
		fmt.Fprintf(os.Stderr, "Error outputting results: %v\n", err)
		os.Exit(1)
	}
}

// initializeUsageMapV2 creates a new UsageMapV2 with initialized maps.
func initializeUsageMapV2() *UsageMapV2 {
	return &UsageMapV2{
		TerraformResources:   make(map[string]*ResourceInfo),
		TerraformActions:     make(map[string]*ResourceInfo),
		TerraformListActions: make(map[string]*ResourceInfo),
		TerraformEphemerals:  make(map[string]*ResourceInfo),
		TerraformDataSources: make(map[string]*ResourceInfo),
		SDKToResourceIndex:   make(map[string][]string),
		Statistics:           Statistics{},
	}
}

// analyzeRepository walks the repository directory tree and analyzes all Go files.
func analyzeRepository(repoPath string, usage *UsageMapV2) error {
	servicesPath := filepath.Join(repoPath, "internal", "services")
	return filepath.Walk(servicesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldSkipFile(path) {
			return nil
		}

		return analyzeFile(path, repoPath, usage)
	})
}

// shouldSkipFile returns true if the file should be excluded from analysis.
func shouldSkipFile(path string) bool {
	return !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")
}

// outputResults serializes usage data to JSON and writes to stdout.
func outputResults(usage *UsageMapV2) error {
	// Finalize statistics
	usage.Statistics = calculateStatistics(usage)

	output, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// analyzeFile parses a Go file and extracts SDK usage patterns.
func analyzeFile(path, repoPath string, usage *UsageMapV2) error {
	node, err := parseGoFile(path)
	if err != nil {
		return nil // Skip files with parse errors
	}

	// Determine entity type and name from file path
	entity := parseEntityFromPath(path, repoPath)
	if entity == nil {
		return nil // Not in a tracked directory
	}

	// Get or create ResourceInfo for this entity
	resourceInfo := getOrCreateResourceInfo(usage, entity)

	// Track this file
	fileName := filepath.Base(path)
	if !slices.Contains(resourceInfo.Files, fileName) {
		resourceInfo.Files = append(resourceInfo.Files, fileName)
	}

	// Collect imports
	fileImports := collectImports(node)
	if len(fileImports) == 0 {
		return nil // No SDK imports
	}

	// Analyze AST and track SDK usage
	analyzeASTForEntity(node, path, fileImports, resourceInfo, usage)

	return nil
}

// Entity represents a Terraform entity (resource, action, data source, etc.)
type Entity struct {
	Type string // "resource", "action", "list-action", "ephemeral", "data-source"
	Name string // e.g., "microsoft365_user"
	Path string // relative path from repo root
}

// parseEntityFromPath extracts entity information from file path
func parseEntityFromPath(path, repoPath string) *Entity {
	relPath := strings.TrimPrefix(path, repoPath+"/")

	// Check each entity type
	if strings.HasPrefix(relPath, "internal/services/resources/") {
		return extractEntityInfo(relPath, "internal/services/resources/", "resource")
	}
	if strings.HasPrefix(relPath, "internal/services/actions/") {
		return extractEntityInfo(relPath, "internal/services/actions/", "action")
	}
	if strings.HasPrefix(relPath, "internal/services/list-resources/") {
		return extractEntityInfo(relPath, "internal/services/list-resources/", "list-action")
	}
	if strings.HasPrefix(relPath, "internal/services/ephemerals/") {
		return extractEntityInfo(relPath, "internal/services/ephemerals/", "ephemeral")
	}
	if strings.HasPrefix(relPath, "internal/services/data-sources/") {
		return extractEntityInfo(relPath, "internal/services/data-sources/", "data-source")
	}

	return nil
}

// extractEntityInfo extracts entity name from path
func extractEntityInfo(relPath, prefix, entityType string) *Entity {
	// Remove prefix
	remainder := strings.TrimPrefix(relPath, prefix)
	
	// Split into parts: domain/graph_version/resource_name/file.go
	parts := strings.Split(remainder, "/")
	if len(parts) < 3 {
		return nil
	}

	// parts[0] = service domain (e.g., "device_management", "users")
	// parts[1] = graph version (e.g., "graph_beta", "graph_v1")
	// parts[2] = resource/action name
	resourceName := parts[2]

	// Convert to Terraform naming: microsoft365_{resource_name}
	tfName := "microsoft365_" + strings.ReplaceAll(resourceName, "_", "_")

	return &Entity{
		Type: entityType,
		Name: tfName,
		Path: prefix + strings.Join(parts[:3], "/"),
	}
}

// getOrCreateResourceInfo gets or creates a ResourceInfo for an entity
func getOrCreateResourceInfo(usage *UsageMapV2, entity *Entity) *ResourceInfo {
	var targetMap map[string]*ResourceInfo

	switch entity.Type {
	case "resource":
		targetMap = usage.TerraformResources
	case "action":
		targetMap = usage.TerraformActions
	case "list-action":
		targetMap = usage.TerraformListActions
	case "ephemeral":
		targetMap = usage.TerraformEphemerals
	case "data-source":
		targetMap = usage.TerraformDataSources
	default:
		return nil
	}

	if targetMap[entity.Name] == nil {
		targetMap[entity.Name] = &ResourceInfo{
			ResourcePath: entity.Path,
			SDKDependencies: SDKDependencies{
				Types:         []string{},
				FieldsUsed:    make(map[string][]string),
				MethodsCalled: []string{},
				EnumsUsed:     []EnumUsage{},
			},
			Files: []string{},
		}
	}

	return targetMap[entity.Name]
}

// parseGoFile parses a Go source file and returns its AST.
func parseGoFile(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, path, nil, parser.ParseComments)
}

// collectImports extracts SDK imports from an AST.
// Returns a map of import aliases to full import paths for the file.
func collectImports(node *ast.File) map[string]string {
	fileImports := make(map[string]string)

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)

		if !isSDKImport(importPath) {
			continue
		}

		alias := getImportAlias(imp, importPath)
		fileImports[alias] = importPath
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

// analyzeASTForEntity walks the AST and extracts SDK usage for an entity
func analyzeASTForEntity(node *ast.File, path string, fileImports map[string]string, resourceInfo *ResourceInfo, usage *UsageMapV2) {
	varTypes := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			processAssignments(x, fileImports, varTypes)
		case *ast.SelectorExpr:
			processSelectorExpr(x, fileImports, varTypes, resourceInfo)
		case *ast.CallExpr:
			processCallExpr(x, path, fileImports, varTypes, resourceInfo, usage)
		case *ast.CompositeLit:
			processCompositeLit(x, fileImports, resourceInfo)
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

		sdkType := extractSDKType(rhs, fileImports)
		if sdkType != "" {
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
	switch expr := rhs.(type) {
	case *ast.CallExpr:
		return extractTypeFromCallExpr(expr, fileImports)
	case *ast.CompositeLit:
		return extractTypeFromCompositeLit(expr, fileImports)
	case *ast.UnaryExpr:
		return extractTypeFromUnaryExpr(expr, fileImports)
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
func processSelectorExpr(sel *ast.SelectorExpr, fileImports map[string]string, varTypes map[string]string, resourceInfo *ResourceInfo) {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if _, exists := fileImports[ident.Name]; exists {
			// Package-level reference - don't track here, tracked in calls
		} else if typeName, exists := varTypes[ident.Name]; exists {
			// Field access on typed object
			trackFieldAccess(typeName, sel.Sel.Name, resourceInfo)
		}
	}
}

// trackFieldAccess records a field access on an SDK type.
func trackFieldAccess(typeName, fieldName string, resourceInfo *ResourceInfo) {
	// Simplify type name for display
	shortType := simplifyTypeName(typeName)
	
	if resourceInfo.SDKDependencies.FieldsUsed[shortType] == nil {
		resourceInfo.SDKDependencies.FieldsUsed[shortType] = []string{}
	}
	if !slices.Contains(resourceInfo.SDKDependencies.FieldsUsed[shortType], fieldName) {
		resourceInfo.SDKDependencies.FieldsUsed[shortType] = append(resourceInfo.SDKDependencies.FieldsUsed[shortType], fieldName)
	}
}

// processCallExpr handles function/method calls and tracks enum parser usage.
func processCallExpr(call *ast.CallExpr, path string, fileImports map[string]string, varTypes map[string]string, resourceInfo *ResourceInfo, usage *UsageMapV2) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	if importPath, exists := fileImports[ident.Name]; exists {
		trackPackageMethod(importPath, sel.Sel.Name, path, resourceInfo, usage)
	} else if typeName, exists := varTypes[ident.Name]; exists {
		trackObjectMethod(typeName, sel.Sel.Name, resourceInfo)
	}
}

// trackPackageMethod records a package-level method call and checks for enum parsers.
func trackPackageMethod(importPath, methodName, path string, resourceInfo *ResourceInfo, usage *UsageMapV2) {
	fullMethodName := fmt.Sprintf("%s.%s", simplifyTypeName(importPath), methodName)
	
	if !slices.Contains(resourceInfo.SDKDependencies.MethodsCalled, fullMethodName) {
		resourceInfo.SDKDependencies.MethodsCalled = append(resourceInfo.SDKDependencies.MethodsCalled, fullMethodName)
	}

	// Track enum parser usage
	if strings.HasPrefix(methodName, "Parse") && strings.Contains(importPath, "models") {
		trackEnumUsage(importPath, methodName, resourceInfo, usage)
	}
}

// trackObjectMethod records a method call on a typed object.
func trackObjectMethod(typeName, methodName string, resourceInfo *ResourceInfo) {
	fullMethodName := fmt.Sprintf("%s.%s", simplifyTypeName(typeName), methodName)
	
	if !slices.Contains(resourceInfo.SDKDependencies.MethodsCalled, fullMethodName) {
		resourceInfo.SDKDependencies.MethodsCalled = append(resourceInfo.SDKDependencies.MethodsCalled, fullMethodName)
	}
}

// trackEnumUsage detects and records enum parser calls.
// Example: ParseRunAsAccountType -> tracks RunAsAccountType enum
func trackEnumUsage(importPath, methodName string, resourceInfo *ResourceInfo, usage *UsageMapV2) {
	enumType := strings.TrimPrefix(methodName, "Parse")
	if enumType == "" {
		return
	}

	fullEnumName := fmt.Sprintf("%s.%s", simplifyTypeName(importPath), enumType)
	
	// Add to resource's enum list
	found := false
	for _, existing := range resourceInfo.SDKDependencies.EnumsUsed {
		if existing.Enum == fullEnumName {
			found = true
			break
		}
	}
	if !found {
		resourceInfo.SDKDependencies.EnumsUsed = append(resourceInfo.SDKDependencies.EnumsUsed, EnumUsage{
			Enum: fullEnumName,
		})
	}
}

// processCompositeLit handles struct literal instantiation and field usage.
// Example: models.User{DisplayName: "foo"}
func processCompositeLit(comp *ast.CompositeLit, fileImports map[string]string, resourceInfo *ResourceInfo) {
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
	trackTypeInstantiation(typeName, resourceInfo)
	trackStructFields(typeName, comp.Elts, resourceInfo)
}

// trackTypeInstantiation records that an SDK type was instantiated.
func trackTypeInstantiation(typeName string, resourceInfo *ResourceInfo) {
	shortType := simplifyTypeName(typeName)
	
	if !slices.Contains(resourceInfo.SDKDependencies.Types, shortType) {
		resourceInfo.SDKDependencies.Types = append(resourceInfo.SDKDependencies.Types, shortType)
	}
}

// trackStructFields extracts and records fields used in struct literals.
func trackStructFields(typeName string, elts []ast.Expr, resourceInfo *ResourceInfo) {
	shortType := simplifyTypeName(typeName)
	
	for _, elt := range elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		fieldIdent, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if resourceInfo.SDKDependencies.FieldsUsed[shortType] == nil {
			resourceInfo.SDKDependencies.FieldsUsed[shortType] = []string{}
		}
		if !slices.Contains(resourceInfo.SDKDependencies.FieldsUsed[shortType], fieldIdent.Name) {
			resourceInfo.SDKDependencies.FieldsUsed[shortType] = append(resourceInfo.SDKDependencies.FieldsUsed[shortType], fieldIdent.Name)
		}
	}
}

// simplifyTypeName converts full SDK type names to short versions
func simplifyTypeName(fullName string) string {
	// Remove the base path, keep only models.TypeName or package.TypeName
	if strings.Contains(fullName, "/models.") {
		parts := strings.Split(fullName, "/models.")
		return "models." + parts[len(parts)-1]
	}
	if strings.Contains(fullName, "/") {
		parts := strings.Split(fullName, "/")
		return parts[len(parts)-1]
	}
	return fullName
}

// calculateStatistics generates final statistics and SDK-to-resource index
func calculateStatistics(usage *UsageMapV2) Statistics {
	// Build SDK-to-resource index
	allTypesUsed := make(map[string]bool)
	allMethodsUsed := make(map[string]bool)
	allEnumsUsed := make(map[string]bool)

	for entityName, info := range usage.TerraformResources {
		indexSDKUsage(entityName, info, usage.SDKToResourceIndex, allTypesUsed, allMethodsUsed, allEnumsUsed)
	}
	for entityName, info := range usage.TerraformActions {
		indexSDKUsage(entityName, info, usage.SDKToResourceIndex, allTypesUsed, allMethodsUsed, allEnumsUsed)
	}
	for entityName, info := range usage.TerraformListActions {
		indexSDKUsage(entityName, info, usage.SDKToResourceIndex, allTypesUsed, allMethodsUsed, allEnumsUsed)
	}
	for entityName, info := range usage.TerraformEphemerals {
		indexSDKUsage(entityName, info, usage.SDKToResourceIndex, allTypesUsed, allMethodsUsed, allEnumsUsed)
	}
	for entityName, info := range usage.TerraformDataSources {
		indexSDKUsage(entityName, info, usage.SDKToResourceIndex, allTypesUsed, allMethodsUsed, allEnumsUsed)
	}

	return Statistics{
		TotalResources:      len(usage.TerraformResources),
		TotalActions:        len(usage.TerraformActions),
		TotalListActions:    len(usage.TerraformListActions),
		TotalEphemerals:     len(usage.TerraformEphemerals),
		TotalDataSources:    len(usage.TerraformDataSources),
		TotalSDKTypesUsed:   len(allTypesUsed),
		TotalSDKMethodsUsed: len(allMethodsUsed),
		TotalEnumsTracked:   len(allEnumsUsed),
	}
}

// indexSDKUsage builds the reverse index from SDK components to entities
func indexSDKUsage(entityName string, info *ResourceInfo, index map[string][]string, typesSet, methodsSet, enumsSet map[string]bool) {
	// Index types
	for _, sdkType := range info.SDKDependencies.Types {
		if !slices.Contains(index[sdkType], entityName) {
			index[sdkType] = append(index[sdkType], entityName)
		}
		typesSet[sdkType] = true
	}

	// Index methods
	for _, method := range info.SDKDependencies.MethodsCalled {
		methodsSet[method] = true
	}

	// Index enums
	for _, enumUsage := range info.SDKDependencies.EnumsUsed {
		if !slices.Contains(index[enumUsage.Enum], entityName) {
			index[enumUsage.Enum] = append(index[enumUsage.Enum], entityName)
		}
		enumsSet[enumUsage.Enum] = true
	}
}
