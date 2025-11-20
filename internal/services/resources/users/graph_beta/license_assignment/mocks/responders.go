package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	userLicenses map[string]map[string]any
}

func init() {
	mockState.userLicenses = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserLicenseAssignmentMock provides mock responses for user license assignment operations
type UserLicenseAssignmentMock struct{}

// loadJSONResponse loads a JSON response from the tests/responses directory
func (m *UserLicenseAssignmentMock) loadJSONResponse(filePath string) (map[string]any, error) {
	fullPath := filepath.Join("tests", "responses", filePath)
	return mocks.LoadJSONResponse(fullPath)
}

// getJSONFileForUserID determines which JSON file to load based on the user ID
func getJSONFileForUserID(userID string, operation string) string {
	switch userID {
	case "00000000-0000-0000-0000-000000000002":
		// Minimal config user
		switch operation {
		case "create":
			return "validate_create/post_user_minimal_success.json"
		case "update":
			return "validate_update/patch_user_success.json"
		case "delete":
			return "validate_delete/get_user_not_found.json"
		}
	case "00000000-0000-0000-0000-000000000003":
		// Maximal config user
		switch operation {
		case "create":
			return "validate_create/post_user_maximal_success.json"
		case "update":
			return "validate_update/patch_user_success.json"
		case "delete":
			return "validate_delete/get_user_not_found.json"
		}
	case "invalid-user-id":
		return "validate_create/post_user_error.json"
	}
	return ""
}

// RegisterMocks registers HTTP mock responses for user license assignment operations
func (m *UserLicenseAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.userLicenses = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for user data
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userID := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.userLicenses[userID]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register GET for license details
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+/licenseDetails$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userID := urlParts[len(urlParts)-2]

			mockState.Lock()
			userData, exists := mockState.userLicenses[userID]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			// Convert assigned licenses to license details format
			var assignedLicenses []any
			if userData["assignedLicenses"] != nil {
				assignedLicenses = userData["assignedLicenses"].([]any)
			}

			licenseDetails := make([]map[string]any, 0, len(assignedLicenses))
			for _, license := range assignedLicenses {
				licenseMap := license.(map[string]any)
				skuID, ok := licenseMap["skuId"].(string)
				if !ok {
					continue
				}

				licenseDetail := map[string]any{
					"id":            userID + "_" + skuID,
					"skuId":         skuID,
					"skuPartNumber": fmt.Sprintf("SKU_%s", skuID[0:8]),
				}
				licenseDetails = append(licenseDetails, licenseDetail)
			}

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#users('%s')/licenseDetails", userID),
				"value":          licenseDetails,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for license assignment (create/update)
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/users/[^/]+/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userID := urlParts[len(urlParts)-2]

			// Check for invalid user ID format
			if userID == "invalid-user-id" {
				jsonFile := getJSONFileForUserID(userID, "create")
				if jsonFile != "" {
					errorResp, err := m.loadJSONResponse(jsonFile)
					if err == nil {
						return httpmock.NewJsonResponse(400, errorResp)
					}
				}
				return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Invalid user ID format"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Determine if this is a create or update operation
			mockState.Lock()
			_, exists := mockState.userLicenses[userID]
			mockState.Unlock()

			var operation string
			if exists {
				operation = "update"
			} else {
				operation = "create"
			}

			// Load the appropriate JSON response
			jsonFile := getJSONFileForUserID(userID, operation)
			if jsonFile == "" {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			userData, err := m.loadJSONResponse(jsonFile)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to load mock response: %s"}}`, err.Error())), nil
			}

			// Process the request to merge add/remove licenses
			if addLicenses, ok := requestBody["addLicenses"].([]any); ok && len(addLicenses) > 0 {
				// For create, just use the JSON file's licenses
				// For update, merge with existing licenses
				if operation == "update" {
					currentLicenses := userData["assignedLicenses"].([]any)

					// Add new licenses from request
					for _, addLicense := range addLicenses {
						licenseObj := addLicense.(map[string]any)
						skuID := licenseObj["skuId"].(string)

						// Check if license already exists
						found := false
						for i, existing := range currentLicenses {
							existingMap := existing.(map[string]any)
							if existingMap["skuId"] == skuID {
								// Update disabled plans
								if disabledPlans, ok := licenseObj["disabledPlans"]; ok {
									existingMap["disabledPlans"] = disabledPlans
									currentLicenses[i] = existingMap
								}
								found = true
								break
							}
						}

						if !found {
							currentLicenses = append(currentLicenses, licenseObj)
						}
					}
					userData["assignedLicenses"] = currentLicenses
				}
			}

			// Process remove licenses
			if removeLicenses, ok := requestBody["removeLicenses"].([]any); ok && len(removeLicenses) > 0 {
				currentLicenses := userData["assignedLicenses"].([]any)
				filteredLicenses := make([]any, 0)

				for _, existing := range currentLicenses {
					existingMap := existing.(map[string]any)
					skuID := existingMap["skuId"].(string)

					// Check if this license should be removed
					shouldRemove := false
					for _, removeLicense := range removeLicenses {
						if removeLicense.(string) == skuID {
							shouldRemove = true
							break
						}
					}

					if !shouldRemove {
						filteredLicenses = append(filteredLicenses, existing)
					}
				}
				userData["assignedLicenses"] = filteredLicenses
			}

			// Store the updated state
			mockState.Lock()
			mockState.userLicenses[userID] = userData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register DELETE for license removal (resource deletion)
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/users/[^/]+/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userID := urlParts[len(urlParts)-2]

			// Parse request body to check if all licenses are being removed
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Check if this is a delete operation (removing all licenses)
			addLicenses, hasAdd := requestBody["addLicenses"].([]any)
			removeLicenses, hasRemove := requestBody["removeLicenses"].([]any)

			if hasRemove && len(removeLicenses) > 0 && (!hasAdd || len(addLicenses) == 0) {
				// This is a delete operation - remove user from state
				mockState.Lock()
				delete(mockState.userLicenses, userID)
				mockState.Unlock()

				// Return empty licenses
				userData := map[string]any{
					"id":                userID,
					"userPrincipalName": "test.user@contoso.com",
					"assignedLicenses":  []any{},
				}
				return httpmock.NewJsonResponse(200, userData)
			}

			// Otherwise, handle as normal update
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`), nil
		})
}

// RegisterErrorMocks registers error mock responses for testing error scenarios
func (m *UserLicenseAssignmentMock) RegisterErrorMocks() {
	// Register all standard mocks first
	m.RegisterMocks()

	// Override with error responses for specific scenarios
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/users/invalid-user-id/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			errorResp, err := m.loadJSONResponse("validate_create/post_user_error.json")
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Invalid user ID format"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResp)
		})
}

// CleanupMockState cleans up the mock state after tests
func (m *UserLicenseAssignmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.userLicenses = make(map[string]map[string]any)
	mockState.Unlock()
}
