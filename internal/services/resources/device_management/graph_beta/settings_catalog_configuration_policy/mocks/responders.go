package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	settingsPolicies map[string]map[string]interface{}
}

func init() {
	mockState.settingsPolicies = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("settings_catalog_configuration_policy", &SettingsCatalogConfigurationPolicyMock{})
}

// SettingsCatalogConfigurationPolicyMock provides mock responses for settings catalog configuration policy operations
type SettingsCatalogConfigurationPolicyMock struct{}

// Ensure interface compliance
var _ mocks.MockRegistrar = (*SettingsCatalogConfigurationPolicyMock)(nil)

// RegisterMocks registers HTTP mock responses
func (m *SettingsCatalogConfigurationPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.settingsPolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	// LIST
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			items := make([]map[string]interface{}, 0, len(mockState.settingsPolicies))
			for _, v := range mockState.settingsPolicies {
				copy := make(map[string]interface{})
				for k, val := range v {
					copy[k] = val
				}
				items = append(items, copy)
			}
			mockState.Unlock()
			resp := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          items,
			}
			return httpmock.NewJsonResponse(200, resp)
		})

	// GET single (expand=assignments)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			policy, ok := mockState.settingsPolicies[id]
			mockState.Unlock()
			if !ok {
				return httpmock.NewJsonResponse(404, map[string]interface{}{"error": map[string]interface{}{"code": "ItemNotFound", "message": "Not found"}})
			}
			copy := make(map[string]interface{})
			for k, v := range policy {
				copy[k] = v
			}
			// If expand=assignments requested, ensure non-nil array
			if strings.Contains(req.URL.RawQuery, "$expand=assignments") {
				if _, has := copy["assignments"]; !has {
					copy["assignments"] = []interface{}{}
				}
			}
			return httpmock.NewJsonResponse(200, copy)
		})

	// POST create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			id := uuid.New().String()
			policy := map[string]interface{}{
				"@odata.type":          "#microsoft.graph.deviceManagementConfigurationPolicy",
				"id":                   id,
				"name":                 body["name"],
				"description":          body["description"],
				"platforms":            body["platforms"],
				"technologies":         body["technologies"],
				"roleScopeTagIds":      valueOrDefault(body["roleScopeTagIds"], []string{"0"}),
				"createdDateTime":      "2024-01-01T00:00:00Z",
				"lastModifiedDateTime": "2024-01-01T00:00:00Z",
				"settingCount":         1,
				"isAssigned":           false,
			}
			// store minimal assignment array
			policy["assignments"] = []interface{}{}

			mockState.Lock()
			mockState.settingsPolicies[id] = policy
			mockState.Unlock()
			return httpmock.NewJsonResponse(201, policy)
		})

	// PUT update by ID (custom request helper uses PUT)
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			existing, ok := mockState.settingsPolicies[id]
			if !ok {
				mockState.Unlock()
				return httpmock.NewJsonResponse(404, map[string]interface{}{"error": map[string]interface{}{"code": "ItemNotFound"}})
			}
			for k, v := range body {
				existing[k] = v
			}
			mockState.settingsPolicies[id] = existing
			mockState.Unlock()

			return factories.SuccessResponse(200, existing)(req)
		})

	// POST assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			var body map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			if existing, ok := mockState.settingsPolicies[id]; ok {
				// Keep raw assignments structure as provided
				if assignments, has := body["assignments"]; has {
					if arr, ok := assignments.([]interface{}); ok {
						existing["assignments"] = arr
					} else {
						existing["assignments"] = []interface{}{}
					}
				} else {
					existing["assignments"] = []interface{}{}
				}
				mockState.settingsPolicies[id] = existing
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// GET settings collection (first page)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			// Return a single basic setting collection response (SDK will PageIterator over it)
			response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_settings_collection_minimal.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// DELETE
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			_, exists := mockState.settingsPolicies[id]
			if exists {
				delete(mockState.settingsPolicies, id)
			}
			mockState.Unlock()
			if !exists {
				return httpmock.NewJsonResponse(404, map[string]interface{}{"error": map[string]interface{}{"code": "ItemNotFound"}})
			}
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears the mock state
func (m *SettingsCatalogConfigurationPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.settingsPolicies {
		delete(mockState.settingsPolicies, id)
	}
}

// RegisterErrorMocks registers mock error scenarios
func (m *SettingsCatalogConfigurationPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.settingsPolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          []map[string]interface{}{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_settings_catalog_policy_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse := map[string]interface{}{"error": map[string]interface{}{"code": "ItemNotFound"}}
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

func (m *SettingsCatalogConfigurationPolicyMock) loadJSONResponse(filePath string) (map[string]interface{}, error) {
	var response map[string]interface{}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(content, &response); err != nil {
		return response, err
	}
	return response, nil
}

func valueOrDefault(v interface{}, def interface{}) interface{} {
	if v == nil {
		return def
	}
	return v
}
