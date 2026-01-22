package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policies    map[string]map[string]any
	settings    map[string]map[string]any
	assignments map[string]map[string]any
}

func init() {
	mockState.policies = make(map[string]map[string]any)
	mockState.settings = make(map[string]map[string]any)
	mockState.assignments = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("settings_catalog_configuration_policy_json", &SettingsCatalogPolicyMock{})
}

type SettingsCatalogPolicyMock struct{}

var _ mocks.MockRegistrar = (*SettingsCatalogPolicyMock)(nil)

// getJSONFileForPolicyName determines which JSON file to load based on the policy's name
// For update tests, the isUpdate parameter indicates if this is an UPDATE operation (vs initial CREATE)
func getJSONFileForPolicyName(name string, isUpdate bool) string {
	nameUpper := strings.ToUpper(name)
	
	// Handle update tests specially - they need different JSON for initial vs updated state
	if strings.Contains(nameUpper, "UPDATE") && strings.Contains(nameUpper, "REVERSE") {
		// Test 19: maximal -> minimal
		if isUpdate {
			// After update: minimal (file explorer)
			return "post_settings_catalog_policy_file_explorer_success.json"
		}
		// Initial state: maximal (local policies)
		return "post_settings_catalog_policy_local_policies_security_maximal_success.json"
	}
	if strings.Contains(nameUpper, "UPDATE") {
		// Test 18: minimal -> maximal
		if isUpdate {
			// After update: maximal (local policies)
			return "post_settings_catalog_policy_local_policies_security_maximal_success.json"
		}
		// Initial state: minimal (file explorer)
		return "post_settings_catalog_policy_file_explorer_success.json"
	}
	
	// Simple policies
	if strings.Contains(nameUpper, "CAMERA") {
		return "post_settings_catalog_policy_camera_success.json"
	}
	if strings.Contains(nameUpper, "TASK MANAGER") {
		return "post_settings_catalog_policy_task_manager_success.json"
	}
	if strings.Contains(nameUpper, "APP PRIVACY") {
		return "post_settings_catalog_policy_app_privacy_success.json"
	}
	if strings.Contains(nameUpper, "CRYPTOGRAPHY") {
		return "post_settings_catalog_policy_cryptography_success.json"
	}
	if strings.Contains(nameUpper, "NOTIFICATIONS") {
		return "post_settings_catalog_policy_notifications_success.json"
	}
	if strings.Contains(nameUpper, "ATTACHMENT MANAGER") {
		return "post_settings_catalog_policy_attachment_manager_success.json"
	}
	if strings.Contains(nameUpper, "CREDENTIAL USER INTERFACE") {
		return "post_settings_catalog_policy_credential_user_interface_success.json"
	}
	
	// Medium complexity
	if strings.Contains(nameUpper, "REMOTE DESKTOP AVD") || strings.Contains(nameUpper, "AVD URL") {
		return "post_settings_catalog_policy_remote_desktop_avd_url_success.json"
	}
	if strings.Contains(nameUpper, "STORAGE SENSE") {
		return "post_settings_catalog_policy_storage_sense_success.json"
	}
	if strings.Contains(nameUpper, "WINDOWS CONNECTION MANAGER") {
		return "post_settings_catalog_policy_windows_connection_manager_success.json"
	}
	if strings.Contains(nameUpper, "AUTOPLAY") {
		return "post_settings_catalog_policy_autoplay_policies_success.json"
	}
	if strings.Contains(nameUpper, "DEFENDER SMARTSCREEN") || strings.Contains(nameUpper, "SMARTSCREEN") {
		return "post_settings_catalog_policy_defender_smartscreen_success.json"
	}
	if strings.Contains(nameUpper, "EDGE") && strings.Contains(nameUpper, "EXTENSIONS") {
		return "post_settings_catalog_policy_edge_extensions_macos_success.json"
	}
	if strings.Contains(nameUpper, "EDGE") && strings.Contains(nameUpper, "SECURITY") {
		return "post_settings_catalog_policy_edge_security_macos_success.json"
	}
	if strings.Contains(nameUpper, "ONEDRIVE") && strings.Contains(nameUpper, "KNOWN FOLDER") {
		return "post_settings_catalog_policy_onedrive_known_folder_macos_success.json"
	}
	if strings.Contains(nameUpper, "OFFICE") && strings.Contains(nameUpper, "CONFIGURATION") {
		return "post_settings_catalog_policy_office_configuration_macos_success.json"
	}
	
	// Complex nested structures
	if strings.Contains(nameUpper, "DEFENDER ANTIVIRUS") && strings.Contains(nameUpper, "SECURITY BASELINE") {
		return "post_settings_catalog_policy_defender_antivirus_baseline_success.json"
	}
	if strings.Contains(nameUpper, "DEFENDER") && strings.Contains(nameUpper, "CONTROLLED FOLDER") {
		return "post_settings_catalog_policy_defender_controlled_folder_success.json"
	}
	if strings.Contains(nameUpper, "DELIVERY OPTIMIZATION") {
		return "post_settings_catalog_policy_delivery_optimization_success.json"
	}
	if strings.Contains(nameUpper, "FILEVAULT") || strings.Contains(nameUpper, "FILE VAULT") {
		return "post_settings_catalog_policy_filevault_macos_success.json"
	}
	if strings.Contains(nameUpper, "IOS") && strings.Contains(nameUpper, "ACCOUNTS") {
		return "post_settings_catalog_policy_ios_accounts_success.json"
	}
	if strings.Contains(nameUpper, "WINRM") || strings.Contains(nameUpper, "WINDOWS REMOTE MANAGEMENT") {
		return "post_settings_catalog_policy_winrm_success.json"
	}
	if strings.Contains(nameUpper, "WIN365") || strings.Contains(nameUpper, "RESOURCE REDIRECTION") {
		return "post_settings_catalog_policy_win365_resource_redirection_success.json"
	}
	if strings.Contains(nameUpper, "FILE EXPLORER") {
		return "post_settings_catalog_policy_file_explorer_success.json"
	}
	if strings.Contains(nameUpper, "LOCAL POLICIES") && strings.Contains(nameUpper, "SECURITY OPTIONS") {
		return "post_settings_catalog_policy_local_policies_security_maximal_success.json"
	}
	
	// Default fallback
	return "post_settings_catalog_policy_minimal_success.json"
}

