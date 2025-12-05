package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	credentials map[string]map[string]map[string]any // blueprintId -> keyId -> credential
}

func init() {
	mockState.credentials = make(map[string]map[string]map[string]any)
	mocks.GlobalRegistry.Register("agent_identity_blueprint_key_credential", &AgentIdentityBlueprintKeyCredentialMock{})
}

type AgentIdentityBlueprintKeyCredentialMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintKeyCredentialMock)(nil)

func getJSONFileForDisplayName(displayName string) string {
	return fmt.Sprintf("post_key_credential_%s_success.json", displayName)
}

func (m *AgentIdentityBlueprintKeyCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]map[string]any)
	mockState.Unlock()

	// Get application with keyCredentials - GET /applications/{id}/microsoft.graph.agentIdentityBlueprint?$select=keyCredentials
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		mockState.Lock()
		defer mockState.Unlock()

		// Build keyCredentials array from stored credentials
		var keyCredentials []map[string]any
		if creds, exists := mockState.credentials[blueprintId]; exists {
			for _, cred := range creds {
				keyCredentials = append(keyCredentials, cred)
			}
		}

		// Return an Application-like response with keyCredentials
		response := map[string]any{
			"@odata.type":    "#microsoft.graph.application",
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#applications(keyCredentials)/$entity",
			"id":             blueprintId,
			"keyCredentials": keyCredentials,
		}

		return httpmock.NewJsonResponse(200, response)
	})

	// Add key credential - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/addKey
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/addKey$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		blueprintId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		keyCredential, ok := requestBody["keyCredential"].(map[string]any)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"keyCredential is required"}}`), nil
		}

		proof, ok := requestBody["proof"].(string)
		if !ok || proof == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"proof is required"}}`), nil
		}

		displayName := ""
		if dn, ok := keyCredential["displayName"].(string); ok {
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

		// Generate unique keyId for each call
		keyId := uuid.New().String()
		responseObj["keyId"] = keyId

		// Set default dates if not provided
		now := time.Now().UTC()
		if responseObj["startDateTime"] == nil {
			responseObj["startDateTime"] = now.Format(time.RFC3339)
		}
		if responseObj["endDateTime"] == nil {
			responseObj["endDateTime"] = now.AddDate(2, 0, 0).Format(time.RFC3339)
		}

		// Use dates from request if provided
		if sd, ok := keyCredential["startDateTime"].(string); ok && sd != "" {
			responseObj["startDateTime"] = sd
		}
		if ed, ok := keyCredential["endDateTime"].(string); ok && ed != "" {
			responseObj["endDateTime"] = ed
		}

		// Store in mock state for Read to find
		mockState.Lock()
		if mockState.credentials[blueprintId] == nil {
			mockState.credentials[blueprintId] = make(map[string]map[string]any)
		}
		mockState.credentials[blueprintId][keyId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Remove key credential - POST /applications/{id}/microsoft.graph.agentIdentityBlueprint/removeKey
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/removeKey$`, func(req *http.Request) (*http.Response, error) {
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

		proof, ok := requestBody["proof"].(string)
		if !ok || proof == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"proof is required"}}`), nil
		}

		mockState.Lock()
		if creds, exists := mockState.credentials[blueprintId]; exists {
			delete(creds, keyId)
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentIdentityBlueprintKeyCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint`,
		httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/addKey$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint/removeKey$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityBlueprintKeyCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]map[string]any)
	mockState.Unlock()
}
