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
	autopilotProfiles map[string]map[string]interface{}
}

func init() {
	mockState.autopilotProfiles = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_autopilot_deployment_profile", &WindowsAutopilotDeploymentProfileMock{})
}

type WindowsAutopilotDeploymentProfileMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopilotDeploymentProfileMock)(nil)

func (m *WindowsAutopilotDeploymentProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.autopilotProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()

	// 1. Create autopilot deployment profile - POST /deviceManagement/windowsAutopilotDeploymentProfiles
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_deployment_profile_success.json")
		var responseObj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			responseObj = make(map[string]interface{})
		}

		// Copy all values from request body
		for k, v := range body {
			responseObj[k] = v
		}
		responseObj["id"] = id

		// Store in mock state
		mockState.Lock()
		if mockState.autopilotProfiles == nil {
			mockState.autopilotProfiles = make(map[string]map[string]interface{})
		}
		mockState.autopilotProfiles[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 2. Read autopilot deployment profile - GET /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		profile, exists := mockState.autopilotProfiles[id]
		mockState.Unlock()
		
		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_deployment_profile_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Determine which response file to use based on profile name
		var responseFile string
		if displayName, ok := profile["displayName"].(string); ok && strings.Contains(displayName, "maximal") {
			responseFile = "../tests/responses/validate_get_maximal/get_windows_autopilot_deployment_profile.json"
		} else {
			responseFile = "../tests/responses/validate_get/get_windows_autopilot_deployment_profile.json"
		}

		// Load base response
		jsonStr, _ := helpers.ParseJSONFile(responseFile)
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with stored profile data
		for k, v := range profile {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 3. Read autopilot deployment profile assignments - GET /deviceManagement/windowsAutopilotDeploymentProfiles/{id}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F-]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2] // Get profile ID from URL
		
		mockState.Lock()
		profile, exists := mockState.autopilotProfiles[id]
		mockState.Unlock()
		
		// Determine which response file to use based on profile name
		var assignmentsFile string
		if exists {
			if displayName, ok := profile["displayName"].(string); ok && strings.Contains(displayName, "maximal") {
				assignmentsFile = "../tests/responses/validate_get_maximal/get_windows_autopilot_deployment_profile_assignments.json"
			} else {
				assignmentsFile = "../tests/responses/validate_get/get_windows_autopilot_deployment_profile_assignments.json"
			}
		} else {
			assignmentsFile = "../tests/responses/validate_get/get_windows_autopilot_deployment_profile_assignments.json"
		}
		
		jsonStr, _ := helpers.ParseJSONFile(assignmentsFile)
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 4. Update autopilot deployment profile - PATCH /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, exists := mockState.autopilotProfiles[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Profile not found"}}`), nil
		}

		// Update existing profile
		for k, v := range body {
			existing[k] = v
		}
		existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.autopilotProfiles[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, existing)
	})

	// 5. Delete autopilot deployment profile - DELETE /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		delete(mockState.autopilotProfiles, id)
		mockState.Unlock()
		
		return httpmock.NewStringResponse(204, ""), nil
	})

	// 6. List autopilot deployment profiles - GET /deviceManagement/windowsAutopilotDeploymentProfiles
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_deployment_profile_list.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		
		if mockState.autopilotProfiles == nil || len(mockState.autopilotProfiles) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			list := make([]map[string]interface{}, 0, len(mockState.autopilotProfiles))
			for _, profile := range mockState.autopilotProfiles {
				list = append(list, profile)
			}
			responseObj["value"] = list
		}
		
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *WindowsAutopilotDeploymentProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.autopilotProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_deployment_profile_list.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []interface{}{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_deployment_profile_error.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_deployment_profile_not_found.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsAutopilotDeploymentProfileMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.autopilotProfiles {
		delete(mockState.autopilotProfiles, id)
	}
}