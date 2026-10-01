package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

const URL = "https://graph.microsoft.com/beta/policies/authorizationPolicy"
const PatchURL = URL + "/authorizationPolicy"

// AuthorizationPolicyMock retains the singleton across Terraform create, update and destroy.
type AuthorizationPolicyMock struct {
	mu      sync.Mutex
	policy  map[string]any
	patches int
}

var _ mocks.MockRegistrar = (*AuthorizationPolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("authorization_policy", &AuthorizationPolicyMock{})
}

// RegisterMocks registers the GET and PATCH endpoints using the captured live response.
func (m *AuthorizationPolicyMock) RegisterMocks() {
	m.CleanupMockState()
	data, err := os.ReadFile(filepath.Join("tests", "responses", "validate_get", "get_authorization_policy_success.json"))
	if err != nil {
		panic(fmt.Sprintf("load authorization policy fixture: %v", err))
	}
	var response struct {
		Value []map[string]any `json:"value"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		panic(fmt.Sprintf("decode authorization policy fixture: %v", err))
	}
	if len(response.Value) != 1 {
		panic("authorization policy fixture must contain one singleton")
	}
	m.mu.Lock()
	m.policy = response.Value[0]
	m.mu.Unlock()

	// Read the singleton policy collection.
	httpmock.RegisterResponder("GET", URL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		response, err := httpmock.NewJsonResponse(200, map[string]any{"value": []any{m.policy}})
		if err != nil {
			return nil, fmt.Errorf("encode authorization policy: %w", err)
		}
		return response, nil
	})

	// Create and update both PATCH the existing singleton.
	httpmock.RegisterResponder("PATCH", PatchURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode authorization policy: %w", err)
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		for key, value := range body {
			if key == "@odata.type" || value == nil {
				continue
			}
			if key == "defaultUserRolePermissions" {
				for name, permission := range value.(map[string]any) {
					m.policy[key].(map[string]any)[name] = permission
				}
			} else if key == "permissionGrantPolicyIdsAssignedToDefaultUserRole" {
				policies := value.([]any)
				for i, policy := range policies {
					policies[i] = strings.Replace(policy.(string), "managePermissionGrants", "ManagePermissionGrants", 1)
				}
				m.policy[key] = policies
			} else {
				m.policy[key] = value
			}
		}
		m.patches++
		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterErrorMocks simulates missing authorization policy permissions.
func (m *AuthorizationPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", URL, httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges"}}`))

	// Create and update both PATCH the existing singleton.
	httpmock.RegisterResponder("PATCH", PatchURL, httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges"}}`))
}

// CleanupMockState releases the singleton mock state between tests.
func (m *AuthorizationPolicyMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policy = nil
	m.patches = 0
}

func (m *AuthorizationPolicyMock) Snapshot() (string, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.policy["allowInvitesFrom"].(string), m.patches
}
func (m *AuthorizationPolicyMock) SetInvites(value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policy["allowInvitesFrom"] = value
}
