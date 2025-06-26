package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	policies map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.policies = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// ConditionalAccessPolicyMock provides mock responses for conditional access policy operations
type ConditionalAccessPolicyMock struct{}

// RegisterMocks registers HTTP mock responses for conditional access policy operations
func (m *ConditionalAccessPolicyMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.policies = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register specific test policies
	registerTestPolicies()

	// Register GET for policy by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/identity/conditionalAccess/policies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.policies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Conditional access policy not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register GET for listing policies
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/identity/conditionalAccess/policies(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			policies := make([]map[string]interface{}, 0, len(mockState.policies))
			for _, policy := range mockState.policies {
				policies = append(policies, policy)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/policies",
				"value":          policies,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating policies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies",
		func(req *http.Request) (*http.Response, error) {
			var policyData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&policyData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := policyData["displayName"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
			}
			if _, ok := policyData["state"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"state is required"}}`), nil
			}

			// Generate ID if not provided
			if policyData["id"] == nil {
				policyData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			policyData["createdDateTime"] = now
			policyData["modifiedDateTime"] = now

			// Ensure required nested structures exist
			ensurePolicyStructures(policyData)

			// Store policy in mock state
			policyId := policyData["id"].(string)
			mockState.Lock()
			mockState.policies[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, policyData)
		})

	// Register PATCH for updating policies
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/identity/conditionalAccess/policies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			policyData, exists := mockState.policies[policyId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Conditional access policy not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update policy data
			mockState.Lock()

			// Check if this is a minimal update (simplified logic)
			isMinimalUpdate := false
			if displayName, ok := updateData["displayName"].(string); ok {
				if strings.Contains(displayName, "Minimal") {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove complex fields for minimal update
				if conditions, ok := policyData["conditions"].(map[string]interface{}); ok {
					// Remove optional condition blocks
					delete(conditions, "platforms")
					delete(conditions, "locations")
					delete(conditions, "devices")
					delete(conditions, "userRiskLevels")
					delete(conditions, "signInRiskLevels")
					delete(conditions, "servicePrincipalRiskLevels")
				}
				// Remove session controls for minimal update
				delete(policyData, "sessionControls")
			}

			// Apply the updates
			for k, v := range updateData {
				policyData[k] = v
			}

			// Update modified timestamp
			policyData["modifiedDateTime"] = time.Now().Format(time.RFC3339)

			// Ensure structures are consistent
			ensurePolicyStructures(policyData)

			mockState.policies[policyId] = policyData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, policyData)
		})

	// Register DELETE for removing policies
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/identity/conditionalAccess/policies/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			policyId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.policies[policyId]
			if exists {
				delete(mockState.policies, policyId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *ConditionalAccessPolicyMock) RegisterErrorMocks() {
	// Register a responder that returns 409 Conflict for duplicate display names
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies",
		func(req *http.Request) (*http.Response, error) {
			var policyData map[string]interface{}
			json.NewDecoder(req.Body).Decode(&policyData)

			if displayName, ok := policyData["displayName"].(string); ok {
				if strings.Contains(displayName, "Error") {
					return httpmock.NewStringResponse(409, `{"error":{"code":"Conflict","message":"A conditional access policy with this display name already exists"}}`), nil
				}
			}

			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`), nil
		})
}

