package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	roleScopeTags map[string]map[string]any
	assignments   map[string][]map[string]any
}

func init() {
	// Initialize mockState
	mockState.roleScopeTags = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("role_scope_tag", &RoleScopeTagMock{})
}

// RoleScopeTagMock provides mock responses for role scope tag operations
type RoleScopeTagMock struct{}

// Ensure RoleScopeTagMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*RoleScopeTagMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for role scope tag operations
// This implements the MockRegistrar interface
func (m *RoleScopeTagMock) RegisterMocks() {
	// POST /deviceManagement/roleScopeTags - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/roleScopeTags",
		m.createRoleScopeTagResponder())

	// GET /deviceManagement/roleScopeTags/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`,
		m.getRoleScopeTagResponder())

	// PATCH /deviceManagement/roleScopeTags/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`,
		m.updateRoleScopeTagResponder())

	// DELETE /deviceManagement/roleScopeTags/{id} - Delete
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`,
		m.deleteRoleScopeTagResponder())

	// GET /deviceManagement/roleScopeTags - List (for validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/roleScopeTags",
		m.listRoleScopeTagsResponder())

	// POST /deviceManagement/roleScopeTags/{id}/assign - Assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)/assign$`,
		m.assignRoleScopeTagResponder())

	// GET /deviceManagement/roleScopeTags/{id}/assignments - Get assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)/assignments$`,
		m.getAssignmentsResponder())
}

// createRoleScopeTagResponder handles POST requests to create role scope tags
func (m *RoleScopeTagMock) createRoleScopeTagResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Check for display name conflicts
		displayName, ok := requestBody["displayName"].(string)
		if ok && displayName != "" {
			mockState.Lock()
			for _, existingTag := range mockState.roleScopeTags {
				if existingDisplayName, exists := existingTag["displayName"].(string); exists && existingDisplayName == displayName {
					mockState.Unlock()
					return httpmock.NewJsonResponse(400, map[string]any{
						"error": map[string]any{
							"code":    "BadRequest",
							"message": "Role scope tag with display name '" + displayName + "' already exists. Display names must be unique",
						},
					})
				}
			}
			mockState.Unlock()
		}

		// Load base response from JSON file - use minimal if no description provided
		var jsonContent string
		var err error
		if description, hasDesc := requestBody["description"]; hasDesc && description != "" {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_maximal.json"))
		} else {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_minimal.json"))
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
		} else {
			response["description"] = ""
		}

		// Store in mock state
		mockState.Lock()
		mockState.roleScopeTags[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getRoleScopeTagResponder handles GET requests to retrieve role scope tags
func (m *RoleScopeTagMock) getRoleScopeTagResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleScopeTags/")

		mockState.Lock()
		roleScopeTag, exists := mockState.roleScopeTags[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_minimal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_maximal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_scope_tag_not_found.json"))
				if err == nil {
					var errorResponse map[string]any
					if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
						return httpmock.NewJsonResponse(404, errorResponse)
					}
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}
		}

		return factories.SuccessResponse(200, roleScopeTag)(req)
	}
}

// updateRoleScopeTagResponder handles PATCH requests to update role scope tags
func (m *RoleScopeTagMock) updateRoleScopeTagResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleScopeTags/")

		mockState.Lock()
		roleScopeTag, exists := mockState.roleScopeTags[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_scope_tag_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
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

		// Check for display name conflicts (excluding current resource)
		if displayName, ok := requestBody["displayName"].(string); ok && displayName != "" {
			mockState.Lock()
			for existingId, existingTag := range mockState.roleScopeTags {
				if existingId != id {
					if existingDisplayName, exists := existingTag["displayName"].(string); exists && existingDisplayName == displayName {
						mockState.Unlock()
						return httpmock.NewJsonResponse(400, map[string]any{
							"error": map[string]any{
								"code":    "BadRequest",
								"message": "Role scope tag with display name '" + displayName + "' already exists. Display names must be unique",
							},
						})
					}
				}
			}
			mockState.Unlock()
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_role_scope_tag_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var updatedTag map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedTag); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range roleScopeTag {
			updatedTag[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedTag[k] = v
		}

		// Update lastModifiedDateTime
		updatedTag["lastModifiedDateTime"] = "2024-01-01T12:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.roleScopeTags[id] = updatedTag
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedTag)(req)
	}
}

// deleteRoleScopeTagResponder handles DELETE requests to remove role scope tags
func (m *RoleScopeTagMock) deleteRoleScopeTagResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/roleScopeTags/")

		mockState.Lock()
		_, exists := mockState.roleScopeTags[id]
		if exists {
			delete(mockState.roleScopeTags, id)
			delete(mockState.assignments, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_scope_tag_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// listRoleScopeTagsResponder handles GET requests to list all role scope tags (for validation)
func (m *RoleScopeTagMock) listRoleScopeTagsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		allTags := make([]map[string]any, 0, len(mockState.roleScopeTags))
		for _, tag := range mockState.roleScopeTags {
			allTags = append(allTags, tag)
		}
		mockState.Unlock()

		response := map[string]any{
			"value": allTags,
		}

		return factories.SuccessResponse(200, response)(req)
	}
}

