package mocks

import (
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of action executions for consistent responses
var mockState struct {
	sync.Mutex
	actionResults map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.actionResults = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("bypass_activation_lock", &BypassActivationLockMock{})
}

// BypassActivationLockMock provides mock responses for bypass activation lock operations
type BypassActivationLockMock struct{}

// Ensure BypassActivationLockMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*BypassActivationLockMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for bypass activation lock operations
// This implements the MockRegistrar interface
func (m *BypassActivationLockMock) RegisterMocks() {
	// POST /deviceManagement/managedDevices/{id}/bypassActivationLock - Managed Device Action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/([^/]+)/bypassActivationLock$`,
		m.bypassActivationLockResponder())

	// GET /deviceManagement/managedDevices/{id} - Device validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/([^/]+)$`,
		m.getDeviceResponder())
}

// bypassActivationLockResponder handles POST requests to bypass activation lock on managed devices
func (m *BypassActivationLockMock) bypassActivationLockResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/managedDevices/")

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "error"):
			errorResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_bypass_activation_lock_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		case strings.Contains(deviceID, "not-found"):
			errorResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "post_bypass_activation_lock_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		default:
			// Load success response
			response, err := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_bypass_activation_lock_success.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}

			// Store in mock state for tracking
			mockState.Lock()
			if mockState.actionResults[deviceID] == nil {
				mockState.actionResults[deviceID] = make(map[string]any)
			}
			mockState.actionResults[deviceID]["bypassActivationLock"] = response
			mockState.Unlock()

			return factories.EmptySuccessResponse(204)(req)
		}
	}
}

// getDeviceResponder handles GET requests for device validation
func (m *BypassActivationLockMock) getDeviceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/managedDevices/")

		// Check for special test device IDs and load appropriate response
		switch {
		case strings.Contains(deviceID, "not-found"):
			errorResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "post_bypass_activation_lock_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		case strings.Contains(deviceID, "error"):
			// Android device - unsupported OS
			deviceResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_device", "get_device_android_unsupported.json"))
			return httpmock.NewJsonResponse(200, deviceResponse)
		case strings.Contains(deviceID, "12345678-1234-1234-1234-123456789abc"):
			// iOS device - supervised
			deviceResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_device", "get_device_ios_supervised.json"))
			return httpmock.NewJsonResponse(200, deviceResponse)
		case strings.Contains(deviceID, "87654321-4321-4321-4321-987654321cba"):
			// iPadOS device - supervised
			deviceResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_device", "get_device_ipados_supervised.json"))
			return httpmock.NewJsonResponse(200, deviceResponse)
		case strings.Contains(deviceID, "11111111-2222-3333-4444-555555555555"):
			// macOS device - DEP enrolled
			deviceResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_device", "get_device_macos_dep.json"))
			return httpmock.NewJsonResponse(200, deviceResponse)
		default:
			// Default device response - iOS supervised
			deviceResponse, _ := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_device", "get_device_ios_supervised.json"))
			return httpmock.NewJsonResponse(200, deviceResponse)
		}
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *BypassActivationLockMock) RegisterErrorMocks() {
	// POST - Device bypass error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/error-id/bypassActivationLock$`,
		factories.ErrorResponse(400, "BadRequest", "Device does not support Activation Lock bypass"))

	// POST - Device not found error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/not-found-id/bypassActivationLock$`,
		factories.ErrorResponse(404, "NotFound", "Device not found"))
}

// CleanupMockState clears all stored mock state
func (m *BypassActivationLockMock) CleanupMockState() {
	mockState.Lock()
	mockState.actionResults = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockBypassActivationLockData returns sample bypass activation lock data for testing
func (m *BypassActivationLockMock) GetMockBypassActivationLockData() map[string]any {
	response, err := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_bypass_activation_lock_success.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#Edm.Null",
			"value":          nil,
		}
	}
	return response
}
