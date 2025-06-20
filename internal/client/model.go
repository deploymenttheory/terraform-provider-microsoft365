package client

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// ProviderData represents the data needed to configure the Microsoft Graph clients.
// This is a simplified version of the provider model that only includes the fields
// needed for client configuration.
type ProviderData struct {
	// Cloud is the Microsoft cloud environment to use (public, dod, gcc, gcchigh, china, etc.)
	Cloud string
	// TenantID is the Microsoft 365 tenant ID
	TenantID string
	// AuthMethod is the authentication method to use (client_secret, client_certificate, etc.)
	AuthMethod string
	// EntraIDOptions contains options for Entra ID authentication
	EntraIDOptions *EntraIDOptions
	// ClientOptions contains options for the Microsoft Graph client
	ClientOptions *ClientOptions
	// TelemetryOptout indicates whether to opt out of telemetry
	TelemetryOptout bool
	// DebugMode indicates whether debug mode is enabled
	DebugMode bool
}

// EntraIDOptions represents the options for Entra ID authentication
type EntraIDOptions struct {
	// ClientID is the client ID (application ID) for the Entra ID application
	ClientID string
	// ClientSecret is the client secret for the Entra ID application
	ClientSecret string
	// ClientCertificate is the path to the client certificate file
	ClientCertificate string
	// ClientCertificatePassword is the password for the client certificate
	ClientCertificatePassword string
	// Username is the username for interactive authentication
	Username string
	// RedirectUrl is the redirect URL for interactive authentication
	RedirectUrl string
	// FederatedTokenFilePath is the path to the federated token file
	FederatedTokenFilePath string
	// ManagedIdentityClientID is the client ID for managed identity authentication
	ManagedIdentityClientID string
	// ManagedIdentityResourceID is the resource ID for managed identity authentication
	ManagedIdentityResourceID string
	// OIDCTokenFilePath is the path to the OIDC token file
	OIDCTokenFilePath string
	// OIDCToken is the OIDC token
	OIDCToken string
	// OIDCRequestToken is the OIDC request token
	OIDCRequestToken string
	// OIDCRequestURL is the OIDC request URL
	OIDCRequestURL string
	// DisableInstanceDiscovery indicates whether to disable instance discovery
	DisableInstanceDiscovery bool
	// SendCertificateChain indicates whether to send the certificate chain
	SendCertificateChain bool
	// AdditionallyAllowedTenants is a list of additionally allowed tenants
	AdditionallyAllowedTenants []string
}

// ClientOptions represents the options for the Microsoft Graph client
type ClientOptions struct {
	// EnableRetry indicates whether to enable retry
	EnableRetry bool
	// MaxRetries is the maximum number of retries
	MaxRetries int64
	// RetryDelaySeconds is the delay between retries in seconds
	RetryDelaySeconds int64
	// EnableRedirect indicates whether to enable redirect
	EnableRedirect bool
	// MaxRedirects is the maximum number of redirects
	MaxRedirects int64
	// EnableCompression indicates whether to enable compression
	EnableCompression bool
	// CustomUserAgent is the custom user agent
	CustomUserAgent string
	// EnableHeadersInspection indicates whether to enable headers inspection
	EnableHeadersInspection bool
	// TimeoutSeconds is the timeout in seconds
	TimeoutSeconds int64
	// UseProxy indicates whether to use a proxy
	UseProxy bool
	// ProxyURL is the proxy URL
	ProxyURL string
	// ProxyUsername is the proxy username
	ProxyUsername string
	// ProxyPassword is the proxy password
	ProxyPassword string
	// EnableChaos indicates whether to enable chaos
	EnableChaos bool
	// ChaosPercentage is the chaos percentage
	ChaosPercentage int64
	// ChaosStatusCode is the chaos status code
	ChaosStatusCode int64
	// ChaosStatusMessage is the chaos status message
	ChaosStatusMessage string
}

// GetClientOptions returns the Azure SDK client options based on the provider data
func (d *ProviderData) GetClientOptions() policy.ClientOptions {
	return policy.ClientOptions{}
}
