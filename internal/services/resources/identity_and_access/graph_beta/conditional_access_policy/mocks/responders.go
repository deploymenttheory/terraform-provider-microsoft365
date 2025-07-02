package mocks

import (
	"encoding/json"
	"io"
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

			// Generate ID if not provided - use predefined IDs for test policies
			if policyData["id"] == nil {
				displayName, _ := policyData["displayName"].(string)

				// Check if this matches a predefined test policy and use its ID
				switch displayName {
				case "Block Legacy Authentication - Minimal":
					policyData["id"] = "minimal-policy-id-12345"
				case "Comprehensive Security Policy - Maximal":
					policyData["id"] = "maximal-policy-id-67890"
				case "Comprehensive Security Policy - Updated from Minimal":
					// This is used in the update test
					policyData["id"] = "minimal-policy-id-12345"
				default:
					// Generate random ID for other policies
					policyData["id"] = uuid.New().String()
				}
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

			// Parse the request body
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Failed to read request body"}}`), nil
			}

			var updateData map[string]interface{}
			if err := json.Unmarshal(body, &updateData); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON in request body"}}`), nil
			}

			// Update the policy with the new data, preserving structure
			mockState.Lock()
			defer mockState.Unlock()

			// Update top-level fields
			if displayName, ok := updateData["displayName"].(string); ok {
				policyData["displayName"] = displayName
			}
			if state, ok := updateData["state"].(string); ok {
				policyData["state"] = state
			}
			if templateId, ok := updateData["templateId"].(string); ok {
				policyData["templateId"] = templateId
			}
			if partialEnablementStrategy, ok := updateData["partialEnablementStrategy"].(string); ok {
				policyData["partialEnablementStrategy"] = partialEnablementStrategy
			}

			// Update conditions
			if updateConditions, ok := updateData["conditions"].(map[string]interface{}); ok {
				currentConditions, _ := policyData["conditions"].(map[string]interface{})
				if currentConditions == nil {
					currentConditions = make(map[string]interface{})
					policyData["conditions"] = currentConditions
				}

				// Update client app types
				if clientAppTypes, ok := updateConditions["clientAppTypes"].([]interface{}); ok {
					currentConditions["clientAppTypes"] = clientAppTypes
				}

				// Update risk levels
				if userRiskLevels, ok := updateConditions["userRiskLevels"].([]interface{}); ok {
					currentConditions["userRiskLevels"] = userRiskLevels
				}
				if signInRiskLevels, ok := updateConditions["signInRiskLevels"].([]interface{}); ok {
					currentConditions["signInRiskLevels"] = signInRiskLevels
				}
				if servicePrincipalRiskLevels, ok := updateConditions["servicePrincipalRiskLevels"].([]interface{}); ok {
					currentConditions["servicePrincipalRiskLevels"] = servicePrincipalRiskLevels
				}

				// Update applications
				if updateApplications, ok := updateConditions["applications"].(map[string]interface{}); ok {
					currentApplications, _ := currentConditions["applications"].(map[string]interface{})
					if currentApplications == nil {
						currentApplications = make(map[string]interface{})
						currentConditions["applications"] = currentApplications
					}

					// Update application fields
					if includeApps, ok := updateApplications["includeApplications"].([]interface{}); ok {
						currentApplications["includeApplications"] = includeApps
					}
					if excludeApps, ok := updateApplications["excludeApplications"].([]interface{}); ok {
						currentApplications["excludeApplications"] = excludeApps
					}
					if includeUserActions, ok := updateApplications["includeUserActions"].([]interface{}); ok {
						currentApplications["includeUserActions"] = includeUserActions
					}
					if includeAuthContext, ok := updateApplications["includeAuthenticationContextClassReferences"].([]interface{}); ok {
						currentApplications["includeAuthenticationContextClassReferences"] = includeAuthContext
					}
					if appFilter, ok := updateApplications["applicationFilter"].(map[string]interface{}); ok {
						currentApplications["applicationFilter"] = appFilter
					} else if updateApplications["applicationFilter"] == nil {
						// Remove application filter if explicitly set to null
						delete(currentApplications, "applicationFilter")
					}
				}

				// Update users
				if updateUsers, ok := updateConditions["users"].(map[string]interface{}); ok {
					currentUsers, _ := currentConditions["users"].(map[string]interface{})
					if currentUsers == nil {
						currentUsers = make(map[string]interface{})
						currentConditions["users"] = currentUsers
					}

					// Update user fields
					if includeUsers, ok := updateUsers["includeUsers"].([]interface{}); ok {
						currentUsers["includeUsers"] = includeUsers
					}
					if excludeUsers, ok := updateUsers["excludeUsers"].([]interface{}); ok {
						currentUsers["excludeUsers"] = excludeUsers
					}
					if includeGroups, ok := updateUsers["includeGroups"].([]interface{}); ok {
						currentUsers["includeGroups"] = includeGroups
					}
					if excludeGroups, ok := updateUsers["excludeGroups"].([]interface{}); ok {
						currentUsers["excludeGroups"] = excludeGroups
					}
					if includeRoles, ok := updateUsers["includeRoles"].([]interface{}); ok {
						currentUsers["includeRoles"] = includeRoles
					}
					if excludeRoles, ok := updateUsers["excludeRoles"].([]interface{}); ok {
						currentUsers["excludeRoles"] = excludeRoles
					}
					if includeGuests, ok := updateUsers["includeGuestsOrExternalUsers"].(map[string]interface{}); ok {
						currentUsers["includeGuestsOrExternalUsers"] = includeGuests
					} else if updateUsers["includeGuestsOrExternalUsers"] == nil {
						delete(currentUsers, "includeGuestsOrExternalUsers")
					}
					if excludeGuests, ok := updateUsers["excludeGuestsOrExternalUsers"].(map[string]interface{}); ok {
						currentUsers["excludeGuestsOrExternalUsers"] = excludeGuests
					} else if updateUsers["excludeGuestsOrExternalUsers"] == nil {
						delete(currentUsers, "excludeGuestsOrExternalUsers")
					}
				}

				// Update platforms
				if updatePlatforms, ok := updateConditions["platforms"].(map[string]interface{}); ok {
					currentPlatforms, _ := currentConditions["platforms"].(map[string]interface{})
					if currentPlatforms == nil {
						currentPlatforms = make(map[string]interface{})
						currentConditions["platforms"] = currentPlatforms
					}

					if includePlatforms, ok := updatePlatforms["includePlatforms"].([]interface{}); ok {
						currentPlatforms["includePlatforms"] = includePlatforms
					}
					if excludePlatforms, ok := updatePlatforms["excludePlatforms"].([]interface{}); ok {
						currentPlatforms["excludePlatforms"] = excludePlatforms
					}
				}

				// Update locations
				if updateLocations, ok := updateConditions["locations"].(map[string]interface{}); ok {
					currentLocations, _ := currentConditions["locations"].(map[string]interface{})
					if currentLocations == nil {
						currentLocations = make(map[string]interface{})
						currentConditions["locations"] = currentLocations
					}

					if includeLocations, ok := updateLocations["includeLocations"].([]interface{}); ok {
						currentLocations["includeLocations"] = includeLocations
					}
					if excludeLocations, ok := updateLocations["excludeLocations"].([]interface{}); ok {
						currentLocations["excludeLocations"] = excludeLocations
					}
				}

				// Update devices
				if updateDevices, ok := updateConditions["devices"].(map[string]interface{}); ok {
					currentDevices, _ := currentConditions["devices"].(map[string]interface{})
					if currentDevices == nil {
						currentDevices = make(map[string]interface{})
						currentConditions["devices"] = currentDevices
					}

					if includeDevices, ok := updateDevices["includeDevices"].([]interface{}); ok {
						currentDevices["includeDevices"] = includeDevices
					}
					if excludeDevices, ok := updateDevices["excludeDevices"].([]interface{}); ok {
						currentDevices["excludeDevices"] = excludeDevices
					}
					if includeStates, ok := updateDevices["includeDeviceStates"].([]interface{}); ok {
						currentDevices["includeDeviceStates"] = includeStates
					}
					if excludeStates, ok := updateDevices["excludeDeviceStates"].([]interface{}); ok {
						currentDevices["excludeDeviceStates"] = excludeStates
					}
					if deviceFilter, ok := updateDevices["deviceFilter"].(map[string]interface{}); ok {
						currentDevices["deviceFilter"] = deviceFilter
					} else if updateDevices["deviceFilter"] == nil {
						delete(currentDevices, "deviceFilter")
					}
				}
			}

			// Update grant controls
			if updateGrantControls, ok := updateData["grantControls"].(map[string]interface{}); ok {
				currentGrantControls, _ := policyData["grantControls"].(map[string]interface{})
				if currentGrantControls == nil {
					currentGrantControls = make(map[string]interface{})
					policyData["grantControls"] = currentGrantControls
				}

				if operator, ok := updateGrantControls["operator"].(string); ok {
					currentGrantControls["operator"] = operator
				}
				if builtInControls, ok := updateGrantControls["builtInControls"].([]interface{}); ok {
					currentGrantControls["builtInControls"] = builtInControls
				}
				if customFactors, ok := updateGrantControls["customAuthenticationFactors"].([]interface{}); ok {
					currentGrantControls["customAuthenticationFactors"] = customFactors
				}
				if termsOfUse, ok := updateGrantControls["termsOfUse"].([]interface{}); ok {
					currentGrantControls["termsOfUse"] = termsOfUse
				}
				if authStrength, ok := updateGrantControls["authenticationStrength"].(map[string]interface{}); ok {
					currentGrantControls["authenticationStrength"] = authStrength
				} else if updateGrantControls["authenticationStrength"] == nil {
					delete(currentGrantControls, "authenticationStrength")
				}
			}

			// Update session controls
			if updateSessionControls, ok := updateData["sessionControls"].(map[string]interface{}); ok {
				currentSessionControls, _ := policyData["sessionControls"].(map[string]interface{})
				if currentSessionControls == nil {
					currentSessionControls = make(map[string]interface{})
					policyData["sessionControls"] = currentSessionControls
				}

				if disableResilience, ok := updateSessionControls["disableResilienceDefaults"].(bool); ok {
					currentSessionControls["disableResilienceDefaults"] = disableResilience
				}

				// Update application enforced restrictions
				if appRestrictions, ok := updateSessionControls["applicationEnforcedRestrictions"].(map[string]interface{}); ok {
					currentSessionControls["applicationEnforcedRestrictions"] = appRestrictions
				} else if updateSessionControls["applicationEnforcedRestrictions"] == nil {
					delete(currentSessionControls, "applicationEnforcedRestrictions")
				}

				// Update cloud app security
				if cloudAppSecurity, ok := updateSessionControls["cloudAppSecurity"].(map[string]interface{}); ok {
					currentSessionControls["cloudAppSecurity"] = cloudAppSecurity
				} else if updateSessionControls["cloudAppSecurity"] == nil {
					delete(currentSessionControls, "cloudAppSecurity")
				}

				// Update sign in frequency
				if signInFrequency, ok := updateSessionControls["signInFrequency"].(map[string]interface{}); ok {
					currentSessionControls["signInFrequency"] = signInFrequency
				} else if updateSessionControls["signInFrequency"] == nil {
					delete(currentSessionControls, "signInFrequency")
				}

				// Update persistent browser
				if persistentBrowser, ok := updateSessionControls["persistentBrowser"].(map[string]interface{}); ok {
					currentSessionControls["persistentBrowser"] = persistentBrowser
				} else if updateSessionControls["persistentBrowser"] == nil {
					delete(currentSessionControls, "persistentBrowser")
				}

				// Update continuous access evaluation
				if cae, ok := updateSessionControls["continuousAccessEvaluation"].(map[string]interface{}); ok {
					currentSessionControls["continuousAccessEvaluation"] = cae
				} else if updateSessionControls["continuousAccessEvaluation"] == nil {
					delete(currentSessionControls, "continuousAccessEvaluation")
				}

				// Update secure sign in session
				if secureSignIn, ok := updateSessionControls["secureSignInSession"].(map[string]interface{}); ok {
					currentSessionControls["secureSignInSession"] = secureSignIn
				} else if updateSessionControls["secureSignInSession"] == nil {
					delete(currentSessionControls, "secureSignInSession")
				}
			}

			// Update modified date time
			policyData["modifiedDateTime"] = time.Now().Format(time.RFC3339)

			// Ensure all required structures and fields exist
			ensurePolicyStructures(policyData)

			// Return the updated policy
			respBody, _ := json.Marshal(policyData)
			return httpmock.NewStringResponse(200, string(respBody)), nil
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
// registerTestPolicies registers predefined test policies
func registerTestPolicies() {
	// Register minimal test policy - only includes fields actually specified in HCL
	minimalPolicy := map[string]interface{}{
		"id":               "minimal-policy-id-12345",
		"displayName":      "Block Legacy Authentication - Minimal",
		"state":            "enabled",
		"createdDateTime":  "2024-01-01T00:00:00Z",
		"modifiedDateTime": "2024-01-01T00:00:00Z",
		"conditions": map[string]interface{}{
			"clientAppTypes":   []interface{}{"exchangeActiveSync", "other"},
			"signInRiskLevels": []interface{}{},
			"applications": map[string]interface{}{
				"includeApplications": []interface{}{"All"},
			},
			"users": map[string]interface{}{
				"includeUsers": []interface{}{"All"},
			},
		},
		"grantControls": map[string]interface{}{
			"operator":        "OR",
			"builtInControls": []interface{}{"block"},
		},
	}

	// Register maximal test policy - includes realistic values for all major features
	maximalPolicy := map[string]interface{}{
		"id":               "maximal-policy-id-67890",
		"displayName":      "Comprehensive Security Policy - Maximal",
		"state":            "enabled",
		"createdDateTime":  "2024-01-01T00:00:00Z",
		"modifiedDateTime": "2024-01-01T00:00:00Z",
		"conditions": map[string]interface{}{
			"clientAppTypes":             []interface{}{"all"},
			"userRiskLevels":             []interface{}{"high"},
			"signInRiskLevels":           []interface{}{"high", "medium"},
			"servicePrincipalRiskLevels": []interface{}{"high", "medium"},
			"applications": map[string]interface{}{
				"includeApplications":                         []interface{}{"All"},
				"excludeApplications":                         []interface{}{"00000002-0000-0ff1-ce00-000000000000"},
				"includeUserActions":                          []interface{}{"urn:user:registersecurityinfo"},
				"includeAuthenticationContextClassReferences": []interface{}{"c00000000-0000-0000-0000-000000000001"},
				"applicationFilter": map[string]interface{}{
					"mode": "exclude",
					"rule": "device.deviceOwnership -eq \"Company\"",
				},
			},
			"users": map[string]interface{}{
				"includeUsers":  []interface{}{"All"},
				"excludeUsers":  []interface{}{"GuestsOrExternalUsers"},
				"includeGroups": []interface{}{"a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
				"excludeGroups": []interface{}{"f1e2d3c4-b5a6-9870-fedc-ba0987654321"},
				"includeRoles":  []interface{}{"62e90394-69f5-4237-9190-012177145e10"},
				"excludeRoles":  []interface{}{"e3973bdf-4987-49ae-837a-ba8e231c7286"},
				"includeGuestsOrExternalUsers": map[string]interface{}{
					"guestOrExternalUserTypes": "internalGuest,b2bCollaborationGuest",
					"externalTenants": map[string]interface{}{
						"membershipKind": "enumerated",
						"members":        []interface{}{"12345678-1234-1234-1234-123456789012"},
					},
				},
			},
			"platforms": map[string]interface{}{
				"includePlatforms": []interface{}{"all"},
				"excludePlatforms": []interface{}{"iOS", "android"},
			},
			"locations": map[string]interface{}{
				"includeLocations": []interface{}{"All"},
				"excludeLocations": []interface{}{"AllTrusted", "11111111-1111-1111-1111-111111111111"},
			},
			"devices": map[string]interface{}{
				"includeDevices":      []interface{}{"All"},
				"excludeDevices":      []interface{}{"22222222-2222-2222-2222-222222222222"},
				"includeDeviceStates": []interface{}{"domainJoined"},
				"excludeDeviceStates": []interface{}{"compliant"},
				"deviceFilter": map[string]interface{}{
					"mode": "include",
					"rule": "device.isCompliant -eq True -and device.deviceOwnership -eq \"Company\"",
				},
			},
		},
		"grantControls": map[string]interface{}{
			"operator":                    "AND",
			"builtInControls":             []interface{}{"mfa", "compliantDevice", "domainJoinedDevice"},
			"customAuthenticationFactors": []interface{}{"33333333-3333-3333-3333-333333333333"},
			"termsOfUse":                  []interface{}{"44444444-4444-4444-4444-444444444444"},
			"authenticationStrength": map[string]interface{}{
				"id":                    "00000000-0000-0000-0000-000000000004",
				"displayName":           "Multifactor authentication",
				"description":           "Combinations of methods that satisfy strong authentication, such as a password + SMS",
				"policyType":            "builtIn",
				"requirementsSatisfied": "mfa",
				"allowedCombinations": []interface{}{
					"password,sms",
					"password,voice",
					"password,hardwareOath",
					"password,softwareOath",
					"password,microsoftAuthenticatorPush",
					"windowsHelloForBusiness",
					"fido2",
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
				"cloudAppSecurityType": "mcasConfigured",
			},
			"signInFrequency": map[string]interface{}{
				"isEnabled":          true,
				"type":               "hours",
				"value":              4,
				"authenticationType": "primaryAndSecondaryAuthentication",
				"frequencyInterval":  "timeBased",
			},
			"persistentBrowser": map[string]interface{}{
				"isEnabled": true,
				"mode":      "always",
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
		// Initialize empty arrays for risk levels if not present
		commonMocks.EnsureField(conditions, "userRiskLevels", []interface{}{})
		commonMocks.EnsureField(conditions, "signInRiskLevels", []interface{}{})
		commonMocks.EnsureField(conditions, "servicePrincipalRiskLevels", []interface{}{})
		commonMocks.EnsureField(conditions, "clientAppTypes", []interface{}{})

		// Ensure applications structure
		if applications, ok := conditions["applications"].(map[string]interface{}); ok {
			commonMocks.EnsureField(applications, "includeApplications", []interface{}{})
			commonMocks.EnsureField(applications, "excludeApplications", []interface{}{})
			commonMocks.EnsureField(applications, "includeUserActions", []interface{}{})
			commonMocks.EnsureField(applications, "includeAuthenticationContextClassReferences", []interface{}{})
		} else {
			conditions["applications"] = map[string]interface{}{
				"includeApplications":                         []interface{}{},
				"excludeApplications":                         []interface{}{},
				"includeUserActions":                          []interface{}{},
				"includeAuthenticationContextClassReferences": []interface{}{},
			}
		}

		// Ensure users structure
		if users, ok := conditions["users"].(map[string]interface{}); ok {
			commonMocks.EnsureField(users, "includeUsers", []interface{}{})
			commonMocks.EnsureField(users, "excludeUsers", []interface{}{})
			commonMocks.EnsureField(users, "includeGroups", []interface{}{})
			commonMocks.EnsureField(users, "excludeGroups", []interface{}{})
			commonMocks.EnsureField(users, "includeRoles", []interface{}{})
			commonMocks.EnsureField(users, "excludeRoles", []interface{}{})
		} else {
			conditions["users"] = map[string]interface{}{
				"includeUsers":  []interface{}{},
				"excludeUsers":  []interface{}{},
				"includeGroups": []interface{}{},
				"excludeGroups": []interface{}{},
				"includeRoles":  []interface{}{},
				"excludeRoles":  []interface{}{},
			}
		}

		// Only ensure platforms structure if it's already present in the request
		if platforms, ok := conditions["platforms"].(map[string]interface{}); ok {
			commonMocks.EnsureField(platforms, "includePlatforms", []interface{}{})
			commonMocks.EnsureField(platforms, "excludePlatforms", []interface{}{})
		}

		// Only ensure locations structure if it's already present in the request
		if locations, ok := conditions["locations"].(map[string]interface{}); ok {
			commonMocks.EnsureField(locations, "includeLocations", []interface{}{})
			commonMocks.EnsureField(locations, "excludeLocations", []interface{}{})
		}

		// Only ensure devices structure if it's already present in the request
		if devices, ok := conditions["devices"].(map[string]interface{}); ok {
			commonMocks.EnsureField(devices, "includeDevices", []interface{}{})
			commonMocks.EnsureField(devices, "excludeDevices", []interface{}{})
			commonMocks.EnsureField(devices, "includeDeviceStates", []interface{}{})
			commonMocks.EnsureField(devices, "excludeDeviceStates", []interface{}{})
		}
	}

	// Ensure grant controls structure
	if grantControls, ok := policyData["grantControls"].(map[string]interface{}); ok {
		commonMocks.EnsureField(grantControls, "builtInControls", []interface{}{})
		commonMocks.EnsureField(grantControls, "customAuthenticationFactors", []interface{}{})
		commonMocks.EnsureField(grantControls, "termsOfUse", []interface{}{})

		// Ensure authentication strength if present
		if authStrength, ok := grantControls["authenticationStrength"].(map[string]interface{}); ok {
			commonMocks.EnsureField(authStrength, "allowedCombinations", []interface{}{})
		}
	}

	// Ensure session controls structure if present
	if _, ok := policyData["sessionControls"].(map[string]interface{}); ok {
		// No arrays to initialize here, but we could add structure checks if needed
	}
}
