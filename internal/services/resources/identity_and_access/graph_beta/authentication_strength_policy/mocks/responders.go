package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	authenticationStrengths map[string]map[string]any
}

func init() {
	mockState.authenticationStrengths = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("authentication_strength", &AuthenticationStrengthMock{})
}

type AuthenticationStrengthMock struct{}

var _ mocks.MockRegistrar = (*AuthenticationStrengthMock)(nil)

func (m *AuthenticationStrengthMock) RegisterMocks() {
	mockState.Lock()
	mockState.authenticationStrengths = make(map[string]map[string]any)
	mockState.Unlock()

	// Create authentication strength policy - POST /identity/conditionalAccess/authenticationStrength/policies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Generate a UUID for the new resource
		newId := uuid.New().String()

		// Load the template response
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_authentication_strength_policy_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response"}}`), nil
		}

		// Update response with request data
		responseObj["id"] = newId
		if displayName, ok := requestBody["displayName"]; ok {
			responseObj["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			responseObj["description"] = description
		}
		if allowedCombinations, ok := requestBody["allowedCombinations"]; ok {
			responseObj["allowedCombinations"] = allowedCombinations
		}
		if combinationConfigurations, ok := requestBody["combinationConfigurations"]; ok {
			// Generate IDs for each combination configuration (mimics API behavior)
			if configsArray, ok := combinationConfigurations.([]any); ok {
				for i, configRaw := range configsArray {
					if configMap, ok := configRaw.(map[string]any); ok {
						configMap["id"] = uuid.New().String()
						configsArray[i] = configMap
					}
				}
				responseObj["combinationConfigurations"] = configsArray
			} else {
				responseObj["combinationConfigurations"] = combinationConfigurations
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.authenticationStrengths[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get authentication strength policy - GET /identity/conditionalAccess/authenticationStrength/policies/{id}
	// Supports $expand query parameter for combinationConfigurations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`, func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		policyId := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Deep copy policy to avoid modifying the stored state
		policyBytes, _ := json.Marshal(policy)
		var policyResponse map[string]any
		json.Unmarshal(policyBytes, &policyResponse)

		// Mimic real API behavior: Remove allowedIssuers from all combinationConfigurations
		// The API accepts this field but never returns it in GET responses
		if configs, ok := policyResponse["combinationConfigurations"].([]any); ok {
			for i, configRaw := range configs {
				if configMap, ok := configRaw.(map[string]any); ok {
					delete(configMap, "allowedIssuers")
					configs[i] = configMap
				}
			}
			policyResponse["combinationConfigurations"] = configs
		}

		return httpmock.NewJsonResponse(200, policyResponse)
	})

	// Update authentication strength policy - PATCH /identity/conditionalAccess/authenticationStrength/policies/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			policy[key] = value
		}
		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.authenticationStrengths[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete authentication strength policy - DELETE /identity/conditionalAccess/authenticationStrength/policies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.authenticationStrengths[policyId]
		if exists {
			delete(mockState.authenticationStrengths, policyId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Update allowed combinations - POST /policies/authenticationStrengthPolicies/{id}/updateAllowedCombinations
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/policies/authenticationStrengthPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/updateAllowedCombinations$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-2]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.authenticationStrengths[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update allowedCombinations
		if allowedCombinations, ok := requestBody["allowedCombinations"]; ok {
			policy["allowedCombinations"] = allowedCombinations
		}
		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.authenticationStrengths[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
	})

	// Update combination configuration - PATCH /identity/conditionalAccess/authenticationStrength/policies/{id}/combinationConfigurations/{configId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/combinationConfigurations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-3]
		configId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
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

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AuthenticationStrengthMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationStrength/policies", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationStrength/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AuthenticationStrengthMock) CleanupMockState() {
	mockState.Lock()
	mockState.authenticationStrengths = make(map[string]map[string]any)
	mockState.Unlock()
}
