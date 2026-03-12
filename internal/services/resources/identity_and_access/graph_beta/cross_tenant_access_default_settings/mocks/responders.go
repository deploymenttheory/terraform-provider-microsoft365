package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	defaultSettings map[string]any
	isPatched       bool
}

func init() {
	mockState.defaultSettings = make(map[string]any)
	mockState.isPatched = false
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("cross_tenant_access_default_settings", &CrossTenantAccessDefaultSettingsMock{})
}

type CrossTenantAccessDefaultSettingsMock struct{}

var _ mocks.MockRegistrar = (*CrossTenantAccessDefaultSettingsMock)(nil)

// loadJSONResponse loads a JSON response file from the tests/responses directory
func loadJSONResponse(filename string) (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_create", filename)
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file %s: %w", filename, err)
	}

	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return responseObj, nil
}

// mergeSettings deeply merges source into target
func mergeSettings(target, source map[string]any) {
	for key, value := range source {
		if sourceMap, ok := value.(map[string]any); ok {
			if targetMap, ok := target[key].(map[string]any); ok {
				mergeSettings(targetMap, sourceMap)
			} else {
				target[key] = value
			}
		} else {
			target[key] = value
		}
	}
}

func (m *CrossTenantAccessDefaultSettingsMock) RegisterMocks() {
	mockState.Lock()
	mockState.defaultSettings = make(map[string]any)
	mockState.isPatched = false
	mockState.Unlock()

	// Register mock dependencies for validation
	m.registerMockUsers()
	m.registerMockGroups()

	// GET /policies/crossTenantAccessPolicy/default - Read default settings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// If we have stored settings (either patched or reset), return them
		if len(mockState.defaultSettings) > 0 {
			resp, err := httpmock.NewJsonResponse(200, mockState.defaultSettings)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
			}
			return resp, nil
		}

		// Otherwise, load and return system defaults from JSON file
		defaultResponse, err := loadJSONResponse("get_cross_tenant_access_default_settings_system_default.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
		}

		// Store the system defaults in mock state for subsequent calls
		mockState.defaultSettings = defaultResponse
		mockState.isPatched = false

		resp, err := httpmock.NewJsonResponse(200, defaultResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})

	// PATCH /policies/crossTenantAccessPolicy/default - Update default settings
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		// Load the base response from JSON file if we don't have state yet
		if len(mockState.defaultSettings) == 0 {
			baseResponse, err := loadJSONResponse("patch_cross_tenant_access_default_settings_success.json")
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
			}
			mockState.defaultSettings = baseResponse
		}

		// Merge request body into current state
		mergeSettings(mockState.defaultSettings, requestBody)
		mockState.defaultSettings["isServiceDefault"] = false
		mockState.isPatched = true

		// PATCH returns 204 No Content
		return httpmock.NewStringResponse(204, ""), nil
	})

	// POST /policies/crossTenantAccessPolicy/default/resetToSystemDefault - Reset to system defaults
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default/resetToSystemDefault", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Load system defaults from JSON file
		systemDefaults, err := loadJSONResponse("get_cross_tenant_access_default_settings_system_default.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
		}

		// Reset mock state to system defaults
		mockState.defaultSettings = systemDefaults
		mockState.isPatched = false

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *CrossTenantAccessDefaultSettingsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.defaultSettings = make(map[string]any)
	mockState.isPatched = false
	mockState.Unlock()

	// GET - Return error
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default",
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))

	// PATCH - Return error
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default",
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))

	// POST - Return error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default/resetToSystemDefault",
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))
}

func (m *CrossTenantAccessDefaultSettingsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.defaultSettings = make(map[string]any)
	mockState.isPatched = false
	httpmock.Reset()
}

// registerMockUsers registers mock user resources for validation
func (m *CrossTenantAccessDefaultSettingsMock) registerMockUsers() {
	// GET /users - List users with filter
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users\?`, func(req *http.Request) (*http.Response, error) {
		// Return a collection response with the requested users
		// For unit tests, we'll return mock users that match the GUIDs in test 09
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#users(id)",
			"value": []any{
				map[string]any{
					"id": "11111111-1111-1111-1111-111111111111",
				},
				map[string]any{
					"id": "22222222-2222-2222-2222-222222222222",
				},
			},
		}

		resp, err := httpmock.NewJsonResponse(200, response)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})
}

// registerMockGroups registers mock group resources for validation
func (m *CrossTenantAccessDefaultSettingsMock) registerMockGroups() {
	// GET /groups - List groups with filter
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups\?`, func(req *http.Request) (*http.Response, error) {
		// Return a collection response with the requested groups
		// For unit tests, we'll return mock groups that match the GUIDs in test 09
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups(id)",
			"value": []any{
				map[string]any{
					"id": "11111111-1111-1111-1111-111111111111",
				},
				map[string]any{
					"id": "22222222-2222-2222-2222-222222222222",
				},
			},
		}

		resp, err := httpmock.NewJsonResponse(200, response)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})
}
