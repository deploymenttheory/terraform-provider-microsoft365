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
	enrollmentStatusPages map[string]map[string]any
}

func init() {
	mockState.enrollmentStatusPages = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_enrollment_status_page", &WindowsEnrollmentStatusPageMock{})
}

type WindowsEnrollmentStatusPageMock struct{}

var _ mocks.MockRegistrar = (*WindowsEnrollmentStatusPageMock)(nil)

func (m *WindowsEnrollmentStatusPageMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentStatusPages = make(map[string]map[string]any)
	mockState.Unlock()

	// Mock the mobile apps endpoint for validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps.*`, func(req *http.Request) (*http.Response, error) {
		// Return mock mobile apps that include the test app IDs used in unit tests
		mockApps := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps",
			"@odata.count":   5,
			"value": []interface{}{
				map[string]any{
					"@odata.type":     "#microsoft.graph.win32LobApp",
					"id":              "12345678-1234-1234-1234-123456789012",
					"displayName":     "Test App 1",
					"description":     "Test application 1 for unit testing",
					"publisher":       "Test Publisher",
					"publishingState": "published",
				},
				map[string]any{
					"@odata.type":     "#microsoft.graph.winGetApp",
					"id":              "87654321-4321-4321-4321-210987654321",
					"displayName":     "Test App 2",
					"description":     "Test application 2 for unit testing",
					"publisher":       "Test Publisher",
					"publishingState": "published",
				},
				map[string]any{
					"@odata.type":     "#microsoft.graph.win32LobApp",
					"id":              "e4938228-aab3-493b-a9d5-8250aa8e9d55",
					"displayName":     "Test App 3",
					"description":     "Test application 3 for unit testing",
					"publisher":       "Test Publisher",
					"publishingState": "published",
				},
				map[string]any{
					"@odata.type":     "#microsoft.graph.win32LobApp",
					"id":              "e83d36e1-3ff2-4567-90d9-940919184ad5",
					"displayName":     "Test App 4",
					"description":     "Test application 4 for unit testing",
					"publisher":       "Test Publisher",
					"publishingState": "published",
				},
				map[string]any{
					"@odata.type":     "#microsoft.graph.win32LobApp",
					"id":              "cd4486df-05cc-42bd-8c34-67ac20e10166",
					"displayName":     "Test App 5",
					"description":     "Test application 5 for unit testing",
					"publisher":       "Test Publisher",
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
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["value"] = []interface{}{}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing enrollment status pages
		list := make([]map[string]any, 0, len(mockState.enrollmentStatusPages))
		for _, v := range mockState.enrollmentStatusPages {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_pages_list.json")
		var responseObj map[string]any
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
			var errObj map[string]any
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

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with actual enrollment status page values
		for k, v := range enrollmentStatusPage {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		// Use standard response template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_success.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]

		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}
		if v, ok := body["priority"]; ok {
			responseObj["priority"] = v
		}
		if v, ok := body["showInstallationProgress"]; ok {
			responseObj["showInstallationProgress"] = v
		}
		if v, ok := body["blockDeviceSetupRetryByUser"]; ok {
			responseObj["blockDeviceSetupRetryByUser"] = v
		}
		if v, ok := body["allowDeviceResetOnInstallFailure"]; ok {
			responseObj["allowDeviceResetOnInstallFailure"] = v
		}
		if v, ok := body["allowLogCollectionOnInstallFailure"]; ok {
			responseObj["allowLogCollectionOnInstallFailure"] = v
		}
		if v, ok := body["customErrorMessage"]; ok {
			responseObj["customErrorMessage"] = v
		}
		if v, ok := body["installProgressTimeoutInMinutes"]; ok {
			responseObj["installProgressTimeoutInMinutes"] = v
		}
		if v, ok := body["allowDeviceUseOnInstallFailure"]; ok {
			responseObj["allowDeviceUseOnInstallFailure"] = v
		}
		if v, ok := body["selectedMobileAppIds"]; ok {
			responseObj["selectedMobileAppIds"] = v
		}
		if v, ok := body["trackInstallProgressForAutopilotOnly"]; ok {
			responseObj["trackInstallProgressForAutopilotOnly"] = v
		}
		if v, ok := body["disableUserStatusTrackingAfterFirstUser"]; ok {
			responseObj["disableUserStatusTrackingAfterFirstUser"] = v
		}
		if v, ok := body["installQualityUpdates"]; ok {
			responseObj["installQualityUpdates"] = v
		}
		if v, ok := body["allowNonBlockingAppInstallation"]; ok {
			responseObj["allowNonBlockingAppInstallation"] = v
		}
		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
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
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		existing, ok := mockState.enrollmentStatusPages[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_enrollment_status_page_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_enrollment_status_page_success.json")
		var responseObj map[string]any
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
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_assign/post_windows_enrollment_status_page_assign_error.json")
			var errObj map[string]any
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
	mockState.enrollmentStatusPages = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_enrollment_status_pages_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []interface{}{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_enrollment_status_page_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_enrollment_status_page_not_found.json")
		var errObj map[string]any
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
