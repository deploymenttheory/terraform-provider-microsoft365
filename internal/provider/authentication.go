package provider

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// createCredential creates an Azure credential based on the provider configuration.
func createCredential(ctx context.Context, data M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
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
			"tenant_id":        data.TenantID.ValueString(),
			"client_id":        data.ClientID.ValueString(),
			"certificate_path": data.ClientCertificateFilePath.ValueString(),
		})

		certs, key, err := helpers.GetCertificatesAndKeyFromCertOrFilePath(data.ClientCertificateFilePath.ValueString(), data.ClientCertificatePassword.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to get certificates and key from path '%s': %s", data.ClientCertificateFilePath.ValueString(), err.Error())
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
