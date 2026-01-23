package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeUsageMapV2(t *testing.T) {
	usage := initializeUsageMapV2()

	assert.NotNil(t, usage)
	assert.NotNil(t, usage.TerraformResources)
	assert.NotNil(t, usage.TerraformActions)
	assert.NotNil(t, usage.TerraformListActions)
	assert.NotNil(t, usage.TerraformEphemerals)
	assert.NotNil(t, usage.TerraformDataSources)
	assert.NotNil(t, usage.SDKToResourceIndex)
	assert.Empty(t, usage.TerraformResources)
}

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Regular Go file",
			path:     "/path/to/resource.go",
			expected: false,
		},
		{
			name:     "Test file",
			path:     "/path/to/resource_test.go",
			expected: true,
		},
		{
			name:     "Non-Go file",
			path:     "/path/to/config.yaml",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipFile(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseEntityFromPath(t *testing.T) {
	repoPath := "/repo"

	tests := []struct {
		name         string
		path         string
		expectedType string
		expectedName string
		shouldBeNil  bool
	}{
		{
			name:         "Resource path",
			path:         "/repo/internal/services/resources/users/graph_beta/user/resource.go",
			expectedType: "resource",
			expectedName: "microsoft365_user",
			shouldBeNil:  false,
		},
		{
			name:         "Action path",
			path:         "/repo/internal/services/actions/device_management/graph_beta/managed_device/invoke.go",
			expectedType: "action",
			expectedName: "microsoft365_managed_device",
			shouldBeNil:  false,
		},
		{
			name:         "List action path",
			path:         "/repo/internal/services/list-resources/device_management/graph_beta/devices/list.go",
			expectedType: "list-action",
			expectedName: "microsoft365_devices",
			shouldBeNil:  false,
		},
		{
			name:         "Ephemeral path",
			path:         "/repo/internal/services/ephemerals/identity/graph_beta/token/ephemeral.go",
			expectedType: "ephemeral",
			expectedName: "microsoft365_token",
			shouldBeNil:  false,
		},
		{
			name:         "Data source path",
			path:         "/repo/internal/services/data-sources/users/graph_beta/user/data.go",
			expectedType: "data-source",
			expectedName: "microsoft365_user",
			shouldBeNil:  false,
		},
		{
			name:        "Common path (not an entity)",
			path:        "/repo/internal/services/common/errors/error.go",
			shouldBeNil: true,
		},
		{
			name:        "Too short path",
			path:        "/repo/internal/services/resources/users/resource.go",
			shouldBeNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := parseEntityFromPath(tt.path, repoPath)

			if tt.shouldBeNil {
				assert.Nil(t, entity)
			} else {
				require.NotNil(t, entity)
				assert.Equal(t, tt.expectedType, entity.Type)
				assert.Equal(t, tt.expectedName, entity.Name)
			}
		})
	}
}

func TestExtractEntityInfo(t *testing.T) {
	tests := []struct {
		name         string
		relPath      string
		prefix       string
		entityType   string
		expectedName string
		shouldBeNil  bool
	}{
		{
			name:         "Valid resource path",
			relPath:      "internal/services/resources/users/graph_beta/user/resource.go",
			prefix:       "internal/services/resources/",
			entityType:   "resource",
			expectedName: "microsoft365_user",
			shouldBeNil:  false,
		},
		{
			name:         "Device management resource",
			relPath:      "internal/services/resources/device_management/graph_beta/settings_catalog/construct.go",
			prefix:       "internal/services/resources/",
			entityType:   "resource",
			expectedName: "microsoft365_settings_catalog",
			shouldBeNil:  false,
		},
		{
			name:        "Path too short",
			relPath:     "internal/services/resources/users/resource.go",
			prefix:      "internal/services/resources/",
			entityType:  "resource",
			shouldBeNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := extractEntityInfo(tt.relPath, tt.prefix, tt.entityType)

			if tt.shouldBeNil {
				assert.Nil(t, entity)
			} else {
				require.NotNil(t, entity)
				assert.Equal(t, tt.entityType, entity.Type)
				assert.Equal(t, tt.expectedName, entity.Name)
			}
		})
	}
}

