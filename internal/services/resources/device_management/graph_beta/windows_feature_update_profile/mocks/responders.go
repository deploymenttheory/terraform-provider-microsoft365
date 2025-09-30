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
	featureProfiles map[string]map[string]any
}

func init() {
	mockState.featureProfiles = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_feature_update_profile", &WindowsFeatureUpdateProfileMock{})
}

type WindowsFeatureUpdateProfileMock struct{}

var _ mocks.MockRegistrar = (*WindowsFeatureUpdateProfileMock)(nil)

func (m *WindowsFeatureUpdateProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.featureProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		list := make([]map[string]any, 0, len(mockState.featureProfiles))
		for _, v := range mockState.featureProfiles {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}
		mockState.Unlock()
		return httpmock.NewJsonResponse(200, map[string]any{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles", "value": list})
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		profile, ok := mockState.featureProfiles[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Return the stored profile as-is (it was created with the right fields)
		c := map[string]any{}
		for k, v := range profile {
			c[k] = v
		}
		return httpmock.NewJsonResponse(200, c)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_feature_update_profile_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		id := uuid.New().String()
		profile := map[string]any{
			"@odata.type":          "#microsoft.graph.windowsFeatureUpdateProfile",
			"id":                   id,
			"displayName":          body["displayName"],
			"featureUpdateVersion": body["featureUpdateVersion"],
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
			"assignments":          []any{},
		}

		// Only add optional fields if they were provided in the request
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

		mockState.Lock()
		mockState.featureProfiles[id] = profile
		mockState.Unlock()
		return httpmock.NewJsonResponse(201, profile)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_feature_update_profile_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}
		mockState.Lock()
		existing, ok := mockState.featureProfiles[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Use success JSON file as template for updates
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_feature_update_profile_success.json")
		if err != nil {
			// Fallback to dynamic response if file not found
			for k, v := range body {
				existing[k] = v
			}
			mockState.featureProfiles[id] = existing
			mockState.Unlock()
			return factories.SuccessResponse(200, existing)(req)
		}

		var templateObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &templateObj)

		// Start with existing data
		for k, v := range existing {
			templateObj[k] = v
		}

		// Apply updates from request
		for k, v := range body {
			templateObj[k] = v
		}

		mockState.featureProfiles[id] = templateObj
		mockState.Unlock()
		return factories.SuccessResponse(200, templateObj)(req)
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		mockState.Lock()
		if existing, ok := mockState.featureProfiles[id]; ok {
			assignments, _ := body["assignments"].([]any)
			if assignments == nil {
				assignments = []any{}
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
	mockState.featureProfiles = make(map[string]map[string]any)
	mockState.Unlock()
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles", "value": []any{}})
	})
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_feature_update_profile_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
		var errObj map[string]any
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
