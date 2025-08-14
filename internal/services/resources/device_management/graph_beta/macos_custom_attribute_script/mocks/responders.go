package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
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

	// Register with global registry
	mocks.GlobalRegistry.Register("macos_custom_attribute_script", &MacOSCustomAttributeScriptMock{})
}


// MacOSCustomAttributeScriptMock provides mock responses for macOS custom attribute script operations
type MacOSCustomAttributeScriptMock struct{}

// Ensure MacOSCustomAttributeScriptMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*MacOSCustomAttributeScriptMock)(nil)

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
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			script, exists := mockState.customAttributeScripts[id]
			mockState.Unlock()

			if !exists {
				// Check for special test IDs
				switch {
				case strings.Contains(id, "minimal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_macos_custom_attribute_script_minimal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var response map[string]interface{}
					if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				case strings.Contains(id, "maximal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_macos_custom_attribute_script_maximal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var response map[string]interface{}
					if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				default:
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
					var errorResponse map[string]interface{}
					_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}

			// Create response copy
			scriptCopy := make(map[string]interface{})
			for k, v := range script {
				scriptCopy[k] = v
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

			return httpmock.NewJsonResponse(200, scriptCopy)
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
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
				var errorResponse map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/patch_macos_custom_attribute_script_minimal.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updatedScript map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &updatedScript); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			mockState.Lock()
			script, exists := mockState.customAttributeScripts[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				var errorResponse map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
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
			mockState.customAttributeScripts[id] = updatedScript
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedScript)(req)
		})

	// Register DELETE for custom attribute script
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.customAttributeScripts[id]
			if exists {
				delete(mockState.customAttributeScripts, id)
			}
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				var errorResponse map[string]interface{}
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register assignment-related endpoints
	m.registerAssignmentMocks()
}

// registerAssignmentMocks registers mock responses for assignment operations
func (m *MacOSCustomAttributeScriptMock) registerAssignmentMocks() {
	// POST assignment for a script
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+/assign$`,
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

	// GET assignments for a script
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			mockState.Lock()
			scriptData, exists := mockState.customAttributeScripts[id]
			mockState.Unlock()

			if !exists {
				response := map[string]interface{}{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCustomAttributeShellScripts('" + id + "')/assignments",
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
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCustomAttributeShellScripts('" + id + "')/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *MacOSCustomAttributeScriptMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored custom attribute scripts
	for id := range mockState.customAttributeScripts {
		delete(mockState.customAttributeScripts, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *MacOSCustomAttributeScriptMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.customAttributeScripts = make(map[string]map[string]interface{})
	mockState.Unlock()

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
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
			var errorResponse map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for custom attribute script not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceCustomAttributeShellScripts/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
			var errorResponse map[string]interface{}
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}