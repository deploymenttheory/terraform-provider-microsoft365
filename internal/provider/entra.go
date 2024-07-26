// REF: https://learn.microsoft.com/en-us/graph/sdks/choose-authentication-providers?tabs=go#on-behalf-of-provider
package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// configureEntraIDClientOptions configures the client options for Entra ID
func configureEntraIDClientOptions(ctx context.Context, useProxy bool, proxyURL string, authorityURL string, telemetryOptout bool) (policy.ClientOptions, error) {
	tflog.Debug(ctx, "Configuring Entra ID client options")

	clientOptions := policy.ClientOptions{}

	if useProxy && proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			return clientOptions, fmt.Errorf("failed to parse the provided proxy URL '%s': %s", proxyURL, err.Error())
		}

		authClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURLParsed),
			},
		}

		clientOptions.Transport = authClient
	}

	clientOptions.Cloud = cloud.Configuration{
		ActiveDirectoryAuthorityHost: authorityURL,
	}

	if telemetryOptout {
		clientOptions.Telemetry.Disabled = true
	}

	clientOptions.Logging = policy.LogOptions{
		IncludeBody: true,
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowedQueryParams: []string{
			"api-version",
		},
	}

	clientOptions.Retry = policy.RetryOptions{
		MaxRetries:    5,
		RetryDelay:    2 * time.Second,
		MaxRetryDelay: 30 * time.Second,
		StatusCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
	tflog.Debug(ctx, "Configured Entra ID client options")

	return clientOptions, nil
}

// obtainCredential creates an Azure credential based on the provider configuration.
func obtainCredential(ctx context.Context, data M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	switch data.AuthMethod.ValueString() {
	case "device_code":
		tflog.Debug(ctx, "Obtaining Device Code Credential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})
		return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID: data.TenantID.ValueString(),
			ClientID: data.ClientID.ValueString(),
			UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
				tflog.Info(ctx, message.Message)
				return nil
			},
			ClientOptions: clientOptions,
		})
	case "client_secret":
		tflog.Debug(ctx, "Obtaining Client Secret Credential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})
		return azidentity.NewClientSecretCredential(data.TenantID.ValueString(), data.ClientID.ValueString(), data.ClientSecret.ValueString(), &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "client_certificate":
		tflog.Debug(ctx, "Obtaining Client Certificate Credential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})

		if data.ClientCertificateFilePath.IsNull() {
			return nil, fmt.Errorf("'client_certificate_file_path' must be provided for client_certificate authentication")
		}

		tflog.Debug(ctx, "Using client certificate file path")
		certData, err := os.ReadFile(data.ClientCertificateFilePath.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to read certificate file: %v", err)
		}

		password := []byte(data.ClientCertificatePassword.ValueString())
		certs, key, err := helpers.ParseCertificateData(ctx, certData, password)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificates: %v", err)
		}

		return azidentity.NewClientCertificateCredential(
			data.TenantID.ValueString(),
			data.ClientID.ValueString(),
			certs,
			key,
			&azidentity.ClientCertificateCredentialOptions{
				ClientOptions: clientOptions,
			})
	case "interactive_browser":
		tflog.Debug(ctx, "Obtaining Interactive Browser Credential", map[string]interface{}{
			"tenant_id":    data.TenantID.ValueString(),
			"client_id":    data.ClientID.ValueString(),
			"redirect_url": data.RedirectURL.ValueString(),
		})
		redirectURL := data.RedirectURL.ValueString()
		return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
			TenantID:      data.TenantID.ValueString(),
			ClientID:      data.ClientID.ValueString(),
			RedirectURL:   redirectURL,
			ClientOptions: clientOptions,
		})
	case "username_password":
		tflog.Debug(ctx, "Obtaining Username / Password Credential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})
		username := data.Username.ValueString()
		password := data.Password.ValueString()
		return azidentity.NewUsernamePasswordCredential(data.TenantID.ValueString(), data.ClientID.ValueString(), username, password, &azidentity.UsernamePasswordCredentialOptions{
			ClientOptions: clientOptions,
		})
	default:
		return nil, fmt.Errorf("unsupported authentication method '%s'", data.AuthMethod.ValueString())
	}
}
