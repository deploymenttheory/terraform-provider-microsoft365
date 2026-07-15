package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	deviceConfigurations map[string]map[string]any
	assignments          map[string][]any
	secretValues         map[string]string
}

func init() {
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.secretValues = make(map[string]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_custom_configuration", &WindowsCustomConfigurationMock{})
}

type WindowsCustomConfigurationMock struct{}

var _ mocks.MockRegistrar = (*WindowsCustomConfigurationMock)(nil)

// maskEncryptedOmaSettings simulates the Graph API behaviour of masking string based OMA setting
// values in responses: omaSettingString values are replaced with "****", flagged as encrypted and
// given a secret reference value id resolvable via getOmaSettingPlainTextValue.
func maskEncryptedOmaSettings(config map[string]any, configId string) {
	omaSettings, ok := config["omaSettings"].([]any)
	if !ok {
		return
	}

	for idx, rawSetting := range omaSettings {
		setting, ok := rawSetting.(map[string]any)
		if !ok {
			continue
		}
		if setting["@odata.type"] != "#microsoft.graph.omaSettingString" {
			continue
		}
		value, ok := setting["value"].(string)
		if !ok || value == "****" {
			continue
		}

		secretReferenceValueId := fmt.Sprintf("%s_secret_%d", configId, idx)
		mockState.secretValues[secretReferenceValueId] = value
		setting["value"] = "****"
		setting["isEncrypted"] = true
		setting["secretReferenceValueId"] = secretReferenceValueId
	}
}

func (m *WindowsCustomConfigurationMock) RegisterMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.secretValues = make(map[string]string)
	mockState.Unlock()

	m.registerDependencyMocks()

	// POST /deviceManagement/deviceConfigurations - Create device configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{"code": "BadRequest", "message": "Invalid request body"},
			})
		}

		if requestBody["@odata.type"] != "#microsoft.graph.windows10CustomConfiguration" {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{"code": "BadRequest", "message": "Unsupported configuration type"},
			})
		}

		id := uuid.New().String()
		requestBody["id"] = id
		requestBody["createdDateTime"] = "2024-01-01T00:00:00Z"
		requestBody["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"
		requestBody["version"] = 1

		mockState.Lock()
		maskEncryptedOmaSettings(requestBody, id)
		mockState.deviceConfigurations[id] = requestBody
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, requestBody)
	})

	// GET /deviceManagement/deviceConfigurations/{id}/getOmaSettingPlainTextValue(secretReferenceValueId='{id}')
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/[^/]+/getOmaSettingPlainTextValue.*$`, func(req *http.Request) (*http.Response, error) {
		rawURL := req.URL.Path
		start := strings.Index(rawURL, "secretReferenceValueId='")
		if start == -1 {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{"code": "BadRequest", "message": "Missing secretReferenceValueId"},
			})
		}
		secretReferenceValueId := rawURL[start+len("secretReferenceValueId='"):]
		if end := strings.Index(secretReferenceValueId, "'"); end != -1 {
			secretReferenceValueId = secretReferenceValueId[:end]
		}

		mockState.Lock()
		value, exists := mockState.secretValues[secretReferenceValueId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewJsonResponse(404, map[string]any{
				"error": map[string]any{"code": "ResourceNotFound", "message": "Secret reference not found"},
			})
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#Edm.String",
			"value":          value,
		})
	})

	// GET /deviceManagement/deviceConfigurations/{id} - Get specific device configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)(\?.*)?$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-1]

		mockState.Lock()
		config, exists := mockState.deviceConfigurations[configId]
		assignments := mockState.assignments[configId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewJsonResponse(404, map[string]any{
				"error": map[string]any{"code": "ResourceNotFound", "message": "Resource not found"},
			})
		}

		responseObj := make(map[string]any, len(config)+1)
		for k, v := range config {
			responseObj[k] = v
		}

		if req.URL.Query().Get("$expand") == "assignments" {
			if assignments == nil {
				assignments = []any{}
			}
			responseObj["assignments"] = assignments
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// PATCH /deviceManagement/deviceConfigurations/{id} - Update device configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, exists := mockState.deviceConfigurations[configId]
		if exists {
			maskEncryptedOmaSettings(requestBody, configId)
			for k, v := range requestBody {
				existing[k] = v
			}
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, existing)
	})

	// DELETE /deviceManagement/deviceConfigurations/{id} - Delete device configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-1]

		mockState.Lock()
		delete(mockState.deviceConfigurations, configId)
		delete(mockState.assignments, configId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// POST /deviceManagement/deviceConfigurations/{id}/assign - Assign device configuration
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assign$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-2]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{"code": "BadRequest", "message": "Invalid request body"},
			})
		}

		storedAssignments := []any{}
		if assignments, ok := requestBody["assignments"].([]any); ok {
			for idx, rawAssignment := range assignments {
				assignment, ok := rawAssignment.(map[string]any)
				if !ok {
					continue
				}
				storedAssignments = append(storedAssignments, map[string]any{
					"@odata.type": "#microsoft.graph.deviceConfigurationAssignment",
					"id":          fmt.Sprintf("%s_assignment_%d", configId, idx),
					"target":      assignment["target"],
				})
			}
		}

		mockState.Lock()
		mockState.assignments[configId] = storedAssignments
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, map[string]any{"value": storedAssignments})
	})
}

func (m *WindowsCustomConfigurationMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.secretValues = make(map[string]string)
	mockState.Unlock()

	m.registerDependencyMocks()

	// POST /deviceManagement/deviceConfigurations - Create device configuration (Error)
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(400, map[string]any{
			"error": map[string]any{"code": "BadRequest", "message": "Error creating windows custom configuration"},
		})
	})

	// GET /deviceManagement/deviceConfigurations/{id} - Not found error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)(\?.*)?$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
	})

	// Other operations also return errors
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error updating windows custom configuration"}}`), nil
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error deleting windows custom configuration"}}`), nil
	})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assign$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Error assigning windows custom configuration"}}`), nil
	})
}

