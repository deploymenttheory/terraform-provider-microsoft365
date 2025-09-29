package mocks

import (
	"github.com/jarcoal/httpmock"
)

// AuthMock provides authentication mock responses
type AuthMock struct{}

// RegisterMocks registers authentication mock responses
func (a *AuthMock) RegisterMocks() {
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}))

	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"tenant_discovery_endpoint": "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/v2.0/.well-known/openid-configuration",
		}))
}

// RegisterErrorMocks registers authentication error responses
func (a *AuthMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(401, map[string]any{
			"error":             "invalid_client",
			"error_description": "Client authentication failed",
		}))
}

func init() {
	// Register the auth mock with the global registry
	GlobalRegistry.Register("auth", &AuthMock{})
}