func (m *SettingsCatalogPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.settings = make(map[string]map[string]any)
	mockState.assignments = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreatePolicyMock()
	m.registerGetPolicyMock()
	m.registerGetSettingsMocks()
	m.registerUpdatePolicyMocks()
	m.registerDeletePolicyMock()
	m.registerAssignmentMocks()
	m.registerMockGroups()
}

// registerCreatePolicyMock registers the POST responder for creating configuration policies
func (m *SettingsCatalogPolicyMock) registerCreatePolicyMock() {
	// Create settings catalog policy - POST /deviceManagement/configurationPolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Invalid request body: %s"}}`, err.Error())), nil
		}

		// Determine which JSON file to load based on name
		name, ok := requestBody["name"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"name is required"}}`), nil
		}

		// For POST (create), isUpdate is always false
		jsonFileName := getJSONFileForPolicyName(name, false)
		responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file %s: %s"}}`, responsesPath, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate a UUID for the new resource and update the response
		newId := uuid.New().String()
		responseObj["id"] = newId
		
		// Update name and description from request
		responseObj["name"] = name
		if desc, ok := requestBody["description"].(string); ok {
			responseObj["description"] = desc
		}
		
		// Extract settings from response and store separately
		if settings, ok := responseObj["settings"].([]any); ok {
			settingsObj := map[string]any{
				"value": settings,
			}
			mockState.Lock()
			mockState.settings[newId] = settingsObj
			mockState.Unlock()
			
			// Remove settings from main response as they're fetched separately
			delete(responseObj, "settings")
		}

		mockState.Lock()
		mockState.policies[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})
}

