package client

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	khttp "github.com/microsoft/kiota-http-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnit_ConfigureGraphClientOptions_UnitTestMode validates that in unit test mode
// (TF_ACC not set), the function returns http.DefaultClient.
func TestUnit_ConfigureGraphClientOptions_UnitTestMode(t *testing.T) {
	os.Unsetenv("TF_ACC")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()
	config := &ProviderData{
		ClientOptions: &ClientOptions{},
	}

	client, err := ConfigureGraphClientOptions(ctx, config)
	require.NoError(t, err)
	assert.Equal(t, http.DefaultClient, client)
}

// TestUnit_AddChaosHandler_Configuration validates chaos handler configuration.
func TestUnit_AddChaosHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Chaos enabled with all options", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableChaos:        true,
			ChaosPercentage:    50,
			ChaosStatusCode:    500,
			ChaosStatusMessage: "Chaos injected",
		}

		result, err := addChaosHandler(ctx, middleware, options)
		require.NoError(t, err)
		assert.Len(t, result, 1, "Should add chaos handler")
	})

	t.Run("Chaos enabled with minimal options", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableChaos:     true,
			ChaosPercentage: 25,
		}

		result, err := addChaosHandler(ctx, middleware, options)
		require.NoError(t, err)
		assert.Len(t, result, 1, "Should add chaos handler")
	})

	t.Run("Chaos disabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableChaos: false,
		}

		result, err := addChaosHandler(ctx, middleware, options)
		require.NoError(t, err)
		assert.Len(t, result, 0, "Should not add chaos handler")
	})
}

// TestUnit_AddRetryHandler_Configuration validates retry handler configuration.
func TestUnit_AddRetryHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Retry enabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableRetry:       true,
			MaxRetries:        5,
			RetryDelaySeconds: 10,
		}

		result := addRetryHandler(ctx, middleware, options)
		assert.Len(t, result, 1, "Should add retry handler")
	})

	t.Run("Retry disabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableRetry: false,
		}

		result := addRetryHandler(ctx, middleware, options)
		assert.Len(t, result, 0, "Should not add retry handler")
	})

	t.Run("Retry with existing middleware", func(t *testing.T) {
		middleware := []khttp.Middleware{khttp.NewRedirectHandler()}
		options := &ClientOptions{
			EnableRetry:       true,
			MaxRetries:        3,
			RetryDelaySeconds: 5,
		}

		result := addRetryHandler(ctx, middleware, options)
		assert.Len(t, result, 2, "Should append retry handler to existing middleware")
	})
}

// TestUnit_AddRedirectHandler_Configuration validates redirect handler configuration.
func TestUnit_AddRedirectHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Redirect enabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableRedirect: true,
			MaxRedirects:   10,
		}

		result := addRedirectHandler(ctx, middleware, options)
		assert.Len(t, result, 1, "Should add redirect handler")
	})

	t.Run("Redirect disabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableRedirect: false,
		}

		result := addRedirectHandler(ctx, middleware, options)
		assert.Len(t, result, 0, "Should not add redirect handler")
	})
}

// TestUnit_AddCompressionHandler_Configuration validates compression handler configuration.
func TestUnit_AddCompressionHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Compression enabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableCompression: true,
		}

		result := addCompressionHandler(ctx, middleware, options)
		assert.Len(t, result, 1, "Should add compression handler")
	})

	t.Run("Compression disabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableCompression: false,
		}

		result := addCompressionHandler(ctx, middleware, options)
		assert.Len(t, result, 0, "Should not add compression handler")
	})
}

// TestUnit_AddUserAgentHandler_Configuration validates user agent handler configuration.
func TestUnit_AddUserAgentHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Custom user agent provided", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			CustomUserAgent: "TestAgent/1.0",
		}

		result := addUserAgentHandler(ctx, middleware, options)
		assert.Len(t, result, 1, "Should add user agent handler")
	})

	t.Run("No custom user agent", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			CustomUserAgent: "",
		}

		result := addUserAgentHandler(ctx, middleware, options)
		assert.Len(t, result, 0, "Should not add user agent handler")
	})
}

