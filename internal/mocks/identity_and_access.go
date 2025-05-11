package mocks

import (
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
)

// ActivateIdentityAndAccessMocks sets up mocks for Identity and Access Management endpoints.
func ActivateIdentityAndAccessMocks() {
	// Conditional Access Policies
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/policies",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test CA Policy",
						"state": "enabled",
						"createdDateTime": "2023-01-01T12:00:00Z"
					}
				]
			}`), nil
		})
}
