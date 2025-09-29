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
	iosMobileAppConfigurations map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.iosMobileAppConfigurations = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("ios_mobile_app_configuration", &IOSMobileAppConfigurationMock{})
}

// IOSMobileAppConfigurationMock provides mock responses for iOS mobile app configuration operations
type IOSMobileAppConfigurationMock struct{}

// Ensure IOSMobileAppConfigurationMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*IOSMobileAppConfigurationMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for iOS mobile app configuration operations
// This implements the MockRegistrar interface
func (m *IOSMobileAppConfigurationMock) RegisterMocks() {
	// POST /deviceAppManagement/mobileAppConfigurations - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations",
		m.createIOSMobileAppConfigurationResponder())

	// GET /deviceAppManagement/mobileAppConfigurations/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.getIOSMobileAppConfigurationResponder())

	// PATCH /deviceAppManagement/mobileAppConfigurations/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.updateIOSMobileAppConfigurationResponder())

	// DELETE /deviceAppManagement/mobileAppConfigurations/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.deleteIOSMobileAppConfigurationResponder())
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

	// Clear all stored iOS mobile app configurations
	for id := range mockState.iosMobileAppConfigurations {
		delete(mockState.iosMobileAppConfigurations, id)
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

// createIOSMobileAppConfigurationResponder handles POST requests to create iOS mobile app configurations
func (m *IOSMobileAppConfigurationMock) createIOSMobileAppConfigurationResponder() httpmock.Responder {
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
		mockState.iosMobileAppConfigurations[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getIOSMobileAppConfigurationResponder handles GET requests to retrieve iOS mobile app configurations
func (m *IOSMobileAppConfigurationMock) getIOSMobileAppConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		iosConfig, exists := mockState.iosMobileAppConfigurations[id]
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

// updateIOSMobileAppConfigurationResponder handles PATCH requests to update iOS mobile app configurations
func (m *IOSMobileAppConfigurationMock) updateIOSMobileAppConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		iosConfig, exists := mockState.iosMobileAppConfigurations[id]
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
		mockState.iosMobileAppConfigurations[id] = updatedConfig
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedConfig)(req)
	}
}

// deleteIOSMobileAppConfigurationResponder handles DELETE requests to delete iOS mobile app configurations
func (m *IOSMobileAppConfigurationMock) deleteIOSMobileAppConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		_, exists := mockState.iosMobileAppConfigurations[id]
		if exists {
			delete(mockState.iosMobileAppConfigurations, id)
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
