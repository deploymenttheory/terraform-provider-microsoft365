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

// TestConvertToClientProviderData_EmptyProviderBlockWithEnvVars tests the specific bug
// where environment variables are set but the provider block is empty (like in GitHub Actions OIDC)
func TestConvertToClientProviderData_EmptyProviderBlockWithEnvVars(t *testing.T) {
	// Clear environment and set test variables (simulating GitHub Actions OIDC scenario)
	clearM365EnvVarsWithRestore(t)
	os.Setenv("M365_CLIENT_ID", "github-actions-client-id")
	os.Setenv("M365_TENANT_ID", "github-actions-tenant-id")
	os.Setenv("M365_AUTH_METHOD", "oidc_github")
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "https://github.actions/token")
	os.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "github-actions-token")

	// Create an empty provider configuration (simulating empty provider block)
	config := M365ProviderModel{
		Cloud:           types.StringNull(),
		TenantID:        types.StringNull(),
		AuthMethod:      types.StringNull(),
		EntraIDOptions:  types.ObjectNull(map[string]attr.Type{}),
		ClientOptions:   types.ObjectNull(map[string]attr.Type{}),
		TelemetryOptout: types.BoolNull(),
		DebugMode:       types.BoolNull(),
	}

	// Step 1: Process the configuration with environment variables
	processedConfig, diags := setProviderConfiguration(context.Background(), config)
	require.False(t, diags.HasError(), "setProviderConfiguration should not return errors")

	// Verify that setProviderConfiguration correctly processes environment variables
	assert.Equal(t, "github-actions-tenant-id", processedConfig.TenantID.ValueString(), "TenantID should be processed from env var")
	assert.Equal(t, "oidc_github", processedConfig.AuthMethod.ValueString(), "AuthMethod should be processed from env var")

	// Debug: Check if the processed EntraIDOptions contains the expected values
	t.Logf("processedConfig.EntraIDOptions.IsNull(): %t", processedConfig.EntraIDOptions.IsNull())
	t.Logf("processedConfig.EntraIDOptions.IsUnknown(): %t", processedConfig.EntraIDOptions.IsUnknown())

	// Extract the processed EntraIDOptions to verify it contains the expected values
	var processedEntraIDOptions EntraIDOptionsModel
	extractDiags := processedConfig.EntraIDOptions.As(context.Background(), &processedEntraIDOptions, basetypes.ObjectAsOptions{})
	require.False(t, extractDiags.HasError(), "Failed to extract processed EntraIDOptions")
	t.Logf("processedEntraIDOptions.ClientID.ValueString(): '%s'", processedEntraIDOptions.ClientID.ValueString())

	// Step 2: Convert to client provider data (this is where the bug occurs)
	clientData := convertToClientProviderData(context.Background(), &processedConfig)

	// This is the bug: clientData.EntraIDOptions.ClientID should be "github-actions-client-id" but it's empty
	t.Logf("TenantID from clientData: '%s'", clientData.TenantID)
	t.Logf("ClientID from clientData.EntraIDOptions: '%s'", clientData.EntraIDOptions.ClientID)

	// These assertions demonstrate the bug
	assert.Equal(t, "github-actions-tenant-id", clientData.TenantID, "TenantID should be correctly converted")

	// This assertion will FAIL, demonstrating the bug
	assert.Equal(t, "github-actions-client-id", clientData.EntraIDOptions.ClientID,
		"ClientID should be correctly converted from processed EntraIDOptions, but it's empty due to the bug")
}

