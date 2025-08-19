package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	deviceConfigurations map[string]map[string]interface{}
	assignments          map[string][]interface{}
}

func init() {
	mockState.deviceConfigurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("macos_device_configuration_templates", &MacosDeviceConfigurationTemplatesMock{})
}

type MacosDeviceConfigurationTemplatesMock struct{}

var _ mocks.MockRegistrar = (*MacosDeviceConfigurationTemplatesMock)(nil)

func (m *MacosDeviceConfigurationTemplatesMock) RegisterMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic dependency mocks
	m.registerDependencyMocks()

	// GET /deviceManagement/deviceConfigurations - List device configurations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_macos_device_configuration_list.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON: %s"}}`, err.Error())), nil
		}
		return httpmock.NewStringResponse(200, jsonStr), nil
	})

	// POST /deviceManagement/deviceConfigurations - Create device configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Determine configuration type and return appropriate response
		odataType, ok := requestBody["@odata.type"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Missing @odata.type"}}`), nil
		}

		var responseFile string
		var id string

		switch odataType {
		case "#microsoft.graph.macOSCustomConfiguration":
			responseFile = "../tests/responses/validate_create/post_macos_custom_configuration_success.json"
			id = "12345678-1234-1234-1234-123456789012"
		case "#microsoft.graph.macOSCustomAppConfiguration":
			responseFile = "../tests/responses/validate_create/post_macos_preference_file_success.json"
			id = "87654321-4321-4321-4321-210987654321"
		case "#microsoft.graph.macOSTrustedRootCertificate":
			responseFile = "../tests/responses/validate_create/post_macos_trusted_certificate_success.json"
			id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
		case "#microsoft.graph.macOSScepCertificateProfile":
			responseFile = "../tests/responses/validate_create/post_macos_scep_certificate_success.json"
			id = "ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb"
		case "#microsoft.graph.macOSPkcsCertificateProfile":
			responseFile = "../tests/responses/validate_create/post_macos_pkcs_certificate_success.json"
			id = "11111111-2222-3333-4444-555555555555"
		default:
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Unsupported configuration type"}}`), nil
		}

		jsonStr, err := helpers.ParseJSONFile(responseFile)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON file '%s': %s"}}`, responseFile, err.Error())), nil
		}

		// Store in mock state
		mockState.Lock()
		requestBody["id"] = id
		mockState.deviceConfigurations[id] = requestBody
		mockState.Unlock()

		return httpmock.NewStringResponse(201, jsonStr), nil
	})

	// GET /deviceManagement/deviceConfigurations/{id} - Get specific device configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)(\?.*)?$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		if len(segments) < 4 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL"}}`), nil
		}

		configId := segments[len(segments)-1]

		// Handle assignments expansion
		if req.URL.Query().Get("$expand") == "assignments" {
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_macos_custom_configuration.json")
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON: %s"}}`, err.Error())), nil
			}
			return httpmock.NewStringResponse(200, jsonStr), nil
		}

		// Return appropriate configuration based on ID
		var responseFile string
		switch configId {
		case "12345678-1234-1234-1234-123456789012":
			responseFile = "../tests/responses/validate_get/get_macos_custom_configuration.json"
		case "87654321-4321-4321-4321-210987654321":
			responseFile = "../tests/responses/validate_get/get_macos_preference_file.json"
		case "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee":
			responseFile = "../tests/responses/validate_get/get_macos_trusted_certificate.json"
		case "ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb":
			responseFile = "../tests/responses/validate_get/get_macos_scep_certificate.json"
		case "11111111-2222-3333-4444-555555555555":
			responseFile = "../tests/responses/validate_get/get_macos_pkcs_certificate.json"
		default:
			responseFile = "../tests/responses/validate_get/get_macos_custom_configuration.json"
		}

		jsonStr, err := helpers.ParseJSONFile(responseFile)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON file '%s': %s"}}`, responseFile, err.Error())), nil
		}

		return httpmock.NewStringResponse(200, jsonStr), nil
	})

	// PATCH /deviceManagement/deviceConfigurations/{id} - Update device configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		if len(segments) < 4 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL"}}`), nil
		}

		configId := segments[len(segments)-1]

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Update mock state
		mockState.Lock()
		if existing, exists := mockState.deviceConfigurations[configId]; exists {
			for k, v := range requestBody {
				existing[k] = v
			}
		}
		mockState.Unlock()

		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/patch_macos_device_configuration_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON: %s"}}`, err.Error())), nil
		}

		return httpmock.NewStringResponse(200, jsonStr), nil
	})

	// DELETE /deviceManagement/deviceConfigurations/{id} - Delete device configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		if len(segments) < 4 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL"}}`), nil
		}

		configId := segments[len(segments)-1]

		// Remove from mock state
		mockState.Lock()
		delete(mockState.deviceConfigurations, configId)
		delete(mockState.assignments, configId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// POST /deviceManagement/deviceConfigurations/{id}/assign - Assign device configuration
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assign$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		if len(segments) < 5 {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL"}}`), nil
		}

		configId := segments[len(segments)-2]

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Store assignments in mock state
		mockState.Lock()
		if assignments, ok := requestBody["assignments"].([]interface{}); ok {
			mockState.assignments[configId] = assignments
		}
		mockState.Unlock()

		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_assign/post_macos_device_configuration_assign_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to parse JSON: %s"}}`, err.Error())), nil
		}

		return httpmock.NewStringResponse(200, jsonStr), nil
	})
}

func (m *MacosDeviceConfigurationTemplatesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic dependency mocks
	m.registerDependencyMocks()

	// POST /deviceManagement/deviceConfigurations - Create device configuration (Error)
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_macos_device_configuration_error.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalError","message":"Failed to parse error response"}}`), nil
		}
		return httpmock.NewStringResponse(400, jsonStr), nil
	})

	// GET /deviceManagement/deviceConfigurations/{id} - Not found error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)(\?.*)?$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_delete/get_macos_device_configuration_not_found.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalError","message":"Failed to parse error response"}}`), nil
		}
		return httpmock.NewStringResponse(404, jsonStr), nil
	})

	// Other operations also return errors
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error updating macOS device configuration template"}}`), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error deleting macOS device configuration template"}}`), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assign$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error assigning macOS device configuration template"}}`), nil
	})
}

func (m *MacosDeviceConfigurationTemplatesMock) CleanupMockState() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.Unlock()
}

// registerDependencyMocks registers mocks for dependencies like groups, role scope tags, and assignment filters
func (m *MacosDeviceConfigurationTemplatesMock) registerDependencyMocks() {
	// Mock role scope tags
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		tagId := segments[len(segments)-1]

		response := map[string]interface{}{
			"@odata.type": "#microsoft.graph.roleScopeTag",
			"id":          tagId,
			"displayName": fmt.Sprintf("Role Scope Tag %s", tagId),
			"description": "Test role scope tag",
		}

		jsonBytes, _ := json.Marshal(response)
		return httpmock.NewStringResponse(200, string(jsonBytes)), nil
	})

	// Mock groups
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		groupId := segments[len(segments)-1]

		response := map[string]interface{}{
			"@odata.type":     "#microsoft.graph.group",
			"id":              groupId,
			"displayName":     fmt.Sprintf("Test Group %s", groupId),
			"description":     "Test group for device configuration",
			"groupTypes":      []string{},
			"securityEnabled": true,
		}

		jsonBytes, _ := json.Marshal(response)
		return httpmock.NewStringResponse(200, string(jsonBytes)), nil
	})

	// Mock assignment filters
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		filterId := segments[len(segments)-1]

		response := map[string]interface{}{
			"@odata.type":   "#microsoft.graph.deviceAndAppManagementAssignmentFilter",
			"id":            filterId,
			"displayName":   fmt.Sprintf("Test Assignment Filter %s", filterId),
			"description":   "Test assignment filter",
			"platform":      "macOS",
			"rule":          "(device.deviceOwnership -eq \"Corporate\")",
			"roleScopeTags": []string{"0"},
		}

		jsonBytes, _ := json.Marshal(response)
		return httpmock.NewStringResponse(200, string(jsonBytes)), nil
	})

	// Mock device configuration assignments endpoint for GET requests
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assignments$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-2]

		mockState.Lock()
		assignments := mockState.assignments[configId]
		mockState.Unlock()

		response := map[string]interface{}{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations('12345678-1234-1234-1234-123456789012')/assignments",
			"value":          assignments,
		}

		if assignments == nil {
			response["value"] = []interface{}{}
		}

		jsonBytes, _ := json.Marshal(response)
		return httpmock.NewStringResponse(200, string(jsonBytes)), nil
	})

	// Register authentication mocks
	httpmock.RegisterResponder("POST", "https://login.microsoftonline.com/common/oauth2/v2.0/token", func(req *http.Request) (*http.Response, error) {
		response := map[string]interface{}{
			"access_token": "mock_access_token_" + uuid.New().String(),
			"token_type":   "Bearer",
			"expires_in":   3600,
		}
		jsonBytes, _ := json.Marshal(response)
		return httpmock.NewStringResponse(200, string(jsonBytes)), nil
	})
}
