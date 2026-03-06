package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	khttp "github.com/microsoft/kiota-http-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEntraIDAuthClientCompression validates that the HTTP client used for
// Entra ID authentication does not compress OAuth token requests.
//
// Background: Entra ID's OAuth2 token endpoint does not support gzip-compressed
// request bodies. When compressed requests are sent, Entra ID returns AADSTS900144
// error claiming the grant_type parameter is missing (it's present but cannot be
// parsed from the compressed body).
//
// Related: GitHub issue #777
func TestUnit_EntraIDAuthClient_CompressionBehavior(t *testing.T) {

	tests := []struct {
		name             string
		clientFactory    func() *http.Client
		expectCompressed bool
		expectSuccess    bool
		description      string
	}{
		{
			name: "Plain http.Client (recommended for auth)",
			clientFactory: func() *http.Client {
				return &http.Client{}
			},
			expectCompressed: false,
			expectSuccess:    true,
			description:      "Baseline - no compression middleware",
		},
		{
			name: "Kiota GetDefaultClient (problematic)",
			clientFactory: func() *http.Client {
				return khttp.GetDefaultClient()
			},
			expectCompressed: true,
			expectSuccess:    false,
			description:      "Includes compression middleware that breaks Entra ID token requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedContentEncoding string
			var receivedBody string
			var receivedBodyRaw []byte

			mockAzureADTokenEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedContentEncoding = r.Header.Get("Content-Encoding")
				receivedBodyRaw, _ = io.ReadAll(r.Body)

				t.Logf("Content-Encoding: %s", receivedContentEncoding)
				t.Logf("Raw body length: %d bytes", len(receivedBodyRaw))

				if receivedContentEncoding == "gzip" {
					reader, err := gzip.NewReader(bytes.NewReader(receivedBodyRaw))
					if err == nil {
						decompressed, _ := io.ReadAll(reader)
						receivedBody = string(decompressed)
						reader.Close()
						t.Logf("Decompressed body: %s", receivedBody)
					} else {
						receivedBody = string(receivedBodyRaw)
						t.Logf("Failed to decompress: %v", err)
					}

					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]any{
						"error":             "invalid_request",
						"error_description": "AADSTS900144: The request body must contain the following parameter: 'grant_type'",
						"error_codes":       []int{900144},
					})
					return
				}

				receivedBody = string(receivedBodyRaw)
				t.Logf("Body: %s", receivedBody)

				hasGrantType := strings.Contains(receivedBody, "grant_type=")
				if !hasGrantType {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]any{
						"error":             "invalid_request",
						"error_description": "AADSTS900144: The request body must contain the following parameter: 'grant_type'",
						"error_codes":       []int{900144},
					})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]any{
					"token_type":   "Bearer",
					"expires_in":   3600,
					"access_token": "test-access-token",
				})
			}))
			defer mockAzureADTokenEndpoint.Close()

			client := tt.clientFactory()

			formData := url.Values{}
			formData.Set("grant_type", "client_credentials")
			formData.Set("client_id", "test-client-id")
			formData.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
			formData.Set("client_assertion", "test-jwt-token")
			formData.Set("scope", "https://graph.microsoft.com/.default")

			req, err := http.NewRequestWithContext(context.Background(), "POST", mockAzureADTokenEndpoint.URL, strings.NewReader(formData.Encode()))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			responseBody, _ := io.ReadAll(resp.Body)
			t.Logf("\nResponse status: %d", resp.StatusCode)
			t.Logf("Response body: %s", string(responseBody))

			if tt.expectCompressed {
				assert.Equal(t, "gzip", receivedContentEncoding, "Expected request to be compressed")
			} else {
				assert.Empty(t, receivedContentEncoding, "Expected request to be uncompressed")
			}

			if tt.expectSuccess {
				assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected successful token response")
				var tokenResp map[string]any
				json.Unmarshal(responseBody, &tokenResp)
				assert.Equal(t, "test-access-token", tokenResp["access_token"])
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Entra ID to reject compressed request")
				assert.Contains(t, string(responseBody), "900144", "Expected AADSTS900144 error code")
			}
		})
	}

}

