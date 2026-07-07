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
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	filteringProfiles map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.filteringProfiles = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("filtering_profile", &FilteringProfileMock{})
}

// FilteringProfileMock provides mock responses for Filtering Profile operations
type FilteringProfileMock struct{}

// Ensure FilteringProfileMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*FilteringProfileMock)(nil)

// RegisterMocks registers HTTP mock responses for Filtering Profile operations
func (m *FilteringProfileMock) RegisterMocks() {
	// POST /networkAccess/filteringProfiles - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringProfiles",
		m.createFilteringProfileResponder())

	// GET /networkAccess/filteringProfiles/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/]+)$`,
		m.getFilteringProfileResponder())

	// PATCH /networkAccess/filteringProfiles/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/]+)$`,
		m.updateFilteringProfileResponder())

	// DELETE /networkAccess/filteringProfiles/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/([^/]+)$`,
		m.deleteFilteringProfileResponder())
}

// createFilteringProfileResponder handles POST requests to create filtering profiles
func (m *FilteringProfileMock) createFilteringProfileResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_profile.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Generate a mock ID
		id := "00000000-0000-0000-0000-000000000001"
		response["id"] = id

		// Add @odata.context
		response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#filteringProfiles/$entity"

		// Update response with request data
		if name, ok := requestBody["name"]; ok {
			response["name"] = name
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		} else {
			response["description"] = nil
		}
		if priority, ok := requestBody["priority"]; ok {
			response["priority"] = priority
		} else {
			response["priority"] = float64(100)
		}
		if state, ok := requestBody["state"]; ok {
			response["state"] = state
		} else {
			response["state"] = "enabled"
		}

		// Store in mock state
		mockState.Lock()
		mockState.filteringProfiles[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getFilteringProfileResponder handles GET requests to retrieve filtering profiles
func (m *FilteringProfileMock) getFilteringProfileResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringProfiles/")

		mockState.Lock()
		profile, exists := mockState.filteringProfiles[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_profile_minimal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_profile_maximal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_profile_not_found.json"))
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

		// Create response copy (deep copy to avoid concurrent modification issues)
		profileCopy := make(map[string]any)
		for k, v := range profile {
			profileCopy[k] = v
		}

		// Ensure @odata.context is present in GET response
		if _, exists := profileCopy["@odata.context"]; !exists {
			profileCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#networkAccess/filteringProfiles/$entity"
		}

		return factories.SuccessResponse(200, profileCopy)(req)
	}
}

// updateFilteringProfileResponder handles PATCH requests to update filtering profiles
func (m *FilteringProfileMock) updateFilteringProfileResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringProfiles/")

		mockState.Lock()
		profile, exists := mockState.filteringProfiles[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_profile_not_found.json"))
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
			jsonContent, errLoad := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_profile_error.json"))
			if errLoad != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_filtering_profile_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		var updatedProfile map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedProfile); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Start with existing data
		for k, v := range profile {
			updatedProfile[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedProfile[k] = v
		}

		// Update lastModifiedDateTime
		updatedProfile["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.filteringProfiles[id] = updatedProfile
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// deleteFilteringProfileResponder handles DELETE requests to remove filtering profiles
func (m *FilteringProfileMock) deleteFilteringProfileResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringProfiles/")

		mockState.Lock()
		_, exists := mockState.filteringProfiles[id]
		if exists {
			delete(mockState.filteringProfiles, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_profile_not_found.json"))
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

// CleanupMockState clears the mock state for clean test runs
func (m *FilteringProfileMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Filtering Profiles
	for id := range mockState.filteringProfiles {
		delete(mockState.filteringProfiles, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *FilteringProfileMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringProfiles",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_profile_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringProfiles/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Filtering profile is in use"))
}

// GetMockFilteringProfileData returns sample filtering profile data for testing
func (m *FilteringProfileMock) GetMockFilteringProfileData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_profile_maximal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}

// GetMockFilteringProfileMinimalData returns minimal filtering profile data for testing
func (m *FilteringProfileMock) GetMockFilteringProfileMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_profile_minimal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}