func TestGetOrCreateResourceInfo(t *testing.T) {
	usage := initializeUsageMapV2()

	entity := &Entity{
		Type: "resource",
		Name: "microsoft365_user",
		Path: "internal/services/resources/users/graph_beta/user",
	}

	// First call should create
	info1 := getOrCreateResourceInfo(usage, entity)
	require.NotNil(t, info1)
	assert.Equal(t, entity.Path, info1.ResourcePath)
	assert.NotNil(t, info1.SDKDependencies.FieldsUsed)

	// Second call should return same instance
	info2 := getOrCreateResourceInfo(usage, entity)
	assert.Equal(t, info1, info2)

	// Verify it's in the correct map
	assert.Len(t, usage.TerraformResources, 1)
	assert.Contains(t, usage.TerraformResources, "microsoft365_user")
}

func TestGetOrCreateResourceInfoDifferentTypes(t *testing.T) {
	usage := initializeUsageMapV2()

	tests := []struct {
		entityType string
		targetMap  *map[string]*ResourceInfo
	}{
		{"resource", &usage.TerraformResources},
		{"action", &usage.TerraformActions},
		{"list-action", &usage.TerraformListActions},
		{"ephemeral", &usage.TerraformEphemerals},
		{"data-source", &usage.TerraformDataSources},
	}

	for _, tt := range tests {
		t.Run(tt.entityType, func(t *testing.T) {
			entity := &Entity{
				Type: tt.entityType,
				Name: "microsoft365_test",
				Path: "internal/services/" + tt.entityType + "/test",
			}

			info := getOrCreateResourceInfo(usage, entity)
			require.NotNil(t, info)
			assert.Len(t, *tt.targetMap, 1)
		})
	}
}

func TestIsSDKImport(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		expected   bool
	}{
		{
			name:       "Microsoft Graph Beta SDK",
			importPath: "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
			expected:   true,
		},
		{
			name:       "Kiota abstractions",
			importPath: "github.com/microsoft/kiota-abstractions-go",
			expected:   true,
		},
		{
			name:       "Standard library",
			importPath: "fmt",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSDKImport(tt.importPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetImportAlias(t *testing.T) {
	t.Run("Explicit alias", func(t *testing.T) {
		imp := &ast.ImportSpec{
			Name: &ast.Ident{Name: "customalias"},
			Path: &ast.BasicLit{Value: `"github.com/microsoftgraph/msgraph-beta-sdk-go/models"`},
		}
		alias := getImportAlias(imp, "github.com/microsoftgraph/msgraph-beta-sdk-go/models")
		assert.Equal(t, "customalias", alias)
	})

	t.Run("Inferred alias from path", func(t *testing.T) {
		imp := &ast.ImportSpec{
			Name: nil,
			Path: &ast.BasicLit{Value: `"github.com/microsoftgraph/msgraph-beta-sdk-go/models"`},
		}
		alias := getImportAlias(imp, "github.com/microsoftgraph/msgraph-beta-sdk-go/models")
		assert.Equal(t, "models", alias)
	})
}

func TestCollectImports(t *testing.T) {
	source := `package test

import (
	"fmt"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/hashicorp/terraform-plugin-framework/types"
)
`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", source, parser.ParseComments)
	require.NoError(t, err)

	fileImports := collectImports(node)

	assert.Len(t, fileImports, 2) // Only SDK imports
	assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models", fileImports["models"])
	assert.Equal(t, "github.com/microsoftgraph/msgraph-sdk-go/users", fileImports["users"])
	assert.NotContains(t, fileImports, "fmt")
	assert.NotContains(t, fileImports, "types")
}

func TestSimplifyTypeName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Models type",
			fullName: "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User",
			expected: "models.User",
		},
		{
			name:     "Package type",
			fullName: "github.com/microsoftgraph/msgraph-beta-sdk-go/users.UserItemRequestBuilder",
			expected: "users.UserItemRequestBuilder",
		},
		{
			name:     "Already simple",
			fullName: "models.User",
			expected: "models.User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := simplifyTypeName(tt.fullName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrackFieldAccess(t *testing.T) {
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			FieldsUsed: make(map[string][]string),
		},
	}

	trackFieldAccess("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", "DisplayName", resourceInfo)
	trackFieldAccess("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", "DisplayName", resourceInfo)
	trackFieldAccess("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", "Email", resourceInfo)

	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed, "models.User")
	assert.Len(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], 2) // No duplicates
	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "DisplayName")
	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "Email")
}

