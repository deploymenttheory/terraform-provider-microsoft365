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
	userSettings map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.userSettings = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserSettingMock provides mock responses for user setting operations
type UserSettingMock struct{}

// RegisterMocks registers HTTP mock responses for user setting operations
func (m *UserSettingMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.userSettings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing user settings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			userSettings := make([]map[string]any, 0, len(mockState.userSettings))
			for _, setting := range mockState.userSettings {
				userSettings = append(userSettings, setting)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/userSettings",
				"value":          userSettings,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual user setting
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userSettingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userSettingData, exists := mockState.userSettings[userSettingId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User setting not found"}}`), nil
			}

			// Always include assignments if they exist (since this mock should simulate the expand=assignments behavior)
			// The real API only includes assignments when expand=assignments is specified, but for testing we'll always include them
			if assignments, hasAssignments := userSettingData["assignments"]; hasAssignments {
				// Ensure assignments are included in the response
				responseCopy := make(map[string]any)
				for k, v := range userSettingData {
					responseCopy[k] = v
				}
				responseCopy["assignments"] = assignments
				return httpmock.NewJsonResponse(200, responseCopy)
			}

			return httpmock.NewJsonResponse(200, userSettingData)
		})

	// Register POST for creating user setting
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new user setting ID
			userSettingId := uuid.New().String()

			// Create user setting data with all fields
			userSettingData := map[string]any{
				"id":                   userSettingId,
				"displayName":          requestBody["displayName"],
				"localAdminEnabled":    getOrDefault(requestBody, "localAdminEnabled", false),
				"resetEnabled":         getOrDefault(requestBody, "resetEnabled", false),
				"selfServiceEnabled":   getOrDefault(requestBody, "selfServiceEnabled", false),
				"createdDateTime":      "2024-01-01T00:00:00Z",
				"lastModifiedDateTime": "2024-01-01T00:00:00Z",
			}

			// Handle nested attributes
			if restorePointSetting, exists := requestBody["restorePointSetting"]; exists {
				userSettingData["restorePointSetting"] = restorePointSetting
			}

			if crossRegionSetting, exists := requestBody["crossRegionDisasterRecoverySetting"]; exists {
				userSettingData["crossRegionDisasterRecoverySetting"] = crossRegionSetting
			}

			if notificationSetting, exists := requestBody["notificationSetting"]; exists {
				userSettingData["notificationSetting"] = notificationSetting
			}

			// Store in mock state
			mockState.Lock()
			mockState.userSettings[userSettingId] = userSettingData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, userSettingData)
		})

	// Register PATCH for updating user setting
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userSettingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userSettingData, exists := mockState.userSettings[userSettingId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User setting not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update user setting data
			mockState.Lock()

			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior

			// For nested attributes, if they're not in the request, remove them
			optionalNestedFields := []string{"restorePointSetting", "notificationSetting", "crossRegionDisasterRecoverySetting"}
			for _, field := range optionalNestedFields {
				if _, hasField := requestBody[field]; !hasField {
					delete(userSettingData, field)
				}
			}

			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(userSettingData, key)
				} else {
					userSettingData[key] = value
				}
			}
			// Ensure the ID is preserved and update timestamp
			userSettingData["id"] = userSettingId
			userSettingData["lastModifiedDateTime"] = "2024-01-01T01:00:00Z"
			mockState.userSettings[userSettingId] = userSettingData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, userSettingData)
		})

	// Register DELETE for removing user setting
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userSettingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.userSettings[userSettingId]
			if exists {
				delete(mockState.userSettings, userSettingId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User setting not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userSettingId := urlParts[len(urlParts)-3] // userSettings/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the user setting
			mockState.Lock()
			if userSettingData, exists := mockState.userSettings[userSettingId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments {
					// Convert assignments to the format returned by Graph API
					assignmentList := assignments.([]interface{})
					graphAssignments := make([]interface{}, len(assignmentList))

					for i, assignment := range assignmentList {
						assignmentMap := assignment.(map[string]any)
						target := assignmentMap["target"].(map[string]any)

						// Create the assignment response format
						graphAssignment := map[string]any{
							"id": target["groupId"], // Use groupId as assignment ID
							"target": map[string]any{
								"@odata.type": "#microsoft.graph.cloudPcManagementGroupAssignmentTarget",
								"groupId":     target["groupId"],
							},
						}
						graphAssignments[i] = graphAssignment
					}

					userSettingData["assignments"] = graphAssignments
					mockState.userSettings[userSettingId] = userSettingData
				}
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register specific user setting mocks for testing
	registerSpecificUserSettingMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserSettingMock) RegisterErrorMocks() {
	// Register error response for creating user setting with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for user setting not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/not-found-setting",
		factories.ErrorResponse(404, "ResourceNotFound", "User setting not found"))
}

// registerSpecificUserSettingMocks registers mocks for specific test scenarios
func registerSpecificUserSettingMocks() {
	// Minimal user setting
	minimalUserSettingId := "11111111-1111-1111-1111-111111111111"
	minimalUserSettingData := map[string]any{
		"id":                 minimalUserSettingId,
		"displayName":        "Test Minimal User Setting",
		"localAdminEnabled":  false,
		"resetEnabled":       false,
		"selfServiceEnabled": false,
		"restorePointSetting": map[string]any{
			"frequencyInHours":   12,
			"frequencyType":      "default",
			"userRestoreEnabled": false,
		},
		"createdDateTime":      "2024-01-01T00:00:00Z",
		"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	}

	mockState.Lock()
	mockState.userSettings[minimalUserSettingId] = minimalUserSettingData
	mockState.Unlock()

	// Maximal user setting
	maximalUserSettingId := "22222222-2222-2222-2222-222222222222"
	maximalUserSettingData := map[string]any{
		"id":                 maximalUserSettingId,
		"displayName":        "Test Maximal User Setting",
		"localAdminEnabled":  true,
		"resetEnabled":       true,
		"selfServiceEnabled": false,
		"restorePointSetting": map[string]any{
			"frequencyInHours":   12,
			"frequencyType":      "default",
			"userRestoreEnabled": true,
		},
		"crossRegionDisasterRecoverySetting": map[string]any{
			"maintainCrossRegionRestorePointEnabled": true,
			"userInitiatedDisasterRecoveryAllowed":   true,
			"disasterRecoveryType":                   "premium",
			"disasterRecoveryNetworkSetting": map[string]any{
				"@odata.type": "#microsoft.graph.cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting",
				"networkType": "microsoftHosted",
			},
		},
		"notificationSetting": map[string]any{
			"restartPromptsDisabled": false,
		},
		"assignments": []interface{}{
			map[string]any{
				"id": "test-group-id-12345",
				"target": map[string]any{
					"@odata.type": "#microsoft.graph.cloudPcManagementGroupAssignmentTarget",
					"groupId":     "test-group-id-12345",
				},
			},
		},
		"createdDateTime":      "2024-01-01T00:00:00Z",
		"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	}

	mockState.Lock()
	mockState.userSettings[maximalUserSettingId] = maximalUserSettingData
	mockState.Unlock()

	// Register specific GET for these user settings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/"+minimalUserSettingId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			userSettingData := mockState.userSettings[minimalUserSettingId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, userSettingData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/userSettings/"+maximalUserSettingId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			userSettingData := mockState.userSettings[maximalUserSettingId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, userSettingData)
		})
}

// getOrDefault returns the value from the map or a default value if the key doesn't exist
func getOrDefault(m map[string]any, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}
