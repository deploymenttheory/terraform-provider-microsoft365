package provider

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

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

// populateProviderData populates the M365ProviderModel with values from the configuration
// or environment variables, using helper functions for default values.
func populateProviderData(ctx context.Context, config M365ProviderModel) (M365ProviderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	entraIDOptions, entraDiags := populateEntraIDOptions(ctx, config.EntraIDOptions)
	diags.Append(entraDiags...)

	clientOptions, clientDiags := populateClientOptions(ctx, config.ClientOptions)
	diags.Append(clientDiags...)

	return M365ProviderModel{
		Cloud:           types.StringValue(helpers.MultiEnvDefaultFunc([]string{"M365_CLOUD", "AZURE_CLOUD"}, config.Cloud.ValueString())),
		TenantID:        types.StringValue(helpers.EnvDefaultFunc("M365_TENANT_ID", config.TenantID.ValueString())),
		AuthMethod:      types.StringValue(helpers.EnvDefaultFunc("M365_AUTH_METHOD", config.AuthMethod.ValueString())),
		EntraIDOptions:  entraIDOptions,
		ClientOptions:   clientOptions,
		TelemetryOptout: types.BoolValue(helpers.EnvDefaultFuncBool("M365_TELEMETRY_OPTOUT", config.TelemetryOptout.ValueBool())),
		DebugMode:       types.BoolValue(helpers.EnvDefaultFuncBool("M365_DEBUG_MODE", config.DebugMode.ValueBool())),
	}, diags
}

func populateEntraIDOptions(ctx context.Context, config types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var entraIDOptions EntraIDOptionsModel

	entraIDSchema := schemaToAttrTypes(EntraIDOptionsSchema())

	if config.IsNull() || config.IsUnknown() {
		return types.ObjectNull(entraIDSchema), diags
	}

	diags.Append(config.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(entraIDSchema), diags
	}

	var defaultAllowedTenants []string
	if !entraIDOptions.AdditionallyAllowedTenants.IsNull() && !entraIDOptions.AdditionallyAllowedTenants.IsUnknown() {
		entraIDOptions.AdditionallyAllowedTenants.ElementsAs(ctx, &defaultAllowedTenants, false)
	}

	allowedTenants := helpers.EnvDefaultFuncStringList("M365_ADDITIONALLY_ALLOWED_TENANTS", defaultAllowedTenants)
	allowedTenantsList, diags := types.ListValueFrom(ctx, types.StringType, allowedTenants)
	if diags.HasError() {
		return types.ObjectNull(entraIDSchema), diags
	}

	return types.ObjectValueMust(entraIDSchema, map[string]attr.Value{
		"client_id":                    types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_ID", entraIDOptions.ClientID.ValueString())),
		"client_secret":                types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_SECRET", entraIDOptions.ClientSecret.ValueString())),
		"client_certificate":           types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_CERTIFICATE_FILE_PATH", entraIDOptions.ClientCertificate.ValueString())),
		"client_certificate_password":  types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_CERTIFICATE_PASSWORD", entraIDOptions.ClientCertificatePassword.ValueString())),
		"send_certificate_chain":       types.BoolValue(helpers.EnvDefaultFuncBool("M365_SEND_CERTIFICATE_CHAIN", entraIDOptions.SendCertificateChain.ValueBool())),
		"username":                     types.StringValue(helpers.EnvDefaultFunc("M365_USERNAME", entraIDOptions.Username.ValueString())),
		"password":                     types.StringValue(helpers.EnvDefaultFunc("M365_PASSWORD", entraIDOptions.Password.ValueString())),
		"disable_instance_discovery":   types.BoolValue(helpers.EnvDefaultFuncBool("M365_DISABLE_INSTANCE_DISCOVERY", entraIDOptions.DisableInstanceDiscovery.ValueBool())),
		"additionally_allowed_tenants": allowedTenantsList,
		"redirect_url":                 types.StringValue(helpers.EnvDefaultFunc("M365_REDIRECT_URI", entraIDOptions.RedirectUrl.ValueString())),
	}), diags
}

func populateClientOptions(ctx context.Context, config types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var clientOptions ClientOptionsModel

	clientSchema := schemaToAttrTypes(ClientOptionsSchema())

	if config.IsNull() || config.IsUnknown() {
		return types.ObjectNull(clientSchema), diags
	}

	diags.Append(config.As(ctx, &clientOptions, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return types.ObjectNull(clientSchema), diags
	}

	return types.ObjectValueMust(clientSchema, map[string]attr.Value{
		"enable_headers_inspection": types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_HEADERS_INSPECTION", clientOptions.EnableHeadersInspection.ValueBool())),
		"enable_retry":              types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_RETRY", clientOptions.EnableRetry.ValueBool())),
		"max_retries":               helpers.EnvDefaultFuncInt64Value("M365_MAX_RETRIES", clientOptions.MaxRetries),
		"retry_delay_seconds":       helpers.EnvDefaultFuncInt64Value("M365_RETRY_DELAY_SECONDS", clientOptions.RetryDelaySeconds),
		"enable_redirect":           types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_REDIRECT", clientOptions.EnableRedirect.ValueBool())),
		"max_redirects":             helpers.EnvDefaultFuncInt64Value("M365_MAX_REDIRECTS", clientOptions.MaxRedirects),
		"enable_compression":        types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_COMPRESSION", clientOptions.EnableCompression.ValueBool())),
		"custom_user_agent":         types.StringValue(helpers.EnvDefaultFunc("M365_CUSTOM_USER_AGENT", clientOptions.CustomUserAgent.ValueString())),
		"use_proxy":                 types.BoolValue(helpers.EnvDefaultFuncBool("M365_USE_PROXY", clientOptions.UseProxy.ValueBool())),
		"proxy_url":                 types.StringValue(helpers.EnvDefaultFunc("M365_PROXY_URL", clientOptions.ProxyURL.ValueString())),
		"proxy_username":            types.StringValue(helpers.EnvDefaultFunc("M365_PROXY_USERNAME", clientOptions.ProxyUsername.ValueString())),
		"proxy_password":            types.StringValue(helpers.EnvDefaultFunc("M365_PROXY_PASSWORD", clientOptions.ProxyPassword.ValueString())),
		"timeout_seconds":           helpers.EnvDefaultFuncInt64Value("M365_TIMEOUT_SECONDS", clientOptions.TimeoutSeconds),
		"enable_chaos":              types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_CHAOS", clientOptions.EnableChaos.ValueBool())),
		"chaos_percentage":          helpers.EnvDefaultFuncInt64Value("M365_CHAOS_PERCENTAGE", clientOptions.ChaosPercentage),
		"chaos_status_code":         helpers.EnvDefaultFuncInt64Value("M365_CHAOS_STATUS_CODE", clientOptions.ChaosStatusCode),
		"chaos_status_message":      types.StringValue(helpers.EnvDefaultFunc("M365_CHAOS_STATUS_MESSAGE", clientOptions.ChaosStatusMessage.ValueString())),
	}), diags
}