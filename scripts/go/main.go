package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// EndpointIndex represents the complete API endpoint index
type EndpointIndex struct {
	Metadata  Metadata           `json:"metadata"`
	Endpoints []ResourceEndpoint `json:"endpoints"`
}

// Metadata contains summary information about the index
type Metadata struct {
	TotalEndpoints int               `json:"total_endpoints"`
	ResourceTypes  map[string]int    `json:"resource_types"`
	SDKVersions    map[string]string `json:"sdk_versions"`
}

// ResourceEndpoint represents a Terraform resource with its API endpoints
type ResourceEndpoint struct {
	ResourceType      string               `json:"resource_type"`
	ServiceDomain     string               `json:"service_domain"`
	APIVersion        string               `json:"api_version"`
	ResourceName      string               `json:"resource_name"`
	TerraformTypeName string               `json:"terraform_type_name"`
	Operations        map[string][]APICall `json:"operations"`
	GraphModelType    string               `json:"graph_model_type,omitempty"`
	TerraformFields   []TerraformField     `json:"terraform_fields,omitempty"`
}

// APICall represents a single API call
type APICall struct {
	SDKMethodChain string `json:"sdk_method_chain"`
	HTTPMethod     string `json:"http_method"`
	APIURLPattern  string `json:"api_url_pattern"`
	LineNumber     int    `json:"line_number"`
	FilePath       string `json:"file_path"`
	GraphModelType string `json:"graph_model_type,omitempty"`
}

// TerraformField represents a field used in the Terraform resource
type TerraformField struct {
	TerraformName  string `json:"terraform_name"`
	GraphFieldName string `json:"graph_field_name"`
	GraphFieldType string `json:"graph_field_type"`
	SetterMethod   string `json:"setter_method"`
	IsRequired     bool   `json:"is_required,omitempty"`
	IsComputed     bool   `json:"is_computed,omitempty"`
}

// Indexer is the main indexer struct
type Indexer struct {
	providerPath string
	fset         *token.FileSet
	endpoints    []ResourceEndpoint
	sdkVersions  map[string]string
	sdkURLCache  map[string]string // Cache of SDK method to URL mappings
}

// NewIndexer creates a new indexer
func NewIndexer(providerPath string) *Indexer {
	return &Indexer{
		providerPath: providerPath,
		fset:         token.NewFileSet(),
		endpoints:    []ResourceEndpoint{},
		sdkVersions:  make(map[string]string),
		sdkURLCache:  make(map[string]string),
	}
}

// Index performs the indexing
func (idx *Indexer) Index() (*EndpointIndex, error) {
	log.Println("Starting indexing...")

	// Parse SDK versions from go.mod
	if err := idx.parseGoMod(); err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	// Index resources
	if err := idx.indexResources(); err != nil {
		return nil, fmt.Errorf("failed to index resources: %w", err)
	}

	// Index datasources
	if err := idx.indexDatasources(); err != nil {
		return nil, fmt.Errorf("failed to index datasources: %w", err)
	}

	// Index actions
	if err := idx.indexActions(); err != nil {
		return nil, fmt.Errorf("failed to index actions: %w", err)
	}

	// Build metadata
	metadata := idx.buildMetadata()

	return &EndpointIndex{
		Metadata:  metadata,
		Endpoints: idx.endpoints,
	}, nil
}

