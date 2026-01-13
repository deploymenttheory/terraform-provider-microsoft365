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
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	filteringPolicies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.filteringPolicies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("filtering_policy", &FilteringPolicyMock{})
}

// FilteringPolicyMock provides mock responses for Filtering Policy operations
type FilteringPolicyMock struct{}

// Ensure FilteringPolicyMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*FilteringPolicyMock)(nil)

// RegisterMocks registers HTTP mock responses for Filtering Policy operations
func (m *FilteringPolicyMock) RegisterMocks() {
	// License check endpoint
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	// POST /networkAccess/filteringPolicies - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringPolicies",
		m.createFilteringPolicyResponder())

	// GET /networkAccess/filteringPolicies/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/([^/]+)$`,
		m.getFilteringPolicyResponder())

	// PATCH /networkAccess/filteringPolicies/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/([^/]+)$`,
		m.updateFilteringPolicyResponder())

	// DELETE /networkAccess/filteringPolicies/{id} - Delete
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/([^/]+)$`,
		m.deleteFilteringPolicyResponder())
}

// createFilteringPolicyResponder handles POST requests to create filtering policies
func (m *FilteringPolicyMock) createFilteringPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_policy.json"))
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
		response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#filteringPolicies/$entity"

		// Update response with request data
		if name, ok := requestBody["name"]; ok {
			response["name"] = name
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		} else {
			response["description"] = nil
		}
		if action, ok := requestBody["action"]; ok {
			response["action"] = action
		}

		// Store in mock state
		mockState.Lock()
		mockState.filteringPolicies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getFilteringPolicyResponder handles GET requests to retrieve filtering policies
func (m *FilteringPolicyMock) getFilteringPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringPolicies/")

		mockState.Lock()
		policy, exists := mockState.filteringPolicies[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_policy_minimal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_policy_maximal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_policy_not_found.json"))
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
		policyCopy := make(map[string]any)
		for k, v := range policy {
			policyCopy[k] = v
		}

		// Ensure @odata.context is present in GET response
		if _, exists := policyCopy["@odata.context"]; !exists {
			policyCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#filteringPolicies/$entity"
		}

		return factories.SuccessResponse(200, policyCopy)(req)
	}
}

// updateFilteringPolicyResponder handles PATCH requests to update filtering policies
func (m *FilteringPolicyMock) updateFilteringPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringPolicies/")

		mockState.Lock()
		policy, exists := mockState.filteringPolicies[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_policy_not_found.json"))
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
			jsonContent, errLoad := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_policy_error.json"))
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
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_filtering_policy_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		var updatedPolicy map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedPolicy); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Start with existing data
		for k, v := range policy {
			updatedPolicy[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedPolicy[k] = v
		}

		// Update lastModifiedDateTime
		updatedPolicy["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.filteringPolicies[id] = updatedPolicy
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// deleteFilteringPolicyResponder handles DELETE requests to remove filtering policies
func (m *FilteringPolicyMock) deleteFilteringPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkAccess/filteringPolicies/")

		mockState.Lock()
		_, exists := mockState.filteringPolicies[id]
		if exists {
			delete(mockState.filteringPolicies, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_filtering_policy_not_found.json"))
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
func (m *FilteringPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Filtering Policies
	for id := range mockState.filteringPolicies {
		delete(mockState.filteringPolicies, id)
	}
}

// getLicenseResponder handles GET requests for license checking
func (m *FilteringPolicyMock) getLicenseResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "license", "get_subscribed_skus_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load license mock"}}`), nil
		}
		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *FilteringPolicyMock) RegisterErrorMocks() {
	// License check endpoint
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringPolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_filtering_policy_error.json"))
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
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Filtering policy is in use"))
}

// GetMockFilteringPolicyData returns sample filtering policy data for testing
func (m *FilteringPolicyMock) GetMockFilteringPolicyData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_policy_maximal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}

// GetMockFilteringPolicyMinimalData returns minimal filtering policy data for testing
func (m *FilteringPolicyMock) GetMockFilteringPolicyMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_filtering_policy_minimal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}
