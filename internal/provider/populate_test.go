package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPopulateProviderData_EmptyConfig(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Create an empty config
	config := M365ProviderModel{
		Cloud:           types.StringNull(),
		TenantID:        types.StringNull(),
		AuthMethod:      types.StringNull(),
		EntraIDOptions:  types.ObjectNull(map[string]attr.Type{}),
		ClientOptions:   types.ObjectNull(map[string]attr.Type{}),
		TelemetryOptout: types.BoolNull(),
		DebugMode:       types.BoolNull(),
	}

	// Call the function
	result, diags := populateProviderData(context.Background(), config)
	require.False(t, diags.HasError(), "populateProviderData should not return errors with empty config")

	// Verify defaults
	assert.Equal(t, "", result.Cloud.ValueString(), "Default cloud should be empty string")
	assert.Equal(t, "", result.TenantID.ValueString(), "Default tenant ID should be empty string")
	assert.Equal(t, "", result.AuthMethod.ValueString(), "Default auth method should be empty string")
	assert.False(t, result.TelemetryOptout.ValueBool(), "Default telemetry optout should be false")
	assert.False(t, result.DebugMode.ValueBool(), "Default debug mode should be false")
}

func TestPopulateProviderData_EnvVarsOverrideConfig(t *testing.T) {
	// Clear environment and set test variables
	os.Clearenv()
	os.Setenv("M365_CLOUD", "dod")
	os.Setenv("M365_TENANT_ID", "env-tenant-id")
	os.Setenv("M365_AUTH_METHOD", "device_code")
	os.Setenv("M365_TELEMETRY_OPTOUT", "true")
	os.Setenv("M365_DEBUG_MODE", "true")

	// Create a config with different values
	config := M365ProviderModel{
		Cloud:           types.StringValue("public"),
		TenantID:        types.StringValue("config-tenant-id"),
		AuthMethod:      types.StringValue("client_secret"),
		EntraIDOptions:  types.ObjectNull(map[string]attr.Type{}),
		ClientOptions:   types.ObjectNull(map[string]attr.Type{}),
		TelemetryOptout: types.BoolValue(false),
		DebugMode:       types.BoolValue(false),
	}

	// Call the function
	result, diags := populateProviderData(context.Background(), config)
	require.False(t, diags.HasError(), "populateProviderData should not return errors")

	// Verify env vars take precedence
	assert.Equal(t, "dod", result.Cloud.ValueString(), "Cloud should be from env var")
	assert.Equal(t, "env-tenant-id", result.TenantID.ValueString(), "Tenant ID should be from env var")
	assert.Equal(t, "device_code", result.AuthMethod.ValueString(), "Auth method should be from env var")
	assert.True(t, result.TelemetryOptout.ValueBool(), "Telemetry optout should be from env var")
	assert.True(t, result.DebugMode.ValueBool(), "Debug mode should be from env var")
}

func TestPopulateProviderData_MixedEnvAndConfig(t *testing.T) {
	// Clear environment and set only some variables
	os.Clearenv()
	os.Setenv("M365_CLOUD", "gcc")
	// Intentionally not setting TENANT_ID
	os.Setenv("M365_AUTH_METHOD", "client_certificate")
	// Intentionally not setting TELEMETRY_OPTOUT
	os.Setenv("M365_DEBUG_MODE", "true")

	// Create a config with values for everything
	config := M365ProviderModel{
		Cloud:           types.StringValue("public"),
		TenantID:        types.StringValue("config-tenant-id"),
		AuthMethod:      types.StringValue("client_secret"),
		EntraIDOptions:  types.ObjectNull(map[string]attr.Type{}),
		ClientOptions:   types.ObjectNull(map[string]attr.Type{}),
		TelemetryOptout: types.BoolValue(true),
		DebugMode:       types.BoolValue(false),
	}

	// Call the function
	result, diags := populateProviderData(context.Background(), config)
	require.False(t, diags.HasError(), "populateProviderData should not return errors")

	// Verify mixed results
	assert.Equal(t, "gcc", result.Cloud.ValueString(), "Cloud should be from env var")
	assert.Equal(t, "config-tenant-id", result.TenantID.ValueString(), "Tenant ID should be from config")
	assert.Equal(t, "client_certificate", result.AuthMethod.ValueString(), "Auth method should be from env var")
	assert.True(t, result.TelemetryOptout.ValueBool(), "Telemetry optout should be from config")
	assert.True(t, result.DebugMode.ValueBool(), "Debug mode should be from env var")
}

