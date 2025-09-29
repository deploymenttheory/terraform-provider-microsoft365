package mocks

import (
	"encoding/json"
	"fmt"
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
	userLicenses map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.userLicenses = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserLicenseAssignmentMock provides mock responses for user license assignment operations
type UserLicenseAssignmentMock struct{}

// RegisterMocks registers HTTP mock responses for user license assignment operations
func (m *UserLicenseAssignmentMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.userLicenses = make(map[string]map[string]any)
	mockState.Unlock()

	// Initialize base user data
	baseUserId := "00000000-0000-0000-0000-000000000001"
	baseUserData := map[string]any{
		"id":                baseUserId,
		"userPrincipalName": "test.user@contoso.com",
		"assignedLicenses": []map[string]any{
			{
				"skuId": "11111111-1111-1111-1111-111111111111",
				"disabledPlans": []string{
					"22222222-2222-2222-2222-222222222222",
				},
			},
		},
	}

	mockState.Lock()
	mockState.userLicenses[baseUserId] = baseUserData
	mockState.Unlock()

	// Register GET for user data
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.userLicenses[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			if userData["assignedLicenses"] == nil {
				userData["assignedLicenses"] = []map[string]any{}
			}

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register GET for license details
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+/licenseDetails$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-2]

			mockState.Lock()
			userData, exists := mockState.userLicenses[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			// Extract assigned licenses from user data
			var assignedLicenses []map[string]any
			if userData["assignedLicenses"] != nil {
				assignedLicenses = userData["assignedLicenses"].([]map[string]any)
			} else {
				assignedLicenses = []map[string]any{}
			}

			// Convert to license details format
			licenseDetails := make([]map[string]any, 0, len(assignedLicenses))
			for _, license := range assignedLicenses {
				skuId, ok := license["skuId"].(string)
				if !ok {
					continue
				}

				// Create license detail with service plans
				licenseDetail := map[string]any{
					"id":            uuid.New().String(),
					"skuId":         skuId,
					"skuPartNumber": fmt.Sprintf("SKU_PART_%s", skuId[0:8]),
					"servicePlans": []map[string]any{
						{
							"servicePlanId":      "33333333-3333-3333-3333-333333333333",
							"servicePlanName":    "ServicePlan1",
							"provisioningStatus": "Success",
							"appliesTo":          "User",
						},
					},
				}

				licenseDetails = append(licenseDetails, licenseDetail)
			}

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#users('%s')/licenseDetails", userId),
				"value":          licenseDetails,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for license assignment
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/users/[^/]+/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-2]

			mockState.Lock()
			userData, exists := mockState.userLicenses[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Process add licenses
			if addLicenses, ok := requestBody["addLicenses"].([]interface{}); ok {
				currentLicenses := make([]map[string]any, 0)

				// Get existing licenses
				if userData["assignedLicenses"] != nil {
					currentLicenses = userData["assignedLicenses"].([]map[string]any)
				}

				// Add new licenses
				for _, addLicense := range addLicenses {
					if licenseObj, ok := addLicense.(map[string]any); ok {
						skuId, hasSkuId := licenseObj["skuId"].(string)
						if !hasSkuId {
							continue
						}

						// Check if this license already exists
						exists := false
						for i, existing := range currentLicenses {
							if existing["skuId"] == skuId {
								// Update existing license
								if disabledPlans, ok := licenseObj["disabledPlans"].([]interface{}); ok {
									disabledPlanStrings := make([]string, 0, len(disabledPlans))
									for _, plan := range disabledPlans {
										if planStr, ok := plan.(string); ok {
											disabledPlanStrings = append(disabledPlanStrings, planStr)
										}
									}
									currentLicenses[i]["disabledPlans"] = disabledPlanStrings
								}
								exists = true
								break
							}
						}

						// Add new license if it doesn't exist
						if !exists {
							newLicense := map[string]any{
								"skuId": skuId,
							}

							if disabledPlans, ok := licenseObj["disabledPlans"].([]interface{}); ok {
								disabledPlanStrings := make([]string, 0, len(disabledPlans))
								for _, plan := range disabledPlans {
									if planStr, ok := plan.(string); ok {
										disabledPlanStrings = append(disabledPlanStrings, planStr)
									}
								}
								newLicense["disabledPlans"] = disabledPlanStrings
							} else {
								newLicense["disabledPlans"] = []string{}
							}

							currentLicenses = append(currentLicenses, newLicense)
						}
					}
				}

				// Process remove licenses
				if removeLicenses, ok := requestBody["removeLicenses"].([]interface{}); ok && len(removeLicenses) > 0 {
					removeLicenseMap := make(map[string]bool)
					for _, licenseId := range removeLicenses {
						if licenseStr, ok := licenseId.(string); ok {
							removeLicenseMap[licenseStr] = true
						}
					}

					// Filter out removed licenses
					filteredLicenses := make([]map[string]any, 0)
					for _, license := range currentLicenses {
						skuId, ok := license["skuId"].(string)
						if !ok || removeLicenseMap[skuId] {
							continue
						}
						filteredLicenses = append(filteredLicenses, license)
					}
					currentLicenses = filteredLicenses
				}

				// Update user data
				mockState.Lock()
				userData["assignedLicenses"] = currentLicenses
				mockState.userLicenses[userId] = userData
				mockState.Unlock()
			}

			// Return a proper response that includes the updated user data
			mockState.Lock()
			updatedUserData := mockState.userLicenses[userId]
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, updatedUserData)
		})

	// Register specific user IDs for testing
	registerSpecificUserMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserLicenseAssignmentMock) RegisterErrorMocks() {
	// Register error response for license assignment
	errorUserId := "99999999-9999-9999-9999-999999999999"
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users/"+errorUserId+"/assignLicense",
		factories.ErrorResponse(400, "BadRequest", "Error assigning license"))

	// Register GET for error user to ensure it exists but will fail on license assignment
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/"+errorUserId,
		func(req *http.Request) (*http.Response, error) {
			userData := map[string]any{
				"id":                errorUserId,
				"userPrincipalName": "error.user@contoso.com",
				"assignedLicenses":  []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, userData)
		})

	// Register error response for user not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/not-found-user",
		factories.ErrorResponse(404, "ResourceNotFound", "User not found"))
}

