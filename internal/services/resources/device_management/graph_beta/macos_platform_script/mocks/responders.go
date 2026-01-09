package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "delete_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Determine which read JSON to load based on stored script data
			scenarioFile := determineReadScenario(scriptData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for GET request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_read", scenarioFile))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load read scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse read scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Only override the ID from stored state, everything else comes from the JSON file
			response["id"] = scriptId

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the script data
				if assignments, hasAssignments := scriptData["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {
						response["assignments"] = assignments
					} else {
						// If assignments array is empty, return empty array (not null)
						response["assignments"] = []any{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					response["assignments"] = []any{}
				}
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating platform script
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_error_invalid_run_as.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Check for validation errors
			if runAsAccount, hasRunAs := requestBody["runAsAccount"].(string); hasRunAs {
				if runAsAccount != "system" && runAsAccount != "user" {
					content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_error_invalid_run_as.json"))
					var errorResponse map[string]any
					json.Unmarshal([]byte(content), &errorResponse)
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}

			if description, hasDesc := requestBody["description"].(string); hasDesc {
				if len(description) > 1500 {
					content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_error_description_length.json"))
					var errorResponse map[string]any
					json.Unmarshal([]byte(content), &errorResponse)
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}

			scriptId := uuid.New().String()

			// Determine scenario JSON to load from validate_create/
			scenarioFile := determineCreateScenario(requestBody)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for POST request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", scenarioFile))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load create scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
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
			if fileName, hasFileName := requestBody["fileName"]; hasFileName {
				response["fileName"] = fileName
			}

			// Store metadata to determine read/update scenarios
			// Convert POST scenario to GET scenario for future reads
			readScenario := strings.Replace(scenarioFile, "post_", "get_", 1)
			scriptMetadata := map[string]any{
				"id":          scriptId,
				"displayName": requestBody["displayName"],
				"scenario":    readScenario,
			}
			if desc, hasDesc := requestBody["description"]; hasDesc {
				scriptMetadata["description"] = desc
			}
			if fileName, hasFileName := requestBody["fileName"]; hasFileName {
				scriptMetadata["fileName"] = fileName
			}
			if runAs, hasRunAs := requestBody["runAsAccount"]; hasRunAs {
				scriptMetadata["runAsAccount"] = runAs
			}

			mockState.Lock()
			mockState.platformScripts[scriptId] = scriptMetadata
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "delete_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Determine update scenario
			scenarioFile := determineUpdateScenario(scriptData, requestBody)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for PATCH request"}}`), nil
			}

			content, updateErr := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_update", scenarioFile))
			if updateErr != nil {
				// If PATCH JSON doesn't exist, try loading the GET JSON (for assignment-only updates)
				getScenario, _ := scriptData["scenario"].(string)
				content, updateErr = helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_read", getScenario))
				if updateErr != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load update scenario JSON: `+updateErr.Error()+`"}}`), nil
				}
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Only override the ID, the JSON file contains all other correct values
			response["id"] = scriptId

			// Only update scenario if this is an actual resource update (not assignment-only)
			// Check if display name changed, which indicates a real resource update
			if displayName, hasName := requestBody["displayName"]; hasName {
				currentDisplayName, _ := scriptData["displayName"].(string)
				newDisplayName, _ := displayName.(string)

				// Only update scenario if display name actually changed
				if newDisplayName != currentDisplayName && newDisplayName != "" {
					// Update scenario marker for future reads - convert patch to get
					readScenario := strings.Replace(scenarioFile, "patch_", "get_", 1)
					scriptData["scenario"] = readScenario
					scriptData["displayName"] = newDisplayName
				}
			}

			mockState.Lock()
			mockState.platformScripts[scriptId] = scriptData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, response)
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "delete_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assigning platform script
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2]

			mockState.Lock()
			scriptData, exists := mockState.platformScripts[scriptId]
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "delete_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in mock state
			if assignments, hasAssignments := requestBody["deviceManagementScriptAssignments"]; hasAssignments {
				scriptData["assignments"] = assignments

				mockState.Lock()
				mockState.platformScripts[scriptId] = scriptData
				mockState.Unlock()
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for platform script assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			scriptId := urlParts[len(urlParts)-2]

			mockState.Lock()
			scriptData, exists := mockState.platformScripts[scriptId]
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "delete_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Check if assignments exist
			if assignments, hasAssignments := scriptData["assignments"]; hasAssignments {
				if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {
					// Determine assignment response based on count
					var scenarioFile string
					if len(assignmentList) == 1 {
						scenarioFile = "post_assign_single_group_success.json"
					} else {
						scenarioFile = "post_assign_all_types_success.json"
					}

					content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_assignments", scenarioFile))
					var response map[string]any
					json.Unmarshal([]byte(content), &response)
					return httpmock.NewJsonResponse(200, response)
				}
			}

			// No assignments - return empty
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_assignments", "get_assignments_empty_success.json"))
			var response map[string]any
			json.Unmarshal([]byte(content), &response)
			return httpmock.NewJsonResponse(200, response)
		})
}

