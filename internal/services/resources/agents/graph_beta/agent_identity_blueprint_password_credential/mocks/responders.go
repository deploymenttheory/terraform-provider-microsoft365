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
	credentials map[string]map[string]any // blueprintId -> keyId -> credential
}

func init() {
	mockState.credentials = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("agent_identity_blueprint_password_credential", &AgentIdentityBlueprintPasswordCredentialMock{})
}

type AgentIdentityBlueprintPasswordCredentialMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintPasswordCredentialMock)(nil)

func getJSONFileForDisplayName(displayName string) string {
	return fmt.Sprintf("post_password_credential_%s_success.json", displayName)
}

func (m *AgentIdentityBlueprintPasswordCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.Unlock()

	// Add password credential - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addPassword
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/addPassword$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		passwordCredential, ok := requestBody["passwordCredential"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"passwordCredential is required"}}`), nil
		}

		displayName := ""
		if dn, ok := passwordCredential["displayName"].(string); ok {
			displayName = dn
		}

		// Load JSON response from file
		jsonFileName := getJSONFileForDisplayName(displayName)
		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/" + jsonFileName)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate unique keyId and secretText for each call
		keyId := uuid.New().String()
		secretText := "generatedSecretText~" + uuid.New().String()[:8]
		responseObj["keyId"] = keyId
		responseObj["secretText"] = secretText
		responseObj["hint"] = secretText[:3]

		// Use dates from request if provided
		if sd, ok := passwordCredential["startDateTime"].(string); ok && sd != "" {
			responseObj["startDateTime"] = sd
		}
		if ed, ok := passwordCredential["endDateTime"].(string); ok && ed != "" {
			responseObj["endDateTime"] = ed
		}

		// Store in mock state
		mockState.Lock()
		if mockState.credentials[blueprintId] == nil {
			mockState.credentials[blueprintId] = make(map[string]any)
		}
		mockState.credentials[blueprintId][keyId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Remove password credential - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removePassword
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/removePassword$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		keyId, ok := requestBody["keyId"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"keyId is required"}}`), nil
		}

		mockState.Lock()
		if creds, exists := mockState.credentials[blueprintId]; exists {
			delete(creds, keyId)
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityBlueprintPasswordCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/addPassword$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/removePassword$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityBlueprintPasswordCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.Unlock()
}
