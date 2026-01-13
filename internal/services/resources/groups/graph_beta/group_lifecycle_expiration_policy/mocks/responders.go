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

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	policies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.policies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("group_lifecycle_policy", &GroupLifecyclePolicyMock{})
}

// GroupLifecyclePolicyMock provides mock responses for group lifecycle policy operations
type GroupLifecyclePolicyMock struct{}

// Ensure GroupLifecyclePolicyMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupLifecyclePolicyMock)(nil)

// RegisterMocks registers HTTP mock responses for group lifecycle policy operations
func (m *GroupLifecyclePolicyMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			policies := make([]map[string]any, 0, len(mockState.policies))
			for _, policy := range mockState.policies {
				policies = append(policies, policy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groupLifecyclePolicies",
				"value":          policies,
			}

			respBody, _ := json.Marshal(response)
			return httpmock.NewStringResponse(200, string(respBody)), nil
		})

	// Register GET for specific policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			policyID := httpmock.MustGetSubmatch(req, 1)

			// Handle special test IDs with external JSON files
			if strings.Contains(policyID, "error") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if policy, exists := mockState.policies[policyID]; exists {
				respBody, _ := json.Marshal(policy)
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})

	// Register POST for creating policies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a new ID for the policy
			newID := uuid.New().String()
			requestBody["id"] = newID

			mockState.Lock()
			mockState.policies[newID] = requestBody
			mockState.Unlock()

			respBody, _ := json.Marshal(requestBody)
			resp := httpmock.NewStringResponse(201, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register PATCH for updating policies
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			policyID := httpmock.MustGetSubmatch(req, 1)

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if existingPolicy, exists := mockState.policies[policyID]; exists {
				// Merge the updates into the existing policy
				for key, value := range requestBody {
					existingPolicy[key] = value
				}
				mockState.policies[policyID] = existingPolicy

				respBody, _ := json.Marshal(existingPolicy)
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})

	// Register DELETE for deleting policies
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			policyID := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			defer mockState.Unlock()

			if _, exists := mockState.policies[policyID]; exists {
				delete(mockState.policies, policyID)
				return httpmock.NewStringResponse(204, ""), nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})
}

// RegisterErrorMocks registers error mock responses
func (m *GroupLifecyclePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()

	// Load error JSON from external file
	errorJSON, err := helpers.ParseJSONFile("tests/responses/validate_error/error_bad_request.json")
	if err != nil {
		// Fallback to inline JSON if file not found
		errorJSON = `{"error":{"code":"BadRequest","message":"Invalid request"}}`
	}

	// Register error responses for POST
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		httpmock.NewStringResponder(400, errorJSON))

	// Register error responses for GET
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/`,
		httpmock.NewStringResponder(400, errorJSON))
}

// CleanupMockState resets the mock state
func (m *GroupLifecyclePolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()
}
