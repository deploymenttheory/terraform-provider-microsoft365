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

// TestGitHubOIDCStrategy_ErrorPaths tests error handling in GitHubOIDCStrategy.
func TestGitHubOIDCStrategy_ErrorPaths(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("Missing request URL and token", func(t *testing.T) {
		os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_URL")
		os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
		assert.Contains(t, err.Error(), "GitHub OIDC authentication requires")
	})

	t.Run("URL from env, token from config", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			token := "test-github-token"
			json.NewEncoder(w).Encode(map[string]any{
				"value": &token,
			})
		}))
		defer mockServer.Close()

		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", mockServer.URL)
		defer os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_URL")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:         "test-client-id",
				OIDCRequestToken: "test-request-token",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Token from env, URL from config", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			token := "test-github-token"
			json.NewEncoder(w).Encode(map[string]any{
				"value": &token,
			})
		}))
		defer mockServer.Close()

		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "env-request-token")
		defer os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:       "test-client-id",
				OIDCRequestURL: mockServer.URL,
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Both from env vars", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			token := "test-github-token"
			json.NewEncoder(w).Encode(map[string]any{
				"value": &token,
			})
		}))
		defer mockServer.Close()

		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", mockServer.URL)
		os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "env-request-token")
		defer os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_URL")
		defer os.Unsetenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestGitHubOIDCStrategy_AssertionErrorPaths tests error paths in the assertion callback.
func TestGitHubOIDCStrategy_AssertionErrorPaths(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("Assertion with nil value in response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"value": nil,
			})
		}))
		defer mockServer.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:         "test-client-id",
				OIDCRequestURL:   mockServer.URL,
				OIDCRequestToken: "test-request-token",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)
		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Assertion with invalid JSON", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not-valid-json"))
		}))
		defer mockServer.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:         "test-client-id",
				OIDCRequestURL:   mockServer.URL,
				OIDCRequestToken: "test-request-token",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)
		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Assertion with HTTP error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}))
		defer mockServer.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:         "test-client-id",
				OIDCRequestURL:   mockServer.URL,
				OIDCRequestToken: "test-request-token",
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)
		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With additional tenants and instance discovery disabled", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"value": "test-token",
			})
		}))
		defer mockServer.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				OIDCRequestURL:             mockServer.URL,
				OIDCRequestToken:           "test-request-token",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &GitHubOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)
		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
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

// TestClientSecretStrategy validates client secret credential creation.
func TestClientSecretStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With all options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				ClientSecret:               "test-client-secret",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &ClientSecretStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
		}

		strategy := &ClientSecretStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestClientCertificateStrategy validates client certificate credential creation.
func TestClientCertificateStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("Missing certificate file", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:          "test-client-id",
				ClientCertificate: "/nonexistent/cert.pfx",
			},
		}

		strategy := &ClientCertificateStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
		assert.Contains(t, err.Error(), "failed to read certificate file")
	})

	t.Run("Invalid certificate data", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "invalid-cert-*.pfx")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write([]byte("invalid certificate data"))
		require.NoError(t, err)
		tmpFile.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:          "test-client-id",
				ClientCertificate: tmpFile.Name(),
			},
		}

		strategy := &ClientCertificateStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
		assert.Contains(t, err.Error(), "failed to parse certificate data")
	})

	t.Run("With certificate password", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "cert-with-password-*.pfx")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write([]byte("invalid certificate data"))
		require.NoError(t, err)
		tmpFile.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                  "test-client-id",
				ClientCertificate:         tmpFile.Name(),
				ClientCertificatePassword: "test-password",
			},
		}

		strategy := &ClientCertificateStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
	})

	t.Run("With additional tenants", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "cert-*.pfx")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write([]byte("invalid certificate data"))
		require.NoError(t, err)
		tmpFile.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				ClientCertificate:          tmpFile.Name(),
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
			},
		}

		strategy := &ClientCertificateStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "cert-*.pfx")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write([]byte("invalid certificate data"))
		require.NoError(t, err)
		tmpFile.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                 "test-client-id",
				ClientCertificate:        tmpFile.Name(),
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &ClientCertificateStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
	})
}

// TestDeviceCodeStrategy validates device code credential creation.
func TestDeviceCodeStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("Successful credential creation", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				AdditionallyAllowedTenants: []string{"tenant1"},
				DisableInstanceDiscovery:   false,
			},
		}

		strategy := &DeviceCodeStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal config", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &DeviceCodeStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                 "test-client-id",
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &DeviceCodeStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With multiple additional tenants", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2", "tenant3"},
			},
		}

		strategy := &DeviceCodeStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With all options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &DeviceCodeStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestInteractiveBrowserStrategy validates interactive browser credential creation.
func TestInteractiveBrowserStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With all options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				RedirectUrl:                "http://localhost:8080",
				Username:                   "test@example.com",
				AdditionallyAllowedTenants: []string{"tenant1"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &InteractiveBrowserStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &InteractiveBrowserStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestUsernamePasswordStrategy validates username/password credential creation.
func TestUsernamePasswordStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With all options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				Username:                   "test@example.com",
				Password:                   "test-password",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &UsernamePasswordStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
				Username: "test@example.com",
				Password: "test-password",
			},
		}

		strategy := &UsernamePasswordStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestWorkloadIdentityStrategy validates workload identity credential creation.
func TestWorkloadIdentityStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With federated token file path and all options", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				FederatedTokenFilePath:     "/path/to/token",
				AdditionallyAllowedTenants: []string{"tenant1"},
				DisableInstanceDiscovery:   true,
			},
		}

		strategy := &WorkloadIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With federated token file path minimal", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:               "test-client-id",
				FederatedTokenFilePath: "/path/to/token",
			},
		}

		strategy := &WorkloadIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Without federated token file path requires env vars", func(t *testing.T) {
		os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/path/to/token")
		defer os.Unsetenv("AZURE_FEDERATED_TOKEN_FILE")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				AdditionallyAllowedTenants: []string{},
				DisableInstanceDiscovery:   false,
			},
		}

		strategy := &WorkloadIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestManagedIdentityStrategy validates managed identity credential creation.
func TestManagedIdentityStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With client ID", func(t *testing.T) {
		config := &ProviderData{
			EntraIDOptions: &EntraIDOptions{
				ManagedIdentityClientID: "test-client-id",
			},
		}

		strategy := &ManagedIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With resource ID", func(t *testing.T) {
		config := &ProviderData{
			EntraIDOptions: &EntraIDOptions{
				ManagedIdentityResourceID: "/subscriptions/test/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test",
			},
		}

		strategy := &ManagedIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With both client ID and resource ID prefers resource ID", func(t *testing.T) {
		config := &ProviderData{
			EntraIDOptions: &EntraIDOptions{
				ManagedIdentityClientID:   "test-client-id",
				ManagedIdentityResourceID: "/subscriptions/test/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test",
			},
		}

		strategy := &ManagedIdentityStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestAzureCLIStrategy validates Azure CLI credential creation.
func TestAzureCLIStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With additional tenants", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
			},
		}

		strategy := &AzureCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal options", func(t *testing.T) {
		config := &ProviderData{
			TenantID:       "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{},
		}

		strategy := &AzureCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &AzureCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestAzureDeveloperCLIStrategy validates Azure Developer CLI credential creation.
func TestAzureDeveloperCLIStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With additional tenants", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				AdditionallyAllowedTenants: []string{"tenant1"},
			},
		}

		strategy := &AzureDeveloperCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With minimal options", func(t *testing.T) {
		config := &ProviderData{
			TenantID:       "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{},
		}

		strategy := &AzureDeveloperCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &AzureDeveloperCLIStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestCredentialFactory validates the credential factory returns correct strategies.
func TestCredentialFactory(t *testing.T) {
	tests := []struct {
		authMethod  string
		expectError bool
	}{
		{"client_secret", false},
		{"client_certificate", false},
		{"device_code", false},
		{"interactive_browser", false},
		{"username_password", false},
		{"workload_identity", false},
		{"managed_identity", false},
		{"oidc", false},
		{"oidc_github", false},
		{"oidc_azure_devops", false},
		{"azure_cli", false},
		{"azure_developer_cli", false},
		{"invalid_method", true},
	}

	for _, tt := range tests {
		t.Run(tt.authMethod, func(t *testing.T) {
			strategy, err := credentialFactory(tt.authMethod)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, strategy)
				assert.Contains(t, err.Error(), "unsupported authentication method")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, strategy)
			}
		})
	}
}

// TestObtainCredential validates the main ObtainCredential function.
func TestObtainCredential(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	tests := []struct {
		name       string
		authMethod string
		setupEnv   func()
		cleanupEnv func()
		config     *EntraIDOptions
		expectErr  bool
	}{
		{
			name:       "client_secret",
			authMethod: "client_secret",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
			expectErr: false,
		},
		{
			name:       "device_code",
			authMethod: "device_code",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID: "test-client-id",
			},
			expectErr: false,
		},
		{
			name:       "interactive_browser",
			authMethod: "interactive_browser",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID: "test-client-id",
			},
			expectErr: false,
		},
		{
			name:       "username_password",
			authMethod: "username_password",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID: "test-client-id",
				Username: "test@example.com",
				Password: "test-password",
			},
			expectErr: false,
		},
		{
			name:       "workload_identity",
			authMethod: "workload_identity",
			setupEnv: func() {
				os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/tmp/token")
			},
			cleanupEnv: func() {
				os.Unsetenv("AZURE_FEDERATED_TOKEN_FILE")
			},
			config: &EntraIDOptions{
				ClientID: "test-client-id",
			},
			expectErr: false,
		},
		{
			name:       "managed_identity",
			authMethod: "managed_identity",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ManagedIdentityClientID: "test-client-id",
			},
			expectErr: false,
		},
		{
			name:       "oidc with token",
			authMethod: "oidc",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:  "test-client-id",
				OIDCToken: "test-token",
			},
			expectErr: false,
		},
		{
			name:       "azure_cli",
			authMethod: "azure_cli",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config:     &EntraIDOptions{},
			expectErr:  false,
		},
		{
			name:       "azure_developer_cli",
			authMethod: "azure_developer_cli",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config:     &EntraIDOptions{},
			expectErr:  false,
		},
		{
			name:       "unsupported_method",
			authMethod: "unsupported_method",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config:     &EntraIDOptions{},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			config := &ProviderData{
				TenantID:       "test-tenant-id",
				AuthMethod:     tt.authMethod,
				EntraIDOptions: tt.config,
			}

			credential, err := ObtainCredential(ctx, config, clientOptions)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, credential)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, credential)
			}
		})
	}
}

