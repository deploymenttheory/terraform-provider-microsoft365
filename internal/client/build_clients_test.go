package client

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

// TestNewGraphClients validates the main client initialization function.
func TestNewGraphClients(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	t.Run("Invalid cloud type", func(t *testing.T) {
		ctx := context.Background()
		diags := &diag.Diagnostics{}

		config := &ProviderData{
			Cloud:    "invalid-cloud",
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			},
			ClientOptions: &ClientOptions{},
		}

		clients := NewGraphClients(ctx, config, diags)

		assert.True(t, diags.HasError(), "Should have diagnostics error for invalid cloud")
		assert.Nil(t, clients)
		assert.Contains(t, diags.Errors()[0].Summary(), "Invalid Microsoft Cloud Type")
	})

	t.Run("Missing authentication credentials", func(t *testing.T) {
		ctx := context.Background()
		diags := &diag.Diagnostics{}

		config := &ProviderData{
			Cloud:    "public",
			TenantID: "test-tenant-id",
			EntraIDOptions: &EntraIDOptions{
				ClientID: "test-client-id",
			},
			ClientOptions: &ClientOptions{},
		}

		clients := NewGraphClients(ctx, config, diags)

		assert.True(t, diags.HasError(), "Should have diagnostics error for missing credentials")
		assert.Nil(t, clients)
	})
}