// parseGoMod extracts SDK versions from go.mod
func (idx *Indexer) parseGoMod() error {
	goModPath := filepath.Join(idx.providerPath, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Extract msgraph SDK versions
	betaPattern := regexp.MustCompile(`github\.com/microsoftgraph/msgraph-beta-sdk-go\s+(v[\d.]+)`)
	v1Pattern := regexp.MustCompile(`github\.com/microsoftgraph/msgraph-sdk-go\s+(v[\d.]+)`)

	if match := betaPattern.FindStringSubmatch(content); len(match) > 1 {
		idx.sdkVersions["msgraph-beta-sdk-go"] = match[1]
	}

	if match := v1Pattern.FindStringSubmatch(content); len(match) > 1 {
		idx.sdkVersions["msgraph-sdk-go"] = match[1]
	}

	log.Printf("Found SDK versions: %v\n", idx.sdkVersions)
	return nil
}

// indexResources indexes all resources
func (idx *Indexer) indexResources() error {
	resourcesPath := filepath.Join(idx.providerPath, "internal", "services", "resources")
	return idx.walkServiceDirectory(resourcesPath, "resource")
}

// indexDatasources indexes all datasources
func (idx *Indexer) indexDatasources() error {
	datasourcesPath := filepath.Join(idx.providerPath, "internal", "services", "datasources")
	return idx.walkServiceDirectory(datasourcesPath, "datasource")
}

// indexActions indexes all actions
func (idx *Indexer) indexActions() error {
	actionsPath := filepath.Join(idx.providerPath, "internal", "services", "actions")
	return idx.walkServiceDirectory(actionsPath, "action")
}

// walkServiceDirectory walks through service directories
func (idx *Indexer) walkServiceDirectory(basePath, resourceType string) error {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s\n", basePath)
		return nil
	}

	log.Printf("Indexing %ss...\n", resourceType)

	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for crud.go files for resources, read.go or datasource.go for datasources,
		// and invoke.go for actions
		fileName := info.Name()
		var shouldProcess bool

		switch resourceType {
		case "resource":
			shouldProcess = fileName == "crud.go"
		case "datasource":
			shouldProcess = fileName == "read.go" || fileName == "datasource.go"
		case "action":
			shouldProcess = fileName == "invoke.go"
		}

		if !shouldProcess {
			return nil
		}

		// Extract resource information from path
		relPath, _ := filepath.Rel(basePath, filepath.Dir(path))
		parts := strings.Split(relPath, string(os.PathSeparator))

		if len(parts) < 2 {
			return nil
		}

		serviceDomain := parts[0]
		apiVersion := parts[1]
		resourceName := strings.Join(parts[2:], "/")

		log.Printf("  Processing %s: %s/%s/%s\n", resourceType, serviceDomain, apiVersion, resourceName)

		// Process the file
		endpoint, err := idx.processFile(path, resourceType, serviceDomain, apiVersion, resourceName)
		if err != nil {
			log.Printf("    Warning: Error processing %s: %v\n", path, err)
			return nil
		}

		if endpoint != nil {
			idx.endpoints = append(idx.endpoints, *endpoint)
		}

		return nil
	})
}

// processFile processes a single Go file
func (idx *Indexer) processFile(filePath, resourceType, serviceDomain, apiVersion, resourceName string) (*ResourceEndpoint, error) {
	// Parse the file
	node, err := parser.ParseFile(idx.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	endpoint := &ResourceEndpoint{
		ResourceType:  resourceType,
		ServiceDomain: serviceDomain,
		APIVersion:    apiVersion,
		ResourceName:  resourceName,
		Operations:    make(map[string][]APICall),
	}

	// Extract terraform type name
	endpoint.TerraformTypeName = idx.extractTerraformTypeName(node, filepath.Dir(filePath))

	// Extract graph model type
	endpoint.GraphModelType = idx.extractGraphModelType(node)

	// Extract API calls
	visitor := &apiCallVisitor{
		fset:       idx.fset,
		filePath:   filePath,
		apiVersion: apiVersion,
		calls:      make(map[string][]APICall),
		indexer:    idx,
	}

	ast.Walk(visitor, node)
	endpoint.Operations = visitor.calls

	// Extract Terraform fields from construct function
	endpoint.TerraformFields = idx.extractTerraformFields(filepath.Dir(filePath), endpoint.GraphModelType)

	// Only return if we found operations
	if len(endpoint.Operations) > 0 {
		return endpoint, nil
	}

	return nil, nil
}

// extractTerraformTypeName extracts the Terraform type name from constants
func (idx *Indexer) extractTerraformTypeName(node *ast.File, dirPath string) string {
	// Look for constants like ResourceName, datasourceName, ActionName
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				nameStr := name.Name
				if nameStr == "ResourceName" || nameStr == "datasourceName" || nameStr == "ActionName" {
					if len(valueSpec.Values) > i {
						if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
							typeName := strings.Trim(lit.Value, `"`)
							if !strings.HasPrefix(typeName, "microsoft365_") {
								typeName = "microsoft365_" + typeName
							}
							return typeName
						}
					}
				}
			}
		}
	}

	// Fallback: use directory name
	return "microsoft365_" + filepath.Base(dirPath)
}

