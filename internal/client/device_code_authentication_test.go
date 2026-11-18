package client

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitDeviceCodeStrategyHeadlessChallengeUsesConfiguredTransport(t *testing.T) {
	const authority = "https://login.microsoftonline.com/11111111-1111-4111-8111-111111111111"
	transport := httpmock.NewMockTransport()
	transport.RegisterResponder("GET", "=~^https://login\\.microsoftonline\\.com/.*/\\.well-known/openid-configuration", httpmock.NewJsonResponderOrPanic(200, map[string]any{
		"authorization_endpoint":        authority + "/oauth2/v2.0/authorize",
		"token_endpoint":                authority + "/oauth2/v2.0/token",
		"device_authorization_endpoint": authority + "/oauth2/v2.0/devicecode",
		"issuer":                        authority + "/v2.0",
		"jwks_uri":                      authority + "/discovery/v2.0/keys",
	}))
	transport.RegisterResponder("POST", authority+"/oauth2/v2.0/devicecode", httpmock.NewJsonResponderOrPanic(200, map[string]any{
		"device_code":      "unit-device-code-not-a-token",
		"user_code":        "UNIT-CODE-ONLY",
		"verification_uri": "https://microsoft.com/devicelogin",
		"expires_in":       900,
		"interval":         1,
		"message":          "Enter UNIT-CODE-ONLY at https://microsoft.com/devicelogin",
	}))
	opened := false
	strategy := &DeviceCodeStrategy{openTerminal: func() (io.WriteCloser, error) { opened = true; return nil, io.ErrClosedPipe }}
	config := &ProviderData{
		TenantID:       "11111111-1111-4111-8111-111111111111",
		EntraIDOptions: &EntraIDOptions{ClientID: "22222222-2222-4222-8222-222222222222", DisableInstanceDiscovery: true},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	credential, err := strategy.GetCredential(ctx, config, policy.ClientOptions{Transport: &http.Client{Transport: transport}, Retry: policy.RetryOptions{MaxRetries: -1}})
	require.NoError(t, err)
	_, err = credential.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{"https://graph.microsoft.com/.default"}})
	require.Error(t, err)
	require.True(t, opened, "the Azure SDK must invoke the configured terminal prompt after the mocked device-code response")
	require.ErrorContains(t, err, "interactive terminal")
	require.NotContains(t, err.Error(), "UNIT-CODE-ONLY")
}