// registerGetPolicyMock registers the GET responder for retrieving a single configuration policy
func (m *SettingsCatalogPolicyMock) registerGetPolicyMock() {
	// Get settings catalog policy - GET /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Policy not found"}}`), nil
		}

		// Return the stored policy data (without settings)
		return httpmock.NewJsonResponse(200, policy)
	})
}

// registerGetSettingsMocks registers GET responders for retrieving policy settings
func (m *SettingsCatalogPolicyMock) registerGetSettingsMocks() {
	// Get settings catalog policy settings - GET /deviceManagement/configurationPolicies/{id}/settings (slash format)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/settings`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		var policyId string
		for i, part := range parts {
			if part == "settings" && i > 0 {
				policyId = parts[i-1]
				break
			}
		}

		mockState.Lock()
		settings, exists := mockState.settings[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings not found"}}`), nil
		}

		// Return the stored settings data
		return httpmock.NewJsonResponse(200, settings)
	})
	
	// Get settings catalog policy settings - GET /deviceManagement/configurationPolicies('id')/settings (OData format)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'\)/settings`, func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		var policyId string
		if strings.Contains(path, "('") && strings.Contains(path, "')") {
			start := strings.Index(path, "('") + 2
			end := strings.Index(path, "')")
			policyId = path[start:end]
		}

		mockState.Lock()
		settings, exists := mockState.settings[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings not found"}}`), nil
		}

		// Return the stored settings data
		return httpmock.NewJsonResponse(200, settings)
	})
}