func (m *WindowsCustomConfigurationMock) CleanupMockState() {
	mockState.Lock()
	mockState.deviceConfigurations = make(map[string]map[string]any)
	mockState.assignments = make(map[string][]any)
	mockState.secretValues = make(map[string]string)
	mockState.Unlock()
}

// registerDependencyMocks registers mocks for dependencies like groups, role scope tags, and assignment filters
func (m *WindowsCustomConfigurationMock) registerDependencyMocks() {
	// Mock role scope tags
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		tagId := segments[len(segments)-1]

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.type": "#microsoft.graph.roleScopeTag",
			"id":          tagId,
			"displayName": fmt.Sprintf("Role Scope Tag %s", tagId),
			"description": "Test role scope tag",
		})
	})

	// Mock groups
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		groupId := segments[len(segments)-1]

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.type":     "#microsoft.graph.group",
			"id":              groupId,
			"displayName":     fmt.Sprintf("Test Group %s", groupId),
			"description":     "Test group for device configuration",
			"groupTypes":      []string{},
			"securityEnabled": true,
		})
	})

	// Mock assignment filters
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		filterId := segments[len(segments)-1]

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.type":   "#microsoft.graph.deviceAndAppManagementAssignmentFilter",
			"id":            filterId,
			"displayName":   fmt.Sprintf("Test Assignment Filter %s", filterId),
			"description":   "Test assignment filter",
			"platform":      "windows10AndLater",
			"rule":          "(device.deviceOwnership -eq \"Corporate\")",
			"roleScopeTags": []string{"0"},
		})
	})

	// Mock device configuration assignments endpoint for GET requests
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)/assignments$`, func(req *http.Request) (*http.Response, error) {
		segments := strings.Split(req.URL.Path, "/")
		configId := segments[len(segments)-2]

		mockState.Lock()
		assignments := mockState.assignments[configId]
		mockState.Unlock()

		if assignments == nil {
			assignments = []any{}
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations('%s')/assignments", configId),
			"value":          assignments,
		})
	})

	// Register authentication mocks
	httpmock.RegisterResponder("POST", "https://login.microsoftonline.com/common/oauth2/v2.0/token", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{
			"access_token": "mock_access_token_" + uuid.New().String(),
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})
}
