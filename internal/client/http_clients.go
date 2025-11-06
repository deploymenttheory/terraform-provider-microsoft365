package client

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AuthenticatedHTTPClient wraps an HTTP client with Microsoft Graph authentication
type AuthenticatedHTTPClient struct {
	client     *http.Client
	credential azcore.TokenCredential
	scope      string
	baseURL    string
}

// NewAuthenticatedHTTPClient creates a new HTTP client with authentication for Microsoft Graph
func NewAuthenticatedHTTPClient(baseClient *http.Client, credential azcore.TokenCredential, scope, baseURL string) *AuthenticatedHTTPClient {
	return &AuthenticatedHTTPClient{
		client:     baseClient,
		credential: credential,
		scope:      scope,
		baseURL:    baseURL,
	}
}
