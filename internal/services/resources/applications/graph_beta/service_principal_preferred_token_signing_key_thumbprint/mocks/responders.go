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
	servicePrincipals map[string]map[string]any // key: servicePrincipalId
}

func init() {
	mockState.servicePrincipals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("service_principal_preferred_token_signing_key_thumbprint", &ServicePrincipalPreferredTokenSigningKeyThumbprintMock{})
}

type ServicePrincipalPreferredTokenSigningKeyThumbprintMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalPreferredTokenSigningKeyThumbprintMock)(nil)

func (m *ServicePrincipalPreferredTokenSigningKeyThumbprintMock) RegisterMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)

	// Seed mock service principal
	mockState.servicePrincipals["11111111-1111-1111-1111-111111111111"] = map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals/$entity",
		"id":             "11111111-1111-1111-1111-111111111111",
		"appId":          "22222222-2222-2222-2222-222222222222",
		"displayName":    "Test Service Principal",
	}
	mockState.Unlock()

	// Get service principal - GET /servicePrincipals/{servicePrincipalId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			servicePrincipalId := parts[len(parts)-1]

			mockState.Lock()
			servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, servicePrincipal)
		})

	// Update service principal - PATCH /servicePrincipals/{servicePrincipalId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			servicePrincipalId := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			// Distinguish an explicit JSON null (clears the property, sent by Delete)
			// from an absent key. Graph normalizes the stored thumbprint to lowercase.
			if raw, keyPresent := requestBody["preferredTokenSigningKeyThumbprint"]; keyPresent {
				if raw == nil {
					delete(servicePrincipal, "preferredTokenSigningKeyThumbprint")
				} else if thumbprint, ok := raw.(string); ok {
					servicePrincipal["preferredTokenSigningKeyThumbprint"] = strings.ToLower(thumbprint)
				}
			}

			mockState.servicePrincipals[servicePrincipalId] = servicePrincipal
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *ServicePrincipalPreferredTokenSigningKeyThumbprintMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ServicePrincipalPreferredTokenSigningKeyThumbprintMock) CleanupMockState() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetThumbprint returns the stored thumbprint for a service principal, for test assertions.
func GetThumbprint(servicePrincipalId string) (string, bool) {
	mockState.Lock()
	defer mockState.Unlock()

	servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
	if !exists {
		return "", false
	}
	thumbprint, ok := servicePrincipal["preferredTokenSigningKeyThumbprint"].(string)
	return thumbprint, ok
}
