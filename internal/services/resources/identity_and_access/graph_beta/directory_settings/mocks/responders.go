package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// MockState tracks the state of the directory settings for consistent responses
// Supports all 8 directory setting templates
type MockState struct {
	sync.Mutex
	SettingsID     string
	TemplateID     string
	SettingsExists bool

	// Store values as a map for flexibility across all templates
	Values map[string]string
}

// DirectorySettingsMock provides mock responses for directory settings operations
type DirectorySettingsMock struct {
	MockState *MockState
}

// Ensure DirectorySettingsMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*DirectorySettingsMock)(nil)

// Template ID to response file mapping
var templateResponseFiles = map[string]struct {
	post string
	get  string
}{
	"08d542b9-071f-4e16-94b0-74abb372e3d9": { // Group.Unified.Guest
		post: "../tests/responses/validate_create/post_group_unified_guest_response.json",
		get:  "../tests/responses/validate_create/get_group_unified_guest_response.json",
	},
	"4bc7f740-180e-4586-adb6-38b2e9024e6b": { // Application
		post: "../tests/responses/validate_create/post_application_response.json",
		get:  "../tests/responses/validate_create/get_application_response.json",
	},
	"5cf42378-d67d-4f36-ba46-e8b86229381d": { // Password Rule Settings
		post: "../tests/responses/validate_create/post_password_rule_settings_response.json",
		get:  "../tests/responses/validate_create/get_password_rule_settings_response.json",
	},
	"62375ab9-6b52-47ed-826b-58e47e0e304b": { // Group.Unified
		post: "../tests/responses/validate_create/post_group_unified_response.json",
		get:  "../tests/responses/validate_create/get_group_unified_response.json",
	},
	"80661d51-be2f-4d46-9713-98a2fcaec5bc": { // Prohibited Names Settings
		post: "../tests/responses/validate_create/post_prohibited_names_settings_response.json",
		get:  "../tests/responses/validate_create/get_prohibited_names_settings_response.json",
	},
	"898f1161-d651-43d1-805c-3b0b388a9fc2": { // Custom Policy Settings
		post: "../tests/responses/validate_create/post_custom_policy_settings_response.json",
		get:  "../tests/responses/validate_create/get_custom_policy_settings_response.json",
	},
	"aad3907d-1d1a-448b-b3ef-7bf7f63db63b": { // Prohibited Names Restricted Settings
		post: "../tests/responses/validate_create/post_prohibited_names_restricted_settings_response.json",
		get:  "../tests/responses/validate_create/get_prohibited_names_restricted_settings_response.json",
	},
	"dffd5d46-495d-40a9-8e21-954ff55e198a": { // Consent Policy Settings
		post: "../tests/responses/validate_create/post_consent_policy_settings_response.json",
		get:  "../tests/responses/validate_create/get_consent_policy_settings_response.json",
	},
}

func init() {
	// Register with global registry
	mocks.GlobalRegistry.Register("directory_settings", &DirectorySettingsMock{})
}