// assignRoleScopeTagResponder handles POST requests to assign role scope tags
func (m *RoleScopeTagMock) assignRoleScopeTagResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL - for /deviceManagement/roleScopeTags/{id}/assign
		urlPath := req.URL.Path
		// Remove the prefix and suffix to get just the ID
		idStartIndex := strings.Index(urlPath, "/deviceManagement/roleScopeTags/") + len("/deviceManagement/roleScopeTags/")
		idEndIndex := strings.Index(urlPath[idStartIndex:], "/assign")
		if idEndIndex == -1 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL format"}}`), nil
		}
		id := urlPath[idStartIndex : idStartIndex+idEndIndex]

		mockState.Lock()
		_, exists := mockState.roleScopeTags[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_role_scope_tag_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Store assignments in mock state
		if assignments, ok := requestBody["assignments"].([]any); ok {
			mockAssignments := make([]map[string]any, 0, len(assignments))
			for i, assignment := range assignments {
				if assignmentMap, ok := assignment.(map[string]any); ok {
					// Add mock assignment ID
					assignmentMap["id"] = uuid.New().String()
					mockAssignments = append(mockAssignments, assignmentMap)
				} else {
					return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Invalid assignment at index %d"}}`, i)), nil
				}
			}

			mockState.Lock()
			mockState.assignments[id] = mockAssignments
			mockState.Unlock()
		}

		// Load and return assignment response
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_assignments", "post_role_scope_tag_assign.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		return factories.SuccessResponse(200, response)(req)
	}
}

// getAssignmentsResponder handles GET requests to retrieve role scope tag assignments
func (m *RoleScopeTagMock) getAssignmentsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL - for /deviceManagement/roleScopeTags/{id}/assignments
		urlPath := req.URL.Path
		// Remove the prefix and suffix to get just the ID
		idStartIndex := strings.Index(urlPath, "/deviceManagement/roleScopeTags/") + len("/deviceManagement/roleScopeTags/")
		idEndIndex := strings.Index(urlPath[idStartIndex:], "/assignments")
		if idEndIndex == -1 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL format"}}`), nil
		}
		id := urlPath[idStartIndex : idStartIndex+idEndIndex]

		mockState.Lock()
		assignments, exists := mockState.assignments[id]
		mockState.Unlock()

		if !exists {
			// Return empty assignments response
			response := map[string]any{
				"value": []any{},
			}
			return factories.SuccessResponse(200, response)(req)
		}

		// Load assignments template and merge with stored assignments
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_assignments", "get_role_scope_tag_assignments.json"))
		var response map[string]any
		if err != nil {
			// Fallback to empty response
			response = map[string]any{
				"value": assignments,
			}
		} else {
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				// Fallback to empty response if parsing fails
				response = map[string]any{
					"value": assignments,
				}
			} else {
				response["value"] = assignments
			}
		}

		return factories.SuccessResponse(200, response)(req)
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *RoleScopeTagMock) RegisterErrorMocks() {
	// GET /deviceManagement/roleScopeTags - List (for validation) - return error
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/roleScopeTags",
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/roleScopeTags",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid role scope tag data",
				},
			})
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Role scope tag is in use"))
}

// CleanupMockState clears all stored mock state
func (m *RoleScopeTagMock) CleanupMockState() {
	mockState.Lock()
	mockState.roleScopeTags = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]map[string]any)
	mockState.Unlock()
}

// GetMockRoleScopeTagData returns sample role scope tag data for testing
func (m *RoleScopeTagMock) GetMockRoleScopeTagData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_maximal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                   "test-role-scope-tag-id",
			"displayName":          "Test Role Scope Tag",
			"description":          "Test role scope tag for unit testing",
			"isBuiltIn":            false,
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		// Fallback to hardcoded response if parsing fails
		return map[string]any{
			"id":                   "test-role-scope-tag-id",
			"displayName":          "Test Role Scope Tag",
			"description":          "Test role scope tag for unit testing",
			"isBuiltIn":            false,
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}
	return response
}

// GetMockRoleScopeTagMinimalData returns minimal role scope tag data for testing
func (m *RoleScopeTagMock) GetMockRoleScopeTagMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_role_scope_tag_minimal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                   "test-minimal-role-scope-tag-id",
			"displayName":          "Test Minimal Role Scope Tag",
			"description":          "",
			"isBuiltIn":            false,
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		// Fallback to hardcoded response if parsing fails
		return map[string]any{
			"id":                   "test-minimal-role-scope-tag-id",
			"displayName":          "Test Minimal Role Scope Tag",
			"description":          "",
			"isBuiltIn":            false,
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}
	return response
}
