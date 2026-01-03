package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	roleDefinitions map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.roleDefinitions = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("role_definition", &RoleDefinitionMock{})
}

// RoleDefinitionMock provides mock responses for role definition operations
type RoleDefinitionMock struct{}

// Ensure RoleDefinitionMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*RoleDefinitionMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for role definition operations
// This implements the MockRegistrar interface
func (m *RoleDefinitionMock) RegisterMocks() {
	// POST /roleManagement/cloudPC/roleDefinitions - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions",
		m.createRoleDefinitionResponder())

	// GET /roleManagement/cloudPC/roleDefinitions/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/([^/]+)$`,
		m.getRoleDefinitionResponder())

	// PATCH /roleManagement/cloudPC/roleDefinitions/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/([^/]+)$`,
		m.updateRoleDefinitionResponder())

	// DELETE /roleManagement/cloudPC/roleDefinitions/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/([^/]+)$`,
		m.deleteRoleDefinitionResponder())

	// GET /roleManagement/cloudPC/roleDefinitions - List (for uniqueness validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions",
		m.listRoleDefinitionsResponder())

}

// createRoleDefinitionResponder handles POST requests to create role definitions
func (m *RoleDefinitionMock) createRoleDefinitionResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file - choose based on request content
		var response map[string]any
		var jsonContent string
		var err error

		if _, hasRolePermissions := requestBody["rolePermissions"]; !hasRolePermissions {
			// No role permissions specified
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_no_permissions.json"))
		} else if isBuiltIn, ok := requestBody["isBuiltIn"].(bool); ok && isBuiltIn {
			// Built-in role definition
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_maximal_builtin.json"))
		} else if isBuiltInRoleDefinition, ok := requestBody["isBuiltInRoleDefinition"].(bool); ok && isBuiltInRoleDefinition {
			// Built-in role definition (alternative field)
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_maximal_builtin.json"))
		} else if description, hasDesc := requestBody["description"]; hasDesc && description != "" {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_maximal.json"))
		} else {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_minimal.json"))
		}

		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Generate a new ID for the created resource
		id := uuid.New().String()
		response["id"] = id

		// Update response with request data
		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		}
		if isBuiltIn, ok := requestBody["isBuiltIn"]; ok {
			response["isBuiltIn"] = isBuiltIn
		}
		if isBuiltInRoleDefinition, ok := requestBody["isBuiltInRoleDefinition"]; ok {
			response["isBuiltInRoleDefinition"] = isBuiltInRoleDefinition
		}
		if builtInRoleName, ok := requestBody["builtInRoleName"]; ok {
			response["builtInRoleName"] = builtInRoleName
		}
		if rolePermissions, ok := requestBody["rolePermissions"]; ok {
			response["rolePermissions"] = rolePermissions
		}
		if roleScopeTagIds, ok := requestBody["roleScopeTagIds"]; ok {
			response["roleScopeTagIds"] = roleScopeTagIds
		}

		// Store in mock state
		mockState.Lock()
		mockState.roleDefinitions[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getRoleDefinitionResponder handles GET requests to retrieve role definitions
func (m *RoleDefinitionMock) getRoleDefinitionResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/roleManagement/cloudPC/roleDefinitions/")

		mockState.Lock()
		roleDefinition, exists := mockState.roleDefinitions[id]
		mockState.Unlock()

		if exists {
			return factories.SuccessResponse(200, roleDefinition)(req)
		}

		// Check for special test IDs
		switch {
		case strings.Contains(id, "minimal"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_minimal.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		case strings.Contains(id, "maximal"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_maximal.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		default:
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_definition_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		}
	}
}

// updateRoleDefinitionResponder handles PATCH requests to update role definitions
func (m *RoleDefinitionMock) updateRoleDefinitionResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/roleManagement/cloudPC/roleDefinitions/")

		mockState.Lock()
		roleDefinition, exists := mockState.roleDefinitions[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_definition_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid JSON",
				},
			})
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_role_definition_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		var updatedRoleDefinition map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedRoleDefinition); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Start with existing data
		for k, v := range roleDefinition {
			updatedRoleDefinition[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedRoleDefinition[k] = v
		}

		// Store updated version
		mockState.Lock()
		mockState.roleDefinitions[id] = updatedRoleDefinition
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedRoleDefinition)(req)
	}
}

