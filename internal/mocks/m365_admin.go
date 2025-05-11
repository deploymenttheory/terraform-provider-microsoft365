package mocks

import (
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
)

// ActivateM365AdminMocks sets up mocks for M365 Admin Center endpoints.
func ActivateM365AdminMocks() {
	// Browser Site Lists
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/admin/edge/browserSiteLists\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/edge/browserSiteLists",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Site List",
						"description": "Test Description",
						"lastModifiedBy": "admin@contoso.com"
					}
				]
			}`), nil
		})

	// Browser Sites
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/admin/edge/browserSiteLists/[^/]+/sites\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/edge/browserSiteLists/{id}/sites",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"url": "https://example.com"
					}
				]
			}`), nil
		})
}
