package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestPopulateProviderData(t *testing.T) {
	tests := []struct {
		name     string
		config   M365ProviderModel
		envVars  map[string]string
		expected M365ProviderModel
	}{
		{
			name: "Empty configuration with environment variables",
			config: M365ProviderModel{
				Cloud:      types.StringNull(),
				TenantID:   types.StringNull(),
				AuthMethod: types.StringNull(),
				EntraIDOptions: types.ObjectNull(map[string]attr.Type{
					"client_id":     types.StringType,
					"client_secret": types.StringType,
				}),
				ClientOptions: types.ObjectNull(map[string]attr.Type{
					"use_proxy":        types.BoolType,
					"proxy_url":        types.StringType,
					"enable_chaos":     types.BoolType,
					"chaos_percentage": types.Int64Type,
				}),
			},
			envVars: map[string]string{
				"M365_CLOUD":            "public",
				"M365_TENANT_ID":        "test-tenant-id",
				"M365_AUTH_METHOD":      "client_secret",
				"M365_CLIENT_ID":        "test-client-id",
				"M365_CLIENT_SECRET":    "test-client-secret",
				"M365_USE_PROXY":        "true",
				"M365_PROXY_URL":        "http://proxy.example.com",
				"M365_ENABLE_CHAOS":     "true",
				"M365_CHAOS_PERCENTAGE": "50",
			},
			expected: M365ProviderModel{
				Cloud:      types.StringValue("public"),
				TenantID:   types.StringValue("test-tenant-id"),
				AuthMethod: types.StringValue("client_secret"),
				EntraIDOptions: types.ObjectValueMust(
					map[string]attr.Type{
						"client_id":     types.StringType,
						"client_secret": types.StringType,
					},
					map[string]attr.Value{
						"client_id":     types.StringValue("test-client-id"),
						"client_secret": types.StringValue("test-client-secret"),
					},
				),
				ClientOptions: types.ObjectValueMust(
					map[string]attr.Type{
						"use_proxy":        types.BoolType,
						"proxy_url":        types.StringType,
						"enable_chaos":     types.BoolType,
						"chaos_percentage": types.Int64Type,
					},
					map[string]attr.Value{
						"use_proxy":        types.BoolValue(true),
						"proxy_url":        types.StringValue("http://proxy.example.com"),
						"enable_chaos":     types.BoolValue(true),
						"chaos_percentage": types.Int64Value(50),
					},
				),
			},
		},
		// Add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			result, diags := populateProviderData(context.Background(), tt.config)
			assert.False(t, diags.HasError(), "populateProviderData returned diagnostics")

			// Compare result with expected
			assert.Equal(t, tt.expected.Cloud, result.Cloud)
			assert.Equal(t, tt.expected.TenantID, result.TenantID)
			assert.Equal(t, tt.expected.AuthMethod, result.AuthMethod)

			// Compare EntraIDOptions
			var resultEntraIDOptions, expectedEntraIDOptions map[string]attr.Value
			tt.expected.EntraIDOptions.As(context.Background(), &expectedEntraIDOptions, basetypes.ObjectAsOptions{})
			result.EntraIDOptions.As(context.Background(), &resultEntraIDOptions, basetypes.ObjectAsOptions{})
			assert.Equal(t, expectedEntraIDOptions, resultEntraIDOptions)

			// Compare ClientOptions
			var resultClientOptions, expectedClientOptions map[string]attr.Value
			tt.expected.ClientOptions.As(context.Background(), &expectedClientOptions, basetypes.ObjectAsOptions{})
			result.ClientOptions.As(context.Background(), &resultClientOptions, basetypes.ObjectAsOptions{})
			assert.Equal(t, expectedClientOptions, resultClientOptions)
		})
	}
}
