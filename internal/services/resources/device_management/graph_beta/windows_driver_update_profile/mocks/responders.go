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
	driverProfiles map[string]map[string]any
}

func init() {
	mockState.driverProfiles = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_driver_update_profile", &WindowsDriverUpdateProfileMock{})
}

type WindowsDriverUpdateProfileMock struct{}

var _ mocks.MockRegistrar = (*WindowsDriverUpdateProfileMock)(nil)

func (m *WindowsDriverUpdateProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.driverProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsDriverUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.driverProfiles) == 0 {
			// Return empty list if no profiles exist
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profiles_list.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["value"] = []any{}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing profiles
		list := make([]map[string]any, 0, len(mockState.driverProfiles))
		for _, v := range mockState.driverProfiles {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profiles_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = list
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsDriverUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		profile, ok := mockState.driverProfiles[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_driver_update_profile_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Get the appropriate template based on profile type
		var jsonTemplate string
		if profile["approvalType"] == "automatic" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profile_automatic.json")
			jsonTemplate = jsonStr
		} else if len(profile["assignments"].([]any)) > 0 {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profile_with_assignments.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profile.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with actual profile values
		for k, v := range profile {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsDriverUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		// Choose the appropriate response template
		var jsonTemplate string
		if body["approvalType"] == "automatic" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_driver_update_profile_automatic_success.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_driver_update_profile_manual_success.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]
		responseObj["approvalType"] = body["approvalType"]

		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}
		if v, ok := body["deploymentDeferralInDays"]; ok {
			responseObj["deploymentDeferralInDays"] = v
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
				responseObj["assignments"] = []any{}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.driverProfiles[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsDriverUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_driver_update_profile_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		existing, ok := mockState.driverProfiles[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_driver_update_profile_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Choose the appropriate response template based on what's being updated
		var jsonTemplate string
		if _, hasRoleScopeTagIds := body["roleScopeTagIds"]; hasRoleScopeTagIds {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_driver_update_profile_tags_success.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_driver_update_profile_success.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

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

		mockState.driverProfiles[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsDriverUpdateProfiles/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_assign/post_windows_driver_update_profile_assign_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		if existing, ok := mockState.driverProfiles[id]; ok {
			assignments, _ := body["assignments"].([]any)
			if assignments == nil {
				assignments = []any{}
			}
			existing["assignments"] = assignments
			mockState.driverProfiles[id] = existing
		}
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsDriverUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.driverProfiles, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsDriverUpdateProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.driverProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsDriverUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_driver_update_profiles_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []any{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsDriverUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_driver_update_profile_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsDriverUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_driver_update_profile_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsDriverUpdateProfileMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.driverProfiles {
		delete(mockState.driverProfiles, id)
	}
}
