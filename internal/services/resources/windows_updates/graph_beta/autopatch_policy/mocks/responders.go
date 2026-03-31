package mocks

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policies map[string]map[string]any
}

func init() {
	mockState.policies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_autopatch_policy", &WindowsAutopatchPolicyMock{})
}

type WindowsAutopatchPolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopatchPolicyMock)(nil)

// getJSONFileForDisplayName determines which JSON file to load based on the policy's display name
func getJSONFileForDisplayName(displayName string) string {
	displayNameUpper := strings.ToUpper(displayName)

	switch {
	case strings.HasPrefix(displayNameUpper, "WAP001"):
		return "post_windows_autopatch_policy_wap001_success.json"
	case strings.HasPrefix(displayNameUpper, "WAP002"):
		return "post_windows_autopatch_policy_wap002_success.json"
	case strings.HasPrefix(displayNameUpper, "WAP003"):
		return "post_windows_autopatch_policy_wap003_success.json"
	default:
		return ""
	}
}

func (m *WindowsAutopatchPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreatePolicyResponder()
	m.registerGetPolicyResponder()
	m.registerUpdatePolicyResponder()
	m.registerDeletePolicyResponder()
}

func (m *WindowsAutopatchPolicyMock) registerCreatePolicyResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/policies", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
		}

		jsonFileName := getJSONFileForDisplayName(displayName)
		if jsonFileName == "" {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Unable to determine JSON file for displayName: %s"}}`, displayName)), nil
		}

		responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		newId := uuid.New().String()
		responseObj["id"] = newId

		// Preserve the description from the request body if provided
		if desc, ok := requestBody["description"]; ok {
			responseObj["description"] = desc
		}
		// Preserve approvalRules from request if provided
		if rules, ok := requestBody["approvalRules"]; ok {
			responseObj["approvalRules"] = rules
		}

		mockState.Lock()
		mockState.policies[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})
}

func (m *WindowsAutopatchPolicyMock) registerGetPolicyResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, policy)
	})
}

func (m *WindowsAutopatchPolicyMock) registerUpdatePolicyResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		maps.Copy(policy, requestBody)
		policy["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.policies[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsAutopatchPolicyMock) registerDeletePolicyResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.policies[policyId]; exists {
			delete(mockState.policies, policyId)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
	})
}

func (m *WindowsAutopatchPolicyMock) RegisterErrorMocks() {
	m.registerCreatePolicyErrorResponder()
	m.registerGetPolicyErrorResponder()
	m.registerUpdatePolicyErrorResponder()
	m.registerDeletePolicyErrorResponder()
}

func (m *WindowsAutopatchPolicyMock) registerCreatePolicyErrorResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/policies",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid policy configuration"}}`))
}

func (m *WindowsAutopatchPolicyMock) registerGetPolicyErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

func (m *WindowsAutopatchPolicyMock) registerUpdatePolicyErrorResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *WindowsAutopatchPolicyMock) registerDeletePolicyErrorResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *WindowsAutopatchPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()
}
