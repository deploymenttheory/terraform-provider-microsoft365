package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// obtainCredential performs the necessary steps to obtain a TokenCredential based on the provider configuration.
// It uses the CredentialFactory and CredentialStrategy to create the appropriate credential type based on the authentication method
// defined within the provider configuraton.
func obtainCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Obtaining credential", map[string]interface{}{
		"auth_method": config.AuthMethod.ValueString(),
	})
	strategy, err := credentialFactory(config.AuthMethod.ValueString())
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

func credentialFactory(authMethod string) (CredentialStrategy, error) {
	tflog.Info(context.Background(), "Creating credential strategy", map[string]interface{}{
		"auth_method": authMethod,
	})
	switch authMethod {
	case "azure_developer_cli":
		return &AzureDeveloperCLIStrategy{}, nil
	case "client_secret":
		return &ClientSecretStrategy{}, nil
	case "client_certificate":
		return &ClientCertificateStrategy{}, nil
	case "device_code":
		return &DeviceCodeStrategy{}, nil
	case "interactive_browser":
		return &InteractiveBrowserStrategy{}, nil
	case "workload_identity":
		return &WorkloadIdentityStrategy{}, nil
	case "managed_identity":
		return &ManagedIdentityStrategy{}, nil
	case "oidc":
		return &OIDCStrategy{}, nil
	case "oidc_github":
		return &GitHubOIDCStrategy{}, nil
	case "oidc_azure_devops":
		return &AzureDevOpsOIDCStrategy{}, nil
	default:
		tflog.Error(context.Background(), "Unsupported authentication method", map[string]interface{}{
			"auth_method": authMethod,
		})
		return nil, fmt.Errorf("unsupported authentication method: %s", authMethod)
	}
}

// AzureDeveloperCLIStrategy implements the credential strategy for Azure Developer CLI authentication
type AzureDeveloperCLIStrategy struct{}

func (s *AzureDeveloperCLIStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating Azure Developer CLI credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
	})

	options := &azidentity.AzureDeveloperCLICredentialOptions{
		TenantID:                   config.TenantID.ValueString(),
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
	}

	return azidentity.NewAzureDeveloperCLICredential(options)
}

// CredentialStrategy defines the interface for credential creation strategies
type CredentialStrategy interface {
	GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error)
}

// ClientSecretStrategy implements the credential strategy for client secret authentication
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

// ClientCertificateStrategy implements the credential strategy for client certificate authentication
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

// DeviceCodeStrategy implements the credential strategy for device code authentication
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

// InteractiveBrowserStrategy implements the credential strategy for interactive browser authentication
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

// WorkloadIdentityStrategy implements the credential strategy for workload identity authentication
type WorkloadIdentityStrategy struct{}

func (s *WorkloadIdentityStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating workload identity credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})

	options := &azidentity.WorkloadIdentityCredentialOptions{
		ClientOptions:              clientOptions,
		ClientID:                   entraIDOptions.ClientID.ValueString(),
		TenantID:                   config.TenantID.ValueString(),
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
	}

	if !entraIDOptions.FederatedTokenFilePath.IsNull() && entraIDOptions.FederatedTokenFilePath.ValueString() != "" {
		options.TokenFilePath = entraIDOptions.FederatedTokenFilePath.ValueString()
		tflog.Debug(ctx, "Using Kubernetes service account token file path for workload identity authentication ", map[string]interface{}{
			"path": options.TokenFilePath,
		})
	}

	return azidentity.NewWorkloadIdentityCredential(options)
}

// ManagedIdentityStrategy implements the credential strategy for managed identity authentication
type ManagedIdentityStrategy struct{}

func (s *ManagedIdentityStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating managed identity credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
	})

	options := &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientOptions,
	}

	if !entraIDOptions.ManagedIdentityID.IsNull() && entraIDOptions.ManagedIdentityID.ValueString() != "" {
		idValue := entraIDOptions.ManagedIdentityID.ValueString()

		if strings.HasPrefix(idValue, "/subscriptions/") {
			options.ID = azidentity.ResourceID(idValue)
			tflog.Debug(ctx, "Using user-assigned managed identity with Resource ID", map[string]interface{}{
				"resource_id": idValue,
			})
		} else {
			options.ID = azidentity.ClientID(idValue)
			tflog.Debug(ctx, "Using user-assigned managed identity with Client ID", map[string]interface{}{
				"client_id": idValue,
			})
		}
	} else {
		tflog.Debug(ctx, "Using system-assigned managed identity")
	}

	return azidentity.NewManagedIdentityCredential(options)
}

// OIDCStrategy implements a minimalist generic credential strategy for OIDC authentication
type OIDCStrategy struct{}

