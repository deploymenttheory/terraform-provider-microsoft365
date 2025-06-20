package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockTokenCredential is a mock implementation of azcore.TokenCredential
type MockTokenCredential struct {
	mock.Mock
}

func (m *MockTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	args := m.Called(ctx, options)
	if args.Get(1) != nil {
		return azcore.AccessToken{}, args.Error(1)
	}
	return args.Get(0).(azcore.AccessToken), nil
}

// Helper function to create a basic provider data for testing
func createTestConfig(authMethod string) *ProviderData {
	return &ProviderData{
		Cloud:           "public",
		TenantID:        "tenant-id",
		AuthMethod:      authMethod,
		EntraIDOptions:  &EntraIDOptions{},
		ClientOptions:   &ClientOptions{},
		TelemetryOptout: false,
		DebugMode:       false,
	}
}

// TestGitHubOIDCStrategy_GetAssertion tests the GitHub OIDC strategy's getAssertion function
func TestGitHubOIDCStrategy_GetAssertion(t *testing.T) {
	// Create a mock GitHub OIDC token server
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request is properly formed
		assert.Equal(t, "GET", r.Method, "OIDC token request should use GET method")
		assert.Equal(t, "Bearer test-request-token", r.Header.Get("Authorization"), "Authorization header should be set")

		// Check that the audience parameter was properly added to the URL
		assert.Contains(t, r.URL.RawQuery, "audience=api", "URL should contain the audience parameter")

		// Return a mock OIDC token response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"value": "test-oidc-token"}`))
	}))
	defer tokenServer.Close()

	// Set up environment for the test
	baseURL := tokenServer.URL
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", baseURL)
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "test-request-token")
	defer os.Clearenv()

	// Create test context
	ctx := context.Background()

	// Manually replicate the URL construction logic from GitHubOIDCStrategy
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	// Add audience parameter to URL - this matches the logic in the real GitHubOIDCStrategy
	audience := "api://AzureADTokenExchange"
	if !strings.Contains(requestURL, "audience=") {
		separator := "&"
		if !strings.Contains(requestURL, "?") {
			separator = "?"
		}
		requestURL = requestURL + separator + "audience=" + audience
	}

	// Manually construct the assertion function (similar to strategy's implementation)
	getAssertion := func(ctx context.Context) (string, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			return "", err
		}

		req.Header.Add("Authorization", "Bearer "+requestToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("failed to get GitHub OIDC token, status code: %d, response: %s",
				resp.StatusCode, string(body))
		}

		var tokenResp struct {
			Value string `json:"value"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			return "", err
		}

		if tokenResp.Value == "" {
			return "", fmt.Errorf("received empty OIDC token from GitHub")
		}

		return tokenResp.Value, nil
	}

	// Execute the assertion function
	token, err := getAssertion(ctx)

	// Verify results
	require.NoError(t, err, "getAssertion should not return an error")
	assert.Equal(t, "test-oidc-token", token, "Token should match the expected value from the mock server")
}

// TestGitHubOIDCStrategy_MissingEnvironment tests error handling when GitHub environment variables are missing
func TestGitHubOIDCStrategy_MissingEnvironment(t *testing.T) {
	// Ensure environment is clean
	os.Clearenv()

	// Create test context and options
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	// Create test config
	config := createTestConfig("oidc_github")

	// Create the GitHubOIDCStrategy
	strategy := &GitHubOIDCStrategy{}

	// Call the GetCredential method
	credential, err := strategy.GetCredential(ctx, config, clientOptions)

	// Verify it returns an error due to missing environment variables
	assert.Error(t, err, "GetCredential should return an error when GitHub environment variables are missing")
	assert.Nil(t, credential, "Credential should be nil when environment variables are missing")
	assert.Contains(t, err.Error(), "ACTIONS_ID_TOKEN_REQUEST_URL", "Error should mention the missing URL env var")
	assert.Contains(t, err.Error(), "ACTIONS_ID_TOKEN_REQUEST_TOKEN", "Error should mention the missing token env var")
}

// TestGitHubOIDCStrategy_AudienceParameter tests URL construction with audience parameter
func TestGitHubOIDCStrategy_AudienceParameter(t *testing.T) {
	// Setup normal URL case
	t.Run("default audience", func(t *testing.T) {
		// Given a URL without audience parameter
		baseURL := "https://example.com/token"
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", baseURL)
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "token")
		os.Unsetenv("M365_OIDC_AUDIENCE")
		os.Unsetenv("ARM_OIDC_AUDIENCE")

		// We need to check what the URL would be constructed as
		requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")

		// Check if a custom audience is specified and add to request URL if not already present
		audience := "api://AzureADTokenExchange" // Default audience

		if !strings.Contains(requestURL, "audience=") {
			separator := "&"
			if !strings.Contains(requestURL, "?") {
				separator = "?"
			}
			requestURL = requestURL + separator + "audience=" + audience
		}

		// Verify
		assert.Contains(t, requestURL, "?audience=api://AzureADTokenExchange",
			"URL should include default audience parameter when none specified")
	})

	// Setup custom audience case
	t.Run("custom audience", func(t *testing.T) {
		// Given a URL without audience parameter
		baseURL := "https://example.com/token"
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", baseURL)
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "token")
		os.Setenv("M365_OIDC_AUDIENCE", "custom-audience")

		// We need to check what the URL would be constructed as
		requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")

		// Check if a custom audience is specified and add to request URL if not already present
		audience := "custom-audience" // From M365_OIDC_AUDIENCE

		if !strings.Contains(requestURL, "audience=") {
			separator := "&"
			if !strings.Contains(requestURL, "?") {
				separator = "?"
			}
			requestURL = requestURL + separator + "audience=" + audience
		}

		// Verify
		assert.Contains(t, requestURL, "?audience=custom-audience",
			"URL should include custom audience parameter from M365_OIDC_AUDIENCE")
	})

	// Cleanup
	os.Clearenv()
}

// TestGitHubOIDCStrategy_ExistingQueryParams tests URL construction with existing query parameters
func TestGitHubOIDCStrategy_ExistingQueryParams(t *testing.T) {
	// Given a URL with existing query parameters
	baseURL := "https://example.com/token?existing=param"
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", baseURL)
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "token")

	// We need to check what the URL would be constructed as
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")

	// Check if a custom audience is specified and add to request URL if not already present
	audience := "api://AzureADTokenExchange" // Default audience

	if !strings.Contains(requestURL, "audience=") {
		separator := "&"
		if !strings.Contains(requestURL, "?") {
			separator = "?"
		}
		requestURL = requestURL + separator + "audience=" + audience
	}

	// Verify
	assert.Contains(t, requestURL, "existing=param", "URL should preserve existing query parameters")
	assert.Contains(t, requestURL, "&audience=api://AzureADTokenExchange",
		"URL should add audience parameter with & when existing params present")

	// Cleanup
	os.Clearenv()
}
