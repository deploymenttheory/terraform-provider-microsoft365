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
	authenticationStrengths map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.authenticationStrengths = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("authentication_strength", &AuthenticationStrengthMock{})
}

// AuthenticationStrengthMock provides mock responses for Authentication Strength Policy operations
type AuthenticationStrengthMock struct{}

// Ensure AuthenticationStrengthMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AuthenticationStrengthMock)(nil)

// RegisterMocks registers HTTP mock responses for Authentication Strength Policy operations
func (m *AuthenticationStrengthMock) RegisterMocks() {
	// GET /identity/conditionalAccess/authenticationStrength/policies - List (for validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies",
		m.listAuthenticationStrengthPoliciesResponder())

	// POST /identity/conditionalAccess/authenticationStrength/policies - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies",
		m.createAuthenticationStrengthPolicyResponder())

	// GET /identity/conditionalAccess/authenticationStrength/policies/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F-]+$`,
		m.getAuthenticationStrengthPolicyResponder())

	// PATCH /identity/conditionalAccess/authenticationStrength/policies/{id} - Update basic fields
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F-]+$`,
		m.updateAuthenticationStrengthPolicyResponder())

	// POST /policies/authenticationStrengthPolicies/{id}/updateAllowedCombinations - Update allowed combinations
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/policies/authenticationStrengthPolicies/[0-9a-fA-F-]+/updateAllowedCombinations$`,
		m.updateAllowedCombinationsResponder())

	// PATCH /identity/conditionalAccess/authenticationStrength/policies/{policyId}/combinationConfigurations/{configId} - Update combination configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F-]+/combinationConfigurations/[0-9a-fA-F-]+$`,
		m.updateCombinationConfigurationResponder())

	// DELETE /identity/conditionalAccess/authenticationStrength/policies/{id} - Delete
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F-]+$`,
		m.deleteAuthenticationStrengthPolicyResponder())
}

// listAuthenticationStrengthPoliciesResponder handles GET requests to list authentication strength policies
func (m *AuthenticationStrengthMock) listAuthenticationStrengthPoliciesResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Return list of policies from mock state
		mockState.Lock()
		defer mockState.Unlock()

		policies := make([]map[string]any, 0, len(mockState.authenticationStrengths))
		for _, policy := range mockState.authenticationStrengths {
			// Create a copy to avoid modification
			policyCopy := make(map[string]any)
			for k, v := range policy {
				policyCopy[k] = v
			}
			policies = append(policies, policyCopy)
		}

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies",
			"value":          policies,
		}

		return factories.SuccessResponse(200, response)(req)
	}
}

// createAuthenticationStrengthPolicyResponder handles POST requests to create authentication strength policies
func (m *AuthenticationStrengthMock) createAuthenticationStrengthPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_authentication_strength_policy_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Generate a mock ID
		id := uuid.New().String()
		response["id"] = id

		// Add @odata.context
		response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies/$entity"

		// Update response with request data
		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		} else {
			response["description"] = nil
		}
		if allowedCombinations, ok := requestBody["allowedCombinations"]; ok {
			response["allowedCombinations"] = allowedCombinations
		}

		// Handle combinationConfigurations - generate IDs for each config
		if combinationConfigurations, ok := requestBody["combinationConfigurations"]; ok {
			if configsArray, ok := combinationConfigurations.([]any); ok {
				for i, configRaw := range configsArray {
					if configMap, ok := configRaw.(map[string]any); ok {
						// Generate ID for this configuration
						configMap["id"] = uuid.New().String()
						// Remove allowedIssuers as API doesn't return it
						delete(configMap, "allowedIssuers")
						configsArray[i] = configMap
					}
				}
				response["combinationConfigurations"] = configsArray
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.authenticationStrengths[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getAuthenticationStrengthPolicyResponder handles GET requests to retrieve authentication strength policies
func (m *AuthenticationStrengthMock) getAuthenticationStrengthPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_authentication_strength_policy_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Create response copy (deep copy to avoid concurrent modification issues)
		policyCopy := make(map[string]any)
		for k, v := range policy {
			policyCopy[k] = v
		}

		// Ensure @odata.context is present in GET response
		if _, exists := policyCopy["@odata.context"]; !exists {
			policyCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationStrength/policies/$entity"
		}

		// Mimic real API behavior: Remove allowedIssuers from all combinationConfigurations
		// The API accepts this field but never returns it in GET responses
		if configs, ok := policyCopy["combinationConfigurations"].([]any); ok {
			for i, configRaw := range configs {
				if configMap, ok := configRaw.(map[string]any); ok {
					delete(configMap, "allowedIssuers")
					configs[i] = configMap
				}
			}
			policyCopy["combinationConfigurations"] = configs
		}

		return factories.SuccessResponse(200, policyCopy)(req)
	}
}

// updateAuthenticationStrengthPolicyResponder handles PATCH requests to update authentication strength policies
func (m *AuthenticationStrengthMock) updateAuthenticationStrengthPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_authentication_strength_policy_not_found.json"))
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

		// Update fields from request
		mockState.Lock()
		for key, value := range requestBody {
			policy[key] = value
		}
		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.authenticationStrengths[id] = policy
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// updateAllowedCombinationsResponder handles POST requests to update allowed combinations
func (m *AuthenticationStrengthMock) updateAllowedCombinationsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL - it's before /updateAllowedCombinations
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-2]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update allowedCombinations
		if allowedCombinations, ok := requestBody["allowedCombinations"]; ok {
			policy["allowedCombinations"] = allowedCombinations
		}
		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.authenticationStrengths[id] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
	}
}

// updateCombinationConfigurationResponder handles PATCH requests to update combination configurations
func (m *AuthenticationStrengthMock) updateCombinationConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract policyId and configId from URL
		pathParts := strings.Split(req.URL.Path, "/")
		configId := pathParts[len(pathParts)-1]
		policyId := pathParts[len(pathParts)-3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Find and update the specific combination configuration
		if configs, ok := policy["combinationConfigurations"].([]any); ok {
			for i, configRaw := range configs {
				if configMap, ok := configRaw.(map[string]any); ok {
					if configMap["id"] == configId {
						// Update the configuration fields from request
						for key, value := range requestBody {
							configMap[key] = value
						}
						// Remove allowedIssuers as API doesn't return it
						delete(configMap, "allowedIssuers")
						configs[i] = configMap
						break
					}
				}
			}
			policy["combinationConfigurations"] = configs
		}

		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.authenticationStrengths[policyId] = policy
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// deleteAuthenticationStrengthPolicyResponder handles DELETE requests to remove authentication strength policies
func (m *AuthenticationStrengthMock) deleteAuthenticationStrengthPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		_, exists := mockState.authenticationStrengths[id]
		if exists {
			delete(mockState.authenticationStrengths, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_authentication_strength_policy_not_found.json"))
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

// CleanupMockState clears the mock state for clean test runs
func (m *AuthenticationStrengthMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Authentication Strength Policies
	for id := range mockState.authenticationStrengths {
		delete(mockState.authenticationStrengths, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *AuthenticationStrengthMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies",
		factories.ErrorResponse(400, "BadRequest", "Invalid authentication strength policy data"))

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Authentication strength policy is in use"))
}

// GetMockAuthenticationStrengthData returns sample authentication strength policy data for testing
func (m *AuthenticationStrengthMock) GetMockAuthenticationStrengthData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_authentication_strength_policy_maximal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse mock response: " + err.Error())
	}
	return response
}

// GetMockAuthenticationStrengthMinimalData returns minimal authentication strength policy data for testing
func (m *AuthenticationStrengthMock) GetMockAuthenticationStrengthMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_authentication_strength_policy_minimal.json"))
	if err != nil {
		panic("Failed to load mock response: " + err.Error())
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse mock response: " + err.Error())
	}
	return response
}
