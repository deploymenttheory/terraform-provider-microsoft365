package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"
)

// Global variable to track if update has happened
var updateOccurred = false

// RegisterDeviceConfigurationAssignmentMocks registers HTTP mocks for device configuration assignment operations
func (m *Mocks) RegisterDeviceConfigurationAssignmentMocks() {
	// Reset the update flag
	updateOccurred = false

	// Register authentication mocks first
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

	// POST Create Assignment - handle different device configuration IDs
	httpmock.RegisterRegexpResponder("POST",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments`),
		func(req *http.Request) (*http.Response, error) {
			// Extract device configuration ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			var deviceConfigId string
			for i, part := range urlParts {
				if part == "deviceConfigurations" && i+1 < len(urlParts) {
					deviceConfigId = urlParts[i+1]
					break
				}
			}

			fmt.Printf("DEBUG MOCK: POST Create for deviceConfig=%s\n", deviceConfigId)

			// Handle error case
			if deviceConfigId == "error-config" {
				return httpmock.NewStringResponse(400, `{"error": {"code": "BadRequest", "message": "Bad Request"}}`), nil
			}

			// Parse the request body
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error": {"message": "Invalid request body"}}`), nil
			}

			// Generate response ID based on target type
			assignmentId := generateAssignmentId(requestBody)

			fmt.Printf("DEBUG MOCK: Generated assignment ID %s for request: %+v\n", assignmentId, requestBody)

			// Create response
			responseBody := map[string]interface{}{
				"id":     assignmentId,
				"target": requestBody["target"],
			}

			return httpmock.NewJsonResponse(201, responseBody)
		})

	// GET Read Assignment - handle different assignment types
	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		func(req *http.Request) (*http.Response, error) {
			// Extract IDs from URL
			urlParts := strings.Split(req.URL.Path, "/")
			var deviceConfigId, assignmentId string
			for i, part := range urlParts {
				if part == "deviceConfigurations" && i+1 < len(urlParts) {
					deviceConfigId = urlParts[i+1]
					if i+3 < len(urlParts) {
						assignmentId = urlParts[i+3]
					}
					break
				}
			}

			fmt.Printf("DEBUG MOCK: GET request for deviceConfig=%s, assignment=%s\n", deviceConfigId, assignmentId)

			// Handle error cases
			if deviceConfigId == "error-config" || assignmentId == "error-id" {
				return httpmock.NewStringResponse(404, `{"error": {"code": "NotFound", "message": "Assignment not found"}}`), nil
			}

			// Return appropriate response based on assignment ID
			response := getAssignmentResponse(assignmentId)
			fmt.Printf("DEBUG MOCK: Returning response for assignment %s: %+v\n", assignmentId, response)
			return httpmock.NewJsonResponse(200, response)
		})

	// Additional specific handlers for known assignment IDs to ensure they work
	specificAssignments := map[string]string{
		"00000000-0000-0000-0000-000000000001": "groupAssignment",
		"00000000-0000-0000-0000-000000000002": "allDevices",
		"00000000-0000-0000-0000-000000000003": "allLicensedUsers",
		"00000000-0000-0000-0000-000000000004": "exclusionGroup",
	}

	for assignmentId, targetType := range specificAssignments {
		httpmock.RegisterResponder("GET",
			fmt.Sprintf("https://graph.microsoft.com/v1.0/deviceManagement/deviceConfigurations/test-config-id/assignments/%s", assignmentId),
			func(assignmentId, targetType string) httpmock.Responder {
				return func(req *http.Request) (*http.Response, error) {
					fmt.Printf("DEBUG MOCK: Specific handler for assignment %s (%s)\n", assignmentId, targetType)
					response := getAssignmentResponse(assignmentId)
					return httpmock.NewJsonResponse(200, response)
				}
			}(assignmentId, targetType))
	}

	// Special handler for update scenario - overrides the generic handler for this specific ID
	httpmock.RegisterResponder("GET",
		"https://graph.microsoft.com/v1.0/deviceManagement/deviceConfigurations/test-config-id/assignments/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			// Check if update has occurred
			if updateOccurred {
				// After update, return allLicensedUsers target
				fmt.Printf("DEBUG MOCK: Returning updated (allLicensedUsers) data after PATCH for ID 00000000-0000-0000-0000-000000000001\n")
				return httpmock.NewJsonResponse(200, map[string]interface{}{
					"id": "00000000-0000-0000-0000-000000000001",
					"target": map[string]interface{}{
						"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
					},
				})
			}

			// Before update, return groupAssignment target
			fmt.Printf("DEBUG MOCK: Returning original (groupAssignment) data before PATCH for ID 00000000-0000-0000-0000-000000000001\n")
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"id": "00000000-0000-0000-0000-000000000001",
				"target": map[string]interface{}{
					"@odata.type": "#microsoft.graph.groupAssignmentTarget",
					"groupId":     "11111111-1111-1111-1111-111111111111",
				},
			})
		})

	// PATCH Update Assignment - handle different assignment types
	httpmock.RegisterRegexpResponder("PATCH",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		func(req *http.Request) (*http.Response, error) {
			// Extract IDs from URL
			urlParts := strings.Split(req.URL.Path, "/")
			var deviceConfigId, assignmentId string
			for i, part := range urlParts {
				if part == "deviceConfigurations" && i+1 < len(urlParts) {
					deviceConfigId = urlParts[i+1]
					if i+3 < len(urlParts) {
						assignmentId = urlParts[i+3]
					}
					break
				}
			}

			fmt.Printf("DEBUG MOCK: PATCH request for deviceConfig=%s, assignment=%s\n", deviceConfigId, assignmentId)

			// Handle error cases
			if deviceConfigId == "error-config" || assignmentId == "error-id" {
				return httpmock.NewStringResponse(400, `{"error": {"code": "BadRequest", "message": "Bad Request"}}`), nil
			}

			// Parse the request body
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error": {"message": "Invalid request body"}}`), nil
			}

			// Generate a new response based on the request target type
			var response map[string]interface{}

			if target, ok := requestBody["target"].(map[string]interface{}); ok {
				if odataType, ok := target["@odata.type"].(string); ok {
					fmt.Printf("DEBUG MOCK: PATCH request with target type: %s\n", odataType)

					// Set the update flag to true if this is the update we're testing for
					if deviceConfigId == "test-config-id" && assignmentId == "00000000-0000-0000-0000-000000000001" {
						updateOccurred = true
						fmt.Printf("DEBUG MOCK: Update flag set to true\n")
					}

					// Create response based on the target type in the request
					response = map[string]interface{}{
						"id":     assignmentId,
						"target": target,
					}
				} else {
					fmt.Printf("DEBUG MOCK: PATCH request missing @odata.type in target\n")
					return httpmock.NewStringResponse(400, `{"error": {"message": "Missing @odata.type in target"}}`), nil
				}
			} else {
				fmt.Printf("DEBUG MOCK: PATCH request missing target\n")
				return httpmock.NewStringResponse(400, `{"error": {"message": "Missing target in request body"}}`), nil
			}

			fmt.Printf("DEBUG MOCK: PATCH response: %+v\n", response)
			return httpmock.NewJsonResponse(200, response)
		})

	// DELETE Assignment
	httpmock.RegisterRegexpResponder("DELETE",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		func(req *http.Request) (*http.Response, error) {
			// Extract IDs from URL
			urlParts := strings.Split(req.URL.Path, "/")
			var deviceConfigId, assignmentId string
			for i, part := range urlParts {
				if part == "deviceConfigurations" && i+1 < len(urlParts) {
					deviceConfigId = urlParts[i+1]
					if i+3 < len(urlParts) {
						assignmentId = urlParts[i+3]
					}
					break
				}
			}

			// Handle error cases
			if deviceConfigId == "error-config" || assignmentId == "error-id" {
				return httpmock.NewStringResponse(403, `{"error": {"code": "Forbidden", "message": "Permission denied"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// GET List Assignments (optional, for completeness)
	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract device configuration ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			var deviceConfigId string
			for i, part := range urlParts {
				if part == "deviceConfigurations" && i+1 < len(urlParts) {
					deviceConfigId = urlParts[i+1]
					break
				}
			}

			// Handle error case
			if deviceConfigId == "error-config" {
				return httpmock.NewStringResponse(404, `{"error": {"code": "NotFound", "message": "Configuration not found"}}`), nil
			}

			// Return list of assignments
			assignments := []map[string]interface{}{
				{
					"id": "00000000-0000-0000-0000-000000000001",
					"target": map[string]interface{}{
						"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
					},
				},
			}

			response := map[string]interface{}{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/v1.0/$metadata#deviceManagement/deviceConfigurations('%s')/assignments", deviceConfigId),
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// generateAssignmentId generates an assignment ID based on the target type
func generateAssignmentId(requestBody map[string]interface{}) string {
	target, ok := requestBody["target"].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG MOCK: No target in request body\n")
		return "00000000-0000-0000-0000-000000000001"
	}

	odataType, ok := target["@odata.type"].(string)
	if !ok {
		fmt.Printf("DEBUG MOCK: No @odata.type in target\n")
		return "00000000-0000-0000-0000-000000000001"
	}

	fmt.Printf("DEBUG MOCK: Target @odata.type is: %s\n", odataType)

	switch odataType {
	case "#microsoft.graph.allDevicesAssignmentTarget":
		fmt.Printf("DEBUG MOCK: Returning ID for allDevices: 00000000-0000-0000-0000-000000000002\n")
		return "00000000-0000-0000-0000-000000000002"
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		fmt.Printf("DEBUG MOCK: Returning ID for allLicensedUsers: 00000000-0000-0000-0000-000000000003\n")
		return "00000000-0000-0000-0000-000000000003"
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		fmt.Printf("DEBUG MOCK: Returning ID for exclusionGroup: 00000000-0000-0000-0000-000000000004\n")
		return "00000000-0000-0000-0000-000000000004"
	case "#microsoft.graph.groupAssignmentTarget":
		fmt.Printf("DEBUG MOCK: Returning ID for groupAssignment: 00000000-0000-0000-0000-000000000001\n")
		return "00000000-0000-0000-0000-000000000001"
	case "#microsoft.graph.configurationManagerCollectionAssignmentTarget":
		fmt.Printf("DEBUG MOCK: Returning ID for configManager: 00000000-0000-0000-0000-000000000006\n")
		return "00000000-0000-0000-0000-000000000006"
	default:
		fmt.Printf("DEBUG MOCK: Unknown @odata.type '%s', returning default ID\n", odataType)
		return "00000000-0000-0000-0000-000000000001"
	}
}

// getAssignmentResponse returns the appropriate assignment response based on ID
func getAssignmentResponse(assignmentId string) map[string]interface{} {
	fmt.Printf("DEBUG MOCK: getAssignmentResponse called with ID: %s\n", assignmentId)

	var response map[string]interface{}

	switch assignmentId {
	case "00000000-0000-0000-0000-000000000001":
		// Group assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.groupAssignmentTarget",
				"groupId":     "11111111-1111-1111-1111-111111111111",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning groupAssignment data for ID %s\n", assignmentId)
	case "00000000-0000-0000-0000-000000000002":
		// All devices assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning allDevices data for ID %s\n", assignmentId)
	case "00000000-0000-0000-0000-000000000003":
		// All licensed users assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning allLicensedUsers data for ID %s\n", assignmentId)
	case "00000000-0000-0000-0000-000000000004":
		// Exclusion group assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
				"groupId":     "22222222-2222-2222-2222-222222222222",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning exclusionGroup data for ID %s\n", assignmentId)
	case "00000000-0000-0000-0000-000000000005":
		// Remove filter support - just return regular group assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.groupAssignmentTarget",
				"groupId":     "11111111-1111-1111-1111-111111111111",
				// Remove filter properties - not supported
			},
		}
		fmt.Printf("DEBUG MOCK: Returning groupAssignment (no filters) data for ID %s\n", assignmentId)
	case "00000000-0000-0000-0000-000000000006":
		// Configuration manager collection assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type":  "#microsoft.graph.configurationManagerCollectionAssignmentTarget",
				"collectionId": "11111111-1111-1111-1111-111111111111",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning configManager data for ID %s\n", assignmentId)
	default:
		// Default group assignment
		response = map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.groupAssignmentTarget",
				"groupId":     "11111111-1111-1111-1111-111111111111",
			},
		}
		fmt.Printf("DEBUG MOCK: Returning DEFAULT groupAssignment data for unknown ID %s\n", assignmentId)
	}

	fmt.Printf("DEBUG MOCK: Final response for ID %s: %+v\n", assignmentId, response)
	return response
}

// RegisterDeviceConfigurationAssignmentErrorMocks registers HTTP mocks for error scenarios
func (m *Mocks) RegisterDeviceConfigurationAssignmentErrorMocks() {
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

	// All operations return errors
	httpmock.RegisterRegexpResponder("POST",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/.*`),
		httpmock.NewStringResponder(400, `{"error": {"code": "BadRequest", "message": "Bad Request"}}`))

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/.*`),
		httpmock.NewStringResponder(403, `{"error": {"code": "Forbidden", "message": "Access denied"}}`))

	httpmock.RegisterRegexpResponder("DELETE",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/.*`),
		httpmock.NewStringResponder(403, `{"error": {"code": "Forbidden", "message": "Access denied"}}`))
}
