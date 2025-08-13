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
	featureProfiles map[string]map[string]interface{}
}

func init() {
	mockState.featureProfiles = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_feature_update_profile", &WindowsFeatureUpdateProfileMock{})
}

type WindowsFeatureUpdateProfileMock struct{}

var _ mocks.MockRegistrar = (*WindowsFeatureUpdateProfileMock)(nil)

func (m *WindowsFeatureUpdateProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.featureProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		list := make([]map[string]interface{}, 0, len(mockState.featureProfiles))
		for _, v := range mockState.featureProfiles {
			c := map[string]interface{}{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}
		mockState.Unlock()
		return httpmock.NewJsonResponse(200, map[string]interface{}{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles", "value": list})
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		profile, ok := mockState.featureProfiles[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}
		c := map[string]interface{}{}
		for k, v := range profile {
			c[k] = v
		}
		return httpmock.NewJsonResponse(200, c)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		id := uuid.New().String()
		profile := map[string]interface{}{"@odata.type": "#microsoft.graph.windowsFeatureUpdateProfile", "id": id, "displayName": body["displayName"], "featureUpdateVersion": body["featureUpdateVersion"]}
		if v, ok := body["description"]; ok {
			profile["description"] = v
		}
		if v, ok := body["installLatestWindows10OnWindows11IneligibleDevice"]; ok {
			profile["installLatestWindows10OnWindows11IneligibleDevice"] = v
		}
		if v, ok := body["installFeatureUpdatesOptional"]; ok {
			profile["installFeatureUpdatesOptional"] = v
		}
		if v, ok := body["roleScopeTagIds"]; ok {
			profile["roleScopeTagIds"] = v
		} else {
			profile["roleScopeTagIds"] = []string{"0"}
		}
		if v, ok := body["rolloutSettings"]; ok {
			profile["rolloutSettings"] = v
		}
		profile["assignments"] = []interface{}{}
		mockState.Lock()
		mockState.featureProfiles[id] = profile
		mockState.Unlock()
		return httpmock.NewJsonResponse(201, profile)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_feature_update_profile_error.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}
		mockState.Lock()
		existing, ok := mockState.featureProfiles[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}
		for k, v := range body {
			existing[k] = v
		}
		mockState.featureProfiles[id] = existing
		mockState.Unlock()
		return factories.SuccessResponse(200, existing)(req)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		mockState.Lock()
		if existing, ok := mockState.featureProfiles[id]; ok {
			assignments, _ := body["assignments"].([]interface{})
			if assignments == nil {
				assignments = []interface{}{}
			}
			existing["assignments"] = assignments
			mockState.featureProfiles[id] = existing
		}
		mockState.Unlock()
		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.featureProfiles, id)
		mockState.Unlock()
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsFeatureUpdateProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.featureProfiles = make(map[string]map[string]interface{})
	mockState.Unlock()
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]interface{}{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles", "value": []interface{}{}})
	})
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_feature_update_profile_error.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
		var errObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsFeatureUpdateProfileMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.featureProfiles {
		delete(mockState.featureProfiles, id)
	}
}
