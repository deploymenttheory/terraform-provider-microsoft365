package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	iosManagedDeviceAppConfigurationPolicies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.iosManagedDeviceAppConfigurationPolicies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("ios_managed_device_app_configuration_policy", &IOSMobileAppConfigurationMock{})
}

// IOSMobileAppConfigurationMock provides mock responses for iOS managed device app configuration policy operations
type IOSMobileAppConfigurationMock struct{}

// Ensure IOSMobileAppConfigurationMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*IOSMobileAppConfigurationMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for iOS managed device app configuration policy operations
// This implements the MockRegistrar interface
func (m *IOSMobileAppConfigurationMock) RegisterMocks() {
	// GET /deviceAppManagement/mobileApps - Validate iOS apps (with filter)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`,
		m.getIOSMobileAppsResponder())

	// POST /deviceAppManagement/mobileAppConfigurations - Create iOS managed device app configuration policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations",
		m.createIOSManagedDeviceAppConfigurationPolicyResponder())

	// GET /deviceAppManagement/mobileAppConfigurations/{id} - Read iOS managed device app configuration policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.getIOSManagedDeviceAppConfigurationPolicyResponder())

	// PATCH /deviceAppManagement/mobileAppConfigurations/{id} - Update iOS managed device app configuration policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.updateIOSManagedDeviceAppConfigurationPolicyResponder())

	// DELETE /deviceAppManagement/mobileAppConfigurations/{id} - Delete iOS managed device app configuration policy
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.deleteIOSManagedDeviceAppConfigurationPolicyResponder())
}

// RegisterErrorMocks sets up mock HTTP responders that return error responses
func (m *IOSMobileAppConfigurationMock) RegisterErrorMocks() {
	// POST /deviceAppManagement/mobileAppConfigurations - Create Error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_ios_mobile_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// GET /deviceAppManagement/mobileAppConfigurations/{id} - Read Error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_ios_mobile_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *IOSMobileAppConfigurationMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored iOS managed device app configuration policies
	for id := range mockState.iosManagedDeviceAppConfigurationPolicies {
		delete(mockState.iosManagedDeviceAppConfigurationPolicies, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *IOSMobileAppConfigurationMock) loadJSONResponse(filePath string) (map[string]any, error) {
	var response map[string]any

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// createIOSManagedDeviceAppConfigurationPolicyResponder handles POST requests to create iOS managed device app configuration policies
func (m *IOSMobileAppConfigurationMock) createIOSManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_ios_mobile_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Generate a new UUID for the created resource
		id := uuid.New().String()

		// Create response with request data
		response := map[string]any{
			"@odata.type": "#microsoft.graph.iosMobileAppConfiguration",
			"id":          id,
			"version":     1,
		}

		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		}
		if targetedMobileApps, ok := requestBody["targetedMobileApps"]; ok {
			response["targetedMobileApps"] = targetedMobileApps
		}
		if roleScopeTagIds, ok := requestBody["roleScopeTagIds"]; ok {
			response["roleScopeTagIds"] = roleScopeTagIds
		} else {
			response["roleScopeTagIds"] = []string{"0"}
		}
		if encodedSettingXml, ok := requestBody["encodedSettingXml"]; ok {
			response["encodedSettingXml"] = encodedSettingXml
		}
		if settings, ok := requestBody["settings"]; ok {
			response["settings"] = settings
		}

		// Store in mock state
		mockState.Lock()
		mockState.iosManagedDeviceAppConfigurationPolicies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getIOSManagedDeviceAppConfigurationPolicyResponder handles GET requests to retrieve iOS managed device app configuration policies
func (m *IOSMobileAppConfigurationMock) getIOSManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		iosConfig, exists := mockState.iosManagedDeviceAppConfigurationPolicies[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_ios_mobile_app_configuration_minimal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			case strings.Contains(id, "maximal"):
				response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_ios_mobile_app_configuration_maximal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			default:
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_ios_mobile_app_configuration_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}
		}

		return factories.SuccessResponse(200, iosConfig)(req)
	}
}

// updateIOSManagedDeviceAppConfigurationPolicyResponder handles PATCH requests to update iOS managed device app configuration policies
func (m *IOSMobileAppConfigurationMock) updateIOSManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		iosConfig, exists := mockState.iosManagedDeviceAppConfigurationPolicies[id]
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_ios_mobile_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_ios_mobile_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Load update template
		updatedConfig, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "get_ios_mobile_app_configuration_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range iosConfig {
			updatedConfig[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedConfig[k] = v
		}

		// Increment version
		if version, ok := updatedConfig["version"].(float64); ok {
			updatedConfig["version"] = int(version) + 1
		}

		// Store updated state
		mockState.Lock()
		mockState.iosManagedDeviceAppConfigurationPolicies[id] = updatedConfig
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedConfig)(req)
	}
}

// deleteIOSManagedDeviceAppConfigurationPolicyResponder handles DELETE requests to delete iOS managed device app configuration policies
func (m *IOSMobileAppConfigurationMock) deleteIOSManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		_, exists := mockState.iosManagedDeviceAppConfigurationPolicies[id]
		if exists {
			delete(mockState.iosManagedDeviceAppConfigurationPolicies, id)
		}
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_ios_mobile_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		// Return 204 No Content for successful deletion
		return httpmock.NewStringResponse(204, ""), nil
	}
}

// getIOSMobileAppsResponder handles GET requests to retrieve iOS mobile apps for validation
func (m *IOSMobileAppConfigurationMock) getIOSMobileAppsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Load the mock response from JSON file
		response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_apps", "get_ios_mobile_apps.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		return factories.SuccessResponse(200, response)(req)
	}
}
