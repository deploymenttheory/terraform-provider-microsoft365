package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	windowsRemediationScripts map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.windowsRemediationScripts = make(map[string]map[string]any)

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
	mockState.windowsRemediationScripts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Windows remediation scripts
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			scripts := make([]map[string]any, 0, len(mockState.windowsRemediationScripts))
			for _, script := range mockState.windowsRemediationScripts {
				// Ensure @odata.type is present
				scriptCopy := make(map[string]any)
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
						if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {
							scriptCopy["assignments"] = assignments
						} else {
							// If assignments array is empty, return empty array (not null)
							scriptCopy["assignments"] = []any{}
						}
					} else {
						// If no assignments stored, return empty array (not null)
						scriptCopy["assignments"] = []any{}
					}
				}

				scripts = append(scripts, scriptCopy)
			}
			mockState.Unlock()

			response := map[string]any{
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
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_remediation_script_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			scriptCopy := make(map[string]any)
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
					if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {

						// Return assignments in Microsoft Graph SDK format (not transformed)
						// The SDK will handle the transformation to Terraform structure
						scriptCopy["assignments"] = assignments

					} else {
						// If assignments array is empty, return empty array (not null)
						scriptCopy["assignments"] = []any{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					scriptCopy["assignments"] = []any{}
				}
			}

			return httpmock.NewJsonResponse(200, scriptCopy)
		})

	// Register POST for creating Windows remediation script
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_remediation_script_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			scriptId := uuid.New().String()

			// Determine scenario JSON to load from validate_create/
			scenarioFile := determineCreateScenario(requestBody)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for POST request"}}`), nil
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load create scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse create scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Override with request-specific values
			response["id"] = scriptId
			if displayName, hasName := requestBody["displayName"]; hasName {
				response["displayName"] = displayName
			}
			if description, hasDesc := requestBody["description"]; hasDesc {
				response["description"] = description
			}
			if publisher, hasPublisher := requestBody["publisher"]; hasPublisher {
				response["publisher"] = publisher
			}
			if runAs32Bit, has := requestBody["runAs32Bit"]; has {
				response["runAs32Bit"] = runAs32Bit
			}
			if runAsAccount, has := requestBody["runAsAccount"]; has {
				response["runAsAccount"] = runAsAccount
			}
			if enforceSignatureCheck, has := requestBody["enforceSignatureCheck"]; has {
				response["enforceSignatureCheck"] = enforceSignatureCheck
			}
			if detectionScriptContent, has := requestBody["detectionScriptContent"]; has {
				response["detectionScriptContent"] = detectionScriptContent
			}
			if remediationScriptContent, has := requestBody["remediationScriptContent"]; has {
				response["remediationScriptContent"] = remediationScriptContent
			}
			if roleScopeTagIds, has := requestBody["roleScopeTagIds"]; has {
				response["roleScopeTagIds"] = roleScopeTagIds
			}

			// Store in mock state
			mockState.Lock()
			mockState.windowsRemediationScripts[scriptId] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Windows remediation script
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			scriptData, exists := mockState.windowsRemediationScripts[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_remediation_script_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_remediation_script_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Determine which update JSON to load based on request body
			scenarioFile := determineUpdateScenario(requestBody, scriptData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for PATCH request"}}`), nil
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load update scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Merge request updates into response
			for k, v := range requestBody {
				response[k] = v
			}
			response["id"] = id
			response["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Preserve existing assignments if not in request
			if _, hasAssignments := requestBody["assignments"]; !hasAssignments {
				if assignments, hasExisting := scriptData["assignments"]; hasExisting {
					response["assignments"] = assignments
				}
			}

			mockState.Lock()
			mockState.windowsRemediationScripts[id] = response
			mockState.Unlock()

			return factories.SuccessResponse(200, response)(req)
		})

	// Register DELETE for Windows remediation script
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/[^/]+$`,
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
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_remediation_script_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
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
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the script
			mockState.Lock()
			if scriptData, exists := mockState.windowsRemediationScripts[scriptId]; exists {
				if assignments, hasAssignments := requestBody["deviceHealthScriptAssignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]any)
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []any{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								// Generate a unique assignment ID
								assignmentId := uuid.New().String()

								// For device health scripts, assignments come with a "target" wrapper
								// Extract the target data from the assignment
								var target map[string]any
								if targetData, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									target = make(map[string]any)
									// Copy target fields
									for k, v := range targetData {
										target[k] = v
									}
								} else {
									continue
								}

								// Handle runSchedule if present
								var runSchedule map[string]any
								if scheduleData, hasSchedule := assignmentMap["runSchedule"].(map[string]any); hasSchedule {
									runSchedule = make(map[string]any)
									// Copy schedule fields with proper structure
									for k, v := range scheduleData {
										runSchedule[k] = v
									}
								}

								// Keep the @odata.type field as-is for SDK processing
								// SDK will process the @odata.type field

								// Keep original Microsoft Graph API field names for SDK processing
								// The SDK will handle the field name mapping to Terraform structure
								graphAssignment := map[string]any{
									"id":     assignmentId,
									"target": target,
								}

								// Add runSchedule if present
								if runSchedule != nil {
									graphAssignment["runSchedule"] = runSchedule
								}

								// Add runRemediationScript if present
								if runRemediationScript, hasRunRemediation := assignmentMap["runRemediationScript"]; hasRunRemediation {
									graphAssignment["runRemediationScript"] = runRemediationScript
								}

								graphAssignments = append(graphAssignments, graphAssignment)
							}
						}
						scriptData["assignments"] = graphAssignments
					} else {
						// Set empty assignments array instead of deleting
						scriptData["assignments"] = []any{}
					}
				} else {
					// Set empty assignments array instead of deleting
					scriptData["assignments"] = []any{}
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
				response := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts('" + id + "')/assignments",
					"value":          []map[string]any{},
				}
				return httpmock.NewJsonResponse(200, response)
			}

			// Get assignments from stored script data
			assignments := []any{}
			if storedAssignments, hasAssignments := scriptData["assignments"]; hasAssignments {
				if assignmentArray, ok := storedAssignments.([]any); ok {
					assignments = assignmentArray
				}
			}

			response := map[string]any{
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

// Removed legacy JSON loader in favor of helpers.ParseJSONFile

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *WindowsRemediationScriptMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.windowsRemediationScripts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Windows remediation scripts (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceHealthScripts",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Windows remediation script with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceHealthScripts",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_remediation_script_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Windows remediation script not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceHealthScripts/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_remediation_script_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

// determineCreateScenario determines which create JSON to load based on request body
func determineCreateScenario(requestBody map[string]any) string {
	// Check displayName for scenario hints
	displayName, hasName := requestBody["displayName"].(string)
	if hasName {
		if strings.Contains(strings.ToLower(displayName), "minimal") {
			return "post_scenario_001_minimal.json"
		}
		if strings.Contains(strings.ToLower(displayName), "maximal") {
			return "post_scenario_002_maximal.json"
		}
		if strings.Contains(strings.ToLower(displayName), "assignment") {
			// Assignments will be added via the assign endpoint, not during creation
			// For now, return minimal for assignment tests
			return "post_scenario_001_minimal.json"
		}
	}

	// Default to minimal
	return "post_scenario_001_minimal.json"
}

// determineReadScenario determines which read JSON to load based on stored script data
func determineReadScenario(scriptData map[string]any) string {
	// Check displayName for scenario hints
	displayName, hasName := scriptData["displayName"].(string)
	if hasName {
		if strings.Contains(strings.ToLower(displayName), "minimal") {
			return "get_scenario_001_minimal.json"
		}
		if strings.Contains(strings.ToLower(displayName), "maximal") {
			return "get_scenario_002_maximal.json"
		}
	}

	// Default to minimal
	return "get_scenario_001_minimal.json"
}

// determineUpdateScenario determines which update JSON to load based on request body
func determineUpdateScenario(requestBody map[string]any, scriptData map[string]any) string {
	// Check displayName for scenario hints
	displayName, hasName := requestBody["displayName"].(string)
	if !hasName {
		// If no displayName in request, check stored data
		displayName, hasName = scriptData["displayName"].(string)
	}

	if hasName {
		if strings.Contains(strings.ToLower(displayName), "minimal") {
			return "patch_scenario_001_minimal.json"
		}
		if strings.Contains(strings.ToLower(displayName), "maximal") {
			return "patch_scenario_002_maximal.json"
		}
	}

	// Default to minimal
	return "patch_scenario_001_minimal.json"
}