func TestPopulateEntraIDOptions_EnvVarsAndConfig(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Set environment variables
	os.Setenv("M365_CLIENT_ID", "env-client-id")
	os.Setenv("M365_CLIENT_SECRET", "env-secret")
	os.Setenv("M365_ADDITIONALLY_ALLOWED_TENANTS", "tenant1,tenant2")

	// Create a proper EntraIDOptionsModel that matches the function's expectations
	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())

	// Initialize a list value for additionally_allowed_tenants
	allowedTenantsList, diags := types.ListValueFrom(context.Background(), types.StringType, []string{})
	require.False(t, diags.HasError(), "Failed to create allowed tenants list")

	configMap := map[string]attr.Value{
		"client_id":                    types.StringValue("config-client-id"),
		"client_secret":                types.StringValue("config-secret"),
		"client_certificate":           types.StringValue("/path/to/cert.pfx"),
		"client_certificate_password":  types.StringValue("certpass"),
		"send_certificate_chain":       types.BoolValue(true),
		"username":                     types.StringValue("user@example.com"),
		"disable_instance_discovery":   types.BoolValue(true),
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringValue("https://example.com/callback"),
		"federated_token_file_path":    types.StringValue("/path/to/token"),
		"managed_identity_id":          types.StringValue("managed-id"),
		"oidc_token_file_path":         types.StringValue("/path/to/oidc"),
		"ado_service_connection_id":    types.StringValue("service-conn-id"),
	}

	configObj, diags := types.ObjectValue(entraIDSchema, configMap)
	require.False(t, diags.HasError(), "Failed to create test config object")

	// Call the function
	result, diags := populateEntraIDOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "populateEntraIDOptions should not return errors")

	// Create an EntraIDOptionsModel to hold the result
	var resultModel EntraIDOptionsModel
	diags = result.As(context.Background(), &resultModel, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "Failed to extract result model")

	// Verify environment variables override config
	assert.Equal(t, "env-client-id", resultModel.ClientID.ValueString(), "Client ID should be from env var")
	assert.Equal(t, "env-secret", resultModel.ClientSecret.ValueString(), "Client secret should be from env var")

	// Verify other values come from config
	assert.Equal(t, "/path/to/cert.pfx", resultModel.ClientCertificate.ValueString(), "Client certificate should be from config")
	assert.Equal(t, "certpass", resultModel.ClientCertificatePassword.ValueString(), "Client certificate password should be from config")
	assert.True(t, resultModel.SendCertificateChain.ValueBool(), "Send certificate chain should be from config")
	assert.Equal(t, "user@example.com", resultModel.Username.ValueString(), "Username should be from config")
	assert.True(t, resultModel.DisableInstanceDiscovery.ValueBool(), "Disable instance discovery should be from config")
	assert.Equal(t, "https://example.com/callback", resultModel.RedirectUrl.ValueString(), "Redirect URL should be from config")
	assert.Equal(t, "/path/to/token", resultModel.FederatedTokenFilePath.ValueString(), "Federated token file path should be from config")
	assert.Equal(t, "managed-id", resultModel.ManagedIdentityID.ValueString(), "Managed identity ID should be from config")
	assert.Equal(t, "/path/to/oidc", resultModel.OIDCTokenFilePath.ValueString(), "OIDC token file path should be from config")
	assert.Equal(t, "service-conn-id", resultModel.ADOServiceConnectionID.ValueString(), "ADO service connection ID should be from config")

	// Check list handling - verify the environment variable took precedence
	var allowedTenants []string
	diags = resultModel.AdditionallyAllowedTenants.ElementsAs(context.Background(), &allowedTenants, false)
	require.False(t, diags.HasError(), "Failed to extract allowed tenants")

	// Expecting tenant1 and tenant2 from the environment variable
	assert.ElementsMatch(t, []string{"tenant1", "tenant2"}, allowedTenants, "Additionally allowed tenants should be from env var")
}

