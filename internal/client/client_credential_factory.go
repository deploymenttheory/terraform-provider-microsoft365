package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ObtainCredential performs the necessary steps to obtain a TokenCredential based on the provider configuration.
// It uses the CredentialFactory and CredentialStrategy to create the appropriate credential type based on the authentication method
// defined within the provider configuraton.
func ObtainCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Obtaining credential", map[string]any{
		"auth_method": config.AuthMethod,
	})
	strategy, err := credentialFactory(config.AuthMethod)
	if err != nil {
		tflog.Error(ctx, "Failed to create credential strategy", map[string]any{
			"error": err,
		})
		return nil, err
	}
	credential, err := strategy.GetCredential(ctx, config, clientOptions)
	if err != nil {
		tflog.Error(ctx, "Failed to get credential", map[string]any{
			"error": err,
		})
		return nil, err
	}
	tflog.Info(ctx, "Successfully obtained credential")
	return credential, nil
}

func credentialFactory(authMethod string) (CredentialStrategy, error) {
	tflog.Info(context.Background(), "Creating credential strategy", map[string]any{
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
		tflog.Error(context.Background(), "Unsupported authentication method", map[string]any{
			"auth_method": authMethod,
		})
		return nil, fmt.Errorf("unsupported authentication method: %s", authMethod)
	}
}

// AzureDeveloperCLIStrategy implements the credential strategy for Azure Developer CLI authentication
type AzureDeveloperCLIStrategy struct{}

func (s *AzureDeveloperCLIStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating Azure Developer CLI credential", map[string]any{
		"tenant_id": config.TenantID,
	})

	options := &azidentity.AzureDeveloperCLICredentialOptions{
		TenantID:                   config.TenantID,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}

	return azidentity.NewAzureDeveloperCLICredential(options)
}

// CredentialStrategy defines the interface for credential creation strategies
type CredentialStrategy interface {
	GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error)
}

// ClientSecretStrategy implements the credential strategy for client secret authentication
type ClientSecretStrategy struct{}

func (s *ClientSecretStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating client secret credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})
	return azidentity.NewClientSecretCredential(
		config.TenantID,
		config.EntraIDOptions.ClientID,
		config.EntraIDOptions.ClientSecret,
		&azidentity.ClientSecretCredentialOptions{
			AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
			DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		})
}

// ClientCertificateStrategy implements the credential strategy for client certificate authentication
type ClientCertificateStrategy struct{}

func (s *ClientCertificateStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating client certificate credential", map[string]any{
		"tenant_id":   config.TenantID,
		"client_id":   config.EntraIDOptions.ClientID,
		"certificate": config.EntraIDOptions.ClientCertificate,
	})

	certData, err := os.ReadFile(config.EntraIDOptions.ClientCertificate)
	if err != nil {
		tflog.Error(ctx, "Failed to read certificate file", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	password := []byte(config.EntraIDOptions.ClientCertificatePassword)
	certs, privateKey, err := helpers.ParseCertificateData(ctx, certData, password)
	if err != nil {
		tflog.Error(ctx, "Failed to parse certificate data", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to parse certificate data: %w", err)
	}

	return azidentity.NewClientCertificateCredential(
		config.TenantID,
		config.EntraIDOptions.ClientID,
		certs,
		privateKey,
		&azidentity.ClientCertificateCredentialOptions{
			AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
			DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
			SendCertificateChain:       config.EntraIDOptions.SendCertificateChain,
		})
}

// DeviceCodeStrategy implements the credential strategy for device code authentication
type DeviceCodeStrategy struct{}

func (s *DeviceCodeStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating device code credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})
	return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: config.TenantID,
		ClientID: config.EntraIDOptions.ClientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			tflog.Info(ctx, "Device code message", map[string]any{
				"message": message.Message,
			})
			fmt.Println(message.Message)
			return nil
		},
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	})
}

// InteractiveBrowserStrategy implements the credential strategy for interactive browser authentication
type InteractiveBrowserStrategy struct{}

