package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	compliancePolicies map[string]map[string]any
}

func init() {
	mockState.compliancePolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_device_compliance_policy", &WindowsDeviceCompliancePolicyMock{})
}

type WindowsDeviceCompliancePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsDeviceCompliancePolicyMock)(nil)

func (m *WindowsDeviceCompliancePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.compliancePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Group validation - called during validateRequest
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		responseObj["id"] = groupId
		responseObj["mailEnabled"] = true

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Create compliance policy - POST /deviceManagement/deviceCompliancePolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_policy_success.json")
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

		// Ensure assignments are preserved if provided
		if assignments, exists := body["assignments"]; exists {
			responseObj["assignments"] = assignments
		}

		// Store in mock state
		mockState.Lock()
		if mockState.compliancePolicies == nil {
			mockState.compliancePolicies = make(map[string]map[string]any)
		}
		mockState.compliancePolicies[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 3. Assign policy - POST /deviceManagement/deviceCompliancePolicies/{id}/assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+/assign$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#Collection(microsoft.graph.deviceCompliancePolicyAssignment)",
			"value":          []interface{}{},
		})
	})

	// 4. Read policy with expand - GET /deviceManagement/deviceCompliancePolicies/{id}?$expand=...
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		policy, exists := mockState.compliancePolicies[id]
		mockState.Unlock()

		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Load base response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policy.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with stored policy data
		for k, v := range policy {
			responseObj[k] = v
		}

		// Handle expand parameters
		queryParams := req.URL.Query()
		expand := queryParams.Get("$expand")
		if strings.Contains(expand, "assignments") {
			// Load assignments from separate JSON file
			assignmentsJsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policy_assignments.json")
			var assignmentsObj map[string]any
			_ = json.Unmarshal([]byte(assignmentsJsonStr), &assignmentsObj)

			if assignmentsValue, exists := assignmentsObj["value"]; exists {
				responseObj["assignments"] = assignmentsValue
			} else {
				responseObj["assignments"] = []interface{}{}
			}
		}
		if strings.Contains(expand, "scheduledActionsForRule") {
			if sched, ok := policy["scheduledActionsForRule"]; ok {
				responseObj["scheduledActionsForRule"] = sched
			}
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 5. Update policy - PATCH /deviceManagement/deviceCompliancePolicies/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, exists := mockState.compliancePolicies[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Policy not found"}}`), nil
		}

		// Update existing policy
		for k, v := range body {
			existing[k] = v
		}
		existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.compliancePolicies[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, existing)
	})

	// 6. Schedule actions for rules - POST /deviceManagement/deviceCompliancePolicies/{id}/scheduleActionsForRules
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+/scheduleActionsForRules$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(204, ""), nil
	})

	// 7. Delete policy - DELETE /deviceManagement/deviceCompliancePolicies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.compliancePolicies, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// 8. List policies - GET /deviceManagement/deviceCompliancePolicies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policies_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		if mockState.compliancePolicies == nil || len(mockState.compliancePolicies) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			list := make([]map[string]any, 0, len(mockState.compliancePolicies))
			for _, policy := range mockState.compliancePolicies {
				list = append(list, policy)
			}
			responseObj["value"] = list
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 9. Get assignments - GET /deviceManagement/deviceCompliancePolicies/{id}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[0-9a-fA-F-]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policy_assignments.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *WindowsDeviceCompliancePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.compliancePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Make groups validation fail during the validation step
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Group not found"}}`), nil
	})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policies_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []interface{}{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_policy_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_policy_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsDeviceCompliancePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.compliancePolicies {
		delete(mockState.compliancePolicies, id)
	}
}
