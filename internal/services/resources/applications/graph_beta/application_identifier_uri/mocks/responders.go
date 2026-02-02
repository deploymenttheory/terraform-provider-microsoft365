package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	applications map[string]map[string]any // key: applicationId, value: application data including identifierUris
}

func init() {
	mockState.applications = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("application_identifier_uri", &ApplicationIdentifierUriMock{})
}

type ApplicationIdentifierUriMock struct{}

var _ mocks.MockRegistrar = (*ApplicationIdentifierUriMock)(nil)

func (m *ApplicationIdentifierUriMock) RegisterMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)

	// Seed mock application
	mockState.applications["11111111-1111-1111-1111-111111111111"] = map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#applications/$entity",
		"id":             "11111111-1111-1111-1111-111111111111",
		"appId":          "22222222-2222-2222-2222-222222222222",
		"displayName":    "Test Application",
		"identifierUris": []any{},
	}
	mockState.Unlock()

	// Get application - GET /applications/{applicationId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			applicationId := parts[len(parts)-1]

			mockState.Lock()
			application, exists := mockState.applications[applicationId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Application not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, application)
		})

	// Update application - PATCH /applications/{applicationId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			applicationId := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			application, exists := mockState.applications[applicationId]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Application not found"}}`), nil
			}

			// Update identifier URIs if provided
			if identifierUris, ok := requestBody["identifierUris"].([]any); ok {
				application["identifierUris"] = identifierUris
			}

			mockState.applications[applicationId] = application
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *ApplicationIdentifierUriMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ApplicationIdentifierUriMock) CleanupMockState() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()
}
