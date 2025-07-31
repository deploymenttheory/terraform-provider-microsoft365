package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	organizationSettings map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.organizationSettings = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// OrganizationSettingsMock provides mock responses for organization settings operations
type OrganizationSettingsMock struct{}

// RegisterMocks registers HTTP mock responses for organization settings operations
func (m *OrganizationSettingsMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.organizationSettings = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for organization settings (singleton)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			// Organization settings is a singleton, so we use a fixed ID
			if orgSettings, exists := mockState.organizationSettings["singleton"]; exists {
				return httpmock.NewJsonResponse(200, orgSettings)
			}

			// If not exists, return default organization settings
			defaultSettings := map[string]interface{}{
				"id":                    "singleton",
				"enableMEMAutoEnroll":   false,
				"enableSingleSignOn":    false,
				"osVersion":             nil,
				"userAccountType":       nil,
				"windowsSettings":       nil,
			}
			mockState.organizationSettings["singleton"] = defaultSettings
			return httpmock.NewJsonResponse(200, defaultSettings)
		})

	// Register PATCH for updating organization settings
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Get existing settings or create default
			orgSettings, exists := mockState.organizationSettings["singleton"]
			if !exists {
				orgSettings = map[string]interface{}{
					"id":                    "singleton",
					"enableMEMAutoEnroll":   false,
					"enableSingleSignOn":    false,
					"osVersion":             nil,
					"userAccountType":       nil,
					"windowsSettings":       nil,
				}
			}

			// Update organization settings data
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(orgSettings, key)
				} else {
					orgSettings[key] = value
				}
			}

			// Ensure the ID is preserved
			orgSettings["id"] = "singleton"
			mockState.organizationSettings["singleton"] = orgSettings

			return httpmock.NewJsonResponse(200, orgSettings)
		})

	// Register specific organization settings mocks for testing
	registerSpecificOrganizationSettingsMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *OrganizationSettingsMock) RegisterErrorMocks() {
	// Register error response for updating organization settings with invalid data
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid settings"))
}

// registerSpecificOrganizationSettingsMocks registers mocks for specific test scenarios
func registerSpecificOrganizationSettingsMocks() {
	// Minimal organization settings
	minimalSettingsData := map[string]interface{}{
		"id":                  "singleton",
		"enableMEMAutoEnroll": false,
		"enableSingleSignOn":  false,
	}

	// Maximal organization settings
	maximalSettingsData := map[string]interface{}{
		"id":                  "singleton",
		"enableMEMAutoEnroll": true,
		"enableSingleSignOn":  true,
		"osVersion":           "windows11",
		"userAccountType":     "administrator",
		"windowsSettings": map[string]interface{}{
			"language": "en-GB",
		},
	}

	mockState.Lock()
	// For testing purposes, we'll store both scenarios but the actual test will determine which one to use
	mockState.organizationSettings["singleton"] = minimalSettingsData
	mockState.Unlock()

	// Store maximal settings for specific tests that need it
	mockState.Lock()
	mockState.organizationSettings["maximal"] = maximalSettingsData
	mockState.Unlock()
}