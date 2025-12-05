package mocks

import (
	"encoding/json"
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
	applications   map[string]map[string]any   // blueprintId -> application data
	identifierUris map[string][]string         // blueprintId -> []identifierUri
	scopes         map[string][]map[string]any // blueprintId -> []scope
	isCreated      map[string]bool             // blueprintId -> has identifier URI been created
}

func init() {
	mockState.applications = make(map[string]map[string]any)
	mockState.identifierUris = make(map[string][]string)
	mockState.scopes = make(map[string][]map[string]any)
	mockState.isCreated = make(map[string]bool)
	mocks.GlobalRegistry.Register("agent_identity_blueprint_identifier_uri", &AgentIdentityBlueprintIdentifierUriMock{})
}

type AgentIdentityBlueprintIdentifierUriMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintIdentifierUriMock)(nil)

func (m *AgentIdentityBlueprintIdentifierUriMock) RegisterMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.identifierUris = make(map[string][]string)
	mockState.scopes = make(map[string][]map[string]any)
	mockState.isCreated = make(map[string]bool)
	mockState.Unlock()

	// GET /applications/{id} - Read application
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		mockState.Lock()
		isCreated := mockState.isCreated[blueprintId]
		mockState.Unlock()

		var jsonFile string
		if isCreated {
			jsonFile = "../tests/responses/validate_read/get_application_success.json"
		} else {
			jsonFile = "../tests/responses/validate_create/get_application_before_create_success.json"
		}

		jsonContent, err := helpers.ParseJSONFile(jsonFile)
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Override the ID to match the request
		responseObj["id"] = blueprintId

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// PATCH /applications/{id} - Update application (add/remove identifier URI)
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		// Check if this is a delete operation (identifierUris is empty)
		if uris, ok := requestBody["identifierUris"].([]any); ok && len(uris) == 0 {
			// Delete operation - load delete response
			jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_delete/patch_application_success.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}

			responseObj["id"] = blueprintId
			mockState.isCreated[blueprintId] = false
			mockState.identifierUris[blueprintId] = []string{}
			mockState.scopes[blueprintId] = []map[string]any{}

			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Create or Update operation
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/patch_application_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		// Update state based on request
		if uris, ok := requestBody["identifierUris"].([]any); ok {
			mockState.identifierUris[blueprintId] = make([]string, 0, len(uris))
			for _, uri := range uris {
				if uriStr, ok := uri.(string); ok {
					mockState.identifierUris[blueprintId] = append(mockState.identifierUris[blueprintId], uriStr)
				}
			}
			responseObj["identifierUris"] = mockState.identifierUris[blueprintId]
		}

		if api, ok := requestBody["api"].(map[string]any); ok {
			if scopes, ok := api["oauth2PermissionScopes"].([]any); ok {
				mockState.scopes[blueprintId] = make([]map[string]any, 0, len(scopes))
				for _, scope := range scopes {
					if scopeMap, ok := scope.(map[string]any); ok {
						// Generate ID if not provided
						if _, hasId := scopeMap["id"]; !hasId {
							scopeMap["id"] = uuid.New().String()
						}
						mockState.scopes[blueprintId] = append(mockState.scopes[blueprintId], scopeMap)
					}
				}
				responseObj["api"] = map[string]any{
					"oauth2PermissionScopes": mockState.scopes[blueprintId],
				}
			}
		}

		responseObj["id"] = blueprintId
		mockState.isCreated[blueprintId] = true

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *AgentIdentityBlueprintIdentifierUriMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`))
}

func (m *AgentIdentityBlueprintIdentifierUriMock) CleanupMockState() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.identifierUris = make(map[string][]string)
	mockState.scopes = make(map[string][]map[string]any)
	mockState.isCreated = make(map[string]bool)
	mockState.Unlock()
}
