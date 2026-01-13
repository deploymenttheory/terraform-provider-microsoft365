package mocks

import (
	"encoding/json"
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
	assignmentFilters map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.assignmentFilters = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("assignment_filter", &AssignmentFilterMock{})
}

// AssignmentFilterMock provides mock responses for assignment filter operations
type AssignmentFilterMock struct{}

// Ensure AssignmentFilterMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AssignmentFilterMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for assignment filter operations
// This implements the MockRegistrar interface
func (m *AssignmentFilterMock) RegisterMocks() {
	// License check endpoint
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	// POST /deviceManagement/assignmentFilters - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/assignmentFilters",
		m.createAssignmentFilterResponder())

	// GET /deviceManagement/assignmentFilters/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`,
		m.getAssignmentFilterResponder())

	// PATCH /deviceManagement/assignmentFilters/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`,
		m.updateAssignmentFilterResponder())

	// DELETE /deviceManagement/assignmentFilters/{id} - Delete
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`,
		m.deleteAssignmentFilterResponder())
}

// createAssignmentFilterResponder handles POST requests to create assignment filters
func (m *AssignmentFilterMock) createAssignmentFilterResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_assignment_filter.json"))
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
		if platform, ok := requestBody["platform"]; ok {
			response["platform"] = platform
		}
		if rule, ok := requestBody["rule"]; ok {
			response["rule"] = rule
		}
		if mgmtType, ok := requestBody["assignmentFilterManagementType"]; ok {
			response["assignmentFilterManagementType"] = mgmtType
		}
		if roleScopeTags, ok := requestBody["roleScopeTags"]; ok {
			response["roleScopeTags"] = roleScopeTags
		}

		// Store in mock state
		mockState.Lock()
		mockState.assignmentFilters[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getAssignmentFilterResponder handles GET requests to retrieve assignment filters
func (m *AssignmentFilterMock) getAssignmentFilterResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/assignmentFilters/")

		mockState.Lock()
		assignmentFilter, exists := mockState.assignmentFilters[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_assignment_filter_minimal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_assignment_filter_maximal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
				if err == nil {
					var errorResponse map[string]any
					if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
						return httpmock.NewJsonResponse(404, errorResponse)
					}
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}
		}

		return factories.SuccessResponse(200, assignmentFilter)(req)
	}
}

// updateAssignmentFilterResponder handles PATCH requests to update assignment filters
func (m *AssignmentFilterMock) updateAssignmentFilterResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/assignmentFilters/")

		mockState.Lock()
		assignmentFilter, exists := mockState.assignmentFilters[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
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
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_assignment_filter_error.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_assignment_filter_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var updatedFilter map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedFilter); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range assignmentFilter {
			updatedFilter[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedFilter[k] = v
		}

		// Update lastModifiedDateTime
		updatedFilter["lastModifiedDateTime"] = "2024-01-01T12:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.assignmentFilters[id] = updatedFilter
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedFilter)(req)
	}
}

// deleteAssignmentFilterResponder handles DELETE requests to remove assignment filters
func (m *AssignmentFilterMock) deleteAssignmentFilterResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/assignmentFilters/")

		mockState.Lock()
		_, exists := mockState.assignmentFilters[id]
		if exists {
			delete(mockState.assignmentFilters, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
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

// getLicenseResponder handles GET requests for license checking
func (m *AssignmentFilterMock) getLicenseResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "license", "get_subscribed_skus_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load license mock"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse license mock"}}`), nil
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *AssignmentFilterMock) RegisterErrorMocks() {
	// License check endpoint
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/assignmentFilters",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_assignment_filter_error.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`), nil
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Assignment filter is in use"))
}

// CleanupMockState clears all stored mock state
func (m *AssignmentFilterMock) CleanupMockState() {
	mockState.Lock()
	mockState.assignmentFilters = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockAssignmentFilterData returns sample assignment filter data for testing
func (m *AssignmentFilterMock) GetMockAssignmentFilterData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_assignment_filter_maximal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse mock response: " + err.Error())
	}
	return response
}

// GetMockAssignmentFilterMinimalData returns minimal assignment filter data for testing
func (m *AssignmentFilterMock) GetMockAssignmentFilterMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_assignment_filter_minimal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse mock response: " + err.Error())
	}
	return response
}
