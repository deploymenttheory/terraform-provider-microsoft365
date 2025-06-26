package client

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// MockAuthProvider implements the required authentication interface for testing
type MockAuthProvider struct{}

// AuthenticateRequest adds a mock authorization header to requests
func (m *MockAuthProvider) AuthenticateRequest(ctx context.Context, request *abstractions.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	if request.Headers != nil {
		request.Headers.Add("Authorization", "Bearer mock-token")
	}
	return nil
}

// MockCredential implements azcore.TokenCredential for testing
type MockCredential struct{}

// GetToken returns a mock access token
func (m *MockCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "mock-access-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

// NewMockGraphClients creates a new instance of MockGraphClients with initialized mock clients
func NewMockGraphClients(httpClient *http.Client) *MockGraphClients {
	// Create mock auth provider
	mockAuthProvider := &MockAuthProvider{}

	// Create mock credential
	mockCredential := &MockCredential{}

	// Create mock adapters
	mockV1Adapter, _ := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		mockAuthProvider,
		nil,
		nil,
		httpClient,
	)

	mockBetaAdapter, _ := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		mockAuthProvider,
		nil,
		nil,
		httpClient,
	)

	// Set base URLs for the mock adapters
	mockV1Adapter.SetBaseUrl("https://graph.microsoft.com/v1.0")
	mockBetaAdapter.SetBaseUrl("https://graph.microsoft.com/beta")

	// Create mock HTTP clients
	mockV1HTTPClient := NewAuthenticatedHTTPClient(httpClient, mockCredential, "https://graph.microsoft.com/.default", "https://graph.microsoft.com/v1.0")
	mockBetaHTTPClient := NewAuthenticatedHTTPClient(httpClient, mockCredential, "https://graph.microsoft.com/.default", "https://graph.microsoft.com/beta")

	// Create and return the mock clients
	return &MockGraphClients{
		MockV1Client:       msgraphsdk.NewGraphServiceClient(mockV1Adapter),
		MockBetaClient:     msgraphbetasdk.NewGraphServiceClient(mockBetaAdapter),
		MockV1HTTPClient:   mockV1HTTPClient,
		MockBetaHTTPClient: mockBetaHTTPClient,
	}
}
