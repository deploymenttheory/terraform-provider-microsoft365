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
	appControlPolicies map[string]map[string]any
}

func init() {
	mockState.appControlPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("app_control_for_business_policy", &AppControlForBusinessPolicyMock{})
}

type AppControlForBusinessPolicyMock struct{}

var _ mocks.MockRegistrar = (*AppControlForBusinessPolicyMock)(nil)

func (m *AppControlForBusinessPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.appControlPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Filter by template for app control policies
		if !strings.Contains(req.URL.RawQuery, "4321b946-b76b-4450-8afd-769c08b16ffc_1") {
			return httpmock.NewJsonResponse(200, map[string]any{
				"value": []any{},
			})
		}

		if len(mockState.appControlPolicies) == 0 {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_policy_list.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing policies
		list := make([]map[string]any, 0, len(mockState.appControlPolicies))
		for _, v := range mockState.appControlPolicies {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_policy_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = list
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		policy, ok := mockState.appControlPolicies[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Get the base template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_base.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with actual policy values
		for k, v := range policy {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/settings$`, func(req *http.Request) (*http.Response, error) {
		// Extract policy ID from URL path
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]

		mockState.Lock()
		policy, ok := mockState.appControlPolicies[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Return the settings from the stored policy
		if settings, ok := policy["settings"]; ok {
			if settingsArray, isArray := settings.([]any); isArray {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/settings",
					"value":          settingsArray,
				})
			}
		}

		// Return the settings template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_settings.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		// Extract policy ID from URL path
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]

		mockState.Lock()
		policy, ok := mockState.appControlPolicies[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Return the assignments from the stored policy
		if assignments, ok := policy["assignments"]; ok {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/assignments",
				"value":          assignments,
			})
		}

		// Return empty assignments if none exist
		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/assignments",
			"value":          []any{},
		})
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_policy_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		id := uuid.New().String()

		// Load the success response template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_policy_success.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Only include fields that were provided in the request
		responseObj["id"] = id
		if v, ok := body["name"]; ok {
			responseObj["name"] = v
		}
		if v, ok := body["description"]; ok {
			responseObj["description"] = v
		}
		responseObj["platforms"] = "windows10"
		responseObj["technologies"] = "mdm"
		responseObj["templateReference"] = map[string]any{
			"templateId":             "4321b946-b76b-4450-8afd-769c08b16ffc_1",
			"templateFamily":         "endpointSecurityApplicationControl",
			"templateDisplayName":    "App Control for Business",
			"templateDisplayVersion": "Version 1",
		}

		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		if v, ok := body["settings"]; ok {
			responseObj["settings"] = v
		}

		// Assignments are handled separately via the /assign endpoint
		responseObj["assignments"] = []any{}

		// Store in mock state
		mockState.Lock()
		mockState.appControlPolicies[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, ok := mockState.appControlPolicies[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Load the update success template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_policy_success.json")
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

		mockState.appControlPolicies[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		if existing, ok := mockState.appControlPolicies[id]; ok {
			assignments, ok := body["assignments"].([]any)
			if !ok {
				assignments = []any{}
			}
			existing["assignments"] = assignments
			mockState.appControlPolicies[id] = existing
		}
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.appControlPolicies, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AppControlForBusinessPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.appControlPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_policy_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_policy_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_policy_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *AppControlForBusinessPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.appControlPolicies {
		delete(mockState.appControlPolicies, id)
	}
}