// extractGraphModelType extracts the Graph model type from imports and function signatures
func (idx *Indexer) extractGraphModelType(node *ast.File) string {
	// Look for function calls like NewDeviceHealthScript() or similar constructors
	// This is a simplified version - could be enhanced

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Look for construct functions
		if strings.HasPrefix(funcDecl.Name.Name, "construct") || funcDecl.Name.Name == "Create" {
			// Try to find model type from return type or usage
			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				// Look for New* calls
				if callExpr, ok := n.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if strings.HasPrefix(selExpr.Sel.Name, "New") {
							return false // Found it, but we need more context
						}
					}
				}
				return true
			})
		}
	}

	return ""
}

// extractTerraformFields extracts Terraform field definitions by analyzing construct.go
func (idx *Indexer) extractTerraformFields(dirPath, graphModelType string) []TerraformField {
	fields := []TerraformField{}

	// Look for construct.go file
	constructPath := filepath.Join(dirPath, "construct.go")
	if _, err := os.Stat(constructPath); os.IsNotExist(err) {
		return fields
	}

	// Parse the construct file
	node, err := parser.ParseFile(idx.fset, constructPath, nil, parser.ParseComments)
	if err != nil {
		return fields
	}

	// Extract model type from constructResource function if not provided
	if graphModelType == "" {
		graphModelType = idx.extractGraphModelFromConstruct(node)
	}

	// Extract field mappings from convert calls
	fieldMappings := idx.extractFieldMappingsFromConstruct(node)

	// Extract nested object type mappings (e.g., Provider -> AgentProvider)
	nestedObjectTypes := idx.extractNestedObjectTypes(node)

	// Get Graph API field names from SDK model
	if graphModelType != "" {
		fields = idx.enrichFieldsWithSDKModel(fieldMappings, graphModelType, nestedObjectTypes)
	} else {
		// Fallback: just use the mappings we found
		for tfName, setterMethod := range fieldMappings {
			fields = append(fields, TerraformField{
				TerraformName: tfName,
				SetterMethod:  setterMethod,
			})
		}
	}

	return fields
}

// extractGraphModelFromConstruct extracts the Graph model type from construct function
func (idx *Indexer) extractGraphModelFromConstruct(node *ast.File) string {
	var modelType string

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for patterns like: requestBody := graphmodels.NewAgentCollection()
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			if len(assignStmt.Rhs) > 0 {
				if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						methodName := selExpr.Sel.Name
						if strings.HasPrefix(methodName, "New") {
							// Extract model name from NewXxx()
							modelType = strings.TrimPrefix(methodName, "New")
							return false
						}
					}
				}
			}
		}
		return true
	})

	return modelType
}

// extractFieldMappingsFromConstruct extracts field mappings from convert function calls
func (idx *Indexer) extractFieldMappingsFromConstruct(node *ast.File) map[string]string {
	mappings := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for convert function calls like:
		// convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
		// Also handle nested fields like: data.Provider.Organization, provider.SetOrganization
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				// Check if this is a convert function
				if strings.HasPrefix(selExpr.Sel.Name, "FrameworkToGraph") {
					if len(callExpr.Args) >= 2 {
						// First arg: data.FieldName or data.NestedObj.FieldName
						if firstArg, ok := callExpr.Args[0].(*ast.SelectorExpr); ok {
							tfFieldName := firstArg.Sel.Name

							// Check if this is a nested field (e.g., data.Provider.Organization)
							var parentFieldName string
							if innerSel, ok := firstArg.X.(*ast.SelectorExpr); ok {
								// This is a nested field access
								parentFieldName = innerSel.Sel.Name
								// Convert each part to snake_case separately, then combine
								parentSnake := toSnakeCase(parentFieldName)
								childSnake := toSnakeCase(tfFieldName)
								tfFieldName = parentSnake + "_" + childSnake
							} else {
								// Simple field, convert to snake_case
								tfFieldName = toSnakeCase(tfFieldName)
							}

							// Second arg: requestBody.SetFieldName or nestedObj.SetFieldName
							if secondArg, ok := callExpr.Args[1].(*ast.SelectorExpr); ok {
								setterMethod := secondArg.Sel.Name
								mappings[tfFieldName] = setterMethod
							}
						}
					}
				}
			}
		}
		return true
	})

	return mappings
}

