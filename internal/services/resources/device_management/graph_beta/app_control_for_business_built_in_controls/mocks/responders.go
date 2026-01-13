package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	appControlPolicies map[string]map[string]any
}

func init() {
	mockState.appControlPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("app_control_for_business_built_in_controls", &AppControlForBusinessBuiltInControlsMock{})
}

// AppControlForBusinessBuiltInControlsMock provides mock responses for App Control operations
type AppControlForBusinessBuiltInControlsMock struct{}

var _ mocks.MockRegistrar = (*AppControlForBusinessBuiltInControlsMock)(nil)

// RegisterMocks registers HTTP mock responses for App Control operations
func (m *AppControlForBusinessBuiltInControlsMock) RegisterMocks() {
	mockState.Lock()
	mockState.appControlPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing App Control policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			// Filter by template for app control policies
			if !strings.Contains(req.URL.RawQuery, "4321b946-b76b-4450-8afd-769c08b16ffc_1") {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
					"value":          []any{},
				})
			}

			mockState.Lock()
			policies := make([]map[string]any, 0, len(mockState.appControlPolicies))
			for _, policy := range mockState.appControlPolicies {
				policies = append(policies, policy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual App Control policy
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.appControlPolicies[policyId]
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Determine which read JSON to load based on stored policy data
			scenarioFile := determineReadScenario(policyData)
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

			// Merge stored state into response
			response["id"] = policyId
			if name, hasName := policyData["name"]; hasName {
				response["name"] = name
			}
			if description, hasDesc := policyData["description"]; hasDesc {
				response["description"] = description
			}
			if roleScopeTagIds, hasRoles := policyData["roleScopeTagIds"]; hasRoles {
				response["roleScopeTagIds"] = roleScopeTagIds
			}
			if settings, hasSettings := policyData["settings"]; hasSettings {
				response["settings"] = settings
			}
			if assignments, hasAssignments := policyData["assignments"]; hasAssignments {
				response["assignments"] = assignments
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for policy settings
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-2]

			mockState.Lock()
			policy, exists := mockState.appControlPolicies[policyId]
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			if settings, ok := policy["settings"]; ok {
				if settingsArray, isArray := settings.([]any); isArray {
					return httpmock.NewJsonResponse(200, map[string]any{
						"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyId + "')/settings",
						"value":          settingsArray,
					})
				}
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyId + "')/settings",
				"value":          []any{},
			})
		})

	// Register GET for policy assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-2]

			mockState.Lock()
			policy, exists := mockState.appControlPolicies[policyId]
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			if assignments, ok := policy["assignments"]; ok {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyId + "')/assignments",
					"value":          assignments,
				})
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyId + "')/assignments",
				"value":          []any{},
			})
		})

	// Register POST for creating App Control policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_app_control_error.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			policyId := uuid.New().String()

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
			response["id"] = policyId
			if name, hasName := requestBody["name"]; hasName {
				response["name"] = name
			}
			if description, hasDesc := requestBody["description"]; hasDesc {
				response["description"] = description
			}
			if roleScopeTagIds, hasRoles := requestBody["roleScopeTagIds"]; hasRoles {
				response["roleScopeTagIds"] = roleScopeTagIds
			}
			if settings, hasSettings := requestBody["settings"]; hasSettings {
				response["settings"] = settings
			}

			// Assignments handled separately via /assign endpoint
			response["assignments"] = []any{}

			mockState.Lock()
			mockState.appControlPolicies[policyId] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PUT for updating App Control policy - matches format: configurationPolicies('{id}')
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\(`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL format: /deviceManagement/configurationPolicies('{id}')
			path := req.URL.Path
			startIdx := strings.Index(path, "('")
			endIdx := strings.Index(path, "')")
			if startIdx == -1 || endIdx == -1 {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL format"}}`), nil
			}
			policyId := path[startIdx+2 : endIdx]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			policyData, exists := mockState.appControlPolicies[policyId]
			if !exists {
				mockState.Unlock()
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Determine update scenario JSON
			scenarioFile := determineUpdateScenario(requestBody, policyData)
			if scenarioFile == "" {
				mockState.Unlock()
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for update request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_update", scenarioFile))
			if err != nil {
				mockState.Unlock()
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load update scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				mockState.Unlock()
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Override with request values
			response["id"] = policyId
			for k, v := range requestBody {
				response[k] = v
				policyData[k] = v
			}

			// Preserve existing values
			for k, v := range policyData {
				if _, hasKey := response[k]; !hasKey {
					response[k] = v
				}
			}

			// Preserve existing assignments if not in request
			if _, hasAssignments := requestBody["assignments"]; !hasAssignments {
				if assignments, hasExisting := policyData["assignments"]; hasExisting {
					response["assignments"] = assignments
				}
			}

			mockState.appControlPolicies[policyId] = response
			mockState.Unlock()

			// PUT returns 204 No Content on success
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-2]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			if policyData, exists := mockState.appControlPolicies[policyId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]any)
					if len(assignmentList) > 0 {
						graphAssignments := []any{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									assignmentId := uuid.New().String()
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
						policyData["assignments"] = graphAssignments
					} else {
						policyData["assignments"] = []any{}
					}
				} else {
					policyData["assignments"] = []any{}
				}
				mockState.appControlPolicies[policyId] = policyData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register DELETE for removing App Control policy
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.appControlPolicies[policyId]
			if exists {
				delete(mockState.appControlPolicies, policyId)
			}
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// hasAdditionalTrustRules checks if settings contain additional_rules_for_trusting_apps
func hasAdditionalTrustRules(settings []any) bool {
	// Convert settings to JSON and check for additional trust rules
	settingsJSON, _ := json.Marshal(settings)
	settingsStr := string(settingsJSON)

	// Look for trust_apps setting ID which indicates additional rules
	return strings.Contains(settingsStr, "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps")
}

// determineCreateScenario determines which create JSON to load based on request body
func determineCreateScenario(requestBody map[string]any) string {
	// Extract policy name to determine test scenario
	name, hasName := requestBody["name"].(string)
	if !hasName {
		return ""
	}

	// Check for assignment tests first
	if strings.Contains(name, "assignments") {
		return "post_test_007_assignments_minimal.json"
	}

	// Check if this has additional trust_apps rules (indicates maximal)
	hasAdditionalRules := false
	if settings, hasSettings := requestBody["settings"].([]any); hasSettings {
		hasAdditionalRules = hasAdditionalTrustRules(settings)
	}

	// Determine test scenario based on name and settings
	switch {
	case strings.Contains(name, "audit-mode"):
		return "post_test_001_audit_mode.json"
	case strings.Contains(name, "enforce-mode"):
		return "post_test_002_enforce_mode.json"
	case strings.Contains(name, "lifecycle"):
		// Check for additional rules
		if hasAdditionalRules {
			return "post_test_004_maximal.json"
		}
		return "post_test_003_minimal.json"
	case strings.Contains(name, "downgrade"):
		// For downgrade tests, start with maximal
		if hasAdditionalRules {
			return "post_test_004_maximal.json"
		}
		return "post_test_003_minimal.json"
	case strings.Contains(name, "maximal"):
		return "post_test_004_maximal.json"
	case strings.Contains(name, "minimal"):
		return "post_test_003_minimal.json"
	default:
		// Check settings to determine if maximal or minimal
		if hasAdditionalRules {
			return "post_test_004_maximal.json"
		}
		return "post_test_003_minimal.json"
	}
}

// determineReadScenario determines which read JSON to load based on stored policy data
func determineReadScenario(policyData map[string]any) string {
	name, hasName := policyData["name"].(string)
	if !hasName {
		return ""
	}

	// Check for assignments
	if assignments, hasAssignments := policyData["assignments"]; hasAssignments {
		if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {
			if len(assignmentList) >= 3 {
				return "get_test_007_assignments_maximal.json"
			}
			return "get_test_007_assignments_minimal.json"
		}
	}

	// Determine test scenario based on name
	switch {
	case strings.Contains(name, "audit-mode"):
		return "get_test_001_audit_mode.json"
	case strings.Contains(name, "enforce-mode"):
		return "get_test_002_enforce_mode.json"
	case strings.Contains(name, "minimal"):
		return "get_test_003_minimal.json"
	case strings.Contains(name, "maximal"):
		return "get_test_004_maximal.json"
	case strings.Contains(name, "lifecycle"), strings.Contains(name, "downgrade"):
		// Check role scope tags and settings to determine state
		hasAdditionalRules := false
		if settings, hasSettings := policyData["settings"].([]any); hasSettings {
			settingsStr := fmt.Sprintf("%v", settings)
			hasAdditionalRules = strings.Contains(settingsStr, "trust_apps")
		}

		if roleScopeTagIds, hasRoles := policyData["roleScopeTagIds"].([]any); hasRoles {
			if len(roleScopeTagIds) >= 3 || hasAdditionalRules {
				return "get_test_004_maximal.json"
			}
		}
		return "get_test_003_minimal.json"
	default:
		return "get_test_003_minimal.json"
	}
}

// determineUpdateScenario determines which update JSON to load based on request and existing data
func determineUpdateScenario(requestBody map[string]any, existingData map[string]any) string {
	// Check if upgrading or downgrading based on settings
	hasAdditionalRules := false

	// Check request body settings first
	if settings, hasSettings := requestBody["settings"].([]any); hasSettings {
		hasAdditionalRules = hasAdditionalTrustRules(settings)
	}

	// If not in request, check existing data
	if !hasAdditionalRules {
		if settings, hasSettings := existingData["settings"].([]any); hasSettings {
			hasAdditionalRules = hasAdditionalTrustRules(settings)
		}
	}

	// Determine if this is maximal or minimal based on additional rules
	if hasAdditionalRules {
		return "put_test_004_maximal.json"
	}

	return "put_test_003_minimal.json"
}

// CleanupMockState clears the mock state for clean test runs
func (m *AppControlForBusinessBuiltInControlsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	for id := range mockState.appControlPolicies {
		delete(mockState.appControlPolicies, id)
	}
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *AppControlForBusinessBuiltInControlsMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_app_control_error.json"))
			var errorResponse map[string]any
			json.Unmarshal([]byte(content), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_app_control_not_found.json"))
			var errorResponse map[string]any
			json.Unmarshal([]byte(content), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
