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
	compliancePolicies map[string]map[string]interface{}
	assignments        map[string][]interface{}
	scheduledActions   map[string][]interface{}
}

func init() {
	mockState.compliancePolicies = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.scheduledActions = make(map[string][]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_device_compliance_policy", &WindowsDeviceCompliancePolicyMock{})
}

type WindowsDeviceCompliancePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsDeviceCompliancePolicyMock)(nil)

func (m *WindowsDeviceCompliancePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.compliancePolicies = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.scheduledActions = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic group mocks for assignment validation
	m.registerGroupMocks()

	// GET /deviceManagement/deviceCompliancePolicies - List policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policies_list.json")
		if err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}

		var responseObj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.compliancePolicies) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			list := make([]map[string]interface{}, 0, len(mockState.compliancePolicies))
			for _, v := range mockState.compliancePolicies {
				c := map[string]interface{}{}
				for k, vv := range v {
					c[k] = vv
				}
				list = append(list, c)
			}
			responseObj["value"] = list
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /deviceManagement/deviceCompliancePolicies/{id} - Get specific policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		policy, ok := mockState.compliancePolicies[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_policy_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policy.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with actual policy values
		for k, v := range policy {
			responseObj[k] = v
		}

		// Check for expand parameter to include assignments
		if strings.Contains(req.URL.RawQuery, "expand=assignments") {
			mockState.Lock()
			assignments, ok := mockState.assignments[id]
			mockState.Unlock()
			if ok && len(assignments) > 0 {
				responseObj["assignments"] = assignments
			} else {
				responseObj["assignments"] = []interface{}{}
			}
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/deviceCompliancePolicies - Create policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_policy_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]

		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}

		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		// Copy all other fields from the request
		for k, v := range body {
			if k != "id" && k != "displayName" && k != "description" && k != "roleScopeTagIds" {
				responseObj[k] = v
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.compliancePolicies[id] = responseObj
		mockState.assignments[id] = []interface{}{}
		mockState.scheduledActions[id] = []interface{}{}
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /deviceManagement/deviceCompliancePolicies/{id} - Update policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, ok := mockState.compliancePolicies[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_policy_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_device_compliance_policy_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with existing values
		for k, v := range existing {
			responseObj[k] = v
		}

		// Apply updates
		for k, v := range body {
			responseObj[k] = v
			existing[k] = v
		}

		// Update last modified time
		responseObj["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

		mockState.compliancePolicies[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /deviceManagement/deviceCompliancePolicies/{id}/assignments - Get assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]

		mockState.Lock()
		storedAssignments, ok := mockState.assignments[id]
		mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_policy_assignments.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		if !ok || len(storedAssignments) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			responseObj["value"] = storedAssignments
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/deviceCompliancePolicies/{id}/assign - Assign policy
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		// The SDK sends assignments as "assignments"
		if assignments, ok := body["assignments"].([]interface{}); ok {
			mockState.assignments[id] = assignments
		} else {
			mockState.assignments[id] = []interface{}{}
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// POST /deviceManagement/deviceCompliancePolicies/{id}/scheduleActionsForRules - Schedule actions for rules
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+/scheduleActionsForRules$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		scheduledActions := mockState.scheduledActions[id]
		if scheduledActions == nil {
			scheduledActions = []interface{}{}
		}

		// Add the new scheduled action
		if ruleName, ok := body["ruleName"].(string); ok {
			action := map[string]interface{}{
				"ruleName": ruleName,
			}

			if scheduledActionConfigurations, ok := body["scheduledActionConfigurations"].([]interface{}); ok {
				action["scheduledActionConfigurations"] = scheduledActionConfigurations
			}

			scheduledActions = append(scheduledActions, action)
		}

		mockState.scheduledActions[id] = scheduledActions
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// DELETE /deviceManagement/deviceCompliancePolicies/{id} - Delete policy
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.compliancePolicies, id)
		delete(mockState.assignments, id)
		delete(mockState.scheduledActions, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsDeviceCompliancePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.compliancePolicies = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.scheduledActions = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic group mocks for assignment validation (successful for error tests)
	m.registerGroupMocks()

	// Error response for creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCompliancePolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_policy_error.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	// Error response for GET operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCompliancePolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_policy_not_found.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsDeviceCompliancePolicyMock) registerGroupMocks() {
	// GET /groups/{id} - Get group (for assignment validation)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["id"] = id

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *WindowsDeviceCompliancePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.compliancePolicies {
		delete(mockState.compliancePolicies, id)
	}
	for id := range mockState.assignments {
		delete(mockState.assignments, id)
	}
	for id := range mockState.scheduledActions {
		delete(mockState.scheduledActions, id)
	}
}
