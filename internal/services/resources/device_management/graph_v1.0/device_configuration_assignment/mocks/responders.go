package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// Global variable to track if update has happened
var updateOccurred = false

// DeviceConfigurationAssignmentMock provides mock responses for device configuration assignment operations
type DeviceConfigurationAssignmentMock struct{}

// RegisterMocks registers HTTP mock responses for device configuration assignment operations
func (m *DeviceConfigurationAssignmentMock) RegisterMocks() {
	// Reset the update flag
	updateOccurred = false

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

			// Set the update flag to true if this is the update we're testing for
			if deviceConfigId == "test-config-id" && assignmentId == "00000000-0000-0000-0000-000000000001" {
				updateOccurred = true
				fmt.Printf("DEBUG MOCK: Update flag set to true\n")
			}

			// For update, return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
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

			fmt.Printf("DEBUG MOCK: DELETE request for deviceConfig=%s, assignment=%s\n", deviceConfigId, assignmentId)

			// Handle error cases
			if deviceConfigId == "error-config" || assignmentId == "error-id" {
				return httpmock.NewStringResponse(400, `{"error": {"code": "BadRequest", "message": "Bad Request"}}`), nil
			}

			// For delete, return a 204 No Content
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses that return errors
func (m *DeviceConfigurationAssignmentMock) RegisterErrorMocks() {
	// POST Create Assignment - error
	httpmock.RegisterRegexpResponder("POST",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// GET Read Assignment - error
	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH Update Assignment - error
	httpmock.RegisterRegexpResponder("PATCH",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// DELETE Assignment - error
	httpmock.RegisterRegexpResponder("DELETE",
		regexp.MustCompile(`https://graph\.microsoft\.com/v1\.0/deviceManagement/deviceConfigurations/([^/]+)/assignments/([^/]+)`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
}

// Helper functions

// generateAssignmentId generates an assignment ID based on the target type
func generateAssignmentId(requestBody map[string]interface{}) string {
	if target, ok := requestBody["target"].(map[string]interface{}); ok {
		if odataType, ok := target["@odata.type"].(string); ok {
			switch odataType {
			case "#microsoft.graph.groupAssignmentTarget":
				return "00000000-0000-0000-0000-000000000001"
			case "#microsoft.graph.allDevicesAssignmentTarget":
				return "00000000-0000-0000-0000-000000000002"
			case "#microsoft.graph.allLicensedUsersAssignmentTarget":
				return "00000000-0000-0000-0000-000000000003"
			case "#microsoft.graph.exclusionGroupAssignmentTarget":
				return "00000000-0000-0000-0000-000000000004"
			}
		}
	}
	return "00000000-0000-0000-0000-000000000001" // Default to group assignment
}

// getAssignmentResponse returns a response for a given assignment ID
func getAssignmentResponse(assignmentId string) map[string]interface{} {
	switch assignmentId {
	case "00000000-0000-0000-0000-000000000001":
		// Group assignment
		return map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.groupAssignmentTarget",
				"groupId":     "11111111-1111-1111-1111-111111111111",
			},
		}
	case "00000000-0000-0000-0000-000000000002":
		// All devices
		return map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
			},
		}
	case "00000000-0000-0000-0000-000000000003":
		// All licensed users
		return map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget",
			},
		}
	case "00000000-0000-0000-0000-000000000004":
		// Exclusion group
		return map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.exclusionGroupAssignmentTarget",
				"groupId":     "22222222-2222-2222-2222-222222222222",
			},
		}
	default:
		// Default to group assignment
		return map[string]interface{}{
			"id": assignmentId,
			"target": map[string]interface{}{
				"@odata.type": "#microsoft.graph.groupAssignmentTarget",
				"groupId":     "11111111-1111-1111-1111-111111111111",
			},
		}
	}
}

func init() {
	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Note: In a real implementation, you would use the global registry from the mocks package
	// For example: mocks.GlobalRegistry.Register("device_configuration_assignment", &DeviceConfigurationAssignmentMock{})
}