func (s *InteractiveBrowserStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating interactive browser credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})
	options := &azidentity.InteractiveBrowserCredentialOptions{
		ClientID:                   config.EntraIDOptions.ClientID,
		TenantID:                   config.TenantID,
		RedirectURL:                config.EntraIDOptions.RedirectUrl,
		LoginHint:                  config.EntraIDOptions.Username,
		ClientOptions:              clientOptions,
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}
	return azidentity.NewInteractiveBrowserCredential(options)
}

// WorkloadIdentityStrategy implements the credential strategy for workload identity authentication
type WorkloadIdentityStrategy struct{}

func (s *WorkloadIdentityStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating workload identity credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})

	options := &azidentity.WorkloadIdentityCredentialOptions{
		ClientOptions:              clientOptions,
		ClientID:                   config.EntraIDOptions.ClientID,
		TenantID:                   config.TenantID,
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}

	if config.EntraIDOptions.FederatedTokenFilePath != "" {
		options.TokenFilePath = config.EntraIDOptions.FederatedTokenFilePath
		tflog.Debug(ctx, "Using Kubernetes service account token file path for workload identity authentication ", map[string]any{
			"path": options.TokenFilePath,
		})
	}

	return azidentity.NewWorkloadIdentityCredential(options)
}

// ManagedIdentityStrategy implements the credential strategy for managed identity authentication
type ManagedIdentityStrategy struct{}

func (s *ManagedIdentityStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating managed identity credential", map[string]any{
		"client_id":   config.EntraIDOptions.ManagedIdentityClientID,
		"resource_id": config.EntraIDOptions.ManagedIdentityResourceID,
	})
	options := &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientOptions,
		ID:            azidentity.ClientID(config.EntraIDOptions.ManagedIdentityClientID),
	}

	if config.EntraIDOptions.ManagedIdentityResourceID != "" {
		options.ID = azidentity.ResourceID(config.EntraIDOptions.ManagedIdentityResourceID)
		tflog.Debug(ctx, "Using resource ID for managed identity authentication", map[string]any{
			"resource_id": config.EntraIDOptions.ManagedIdentityResourceID,
		})
	}

	return azidentity.NewManagedIdentityCredential(options)
}

// OIDCStrategy implements the credential strategy for OIDC authentication
type OIDCStrategy struct{}

func (s *OIDCStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating OIDC credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})

	options := &azidentity.ClientAssertionCredentialOptions{
		ClientOptions:              clientOptions,
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}

	var assertion func(ctx context.Context) (string, error)

	if config.EntraIDOptions.OIDCTokenFilePath != "" {
		tflog.Debug(ctx, "Using OIDC token file", map[string]any{
			"path": config.EntraIDOptions.OIDCTokenFilePath,
		})
		assertion = func(ctx context.Context) (string, error) {
			token, err := os.ReadFile(config.EntraIDOptions.OIDCTokenFilePath)
			if err != nil {
				return "", fmt.Errorf("failed to read OIDC token file: %w", err)
			}
			return string(token), nil
		}
	} else if config.EntraIDOptions.OIDCToken != "" {
		tflog.Debug(ctx, "Using OIDC token from configuration")
		assertion = func(ctx context.Context) (string, error) {
			return config.EntraIDOptions.OIDCToken, nil
		}
	} else if config.EntraIDOptions.OIDCRequestToken != "" && config.EntraIDOptions.OIDCRequestURL != "" {
		tflog.Debug(ctx, "Using OIDC token exchange", map[string]any{
			"url": config.EntraIDOptions.OIDCRequestURL,
		})
		assertion = func(ctx context.Context) (string, error) {
			return exchangeOIDCToken(ctx, config.EntraIDOptions.OIDCRequestURL, config.EntraIDOptions.OIDCRequestToken)
		}
	} else {
		return nil, fmt.Errorf("OIDC authentication requires either oidc_token_file_path, oidc_token, or both oidc_request_token and oidc_request_url to be set")
	}

	return azidentity.NewClientAssertionCredential(
		config.TenantID,
		config.EntraIDOptions.ClientID,
		assertion,
		options,
	)
}

// GitHubOIDCStrategy implements the credential strategy for GitHub OIDC authentication
type GitHubOIDCStrategy struct{}

