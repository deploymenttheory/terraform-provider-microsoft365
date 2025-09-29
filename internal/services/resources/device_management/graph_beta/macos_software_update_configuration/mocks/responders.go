package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	softwareUpdateConfigurations map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.softwareUpdateConfigurations = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// MacOSSoftwareUpdateConfigurationMock provides mock responses for macOS software update configuration operations
type MacOSSoftwareUpdateConfigurationMock struct{}

// RegisterMocks registers HTTP mock responses for macOS software update configuration operations
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.softwareUpdateConfigurations = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing software update configurations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			configs := make([]map[string]any, 0, len(mockState.softwareUpdateConfigurations))
			for _, config := range mockState.softwareUpdateConfigurations {
				configs = append(configs, config)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          configs,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual software update configuration
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.softwareUpdateConfigurations[configId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Software update configuration not found"}}`), nil
			}

			// Create response copy
			responseCopy := make(map[string]any)
			for k, v := range configData {
				responseCopy[k] = v
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the config data
				if assignments, hasAssignments := configData["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
						responseCopy["assignments"] = assignments
					} else {
						// If assignments array is empty, return empty array (not null)
						responseCopy["assignments"] = []interface{}{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					responseCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, responseCopy)
		})

	// Register POST for creating software update configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new configuration ID
			configId := uuid.New().String()

			// Create configuration data - only include fields that were provided or have defaults
			configData := map[string]any{
				"@odata.type":              "#microsoft.graph.macOSSoftwareUpdateConfiguration",
				"id":                       configId,
				"displayName":              requestBody["displayName"],
				"updateScheduleType":       requestBody["updateScheduleType"],
				"criticalUpdateBehavior":   requestBody["criticalUpdateBehavior"],
				"configDataUpdateBehavior": requestBody["configDataUpdateBehavior"],
				"firmwareUpdateBehavior":   requestBody["firmwareUpdateBehavior"],
				"allOtherUpdateBehavior":   requestBody["allOtherUpdateBehavior"],
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				configData["description"] = description
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				configData["roleScopeTagIds"] = roleScopeTagIds
			} else {
				configData["roleScopeTagIds"] = []string{"0"} // Default value
			}
			if updateTimeWindowUtcOffsetInMinutes, exists := requestBody["updateTimeWindowUtcOffsetInMinutes"]; exists {
				configData["updateTimeWindowUtcOffsetInMinutes"] = updateTimeWindowUtcOffsetInMinutes
			}
			if customUpdateTimeWindows, exists := requestBody["customUpdateTimeWindows"]; exists {
				configData["customUpdateTimeWindows"] = customUpdateTimeWindows
			}
			if maxUserDeferralsCount, exists := requestBody["maxUserDeferralsCount"]; exists {
				configData["maxUserDeferralsCount"] = maxUserDeferralsCount
			}
			if priority, exists := requestBody["priority"]; exists {
				configData["priority"] = priority
			}

			// Add computed fields that are always returned by the API
			configData["createdDateTime"] = "2024-01-01T00:00:00Z"
			configData["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Initialize assignments as empty array
			configData["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.softwareUpdateConfigurations[configId] = configData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, configData)
		})

	// Register PATCH for updating software update configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.softwareUpdateConfigurations[configId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Software update configuration not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update configuration data
			mockState.Lock()

			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior

			// For optional fields, if they're not in the request, remove them
			optionalFields := []string{"description", "updateTimeWindowUtcOffsetInMinutes", "customUpdateTimeWindows", "maxUserDeferralsCount", "priority"}
			for _, field := range optionalFields {
				if _, hasField := requestBody[field]; !hasField {
					delete(configData, field)
				}
			}

			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(configData, key)
				} else {
					configData[key] = value
				}
			}
			// Ensure the ID and @odata.type are preserved and update timestamp
			configData["id"] = configId
			configData["@odata.type"] = "#microsoft.graph.macOSSoftwareUpdateConfiguration"
			configData["lastModifiedDateTime"] = "2024-01-01T01:00:00Z"
			mockState.softwareUpdateConfigurations[configId] = configData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, configData)
		})

	// Register DELETE for removing software update configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.softwareUpdateConfigurations[configId]
			if exists {
				delete(mockState.softwareUpdateConfigurations, configId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Software update configuration not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-2] // deviceConfigurations/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the configuration
			mockState.Lock()
			if configData, exists := mockState.softwareUpdateConfigurations[configId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									// Generate a unique assignment ID
									assignmentId := uuid.New().String()

									// Create assignment in the format the API returns
									// The API returns the target exactly as submitted but with additional metadata
									targetCopy := make(map[string]any)
									for k, v := range target {
										targetCopy[k] = v
									}

									graphAssignment := map[string]any{
										"id":     assignmentId,
										"target": targetCopy,
									}
									graphAssignments = append(graphAssignments, graphAssignment)
								}
							}
						}
						configData["assignments"] = graphAssignments
					} else {
						// Set empty assignments array instead of deleting
						configData["assignments"] = []interface{}{}
					}
				} else {
					// Set empty assignments array instead of deleting
					configData["assignments"] = []interface{}{}
				}
				mockState.softwareUpdateConfigurations[configId] = configData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/assignments",
				"value":          []map[string]any{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Dynamic mocks will handle all test cases
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterErrorMocks() {
	// Register GET for listing software update configurations (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating software update configuration with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for software update configuration not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/not-found-config",
		factories.ErrorResponse(404, "ResourceNotFound", "Software update configuration not found"))
}