// RegisterMocks registers HTTP mock responses for directory settings operations
func (m *DirectorySettingsMock) RegisterMocks() {
	// Initialize mock state
	m.MockState = &MockState{
		SettingsID:     uuid.New().String(),
		TemplateID:     "",
		SettingsExists: false,
		Values:         make(map[string]string),
	}

	// Register GET for directory settings (list - to check if settings exist)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/settings$`,
		func(req *http.Request) (*http.Response, error) {
			m.MockState.Lock()
			defer m.MockState.Unlock()

			// If settings don't exist, return empty list
			if !m.MockState.SettingsExists {
				resp := httpmock.NewStringResponse(200, `{"@odata.context":"https://graph.microsoft.com/beta/$metadata#settings","value":[]}`)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			// Load the appropriate JSON response based on template ID
			files, ok := templateResponseFiles[m.MockState.TemplateID]
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"message":"Unknown template ID"}}`), nil
			}

			jsonStr, _ := helpers.ParseJSONFile(files.get)
			var settingsObject map[string]any
			json.Unmarshal([]byte(jsonStr), &settingsObject)

			// Update with current mock state
			settingsObject["id"] = m.MockState.SettingsID
			m.updateSettingsValues(settingsObject)

			// Wrap in list response
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#settings",
				"value":          []any{settingsObject},
			}

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register GET for individual settings object
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/settings/([a-fA-F0-9\-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			settingsID := httpmock.MustGetSubmatch(req, 1)

			m.MockState.Lock()
			defer m.MockState.Unlock()

			// Check if settings ID matches
			if settingsID != m.MockState.SettingsID || !m.MockState.SettingsExists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			// Load the appropriate JSON response based on template ID
			files, ok := templateResponseFiles[m.MockState.TemplateID]
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"message":"Unknown template ID"}}`), nil
			}

			jsonStr, _ := helpers.ParseJSONFile(files.get)
			var response map[string]any
			json.Unmarshal([]byte(jsonStr), &response)

			// Update with current mock state
			response["id"] = m.MockState.SettingsID
			m.updateSettingsValues(response)

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register POST for creating settings
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/settings$`,
		func(req *http.Request) (*http.Response, error) {
			m.MockState.Lock()
			defer m.MockState.Unlock()

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			// Extract template ID from request
			if templateID, ok := requestBody["templateId"].(string); ok {
				m.MockState.TemplateID = templateID
			}

			// Extract values from request
			m.extractSettingsFromRequest(requestBody)
			m.MockState.SettingsExists = true

			// Load the appropriate JSON response based on template ID
			files, ok := templateResponseFiles[m.MockState.TemplateID]
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"message":"Unknown template ID"}}`), nil
			}

			jsonStr, _ := helpers.ParseJSONFile(files.post)
			var response map[string]any
			json.Unmarshal([]byte(jsonStr), &response)

			// Update with current mock state
			response["id"] = m.MockState.SettingsID
			m.updateSettingsValues(response)

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(201, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register PATCH for updating settings
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/settings/([a-fA-F0-9\-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			settingsID := httpmock.MustGetSubmatch(req, 1)

			m.MockState.Lock()
			defer m.MockState.Unlock()

			// Check if settings ID matches
			if settingsID != m.MockState.SettingsID {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			// Update values from request
			m.extractSettingsFromRequest(requestBody)
			m.MockState.SettingsExists = true

			// Load the appropriate JSON response based on template ID
			files, ok := templateResponseFiles[m.MockState.TemplateID]
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"message":"Unknown template ID"}}`), nil
			}

			jsonStr, _ := helpers.ParseJSONFile(files.get)
			var response map[string]any
			json.Unmarshal([]byte(jsonStr), &response)

			// Update with current mock state
			response["id"] = m.MockState.SettingsID
			m.updateSettingsValues(response)

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register DELETE for deleting settings
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/settings/([a-fA-F0-9\-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			settingsID := httpmock.MustGetSubmatch(req, 1)

			m.MockState.Lock()
			defer m.MockState.Unlock()

			// Check if settings ID matches
			if settingsID != m.MockState.SettingsID {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			m.MockState.SettingsExists = false

			resp := httpmock.NewStringResponse(204, "")
			return resp, nil
		})
}

// updateSettingsValues updates the settings values in the response with current mock state
func (m *DirectorySettingsMock) updateSettingsValues(settings map[string]any) {
	if values, ok := settings["values"].([]any); ok {
		for _, v := range values {
			if settingValue, ok := v.(map[string]any); ok {
				if name, ok := settingValue["name"].(string); ok {
					// Update the value if we have it in our mock state
					if mockValue, exists := m.MockState.Values[name]; exists {
						settingValue["value"] = mockValue
					}
				}
			}
		}
	}
}

// extractSettingsFromRequest extracts settings values from the request body
func (m *DirectorySettingsMock) extractSettingsFromRequest(requestBody map[string]any) {
	if values, ok := requestBody["values"].([]any); ok {
		for _, v := range values {
			if settingValue, ok := v.(map[string]any); ok {
				if name, ok := settingValue["name"].(string); ok {
					if value, ok := settingValue["value"].(string); ok {
						m.MockState.Values[name] = value
					}
				}
			}
		}
	}
}

// RegisterErrorMocks registers HTTP mock responses that simulate error conditions
func (m *DirectorySettingsMock) RegisterErrorMocks() {
	// Initialize mock state
	m.MockState = &MockState{
		SettingsID:     uuid.New().String(),
		TemplateID:     "",
		SettingsExists: false,
		Values:         make(map[string]string),
	}

	// Register GET for directory settings (returns empty - no existing settings)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/settings$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `{"@odata.context":"https://graph.microsoft.com/beta/$metadata#settings","value":[]}`)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register POST for creating settings - return error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/settings$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
			return httpmock.NewStringResponse(400, jsonStr), nil
		})

	// Register PATCH for updating settings - return error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/settings/([a-fA-F0-9\-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
			return httpmock.NewStringResponse(400, jsonStr), nil
		})
}

// CleanupMockState resets the mock state
func (m *DirectorySettingsMock) CleanupMockState() {
	if m.MockState != nil {
		m.MockState.Lock()
		defer m.MockState.Unlock()
		m.MockState.SettingsExists = false
		m.MockState.TemplateID = ""
		m.MockState.Values = make(map[string]string)
	}
}
