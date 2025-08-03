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
	deviceCategories map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.deviceCategories = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("device_category", &DeviceCategoryMock{})
}

// DeviceCategoryMock provides mock responses for device category operations
type DeviceCategoryMock struct{}

// Ensure DeviceCategoryMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*DeviceCategoryMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for device category operations
// This implements the MockRegistrar interface
func (m *DeviceCategoryMock) RegisterMocks() {
	// POST /deviceManagement/deviceCategories - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCategories",
		m.createDeviceCategoryResponder())

	// GET /deviceManagement/deviceCategories/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/([^/]+)$`,
		m.getDeviceCategoryResponder())

	// PATCH /deviceManagement/deviceCategories/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/([^/]+)$`,
		m.updateDeviceCategoryResponder())

	// DELETE /deviceManagement/deviceCategories/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/([^/]+)$`,
		m.deleteDeviceCategoryResponder())
}

// RegisterErrorMocks sets up mock HTTP responders that return error responses
func (m *DeviceCategoryMock) RegisterErrorMocks() {
	// POST /deviceManagement/deviceCategories - Create Error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCategories",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_device_category_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// GET /deviceManagement/deviceCategories/{id} - Read Error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_device_category_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *DeviceCategoryMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored device categories
	for id := range mockState.deviceCategories {
		delete(mockState.deviceCategories, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *DeviceCategoryMock) loadJSONResponse(filePath string) (map[string]interface{}, error) {
	var response map[string]interface{}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// createDeviceCategoryResponder handles POST requests to create device categories
func (m *DeviceCategoryMock) createDeviceCategoryResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_device_category_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Generate a new UUID for the created resource
		id := uuid.New().String()

		// Create response with request data
		response := map[string]interface{}{
			"id": id,
		}

		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		}
		if roleScopeTagIds, ok := requestBody["roleScopeTagIds"]; ok {
			response["roleScopeTagIds"] = roleScopeTagIds
		} else {
			response["roleScopeTagIds"] = []string{"0"}
		}

		// Store in mock state
		mockState.Lock()
		mockState.deviceCategories[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getDeviceCategoryResponder handles GET requests to retrieve device categories
func (m *DeviceCategoryMock) getDeviceCategoryResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/deviceCategories/")

		mockState.Lock()
		deviceCategory, exists := mockState.deviceCategories[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_device_category_minimal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			case strings.Contains(id, "maximal"):
				response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_device_category_maximal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			default:
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_device_category_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}
		}

		return factories.SuccessResponse(200, deviceCategory)(req)
	}
}

// updateDeviceCategoryResponder handles PATCH requests to update device categories
func (m *DeviceCategoryMock) updateDeviceCategoryResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/deviceCategories/")

		mockState.Lock()
		deviceCategory, exists := mockState.deviceCategories[id]
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_device_category_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_device_category_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Load update template
		updatedCategory, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "get_device_category_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range deviceCategory {
			updatedCategory[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedCategory[k] = v
		}

		// Store updated state
		mockState.Lock()
		mockState.deviceCategories[id] = updatedCategory
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedCategory)(req)
	}
}

// deleteDeviceCategoryResponder handles DELETE requests to delete device categories
func (m *DeviceCategoryMock) deleteDeviceCategoryResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/deviceCategories/")

		mockState.Lock()
		_, exists := mockState.deviceCategories[id]
		if exists {
			delete(mockState.deviceCategories, id)
		}
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_device_category_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		// Return 204 No Content for successful deletion
		return httpmock.NewStringResponse(204, ""), nil
	}
}
