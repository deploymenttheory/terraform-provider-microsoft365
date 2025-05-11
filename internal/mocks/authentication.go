package mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

// activateAuthenticationMocks sets up mocks for authentication and common Graph API endpoints.
func activateAuthenticationMocks() {
	// Mock for Graph API auth token
	httpmock.RegisterResponder("POST", `=~^https://login\.microsoftonline\.com/[^/]+/oauth2/v2\.0/token\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"token_type": "Bearer",
				"expires_in": 3599,
				"access_token": "mock-access-token"
			}`), nil
		})

	// Mock for Graph API version information
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/\$metadata\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/v1.0/$metadata",
				"value": [
					{
						"name": "microsoft.graph",
						"url": "https://graph.microsoft.com/v1.0"
					}
				]
			}`), nil
		})

	// Mock for Graph Beta API version information
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/\$metadata\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata",
				"value": [
					{
						"name": "microsoft.graph",
						"url": "https://graph.microsoft.com/beta"
					}
				]
			}`), nil
		})
}
