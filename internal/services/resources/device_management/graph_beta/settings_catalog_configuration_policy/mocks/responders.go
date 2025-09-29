package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	settingsCatalogConfigurationPolicies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.settingsCatalogConfigurationPolicies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("settings_catalog_configuration_policy", &SettingsCatalogConfigurationPolicyMock{})
}

// SettingsCatalogConfigurationPolicyMock provides mock responses for settings catalog configuration policy operations
type SettingsCatalogConfigurationPolicyMock struct{}

// Ensure SettingsCatalogConfigurationPolicyMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*SettingsCatalogConfigurationPolicyMock)(nil)

// RegisterMocks registers HTTP mock responses for settings catalog configuration policy operations
func (m *SettingsCatalogConfigurationPolicyMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.settingsCatalogConfigurationPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing configuration policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			policies := make([]map[string]any, 0, len(mockState.settingsCatalogConfigurationPolicies))
			for _, policy := range mockState.settingsCatalogConfigurationPolicies {
				// Ensure @odata.type is present
				policyCopy := make(map[string]any)
				for k, v := range policy {
					policyCopy[k] = v
				}
				if _, hasODataType := policyCopy["@odata.type"]; !hasODataType {
					policyCopy["@odata.type"] = "#microsoft.graph.deviceManagementConfigurationPolicy"
				}

				// Check if expand=assignments is requested
				expandParam := req.URL.Query().Get("$expand")
				if strings.Contains(expandParam, "assignments") {
					// Include assignments if they exist in the policy data
					if assignments, hasAssignments := policy["assignments"]; hasAssignments && assignments != nil {
						if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
							policyCopy["assignments"] = assignments
						} else {
							// If assignments array is empty, return empty array (not null)
							policyCopy["assignments"] = []interface{}{}
						}
					} else {
						// If no assignments stored, return empty array (not null)
						policyCopy["assignments"] = []interface{}{}
					}
				}

				policies = append(policies, policyCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual configuration policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			policy, exists := mockState.settingsCatalogConfigurationPolicies[id]
			mockState.Unlock()

			if !exists {
				// Check for special test IDs and load corresponding response files
				switch id {
				case "00000000-0000-0000-0000-000000000001":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_minimal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000002":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_maximal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000003":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_all_assignments.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000004":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_group_assignments.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000005":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_all_devices.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000006":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_all_users.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				case "00000000-0000-0000-0000-000000000007":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_settings_catalog_configuration_policy_exclusion.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return factories.SuccessResponse(200, response)(req)
				default:
					errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_settings_catalog_configuration_policy_not_found.json"))
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}

			// Create response copy
			policyCopy := make(map[string]any)
			for k, v := range policy {
				policyCopy[k] = v
			}
			if _, hasODataType := policyCopy["@odata.type"]; !hasODataType {
				policyCopy["@odata.type"] = "#microsoft.graph.deviceManagementConfigurationPolicy"
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the policy data
				if assignments, hasAssignments := policy["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
						// Return assignments in Microsoft Graph SDK format (not transformed)
						// The SDK will handle the transformation to Terraform structure
						policyCopy["assignments"] = assignments
					} else {
						// If assignments array is empty, return empty array (not null)
						policyCopy["assignments"] = []interface{}{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					policyCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, policyCopy)
		})

	// Register POST for creating configuration policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a unique ID for the new policy
			id := uuid.New().String()

			// Create the policy object with required fields
			policy := map[string]any{
				"@odata.type":          "#microsoft.graph.deviceManagementConfigurationPolicy",
				"id":                   id,
				"name":                 requestBody["name"],
				"platforms":            requestBody["platforms"],
				"technologies":         requestBody["technologies"],
				"isAssigned":           false,
				"settingsCount":        0,
				"createdDateTime":      "2024-01-15T10:30:00Z",
				"lastModifiedDateTime": "2024-01-15T10:30:00Z",
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				policy["description"] = description
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				policy["roleScopeTagIds"] = roleScopeTagIds
			} else {
				policy["roleScopeTagIds"] = []string{"0"} // Default value
			}
			if templateReference, exists := requestBody["templateReference"]; exists {
				if templateRefMap, ok := templateReference.(map[string]any); ok {
					// Ensure all required fields are present with defaults
					templateRef := map[string]any{
						"templateId":             "",
						"templateFamily":         "none",
						"templateDisplayName":    "",
						"templateDisplayVersion": "",
					}
					// Override with any provided values
					for k, v := range templateRefMap {
						templateRef[k] = v
					}
					// Set default templateFamily if not provided or empty
					if templateFamily, hasFamily := templateRef["templateFamily"]; !hasFamily || templateFamily == "" {
						templateRef["templateFamily"] = "none"
					}
					policy["templateReference"] = templateRef
				} else {
					// Fallback for invalid template reference
					policy["templateReference"] = map[string]any{
						"templateId":             "",
						"templateFamily":         "none",
						"templateDisplayName":    "",
						"templateDisplayVersion": "",
					}
				}
			} else {
				// Provide a default template reference for minimal cases
				policy["templateReference"] = map[string]any{
					"templateId":             "",
					"templateFamily":         "none",
					"templateDisplayName":    "",
					"templateDisplayVersion": "",
				}
			}
			if settings, exists := requestBody["settings"]; exists {
				if settingsList, ok := settings.([]interface{}); ok {
					// Ensure each setting has an ID for proper state mapping
					processedSettings := make([]interface{}, 0, len(settingsList))
					for i, setting := range settingsList {
						if settingMap, ok := setting.(map[string]any); ok {
							// Create a copy of the setting
							settingCopy := make(map[string]any)
							for k, v := range settingMap {
								settingCopy[k] = v
							}
							// Ensure the setting has an ID if it doesn't already
							if _, hasId := settingCopy["id"]; !hasId {
								settingCopy["id"] = strconv.Itoa(i)
							}
							processedSettings = append(processedSettings, settingCopy)
						} else {
							processedSettings = append(processedSettings, setting)
						}
					}
					policy["settings"] = processedSettings
					policy["settingsCount"] = len(processedSettings)
				} else {
					policy["settings"] = settings
				}
			}

			// Initialize assignments as empty array
			policy["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.settingsCatalogConfigurationPolicies[id] = policy
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, policy)
		})

	// Register PUT for updating configuration policy using custom URL format configurationPolicies('id')
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[^']+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID between configurationPolicies('...')
			path := req.URL.Path
			startMarker := "configurationPolicies('"
			startIdx := strings.Index(path, startMarker)
			if startIdx == -1 {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL format"}}`), nil
			}
			startIdx += len(startMarker)
			endIdx := strings.Index(path[startIdx:], "')")
			if endIdx == -1 {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid URL format"}}`), nil
			}
			id := path[startIdx : startIdx+endIdx]

			// Decode request body
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update stored state if exists
			mockState.Lock()
			defer mockState.Unlock()
			if existing, ok := mockState.settingsCatalogConfigurationPolicies[id]; ok {
				// Preserve or assign IDs to settings if provided in request
				if settings, hasSettings := requestBody["settings"]; hasSettings {
					if settingsList, okList := settings.([]interface{}); okList {
						processedSettings := make([]interface{}, 0, len(settingsList))
						var existingSettings []interface{}
						if ex, okEx := existing["settings"].([]interface{}); okEx {
							existingSettings = ex
						}
						for i, setting := range settingsList {
							if settingMap, okMap := setting.(map[string]any); okMap {
								settingCopy := make(map[string]any)
								for k, v := range settingMap {
									settingCopy[k] = v
								}
								// If id missing, try to carry over from existing by index, else assign index
								if _, hasId := settingCopy["id"]; !hasId {
									assignedId := strconv.Itoa(i)
									if i < len(existingSettings) {
										if exMap, okExMap := existingSettings[i].(map[string]any); okExMap {
											if exId, okExId := exMap["id"].(string); okExId && exId != "" {
												assignedId = exId
											}
										}
									}
									settingCopy["id"] = assignedId
								}
								processedSettings = append(processedSettings, settingCopy)
							} else {
								processedSettings = append(processedSettings, setting)
							}
						}
						requestBody["settings"] = processedSettings
					}
				}

				// Merge request body into existing policy
				for k, v := range requestBody {
					existing[k] = v
				}
				mockState.settingsCatalogConfigurationPolicies[id] = existing
			} else {
				// If not found, return 404 consistent with other responders
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_settings_catalog_configuration_policy_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Custom PUT expects 204 No Content on success
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register PATCH for updating configuration policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Load update template
			updatedPolicy, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "patch_settings_catalog_configuration_policy_updated.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}

			mockState.Lock()
			policy, exists := mockState.settingsCatalogConfigurationPolicies[id]
			if !exists {
				mockState.Unlock()
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_settings_catalog_configuration_policy_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range policy {
				updatedPolicy[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedPolicy[k] = v
			}

			// Update last modified time
			updatedPolicy["lastModifiedDateTime"] = "2024-01-15T11:00:00Z"

			// Store updated state
			mockState.settingsCatalogConfigurationPolicies[id] = updatedPolicy
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedPolicy)(req)
		})

	// Register DELETE for configuration policy
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.settingsCatalogConfigurationPolicies[id]
			if exists {
				delete(mockState.settingsCatalogConfigurationPolicies, id)
			}
			mockState.Unlock()

			if !exists {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_settings_catalog_configuration_policy_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register settings-related endpoints
	m.registerSettingsMocks()

	// Register assignment-related endpoints
	m.registerAssignmentMocks()
}

// registerAssignmentMocks registers mock responses for assignment operations
func (m *SettingsCatalogConfigurationPolicyMock) registerAssignmentMocks() {
	// POST assignment for a configuration policy
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-2] // configurationPolicies/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the policy
			mockState.Lock()
			if policyData, exists := mockState.settingsCatalogConfigurationPolicies[policyId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								// Generate a unique assignment ID
								assignmentId := uuid.New().String()

								// For configuration policies, assignments come with a "target" wrapper
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

								// Handle assignment filters
								var filterInfo map[string]any
								if filterData, hasFilter := assignmentMap["deviceAndAppManagementAssignmentFilterType"]; hasFilter {
									filterInfo = make(map[string]any)
									filterInfo["deviceAndAppManagementAssignmentFilterType"] = filterData
									if filterId, hasFilterId := assignmentMap["deviceAndAppManagementAssignmentFilterId"]; hasFilterId {
										filterInfo["deviceAndAppManagementAssignmentFilterId"] = filterId
									}
								}

								// Keep original Microsoft Graph API field names for SDK processing
								// The SDK will handle the field name mapping to Terraform structure
								graphAssignment := map[string]any{
									"id":          assignmentId,
									"target":      target,
									"source":      "direct",
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyAssignment",
								}

								// Add filter information if present
								for k, v := range filterInfo {
									graphAssignment[k] = v
								}

								graphAssignments = append(graphAssignments, graphAssignment)
							}
						}
						policyData["assignments"] = graphAssignments
						policyData["isAssigned"] = len(graphAssignments) > 0
					} else {
						// Set empty assignments array instead of deleting
						policyData["assignments"] = []interface{}{}
						policyData["isAssigned"] = false
					}
				} else {
					// Set empty assignments array instead of deleting
					policyData["assignments"] = []interface{}{}
					policyData["isAssigned"] = false
				}
				mockState.settingsCatalogConfigurationPolicies[policyId] = policyData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// GET assignments for a configuration policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogConfigurationPolicies[id]
			mockState.Unlock()

			if !exists {
				response := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/assignments",
					"value":          []map[string]any{},
				}
				return httpmock.NewJsonResponse(200, response)
			}

			// Get assignments from stored policy data
			assignments := []interface{}{}
			if storedAssignments, hasAssignments := policyData["assignments"]; hasAssignments {
				if assignmentArray, ok := storedAssignments.([]interface{}); ok {
					assignments = assignmentArray
				}
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// registerSettingsMocks registers mock responses for settings operations
func (m *SettingsCatalogConfigurationPolicyMock) registerSettingsMocks() {
	// GET settings for a configuration policy
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-2]

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogConfigurationPolicies[id]
			mockState.Unlock()

			if !exists {
				response := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/settings",
					"value":          []map[string]any{},
				}
				return httpmock.NewJsonResponse(200, response)
			}

			// Get settings from stored policy data
			settings := []interface{}{}
			if storedSettings, hasSettings := policyData["settings"]; hasSettings {
				if settingsArray, ok := storedSettings.([]interface{}); ok {
					settings = settingsArray
				}
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/settings",
				"value":          settings,
			}

			return httpmock.NewJsonResponse(200, response)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *SettingsCatalogConfigurationPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored configuration policies
	for id := range mockState.settingsCatalogConfigurationPolicies {
		delete(mockState.settingsCatalogConfigurationPolicies, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *SettingsCatalogConfigurationPolicyMock) loadJSONResponse(filePath string) (map[string]any, error) {
	var response map[string]any

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *SettingsCatalogConfigurationPolicyMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.settingsCatalogConfigurationPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing configuration policies (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating configuration policy with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		})

	// Register error response for configuration policy not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_settings_catalog_configuration_policy_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
