package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	servicePrincipals map[string]map[string]any
	deletedItems      map[string]map[string]any
}

func init() {
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("agent_identity_blueprint_service_principal", &AgentIdentityBlueprintServicePrincipalMock{})
}

type AgentIdentityBlueprintServicePrincipalMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintServicePrincipalMock)(nil)

func (m *AgentIdentityBlueprintServicePrincipalMock) RegisterMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()

	// Create agent identity blueprint service principal - POST /servicePrincipals
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Verify @odata.type is set correctly
		odataType, ok := requestBody["@odata.type"].(string)
		if !ok || odataType != "#microsoft.graph.agentIdentityBlueprintPrincipal" {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Verify appId is provided
		appId, ok := requestBody["appId"].(string)
		if !ok {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
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
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
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
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
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

	// Delete agent identity blueprint service principal (soft delete) - DELETE /servicePrincipals/{id}
	// Moves item to deletedItems collection instead of permanently deleting
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		servicePrincipalId := parts[len(parts)-1]

		mockState.Lock()
		servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
		if exists {
			// Move to deletedItems (soft delete behavior)
			mockState.deletedItems[servicePrincipalId] = servicePrincipal
			delete(mockState.servicePrincipals, servicePrincipalId)
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
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		resourceId := parts[len(parts)-1]

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

func (m *AgentIdentityBlueprintServicePrincipalMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/servicePrincipals",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentIdentityBlueprintServicePrincipalMock) CleanupMockState() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
	mockState.Unlock()
}
