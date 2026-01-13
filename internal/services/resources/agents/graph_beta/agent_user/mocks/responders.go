package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	agentUsers   map[string]map[string]any
	sponsors     map[string][]string       // agentUserId -> []sponsorId
	deletedItems map[string]map[string]any // soft-deleted users awaiting hard delete
}

func init() {
	mockState.agentUsers = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("agent_user", &AgentUserMock{})
}

type AgentUserMock struct{}

var _ mocks.MockRegistrar = (*AgentUserMock)(nil)

func (m *AgentUserMock) RegisterMocks() {
	mockState.Lock()
	mockState.agentUsers = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	// Create agent user - POST /users
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Verify required fields
		displayName, ok := requestBody["displayName"].(string)
		if !ok || displayName == "" {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_display_name_required.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Load JSON response from file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_agent_user_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate UUID for the new resource
		newId := uuid.New().String()

		// Copy all fields from request to response
		for key, value := range requestBody {
			if key != "sponsors@odata.bind" && key != "@odata.type" {
				responseObj[key] = value
			}
		}

		// Ensure required fields are set (override any from request)
		responseObj["id"] = newId
		responseObj["displayName"] = displayName
		responseObj["@odata.type"] = "#microsoft.graph.agentUser"

		// Extract and store sponsors from sponsors@odata.bind
		var sponsorIds []string
		if sponsorBinds, ok := requestBody["sponsors@odata.bind"].([]any); ok {
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			for _, bind := range sponsorBinds {
				if bindStr, ok := bind.(string); ok {
					if match := uuidRegex.FindString(bindStr); match != "" {
						sponsorIds = append(sponsorIds, match)
					}
				}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.agentUsers[newId] = responseObj
		mockState.sponsors[newId] = sponsorIds
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get agent user - GET /users/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[len(parts)-1]

		mockState.Lock()
		storedUser, exists := mockState.agentUsers[agentUserId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		// Load JSON response template from file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_agent_user_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Overlay stored values onto template
		for key, value := range storedUser {
			responseObj[key] = value
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Get sponsors - GET /users/{id}/sponsors
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/sponsors$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[3]

		mockState.Lock()
		sponsorIds := mockState.sponsors[agentUserId]
		mockState.Unlock()

		// Build sponsors response
		sponsorValues := make([]map[string]any, len(sponsorIds))
		for i, sponsorId := range sponsorIds {
			sponsorValues[i] = map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          sponsorId,
			}
		}

		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          sponsorValues,
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Update agent user - PATCH /users/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[len(parts)-1]

		mockState.Lock()
		agentUser, exists := mockState.agentUsers[agentUserId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Merge request body into stored state
		for key, value := range requestBody {
			agentUser[key] = value
		}

		// Preserve ID and type
		agentUser["id"] = agentUserId
		agentUser["@odata.type"] = "#microsoft.graph.agentUser"

		// Update mock state
		mockState.Lock()
		mockState.agentUsers[agentUserId] = agentUser
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete agent user (soft delete) - DELETE /users/{id}
	// Moves user to deletedItems collection instead of permanently deleting
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[len(parts)-1]

		mockState.Lock()
		user, exists := mockState.agentUsers[agentUserId]
		if exists {
			// Move to deletedItems (soft delete behavior)
			mockState.deletedItems[agentUserId] = user
			delete(mockState.agentUsers, agentUserId)
			// Keep sponsors in case of restore (not implemented, but realistic behavior)
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Add sponsor - POST /users/{id}/sponsors/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/sponsors/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		if odataId, ok := requestBody["@odata.id"].(string); ok {
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			if match := uuidRegex.FindString(odataId); match != "" {
				mockState.Lock()
				mockState.sponsors[agentUserId] = append(mockState.sponsors[agentUserId], match)
				mockState.Unlock()
			}
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove sponsor - DELETE /users/{id}/sponsors/{sponsorId}/$ref
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/sponsors/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentUserId := parts[3]
		sponsorId := parts[5] // Changed from 6 to 5 since we removed microsoft.graph.agentUser

		mockState.Lock()
		if sponsors, exists := mockState.sponsors[agentUserId]; exists {
			var newSponsors []string
			for _, s := range sponsors {
				if s != sponsorId {
					newSponsors = append(newSponsors, s)
				}
			}
			mockState.sponsors[agentUserId] = newSponsors
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Get deleted item - GET /directory/deletedItems/{id}
	// Used for soft delete verification (polling until resource appears in deleted items)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		resourceId := parts[len(parts)-1]

		mockState.Lock()
		deletedItem, exists := mockState.deletedItems[resourceId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		// Return the deleted item
		return httpmock.NewJsonResponse(200, deletedItem)
	})

	// Permanent delete from deleted items - DELETE /directory/deletedItems/{id}
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		resourceId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.deletedItems[resourceId]
		if exists {
			delete(mockState.deletedItems, resourceId)
			delete(mockState.sponsors, resourceId) // Also remove sponsors on hard delete
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentUserMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/sponsors$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentUserMock) CleanupMockState() {
	mockState.Lock()
	mockState.agentUsers = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()
}
