package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"golang.org/x/exp/rand"
)

// ConfigureGraphClientOptions configures the Graph client options based on the provided configuration
func ConfigureGraphClientOptions(ctx context.Context, config *ProviderData) (*http.Client, error) {

	// In unit test mode (TF_ACC not set), use http.DefaultClient so httpmock can intercept
	if os.Getenv("TF_ACC") == "" {
		tflog.Debug(ctx, "Unit test mode detected, using http.DefaultClient for httpmock interception")
		return http.DefaultClient, nil
	}
	tflog.Info(ctx, "Configuring Graph client options")

	defaultClientOptions := msgraphsdk.GetDefaultClientOptions()
	tflog.Debug(ctx, "Obtained default client options")

	tflog.Debug(ctx, "Initialized default middleware")
	defaultMiddleware := msgraphgocore.GetDefaultMiddlewaresWithOptions(&defaultClientOptions)

	// Customize middleware based on client options
	var err error
	defaultMiddleware, err = addChaosHandler(ctx, defaultMiddleware, config.ClientOptions)
	if err != nil {
		tflog.Error(ctx, "Failed to add chaos handler", map[string]interface{}{"error": err})
		return nil, err
	}

	tflog.Debug(ctx, "Adding custom middleware handlers")
	defaultMiddleware = addRetryHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addRedirectHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addCompressionHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addUserAgentHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addHeadersInspectionHandler(ctx, defaultMiddleware, config.ClientOptions)

	httpClient, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, defaultMiddleware)
	if err != nil {
		tflog.Error(ctx, "Failed to configure HTTP client with proxy and middleware", map[string]interface{}{"error": err})
		return nil, err
	}

	configureTimeout(ctx, httpClient, config.ClientOptions)
	tflog.Info(ctx, "Configured HTTP client timeout", map[string]interface{}{"timeoutSeconds": config.ClientOptions.TimeoutSeconds})

	tflog.Info(ctx, "Successfully configured Graph client options")
	return httpClient, nil
}

// addChaosHandler adds a chaos handler to the middleware if enabled in the options
func addChaosHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) ([]khttp.Middleware, error) {
	if options.EnableChaos {
		tflog.Debug(ctx, "Configuring chaos handler", map[string]interface{}{
			"chaosPercentage":    options.ChaosPercentage,
			"chaosStatusCode":    options.ChaosStatusCode,
			"chaosStatusMessage": options.ChaosStatusMessage,
		})

		chaosOptions := &khttp.ChaosHandlerOptions{
			ChaosStrategy:   khttp.Random,
			ChaosPercentage: int(options.ChaosPercentage),
			Headers: map[string][]string{
				"X-Chaos-Injected": {"true"},
			},
		}

		if options.ChaosStatusCode > 0 {
			chaosOptions.StatusCode = int(options.ChaosStatusCode)
		}

		if options.ChaosStatusMessage != "" {
			chaosOptions.StatusMessage = options.ChaosStatusMessage
		}

		chaosHandler, err := khttp.NewChaosHandlerWithOptions(chaosOptions)
		if err != nil {
			tflog.Debug(ctx, "Failed to create chaos handler", map[string]interface{}{"error": err})
			return nil, fmt.Errorf("failed to create chaos handler: %v", err)
		}
		middleware = append(middleware, chaosHandler)
		tflog.Debug(ctx, "Chaos handler added to middleware")
	} else {
		tflog.Debug(ctx, "Chaos handler not enabled")
	}
	return middleware, nil
}

// addRetryHandler adds a retry handler to the middleware if enabled in the options
func addRetryHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) []khttp.Middleware {
	if options.EnableRetry {
		tflog.Debug(ctx, "Configuring retry handler", map[string]interface{}{
			"maxRetries":        options.MaxRetries,
			"retryDelaySeconds": options.RetryDelaySeconds,
		})

		retryOptions := khttp.RetryHandlerOptions{
			MaxRetries: int(options.MaxRetries),
			ShouldRetry: func(delay time.Duration, executionCount int, req *http.Request, resp *http.Response) bool {
				if executionCount >= int(options.MaxRetries) {
					return false
				}
				if resp.StatusCode >= 500 && resp.StatusCode < 600 {
					baseDelay := time.Duration(options.RetryDelaySeconds) * time.Second
					exponentialBackoff := baseDelay * time.Duration(math.Pow(2, float64(executionCount)))
					jitter := time.Duration(rand.Int63n(int64(baseDelay))) // Random jitter between 0 and base delay
					delayWithJitter := exponentialBackoff + jitter

					tflog.Debug(ctx, "Retrying request", map[string]interface{}{
						"attempt":    executionCount,
						"statusCode": resp.StatusCode,
						"delay":      delayWithJitter,
						"baseDelay":  exponentialBackoff,
						"jitter":     jitter,
					})

					time.Sleep(delayWithJitter)
					return true
				}
				return false
			},
		}

		retryHandler := khttp.NewRetryHandlerWithOptions(retryOptions)
		middleware = append(middleware, retryHandler)
		tflog.Debug(ctx, "Retry handler with jitter added to middleware")
	} else {
		tflog.Debug(ctx, "Retry handler not enabled")
	}
	return middleware
}

