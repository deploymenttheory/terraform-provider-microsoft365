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
	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	userLicenses map[string]map[string]any // userID -> license data including assignedLicenses array
}

func init() {
	mockState.userLicenses = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	commonMocks.GlobalRegistry.Register("user_license_assignment", &UserLicenseAssignmentMock{})
}

// UserLicenseAssignmentMock provides mock responses for user license assignment operations
type UserLicenseAssignmentMock struct{}

// Ensure UserLicenseAssignmentMock implements MockRegistrar interface
var _ commonMocks.MockRegistrar = (*UserLicenseAssignmentMock)(nil)

// loadJSONResponse loads a JSON response from the tests/responses directory
func (m *UserLicenseAssignmentMock) loadJSONResponse(filePath string) (map[string]any, error) {
	fullPath := filepath.Join("..", "tests", "responses", filePath)
	jsonContent, err := helpers.ParseJSONFile(fullPath)
	if err != nil {
		return nil, err
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, err
	}
	return response, nil
}

// getJSONFileForUserID determines which JSON file to load based on the user ID
func getJSONFileForUserID(userID string, operation string) string {
	switch userID {
	case "00000000-0000-0000-0000-000000000002":
		// Minimal config user
		switch operation {
		case constants.TfOperationCreate:
			return "validate_create/post_user_minimal_success.json"
		case constants.TfOperationUpdate:
			return "validate_update/patch_user_success.json"
		case constants.TfTfOperationDelete:
			return "validate_delete/get_user_not_found.json"
		}
	case "00000000-0000-0000-0000-000000000003":
		// Maximal config user
		switch operation {
		case constants.TfOperationCreate:
			return "validate_create/post_user_maximal_success.json"
		case constants.TfOperationUpdate:
			return "validate_update/patch_user_success.json"
		case constants.TfTfOperationDelete:
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

				// Create mock service plans
				servicePlans := []map[string]any{
					{
						"servicePlanId":      fmt.Sprintf("00000000-0000-0000-0000-%012s", skuID[24:36]),
						"servicePlanName":    "Exchange Online",
						"provisioningStatus": "Success",
						"appliesTo":          "User",
					},
				}

				licenseDetail := map[string]any{
					"id":            userID + "_" + skuID,
					"skuId":         skuID,
					"skuPartNumber": fmt.Sprintf("SKU_%s", skuID[0:8]),
					"servicePlans":  servicePlans,
				}
				licenseDetails = append(licenseDetails, licenseDetail)
			}

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#users('%s')/licenseDetails", userID),
				"value":          licenseDetails,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for license assignment (create/update/delete)
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/users/[^/]+/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userID := urlParts[len(urlParts)-2]

			// Check for invalid user ID format
			if userID == "invalid-user-id" {
				jsonFile := getJSONFileForUserID(userID, constants.TfOperationCreate)
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

			mockState.Lock()
			defer mockState.Unlock()

			// Get or create user data
			userData, exists := mockState.userLicenses[userID]
			if !exists {
				// Load initial user data
				jsonFile := getJSONFileForUserID(userID, constants.TfOperationCreate)
				if jsonFile == "" {
					return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
				}

				loadedData, err := m.loadJSONResponse(jsonFile)
				if err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to load mock response: %s"}}`, err.Error())), nil
				}
				userData = loadedData
				mockState.userLicenses[userID] = userData
			}

			// Get current licenses
			currentLicenses, ok := userData["assignedLicenses"].([]any)
			if !ok {
				currentLicenses = []any{}
			}

			// Process add licenses
			addLicenses, hasAdd := requestBody["addLicenses"].([]any)
			if hasAdd && len(addLicenses) > 0 {
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
							} else {
								existingMap["disabledPlans"] = []any{}
							}
							currentLicenses[i] = existingMap
							found = true
							break
						}
					}

					if !found {
						normalizedLicense := map[string]any{
							"skuId": skuID,
						}
						if disabledPlans, ok := licenseObj["disabledPlans"]; ok {
							normalizedLicense["disabledPlans"] = disabledPlans
						} else {
							normalizedLicense["disabledPlans"] = []any{}
						}
						currentLicenses = append(currentLicenses, normalizedLicense)
					}
				}
			}

			// Process remove licenses
			removeLicenses, hasRemove := requestBody["removeLicenses"].([]any)
			if hasRemove && len(removeLicenses) > 0 {
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
				currentLicenses = filteredLicenses
			}

			userData["assignedLicenses"] = currentLicenses
			return httpmock.NewJsonResponse(200, userData)
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