// extractNestedObjectTypes extracts nested object types from construct function
// Returns a map of field name (in snake_case) to nested model type (e.g., "provider" -> "AgentProvider")
func (idx *Indexer) extractNestedObjectTypes(node *ast.File) map[string]string {
	nestedTypes := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for assignments like: provider := graphmodels.NewAgentProvider()
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) > 0 && len(assignStmt.Rhs) > 0 {
				// Get variable name (e.g., "provider")
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					varName := ident.Name

					// Get the constructor call (e.g., graphmodels.NewAgentProvider())
					if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
						if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							methodName := selExpr.Sel.Name
							if strings.HasPrefix(methodName, "New") {
								// Extract model type (e.g., "NewAgentProvider" -> "AgentProvider")
								modelType := strings.TrimPrefix(methodName, "New")

								// Convert variable name to snake_case
								snakeName := toSnakeCase(varName)
								nestedTypes[snakeName] = modelType
							}
						}
					}
				}
			}
		}
		return true
	})

	return nestedTypes
}

// enrichFieldsWithSDKModel enriches field data with SDK model information
func (idx *Indexer) enrichFieldsWithSDKModel(fieldMappings map[string]string, modelType string, nestedObjectTypes map[string]string) []TerraformField {
	fields := []TerraformField{}

	// Find SDK model file
	sdkModelPath := idx.findSDKModelFile(modelType)
	if sdkModelPath == "" {
		// Fallback: return basic mappings
		for tfName, setterMethod := range fieldMappings {
			fields = append(fields, TerraformField{
				TerraformName: tfName,
				SetterMethod:  setterMethod,
			})
		}
		return fields
	}

	// Parse SDK model file to get JSON field names for the main model
	jsonFieldNames := idx.extractJSONFieldNamesFromSDK(sdkModelPath)

	// Also get JSON field names for nested object types
	nestedJsonFields := make(map[string]map[string]jsonFieldInfo)
	for nestedFieldName, nestedModelType := range nestedObjectTypes {
		nestedSDKPath := idx.findSDKModelFile(nestedModelType)
		if nestedSDKPath != "" {
			nestedJsonFields[nestedFieldName] = idx.extractJSONFieldNamesFromSDK(nestedSDKPath)
		}
	}

	// Match setter methods to JSON field names
	for tfName, setterMethod := range fieldMappings {
		field := TerraformField{
			TerraformName: tfName,
			SetterMethod:  setterMethod,
		}

		// Check if this is a nested field (e.g., provider_organization)
		var targetJsonFields map[string]jsonFieldInfo
		var parentFieldName string

		for nestedName := range nestedObjectTypes {
			if strings.HasPrefix(tfName, nestedName+"_") {
				// This field belongs to a nested object
				targetJsonFields = nestedJsonFields[nestedName]
				parentFieldName = nestedName
				break
			}
		}

		// If not nested, use main model's fields
		if targetJsonFields == nil {
			targetJsonFields = jsonFieldNames
		}

		// Extract field name from setter (e.g., SetDisplayName -> DisplayName)
		if strings.HasPrefix(setterMethod, "Set") {
			fieldName := strings.TrimPrefix(setterMethod, "Set")

			// Get JSON field name from SDK model
			if jsonInfo, ok := targetJsonFields[fieldName]; ok {
				field.GraphFieldName = jsonInfo.jsonName
				field.GraphFieldType = jsonInfo.fieldType

				// For nested fields, prefix the JSON field name with parent
				if parentFieldName != "" {
					// Convert parent field name from snake_case to camelCase for JSON path
					parentCamelCase := snakeToCamelCase(parentFieldName)
					field.GraphFieldName = parentCamelCase + "." + jsonInfo.jsonName
				}
			} else {
				// Fallback: convert to camelCase
				fallbackName := strings.ToLower(fieldName[0:1]) + fieldName[1:]
				if parentFieldName != "" {
					parentCamelCase := snakeToCamelCase(parentFieldName)
					field.GraphFieldName = parentCamelCase + "." + fallbackName
				} else {
					field.GraphFieldName = fallbackName
				}
			}
		} else {
			// For non-Set methods, try to match directly
			if jsonInfo, ok := targetJsonFields[setterMethod]; ok {
				field.GraphFieldName = jsonInfo.jsonName
				field.GraphFieldType = jsonInfo.fieldType
				if parentFieldName != "" {
					field.GraphFieldName = parentFieldName + "." + jsonInfo.jsonName
				}
			}
		}

		fields = append(fields, field)
	}

	return fields
}