// TestNewGraphClients_ConfigureEntraIDClientOptionsError validates error handling
// when Entra ID client options configuration fails.
func TestNewGraphClients_ConfigureEntraIDClientOptionsError(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()
	diags := &diag.Diagnostics{}

	config := &ProviderData{
		Cloud:    "public",
		TenantID: "test-tenant-id",
		EntraIDOptions: &EntraIDOptions{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
		ClientOptions: &ClientOptions{
			UseProxy: true,
			ProxyURL: "://invalid-proxy-url",
		},
	}

	clients := NewGraphClients(ctx, config, diags)

	assert.True(t, diags.HasError(), "Should have diagnostics error for invalid proxy URL")
	assert.Nil(t, clients)
	assert.Contains(t, diags.Errors()[0].Summary(), "Unable to configure client options")
}

// TestNewGraphClients_ConfigureGraphClientOptionsError validates error handling
// when Graph client options configuration fails.
func TestNewGraphClients_ConfigureGraphClientOptionsError(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()
	diags := &diag.Diagnostics{}

	config := &ProviderData{
		Cloud:    "public",
		TenantID: "test-tenant-id",
		EntraIDOptions: &EntraIDOptions{
			ClientID:                 "test-client-id",
			ClientSecret:             "test-client-secret",
			DisableInstanceDiscovery: true,
		},
		ClientOptions: &ClientOptions{
			UseProxy: true,
			ProxyURL: "://invalid-proxy-url",
		},
	}

	clients := NewGraphClients(ctx, config, diags)

	assert.True(t, diags.HasError(), "Should have diagnostics error for invalid proxy URL")
	assert.Nil(t, clients)
}

// TestNewGraphClients_DifferentClouds validates client creation with different cloud configurations.
func TestNewGraphClients_DifferentClouds(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()

	clouds := []string{"public", "gcc", "gcchigh", "dod", "china"}

	for _, cloud := range clouds {
		t.Run(cloud, func(t *testing.T) {
			diags := &diag.Diagnostics{}
			config := &ProviderData{
				Cloud:      cloud,
				TenantID:   "test-tenant-id",
				AuthMethod: "client_secret",
				EntraIDOptions: &EntraIDOptions{
					ClientID:                 "test-client-id",
					ClientSecret:             "test-client-secret",
					DisableInstanceDiscovery: true,
				},
				ClientOptions: &ClientOptions{
					TimeoutSeconds: 30,
				},
			}

			clients := NewGraphClients(ctx, config, diags)

			assert.NotNil(t, clients, "Clients should be created for cloud: "+cloud)
			if clients != nil {
				assert.NotNil(t, clients.GetKiotaGraphV1Client())
				assert.NotNil(t, clients.GetKiotaGraphBetaClient())
				assert.NotNil(t, clients.GetGraphV1Client())
				assert.NotNil(t, clients.GetGraphBetaClient())
			}
		})
	}
}

// TestNewGraphClients_WithClientOptions validates client creation with various client options.
func TestNewGraphClients_WithClientOptions(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()

	tests := []struct {
		name          string
		clientOptions *ClientOptions
	}{
		{
			name: "With retry enabled",
			clientOptions: &ClientOptions{
				EnableRetry:       true,
				MaxRetries:        5,
				RetryDelaySeconds: 10,
				TimeoutSeconds:    60,
			},
		},
		{
			name: "With redirect enabled",
			clientOptions: &ClientOptions{
				EnableRedirect: true,
				MaxRedirects:   10,
				TimeoutSeconds: 60,
			},
		},
		{
			name: "With compression enabled",
			clientOptions: &ClientOptions{
				EnableCompression: true,
				TimeoutSeconds:    60,
			},
		},
		{
			name: "With custom user agent",
			clientOptions: &ClientOptions{
				CustomUserAgent: "TestAgent/1.0",
				TimeoutSeconds:  60,
			},
		},
		{
			name: "With headers inspection",
			clientOptions: &ClientOptions{
				EnableHeadersInspection: true,
				TimeoutSeconds:          60,
			},
		},
		{
			name: "With all options enabled",
			clientOptions: &ClientOptions{
				EnableRetry:             true,
				MaxRetries:              3,
				RetryDelaySeconds:       5,
				EnableRedirect:          true,
				MaxRedirects:            10,
				EnableCompression:       true,
				CustomUserAgent:         "TestAgent/1.0",
				EnableHeadersInspection: true,
				TimeoutSeconds:          120,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := &diag.Diagnostics{}
			config := &ProviderData{
				Cloud:      "public",
				TenantID:   "test-tenant-id",
				AuthMethod: "client_secret",
				EntraIDOptions: &EntraIDOptions{
					ClientID:                 "test-client-id",
					ClientSecret:             "test-client-secret",
					DisableInstanceDiscovery: true,
				},
				ClientOptions: tt.clientOptions,
			}

			clients := NewGraphClients(ctx, config, diags)

			assert.NotNil(t, clients, "Clients should be created")
			if clients != nil {
				assert.NotNil(t, clients.GetKiotaGraphV1Client())
				assert.NotNil(t, clients.GetKiotaGraphBetaClient())
				assert.NotNil(t, clients.GetGraphV1Client())
				assert.NotNil(t, clients.GetGraphBetaClient())
			}
		})
	}
}

// TestNewGraphClients_SuccessfulIntegration validates successful client creation
// with various authentication methods.
func TestNewGraphClients_SuccessfulIntegration(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	defer os.Unsetenv("TF_ACC")

	ctx := context.Background()

	tests := []struct {
		name       string
		authMethod string
		setupEnv   func()
		cleanupEnv func()
		config     *EntraIDOptions
	}{
		{
			name:       "Azure CLI",
			authMethod: "azure_cli",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Azure Developer CLI",
			authMethod: "azure_developer_cli",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Managed Identity with Client ID",
			authMethod: "managed_identity",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ManagedIdentityClientID:  "test-client-id",
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Managed Identity with Resource ID",
			authMethod: "managed_identity",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ManagedIdentityResourceID: "/subscriptions/test/resourceGroups/test/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test",
				DisableInstanceDiscovery:  true,
			},
		},
		{
			name:       "Workload Identity",
			authMethod: "workload_identity",
			setupEnv: func() {
				os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/tmp/token")
			},
			cleanupEnv: func() {
				os.Unsetenv("AZURE_FEDERATED_TOKEN_FILE")
			},
			config: &EntraIDOptions{
				ClientID:                 "test-client-id",
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Client Secret",
			authMethod: "client_secret",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:                 "test-client-id",
				ClientSecret:             "test-secret",
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Device Code",
			authMethod: "device_code",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:                 "test-client-id",
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Interactive Browser",
			authMethod: "interactive_browser",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:                 "test-client-id",
				DisableInstanceDiscovery: true,
			},
		},
		{
			name:       "Username Password",
			authMethod: "username_password",
			setupEnv:   func() {},
			cleanupEnv: func() {},
			config: &EntraIDOptions{
				ClientID:                 "test-client-id",
				Username:                 "test@example.com",
				Password:                 "test-password",
				DisableInstanceDiscovery: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			diags := &diag.Diagnostics{}
			config := &ProviderData{
				Cloud:          "public",
				TenantID:       "test-tenant-id",
				AuthMethod:     tt.authMethod,
				EntraIDOptions: tt.config,
				ClientOptions: &ClientOptions{
					TimeoutSeconds: 30,
				},
			}

			clients := NewGraphClients(ctx, config, diags)

			if diags.HasError() {
				t.Logf("Diagnostics errors: %v", diags.Errors())
			}

			assert.NotNil(t, clients, "Clients should be created")
			if clients != nil {
				assert.NotNil(t, clients.GetKiotaGraphV1Client())
				assert.NotNil(t, clients.GetKiotaGraphBetaClient())
				assert.NotNil(t, clients.GetGraphV1Client())
				assert.NotNil(t, clients.GetGraphBetaClient())
			}
		})
	}
}
