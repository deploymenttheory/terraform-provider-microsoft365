package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	settings map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.settings = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// GroupSettingsMock provides mock responses for group settings operations
type GroupSettingsMock struct{}

// RegisterMocks registers HTTP mock responses for group settings operations
func (m *GroupSettingsMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.settings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register specific test settings
	registerTestSettings()

	// Register GET for group setting by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			settingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			settingData, exists := mockState.settings[settingId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Setting not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, settingData)
		})

	// Register GET for listing group settings
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			settings := make([]map[string]any, 0, len(mockState.settings))
			for _, setting := range mockState.settings {
				settings = append(settings, setting)
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups('12345678-1234-1234-1234-123456789012')/settings",
				"value":          settings,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating group settings
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			var settingData map[string]any
			err := json.NewDecoder(req.Body).Decode(&settingData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := settingData["templateId"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"templateId is required"}}`), nil
			}
			if _, ok := settingData["values"].([]any); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"values is required"}}`), nil
			}

			// Generate ID if not provided
			if settingData["id"] == nil {
				settingData["id"] = "test-setting-id"
			}

			// Set computed fields based on template ID
			now := time.Now().Format(time.RFC3339)
			settingData["createdDateTime"] = now

			// Set display name based on template ID
			templateId := settingData["templateId"].(string)
			switch templateId {
			case "08d542b9-071f-4e16-94b0-74abb372e3d9":
				settingData["displayName"] = "Group.Unified.Guest"
			case "62375ab9-6b52-47ed-826b-58e47e0e304b":
				settingData["displayName"] = "Group.Unified"
			default:
				settingData["displayName"] = "Unknown Template"
			}

			// Store setting in mock state
			settingId := settingData["id"].(string)
			mockState.Lock()
			mockState.settings[settingId] = settingData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, settingData)
		})

	// Register PATCH for updating group settings
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			settingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			settingData, exists := mockState.settings[settingId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Setting not found"}}`), nil
			}

			var updateData map[string]any
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update setting data
			mockState.Lock()

			// Special handling for updates that remove fields
			// If we're updating from maximal to minimal, we need to remove fields not in the minimal config
			isMinimalUpdate := false
			if _, hasTemplateId := updateData["templateId"]; hasTemplateId {
				if _, hasDisplayName := updateData["displayName"]; !hasDisplayName {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove fields that are not part of minimal configuration
				fieldsToRemove := []string{
					"displayName", "description",
				}

				for _, field := range fieldsToRemove {
					delete(settingData, field)
				}
			}

			// Apply the updates
			for k, v := range updateData {
				settingData[k] = v
			}

			mockState.settings[settingId] = settingData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, settingData)
		})

	// Register DELETE for removing group settings
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			settingId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.settings[settingId]
			if exists {
				delete(mockState.settings, settingId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *GroupSettingsMock) RegisterErrorMocks() {
	// Register error response for group setting creation
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings$`,
		factories.ErrorResponse(400, "BadRequest", "Error creating group setting"))

	// Register error response for group setting not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/12345678-1234-1234-1234-123456789012/settings/not-found-setting",
		factories.ErrorResponse(404, "ResourceNotFound", "Setting not found"))

	// Register error response for duplicate template ID
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/[^/]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			var settingData map[string]any
			err := json.NewDecoder(req.Body).Decode(&settingData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			if templateId, ok := settingData["templateId"].(string); ok && templateId == "error-template" {
				return factories.ErrorResponse(400, "BadRequest", "Setting with this templateId already exists")(req)
			}

			// Fallback to normal creation flow
			return nil, nil
		})
}

// registerTestSettings registers predefined test settings
func registerTestSettings() {
	// Minimal setting with only required attributes
	minimalSettingId := "test-setting-id"
	minimalSettingData := map[string]any{
		"id":          minimalSettingId,
		"templateId":  "08d542b9-071f-4e16-94b0-74abb372e3d9",
		"displayName": "Group.Unified.Guest",
		"values": []map[string]any{
			{
				"name":  "AllowToAddGuests",
				"value": "false",
			},
		},
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Maximal setting with all attributes
	maximalSettingId := "test-setting-id"
	maximalSettingData := map[string]any{
		"id":          maximalSettingId,
		"templateId":  "62375ab9-6b52-47ed-826b-58e47e0e304b",
		"displayName": "Group.Unified",
		"values": []map[string]any{
			{
				"name":  "ClassificationList",
				"value": "Confidential,Secret,Top Secret",
			},
			{
				"name":  "DefaultClassification",
				"value": "Confidential",
			},
			{
				"name":  "AllowGuestsToBeGroupOwner",
				"value": "false",
			},
			{
				"name":  "AllowGuestsToAccessGroups",
				"value": "true",
			},
			{
				"name":  "AllowToAddGuests",
				"value": "true",
			},
			{
				"name":  "UsageGuidelinesUrl",
				"value": "https://contoso.com/marketing-group-guidelines",
			},
		},
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Store settings in mock state
	mockState.Lock()
	mockState.settings[minimalSettingId] = minimalSettingData
	mockState.settings[maximalSettingId] = maximalSettingData
	mockState.Unlock()
}

// SetupImportTest sets up the mock state for an import test
func (m *GroupSettingsMock) SetupImportTest(groupId, settingId string) {
	// Create a setting for import testing
	importSettingData := map[string]any{
		"id":          settingId,
		"templateId":  "08d542b9-071f-4e16-94b0-74abb372e3d9",
		"displayName": "Group.Unified.Guest",
		"values": []map[string]any{
			{
				"name":  "AllowToAddGuests",
				"value": "false",
			},
		},
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Store setting in mock state
	mockState.Lock()
	mockState.settings[settingId] = importSettingData
	mockState.Unlock()
}
