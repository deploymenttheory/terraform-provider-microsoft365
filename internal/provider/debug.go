package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func logDebugInfo(ctx context.Context, req provider.ConfigureRequest, data M365ProviderModel) {
	if !data.DebugMode.ValueBool() {
		return
	}

	fmt.Println("\n==== M365ProviderModel Debug Information ====")

	var config M365ProviderModel
	req.Config.Get(ctx, &config)

	logValueSource("Cloud", []string{"M365_CLOUD", "AZURE_CLOUD"}, config.Cloud, data.Cloud)
	logValueSource("Tenant ID", []string{"M365_TENANT_ID"}, config.TenantID, data.TenantID)
	logValueSource("Auth Method", []string{"M365_AUTH_METHOD"}, config.AuthMethod, data.AuthMethod)
	logValueSource("Client ID", []string{"M365_CLIENT_ID"}, config.ClientID, data.ClientID)
	logValueSource("Client Secret", []string{"M365_CLIENT_SECRET"}, config.ClientSecret, data.ClientSecret)
	logValueSource("Client Certificate Base64", []string{"M365_CLIENT_CERTIFICATE_BASE64"}, config.ClientCertificateBase64, data.ClientCertificateBase64)
	logValueSource("Client Certificate File Path", []string{"M365_CLIENT_CERTIFICATE_FILE_PATH"}, config.ClientCertificateFilePath, data.ClientCertificateFilePath)
	logValueSource("Client Certificate Password", []string{"M365_CLIENT_CERTIFICATE_PASSWORD"}, config.ClientCertificatePassword, data.ClientCertificatePassword)
	logValueSource("Username", []string{"M365_USERNAME"}, config.Username, data.Username)
	logValueSource("Password", []string{"M365_PASSWORD"}, config.Password, data.Password)
	logValueSource("Redirect URL", []string{"M365_REDIRECT_URL"}, config.RedirectURL, data.RedirectURL)
	logBoolValueSource("Use Proxy", []string{"M365_USE_PROXY"}, config.UseProxy, data.UseProxy)
	logValueSource("Proxy URL", []string{"M365_PROXY_URL"}, config.ProxyURL, data.ProxyURL)
	logBoolValueSource("Enable Chaos", []string{"M365_ENABLE_CHAOS"}, config.EnableChaos, data.EnableChaos)
	logBoolValueSource("Telemetry Optout", []string{"M365_TELEMETRY_OPTOUT"}, config.TelemetryOptout, data.TelemetryOptout)
	logBoolValueSource("Debug Mode", []string{"M365_DEBUG_MODE"}, config.DebugMode, data.DebugMode)

	fmt.Println("========================================")
}

func logValueSource(name string, envVars []string, configValue, dataValue types.String) {
	var source, value string
	for _, env := range envVars {
		if v := os.Getenv(env); v != "" {
			source = "Environment Variable"
			value = v
			break
		}
	}
	if source == "" {
		if !configValue.IsNull() && !configValue.IsUnknown() {
			source = "HCL Configuration"
			value = configValue.ValueString()
		} else {
			source = "HCL Default"
			value = dataValue.ValueString()
		}
	}
	fmt.Printf("%s: %s (Source: %s)\n", name, maskSensitiveValue(value), source)
}

func logBoolValueSource(name string, envVars []string, configValue, dataValue types.Bool) {
	var source string
	var value bool
	for _, env := range envVars {
		if v := os.Getenv(env); v != "" {
			source = "Environment Variable"
			value = v == "true" || v == "1"
			break
		}
	}
	if source == "" {
		if !configValue.IsNull() && !configValue.IsUnknown() {
			source = "HCL Configuration"
			value = configValue.ValueBool()
		} else {
			source = "HCL Default"
			value = dataValue.ValueBool()
		}
	}
	fmt.Printf("%s: %t (Source: %s)\n", name, value, source)
}

func maskSensitiveValue(value string) string {
	if value != "" {
		return value
	}
	return ""
}
