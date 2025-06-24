package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	configurations map[string]map[string]interface{}
	assignments    map[string][]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.configurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// MacOSSoftwareUpdateConfigurationMock provides mock responses for macOS software update configuration operations
type MacOSSoftwareUpdateConfigurationMock struct{}

// RegisterMocks registers HTTP mock responses for macOS software update configuration operations
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]map[string]interface{})
	mockState.Unlock()

	// Register specific test configurations
	registerTestConfigurations()

	// Register GET for configuration by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.configurations[configId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Configuration not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, configData)
		})

	// Register GET for listing configurations
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			configs := make([]map[string]interface{}, 0, len(mockState.configurations))
			for _, config := range mockState.configurations {
				configs = append(configs, config)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          configs,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating configurations
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			var configData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&configData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := configData["displayName"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
			}

			// Generate ID if not provided
			if configData["id"] == nil {
				configData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			configData["createdDateTime"] = now
			configData["lastModifiedDateTime"] = now
			configData["version"] = 1

			// Set @odata.type for macOS software update configuration
			configData["@odata.type"] = "#microsoft.graph.macOSSoftwareUpdateConfiguration"

			// Ensure collection fields are initialized
			if configData["roleScopeTagIds"] == nil {
				configData["roleScopeTagIds"] = []string{"0"}
			}

			// Store configuration in mock state
			configId := configData["id"].(string)
			mockState.Lock()
			mockState.configurations[configId] = configData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, configData)
		})

	// Register PATCH for updating configurations
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.configurations[configId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Configuration not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update configuration data
			mockState.Lock()

			// Update last modified time
			now := time.Now().Format(time.RFC3339)
			configData["lastModifiedDateTime"] = now

			// Increment version
			if version, ok := configData["version"].(float64); ok {
				configData["version"] = version + 1
			} else {
				configData["version"] = 1
			}

			// Special handling for updates that remove fields
			// If we're updating from maximal to minimal, we need to remove fields not in the minimal config
			// Check if this is a minimal update by looking for key indicators
			isMinimalUpdate := false
			if displayName, ok := updateData["displayName"].(string); ok {
				if displayName == "Minimal macOS Software Update Configuration" {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove fields that are not part of minimal configuration
				fieldsToRemove := []string{
					"description", "maxUserDeferralsCount", "priority",
				}

				for _, field := range fieldsToRemove {
					delete(configData, field)
				}
			}

			// Apply the updates
			for k, v := range updateData {
				configData[k] = v
			}

			mockState.configurations[configId] = configData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, configData)
		})

	// Register DELETE for removing configurations
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.configurations[configId]
			if exists {
				delete(mockState.configurations, configId)
				delete(mockState.assignments, configId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-2]

			mockState.Lock()
			assignments, exists := mockState.assignments[configId]
			mockState.Unlock()

			if !exists {
				assignments = []map[string]interface{}{}
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations('configId')/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-2]

			var assignmentData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&assignmentData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Process assignments from the request
			assignments := []map[string]interface{}{}
			if assignmentsArray, ok := assignmentData["assignments"].([]interface{}); ok {
				for _, assignment := range assignmentsArray {
					if assignmentMap, ok := assignment.(map[string]interface{}); ok {
						assignmentId := uuid.New().String()
						assignmentMap["id"] = assignmentId
						assignments = append(assignments, assignmentMap)
					}
				}
			}

			// Store assignments in mock state
			mockState.Lock()
			mockState.assignments[configId] = assignments
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterErrorMocks() {
	// Register error response for configuration creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		factories.ErrorResponse(400, "BadRequest", "Error creating configuration"))

	// Register error response for configuration not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/not-found-config",
		factories.ErrorResponse(404, "ResourceNotFound", "Configuration not found"))

	// Register error response for duplicate configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			var configData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&configData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			if name, ok := configData["displayName"].(string); ok && name == "Error macOS Software Update Configuration" {
				return factories.ErrorResponse(400, "BadRequest", "Configuration with this name already exists")(req)
			}

			// Fallback to normal creation flow
			return nil, nil
		})
}

// registerTestConfigurations registers predefined test configurations
func registerTestConfigurations() {
	// Minimal configuration with only required attributes
	minimalConfigId := "00000000-0000-0000-0000-000000000001"
	minimalConfigData := map[string]interface{}{
		"id":                                 minimalConfigId,
		"displayName":                        "Minimal macOS Software Update Configuration",
		"@odata.type":                        "#microsoft.graph.macOSSoftwareUpdateConfiguration",
		"criticalUpdateBehavior":             "default",
		"configDataUpdateBehavior":           "default",
		"firmwareUpdateBehavior":             "default",
		"allOtherUpdateBehavior":             "default",
		"updateScheduleType":                 "alwaysUpdate",
		"updateTimeWindowUtcOffsetInMinutes": float64(0),
		"roleScopeTagIds":                    []string{"0"},
		"createdDateTime":                    "2023-01-01T00:00:00Z",
		"lastModifiedDateTime":               "2023-01-01T00:00:00Z",
		"version":                            float64(1),
	}

	// Maximal configuration with all attributes
	maximalConfigId := "00000000-0000-0000-0000-000000000002"
	maximalConfigData := map[string]interface{}{
		"id":                                 maximalConfigId,
		"displayName":                        "Maximal macOS Software Update Configuration",
		"description":                        "This is a comprehensive configuration with all fields populated",
		"@odata.type":                        "#microsoft.graph.macOSSoftwareUpdateConfiguration",
		"criticalUpdateBehavior":             "installASAP",
		"configDataUpdateBehavior":           "notifyOnly",
		"firmwareUpdateBehavior":             "downloadOnly",
		"allOtherUpdateBehavior":             "installLater",
		"updateScheduleType":                 "updateDuringTimeWindows",
		"updateTimeWindowUtcOffsetInMinutes": float64(60),
		"maxUserDeferralsCount":              float64(3),
		"priority":                           "high",
		"roleScopeTagIds":                    []string{"0", "1"},
		"createdDateTime":                    "2023-01-01T00:00:00Z",
		"lastModifiedDateTime":               "2023-01-01T00:00:00Z",
		"version":                            float64(1),
	}

	// Store configurations in mock state
	mockState.Lock()
	mockState.configurations[minimalConfigId] = minimalConfigData
	mockState.configurations[maximalConfigId] = maximalConfigData

	// Set up minimal assignments
	mockState.assignments[minimalConfigId] = []map[string]interface{}{
		{
			"id": "00000000-0000-0000-0000-000000000101",
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
				"deviceAndAppManagementAssignmentFilterId":   nil,
				"deviceAndAppManagementAssignmentFilterType": "none",
			},
		},
	}

	// Set up maximal assignments
	mockState.assignments[maximalConfigId] = []map[string]interface{}{
		{
			"id": "00000000-0000-0000-0000-000000000102",
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
				"deviceAndAppManagementAssignmentFilterId":   nil,
				"deviceAndAppManagementAssignmentFilterType": "none",
			},
		},
	}
	mockState.Unlock()
}