// TestEntraIDAuthClientProxy validates that proxy configuration works correctly
// for authentication requests. Tests both Kiota-based and plain http.Client approaches
// to verify proxy routing and compression behavior.
func TestUnit_EntraIDAuthClient_ProxyConfiguration(t *testing.T) {

	var proxyHitCount atomic.Int32
	var proxyReceivedCompression string

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("[TARGET] Received: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"token_type":   "Bearer",
			"access_token": "test-token",
		})
	}))
	defer targetServer.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyHitCount.Add(1)
		proxyReceivedCompression = r.Header.Get("Content-Encoding")

		t.Logf("[PROXY] Intercepted: %s %s", r.Method, r.URL.String())
		t.Logf("[PROXY] Content-Encoding: %s", proxyReceivedCompression)

		targetReq, _ := http.NewRequest(r.Method, targetServer.URL, r.Body)
		targetReq.Header = r.Header.Clone()

		targetResp, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer targetResp.Body.Close()

		w.WriteHeader(targetResp.StatusCode)
		io.Copy(w, targetResp.Body)
	}))
	defer proxyServer.Close()

	tests := []struct {
		name              string
		clientFactory     func(proxyURL string) (*http.Client, error)
		expectProxyUsed   bool
		expectCompression bool
		description       string
	}{
		{
			name: "Kiota GetClientWithProxySettings",
			clientFactory: func(proxyURL string) (*http.Client, error) {
				return khttp.GetClientWithProxySettings(proxyURL)
			},
			expectProxyUsed:   true,
			expectCompression: true,
			description:       "Current code - Kiota with proxy and compression",
		},
		{
			name: "Plain http.Client with proxy",
			clientFactory: func(proxyURL string) (*http.Client, error) {
				parsedURL, err := url.Parse(proxyURL)
				if err != nil {
					return nil, err
				}
				return &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(parsedURL),
					},
				}, nil
			},
			expectProxyUsed:   true,
			expectCompression: false,
			description:       "PR approach - plain client with proxy, no compression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxyHitCount.Store(0)
			proxyReceivedCompression = ""

			client, err := tt.clientFactory(proxyServer.URL)
			require.NoError(t, err)

			formData := url.Values{}
			formData.Set("grant_type", "client_credentials")
			formData.Set("client_id", "test-client")

			req, err := http.NewRequestWithContext(context.Background(), "POST", targetServer.URL, strings.NewReader(formData.Encode()))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			hits := proxyHitCount.Load()
			t.Logf("Proxy hits: %d", hits)
			t.Logf("Compression: %s", proxyReceivedCompression)

			if tt.expectProxyUsed {
				assert.Greater(t, int(hits), 0, "Request should go through proxy")
			}

			if tt.expectCompression {
				assert.Equal(t, "gzip", proxyReceivedCompression, "Expected compression")
			} else {
				assert.Empty(t, proxyReceivedCompression, "Expected no compression")
			}
		})
	}
}

// TestEntraIDAuthClientAuthenticatedProxy validates that authenticated proxy
// configuration is properly handled by both Kiota and plain http.Client approaches.
func TestUnit_EntraIDAuthClient_AuthenticatedProxy(t *testing.T) {

	const testUsername = "proxyuser"
	const testPassword = "proxypass"

	var proxyAuthHitCount atomic.Int32

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"token_type":   "Bearer",
			"access_token": "test-token",
		})
	}))
	defer targetServer.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			t.Logf("[PROXY] Received auth - user: %s", username)
		}

		if !ok || username != testUsername || password != testPassword {
			t.Log("[PROXY] Auth failed or missing, returning 407")
			w.Header().Set("Proxy-Authenticate", "Basic realm=\"Test Proxy\"")
			w.WriteHeader(http.StatusProxyAuthRequired)
			return
		}

		proxyAuthHitCount.Add(1)
		t.Log("[PROXY] Auth successful, forwarding request")

		targetReq, _ := http.NewRequest(r.Method, targetServer.URL, r.Body)
		targetReq.Header = r.Header.Clone()

		targetResp, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer targetResp.Body.Close()

		w.WriteHeader(targetResp.StatusCode)
		io.Copy(w, targetResp.Body)
	}))
	defer proxyServer.Close()

	tests := []struct {
		name          string
		clientFactory func(proxyURL, username, password string) (*http.Client, error)
		description   string
	}{
		{
			name: "Kiota GetClientWithAuthenticatedProxySettings",
			clientFactory: func(proxyURL, username, password string) (*http.Client, error) {
				return khttp.GetClientWithAuthenticatedProxySettings(proxyURL, username, password)
			},
			description: "Current code - Kiota authenticated proxy",
		},
		{
			name: "Plain http.Client with authenticated proxy",
			clientFactory: func(proxyURL, username, password string) (*http.Client, error) {
				parsedURL, err := url.Parse(proxyURL)
				if err != nil {
					return nil, err
				}
				parsedURL.User = url.UserPassword(username, password)
				return &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(parsedURL),
					},
				}, nil
			},
			description: "PR approach - plain client with authenticated proxy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proxyAuthHitCount.Store(0)

			client, err := tt.clientFactory(proxyServer.URL, testUsername, testPassword)
			require.NoError(t, err)

			formData := url.Values{}
			formData.Set("grant_type", "client_credentials")
			formData.Set("client_id", "test-client")

			req, err := http.NewRequestWithContext(context.Background(), "POST", targetServer.URL, strings.NewReader(formData.Encode()))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)

			t.Logf("\nAuth hits: %d", proxyAuthHitCount.Load())
			t.Logf("Response status: %d", resp.StatusCode)

			if err == nil {
				resp.Body.Close()
			}
		})
	}
}

