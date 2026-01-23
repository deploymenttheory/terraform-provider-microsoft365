package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeUsageMap(t *testing.T) {
	usage := initializeUsageMap()

	assert.NotNil(t, usage)
	assert.NotNil(t, usage.Packages)
	assert.NotNil(t, usage.Imports)
	assert.NotNil(t, usage.Types)
	assert.NotNil(t, usage.Methods)
	assert.NotNil(t, usage.Fields)
	assert.NotNil(t, usage.Enums)
	assert.Empty(t, usage.Packages)
	assert.Empty(t, usage.Enums)
}

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Regular Go file",
			path:     "/path/to/file.go",
			expected: false,
		},
		{
			name:     "Test file",
			path:     "/path/to/file_test.go",
			expected: true,
		},
		{
			name:     "Non-Go file",
			path:     "/path/to/file.txt",
			expected: true,
		},
		{
			name:     "Go file in subdirectory",
			path:     "/path/to/subdir/resource.go",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipFile(tt.path)
			assert.Equal(t, tt.expected, result)
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
			name:       "Microsoft Graph SDK",
			importPath: "github.com/microsoftgraph/msgraph-sdk-go/users",
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
		{
			name:       "Third-party package",
			importPath: "github.com/hashicorp/terraform-plugin-framework",
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

	t.Run("Complex path", func(t *testing.T) {
		imp := &ast.ImportSpec{
			Name: nil,
			Path: &ast.BasicLit{Value: `"github.com/microsoft/kiota-abstractions-go/serialization"`},
		}
		alias := getImportAlias(imp, "github.com/microsoft/kiota-abstractions-go/serialization")
		assert.Equal(t, "serialization", alias)
	})
}

func TestTrackImport(t *testing.T) {
	usage := initializeUsageMap()
	importPath := "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	filePath := "/path/to/file.go"

	trackImport(importPath, filePath, usage)

	assert.Equal(t, 1, usage.Packages[importPath])
	assert.Contains(t, usage.Imports[importPath], filePath)

	// Track same import from another file
	trackImport(importPath, "/path/to/other.go", usage)
	assert.Equal(t, 2, usage.Packages[importPath])
	assert.Len(t, usage.Imports[importPath], 2)
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

func TestTrackFieldAccess(t *testing.T) {
	usage := initializeUsageMap()
	typeName := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User"
	
	trackFieldAccess(typeName, "DisplayName", usage)
	trackFieldAccess(typeName, "DisplayName", usage)
	trackFieldAccess(typeName, "Email", usage)

	assert.Equal(t, 2, usage.Fields[typeName]["DisplayName"])
	assert.Equal(t, 1, usage.Fields[typeName]["Email"])
}

func TestTrackObjectMethod(t *testing.T) {
	usage := initializeUsageMap()
	typeName := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User"
	
	trackObjectMethod(typeName, "SetDisplayName", usage)
	trackObjectMethod(typeName, "SetDisplayName", usage)
	trackObjectMethod(typeName, "GetEmail", usage)

	assert.Equal(t, 2, usage.Methods[typeName+".SetDisplayName"])
	assert.Equal(t, 1, usage.Methods[typeName+".GetEmail"])
}

func TestTrackEnumUsage(t *testing.T) {
	usage := initializeUsageMap()
	importPath := "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	filePath := "/path/to/file.go"

	t.Run("Parse function creates enum entry", func(t *testing.T) {
		trackEnumUsage(importPath, "ParseRunAsAccountType", filePath, usage)
		
		expectedEnum := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.RunAsAccountType"
		assert.Contains(t, usage.Enums, expectedEnum)
		assert.Contains(t, usage.Enums[expectedEnum], filePath)
	})

	t.Run("Same enum from multiple files", func(t *testing.T) {
		trackEnumUsage(importPath, "ParseRunAsAccountType", "/other/file.go", usage)
		
		expectedEnum := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.RunAsAccountType"
		assert.Len(t, usage.Enums[expectedEnum], 2)
	})

	t.Run("Same file not duplicated", func(t *testing.T) {
		trackEnumUsage(importPath, "ParseRunAsAccountType", filePath, usage)
		
		expectedEnum := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.RunAsAccountType"
		assert.Len(t, usage.Enums[expectedEnum], 2) // Still 2, not 3
	})

	t.Run("Invalid parse function ignored", func(t *testing.T) {
		initialLen := len(usage.Enums)
		trackEnumUsage(importPath, "Parse", filePath, usage) // Empty enum type
		assert.Len(t, usage.Enums, initialLen)
	})
}

func TestTrackTypeInstantiation(t *testing.T) {
	usage := initializeUsageMap()
	typeName := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User"

	trackTypeInstantiation(typeName, usage)
	trackTypeInstantiation(typeName, usage)

	assert.Equal(t, 2, usage.Types[typeName]["_instantiated"])
}

func TestTrackStructFields(t *testing.T) {
	usage := initializeUsageMap()
	typeName := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User"

	// Create AST elements for struct fields
	elts := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "DisplayName"},
			Value: &ast.BasicLit{Value: `"John"`},
		},
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "Email"},
			Value: &ast.BasicLit{Value: `"john@example.com"`},
		},
		&ast.KeyValueExpr{
			Key:   &ast.Ident{Name: "DisplayName"}, // Duplicate
			Value: &ast.BasicLit{Value: `"Jane"`},
		},
	}

	trackStructFields(typeName, elts, usage)

	assert.Equal(t, 2, usage.Fields[typeName]["DisplayName"])
	assert.Equal(t, 1, usage.Fields[typeName]["Email"])
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

	t.Run("Unknown package", func(t *testing.T) {
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "unknown"},
				Sel: &ast.Ident{Name: "Type"},
			},
		}

		typeName := extractTypeFromCompositeLit(comp, fileImports)
		assert.Equal(t, "", typeName)
	})
}

