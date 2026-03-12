package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

const crossTenantAccessPolicyURL = "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy"

// mockState holds the in-memory state for the singleton crossTenantAccessPolicy.
// There is no ID — the policy is always addressed by the fixed URL.
var mockState struct {
	sync.Mutex
	policy map[string]any
}

func init() {
	mockState.policy = defaultPolicyState()

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	mocks.GlobalRegistry.Register("cross_tenant_access_policy", &CrossTenantAccessPolicyMock{})
}

// defaultPolicyState returns the service-default policy, matching what the API returns
// when the tenant has never modified their cross-tenant access policy.
func defaultPolicyState() map[string]any {
	return map[string]any{
		"@odata.context":        "https://graph.microsoft.com/beta/$metadata#policies/crossTenantAccessPolicy",
		"@odata.type":           "#microsoft.graph.crossTenantAccessPolicy",
		"displayName":           "CrossTenantAccessPolicy",
		"allowedCloudEndpoints": []any{},
	}
}

// CrossTenantAccessPolicyMock provides mock responses for crossTenantAccessPolicy operations.
type CrossTenantAccessPolicyMock struct{}

// Ensure CrossTenantAccessPolicyMock implements MockRegistrar interface.
var _ mocks.MockRegistrar = (*CrossTenantAccessPolicyMock)(nil)

// RegisterMocks registers HTTP mock responses for crossTenantAccessPolicy operations.
// This resource is a singleton — all requests target the fixed URL with no resource ID.
func (m *CrossTenantAccessPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.policy = defaultPolicyState()
	mockState.Unlock()

	// GET /beta/policies/crossTenantAccessPolicy - Read
	httpmock.RegisterResponder("GET", crossTenantAccessPolicyURL,
		m.getPolicyResponder())

	// PATCH /beta/policies/crossTenantAccessPolicy - Create (via PATCH) and Update
	httpmock.RegisterResponder("PATCH", crossTenantAccessPolicyURL,
		m.patchPolicyResponder())
}

// getPolicyResponder returns the current singleton policy state.
func (m *CrossTenantAccessPolicyMock) getPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		policyCopy := make(map[string]any, len(mockState.policy))
		for k, v := range mockState.policy {
			policyCopy[k] = v
		}
		mockState.Unlock()

		return factories.SuccessResponse(200, policyCopy)(req)
	}
}

// patchPolicyResponder applies PATCH request fields to the singleton policy state and returns 204.
func (m *CrossTenantAccessPolicyMock) patchPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		for key, value := range requestBody {
			mockState.policy[key] = value
		}
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions.
func (m *CrossTenantAccessPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", crossTenantAccessPolicyURL,
		factories.ErrorResponse(403, "Forbidden", "Access denied to cross-tenant access policy"))

	httpmock.RegisterResponder("PATCH", crossTenantAccessPolicyURL,
		factories.ErrorResponse(403, "Forbidden", "Insufficient privileges to update cross-tenant access policy"))
}

// CleanupMockState resets the singleton policy state to service defaults.
func (m *CrossTenantAccessPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policy = defaultPolicyState()
}