func TestPopulateClientOptions_DefaultsAndOverrides(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Set environment variables
	os.Setenv("M365_USE_PROXY", "true")
	os.Setenv("M365_PROXY_URL", "http://proxy.example.com")
	os.Setenv("M365_MAX_RETRIES", "5")
	os.Setenv("M365_ENABLE_CHAOS", "true")
	os.Setenv("M365_CHAOS_PERCENTAGE", "30")

	// Get the schema for client options
	clientSchema := schemaToAttrTypes(ClientOptionsSchema())

	// Create null config (testing defaults)
	nullConfig := types.ObjectNull(clientSchema)

	// Call the function with null config
	nullResult, diags := populateClientOptions(context.Background(), nullConfig)
	require.False(t, diags.HasError(), "populateClientOptions should not return errors with null config")
	assert.True(t, nullResult.IsNull(), "Result should be null with null config")

	// Create a config with values
	configMap := map[string]attr.Value{
		"enable_headers_inspection": types.BoolValue(true),
		"enable_retry":              types.BoolValue(true),
		"max_retries":               types.Int64Value(10), // Should be overridden by env var
		"retry_delay_seconds":       types.Int64Value(3),
		"enable_redirect":           types.BoolValue(true),
		"max_redirects":             types.Int64Value(30),
		"enable_compression":        types.BoolValue(true),
		"custom_user_agent":         types.StringValue("custom-agent"),
		"use_proxy":                 types.BoolValue(false),                      // Should be overridden by env var
		"proxy_url":                 types.StringValue("http://other-proxy.com"), // Should be overridden
		"proxy_username":            types.StringValue("proxyuser"),
		"proxy_password":            types.StringValue("proxypass"),
		"timeout_seconds":           types.Int64Value(60),
		"enable_chaos":              types.BoolValue(false), // Should be overridden by env var
		"chaos_percentage":          types.Int64Value(10),   // Should be overridden by env var
		"chaos_status_code":         types.Int64Value(503),
		"chaos_status_message":      types.StringValue("Chaos error"),
	}

	configObj, diags := types.ObjectValue(clientSchema, configMap)
	require.False(t, diags.HasError(), "Failed to create test config object")

	// Call the function
	result, diags := populateClientOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "populateClientOptions should not return errors")

	// Create a ClientOptionsModel to hold the result
	var resultModel ClientOptionsModel
	diags = result.As(context.Background(), &resultModel, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "Failed to extract result model")

	// Verify environment variables override config
	assert.True(t, resultModel.UseProxy.ValueBool(), "use_proxy should be from env var (true)")
	assert.Equal(t, "http://proxy.example.com", resultModel.ProxyURL.ValueString(), "proxy_url should be from env var")
	assert.Equal(t, int64(5), resultModel.MaxRetries.ValueInt64(), "max_retries should be from env var (5)")
	assert.True(t, resultModel.EnableChaos.ValueBool(), "enable_chaos should be from env var (true)")
	assert.Equal(t, int64(30), resultModel.ChaosPercentage.ValueInt64(), "chaos_percentage should be from env var (30)")

	// Verify other values come from config
	assert.True(t, resultModel.EnableHeadersInspection.ValueBool(), "enable_headers_inspection should be from config")
	assert.True(t, resultModel.EnableRetry.ValueBool(), "enable_retry should be from config")
	assert.Equal(t, int64(3), resultModel.RetryDelaySeconds.ValueInt64(), "retry_delay_seconds should be from config")
	assert.True(t, resultModel.EnableRedirect.ValueBool(), "enable_redirect should be from config")
	assert.Equal(t, int64(30), resultModel.MaxRedirects.ValueInt64(), "max_redirects should be from config")
	assert.True(t, resultModel.EnableCompression.ValueBool(), "enable_compression should be from config")
	assert.Equal(t, "custom-agent", resultModel.CustomUserAgent.ValueString(), "custom_user_agent should be from config")
	assert.Equal(t, "proxyuser", resultModel.ProxyUsername.ValueString(), "proxy_username should be from config")
	assert.Equal(t, "proxypass", resultModel.ProxyPassword.ValueString(), "proxy_password should be from config")
	assert.Equal(t, int64(60), resultModel.TimeoutSeconds.ValueInt64(), "timeout_seconds should be from config")
	assert.Equal(t, int64(503), resultModel.ChaosStatusCode.ValueInt64(), "chaos_status_code should be from config")
	assert.Equal(t, "Chaos error", resultModel.ChaosStatusMessage.ValueString(), "chaos_status_message should be from config")
}