// registerSpecificUserMocks registers mocks for specific test scenarios
func registerSpecificUserMocks() {
	// Minimal user with no licenses
	minimalUserId := "00000000-0000-0000-0000-000000000002"
	minimalUserData := map[string]any{
		"id":                minimalUserId,
		"userPrincipalName": "minimal.user@contoso.com",
		"assignedLicenses":  []map[string]any{},
	}

	mockState.Lock()
	mockState.userLicenses[minimalUserId] = minimalUserData
	mockState.Unlock()

	// Maximal user with multiple licenses
	maximalUserId := "00000000-0000-0000-0000-000000000003"
	maximalUserData := map[string]any{
		"id":                maximalUserId,
		"userPrincipalName": "maximal.user@contoso.com",
		"assignedLicenses": []map[string]any{
			{
				"skuId": "44444444-4444-4444-4444-444444444444",
				"disabledPlans": []string{
					"55555555-5555-5555-5555-555555555555",
					"66666666-6666-6666-6666-666666666666",
				},
			},
			{
				"skuId":         "77777777-7777-7777-7777-777777777777",
				"disabledPlans": []string{},
			},
		},
	}

	mockState.Lock()
	mockState.userLicenses[maximalUserId] = maximalUserData
	mockState.Unlock()

	// Register specific GET for these users
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/"+minimalUserId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			userData := mockState.userLicenses[minimalUserId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, userData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/"+maximalUserId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			userData := mockState.userLicenses[maximalUserId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, userData)
		})
}
