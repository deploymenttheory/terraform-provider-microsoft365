package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

const URL = "https://graph.microsoft.com/beta/networkAccess/settings/crossTenantAccess"

type NetworkCrossTenantAccessSettingsMock struct {
	mu      sync.Mutex
	status  string
	patches int
}

var _ mocks.MockRegistrar = (*NetworkCrossTenantAccessSettingsMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_cross_tenant_access_settings", &NetworkCrossTenantAccessSettingsMock{})
}

func (m *NetworkCrossTenantAccessSettingsMock) RegisterMocks() {
	m.CleanupMockState()
	httpmock.RegisterResponder("GET", URL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		response, err := httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context":             "https://graph.microsoft.com/beta/$metadata#networkAccess/settings/crossTenantAccess/$entity",
			"networkPacketTaggingStatus": m.status,
			"dataPlaneTaggingOptions":    "none",
		})
		if err != nil {
			return nil, fmt.Errorf("encode cross-tenant settings mock response: %w", err)
		}
		return response, nil
	})
	httpmock.RegisterResponder("PATCH", URL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode cross-tenant settings mock request: %w", err)
		}
		for key := range body {
			if key != "networkPacketTaggingStatus" && key != "@odata.type" {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Unexpected property"}}`), nil
			}
		}
		status, _ := body["networkPacketTaggingStatus"].(string)
		if status != "enabled" && status != "disabled" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid status"}}`), nil
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		m.status = status
		m.patches++
		return httpmock.NewStringResponse(204, ""), nil
	})
	for _, method := range []string{"POST", "DELETE"} {
		httpmock.RegisterResponder(method, URL, func(req *http.Request) (*http.Response, error) {
			return nil, io.ErrUnexpectedEOF
		})
	}
}

func (m *NetworkCrossTenantAccessSettingsMock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "PATCH"} {
		httpmock.RegisterResponder(method, URL, httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges"}}`))
	}
}

func (m *NetworkCrossTenantAccessSettingsMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = "disabled"
	m.patches = 0
}

func (m *NetworkCrossTenantAccessSettingsMock) Snapshot() (string, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status, m.patches
}

func (m *NetworkCrossTenantAccessSettingsMock) SetStatus(status string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = status
}
