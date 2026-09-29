package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// MacosDeviceCompliancePolicyMock simulates the compliance policy lifecycle.
// Graph does not echo ruleName, and navigation properties require explicit expansion.
type MacosDeviceCompliancePolicyMock struct {
	sync.Mutex
	policies map[string]map[string]any
}

var _ mocks.MockRegistrar = (*MacosDeviceCompliancePolicyMock)(nil)

func init() {
	mocks.GlobalRegistry.Register(
		"macos_device_compliance_policy",
		&MacosDeviceCompliancePolicyMock{},
	)
}

// RegisterMocks registers stateful responders for macOS compliance policies.
func (m *MacosDeviceCompliancePolicyMock) RegisterMocks() {
	m.policies = make(map[string]map[string]any)
	const base = "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies"
	const item = `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+`
	httpmock.RegisterResponder("POST", base, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode policy: %w", err)
		}
		fixture, err := helpers.ParseJSONFile(
			"../tests/responses/validate_get/get_macos_device_compliance_policy.json",
		)
		if err != nil {
			return nil, fmt.Errorf("load policy fixture: %w", err)
		}
		var policy map[string]any
		if err := json.Unmarshal([]byte(fixture), &policy); err != nil {
			return nil, fmt.Errorf("decode policy fixture: %w", err)
		}
		for key, value := range body {
			policy[key] = value
		}
		policy["id"] = uuid.NewString()
		m.Lock()
		defer m.Unlock()
		m.policies[policy["id"].(string)] = policy
		return httpmock.NewJsonResponse(201, policy)
	})
	httpmock.RegisterResponder("GET", item+`$`, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		policy, ok := m.policies[policyID(req)]
		if !ok {
			return httpmock.NewStringResponse(
				404,
				`{"error":{"code":"ResourceNotFound","message":"Policy not found"}}`,
			), nil
		}
		// Clone the stored request so that reads cannot mutate configured values.
		raw, err := json.Marshal(policy)
		if err != nil {
			return nil, fmt.Errorf("encode policy: %w", err)
		}
		var response map[string]any
		if err := json.Unmarshal(raw, &response); err != nil {
			return nil, fmt.Errorf("clone policy: %w", err)
		}
		expand := req.URL.Query().Get("$expand")
		if !strings.Contains(expand, "assignments") {
			delete(response, "assignments")
		}
		if !strings.Contains(expand, "scheduledActionsForRule") {
			delete(response, "scheduledActionsForRule")
		} else {
			for _, value := range response["scheduledActionsForRule"].([]any) {
				rule := value.(map[string]any)
				rule["ruleName"] = nil
				if !strings.Contains(expand, "$expand=scheduledActionConfigurations") {
					delete(rule, "scheduledActionConfigurations")
					continue
				}
				for _, value := range rule["scheduledActionConfigurations"].([]any) {
					action := value.(map[string]any)
					if _, ok := action["gracePeriodHours"]; !ok {
						action["gracePeriodHours"] = 0
					}
					if value, ok := action["notificationTemplateId"]; !ok || value == "" {
						action["notificationTemplateId"] = "00000000-0000-0000-0000-000000000000"
					}
					if _, ok := action["notificationMessageCCList"]; !ok {
						action["notificationMessageCCList"] = []any{}
					}
				}
			}
		}
		return httpmock.NewJsonResponse(200, response)
	})
	httpmock.RegisterResponder("PATCH", item+`$`, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode update: %w", err)
		}
		m.Lock()
		defer m.Unlock()
		policy, ok := m.policies[policyID(req)]
		if !ok {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound"}}`), nil
		}
		for key, value := range body {
			policy[key] = value
		}
		return httpmock.NewStringResponse(204, ""), nil
	})
	for _, action := range []string{"assign", "scheduleActionsForRules"} {
		httpmock.RegisterResponder(
			"POST",
			item+`/`+action+`$`,
			func(req *http.Request) (*http.Response, error) {
				var body map[string]any
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					return nil, fmt.Errorf("decode action: %w", err)
				}
				m.Lock()
				defer m.Unlock()
				policy, ok := m.policies[policyID(req)]
				if !ok {
					return httpmock.NewStringResponse(
						404,
						`{"error":{"code":"ResourceNotFound"}}`,
					), nil
				}
				if strings.HasSuffix(req.URL.Path, "/assign") {
					policy["assignments"] = body["assignments"]
					return httpmock.NewJsonResponse(
						200,
						map[string]any{"value": body["assignments"]},
					)
				}
				policy["scheduledActionsForRule"] = body["deviceComplianceScheduledActionForRules"]
				return httpmock.NewStringResponse(204, ""), nil
			},
		)
	}
	httpmock.RegisterResponder("DELETE", item+`$`, func(req *http.Request) (*http.Response, error) {
		m.Lock()
		defer m.Unlock()
		delete(m.policies, policyID(req))
		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterErrorMocks registers a failed policy creation response.
func (m *MacosDeviceCompliancePolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies",
		httpmock.NewStringResponder(
			400,
			`{"error":{"code":"BadRequest","message":"Invalid compliance policy"}}`,
		),
	)
}

// CleanupMockState clears policies between tests.
func (m *MacosDeviceCompliancePolicyMock) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.policies = make(map[string]map[string]any)
}

func policyID(req *http.Request) string {
	return strings.Split(req.URL.Path, "/")[4]
}