// registerTestPolicies registers predefined test policies
func registerTestPolicies() {
	// Register minimal test policy
	minimalPolicy := map[string]interface{}{
		"id":               "minimal-policy-id-12345",
		"displayName":      "Block Legacy Authentication - Minimal",
		"state":            "enabled",
		"createdDateTime":  "2024-01-01T00:00:00Z",
		"modifiedDateTime": "2024-01-01T00:00:00Z",
		"conditions": map[string]interface{}{
			"clientAppTypes":   []string{"exchangeActiveSync", "other"},
			"userRiskLevels":   []string{},
			"signInRiskLevels": []string{},
			"applications": map[string]interface{}{
				"includeApplications": []string{"All"},
				"excludeApplications": []string{},
				"includeUserActions":  []string{},
			},
			"users": map[string]interface{}{
				"includeUsers":  []string{"All"},
				"excludeUsers":  []string{},
				"includeGroups": []string{},
				"excludeGroups": []string{},
			},
		},
		"grantControls": map[string]interface{}{
			"operator":        "OR",
			"builtInControls": []string{"block"},
		},
	}

	// Register maximal test policy
	maximalPolicy := map[string]interface{}{
		"id":               "maximal-policy-id-67890",
		"displayName":      "Comprehensive Security Policy - Maximal",
		"state":            "enabled",
		"createdDateTime":  "2024-01-01T00:00:00Z",
		"modifiedDateTime": "2024-01-01T00:00:00Z",
		"conditions": map[string]interface{}{
			"clientAppTypes":   []string{"all"},
			"userRiskLevels":   []string{"high"},
			"signInRiskLevels": []string{"high", "medium"},
			"applications": map[string]interface{}{
				"includeApplications": []string{"All"},
				"excludeApplications": []string{"00000002-0000-0ff1-ce00-000000000000"},
				"includeUserActions":  []string{},
				"applicationFilter": map[string]interface{}{
					"mode": "exclude",
					"rule": "device.deviceOwnership -eq \"Company\"",
				},
			},
			"users": map[string]interface{}{
				"includeUsers":  []string{"All"},
				"excludeUsers":  []string{"GuestsOrExternalUsers"},
				"includeGroups": []string{},
				"excludeGroups": []string{},
				"includeRoles":  []string{},
				"excludeRoles":  []string{"62e90394-69f5-4237-9190-012177145e10"},
			},
			"platforms": map[string]interface{}{
				"includePlatforms": []string{"all"},
				"excludePlatforms": []string{},
			},
			"locations": map[string]interface{}{
				"includeLocations": []string{"All"},
				"excludeLocations": []string{"AllTrusted"},
			},
			"devices": map[string]interface{}{
				"includeDevices":      []string{},
				"excludeDevices":      []string{},
				"includeDeviceStates": []string{},
				"excludeDeviceStates": []string{},
				"deviceFilter": map[string]interface{}{
					"mode": "include",
					"rule": "device.isCompliant -eq True",
				},
			},
		},
		"grantControls": map[string]interface{}{
			"operator":                    "AND",
			"builtInControls":             []string{"mfa", "compliantDevice"},
			"customAuthenticationFactors": []string{},
			"termsOfUse":                  []string{},
			"authenticationStrength": map[string]interface{}{
				"id":                    "00000000-0000-0000-0000-000000000004",
				"displayName":           "Multifactor authentication",
				"description":           "Combinations of methods that satisfy strong authentication, such as a password + SMS",
				"policyType":            "builtIn",
				"requirementsSatisfied": "mfa",
				"allowedCombinations": []string{
					"password,sms",
					"password,voice",
					"password,hardwareOath",
					"password,softwareOath",
					"password,microsoftAuthenticatorPush",
				},
			},
		},
		"sessionControls": map[string]interface{}{
			"disableResilienceDefaults": false,
			"applicationEnforcedRestrictions": map[string]interface{}{
				"isEnabled": true,
			},
			"cloudAppSecurity": map[string]interface{}{
				"isEnabled":            true,
				"cloudAppSecurityType": "monitorOnly",
			},
			"signInFrequency": map[string]interface{}{
				"isEnabled":          true,
				"type":               "hours",
				"value":              4,
				"authenticationType": "primaryAndSecondaryAuthentication",
				"frequencyInterval":  "timeBased",
			},
			"persistentBrowser": map[string]interface{}{
				"isEnabled": false,
				"mode":      "never",
			},
			"continuousAccessEvaluation": map[string]interface{}{
				"mode": "strict",
			},
			"secureSignInSession": map[string]interface{}{
				"isEnabled": true,
			},
		},
	}

	mockState.Lock()
	mockState.policies[minimalPolicy["id"].(string)] = minimalPolicy
	mockState.policies[maximalPolicy["id"].(string)] = maximalPolicy
	mockState.Unlock()
}

// ensurePolicyStructures ensures required nested structures exist in policy data
func ensurePolicyStructures(policyData map[string]interface{}) {
	// Ensure conditions structure
	if conditions, ok := policyData["conditions"].(map[string]interface{}); ok {
		// Ensure applications structure
		if applications, ok := conditions["applications"].(map[string]interface{}); ok {
			commonMocks.EnsureField(applications, "includeApplications", []string{})
			commonMocks.EnsureField(applications, "excludeApplications", []string{})
			commonMocks.EnsureField(applications, "includeUserActions", []string{})
		}

		// Ensure users structure
		if users, ok := conditions["users"].(map[string]interface{}); ok {
			commonMocks.EnsureField(users, "includeUsers", []string{})
			commonMocks.EnsureField(users, "excludeUsers", []string{})
			commonMocks.EnsureField(users, "includeGroups", []string{})
			commonMocks.EnsureField(users, "excludeGroups", []string{})
		}

		// Ensure optional structures exist if referenced
		if platforms, ok := conditions["platforms"].(map[string]interface{}); ok {
			commonMocks.EnsureField(platforms, "includePlatforms", []string{})
			commonMocks.EnsureField(platforms, "excludePlatforms", []string{})
		}

		if devices, ok := conditions["devices"].(map[string]interface{}); ok {
			commonMocks.EnsureField(devices, "includeDevices", []string{})
			commonMocks.EnsureField(devices, "excludeDevices", []string{})
			commonMocks.EnsureField(devices, "includeDeviceStates", []string{})
			commonMocks.EnsureField(devices, "excludeDeviceStates", []string{})
		}

		if locations, ok := conditions["locations"].(map[string]interface{}); ok {
			commonMocks.EnsureField(locations, "includeLocations", []string{})
			commonMocks.EnsureField(locations, "excludeLocations", []string{})
		}

		// Ensure array fields
		commonMocks.EnsureField(conditions, "clientAppTypes", []string{})
		commonMocks.EnsureField(conditions, "userRiskLevels", []string{})
		commonMocks.EnsureField(conditions, "signInRiskLevels", []string{})
		commonMocks.EnsureField(conditions, "servicePrincipalRiskLevels", []string{})
	}

	// Ensure grant controls structure
	if grantControls, ok := policyData["grantControls"].(map[string]interface{}); ok {
		commonMocks.EnsureField(grantControls, "builtInControls", []string{})
		commonMocks.EnsureField(grantControls, "customAuthenticationFactors", []string{})
		commonMocks.EnsureField(grantControls, "termsOfUse", []string{})
	}
}
