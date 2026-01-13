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
	roleAssignments map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.roleAssignments = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("role_assignment", &RoleAssignmentMock{})
}

// RoleAssignmentMock provides mock responses for role assignment operations
type RoleAssignmentMock struct{}

// Ensure RoleAssignmentMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*RoleAssignmentMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for role assignment operations
// This implements the MockRegistrar interface
func (m *RoleAssignmentMock) RegisterMocks() {
	// POST /deviceManagement/roleAssignments - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/roleAssignments",
		m.createRoleAssignmentResponder())

	// GET /deviceManagement/roleAssignments/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleAssignments/([^/]+)$`,
		m.getRoleAssignmentResponder())

	// PATCH /deviceManagement/roleAssignments/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleAssignments/([^/]+)$`,
		m.updateRoleAssignmentResponder())

	// DELETE /deviceManagement/roleAssignments/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleAssignments/([^/]+)$`,
		m.deleteRoleAssignmentResponder())
}

// RegisterErrorMocks sets up error mock responses
// This implements the MockRegistrar interface
func (m *RoleAssignmentMock) RegisterErrorMocks() {
	// Error responses for testing error scenarios
}

// CleanupMockState clears all stored mock state
func (m *RoleAssignmentMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.roleAssignments = make(map[string]map[string]any)
}

// createRoleAssignmentResponder handles POST requests to create role assignments
func (m *RoleAssignmentMock) createRoleAssignmentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file - choose based on request content
		var jsonContent string
		var err error

		// Determine which response to load based on scope type
		if scopeType, exists := requestBody["scopeType"]; exists {
			switch scopeType {
			case "AllDevices":
				jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_all_devices.json"))
			case "AllLicensedUsers":
				jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_minimal.json"))
			default:
				jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_minimal.json"))
			}
		} else if resourceScopes, exists := requestBody["resourceScopes"]; exists {
			if scopes, ok := resourceScopes.([]any); ok && len(scopes) > 0 {
				jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_maximal.json"))
			} else {
				jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_minimal.json"))
			}
		} else {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_minimal.json"))
		}

		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
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
		if members, ok := requestBody["members"]; ok {
			response["members"] = members
		}
		if resourceScopes, ok := requestBody["resourceScopes"]; ok {
			response["resourceScopes"] = resourceScopes
		}
		if scopeType, ok := requestBody["scopeType"]; ok {
			response["scopeType"] = scopeType
		}
		// Preserve role definition odata.bind from request
		if roleDefBind, ok := requestBody["roleDefinition@odata.bind"]; ok {
			response["roleDefinition@odata.bind"] = roleDefBind
		}

		// Store in mock state
		mockState.Lock()
		mockState.roleAssignments[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getRoleAssignmentResponder handles GET requests to retrieve role assignments
func (m *RoleAssignmentMock) getRoleAssignmentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleAssignments/")

		mockState.Lock()
		roleAssignment, exists := mockState.roleAssignments[id]
		mockState.Unlock()

		if exists {
			return factories.SuccessResponse(200, roleAssignment)(req)
		}

		// Check for special test IDs
		switch {
		case strings.Contains(id, "minimal"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_minimal.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		case strings.Contains(id, "maximal"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_maximal.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		case strings.Contains(id, "all-devices"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_assignment_all_devices.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Role assignment not found"}}`), nil
		}
	}
}

// updateRoleAssignmentResponder handles PATCH requests to update role assignments
func (m *RoleAssignmentMock) updateRoleAssignmentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleAssignments/")

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		roleAssignment, exists := mockState.roleAssignments[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Role assignment not found"}}`), nil
		}

		// Update the stored role assignment
		for key, value := range requestBody {
			if key != "id" { // Don't allow ID changes
				roleAssignment[key] = value
			}
		}

		mockState.roleAssignments[id] = roleAssignment
		mockState.Unlock()

		return factories.SuccessResponse(200, roleAssignment)(req)
	}
}

// deleteRoleAssignmentResponder handles DELETE requests to delete role assignments
func (m *RoleAssignmentMock) deleteRoleAssignmentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleAssignments/")

		mockState.Lock()
		_, exists := mockState.roleAssignments[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Role assignment not found"}}`), nil
		}

		delete(mockState.roleAssignments, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	}
}
