package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CredentialFactory creates the appropriate CredentialStrategy based on the authentication method
func CredentialFactory(authMethod string) (CredentialStrategy, error) {
	tflog.Info(context.Background(), "Creating credential strategy", map[string]interface{}{
		"auth_method": authMethod,
	})
	switch authMethod {
	case "client_secret":
		return &ClientSecretStrategy{}, nil
	case "client_certificate":
		return &ClientCertificateStrategy{}, nil
	case "username_password":
		return &UsernamePasswordStrategy{}, nil
	case "device_code":
		return &DeviceCodeStrategy{}, nil
	case "interactive_browser":
		return &InteractiveBrowserStrategy{}, nil
	default:
		tflog.Error(context.Background(), "Unsupported authentication method", map[string]interface{}{
			"auth_method": authMethod,
		})
		return nil, fmt.Errorf("unsupported authentication method: %s", authMethod)
	}
}

// obtainCredential is now a wrapper that uses the CredentialFactory and CredentialStrategy
func obtainCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Obtaining credential", map[string]interface{}{
		"auth_method": config.AuthMethod.ValueString(),
	})
	strategy, err := CredentialFactory(config.AuthMethod.ValueString())
	if err != nil {
		tflog.Error(ctx, "Failed to create credential strategy", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}
	credential, err := strategy.GetCredential(ctx, config, clientOptions)
	if err != nil {
		tflog.Error(ctx, "Failed to get credential", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}
	tflog.Info(ctx, "Successfully obtained credential")
	return credential, nil
}

// CredentialStrategy defines the interface for credential creation strategies
type CredentialStrategy interface {
	GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error)
}

// ClientSecretStrategy implements CredentialStrategy for client secret authentication
type ClientSecretStrategy struct{}

func (s *ClientSecretStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating client secret credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})
	return azidentity.NewClientSecretCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		entraIDOptions.ClientSecret.ValueString(),
		&azidentity.ClientSecretCredentialOptions{
			AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
			DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		})
}

// ClientCertificateStrategy implements CredentialStrategy for client certificate authentication
type ClientCertificateStrategy struct{}

func (s *ClientCertificateStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating client certificate credential", map[string]interface{}{
		"tenant_id":   config.TenantID.ValueString(),
		"client_id":   entraIDOptions.ClientID.ValueString(),
		"certificate": entraIDOptions.ClientCertificate.ValueString(),
	})

	certData, err := os.ReadFile(entraIDOptions.ClientCertificate.ValueString())
	if err != nil {
		tflog.Error(ctx, "Failed to read certificate file", map[string]interface{}{
			"error": err,
		})
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	password := []byte(entraIDOptions.ClientCertificatePassword.ValueString())
	certs, privateKey, err := helpers.ParseCertificateData(ctx, certData, password)
	if err != nil {
		tflog.Error(ctx, "Failed to parse certificate data", map[string]interface{}{
			"error": err,
		})
		return nil, fmt.Errorf("failed to parse certificate data: %w", err)
	}

	return azidentity.NewClientCertificateCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		certs,
		privateKey,
		&azidentity.ClientCertificateCredentialOptions{
			AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
			DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
			SendCertificateChain:       entraIDOptions.SendCertificateChain.ValueBool(),
		})
}

// UsernamePasswordStrategy implements CredentialStrategy for username/password authentication
type UsernamePasswordStrategy struct{}

func (s *UsernamePasswordStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating username/password credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
		"username":  entraIDOptions.Username.ValueString(),
	})
	return azidentity.NewUsernamePasswordCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		entraIDOptions.Username.ValueString(),
		entraIDOptions.Password.ValueString(),
		&azidentity.UsernamePasswordCredentialOptions{
			AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
			DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		})
}

// DeviceCodeStrategy implements CredentialStrategy for device code authentication
type DeviceCodeStrategy struct{}

func (s *DeviceCodeStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating device code credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})
	return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: config.TenantID.ValueString(),
		ClientID: entraIDOptions.ClientID.ValueString(),
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			tflog.Info(ctx, "Device code message", map[string]interface{}{
				"message": message.Message,
			})
			fmt.Println(message.Message)
			return nil
		},
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
	})
}

// InteractiveBrowserStrategy implements CredentialStrategy for interactive browser authentication
type InteractiveBrowserStrategy struct{}

func (s *InteractiveBrowserStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating interactive browser credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})
	options := &azidentity.InteractiveBrowserCredentialOptions{
		ClientID:                   entraIDOptions.ClientID.ValueString(),
		TenantID:                   config.TenantID.ValueString(),
		RedirectURL:                entraIDOptions.RedirectUrl.ValueString(),
		LoginHint:                  entraIDOptions.Username.ValueString(),
		ClientOptions:              clientOptions,
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
	}
	return azidentity.NewInteractiveBrowserCredential(options)
}

// Helper function to convert types.List to []string for AdditionallyAllowedTenants
func getAdditionallyAllowedTenants(tenants types.List) []string {
	var result []string
	for _, tenant := range tenants.Elements() {
		if strVal, ok := tenant.(types.String); ok {
			result = append(result, strVal.ValueString())
		}
	}
	return result
}
