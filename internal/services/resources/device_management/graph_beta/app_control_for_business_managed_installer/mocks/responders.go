package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of the Windows Management App for consistent responses
var mockState struct {
	sync.Mutex
	windowsManagementApp map[string]interface{}
}

func init() {
	// Initialize mockState with default Windows Management App data
	mockState.windowsManagementApp = map[string]interface{}{
		"@odata.context":                     "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
		"id":                                 "54fac284-7866-43e5-860a-9c8e10fa3d7d",
		"availableVersion":                   "1.93.102.0",
		"managedInstaller":                   "disabled",
		"managedInstallerConfiguredDateTime": nil,
	}

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("app_control_for_business_managed_installer", &AppControlForBusinessManagedInstallerMock{})
}

// AppControlForBusinessManagedInstallerMock provides mock responses for managed installer operations
type AppControlForBusinessManagedInstallerMock struct{}

// Ensure AppControlForBusinessManagedInstallerMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AppControlForBusinessManagedInstallerMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for managed installer operations
// This implements the MockRegistrar interface
func (m *AppControlForBusinessManagedInstallerMock) RegisterMocks() {
	// GET /deviceAppManagement/windowsManagementApp - Read
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp",
		m.getWindowsManagementAppResponder())

	// POST /deviceAppManagement/windowsManagementApp/setAsManagedInstaller - Toggle
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller",
		m.setAsManagedInstallerResponder())
}

// getWindowsManagementAppResponder handles GET requests to retrieve the Windows Management App
func (m *AppControlForBusinessManagedInstallerMock) getWindowsManagementAppResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		currentStatus := mockState.windowsManagementApp["managedInstaller"].(string)
		mockState.Unlock()

		var response map[string]any
		var jsonStr string
		var err error

		// Load appropriate JSON response based on current state
		if currentStatus == "enabled" {
			jsonStr, _ = helpers.ParseJSONFile("../tests/responses/validate_create/get_windows_management_app_enabled.json")
		} else {
			jsonStr, _ = helpers.ParseJSONFile("../tests/responses/validate_create/get_windows_management_app_disabled.json")
		}

		if err = json.Unmarshal([]byte(jsonStr), &response); err != nil {
			// Fallback to mock state if JSON loading fails
			response = make(map[string]any)
			mockState.Lock()
			for k, v := range mockState.windowsManagementApp {
				response[k] = v
			}
			mockState.Unlock()
		}

		return httpmock.NewJsonResponse(200, response)
	}
}

// setAsManagedInstallerResponder handles POST requests to toggle the managed installer status
func (m *AppControlForBusinessManagedInstallerMock) setAsManagedInstallerResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		// Toggle the managed installer status
		currentStatus, ok := mockState.windowsManagementApp["managedInstaller"].(string)
		if !ok {
			currentStatus = "disabled"
		}

		if currentStatus == "disabled" {
			mockState.windowsManagementApp["managedInstaller"] = "enabled"
			mockState.windowsManagementApp["managedInstallerConfiguredDateTime"] = "2024-01-01T12:00:00Z"
		} else {
			mockState.windowsManagementApp["managedInstaller"] = "disabled"
			mockState.windowsManagementApp["managedInstallerConfiguredDateTime"] = nil
		}

		// Return 204 No Content for successful POST
		return httpmock.NewStringResponse(204, ""), nil
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *AppControlForBusinessManagedInstallerMock) RegisterErrorMocks() {
	// GET - Read error using JSON response file
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_management_app_not_found.json")
			var errorObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorObj)
			return httpmock.NewJsonResponse(404, errorObj)
		})

	// POST - setAsManagedInstaller error using JSON response file
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagementApp/setAsManagedInstaller",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_set_managed_installer_error.json")
			var errorObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorObj)
			return httpmock.NewJsonResponse(500, errorObj)
		})
}

// CleanupMockState resets the mock state to default values
func (m *AppControlForBusinessManagedInstallerMock) CleanupMockState() {
	mockState.Lock()
	mockState.windowsManagementApp = map[string]interface{}{
		"@odata.context":                     "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
		"id":                                 "54fac284-7866-43e5-860a-9c8e10fa3d7d",
		"availableVersion":                   "1.93.102.0",
		"managedInstaller":                   "disabled",
		"managedInstallerConfiguredDateTime": nil,
	}
	mockState.Unlock()
}

// GetMockWindowsManagementAppDisabledData returns sample Windows Management App data with managed installer disabled
func (m *AppControlForBusinessManagedInstallerMock) GetMockWindowsManagementAppDisabledData() map[string]interface{} {
	return map[string]interface{}{
		"@odata.context":                     "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
		"id":                                 "54fac284-7866-43e5-860a-9c8e10fa3d7d",
		"availableVersion":                   "1.93.102.0",
		"managedInstaller":                   "disabled",
		"managedInstallerConfiguredDateTime": nil,
	}
}

// GetMockWindowsManagementAppEnabledData returns sample Windows Management App data with managed installer enabled
func (m *AppControlForBusinessManagedInstallerMock) GetMockWindowsManagementAppEnabledData() map[string]interface{} {
	return map[string]interface{}{
		"@odata.context":                     "https://graph.microsoft.com/beta/$metadata#microsoft.graph.windowsManagementApp",
		"id":                                 "54fac284-7866-43e5-860a-9c8e10fa3d7d",
		"availableVersion":                   "1.93.102.0",
		"managedInstaller":                   "enabled",
		"managedInstallerConfiguredDateTime": "2024-01-01T12:00:00Z",
	}
}
