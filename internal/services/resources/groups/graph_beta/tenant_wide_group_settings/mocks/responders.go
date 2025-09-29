package mocks

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

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

// MockTenantWideGroupSettingsResponders sets up all the HTTP mock responders for tenant-wide group settings testing
func MockTenantWideGroupSettingsResponders() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.settings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register 500 error for error testing - register this first for higher priority
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/settings",
		func(req *http.Request) (*http.Response, error) {
			// Read the request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close() // Close the body after reading

			var requestBody map[string]any
			_ = json.Unmarshal(body, &requestBody)

			// Check if this is an error test by looking at the values
			values, ok := requestBody["values"].([]interface{})
			if !ok {
				return nil, nil // Not an error test, continue to next responder
			}

			// Check if any value is the error test case
			for _, val := range values {
				value, ok := val.(map[string]any)
				if !ok {
					continue
				}

				if value["name"] == "EnableGroupCreation" && value["value"] == "error" {
					errorResponse := map[string]any{
						"error": map[string]any{
							"code":    "InternalServerError",
							"message": "Internal server error occurred",
						},
					}
					return httpmock.NewJsonResponse(500, errorResponse)
				}
			}

			// If not an error, continue with normal processing
			return nil, nil
		})

	// Register POST /settings (Create)
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/settings",
		func(req *http.Request) (*http.Response, error) {
			// Read the request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close() // Close the body after reading

			var requestBody map[string]any
			_ = json.Unmarshal(body, &requestBody)

			// Check if this is a minimal or maximal request by looking at the values
			values := requestBody["values"].([]interface{})
			isMinimal := true

			// Check if this is a maximal request by looking for specific values
			for _, val := range values {
				value := val.(map[string]any)
				if value["name"] == "AllowGuestsToBeGroupOwner" ||
					value["name"] == "ClassificationDescriptions" ||
					value["name"] == "CustomBlockedWordsList" {
					isMinimal = false
					break
				}
			}

			var response map[string]any

			if isMinimal {
				// Minimal response with exactly the values from the minimal config
				response = map[string]any{
					"@odata.type": "#microsoft.graph.directorySetting",
					"id":          "test-tenant-setting-id",
					"displayName": "Group.Unified",
					"templateId":  "62375ab9-6b52-47ed-826b-58e47e0e304b",
					"values": []map[string]any{
						{
							"name":  "EnableGroupCreation",
							"value": "true",
						},
						{
							"name":  "AllowGuestsToAccessGroups",
							"value": "true",
						},
						{
							"name":  "AllowToAddGuests",
							"value": "true",
						},
					},
				}
			} else {
				// Maximal response with all the values from the maximal config
				response = map[string]any{
					"@odata.type": "#microsoft.graph.directorySetting",
					"id":          "test-tenant-setting-id",
					"displayName": "Group.Unified",
					"templateId":  "62375ab9-6b52-47ed-826b-58e47e0e304b",
					"values": []map[string]any{
						{
							"name":  "EnableGroupCreation",
							"value": "false",
						},
						{
							"name":  "GroupCreationAllowedGroupId",
							"value": "12345678-1234-1234-1234-123456789012",
						},
						{
							"name":  "PrefixSuffixNamingRequirement",
							"value": "[Contoso]-[GroupName]",
						},
						{
							"name":  "CustomBlockedWordsList",
							"value": "CEO,Legal,HR",
						},
						{
							"name":  "EnableMSStandardBlockedWords",
							"value": "true",
						},
						{
							"name":  "ClassificationList",
							"value": "Public,Internal,Confidential",
						},
						{
							"name":  "ClassificationDescriptions",
							"value": "Public:Public data,Internal:Internal data,Confidential:Confidential data",
						},
						{
							"name":  "DefaultClassification",
							"value": "Internal",
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
							"name":  "GuestUsageGuidelinesUrl",
							"value": "https://contoso.com/guestpolicies",
						},
						{
							"name":  "UsageGuidelinesUrl",
							"value": "https://contoso.com/groupguidelines",
						},
						{
							"name":  "EnableMIPLabels",
							"value": "true",
						},
						{
							"name":  "NewUnifiedGroupWritebackDefault",
							"value": "true",
						},
					},
				}
			}

			// Store in mock state
			settingId := "test-tenant-setting-id"
			mockState.Lock()
			mockState.settings[settingId] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register GET /settings/{setting-id} (Read)
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/settings/.*$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract the setting ID from the URL
			path := req.URL.Path
			settingId := path[len("/beta/settings/"):]

			mockState.Lock()
			setting, exists := mockState.settings[settingId]
			mockState.Unlock()

			if !exists {
				// Return a 404 if the setting doesn't exist in state
				errorResponse := map[string]any{
					"error": map[string]any{
						"code":    "NotFound",
						"message": "The specified object was not found",
					},
				}
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Return the setting if it exists
			return httpmock.NewJsonResponse(200, setting)
		})

	// Register PATCH /settings/{setting-id} (Update)
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/settings/.*$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract the setting ID from the URL
			path := req.URL.Path
			settingId := path[len("/beta/settings/"):]

			mockState.Lock()
			setting, exists := mockState.settings[settingId]
			mockState.Unlock()

			if !exists {
				errorResponse := map[string]any{
					"error": map[string]any{
						"code":    "ResourceNotFound",
						"message": "Setting not found",
					},
				}
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Read the request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close() // Close the body after reading

			var updateData map[string]any
			_ = json.Unmarshal(body, &updateData)

			// Update the setting
			mockState.Lock()
			for k, v := range updateData {
				setting[k] = v
			}
			mockState.settings[settingId] = setting
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, setting)
		})

	// Register DELETE /settings/{setting-id} (Delete)
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/settings/.*$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract the setting ID from the URL
			path := req.URL.Path
			settingId := path[len("/beta/settings/"):]

			mockState.Lock()
			// Ensure the resource is removed from the state
			delete(mockState.settings, settingId)
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// SetupTenantWideGroupSettingsMocks initializes httpmock and sets up all responders
func SetupTenantWideGroupSettingsMocks() {
	MockTenantWideGroupSettingsResponders()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func RegisterErrorMocks() {
	// Register error response for settings creation with "error" value
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/settings",
		func(req *http.Request) (*http.Response, error) {
			// Read the request body
			body, _ := io.ReadAll(req.Body)
			req.Body.Close() // Close the body after reading

			var requestBody map[string]any
			_ = json.Unmarshal(body, &requestBody)

			// Check if this is an error test by looking at the values
			values, ok := requestBody["values"].([]interface{})
			if !ok {
				return nil, nil // Not an error test, continue to next responder
			}

			// Check if any value is the error test case
			for _, val := range values {
				value, ok := val.(map[string]any)
				if !ok {
					continue
				}

				if value["name"] == "EnableGroupCreation" && value["value"] == "error" {
					errorResponse := map[string]any{
						"error": map[string]any{
							"code":    "InternalServerError",
							"message": "Internal server error occurred",
						},
					}
					return httpmock.NewJsonResponse(500, errorResponse)
				}
			}

			// If not an error, continue with normal processing
			return nil, nil
		})
}

// TeardownTenantWideGroupSettingsMocks deactivates httpmock
func TeardownTenantWideGroupSettingsMocks() {
	httpmock.DeactivateAndReset()
}

// MockTenantWideGroupSettingsError sets up an error response for testing error handling
func MockTenantWideGroupSettingsError(method, url string, statusCode int, errorCode, errorMessage string) {
	httpmock.RegisterResponder(method, url,
		func(req *http.Request) (*http.Response, error) {
			errorResponse := map[string]any{
				"error": map[string]any{
					"code":    errorCode,
					"message": errorMessage,
				},
			}
			return httpmock.NewJsonResponse(statusCode, errorResponse)
		})
}

// MockTenantWideGroupSettingsSuccess sets up a success response for testing
func MockTenantWideGroupSettingsSuccess(method, url string, statusCode int, responseData map[string]any) {
	httpmock.RegisterResponder(method, url,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(statusCode, responseData)
		})
}