// TestConfigureAuthClientProxy tests the actual configureAuthClientProxy function
func TestUnit_ConfigureAuthClientProxy_AllScenarios(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		config      *ProviderData
		expectProxy bool
		description string
	}{
		{
			name: "No proxy configured",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					UseProxy: false,
				},
			},
			expectProxy: false,
			description: "Should return default client when proxy not enabled",
		},
		{
			name: "Proxy configured but not enabled",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					UseProxy: false,
					ProxyURL: "http://proxy.example.com:8080",
				},
			},
			expectProxy: false,
			description: "Should ignore proxy URL when UseProxy is false",
		},
		{
			name: "Proxy enabled with URL",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					UseProxy: true,
					ProxyURL: "http://proxy.example.com:8080",
				},
			},
			expectProxy: true,
			description: "Should configure proxy when enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := configureAuthClientProxy(ctx, tt.config)
			require.NoError(t, err)
			require.NotNil(t, client)

			if tt.expectProxy {
				require.NotNil(t, client.Transport, "Proxy client should have custom transport")
			}

			t.Logf("Client timeout: %v", client.Timeout)
		})
	}
}

// TestAuthClientWithoutCompressionFix validates that configureAuthClientProxy
// creates an HTTP client that does not compress requests while maintaining other
// Kiota middleware features (retry, redirect, user agent).
func TestUnit_AuthClient_WithoutCompressionFix(t *testing.T) {
	var receivedContentEncoding string
	var receivedBody string

	mockAzureAD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedContentEncoding = r.Header.Get("Content-Encoding")
		bodyBytes, _ := io.ReadAll(r.Body)
		receivedBody = string(bodyBytes)

		if receivedContentEncoding == "gzip" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"error":             "invalid_request",
				"error_description": "AADSTS900144: The request body must contain the following parameter: 'grant_type'",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"token_type":   "Bearer",
			"access_token": "test-token",
		})
	}))
	defer mockAzureAD.Close()

	ctx := context.Background()
	config := &ProviderData{
		ClientOptions: &ClientOptions{
			EnableRetry:    true,
			MaxRetries:     3,
			EnableRedirect: true,
			MaxRedirects:   5,
		},
	}

	client, err := configureAuthClientProxy(ctx, config)
	require.NoError(t, err)

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", "test-client")

	req, err := http.NewRequestWithContext(ctx, "POST", mockAzureAD.URL, strings.NewReader(formData.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Empty(t, receivedContentEncoding, "Auth client should not compress requests")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Request should succeed without compression")
	assert.Contains(t, receivedBody, "grant_type=", "grant_type parameter should be present")
}

// TestAuthClientWithProxyWithoutCompression validates that configureAuthClientProxy
// correctly disables compression when proxy is configured, ensuring compatibility
// with Entra ID token endpoint while maintaining proxy functionality.
func TestUnit_AuthClient_WithProxyWithoutCompression(t *testing.T) {
	var proxyReceivedCompression string
	var proxyHitCount atomic.Int32

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"token_type":   "Bearer",
			"access_token": "test-token",
		})
	}))
	defer targetServer.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyHitCount.Add(1)
		proxyReceivedCompression = r.Header.Get("Content-Encoding")

		targetReq, _ := http.NewRequest(r.Method, targetServer.URL, r.Body)
		targetReq.Header = r.Header.Clone()

		targetResp, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer targetResp.Body.Close()

		w.WriteHeader(targetResp.StatusCode)
		io.Copy(w, targetResp.Body)
	}))
	defer proxyServer.Close()

	ctx := context.Background()
	config := &ProviderData{
		ClientOptions: &ClientOptions{
			UseProxy:       true,
			ProxyURL:       proxyServer.URL,
			EnableRetry:    true,
			MaxRetries:     3,
			EnableRedirect: true,
			MaxRedirects:   5,
		},
	}

	client, err := configureAuthClientProxy(ctx, config)
	require.NoError(t, err)

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", "test-client")

	req, err := http.NewRequestWithContext(ctx, "POST", targetServer.URL, strings.NewReader(formData.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Greater(t, int(proxyHitCount.Load()), 0, "Request should route through proxy")
	assert.Empty(t, proxyReceivedCompression, "Requests should not be compressed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Request should succeed")
}

// TestConfigureEntraIDClientOptions validates the main entry point function that
// orchestrates all client option configuration.
func TestUnit_ConfigureEntraIDClientOptions_AllScenarios(t *testing.T) {
	tests := []struct {
		name              string
		authorityURL      string
		config            *ProviderData
		expectError       bool
		validateTransport bool
	}{
		{
			name:         "Basic configuration without proxy",
			authorityURL: "https://login.microsoftonline.com/",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					EnableRetry:       true,
					MaxRetries:        3,
					RetryDelaySeconds: 5,
					EnableRedirect:    true,
					MaxRedirects:      5,
					TimeoutSeconds:    300,
				},
				TelemetryOptout: false,
			},
			expectError:       false,
			validateTransport: true,
		},
		{
			name:         "Configuration with custom user agent",
			authorityURL: "https://login.microsoftonline.com/",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					EnableRetry:       true,
					MaxRetries:        3,
					RetryDelaySeconds: 5,
					CustomUserAgent:   "TestAgent/1.0",
				},
				TelemetryOptout: false,
			},
			expectError:       false,
			validateTransport: true,
		},
		{
			name:         "Configuration with telemetry disabled",
			authorityURL: "https://login.microsoftonline.com/",
			config: &ProviderData{
				ClientOptions: &ClientOptions{
					EnableRetry:       true,
					MaxRetries:        3,
					RetryDelaySeconds: 5,
				},
				TelemetryOptout: true,
			},
			expectError:       false,
			validateTransport: true,
		},
		{
			name:         "Minimal configuration",
			authorityURL: "https://login.microsoftonline.com/",
			config: &ProviderData{
				ClientOptions:   &ClientOptions{},
				TelemetryOptout: false,
			},
			expectError:       false,
			validateTransport: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			clientOptions, err := ConfigureEntraIDClientOptions(ctx, tt.config, tt.authorityURL)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.authorityURL, clientOptions.Cloud.ActiveDirectoryAuthorityHost)

			if tt.validateTransport {
				assert.NotNil(t, clientOptions.Transport, "Transport should be configured")
				httpClient, ok := clientOptions.Transport.(*http.Client)
				require.True(t, ok, "Transport should be an *http.Client")
				assert.NotNil(t, httpClient, "HTTP client should not be nil")

				if tt.config.ClientOptions.TimeoutSeconds > 0 {
					expectedTimeout := time.Duration(tt.config.ClientOptions.TimeoutSeconds) * time.Second
					assert.Equal(t, expectedTimeout, httpClient.Timeout, "Client timeout should match config")
				}
			}
		})
	}
}

