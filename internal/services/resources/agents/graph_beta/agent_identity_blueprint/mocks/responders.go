package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	blueprints   map[string]map[string]any
	sponsors     map[string][]string
	owners       map[string][]string
	deletedItems map[string]map[string]any
}

func init() {
	mockState.blueprints = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("agent_identity_blueprint", &AgentIdentityBlueprintMock{})
}

type AgentIdentityBlueprintMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintMock)(nil)

func getCreateJSONFileForDisplayName(displayName string) string {
	return fmt.Sprintf("../tests/responses/validate_create/post_agent_identity_blueprint_%s_success.json", displayName)
}

func getUpdateJSONFileForDisplayName(displayName string) string {
	return fmt.Sprintf("../tests/responses/validate_update/post_agent_identity_blueprint_%s_success.json", displayName)
}

func (m *AgentIdentityBlueprintMock) RegisterMocks() {
	mockState.Lock()
	mockState.blueprints = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	// Mock user validation - GET /users/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		userId := parts[len(parts)-1]

		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_user_success.json")
		if err != nil {
			// Fallback if JSON file doesn't exist
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          userId,
			})
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          userId,
			})
		}
		responseObj["id"] = userId
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Create agent identity blueprint - POST /applications
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/applications", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Verify @odata.type is set correctly
		odataType, ok := requestBody["@odata.type"].(string)
		if !ok || odataType != "#microsoft.graph.agentIdentityBlueprint" {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Determine which JSON file to load based on displayName
		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		jsonFilePath := getCreateJSONFileForDisplayName(displayName)

		// Load JSON response from file
		jsonContent, err := helpers.ParseJSONFile(jsonFilePath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file %s: %s"}}`, jsonFilePath, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate UUIDs for the new resource
		newId := uuid.New().String()
		newAppId := uuid.New().String()
		responseObj["id"] = newId
		responseObj["appId"] = newAppId

		// Extract sponsors and owners from request
		var sponsorIds []string
		var ownerIds []string

		if sponsorsBinding, ok := requestBody["sponsors@odata.bind"].([]any); ok {
			for _, s := range sponsorsBinding {
				if url, ok := s.(string); ok {
					parts := strings.Split(url, "/")
					if len(parts) > 0 {
						sponsorIds = append(sponsorIds, parts[len(parts)-1])
					}
				}
			}
		}

		if ownersBinding, ok := requestBody["owners@odata.bind"].([]any); ok {
			for _, o := range ownersBinding {
				if url, ok := o.(string); ok {
					parts := strings.Split(url, "/")
					if len(parts) > 0 {
						ownerIds = append(ownerIds, parts[len(parts)-1])
					}
				}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.blueprints[newId] = responseObj
		mockState.sponsors[newId] = sponsorIds
		mockState.owners[newId] = ownerIds
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get agent identity blueprint - GET /applications/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[len(parts)-1]

		mockState.Lock()
		blueprint, exists := mockState.blueprints[blueprintId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewJsonResponse(200, blueprint)
	})

	// Get sponsors - GET /applications/{id}/microsoft.graph.agentIdentityBlueprint/sponsors
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/sponsors$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		mockState.Lock()
		sponsorIds, exists := mockState.sponsors[blueprintId]
		mockState.Unlock()

		if !exists {
			sponsorIds = []string{}
		}

		// Try to load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_sponsors_success.json")
		if err == nil {
			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err == nil {
				// Override the value with actual sponsor IDs from mock state
				values := make([]map[string]any, len(sponsorIds))
				for i, id := range sponsorIds {
					values[i] = map[string]any{
						"@odata.type": "#microsoft.graph.user",
						"id":          id,
					}
				}
				responseObj["value"] = values
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Fallback if JSON file doesn't exist
		values := make([]map[string]any, len(sponsorIds))
		for i, id := range sponsorIds {
			values[i] = map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          id,
			}
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          values,
		})
	})

	// Get owners - GET /applications/{id}/owners
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		mockState.Lock()
		ownerIds, exists := mockState.owners[blueprintId]
		mockState.Unlock()

		if !exists {
			ownerIds = []string{}
		}

		// Try to load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_owners_success.json")
		if err == nil {
			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err == nil {
				// Override the value with actual owner IDs from mock state
				values := make([]map[string]any, len(ownerIds))
				for i, id := range ownerIds {
					values[i] = map[string]any{
						"@odata.type": "#microsoft.graph.user",
						"id":          id,
					}
				}
				responseObj["value"] = values
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Fallback if JSON file doesn't exist
		values := make([]map[string]any, len(ownerIds))
		for i, id := range ownerIds {
			values[i] = map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          id,
			}
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          values,
		})
	})

	// Update agent identity blueprint - PATCH /applications/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		mockState.Lock()
		blueprint, exists := mockState.blueprints[blueprintId]
		if !exists {
			mockState.Unlock()
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		// If displayName is being updated, try to load from validate_update
		if displayName, ok := requestBody["displayName"].(string); ok {
			jsonFilePath := getUpdateJSONFileForDisplayName(displayName)
			jsonContent, err := helpers.ParseJSONFile(jsonFilePath)
			if err == nil {
				var updateResponse map[string]any
				if err := json.Unmarshal([]byte(jsonContent), &updateResponse); err == nil {
					// Merge update response into existing blueprint
					for key, value := range updateResponse {
						if key != "id" && key != "appId" {
							blueprint[key] = value
						}
					}
				}
			}
		}

		// Update fields from request
		for key, value := range requestBody {
			blueprint[key] = value
		}

		mockState.blueprints[blueprintId] = blueprint
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Add sponsor - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/sponsors/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/sponsors/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		odataId, ok := requestBody["@odata.id"].(string)
		if !ok {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		idParts := strings.Split(odataId, "/")
		sponsorId := idParts[len(idParts)-1]

		mockState.Lock()
		mockState.sponsors[blueprintId] = append(mockState.sponsors[blueprintId], sponsorId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove sponsor - DELETE /applications/{id}/microsoft.graph.agentIdentityBlueprint/sponsors/{sponsorId}/$ref
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/sponsors/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]
		sponsorId := parts[6]

		mockState.Lock()
		sponsors := mockState.sponsors[blueprintId]
		for i, id := range sponsors {
			if id == sponsorId {
				mockState.sponsors[blueprintId] = append(sponsors[:i], sponsors[i+1:]...)
				break
			}
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Add owner - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/owners/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/owners/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		odataId, ok := requestBody["@odata.id"].(string)
		if !ok {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		idParts := strings.Split(odataId, "/")
		ownerId := idParts[len(idParts)-1]

		mockState.Lock()
		mockState.owners[blueprintId] = append(mockState.owners[blueprintId], ownerId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove owner - DELETE /applications/{id}/microsoft.graph.agentIdentityBlueprint/owners/{ownerId}/$ref
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/owners/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]
		ownerId := parts[6]

		mockState.Lock()
		owners := mockState.owners[blueprintId]
		for i, id := range owners {
			if id == ownerId {
				mockState.owners[blueprintId] = append(owners[:i], owners[i+1:]...)
				break
			}
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete agent identity blueprint (soft delete) - DELETE /applications/{id}
	// Moves item to deletedItems collection instead of permanently deleting
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[len(parts)-1]

		mockState.Lock()
		blueprint, exists := mockState.blueprints[blueprintId]
		if exists {
			// Move to deletedItems (soft delete behavior)
			mockState.deletedItems[blueprintId] = blueprint
			delete(mockState.blueprints, blueprintId)
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

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

		return httpmock.NewJsonResponse(200, deletedItem)
	})

	// Permanent delete from deleted items - DELETE /directory/deletedItems/{id}
	// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		resourceId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.deletedItems[resourceId]
		if exists {
			delete(mockState.deletedItems, resourceId)
			delete(mockState.sponsors, resourceId)
			delete(mockState.owners, resourceId)
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityBlueprintMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/applications",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentIdentityBlueprintMock) CleanupMockState() {
	mockState.Lock()
	mockState.blueprints = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()
}
