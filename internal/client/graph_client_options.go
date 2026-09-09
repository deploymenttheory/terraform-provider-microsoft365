package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// ConfigureGraphClientOptions configures the Graph client options based on the provided configuration
func ConfigureGraphClientOptions(ctx context.Context, config *ProviderData) (*http.Client, error) {
	// The unit-test harness explicitly sets TF_ACC=0 so httpmock can intercept the
	// default client. An unset TF_ACC is the normal provider runtime and must use
	// the configured Kiota middleware pipeline.
	if os.Getenv("TF_ACC") == "0" {
		tflog.Debug(
			ctx,
			"Unit test mode detected, using http.DefaultClient for httpmock interception",
		)
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
		tflog.Error(ctx, "Failed to add chaos handler", map[string]any{"error": err})
		return nil, err
	}

	tflog.Debug(ctx, "Adding custom middleware handlers")
	defaultMiddleware = addRetryHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addRedirectHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addCompressionHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addUserAgentHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = addHeadersInspectionHandler(ctx, defaultMiddleware, config.ClientOptions)
	defaultMiddleware = ensureCompressionPrecedesRetry(defaultMiddleware)

	httpClient, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, defaultMiddleware)
	if err != nil {
		tflog.Error(
			ctx,
			"Failed to configure HTTP client with proxy and middleware",
			map[string]any{"error": err},
		)
		return nil, err
	}

	configureTimeout(ctx, httpClient, config.ClientOptions)
	tflog.Info(
		ctx,
		"Configured HTTP client timeout",
		map[string]any{"timeoutSeconds": config.ClientOptions.TimeoutSeconds},
	)

	tflog.Info(ctx, "Successfully configured Graph client options")
	return httpClient, nil
}

// addChaosHandler adds a chaos handler to the middleware if enabled in the options
func addChaosHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) ([]khttp.Middleware, error) {
	if options.EnableChaos {
		tflog.Debug(ctx, "Configuring chaos handler", map[string]any{
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
			tflog.Debug(ctx, "Failed to create chaos handler", map[string]any{"error": err})
			return nil, fmt.Errorf("failed to create chaos handler: %w", err)
		}
		middleware = append(middleware, chaosHandler)
		tflog.Debug(ctx, "Chaos handler added to middleware")
	} else {
		tflog.Debug(ctx, "Chaos handler not enabled")
	}
	return middleware, nil
}

// addRetryHandler adds a retry handler to the middleware if enabled in the options
func addRetryHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) []khttp.Middleware {
	if !options.EnableRetry {
		tflog.Debug(ctx, "Retry handler not enabled")
		return removeMiddleware(middleware, func(existing khttp.Middleware) bool {
			_, ok := existing.(*khttp.RetryHandler)
			return ok
		})
	}

	tflog.Debug(ctx, "Configuring retry handler", map[string]any{
		"maxRetries":        options.MaxRetries,
		"retryDelaySeconds": options.RetryDelaySeconds,
	})

	retryOptions := khttp.RetryHandlerOptions{
		MaxRetries:   int(options.MaxRetries),
		DelaySeconds: int(options.RetryDelaySeconds),
		ShouldRetry: func(delay time.Duration, executionCount int, req *http.Request, resp *http.Response) bool {
			tflog.Debug(
				ctx,
				"Kiota retry handler accepted a retryable response",
				map[string]any{
					"attempt":          executionCount,
					"statusCode":       resp.StatusCode,
					"cumulative_delay": delay,
					"method":           req.Method,
				},
			)
			return true
		},
	}

	retryHandler := khttp.NewRetryHandlerWithOptions(retryOptions)
	for i, existing := range middleware {
		if _, ok := existing.(*khttp.RetryHandler); ok {
			middleware[i] = retryHandler
			tflog.Debug(ctx, "Configured existing Kiota retry handler")
			return middleware
		}
	}
	middleware = append(middleware, retryHandler)
	tflog.Debug(ctx, "Kiota retry handler added to middleware")
	return middleware
}

// addRedirectHandler adds a redirect handler to the middleware if enabled in the options
func addRedirectHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) []khttp.Middleware {
	if !options.EnableRedirect {
		tflog.Debug(ctx, "Redirect handler not enabled")
		return removeMiddleware(middleware, func(existing khttp.Middleware) bool {
			_, ok := existing.(*khttp.RedirectHandler)
			return ok
		})
	}

	tflog.Debug(
		ctx,
		"Configuring redirect handler",
		map[string]any{"maxRedirects": options.MaxRedirects},
	)
	redirectOptions := khttp.RedirectHandlerOptions{
		MaxRedirects: int(options.MaxRedirects),
		ShouldRedirect: func(req *http.Request, resp *http.Response) bool {
			return resp.StatusCode >= 300 && resp.StatusCode < 400
		},
	}
	redirectHandler := khttp.NewRedirectHandlerWithOptions(redirectOptions)
	for i, existing := range middleware {
		if _, ok := existing.(*khttp.RedirectHandler); ok {
			middleware[i] = redirectHandler
			tflog.Debug(ctx, "Configured existing Kiota redirect handler")
			return middleware
		}
	}
	middleware = append(middleware, redirectHandler)
	tflog.Debug(ctx, "Redirect handler added to middleware")
	return middleware
}