func TestTrackTypeInstantiation(t *testing.T) {
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types: []string{},
		},
	}

	trackTypeInstantiation("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", resourceInfo)
	trackTypeInstantiation("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", resourceInfo) // Duplicate

	assert.Len(t, resourceInfo.SDKDependencies.Types, 1) // No duplicates
	assert.Contains(t, resourceInfo.SDKDependencies.Types, "models.User")
}

func TestTrackEnumUsage(t *testing.T) {
	usage := initializeUsageMapV2()
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			EnumsUsed: []EnumUsage{},
		},
	}

	trackEnumUsage("github.com/microsoftgraph/msgraph-beta-sdk-go/models", "ParseRunAsAccountType", resourceInfo, usage)
	trackEnumUsage("github.com/microsoftgraph/msgraph-beta-sdk-go/models", "ParseRunAsAccountType", resourceInfo, usage) // Duplicate

	assert.Len(t, resourceInfo.SDKDependencies.EnumsUsed, 1) // No duplicates
	assert.Equal(t, "models.RunAsAccountType", resourceInfo.SDKDependencies.EnumsUsed[0].Enum)
}

func TestExtractVarName(t *testing.T) {
	t.Run("Valid identifier", func(t *testing.T) {
		expr := &ast.Ident{Name: "user"}
		name := extractVarName(expr)
		assert.Equal(t, "user", name)
	})

	t.Run("Non-identifier expression", func(t *testing.T) {
		expr := &ast.BasicLit{Value: "123"}
		name := extractVarName(expr)
		assert.Equal(t, "", name)
	})
}

func TestExtractTypeFromCallExpr(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}

	t.Run("Valid SDK call", func(t *testing.T) {
		// models.NewUser()
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "NewUser"},
			},
		}

		typeName := extractTypeFromCallExpr(call, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.NewUser", typeName)
	})

	t.Run("Unknown package", func(t *testing.T) {
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "unknown"},
				Sel: &ast.Ident{Name: "NewSomething"},
			},
		}

		typeName := extractTypeFromCallExpr(call, fileImports)
		assert.Equal(t, "", typeName)
	})
}

func TestExtractTypeFromCompositeLit(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}

	t.Run("Valid SDK struct literal", func(t *testing.T) {
		// models.User{}
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "User"},
			},
		}

		typeName := extractTypeFromCompositeLit(comp, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})
}

func TestExtractTypeFromUnaryExpr(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}

	t.Run("Pointer to struct literal", func(t *testing.T) {
		// &models.User{}
		unary := &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "models"},
					Sel: &ast.Ident{Name: "User"},
				},
			},
		}

		typeName := extractTypeFromUnaryExpr(unary, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})
}

func TestTrackStructFields(t *testing.T) {
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			FieldsUsed: make(map[string][]string),
		},
	}

	elts := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "DisplayName"},
			Value: &ast.BasicLit{Value: `"John"`},
		},
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "Email"},
			Value: &ast.BasicLit{Value: `"john@example.com"`},
		},
	}

	trackStructFields("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", elts, resourceInfo)

	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed, "models.User")
	assert.Len(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], 2)
	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "DisplayName")
	assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "Email")
}

func TestProcessAssignments(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}
	varTypes := make(map[string]string)

	// user := models.NewUser()
	stmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "user"},
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "models"},
					Sel: &ast.Ident{Name: "NewUser"},
				},
			},
		},
	}

	processAssignments(stmt, fileImports, varTypes)

	assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.NewUser", varTypes["user"])
}

