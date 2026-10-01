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

const URL = "https://graph.microsoft.com/beta/policies/externalIdentitiesPolicy"

// ExternalIdentitiesPolicyMock retains the singleton across create, update, and destroy.
type ExternalIdentitiesPolicyMock struct {
	mu      sync.Mutex
	policy  map[string]any
	patches int
}

var _ mocks.MockRegistrar = (*ExternalIdentitiesPolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("external_identities_policy", &ExternalIdentitiesPolicyMock{})
}

// RegisterMocks registers the direct singleton GET and PATCH endpoints.
func (m *ExternalIdentitiesPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	data, err := os.ReadFile(
		"tests/responses/validate_get/get_external_identities_policy_success.json",
	)
	if err != nil {
		panic(fmt.Sprintf("load external identities policy fixture: %v", err))
	}
	var policy map[string]any
	if err := json.Unmarshal(data, &policy); err != nil {
		panic(fmt.Sprintf("decode external identities policy fixture: %v", err))
	}
	m.mu.Lock()
	m.policy = policy
	m.mu.Unlock()
	httpmock.RegisterResponder("GET", URL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		response, err := httpmock.NewJsonResponse(200, m.policy)
		if err != nil {
			return nil, fmt.Errorf("encode external identities policy: %w", err)
		}
		return response, nil
	})
	httpmock.RegisterResponder("PATCH", URL, func(req *http.Request) (*http.Response, error) {
		var body struct {
			AllowExternalIdentitiesToLeave    *bool `json:"allowExternalIdentitiesToLeave"`
			AllowDeletedIdentitiesDataRemoval *bool `json:"allowDeletedIdentitiesDataRemoval"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode external identities policy: %w", err)
		}
		if body.AllowExternalIdentitiesToLeave == nil ||
			body.AllowDeletedIdentitiesDataRemoval == nil {
			return httpmock.NewStringResponse(
				400,
				`{"error":{"code":"BadRequest","message":"both settings must be Booleans"}}`,
			), nil
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		m.policy["allowExternalIdentitiesToLeave"] = *body.AllowExternalIdentitiesToLeave
		m.policy["allowDeletedIdentitiesDataRemoval"] = *body.AllowDeletedIdentitiesDataRemoval
		m.patches++
		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterErrorMocks reproduces the missing application permission observed with curl.
func (m *ExternalIdentitiesPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"PATCH",
		URL,
		httpmock.NewStringResponder(
			403,
			`{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`,
		),
	)
}

// CleanupMockState resets the singleton between tests.
func (m *ExternalIdentitiesPolicyMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policy = nil
	m.patches = 0
}

func (m *ExternalIdentitiesPolicyMock) Snapshot() (bool, bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.policy["allowExternalIdentitiesToLeave"].(bool), m.policy["allowDeletedIdentitiesDataRemoval"].(bool), m.patches
}
