package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policies map[string]map[string]any
}

func init() {
	mockState.policies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_autopilot_device_preparation_policy_datasource", &WindowsAutopilotDevicePreparationPolicyMock{})
}

type WindowsAutopilotDevicePreparationPolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopilotDevicePreparationPolicyMock)(nil)

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()

	RegisterGetByPolicyIdMock()
	RegisterAssignmentsMock()
	RegisterListAndFilterMocks()
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation"}}`))
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()
}

// RegisterAssignmentsMock handles the assignments navigation property.
func RegisterAssignmentsMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/assignments(\?.*)?$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"value": []any{
					map[string]any{
						"id": "11111111-1111-1111-1111-111111111111_aaaaaaaa-0000-1111-2222-333333333333",
						"target": map[string]any{
							"@odata.type": "#microsoft.graph.groupAssignmentTarget",
							"groupId":     "aaaaaaaa-0000-1111-2222-333333333333",
							"deviceAndAppManagementAssignmentFilterId":   "bbbbbbbb-0000-1111-2222-444444444444",
							"deviceAndAppManagementAssignmentFilterType": "include",
						},
					},
				},
			})
		})
}

// RegisterGetByPolicyIdMock handles direct lookups by policy id.
func RegisterGetByPolicyIdMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+(\?.*)?$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(strings.Split(req.URL.Path, "?")[0], "/")
			policyId := parts[len(parts)-1]

			switch policyId {
			case "11111111-1111-1111-1111-111111111111":
				jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_policy_by_id.json")
				if err != nil {
					return httpmock.NewStringResponse(500, err.Error()), nil
				}
				var responseObj map[string]any
				if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
					return httpmock.NewStringResponse(500, err.Error()), nil
				}
				return httpmock.NewJsonResponse(200, responseObj)

			case "22222222-2222-2222-2222-222222222222":
				return httpmock.NewJsonResponse(200, automaticModePolicy())

			// A settings catalog policy that is not a device preparation policy.
			case "33333333-3333-3333-3333-333333333333":
				return httpmock.NewJsonResponse(200, unrelatedSettingsCatalogPolicy())

			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Configuration policy not found"}}`), nil
			}
		})
}

// RegisterListAndFilterMocks handles list, name filter and custom OData filter lookups.
func RegisterListAndFilterMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies(\?.*)?$`,
		func(req *http.Request) (*http.Response, error) {
			filter := req.URL.Query().Get("$filter")

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_policies_all.json")
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}

			if filter == "" {
				return httpmock.NewJsonResponse(200, responseObj)
			}

			all, _ := responseObj["value"].([]any)

			// Name filter. The fixture deliberately contains an unrelated settings catalog policy
			// sharing the same name so the template reference guard is exercised.
			if strings.Contains(filter, "name eq") {
				name := extractQuotedValue(filter)
				matched := []any{}
				for _, item := range all {
					policy, ok := item.(map[string]any)
					if !ok {
						continue
					}
					if policy["name"] == name {
						matched = append(matched, policy)
					}
				}
				return httpmock.NewJsonResponse(200, map[string]any{"value": matched})
			}

			// Custom OData query, only isAssigned is modelled here.
			if strings.Contains(filter, "isAssigned eq true") {
				matched := []any{}
				for _, item := range all {
					policy, ok := item.(map[string]any)
					if !ok {
						continue
					}
					if policy["isAssigned"] == true {
						matched = append(matched, policy)
					}
				}
				return httpmock.NewJsonResponse(200, map[string]any{"value": matched})
			}

			return httpmock.NewJsonResponse(200, map[string]any{"value": []any{}})
		})
}

// extractQuotedValue returns the first single quoted value in an OData filter expression.
func extractQuotedValue(filter string) string {
	start := strings.Index(filter, "'")
	if start == -1 {
		return ""
	}
	end := strings.LastIndex(filter, "'")
	if end <= start {
		return ""
	}

	return strings.ReplaceAll(filter[start+1:end], "''", "'")
}

func automaticModePolicy() map[string]any {
	return map[string]any{
		"id":                   "22222222-2222-2222-2222-222222222222",
		"name":                 "Autopilot Device Preparation - Automatic",
		"description":          "Automatic mode Autopilot device preparation policy",
		"createdDateTime":      "2024-10-02T14:31:05.0000000Z",
		"lastModifiedDateTime": "2024-11-09T08:44:19.0000000Z",
		"platforms":            "windows10",
		"technologies":         "enrollment",
		"isAssigned":           false,
		"settingCount":         5,
		"roleScopeTagIds":      []any{"0", "9"},
		"templateReference": map[string]any{
			"templateId":             "a6157a7f-aa00-42d9-ac82-7d2479f545db_1",
			"templateFamily":         "enrollmentConfiguration",
			"templateDisplayName":    "Windows Autopilot device preparation",
			"templateDisplayVersion": "1",
		},
	}
}

func unrelatedSettingsCatalogPolicy() map[string]any {
	return map[string]any{
		"id":                   "33333333-3333-3333-3333-333333333333",
		"name":                 "Autopilot Device Preparation - User Driven",
		"description":          "An unrelated settings catalog policy that must never be returned",
		"createdDateTime":      "2024-09-12T10:05:00.0000000Z",
		"lastModifiedDateTime": "2024-09-12T10:05:00.0000000Z",
		"platforms":            "windows10",
		"technologies":         "mdm",
		"isAssigned":           false,
		"settingCount":         2,
		"roleScopeTagIds":      []any{"0"},
		"templateReference": map[string]any{
			"templateId":     "",
			"templateFamily": "none",
		},
	}
}