func TestCalculateStatistics(t *testing.T) {
	usage := initializeUsageMapV2()

	// Add some resources
	usage.TerraformResources["microsoft365_user"] = &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types:         []string{"models.User", "models.AssignedLicense"},
			MethodsCalled: []string{"models.NewUser"},
			EnumsUsed: []EnumUsage{
				{Enum: "models.DayOfWeek"},
			},
		},
	}

	usage.TerraformActions["microsoft365_lock_device"] = &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types:         []string{"models.Device"},
			MethodsCalled: []string{"devicemanagement.Post"},
		},
	}

	stats := calculateStatistics(usage)

	assert.Equal(t, 1, stats.TotalResources)
	assert.Equal(t, 1, stats.TotalActions)
	assert.Equal(t, 0, stats.TotalListActions)
	assert.Equal(t, 0, stats.TotalEphemerals)
	assert.Equal(t, 0, stats.TotalDataSources)
	assert.Equal(t, 3, stats.TotalSDKTypesUsed) // User, AssignedLicense, Device
	assert.Equal(t, 2, stats.TotalSDKMethodsUsed)
	assert.Equal(t, 1, stats.TotalEnumsTracked)
}

func TestIndexSDKUsage(t *testing.T) {
	index := make(map[string][]string)
	typesSet := make(map[string]bool)
	methodsSet := make(map[string]bool)
	enumsSet := make(map[string]bool)

	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types:         []string{"models.User"},
			MethodsCalled: []string{"models.NewUser"},
			EnumsUsed: []EnumUsage{
				{Enum: "models.DayOfWeek"},
			},
		},
	}

	indexSDKUsage("microsoft365_user", resourceInfo, index, typesSet, methodsSet, enumsSet)

	// Check reverse index
	assert.Contains(t, index, "models.User")
	assert.Contains(t, index["models.User"], "microsoft365_user")
	assert.Contains(t, index, "models.DayOfWeek")
	assert.Contains(t, index["models.DayOfWeek"], "microsoft365_user")

	// Check sets
	assert.True(t, typesSet["models.User"])
	assert.True(t, methodsSet["models.NewUser"])
	assert.True(t, enumsSet["models.DayOfWeek"])
}

func TestIntegrationFullAnalysis(t *testing.T) {
	source := `package test

import (
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func CreateUser() {
	user := models.NewUser()
	user.SetDisplayName("John")
	
	config := models.DeviceConfiguration{
		DisplayName: "Test Config",
	}
	
	accountType := models.ParseRunAsAccountType("system")
}
`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", source, parser.ParseComments)
	require.NoError(t, err)

	usage := initializeUsageMapV2()
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types:         []string{},
			FieldsUsed:    make(map[string][]string),
			MethodsCalled: []string{},
			EnumsUsed:     []EnumUsage{},
		},
		Files: []string{},
	}

	fileImports := collectImports(node)
	analyzeASTForEntity(node, "/test.go", fileImports, resourceInfo, usage)

	t.Run("Tracks types", func(t *testing.T) {
		// Only struct literals (DeviceConfiguration{}) get tracked as types
		// Constructor calls (NewUser()) are tracked as methods
		assert.Contains(t, resourceInfo.SDKDependencies.Types, "models.DeviceConfiguration")
		assert.Len(t, resourceInfo.SDKDependencies.Types, 1)
	})

	t.Run("Tracks fields", func(t *testing.T) {
		assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed, "models.DeviceConfiguration")
		assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.DeviceConfiguration"], "DisplayName")
	})

	t.Run("Tracks methods", func(t *testing.T) {
		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.NewUser")
		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.NewUser.SetDisplayName")
	})

	t.Run("Tracks enums", func(t *testing.T) {
		found := false
		for _, enum := range resourceInfo.SDKDependencies.EnumsUsed {
			if enum.Enum == "models.RunAsAccountType" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should track RunAsAccountType enum")
	})
}

