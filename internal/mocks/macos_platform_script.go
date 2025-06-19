// Extension to your existing mocks package for macOS Platform Script support

package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

// Add this method to your existing Mocks struct
func (m *Mocks) RegisterMacOSPlatformScriptMocks() {
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
}

// Update your main Activate method to include macOS script mocks
func (m *Mocks) ActivateWithMacOSScripts() {
	m.Activate()
	m.RegisterMacOSPlatformScriptMocks()
}