type jsonFieldInfo struct {
	jsonName  string
	fieldType string
}

// extractJSONFieldNamesFromSDK extracts JSON field names from SDK model file
func (idx *Indexer) extractJSONFieldNamesFromSDK(sdkModelPath string) map[string]jsonFieldInfo {
	jsonFields := make(map[string]jsonFieldInfo)

	data, err := os.ReadFile(sdkModelPath)
	if err != nil {
		return jsonFields
	}

	content := string(data)

	// Look for Serialize method which writes JSON field names
	// Pattern: writer.WriteStringValue("displayName", m.GetDisplayName())
	pattern := regexp.MustCompile(`writer\.Write\w+Value\("([^"]+)",\s*m\.Get(\w+)\(\)`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			jsonName := match[1]
			fieldName := match[2]
			jsonFields[fieldName] = jsonFieldInfo{
				jsonName: jsonName,
			}
		}
	}

	// Also extract field types from getter methods
	// Pattern: func (m *Model) GetFieldName()(*string) or (*time.Time) etc
	// Note: There may be optional whitespace between )(
	typePattern := regexp.MustCompile(`func \(m \*\w+\) Get(\w+)\(\)\s*\(([^)]+)\)`)
	typeMatches := typePattern.FindAllStringSubmatch(content, -1)

	for _, match := range typeMatches {
		if len(match) >= 3 {
			fieldName := match[1]
			fieldType := match[2]

			// Clean up type (remove pointer and package prefixes)
			fieldType = strings.TrimPrefix(fieldType, "*")
			if idx := strings.LastIndex(fieldType, "."); idx != -1 {
				fieldType = fieldType[idx+1:]
			}

			if info, ok := jsonFields[fieldName]; ok {
				info.fieldType = fieldType
				jsonFields[fieldName] = info
			} else {
				jsonFields[fieldName] = jsonFieldInfo{
					fieldType: fieldType,
				}
			}
		}
	}

	return jsonFields
}

// findSDKModelFile finds the SDK model file for a given model type
func (idx *Indexer) findSDKModelFile(modelType string) string {
	// Get SDK paths from environment
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = filepath.Join(os.Getenv("HOME"), "go")
	}

	// Try both beta and v1.0 SDKs
	sdkPackages := []string{
		"github.com/microsoftgraph/msgraph-beta-sdk-go",
		"github.com/microsoftgraph/msgraph-sdk-go",
	}

	for _, sdkPackage := range sdkPackages {
		// Get SDK version from go.mod
		goModPath := filepath.Join(idx.providerPath, "go.mod")
		goModData, err := os.ReadFile(goModPath)
		if err != nil {
			continue
		}

		var sdkVersion string
		pattern := regexp.MustCompile(regexp.QuoteMeta(sdkPackage) + `\s+(v[\d.]+)`)
		if match := pattern.FindStringSubmatch(string(goModData)); len(match) > 1 {
			sdkVersion = match[1]
		}

		if sdkVersion == "" {
			continue
		}

		// Build path to SDK models
		sdkPath := filepath.Join(goPath, "pkg", "mod",
			strings.Replace(sdkPackage, "/", string(filepath.Separator), -1)+"@"+sdkVersion,
			"models")

		// Convert model type to snake_case for file name
		modelFileName := toSnakeCase(modelType) + ".go"
		modelFilePath := filepath.Join(sdkPath, modelFileName)

		if _, err := os.Stat(modelFilePath); err == nil {
			return modelFilePath
		}
	}

	return ""
}

// apiCallVisitor is an AST visitor that extracts API calls
type apiCallVisitor struct {
	fset       *token.FileSet
	filePath   string
	apiVersion string
	calls      map[string][]APICall
	currentOp  string
	indexer    *Indexer
}