// deleteRoleDefinitionResponder handles DELETE requests to remove role definitions
func (m *RoleDefinitionMock) deleteRoleDefinitionResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/roleManagement/cloudPC/roleDefinitions/")

		mockState.Lock()
		_, exists := mockState.roleDefinitions[id]
		if exists {
			delete(mockState.roleDefinitions, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_definition_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *RoleDefinitionMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid role definition data",
				},
			})
		})

	// GET - List error (for uniqueness validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, map[string]any{
				"error": map[string]any{
					"code":    "InternalServerError",
					"message": "Failed to retrieve existing role definitions",
				},
			})
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/roleManagement/cloudPC/roleDefinitions/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Role definition is in use"))

}

// CleanupMockState clears all stored mock state
func (m *RoleDefinitionMock) CleanupMockState() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockRoleDefinitionData returns sample role definition data for testing
func (m *RoleDefinitionMock) GetMockRoleDefinitionData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_maximal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                      "test-role-definition-id",
			"displayName":             "Test Role Definition",
			"description":             "Test role definition for unit testing",
			"isBuiltIn":               false,
			"isBuiltInRoleDefinition": false,
			"roleScopeTagIds":         []string{"0"},
			"rolePermissions": []map[string]any{
				{
					"allowedResourceActions": []string{
						"microsoft.management/managedDevices/read",
						"microsoft.management/managedDevices/write",
					},
				},
			},
		}
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}

