package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	assignmentFilters map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.assignmentFilters = make(map[string]map[string]interface{})

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
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`,
		m.deleteAssignmentFilterResponder())
}

// loadJSONResponse loads a JSON response file and returns its contents
func (m *AssignmentFilterMock) loadJSONResponse(filepath string) (map[string]interface{}, error) {
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// createAssignmentFilterResponder handles POST requests to create assignment filters
func (m *AssignmentFilterMock) createAssignmentFilterResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		response, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "post_assignment_filter.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
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
				response, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "get_assignment_filter_minimal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			case strings.Contains(id, "maximal"):
				response, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "get_assignment_filter_maximal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			default:
				errorResponse, _ := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
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
			errorResponse, _ := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "post_assignment_filter_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Load update template
		updatedFilter, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_update", "get_assignment_filter_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
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
			errorResponse, _ := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_delete", "get_assignment_filter_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *AssignmentFilterMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/assignmentFilters",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "post_assignment_filter_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Assignment filter is in use"))
}

// CleanupMockState clears all stored mock state
func (m *AssignmentFilterMock) CleanupMockState() {
	mockState.Lock()
	mockState.assignmentFilters = make(map[string]map[string]interface{})
	mockState.Unlock()
}

// GetMockAssignmentFilterData returns sample assignment filter data for testing
func (m *AssignmentFilterMock) GetMockAssignmentFilterData() map[string]interface{} {
	response, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "get_assignment_filter_maximal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]interface{}{
			"id":                             "test-assignment-filter-id",
			"displayName":                    "Test Assignment Filter",
			"description":                    "Test assignment filter for unit testing",
			"platform":                       "windows10AndLater",
			"rule":                           "(device.osVersion -startsWith \"10.0\")",
			"assignmentFilterManagementType": "devices",
			"createdDateTime":                "2024-01-01T00:00:00Z",
			"lastModifiedDateTime":           "2024-01-01T00:00:00Z",
			"roleScopeTags":                  []string{"0"},
		}
	}
	return response
}

// GetMockAssignmentFilterMinimalData returns minimal assignment filter data for testing
func (m *AssignmentFilterMock) GetMockAssignmentFilterMinimalData() map[string]interface{} {
	response, err := m.loadJSONResponse(filepath.Join("mocks", "responses", "validate_create", "get_assignment_filter_minimal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]interface{}{
			"id":                   "test-minimal-assignment-filter-id",
			"displayName":          "Test Minimal Assignment Filter",
			"platform":             "windows10AndLater",
			"rule":                 "(device.osVersion -startsWith \"10.0\")",
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
			"roleScopeTags":        []string{"0"},
		}
	}
	return response
}
