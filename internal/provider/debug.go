package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func logDebugInfo(ctx context.Context, req provider.ConfigureRequest, data M365ProviderModel) {
	if !data.DebugMode.ValueBool() {
		return
	}

	tflog.Info(ctx, "==== M365ProviderModel Debug Information ====")

	logEnvironmentVariables(ctx)
	logSchemaValues(ctx, req)
	logProviderDataModel(ctx, data)

	tflog.Info(ctx, "========================================")
}

func logEnvironmentVariables(ctx context.Context) {
	tflog.Info(ctx, "==== Environment Variables ====")
	envVars := []string{
		"M365_CLOUD", "M365_TENANT_ID", "M365_AUTH_METHOD", "M365_CLIENT_ID",
		"M365_CLIENT_SECRET", "M365_CLIENT_CERTIFICATE_BASE64", "M365_CLIENT_CERTIFICATE_FILE_PATH",
		"M365_CLIENT_CERTIFICATE_PASSWORD", "M365_USERNAME", "M365_PASSWORD",
		"M365_REDIRECT_URL", "M365_USE_PROXY", "M365_PROXY_URL", "M365_ENABLE_CHAOS",
		"M365_TELEMETRY_OPTOUT", "M365_DEBUG_MODE",
	}

	for _, env := range envVars {
		value := os.Getenv(env)
		if isSecretValue(env) && value != "" {
			value = "[REDACTED]"
		}
		tflog.Info(ctx, fmt.Sprintf("%s: %s", env, value))
	}
}

func logSchemaValues(ctx context.Context, req provider.ConfigureRequest) {
	tflog.Info(ctx, "==== Values Set in Schema ====")
	var config M365ProviderModel
	diags := req.Config.Get(ctx, &config)
	if diags.HasError() {
		tflog.Error(ctx, "Error retrieving schema values", map[string]interface{}{"diagnostics": diags.Errors()})
		return
	}

	logSchemaValue(ctx, "Tenant ID", config.TenantID)
	logSchemaValue(ctx, "Auth Method", config.AuthMethod)
	logSchemaValue(ctx, "Client ID", config.ClientID)
	logSchemaValue(ctx, "Client Secret", config.ClientSecret)
	logSchemaValue(ctx, "Client Certificate Base64", config.ClientCertificateBase64)
	logSchemaValue(ctx, "Client Certificate File Path", config.ClientCertificateFilePath)
	logSchemaValue(ctx, "Client Certificate Password", config.ClientCertificatePassword)
	logSchemaValue(ctx, "Username", config.Username)
	logSchemaValue(ctx, "Password", config.Password)
	logSchemaValue(ctx, "Redirect URL", config.RedirectURL)
	logSchemaValue(ctx, "Use Proxy", config.UseProxy)
	logSchemaValue(ctx, "Proxy URL", config.ProxyURL)
	logSchemaValue(ctx, "Cloud", config.Cloud)
	logSchemaValue(ctx, "Enable Chaos", config.EnableChaos)
	logSchemaValue(ctx, "Telemetry Optout", config.TelemetryOptout)
	logSchemaValue(ctx, "Debug Mode", config.DebugMode)
}

func logProviderDataModel(ctx context.Context, data M365ProviderModel) {
	tflog.Info(ctx, "==== Values Mapped to Provider Data Model ====")
	tflog.Info(ctx, fmt.Sprintf("Tenant ID Length: %d", len(data.TenantID.ValueString())))
	tflog.Info(ctx, fmt.Sprintf("Auth Method: %s", data.AuthMethod.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Client ID Length: %d", len(data.ClientID.ValueString())))
	tflog.Info(ctx, fmt.Sprintf("Client Secret Length: %d", len(data.ClientSecret.ValueString())))
	tflog.Info(ctx, fmt.Sprintf("Client Certificate Base64 Length: %d", len(data.ClientCertificateBase64.ValueString())))
	tflog.Info(ctx, fmt.Sprintf("Client Certificate File Path: %s", data.ClientCertificateFilePath.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Client Certificate Password Set: %t", data.ClientCertificatePassword.ValueString() != ""))
	tflog.Info(ctx, fmt.Sprintf("Username Set: %t", data.Username.ValueString() != ""))
	tflog.Info(ctx, fmt.Sprintf("Password Set: %t", data.Password.ValueString() != ""))
	tflog.Info(ctx, fmt.Sprintf("Redirect URL: %s", data.RedirectURL.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Use Proxy: %t", data.UseProxy.ValueBool()))
	tflog.Info(ctx, fmt.Sprintf("Proxy URL: %s", data.ProxyURL.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Cloud: %s", data.Cloud.ValueString()))
	tflog.Info(ctx, fmt.Sprintf("Enable Chaos: %t", data.EnableChaos.ValueBool()))
	tflog.Info(ctx, fmt.Sprintf("Telemetry Optout: %t", data.TelemetryOptout.ValueBool()))
	tflog.Info(ctx, fmt.Sprintf("Debug Mode: %t", data.DebugMode.ValueBool()))
}

func logSchemaValue(ctx context.Context, name string, value interface{}) {
	switch v := value.(type) {
	case types.String:
		tflog.Info(ctx, fmt.Sprintf("%s: %t", name, !v.IsNull() && !v.IsUnknown()))
	case types.Bool:
		tflog.Info(ctx, fmt.Sprintf("%s: %t", name, !v.IsNull() && !v.IsUnknown()))
	default:
		tflog.Info(ctx, fmt.Sprintf("%s: Unknown type", name))
	}
}

func isSecretValue(envVar string) bool {
	secretVars := []string{"M365_CLIENT_SECRET", "M365_CLIENT_CERTIFICATE_BASE64", "M365_CLIENT_CERTIFICATE_PASSWORD", "M365_PASSWORD"}
	for _, secretVar := range secretVars {
		if envVar == secretVar {
			return true
		}
	}
	return false
}
