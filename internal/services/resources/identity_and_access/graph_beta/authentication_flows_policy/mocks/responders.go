package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const URL = "https://graph.microsoft.com/beta/policies/authenticationFlowsPolicy"

// AuthenticationFlowsPolicyMock retains the singleton across create, update, and destroy.
type AuthenticationFlowsPolicyMock struct {
	mu      sync.Mutex
	policy  map[string]any
	patches int
}

var _ mocks.MockRegistrar = (*AuthenticationFlowsPolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("authentication_flows_policy", &AuthenticationFlowsPolicyMock{})
}

// RegisterMocks registers the direct singleton GET and PATCH endpoints.
func (m *AuthenticationFlowsPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	data, err := os.ReadFile(
		"tests/responses/validate_get/get_authentication_flows_policy_success.json",
	)
	if err != nil {
		panic(fmt.Sprintf("load authentication flows policy fixture: %v", err))
	}
	var policy map[string]any
	if err := json.Unmarshal(data, &policy); err != nil {
		panic(fmt.Sprintf("decode authentication flows policy fixture: %v", err))
	}
	m.mu.Lock()
	m.policy = policy
	m.mu.Unlock()
	httpmock.RegisterResponder("GET", URL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		response, err := httpmock.NewJsonResponse(200, m.policy)
		if err != nil {
			return nil, fmt.Errorf("encode authentication flows policy: %w", err)
		}
		return response, nil
	})
	httpmock.RegisterResponder("PATCH", URL, func(req *http.Request) (*http.Response, error) {
		var body struct {
			SelfServiceSignUp struct {
				IsEnabled *bool `json:"isEnabled"`
			} `json:"selfServiceSignUp"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode authentication flows policy: %w", err)
		}
		if body.SelfServiceSignUp.IsEnabled == nil {
			return httpmock.NewStringResponse(
				400,
				`{"error":{"code":"BadRequest","message":"isEnabled must be a Boolean"}}`,
			), nil
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		m.policy["selfServiceSignUp"] = map[string]any{
			"isEnabled": *body.SelfServiceSignUp.IsEnabled,
		}
		m.patches++
		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterErrorMocks reproduces the missing application permission observed with curl.
func (m *AuthenticationFlowsPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"PATCH",
		URL,
		httpmock.NewStringResponder(
			403,
			`{"error":{"code":"AADB2C","message":"The application does not have any of the required application permissions (Policy.ReadWrite.AuthenticationFlows) to access the resource."}}`,
		),
	)
}

// CleanupMockState resets the singleton between tests.
func (m *AuthenticationFlowsPolicyMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policy = nil
	m.patches = 0
}

func (m *AuthenticationFlowsPolicyMock) Snapshot() (bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.policy["selfServiceSignUp"].(map[string]any)["isEnabled"].(bool), m.patches
}
