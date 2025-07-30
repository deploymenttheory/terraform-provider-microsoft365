package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	scriptUpdated map[string]bool
	// Store the actual script data to ensure consistency
	scriptData map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.scriptUpdated = make(map[string]bool)
	mockState.scriptData = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// MacOSPlatformScriptMock provides mock responses for macOS platform script operations
type MacOSPlatformScriptMock struct{}

// RegisterMocks registers HTTP mock responses for macOS platform script operations
func (m *MacOSPlatformScriptMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.scriptUpdated = make(map[string]bool)
	mockState.scriptData = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Initialize base script data
	baseScriptID := "00000000-0000-0000-0000-000000000001"
	baseScriptData := map[string]interface{}{
		"id":                          baseScriptID,
		"displayName":                 "Test macOS Script",
		"description":                 "Test description",
		"runAsAccount":                "system",
		"fileName":                    "test-script.sh",
		"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn", // Base64 encoded
		"createdDateTime":             "2023-11-01T10:30:00Z",
		"lastModifiedDateTime":        "2023-11-01T10:30:00Z",
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

	mockState.Lock()
	mockState.scriptData[baseScriptID] = baseScriptData
	mockState.Unlock()

	// Register GET Read - Basic/Default script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			isUpdated := mockState.scriptUpdated["00000000-0000-0000-0000-000000000001"]
			scriptData := mockState.scriptData["00000000-0000-0000-0000-000000000001"]
			mockState.Unlock()

			if isUpdated {
				// Return the updated script data (no need to override assignments)
				return httpmock.NewJsonResponse(200, scriptData)
			}

			// Return the original script data
			return httpmock.NewJsonResponse(200, scriptData)
		})

	// Register GET Read - Updated script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000003",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000003",
				"displayName":                 "Updated macOS Script",
				"description":                 "Updated description",
				"runAsAccount":                "user",
				"fileName":                    "updated-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gVXBkYXRlZCBXb3JsZCc=", // Base64 encoded
				"createdDateTime":             "2023-11-01T10:30:00Z",
				"lastModifiedDateTime":        "2023-11-01T10:30:00Z",
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

	// Register GET Read - Minimal script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000005",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                   "00000000-0000-0000-0000-000000000005",
				"displayName":          "Minimal macOS Script",
				"runAsAccount":         "system",
				"fileName":             "minimal-script.sh",
				"scriptContent":        "IyEvYmluL2Jhc2gKZWNobyAnTWluaW1hbCBTY3JpcHQn", // Base64 encoded
				"createdDateTime":      "2023-11-01T10:30:00Z",
				"lastModifiedDateTime": "2023-11-01T10:30:00Z",
				"roleScopeTagIds":      []string{"0"},
				"assignments":          []map[string]interface{}{},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET Read - Maximal script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000006",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                          "00000000-0000-0000-0000-000000000006",
				"displayName":                 "Maximal macOS Script",
				"description":                 "This is a comprehensive script with all fields populated",
				"runAsAccount":                "user",
				"fileName":                    "maximal-script.sh",
				"scriptContent":               "IyEvYmluL2Jhc2gKZWNobyAnTWF4aW1hbCBTY3JpcHQgQ29uZmlndXJhdGlvbic=", // Base64 encoded
				"createdDateTime":             "2023-11-01T10:30:00Z",
				"lastModifiedDateTime":        "2023-11-01T10:30:00Z",
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

	// Register GET Read - Group Assignment script
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000008",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                   "00000000-0000-0000-0000-000000000008",
				"displayName":          "Group Assignment Script",
				"description":          "Script with group assignments",
				"runAsAccount":         "system",
				"fileName":             "group-script.sh",
				"scriptContent":        "IyEvYmluL2Jhc2gKZWNobyAnR3JvdXAgQXNzaWdubWVudCBTY3JpcHQn", // Base64 encoded
				"createdDateTime":      "2023-11-01T10:30:00Z",
				"lastModifiedDateTime": "2023-11-01T10:30:00Z",
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

	// Register POST Create
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
			response := make(map[string]interface{})
			for k, v := range requestBody {
				response[k] = v
			}
			response["id"] = id
			response["createdDateTime"] = "2023-11-01T10:30:00Z"
			response["lastModifiedDateTime"] = "2023-11-01T10:30:00Z"

			if _, ok := response["roleScopeTagIds"]; !ok {
				response["roleScopeTagIds"] = []string{"0"}
			}

			// Store in our state
			mockState.Lock()
			mockState.scriptData[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH Update
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)`),
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Debug log the request body
			fmt.Printf("PATCH request body for ID %s: %+v\n", id, requestBody)

			// Update the script data in our state
			mockState.Lock()
			if currentData, exists := mockState.scriptData[id]; exists {
				// Update fields from the request
				for k, v := range requestBody {
					currentData[k] = v
				}

				// Mark as updated
				mockState.scriptUpdated[id] = true

				// Update the lastModifiedDateTime
				currentData["lastModifiedDateTime"] = "2023-11-01T10:30:00Z"
			}
			mockState.Unlock()

			// For update, return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register DELETE
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)`),
		func(req *http.Request) (*http.Response, error) {
			// For delete, return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET List
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts",
				"value": []map[string]interface{}{
					{
						"id":                   "00000000-0000-0000-0000-000000000001",
						"displayName":          "Test macOS Script",
						"description":          "Test description",
						"runAsAccount":         "system",
						"fileName":             "test-script.sh",
						"createdDateTime":      "2023-11-01T10:30:00Z",
						"lastModifiedDateTime": "2023-11-01T10:30:00Z",
						"roleScopeTagIds":      []string{"0"},
					},
					{
						"id":                   "00000000-0000-0000-0000-000000000003",
						"displayName":          "Updated macOS Script",
						"description":          "Updated description",
						"runAsAccount":         "user",
						"fileName":             "updated-script.sh",
						"createdDateTime":      "2023-11-01T10:30:00Z",
						"lastModifiedDateTime": "2023-11-01T10:30:00Z",
						"roleScopeTagIds":      []string{"0"},
					},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST Assign for all script IDs
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)/assign`),
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Debug log the request body
			fmt.Printf("Assign request body for ID %s: %+v\n", id, requestBody)

			// Update assignments in our state
			mockState.Lock()
			if currentData, exists := mockState.scriptData[id]; exists {
				assignments := make([]map[string]interface{}, 0)

				// Process deviceManagementScriptAssignments from request
				if deviceManagementScriptAssignments, ok := requestBody["deviceManagementScriptAssignments"].([]interface{}); ok {
					for _, assignment := range deviceManagementScriptAssignments {
						if assignmentMap, ok := assignment.(map[string]interface{}); ok {
							if target, ok := assignmentMap["target"].(map[string]interface{}); ok {
								assignmentID := fmt.Sprintf("%s_assign_%s", id, target["@odata.type"])
								assignments = append(assignments, map[string]interface{}{
									"id":     assignmentID,
									"target": target,
								})
							}
						}
					}
				}

				// Store the assignments in the script data
				currentData["assignments"] = assignments

				// Mark the script as updated
				mockState.scriptUpdated[id] = true
			}
			mockState.Unlock()

			// For assign, return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET Assignments for all script IDs
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)/assignments`),
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL to determine which assignments to return
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			mockState.Lock()
			scriptData := mockState.scriptData[id]
			mockState.Unlock()

			var assignments []map[string]interface{}

			// Use stored assignments if available, regardless of update status
			if scriptData != nil {
				if scriptAssignments, ok := scriptData["assignments"].([]map[string]interface{}); ok {
					assignments = scriptAssignments
				} else if scriptAssignmentsInterface, ok := scriptData["assignments"].([]interface{}); ok {
					// Convert from []interface{} to []map[string]interface{}
					for _, a := range scriptAssignmentsInterface {
						if aMap, ok := a.(map[string]interface{}); ok {
							assignments = append(assignments, aMap)
						}
					}
				}
			}

			// If no assignments found in state, use defaults based on ID
			if len(assignments) == 0 {
				switch id {
				case "00000000-0000-0000-0000-000000000003":
					// Updated script - all devices (from the update config)
					assignments = []map[string]interface{}{
						{
							"id": fmt.Sprintf("%s_adadadad-808e-44e2-905a-0b7873a8a531", id),
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
							},
						},
					}
				case "00000000-0000-0000-0000-000000000006":
					// Maximal script - all devices
					assignments = []map[string]interface{}{
						{
							"id": fmt.Sprintf("%s_adadadad-808e-44e2-905a-0b7873a8a531", id),
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
							},
						},
					}
				case "00000000-0000-0000-0000-000000000008":
					// Group assignment script
					assignments = []map[string]interface{}{
						{
							"id": fmt.Sprintf("%s_adadadad-808e-44e2-905a-0b7873a8a531", id),
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.groupAssignmentTarget",
								"groupId":     "11111111-1111-1111-1111-111111111111",
							},
						},
						{
							"id": fmt.Sprintf("%s_bbbbbbbb-808e-44e2-905a-0b7873a8a531", id),
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
								"groupId":     "22222222-2222-2222-2222-222222222222",
							},
						},
					}
				default:
					// Default for basic script - all users
					if id == "00000000-0000-0000-0000-000000000001" {
						assignments = []map[string]interface{}{
							{
								"id": fmt.Sprintf("%s_adadadad-808e-44e2-905a-0b7873a8a531", id),
								"target": map[string]interface{}{
									"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
								},
							},
						}
					} else {
						// Empty assignments for other scripts
						assignments = []map[string]interface{}{}
					}
				}
			}

			responseBody := map[string]interface{}{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts('%s')/assignments", id),
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, responseBody)
		})
}

// RegisterErrorMocks registers HTTP mock responses that return errors
func (m *MacOSPlatformScriptMock) RegisterErrorMocks() {
	// Register error responses for each operation

	// GET - Read error
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/00000000-0000-0000-0000-000000000001",
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// DELETE - Delete error
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// GET - List error
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// POST - Assign error
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)/assign`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// GET - Assignments error
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/([0-9a-f-]+)/assignments`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
}
