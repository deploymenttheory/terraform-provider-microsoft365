package mocks

import (
	"encoding/json"
	"maps"
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
	b2bManagementPolicies map[string]map[string]any
}

func init() {
	mockState.b2bManagementPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("b2b_management_policy", &B2bManagementPolicyMock{})
}

// B2bManagementPolicyMock provides mock responses for B2B Management Policy operations
type B2bManagementPolicyMock struct{}

var _ mocks.MockRegistrar = (*B2bManagementPolicyMock)(nil)

// descriptionNotFoundResponse reproduces the live Microsoft Graph response to any write that
// carries the description property.
func descriptionNotFoundResponse() *http.Response {
	return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource '' does not exist or one of its queried reference-property objects are not present."}}`)
}

// RegisterMocks registers HTTP mock responses for B2B Management Policy operations
func (m *B2bManagementPolicyMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/b2bManagementPolicies",
		m.createB2bManagementPolicyResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		m.getB2bManagementPolicyResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		m.updateB2bManagementPolicyResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		m.deleteB2bManagementPolicyResponder())
}

func (m *B2bManagementPolicyMock) createB2bManagementPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		if _, ok := requestBody["description"]; ok {
			return descriptionNotFoundResponse(), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_b2b_management_policy_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		id := uuid.New().String()
		response["id"] = id

		for _, key := range []string{"displayName", "definition", "isOrganizationDefault"} {
			if value, ok := requestBody[key]; ok {
				response[key] = value
			}
		}

		mockState.Lock()
		mockState.b2bManagementPolicies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *B2bManagementPolicyMock) getB2bManagementPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		policy, exists := mockState.b2bManagementPolicies[id]
		var policyCopy map[string]any
		if exists {
			policyCopy = maps.Clone(policy)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Directory_ObjectNotFound","message":"Unable to read the company information from the directory."}}`), nil
		}

		policyCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#policies/b2bManagementPolicies/$entity"

		return factories.SuccessResponse(200, policyCopy)(req)
	}
}

func (m *B2bManagementPolicyMock) updateB2bManagementPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		if _, ok := requestBody["description"]; ok {
			return descriptionNotFoundResponse(), nil
		}

		mockState.Lock()
		policy, exists := mockState.b2bManagementPolicies[id]
		if exists {
			for _, key := range []string{"displayName", "definition", "isOrganizationDefault"} {
				if value, ok := requestBody[key]; ok {
					policy[key] = value
				}
			}
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *B2bManagementPolicyMock) deleteB2bManagementPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		mockState.Lock()
		_, exists := mockState.b2bManagementPolicies[id]
		delete(mockState.b2bManagementPolicies, id)
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterEventualConsistencyMocks overrides the GET responder so that the first staleReadCount
// reads after each write return the previous representation of the policy, reproducing the
// replica flapping observed live after POST and PATCH. Call after RegisterMocks.
func (m *B2bManagementPolicyMock) RegisterEventualConsistencyMocks(staleReadCount int) {
	var mu sync.Mutex
	previous := map[string]map[string]any{}
	staleRemaining := map[string]int{}

	get := m.getB2bManagementPolicyResponder()
	patch := m.updateB2bManagementPolicyResponder()

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			pathParts := strings.Split(req.URL.Path, "/")
			id := pathParts[len(pathParts)-1]

			mockState.Lock()
			snapshot := maps.Clone(mockState.b2bManagementPolicies[id])
			mockState.Unlock()

			resp, err := patch(req)
			if err == nil && resp.StatusCode == http.StatusNoContent {
				mu.Lock()
				previous[id] = snapshot
				staleRemaining[id] = staleReadCount
				mu.Unlock()
			}
			return resp, err
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			pathParts := strings.Split(req.URL.Path, "/")
			id := pathParts[len(pathParts)-1]

			mu.Lock()
			if staleRemaining[id] > 0 {
				staleRemaining[id]--
				stale := maps.Clone(previous[id])
				mu.Unlock()
				stale["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#policies/b2bManagementPolicies/$entity"
				return factories.SuccessResponse(200, stale)(req)
			}
			mu.Unlock()
			return get(req)
		})
}

// RegisterStaleNotFoundMocks overrides the GET responder so that the next notFoundCount reads
// return the 404 Directory_ObjectNotFound served by a stale Microsoft Entra replica, then delegate
// to the normal responder. Call after RegisterMocks, e.g. from a test step's PreConfig.
func (m *B2bManagementPolicyMock) RegisterStaleNotFoundMocks(notFoundCount int) {
	var mu sync.Mutex
	remaining := notFoundCount
	get := m.getB2bManagementPolicyResponder()

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			if remaining > 0 {
				remaining--
				mu.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"Directory_ObjectNotFound","message":"Unable to read the company information from the directory."}}`), nil
			}
			mu.Unlock()
			return get(req)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *B2bManagementPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	clear(mockState.b2bManagementPolicies)
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *B2bManagementPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/b2bManagementPolicies",
		factories.ErrorResponse(400, "Request_BadRequest", "Property definition has an invalid value."))
}
