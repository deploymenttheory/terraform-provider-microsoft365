package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	credentials map[string]map[string]any // key: credentialId, value: credential data
	blueprints  map[string][]string       // key: blueprintId, value: list of credential IDs
}

func init() {
	mockState.credentials = make(map[string]map[string]any)
	mockState.blueprints = make(map[string][]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("agent_identity_blueprint_federated_identity_credential", &AgentIdentityBlueprintFederatedIdentityCredentialMock{})
}

type AgentIdentityBlueprintFederatedIdentityCredentialMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintFederatedIdentityCredentialMock)(nil)

func getJSONFileForName(name string) string {
	return fmt.Sprintf("post_federated_identity_credential_%s_success.json", name)
}

func (m *AgentIdentityBlueprintFederatedIdentityCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.blueprints = make(map[string][]string)
	mockState.Unlock()

	// Create federated identity credential - POST /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Verify name is provided
		name, ok := requestBody["name"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"name is required"}}`), nil
		}

		// Determine which JSON file to load based on name
		jsonFileName := getJSONFileForName(name)

		// Load JSON response from file
		responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate UUID for the new credential
		newId := uuid.New().String()
		responseObj["id"] = newId

		// Update from request body
		for key, value := range requestBody {
			responseObj[key] = value
		}

		// Store in mock state
		mockState.Lock()
		mockState.credentials[newId] = responseObj
		if mockState.blueprints[blueprintId] == nil {
			mockState.blueprints[blueprintId] = []string{}
		}
		mockState.blueprints[blueprintId] = append(mockState.blueprints[blueprintId], newId)
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get federated identity credential - GET /applications/{blueprintId}/federatedIdentityCredentials/{credentialId}
	// Note: Read uses standard endpoint (not cast endpoint) - cast endpoint doesn't work for reading
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		credentialId := parts[len(parts)-1]

		mockState.Lock()
		credential, exists := mockState.credentials[credentialId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, credential)
	})

	// Update federated identity credential - PATCH /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		credentialId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		credential, exists := mockState.credentials[credentialId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			credential[key] = value
		}

		mockState.credentials[credentialId] = credential
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete federated identity credential - DELETE /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]
		credentialId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.credentials[credentialId]
		if exists {
			delete(mockState.credentials, credentialId)

			// Remove from blueprint's credential list
			if credList, ok := mockState.blueprints[blueprintId]; ok {
				for i, id := range credList {
					if id == credentialId {
						mockState.blueprints[blueprintId] = append(credList[:i], credList[i+1:]...)
						break
					}
				}
			}
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityBlueprintFederatedIdentityCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder(constants.TfOperationDelete, `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityBlueprintFederatedIdentityCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.blueprints = make(map[string][]string)
	mockState.Unlock()
}
