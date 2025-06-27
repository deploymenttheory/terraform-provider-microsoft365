package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	policies map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.policies = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// GroupLifecyclePolicyMock provides mock responses for group lifecycle policy operations
type GroupLifecyclePolicyMock struct{}

// RegisterMocks registers HTTP mock responses for group lifecycle policy operations
func (m *GroupLifecyclePolicyMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.policies = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register specific test policies
	registerTestPolicies()

	// Register GET for policy by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groupLifecyclePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.policies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group lifecycle policy not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register GET for listing policies
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groupLifecyclePolicies(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			policies := make([]map[string]interface{}, 0, len(mockState.policies))
			for _, policy := range mockState.policies {
				policies = append(policies, policy)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groupLifecyclePolicies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating policies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			var policyData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&policyData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := policyData["groupLifetimeInDays"].(float64); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"groupLifetimeInDays is required"}}`), nil
			}
			if _, ok := policyData["managedGroupTypes"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"managedGroupTypes is required"}}`), nil
			}

			// Generate ID if not provided
			if policyData["id"] == nil {
				policyData["id"] = uuid.New().String()
			}

			// Store policy in mock state
			policyId := policyData["id"].(string)
			mockState.Lock()
			mockState.policies[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, policyData)
		})

	// Register PATCH for updating policies
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/groupLifecyclePolicies/[^/]+/?$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.policies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group lifecycle policy not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Only update fields present in the PATCH body, and remove optional fields if omitted
			mockState.Lock()
			for key, value := range updateData {
				policyData[key] = value
			}
			// Remove optional fields if not present in PATCH body
			optionalFields := []string{"alternateNotificationEmails"}
			for _, field := range optionalFields {
				if _, present := updateData[field]; !present {
					delete(policyData, field)
				}
			}
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register DELETE for deleting policies
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/groupLifecyclePolicies/[^/]+/?$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.policies[policyId]
			if exists {
				delete(mockState.policies, policyId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group lifecycle policy not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers error responses for testing error scenarios
func (m *GroupLifecyclePolicyMock) RegisterErrorMocks() {
	// Register error responses for POST
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`), nil
		})

	// Register error responses for GET
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groupLifecyclePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`), nil
		})
}

// registerTestPolicies registers predefined test policies
func registerTestPolicies() {
	// Register a test policy for import testing
	testPolicy := map[string]interface{}{
		"id":                          "test-policy-id",
		"groupLifetimeInDays":         float64(180),
		"managedGroupTypes":           "All",
		"alternateNotificationEmails": "admin@example.com",
	}

	mockState.Lock()
	mockState.policies["test-policy-id"] = testPolicy
	mockState.Unlock()
}
