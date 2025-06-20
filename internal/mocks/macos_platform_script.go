// Extension to your existing mocks package for macOS Platform Script support

package mocks

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

// Add this method to your existing Mocks struct
func (m *Mocks) RegisterMacOSPlatformScriptMocks() {
	// Register authentication mocks
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}))

	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"tenant_discovery_endpoint": "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/v2.0/.well-known/openid-configuration",
		}))

	// Basic CRUD operations for macOS Platform Scripts

	// POST Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":           "00000000-0000-0000-0000-000000000001",
			"displayName":  "Test macOS Script",
			"runAsAccount": "system",
			"fileName":     "test-script.sh",
		}))

	// POST Assign
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"value": "Assignment completed successfully",
		}))

	// GET Read
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":           "00000000-0000-0000-0000-000000000001",
				"displayName":  "Test macOS Script",
				"runAsAccount": "system",
				"fileName":     "test-script.sh",
			}

			if req.URL.Query().Get("$expand") == "assignments" {
				response["assignments"] = map[string]interface{}{
					"value": []map[string]interface{}{
						{
							"id": "assignment1",
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
							},
						},
					},
				}
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// PATCH Update
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"id":           "00000000-0000-0000-0000-000000000001",
			"displayName":  "Updated macOS Script",
			"runAsAccount": "user",
			"fileName":     "updated-script.sh",
		}))

	// DELETE
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewStringResponder(204, ""))

	// Register a catch-all responder for unexpected API calls
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		fmt.Printf("⚠️ Unexpected API call: %s %s\n", req.Method, req.URL.String())
		return httpmock.NewStringResponse(404, "Not Found - Mock not registered for this request"), nil
	})
}

// RegisterMacOSPlatformScriptErrorMocks registers error responses for macOS platform script operations
func (m *Mocks) RegisterMacOSPlatformScriptErrorMocks() {
	// Error on Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		httpmock.NewJsonResponderOrPanic(403, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "Forbidden",
				"message": "Access denied",
			},
		}))

	// Error on Update
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewJsonResponderOrPanic(400, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "BadRequest",
				"message": "Invalid property value",
			},
		}))

	// Error on Delete
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		httpmock.NewJsonResponderOrPanic(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Resource not found",
			},
		}))

	// Error on Assign
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001/assign",
		httpmock.NewJsonResponderOrPanic(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": "An unexpected error occurred",
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