// TestConvertToClientProviderData_ExplicitProviderBlockWithEnvVars tests that the bug doesn't occur
// when the provider block explicitly sets entra_id_options
func TestConvertToClientProviderData_ExplicitProviderBlockWithEnvVars(t *testing.T) {
	// Clear environment and set test variables
	clearM365EnvVarsWithRestore(t)
	os.Setenv("M365_CLIENT_ID", "env-client-id")
	os.Setenv("M365_TENANT_ID", "env-tenant-id")
	os.Setenv("M365_AUTH_METHOD", "oidc_github")

	// Create EntraID options with explicit client_id (this will be overridden by env var)
	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())
	allowedTenantsList, diags := types.ListValueFrom(context.Background(), types.StringType, []string{})
	require.False(t, diags.HasError(), "Failed to create allowed tenants list")

	entraIDConfigMap := map[string]attr.Value{
		"client_id":                    types.StringValue("explicit-client-id"),
		"client_secret":                types.StringNull(),
		"client_certificate":           types.StringNull(),
		"client_certificate_password":  types.StringNull(),
		"send_certificate_chain":       types.BoolNull(),
		"username":                     types.StringNull(),
		"disable_instance_discovery":   types.BoolNull(),
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringNull(),
		"federated_token_file_path":    types.StringNull(),
		"managed_identity_id":          types.StringNull(),
		"oidc_token_file_path":         types.StringNull(),
		"oidc_request_token":           types.StringNull(),
		"oidc_request_url":             types.StringNull(),
		"ado_service_connection_id":    types.StringNull(),
	}

	entraIDConfigObj, diags := types.ObjectValue(entraIDSchema, entraIDConfigMap)
	require.False(t, diags.HasError(), "Failed to create EntraID config object")

	// Create a provider configuration with explicit entra_id_options
	config := M365ProviderModel{
		Cloud:           types.StringNull(),
		TenantID:        types.StringNull(),
		AuthMethod:      types.StringNull(),
		EntraIDOptions:  entraIDConfigObj,
		ClientOptions:   types.ObjectNull(map[string]attr.Type{}),
		TelemetryOptout: types.BoolNull(),
		DebugMode:       types.BoolNull(),
	}

	// Step 1: Process the configuration with environment variables
	processedConfig, diags := setProviderConfiguration(context.Background(), config)
	require.False(t, diags.HasError(), "setProviderConfiguration should not return errors")

	// Step 2: Convert to client provider data
	clientData := convertToClientProviderData(context.Background(), &processedConfig)

	t.Logf("TenantID from clientData: '%s'", clientData.TenantID)
	t.Logf("ClientID from clientData.EntraIDOptions: '%s'", clientData.EntraIDOptions.ClientID)

	// These assertions should pass because the EntraIDOptions were explicitly provided
	assert.Equal(t, "env-tenant-id", clientData.TenantID, "TenantID should be correctly converted")
	assert.Equal(t, "env-client-id", clientData.EntraIDOptions.ClientID,
		"ClientID should be correctly converted when entra_id_options is explicitly provided")
}

// clearM365EnvVarsWithRestore clears M365 related environment variables and restores them after the test.
func clearM365EnvVarsWithRestore(t *testing.T) {
	t.Helper()
	originalEnv := make(map[string]string)

	// Define all environment variables that might be used in tests
	envVarsToClean := []string{
		// M365 prefixed variables
		"M365_CLOUD", "M365_TENANT_ID", "M365_CLIENT_ID", "M365_CLIENT_SECRET",
		"M365_AUTH_METHOD", "M365_REDIRECT_URL", "M365_USE_PROXY", "M365_PROXY_URL",
		"M365_ENABLE_CHAOS", "M365_TELEMETRY_OPTOUT", "M365_DEBUG_MODE",
		"M365_CLIENT_CERTIFICATE_FILE_PATH", "M365_CLIENT_CERTIFICATE_PASSWORD",
		"M365_USERNAME", "M365_SEND_CERTIFICATE_CHAIN", "M365_DISABLE_INSTANCE_DISCOVERY",
		"M365_ADDITIONALLY_ALLOWED_TENANTS", "M365_MANAGED_IDENTITY_ID",
		"M365_OIDC_TOKEN_FILE_PATH", "M365_OIDC_REQUEST_URL", "M365_OIDC_REQUEST_TOKEN",
		// Azure and GitHub Actions variables
		"AZURE_CLOUD", "AZURE_CLIENT_ID", "AZURE_FEDERATED_TOKEN_FILE",
		"ACTIONS_ID_TOKEN_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		"ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "ARM_OIDC_AZURE_SERVICE_CONNECTION_ID",
	}

	// Save current values and clear them
	for _, envVar := range envVarsToClean {
		if value := os.Getenv(envVar); value != "" {
			originalEnv[envVar] = value
		}
		os.Unsetenv(envVar)
	}

	t.Cleanup(func() {
		// Clear any variables that might have been set during the test
		for _, envVar := range envVarsToClean {
			os.Unsetenv(envVar)
		}
		// Restore original values
		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	})
}
