package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnit_MockAuthProvider_AuthenticateRequest validates the mock auth provider.
func TestUnit_MockAuthProvider_AuthenticateRequest(t *testing.T) {
	ctx := context.Background()
	provider := &MockAuthProvider{}

	req := abstractions.NewRequestInformation()
	req.Headers = abstractions.NewRequestHeaders()

	err := provider.AuthenticateRequest(ctx, req, nil)
	require.NoError(t, err)

	authHeaders := req.Headers.Get("Authorization")
	assert.Contains(t, authHeaders, "Bearer mock-token")
}

// TestUnit_MockCredential_GetToken validates the mock credential.
func TestUnit_MockCredential_GetToken(t *testing.T) {
	ctx := context.Background()
	cred := &MockCredential{}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})

	require.NoError(t, err)
	assert.Equal(t, "mock-access-token", token.Token)
	assert.True(t, token.ExpiresOn.After(time.Now()), "Token should not be expired")
}

// TestUnit_NewMockGraphClients_Initialization validates mock client creation.
func TestUnit_NewMockGraphClients_Initialization(t *testing.T) {
	httpClient := &http.Client{}

	mockClients := NewMockGraphClients(httpClient)

	require.NotNil(t, mockClients)
	assert.NotNil(t, mockClients.MockV1Client, "V1 client should be initialized")
	assert.NotNil(t, mockClients.MockBetaClient, "Beta client should be initialized")
	assert.NotNil(t, mockClients.MockV1HTTPClient, "V1 HTTP client should be initialized")
	assert.NotNil(t, mockClients.MockBetaHTTPClient, "Beta HTTP client should be initialized")

	assert.Equal(t, "https://graph.microsoft.com/v1.0", mockClients.MockV1HTTPClient.GetBaseURL())
	assert.Equal(t, "https://graph.microsoft.com/beta", mockClients.MockBetaHTTPClient.GetBaseURL())
}

// TestUnit_MockGraphClients_Getters validates all getter methods for mock clients.
func TestUnit_MockGraphClients_Getters(t *testing.T) {
	httpClient := &http.Client{}
	mockClients := NewMockGraphClients(httpClient)

	t.Run("GetKiotaGraphV1Client", func(t *testing.T) {
		result := mockClients.GetKiotaGraphV1Client()
		assert.Equal(t, mockClients.MockV1Client, result)
	})

	t.Run("GetKiotaGraphBetaClient", func(t *testing.T) {
		result := mockClients.GetKiotaGraphBetaClient()
		assert.Equal(t, mockClients.MockBetaClient, result)
	})

	t.Run("GetGraphV1Client", func(t *testing.T) {
		result := mockClients.GetGraphV1Client()
		assert.Equal(t, mockClients.MockV1HTTPClient, result)
		assert.Equal(t, "https://graph.microsoft.com/v1.0", result.GetBaseURL())
	})

	t.Run("GetGraphBetaClient", func(t *testing.T) {
		result := mockClients.GetGraphBetaClient()
		assert.Equal(t, mockClients.MockBetaHTTPClient, result)
		assert.Equal(t, "https://graph.microsoft.com/beta", result.GetBaseURL())
	})
}
