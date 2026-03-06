package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTokenCredential implements azcore.TokenCredential for testing
type mockTokenCredential struct {
	token string
	err   error
}

func (m *mockTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if m.err != nil {
		return azcore.AccessToken{}, m.err
	}
	return azcore.AccessToken{Token: m.token}, nil
}

// TestUnit_NewAuthenticatedHTTPClient_Initialization validates that NewAuthenticatedHTTPClient creates
// a properly configured client with all required fields.
func TestUnit_NewAuthenticatedHTTPClient_Initialization(t *testing.T) {
	baseClient := &http.Client{}
	cred := &mockTokenCredential{token: "test-token"}
	scope := "https://graph.microsoft.com/.default"
	baseURL := "https://graph.microsoft.com/v1.0"

	client := NewAuthenticatedHTTPClient(baseClient, cred, scope, baseURL)

	require.NotNil(t, client)
	assert.Equal(t, baseClient, client.GetClient())
	assert.Equal(t, baseURL, client.GetBaseURL())
	assert.Equal(t, scope, client.scope)
	assert.Equal(t, cred, client.credential)
}
