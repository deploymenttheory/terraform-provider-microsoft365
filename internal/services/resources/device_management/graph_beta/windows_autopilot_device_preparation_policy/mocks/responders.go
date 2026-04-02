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

var mockState struct {
	sync.Mutex
	autopilotPolicies map[string]map[string]any
	policySettings    map[string][]any // policy ID -> settings array stored from POST/PUT body
	policyAssignments map[string][]any // policy ID -> assignments in GET response format
}

func init() {
	mockState.autopilotPolicies = make(map[string]map[string]any)
	mockState.policySettings = make(map[string][]any)
	mockState.policyAssignments = make(map[string][]any)
	httpmock.RegisterNoResponder(
		httpmock.NewStringResponder(
			404,
			`{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`,
		),
	)
	mocks.GlobalRegistry.Register(
		"windows_autopilot_device_preparation_policy",
		&WindowsAutopilotDevicePreparationPolicyMock{},
	)
}

type WindowsAutopilotDevicePreparationPolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopilotDevicePreparationPolicyMock)(nil)

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.autopilotPolicies = make(map[string]map[string]any)
	mockState.policySettings = make(map[string][]any)
	mockState.policyAssignments = make(map[string][]any)
	mockState.Unlock()

	// 0a. Mobile apps - called during validateAllowedApps
	// Returns the apps referenced in test fixtures (IDs 001-003) with matching types.
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_mobile_apps.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 0b. Device management scripts - called during validateAllowedScripts
	// Returns the scripts referenced in test fixtures (IDs 004-006).
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceManagementScripts",
				"value": []map[string]any{
					{
						"id":          "00000000-0000-0000-0000-000000000004",
						"displayName": "Test Script 4",
						"description": "",
					},
					{
						"id":          "00000000-0000-0000-0000-000000000005",
						"displayName": "Test Script 5",
						"description": "",
					},
					{
						"id":          "00000000-0000-0000-0000-000000000006",
						"displayName": "Test Script 6",
						"description": "",
					},
				},
			})
		},
	)

	// 1a. Group owners - called during validateSecurityGroupOwnership
	// Returns the Intune Provisioning Client as an owner so validation passes.
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+/owners$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_owners.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 1b. Group details - GET /groups/{id}
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			groupId := parts[len(parts)-1]

			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)

			responseObj["id"] = groupId
			responseObj["securityEnabled"] = true

			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 2. Create configuration policy - POST /deviceManagement/configurationPolicies
	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(
					400,
					`{"error":{"code":"BadRequest","message":"Invalid request body"}}`,
				), nil
			}

			id := uuid.New().String()
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_create/post_windows_autopilot_device_preparation_policy_success.json",
			)
			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
				responseObj = make(map[string]any)
			}

			// Copy all values from request body
			for k, v := range body {
				responseObj[k] = v
			}
			responseObj["id"] = id
			if _, ok := body["roleScopeTagIds"]; !ok {
				responseObj["roleScopeTagIds"] = []string{"0"}
			}

			mockState.Lock()
			if mockState.autopilotPolicies == nil {
				mockState.autopilotPolicies = make(map[string]map[string]any)
			}
			mockState.autopilotPolicies[id] = responseObj

			// Store settings from request body so GET /settings returns the same values
			if settings, ok := body["settings"]; ok {
				if settingsSlice, ok := settings.([]any); ok {
					mockState.policySettings[id] = settingsSlice
				}
			}
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, responseObj)
		},
	)

	// 3. Set enrollment time device membership target - POST /deviceManagement/configurationPolicies/{id}/setEnrollmentTimeDeviceMembershipTarget
	httpmock.RegisterResponder(
		"POST",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/setEnrollmentTimeDeviceMembershipTarget$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	// 4. Assign policy - POST /deviceManagement/configurationPolicies/{id}/assign
	// Stores assignments dynamically so GET /assignments returns the actual configured assignments.
	httpmock.RegisterResponder(
		"POST",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			policyId := parts[len(parts)-2]

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err == nil {
				if assignments, ok := body["assignments"]; ok {
					if assignmentSlice, ok := assignments.([]any); ok {
						responseAssignments := make([]any, 0, len(assignmentSlice))
						for _, a := range assignmentSlice {
							if aMap, ok := a.(map[string]any); ok {
								responseAssignment := map[string]any{
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyAssignment",
									"id":          uuid.New().String(),
								}
								if target, ok := aMap["target"]; ok {
									if targetMap, ok := target.(map[string]any); ok {
										// Ensure filter type defaults are present
										if _, ok := targetMap["deviceAndAppManagementAssignmentFilterType"]; !ok {
											targetMap["deviceAndAppManagementAssignmentFilterType"] = "none"
										}
										responseAssignment["target"] = targetMap
									}
								}
								responseAssignments = append(responseAssignments, responseAssignment)
							}
						}
						mockState.Lock()
						mockState.policyAssignments[policyId] = responseAssignments
						mockState.Unlock()
					}
				}
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#Collection(microsoft.graph.deviceManagementConfigurationPolicyAssignment)",
				"value":          []any{},
			})
		},
	)

	// 5. Read policy - GET /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policy, exists := mockState.autopilotPolicies[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile(
					"../tests/responses/validate_delete/get_windows_autopilot_device_preparation_policy_not_found.json",
				)
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}

			// Load base response
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)

			// Override with stored policy data
			for k, v := range policy {
				responseObj[k] = v
			}

			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 6. Read policy settings - GET /deviceManagement/configurationPolicies/{id}/settings
	// Returns dynamically stored settings from the POST/PUT body so the read-back matches what was configured.
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			mockState.Lock()
			settings, hasSettings := mockState.policySettings[id]
			mockState.Unlock()

			if hasSettings && len(settings) > 0 {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": fmt.Sprintf(
						"https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('%s')/settings",
						id,
					),
					"value": settings,
				})
			}

			// Fallback to static file
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_settings.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 7. Read policy assignments - GET /deviceManagement/configurationPolicies/{id}/assignments
	// Returns dynamically stored assignments from the POST /assign body so the read-back matches what was configured.
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			mockState.Lock()
			assignments, hasAssignments := mockState.policyAssignments[id]
			mockState.Unlock()

			if hasAssignments {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('...')/assignments",
					"value":          assignments,
				})
			}

			// Fallback to static file
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_assignments_empty.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// 8. Update policy - PUT /deviceManagement/configurationPolicies/{id}
	// The resource uses PUT (not PATCH) because the Graph API does not support PATCH
	// on the 'settings' navigation property of deviceManagementConfigurationPolicy.
	httpmock.RegisterResponder(
		"PUT",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(
					400,
					`{"error":{"code":"BadRequest","message":"Invalid request body"}}`,
				), nil
			}

			mockState.Lock()
			existing, exists := mockState.autopilotPolicies[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(
					404,
					`{"error":{"code":"NotFound","message":"Policy not found"}}`,
				), nil
			}

			// Replace existing policy (PUT semantics)
			for k, v := range body {
				existing[k] = v
			}
			existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
			mockState.autopilotPolicies[id] = existing

			// Update stored settings from PUT body
			if settings, ok := body["settings"]; ok {
				if settingsSlice, ok := settings.([]any); ok {
					mockState.policySettings[id] = settingsSlice
				}
			}
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, existing)
		},
	)

	// 9. Delete policy - DELETE /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder(
		"DELETE",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.autopilotPolicies, id)
			delete(mockState.policySettings, id)
			delete(mockState.policyAssignments, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	// 10. List policies - GET /deviceManagement/configurationPolicies
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_list.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)

			if mockState.autopilotPolicies == nil || len(mockState.autopilotPolicies) == 0 {
				responseObj["value"] = []any{}
			} else {
				list := make([]map[string]any, 0, len(mockState.autopilotPolicies))
				for _, policy := range mockState.autopilotPolicies {
					list = append(list, policy)
				}
				responseObj["value"] = list
			}

			return httpmock.NewJsonResponse(200, responseObj)
		},
	)
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.autopilotPolicies = make(map[string]map[string]any)
	mockState.policySettings = make(map[string][]any)
	mockState.policyAssignments = make(map[string][]any)
	mockState.Unlock()

	// Mobile apps - allow validateAllowedApps to pass so POST is reached
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_mobile_apps.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	// Make groups validation fail during the validation step
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				404,
				`{"error":{"code":"NotFound","message":"Group not found"}}`,
			), nil
		},
	)

	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_list.json",
			)
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["value"] = []any{}
			return httpmock.NewJsonResponse(200, responseObj)
		},
	)

	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_create/post_windows_autopilot_device_preparation_policy_error.json",
			)
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		},
	)

	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile(
				"../tests/responses/validate_delete/get_windows_autopilot_device_preparation_policy_not_found.json",
			)
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		},
	)
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.autopilotPolicies {
		delete(mockState.autopilotPolicies, id)
	}
	for id := range mockState.policySettings {
		delete(mockState.policySettings, id)
	}
	for id := range mockState.policyAssignments {
		delete(mockState.policyAssignments, id)
	}
}
