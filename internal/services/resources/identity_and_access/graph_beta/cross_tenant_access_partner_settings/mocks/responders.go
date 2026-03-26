package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	partnerSettings map[string]map[string]any
	deletedTenants  map[string]bool
}

func init() {
	mockState.partnerSettings = make(map[string]map[string]any)
	mockState.deletedTenants = make(map[string]bool)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("cross_tenant_access_partner_settings", &CrossTenantAccessPartnerSettingsMock{})
}

type CrossTenantAccessPartnerSettingsMock struct{}

var _ mocks.MockRegistrar = (*CrossTenantAccessPartnerSettingsMock)(nil)

func loadJSONResponse(filename string) (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_create", filename)
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file %s: %w", filename, err)
	}

	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return responseObj, nil
}

func mergeSettings(target, source map[string]any) {
	for key, value := range source {
		if sourceMap, ok := value.(map[string]any); ok {
			if targetMap, ok := target[key].(map[string]any); ok {
				mergeSettings(targetMap, sourceMap)
			} else {
				target[key] = value
			}
		} else {
			target[key] = value
		}
	}
}

func (m *CrossTenantAccessPartnerSettingsMock) RegisterMocks() {
	mockState.Lock()
	mockState.partnerSettings = make(map[string]map[string]any)
	mockState.deletedTenants = make(map[string]bool)
	mockState.Unlock()

	m.registerTenantValidationMock()
	m.registerMockUsers()
	m.registerMockGroups()

	// POST /policies/crossTenantAccessPolicy/partners - Create partner configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/partners", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Invalid request body: %s"}}`, err.Error())), nil
		}

		tenantID, ok := requestBody["tenantId"].(string)
		if !ok || tenantID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"tenantId is required"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		// Load base response from JSON file
		baseResponse, err := loadJSONResponse("post_partner_settings_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
		}

		// Merge request body into base response
		mergeSettings(baseResponse, requestBody)
		baseResponse["tenantId"] = tenantID

		mockState.partnerSettings[tenantID] = baseResponse

		resp, err := httpmock.NewJsonResponse(201, baseResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})

	// GET /policies/crossTenantAccessPolicy/partners/{tenantId} - Read partner configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		tenantID := httpmock.MustGetSubmatch(req, 1)

		mockState.Lock()
		defer mockState.Unlock()

		// Return 404 for explicitly deleted tenants (simulates DELETE propagation)
		if mockState.deletedTenants[tenantID] {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Partner configuration not found"}}`), nil
		}

		if settings, exists := mockState.partnerSettings[tenantID]; exists {
			resp, err := httpmock.NewJsonResponse(200, settings)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
			}
			return resp, nil
		}

		// If not in state, load from JSON file
		baseResponse, err := loadJSONResponse("get_partner_settings_success.json")
		if err != nil {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Partner configuration not found"}}`), nil
		}

		baseResponse["tenantId"] = tenantID
		mockState.partnerSettings[tenantID] = baseResponse

		resp, err := httpmock.NewJsonResponse(200, baseResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})

	// PATCH /policies/crossTenantAccessPolicy/partners/{tenantId} - Update partner configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		tenantID := httpmock.MustGetSubmatch(req, 1)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Invalid request body: %s"}}`, err.Error())), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		if settings, exists := mockState.partnerSettings[tenantID]; exists {
			mergeSettings(settings, requestBody)
			return httpmock.NewStringResponse(204, ""), nil
		}

		// If not in state, load base and merge
		baseResponse, err := loadJSONResponse("get_partner_settings_success.json")
		if err != nil {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Partner configuration not found"}}`), nil
		}

		baseResponse["tenantId"] = tenantID
		mergeSettings(baseResponse, requestBody)
		mockState.partnerSettings[tenantID] = baseResponse

		return httpmock.NewStringResponse(204, ""), nil
	})

	// DELETE /policies/crossTenantAccessPolicy/partners/{tenantId} - Delete partner configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		tenantID := httpmock.MustGetSubmatch(req, 1)

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.partnerSettings[tenantID]; exists {
			delete(mockState.partnerSettings, tenantID)
			mockState.deletedTenants[tenantID] = true
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Partner configuration not found"}}`), nil
	})

	// DELETE /directory/deletedItems/{tenantId} - Hard delete partner configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([^/]+)$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *CrossTenantAccessPartnerSettingsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.partnerSettings = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerTenantValidationMock()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/partners",
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`,
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`,
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/policies/crossTenantAccessPolicy/partners/([^/]+)$`,
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Mock error for testing"}}`))
}

func (m *CrossTenantAccessPartnerSettingsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.partnerSettings = make(map[string]map[string]any)
	mockState.deletedTenants = make(map[string]bool)
	httpmock.Reset()
}

// registerTenantValidationMock registers the mock for tenant validation API call
func (m *CrossTenantAccessPartnerSettingsMock) registerTenantValidationMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/tenantRelationships/findTenantInformationByTenantId\(tenantId='([^']+)'\)$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context":      "https://graph.microsoft.com/beta/$metadata#microsoft.graph.tenantInformation",
				"tenantId":            "12345678-1234-1234-1234-123456789012",
				"displayName":         "Partner Organization",
				"defaultDomainName":   "partner.onmicrosoft.com",
				"federationBrandName": "Partner Org",
			}
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
			}
			return resp, nil
		})
}

// registerMockUsers registers mock user resources for validation
func (m *CrossTenantAccessPartnerSettingsMock) registerMockUsers() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users\?`, func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#users(id)",
			"value": []any{
				map[string]any{
					"id": "11111111-1111-1111-1111-111111111111",
				},
				map[string]any{
					"id": "22222222-2222-2222-2222-222222222222",
				},
			},
		}

		resp, err := httpmock.NewJsonResponse(200, response)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})
}

// registerMockGroups registers mock group resources for validation
func (m *CrossTenantAccessPartnerSettingsMock) registerMockGroups() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups\?`, func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups(id)",
			"value": []any{
				map[string]any{
					"id": "33333333-3333-3333-3333-333333333333",
				},
				map[string]any{
					"id": "44444444-4444-4444-4444-444444444444",
				},
			},
		}

		resp, err := httpmock.NewJsonResponse(200, response)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to create JSON response: %s"}}`, err.Error())), nil
		}
		return resp, nil
	})
}