// CleanupMockState cleans up the mock state after tests
func (m *MacOSPlatformScriptMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.platformScripts = make(map[string]map[string]any)
}

// RegisterErrorMocks registers error response mocks for negative testing
func (m *MacOSPlatformScriptMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.platformScripts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register error response for invalid run_as_account
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid run_as_account value"))

	// Register error response for platform script not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/not-found-script",
		factories.ErrorResponse(404, "ResourceNotFound", "Platform script not found"))
}

// determineCreateScenario determines which create scenario JSON to load based on request body
func determineCreateScenario(requestBody map[string]any) string {
	displayName, _ := requestBody["displayName"].(string)

	// Match based on display name patterns from test terraform files
	switch {
	case strings.Contains(displayName, "unit-test-minimal-macos-script"):
		return "post_scenario_01_minimal_success.json"
	case strings.Contains(displayName, "unit-test-maximal-macos-script"):
		return "post_scenario_02_maximal_success.json"
	case strings.Contains(displayName, "unit-test-update-test-script"):
		return "post_scenario_03_step_01_minimal_success.json"
	case strings.Contains(displayName, "unit-test-downgrade-test-script"):
		return "post_scenario_04_step_01_maximal_success.json"
	case strings.Contains(displayName, "unit-test-add-minimal-assignment"):
		return "post_scenario_05_step_01_no_assignments_success.json"
	case strings.Contains(displayName, "unit-test-add-maximal-assignments"):
		return "post_scenario_06_step_01_no_assignments_success.json"
	case strings.Contains(displayName, "unit-test-assignment-update"):
		return "post_scenario_07_step_01_no_assignments_success.json"
	case strings.Contains(displayName, "unit-test-assignment-downgrade"):
		return "post_scenario_08_step_01_no_assignments_success.json"
	default:
		return "post_scenario_01_minimal_success.json"
	}
}

// determineReadScenario determines which read scenario JSON to load based on stored data
func determineReadScenario(scriptData map[string]any) string {
	if scenario, hasScenario := scriptData["scenario"].(string); hasScenario {
		return scenario
	}

	// Fallback to minimal
	return "get_scenario_01_minimal_success.json"
}

// determineUpdateScenario determines which update scenario JSON to load based on current state and update request
func determineUpdateScenario(scriptData map[string]any, requestBody map[string]any) string {
	displayName, _ := requestBody["displayName"].(string)

	// Match update scenarios
	switch {
	case strings.Contains(displayName, "unit-test-update-test-script-updated"):
		return "patch_scenario_03_step_02_maximal_success.json"
	case strings.Contains(displayName, "unit-test-downgrade-test-script-minimal"):
		return "patch_scenario_04_step_02_minimal_success.json"
	default:
		// For assignment-only updates, keep the current scenario
		// The display name won't change, so use the current read scenario
		if currentScenario, ok := scriptData["scenario"].(string); ok {
			// Convert get_ to patch_ if a PATCH JSON exists, otherwise keep current
			patchScenario := strings.Replace(currentScenario, "get_", "patch_", 1)
			return patchScenario
		}
		// Fallback to scenario 03 step 2
		return "patch_scenario_03_step_02_maximal_success.json"
	}
}
