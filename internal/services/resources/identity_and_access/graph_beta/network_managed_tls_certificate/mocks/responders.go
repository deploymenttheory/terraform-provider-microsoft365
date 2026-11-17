package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
)

const managedTLSCertificateID = "00000000-0000-0000-0000-000000000301"

var mockState struct {
	sync.Mutex
	certificates map[string]map[string]any
	// enabledStatus, when set, is the status an enable request leaves behind
	// instead of active (for example enrolling for a slow activation).
	enabledStatus string
}

func init() {
	mockState.certificates = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("network_managed_tls_certificate", &ManagedTLSCertificateMock{})
}

type ManagedTLSCertificateMock struct{}

var _ mocks.MockRegistrar = (*ManagedTLSCertificateMock)(nil)

func (m *ManagedTLSCertificateMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates", m.createResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates/([^/]+)$`, m.getResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates/([^/]+)$`, m.updateResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates/([^/]+)$`, m.deleteResponder())
}

func (m *ManagedTLSCertificateMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/networkaccess/tls/managedCertificateAuthorityCertificates",
		func(req *http.Request) (*http.Response, error) {
			return jsonFixtureResponse(req, 400, filepath.Join("..", "tests", "responses", "validate_create", "post_managed_tls_certificate_error.json"))
		},
	)
}

func (m *ManagedTLSCertificateMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.certificates = make(map[string]map[string]any)
	mockState.enabledStatus = ""
}

// SetEnabledStatus makes enabling leave the certificate in status, such as
// enrolling (slow activation) or expired (terminal). Empty restores active.
func (m *ManagedTLSCertificateMock) SetEnabledStatus(status string) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.enabledStatus = status
}

// maskEnums hides evolvable enum members unless the client opts in, as Graph does.
func maskEnums(req *http.Request, response map[string]any) {
	if response["status"] == "disabled" && !strings.Contains(req.Header.Get("Prefer"), "include-unknown-enum-members") {
		response["status"] = "unknownFutureValue"
	}
}

func (m *ManagedTLSCertificateMock) createResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}
		if _, exists := requestBody["status"]; exists {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"create payload must match the portal request and omit status"}}`), nil
		}

		response, err := loadJSONFixture(filepath.Join("..", "tests", "responses", "validate_create", "post_managed_tls_certificate.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		for _, key := range []string{"name", "commonName", "organizationName", "validityMonths"} {
			if value, ok := requestBody[key]; ok {
				response[key] = value
			}
		}

		mockState.Lock()
		mockState.certificates[managedTLSCertificateID] = response
		mockState.Unlock()
		created := cloneMap(response)
		maskEnums(req, created)
		return factories.SuccessResponse(201, created)(req)
	}
}

func (m *ManagedTLSCertificateMock) getResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/tls/managedCertificateAuthorityCertificates/")
		mockState.Lock()
		certificate, exists := mockState.certificates[id]
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_managed_tls_certificate_not_found.json"))
		}
		response := cloneMap(certificate)
		maskEnums(req, response)
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *ManagedTLSCertificateMock) updateResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/tls/managedCertificateAuthorityCertificates/")
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()
		certificate, exists := mockState.certificates[id]
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_managed_tls_certificate_not_found.json"))
		}
		if status, ok := requestBody["status"]; ok {
			switch {
			case status == "enabled" && mockState.enabledStatus != "":
				certificate["status"] = mockState.enabledStatus
			case status == "enabled":
				certificate["status"] = "active"
			default:
				certificate["status"] = status
			}
		}
		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *ManagedTLSCertificateMock) deleteResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/tls/managedCertificateAuthorityCertificates/")
		mockState.Lock()
		_, exists := mockState.certificates[id]
		if exists {
			delete(mockState.certificates, id)
		}
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_managed_tls_certificate_not_found.json"))
		}
		return factories.EmptySuccessResponse(204)(req)
	}
}

func loadJSONFixture(path string) (map[string]any, error) {
	jsonContent, err := helpers.ParseJSONFile(path)
	if err != nil {
		return nil, err
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, err
	}
	return response, nil
}

func jsonFixtureResponse(req *http.Request, status int, path string) (*http.Response, error) {
	response, err := loadJSONFixture(path)
	if err != nil {
		return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
	}
	return factories.SuccessResponse(status, response)(req)
}

func cloneMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
