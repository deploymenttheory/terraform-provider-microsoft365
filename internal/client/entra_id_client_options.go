package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
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

func initializeAuthClientOptions(ctx context.Context, authorityURL string) policy.ClientOptions {
	options := policy.ClientOptions{
		Cloud: cloud.Configuration{
			ActiveDirectoryAuthorityHost: authorityURL,
		},
	}
	tflog.Debug(ctx, "Authentication Client options initialized with authority URL: "+authorityURL)
	return options
}

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

			for _, code := range clientOptions.Retry.StatusCodes {
				if resp.StatusCode == code {
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
			}
			return false
		},
	}

	tflog.Debug(ctx, fmt.Sprintf("Retry options set: MaxRetries=%d, BaseRetryDelay=%v",
		clientOptions.Retry.MaxRetries, time.Duration(config.ClientOptions.RetryDelaySeconds)*time.Second))
}

func configureTelemetryOptions(ctx context.Context, clientOptions *policy.ClientOptions, config *ProviderData) {
	clientOptions.Telemetry = policy.TelemetryOptions{
		ApplicationID: config.ClientOptions.CustomUserAgent,
		Disabled:      config.TelemetryOptout,
	}
	tflog.Debug(ctx, fmt.Sprintf("Telemetry options set: ApplicationID=%s, Disabled=%t",
		clientOptions.Telemetry.ApplicationID, clientOptions.Telemetry.Disabled))
}

func configureAuthTimeout(ctx context.Context, clientOptions *policy.ClientOptions, config *ProviderData) {
	if config.ClientOptions.TimeoutSeconds > 0 {
		clientOptions.Retry.TryTimeout = time.Duration(config.ClientOptions.TimeoutSeconds) * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Auth timeout set to %v", clientOptions.Retry.TryTimeout))
	} else {
		tflog.Debug(ctx, "No custom auth timeout configured")
	}
}

func configureAuthClientProxy(ctx context.Context, config *ProviderData) (*http.Client, error) {
	if config.ClientOptions.UseProxy && config.ClientOptions.ProxyURL != "" {
		return configureProxyHTTPClient(ctx, config.ClientOptions)
	}

	tflog.Info(ctx, "Using default HTTP client without proxy")
	return khttp.GetDefaultClient(), nil
}

func configureProxyHTTPClient(ctx context.Context, clientOptions *ClientOptions) (*http.Client, error) {
	tflog.Info(ctx, "Attempting to configure proxy with URL: "+clientOptions.ProxyURL)

	var httpClient *http.Client
	var err error

	if clientOptions.ProxyUsername != "" && clientOptions.ProxyPassword != "" {
		tflog.Info(ctx, "Configuring authenticated proxy")
		httpClient, err = khttp.GetClientWithAuthenticatedProxySettings(
			clientOptions.ProxyURL,
			clientOptions.ProxyUsername,
			clientOptions.ProxyPassword,
		)
	} else {
		tflog.Info(ctx, "Configuring unauthenticated proxy")
		httpClient, err = khttp.GetClientWithProxySettings(
			clientOptions.ProxyURL,
		)
	}

	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Failed to create HTTP client with proxy settings: %v", err))
		return nil, fmt.Errorf("unable to create HTTP client with proxy settings: %w", err)
	}

	tflog.Debug(ctx, "Proxy settings configured successfully")
	return httpClient, nil
}