// addCompressionHandler adds a compression handler to the middleware if enabled in the options
func addCompressionHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) []khttp.Middleware {
	if !options.EnableCompression {
		tflog.Debug(ctx, "Compression handler not enabled")
		return removeMiddleware(middleware, func(existing khttp.Middleware) bool {
			_, ok := existing.(*khttp.CompressionHandler)
			return ok
		})
	}

	tflog.Debug(ctx, "Configuring compression handler")
	compressionOptions := khttp.NewCompressionOptionsReference(true)
	compressionHandler := khttp.NewCompressionHandlerWithOptions(*compressionOptions)
	for i, existing := range middleware {
		if _, ok := existing.(*khttp.CompressionHandler); ok {
			middleware[i] = compressionHandler
			tflog.Debug(ctx, "Configured existing Kiota compression handler")
			return middleware
		}
	}
	middleware = append(middleware, compressionHandler)
	tflog.Debug(ctx, "Compression handler added to middleware")
	return middleware
}

func removeMiddleware(
	middleware []khttp.Middleware,
	shouldRemove func(khttp.Middleware) bool,
) []khttp.Middleware {
	result := make([]khttp.Middleware, 0, len(middleware))
	for _, existing := range middleware {
		if !shouldRemove(existing) {
			result = append(result, existing)
		}
	}
	return result
}

// ensureCompressionPrecedesRetry keeps Kiota's compression handler in the retry scope.
// This lets Kiota resend the compressed request body after a throttled response.
func ensureCompressionPrecedesRetry(middleware []khttp.Middleware) []khttp.Middleware {
	compressionIndex := -1
	retryIndex := -1
	for i, handler := range middleware {
		switch handler.(type) {
		case *khttp.CompressionHandler:
			compressionIndex = i
		case *khttp.RetryHandler:
			retryIndex = i
		}
	}
	if compressionIndex < 0 || retryIndex < 0 || compressionIndex < retryIndex {
		return middleware
	}

	compressionHandler := middleware[compressionIndex]
	copy(middleware[retryIndex+1:compressionIndex+1], middleware[retryIndex:compressionIndex])
	middleware[retryIndex] = compressionHandler
	return middleware
}

// addUserAgentHandler adds a user agent handler to the middleware if a custom user agent is specified
func addUserAgentHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) []khttp.Middleware {
	if options.CustomUserAgent != "" {
		tflog.Debug(
			ctx,
			"Configuring user agent handler",
			map[string]any{"customUserAgent": options.CustomUserAgent},
		)
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
func addHeadersInspectionHandler(
	ctx context.Context,
	middleware []khttp.Middleware,
	options *ClientOptions,
) []khttp.Middleware {
	if options.EnableHeadersInspection {
		tflog.Debug(ctx, "Configuring headers inspection handler")
		headersInspectionOptions := khttp.NewHeadersInspectionOptions()
		headersInspectionOptions.InspectRequestHeaders = true
		headersInspectionOptions.InspectResponseHeaders = true
		headersInspectionOptions.RequestHeaders = &abstractions.RequestHeaders{}
		headersInspectionOptions.ResponseHeaders = &abstractions.ResponseHeaders{}
		headersInspectionHandler := khttp.NewHeadersInspectionHandlerWithOptions(
			*headersInspectionOptions,
		)
		middleware = append(middleware, headersInspectionHandler)
		tflog.Debug(ctx, "Headers inspection handler added to middleware")
	}
	return middleware
}

// configureHTTPClientWithProxyAndMiddleware creates and configures an HTTP client with proxy settings and middleware
func configureHTTPClientWithProxyAndMiddleware(
	ctx context.Context,
	config *ProviderData,
	middleware []khttp.Middleware,
) (*http.Client, error) {
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
			tflog.Debug(
				ctx,
				"Failed to create HTTP client with proxy settings",
				map[string]any{"error": err},
			)
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
		tflog.Debug(
			ctx,
			"Set HTTP client timeout",
			map[string]any{"timeoutSeconds": options.TimeoutSeconds},
		)
	} else {
		tflog.Debug(ctx, "No custom timeout set for HTTP client")
	}
}
