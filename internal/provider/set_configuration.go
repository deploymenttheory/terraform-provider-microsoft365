package provider

import (
	"context"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// setProviderConfiguration populates the M365ProviderModel with values from the configuration
// or environment variables, using helper functions for default values.
func setProviderConfiguration(ctx context.Context, config M365ProviderModel) (M365ProviderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	entraIDOptions, entraDiags := setEntraIDOptions(ctx, config.EntraIDOptions)
	diags.Append(entraDiags...)

	clientOptions, clientDiags := setClientOptions(ctx, config.ClientOptions)
	diags.Append(clientDiags...)

	return M365ProviderModel{
		Cloud:           types.StringValue(helpers.GetFirstEnvString([]string{"M365_CLOUD", "AZURE_CLOUD"}, config.Cloud.ValueString())),
		TenantID:        types.StringValue(helpers.GetEnvString("M365_TENANT_ID", config.TenantID.ValueString())),
		AuthMethod:      types.StringValue(helpers.GetEnvString("M365_AUTH_METHOD", config.AuthMethod.ValueString())),
		EntraIDOptions:  entraIDOptions,
		ClientOptions:   clientOptions,
		TelemetryOptout: types.BoolValue(helpers.GetEnvBool("M365_TELEMETRY_OPTOUT", config.TelemetryOptout.ValueBool())),
		DebugMode:       types.BoolValue(helpers.GetEnvBool("M365_DEBUG_MODE", config.DebugMode.ValueBool())),
	}, diags
}

// unmarshalConfigOrCheckEnvVars attempts to unmarshal the config object into the target model.
// If config is null/unknown and checkEnvVars is true, it checks if any of the provided
// environment variables are set. Returns the unmarshaled model, whether to proceed, and any diagnostics.
func unmarshalConfigOrCheckEnvVars(
	ctx context.Context,
	config types.Object,
	target any,
	checkEnvVars bool,
	envVarsToCheck []string,
) (shouldProceed bool, diags diag.Diagnostics) {
	if !config.IsNull() && !config.IsUnknown() {
		diags.Append(config.As(ctx, target, basetypes.ObjectAsOptions{})...)
		return !diags.HasError(), diags
	}

	if !checkEnvVars {
		return false, diags
	}

	// Check if any relevant environment variables are set
	for _, envVar := range envVarsToCheck {
		if os.Getenv(envVar) != "" {
			return true, diags
		}
	}

	return false, diags
}

// setEntraIDOptions sets the EntraIDOptionsModel from the configuration or environment variables
func setEntraIDOptions(ctx context.Context, config types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var entraIDOptions EntraIDOptionsModel
	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())

	envVarsToCheck := []string{
		"M365_CLIENT_ID", "M365_CLIENT_SECRET", "M365_CLIENT_CERTIFICATE_FILE_PATH",
		"M365_CLIENT_CERTIFICATE_PASSWORD", "M365_USERNAME", "M365_SEND_CERTIFICATE_CHAIN",
		"M365_DISABLE_INSTANCE_DISCOVERY", "M365_ADDITIONALLY_ALLOWED_TENANTS",
		"M365_REDIRECT_URI", "AZURE_FEDERATED_TOKEN_FILE", "M365_MANAGED_IDENTITY_ID",
		"AZURE_CLIENT_ID", "M365_OIDC_TOKEN_FILE_PATH", "M365_OIDC_REQUEST_URL",
		"ACTIONS_ID_TOKEN_REQUEST_URL", "M365_OIDC_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN", "ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID",
		"ARM_OIDC_AZURE_SERVICE_CONNECTION_ID",
	}

	shouldProceed, unmarshalDiags := unmarshalConfigOrCheckEnvVars(
		ctx, config, &entraIDOptions, true, envVarsToCheck,
	)
	diags.Append(unmarshalDiags...)

	if !shouldProceed {
		return types.ObjectNull(entraIDSchema), diags
	}

	// Handle additionally_allowed_tenants list
	var defaultAllowedTenants []string
	if !entraIDOptions.AdditionallyAllowedTenants.IsNull() && !entraIDOptions.AdditionallyAllowedTenants.IsUnknown() {
		entraIDOptions.AdditionallyAllowedTenants.ElementsAs(ctx, &defaultAllowedTenants, false)
	}
	allowedTenants := helpers.GetEnvStringSlice("M365_ADDITIONALLY_ALLOWED_TENANTS", defaultAllowedTenants)
	allowedTenantsList, listDiags := types.ListValueFrom(ctx, types.StringType, allowedTenants)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ObjectNull(entraIDSchema), diags
	}

	// Handle OIDC fields with multiple env var sources
	oidcRequestURL := helpers.GetFirstEnvString(
		[]string{"M365_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"},
		entraIDOptions.OIDCRequestURL.ValueString(),
	)
	oidcRequestToken := helpers.GetFirstEnvString(
		[]string{"M365_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"},
		entraIDOptions.OIDCRequestToken.ValueString(),
	)

	tflog.Info(ctx, "OIDC configuration debug", map[string]any{
		"oidc_request_url":       oidcRequestURL,
		"oidc_request_token_set": oidcRequestToken != "",
		"env_actions_url":        os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
		"env_actions_token_set":  os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN") != "",
	})

	return types.ObjectValueMust(entraIDSchema, map[string]attr.Value{
		"client_id":                    types.StringValue(helpers.GetEnvString("M365_CLIENT_ID", entraIDOptions.ClientID.ValueString())),
		"client_secret":                types.StringValue(helpers.GetEnvString("M365_CLIENT_SECRET", entraIDOptions.ClientSecret.ValueString())),
		"client_certificate":           types.StringValue(helpers.GetEnvString("M365_CLIENT_CERTIFICATE_FILE_PATH", entraIDOptions.ClientCertificate.ValueString())),
		"client_certificate_password":  types.StringValue(helpers.GetEnvString("M365_CLIENT_CERTIFICATE_PASSWORD", entraIDOptions.ClientCertificatePassword.ValueString())),
		"username":                     types.StringValue(helpers.GetEnvString("M365_USERNAME", entraIDOptions.Username.ValueString())),
		"send_certificate_chain":       types.BoolValue(helpers.GetEnvBool("M365_SEND_CERTIFICATE_CHAIN", entraIDOptions.SendCertificateChain.ValueBool())),
		"disable_instance_discovery":   types.BoolValue(helpers.GetEnvBool("M365_DISABLE_INSTANCE_DISCOVERY", entraIDOptions.DisableInstanceDiscovery.ValueBool())),
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringValue(helpers.GetEnvString("M365_REDIRECT_URI", entraIDOptions.RedirectUrl.ValueString())),
		"federated_token_file_path":    types.StringValue(helpers.GetEnvString("AZURE_FEDERATED_TOKEN_FILE", entraIDOptions.FederatedTokenFilePath.ValueString())),
		"managed_identity_id":          types.StringValue(helpers.GetFirstEnvString([]string{"M365_MANAGED_IDENTITY_ID", "AZURE_CLIENT_ID"}, entraIDOptions.ManagedIdentityID.ValueString())),
		"oidc_token_file_path":         types.StringValue(helpers.GetEnvString("M365_OIDC_TOKEN_FILE_PATH", entraIDOptions.OIDCTokenFilePath.ValueString())),
		"oidc_request_token":           types.StringValue(oidcRequestToken),
		"oidc_request_url":             types.StringValue(oidcRequestURL),
		"ado_service_connection_id":    types.StringValue(helpers.GetFirstEnvString([]string{"ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "ARM_OIDC_AZURE_SERVICE_CONNECTION_ID"}, entraIDOptions.ADOServiceConnectionID.ValueString())),
	}), diags
}

// setClientOptions sets the ClientOptionsModel from the configuration or environment variables
func setClientOptions(ctx context.Context, config types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var clientOptions ClientOptionsModel
	clientSchema := schemaToAttrTypes(ClientOptionsSchema())

	shouldProceed, unmarshalDiags := unmarshalConfigOrCheckEnvVars(
		ctx, config, &clientOptions, false, nil,
	)
	diags.Append(unmarshalDiags...)

	if !shouldProceed {
		return types.ObjectNull(clientSchema), diags
	}

	return types.ObjectValueMust(clientSchema, map[string]attr.Value{
		"enable_headers_inspection": types.BoolValue(helpers.GetEnvBool("M365_ENABLE_HEADERS_INSPECTION", clientOptions.EnableHeadersInspection.ValueBool())),
		"enable_retry":              types.BoolValue(helpers.GetEnvBool("M365_ENABLE_RETRY", clientOptions.EnableRetry.ValueBool())),
		"max_retries":               helpers.GetEnvInt64("M365_MAX_RETRIES", clientOptions.MaxRetries),
		"retry_delay_seconds":       helpers.GetEnvInt64("M365_RETRY_DELAY_SECONDS", clientOptions.RetryDelaySeconds),
		"enable_redirect":           types.BoolValue(helpers.GetEnvBool("M365_ENABLE_REDIRECT", clientOptions.EnableRedirect.ValueBool())),
		"max_redirects":             helpers.GetEnvInt64("M365_MAX_REDIRECTS", clientOptions.MaxRedirects),
		"enable_compression":        types.BoolValue(helpers.GetEnvBool("M365_ENABLE_COMPRESSION", clientOptions.EnableCompression.ValueBool())),
		"custom_user_agent":         types.StringValue(helpers.GetEnvString("M365_CUSTOM_USER_AGENT", clientOptions.CustomUserAgent.ValueString())),
		"use_proxy":                 types.BoolValue(helpers.GetEnvBool("M365_USE_PROXY", clientOptions.UseProxy.ValueBool())),
		"proxy_url":                 types.StringValue(helpers.GetEnvString("M365_PROXY_URL", clientOptions.ProxyURL.ValueString())),
		"proxy_username":            types.StringValue(helpers.GetEnvString("M365_PROXY_USERNAME", clientOptions.ProxyUsername.ValueString())),
		"proxy_password":            types.StringValue(helpers.GetEnvString("M365_PROXY_PASSWORD", clientOptions.ProxyPassword.ValueString())),
		"timeout_seconds":           helpers.GetEnvInt64("M365_TIMEOUT_SECONDS", clientOptions.TimeoutSeconds),
		"enable_chaos":              types.BoolValue(helpers.GetEnvBool("M365_ENABLE_CHAOS", clientOptions.EnableChaos.ValueBool())),
		"chaos_percentage":          helpers.GetEnvInt64("M365_CHAOS_PERCENTAGE", clientOptions.ChaosPercentage),
		"chaos_status_code":         helpers.GetEnvInt64("M365_CHAOS_STATUS_CODE", clientOptions.ChaosStatusCode),
		"chaos_status_message":      types.StringValue(helpers.GetEnvString("M365_CHAOS_STATUS_MESSAGE", clientOptions.ChaosStatusMessage.ValueString())),
	}), diags
}

// schemaToAttrTypes converts a map of schema.Attribute to a map of attr.Type
func schemaToAttrTypes(schemaMap map[string]schema.Attribute) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type)
	for k, v := range schemaMap {
		if v.GetType() != nil {
			attrTypes[k] = v.GetType()
		}
	}
	return attrTypes
}
