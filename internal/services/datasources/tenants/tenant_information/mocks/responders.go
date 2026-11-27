package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

type TenantInformationMock struct{}

func (m *TenantInformationMock) RegisterMocks() {
	body := "{\"tenantId\":\"6babcaad-604b-40ac-a9d7-9fd97c0b779f\",\"displayName\":\"Deployment Theory\",\"defaultDomainName\":\"deploymenttheory.com\",\"federationBrandName\":null}"
	for _, url := range []string{
		"https://graph.microsoft.com/beta/tenantRelationships/findTenantInformationByTenantId(tenantId='6babcaad-604b-40ac-a9d7-9fd97c0b779f')",
		"https://graph.microsoft.com/beta/tenantRelationships/findTenantInformationByDomainName(domainName='deploymenttheory.com')",
	} {
		httpmock.RegisterResponder(http.MethodGet, url, httpmock.NewStringResponder(http.StatusOK, body).HeaderSet(http.Header{"Content-Type": []string{"application/json"}}))
	}
}

func (m *TenantInformationMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		http.MethodGet,
		"=~^https://graph\\.microsoft\\.com/beta/tenantRelationships/.*",
		httpmock.NewStringResponder(
			http.StatusForbidden,
			"{\"error\":{\"code\":\"Authorization_RequestDenied\",\"message\":\"Insufficient privileges to complete the operation.\"}}",
		),
	)
}