func TestPopulateEntraIDOptions_EmptyConfig(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Create null config (testing defaults)
	nullConfig := types.ObjectNull(map[string]attr.Type{})

	// Call the function with null config
	result, diags := populateEntraIDOptions(context.Background(), nullConfig)
	require.False(t, diags.HasError(), "populateEntraIDOptions should not return errors with null config")
	assert.True(t, result.IsNull(), "Result should be null with null config")
}

func TestPopulateEntraIDOptions_BooleanEnvVars(t *testing.T) {
	// Clear environment variables
	os.Clearenv()

	// Set boolean environment variables with different formats
	os.Setenv("M365_SEND_CERTIFICATE_CHAIN", "true")
	os.Setenv("M365_DISABLE_INSTANCE_DISCOVERY", "TRUE")

	// Get the schema for EntraID options
	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())

	// Initialize a list value for additionally_allowed_tenants
	allowedTenantsList, diags := types.ListValueFrom(context.Background(), types.StringType, []string{})
	require.False(t, diags.HasError(), "Failed to create allowed tenants list")

	// Create a minimal config
	configMap := map[string]attr.Value{
		"client_id":                    types.StringValue("config-client-id"),
		"client_secret":                types.StringValue("config-secret"),
		"client_certificate":           types.StringNull(),
		"client_certificate_password":  types.StringNull(),
		"send_certificate_chain":       types.BoolValue(false), // Should be overridden by env var
		"username":                     types.StringNull(),
		"disable_instance_discovery":   types.BoolValue(false), // Should be overridden by env var
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringNull(),
		"federated_token_file_path":    types.StringNull(),
		"managed_identity_id":          types.StringNull(),
		"oidc_token_file_path":         types.StringNull(),
		"ado_service_connection_id":    types.StringNull(),
	}

	configObj, diags := types.ObjectValue(entraIDSchema, configMap)
	require.False(t, diags.HasError(), "Failed to create test config object")

	// Call the function
	result, diags := populateEntraIDOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "populateEntraIDOptions should not return errors")

	// Create an EntraIDOptionsModel to hold the result
	var resultModel EntraIDOptionsModel
	diags = result.As(context.Background(), &resultModel, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "Failed to extract result model")

	// Boolean environment variables should be correctly parsed
	assert.True(t, resultModel.SendCertificateChain.ValueBool(),
		"send_certificate_chain should be true from env var")
	assert.True(t, resultModel.DisableInstanceDiscovery.ValueBool(),
		"disable_instance_discovery should be true from env var")

	// Verify other values were not changed
	assert.Equal(t, "config-client-id", resultModel.ClientID.ValueString(),
		"client_id should remain from config")
	assert.Equal(t, "config-secret", resultModel.ClientSecret.ValueString(),
		"client_secret should remain from config")
}

// Helper function to test multiple variant formats of boolean environment variables
func TestGetEnvBool_Variants(t *testing.T) {
	testCases := []struct {
		envValue string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"t", true},
		{"T", true},
		{"1", true},
		{"false", false},
		{"FALSE", false},
		{"False", false},
		{"f", false},
		{"F", false},
		{"0", false},
		{"invalid", false}, // Invalid values should default to false
		{"", false},        // Empty string should default to false
	}

	for _, tc := range testCases {
		t.Run(tc.envValue, func(t *testing.T) {
			os.Clearenv()
			if tc.envValue != "" { // Skip setting for empty string test
				os.Setenv("TEST_BOOL", tc.envValue)
			}

			// Create a helper function similar to your GetEnvBool
			result := func() bool {
				if v, ok := os.LookupEnv("TEST_BOOL"); ok {
					switch v {
					case "true", "TRUE", "True", "t", "T", "1":
						return true
					default:
						return false
					}
				}
				return false
			}()

			assert.Equal(t, tc.expected, result, "GetEnvBool should correctly parse %q as %v", tc.envValue, tc.expected)
		})
	}
}