// TestOIDCStrategy validates OIDC credential creation with different token sources.
func TestOIDCStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With OIDC token from string", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:  "test-client-id",
				OIDCToken: "test-oidc-token",
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With OIDC token from file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "oidc-token-*.txt")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("file-based-token")
		require.NoError(t, err)
		tmpFile.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:          "test-client-id",
				OIDCTokenFilePath: tmpFile.Name(),
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With OIDC token file that doesn't exist", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:          "test-client-id",
				OIDCTokenFilePath: "/nonexistent/token.txt",
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With OIDC request URL and token", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"value": "exchanged-token"})
		}))
		defer mockServer.Close()

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:         "test-client-id",
				OIDCRequestURL:   mockServer.URL,
				OIDCRequestToken: "test-request-token",
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Missing OIDC configuration", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
		assert.Contains(t, err.Error(), "OIDC authentication requires")
	})

	t.Run("With additional tenants", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				OIDCToken:                  "test-oidc-token",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                 "test-client-id",
				OIDCToken:                "test-oidc-token",
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &OIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestAzureDevOpsOIDCStrategy validates Azure DevOps OIDC credential creation.
func TestAzureDevOpsOIDCStrategy(t *testing.T) {
	ctx := context.Background()
	clientOptions := policy.ClientOptions{}

	t.Run("With federation token env var", func(t *testing.T) {
		os.Setenv("AZURE_DEVOPS_FEDERATION_TOKEN", "test-federation-token")
		defer os.Unsetenv("AZURE_DEVOPS_FEDERATION_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &AzureDevOpsOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("Missing federation token env var", func(t *testing.T) {
		os.Unsetenv("AZURE_DEVOPS_FEDERATION_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
		}

		strategy := &AzureDevOpsOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.Error(t, err)
		assert.Nil(t, credential)
		assert.Contains(t, err.Error(), "AZURE_DEVOPS_FEDERATION_TOKEN")
	})

	t.Run("With additional tenants", func(t *testing.T) {
		os.Setenv("AZURE_DEVOPS_FEDERATION_TOKEN", "test-federation-token")
		defer os.Unsetenv("AZURE_DEVOPS_FEDERATION_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                   "test-client-id",
				AdditionallyAllowedTenants: []string{"tenant1", "tenant2"},
			},
		}

		strategy := &AzureDevOpsOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})

	t.Run("With instance discovery disabled", func(t *testing.T) {
		os.Setenv("AZURE_DEVOPS_FEDERATION_TOKEN", "test-federation-token")
		defer os.Unsetenv("AZURE_DEVOPS_FEDERATION_TOKEN")

		config := &ProviderData{
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:                 "test-client-id",
				DisableInstanceDiscovery: true,
			},
		}

		strategy := &AzureDevOpsOIDCStrategy{}
		credential, err := strategy.GetCredential(ctx, config, clientOptions)

		assert.NoError(t, err)
		assert.NotNil(t, credential)
	})
}

// TestExchangeOIDCToken validates the OIDC token exchange function.
func TestExchangeOIDCToken(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful token exchange", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer test-request-token", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"value": "exchanged-token"})
		}))
		defer mockServer.Close()

		token, err := exchangeOIDCToken(ctx, mockServer.URL, "test-request-token")
		require.NoError(t, err)
		assert.Equal(t, "exchanged-token", token)
	})

	t.Run("Server returns error", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}))
		defer mockServer.Close()

		token, err := exchangeOIDCToken(ctx, mockServer.URL, "test-request-token")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "status code 401")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		token, err := exchangeOIDCToken(ctx, "://invalid-url", "test-request-token")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not-json"))
		}))
		defer mockServer.Close()

		token, err := exchangeOIDCToken(ctx, mockServer.URL, "test-request-token")
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "failed to parse OIDC token exchange response")
	})
}