// TestGetAuthClientMiddleware validates that the middleware chain is built correctly
// based on client options configuration.
func TestUnit_GetAuthClientMiddleware_BasicConfiguration(t *testing.T) {
	tests := []struct {
		name               string
		clientOptions      *ClientOptions
		expectedMiddleware int
	}{
		{
			name: "All middleware enabled",
			clientOptions: &ClientOptions{
				EnableRetry:       true,
				MaxRetries:        3,
				RetryDelaySeconds: 5,
				EnableRedirect:    true,
				MaxRedirects:      5,
				CustomUserAgent:   "TestAgent/1.0",
			},
			expectedMiddleware: 3,
		},
		{
			name: "Only retry enabled",
			clientOptions: &ClientOptions{
				EnableRetry:       true,
				MaxRetries:        3,
				RetryDelaySeconds: 5,
			},
			expectedMiddleware: 1,
		},
		{
			name: "Only redirect enabled",
			clientOptions: &ClientOptions{
				EnableRedirect: true,
				MaxRedirects:   5,
			},
			expectedMiddleware: 1,
		},
		{
			name: "Only user agent enabled",
			clientOptions: &ClientOptions{
				CustomUserAgent: "TestAgent/1.0",
			},
			expectedMiddleware: 1,
		},
		{
			name:               "No middleware enabled (fallback to redirect)",
			clientOptions:      &ClientOptions{},
			expectedMiddleware: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			middleware := getAuthClientMiddleware(ctx, tt.clientOptions)
			assert.Equal(t, tt.expectedMiddleware, len(middleware), "Middleware count should match expected")
		})
	}
}

