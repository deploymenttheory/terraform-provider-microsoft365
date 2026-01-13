package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	users        map[string]map[string]any
	deletedItems map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.users = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserMock provides mock responses for user operations
type UserMock struct{}

// RegisterMocks registers HTTP mock responses for user operations
func (m *UserMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	// Register specific test users
	registerTestUsers()

	// Register GET for user by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register GET for user manager
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+/manager$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-2]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			// Check if user has a manager set
			managerId, hasManager := userData["managerId"].(string)
			if !hasManager || managerId == "" {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			// Return minimal manager info
			managerData := map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          managerId,
			}

			return httpmock.NewJsonResponse(200, managerData)
		})

	// Register GET for listing users
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			users := make([]map[string]any, 0, len(mockState.users))
			for _, user := range mockState.users {
				users = append(users, user)
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
				"value":          users,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating users
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			var userData map[string]any
			err := json.NewDecoder(req.Body).Decode(&userData)
			if err != nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			// Validate required fields
			displayName, hasDisplayName := userData["displayName"].(string)
			if !hasDisplayName {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}
			if _, ok := userData["userPrincipalName"].(string); !ok {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}
			if passwordProfile, ok := userData["passwordProfile"].(map[string]any); !ok || passwordProfile["password"] == nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			// Determine which fixture to use based on displayName
			var fixtureData map[string]any
			var fixtureErr error

			switch displayName {
			case "Minimal User":
				fixtureData, fixtureErr = loadFixture("validate_create/post_user_minimal_success.json")
			case "unit-test-user-maximal":
				fixtureData, fixtureErr = loadFixture("validate_create/post_user_maximal_success.json")
			case "unit-test-user-custom-sec-att":
				fixtureData, fixtureErr = loadFixture("validate_create/post_user_custom_sec_att_success.json")
			default:
				fixtureData = nil
				fixtureErr = nil
			}

			if fixtureErr == nil && fixtureData != nil {
				// Use fixture data but update dynamic fields
				userId := uuid.New().String()
				fixtureData["id"] = userId
				fixtureData["createdDateTime"] = time.Now().Format(time.RFC3339)

				// Handle manager@odata.bind if present in request
				if managerBindURL, ok := userData["manager@odata.bind"].(string); ok {
					parts := strings.Split(managerBindURL, "/")
					if len(parts) > 0 {
						fixtureData["managerId"] = parts[len(parts)-1]
					}
				}

				mockState.Lock()
				mockState.users[userId] = fixtureData
				mockState.Unlock()

				return httpmock.NewJsonResponse(201, fixtureData)
			}

			// Fallback to dynamic creation
			if userData["id"] == nil {
				userData["id"] = uuid.New().String()
			}

			now := time.Now().Format(time.RFC3339)
			userData["createdDateTime"] = now

			// Handle manager@odata.bind - extract ID and store as managerId
			if managerBindURL, ok := userData["manager@odata.bind"].(string); ok {
				parts := strings.Split(managerBindURL, "/")
				if len(parts) > 0 {
					managerId := parts[len(parts)-1]
					userData["managerId"] = managerId
				}
				delete(userData, "manager@odata.bind")
			}

			// Ensure collection fields are initialized
			commonMocks.EnsureField(userData, "businessPhones", []string{})
			commonMocks.EnsureField(userData, "otherMails", []string{})
			commonMocks.EnsureField(userData, "proxyAddresses", []string{})

			// Remove password from response (write-only field)
			if passwordProfile, ok := userData["passwordProfile"].(map[string]any); ok {
				delete(passwordProfile, "password")
			}

			// Store user in mock state
			userId := userData["id"].(string)
			mockState.Lock()
			mockState.users[userId] = userData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, userData)
		})

	// Register PATCH for updating users
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			var updateData map[string]any
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			// Determine fixture to use based on displayName in update
			displayName, hasDisplayName := updateData["displayName"].(string)
			if !hasDisplayName {
				displayName, _ = userData["displayName"].(string)
			}

			var fixtureData map[string]any
			var fixtureErr error

			switch displayName {
			case "Minimal User":
				fixtureData, fixtureErr = loadFixture("validate_update/patch_user_minimal_success.json")
			case "unit-test-user-maximal":
				fixtureData, fixtureErr = loadFixture("validate_update/patch_user_maximal_success.json")
			case "unit-test-user-custom-sec-att":
				fixtureData, fixtureErr = loadFixture("validate_update/patch_user_custom_sec_att_success.json")
			default:
				fixtureData = nil
				fixtureErr = nil
			}

			if fixtureErr == nil && fixtureData != nil {
				// Use fixture data but preserve the original ID
				fixtureData["id"] = userId

				// Handle manager@odata.bind if present in update
				if managerBindURL, ok := updateData["manager@odata.bind"].(string); ok {
					parts := strings.Split(managerBindURL, "/")
					if len(parts) > 0 {
						fixtureData["managerId"] = parts[len(parts)-1]
					}
				}

				mockState.Lock()
				mockState.users[userId] = fixtureData
				mockState.Unlock()

				return httpmock.NewJsonResponse(200, fixtureData)
			}

			// Fallback to dynamic update
			mockState.Lock()

			// Special handling for updates that remove fields
			isMinimalUpdate := false
			if _, hasDisplayName := updateData["displayName"]; hasDisplayName {
				if _, hasGivenName := updateData["givenName"]; !hasGivenName {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				fieldsToRemove := []string{
					"givenName", "surname", "jobTitle", "department", "companyName",
					"officeLocation", "city", "state", "country", "postalCode",
					"mobilePhone", "mail", "mailNickname", "usageLocation",
				}

				for _, field := range fieldsToRemove {
					delete(userData, field)
				}

				userData["businessPhones"] = []string{}
				userData["otherMails"] = []string{}
				userData["proxyAddresses"] = []string{}
			}

			// Apply the updates
			for k, v := range updateData {
				if k == "manager@odata.bind" {
					if managerBindURL, ok := v.(string); ok {
						parts := strings.Split(managerBindURL, "/")
						if len(parts) > 0 {
							managerId := parts[len(parts)-1]
							userData["managerId"] = managerId
						}
					}
				} else {
					userData[k] = v
				}
			}

			mockState.users[userId] = userData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register DELETE for removing users (soft delete)
	// Moves item to deletedItems collection instead of permanently deleting
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			if exists {
				// Move to deletedItems (soft delete behavior)
				mockState.deletedItems[userId] = userData
				delete(mockState.users, userId)
			}
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Get deleted item - GET /directory/deletedItems/{id}
	// Used for soft delete verification (polling until resource appears in deleted items)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			resourceId := urlParts[len(urlParts)-1]

			mockState.Lock()
			deletedItem, exists := mockState.deletedItems[resourceId]
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			// Load the deleted item response format from fixture
			fixtureData, err := loadFixture("validate_delete/get_deleted_item_success.json")
			if err == nil && fixtureData != nil {
				// Update with actual deleted item data
				fixtureData["id"] = resourceId
				if displayName, ok := deletedItem["displayName"]; ok {
					fixtureData["displayName"] = displayName
				}
				if upn, ok := deletedItem["userPrincipalName"]; ok {
					fixtureData["userPrincipalName"] = upn
				}
				return httpmock.NewJsonResponse(200, fixtureData)
			}

			return httpmock.NewJsonResponse(200, deletedItem)
		})

	// Permanent delete from deleted items - DELETE /directory/deletedItems/{id}
	// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			resourceId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.deletedItems[resourceId]
			if exists {
				delete(mockState.deletedItems, resourceId)
			}
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	// Register error response for user creation - always return error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		httpmock.NewStringResponder(400, errorBadRequest))

	// Register error response for user not found
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		httpmock.NewStringResponder(404, errorNotFound))

	// Register error response for PATCH
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		httpmock.NewStringResponder(400, errorBadRequest))

	// Register error response for DELETE
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		httpmock.NewStringResponder(400, errorBadRequest))

	// Register error response for GET deleted items
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[^/]+$`,
		httpmock.NewStringResponder(404, errorNotFound))

	// Register error response for DELETE deleted items
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[^/]+$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

// loadFixture loads a JSON fixture file from the tests/responses directory using the secure helpers package
func loadFixture(filename string) (map[string]any, error) {
	// Path relative to the mocks directory: ../tests/responses/
	fixturesPath := "../tests/responses/" + filename

	// Use the secure JSON parser from helpers package
	jsonContent, err := helpers.ParseJSONFile(fixturesPath)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// registerTestUsers registers predefined test users from JSON fixtures
func registerTestUsers() {
	minimalUserData, _ := loadFixture("validate_read/get_user_minimal_success.json")
	maximalUserData, _ := loadFixture("validate_read/get_user_maximal_success.json")
	customSecAttUserData, _ := loadFixture("validate_read/get_user_custom_sec_att_success.json")

	mockState.Lock()
	defer mockState.Unlock()

	if minimalUserData != nil {
		mockState.users[minimalUserData["id"].(string)] = minimalUserData
	}
	if maximalUserData != nil {
		mockState.users[maximalUserData["id"].(string)] = maximalUserData
	}
	if customSecAttUserData != nil {
		mockState.users[customSecAttUserData["id"].(string)] = customSecAttUserData
	}
}

// CleanupMockState clears the mock state
func (m *UserMock) CleanupMockState() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()
}
