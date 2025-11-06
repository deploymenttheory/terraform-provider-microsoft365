package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
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
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.filteringPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for individual Filtering Policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/[a-zA-Z0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			policy, exists := mockState.filteringPolicies[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			// Create response copy (deep copy to avoid concurrent modification issues)
			policyCopy := make(map[string]any)
			for k, v := range policy {
				policyCopy[k] = v
			}

			return httpmock.NewJsonResponse(200, policyCopy)
		})

	// Register POST for creating Filtering Policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringPolicies",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a mock ID
			id := "00000000-0000-0000-0000-000000000001"

			response := map[string]any{
				"id":                   id,
				"createdDateTime":      "2024-01-01T00:00:00Z",
				"lastModifiedDateTime": "2024-01-01T00:00:00Z",
				"version":              "1.0",
			}

			// Copy fields from request
			if name, ok := requestBody["name"].(string); ok {
				response["name"] = name
			}
			if description, ok := requestBody["description"].(string); ok {
				response["description"] = description
			}
			if action, ok := requestBody["action"].(string); ok {
				response["action"] = action
			}

			// Store in mock state
			mockState.Lock()
			mockState.filteringPolicies[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Filtering Policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/[a-zA-Z0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			policy, exists := mockState.filteringPolicies[id]
			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			// Create a new map for updated policy to avoid race conditions
			updatedPolicy := make(map[string]any)
			for k, v := range policy {
				updatedPolicy[k] = v
			}

			// Update fields from request
			if name, ok := requestBody["name"].(string); ok {
				updatedPolicy["name"] = name
			}
			if description, ok := requestBody["description"].(string); ok {
				updatedPolicy["description"] = description
			}
			if action, ok := requestBody["action"].(string); ok {
				updatedPolicy["action"] = action
			}
			updatedPolicy["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
			updatedPolicy["version"] = "1.1" // Increment version on update

			// Note: 'state' and 'priority' are not included here as they are properties used when
			// linking policies to security profiles, not direct properties of filtering policies.

			mockState.filteringPolicies[id] = updatedPolicy

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register DELETE for Filtering Policy
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/[a-zA-Z0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.filteringPolicies[id]
			if exists {
				delete(mockState.filteringPolicies, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
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

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *FilteringPolicyMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.filteringPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register error response for creating Filtering Policy with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkAccess/filteringPolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid filtering policy data"}}`), nil
		})

	// Register error response for Filtering Policy not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkAccess/filteringPolicies/[a-zA-Z0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})
}