// listRoleDefinitionsResponder handles GET requests to list role definitions (for uniqueness validation)
func (m *RoleDefinitionMock) listRoleDefinitionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Return empty collection for unit tests to avoid name conflicts
		response := map[string]any{
			"value": []any{},
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// listResourceOperationsResponder handles GET requests to list resource operations (for role permission validation)
func (m *RoleDefinitionMock) listResourceOperationsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Return a list of valid resource operations that match what's used in tests
		response := map[string]any{
			"value": []map[string]any{
				// Legacy format for backward compatibility
				{
					"id":          "microsoft.management/managedDevices/read",
					"actionName":  "Read managed devices",
					"description": "Allows reading of managed device information",
				},
				{
					"id":          "microsoft.management/managedDevices/write",
					"actionName":  "Write managed devices",
					"description": "Allows modification of managed device settings",
				},
				{
					"id":          "microsoft.management/managedDevices/delete",
					"actionName":  "Delete managed devices",
					"description": "Allows deletion of managed devices",
				},
				{
					"id":          "microsoft.management/deviceConfigurations/read",
					"actionName":  "Read device configurations",
					"description": "Allows reading of device configuration policies",
				},
				{
					"id":          "microsoft.management/deviceConfigurations/write",
					"actionName":  "Write device configurations",
					"description": "Allows modification of device configuration policies",
				},
				// Proper Microsoft Intune format
				{
					"id":          "Microsoft.Intune_ManagedDevices_Read",
					"actionName":  "Read managed devices",
					"description": "Allows reading of managed device information",
				},
				{
					"id":          "Microsoft.Intune_ManagedDevices_Update",
					"actionName":  "Update managed devices",
					"description": "Allows modification of managed device settings",
				},
				{
					"id":          "Microsoft.Intune_ManagedDevices_Delete",
					"actionName":  "Delete managed devices",
					"description": "Allows deletion of managed devices",
				},
				{
					"id":          "Microsoft.Intune_DeviceConfigurations_Read",
					"actionName":  "Read device configurations",
					"description": "Allows reading of device configuration policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceConfigurations_Create",
					"actionName":  "Create device configurations",
					"description": "Allows creation of device configuration policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceConfigurations_Update",
					"actionName":  "Update device configurations",
					"description": "Allows modification of device configuration policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceConfigurations_Delete",
					"actionName":  "Delete device configurations",
					"description": "Allows deletion of device configuration policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceConfigurations_Assign",
					"actionName":  "Assign device configurations",
					"description": "Allows assignment of device configuration policies",
				},
				{
					"id":          "Microsoft.Intune_Audit_Read",
					"actionName":  "Read audit logs",
					"description": "Allows reading of audit information",
				},
				{
					"id":          "Microsoft.Intune_Organization_Read",
					"actionName":  "Read organization",
					"description": "Allows reading of organization information",
				},
				// Additional permissions for built-in roles
				{
					"id":          "Microsoft.Intune_DeviceCompliancePolices_Read",
					"actionName":  "Read device compliance policies",
					"description": "Allows reading of device compliance policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceCompliancePolices_Create",
					"actionName":  "Create device compliance policies",
					"description": "Allows creation of device compliance policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceCompliancePolices_Update",
					"actionName":  "Update device compliance policies",
					"description": "Allows modification of device compliance policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceCompliancePolices_Delete",
					"actionName":  "Delete device compliance policies",
					"description": "Allows deletion of device compliance policies",
				},
				{
					"id":          "Microsoft.Intune_DeviceCompliancePolices_Assign",
					"actionName":  "Assign device compliance policies",
					"description": "Allows assignment of device compliance policies",
				},
				{
					"id":          "Microsoft.Intune_MobileApps_Read",
					"actionName":  "Read mobile apps",
					"description": "Allows reading of mobile applications",
				},
				{
					"id":          "Microsoft.Intune_MobileApps_Create",
					"actionName":  "Create mobile apps",
					"description": "Allows creation of mobile applications",
				},
				{
					"id":          "Microsoft.Intune_MobileApps_Update",
					"actionName":  "Update mobile apps",
					"description": "Allows modification of mobile applications",
				},
				{
					"id":          "Microsoft.Intune_MobileApps_Delete",
					"actionName":  "Delete mobile apps",
					"description": "Allows deletion of mobile applications",
				},
				{
					"id":          "Microsoft.Intune_MobileApps_Assign",
					"actionName":  "Assign mobile apps",
					"description": "Allows assignment of mobile applications",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_Wipe",
					"actionName":  "Wipe devices",
					"description": "Allows remote wiping of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_Retire",
					"actionName":  "Retire devices",
					"description": "Allows remote retiring of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_SyncDevice",
					"actionName":  "Sync devices",
					"description": "Allows remote syncing of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_RebootNow",
					"actionName":  "Reboot devices",
					"description": "Allows remote rebooting of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_ShutDown",
					"actionName":  "Shutdown devices",
					"description": "Allows remote shutdown of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_RemoteLock",
					"actionName":  "Remote lock devices",
					"description": "Allows remote locking of devices",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_ResetPasscode",
					"actionName":  "Reset device passcode",
					"description": "Allows resetting device passcodes",
				},
				{
					"id":          "Microsoft.Intune_RemoteTasks_LocateDevice",
					"actionName":  "Locate devices",
					"description": "Allows locating devices remotely",
				},
			},
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// GetMockRoleDefinitionMinimalData returns minimal role definition data for testing
func (m *RoleDefinitionMock) GetMockRoleDefinitionMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_definition_minimal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                      "test-minimal-role-definition-id",
			"displayName":             "Test Minimal Role Definition",
			"description":             "",
			"isBuiltIn":               false,
			"isBuiltInRoleDefinition": false,
			"roleScopeTagIds":         []string{"0"},
			"rolePermissions": []map[string]any{
				{
					"allowedResourceActions": []string{
						"microsoft.management/managedDevices/read",
						"microsoft.management/managedDevices/write",
					},
				},
			},
		}
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}
