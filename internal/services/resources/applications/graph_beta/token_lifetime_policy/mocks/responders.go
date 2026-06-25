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
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	tokenLifetimePolicies map[string]map[string]any
}

func init() {
	mockState.tokenLifetimePolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("token_lifetime_policy", &TokenLifetimePolicyMock{})
}

// TokenLifetimePolicyMock provides mock responses for Token Lifetime Policy operations
type TokenLifetimePolicyMock struct{}

var _ mocks.MockRegistrar = (*TokenLifetimePolicyMock)(nil)

// RegisterMocks registers HTTP mock responses for Token Lifetime Policy operations
func (m *TokenLifetimePolicyMock) RegisterMocks() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/policies/tokenLifetimePolicies",
		m.listTokenLifetimePoliciesResponder())
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/tokenLifetimePolicies",
		m.createTokenLifetimePolicyResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/[0-9a-fA-F-]+$`,
		m.getTokenLifetimePolicyResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/[0-9a-fA-F-]+$`,
		m.updateTokenLifetimePolicyResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/[0-9a-fA-F-]+$`,
		m.deleteTokenLifetimePolicyResponder())
}

func (m *TokenLifetimePolicyMock) listTokenLifetimePoliciesResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		policies := make([]map[string]any, 0, len(mockState.tokenLifetimePolicies))
		for _, policy := range mockState.tokenLifetimePolicies {
			policyCopy := make(map[string]any)
			for k, v := range policy {
				policyCopy[k] = v
			}
			policies = append(policies, policyCopy)
		}

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/tokenLifetimePolicies",
			"value":          policies,
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *TokenLifetimePolicyMock) createTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_token_lifetime_policy_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		id := uuid.New().String()
		response["id"] = id
		response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#policies/tokenLifetimePolicies/$entity"

		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		} else {
			response["description"] = nil
		}
		if definition, ok := requestBody["definition"]; ok {
			response["definition"] = definition
		}
		if isOrgDefault, ok := requestBody["isOrganizationDefault"]; ok {
			response["isOrganizationDefault"] = isOrgDefault
		}

		mockState.Lock()
		mockState.tokenLifetimePolicies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *TokenLifetimePolicyMock) getTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.tokenLifetimePolicies[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource 'tokenLifetimePolicy' with id '" + id + "' not found."}}`), nil
		}

		policyCopy := make(map[string]any)
		for k, v := range policy {
			policyCopy[k] = v
		}

		if _, exists := policyCopy["@odata.context"]; !exists {
			policyCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#policies/tokenLifetimePolicies/$entity"
		}

		return factories.SuccessResponse(200, policyCopy)(req)
	}
}

func (m *TokenLifetimePolicyMock) updateTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.tokenLifetimePolicies[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		for key, value := range requestBody {
			policy[key] = value
		}
		mockState.tokenLifetimePolicies[id] = policy
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *TokenLifetimePolicyMock) deleteTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		_, exists := mockState.tokenLifetimePolicies[id]
		if exists {
			delete(mockState.tokenLifetimePolicies, id)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// CleanupMockState clears the mock state for clean test runs
func (m *TokenLifetimePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.tokenLifetimePolicies {
		delete(mockState.tokenLifetimePolicies, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *TokenLifetimePolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/tokenLifetimePolicies",
		factories.ErrorResponse(400, "BadRequest", "Invalid token lifetime policy data"))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/tokenLifetimePolicies/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Token lifetime policy is in use"))
}