// addRedirectHandler adds a redirect handler to the middleware if enabled in the options
func addRedirectHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) []khttp.Middleware {
	if options.EnableRedirect {
		tflog.Debug(ctx, "Configuring redirect handler", map[string]interface{}{"maxRedirects": options.MaxRedirects})
		redirectOptions := khttp.RedirectHandlerOptions{
			MaxRedirects: int(options.MaxRedirects),
			ShouldRedirect: func(req *http.Request, resp *http.Response) bool {
				return resp.StatusCode >= 300 && resp.StatusCode < 400
			},
		}
		redirectHandler := khttp.NewRedirectHandlerWithOptions(redirectOptions)
		middleware = append(middleware, redirectHandler)
		tflog.Debug(ctx, "Redirect handler added to middleware")
	} else {
		tflog.Debug(ctx, "Redirect handler not enabled")
	}
	return middleware
}

// addCompressionHandler adds a compression handler to the middleware if enabled in the options
func addCompressionHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) []khttp.Middleware {
	if options.EnableCompression {
		tflog.Debug(ctx, "Configuring compression handler")
		compressionOptions := khttp.NewCompressionOptionsReference(true)
		compressionHandler := khttp.NewCompressionHandlerWithOptions(*compressionOptions)
		middleware = append(middleware, compressionHandler)
		tflog.Debug(ctx, "Compression handler added to middleware")
	} else {
		tflog.Debug(ctx, "Compression handler not enabled")
	}
	return middleware
}

// addUserAgentHandler adds a user agent handler to the middleware if a custom user agent is specified
func addUserAgentHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) []khttp.Middleware {
	if options.CustomUserAgent != "" {
		tflog.Debug(ctx, "Configuring user agent handler", map[string]interface{}{"customUserAgent": options.CustomUserAgent})
		userAgentOptions := khttp.NewUserAgentHandlerOptions()
		userAgentOptions.ProductName = options.CustomUserAgent
		userAgentHandler := khttp.NewUserAgentHandlerWithOptions(userAgentOptions)
		middleware = append(middleware, userAgentHandler)
		tflog.Debug(ctx, "User agent handler added to middleware")
	} else {
		tflog.Debug(ctx, "Custom user agent not specified")
	}
	return middleware
}

// addHeadersInspectionHandler adds a headers inspection handler to the middleware if enabled in the options
func addHeadersInspectionHandler(ctx context.Context, middleware []khttp.Middleware, options *ClientOptions) []khttp.Middleware {
	if options.EnableHeadersInspection {
		tflog.Debug(ctx, "Configuring headers inspection handler")
		headersInspectionOptions := khttp.NewHeadersInspectionOptions()
		headersInspectionOptions.InspectRequestHeaders = true
		headersInspectionOptions.InspectResponseHeaders = true
		headersInspectionOptions.RequestHeaders = &abstractions.RequestHeaders{}
		headersInspectionOptions.ResponseHeaders = &abstractions.ResponseHeaders{}
		headersInspectionHandler := khttp.NewHeadersInspectionHandlerWithOptions(*headersInspectionOptions)
		middleware = append(middleware, headersInspectionHandler)
		tflog.Debug(ctx, "Headers inspection handler added to middleware")
	}
	return middleware
}

// configureHTTPClientWithProxyAndMiddleware creates and configures an HTTP client with proxy settings and middleware
func configureHTTPClientWithProxyAndMiddleware(ctx context.Context, config *ProviderData, middleware []khttp.Middleware) (*http.Client, error) {
	tflog.Debug(ctx, "Configuring HTTP client with proxy and middleware")
	var httpClient *http.Client
	var err error

	if config.ClientOptions.UseProxy && config.ClientOptions.ProxyURL != "" {
		if config.ClientOptions.ProxyUsername != "" && config.ClientOptions.ProxyPassword != "" {
			tflog.Debug(ctx, "Configuring authenticated proxy")
			httpClient, err = khttp.GetClientWithAuthenticatedProxySettings(
				config.ClientOptions.ProxyURL,
				config.ClientOptions.ProxyUsername,
				config.ClientOptions.ProxyPassword,
				middleware...,
			)
		} else {
			tflog.Debug(ctx, "Configuring unauthenticated proxy")
			httpClient, err = khttp.GetClientWithProxySettings(
				config.ClientOptions.ProxyURL,
				middleware...,
			)
		}
		if err != nil {
			tflog.Debug(ctx, "Failed to create HTTP client with proxy settings", map[string]interface{}{"error": err})
			return nil, fmt.Errorf("unable to create HTTP client with proxy settings: %w", err)
		}
	} else {
		tflog.Debug(ctx, "Using default HTTP client")
		httpClient = khttp.GetDefaultClient(middleware...)
	}

	tflog.Debug(ctx, "HTTP client configured successfully")
	return httpClient, nil
}

// configureTimeout sets the timeout for the HTTP client based on the options
func configureTimeout(ctx context.Context, client *http.Client, options *ClientOptions) {
	if options.TimeoutSeconds > 0 {
		client.Timeout = time.Duration(options.TimeoutSeconds) * time.Second
		tflog.Debug(ctx, "Set HTTP client timeout", map[string]interface{}{"timeoutSeconds": options.TimeoutSeconds})
	} else {
		tflog.Debug(ctx, "No custom timeout set for HTTP client")
	}
}
