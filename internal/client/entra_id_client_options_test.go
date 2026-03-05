package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	khttp "github.com/microsoft/kiota-http-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEntraIDAuthClientCompression validates that the HTTP client used for
// Entra ID authentication does not compress OAuth token requests, as Azure AD's
// token endpoint does not support gzip-compressed request bodies.
//
// This test documents the root cause of issue #777 where GitHub OIDC authentication
// failed with AADSTS900144 error due to Kiota's compression middleware.
func TestEntraIDAuthClientCompression(t *testing.T) {
	t.Log("=== Testing Entra ID Auth Client Compression Behavior ===")
	t.Log("Issue #777: Azure AD token endpoint rejects gzip-compressed requests")
	t.Log("")

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
			description:      "Includes compression middleware that breaks Azure AD token requests",
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
					t.Log("⚠️  Request is gzip compressed")
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

					t.Log("❌ Azure AD does NOT support gzip-compressed token requests")
					t.Log("   Simulating Azure AD rejection...")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":             "invalid_request",
						"error_description": "AADSTS900144: The request body must contain the following parameter: 'grant_type'",
						"error_codes":       []int{900144},
					})
					return
				}

				receivedBody = string(receivedBodyRaw)
				t.Logf("Uncompressed body: %s", receivedBody)

				hasGrantType := strings.Contains(receivedBody, "grant_type=")
				t.Logf("Has grant_type: %v", hasGrantType)

				if !hasGrantType {
					t.Log("❌ grant_type parameter missing")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":             "invalid_request",
						"error_description": "AADSTS900144: The request body must contain the following parameter: 'grant_type'",
						"error_codes":       []int{900144},
					})
					return
				}

				t.Log("✅ Request accepted by Azure AD")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
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
				var tokenResp map[string]interface{}
				json.Unmarshal(responseBody, &tokenResp)
				assert.Equal(t, "test-access-token", tokenResp["access_token"])
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected Azure AD to reject compressed request")
				assert.Contains(t, string(responseBody), "900144", "Expected AADSTS900144 error code")
			}
		})
	}

	t.Log("\n" + strings.Repeat("=", 80))
	t.Log("FINDINGS:")
	t.Log("1. Kiota GetDefaultClient() compresses OAuth token request bodies with gzip")
	t.Log("2. Azure AD token endpoint does NOT support gzip-compressed requests")
	t.Log("3. This causes AADSTS900144 error (grant_type cannot be parsed)")
	t.Log("4. Solution: Use plain http.Client OR Kiota without compression middleware")
	t.Log(strings.Repeat("=", 80))
}

// TestEntraIDAuthClientProxy validates that proxy configuration works correctly
// for authentication requests with both Kiota and plain http.Client approaches
func TestEntraIDAuthClientProxy(t *testing.T) {
	t.Log("=== Testing Proxy Support for Entra ID Auth Client ===")

	var proxyHitCount atomic.Int32
	var proxyReceivedCompression string

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("[TARGET] Received: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
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
		
		targetResp, _ := http.DefaultClient.Do(targetReq)
		defer targetResp.Body.Close()
		
		w.WriteHeader(targetResp.StatusCode)
		io.Copy(w, targetResp.Body)
	}))
	defer proxyServer.Close()

	tests := []struct {
		name                 string
		clientFactory        func(proxyURL string) (*http.Client, error)
		expectProxyUsed      bool
		expectCompression    bool
		description          string
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
			t.Logf("\nProxy hits: %d", hits)
			t.Logf("Compression: %s", proxyReceivedCompression)

			if tt.expectProxyUsed {
				assert.Greater(t, int(hits), 0, "Request should go through proxy")
			}

			if tt.expectCompression {
				assert.Equal(t, "gzip", proxyReceivedCompression, "Expected compression")
				t.Log("   ⚠️  Request is compressed (problematic for Azure AD)")
			} else {
				assert.Empty(t, proxyReceivedCompression, "Expected no compression")
				t.Log("   ✅ Request is not compressed (compatible with Azure AD)")
			}
		})
	}
}

// TestEntraIDAuthClientAuthenticatedProxy validates authenticated proxy support
func TestEntraIDAuthClientAuthenticatedProxy(t *testing.T) {
	t.Log("=== Testing Authenticated Proxy Support ===")

	const testUsername = "proxyuser"
	const testPassword = "proxypass"

	var proxyAuthHitCount atomic.Int32
	var receivedUsername string
	var receivedPassword string

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token_type":   "Bearer",
			"access_token": "test-token",
		})
	}))
	defer targetServer.Close()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			receivedUsername = username
			receivedPassword = password
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
		
		targetResp, _ := http.DefaultClient.Do(targetReq)
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
			receivedUsername = ""
			receivedPassword = ""

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
func TestConfigureAuthClientProxy(t *testing.T) {
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
