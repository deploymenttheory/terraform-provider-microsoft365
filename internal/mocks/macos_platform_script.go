// Extension to your existing mocks package for macOS Platform Script support

package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"
)

// RegisterMacOSPlatformScriptMocks registers mock handlers for macOS platform script operations
func (m *Mocks) RegisterMacOSPlatformScriptMocks() {
	// Register authentication mocks
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}))

	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"tenant_discovery_endpoint": "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/v2.0/.well-known/openid-configuration",
		}))

	// GET Read - Basic/Default script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000001",
				"displayName":                 "Test macOS Script",
				"description":                 "Test description",
				"runAsAccount":                "system",
				"fileName":                    "test-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn", // Base64 encoded
				"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime":        "2023-11-01T10:30:00.0000000Z",
				"roleScopeTagIds":             []string{"0"},
				"blockExecutionNotifications": true,
				"executionFrequency":          "P1D",
				"retryCount":                  3,
				"assignments": []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
						},
					},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Updated script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000003",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000003",
				"displayName":                 "Updated macOS Script",
				"description":                 "Updated description",
				"runAsAccount":                "user",
				"fileName":                    "updated-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gVXBkYXRlZCBXb3JsZCc=", // Base64 encoded
				"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime":        "2023-11-02T15:45:00.0000000Z",
				"roleScopeTagIds":             []string{"0"},
				"blockExecutionNotifications": false,
				"executionFrequency":          "P1W", // Exactly match what's in the test config
				"retryCount":                  5,
				"assignments": []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000004",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
						},
					},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Minimal script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000005",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                   "00000000-0000-0000-0000-000000000005",
				"displayName":          "Minimal macOS Script",
				"runAsAccount":         "system",
				"fileName":             "minimal-script.sh",
				"scriptContent":        "IyEvYmluL2Jhc2gKZWNobyAnTWluaW1hbCBTY3JpcHQn", // Base64 encoded
				"createdDateTime":      "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime": "2023-11-01T10:30:00.0000000Z",
				"roleScopeTagIds":      []string{"0"},
				"assignments":          []map[string]interface{}{},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Maximal script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000006",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000006",
				"displayName":                 "Maximal macOS Script",
				"description":                 "This is a comprehensive script with all fields populated",
				"runAsAccount":                "user",
				"fileName":                    "maximal-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnTWF4aW1hbCBTY3JpcHQgQ29uZmlndXJhdGlvbic=", // Base64 encoded
				"createdDateTime":             "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime":        "2023-11-01T10:30:00.0000000Z",
				"roleScopeTagIds":             []string{"0", "1"},
				"blockExecutionNotifications": true,
				"executionFrequency":          "P4W",
				"retryCount":                  10,
				"assignments": []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000007",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
						},
					},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Group Assignment script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000008",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                   "00000000-0000-0000-0000-000000000008",
				"displayName":          "Group Assignment Script",
				"description":          "Script with group assignments",
				"runAsAccount":         "system",
				"fileName":             "group-script.sh",
				"scriptContent":        "IyEvYmluL2Jhc2gKZWNobyAnR3JvdXAgQXNzaWdubWVudCBTY3JpcHQn", // Base64 encoded
				"createdDateTime":      "2023-11-01T10:30:00.0000000Z",
				"lastModifiedDateTime": "2023-11-01T10:30:00.0000000Z",
				"roleScopeTagIds":      []string{"0"},
				"assignments": []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000009",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.groupAssignmentTarget",
							"groupId":     "11111111-1111-1111-1111-111111111111",
						},
					},
					{
						"id": "00000000-0000-0000-0000-000000000010",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
							"groupId":     "22222222-2222-2222-2222-222222222222",
						},
					},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// POST Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Generate a response based on the request
			id := "00000000-0000-0000-0000-000000000001" // Default ID
			displayName := requestBody["displayName"].(string)

			// Assign different IDs based on display name for different test cases
			if displayName == "Minimal macOS Script" {
				id = "00000000-0000-0000-0000-000000000005"
			} else if displayName == "Maximal macOS Script" {
				id = "00000000-0000-0000-0000-000000000006"
			} else if displayName == "Group Assignment Script" {
				id = "00000000-0000-0000-0000-000000000008"
			} else if displayName == "Updated macOS Script" {
				id = "00000000-0000-0000-0000-000000000003"
			}

			// Create response with the same fields as the request plus ID and dates
			response := requestBody
			response["id"] = id
			response["createdDateTime"] = "2023-11-01T10:30:00.0000000Z"
			response["lastModifiedDateTime"] = "2023-11-01T10:30:00.0000000Z"

			if _, ok := response["roleScopeTagIds"]; !ok {
				response["roleScopeTagIds"] = []string{"0"}
			}

			return httpmock.NewJsonResponse(201, response)
		})

	// PATCH Update
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[0-9a-f-]+`),
		func(req *http.Request) (*http.Response, error) {
			// Parse the request body to get the updated fields
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Return success with no content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// DELETE
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[0-9a-f-]+`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// POST Assign
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[0-9a-f-]+/assign`),
		httpmock.NewJsonResponderOrPanic(204, nil))

	// GET Assignments
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[0-9a-f-]+/assignments`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the script ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2]

			var assignments []map[string]interface{}

			// Return different assignments based on script ID
			switch scriptId {
			case "00000000-0000-0000-0000-000000000001":
				// Basic script - all users assignment
				assignments = []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
						},
					},
				}
			case "00000000-0000-0000-0000-000000000003":
				// Updated script - all devices assignment
				assignments = []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000004",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
						},
					},
				}
			case "00000000-0000-0000-0000-000000000005":
				// Minimal script - no assignments
				assignments = []map[string]interface{}{}
			case "00000000-0000-0000-0000-000000000006":
				// Maximal script - all devices assignment
				assignments = []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000007",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
						},
					},
				}
			case "00000000-0000-0000-0000-000000000008":
				// Group assignment script - include and exclude groups
				assignments = []map[string]interface{}{
					{
						"id": "00000000-0000-0000-0000-000000000009",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.groupAssignmentTarget",
							"groupId":     "11111111-1111-1111-1111-111111111111",
						},
					},
					{
						"id": "00000000-0000-0000-0000-000000000010",
						"target": map[string]interface{}{
							"@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
							"groupId":     "22222222-2222-2222-2222-222222222222",
						},
					},
				}
			default:
				assignments = []map[string]interface{}{}
			}

			response := map[string]interface{}{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts('%s')/assignments", scriptId),
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// RegisterMacOSPlatformScriptErrorMocks registers mock handlers that return errors for testing error handling
func (m *Mocks) RegisterMacOSPlatformScriptErrorMocks() {
	// POST Create with error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewJsonResponderOrPanic(403, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "Forbidden",
				"message": "Access denied. You do not have permission to perform this action or access this resource.",
			},
		}))
}

// Update your main Activate method to include macOS script mocks
func (m *Mocks) ActivateWithMacOSScripts() {
	m.Activate()
	m.RegisterMacOSPlatformScriptMocks()
}

// ActivateWithMacOSScriptErrors activates mocks with error responses for macOS platform scripts
func (m *Mocks) ActivateWithMacOSScriptErrors() {
	m.Activate()
	m.RegisterMacOSPlatformScriptErrorMocks()
}
