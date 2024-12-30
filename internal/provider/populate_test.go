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
	entraIDSchema := map[string]attr.Type{
		"client_id":                    types.StringType,
		"client_secret":                types.StringType,
		"client_certificate":           types.StringType,
		"client_certificate_password":  types.StringType,
		"send_certificate_chain":       types.BoolType,
		"username":                     types.StringType,
		"password":                     types.StringType,
		"disable_instance_discovery":   types.BoolType,
		"additionally_allowed_tenants": types.ListType{ElemType: types.StringType},
		"redirect_url":                 types.StringType,
	}

	clientSchema := map[string]attr.Type{
		"enable_headers_inspection": types.BoolType,
		"enable_retry":              types.BoolType,
		"max_retries":               types.Int64Type,
		"retry_delay_seconds":       types.Int64Type,
		"enable_redirect":           types.BoolType,
		"max_redirects":             types.Int64Type,
		"enable_compression":        types.BoolType,
		"custom_user_agent":         types.StringType,
		"use_proxy":                 types.BoolType,
		"proxy_url":                 types.StringType,
		"proxy_username":            types.StringType,
		"proxy_password":            types.StringType,
		"timeout_seconds":           types.Int64Type,
		"enable_chaos":              types.BoolType,
		"chaos_percentage":          types.Int64Type,
		"chaos_status_code":         types.Int64Type,
		"chaos_status_message":      types.StringType,
	}

	// Empty values maps for EntraID and Client options
	// emptyEntraIDValues := map[string]attr.Value{
	// 	"client_id":                    types.StringNull(),
	// 	"client_secret":                types.StringNull(),
	// 	"client_certificate":           types.StringNull(),
	// 	"client_certificate_password":  types.StringNull(),
	// 	"send_certificate_chain":       types.BoolNull(),
	// 	"username":                     types.StringNull(),
	// 	"password":                     types.StringNull(),
	// 	"disable_instance_discovery":   types.BoolNull(),
	// 	"additionally_allowed_tenants": types.ListNull(types.StringType),
	// 	"redirect_url":                 types.StringNull(),
	// }

	emptyClientValues := map[string]attr.Value{
		"enable_headers_inspection": types.BoolNull(),
		"enable_retry":              types.BoolNull(),
		"max_retries":               types.Int64Null(),
		"retry_delay_seconds":       types.Int64Null(),
		"enable_redirect":           types.BoolNull(),
		"max_redirects":             types.Int64Null(),
		"enable_compression":        types.BoolNull(),
		"custom_user_agent":         types.StringNull(),
		"use_proxy":                 types.BoolNull(),
		"proxy_url":                 types.StringNull(),
		"proxy_username":            types.StringNull(),
		"proxy_password":            types.StringNull(),
		"timeout_seconds":           types.Int64Null(),
		"enable_chaos":              types.BoolNull(),
		"chaos_percentage":          types.Int64Null(),
		"chaos_status_code":         types.Int64Null(),
		"chaos_status_message":      types.StringNull(),
	}

	tests := []struct {
		name     string
		config   M365ProviderModel
		envVars  map[string]string
		expected M365ProviderModel
	}{
		{
			name: "All values from environment variables",
			config: M365ProviderModel{
				Cloud:           types.StringNull(),
				TenantID:        types.StringNull(),
				AuthMethod:      types.StringNull(),
				TelemetryOptout: types.BoolNull(),
				DebugMode:       types.BoolNull(),
				EntraIDOptions:  types.ObjectNull(entraIDSchema),
				ClientOptions:   types.ObjectNull(clientSchema),
			},
			envVars: map[string]string{
				"M365_CLOUD":            "public",
				"M365_TENANT_ID":        "env-tenant-id",
				"M365_AUTH_METHOD":      "client_secret",
				"M365_CLIENT_ID":        "env-client-id",
				"M365_CLIENT_SECRET":    "env-secret",
				"M365_TELEMETRY_OPTOUT": "true",
				"M365_DEBUG_MODE":       "true",
				"M365_USE_PROXY":        "true",
				"M365_PROXY_URL":        "http://proxy.com",
			},
			expected: M365ProviderModel{
				Cloud:           types.StringValue("public"),
				TenantID:        types.StringValue("env-tenant-id"),
				AuthMethod:      types.StringValue("client_secret"),
				TelemetryOptout: types.BoolValue(true),
				DebugMode:       types.BoolValue(true),
				EntraIDOptions: types.ObjectValueMust(
					entraIDSchema,
					map[string]attr.Value{
						"client_id":                    types.StringValue("env-client-id"),
						"client_secret":                types.StringValue("env-secret"),
						"client_certificate":           types.StringNull(),
						"client_certificate_password":  types.StringNull(),
						"send_certificate_chain":       types.BoolNull(),
						"username":                     types.StringNull(),
						"password":                     types.StringNull(),
						"disable_instance_discovery":   types.BoolNull(),
						"additionally_allowed_tenants": types.ListNull(types.StringType),
						"redirect_url":                 types.StringNull(),
					},
				),
				ClientOptions: types.ObjectValueMust(
					clientSchema,
					map[string]attr.Value{
						"enable_headers_inspection": types.BoolNull(),
						"enable_retry":              types.BoolNull(),
						"max_retries":               types.Int64Null(),
						"retry_delay_seconds":       types.Int64Null(),
						"enable_redirect":           types.BoolNull(),
						"max_redirects":             types.Int64Null(),
						"enable_compression":        types.BoolNull(),
						"custom_user_agent":         types.StringNull(),
						"use_proxy":                 types.BoolValue(true),
						"proxy_url":                 types.StringValue("http://proxy.com"),
						"proxy_username":            types.StringNull(),
						"proxy_password":            types.StringNull(),
						"timeout_seconds":           types.Int64Null(),
						"enable_chaos":              types.BoolNull(),
						"chaos_percentage":          types.Int64Null(),
						"chaos_status_code":         types.Int64Null(),
						"chaos_status_message":      types.StringNull(),
					},
				),
			},
		},
		{
			name: "All values from configuration",
			config: M365ProviderModel{
				Cloud:           types.StringValue("public"),
				TenantID:        types.StringValue("config-tenant-id"),
				AuthMethod:      types.StringValue("client_certificate"),
				TelemetryOptout: types.BoolValue(true),
				DebugMode:       types.BoolValue(true),
				EntraIDOptions: types.ObjectValueMust(
					entraIDSchema,
					map[string]attr.Value{
						"client_id":                    types.StringValue("config-client-id"),
						"client_secret":                types.StringNull(),
						"client_certificate":           types.StringValue("/path/to/cert.pfx"),
						"client_certificate_password":  types.StringValue("certpass"),
						"send_certificate_chain":       types.BoolValue(true),
						"username":                     types.StringNull(),
						"password":                     types.StringNull(),
						"disable_instance_discovery":   types.BoolNull(),
						"additionally_allowed_tenants": types.ListNull(types.StringType),
						"redirect_url":                 types.StringNull(),
					},
				),
				ClientOptions: types.ObjectValueMust(
					clientSchema,
					emptyClientValues,
				),
			},
			envVars: map[string]string{},
			expected: M365ProviderModel{
				Cloud:           types.StringValue("public"),
				TenantID:        types.StringValue("config-tenant-id"),
				AuthMethod:      types.StringValue("client_certificate"),
				TelemetryOptout: types.BoolValue(true),
				DebugMode:       types.BoolValue(true),
				EntraIDOptions: types.ObjectValueMust(
					entraIDSchema,
					map[string]attr.Value{
						"client_id":                    types.StringValue("config-client-id"),
						"client_secret":                types.StringNull(),
						"client_certificate":           types.StringValue("/path/to/cert.pfx"),
						"client_certificate_password":  types.StringValue("certpass"),
						"send_certificate_chain":       types.BoolValue(true),
						"username":                     types.StringNull(),
						"password":                     types.StringNull(),
						"disable_instance_discovery":   types.BoolNull(),
						"additionally_allowed_tenants": types.ListNull(types.StringType),
						"redirect_url":                 types.StringNull(),
					},
				),
				ClientOptions: types.ObjectValueMust(
					clientSchema,
					emptyClientValues,
				),
			},
		},
		{
			name: "Environment variables override configuration",
			config: M365ProviderModel{
				Cloud:           types.StringValue("public"),
				TenantID:        types.StringValue("config-tenant-id"),
				AuthMethod:      types.StringValue("client_certificate"),
				TelemetryOptout: types.BoolValue(true),
				DebugMode:       types.BoolValue(true),
				EntraIDOptions: types.ObjectValueMust(
					entraIDSchema,
					map[string]attr.Value{
						"client_id":                    types.StringValue("config-client-id"),
						"client_secret":                types.StringNull(),
						"client_certificate":           types.StringValue("/path/to/cert.pfx"),
						"client_certificate_password":  types.StringValue("certpass"),
						"send_certificate_chain":       types.BoolValue(true),
						"username":                     types.StringNull(),
						"password":                     types.StringNull(),
						"disable_instance_discovery":   types.BoolNull(),
						"additionally_allowed_tenants": types.ListNull(types.StringType),
						"redirect_url":                 types.StringNull(),
					},
				),
				ClientOptions: types.ObjectValueMust(
					clientSchema,
					emptyClientValues,
				),
			},
			envVars: map[string]string{
				"M365_CLOUD":            "gcc",
				"M365_TENANT_ID":        "env-tenant-id",
				"M365_AUTH_METHOD":      "client_secret",
				"M365_CLIENT_ID":        "env-client-id",
				"M365_CLIENT_SECRET":    "env-secret",
				"M365_TELEMETRY_OPTOUT": "false",
				"M365_DEBUG_MODE":       "false",
			},
			expected: M365ProviderModel{
				Cloud:           types.StringValue("gcc"),
				TenantID:        types.StringValue("env-tenant-id"),
				AuthMethod:      types.StringValue("client_secret"),
				TelemetryOptout: types.BoolValue(false),
				DebugMode:       types.BoolValue(false),
				EntraIDOptions: types.ObjectValueMust(
					entraIDSchema,
					map[string]attr.Value{
						"client_id":                    types.StringValue("env-client-id"),
						"client_secret":                types.StringValue("env-secret"),
						"client_certificate":           types.StringNull(),
						"client_certificate_password":  types.StringNull(),
						"send_certificate_chain":       types.BoolNull(),
						"username":                     types.StringNull(),
						"password":                     types.StringNull(),
						"disable_instance_discovery":   types.BoolNull(),
						"additionally_allowed_tenants": types.ListNull(types.StringType),
						"redirect_url":                 types.StringNull(),
					},
				),
				ClientOptions: types.ObjectValueMust(
					clientSchema,
					emptyClientValues,
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			result, diags := populateProviderData(context.Background(), tt.config)
			assert.False(t, diags.HasError(), "populateProviderData returned diagnostics")

			// Compare basic fields
			assert.Equal(t, tt.expected.Cloud, result.Cloud, "Cloud mismatch")
			assert.Equal(t, tt.expected.TenantID, result.TenantID, "TenantID mismatch")
			assert.Equal(t, tt.expected.AuthMethod, result.AuthMethod, "AuthMethod mismatch")
			assert.Equal(t, tt.expected.TelemetryOptout, result.TelemetryOptout, "TelemetryOptout mismatch")
			assert.Equal(t, tt.expected.DebugMode, result.DebugMode, "DebugMode mismatch")

			// Compare EntraIDOptions
			var expectedEntraOpts, resultEntraOpts map[string]attr.Value
			tt.expected.EntraIDOptions.As(context.Background(), &expectedEntraOpts, basetypes.ObjectAsOptions{})
			result.EntraIDOptions.As(context.Background(), &resultEntraOpts, basetypes.ObjectAsOptions{})
			assert.Equal(t, expectedEntraOpts, resultEntraOpts, "EntraIDOptions mismatch")

			// Compare ClientOptions
			var expectedClientOpts, resultClientOpts map[string]attr.Value
			tt.expected.ClientOptions.As(context.Background(), &expectedClientOpts, basetypes.ObjectAsOptions{})
			result.ClientOptions.As(context.Background(), &resultClientOpts, basetypes.ObjectAsOptions{})
			assert.Equal(t, expectedClientOpts, resultClientOpts, "ClientOptions mismatch")
		})
	}
}
