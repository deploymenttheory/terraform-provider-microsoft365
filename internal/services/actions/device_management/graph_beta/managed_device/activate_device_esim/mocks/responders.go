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
	mocks.GlobalRegistry.Register("activate_device_esim", &ActivateDeviceEsimMock{})
}

// ActivateDeviceEsimMock provides mock responses for activate device esim operations
type ActivateDeviceEsimMock struct{}

// Ensure ActivateDeviceEsimMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*ActivateDeviceEsimMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for activate device esim operations
// This implements the MockRegistrar interface
func (m *ActivateDeviceEsimMock) RegisterMocks() {
	// POST /deviceManagement/managedDevices/{id}/activateDeviceEsim - Managed Device Action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/([^/]+)/activateDeviceEsim$`,
		m.activateManagedDeviceEsimResponder())

	// POST /deviceManagement/comanagedDevices/{id}/activateDeviceEsim - Co-managed Device Action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/([^/]+)/activateDeviceEsim$`,
		m.activateComanagedDeviceEsimResponder())
}

// activateManagedDeviceEsimResponder handles POST requests to activate eSIM on managed devices
func (m *ActivateDeviceEsimMock) activateManagedDeviceEsimResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/managedDevices/")

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Validate carrier URL is provided
		carrierURL, hasCarrierURL := requestBody["carrierUrl"]
		if !hasCarrierURL || carrierURL == "" {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "error"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_activate_device_esim_not_found.json"))
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
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_success.json"))
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
			mockState.actionResults[deviceID]["managedDevice"] = response
			mockState.Unlock()

			return factories.EmptySuccessResponse(204)(req)
		}
	}
}

// activateComanagedDeviceEsimResponder handles POST requests to activate eSIM on co-managed devices
func (m *ActivateDeviceEsimMock) activateComanagedDeviceEsimResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract device ID from URL
		deviceID := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/comanagedDevices/")

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Validate carrier URL is provided
		carrierURL, hasCarrierURL := requestBody["carrierUrl"]
		if !hasCarrierURL || carrierURL == "" {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		}

		// Check for special test device IDs
		switch {
		case strings.Contains(deviceID, "error"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var errorResponse map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &errorResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResponse)
		case strings.Contains(deviceID, "not-found"):
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "post_activate_device_esim_not_found.json"))
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
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_success.json"))
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
			mockState.actionResults[deviceID]["comanagedDevice"] = response
			mockState.Unlock()

			return factories.EmptySuccessResponse(204)(req)
		}
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *ActivateDeviceEsimMock) RegisterErrorMocks() {
	// POST - Managed device activation error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/error-id/activateDeviceEsim$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid carrier URL"))

	// POST - Co-managed device activation error
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/error-id/activateDeviceEsim$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid carrier URL"))

	// POST - Device not found error (managed)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/not-found-id/activateDeviceEsim$`,
		factories.ErrorResponse(404, "NotFound", "Device not found"))

	// POST - Device not found error (co-managed)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/comanagedDevices/not-found-id/activateDeviceEsim$`,
		factories.ErrorResponse(404, "NotFound", "Device not found"))
}

// CleanupMockState clears all stored mock state
func (m *ActivateDeviceEsimMock) CleanupMockState() {
	mockState.Lock()
	mockState.actionResults = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockActivateDeviceEsimData returns sample activate device esim data for testing
func (m *ActivateDeviceEsimMock) GetMockActivateDeviceEsimData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_activate_device_esim_success.json"))
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
