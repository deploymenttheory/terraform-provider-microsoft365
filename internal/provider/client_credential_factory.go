package provider

import (
	"context"
	"encoding/base64"
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
	case "workload_identity":
		return &WorkloadIdentityStrategy{}, nil
	case "managed_identity":
		return &ManagedIdentityStrategy{}, nil
	case "oidc":
		return &OIDCStrategy{}, nil
	default:
		tflog.Error(context.Background(), "Unsupported authentication method", map[string]interface{}{
			"auth_method": authMethod,
		})
		return nil, fmt.Errorf("unsupported authentication method: %s", authMethod)
	}
}

// obtainCredential performs the necessary steps to obtain a TokenCredential based on the provider configuration.
// It uses the CredentialFactory and CredentialStrategy to create the appropriate credential type based on the authentication method
// defined within the provider configuraton.
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

// UsernamePasswordStrategy implements the credential strategy for username/password authentication
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

	// Check if a specific managed identity is specified
	if !entraIDOptions.ManagedIdentityID.IsNull() && entraIDOptions.ManagedIdentityID.ValueString() != "" {
		idValue := entraIDOptions.ManagedIdentityID.ValueString()

		// Determine ID type based on format - we let the schema validation handle the format check
		if strings.HasPrefix(idValue, "/subscriptions/") {
			// Resource ID format
			options.ID = azidentity.ResourceID(idValue)
			tflog.Debug(ctx, "Using user-assigned managed identity with Resource ID", map[string]interface{}{
				"resource_id": idValue,
			})
		} else {
			// Client ID format (GUID) - already validated by schema
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

// OIDCStrategy implements the credential strategy for OIDC authentication with federated credentials
type OIDCStrategy struct{}

func (s *OIDCStrategy) GetCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions) (azcore.TokenCredential, error) {
	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	tflog.Info(ctx, "Creating OIDC credential", map[string]interface{}{
		"tenant_id": config.TenantID.ValueString(),
		"client_id": entraIDOptions.ClientID.ValueString(),
	})

	// Option 1: Direct OIDC token
	if oidcToken := helpers.GetEnvString("ARM_OIDC_TOKEN", ""); oidcToken != "" {
		tflog.Debug(ctx, "Using direct OIDC token from environment variable")
		getAssertion := func(context.Context) (string, error) {
			return oidcToken, nil
		}

		return createClientAssertionCredential(ctx, config, clientOptions, getAssertion)
	}

	// Option 2: GitHub Actions OIDC
	// GitHub automatically sets these environment variables when OIDC is enabled
	requestURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	requestToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")

	if requestURL != "" && requestToken != "" {
		tflog.Debug(ctx, "Using GitHub Actions OIDC token request")

		// Check if a custom audience is specified
		audience := helpers.GetEnvString("ARM_OIDC_AUDIENCE", "api://AzureADTokenExchange")
		if audience != "" {
			// Add audience parameter to request URL if not already present
			if !strings.Contains(requestURL, "audience=") {
				separator := "&"
				if !strings.Contains(requestURL, "?") {
					separator = "?"
				}
				requestURL = requestURL + separator + "audience=" + url.QueryEscape(audience)
			}
			tflog.Debug(ctx, "Using custom audience for OIDC token", map[string]interface{}{
				"audience": audience,
			})
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

			// If debug mode is enabled, decode and log token claims (without sensitive parts)
			if config.DebugMode.ValueBool() {
				logTokenDetails(ctx, tokenResp.Value)
			}

			return tokenResp.Value, nil
		}

		return createClientAssertionCredential(ctx, config, clientOptions, getAssertion)
	}

	// Option 3: Azure DevOps OIDC
	adoRequestToken := helpers.GetFirstEnvString([]string{"SYSTEM_ACCESSTOKEN"}, "")
	serviceConnectionID := helpers.GetFirstEnvString(
		[]string{"ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "ARM_OIDC_AZURE_SERVICE_CONNECTION_ID"},
		entraIDOptions.ADOServiceConnectionID.ValueString())

	if adoRequestToken != "" && serviceConnectionID != "" {
		tflog.Debug(ctx, "Using Azure DevOps OIDC", map[string]interface{}{
			"service_connection_id": serviceConnectionID,
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

	// Option 4: Manual OIDC token file
	if !entraIDOptions.OIDCTokenFilePath.IsNull() && entraIDOptions.OIDCTokenFilePath.ValueString() != "" {
		tokenPath := entraIDOptions.OIDCTokenFilePath.ValueString()
		tflog.Debug(ctx, "Using OIDC token file", map[string]interface{}{
			"path": tokenPath,
		})

		getAssertion := func(ctx context.Context) (string, error) {
			tokenBytes, err := os.ReadFile(tokenPath)
			if err != nil {
				return "", fmt.Errorf("failed to read OIDC token file: %w", err)
			}

			token := strings.TrimSpace(string(tokenBytes))
			if token == "" {
				return "", fmt.Errorf("OIDC token file is empty: %s", tokenPath)
			}

			// If debug mode is enabled, decode and log token claims (without sensitive parts)
			if config.DebugMode.ValueBool() {
				logTokenDetails(ctx, token)
			}

			return token, nil
		}

		return createClientAssertionCredential(ctx, config, clientOptions, getAssertion)
	}

	return nil, fmt.Errorf("no OIDC token source configured - requires one of: " +
		"GitHub Actions OIDC environment (ensure 'permissions: id-token: write' is set in your workflow), " +
		"Azure DevOps pipeline with service connection, " +
		"ARM_OIDC_TOKEN environment variable, " +
		"or oidc_token_file_path configuration")
}

// Helper function to create ClientAssertionCredential
func createClientAssertionCredential(ctx context.Context, config *M365ProviderModel, clientOptions policy.ClientOptions,
	getAssertion func(context.Context) (string, error)) (azcore.TokenCredential, error) {

	var entraIDOptions EntraIDOptionsModel
	config.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

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

// Helper function to log token details for debugging (without logging the full token)
func logTokenDetails(ctx context.Context, token string) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		tflog.Debug(ctx, "Invalid JWT token format - expected 3 parts")
		return
	}

	// Decode the header (first part)
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		tflog.Debug(ctx, "Failed to decode JWT header", map[string]interface{}{"error": err.Error()})
		return
	}

	// Decode the payload (second part)
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		tflog.Debug(ctx, "Failed to decode JWT payload", map[string]interface{}{"error": err.Error()})
		return
	}

	var header, payload map[string]interface{}

	if err := json.Unmarshal(headerJSON, &header); err != nil {
		tflog.Debug(ctx, "Failed to parse JWT header", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		tflog.Debug(ctx, "Failed to parse JWT payload", map[string]interface{}{"error": err.Error()})
		return
	}

	// Log important fields for debugging without logging the entire token
	tflog.Debug(ctx, "OIDC Token Header", map[string]interface{}{
		"alg": header["alg"],
		"typ": header["typ"],
		"kid": header["kid"],
	})

	// Extract and log select non-sensitive claims
	debugInfo := map[string]interface{}{}

	safeClaimsToLog := []string{"iss", "aud", "sub", "exp", "nbf", "iat", "jti"}

	// Add GitHub-specific claims when present
	githubClaimsToLog := []string{
		"repository", "repository_owner", "repository_visibility",
		"workflow", "ref", "ref_type", "actor", "event_name",
		"environment", "job_workflow_ref", "runner_environment",
	}
	safeClaimsToLog = append(safeClaimsToLog, githubClaimsToLog...)

	for _, claim := range safeClaimsToLog {
		if value, ok := payload[claim]; ok {
			debugInfo[claim] = value
		}
	}

	tflog.Debug(ctx, "OIDC Token Claims", debugInfo)
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