// registerUpdatePolicyMocks registers PUT responders for updating configuration policies
func (m *SettingsCatalogPolicyMock) registerUpdatePolicyMocks() {
	// Update settings catalog policy - PUT /deviceManagement/configurationPolicies/{id} (slash format)
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request (excluding settings)
		for key, value := range requestBody {
			if key != "settings" {
				policy[key] = value
			}
		}
		policy["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		
		// For update operations, load the updated JSON to get properly formatted settings
		if _, ok := requestBody["settings"]; ok {
			policyName, _ := policy["name"].(string)
			jsonFileName := getJSONFileForPolicyName(policyName, true)
			responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
			
			if jsonData, err := os.ReadFile(responsesPath); err == nil {
				var updateResponseObj map[string]any
				if err := json.Unmarshal(jsonData, &updateResponseObj); err == nil {
					if updatedSettings, ok := updateResponseObj["settings"].([]any); ok {
						policy["settingCount"] = len(updatedSettings)
						mockState.settings[policyId] = map[string]any{
							"value": updatedSettings,
						}
					}
				}
			}
		}
		
		mockState.policies[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
	
	// Update settings catalog policy - PUT /deviceManagement/configurationPolicies('id') (OData format)
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'\)$`, func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		var policyId string
		if strings.Contains(path, "('") && strings.Contains(path, "')") {
			start := strings.Index(path, "('") + 2
			end := strings.Index(path, "')")
			policyId = path[start:end]
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request (excluding settings)
		for key, value := range requestBody {
			if key != "settings" {
				policy[key] = value
			}
		}
		policy["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		
		// For update operations, load the updated JSON to get properly formatted settings
		if _, ok := requestBody["settings"]; ok {
			policyName, _ := policy["name"].(string)
			jsonFileName := getJSONFileForPolicyName(policyName, true)
			responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
			
			if jsonData, err := os.ReadFile(responsesPath); err == nil {
				var updateResponseObj map[string]any
				if err := json.Unmarshal(jsonData, &updateResponseObj); err == nil {
					if updatedSettings, ok := updateResponseObj["settings"].([]any); ok {
						policy["settingCount"] = len(updatedSettings)
						mockState.settings[policyId] = map[string]any{
							"value": updatedSettings,
						}
					}
				}
			}
		}
		
		mockState.policies[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

// registerAssignmentMocks registers responders for assignment operations
func (m *SettingsCatalogPolicyMock) registerAssignmentMocks() {
	// Assign settings catalog policy - POST /deviceManagement/configurationPolicies/{id}/assign
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		var policyId string
		for i, part := range parts {
			if part == "assign" && i > 0 {
				policyId = parts[i-1]
				break
			}
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Process assignments - transform from Terraform format to Graph API format
		if assignmentsArray, ok := requestBody["assignments"].([]any); ok {
			graphAssignments := []any{}
			for _, assignment := range assignmentsArray {
				if assignmentMap, ok := assignment.(map[string]any); ok {
					// Generate a unique assignment ID
					assignmentId := uuid.New().String()
					
					// Extract target data
					var target map[string]any
					if targetData, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
						target = make(map[string]any)
						for k, v := range targetData {
							target[k] = v
						}
					}
					
					// Create Graph API formatted assignment
					graphAssignment := map[string]any{
						"id":          assignmentId,
						"target":      target,
						"source":      "direct",
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyAssignment",
					}
					
					// Add filter information if present
					if filterType, hasFilterType := assignmentMap["deviceAndAppManagementAssignmentFilterType"]; hasFilterType {
						graphAssignment["deviceAndAppManagementAssignmentFilterType"] = filterType
					}
					if filterId, hasFilterId := assignmentMap["deviceAndAppManagementAssignmentFilterId"]; hasFilterId {
						graphAssignment["deviceAndAppManagementAssignmentFilterId"] = filterId
					}
					
					graphAssignments = append(graphAssignments, graphAssignment)
				}
			}
			
			// Store assignments wrapped in value array
			mockState.assignments[policyId] = map[string]any{
				"value": graphAssignments,
			}
			
			// Update policy's isAssigned field
			policy["isAssigned"] = len(graphAssignments) > 0
			mockState.policies[policyId] = policy
		} else {
			// No assignments - set to empty
			mockState.assignments[policyId] = map[string]any{
				"value": []any{},
			}
			policy["isAssigned"] = false
			mockState.policies[policyId] = policy
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
	})
	
	// Get settings catalog policy assignments - GET /deviceManagement/configurationPolicies/{id}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/assignments$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		var policyId string
		for i, part := range parts {
			if part == "assignments" && i > 0 {
				policyId = parts[i-1]
				break
			}
		}

		mockState.Lock()
		assignments, exists := mockState.assignments[policyId]
		mockState.Unlock()

		if !exists {
			// Return empty assignments if none exist
			return httpmock.NewJsonResponse(200, map[string]any{
				"value": []any{},
			})
		}

		// Return the stored assignments (already wrapped in value array)
		return httpmock.NewJsonResponse(200, assignments)
	})
}

// registerDeletePolicyMock registers the DELETE responder for removing configuration policies
func (m *SettingsCatalogPolicyMock) registerDeletePolicyMock() {
	// Delete settings catalog policy - DELETE /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.policies[policyId]; exists {
			delete(mockState.policies, policyId)
			delete(mockState.settings, policyId)
			delete(mockState.assignments, policyId)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
	})
}

func (m *SettingsCatalogPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *SettingsCatalogPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.settings = make(map[string]map[string]any)
	mockState.assignments = make(map[string]map[string]any)
	mockState.Unlock()
}

// registerMockGroups registers mock group resources for unit tests
func (m *SettingsCatalogPolicyMock) registerMockGroups() {
	// Mock group creation - returns a UUID for any group name
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		newId := uuid.New().String()
		response := map[string]any{
			"id":          newId,
			"displayName": requestBody["displayName"],
		}

		return httpmock.NewJsonResponse(201, response)
	})

	// Mock group GET - return a valid response for any group ID
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		response := map[string]any{
			"id":          groupId,
			"displayName": "Mock Group",
		}

		return httpmock.NewJsonResponse(200, response)
	})
	
	// Mock group DELETE - for acceptance tests
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(204, ""), nil
	})
}
