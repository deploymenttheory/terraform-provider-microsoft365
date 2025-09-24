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

	// 1. Group validation - called during validateRequest
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]
		
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		
		responseObj["id"] = groupId
		responseObj["displayName"] = "Test Group " + groupId[:8]
		
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Create Windows Autopilot Deployment Profile - POST /deviceManagement/windowsAutopilotDeploymentProfiles
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		json.NewDecoder(req.Body).Decode(&requestBody)

		profileId := uuid.New().String()
		
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_deployment_profile_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		
		responseObj["id"] = profileId
		if displayName, ok := requestBody["displayName"].(string); ok {
			responseObj["displayName"] = displayName
		}
		if description, ok := requestBody["description"].(string); ok {
			responseObj["description"] = description
		}

		mockState.Lock()
		mockState.autopilotProfiles[profileId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 3. Get Windows Autopilot Deployment Profile - GET /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		mockState.Lock()
		profile, exists := mockState.autopilotProfiles[profileId]
		mockState.Unlock()

		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_deployment_profile_not_found.json")
			var errorObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errorObj)
			return httpmock.NewJsonResponse(404, errorObj)
		}

		return httpmock.NewJsonResponse(200, profile)
	})

	// 4. Get Windows Autopilot Deployment Profile Assignments - GET /deviceManagement/windowsAutopilotDeploymentProfiles/{id}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/assignments$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_autopilot_deployment_profile_assignments.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 5. Create Assignment - POST /deviceManagement/windowsAutopilotDeploymentProfiles/{id}/assignments
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/assignments$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		json.NewDecoder(req.Body).Decode(&requestBody)

		assignmentId := uuid.New().String()
		responseObj := map[string]interface{}{
			"id": assignmentId,
			"target": requestBody["target"],
		}

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 6. Delete Assignment - DELETE /deviceManagement/windowsAutopilotDeploymentProfiles/{id}/assignments/{assignmentId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/assignments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(204, ""), nil
	})

	// 7. Update Windows Autopilot Deployment Profile - PATCH /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		var requestBody map[string]interface{}
		json.NewDecoder(req.Body).Decode(&requestBody)

		mockState.Lock()
		if profile, exists := mockState.autopilotProfiles[profileId]; exists {
			for k, v := range requestBody {
				profile[k] = v
			}
			mockState.autopilotProfiles[profileId] = profile
		}
		mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_autopilot_deployment_profile_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["id"] = profileId

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 8. Delete Windows Autopilot Deployment Profile - DELETE /deviceManagement/windowsAutopilotDeploymentProfiles/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.autopilotProfiles, profileId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsAutopilotDeploymentProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.autopilotProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_autopilot_deployment_profile_error.json")
		var errorObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errorObj)
		return httpmock.NewJsonResponse(400, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsAutopilotDeploymentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_autopilot_deployment_profile_not_found.json")
		var errorObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errorObj)
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *WindowsAutopilotDeploymentProfileMock) CleanupMockState() {
	mockState.Lock()
	mockState.autopilotProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()
}