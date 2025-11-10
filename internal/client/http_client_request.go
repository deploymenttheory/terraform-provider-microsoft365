package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Do performs an HTTP request with authentication
func (c *AuthenticatedHTTPClient) Do(req *http.Request) (*http.Response, error) {

	token, err := c.credential.GetToken(req.Context(), policy.TokenRequestOptions{
		Scopes: []string{c.scope},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Accept", "application/json")

	// Set default headers for Graph API
	if req.Header.Get("Content-Type") == "" && (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add consistency level header for certain operations
	if req.Method == "GET" {
		req.Header.Set("ConsistencyLevel", "eventual")
	}

	return c.client.Do(req)
}

// doHTTPRequestWithRetry performs an HTTP request with exponential backoff retry logic for 429 rate limit errors
// It uses the AuthenticatedHTTPClient to preserve authentication
func DoWithRetry(ctx context.Context, httpClient *AuthenticatedHTTPClient, req *http.Request, maxRetries int) (*http.Response, error) {
	var lastResp *http.Response
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Clone the request body for retries if needed
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, lastErr = io.ReadAll(req.Body)
			if lastErr != nil {
				return nil, fmt.Errorf("failed to read request body: %w", lastErr)
			}
			req.Body.Close()
		}

		// Create new request with fresh body
		if len(bodyBytes) > 0 {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		tflog.Debug(ctx, fmt.Sprintf("Executing HTTP request (attempt %d/%d)", attempt+1, maxRetries+1), map[string]any{
			"method": req.Method,
			"url":    req.URL.String(),
		})

		httpResp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			tflog.Warn(ctx, "HTTP request failed with error", map[string]any{
				"attempt": attempt + 1,
				"error":   err.Error(),
			})
			continue
		}

		// If not a 429 error, return immediately
		if httpResp.StatusCode != http.StatusTooManyRequests {
			return httpResp, nil
		}

		// Handle 429 rate limit error
		lastResp = httpResp

		errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
		retryDelay := errors.GetHTTPRetryDelay(ctx, errorInfo, attempt)

		tflog.Warn(ctx, "Rate limit exceeded (429), will retry", map[string]any{
			"attempt":          attempt + 1,
			"max_retries":      maxRetries + 1,
			"retry_after":      errorInfo.RetryAfter,
			"retry_delay":      retryDelay.String(),
			"throttled_reason": errorInfo.ThrottledReason,
			"request_id":       errorInfo.RequestID,
		})

		// Close the response body before retrying
		httpResp.Body.Close()

		// Don't sleep after the last attempt
		if attempt < maxRetries {
			// Check if context is still valid before sleeping
			if ctx.Err() != nil {
				tflog.Warn(ctx, "Context cancelled before retry delay", map[string]any{
					"attempt":     attempt + 1,
					"context_err": ctx.Err().Error(),
				})
				return nil, ctx.Err()
			}

			// Sleep for the retry delay
			// Note: We use time.Sleep instead of select with ctx.Done() to ensure
			// we respect the API's Retry-After header. If the context is cancelled
			// during the sleep, the next attempt will detect it immediately.
			tflog.Debug(ctx, fmt.Sprintf("Waiting %s before retry attempt %d", retryDelay, attempt+2), map[string]any{
				"retry_delay":  retryDelay.String(),
				"next_attempt": attempt + 2,
			})
			time.Sleep(retryDelay)

			// Check context validity immediately after sleep
			if ctx.Err() != nil {
				tflog.Warn(ctx, "Context cancelled during retry delay", map[string]any{
					"attempt":     attempt + 1,
					"context_err": ctx.Err().Error(),
				})
				return nil, ctx.Err()
			}
		}
	}

	// If we got a 429 response, return it
	if lastResp != nil {
		return lastResp, nil
	}

	// Return the last error if no response was received
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("max retries exceeded")
}

// GetBaseURL returns the base URL for this client
func (c *AuthenticatedHTTPClient) GetBaseURL() string {
	return c.baseURL
}

// GetClient returns the underlying HTTP client
func (c *AuthenticatedHTTPClient) GetClient() *http.Client {
	return c.client
}