func (v *apiCallVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// Look for CRUD method declarations
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		methodName := funcDecl.Name.Name
		if methodName == "Create" || methodName == "Read" || methodName == "Update" || methodName == "Delete" || methodName == "Invoke" {
			v.currentOp = methodName
		}
	}

	// Look for selector expressions (method calls)
	if selExpr, ok := node.(*ast.SelectorExpr); ok {
		if v.currentOp != "" {
			// Check if this is an HTTP method call
			httpMethod := selExpr.Sel.Name
			if httpMethod == "Get" || httpMethod == "Post" || httpMethod == "Patch" || httpMethod == "Put" || httpMethod == "Delete" {
				// Try to extract the full method chain
				chain := v.extractMethodChain(selExpr)
				if strings.Contains(chain, "client") {
					apiCall := v.buildAPICall(chain, httpMethod, selExpr)
					v.calls[v.currentOp] = append(v.calls[v.currentOp], apiCall)
				}
			}
		}
	}

	return v
}

// extractMethodChain extracts the full method chain from a selector expression
func (v *apiCallVisitor) extractMethodChain(expr ast.Expr) string {
	var parts []string

	current := expr
	for {
		switch e := current.(type) {
		case *ast.SelectorExpr:
			parts = append([]string{e.Sel.Name + "()"}, parts...)
			current = e.X
		case *ast.CallExpr:
			// Handle function calls in the chain
			if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
				parts = append([]string{sel.Sel.Name + "()"}, parts...)
				current = sel.X
			} else {
				return strings.Join(parts, ".")
			}
		case *ast.Ident:
			parts = append([]string{e.Name}, parts...)
			return strings.Join(parts, ".")
		default:
			return strings.Join(parts, ".")
		}
	}
}

// buildAPICall builds an APICall from the extracted information
func (v *apiCallVisitor) buildAPICall(chain, httpMethod string, expr ast.Expr) APICall {
	pos := v.fset.Position(expr.Pos())

	// Build SDK method chain
	sdkChain := chain

	// Build API URL pattern
	apiURL := v.buildAPIURL(chain, httpMethod)

	// Make file path relative
	relPath := strings.Replace(v.filePath, os.Getenv("PWD")+"/", "", 1)
	relPath = strings.Replace(relPath, "/Users/dafyddwatkins/GitHub/deploymenttheory/", "", 1)
	if !strings.HasPrefix(relPath, "terraform-provider-microsoft365/") {
		relPath = "terraform-provider-microsoft365/" + strings.TrimPrefix(relPath, "../")
	}

	return APICall{
		SDKMethodChain: sdkChain,
		HTTPMethod:     strings.ToUpper(httpMethod),
		APIURLPattern:  apiURL,
		LineNumber:     pos.Line,
		FilePath:       relPath,
	}
}

// buildAPIURL builds the API URL from the SDK chain by inspecting SDK source
func (v *apiCallVisitor) buildAPIURL(chain, httpMethod string) string {
	version := "beta"
	if strings.Contains(v.apiVersion, "v1.0") {
		version = "v1.0"
	}

	baseURL := fmt.Sprintf("https://graph.microsoft.com/%s", version)

	// Try to extract URL from SDK source first
	if url, ok := v.indexer.getSDKURLForChain(chain, version); ok {
		return url
	}

	// Fallback to parsing the chain
	parts := strings.Split(chain, ".")
	var urlParts []string

	for _, part := range parts {
		part = strings.TrimSuffix(part, "()")

		// Skip client references
		if part == "client" || part == "r" || part == "d" || part == "a" {
			continue
		}

		// Skip HTTP methods
		if part == httpMethod {
			continue
		}

		// Convert SDK method names to URL paths
		urlPart := v.indexer.convertSDKMethodToURL(part, version)
		if urlPart != "" {
			urlParts = append(urlParts, urlPart)
		}
	}

	if len(urlParts) == 0 {
		return baseURL
	}

	return baseURL + "/" + strings.Join(urlParts, "/")
}

