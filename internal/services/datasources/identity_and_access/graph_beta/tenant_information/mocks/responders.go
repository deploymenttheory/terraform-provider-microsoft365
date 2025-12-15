package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	tenantInfo map[string]map[string]any
}

func init() {
	mockState.tenantInfo = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("tenant_information", &TenantInformationMock{})
}

type TenantInformationMock struct{}

var _ mocks.MockRegistrar = (*TenantInformationMock)(nil)

func (m *TenantInformationMock) RegisterMocks() {
	mockState.Lock()
	mockState.tenantInfo = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Find tenant information by tenant ID
	// GET /tenantRelationships/findTenantInformationByTenantId(tenantId='{id}')
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/tenantRelationships/findTenantInformationByTenantId\(tenantId='[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'\)$`, func(req *http.Request) (*http.Response, error) {
		// Extract tenant ID from URL
		urlPath := req.URL.Path
		if strings.Contains(urlPath, "6babcaad-604b-40ac-a9d7-9fd97c0b779f") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_tenant_info_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Tenant not found"}}`), nil
	})

	// 2. Find tenant information by domain name
	// GET /tenantRelationships/findTenantInformationByDomainName(domainName='{name}')
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/tenantRelationships/findTenantInformationByDomainName\(domainName='[^']+'\)$`, func(req *http.Request) (*http.Response, error) {
		// Extract domain name from URL
		urlPath := req.URL.Path
		if strings.Contains(urlPath, "deploymenttheory.com") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_tenant_info_by_domain.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Tenant not found for domain"}}`), nil
	})
}

func (m *TenantInformationMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.tenantInfo = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/tenantRelationships/findTenantInformationByTenantId\(tenantId='[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'\)$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/tenantRelationships/findTenantInformationByDomainName\(domainName='[^']+'\)$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})
}

func (m *TenantInformationMock) CleanupMockState() {
	mockState.Lock()
	mockState.tenantInfo = make(map[string]map[string]any)
	mockState.Unlock()
}
