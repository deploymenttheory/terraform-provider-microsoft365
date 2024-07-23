package provider

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// configureEntraIDClientOptions configures the client options for Entra ID
func configureEntraIDClientOptions(useProxy bool, proxyURL string, authorityURL string, telemetryOptout bool) (policy.ClientOptions, error) {
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

	return clientOptions, nil
}

// obtainCredential creates an Azure credential based on the provider configuration.
func obtainCredential(ctx context.Context, data M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	switch data.AuthMethod.ValueString() {
	case "device_code":
		tflog.Debug(ctx, "Creating DeviceCodeCredential", map[string]interface{}{
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
		tflog.Debug(ctx, "Creating ClientSecretCredential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})
		return azidentity.NewClientSecretCredential(data.TenantID.ValueString(), data.ClientID.ValueString(), data.ClientSecret.ValueString(), &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "client_certificate":
		tflog.Debug(ctx, "Creating ClientCertificateCredential", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})

		var certs []*x509.Certificate
		var key interface{}
		var err error

		if !data.ClientCertificate.IsNull() {
			tflog.Debug(ctx, "Using base64 encoded client certificate")
			certs, key, err = helpers.GetCertificatesAndKeyFromCertOrFilePath(data.ClientCertificate.ValueString(), data.ClientCertificatePassword.ValueString())
		} else if !data.ClientCertificateFilePath.IsNull() {
			tflog.Debug(ctx, "Using client certificate file path")
			certs, key, err = helpers.GetCertificatesAndKeyFromCertOrFilePath(data.ClientCertificateFilePath.ValueString(), data.ClientCertificatePassword.ValueString())
		} else {
			return nil, fmt.Errorf("either 'client_certificate' or 'client_certificate_file_path' must be provided for client_certificate authentication")
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get certificates and key: %s", err.Error())
		}

		return azidentity.NewClientCertificateCredential(data.TenantID.ValueString(), data.ClientID.ValueString(), certs, key, &azidentity.ClientCertificateCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "on_behalf_of":
		tflog.Debug(ctx, "Creating OnBehalfOfCredentialWithSecret", map[string]interface{}{
			"tenant_id": data.TenantID.ValueString(),
			"client_id": data.ClientID.ValueString(),
		})
		userAssertion := data.UserAssertion.ValueString()
		return azidentity.NewOnBehalfOfCredentialWithSecret(data.TenantID.ValueString(), data.ClientID.ValueString(), userAssertion, data.ClientSecret.ValueString(), &azidentity.OnBehalfOfCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "interactive_browser":
		tflog.Debug(ctx, "Creating InteractiveBrowserCredential", map[string]interface{}{
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
		tflog.Debug(ctx, "Creating UsernamePasswordCredential", map[string]interface{}{
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
