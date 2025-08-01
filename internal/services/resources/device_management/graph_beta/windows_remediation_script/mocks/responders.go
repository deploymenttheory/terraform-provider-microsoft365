package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	windowsRemediationScripts map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.windowsRemediationScripts = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// WindowsRemediationScriptMock provides mock responses for Windows remediation script operations
type WindowsRemediationScriptMock struct{}

// RegisterMocks registers HTTP mock responses for Windows remediation script operations
func (m *WindowsRemediationScriptMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.windowsRemediationScripts = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing Windows remediation scripts
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			scripts := make([]map[string]interface{}, 0, len(mockState.windowsRemediationScripts))
			for _, script := range mockState.windowsRemediationScripts {
				// Ensure @odata.type is present
				scriptCopy := make(map[string]interface{})
				for k, v := range script {
					scriptCopy[k] = v
				}
				if _, hasODataType := scriptCopy["@odata.type"]; !hasODataType {
					scriptCopy["@odata.type"] = "#microsoft.graph.deviceHealthScript"
				}

				// Check if expand=assignments is requested
				expandParam := req.URL.Query().Get("$expand")
				if strings.Contains(expandParam, "assignments") {
					// Include assignments if they exist in the script data
					if assignments, hasAssignments := script["assignments"]; hasAssignments && assignments != nil {
						if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
							scriptCopy["assignments"] = assignments
						} else {
							// If assignments array is empty, return empty array (not null)
							scriptCopy["assignments"] = []interface{}{}
						}
					} else {
						// If no assignments stored, return empty array (not null)
						scriptCopy["assignments"] = []interface{}{}
					}
				}

				scripts = append(scripts, scriptCopy)
			}
			mockState.Unlock()

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts",
				"value":          scripts,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Windows remediation script
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			script, exists := mockState.windowsRemediationScripts[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified device health script was not found"}}`), nil
			}

			// Create response copy
			scriptCopy := make(map[string]interface{})
			for k, v := range script {
				scriptCopy[k] = v
			}
			if _, hasODataType := scriptCopy["@odata.type"]; !hasODataType {
				scriptCopy["@odata.type"] = "#microsoft.graph.deviceHealthScript"
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the script data
				if assignments, hasAssignments := script["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {

						// Return assignments in Microsoft Graph SDK format (not transformed)
						// The SDK will handle the transformation to Terraform structure
						scriptCopy["assignments"] = assignments

					} else {
						// If assignments array is empty, return empty array (not null)
						scriptCopy["assignments"] = []interface{}{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					scriptCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, scriptCopy)
		})

	// Register POST for creating Windows remediation script
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a unique ID for the new script
			id := uuid.New().String()

			// Create the script object with required fields
			script := map[string]interface{}{
				"@odata.type":              "#microsoft.graph.deviceHealthScript",
				"id":                       id,
				"displayName":              requestBody["displayName"],
				"publisher":                requestBody["publisher"],
				"runAs32Bit":               getOptionalBool(requestBody, "runAs32Bit", false),
				"runAsAccount":             requestBody["runAsAccount"],
				"enforceSignatureCheck":    getOptionalBool(requestBody, "enforceSignatureCheck", false),
				"detectionScriptContent":   requestBody["detectionScriptContent"],
				"remediationScriptContent": requestBody["remediationScriptContent"],
				"version":                  "1.0",
				"isGlobalScript":           false,
				"deviceHealthScriptType":   "deviceHealthScript",
				"createdDateTime":          "2024-01-01T00:00:00Z",
				"lastModifiedDateTime":     "2024-01-01T00:00:00Z",
				"highestAvailableVersion":  "1.0",
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				script["description"] = description
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				script["roleScopeTagIds"] = roleScopeTagIds
			} else {
				script["roleScopeTagIds"] = []string{"0"} // Default value
			}

			// Initialize assignments as empty array
			script["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.windowsRemediationScripts[id] = script
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, script)
		})

	// Register PATCH for updating Windows remediation script
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			script, exists := mockState.windowsRemediationScripts[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified device health script was not found"}}`), nil
			}

			// For optional fields that can be null, we need to explicitly remove them if they're not in the PATCH request
			// This simulates how Microsoft Graph API handles PATCH requests for optional fields
			optionalFields := []string{"description"}

			// Check if optional fields are missing from the request and remove them from the stored state
			for _, field := range optionalFields {
				if _, exists := requestBody[field]; !exists {
					if _, existsInState := script[field]; existsInState {
						delete(script, field)
					}
				}
			}

			// Update the script with new values
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(script, key)
				} else {
					script[key] = value
				}
			}
			script["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Ensure @odata.type is present
			if _, hasODataType := script["@odata.type"]; !hasODataType {
				script["@odata.type"] = "#microsoft.graph.deviceHealthScript"
			}

			mockState.windowsRemediationScripts[id] = script
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, script)
		})

	// Register DELETE for Windows remediation script
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.windowsRemediationScripts[id]
			if exists {
				delete(mockState.windowsRemediationScripts, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified device health script was not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register assignment-related endpoints
	m.registerAssignmentMocks()
}

// registerAssignmentMocks registers mock responses for assignment operations
func (m *WindowsRemediationScriptMock) registerAssignmentMocks() {
	// POST assignment for a script
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2] // deviceHealthScripts/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the script
			mockState.Lock()
			if scriptData, exists := mockState.windowsRemediationScripts[scriptId]; exists {
				if assignments, hasAssignments := requestBody["deviceHealthScriptAssignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]interface{}); ok {
								// Generate a unique assignment ID
								assignmentId := uuid.New().String()

								// For device health scripts, assignments come with a "target" wrapper
								// Extract the target data from the assignment
								var target map[string]interface{}
								if targetData, hasTarget := assignmentMap["target"].(map[string]interface{}); hasTarget {
									target = make(map[string]interface{})
									// Copy target fields
									for k, v := range targetData {
										target[k] = v
									}
								} else {
									continue
								}

								// Handle runSchedule if present
								var runSchedule map[string]interface{}
								if scheduleData, hasSchedule := assignmentMap["runSchedule"].(map[string]interface{}); hasSchedule {
									runSchedule = make(map[string]interface{})
									// Copy schedule fields with proper structure
									for k, v := range scheduleData {
										runSchedule[k] = v
									}
								}

								// Keep the @odata.type field as-is for SDK processing
								// SDK will process the @odata.type field

								// Keep original Microsoft Graph API field names for SDK processing
								// The SDK will handle the field name mapping to Terraform structure
								graphAssignment := map[string]interface{}{
									"id":     assignmentId,
									"target": target,
								}

								// Add runSchedule if present
								if runSchedule != nil {
									graphAssignment["runSchedule"] = runSchedule
								}

								graphAssignments = append(graphAssignments, graphAssignment)
							}
						}
						scriptData["assignments"] = graphAssignments
					} else {
						// Set empty assignments array instead of deleting
						scriptData["assignments"] = []interface{}{}
					}
				} else {
					// Set empty assignments array instead of deleting
					scriptData["assignments"] = []interface{}{}
				}
				mockState.windowsRemediationScripts[scriptId] = scriptData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// GET assignments for a script
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			mockState.Lock()
			scriptData, exists := mockState.windowsRemediationScripts[id]
			mockState.Unlock()

			if !exists {
				response := map[string]interface{}{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts('" + id + "')/assignments",
					"value":          []map[string]interface{}{},
				}
				return httpmock.NewJsonResponse(200, response)
			}

			// Get assignments from stored script data
			assignments := []interface{}{}
			if storedAssignments, hasAssignments := scriptData["assignments"]; hasAssignments {
				if assignmentArray, ok := storedAssignments.([]interface{}); ok {
					assignments = assignmentArray
				}
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts('" + id + "')/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *WindowsRemediationScriptMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.windowsRemediationScripts = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing Windows remediation scripts (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts",
				"value":          []map[string]interface{}{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Windows remediation script with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for Windows remediation script not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts/not-found-script",
		factories.ErrorResponse(404, "ResourceNotFound", "Device health script not found"))
}

func getOptionalBool(data map[string]interface{}, key string, defaultValue bool) bool {
	if value, exists := data[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}
