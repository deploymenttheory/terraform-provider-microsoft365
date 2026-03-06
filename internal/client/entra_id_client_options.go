package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	khttp "github.com/microsoft/kiota-http-go"
	"golang.org/x/exp/rand"
)

// ConfigureEntraIDClientOptions configures the Entra ID client options based on the provided configuration
func ConfigureEntraIDClientOptions(ctx context.Context, config *ProviderData, authorityURL string) (policy.ClientOptions, error) {
	tflog.Info(ctx, "Starting Entra ID client options configuration")
	tflog.Info(ctx, "Authority URL: "+authorityURL)

	tflog.Info(ctx, "Initializing authentication client options")
	clientOptions := initializeAuthClientOptions(ctx, authorityURL)

	tflog.Info(ctx, "Configuring Retry options")
	configureRetryOptions(ctx, &clientOptions, config)

	tflog.Info(ctx, "Configuring Telemetry options")
	configureTelemetryOptions(ctx, &clientOptions, config)

	tflog.Info(ctx, "Configuring authentication timeout")
	configureAuthTimeout(ctx, &clientOptions, config)

	httpClient, err := configureAuthClientProxy(ctx, config)
	if err != nil {
		tflog.Error(ctx, "Failed to configure HTTP client")
		return clientOptions, err
	}

	if config.ClientOptions.TimeoutSeconds > 0 {
		httpClient.Timeout = time.Duration(config.ClientOptions.TimeoutSeconds) * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Auth client timeout set to %v", httpClient.Timeout))
	}

	clientOptions.Transport = httpClient

	tflog.Info(ctx, "Entra ID client options configuration completed successfully")
	return clientOptions, nil
}

// initializeAuthClientOptions creates the base policy.ClientOptions with the authority URL.
func initializeAuthClientOptions(ctx context.Context, authorityURL string) policy.ClientOptions {
	options := policy.ClientOptions{
		Cloud: cloud.Configuration{
			ActiveDirectoryAuthorityHost: authorityURL,
		},
	}
	tflog.Debug(ctx, "Authentication Client options initialized with authority URL: "+authorityURL)
	return options
}

// configureRetryOptions sets up retry behavior for authentication requests with exponential
// backoff and jitter for transient failures.
func configureRetryOptions(ctx context.Context, clientOptions *policy.ClientOptions, config *ProviderData) {
	maxRetries := int32(config.ClientOptions.MaxRetries)
	baseDelay := time.Duration(config.ClientOptions.RetryDelaySeconds) * time.Second

	clientOptions.Retry = policy.RetryOptions{
		MaxRetries:    maxRetries,
		RetryDelay:    baseDelay,
		MaxRetryDelay: baseDelay * 10,
		StatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
		ShouldRetry: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			if resp == nil {
				return false
			}

			if slices.Contains(clientOptions.Retry.StatusCodes, resp.StatusCode) {
				executionCount := resp.Request.Context().Value("RetryCount").(int32)
					exponentialBackoff := baseDelay * time.Duration(math.Pow(2, float64(executionCount)))
					jitter := time.Duration(rand.Int63n(int64(baseDelay)))
					delayWithJitter := exponentialBackoff + jitter

					if delayWithJitter > clientOptions.Retry.MaxRetryDelay {
						delayWithJitter = clientOptions.Retry.MaxRetryDelay
					}

					tflog.Debug(ctx, fmt.Sprintf("Retrying request due to status code %d. Delay with jitter: %v (base: %v, jitter: %v)", resp.StatusCode, delayWithJitter, exponentialBackoff, jitter))

					time.Sleep(delayWithJitter)
					return true
			}
			return false
		},
	}

	tflog.Debug(ctx, fmt.Sprintf("Retry options set: MaxRetries=%d, BaseRetryDelay=%v",
		clientOptions.Retry.MaxRetries, time.Duration(config.ClientOptions.RetryDelaySeconds)*time.Second))
}

// configureTelemetryOptions configures telemetry settings for authentication requests.
func configureTelemetryOptions(ctx context.Context, clientOptions *policy.ClientOptions, config *ProviderData) {
	clientOptions.Telemetry = policy.TelemetryOptions{
		ApplicationID: config.ClientOptions.CustomUserAgent,
		Disabled:      config.TelemetryOptout,
	}
	tflog.Debug(ctx, fmt.Sprintf("Telemetry options set: ApplicationID=%s, Disabled=%t",
		clientOptions.Telemetry.ApplicationID, clientOptions.Telemetry.Disabled))
}

