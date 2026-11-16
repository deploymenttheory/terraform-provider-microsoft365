package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const policyURL = "https://graph.microsoft.com/beta/policies/identitySecurityDefaultsEnforcementPolicy"

// conditionalAccessConflictMessage mirrors the 400 Graph returns when security defaults are enabled while
// Conditional Access policies are enabled in the tenant.
const conditionalAccessConflictMessage = "Security Defaults cannot be enabled when Conditional Access policies are enabled."

var mockState struct {
	sync.Mutex
	policy                     map[string]any
	patches                    []map[string]any
	seedErr                    error
	conditionalAccessIsEnabled bool
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("identity_and_access_security_defaults_enforcement_policy", &SecurityDefaultsEnforcementPolicyMock{})
}

// SecurityDefaultsEnforcementPolicyMock provides stateful httpmock responders for the singleton security defaults enforcement policy.
type SecurityDefaultsEnforcementPolicyMock struct{}

var _ mocks.MockRegistrar = (*SecurityDefaultsEnforcementPolicyMock)(nil)

// loadSeedPolicy loads the policy that GET returns before any PATCH.
func loadSeedPolicy() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_get", "get_security_defaults_enforcement_policy_success.json")
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var policy map[string]any
	if err := json.Unmarshal(jsonData, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return policy, nil
}

func errorResponse(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("build error response: %w", err)
	}
	return resp, nil
}

func (m *SecurityDefaultsEnforcementPolicyMock) RegisterMocks() {
	policy, err := loadSeedPolicy()
	mockState.Lock()
	mockState.policy = policy
	mockState.patches = nil
	mockState.seedErr = err
	mockState.conditionalAccessIsEnabled = false
	mockState.Unlock()

	// Get security defaults enforcement policy - GET /policies/identitySecurityDefaultsEnforcementPolicy
	httpmock.RegisterResponder("GET", policyURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.seedErr != nil {
			return errorResponse(500, "InternalServerError", mockState.seedErr.Error())
		}
		return httpmock.NewJsonResponse(200, mockState.policy)
	})

	// Update security defaults enforcement policy - PATCH /policies/identitySecurityDefaultsEnforcementPolicy
	httpmock.RegisterResponder("PATCH", policyURL, func(req *http.Request) (*http.Response, error) {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		var requestBody map[string]any
		if err := json.Unmarshal(content, &requestBody); err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}

		isEnabled, hasIsEnabled := requestBody["isEnabled"]
		if hasIsEnabled {
			if _, isBool := isEnabled.(bool); !isBool {
				return errorResponse(400, "BadRequest", "Invalid value for property 'isEnabled'")
			}
		}

		mockState.Lock()
		defer mockState.Unlock()

		if isEnabled == true && mockState.conditionalAccessIsEnabled {
			return errorResponse(400, "BadRequest", conditionalAccessConflictMessage)
		}

		var recorded map[string]any
		_ = json.Unmarshal(content, &recorded)
		mockState.patches = append(mockState.patches, recorded)

		if hasIsEnabled {
			mockState.policy["isEnabled"] = isEnabled
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterConditionalAccessConflictMocks registers the stateful responders for a tenant that has enabled
// Conditional Access policies, where Graph rejects enabling security defaults with 400.
func (m *SecurityDefaultsEnforcementPolicyMock) RegisterConditionalAccessConflictMocks() {
	m.RegisterMocks()
	mockState.Lock()
	mockState.conditionalAccessIsEnabled = true
	mockState.Unlock()
}

func (m *SecurityDefaultsEnforcementPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterResponder("GET", policyURL, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterResponder("PATCH", policyURL, httpmock.NewStringResponder(403, forbidden))
}

// RegisterNotFoundMocks returns 404 for every request to the policy.
func (m *SecurityDefaultsEnforcementPolicyMock) RegisterNotFoundMocks() {
	notFound := `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`
	httpmock.RegisterResponder("GET", policyURL, httpmock.NewStringResponder(404, notFound))
	httpmock.RegisterResponder("PATCH", policyURL, httpmock.NewStringResponder(404, notFound))
}

func (m *SecurityDefaultsEnforcementPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policy = nil
	mockState.patches = nil
	mockState.seedErr = nil
	mockState.conditionalAccessIsEnabled = false
	mockState.Unlock()
}

// PatchRequests returns the request bodies received by the PATCH responder, oldest first.
func (m *SecurityDefaultsEnforcementPolicyMock) PatchRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.patches...)
}

// Policy returns a copy of the policy the mock currently stores.
func (m *SecurityDefaultsEnforcementPolicyMock) Policy() map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	policy := map[string]any{}
	for key, value := range mockState.policy {
		policy[key] = value
	}
	return policy
}
