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
	organizationSettings map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.organizationSettings = make(map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// CloudPcOrganizationSettingsMock provides mock responses for Cloud PC organization settings operations
type CloudPcOrganizationSettingsMock struct{}

// RegisterMocks registers HTTP mock responses for Cloud PC organization settings operations
func (m *CloudPcOrganizationSettingsMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.organizationSettings = make(map[string]interface{})
	mockState.Unlock()

	// Register GET for Cloud PC organization settings (singleton)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			settings := mockState.organizationSettings
			mockState.Unlock()

			if len(settings) == 0 {
				// Return default settings if none exist
				defaultSettings := map[string]interface{}{
					"@odata.type":          "#microsoft.graph.cloudPcOrganizationSettings",
					"id":                   "default",
					"enableMEMAutoEnroll":  false,
					"enableSingleSignOn":   false,
					"osVersion":            "windows10",
					"userAccountType":      "standardUser",
					"windowsSettings": map[string]interface{}{
						"language": "en-US",
					},
				}
				return httpmock.NewJsonResponse(200, defaultSettings)
			}

			return httpmock.NewJsonResponse(200, settings)
		})

	// Register PATCH for updating Cloud PC organization settings
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Create or update the organization settings
			settings := map[string]interface{}{
				"@odata.type": "#microsoft.graph.cloudPcOrganizationSettings",
				"id":          "default",
			}

			// Update the settings with new values
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(settings, key)
				} else {
					settings[key] = value
				}
			}

			// Store in mock state
			mockState.Lock()
			mockState.organizationSettings = settings
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, settings)
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *CloudPcOrganizationSettingsMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.organizationSettings = make(map[string]interface{})
	mockState.Unlock()

	// Register error response for getting Cloud PC organization settings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		factories.ErrorResponse(500, "InternalServerError", "Internal server error occurred"))

	// Register error response for updating Cloud PC organization settings with invalid data
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/organizationSettings",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid OS version"))
}