// getSDKURLForChain extracts URL from SDK source code
func (idx *Indexer) getSDKURLForChain(chain, version string) (string, bool) {
	// Parse chain to identify SDK package and method
	parts := strings.Split(chain, ".")
	if len(parts) < 2 {
		return "", false
	}

	// Try to find URL in cache
	cacheKey := version + ":" + chain
	if url, ok := idx.sdkURLCache[cacheKey]; ok {
		return url, true
	}

	// Extract the URL by parsing SDK files
	url := idx.extractURLFromSDK(parts, version)
	if url != "" {
		url = fmt.Sprintf("https://graph.microsoft.com/%s/%s", version, url)
		idx.sdkURLCache[cacheKey] = url
		return url, true
	}

	return "", false
}

// extractURLFromSDK parses SDK source files to extract URL templates
func (idx *Indexer) extractURLFromSDK(methodChain []string, version string) string {
	// Determine SDK module path based on version
	goModPath := filepath.Join(idx.providerPath, "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	// Find SDK package in GOPATH or module cache
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = filepath.Join(os.Getenv("HOME"), "go")
	}

	sdkPackage := "github.com/microsoftgraph/msgraph-sdk-go"
	if version == "beta" {
		sdkPackage = "github.com/microsoftgraph/msgraph-beta-sdk-go"
	}

	// Get SDK version from go.mod
	var sdkVersion string
	pattern := regexp.MustCompile(regexp.QuoteMeta(sdkPackage) + `\s+(v[\d.]+)`)
	if match := pattern.FindStringSubmatch(string(goModData)); len(match) > 1 {
		sdkVersion = match[1]
	}

	if sdkVersion == "" {
		return ""
	}

	// Try module cache first
	sdkPath := filepath.Join(goPath, "pkg", "mod", strings.Replace(sdkPackage, "/", string(filepath.Separator), -1)+"@"+sdkVersion)

	// Check if SDK path exists
	if _, err := os.Stat(sdkPath); os.IsNotExist(err) {
		return ""
	}

	// Parse method chain to find the right SDK file
	return idx.findURLInSDKPath(sdkPath, methodChain)
}

// findURLInSDKPath searches SDK files for URL templates
func (idx *Indexer) findURLInSDKPath(sdkPath string, methodChain []string) string {
	// Find the last non-HTTP-method in the chain
	// The SDK structure: client().Resource().ByResourceId().SubResource().HttpMethod()
	// We want to extract the URL from the last resource's request builder before the HTTP method

	var currentPath string

	// Filter out client references and HTTP methods
	var relevantMethods []string
	for _, method := range methodChain {
		method = strings.TrimSuffix(method, "()")
		if method == "client" || method == "r" || method == "d" || method == "a" {
			continue
		}
		// Skip HTTP methods
		if method == "Get" || method == "Post" || method == "Patch" || method == "Put" || method == "Delete" {
			continue
		}
		relevantMethods = append(relevantMethods, method)
	}

	if len(relevantMethods) == 0 {
		return ""
	}

	// Build the path to the last method's request builder file
	// Navigate through the SDK structure
	var parentMethod string
	for i, method := range relevantMethods {
		snakeCase := toSnakeCase(method)

		if i == len(relevantMethods)-1 {
			// This is the last method, this is where we extract the URL

			// Check if this is a ByXxxId item request builder
			if strings.HasPrefix(method, "By") && strings.HasSuffix(method, "Id") {
				// Derive the item name from the parent method
				// e.g., for Users().ByUserId(), parent is "Users", item file is "user_item_request_builder.go"
				if parentMethod != "" {
					// Singularize parent method name (Users -> User)
					itemName := strings.TrimSuffix(parentMethod, "s") // Simple singularization
					itemFileName := toSnakeCase(itemName) + "_item_request_builder.go"
					builderFile := filepath.Join(sdkPath, currentPath, itemFileName)

					if url := extractURLTemplateFromFile(builderFile); url != "" {
						return url
					}
				}

				// Fallback: try generic pattern
				pattern := filepath.Join(sdkPath, currentPath, "*item_request_builder.go")
				matches, _ := filepath.Glob(pattern)
				if len(matches) > 0 {
					return extractURLTemplateFromFile(matches[0])
				}
			} else {
				// Look for regular request builder
				builderFile := filepath.Join(sdkPath, currentPath, snakeCase+"_request_builder.go")
				if url := extractURLTemplateFromFile(builderFile); url != "" {
					return url
				}
			}
		} else {
			// Navigate to the next level
			parentMethod = method
			currentPath = filepath.Join(currentPath, snakeCase)
		}
	}

	return ""
}

