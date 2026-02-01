package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	credentials  map[string]map[string]any // key: credentialId, value: credential data
	applications map[string][]string       // key: applicationId, value: list of credential IDs
}

func init() {
	mockState.credentials = make(map[string]map[string]any)
	mockState.applications = make(map[string][]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("application_federated_identity_credential", &ApplicationFederatedIdentityCredentialMock{})
}

type ApplicationFederatedIdentityCredentialMock struct{}

var _ mocks.MockRegistrar = (*ApplicationFederatedIdentityCredentialMock)(nil)

func getJSONFileForName(name string) string {
	return fmt.Sprintf("post_federated_identity_credential_%s_success.json", name)
}

func (m *ApplicationFederatedIdentityCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.applications = make(map[string][]string)
	mockState.Unlock()

	// Create federated identity credential - POST /applications/{applicationId}/federatedIdentityCredentials
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]

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
		if mockState.applications[applicationId] == nil {
			mockState.applications[applicationId] = []string{}
		}
		mockState.applications[applicationId] = append(mockState.applications[applicationId], newId)
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get federated identity credential - GET /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
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

	// Update federated identity credential - PATCH /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
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

	// Delete federated identity credential - DELETE /applications/{applicationId}/federatedIdentityCredentials/{credentialId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]
		credentialId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.credentials[credentialId]
		if exists {
			delete(mockState.credentials, credentialId)

			// Remove from application's credential list
			if credList, ok := mockState.applications[applicationId]; ok {
				for i, id := range credList {
					if id == credentialId {
						mockState.applications[applicationId] = append(credList[:i], credList[i+1:]...)
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

func (m *ApplicationFederatedIdentityCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/federatedIdentityCredentials/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ApplicationFederatedIdentityCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.credentials = make(map[string]map[string]any)
	mockState.applications = make(map[string][]string)
	mockState.Unlock()
}