func (s *GitHubOIDCStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating GitHub OIDC credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})

	options := &azidentity.ClientAssertionCredentialOptions{
		ClientOptions:              clientOptions,
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}

	// GitHub Actions provides the ID token via the ACTIONS_ID_TOKEN_REQUEST_URL
	// and ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variables
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	if requestURL == "" || requestToken == "" {
		return nil, fmt.Errorf("GitHub OIDC authentication requires the ACTIONS_ID_TOKEN_REQUEST_URL and ACTIONS_ID_TOKEN_REQUEST_TOKEN environment variables to be set. Got: ACTIONS_ID_TOKEN_REQUEST_URL='%s', ACTIONS_ID_TOKEN_REQUEST_TOKEN='%s'", requestURL, requestToken)
	}

	assertion := func(ctx context.Context) (string, error) {
		tflog.Debug(ctx, "Requesting GitHub OIDC token", map[string]any{
			"url":      requestURL,
			"audience": "api://AzureADTokenExchange",
		})

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
		if err != nil {
			return "", fmt.Errorf("getAssertion: failed to build request")
		}

		query, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			return "", fmt.Errorf("getAssertion: cannot parse URL query")
		}

		if query.Get("audience") == "" {
			query.Set("audience", "api://AzureADTokenExchange")
			req.URL.RawQuery = query.Encode()
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestToken))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("getAssertion: cannot request token: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return "", fmt.Errorf("getAssertion: cannot parse response: %v", err)
		}

		if c := resp.StatusCode; c < 200 || c > 299 {
			return "", fmt.Errorf("getAssertion: received HTTP status %d with response: %s", resp.StatusCode, body)
		}

		var tokenRes struct {
			Count *int    `json:"count"`
			Value *string `json:"value"`
		}
		if err := json.Unmarshal(body, &tokenRes); err != nil {
			return "", fmt.Errorf("getAssertion: cannot unmarshal response: %v", err)
		}

		if tokenRes.Value == nil {
			return "", fmt.Errorf("getAssertion: nil JWT assertion received from OIDC provider")
		}

		return *tokenRes.Value, nil
	}

	return azidentity.NewClientAssertionCredential(
		config.TenantID,
		config.EntraIDOptions.ClientID,
		assertion,
		options,
	)
}

// AzureDevOpsOIDCStrategy implements the credential strategy for Azure DevOps OIDC authentication
type AzureDevOpsOIDCStrategy struct{}

func (s *AzureDevOpsOIDCStrategy) GetCredential(ctx context.Context, config *ProviderData, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	tflog.Info(ctx, "Creating Azure DevOps OIDC credential", map[string]any{
		"tenant_id": config.TenantID,
		"client_id": config.EntraIDOptions.ClientID,
	})

	options := &azidentity.ClientAssertionCredentialOptions{
		ClientOptions:              clientOptions,
		DisableInstanceDiscovery:   config.EntraIDOptions.DisableInstanceDiscovery,
		AdditionallyAllowedTenants: config.EntraIDOptions.AdditionallyAllowedTenants,
	}

	// Azure DevOps provides the ID token via the AZURE_DEVOPS_FEDERATION_TOKEN environment variable
	federationToken := os.Getenv("AZURE_DEVOPS_FEDERATION_TOKEN")

	if federationToken == "" {
		return nil, fmt.Errorf("azure DevOps OIDC authentication requires the AZURE_DEVOPS_FEDERATION_TOKEN environment variable to be set")
	}

	// Create the assertion callback
	assertion := func(ctx context.Context) (string, error) {
		return federationToken, nil
	}

	return azidentity.NewClientAssertionCredential(
		config.TenantID,
		config.EntraIDOptions.ClientID,
		assertion,
		options,
	)
}

// exchangeOIDCToken exchanges an OIDC token for another token
func exchangeOIDCToken(ctx context.Context, requestURL, requestToken string) (string, error) {
	tflog.Debug(ctx, "Exchanging OIDC token", map[string]any{
		"url": requestURL,
	})

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for OIDC token exchange: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", requestToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange OIDC token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to exchange OIDC token: status code %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse OIDC token exchange response: %w", err)
	}

	return tokenResp.Value, nil
}
