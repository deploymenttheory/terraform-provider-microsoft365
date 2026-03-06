package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnit_AuthenticatedHTTPClient_Do validates the Do method adds proper authentication
// and headers to requests.
func TestUnit_AuthenticatedHTTPClient_Do(t *testing.T) {
	var receivedAuthHeader string
	var receivedAcceptHeader string
	var receivedContentType string
	var receivedConsistencyLevel string
	var receivedMethod string

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		receivedAcceptHeader = r.Header.Get("Accept")
		receivedContentType = r.Header.Get("Content-Type")
		receivedConsistencyLevel = r.Header.Get("ConsistencyLevel")
		receivedMethod = r.Method

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer mockServer.Close()

	tests := []struct {
		name                     string
		method                   string
		body                     string
		setContentType           bool
		expectedContentType      string
		expectedConsistencyLevel string
	}{
		{
			name:                     "GET request adds ConsistencyLevel",
			method:                   "GET",
			body:                     "",
			setContentType:           false,
			expectedContentType:      "",
			expectedConsistencyLevel: "eventual",
		},
		{
			name:                     "POST request adds default Content-Type",
			method:                   "POST",
			body:                     `{"test": "data"}`,
			setContentType:           false,
			expectedContentType:      "application/json",
			expectedConsistencyLevel: "",
		},
		{
			name:                     "POST request preserves custom Content-Type",
			method:                   "POST",
			body:                     "test=data",
			setContentType:           true,
			expectedContentType:      "application/x-www-form-urlencoded",
			expectedConsistencyLevel: "",
		},
		{
			name:                     "PUT request adds default Content-Type",
			method:                   "PUT",
			body:                     `{"test": "data"}`,
			setContentType:           false,
			expectedContentType:      "application/json",
			expectedConsistencyLevel: "",
		},
		{
			name:                     "PATCH request adds default Content-Type",
			method:                   "PATCH",
			body:                     `{"test": "data"}`,
			setContentType:           false,
			expectedContentType:      "application/json",
			expectedConsistencyLevel: "",
		},
		{
			name:                     "DELETE request no Content-Type",
			method:                   "DELETE",
			body:                     "",
			setContentType:           false,
			expectedContentType:      "",
			expectedConsistencyLevel: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := &mockTokenCredential{token: "test-access-token"}
			client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			}

			req, err := http.NewRequest(tt.method, mockServer.URL+"/test", bodyReader)
			require.NoError(t, err)

			if tt.setContentType {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, "Bearer test-access-token", receivedAuthHeader, "Authorization header should be set")
			assert.Equal(t, "application/json", receivedAcceptHeader, "Accept header should be set")
			assert.Equal(t, tt.method, receivedMethod, "Method should match")

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, receivedContentType, "Content-Type should match expected")
			}

			if tt.expectedConsistencyLevel != "" {
				assert.Equal(t, tt.expectedConsistencyLevel, receivedConsistencyLevel, "ConsistencyLevel should match expected")
			}
		})
	}
}

// TestUnit_AuthenticatedHTTPClient_Do_TokenError validates error handling when token acquisition fails.
func TestUnit_AuthenticatedHTTPClient_Do_TokenError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	cred := &mockTokenCredential{err: fmt.Errorf("token acquisition failed")}
	client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

	req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get access token")
}

// TestDoWithRetry validates retry logic for 429 rate limit errors.
func TestDoWithRetry(t *testing.T) {
	t.Run("Success on first attempt", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := DoWithRetry(ctx, client, req, 3)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Non-429 error returns immediately", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := DoWithRetry(ctx, client, req, 3)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, 1, attemptCount, "Should not retry for non-429 errors")
	})

	t.Run("Retry on 429 then success", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount < 2 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"code":    "TooManyRequests",
						"message": "Rate limit exceeded",
					},
				})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := DoWithRetry(ctx, client, req, 3)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 2, attemptCount, "Should have retried once")
	})

	t.Run("Max retries exceeded returns 429", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "TooManyRequests",
					"message": "Rate limit exceeded",
				},
			})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := DoWithRetry(ctx, client, req, 2)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		assert.Equal(t, 3, attemptCount, "Should have attempted 3 times (initial + 2 retries)")
	})

	t.Run("Context cancellation during retry", func(t *testing.T) {
		attemptCount := 0
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"code":    "TooManyRequests",
					"message": "Rate limit exceeded",
				},
			})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		resp, err := DoWithRetry(ctx, client, req, 5)

		if err != nil {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "context")
		} else {
			defer resp.Body.Close()
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	})

	t.Run("Request with body can be retried", func(t *testing.T) {
		attemptCount := 0
		var receivedBodies []string

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			bodyBytes, _ := io.ReadAll(r.Body)
			receivedBodies = append(receivedBodies, string(bodyBytes))

			if attemptCount < 2 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"code":    "TooManyRequests",
						"message": "Rate limit exceeded",
					},
				})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		}))
		defer mockServer.Close()

		cred := &mockTokenCredential{token: "test-token"}
		client := NewAuthenticatedHTTPClient(&http.Client{}, cred, "https://graph.microsoft.com/.default", mockServer.URL)

		body := `{"test": "data"}`
		req, err := http.NewRequest("POST", mockServer.URL+"/test", strings.NewReader(body))
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := DoWithRetry(ctx, client, req, 3)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 2, attemptCount, "Should have retried once")
		assert.Len(t, receivedBodies, 2, "Should have received body twice")
		assert.Equal(t, body, receivedBodies[0], "First attempt should have correct body")
		assert.Equal(t, body, receivedBodies[1], "Retry should have correct body")
	})
}

// TestUnit_AuthenticatedHTTPClient_GetBaseURL validates the GetBaseURL method.
func TestUnit_AuthenticatedHTTPClient_GetBaseURL(t *testing.T) {
	baseURL := "https://graph.microsoft.com/beta"
	client := NewAuthenticatedHTTPClient(&http.Client{}, &mockTokenCredential{}, "", baseURL)

	assert.Equal(t, baseURL, client.GetBaseURL())
}

// TestUnit_AuthenticatedHTTPClient_GetClient validates the GetClient method.
func TestUnit_AuthenticatedHTTPClient_GetClient(t *testing.T) {
	baseClient := &http.Client{Timeout: 30 * time.Second}
	client := NewAuthenticatedHTTPClient(baseClient, &mockTokenCredential{}, "", "")

	assert.Equal(t, baseClient, client.GetClient())
	assert.Equal(t, 30*time.Second, client.GetClient().Timeout)
}
