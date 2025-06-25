package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	settingsCatalogs map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.settingsCatalogs = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// SettingsCatalogMock provides mock responses for settings catalog operations
type SettingsCatalogMock struct{}

// RegisterMocks registers HTTP mock responses for settings catalog operations
func (m *SettingsCatalogMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.settingsCatalogs = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register test settings catalog policies
	registerTestSettingsCatalogs()

	// Register GET for settings catalog by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogs[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings catalog policy not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register GET for listing settings catalogs
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			policies := make([]map[string]interface{}, 0, len(mockState.settingsCatalogs))
			for _, policy := range mockState.settingsCatalogs {
				policies = append(policies, policy)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating settings catalogs
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var policyData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&policyData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := policyData["name"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"name is required"}}`), nil
			}
			if _, ok := policyData["description"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"description is required"}}`), nil
			}
			if _, ok := policyData["platforms"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"platforms is required"}}`), nil
			}

			// Generate ID if not provided
			if policyData["id"] == nil {
				policyData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			policyData["createdDateTime"] = now
			policyData["lastModifiedDateTime"] = now

			// Ensure settings field is initialized
			if policyData["settings"] == nil {
				policyData["settings"] = []map[string]interface{}{}
			}

			// Store policy in mock state
			policyId := policyData["id"].(string)
			mockState.Lock()
			mockState.settingsCatalogs[policyId] = policyData
			mockState.Unlock()

			// Register assignments endpoint for this policy
			registerAssignmentsEndpoints(policyId)

			return httpmock.NewJsonResponse(201, policyData)
		})

	// Register PATCH for updating settings catalogs
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogs[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings catalog policy not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update policy data
			mockState.Lock()

			// Special handling for updates that remove fields
			// If we're updating from maximal to minimal, we need to remove fields not in the minimal config
			// Check if this is a minimal update by looking for key indicators
			isMinimalUpdate := false
			if _, hasName := updateData["name"]; hasName {
				if _, hasSettings := updateData["settings"]; !hasSettings {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove fields that are not part of minimal configuration
				policyData["settings"] = []map[string]interface{}{}
				delete(policyData, "technologies")
				// Keep assignments as they're managed separately
			}

			// Apply the updates
			for k, v := range updateData {
				policyData[k] = v
			}

			// Update last modified time
			policyData["lastModifiedDateTime"] = time.Now().Format(time.RFC3339)

			mockState.settingsCatalogs[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register DELETE for removing settings catalogs
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.settingsCatalogs[policyId]
			if exists {
				delete(mockState.settingsCatalogs, policyId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// registerAssignmentsEndpoints registers endpoints for handling assignments for a specific policy
func registerAssignmentsEndpoints(policyId string) {
	// GET assignments
	httpmock.RegisterResponder("GET",
		`=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/`+policyId+`/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			policyData, exists := mockState.settingsCatalogs[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings catalog policy not found"}}`), nil
			}

			// Return assignments if they exist
			assignments := []map[string]interface{}{}
			if assignmentsData, ok := policyData["assignments"].([]map[string]interface{}); ok {
				assignments = assignmentsData
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies(" + policyId + ")/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// POST assignments
	httpmock.RegisterResponder("POST",
		`=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/`+policyId+`/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			var assignmentData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&assignmentData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogs[policyId]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings catalog policy not found"}}`), nil
			}

			// Add ID to assignment if not present
			if assignmentData["id"] == nil {
				assignmentData["id"] = uuid.New().String()
			}

			// Initialize assignments array if it doesn't exist
			if policyData["assignments"] == nil {
				policyData["assignments"] = []map[string]interface{}{}
			}

			// Add the assignment
			assignments := policyData["assignments"].([]map[string]interface{})
			assignments = append(assignments, assignmentData)
			policyData["assignments"] = assignments

			mockState.settingsCatalogs[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, assignmentData)
		})

	// DELETE assignments
	httpmock.RegisterResponder("DELETE",
		`=~^https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/`+policyId+`/assignments/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.settingsCatalogs[policyId]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Settings catalog policy not found"}}`), nil
			}

			// Remove the assignment with matching ID
			if assignments, ok := policyData["assignments"].([]map[string]interface{}); ok {
				updatedAssignments := []map[string]interface{}{}
				for _, assignment := range assignments {
					if id, ok := assignment["id"].(string); !ok || id != assignmentId {
						updatedAssignments = append(updatedAssignments, assignment)
					}
				}
				policyData["assignments"] = updatedAssignments
				mockState.settingsCatalogs[policyId] = policyData
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *SettingsCatalogMock) RegisterErrorMocks() {
	// Register error response for policy creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		factories.ErrorResponse(400, "BadRequest", "Settings catalog policy creation failed"))

	// Register error response for policy not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/not-found-policy",
		factories.ErrorResponse(404, "ResourceNotFound", "Settings catalog policy not found"))
}

