package provider

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	khttp "github.com/microsoft/kiota-http-go"
	"golang.org/x/exp/rand"
)

func configureEntraIDClientOptions(ctx context.Context, config *M365ProviderModel, authorityURL string) (policy.ClientOptions, error) {
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

	var clientOptionsModel ClientOptionsModel
	config.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	if clientOptionsModel.TimeoutSeconds.ValueInt64() > 0 {
		httpClient.Timeout = time.Duration(clientOptionsModel.TimeoutSeconds.ValueInt64()) * time.Second
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

func configureRetryOptions(ctx context.Context, clientOptions *policy.ClientOptions, config *M365ProviderModel) {
	var clientOptionsModel ClientOptionsModel
	config.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	maxRetries := int32(clientOptionsModel.MaxRetries.ValueInt64())
	baseDelay := time.Duration(clientOptionsModel.RetryDelaySeconds.ValueInt64()) * time.Second

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
		clientOptions.Retry.MaxRetries, time.Duration(clientOptionsModel.RetryDelaySeconds.ValueInt64())*time.Second))
}

func configureTelemetryOptions(ctx context.Context, clientOptions *policy.ClientOptions, config *M365ProviderModel) {
	var clientOptionsModel ClientOptionsModel
	config.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	clientOptions.Telemetry = policy.TelemetryOptions{
		ApplicationID: clientOptionsModel.CustomUserAgent.ValueString(),
		Disabled:      config.TelemetryOptout.ValueBool(),
	}
	tflog.Debug(ctx, fmt.Sprintf("Telemetry options set: ApplicationID=%s, Disabled=%t",
		clientOptions.Telemetry.ApplicationID, clientOptions.Telemetry.Disabled))
}

func configureAuthTimeout(ctx context.Context, clientOptions *policy.ClientOptions, config *M365ProviderModel) {
	var clientOptionsModel ClientOptionsModel
	config.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	if clientOptionsModel.TimeoutSeconds.ValueInt64() > 0 {
		clientOptions.Retry.TryTimeout = time.Duration(clientOptionsModel.TimeoutSeconds.ValueInt64()) * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Auth timeout set to %v", clientOptions.Retry.TryTimeout))
	} else {
		tflog.Debug(ctx, "No custom auth timeout configured")
	}
}

func configureAuthClientProxy(ctx context.Context, config *M365ProviderModel) (*http.Client, error) {
	var clientOptionsModel ClientOptionsModel
	config.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	if clientOptionsModel.UseProxy.ValueBool() && clientOptionsModel.ProxyURL.ValueString() != "" {
		return configureProxyHTTPClient(ctx, &clientOptionsModel)
	}

	tflog.Info(ctx, "Using default HTTP client without proxy")
	return khttp.GetDefaultClient(), nil
}

func configureProxyHTTPClient(ctx context.Context, clientOptions *ClientOptionsModel) (*http.Client, error) {
	tflog.Info(ctx, "Attempting to configure proxy with URL: "+clientOptions.ProxyURL.ValueString())

	var httpClient *http.Client
	var err error

	if clientOptions.ProxyUsername.ValueString() != "" && clientOptions.ProxyPassword.ValueString() != "" {
		tflog.Info(ctx, "Configuring authenticated proxy")
		httpClient, err = khttp.GetClientWithAuthenticatedProxySettings(
			clientOptions.ProxyURL.ValueString(),
			clientOptions.ProxyUsername.ValueString(),
			clientOptions.ProxyPassword.ValueString(),
		)
	} else {
		tflog.Info(ctx, "Configuring unauthenticated proxy")
		httpClient, err = khttp.GetClientWithProxySettings(
			clientOptions.ProxyURL.ValueString(),
		)
	}

	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Failed to create HTTP client with proxy settings: %v", err))
		return nil, fmt.Errorf("unable to create HTTP client with proxy settings: %w", err)
	}

	tflog.Debug(ctx, "Proxy settings configured successfully")
	return httpClient, nil
}
