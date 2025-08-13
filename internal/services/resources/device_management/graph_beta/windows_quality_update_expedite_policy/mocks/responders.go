package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	expeditePolicies map[string]map[string]interface{}
}

func init() {
	mockState.expeditePolicies = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_quality_update_expedite_policy", &WindowsQualityUpdateExpeditePolicyMock{})
}

type WindowsQualityUpdateExpeditePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsQualityUpdateExpeditePolicyMock)(nil)

func (m *WindowsQualityUpdateExpeditePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.expeditePolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]interface{}, 0, len(mockState.expeditePolicies))
			for _, v := range mockState.expeditePolicies {
				c := map[string]interface{}{}
				for k, vv := range v {
					c[k] = vv
				}
				list = append(list, c)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]interface{}{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdateProfiles", "value": list})
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			policy, ok := mockState.expeditePolicies[id]
			mockState.Unlock()
			if !ok {
				if strings.Contains(id, "minimal") {
					// Return a minimal-like shape
					resp := map[string]interface{}{"@odata.type": "#microsoft.graph.windowsQualityUpdateProfile", "id": id, "displayName": "Test Minimal Windows Quality Update Expedite Policy - Unique", "roleScopeTagIds": []interface{}{"0"}, "assignments": []interface{}{}}
					return factories.SuccessResponse(200, resp)(req)
				}
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			c := map[string]interface{}{}
			for k, v := range policy {
				c[k] = v
			}
			return httpmock.NewJsonResponse(200, c)
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			id := uuid.New().String()
			policy := map[string]interface{}{"@odata.type": "#microsoft.graph.windowsQualityUpdateProfile", "id": id, "displayName": body["displayName"]}
			if v, ok := body["description"]; ok {
				policy["description"] = v
			}
			if v, ok := body["roleScopeTagIds"]; ok {
				policy["roleScopeTagIds"] = v
			} else {
				policy["roleScopeTagIds"] = []string{"0"}
			}
			// expedite settings are optional
			if v, ok := body["expeditedUpdateSettings"]; ok {
				policy["expeditedUpdateSettings"] = v
			}
			policy["assignments"] = []interface{}{}
			mockState.Lock()
			mockState.expeditePolicies[id] = policy
			mockState.Unlock()
			return httpmock.NewJsonResponse(201, policy)
		})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_expedite_policy_error.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(400, errObj)
			}
			mockState.Lock()
			existing, ok := mockState.expeditePolicies[id]
			if !ok {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			for k, v := range body {
				existing[k] = v
			}
			mockState.expeditePolicies[id] = existing
			mockState.Unlock()
			return factories.SuccessResponse(200, existing)(req)
		})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			mockState.Lock()
			if existing, ok := mockState.expeditePolicies[id]; ok {
				assignments, _ := body["assignments"].([]interface{})
				if assignments == nil {
					assignments = []interface{}{}
				}
				existing["assignments"] = assignments
				mockState.expeditePolicies[id] = existing
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.expeditePolicies, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsQualityUpdateExpeditePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.expeditePolicies = make(map[string]map[string]interface{})
	mockState.Unlock()
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdateProfiles", "value": []interface{}{}})
	})
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_expedite_policy_error.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsQualityUpdateExpeditePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.expeditePolicies {
		delete(mockState.expeditePolicies, id)
	}
}
