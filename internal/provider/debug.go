package provider

import (
	"context"
	"fmt"
	"os"
)

func logDebugInfo(ctx context.Context, data M365ProviderModel) {
	if data.DebugMode.ValueBool() {
		fmt.Println("\n==== M365ProviderModel Debug Information ====")

		fmt.Println("\n==== Environment Variables ====")
		envVars := []string{
			"M365_CLOUD", "M365_TENANT_ID", "M365_AUTH_METHOD", "M365_CLIENT_ID",
			"M365_CLIENT_SECRET", "M365_CLIENT_CERTIFICATE_BASE64", "M365_CLIENT_CERTIFICATE_FILE_PATH",
			"M365_CLIENT_CERTIFICATE_PASSWORD", "M365_USERNAME", "M365_PASSWORD",
			"M365_REDIRECT_URL", "M365_USE_PROXY", "M365_PROXY_URL", "M365_ENABLE_CHAOS",
			"M365_TELEMETRY_OPTOUT", "M365_DEBUG_MODE",
		}

		for _, env := range envVars {
			value := os.Getenv(env)
			if env == "M365_CLIENT_SECRET" || env == "M365_CLIENT_CERTIFICATE_BASE64" ||
				env == "M365_CLIENT_CERTIFICATE_PASSWORD" || env == "M365_PASSWORD" {
				if value != "" {
					value = "[REDACTED]"
				}
			}
			fmt.Printf("%s: %s\n", env, value)
		}

		fmt.Println("\n==== Values Set in Schema ====")
		fmt.Printf("Tenant ID: %t\n", !data.TenantID.IsNull() && !data.TenantID.IsUnknown())
		fmt.Printf("Auth Method: %t\n", !data.AuthMethod.IsNull() && !data.AuthMethod.IsUnknown())
		fmt.Printf("Client ID: %t\n", !data.ClientID.IsNull() && !data.ClientID.IsUnknown())
		fmt.Printf("Client Secret: %t\n", !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown())
		fmt.Printf("Client Certificate Base64: %t\n", !data.ClientCertificateBase64.IsNull() && !data.ClientCertificateBase64.IsUnknown())
		fmt.Printf("Client Certificate File Path: %t\n", !data.ClientCertificateFilePath.IsNull() && !data.ClientCertificateFilePath.IsUnknown())
		fmt.Printf("Client Certificate Password: %t\n", !data.ClientCertificatePassword.IsNull() && !data.ClientCertificatePassword.IsUnknown())
		fmt.Printf("Username: %t\n", !data.Username.IsNull() && !data.Username.IsUnknown())
		fmt.Printf("Password: %t\n", !data.Password.IsNull() && !data.Password.IsUnknown())
		fmt.Printf("Redirect URL: %t\n", !data.RedirectURL.IsNull() && !data.RedirectURL.IsUnknown())
		fmt.Printf("Use Proxy: %t\n", !data.UseProxy.IsNull() && !data.UseProxy.IsUnknown())
		fmt.Printf("Proxy URL: %t\n", !data.ProxyURL.IsNull() && !data.ProxyURL.IsUnknown())
		fmt.Printf("Cloud: %t\n", !data.Cloud.IsNull() && !data.Cloud.IsUnknown())
		fmt.Printf("Enable Chaos: %t\n", !data.EnableChaos.IsNull() && !data.EnableChaos.IsUnknown())
		fmt.Printf("Telemetry Optout: %t\n", !data.TelemetryOptout.IsNull() && !data.TelemetryOptout.IsUnknown())
		fmt.Printf("Debug Mode: %t\n", !data.DebugMode.IsNull() && !data.DebugMode.IsUnknown())

		fmt.Println("\n==== Values Mapped to Provider Data Model ====")
		fmt.Printf("Tenant ID Length: %d\n", len(data.TenantID.ValueString()))
		fmt.Printf("Auth Method: %s\n", data.AuthMethod.ValueString())
		fmt.Printf("Client ID Length: %d\n", len(data.ClientID.ValueString()))
		fmt.Printf("Client Secret Length: %d\n", len(data.ClientSecret.ValueString()))
		fmt.Printf("Client Certificate Base64 Length: %d\n", len(data.ClientCertificateBase64.ValueString()))
		fmt.Printf("Client Certificate File Path: %s\n", data.ClientCertificateFilePath.ValueString())
		fmt.Printf("Client Certificate Password Set: %t\n", data.ClientCertificatePassword.ValueString() != "")
		fmt.Printf("Username Set: %t\n", data.Username.ValueString() != "")
		fmt.Printf("Password Set: %t\n", data.Password.ValueString() != "")
		fmt.Printf("Redirect URL: %s\n", data.RedirectURL.ValueString())
		fmt.Printf("Use Proxy: %t\n", data.UseProxy.ValueBool())
		fmt.Printf("Proxy URL: %s\n", data.ProxyURL.ValueString())
		fmt.Printf("Cloud: %s\n", data.Cloud.ValueString())
		fmt.Printf("Enable Chaos: %t\n", data.EnableChaos.ValueBool())
		fmt.Printf("Telemetry Optout: %t\n", data.TelemetryOptout.ValueBool())
		fmt.Printf("Debug Mode: %t\n", data.DebugMode.ValueBool())

		fmt.Println("========================================")
	}
}
