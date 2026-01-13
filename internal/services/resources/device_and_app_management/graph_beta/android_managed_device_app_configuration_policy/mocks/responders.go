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
	androidManagedDeviceAppConfigurationPolicies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.androidManagedDeviceAppConfigurationPolicies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("android_managed_device_app_configuration_policy", &AndroidManagedDeviceAppConfigurationMock{})
}

// AndroidManagedDeviceAppConfigurationMock provides mock responses for Android managed device app configuration policy operations
type AndroidManagedDeviceAppConfigurationMock struct{}

// Ensure AndroidManagedDeviceAppConfigurationMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AndroidManagedDeviceAppConfigurationMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for Android managed device app configuration policy operations
// This implements the MockRegistrar interface
func (m *AndroidManagedDeviceAppConfigurationMock) RegisterMocks() {
	// GET /deviceAppManagement/mobileApps - Validate Android apps (with filter)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`,
		m.getAndroidManagedStoreAppsResponder())

	// POST /deviceAppManagement/mobileAppConfigurations - Create Android managed device app configuration policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations",
		m.createAndroidManagedDeviceAppConfigurationPolicyResponder())

	// GET /deviceAppManagement/mobileAppConfigurations/{id} - Read Android managed device app configuration policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.getAndroidManagedDeviceAppConfigurationPolicyResponder())

	// PATCH /deviceAppManagement/mobileAppConfigurations/{id} - Update Android managed device app configuration policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.updateAndroidManagedDeviceAppConfigurationPolicyResponder())

	// DELETE /deviceAppManagement/mobileAppConfigurations/{id} - Delete Android managed device app configuration policy
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		m.deleteAndroidManagedDeviceAppConfigurationPolicyResponder())
}

// RegisterErrorMocks sets up mock HTTP responders that return error responses
func (m *AndroidManagedDeviceAppConfigurationMock) RegisterErrorMocks() {
	// POST /deviceAppManagement/mobileAppConfigurations - Create Error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppConfigurations",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_android_managed_store_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// GET /deviceAppManagement/mobileAppConfigurations/{id} - Read Error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppConfigurations/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_android_managed_store_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *AndroidManagedDeviceAppConfigurationMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Android managed device app configuration policies
	for id := range mockState.androidManagedDeviceAppConfigurationPolicies {
		delete(mockState.androidManagedDeviceAppConfigurationPolicies, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *AndroidManagedDeviceAppConfigurationMock) loadJSONResponse(filePath string) (map[string]any, error) {
	var response map[string]any

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// createAndroidManagedDeviceAppConfigurationPolicyResponder handles POST requests to create Android managed device app configuration policies
func (m *AndroidManagedDeviceAppConfigurationMock) createAndroidManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_android_managed_store_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Generate a new UUID for the created resource
		id := uuid.New().String()

		// Create response with request data
		response := map[string]any{
			"@odata.type":          "#microsoft.graph.androidManagedStoreAppConfiguration",
			"id":                   id,
			"version":              1,
			"appSupportsOemConfig": false,
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
		if packageId, ok := requestBody["packageId"]; ok {
			response["packageId"] = packageId
		}
		if payloadJson, ok := requestBody["payloadJson"]; ok {
			response["payloadJson"] = payloadJson
		}
		if profileApplicability, ok := requestBody["profileApplicability"]; ok {
			response["profileApplicability"] = profileApplicability
		}
		if connectedAppsEnabled, ok := requestBody["connectedAppsEnabled"]; ok {
			response["connectedAppsEnabled"] = connectedAppsEnabled
		}
		if permissionActions, ok := requestBody["permissionActions"]; ok {
			response["permissionActions"] = permissionActions
		}

		// Store in mock state
		mockState.Lock()
		mockState.androidManagedDeviceAppConfigurationPolicies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getAndroidManagedDeviceAppConfigurationPolicyResponder handles GET requests to retrieve Android managed device app configuration policies
func (m *AndroidManagedDeviceAppConfigurationMock) getAndroidManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		androidConfig, exists := mockState.androidManagedDeviceAppConfigurationPolicies[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs for different app configurations
			var responseFile string
			switch {
			case strings.Contains(id, "minimal"):
				responseFile = "get_android_managed_store_app_configuration_minimal.json"
			case strings.Contains(id, "authenticator"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_authenticator_maximal.json"
			case strings.Contains(id, "copilot"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_365_copilot_maximal.json"
			case strings.Contains(id, "home_screen"):
				responseFile = "get_android_managed_store_app_configuration_managed_home_screen_maximal.json"
			case strings.Contains(id, "defender"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_defender_antivirus_maximal.json"
			case strings.Contains(id, "edge"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_edge_browser_maximal.json"
			case strings.Contains(id, "excel"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_excel_maximal.json"
			case strings.Contains(id, "onedrive"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_onedrive_maximal.json"
			case strings.Contains(id, "onenote"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_onenote_maximal.json"
			case strings.Contains(id, "outlook"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_outlook_maximal.json"
			case strings.Contains(id, "powerpoint"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_powerpoint_maximal.json"
			case strings.Contains(id, "teams"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_teams_maximal.json"
			case strings.Contains(id, "word"):
				responseFile = "get_android_managed_store_app_configuration_microsoft_word_maximal.json"
			default:
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_android_managed_store_app_configuration_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			if responseFile != "" {
				response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", responseFile))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			}
		}

		return factories.SuccessResponse(200, androidConfig)(req)
	}
}

// updateAndroidManagedDeviceAppConfigurationPolicyResponder handles PATCH requests to update Android managed device app configuration policies
func (m *AndroidManagedDeviceAppConfigurationMock) updateAndroidManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		androidConfig, exists := mockState.androidManagedDeviceAppConfigurationPolicies[id]
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_android_managed_store_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_android_managed_store_app_configuration_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Load update template
		updatedConfig, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "get_android_managed_store_app_configuration_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range androidConfig {
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
		mockState.androidManagedDeviceAppConfigurationPolicies[id] = updatedConfig
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedConfig)(req)
	}
}

// deleteAndroidManagedDeviceAppConfigurationPolicyResponder handles DELETE requests to delete Android managed device app configuration policies
func (m *AndroidManagedDeviceAppConfigurationMock) deleteAndroidManagedDeviceAppConfigurationPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceAppManagement/mobileAppConfigurations/")

		mockState.Lock()
		_, exists := mockState.androidManagedDeviceAppConfigurationPolicies[id]
		if exists {
			delete(mockState.androidManagedDeviceAppConfigurationPolicies, id)
		}
		mockState.Unlock()

		if !exists {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_android_managed_store_app_configuration_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		// Return 204 No Content for successful deletion
		return httpmock.NewStringResponse(204, ""), nil
	}
}

// getAndroidManagedStoreAppsResponder handles GET requests to retrieve Android managed store apps for validation
func (m *AndroidManagedDeviceAppConfigurationMock) getAndroidManagedStoreAppsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Load the mock response from JSON file
		response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_apps", "get_android_managed_store_apps.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		return factories.SuccessResponse(200, response)(req)
	}
}