// extractURLTemplateFromFile extracts URL template from SDK request builder file
func extractURLTemplateFromFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	content := string(data)

	// Look for NewBaseRequestBuilder call with URL template
	// Pattern: NewBaseRequestBuilder(requestAdapter, "{+baseurl}/path/to/resource{?params}", pathParameters)
	pattern := regexp.MustCompile(`NewBaseRequestBuilder\([^,]+,\s*"([^"]+)"`)
	if match := pattern.FindStringSubmatch(content); len(match) > 1 {
		template := match[1]
		// Clean up template - remove {+baseurl} and query params
		template = strings.Replace(template, "{+baseurl}/", "", 1)
		template = strings.Replace(template, "{+baseurl}", "", 1)
		// Remove query parameters
		if idx := strings.Index(template, "{?"); idx != -1 {
			template = template[:idx]
		}
		return template
	}

	return ""
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// snakeToCamelCase converts snake_case to camelCase
func snakeToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 && len(parts[i]) > 0 {
			// Capitalize first letter of each word after the first
			parts[i] = strings.ToUpper(parts[i][0:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// convertSDKMethodToURL converts SDK method names to URL paths (fallback)
func (idx *Indexer) convertSDKMethodToURL(method, version string) string {
	// Try cache first
	cacheKey := version + ":" + method
	if url, ok := idx.sdkURLCache[cacheKey]; ok {
		return url
	}

	// For ByXxxId methods, return parameter placeholder
	if strings.HasPrefix(method, "By") && strings.HasSuffix(method, "Id") {
		paramName := strings.TrimSuffix(strings.TrimPrefix(method, "By"), "Id")
		return "{" + toSnakeCase(paramName) + "Id}"
	}

	// Convert PascalCase to camelCase
	if len(method) > 0 && method[0] >= 'A' && method[0] <= 'Z' {
		return strings.ToLower(method[0:1]) + method[1:]
	}

	return method
}

// buildMetadata builds the metadata summary
func (idx *Indexer) buildMetadata() Metadata {
	counts := make(map[string]int)

	for _, ep := range idx.endpoints {
		counts[ep.ResourceType]++
	}

	return Metadata{
		TotalEndpoints: len(idx.endpoints),
		ResourceTypes:  counts,
		SDKVersions:    idx.sdkVersions,
	}
}

func main() {
	// Get provider path
	providerPath := "/Users/dafyddwatkins/GitHub/deploymenttheory/terraform-provider-microsoft365"
	if len(os.Args) > 1 {
		providerPath = os.Args[1]
	}

	// Create indexer
	indexer := NewIndexer(providerPath)

	// Run indexing
	index, err := indexer.Index()
	if err != nil {
		log.Fatalf("Indexing failed: %v", err)
	}

	// Sort endpoints for consistent output
	sort.Slice(index.Endpoints, func(i, j int) bool {
		if index.Endpoints[i].ResourceType != index.Endpoints[j].ResourceType {
			return index.Endpoints[i].ResourceType < index.Endpoints[j].ResourceType
		}
		if index.Endpoints[i].ServiceDomain != index.Endpoints[j].ServiceDomain {
			return index.Endpoints[i].ServiceDomain < index.Endpoints[j].ServiceDomain
		}
		return index.Endpoints[i].ResourceName < index.Endpoints[j].ResourceName
	})

	// Write output
	outputPath := filepath.Join(providerPath, "api_endpoint_index_go.json")
	output, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("INDEXING COMPLETE")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total endpoints indexed: %d\n", index.Metadata.TotalEndpoints)
	for rt, count := range index.Metadata.ResourceTypes {
		fmt.Printf("  %s: %d\n", rt, count)
	}
	fmt.Printf("\nSDK Versions:\n")
	for sdk, version := range index.Metadata.SDKVersions {
		fmt.Printf("  %s: %s\n", sdk, version)
	}
	fmt.Printf("\nOutput saved to: %s\n", outputPath)

	// Print sample
	if len(index.Endpoints) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("SAMPLE OUTPUT (first endpoint):")
		fmt.Println(strings.Repeat("=", 80))
		sample, _ := json.MarshalIndent(index.Endpoints[0], "", "  ")
		fmt.Println(string(sample))
	}
}