func TestParseGoFile(t *testing.T) {
	t.Run("Valid Go source", func(t *testing.T) {
		// Create a temporary test file
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

	usage := initializeUsageMap()
	fileImports := collectImports(node, "/test.go", usage)

	t.Run("Returns correct aliases", func(t *testing.T) {
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models", fileImports["models"])
		assert.Equal(t, "github.com/microsoftgraph/msgraph-sdk-go/users", fileImports["users"])
		assert.NotContains(t, fileImports, "fmt") // Standard library filtered
		assert.NotContains(t, fileImports, "types") // Non-SDK package filtered
	})

	t.Run("Tracks SDK packages", func(t *testing.T) {
		assert.Equal(t, 1, usage.Packages["github.com/microsoftgraph/msgraph-beta-sdk-go/models"])
		assert.Equal(t, 1, usage.Packages["github.com/microsoftgraph/msgraph-sdk-go/users"])
		assert.NotContains(t, usage.Packages, "fmt")
	})

	t.Run("Tracks import locations", func(t *testing.T) {
		assert.Contains(t, usage.Imports["github.com/microsoftgraph/msgraph-beta-sdk-go/models"], "/test.go")
	})
}

func TestExtractSDKType(t *testing.T) {
	fileImports := map[string]string{
		"models": "github.com/microsoftgraph/msgraph-beta-sdk-go/models",
	}

	t.Run("Call expression", func(t *testing.T) {
		// models.NewUser()
		call := &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "NewUser"},
			},
		}

		typeName := extractSDKType(call, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.NewUser", typeName)
	})

	t.Run("Composite literal", func(t *testing.T) {
		// models.User{}
		comp := &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "models"},
				Sel: &ast.Ident{Name: "User"},
			},
		}

		typeName := extractSDKType(comp, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})

	t.Run("Unary expression (pointer)", func(t *testing.T) {
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

		typeName := extractSDKType(unary, fileImports)
		assert.Equal(t, "github.com/microsoftgraph/msgraph-beta-sdk-go/models.User", typeName)
	})

	t.Run("Non-SDK expression", func(t *testing.T) {
		expr := &ast.BasicLit{Value: "123"}
		typeName := extractSDKType(expr, fileImports)
		assert.Equal(t, "", typeName)
	})
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

func TestTrackEnumUsageIntegration(t *testing.T) {
	usage := initializeUsageMap()
	importPath := "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	testCases := []struct {
		methodName   string
		expectedEnum string
	}{
		{
			methodName:   "ParseRunAsAccountType",
			expectedEnum: "github.com/microsoftgraph/msgraph-beta-sdk-go/models.RunAsAccountType",
		},
		{
			methodName:   "ParseCloudPcRegionGroup",
			expectedEnum: "github.com/microsoftgraph/msgraph-beta-sdk-go/models.CloudPcRegionGroup",
		},
		{
			methodName:   "ParseInstallIntent",
			expectedEnum: "github.com/microsoftgraph/msgraph-beta-sdk-go/models.InstallIntent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.methodName, func(t *testing.T) {
			trackEnumUsage(importPath, tc.methodName, "/test.go", usage)
			assert.Contains(t, usage.Enums, tc.expectedEnum)
			assert.Contains(t, usage.Enums[tc.expectedEnum], "/test.go")
		})
	}
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

	t.Run("Non-pointer operator", func(t *testing.T) {
		unary := &ast.UnaryExpr{
			Op: token.NOT,
			X:  &ast.Ident{Name: "something"},
		}

		typeName := extractTypeFromUnaryExpr(unary, fileImports)
		assert.Equal(t, "", typeName)
	})
}

func TestTrackPackageMethod(t *testing.T) {
	usage := initializeUsageMap()
	importPath := "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	t.Run("Regular method", func(t *testing.T) {
		trackPackageMethod(importPath, "NewUser", "/test.go", usage)
		
		expectedMethod := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.NewUser"
		assert.Equal(t, 1, usage.Methods[expectedMethod])
	})

	t.Run("Enum parser method also tracks enum", func(t *testing.T) {
		trackPackageMethod(importPath, "ParseRunAsAccountType", "/test.go", usage)
		
		// Method tracked
		expectedMethod := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.ParseRunAsAccountType"
		assert.Equal(t, 1, usage.Methods[expectedMethod])
		
		// Enum tracked
		expectedEnum := "github.com/microsoftgraph/msgraph-beta-sdk-go/models.RunAsAccountType"
		assert.Contains(t, usage.Enums, expectedEnum)
	})
}