func TestProcessSelectorExpr(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}
	varTypes := map[string]string{
		"user": "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User",
	}
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			FieldsUsed: make(map[string][]string),
		},
	}

	t.Run("Field access on typed variable", func(t *testing.T) {
		// user.DisplayName
		sel := &ast.SelectorExpr{
			X:   &ast.Ident{Name: "user"},
			Sel: &ast.Ident{Name: "DisplayName"},
		}

		processSelectorExpr(sel, fileImports, varTypes, resourceInfo)

		assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed, "models.User")
		assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "DisplayName")
	})

	t.Run("Package-level reference ignored", func(t *testing.T) {
		initialLen := len(resourceInfo.SDKDependencies.FieldsUsed)

		// models.User (package reference, not field access)
		sel := &ast.SelectorExpr{
			X:   &ast.Ident{Name: "models"},
			Sel: &ast.Ident{Name: "User"},
		}

		processSelectorExpr(sel, fileImports, varTypes, resourceInfo)

		// Should not add new fields
		assert.Equal(t, initialLen, len(resourceInfo.SDKDependencies.FieldsUsed))
	})

	t.Run("Unknown variable ignored", func(t *testing.T) {
		sel := &ast.SelectorExpr{
			X:   &ast.Ident{Name: "unknown"},
			Sel: &ast.Ident{Name: "Field"},
		}

		processSelectorExpr(sel, fileImports, varTypes, resourceInfo)
		// Should not panic
	})

	t.Run("Non-identifier selector", func(t *testing.T) {
		sel := &ast.SelectorExpr{
			X:   &ast.BasicLit{Value: "123"},
			Sel: &ast.Ident{Name: "Field"},
		}

		processSelectorExpr(sel, fileImports, varTypes, resourceInfo)
		// Should not panic
	})
}

func TestProcessCallExpr(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}
	varTypes := map[string]string{
		"user": "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User",
	}
	usage := initializeUsageMapV2()
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			MethodsCalled: []string{},
			EnumsUsed:     []EnumUsage{},
		},
	}

	t.Run("Package method call", func(t *testing.T) {
		// models.NewUser()
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "NewUser"},
			},
		}

		processCallExpr(call, "/test.go", fileImports, varTypes, resourceInfo, usage)

		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.NewUser")
	})

	t.Run("Object method call", func(t *testing.T) {
		// user.SetDisplayName()
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "user"},
				Sel: &ast.Ident{Name: "SetDisplayName"},
			},
		}

		processCallExpr(call, "/test.go", fileImports, varTypes, resourceInfo, usage)

		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.User.SetDisplayName")
	})

	t.Run("Enum parser call", func(t *testing.T) {
		// models.ParseInstallIntent()
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "ParseInstallIntent"},
			},
		}

		processCallExpr(call, "/test.go", fileImports, varTypes, resourceInfo, usage)

		found := false
		for _, enum := range resourceInfo.SDKDependencies.EnumsUsed {
			if enum.Enum == "models.InstallIntent" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Non-selector call ignored", func(t *testing.T) {
		call := &ast.CallExpr{
			Fun: &ast.Ident{Name: "someFunc"},
		}

		processCallExpr(call, "/test.go", fileImports, varTypes, resourceInfo, usage)
		// Should not panic
	})

	t.Run("Non-identifier in selector", func(t *testing.T) {
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.BasicLit{Value: "123"},
				Sel: &ast.Ident{Name: "Method"},
			},
		}

		processCallExpr(call, "/test.go", fileImports, varTypes, resourceInfo, usage)
		// Should not panic
	})
}

func TestProcessCompositeLit(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			Types:      []string{},
			FieldsUsed: make(map[string][]string),
		},
	}

	t.Run("Valid struct literal with fields", func(t *testing.T) {
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "User"},
			},
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key:   &ast.Ident{Name: "DisplayName"},
					Value: &ast.BasicLit{Value: `"John"`},
				},
			},
		}

		processCompositeLit(comp, fileImports, resourceInfo)

		assert.Contains(t, resourceInfo.SDKDependencies.Types, "models.User")
		assert.Contains(t, resourceInfo.SDKDependencies.FieldsUsed["models.User"], "DisplayName")
	})

	t.Run("Non-selector type ignored", func(t *testing.T) {
		comp := &ast.CompositeLit{
			Type: &ast.Ident{Name: "SomeType"},
		}

		processCompositeLit(comp, fileImports, resourceInfo)
		// Should not panic
	})

	t.Run("Non-identifier package", func(t *testing.T) {
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.BasicLit{Value: "123"},
				Sel: &ast.Ident{Name: "Type"},
			},
		}

		processCompositeLit(comp, fileImports, resourceInfo)
		// Should not panic
	})

	t.Run("Unknown package", func(t *testing.T) {
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "unknown"},
				Sel: &ast.Ident{Name: "Type"},
			},
		}

		processCompositeLit(comp, fileImports, resourceInfo)
		// Should not panic
	})
}

