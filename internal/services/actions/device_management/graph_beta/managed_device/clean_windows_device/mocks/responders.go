package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
	mocks.GlobalRegistry.Register("clean_windows_device", &CleanWindowsDeviceMock{})
}

// CleanWindowsDeviceMock provides mock responses for clean Windows device operations
type CleanWindowsDeviceMock struct{}

// Ensure CleanWindowsDeviceMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*CleanWindowsDeviceMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for clean Windows device operations
// This implements the MockRegistrar interface
func (m *CleanWindowsDeviceMock) RegisterMocks() {
	// POST /deviceManagement/managedDevices/{id}/cleanWindowsDevice - Managed Device Action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/([^/]+)/cleanWindowsDevice$`,
		m.cleanManagedDeviceResponder())

	// POST /deviceManagement/comanagedDevices/{id}/cleanWindowsDevice - Co-managed Device Action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/([^/]+)/cleanWindowsDevice$`,
		m.cleanComanagedDeviceResponder())

	// GET /deviceManagement/managedDevices/{id} - Device validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/([^/]+)$`,
		m.getDeviceResponder())

	// GET /deviceManagement/comanagedDevices/{id} - Co-managed device validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/([^/]+)$`,
		m.getComanagedDeviceResponder())
}

// cleanManagedDeviceResponder handles POST requests to clean managed devices
func (m *CleanWindowsDeviceMock) cleanManagedDeviceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/managedDevices/")

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "error"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_clean_windows_device_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_clean_windows_device_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		default:
			// Load success response
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_clean_windows_device_success.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}

			// Store in mock state for tracking
			mockState.Lock()
			if mockState.actionResults[deviceID] == nil {
				mockState.actionResults[deviceID] = make(map[string]any)
			}
			mockState.actionResults[deviceID]["cleanManagedDevice"] = response
			mockState.Unlock()

			return factories.EmptySuccessResponse(204)(req)
		}
	}
}

// cleanComanagedDeviceResponder handles POST requests to clean co-managed devices
func (m *CleanWindowsDeviceMock) cleanComanagedDeviceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/comanagedDevices/")

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "error"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_clean_windows_device_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_clean_windows_device_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		default:
			// Load success response
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_clean_windows_device_success.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}

			// Store in mock state for tracking
			mockState.Lock()
			if mockState.actionResults[deviceID] == nil {
				mockState.actionResults[deviceID] = make(map[string]any)
			}
			mockState.actionResults[deviceID]["cleanComanagedDevice"] = response
			mockState.Unlock()

			return factories.EmptySuccessResponse(204)(req)
		}
	}
}

// getDeviceResponder handles GET requests for managed device validation
func (m *CleanWindowsDeviceMock) getDeviceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/managedDevices/")

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_clean_windows_device_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		case strings.Contains(deviceID, "12345678-1234-1234-1234-123456789abc"):
			// Windows 10 device
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_device", "get_device_windows10.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var deviceResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &deviceResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, deviceResponse)
		case strings.Contains(deviceID, "87654321-4321-4321-4321-987654321cba"):
			// Windows 11 device
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_device", "get_device_windows11.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var deviceResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &deviceResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, deviceResponse)
		default:
			// Default Windows 10 device
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_device", "get_device_windows10.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var deviceResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &deviceResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, deviceResponse)
		}
	}
}

// getComanagedDeviceResponder handles GET requests for co-managed device validation
func (m *CleanWindowsDeviceMock) getComanagedDeviceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/comanagedDevices/")

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_clean_windows_device_not_found.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(404, errorResponse)
		case strings.Contains(deviceID, "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"):
			// Co-managed device
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_device", "get_device_comanaged.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var deviceResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &deviceResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, deviceResponse)
		default:
			// Default co-managed device
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_device", "get_device_comanaged.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var deviceResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &deviceResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(200, deviceResponse)
		}
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *CleanWindowsDeviceMock) RegisterErrorMocks() {
	// POST - Managed device clean error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/error-id/cleanWindowsDevice$`,
		factories.ErrorResponse(400, "BadRequest", "Device does not support clean operation"))

	// POST - Co-managed device clean error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/error-id/cleanWindowsDevice$`,
		factories.ErrorResponse(400, "BadRequest", "Device does not support clean operation"))

	// POST - Device not found error (managed)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/not-found-id/cleanWindowsDevice$`,
		factories.ErrorResponse(404, "NotFound", "Device not found"))

	// POST - Device not found error (co-managed)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/not-found-id/cleanWindowsDevice$`,
		factories.ErrorResponse(404, "NotFound", "Device not found"))
}

// CleanupMockState clears all stored mock state
func (m *CleanWindowsDeviceMock) CleanupMockState() {
	mockState.Lock()
	mockState.actionResults = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockCleanWindowsDeviceData returns sample clean Windows device data for testing
func (m *CleanWindowsDeviceMock) GetMockCleanWindowsDeviceData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_clean_windows_device_success.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#Edm.Null",
			"value":          nil,
		}
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		panic("Failed to parse JSON response: " + err.Error())
	}
	return response
}
