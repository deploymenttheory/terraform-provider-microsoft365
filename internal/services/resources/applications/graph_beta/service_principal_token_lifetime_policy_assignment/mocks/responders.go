package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	// assignments maps servicePrincipalID to tokenLifetimePolicyID
	assignments map[string]string
}

func init() {
	mockState.assignments = make(map[string]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("service_principal_token_lifetime_policy_assignment", &ServicePrincipalTokenLifetimePolicyAssignmentMock{})
}

// ServicePrincipalTokenLifetimePolicyAssignmentMock provides mock responses for SP token lifetime policy assignment operations
type ServicePrincipalTokenLifetimePolicyAssignmentMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalTokenLifetimePolicyAssignmentMock)(nil)

// RegisterMocks registers HTTP mock responses for SP token lifetime policy assignment operations
func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) RegisterMocks() {
	// POST /servicePrincipals/{id}/tokenLifetimePolicies/$ref - Assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/tokenLifetimePolicies/\$ref$`,
		m.assignTokenLifetimePolicyResponder())

	// GET /servicePrincipals/{id}/tokenLifetimePolicies - List assigned policies
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/tokenLifetimePolicies$`,
		m.listTokenLifetimePoliciesResponder())

	// DELETE /servicePrincipals/{id}/tokenLifetimePolicies/{policyId}/$ref - Remove
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/tokenLifetimePolicies/[0-9a-fA-F-]+/\$ref$`,
		m.removeTokenLifetimePolicyResponder())
}

func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) assignTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		// Path: /beta/servicePrincipals/{spID}/tokenLifetimePolicies/$ref
		var spID string
		for i, part := range pathParts {
			if part == "servicePrincipals" && i+1 < len(pathParts) {
				spID = pathParts[i+1]
				break
			}
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Extract policy ID from @odata.id
		odataID, ok := requestBody["@odata.id"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Missing @odata.id"}}`), nil
		}

		odataParts := strings.Split(odataID, "/")
		policyID := odataParts[len(odataParts)-1]

		mockState.Lock()
		mockState.assignments[spID] = policyID
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) listTokenLifetimePoliciesResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		var spID string
		for i, part := range pathParts {
			if part == "servicePrincipals" && i+1 < len(pathParts) {
				spID = pathParts[i+1]
				break
			}
		}

		mockState.Lock()
		policyID, exists := mockState.assignments[spID]
		mockState.Unlock()

		if !exists {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/tokenLifetimePolicies",
				"value":          []any{},
			}
			return factories.SuccessResponse(200, response)(req)
		}

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/tokenLifetimePolicies",
			"value": []any{
				map[string]any{
					"id":                    policyID,
					"@odata.type":           "#microsoft.graph.tokenLifetimePolicy",
					"displayName":           "test-token-lifetime-policy",
					"isOrganizationDefault": false,
					"definition":            []string{"{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"},
				},
			},
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) removeTokenLifetimePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{spID}/tokenLifetimePolicies/{policyID}/$ref
		pathParts := strings.Split(req.URL.Path, "/")
		var spID, policyID string
		for i, part := range pathParts {
			if part == "servicePrincipals" && i+1 < len(pathParts) {
				spID = pathParts[i+1]
			}
			if part == "tokenLifetimePolicies" && i+1 < len(pathParts) {
				policyID = pathParts[i+1]
			}
		}

		mockState.Lock()
		existingPolicyID, exists := mockState.assignments[spID]
		if exists && existingPolicyID == policyID {
			delete(mockState.assignments, spID)
		} else {
			exists = false
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// CleanupMockState clears the mock state for clean test runs
func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.assignments {
		delete(mockState.assignments, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *ServicePrincipalTokenLifetimePolicyAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/error-id/tokenLifetimePolicies/\$ref$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
}
