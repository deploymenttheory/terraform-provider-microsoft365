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
	enrollmentStatusPages map[string]map[string]interface{}
}

func init() {
	mockState.enrollmentStatusPages = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_enrollment_status_page", &WindowsEnrollmentStatusPageMock{})
}

type WindowsEnrollmentStatusPageMock struct{}

var _ mocks.MockRegistrar = (*WindowsEnrollmentStatusPageMock)(nil)

func (m *WindowsEnrollmentStatusPageMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentStatusPages = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Mock the mobile apps endpoint for validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps.*`, func(req *http.Request) (*http.Response, error) {
		// Return mock mobile apps that include the test app IDs used in unit tests
		mockApps := map[string]interface{}{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps",
			"@odata.count": 2,
			"value": []interface{}{
				map[string]interface{}{
					"@odata.type": "#microsoft.graph.win32LobApp",
					"id": "12345678-1234-1234-1234-123456789012",
					"displayName": "Test App 1",
					"description": "Test application 1 for unit testing",
					"publisher": "Test Publisher",
					"publishingState": "published",
				},
				map[string]interface{}{
					"@odata.type": "#microsoft.graph.winGetApp",
					"id": "87654321-4321-4321-4321-210987654321",
					"displayName": "Test App 2",
					"description": "Test application 2 for unit testing",
					"publisher": "Test Publisher",
					"publishingState": "published",
				},
			},
		}
		return httpmock.NewJsonResponse(200, mockApps)
	})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.enrollmentStatusPages) == 0 {
			// Return empty list if no enrollment status pages exist
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_pages_list.json")
			var responseObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["value"] = []interface{}{}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing enrollment status pages
		list := make([]map[string]interface{}, 0, len(mockState.enrollmentStatusPages))
		for _, v := range mockState.enrollmentStatusPages {
			c := map[string]interface{}{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_pages_list.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = list
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		enrollmentStatusPage, ok := mockState.enrollmentStatusPages[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_enrollment_status_page_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Get the appropriate template based on configuration
		var jsonTemplate string
		if len(enrollmentStatusPage["assignments"].([]interface{})) > 0 {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_page_with_assignments.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_page.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with actual enrollment status page values
		for k, v := range enrollmentStatusPage {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		// Use standard response template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_success.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]

		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}
		if v, ok := body["priority"]; ok {
			responseObj["priority"] = v
		} else {
			responseObj["priority"] = 0
		}
		if v, ok := body["showInstallationProgress"]; ok {
			responseObj["showInstallationProgress"] = v
		} else {
			responseObj["showInstallationProgress"] = true
		}
		if v, ok := body["blockDeviceSetupRetryByUser"]; ok {
			responseObj["blockDeviceSetupRetryByUser"] = v
		} else {
			responseObj["blockDeviceSetupRetryByUser"] = false
		}
		if v, ok := body["allowDeviceResetOnInstallFailure"]; ok {
			responseObj["allowDeviceResetOnInstallFailure"] = v
		} else {
			responseObj["allowDeviceResetOnInstallFailure"] = false
		}
		if v, ok := body["allowLogCollectionOnInstallFailure"]; ok {
			responseObj["allowLogCollectionOnInstallFailure"] = v
		} else {
			responseObj["allowLogCollectionOnInstallFailure"] = false
		}
		if v, ok := body["customErrorMessage"]; ok {
			responseObj["customErrorMessage"] = v
		}
		if v, ok := body["installProgressTimeoutInMinutes"]; ok {
			responseObj["installProgressTimeoutInMinutes"] = v
		} else {
			responseObj["installProgressTimeoutInMinutes"] = 60
		}
		if v, ok := body["allowDeviceUseOnInstallFailure"]; ok {
			responseObj["allowDeviceUseOnInstallFailure"] = v
		} else {
			responseObj["allowDeviceUseOnInstallFailure"] = false
		}
		if v, ok := body["selectedMobileAppIds"]; ok {
			responseObj["selectedMobileAppIds"] = v
		} else {
			responseObj["selectedMobileAppIds"] = []interface{}{}
		}
		if v, ok := body["trackInstallProgressForAutopilotOnly"]; ok {
			responseObj["trackInstallProgressForAutopilotOnly"] = v
		} else {
			responseObj["trackInstallProgressForAutopilotOnly"] = false
		}
		if v, ok := body["disableUserStatusTrackingAfterFirstUser"]; ok {
			responseObj["disableUserStatusTrackingAfterFirstUser"] = v
		} else {
			responseObj["disableUserStatusTrackingAfterFirstUser"] = false
		}
		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		// Handle assignments if provided
		if v, ok := body["assignments"]; ok {
			responseObj["assignments"] = v
		} else {
			// Ensure assignments field exists (empty array by default)
			if _, exists := responseObj["assignments"]; !exists {
				responseObj["assignments"] = []interface{}{}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.enrollmentStatusPages[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_error.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		existing, ok := mockState.enrollmentStatusPages[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_enrollment_status_page_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_enrollment_status_page_success.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

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

		mockState.enrollmentStatusPages[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_assign/post_windows_enrollment_status_page_assign_error.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		if existing, ok := mockState.enrollmentStatusPages[id]; ok {
			assignments, _ := body["enrollmentConfigurationAssignments"].([]interface{})
			if assignments == nil {
				assignments = []interface{}{}
			}
			existing["assignments"] = assignments
			mockState.enrollmentStatusPages[id] = existing
		}
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.enrollmentStatusPages, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsEnrollmentStatusPageMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrollmentStatusPages = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_pages_list.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []interface{}{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_error.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_enrollment_status_page_not_found.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsEnrollmentStatusPageMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.enrollmentStatusPages {
		delete(mockState.enrollmentStatusPages, id)
	}
}