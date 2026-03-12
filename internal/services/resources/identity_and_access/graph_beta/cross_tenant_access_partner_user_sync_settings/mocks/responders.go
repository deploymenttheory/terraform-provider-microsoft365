package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

const baseURL = "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/partners"

// mockState holds the in-memory state for partner user sync settings.
var mockState struct {
	sync.Mutex
	settings map[string]map[string]any // keyed by tenant ID
}

func init() {
	mockState.settings = make(map[string]map[string]any)

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	mocks.GlobalRegistry.Register("cross_tenant_access_partner_user_sync_settings", &CrossTenantAccessPartnerUserSyncSettingsMock{})
}

// CrossTenantAccessPartnerUserSyncSettingsMock provides mock responses for partner user sync settings operations.
type CrossTenantAccessPartnerUserSyncSettingsMock struct{}

// Ensure CrossTenantAccessPartnerUserSyncSettingsMock implements MockRegistrar interface.
var _ mocks.MockRegistrar = (*CrossTenantAccessPartnerUserSyncSettingsMock)(nil)

// RegisterMocks registers HTTP mock responses for partner user sync settings operations.
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) RegisterMocks() {
	mockState.Lock()
	mockState.settings = make(map[string]map[string]any)
	mockState.Unlock()

	// PUT /beta/policies/crossTenantAccessPolicy/partners/{tenantId}/identitySynchronization - Create
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		m.putSettingsResponder())

	// GET /beta/policies/crossTenantAccessPolicy/partners/{tenantId}/identitySynchronization - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		m.getSettingsResponder())

	// DELETE /beta/policies/crossTenantAccessPolicy/partners/{tenantId}/identitySynchronization - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		m.deleteSettingsResponder())
}

// loadJSONResponse loads a JSON response from the tests/responses directory
func loadJSONResponse(filename string) (map[string]any, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	responsePath := filepath.Join(filepath.Dir(currentFile), "..", "tests", "responses", "validate_create", filename)

	data, err := os.ReadFile(responsePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read response file %s: %w", filename, err)
	}

	var response map[string]any
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// putSettingsResponder handles PUT requests to create/update partner user sync settings
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) putSettingsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Extract tenant ID from URL
		tenantID := "12345678-1234-1234-1234-123456789012" // Default for tests

		// Load response template
		response, err := loadJSONResponse("put_partner_user_sync_settings_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
		}

		// Merge request body into response
		for key, value := range requestBody {
			response[key] = value
		}
		response["tenantId"] = tenantID

		// Store in mock state
		mockState.Lock()
		mockState.settings[tenantID] = response
		mockState.Unlock()

		resp, err := httpmock.NewJsonResponse(200, response)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	}
}

// getSettingsResponder handles GET requests to read partner user sync settings
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) getSettingsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract tenant ID from URL
		tenantID := "12345678-1234-1234-1234-123456789012" // Default for tests

		mockState.Lock()
		settings, exists := mockState.settings[tenantID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Identity synchronization not found"}}`), nil
		}

		resp, err := httpmock.NewJsonResponse(200, settings)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	}
}

// deleteSettingsResponder handles DELETE requests to remove partner user sync settings
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) deleteSettingsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract tenant ID from URL
		tenantID := "12345678-1234-1234-1234-123456789012" // Default for tests

		mockState.Lock()
		delete(mockState.settings, tenantID)
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions.
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		factories.ErrorResponse(403, "Forbidden", "Insufficient privileges to configure identity synchronization"))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied to identity synchronization"))

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)/identitySynchronization$`,
		factories.ErrorResponse(403, "Forbidden", "Insufficient privileges to delete identity synchronization"))
}

// CleanupMockState resets the mock state.
func (m *CrossTenantAccessPartnerUserSyncSettingsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.settings = make(map[string]map[string]any)
}
