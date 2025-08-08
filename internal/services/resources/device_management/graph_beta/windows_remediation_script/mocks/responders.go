package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
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

	// Register with global registry
	mocks.GlobalRegistry.Register("windows_remediation_script", &WindowsRemediationScriptMock{})
}

// WindowsRemediationScriptMock provides mock responses for Windows remediation script operations
type WindowsRemediationScriptMock struct{}

// Ensure WindowsRemediationScriptMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*WindowsRemediationScriptMock)(nil)

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
				// Check for special test IDs
				switch {
				case strings.Contains(id, "minimal"):
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_windows_remediation_script_minimal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				case strings.Contains(id, "maximal"):
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_windows_remediation_script_maximal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				default:
					errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_remediation_script_not_found.json"))
					return httpmock.NewJsonResponse(404, errorResponse)
				}
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
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_windows_remediation_script_error.json"))
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			updatedScript, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "get_windows_remediation_script_updated.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}

			mockState.Lock()
			script, exists := mockState.windowsRemediationScripts[id]
			if !exists {
				mockState.Unlock()
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_remediation_script_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range script {
				updatedScript[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedScript[k] = v
			}

			// Store updated state
			mockState.windowsRemediationScripts[id] = updatedScript
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedScript)(req)
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
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_remediation_script_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
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

// CleanupMockState clears the mock state for clean test runs
func (m *WindowsRemediationScriptMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Windows remediation scripts
	for id := range mockState.windowsRemediationScripts {
		delete(mockState.windowsRemediationScripts, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *WindowsRemediationScriptMock) loadJSONResponse(filePath string) (map[string]interface{}, error) {
	var response map[string]interface{}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
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
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_windows_remediation_script_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Windows remediation script not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_remediation_script_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

func getOptionalBool(data map[string]interface{}, key string, defaultValue bool) bool {
	if value, exists := data[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}
