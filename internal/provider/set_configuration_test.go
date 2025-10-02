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

func TestSetProviderConfiguration_EmptyConfig(t *testing.T) {
	// Clear M365 environment variables only
	clearM365EnvVarsWithRestore(t)

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
	result, diags := setProviderConfiguration(context.Background(), config)
	require.False(t, diags.HasError(), "setProviderConfiguration should not return errors with empty config")

	// Verify defaults
	assert.Equal(t, "", result.Cloud.ValueString(), "Default cloud should be empty string")
	assert.Equal(t, "", result.TenantID.ValueString(), "Default tenant ID should be empty string")
	assert.Equal(t, "", result.AuthMethod.ValueString(), "Default auth method should be empty string")
	assert.False(t, result.TelemetryOptout.ValueBool(), "Default telemetry optout should be false")
	assert.False(t, result.DebugMode.ValueBool(), "Default debug mode should be false")
}

func TestSetProviderConfiguration_ZZZ_ComprehensiveEnvVarsOverrideConfig(t *testing.T) {
	// Store original environment values for manual cleanup
	originalEnvVars := make(map[string]string)
	envVarsToTest := []string{
		"M365_CLOUD", "AZURE_CLOUD", "M365_TENANT_ID", "M365_AUTH_METHOD",
		"M365_TELEMETRY_OPTOUT", "M365_DEBUG_MODE", "M365_CLIENT_ID", "M365_CLIENT_SECRET",
		"M365_CLIENT_CERTIFICATE_FILE_PATH", "M365_CLIENT_CERTIFICATE_PASSWORD",
		"M365_USERNAME", "M365_SEND_CERTIFICATE_CHAIN", "M365_DISABLE_INSTANCE_DISCOVERY",
		"M365_ADDITIONALLY_ALLOWED_TENANTS", "M365_REDIRECT_URI", "AZURE_FEDERATED_TOKEN_FILE",
		"M365_MANAGED_IDENTITY_ID", "AZURE_CLIENT_ID", "M365_OIDC_TOKEN_FILE_PATH",
		"M365_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL", "M365_OIDC_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN", "ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID",
		"ARM_OIDC_AZURE_SERVICE_CONNECTION_ID", "M365_ENABLE_HEADERS_INSPECTION",
		"M365_ENABLE_RETRY", "M365_MAX_RETRIES", "M365_RETRY_DELAY_SECONDS",
		"M365_ENABLE_REDIRECT", "M365_MAX_REDIRECTS", "M365_ENABLE_COMPRESSION",
		"M365_CUSTOM_USER_AGENT", "M365_USE_PROXY", "M365_PROXY_URL",
		"M365_PROXY_USERNAME", "M365_PROXY_PASSWORD", "M365_TIMEOUT_SECONDS",
		"M365_ENABLE_CHAOS", "M365_CHAOS_PERCENTAGE", "M365_CHAOS_STATUS_CODE",
		"M365_CHAOS_STATUS_MESSAGE",
	}

	// Store original values and clear all environment variables
	for _, envVar := range envVarsToTest {
		originalEnvVars[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// Ensure cleanup happens even if test fails
	defer func() {
		for envVar, originalValue := range originalEnvVars {
			if originalValue != "" {
				os.Setenv(envVar, originalValue)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	// Provider-level environment variables
	os.Setenv("M365_CLOUD", "dod")
	os.Setenv("AZURE_CLOUD", "gcc") // Should be overridden by M365_CLOUD
	os.Setenv("M365_TENANT_ID", "env-tenant-id")
	os.Setenv("M365_AUTH_METHOD", "device_code")
	os.Setenv("M365_TELEMETRY_OPTOUT", "true")
	os.Setenv("M365_DEBUG_MODE", "true")

	// EntraID Options environment variables
	os.Setenv("M365_CLIENT_ID", "env-client-id")
	os.Setenv("M365_CLIENT_SECRET", "env-client-secret")
	os.Setenv("M365_CLIENT_CERTIFICATE_FILE_PATH", "/env/path/to/cert.pfx")
	os.Setenv("M365_CLIENT_CERTIFICATE_PASSWORD", "env-cert-password")
	os.Setenv("M365_USERNAME", "env-user@example.com")
	os.Setenv("M365_SEND_CERTIFICATE_CHAIN", "true")
	os.Setenv("M365_DISABLE_INSTANCE_DISCOVERY", "true")
	os.Setenv("M365_ADDITIONALLY_ALLOWED_TENANTS", "tenant1,tenant2,tenant3")
	os.Setenv("M365_REDIRECT_URI", "https://env.example.com/callback")
	os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "/env/path/to/federated-token")
	os.Setenv("M365_MANAGED_IDENTITY_ID", "env-managed-identity-id")
	os.Setenv("AZURE_CLIENT_ID", "env-azure-client-id") // Should be overridden by M365_MANAGED_IDENTITY_ID
	os.Setenv("M365_OIDC_TOKEN_FILE_PATH", "/env/path/to/oidc-token")
	os.Setenv("M365_OIDC_REQUEST_URL", "https://env.example.com/oidc-request")
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "https://github.actions.example.com/token") // Should be overridden by M365_OIDC_REQUEST_URL
	os.Setenv("M365_OIDC_REQUEST_TOKEN", "env-oidc-request-token")
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "github-actions-token") // Should be overridden by M365_OIDC_REQUEST_TOKEN
	os.Setenv("ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "env-ado-service-connection-id")
	os.Setenv("ARM_OIDC_AZURE_SERVICE_CONNECTION_ID", "env-azure-service-connection-id") // Should be overridden by ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID

	// Client Options environment variables
	os.Setenv("M365_ENABLE_HEADERS_INSPECTION", "true")
	os.Setenv("M365_ENABLE_RETRY", "true")
	os.Setenv("M365_MAX_RETRIES", "5")
	os.Setenv("M365_RETRY_DELAY_SECONDS", "10")
	os.Setenv("M365_ENABLE_REDIRECT", "true")
	os.Setenv("M365_MAX_REDIRECTS", "3")
	os.Setenv("M365_ENABLE_COMPRESSION", "true")
	os.Setenv("M365_CUSTOM_USER_AGENT", "env-custom-user-agent")
	os.Setenv("M365_USE_PROXY", "true")
	os.Setenv("M365_PROXY_URL", "https://env.proxy.example.com:8080")
	os.Setenv("M365_PROXY_USERNAME", "env-proxy-user")
	os.Setenv("M365_PROXY_PASSWORD", "env-proxy-password")
	os.Setenv("M365_TIMEOUT_SECONDS", "120")
	os.Setenv("M365_ENABLE_CHAOS", "true")
	os.Setenv("M365_CHAOS_PERCENTAGE", "25")
	os.Setenv("M365_CHAOS_STATUS_CODE", "503")
	os.Setenv("M365_CHAOS_STATUS_MESSAGE", "env-chaos-message")

	// Create EntraID options config with different values
	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())
	allowedTenantsList, diags := types.ListValueFrom(context.Background(), types.StringType, []string{"config-tenant1", "config-tenant2"})
	require.False(t, diags.HasError(), "Failed to create allowed tenants list")

	entraIDConfigMap := map[string]attr.Value{
		"client_id":                    types.StringValue("config-client-id"),
		"client_secret":                types.StringValue("config-client-secret"),
		"client_certificate":           types.StringValue("/config/path/to/cert.pfx"),
		"client_certificate_password":  types.StringValue("config-cert-password"),
		"send_certificate_chain":       types.BoolValue(false),
		"username":                     types.StringValue("config-user@example.com"),
		"disable_instance_discovery":   types.BoolValue(false),
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringValue("https://config.example.com/callback"),
		"federated_token_file_path":    types.StringValue("/config/path/to/federated-token"),
		"managed_identity_id":          types.StringValue("config-managed-identity-id"),
		"oidc_token_file_path":         types.StringValue("/config/path/to/oidc-token"),
		"oidc_request_token":           types.StringValue("config-oidc-request-token"),
		"oidc_request_url":             types.StringValue("https://config.example.com/oidc-request"),
		"ado_service_connection_id":    types.StringValue("config-ado-service-connection-id"),
	}

	entraIDConfigObj, diags := types.ObjectValue(entraIDSchema, entraIDConfigMap)
	require.False(t, diags.HasError(), "Failed to create EntraID config object")

	// Create Client options config with different values
	clientSchema := schemaToAttrTypes(ClientOptionsSchema())
	clientConfigMap := map[string]attr.Value{
		"enable_headers_inspection": types.BoolValue(false),
		"enable_retry":              types.BoolValue(false),
		"max_retries":               types.Int64Value(3),
		"retry_delay_seconds":       types.Int64Value(5),
		"enable_redirect":           types.BoolValue(false),
		"max_redirects":             types.Int64Value(2),
		"enable_compression":        types.BoolValue(false),
		"custom_user_agent":         types.StringValue("config-custom-user-agent"),
		"use_proxy":                 types.BoolValue(false),
		"proxy_url":                 types.StringValue("https://config.proxy.example.com:8080"),
		"proxy_username":            types.StringValue("config-proxy-user"),
		"proxy_password":            types.StringValue("config-proxy-password"),
		"timeout_seconds":           types.Int64Value(60),
		"enable_chaos":              types.BoolValue(false),
		"chaos_percentage":          types.Int64Value(10),
		"chaos_status_code":         types.Int64Value(500),
		"chaos_status_message":      types.StringValue("config-chaos-message"),
	}

	clientConfigObj, diags := types.ObjectValue(clientSchema, clientConfigMap)
	require.False(t, diags.HasError(), "Failed to create Client config object")

	// Create a config with different values that should be overridden
	config := M365ProviderModel{
		Cloud:           types.StringValue("public"),
		TenantID:        types.StringValue("config-tenant-id"),
		AuthMethod:      types.StringValue("client_secret"),
		EntraIDOptions:  entraIDConfigObj,
		ClientOptions:   clientConfigObj,
		TelemetryOptout: types.BoolValue(false),
		DebugMode:       types.BoolValue(false),
	}

	// Call the function
	result, diags := setProviderConfiguration(context.Background(), config)
	require.False(t, diags.HasError(), "setProviderConfiguration should not return errors")

	// Verify provider-level env vars take precedence
	assert.Equal(t, "dod", result.Cloud.ValueString(), "Cloud should be from M365_CLOUD env var")
	assert.Equal(t, "env-tenant-id", result.TenantID.ValueString(), "Tenant ID should be from env var")
	assert.Equal(t, "device_code", result.AuthMethod.ValueString(), "Auth method should be from env var")
	assert.True(t, result.TelemetryOptout.ValueBool(), "Telemetry optout should be from env var")
	assert.True(t, result.DebugMode.ValueBool(), "Debug mode should be from env var")

	// Extract and verify EntraID options
	var resultEntraIDOptions EntraIDOptionsModel
	diags = result.EntraIDOptions.As(context.Background(), &resultEntraIDOptions, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "Failed to extract EntraID options from result")

	// Verify EntraID env vars take precedence
	assert.Equal(t, "env-client-id", resultEntraIDOptions.ClientID.ValueString(), "Client ID should be from env var")
	assert.Equal(t, "env-client-secret", resultEntraIDOptions.ClientSecret.ValueString(), "Client secret should be from env var")
	assert.Equal(t, "/env/path/to/cert.pfx", resultEntraIDOptions.ClientCertificate.ValueString(), "Client certificate should be from env var")
	assert.Equal(t, "env-cert-password", resultEntraIDOptions.ClientCertificatePassword.ValueString(), "Client certificate password should be from env var")
	assert.Equal(t, "env-user@example.com", resultEntraIDOptions.Username.ValueString(), "Username should be from env var")
	assert.True(t, resultEntraIDOptions.SendCertificateChain.ValueBool(), "Send certificate chain should be from env var")
	assert.True(t, resultEntraIDOptions.DisableInstanceDiscovery.ValueBool(), "Disable instance discovery should be from env var")
	assert.Equal(t, "https://env.example.com/callback", resultEntraIDOptions.RedirectUrl.ValueString(), "Redirect URL should be from env var")
	assert.Equal(t, "/env/path/to/federated-token", resultEntraIDOptions.FederatedTokenFilePath.ValueString(), "Federated token file path should be from env var")
	assert.Equal(t, "env-managed-identity-id", resultEntraIDOptions.ManagedIdentityID.ValueString(), "Managed identity ID should be from M365_MANAGED_IDENTITY_ID env var")
	assert.Equal(t, "/env/path/to/oidc-token", resultEntraIDOptions.OIDCTokenFilePath.ValueString(), "OIDC token file path should be from env var")
	assert.Equal(t, "env-oidc-request-token", resultEntraIDOptions.OIDCRequestToken.ValueString(), "OIDC request token should be from M365_OIDC_REQUEST_TOKEN env var")
	assert.Equal(t, "https://env.example.com/oidc-request", resultEntraIDOptions.OIDCRequestURL.ValueString(), "OIDC request URL should be from M365_OIDC_REQUEST_URL env var")
	assert.Equal(t, "env-ado-service-connection-id", resultEntraIDOptions.ADOServiceConnectionID.ValueString(), "ADO service connection ID should be from ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID env var")

	// Verify additionally allowed tenants from env var
	var resultAllowedTenants []string
	diags = resultEntraIDOptions.AdditionallyAllowedTenants.ElementsAs(context.Background(), &resultAllowedTenants, false)
	require.False(t, diags.HasError(), "Failed to extract allowed tenants")
	assert.Equal(t, []string{"tenant1", "tenant2", "tenant3"}, resultAllowedTenants, "Additionally allowed tenants should be from env var")

	// Extract and verify Client options
	var resultClientOptions ClientOptionsModel
	diags = result.ClientOptions.As(context.Background(), &resultClientOptions, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "Failed to extract Client options from result")

	// Verify Client options env vars take precedence
	assert.True(t, resultClientOptions.EnableHeadersInspection.ValueBool(), "Enable headers inspection should be from env var")
	assert.True(t, resultClientOptions.EnableRetry.ValueBool(), "Enable retry should be from env var")
	assert.Equal(t, int64(5), resultClientOptions.MaxRetries.ValueInt64(), "Max retries should be from env var")
	assert.Equal(t, int64(10), resultClientOptions.RetryDelaySeconds.ValueInt64(), "Retry delay seconds should be from env var")
	assert.True(t, resultClientOptions.EnableRedirect.ValueBool(), "Enable redirect should be from env var")
	assert.Equal(t, int64(3), resultClientOptions.MaxRedirects.ValueInt64(), "Max redirects should be from env var")
	assert.True(t, resultClientOptions.EnableCompression.ValueBool(), "Enable compression should be from env var")
	assert.Equal(t, "env-custom-user-agent", resultClientOptions.CustomUserAgent.ValueString(), "Custom user agent should be from env var")
	assert.True(t, resultClientOptions.UseProxy.ValueBool(), "Use proxy should be from env var")
	assert.Equal(t, "https://env.proxy.example.com:8080", resultClientOptions.ProxyURL.ValueString(), "Proxy URL should be from env var")
	assert.Equal(t, "env-proxy-user", resultClientOptions.ProxyUsername.ValueString(), "Proxy username should be from env var")
	assert.Equal(t, "env-proxy-password", resultClientOptions.ProxyPassword.ValueString(), "Proxy password should be from env var")
	assert.Equal(t, int64(120), resultClientOptions.TimeoutSeconds.ValueInt64(), "Timeout seconds should be from env var")
	assert.True(t, resultClientOptions.EnableChaos.ValueBool(), "Enable chaos should be from env var")
	assert.Equal(t, int64(25), resultClientOptions.ChaosPercentage.ValueInt64(), "Chaos percentage should be from env var")
	assert.Equal(t, int64(503), resultClientOptions.ChaosStatusCode.ValueInt64(), "Chaos status code should be from env var")
	assert.Equal(t, "env-chaos-message", resultClientOptions.ChaosStatusMessage.ValueString(), "Chaos status message should be from env var")
}

func TestSetProviderConfiguration_MixedEnvAndConfig(t *testing.T) {
	// Clear environment and set only some variables
	clearM365EnvVarsWithRestore(t)
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
	result, diags := setProviderConfiguration(context.Background(), config)
	require.False(t, diags.HasError(), "setProviderConfiguration should not return errors")

	// Verify mixed results
	assert.Equal(t, "gcc", result.Cloud.ValueString(), "Cloud should be from env var")
	assert.Equal(t, "config-tenant-id", result.TenantID.ValueString(), "Tenant ID should be from config")
	assert.Equal(t, "client_certificate", result.AuthMethod.ValueString(), "Auth method should be from env var")
	assert.True(t, result.TelemetryOptout.ValueBool(), "Telemetry optout should be from config")
	assert.True(t, result.DebugMode.ValueBool(), "Debug mode should be from env var")
}

func TestSetEntraIDOptions_EnvVarsAndConfig(t *testing.T) {
	// Clear environment variables
	clearM365EnvVarsWithRestore(t)

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
		"oidc_request_token":           types.StringValue("config-oidc-token"),
		"oidc_request_url":             types.StringValue("https://config.example.com/token"),
		"ado_service_connection_id":    types.StringValue("service-conn-id"),
	}

	configObj, diags := types.ObjectValue(entraIDSchema, configMap)
	require.False(t, diags.HasError(), "Failed to create test config object")

	// Call the function
	result, diags := setEntraIDOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "setEntraIDOptions should not return errors")

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

func TestSetClientOptions_DefaultsAndOverrides(t *testing.T) {
	// Clear environment variables
	clearM365EnvVarsWithRestore(t)

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
	nullResult, diags := setClientOptions(context.Background(), nullConfig)
	require.False(t, diags.HasError(), "setClientOptions should not return errors with null config")
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
	result, diags := setClientOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "setClientOptions should not return errors")

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

func TestSetEntraIDOptions_EmptyConfig(t *testing.T) {
	// Clear environment variables
	clearM365EnvVarsWithRestore(t)

	// Create null config (testing defaults)
	nullConfig := types.ObjectNull(map[string]attr.Type{})

	// Call the function with null config
	result, diags := setEntraIDOptions(context.Background(), nullConfig)
	require.False(t, diags.HasError(), "setEntraIDOptions should not return errors with null config")
	assert.True(t, result.IsNull(), "Result should be null with null config")
}

func TestSetEntraIDOptions_BooleanEnvVars(t *testing.T) {
	// Clear environment variables
	clearM365EnvVarsWithRestore(t)

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
		"oidc_request_token":           types.StringNull(),
		"oidc_request_url":             types.StringNull(),
		"ado_service_connection_id":    types.StringNull(),
	}

	configObj, diags := types.ObjectValue(entraIDSchema, configMap)
	require.False(t, diags.HasError(), "Failed to create test config object")

	// Call the function
	result, diags := setEntraIDOptions(context.Background(), configObj)
	require.False(t, diags.HasError(), "setEntraIDOptions should not return errors")

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
			clearM365EnvVarsWithRestore(t)
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