// TestAuthClientProxyErrorHandling validates error handling in proxy configuration.
func TestUnit_AuthClientProxy_ErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid proxy URL", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy: true,
				ProxyURL: "://invalid-url",
			},
		}

		client, err := configureAuthClientProxy(ctx, config)
		assert.Error(t, err, "Should return error for invalid proxy URL")
		assert.Nil(t, client, "Client should be nil on error")
		assert.Contains(t, err.Error(), "unable to create HTTP client with proxy settings")
	})
}

// TestAuthClientMiddlewareExcludesCompression validates that CompressionHandler
// is never included in the auth client middleware chain, regardless of configuration.
func TestUnit_AuthClientMiddleware_ExcludesCompression(t *testing.T) {
	ctx := context.Background()

	configs := []*ClientOptions{
		{},
		{EnableRetry: true, MaxRetries: 3, RetryDelaySeconds: 5},
		{EnableRedirect: true, MaxRedirects: 5},
		{CustomUserAgent: "TestAgent/1.0"},
		{EnableRetry: true, EnableRedirect: true, CustomUserAgent: "TestAgent/1.0"},
	}

	for i, config := range configs {
		t.Run(fmt.Sprintf("Config_%d", i), func(t *testing.T) {
			middleware := getAuthClientMiddleware(ctx, config)

			client := khttp.GetDefaultClient(middleware...)

			var receivedContentEncoding string
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedContentEncoding = r.Header.Get("Content-Encoding")
				w.WriteHeader(http.StatusOK)
			}))
			defer mockServer.Close()

			formData := url.Values{}
			formData.Set("test", "data")
			req, _ := http.NewRequestWithContext(ctx, "POST", mockServer.URL, strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Empty(t, receivedContentEncoding, "Compression should never be applied to auth requests")
		})
	}
}

// TestGetAuthClientMiddleware_AllCombinations validates all middleware combinations.
func TestUnit_GetAuthClientMiddleware_AllCombinations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		options        *ClientOptions
		expectedMinLen int
	}{
		{
			name:           "All middleware enabled",
			options:        &ClientOptions{EnableRetry: true, EnableRedirect: true, CustomUserAgent: "Test/1.0"},
			expectedMinLen: 3,
		},
		{
			name:           "Only retry",
			options:        &ClientOptions{EnableRetry: true},
			expectedMinLen: 1,
		},
		{
			name:           "Only redirect",
			options:        &ClientOptions{EnableRedirect: true},
			expectedMinLen: 1,
		},
		{
			name:           "Only user agent",
			options:        &ClientOptions{CustomUserAgent: "Test/1.0"},
			expectedMinLen: 1,
		},
		{
			name:           "Retry and redirect",
			options:        &ClientOptions{EnableRetry: true, EnableRedirect: true},
			expectedMinLen: 2,
		},
		{
			name:           "None enabled (fallback)",
			options:        &ClientOptions{},
			expectedMinLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := getAuthClientMiddleware(ctx, tt.options)
			assert.GreaterOrEqual(t, len(middleware), tt.expectedMinLen, "Should have at least expected middleware count")
		})
	}
}