// configureAuthTimeout sets the timeout for authentication requests.
func configureAuthTimeout(ctx context.Context, clientOptions *policy.ClientOptions, config *ProviderData) {
	if config.ClientOptions.TimeoutSeconds > 0 {
		clientOptions.Retry.TryTimeout = time.Duration(config.ClientOptions.TimeoutSeconds) * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Auth timeout set to %v", clientOptions.Retry.TryTimeout))
	} else {
		tflog.Debug(ctx, "No custom auth timeout configured")
	}
}

// configureAuthClientProxy creates an HTTP client for Entra ID authentication requests.
// Returns a Kiota client with custom middleware (excluding compression) for both proxy
// and non-proxy scenarios.
func configureAuthClientProxy(ctx context.Context, config *ProviderData) (*http.Client, error) {
	if config.ClientOptions.UseProxy && config.ClientOptions.ProxyURL != "" {
		return configureAuthProxyHTTPClient(ctx, config.ClientOptions)
	}

	tflog.Info(ctx, "Creating auth HTTP client without compression middleware")
	middleware := getAuthClientMiddleware(ctx, config.ClientOptions)
	return khttp.GetDefaultClient(middleware...), nil
}

// configureAuthProxyHTTPClient creates an HTTP client for authentication requests when
// proxy is configured. Uses Kiota's proxy helpers with custom middleware to exclude compression.
func configureAuthProxyHTTPClient(ctx context.Context, clientOptions *ClientOptions) (*http.Client, error) {
	tflog.Info(ctx, "Configuring proxy for auth client with URL: "+clientOptions.ProxyURL)

	middleware := getAuthClientMiddleware(ctx, clientOptions)

	var httpClient *http.Client
	var err error

	if clientOptions.ProxyUsername != "" && clientOptions.ProxyPassword != "" {
		tflog.Info(ctx, "Configuring authenticated proxy for auth client")
		httpClient, err = khttp.GetClientWithAuthenticatedProxySettings(
			clientOptions.ProxyURL,
			clientOptions.ProxyUsername,
			clientOptions.ProxyPassword,
			middleware...,
		)
	} else {
		tflog.Info(ctx, "Configuring unauthenticated proxy for auth client")
		httpClient, err = khttp.GetClientWithProxySettings(
			clientOptions.ProxyURL,
			middleware...,
		)
	}

	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Failed to create HTTP client with proxy settings: %v", err))
		return nil, fmt.Errorf("unable to create HTTP client with proxy settings: %w", err)
	}

	tflog.Debug(ctx, "Proxy settings configured successfully for auth client")
	return httpClient, nil
}

// getAuthClientMiddleware builds a custom middleware chain for authentication HTTP clients.
// By explicitly providing middleware to khttp.GetDefaultClient(), we prevent Kiota from
// adding its default middleware set (which includes CompressionHandler). This ensures
// auth requests are not gzip-compressed, as Entra ID's token endpoint does not support
// compressed request bodies.
func getAuthClientMiddleware(ctx context.Context, clientOptions *ClientOptions) []khttp.Middleware {
	var middleware []khttp.Middleware

	if clientOptions.EnableRetry {
		tflog.Debug(ctx, "Adding retry handler to auth client middleware")
		retryOptions := khttp.RetryHandlerOptions{
			MaxRetries: int(clientOptions.MaxRetries),
			ShouldRetry: func(delay time.Duration, executionCount int, req *http.Request, resp *http.Response) bool {
				if executionCount >= int(clientOptions.MaxRetries) {
					return false
				}
				if resp.StatusCode >= 500 && resp.StatusCode < 600 {
					baseDelay := time.Duration(clientOptions.RetryDelaySeconds) * time.Second
					exponentialBackoff := baseDelay * time.Duration(math.Pow(2, float64(executionCount)))
					jitter := time.Duration(rand.Int63n(int64(baseDelay)))
					delayWithJitter := exponentialBackoff + jitter
					time.Sleep(delayWithJitter)
					return true
				}
				return false
			},
		}
		middleware = append(middleware, khttp.NewRetryHandlerWithOptions(retryOptions))
	}

	if clientOptions.EnableRedirect {
		tflog.Debug(ctx, "Adding redirect handler to auth client middleware")
		redirectOptions := khttp.RedirectHandlerOptions{
			MaxRedirects: int(clientOptions.MaxRedirects),
		}
		middleware = append(middleware, khttp.NewRedirectHandlerWithOptions(redirectOptions))
	}

	if clientOptions.CustomUserAgent != "" {
		tflog.Debug(ctx, "Adding user agent handler to auth client middleware")
		userAgentOptions := khttp.NewUserAgentHandlerOptions()
		userAgentOptions.ProductName = clientOptions.CustomUserAgent
		middleware = append(middleware, khttp.NewUserAgentHandlerWithOptions(userAgentOptions))
	}

	tflog.Info(ctx, "Auth client middleware configured without compression (Entra ID token endpoint does not support gzip)")
	return middleware
}
