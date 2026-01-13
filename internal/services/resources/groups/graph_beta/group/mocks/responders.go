package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	groups       map[string]map[string]any
	deletedItems map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.groups = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("group", &GroupMock{})
}

// GroupMock provides mock responses for group operations
type GroupMock struct{}

// Ensure GroupMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupMock)(nil)

// RegisterMocks registers HTTP mock responses for group operations
func (m *GroupMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing groups
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groups := make([]map[string]any, 0, len(mockState.groups))
			for _, group := range mockState.groups {
				groups = append(groups, group)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value":          groups,
			}

			respBody, _ := json.Marshal(response)
			return httpmock.NewStringResponse(200, string(respBody)), nil
		})

	// Register GET for specific group with special test IDs
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			// Handle special test IDs with external JSON files
			if strings.Contains(groupID, "error") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if group, exists := mockState.groups[groupID]; exists {
				respBody, _ := json.Marshal(group)
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		})

	// Register POST for creating groups
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			// Generate a new ID for the group
			newID := uuid.New().String()
			requestBody["id"] = newID

			// Add additional fields that would be set by the API
			requestBody["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#groups/$entity"
			requestBody["createdDateTime"] = "2023-01-01T00:00:00Z"

			// Handle special test cases based on display name
			if displayName, ok := requestBody["displayName"].(string); ok {
				var jsonFile string
				switch displayName {
				case "acc-security-group-with-assigned-membership-type":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_1.json"
				case "acc-security-group-with-dynamic-user-membership-type":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_2.json"
				case "acc-security-group-with-dynamic-device-membership-type":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_3.json"
				case "acc-security-group-with-entra-role-assignment":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_4.json"
				case "acc-m365-group-with-dynamic-user-membership-type":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_5.json"
				case "acc-m365-group-with-assigned-membership-type":
					jsonFile = "../tests/responses/validate_create/post_group_scenario_6.json"
				}

				if jsonFile != "" {
					jsonStr, err := helpers.ParseJSONFile(jsonFile)
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					// Parse the JSON and store it in mockState for subsequent GET requests
					var groupData map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &groupData); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					mockState.Lock()
					if id, ok := groupData["id"].(string); ok {
						mockState.groups[id] = groupData
					}
					mockState.Unlock()

					// Return with proper JSON content type
					resp := httpmock.NewStringResponse(201, jsonStr)
					resp.Header.Set("Content-Type", "application/json")
					return resp, nil
				}
			}

			mockState.Lock()
			mockState.groups[newID] = requestBody
			mockState.Unlock()

			respBody, _ := json.Marshal(requestBody)
			resp := httpmock.NewStringResponse(201, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register PATCH for updating groups
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if group, exists := mockState.groups[groupID]; exists {
				// Update the existing group
				for k, v := range requestBody {
					group[k] = v
				}
				return httpmock.NewStringResponse(204, ""), nil
			}

			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		})

	// Register DELETE for deleting groups (soft delete)
	// Moves item to deletedItems collection instead of permanently deleting
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			defer mockState.Unlock()

			if group, exists := mockState.groups[groupID]; exists {
				// Move to deletedItems (soft delete behavior)
				mockState.deletedItems[groupID] = group
				delete(mockState.groups, groupID)
				return httpmock.NewStringResponse(204, ""), nil
			}

			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		})

	// Get deleted item - GET /directory/deletedItems/{id}
	// Used for soft delete verification (polling until resource appears in deleted items)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			resourceID := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			defer mockState.Unlock()

			if deletedItem, exists := mockState.deletedItems[resourceID]; exists {
				respBody, _ := json.Marshal(deletedItem)
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		})

	// Permanent delete from deleted items - DELETE /directory/deletedItems/{id}
	// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			resourceID := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			defer mockState.Unlock()

			if _, exists := mockState.deletedItems[resourceID]; exists {
				delete(mockState.deletedItems, resourceID)
				return httpmock.NewStringResponse(204, ""), nil
			}

			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses that return errors for testing error handling
func (m *GroupMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
			return httpmock.NewStringResponse(400, jsonStr), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
			return httpmock.NewStringResponse(404, jsonStr), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		httpmock.NewStringResponder(400, errorBadRequest))

	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		httpmock.NewStringResponder(400, errorBadRequest))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([a-fA-F0-9\-]+)`,
		httpmock.NewStringResponder(404, errorNotFound))

	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([a-fA-F0-9\-]+)`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

// CleanupMockState cleans up the mock state (called after each test)
func (m *GroupMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.groups = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
}
