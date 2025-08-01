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
	customAttributeScripts map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.customAttributeScripts = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// MacOSCustomAttributeScriptMock provides mock responses for macOS custom attribute script operations
type MacOSCustomAttributeScriptMock struct{}

// RegisterMocks registers HTTP mock responses for macOS custom attribute script operations
func (m *MacOSCustomAttributeScriptMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.customAttributeScripts = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing custom attribute scripts
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			scripts := make([]map[string]interface{}, 0, len(mockState.customAttributeScripts))
			for _, script := range mockState.customAttributeScripts {
				scripts = append(scripts, script)
			}
			mockState.Unlock()

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCustomAttributeShellScripts",
				"value":          scripts,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual custom attribute script
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			scriptData, exists := mockState.customAttributeScripts[scriptId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Custom attribute script not found"}}`), nil
			}

			// Create response copy
			responseCopy := make(map[string]interface{})
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

	// Register POST for creating custom attribute script
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new custom attribute script ID
			scriptId := uuid.New().String()

			// Create custom attribute script data - only include fields that were provided or have defaults
			scriptData := map[string]interface{}{
				"id":                  scriptId,
				"displayName":         requestBody["displayName"],
				"customAttributeName": requestBody["customAttributeName"],
				"customAttributeType": requestBody["customAttributeType"],
				"fileName":            requestBody["fileName"],
				"scriptContent":       requestBody["scriptContent"],
				"runAsAccount":        requestBody["runAsAccount"],
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

			// Add computed fields that are always returned by the API
			scriptData["createdDateTime"] = "2024-01-01T00:00:00Z"
			scriptData["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Initialize assignments as empty array
			scriptData["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.customAttributeScripts[scriptId] = scriptData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, scriptData)
		})

	// Register PATCH for updating custom attribute script
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			scriptData, exists := mockState.customAttributeScripts[scriptId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Custom attribute script not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update custom attribute script data
			mockState.Lock()
			
			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior
			
			// For optional fields, if they're not in the request, remove them
			optionalFields := []string{"description"}
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
			// Ensure the ID is preserved and update timestamp
			scriptData["id"] = scriptId
			scriptData["lastModifiedDateTime"] = "2024-01-01T01:00:00Z"
			mockState.customAttributeScripts[scriptId] = scriptData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, scriptData)
		})

	// Register DELETE for removing custom attribute script
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.customAttributeScripts[scriptId]
			if exists {
				delete(mockState.customAttributeScripts, scriptId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Custom attribute script not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2] // deviceCustomAttributeShellScripts/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the script
			mockState.Lock()
			if scriptData, exists := mockState.customAttributeScripts[scriptId]; exists {
				if assignments, hasAssignments := requestBody["deviceManagementScriptAssignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]interface{}); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]interface{}); hasTarget {
									// Generate a unique assignment ID
									assignmentId := uuid.New().String()
									
									// Create assignment in the format the API returns
									// The API returns the target exactly as submitted but with additional metadata
									targetCopy := make(map[string]interface{})
									for k, v := range target {
										targetCopy[k] = v
									}
									
									graphAssignment := map[string]interface{}{
										"id": assignmentId,
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
				mockState.customAttributeScripts[scriptId] = scriptData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCustomAttributeShellScripts/assignments",
				"value":          []map[string]interface{}{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Dynamic mocks will handle all test cases
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MacOSCustomAttributeScriptMock) RegisterErrorMocks() {
	// Register GET for listing custom attribute scripts (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCustomAttributeShellScripts",
				"value":          []map[string]interface{}{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating custom attribute script with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for custom attribute script not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceCustomAttributeShellScripts/not-found-script",
		factories.ErrorResponse(404, "ResourceNotFound", "Custom attribute script not found"))
}