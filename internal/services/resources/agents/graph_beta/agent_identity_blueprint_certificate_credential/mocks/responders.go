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
	applications map[string]map[string]any
}

func init() {
	mockState.applications = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("agent_identity_blueprint_certificate_credential", &AgentIdentityBlueprintCertificateCredentialMock{})
}

type AgentIdentityBlueprintCertificateCredentialMock struct{}

var _ mocks.MockRegistrar = (*AgentIdentityBlueprintCertificateCredentialMock)(nil)

func (m *AgentIdentityBlueprintCertificateCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()

	// Initialize the test application with empty keyCredentials
	initializeTestApplication()

	// GET application - used to fetch existing keyCredentials before PATCH
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[len(parts)-1]

		mockState.Lock()
		application, exists := mockState.applications[applicationId]
		mockState.Unlock()

		if !exists {
			// Load from validate_read JSON for known test IDs
			jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_application_success.json")
			if err != nil {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON: %s"}}`, err.Error())), nil
			}

			responseObj["id"] = applicationId

			// Store in mock state
			mockState.Lock()
			mockState.applications[applicationId] = responseObj
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, responseObj)
		}

		return httpmock.NewJsonResponse(200, application)
	})

	// PATCH application with agentIdentityBlueprint cast - used to add/update/delete certificate
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		// Path is /applications/{id}/microsoft.graph.agentIdentityBlueprint
		applicationId := parts[len(parts)-2]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		application, exists := mockState.applications[applicationId]
		if !exists {
			application = make(map[string]any)
			application["id"] = applicationId
			application["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#applications/$entity"
		}
		mockState.Unlock()

		// Handle keyCredentials in the request
		if keyCredentials, ok := requestBody["keyCredentials"].([]any); ok {
			// Process each key credential and add generated keyId if missing
			processedCredentials := make([]any, 0, len(keyCredentials))
			for _, cred := range keyCredentials {
				if credMap, ok := cred.(map[string]any); ok {
					// Generate keyId if not present
					if _, hasKeyId := credMap["keyId"]; !hasKeyId {
						credMap["keyId"] = uuid.New().String()
					}
					// Add computed fields
					if _, hasThumbprint := credMap["customKeyIdentifier"]; !hasThumbprint {
						credMap["customKeyIdentifier"] = "VGVzdFRodW1icHJpbnQ=" // Base64 encoded test thumbprint
					}
					processedCredentials = append(processedCredentials, credMap)
				}
			}
			application["keyCredentials"] = processedCredentials
		}

		// Update mock state
		mockState.Lock()
		mockState.applications[applicationId] = application
		mockState.Unlock()

		// Return 204 No Content for successful PATCH
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func initializeTestApplication() {
	// Load initial application state from JSON
	jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/get_application_before_create_success.json")
	if err != nil {
		return
	}

	var application map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &application); err != nil {
		return
	}

	applicationId, ok := application["id"].(string)
	if !ok {
		applicationId = "11111111-1111-1111-1111-111111111111"
		application["id"] = applicationId
	}

	mockState.Lock()
	mockState.applications[applicationId] = application
	mockState.Unlock()
}

func (m *AgentIdentityBlueprintCertificateCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/microsoft\.graph\.agentIdentityBlueprint$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *AgentIdentityBlueprintCertificateCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()
}
