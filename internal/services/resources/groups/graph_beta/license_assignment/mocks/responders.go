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
	groupLicenses map[string]map[string]any // groupID -> license data including assignedLicenses array
}

func init() {
	mockState.groupLicenses = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	commonMocks.GlobalRegistry.Register("group_license_assignment", &GroupLicenseAssignmentMock{})
}

// GroupLicenseAssignmentMock provides mock responses for group license assignment operations
type GroupLicenseAssignmentMock struct{}

// Ensure GroupLicenseAssignmentMock implements MockRegistrar interface
var _ commonMocks.MockRegistrar = (*GroupLicenseAssignmentMock)(nil)

// loadJSONResponse loads a JSON response from the tests/responses directory
func (m *GroupLicenseAssignmentMock) loadJSONResponse(filePath string) (map[string]any, error) {
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

// getJSONFileForGroupID determines which JSON file to load based on the group ID
func getJSONFileForGroupID(groupID string, operation string) string {
	switch groupID {
	case "00000000-0000-0000-0000-000000000002":
		// Minimal config group
		switch operation {
		case constants.TfOperationCreate:
			return "validate_create/post_group_minimal_success.json"
		case constants.TfOperationUpdate:
			return "validate_update/patch_group_success.json"
		case constants.TfTfOperationDelete:
			return "validate_delete/get_group_not_found.json"
		}
	case "00000000-0000-0000-0000-000000000003":
		// Maximal config group
		switch operation {
		case constants.TfOperationCreate:
			return "validate_create/post_group_maximal_success.json"
		case constants.TfOperationUpdate:
			return "validate_update/patch_group_success.json"
		case constants.TfTfOperationDelete:
			return "validate_delete/get_group_not_found.json"
		}
	case "invalid-group-id":
		return "validate_create/post_group_error.json"
	}
	return ""
}

// RegisterMocks registers HTTP mock responses for group license assignment operations
func (m *GroupLicenseAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.groupLicenses = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for group data
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupID := urlParts[len(urlParts)-1]

			mockState.Lock()
			groupData, exists := mockState.groupLicenses[groupID]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, groupData)
		})

	// Register POST for license assignment (create/update/delete)
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/[^/]+/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupID := urlParts[len(urlParts)-2]

			// Check for invalid group ID format
			if groupID == "invalid-group-id" {
				jsonFile := getJSONFileForGroupID(groupID, constants.TfOperationCreate)
				if jsonFile != "" {
					errorResp, err := m.loadJSONResponse(jsonFile)
					if err == nil {
						return httpmock.NewJsonResponse(400, errorResp)
					}
				}
				return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Invalid group ID format"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Get or create group data
			groupData, exists := mockState.groupLicenses[groupID]
			if !exists {
				// Load initial group data
				jsonFile := getJSONFileForGroupID(groupID, constants.TfOperationCreate)
				if jsonFile == "" {
					return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
				}

				loadedData, err := m.loadJSONResponse(jsonFile)
				if err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"Failed to load mock response: %s"}}`, err.Error())), nil
				}
				groupData = loadedData
				mockState.groupLicenses[groupID] = groupData
			}

			// Get current licenses
			currentLicenses, ok := groupData["assignedLicenses"].([]any)
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

			groupData["assignedLicenses"] = currentLicenses
			return httpmock.NewJsonResponse(200, groupData)
		})
}

// RegisterErrorMocks registers error mock responses for testing error scenarios
func (m *GroupLicenseAssignmentMock) RegisterErrorMocks() {
	// Register all standard mocks first
	m.RegisterMocks()

	// Override with error responses for specific scenarios
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/invalid-group-id/assignLicense$`,
		func(req *http.Request) (*http.Response, error) {
			errorResp, err := m.loadJSONResponse("validate_create/post_group_error.json")
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Invalid group ID format"}}`), nil
			}
			return httpmock.NewJsonResponse(400, errorResp)
		})
}

// CleanupMockState cleans up the mock state after tests
func (m *GroupLicenseAssignmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.groupLicenses = make(map[string]map[string]any)
	mockState.Unlock()
}
