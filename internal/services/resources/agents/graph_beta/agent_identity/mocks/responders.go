package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	agentIdentities map[string]map[string]any
	sponsors        map[string][]string // agentIdentityId -> []sponsorId
	owners          map[string][]string // agentIdentityId -> []ownerId
}

func init() {
	mockState.agentIdentities = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mocks.GlobalRegistry.Register("agent_identity", &AgentIdentityMock{})
}

type AgentIdentityMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityMock)(nil)

func (m *AgentIdentityMock) RegisterMocks() {
	mockState.Lock()
	mockState.agentIdentities = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mockState.Unlock()

	// Create agent identity - POST /servicePrincipals/microsoft.graph.agentIdentity
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals/microsoft.graph.agentIdentity", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Verify required fields
		displayName, ok := requestBody["displayName"].(string)
		if !ok || displayName == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
		}

		// Load JSON response from file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_agent_identity_success.json")
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
		responseObj["id"] = newId
		responseObj["displayName"] = displayName
		responseObj["@odata.type"] = "#microsoft.graph.agentIdentity"
		responseObj["servicePrincipalType"] = "ServiceIdentity"
		responseObj["accountEnabled"] = true

		// Copy agentIdentityBlueprintId from request
		if blueprintId, ok := requestBody["agentIdentityBlueprintId"].(string); ok {
			responseObj["agentIdentityBlueprintId"] = blueprintId
		}

		// Copy tags from request if present
		if tags, ok := requestBody["tags"].([]any); ok {
			responseObj["tags"] = tags
		}

		// Extract and store sponsors from sponsors@odata.bind
		var sponsorIds []string
		if sponsorBinds, ok := requestBody["sponsors@odata.bind"].([]any); ok {
			// Extract sponsor IDs from URLs like "https://graph.microsoft.com/beta/users/{id}"
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			for _, bind := range sponsorBinds {
				if bindStr, ok := bind.(string); ok {
					if match := uuidRegex.FindString(bindStr); match != "" {
						sponsorIds = append(sponsorIds, match)
					}
				}
			}
		}

		// Extract and store owners from owners@odata.bind
		var ownerIds []string
		if ownerBinds, ok := requestBody["owners@odata.bind"].([]any); ok {
			// Extract owner IDs from URLs like "https://graph.microsoft.com/beta/users/{id}"
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			for _, bind := range ownerBinds {
				if bindStr, ok := bind.(string); ok {
					if match := uuidRegex.FindString(bindStr); match != "" {
						ownerIds = append(ownerIds, match)
					}
				}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.agentIdentities[newId] = responseObj
		mockState.sponsors[newId] = sponsorIds
		mockState.owners[newId] = ownerIds
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get agent identity - GET /servicePrincipals/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[len(parts)-1]

		mockState.Lock()
		agentIdentity, exists := mockState.agentIdentities[agentIdentityId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, agentIdentity)
	})

	// Get sponsors - GET /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentity/sponsors$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors
		// Split: ["", "beta", "servicePrincipals", "{id}", "microsoft.graph.agentIdentity", "sponsors"]
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]

		mockState.Lock()
		sponsorIds, exists := mockState.sponsors[agentIdentityId]
		mockState.Unlock()

		if !exists {
			sponsorIds = []string{}
		}

		// Build response with stored sponsors
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

	// Update agent identity - PATCH /servicePrincipals/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[len(parts)-1]

		mockState.Lock()
		agentIdentity, exists := mockState.agentIdentities[agentIdentityId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Merge request body into stored state
		for key, value := range requestBody {
			agentIdentity[key] = value
		}

		// Preserve ID and type
		agentIdentity["id"] = agentIdentityId
		agentIdentity["@odata.type"] = "#microsoft.graph.agentIdentity"

		// Update mock state
		mockState.Lock()
		mockState.agentIdentities[agentIdentityId] = agentIdentity
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, agentIdentity)
	})

	// Get owners - GET /servicePrincipals/{id}/owners
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/owners
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]

		mockState.Lock()
		ownerIds, exists := mockState.owners[agentIdentityId]
		mockState.Unlock()

		if !exists {
			ownerIds = []string{}
		}

		// Build response with stored owners
		ownerValues := make([]map[string]any, len(ownerIds))
		for i, ownerId := range ownerIds {
			ownerValues[i] = map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          ownerId,
			}
		}

		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          ownerValues,
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Delete agent identity - DELETE /servicePrincipals/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.agentIdentities[agentIdentityId]
		if exists {
			delete(mockState.agentIdentities, agentIdentityId)
			delete(mockState.sponsors, agentIdentityId)
			delete(mockState.owners, agentIdentityId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Add sponsor - POST /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentity/sponsors/\$ref$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/$ref
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]

		// Parse request body to get sponsor ID
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Extract sponsor ID from @odata.id
		if odataId, ok := requestBody["@odata.id"].(string); ok {
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			if match := uuidRegex.FindString(odataId); match != "" {
				mockState.Lock()
				mockState.sponsors[agentIdentityId] = append(mockState.sponsors[agentIdentityId], match)
				mockState.Unlock()
			}
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove sponsor - DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/{sponsorId}/$ref
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentity/sponsors/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/\$ref$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/{sponsorId}/$ref
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]
		sponsorId := parts[6]

		// Remove sponsor from state
		mockState.Lock()
		if sponsors, exists := mockState.sponsors[agentIdentityId]; exists {
			var newSponsors []string
			for _, s := range sponsors {
				if s != sponsorId {
					newSponsors = append(newSponsors, s)
				}
			}
			mockState.sponsors[agentIdentityId] = newSponsors
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Add owner - POST /servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentity/owners/\$ref$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/$ref
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]

		// Parse request body to get owner ID
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Extract owner ID from @odata.id
		if odataId, ok := requestBody["@odata.id"].(string); ok {
			uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
			if match := uuidRegex.FindString(odataId); match != "" {
				mockState.Lock()
				mockState.owners[agentIdentityId] = append(mockState.owners[agentIdentityId], match)
				mockState.Unlock()
			}
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove owner - DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/{ownerId}/$ref
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentity/owners/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/\$ref$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/servicePrincipals/{id}/microsoft.graph.agentIdentity/owners/{ownerId}/$ref
		parts := strings.Split(req.URL.Path, "/")
		agentIdentityId := parts[3]
		ownerId := parts[6]

		// Remove owner from state
		mockState.Lock()
		if owners, exists := mockState.owners[agentIdentityId]; exists {
			var newOwners []string
			for _, o := range owners {
				if o != ownerId {
					newOwners = append(newOwners, o)
				}
			}
			mockState.owners[agentIdentityId] = newOwners
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Permanent delete from deleted items - DELETE /directory/deletedItems/{id}
	// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		// Path: /beta/directory/deletedItems/{id}
		// The item was already removed from agentIdentities in the soft delete step
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals/microsoft.graph.agentIdentity",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityMock) CleanupMockState() {
	mockState.Lock()
	mockState.agentIdentities = make(map[string]map[string]any)
	mockState.sponsors = make(map[string][]string)
	mockState.owners = make(map[string][]string)
	mockState.Unlock()
}
