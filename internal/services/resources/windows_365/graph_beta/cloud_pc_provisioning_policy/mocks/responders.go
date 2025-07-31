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
	provisioningPolicies map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.provisioningPolicies = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// CloudPcProvisioningPolicyMock provides mock responses for cloud pc provisioning policy operations
type CloudPcProvisioningPolicyMock struct{}

// RegisterMocks registers HTTP mock responses for cloud pc provisioning policy operations
func (m *CloudPcProvisioningPolicyMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.provisioningPolicies = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing provisioning policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			policies := make([]map[string]interface{}, 0, len(mockState.provisioningPolicies))
			for _, policy := range mockState.provisioningPolicies {
				policies = append(policies, policy)
			}
			mockState.Unlock()

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/provisioningPolicies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual provisioning policy
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.provisioningPolicies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Provisioning policy not found"}}`), nil
			}

			// Create response copy
			responseCopy := make(map[string]interface{})
			for k, v := range policyData {
				responseCopy[k] = v
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the policy data
				if assignments, hasAssignments := policyData["assignments"]; hasAssignments && assignments != nil {
					responseCopy["assignments"] = assignments
				} else {
					// If no assignments stored, return empty array (not null)
					responseCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, responseCopy)
		})

	// Register POST for creating provisioning policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new provisioning policy ID
			policyId := uuid.New().String()

			// Create provisioning policy data - only include fields that were provided or have defaults
			policyData := map[string]interface{}{
				"id":          policyId,
				"displayName": requestBody["displayName"],
				"imageId":     requestBody["imageId"],
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				policyData["description"] = description
			}
			if cloudPcNamingTemplate, exists := requestBody["cloudPcNamingTemplate"]; exists {
				policyData["cloudPcNamingTemplate"] = cloudPcNamingTemplate
			}
			if provisioningType, exists := requestBody["provisioningType"]; exists {
				policyData["provisioningType"] = provisioningType
			} else {
				policyData["provisioningType"] = "dedicated" // Default value
			}
			if enableSingleSignOn, exists := requestBody["enableSingleSignOn"]; exists {
				policyData["enableSingleSignOn"] = enableSingleSignOn
			}
			if localAdminEnabled, exists := requestBody["localAdminEnabled"]; exists {
				policyData["localAdminEnabled"] = localAdminEnabled
			}
			if imageType, exists := requestBody["imageType"]; exists {
				policyData["imageType"] = imageType
			} else {
				policyData["imageType"] = "gallery" // Default value
			}
			if managedBy, exists := requestBody["managedBy"]; exists {
				policyData["managedBy"] = managedBy
			} else {
				policyData["managedBy"] = "windows365" // Default value
			}
			if scopeIds, exists := requestBody["scopeIds"]; exists {
				policyData["scopeIds"] = scopeIds
			} else {
				policyData["scopeIds"] = []string{"0"} // Default value
			}

			// Add computed fields that are always returned by the API
			policyData["gracePeriodInHours"] = 4

			// Handle nested attributes - always set if provided (including empty arrays)
			if domainJoinConfigs, exists := requestBody["domainJoinConfigurations"]; exists {
				policyData["domainJoinConfigurations"] = domainJoinConfigs
			}

			if windowsSetting, exists := requestBody["windowsSetting"]; exists {
				policyData["windowsSetting"] = windowsSetting
			}

			if microsoftManagedDesktop, exists := requestBody["microsoftManagedDesktop"]; exists {
				policyData["microsoftManagedDesktop"] = microsoftManagedDesktop
			}

			if autopatch, exists := requestBody["autopatch"]; exists {
				policyData["autopatch"] = autopatch
			}

			if autopilotConfig, exists := requestBody["autopilotConfiguration"]; exists {
				// Ensure all autopilot fields are preserved
				policyData["autopilotConfiguration"] = autopilotConfig
			}

			if applyToExisting, exists := requestBody["applyToExistingCloudPcs"]; exists {
				policyData["applyToExistingCloudPcs"] = applyToExisting
			}

			// Store in mock state
			mockState.Lock()
			mockState.provisioningPolicies[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, policyData)
		})

	// Register PATCH for updating provisioning policy
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.provisioningPolicies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Provisioning policy not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update provisioning policy data
			mockState.Lock()

			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior

			// For nested attributes, if they're not in the request, remove them
			optionalNestedFields := []string{"windowsSetting", "microsoftManagedDesktop", "autopatch", "autopilotConfiguration", "domainJoinConfigurations", "applyToExistingCloudPcs"}
			for _, field := range optionalNestedFields {
				if _, hasField := requestBody[field]; !hasField {
					delete(policyData, field)
				}
			}

			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(policyData, key)
				} else {
					policyData[key] = value
				}
			}
			// Ensure the ID is preserved
			policyData["id"] = policyId
			mockState.provisioningPolicies[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register DELETE for removing provisioning policy
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.provisioningPolicies[policyId]
			if exists {
				delete(mockState.provisioningPolicies, policyId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Provisioning policy not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-3] // provisioningPolicies/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the policy
			mockState.Lock()
			if policyData, exists := mockState.provisioningPolicies[policyId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Create assignment in the format the API returns
						graphAssignments := []interface{}{
							map[string]interface{}{
								"id": "44444444-4444-4444-4444-444444444444",
								"target": map[string]interface{}{
									"@odata.type": "#microsoft.graph.cloudPcManagementGroupAssignmentTarget",
									"additionalData": map[string]interface{}{
										"groupId": "44444444-4444-4444-4444-444444444444",
									},
								},
							},
						}
						policyData["assignments"] = graphAssignments
					} else {
						delete(policyData, "assignments")
					}
				} else {
					// Remove assignments if none provided
					delete(policyData, "assignments")
				}
				mockState.provisioningPolicies[policyId] = policyData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/provisioningPolicies/assignments",
				"value":          []map[string]interface{}{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for apply operations
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/[^/]+/apply$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Dynamic mocks will handle all test cases
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *CloudPcProvisioningPolicyMock) RegisterErrorMocks() {
	// Register GET for listing provisioning policies (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/provisioningPolicies",
				"value":          []map[string]interface{}{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating provisioning policy with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for provisioning policy not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/not-found-policy",
		factories.ErrorResponse(404, "ResourceNotFound", "Provisioning policy not found"))
}

// getOrDefault returns the value from the map or a default value if the key doesn't exist
func getOrDefault(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}
