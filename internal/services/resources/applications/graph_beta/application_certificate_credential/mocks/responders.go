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
	mocks.GlobalRegistry.Register("application_certificate_credential", &ApplicationCertificateCredentialMock{})
}

type ApplicationCertificateCredentialMock struct{}

var _ mocks.MockRegistrar = (*ApplicationCertificateCredentialMock)(nil)

func (m *ApplicationCertificateCredentialMock) RegisterMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()

	// Initialize the test applications with empty keyCredentials
	initializeTestApplications()

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

	// PATCH application - used to add/update/delete certificate
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[len(parts)-1]

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

func initializeTestApplications() {
	// Initialize test applications for different encoding types
	testApplicationIds := map[string]string{
		"11111111-1111-1111-1111-111111111111": "unit-test-application-pem",
		"22222222-2222-2222-2222-222222222222": "unit-test-application-base64",
		"33333333-3333-3333-3333-333333333333": "unit-test-application-der",
		"44444444-4444-4444-4444-444444444444": "unit-test-application-hex",
		"55555555-5555-5555-5555-555555555555": "unit-test-application-replace",
	}

	// Load initial application state from JSON
	jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/get_application_before_create_success.json")
	if err != nil {
		return
	}

	var baseApplication map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &baseApplication); err != nil {
		return
	}

	mockState.Lock()
	for applicationId, displayName := range testApplicationIds {
		application := make(map[string]any)
		for k, v := range baseApplication {
			application[k] = v
		}
		application["id"] = applicationId
		application["displayName"] = displayName
		application["keyCredentials"] = []any{} // Start with empty keyCredentials
		mockState.applications[applicationId] = application
	}
	mockState.Unlock()
}

func (m *ApplicationCertificateCredentialMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ApplicationCertificateCredentialMock) CleanupMockState() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()
}

// AddPreExistingCertificates adds pre-existing certificates to an application for testing
func (m *ApplicationCertificateCredentialMock) AddPreExistingCertificates(applicationID string, certificates []map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()

	if app, exists := mockState.applications[applicationID]; exists {
		existingCreds := []any{}
		if creds, ok := app["keyCredentials"].([]any); ok {
			existingCreds = creds
		}
		
		// Add new pre-existing certificates
		for _, cert := range certificates {
			existingCreds = append(existingCreds, cert)
		}
		
		app["keyCredentials"] = existingCreds
		mockState.applications[applicationID] = app
	}
}

// GetCertificateCount returns the number of certificates on an application
func (m *ApplicationCertificateCredentialMock) GetCertificateCount(applicationID string) int {
	mockState.Lock()
	defer mockState.Unlock()

	if app, exists := mockState.applications[applicationID]; exists {
		if creds, ok := app["keyCredentials"].([]any); ok {
			return len(creds)
		}
	}
	return 0
}

// GetCertificates returns all certificates on an application
func (m *ApplicationCertificateCredentialMock) GetCertificates(applicationID string) []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()

	var result []map[string]any
	if app, exists := mockState.applications[applicationID]; exists {
		if creds, ok := app["keyCredentials"].([]any); ok {
			for _, cred := range creds {
				if credMap, ok := cred.(map[string]any); ok {
					result = append(result, credMap)
				}
			}
		}
	}
	return result
}
