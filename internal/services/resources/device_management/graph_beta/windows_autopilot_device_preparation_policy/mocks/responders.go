package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	autopilotPolicies map[string]map[string]any
}

func init() {
	mockState.autopilotPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_autopilot_device_preparation_policy", &WindowsAutopilotDevicePreparationPolicyMock{})
}

type WindowsAutopilotDevicePreparationPolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopilotDevicePreparationPolicyMock)(nil)

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.autopilotPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Group validation - called during validateRequest
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		responseObj["id"] = groupId
		responseObj["securityEnabled"] = true

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Create configuration policy - POST /deviceManagement/configurationPolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_device_preparation_policy_success.json")
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

		// Store in mock state
		mockState.Lock()
		if mockState.autopilotPolicies == nil {
			mockState.autopilotPolicies = make(map[string]map[string]any)
		}
		mockState.autopilotPolicies[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 3. Set enrollment time device membership target - POST /deviceManagement/configurationPolicies/{id}/setEnrollmentTimeDeviceMembershipTarget
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/setEnrollmentTimeDeviceMembershipTarget$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(204, ""), nil
	})

	// 4. Assign policy - POST /deviceManagement/configurationPolicies/{id}/assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/assign$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#Collection(microsoft.graph.deviceManagementConfigurationPolicyAssignment)",
			"value":          []any{},
		})
	})

	// 5. Read policy - GET /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		policy, exists := mockState.autopilotPolicies[id]
		mockState.Unlock()

		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_device_preparation_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Load base response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with stored policy data
		for k, v := range policy {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 6. Read policy settings - GET /deviceManagement/configurationPolicies/{id}/settings
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/settings$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_settings.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 7. Read policy assignments - GET /deviceManagement/configurationPolicies/{id}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_assignments.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 8. Update policy - PATCH /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, exists := mockState.autopilotPolicies[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Policy not found"}}`), nil
		}

		// Update existing policy
		for k, v := range body {
			existing[k] = v
		}
		existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.autopilotPolicies[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, existing)
	})

	// 9. Delete policy - DELETE /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.autopilotPolicies, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// 10. List policies - GET /deviceManagement/configurationPolicies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_list.json")
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
	})
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.autopilotPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Make groups validation fail during the validation step
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Group not found"}}`), nil
	})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_device_preparation_policy_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []any{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_device_preparation_policy_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_device_preparation_policy_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsAutopilotDevicePreparationPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.autopilotPolicies {
		delete(mockState.autopilotPolicies, id)
	}
}
