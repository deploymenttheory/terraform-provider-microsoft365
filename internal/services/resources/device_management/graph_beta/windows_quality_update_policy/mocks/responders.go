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
	qualityPolicies map[string]map[string]interface{}
}

func init() {
	mockState.qualityPolicies = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_quality_update_policy", &WindowsQualityUpdatePolicyMock{})
}

type WindowsQualityUpdatePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsQualityUpdatePolicyMock)(nil)

func (m *WindowsQualityUpdatePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.qualityPolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	// List
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]interface{}, 0, len(mockState.qualityPolicies))
			for _, v := range mockState.qualityPolicies {
				copy := map[string]interface{}{}
				for k, vv := range v {
					copy[k] = vv
				}
				list = append(list, copy)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdatePolicies",
				"value":          list,
			})
		})

	// Get by id
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			policy, ok := mockState.qualityPolicies[id]
			mockState.Unlock()
			if !ok {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			copy := map[string]interface{}{}
			for k, v := range policy {
				copy[k] = v
			}
			return httpmock.NewJsonResponse(200, copy)
		})

	// Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			id := uuid.New().String()
			policy := map[string]interface{}{
				"@odata.type": "#microsoft.graph.windowsQualityUpdatePolicy",
				"id":          id,
				"displayName": body["displayName"],
			}
			if v, ok := body["description"]; ok {
				policy["description"] = v
			}
			if v, ok := body["hotpatchEnabled"]; ok {
				policy["hotpatchEnabled"] = v
			}
			if v, ok := body["roleScopeTagIds"]; ok {
				policy["roleScopeTagIds"] = v
			} else {
				policy["roleScopeTagIds"] = []string{"0"}
			}
			policy["assignments"] = []interface{}{}
			mockState.Lock()
			mockState.qualityPolicies[id] = policy
			mockState.Unlock()
			return httpmock.NewJsonResponse(201, policy)
		})

	// Patch
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(400, errObj)
			}
			mockState.Lock()
			existing, ok := mockState.qualityPolicies[id]
			if !ok {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errObj map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			for k, v := range body {
				existing[k] = v
			}
			mockState.qualityPolicies[id] = existing
			mockState.Unlock()
			return factories.SuccessResponse(200, existing)(req)
		})

	// Assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			mockState.Lock()
			if existing, ok := mockState.qualityPolicies[id]; ok {
				assignments, _ := body["assignments"].([]interface{})
				if assignments == nil {
					assignments = []interface{}{}
				}
				existing["assignments"] = assignments
				mockState.qualityPolicies[id] = existing
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.qualityPolicies, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsQualityUpdatePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.qualityPolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdatePolicies",
				"value":          []interface{}{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
			var errObj map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		})
}

func (m *WindowsQualityUpdatePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.qualityPolicies {
		delete(mockState.qualityPolicies, id)
	}
}