// registerTestSettingsCatalogs registers predefined test settings catalogs
func registerTestSettingsCatalogs() {
	// Minimal settings catalog with only required attributes
	minimalPolicyId := "00000000-0000-0000-0000-000000000001"
	minimalPolicyData := map[string]interface{}{
		"id":                   minimalPolicyId,
		"name":                 "Minimal Settings Catalog",
		"description":          "Minimal settings catalog policy",
		"platforms":            "windows10",
		"technologies":         "mdm",
		"createdDateTime":      "2023-01-01T00:00:00Z",
		"lastModifiedDateTime": "2023-01-01T00:00:00Z",
		"settings":             []map[string]interface{}{},
		"assignments":          []map[string]interface{}{},
	}

	// Maximal settings catalog with all attributes
	maximalPolicyId := "00000000-0000-0000-0000-000000000002"
	maximalPolicyData := map[string]interface{}{
		"id":                   maximalPolicyId,
		"name":                 "Maximal Settings Catalog",
		"description":          "Maximal settings catalog policy with all options",
		"platforms":            "windows10",
		"technologies":         "mdm",
		"createdDateTime":      "2023-01-01T00:00:00Z",
		"lastModifiedDateTime": "2023-01-01T00:00:00Z",
		"settings": []map[string]interface{}{
			{
				"@odata.type": "#microsoft.graph.deviceManagementConfigurationSetting",
				"settingInstance": map[string]interface{}{
					"@odata.type":         "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
					"settingDefinitionId": "device_vendor_msft_policy_config_defender_allowarchivescanning",
					"choiceSettingValue": map[string]interface{}{
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
						"value":       "device_vendor_msft_policy_config_defender_allowarchivescanning_1",
						"children":    []map[string]interface{}{},
					},
				},
			},
			{
				"@odata.type": "#microsoft.graph.deviceManagementConfigurationSetting",
				"settingInstance": map[string]interface{}{
					"@odata.type":         "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
					"settingDefinitionId": "device_vendor_msft_policy_config_defender_allowbehaviormonitoring",
					"choiceSettingValue": map[string]interface{}{
						"@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
						"value":       "device_vendor_msft_policy_config_defender_allowbehaviormonitoring_1",
						"children":    []map[string]interface{}{},
					},
				},
			},
		},
		"assignments": []map[string]interface{}{
			{
				"id":     "00000000-0000-0000-0000-000000000003",
				"intent": "apply",
				"target": map[string]interface{}{
					"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget",
					"deviceAndAppManagementAssignmentFilterId":   nil,
					"deviceAndAppManagementAssignmentFilterType": "none",
				},
			},
		},
	}

	// Store in mock state
	mockState.Lock()
	mockState.settingsCatalogs[minimalPolicyId] = minimalPolicyData
	mockState.settingsCatalogs[maximalPolicyId] = maximalPolicyData
	mockState.Unlock()

	// Register assignments endpoints for these policies
	registerAssignmentsEndpoints(minimalPolicyId)
	registerAssignmentsEndpoints(maximalPolicyId)
}