func TestTrackStructFieldsEdgeCases(t *testing.T) {
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			FieldsUsed: make(map[string][]string),
		},
	}

	t.Run("Non-KeyValueExpr elements ignored", func(t *testing.T) {
		elts := []ast.Expr{
			&ast.BasicLit{Value: `"value"`},
			&ast.Ident{Name: "field"},
		}

		trackStructFields("models.Test", elts, resourceInfo)
		// Should not panic
	})

	t.Run("Non-identifier keys ignored", func(t *testing.T) {
		elts := []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.BasicLit{Value: `"key"`},
				Value: &ast.BasicLit{Value: `"value"`},
			},
		}

		trackStructFields("models.Test", elts, resourceInfo)
		// Should not panic
	})

	t.Run("Empty elements list", func(t *testing.T) {
		elts := []ast.Expr{}

		trackStructFields("models.Test", elts, resourceInfo)
		// Should not panic
	})

	t.Run("Duplicate fields not added", func(t *testing.T) {
		elts := []ast.Expr{
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Name"},
				Value: &ast.BasicLit{Value: `"test"`},
			},
			&ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "Name"}, // Duplicate
				Value: &ast.BasicLit{Value: `"test2"`},
			},
		}

		trackStructFields("models.Config", elts, resourceInfo)
		
		assert.Len(t, resourceInfo.SDKDependencies.FieldsUsed["models.Config"], 1)
	})
}

func TestExtractSDKType(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}

	t.Run("Call expression", func(t *testing.T) {
		// models.NewUser()
		expr := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "NewUser"},
			},
		}

		typeName := extractSDKType(expr, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.NewUser", typeName)
	})

	t.Run("Composite literal", func(t *testing.T) {
		// models.User{}
		expr := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "User"},
			},
		}

		typeName := extractSDKType(expr, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})

	t.Run("Unary expression", func(t *testing.T) {
		// &models.User{}
		expr := &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "models"},
					Sel: &ast.Ident{Name: "User"},
				},
			},
		}

		typeName := extractSDKType(expr, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})

	t.Run("Unknown expression type", func(t *testing.T) {
		expr := &ast.BasicLit{Value: "123"}
		
		typeName := extractSDKType(expr, fileImports)
		assert.Equal(t, "", typeName)
	})
}

