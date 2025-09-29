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
	platformScripts map[string]map[string]any
}

func init() {
	mockState.platformScripts = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_platform_script", &WindowsPlatformScriptMock{})
}

type WindowsPlatformScriptMock struct{}

var _ mocks.MockRegistrar = (*WindowsPlatformScriptMock)(nil)

func (m *WindowsPlatformScriptMock) RegisterMocks() {
	mockState.Lock()
	mockState.platformScripts = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.platformScripts) == 0 {
			// Return empty list if no scripts exist
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_scripts_list.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["value"] = []interface{}{}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing scripts
		list := make([]map[string]any, 0, len(mockState.platformScripts))
		for _, v := range mockState.platformScripts {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_scripts_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = list
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		script, ok := mockState.platformScripts[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_platform_script_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Get the appropriate template based on script type
		var jsonTemplate string
		if script["runAsAccount"] == "user" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_script_user.json")
			jsonTemplate = jsonStr
		} else if len(script["assignments"].([]interface{})) > 0 {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_script_with_assignments.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_script.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with actual script values
		for k, v := range script {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		// Choose the appropriate response template
		var jsonTemplate string
		if body["runAsAccount"] == "user" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_platform_script_user_success.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_platform_script_system_success.json")
			jsonTemplate = jsonStr
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Only include fields that were provided in the request
		responseObj["id"] = id
		if v, ok := body["displayName"]; ok {
			responseObj["displayName"] = v
		}
		if v, ok := body["runAsAccount"]; ok {
			responseObj["runAsAccount"] = v
		}
		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}
		if v, ok := body["fileName"]; ok {
			responseObj["fileName"] = v
		}
		if v, ok := body["scriptContent"]; ok {
			responseObj["scriptContent"] = v
		}
		if v, ok := body["enforceSignatureCheck"]; ok {
			responseObj["enforceSignatureCheck"] = v
		}
		if v, ok := body["runAs32Bit"]; ok {
			responseObj["runAs32Bit"] = v
		}
		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		// Assignments are handled separately via the /assign endpoint
		// Always start with empty assignments for new resources
		responseObj["assignments"] = []interface{}{}

		// Store in mock state
		mockState.Lock()
		mockState.platformScripts[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_platform_script_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		existing, ok := mockState.platformScripts[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_platform_script_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Choose the appropriate response template based on what's being updated
		var jsonTemplate string
		if _, hasRoleScopeTagIds := body["roleScopeTagIds"]; hasRoleScopeTagIds {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_platform_script_tags_success.json")
			jsonTemplate = jsonStr
		} else {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_platform_script_success.json")
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

		mockState.platformScripts[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_assign/post_windows_platform_script_assign_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		mockState.Lock()
		if existing, ok := mockState.platformScripts[id]; ok {
			// The assignment constructor sends assignments via deviceManagementScriptAssignments field
			assignments, ok := body["deviceManagementScriptAssignments"].([]interface{})
			if !ok {
				// Fallback to assignments field for compatibility
				assignments, _ = body["assignments"].([]interface{})
			}
			if assignments == nil {
				assignments = []interface{}{}
			}
			existing["assignments"] = assignments
			mockState.platformScripts[id] = existing
		}
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.platformScripts, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsPlatformScriptMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.platformScripts = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_platform_scripts_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = []interface{}{}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_platform_script_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_platform_script_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsPlatformScriptMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.platformScripts {
		delete(mockState.platformScripts, id)
	}
}
