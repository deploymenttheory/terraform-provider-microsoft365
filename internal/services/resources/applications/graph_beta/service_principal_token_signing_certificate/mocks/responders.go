package mocks

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	servicePrincipals map[string]map[string]any // key: servicePrincipalId
	certCounter       int
}

func init() {
	mockState.servicePrincipals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("service_principal_token_signing_certificate", &ServicePrincipalTokenSigningCertificateMock{})
}

type ServicePrincipalTokenSigningCertificateMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalTokenSigningCertificateMock)(nil)

func (m *ServicePrincipalTokenSigningCertificateMock) RegisterMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.certCounter = 0

	// Seed mock service principal with an unrelated credential pair: delete must remove
	// only the certificate's own credentials and retain everything else.
	mockState.servicePrincipals["11111111-1111-1111-1111-111111111111"] = map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals/$entity",
		"id":             "11111111-1111-1111-1111-111111111111",
		"appId":          "22222222-2222-2222-2222-222222222222",
		"displayName":    "Test Service Principal",
		"keyCredentials": []any{
			map[string]any{
				"customKeyIdentifier": base64.StdEncoding.EncodeToString([]byte("unrelated-credential")),
				"displayName":         "CN=Unrelated Credential",
				"endDateTime":         "2030-01-01T00:00:00.0000000Z",
				"key":                 nil,
				"keyId":               "99999999-9999-9999-9999-999999999999",
				"startDateTime":       "2026-01-01T00:00:00.1234567Z",
				"type":                "AsymmetricX509Cert",
				"usage":               "Verify",
			},
		},
		"passwordCredentials": []any{
			map[string]any{
				"customKeyIdentifier": base64.StdEncoding.EncodeToString([]byte("unrelated-credential")),
				"displayName":         "Unrelated Password",
				"endDateTime":         "2030-01-01T00:00:00.0000000Z",
				"hint":                nil,
				"keyId":               "88888888-8888-8888-8888-888888888888",
				"secretText":          nil,
				"startDateTime":       "2026-01-01T00:00:00.1234567Z",
			},
		},
	}
	mockState.Unlock()

	// Generate token signing certificate - POST /servicePrincipals/{id}/addTokenSigningCertificate
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/addTokenSigningCertificate$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			servicePrincipalId := parts[len(parts)-2]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			displayName := "CN=Microsoft Azure Federated SSO Certificate"
			if name, ok := requestBody["displayName"].(string); ok && name != "" {
				displayName = name
			}
			startDateTime := "2026-01-01T00:00:00Z"
			endDateTime := "2029-01-01T00:00:00Z"
			if end, ok := requestBody["endDateTime"].(string); ok && end != "" {
				if parsed, err := time.Parse(time.RFC3339, end); err == nil {
					endDateTime = parsed.UTC().Format(time.RFC3339)
				}
			}

			mockState.certCounter++
			signKeyId := uuid.New().String()
			verifyKeyId := uuid.New().String()
			// The password credential shares the Sign key credential's keyId (live API behavior)
			passwordKeyId := signKeyId
			thumbprint := fmt.Sprintf("%040x", mockState.certCounter)
			// customKeyIdentifier is the raw SHA-1 thumbprint bytes, base64-encoded in JSON
			// (live API behavior)
			rawThumbprint, _ := hex.DecodeString(thumbprint)
			customKeyIdentifier := base64.StdEncoding.EncodeToString(rawThumbprint)
			key := base64.StdEncoding.EncodeToString([]byte("mock-certificate-" + thumbprint))

			signCredential := map[string]any{
				"customKeyIdentifier": customKeyIdentifier,
				"displayName":         displayName,
				"endDateTime":         endDateTime,
				"key":                 key,
				"keyId":               signKeyId,
				"startDateTime":       startDateTime,
				"type":                "AsymmetricX509Cert",
				"usage":               "Sign",
			}
			verifyCredential := map[string]any{
				"customKeyIdentifier": customKeyIdentifier,
				"displayName":         displayName,
				"endDateTime":         endDateTime,
				"key":                 key,
				"keyId":               verifyKeyId,
				"startDateTime":       startDateTime,
				"type":                "AsymmetricX509Cert",
				"usage":               "Verify",
			}
			passwordCredential := map[string]any{
				"customKeyIdentifier": customKeyIdentifier,
				"displayName":         displayName,
				"endDateTime":         endDateTime,
				"keyId":               passwordKeyId,
				"startDateTime":       startDateTime,
			}

			servicePrincipal["keyCredentials"] = append(servicePrincipal["keyCredentials"].([]any), signCredential, verifyCredential)
			servicePrincipal["passwordCredentials"] = append(servicePrincipal["passwordCredentials"].([]any), passwordCredential)
			mockState.servicePrincipals[servicePrincipalId] = servicePrincipal

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context":      "https://graph.microsoft.com/beta/$metadata#microsoft.graph.selfSignedCertificate",
				"customKeyIdentifier": customKeyIdentifier,
				"displayName":         displayName,
				"endDateTime":         endDateTime,
				"key":                 key,
				"keyId":               signKeyId,
				"startDateTime":       startDateTime,
				"thumbprint":          thumbprint,
				"type":                "AsymmetricX509Cert",
				"usage":               "Sign",
			})
		})

	// Get service principal - GET /servicePrincipals/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
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

	// Update service principal - PATCH /servicePrincipals/{id} (credential removal on delete)
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			servicePrincipalId := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			if keyCredentials, ok := requestBody["keyCredentials"].([]any); ok {
				servicePrincipal["keyCredentials"] = keyCredentials
			}
			if passwordCredentials, ok := requestBody["passwordCredentials"].([]any); ok {
				servicePrincipal["passwordCredentials"] = passwordCredentials
			}

			mockState.servicePrincipals[servicePrincipalId] = servicePrincipal

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *ServicePrincipalTokenSigningCertificateMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+/addTokenSigningCertificate$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ServicePrincipalTokenSigningCertificateMock) CleanupMockState() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.certCounter = 0
	mockState.Unlock()
}

// GetCredentialCounts returns the number of key and password credentials stored on a
// service principal, for test assertions on delete behavior.
func GetCredentialCounts(servicePrincipalId string) (keyCredentials int, passwordCredentials int, ok bool) {
	mockState.Lock()
	defer mockState.Unlock()

	servicePrincipal, exists := mockState.servicePrincipals[servicePrincipalId]
	if !exists {
		return 0, 0, false
	}
	keys, _ := servicePrincipal["keyCredentials"].([]any)
	passwords, _ := servicePrincipal["passwordCredentials"].([]any)
	return len(keys), len(passwords), true
}