func TestTrackPackageMethod(t *testing.T) {
	usage := initializeUsageMapV2()
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			MethodsCalled: []string{},
			EnumsUsed:     []EnumUsage{},
		},
	}

	t.Run("Regular method", func(t *testing.T) {
		trackPackageMethod("github.com/microsoftgraph/msgraph-beta-sdk-go/models", "NewUser", "/test.go", resourceInfo, usage)

		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.NewUser")
	})

	t.Run("Duplicate method not added twice", func(t *testing.T) {
		initialLen := len(resourceInfo.SDKDependencies.MethodsCalled)
		
		trackPackageMethod("github.com/microsoftgraph/msgraph-beta-sdk-go/models", "NewUser", "/test.go", resourceInfo, usage)

		assert.Len(t, resourceInfo.SDKDependencies.MethodsCalled, initialLen) // No duplicate
	})

	t.Run("Enum parser also tracks enum", func(t *testing.T) {
		trackPackageMethod("github.com/microsoftgraph/msgraph-beta-sdk-go/models", "ParseRunAsAccountType", "/test.go", resourceInfo, usage)

		// Method tracked
		assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.ParseRunAsAccountType")

		// Enum tracked
		found := false
		for _, enum := range resourceInfo.SDKDependencies.EnumsUsed {
			if enum.Enum == "models.RunAsAccountType" {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}

func TestTrackObjectMethod(t *testing.T) {
	resourceInfo := &ResourceInfo{
		SDKDependencies: SDKDependencies{
			MethodsCalled: []string{},
		},
	}

	trackObjectMethod("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", "SetDisplayName", resourceInfo)
	trackObjectMethod("github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", "SetDisplayName", resourceInfo) // Duplicate

	assert.Len(t, resourceInfo.SDKDependencies.MethodsCalled, 1) // No duplicates
	assert.Contains(t, resourceInfo.SDKDependencies.MethodsCalled, "models.User.SetDisplayName")
}

func TestOutputResults(t *testing.T) {
	t.Run("Valid usage map produces JSON", func(t *testing.T) {
		usage := initializeUsageMapV2()
		usage.TerraformResources["microsoft365_user"] = &ResourceInfo{
			ResourcePath: "internal/services/resources/users/graph_beta/user",
			SDKDependencies: SDKDependencies{
				Types: []string{"models.User"},
			},
		}

		err := outputResults(usage)
		assert.NoError(t, err)
	})

	t.Run("Empty usage map", func(t *testing.T) {
		usage := initializeUsageMapV2()

		err := outputResults(usage)
		assert.NoError(t, err)
	})
}

func TestComplexMultiResourceScenario(t *testing.T) {
	usage := initializeUsageMapV2()

	// Simulate multiple resources using the same SDK type
	userResource := &ResourceInfo{
		ResourcePath: "internal/services/resources/users/graph_beta/user",
		SDKDependencies: SDKDependencies{
			Types:         []string{"models.User"},
			MethodsCalled: []string{"models.NewUser", "users.UserItemRequestBuilder.Get"},
			EnumsUsed: []EnumUsage{
				{Enum: "models.DayOfWeek"},
			},
		},
	}

	groupResource := &ResourceInfo{
		ResourcePath: "internal/services/resources/groups/graph_beta/group",
		SDKDependencies: SDKDependencies{
			Types:         []string{"models.Group", "models.User"},
			MethodsCalled: []string{"models.NewGroup"},
		},
	}

	usage.TerraformResources["microsoft365_user"] = userResource
	usage.TerraformResources["microsoft365_group"] = groupResource

	stats := calculateStatistics(usage)

	t.Run("Statistics calculated correctly", func(t *testing.T) {
		assert.Equal(t, 2, stats.TotalResources)
		assert.Equal(t, 0, stats.TotalActions)
		assert.Equal(t, 2, stats.TotalSDKTypesUsed) // User, Group (User is deduplicated)
		assert.Equal(t, 1, stats.TotalEnumsTracked)
	})

	t.Run("SDK to resource index built", func(t *testing.T) {
		assert.Contains(t, usage.SDKToResourceIndex, "models.User")
		assert.Len(t, usage.SDKToResourceIndex["models.User"], 2) // Both resources use User
		assert.Contains(t, usage.SDKToResourceIndex["models.User"], "microsoft365_user")
		assert.Contains(t, usage.SDKToResourceIndex["models.User"], "microsoft365_group")

		assert.Contains(t, usage.SDKToResourceIndex, "models.Group")
		assert.Len(t, usage.SDKToResourceIndex["models.Group"], 1)
		assert.Contains(t, usage.SDKToResourceIndex["models.Group"], "microsoft365_group")
	})
}

func TestParseGoFile(t *testing.T) {
	t.Run("Valid Go source", func(t *testing.T) {
		source := `package test
import "fmt"

func main() {
	fmt.Println("test")
}
`
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, "test.go", source, parser.ParseComments)

		require.NoError(t, err)
		assert.NotNil(t, node)
		assert.Equal(t, "test", node.Name.Name)
	})

	t.Run("Invalid Go source returns error", func(t *testing.T) {
		source := `package test
this is not valid go code {{{}
`
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "invalid.go", source, parser.ParseComments)

		assert.Error(t, err)
	})
}