// TestUnit_AddHeadersInspectionHandler_Configuration validates headers inspection handler configuration.
func TestUnit_AddHeadersInspectionHandler_Configuration(t *testing.T) {
	ctx := context.Background()

	t.Run("Headers inspection enabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableHeadersInspection: true,
		}

		result := addHeadersInspectionHandler(ctx, middleware, options)
		assert.Len(t, result, 1, "Should add headers inspection handler")
	})

	t.Run("Headers inspection disabled", func(t *testing.T) {
		middleware := []khttp.Middleware{}
		options := &ClientOptions{
			EnableHeadersInspection: false,
		}

		result := addHeadersInspectionHandler(ctx, middleware, options)
		assert.Len(t, result, 0, "Should not add headers inspection handler")
	})
}

// TestConfigureHTTPClientWithProxyAndMiddleware validates proxy and middleware configuration.
func TestConfigureHTTPClientWithProxyAndMiddleware(t *testing.T) {
	ctx := context.Background()

	t.Run("No proxy configured", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy: false,
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Proxy configured but not enabled", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy: false,
				ProxyURL: "http://proxy.example.com:8080",
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Unauthenticated proxy", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy: true,
				ProxyURL: "http://proxy.example.com:8080",
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Authenticated proxy", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy:      true,
				ProxyURL:      "http://proxy.example.com:8080",
				ProxyUsername: "testuser",
				ProxyPassword: "testpass",
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Invalid proxy URL", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy: true,
				ProxyURL: "://invalid-url",
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("Invalid authenticated proxy URL", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				UseProxy:      true,
				ProxyURL:      "://invalid-url",
				ProxyUsername: "testuser",
				ProxyPassword: "testpass",
			},
		}
		middleware := []khttp.Middleware{}

		client, err := configureHTTPClientWithProxyAndMiddleware(ctx, config, middleware)
		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

// TestConfigureTimeout validates timeout configuration for HTTP clients.
func TestConfigureTimeout(t *testing.T) {
	ctx := context.Background()

	t.Run("Custom timeout configured", func(t *testing.T) {
		client := &http.Client{}
		options := &ClientOptions{
			TimeoutSeconds: 120,
		}

		configureTimeout(ctx, client, options)
		assert.Equal(t, 120*time.Second, client.Timeout)
	})

	t.Run("No custom timeout", func(t *testing.T) {
		client := &http.Client{}
		options := &ClientOptions{
			TimeoutSeconds: 0,
		}

		originalTimeout := client.Timeout
		configureTimeout(ctx, client, options)
		assert.Equal(t, originalTimeout, client.Timeout, "Timeout should remain unchanged")
	})
}

// TestConfigureGraphClientOptions_Integration validates the full ConfigureGraphClientOptions
// function with TF_ACC set (non-unit test mode).
func TestConfigureGraphClientOptions_Integration(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()

	t.Run("Basic configuration without proxy", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				EnableRetry:       true,
				MaxRetries:        3,
				RetryDelaySeconds: 5,
				EnableRedirect:    true,
				MaxRedirects:      10,
				TimeoutSeconds:    60,
			},
		}

		client, err := ConfigureGraphClientOptions(ctx, config)
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 60*time.Second, client.Timeout)
	})

	t.Run("With all middleware enabled", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				EnableRetry:             true,
				MaxRetries:              5,
				RetryDelaySeconds:       10,
				EnableRedirect:          true,
				MaxRedirects:            15,
				EnableCompression:       true,
				CustomUserAgent:         "TestAgent/1.0",
				EnableHeadersInspection: true,
				TimeoutSeconds:          120,
			},
		}

		client, err := ConfigureGraphClientOptions(ctx, config)
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, 120*time.Second, client.Timeout)
	})

	t.Run("With chaos handler", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				EnableChaos:        true,
				ChaosPercentage:    25,
				ChaosStatusCode:    500,
				ChaosStatusMessage: "Chaos test",
			},
		}

		client, err := ConfigureGraphClientOptions(ctx, config)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Invalid chaos percentage", func(t *testing.T) {
		config := &ProviderData{
			ClientOptions: &ClientOptions{
				EnableChaos:     true,
				ChaosPercentage: 150,
			},
		}

		client, err := ConfigureGraphClientOptions(ctx, config)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "failed to create chaos handler")
	})
}
