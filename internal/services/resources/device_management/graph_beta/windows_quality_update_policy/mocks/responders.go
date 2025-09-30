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
	qualityPolicies map[string]map[string]any
}

func init() {
	mockState.qualityPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_quality_update_policy", &WindowsQualityUpdatePolicyMock{})
}

type WindowsQualityUpdatePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsQualityUpdatePolicyMock)(nil)

func (m *WindowsQualityUpdatePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.qualityPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// List
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]any, 0, len(mockState.qualityPolicies))
			for _, v := range mockState.qualityPolicies {
				copy := map[string]any{}
				for k, vv := range v {
					copy[k] = vv
				}
				list = append(list, copy)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{
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
				// For convenience during unit tests, return predefined JSON shapes
				if strings.Contains(id, "minimal") {
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/get_windows_quality_update_policy_minimal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var resp map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					resp["id"] = id
					return factories.SuccessResponse(200, resp)(req)
				} else if strings.Contains(id, "maximal") {
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/get_windows_quality_update_policy_maximal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var resp map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					resp["id"] = id
					return factories.SuccessResponse(200, resp)(req)
				}
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			copy := map[string]any{}
			for k, v := range policy {
				copy[k] = v
			}
			return httpmock.NewJsonResponse(200, copy)
		})

	// Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			id := uuid.New().String()
			policy := map[string]any{
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
			policy["assignments"] = []any{}
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
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(400, errObj)
			}
			mockState.Lock()
			existing, ok := mockState.qualityPolicies[id]
			if !ok {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}
			for k, v := range body {
				existing[k] = v
			}
			mockState.qualityPolicies[id] = existing
			mockState.Unlock()
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/get_windows_quality_update_policy_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updated map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &updated); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}
			for k, v := range existing {
				updated[k] = v
			}
			return factories.SuccessResponse(200, updated)(req)
		})

	// Assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			mockState.Lock()
			if existing, ok := mockState.qualityPolicies[id]; ok {
				assignments, _ := body["assignments"].([]any)
				if assignments == nil {
					assignments = []any{}
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
	mockState.qualityPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdatePolicies",
				"value":          []any{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
			var errObj map[string]any
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
