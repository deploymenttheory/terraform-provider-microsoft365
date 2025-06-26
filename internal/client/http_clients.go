package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

// Do performs an HTTP request with authentication
func (c *AuthenticatedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Get token from credential
	token, err := c.credential.GetToken(req.Context(), policy.TokenRequestOptions{
		Scopes: []string{c.scope},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+token.Token)

	// Set standard Microsoft Graph API headers
	req.Header.Set("Accept", "application/json")

	// Set default headers for Graph API
	if req.Header.Get("Content-Type") == "" && (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add consistency level header for certain operations
	if req.Method == "GET" {
		req.Header.Set("ConsistencyLevel", "eventual")
	}

	// Perform the request
	return c.client.Do(req)
}

// Get performs a GET request
func (c *AuthenticatedHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post performs a POST request with JSON body
func (c *AuthenticatedHTTPClient) Post(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// GetBaseURL returns the base URL for this client
func (c *AuthenticatedHTTPClient) GetBaseURL() string {
	return c.baseURL
}

// GetClient returns the underlying HTTP client
func (c *AuthenticatedHTTPClient) GetClient() *http.Client {
	return c.client
}