func (s *OIDCStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating generic OIDC credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})

	var getAssertion func(context.Context) (string, error)

	if !entraIDOptions.OIDCTokenFilePath.IsNull() && entraIDOptions.OIDCTokenFilePath.ValueString() != "" {
		tokenPath := entraIDOptions.OIDCTokenFilePath.ValueString()
		tflog.Debug(ctx, "Using OIDC token file", map[string]interface{}{
			"path": tokenPath,
		})

		getAssertion = func(ctx context.Context) (string, error) {
			tokenBytes, err := os.ReadFile(tokenPath)
			if err != nil {
				return "", fmt.Errorf("failed to read OIDC token file: %w", err)
			}

			token := strings.TrimSpace(string(tokenBytes))
			if token == "" {
				return "", fmt.Errorf("OIDC token file is empty: %s", tokenPath)
			}

			return token, nil
		}
	}

	options := &azidentity.ClientAssertionCredentialOptions{
		ClientOptions:              clientOptions,
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
	}

	return azidentity.NewClientAssertionCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		getAssertion,
		options)
}

// GitHubOIDCStrategy implements the credential strategy for GitHub Actions OIDC authentication
type GitHubOIDCStrategy struct{}

func (s *GitHubOIDCStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating GitHub OIDC credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})

	// Check if GitHub Actions environment variables are set
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	if requestURL == "" || requestToken == "" {
		return nil, fmt.Errorf("GitHub Actions OIDC environment variables ACTIONS_ID_TOKEN_REQUEST_URL and ACTIONS_ID_TOKEN_REQUEST_TOKEN are required. " +
			"Ensure your workflow has 'permissions: id-token: write' set")
	}

	// Check if a custom audience is specified and add to request URL if not already present
	audience := helpers.GetEnvString("M365_OIDC_AUDIENCE",
		helpers.GetEnvString("ARM_OIDC_AUDIENCE", "api://AzureADTokenExchange"))

	tflog.Debug(ctx, "Using audience for GitHub OIDC token", map[string]interface{}{
		"audience": audience,
	})

	if !strings.Contains(requestURL, "audience=") {
		separator := "&"
		if !strings.Contains(requestURL, "?") {
			separator = "?"
		}
		requestURL = requestURL + separator + "audience=" + url.QueryEscape(audience)
	}

	getAssertion := func(ctx context.Context) (string, error) {
		// GitHub Actions provides an endpoint to request the JWT token
		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request for GitHub OIDC token: %w", err)
		}

		req.Header.Add("Authorization", "Bearer "+requestToken)
		req.Header.Add("Accept", "application/json; api-version=2.0")
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to request GitHub OIDC token: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("failed to get GitHub OIDC token, status code: %d, response: %s",
				resp.StatusCode, string(body))
		}

		var tokenResp struct {
			Value string `json:"value"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			return "", fmt.Errorf("failed to decode GitHub OIDC token response: %w", err)
		}

		if tokenResp.Value == "" {
			return "", fmt.Errorf("received empty OIDC token from GitHub")
		}

		return tokenResp.Value, nil
	}

	options := &azidentity.ClientAssertionCredentialOptions{
		ClientOptions:              clientOptions,
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
	}

	return azidentity.NewClientAssertionCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		getAssertion,
		options)
}

// AzureDevOpsOIDCStrategy implements the credential strategy for Azure DevOps OIDC authentication
type AzureDevOpsOIDCStrategy struct{}

func (s *AzureDevOpsOIDCStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating Azure DevOps OIDC credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})

	serviceConnectionID := entraIDOptions.ADOServiceConnectionID.ValueString()
	if serviceConnectionID == "" {
		return nil, fmt.Errorf("azure devops service connection id is required for azure devops OIDC authentication")
	}

	adoRequestToken := os.Getenv("SYSTEM_ACCESSTOKEN")
	if adoRequestToken == "" {
		return nil, fmt.Errorf("SYSTEM_ACCESSTOKEN environment variable is required for azure devops OIDC authentication")
	}

	oidcRequestURI := os.Getenv("SYSTEM_OIDCREQUESTURI")
	if oidcRequestURI == "" {
		return nil, fmt.Errorf("SYSTEM_OIDCREQUESTURI environment variable is required for azure devops OIDC authentication")
	}

	tflog.Debug(ctx, "Using Azure DevOps OIDC request URI", map[string]interface{}{
		"oidc_request_uri": oidcRequestURI,
	})

	options := &azidentity.AzurePipelinesCredentialOptions{
		ClientOptions:              clientOptions,
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
	}

	return azidentity.NewAzurePipelinesCredential(
		config.TenantID.ValueString(),
		entraIDOptions.ClientID.ValueString(),
		serviceConnectionID,
		adoRequestToken,
		options)
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
