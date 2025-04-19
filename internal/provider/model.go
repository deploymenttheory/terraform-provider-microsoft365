package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// M365ProviderModel describes the provider data model.
type M365ProviderModel struct {
	Cloud           types.String `tfsdk:"cloud"`
	TenantID        types.String `tfsdk:"tenant_id"`
	AuthMethod      types.String `tfsdk:"auth_method"`
	EntraIDOptions  types.Object `tfsdk:"entra_id_options"`
	ClientOptions   types.Object `tfsdk:"client_options"`
	TelemetryOptout types.Bool   `tfsdk:"telemetry_optout"`
	DebugMode       types.Bool   `tfsdk:"debug_mode"`
}

// ClientOptionsModel describes the client options
type ClientOptionsModel struct {
	EnableHeadersInspection types.Bool   `tfsdk:"enable_headers_inspection"`
	EnableRetry             types.Bool   `tfsdk:"enable_retry"`
	MaxRetries              types.Int64  `tfsdk:"max_retries"`
	RetryDelaySeconds       types.Int64  `tfsdk:"retry_delay_seconds"`
	EnableRedirect          types.Bool   `tfsdk:"enable_redirect"`
	MaxRedirects            types.Int64  `tfsdk:"max_redirects"`
	EnableCompression       types.Bool   `tfsdk:"enable_compression"`
	CustomUserAgent         types.String `tfsdk:"custom_user_agent"`
	UseProxy                types.Bool   `tfsdk:"use_proxy"`
	ProxyURL                types.String `tfsdk:"proxy_url"`
	ProxyUsername           types.String `tfsdk:"proxy_username"`
	ProxyPassword           types.String `tfsdk:"proxy_password"`
	TimeoutSeconds          types.Int64  `tfsdk:"timeout_seconds"`
	EnableChaos             types.Bool   `tfsdk:"enable_chaos"`
	ChaosPercentage         types.Int64  `tfsdk:"chaos_percentage"`
	ChaosStatusCode         types.Int64  `tfsdk:"chaos_status_code"`
	ChaosStatusMessage      types.String `tfsdk:"chaos_status_message"`
}

// EntraIDOptionsModel describes the Entra ID options
type EntraIDOptionsModel struct {
	ClientID                   types.String `tfsdk:"client_id"`
	ClientSecret               types.String `tfsdk:"client_secret"`
	ClientCertificate          types.String `tfsdk:"client_certificate"`
	ClientCertificatePassword  types.String `tfsdk:"client_certificate_password"`
	SendCertificateChain       types.Bool   `tfsdk:"send_certificate_chain"`
	Username                   types.String `tfsdk:"username"` // For Interactive Browser Credential
	DisableInstanceDiscovery   types.Bool   `tfsdk:"disable_instance_discovery"`
	AdditionallyAllowedTenants types.List   `tfsdk:"additionally_allowed_tenants"`
	RedirectUrl                types.String `tfsdk:"redirect_url"`
	FederatedTokenFilePath     types.String `tfsdk:"federated_token_file_path"` // For workload identity
	ManagedIdentityID          types.String `tfsdk:"managed_identity_id"`       // For managed identity
	OIDCTokenFilePath          types.String `tfsdk:"oidc_token_file_path"`      // For OIDC authentication
	ADOServiceConnectionID     types.String `tfsdk:"ado_service_connection_id"` // For Azure DevOps OIDC
}
