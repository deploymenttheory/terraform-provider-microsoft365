package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of mailbox settings for consistent responses
var mockState struct {
	sync.Mutex
	mailboxSettings map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.mailboxSettings = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserMailboxSettingsMock provides mock responses for mailbox settings operations
type UserMailboxSettingsMock struct{}

// RegisterMocks registers HTTP mock responses for mailbox settings operations
func (m *UserMailboxSettingsMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.mailboxSettings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register predefined test mailbox settings
	registerTestMailboxSettings()

	// Register GET for mailbox settings by user ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+/mailboxSettings$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-2] // Second to last part is the user ID

			mockState.Lock()
			settingsData, exists := mockState.mailboxSettings[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, settingsData)
		})

	// Register PATCH for updating mailbox settings
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/users/[^/]+/mailboxSettings$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-2] // Second to last part is the user ID

			mockState.Lock()
			settingsData, exists := mockState.mailboxSettings[userId]
			if !exists {
				// Initialize if doesn't exist (first time creation)
				settingsData = map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#users('" + userId + "')/mailboxSettings",
				}
			}
			mockState.Unlock()

			var updateData map[string]any
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update mailbox settings data
			mockState.Lock()
			for k, v := range updateData {
				settingsData[k] = v
			}
			mockState.mailboxSettings[userId] = settingsData
			mockState.Unlock()

			// PATCH returns 204 No Content for mailbox settings
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserMailboxSettingsMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.mailboxSettings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register error response for user not found
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+/mailboxSettings$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
		})

	// Register error response for PATCH
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/users/[^/]+/mailboxSettings$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid mailbox settings data"}}`), nil
		})
}

// loadFixture loads a JSON fixture file from the tests/responses directory using the secure helpers package
func loadFixture(filename string) (map[string]any, error) {
	// Path relative to the mocks directory: ../tests/responses/
	fixturesPath := "../tests/responses/" + filename

	// Use the secure JSON parser from helpers package
	jsonContent, err := helpers.ParseJSONFile(fixturesPath)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// registerTestMailboxSettings registers predefined test mailbox settings from JSON fixtures
func registerTestMailboxSettings() {
	// Load minimal mailbox settings from fixture
	minimalSettingsData, err := loadFixture("mailbox_settings_minimal.json")
	if err != nil {
		// Fallback to inline data if fixture loading fails
		minimalSettingsData = map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#users('00000000-0000-0000-0000-000000000001')/mailboxSettings",
			"timeZone":       "UTC",
			"dateFormat":     "MM/dd/yyyy",
			"timeFormat":     "hh:mm tt",
		}
	}

	// Load maximal mailbox settings from fixture
	maximalSettingsData, err := loadFixture("mailbox_settings_maximal.json")
	if err != nil {
		// Fallback to inline data if fixture loading fails
		maximalSettingsData = map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#users('00000000-0000-0000-0000-000000000002')/mailboxSettings",
			"automaticRepliesSetting": map[string]any{
				"status":           "scheduled",
				"externalAudience": "all",
				"scheduledStartDateTime": map[string]any{
					"dateTime": "2016-03-14T07:00:00.0000000",
					"timeZone": "UTC",
				},
				"scheduledEndDateTime": map[string]any{
					"dateTime": "2016-03-28T07:00:00.0000000",
					"timeZone": "UTC",
				},
				"internalReplyMessage": "<html>\n<body>\n<p>I'm at our company's worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n",
				"externalReplyMessage": "<html>\n<body>\n<p>I'm at the Contoso worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n",
			},
			"timeZone": "UTC",
			"language": map[string]any{
				"locale":      "en-US",
				"displayName": "English (United States)",
			},
			"workingHours": map[string]any{
				"daysOfWeek": []any{"monday", "tuesday", "wednesday", "thursday", "friday"},
				"startTime":  "08:00:00.0000000",
				"endTime":    "17:00:00.0000000",
				"timeZone": map[string]any{
					"name": "Pacific Standard Time",
				},
			},
			"userPurpose":                           "user",
			"dateFormat":                            "MM/dd/yyyy",
			"timeFormat":                            "hh:mm tt",
			"delegateMeetingMessageDeliveryOptions": "sendToDelegateOnly",
		}
	}

	// Store mailbox settings in mock state
	mockState.Lock()
	mockState.mailboxSettings["00000000-0000-0000-0000-000000000001"] = minimalSettingsData
	mockState.mailboxSettings["00000000-0000-0000-0000-000000000002"] = maximalSettingsData
	mockState.Unlock()
}

// CleanupMockState clears the mock state
func (m *UserMailboxSettingsMock) CleanupMockState() {
	mockState.Lock()
	mockState.mailboxSettings = make(map[string]map[string]any)
	mockState.Unlock()
}
