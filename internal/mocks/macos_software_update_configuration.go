package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

const (
	macOSSoftwareUpdateConfigurationID = "00000000-0000-0000-0000-000000000001"
)

// RegisterMacOSSoftwareUpdateConfigurationMocks registers HTTP mocks for macOS software update configuration operations
func (m *Mocks) RegisterMacOSSoftwareUpdateConfigurationMocks() {
	// Create - POST
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Check if this is a macOS software update configuration
			odataType, ok := requestBody["@odata.type"].(string)
			if !ok || odataType != "#microsoft.graph.macOSSoftwareUpdateConfiguration" {
				return httpmock.NewStringResponse(400, "Invalid @odata.type"), nil
			}

			// Create response with ID
			requestBody["id"] = macOSSoftwareUpdateConfigurationID
			responseBody, err := json.Marshal(requestBody)
			if err != nil {
				return httpmock.NewStringResponse(500, "Error creating response"), nil
			}

			return httpmock.NewStringResponse(201, string(responseBody)), nil
		})

	// Read - GET by ID
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			// Basic macOS software update configuration
			responseBody := map[string]interface{}{
				"@odata.type":                        "#microsoft.graph.macOSSoftwareUpdateConfiguration",
				"id":                                 macOSSoftwareUpdateConfigurationID,
				"displayName":                        "Test macOS Software Update Configuration",
				"description":                        "Test description",
				"roleScopeTagIds":                    []string{"0"},
				"criticalUpdateBehavior":             "default",
				"configDataUpdateBehavior":           "default",
				"firmwareUpdateBehavior":             "default",
				"allOtherUpdateBehavior":             "default",
				"updateScheduleType":                 "alwaysUpdate",
				"updateTimeWindowUtcOffsetInMinutes": 0,
				"maxUserDeferralsCount":              0,
				"priority":                           "low",
				"customUpdateTimeWindows":            []interface{}{},
			}

			return httpmock.NewJsonResponse(200, responseBody)
		})

	// Update - PATCH
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// For update, just return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Delete - DELETE
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// List - GET
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			// Check if the request is filtering for macOS software update configurations
			if strings.Contains(req.URL.RawQuery, "microsoft.graph.macOSSoftwareUpdateConfiguration") {
				responseBody := map[string]interface{}{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
					"value": []map[string]interface{}{
						{
							"@odata.type":                        "#microsoft.graph.macOSSoftwareUpdateConfiguration",
							"id":                                 macOSSoftwareUpdateConfigurationID,
							"displayName":                        "Test macOS Software Update Configuration",
							"description":                        "Test description",
							"roleScopeTagIds":                    []string{"0"},
							"criticalUpdateBehavior":             "default",
							"configDataUpdateBehavior":           "default",
							"firmwareUpdateBehavior":             "default",
							"allOtherUpdateBehavior":             "default",
							"updateScheduleType":                 "alwaysUpdate",
							"updateTimeWindowUtcOffsetInMinutes": 0,
							"maxUserDeferralsCount":              0,
							"priority":                           "low",
							"customUpdateTimeWindows":            []interface{}{},
						},
					},
				}
				return httpmock.NewJsonResponse(200, responseBody)
			}

			// Return empty list for other queries
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          []interface{}{},
			})
		})

	// Assignments - GET
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s/assignments", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			responseBody := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations('00000000-0000-0000-0000-000000000001')/assignments",
				"value": []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000001_adadadad-808e-44e2-905a-0b7873a8a531",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
						},
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseBody)
		})

	// Create Assignment - POST
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s/assignments", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// For assignment creation, return a 201 Created
			responseBody := map[string]interface{}{
				"id": "00000000-0000-0000-0000-000000000001_adadadad-808e-44e2-905a-0b7873a8a531",
				"target": map[string]interface{}{
					"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
				},
			}
			return httpmock.NewJsonResponse(201, responseBody)
		})

	// Delete Assignment - DELETE
	httpmock.RegisterResponder("DELETE", "=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/.*/assignments/.*$",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterMacOSSoftwareUpdateConfigurationErrorMocks registers HTTP mocks that return errors for macOS software update configuration operations
func (m *Mocks) RegisterMacOSSoftwareUpdateConfigurationErrorMocks() {
	// Create - POST with error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(403, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	// Read - GET by ID with error
	httpmock.RegisterResponder("GET", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(403, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	// Update - PATCH with error
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(403, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})

	// Delete - DELETE with error
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/%s", macOSSoftwareUpdateConfigurationID),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(403, `{"error":{"code":"Forbidden","message":"Access denied"}}`), nil
		})
}
