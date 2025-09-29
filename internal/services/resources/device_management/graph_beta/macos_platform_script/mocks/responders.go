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
	platformScripts map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.platformScripts = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// MacOSPlatformScriptMock provides mock responses for macOS platform script operations
type MacOSPlatformScriptMock struct{}

// RegisterMocks registers HTTP mock responses for macOS platform script operations
func (m *MacOSPlatformScriptMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.platformScripts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing platform scripts
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			scripts := make([]map[string]any, 0, len(mockState.platformScripts))
			for _, script := range mockState.platformScripts {
				scripts = append(scripts, script)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts",
				"value":          scripts,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual platform script
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			scriptData, exists := mockState.platformScripts[scriptId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Platform script not found"}}`), nil
			}

			// Create response copy
			responseCopy := make(map[string]any)
			for k, v := range scriptData {
				responseCopy[k] = v
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the script data
				if assignments, hasAssignments := scriptData["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
						responseCopy["assignments"] = assignments
					} else {
						// If assignments array is empty, return empty array (not null)
						responseCopy["assignments"] = []interface{}{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					responseCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, responseCopy)
		})

	// Register POST for creating platform script
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new platform script ID
			scriptId := uuid.New().String()

			// Create platform script data - only include fields that were provided or have defaults
			scriptData := map[string]any{
				"id":            scriptId,
				"displayName":   requestBody["displayName"],
				"fileName":      requestBody["fileName"],
				"scriptContent": requestBody["scriptContent"],
				"runAsAccount":  requestBody["runAsAccount"],
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				scriptData["description"] = description
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				scriptData["roleScopeTagIds"] = roleScopeTagIds
			} else {
				scriptData["roleScopeTagIds"] = []string{"0"} // Default value
			}
			if blockExecutionNotifications, exists := requestBody["blockExecutionNotifications"]; exists {
				scriptData["blockExecutionNotifications"] = blockExecutionNotifications
			}
			if executionFrequency, exists := requestBody["executionFrequency"]; exists {
				scriptData["executionFrequency"] = executionFrequency
			}
			if retryCount, exists := requestBody["retryCount"]; exists {
				scriptData["retryCount"] = retryCount
			}

			// Initialize assignments as empty array
			scriptData["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.platformScripts[scriptId] = scriptData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, scriptData)
		})

	// Register PATCH for updating platform script
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			scriptData, exists := mockState.platformScripts[scriptId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Platform script not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update platform script data
			mockState.Lock()

			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior

			// For optional fields, if they're not in the request, remove them
			optionalFields := []string{"description", "blockExecutionNotifications", "executionFrequency", "retryCount"}
			for _, field := range optionalFields {
				if _, hasField := requestBody[field]; !hasField {
					delete(scriptData, field)
				}
			}

			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(scriptData, key)
				} else {
					scriptData[key] = value
				}
			}
			// Ensure the ID is preserved
			scriptData["id"] = scriptId
			mockState.platformScripts[scriptId] = scriptData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, scriptData)
		})

	// Register DELETE for removing platform script
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.platformScripts[scriptId]
			if exists {
				delete(mockState.platformScripts, scriptId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Platform script not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2] // deviceShellScripts/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the script
			mockState.Lock()
			if scriptData, exists := mockState.platformScripts[scriptId]; exists {
				if assignments, hasAssignments := requestBody["deviceManagementScriptAssignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									// Generate a unique assignment ID
									assignmentId := uuid.New().String()

									// Create assignment in the format the API returns
									// The API returns the target exactly as submitted but with additional metadata
									targetCopy := make(map[string]any)
									for k, v := range target {
										targetCopy[k] = v
									}

									graphAssignment := map[string]any{
										"id":     assignmentId,
										"target": targetCopy,
									}
									graphAssignments = append(graphAssignments, graphAssignment)
								}
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
				mockState.platformScripts[scriptId] = scriptData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts/assignments",
				"value":          []map[string]any{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Dynamic mocks will handle all test cases
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MacOSPlatformScriptMock) RegisterErrorMocks() {
	// Register GET for listing platform scripts (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating platform script with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for platform script not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/not-found-script",
		factories.ErrorResponse(404, "ResourceNotFound", "Platform script not found"))
}
