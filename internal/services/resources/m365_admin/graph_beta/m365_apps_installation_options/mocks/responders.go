package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
)

// M365AppsInstallationOptionsMock contains methods to register mocks for the M365 Apps Installation Options resource
type M365AppsInstallationOptionsMock struct{}

// RegisterMocks registers HTTP mocks for the M365 Apps Installation Options resource
func (m *M365AppsInstallationOptionsMock) RegisterMocks() {
	// Register GET mock for the resource - beta endpoint
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/admin/microsoft365Apps/installationOptions",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{
				"id": "microsoft365AppsInstallationOptions",
				"updateChannel": "current",
				"appsForWindows": {
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true
				},
				"appsForMac": {
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true
				}
			}`,
		),
	)

	// Register GET mock for the resource - v1.0 endpoint
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/v1.0/admin/microsoft365Apps/installationOptions",
		httpmock.NewStringResponder(
			http.StatusOK,
			`{
				"id": "microsoft365AppsInstallationOptions",
				"updateChannel": "current",
				"appsForWindows": {
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true
				},
				"appsForMac": {
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true
				}
			}`,
		),
	)

	// Register PATCH mock for updating the resource - beta endpoint
	httpmock.RegisterResponder(
		"PATCH",
		"https://graph.microsoft.com/beta/admin/microsoft365Apps/installationOptions",
		func(req *http.Request) (*http.Response, error) {
			// Parse the request body to get the updated values
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error": {"message": "Invalid request body"}}`), nil
			}

			// Create a response with the updated values
			responseBody := map[string]any{
				"id": "microsoft365AppsInstallationOptions",
			}

			// Copy update_channel if provided
			if updateChannel, ok := requestBody["updateChannel"].(string); ok {
				responseBody["updateChannel"] = updateChannel
			} else {
				responseBody["updateChannel"] = "current"
			}

			// Copy apps_for_windows if provided
			if appsForWindows, ok := requestBody["appsForWindows"].(map[string]any); ok {
				responseBody["appsForWindows"] = appsForWindows
			} else {
				responseBody["appsForWindows"] = map[string]any{
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true,
				}
			}

			// Copy apps_for_mac if provided
			if appsForMac, ok := requestBody["appsForMac"].(map[string]any); ok {
				responseBody["appsForMac"] = appsForMac
			} else {
				responseBody["appsForMac"] = map[string]any{
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true,
				}
			}

			// Convert the response to JSON
			responseJSON, err := json.Marshal(responseBody)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error": {"message": "Failed to create response"}}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, string(responseJSON)), nil
		},
	)

	// Register PATCH mock for updating the resource - v1.0 endpoint
	httpmock.RegisterResponder(
		"PATCH",
		"https://graph.microsoft.com/v1.0/admin/microsoft365Apps/installationOptions",
		func(req *http.Request) (*http.Response, error) {
			// Parse the request body to get the updated values
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error": {"message": "Invalid request body"}}`), nil
			}

			// Create a response with the updated values
			responseBody := map[string]any{
				"id": "microsoft365AppsInstallationOptions",
			}

			// Copy update_channel if provided
			if updateChannel, ok := requestBody["updateChannel"].(string); ok {
				responseBody["updateChannel"] = updateChannel
			} else {
				responseBody["updateChannel"] = "current"
			}

			// Copy apps_for_windows if provided
			if appsForWindows, ok := requestBody["appsForWindows"].(map[string]any); ok {
				responseBody["appsForWindows"] = appsForWindows
			} else {
				responseBody["appsForWindows"] = map[string]any{
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true,
				}
			}

			// Copy apps_for_mac if provided
			if appsForMac, ok := requestBody["appsForMac"].(map[string]any); ok {
				responseBody["appsForMac"] = appsForMac
			} else {
				responseBody["appsForMac"] = map[string]any{
					"isMicrosoft365AppsEnabled": true,
					"isSkypeForBusinessEnabled": true,
				}
			}

			// Convert the response to JSON
			responseJSON, err := json.Marshal(responseBody)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error": {"message": "Failed to create response"}}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, string(responseJSON)), nil
		},
	)
}

// RegisterErrorMocks registers HTTP mocks that return errors
func (m *M365AppsInstallationOptionsMock) RegisterErrorMocks() {
	// Register GET mock that returns an error
	httpmock.RegisterRegexpResponder(
		"GET",
		regexp.MustCompile(`https://graph.microsoft.com/beta/admin/microsoft365Apps/installationOptions`),
		httpmock.NewStringResponder(
			http.StatusForbidden,
			`{
				"error": {
					"code": "Forbidden",
					"message": "Access denied. You do not have permission to perform this action or access this resource.",
					"innerError": {
						"date": "2023-01-01T12:00:00",
						"request-id": "00000000-0000-0000-0000-000000000000"
					}
				}
			}`,
		),
	)

	// Register PATCH mock that returns an error
	httpmock.RegisterRegexpResponder(
		"PATCH",
		regexp.MustCompile(`https://graph.microsoft.com/beta/admin/microsoft365Apps/installationOptions`),
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{
				"error": {
					"code": "BadRequest",
					"message": "Invalid value specified for property 'updateChannel' of resource 'microsoft365AppsInstallationOptions'.",
					"innerError": {
						"date": "2023-01-01T12:00:00",
						"request-id": "00000000-0000-0000-0000-000000000000"
					}
				}
			}`,
		),
	)
}

// GetMockResponse returns a mock response for the given ID
func (m *M365AppsInstallationOptionsMock) GetMockResponse(id string) string {
	if id == "error" {
		return `{
			"error": {
				"code": "NotFound",
				"message": "Resource not found",
				"innerError": {
					"date": "2023-01-01T12:00:00",
					"request-id": "00000000-0000-0000-0000-000000000000"
				}
			}
		}`
	}

	return fmt.Sprintf(`{
		"id": "microsoft365AppsInstallationOptions",
		"updateChannel": "current",
		"appsForWindows": {
			"isMicrosoft365AppsEnabled": true,
			"isSkypeForBusinessEnabled": true
		},
		"appsForMac": {
			"isMicrosoft365AppsEnabled": true,
			"isSkypeForBusinessEnabled": true
		}
	}`)
}
