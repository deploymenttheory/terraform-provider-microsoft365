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
	servicePrincipals map[string]map[string]any
}

func init() {
	mockState.servicePrincipals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("agent_identity_blueprint_service_principal", &AgentIdentityBlueprintServicePrincipalMock{})
}

type AgentIdentityBlueprintServicePrincipalMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintServicePrincipalMock)(nil)

func (m *AgentIdentityBlueprintServicePrincipalMock) RegisterMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()

	// Create agent identity blueprint service principal - POST /servicePrincipals
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Verify @odata.type is set correctly
		odataType, ok := requestBody["@odata.type"].(string)
		if !ok || odataType != "#microsoft.graph.agentIdentityBlueprintPrincipal" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"@odata.type must be #microsoft.graph.agentIdentityBlueprintPrincipal"}}`), nil
		}

		// Verify appId is provided
		appId, ok := requestBody["appId"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"appId is required"}}`), nil
		}

		// Load JSON response from file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_service_principal_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate UUIDs for the new resource
		newId := uuid.New().String()
		responseObj["id"] = newId
		responseObj["appId"] = appId
		responseObj["@odata.type"] = "#microsoft.graph.agentIdentityBlueprintPrincipal"

		// Store in mock state
		mockState.Lock()
		mockState.servicePrincipals[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get agent identity blueprint service principal - GET /servicePrincipals/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		servicePrincipalId := parts[len(parts)-1]

		mockState.Lock()
		servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, servicePrincipal)
	})

	// Update agent identity blueprint service principal - PATCH /servicePrincipals/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		servicePrincipalId := parts[len(parts)-1]

		mockState.Lock()
		servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Load JSON response from file
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_update/patch_service_principal_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Merge request body into stored state and response
		for key, value := range requestBody {
			servicePrincipal[key] = value
			responseObj[key] = value
		}

		// Preserve ID and type
		responseObj["id"] = servicePrincipalId
		responseObj["@odata.type"] = "#microsoft.graph.agentIdentityBlueprintPrincipal"

		// Update mock state
		mockState.Lock()
		mockState.servicePrincipals[servicePrincipalId] = servicePrincipal
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Delete agent identity blueprint service principal - DELETE /servicePrincipals/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		servicePrincipalId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.servicePrincipals[servicePrincipalId]
		if exists {
			delete(mockState.servicePrincipals, servicePrincipalId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityBlueprintServicePrincipalMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityBlueprintServicePrincipalMock) CleanupMockState() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()
}
