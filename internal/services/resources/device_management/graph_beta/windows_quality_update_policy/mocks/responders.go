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

var mockState struct {
	sync.Mutex
	qualityPolicies map[string]map[string]any
}

func init() {
	mockState.qualityPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_quality_update_policy", &WindowsQualityUpdatePolicyMock{})
}

type WindowsQualityUpdatePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsQualityUpdatePolicyMock)(nil)

func (m *WindowsQualityUpdatePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.qualityPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// List
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]any, 0, len(mockState.qualityPolicies))
			for _, v := range mockState.qualityPolicies {
				copy := map[string]any{}
				for k, vv := range v {
					copy[k] = vv
				}
				list = append(list, copy)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdatePolicies",
				"value":          list,
			})
		})

	// Get by id
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policyData, exists := mockState.qualityPolicies[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Determine scenario JSON to load from validate_read/
			scenarioFile := determineReadScenario(policyData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for GET request"}}`), nil
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load read scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse read scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Override with stored values
			for k, v := range policyData {
				response[k] = v
			}

			// Only include hotpatchEnabled if it was explicitly set
			if _, has := policyData["hotpatchEnabled"]; !has {
				delete(response, "hotpatchEnabled")
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			id := uuid.New().String()

			// Determine scenario JSON to load from validate_create/
			scenarioFile := determineCreateScenario(body)
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
			response["id"] = id
			if displayName, hasName := body["displayName"]; hasName {
				response["displayName"] = displayName
			}
			if description, hasDesc := body["description"]; hasDesc {
				response["description"] = description
			}
			// Only include hotpatchEnabled if it was explicitly set in the request
			if hotpatchEnabled, has := body["hotpatchEnabled"]; has {
				response["hotpatchEnabled"] = hotpatchEnabled
			} else {
				// Remove hotpatchEnabled from response if not in request
				delete(response, "hotpatchEnabled")
			}
			if roleScopeTagIds, has := body["roleScopeTagIds"]; has {
				response["roleScopeTagIds"] = roleScopeTagIds
			}

			// Store in mock state
			mockState.Lock()
			mockState.qualityPolicies[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Patch
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policyData, exists := mockState.qualityPolicies[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Merge request body into existing policy data
			for k, v := range requestBody {
				policyData[k] = v
			}

			// Determine scenario JSON to load from validate_update/
			scenarioFile := determineUpdateScenario(requestBody, policyData)
			if scenarioFile == "" {
				// If no specific update scenario, just return the merged data
				// Only include hotpatchEnabled if it was explicitly set
				if _, has := policyData["hotpatchEnabled"]; !has {
					delete(policyData, "hotpatchEnabled")
				}
				mockState.Lock()
				mockState.qualityPolicies[id] = policyData
				mockState.Unlock()
				return httpmock.NewJsonResponse(200, policyData)
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/" + scenarioFile)
			if err != nil {
				// Fall back to returning the merged data if JSON not found
				mockState.Lock()
				mockState.qualityPolicies[id] = policyData
				mockState.Unlock()
				return httpmock.NewJsonResponse(200, policyData)
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Update response with values from the request body
			for k, v := range requestBody {
				response[k] = v
			}

			// Preserve the ID from stored state
			response["id"] = policyData["id"]

			// Only include hotpatchEnabled if it was explicitly set in request
			if _, has := requestBody["hotpatchEnabled"]; !has {
				delete(response, "hotpatchEnabled")
			}

			mockState.Lock()
			mockState.qualityPolicies[id] = response
			mockState.Unlock()

			return factories.SuccessResponse(200, response)(req)
		})

	// Assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			mockState.Lock()
			if existing, ok := mockState.qualityPolicies[id]; ok {
				assignments, _ := body["assignments"].([]any)
				if assignments == nil {
					assignments = []any{}
				}
				existing["assignments"] = assignments
				mockState.qualityPolicies[id] = existing
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.qualityPolicies, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsQualityUpdatePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.qualityPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdatePolicies",
				"value":          []any{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdatePolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_policy_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdatePolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		})
}

func (m *WindowsQualityUpdatePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.qualityPolicies {
		delete(mockState.qualityPolicies, id)
	}
}

// determineCreateScenario determines which create JSON to load based on request body
func determineCreateScenario(requestBody map[string]any) string {
	displayName, hasName := requestBody["displayName"].(string)
	if !hasName {
		return "post_001_scenario_minimal.json"
	}

	name := strings.ToLower(displayName)

	// Check for specific test scenarios based on test number patterns
	if strings.Contains(name, "001") {
		return "post_001_scenario_minimal.json"
	}
	if strings.Contains(name, "002") {
		return "post_002_scenario_maximal.json"
	}
	if strings.Contains(name, "003") {
		return "post_003_lifecycle_step_1.json"
	}
	if strings.Contains(name, "004") {
		return "post_004_lifecycle_step_1.json"
	}
	if strings.Contains(name, "005") {
		return "post_005_assignments_minimal.json"
	}
	if strings.Contains(name, "006") {
		return "post_006_assignments_maximal.json"
	}
	if strings.Contains(name, "007") {
		return "post_007_assignments_lifecycle_step_1.json"
	}
	if strings.Contains(name, "008") {
		return "post_008_assignments_lifecycle_step_1.json"
	}

	// Fallback to checking keywords
	if strings.Contains(name, "minimal") {
		return "post_001_scenario_minimal.json"
	}
	if strings.Contains(name, "maximal") {
		return "post_002_scenario_maximal.json"
	}

	// Default to minimal
	return "post_001_scenario_minimal.json"
}

// determineReadScenario determines which read JSON to load based on stored policy data
func determineReadScenario(policyData map[string]any) string {
	displayName, hasName := policyData["displayName"].(string)
	if !hasName {
		return "get_001_scenario_minimal.json"
	}

	name := strings.ToLower(displayName)

	// Check for lifecycle step indicators first (step_1 vs step_2)
	// This helps distinguish between initial and updated states

	// Test 003 lifecycle
	if strings.Contains(name, "003") {
		// Check if it's been updated (has description or hotpatch enabled)
		if desc, hasDesc := policyData["description"]; hasDesc && desc != nil && desc != "" {
			return "get_003_lifecycle_step_2.json"
		}
		if hotpatch, hasHotpatch := policyData["hotpatchEnabled"]; hasHotpatch && hotpatch == true {
			return "get_003_lifecycle_step_2.json"
		}
		return "get_003_lifecycle_step_1.json"
	}

	// Test 004 lifecycle
	if strings.Contains(name, "004") {
		// Check if it's been updated (description removed or hotpatch disabled)
		if desc, hasDesc := policyData["description"]; !hasDesc || desc == nil || desc == "" {
			// Check roleScopeTagIds - step 2 has only ["0"]
			if tags, hasTags := policyData["roleScopeTagIds"].([]any); hasTags && len(tags) == 1 {
				return "get_004_lifecycle_step_2.json"
			}
		}
		return "get_004_lifecycle_step_1.json"
	}

	// Test 007 assignments lifecycle
	if strings.Contains(name, "007") {
		// Check assignments count
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments {
			if len(assignments) > 1 {
				return "get_007_assignments_lifecycle_step_2.json"
			}
		}
		return "get_007_assignments_lifecycle_step_1.json"
	}

	// Test 008 assignments lifecycle
	if strings.Contains(name, "008") {
		// Check assignments count
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments {
			if len(assignments) == 1 {
				return "get_008_assignments_lifecycle_step_2.json"
			}
		}
		return "get_008_assignments_lifecycle_step_1.json"
	}

	// Check for specific test scenarios
	if strings.Contains(name, "001") {
		return "get_001_scenario_minimal.json"
	}
	if strings.Contains(name, "002") {
		return "get_002_scenario_maximal.json"
	}
	if strings.Contains(name, "005") {
		return "get_005_assignments_minimal.json"
	}
	if strings.Contains(name, "006") {
		return "get_006_assignments_maximal.json"
	}

	// Default to minimal
	return "get_001_scenario_minimal.json"
}

// determineUpdateScenario determines which update JSON to load based on request body and stored data
func determineUpdateScenario(requestBody map[string]any, policyData map[string]any) string {
	displayName, hasName := requestBody["displayName"].(string)
	if !hasName {
		// If no displayName in request, check stored data
		displayName, hasName = policyData["displayName"].(string)
	}

	if !hasName {
		return ""
	}

	name := strings.ToLower(displayName)

	// Check for lifecycle test updates (step 2)
	if strings.Contains(name, "003") {
		// Step 2 update should have description or hotpatch enabled
		if _, hasDesc := requestBody["description"]; hasDesc {
			return "patch_003_lifecycle_step_2.json"
		}
		if hotpatch, hasHotpatch := requestBody["hotpatchEnabled"]; hasHotpatch && hotpatch == true {
			return "patch_003_lifecycle_step_2.json"
		}
	}

	if strings.Contains(name, "004") {
		// Step 2 update should remove/clear description or disable hotpatch or reduce role scope tags
		if desc, hasDesc := requestBody["description"]; hasDesc && (desc == nil || desc == "") {
			return "patch_004_lifecycle_step_2.json"
		}
		if hotpatch, hasHotpatch := requestBody["hotpatchEnabled"]; hasHotpatch && hotpatch == false {
			return "patch_004_lifecycle_step_2.json"
		}
		// Check if roleScopeTagIds is being reduced to 1 element
		if tags, hasTags := requestBody["roleScopeTagIds"].([]any); hasTags && len(tags) == 1 {
			return "patch_004_lifecycle_step_2.json"
		}
	}

	// For assignment lifecycle tests, updates happen via the assign endpoint, not PATCH
	// But we still provide the JSON files for completeness
	if strings.Contains(name, "007") {
		// Check if assignments are being added (step 2)
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments && len(assignments) > 1 {
			return "patch_007_assignments_lifecycle_step_2.json"
		}
	}

	if strings.Contains(name, "008") {
		// Check if assignments are being reduced (step 2)
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments && len(assignments) == 1 {
			return "patch_008_assignments_lifecycle_step_2.json"
		}
	}

	// Return empty string if no specific update scenario matched
	return ""
}
