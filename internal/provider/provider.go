package provider

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &M365Provider{}

// M365Provider defines the provider implementation.
type M365Provider struct {
	version      string
	clients      client.GraphClientInterface
	unitTestMode bool
}

func (p *M365Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365"
	resp.Version = p.version
}

func (p *M365Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Specifies the Microsoft cloud environment for authentication and API requests. " +
					"This setting determines the endpoints used for Microsoft Graph and Graph Beta APIs. " +
					"Valid values:\n" +
					"- `public`: Microsoft Azure Public Cloud (default)\n" +
					"- `dod`: US Department of Defense (DoD) Cloud\n" +
					"- `gcc`: US Government Cloud\n" +
					"- `gcchigh`: US Government High Cloud\n" +
					"- `china`: China Cloud\n" +
					"- `ex`: EagleX Cloud\n" +
					"- `rx`: Secure Cloud (RX)\n\n" +
					"Can be set using the `M365_CLOUD` environment variable.",
				Validators: []validator.String{
					stringvalidator.OneOf("public", "dod", "gcc", "gcchigh", "china", "ex", "rx"),
				},
			},
			"tenant_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "The Microsoft 365 tenant ID for the Entra ID (formerly Azure AD) application. " +
					"This GUID uniquely identifies your Entra ID instance." +
					"Can be set using the `M365_TENANT_ID` environment variable.\n\n" +
					"To find your tenant ID:\n" +
					"1. Log in to the [Azure portal](https://portal.azure.com)\n" +
					"2. Navigate to 'Microsoft Entra ID' (formerly Azure Active Directory)\n" +
					"3. In the Overview page, look for 'Tenant ID'\n\n" +
					"Alternatively, you can use PowerShell:\n" +
					"```powershell\n" +
					"Connect-AzAccount\n" +
					"(Get-AzContext).Tenant.Id\n" +
					"```\n\n" +
					"Can also be set using the `M365_TENANT_ID` environment variable.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"auth_method": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options:\n" +
					"- `azure_developer_cli`: Uses the identity logged into the Azure Developer CLI (azd) for authentication. Ideal for local Terraform development when you're already authenticated with azd.\n" +
					"- `device_code`: Uses a device code flow for authentication.\n" +
					"- `client_secret`: Uses a client ID and secret for authentication.\n" +
					"- `client_certificate`: Uses a client certificate (.pfx) for authentication.\n" +
					"- `interactive_browser`: Opens a browser for interactive login.\n" +
					"- `workload_identity`: Uses workload identity federation for Kubernetes pods, enabling them to authenticate via a service account token file.\n" +
					"- `managed_identity`: Uses Azure managed identity for authentication when Terraform is running on an Azure resource (like a VM, Azure Container Instance, or App Service) that has been assigned a managed identity.\n" +
					"- `oidc`: Uses generic OpenID Connect (OIDC) authentication with a JWT token from a file or environment variable.\n" +
					"- `oidc_github`: Uses GitHub Actions-specific OIDC authentication, with support for subject claims that specify repositories, branches, tags, pull requests, and environments for fine-grained trust configurations.\n" +
					"- `oidc_azure_devops`: Uses Azure DevOps-specific OIDC authentication with service connections, supporting federated credentials for secure pipeline-to-cloud authentication without storing secrets.\n" +
					"Each method requires different credentials to be provided.\n" +
					"Can also be set using the `M365_AUTH_METHOD` environment variable.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"azure_developer_cli",
						"client_secret",
						"client_certificate",
						"interactive_browser",
						"device_code",
						"workload_identity",
						"managed_identity",
						"oidc",
						"oidc_github",
						"oidc_azure_devops",
					),
				},
			},
			"entra_id_options": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Configuration options for Entra ID authentication.",
				Attributes:  EntraIDOptionsSchema(),
			},
			"client_options": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Configuration options for the Microsoft Graph client.",
				Attributes:  ClientOptionsSchema(),
			},
			"telemetry_optout": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Controls the collection of telemetry data for the Microsoft 365 provider by Microsoft Services.\n\n" +
					"**Usage:**\n" +
					"- Set to `true` to disable all telemetry collection\n" +
					"- Set to `false` (default) to allow telemetry collection\n\n" +
					"**Behavior:**\n" +
					"- When set to `true`, it prevents the addition of any telemetry data to API requests\n" +
					"- This affects the User-Agent string and other potential telemetry mechanisms\n\n" +
					"**Privacy:**\n" +
					"- Telemetry, when enabled, may include provider version, Terraform version, and general usage patterns\n" +
					"- No personally identifiable information (PII) or sensitive data is collected\n\n" +
					"**Recommendations:**\n" +
					"- For development or non-sensitive environments, consider leaving telemetry enabled to support product improvement\n" +
					"- For production or sensitive environments, you may choose to opt out\n\n" +
					"Can be set using the `M365_TELEMETRY_OPTOUT` environment variable.",
			},
			"debug_mode": schema.BoolAttribute{
				Optional: true,
				Description: "Flag to enable debug mode for the provider." +
					"Can also be set using the `M365_DEBUG_MODE` environment variable.",
				MarkdownDescription: "Flag to enable debug mode for the provider.\n\n" +
					"This setting enables additional logging and diagnostics for the provider.\n\n" +
					"Can also be set using the `M365_DEBUG_MODE` environment variable.",
			},
		},
	}
}

func EntraIDOptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_id": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
			Description: "The client ID (application ID) for the Entra ID application. " +
				"This GUID is generated when you register an application in Entra ID. " +
				"Can be set using the `M365_CLIENT_ID` environment variable.",
			MarkdownDescription: "The client ID (application ID) for the Entra ID (formerly Azure AD) application. " +
				"This GUID is generated when you register an application in Entra ID.\n\n" +
				"To find or create a client ID:\n" +
				"1. Log in to the [Azure portal](https://portal.azure.com)\n" +
				"2. Navigate to 'Microsoft Entra ID' > 'App registrations'\n" +
				"3. Select your application or create a new one\n" +
				"4. The client ID is listed as 'Application (client) ID' in the Overview page\n\n" +
				"Using Azure CLI:\n" +
				"```bash\n" +
				"az ad app list --query \"[].{appId:appId, displayName:displayName}\"\n" +
				"```\n\n" +
				"Using Microsoft Graph PowerShell:\n" +
				"```powershell\n" +
				"Get-MgApplication -Filter \"displayName eq 'Your App Name'\" | Select-Object AppId, DisplayName\n" +
				"```\n\n" +
				"Example usage:\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  client_id = \"my_client_id\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CLIENT_ID` environment variable.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(constants.GuidRegex),
					"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
				),
			},
		},
		"client_secret": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
			MarkdownDescription: "Used for the 'client_secret' authentication method.\n\n" +
				"The client secret for the Entra ID application. Required for client credentials authentication. " +
				"This secret is generated in Entra ID and has an expiration date.\n\n" +
				"Important notes:\n" +
				"- Client secrets are sensitive and should be handled securely\n" +
				"- Secrets have an expiration date and need to be rotated periodically\n" +
				"- Use managed identities or certificate-based authentication when possible for improved security\n\n" +
				"To create a client secret:\n" +
				"1. Log in to the [Azure portal](https://portal.azure.com)\n" +
				"2. Navigate to 'Microsoft Entra ID' > 'App registrations'\n" +
				"3. Select your application\n" +
				"4. Go to 'Certificates & secrets' > 'Client secrets'\n" +
				"5. Click 'New client secret' and set a description and expiration\n" +
				"6. Copy the secret value immediately (it won't be shown again)\n\n" +
				"Using Azure CLI:\n" +
				"```bash\n" +
				"az ad app credential reset --id <app-id> --append\n" +
				"```\n\n" +
				"Using Microsoft Graph PowerShell:\n" +
				"```powershell\n" +
				"$credential = @{\n" +
				"    displayName = 'My Secret'\n" +
				"    endDateTime = (Get-Date).AddMonths(6)\n" +
				"}\n" +
				"New-MgApplicationPassword -ApplicationId <app-id> -PasswordCredential $credential\n" +
				"```\n\n" +
				"Example usage:\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  client_secret = \"my_client_secret\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CLIENT_SECRET` environment variable.",
		},
		"client_certificate": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
			MarkdownDescription: "Used for the 'client_certificate' authentication method.\n\n" +
				"The path to the Client Certificate file associated with the Service " +
				"Principal for use when authenticating as a Service Principal using a Client Certificate. " +
				"Supports PKCS#12 (.pfx or .p12) file format. The file should contain the certificate, " +
				"private key with an RSA type, and optionally a password which can be defined in client_certificate_password.\n\n" +
				"The path to the client certificate file for certificate-based authentication with Entra ID (formerly Azure AD). " +
				"This method is more secure than client secret-based authentication.\n\n" +
				"Requirements:\n" +
				"- File format: PKCS#12 (.pfx or .p12)\n" +
				"- Contents: Certificate, private key, and optionally a certificate chain\n" +
				"- The public key certificate must be uploaded to Entra ID\n\n" +
				"Steps to set up certificate authentication:\n" +
				"1. Generate a self-signed certificate or obtain one from a trusted Certificate Authority\n" +
				"2. Convert the certificate to PKCS#12 format if necessary\n" +
				"3. Upload the public key to Entra ID:\n" +
				"   - Go to Azure Portal > 'Microsoft Entra ID' > 'App registrations' > [Your App] > 'Certificates & secrets'\n" +
				"   - Upload the public key portion of your certificate\n" +
				"4. Provide the path to the PKCS#12 file in this attribute\n\n" +
				"Using OpenSSL to create a self-signed certificate:\n" +
				"```bash\n" +
				"openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365\n" +
				"openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem\n" +
				"```\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  client_certificate        = \"/path/to/cert.pfx\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.",
		},
		"client_certificate_password": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
			MarkdownDescription: "Used for the 'client_certificate' authentication method.\n\n" +
				"The password to decrypt the PKCS#12 (.pfx or .p12) client certificate file. " +
				"Required only if the certificate file is password-protected.\n\n" +
				"Important notes:\n" +
				"- This password is used to encrypt the private key in the certificate file\n" +
				"- It's not related to any Entra ID settings, but to the certificate file itself\n" +
				"- If your PKCS#12 file was created without a password, leave this field empty\n" +
				"- Treat this password with the same level of security as the certificate itself\n\n" +
				"When creating a PKCS#12 file with OpenSSL, you'll be prompted for this password:\n" +
				"```bash\n" +
				"openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem\n" +
				"```\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  client_certificate_password = \"certpassword\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.",
		},
		"send_certificate_chain": schema.BoolAttribute{
			Optional: true,
			MarkdownDescription: "Used for the 'client_certificate' authentication method.\n\n" +
				"Controls whether the credential sends the public certificate chain in the x5c header" +
				"of each token request's JWT. This is required for Subject Name/Issuer (SNI) authentication" +
				"and can be used in certain advanced scenarios. Defaults to false. Enable this if your" +
				"application uses certificate chain validation or if specifically instructed by Azure support.\n\n" +
				"**Key points:**\n" +
				"- Default value: `false`\n" +
				"- Set to `true` to enable sending the certificate chain\n\n" +
				"**Use cases:**\n" +
				"1. **Subject Name/Issuer (SNI) Authentication:** Required for SNI authentication scenarios.\n" +
				"2. **Enhanced Security:** Provides additional validation capabilities on the Entra ID side.\n" +
				"3. **Compatibility:** May be necessary for certain Azure services or configurations.\n\n" +
				"**How it works:**\n" +
				"- When enabled, the full X.509 certificate chain is included in the token request.\n" +
				"- This allows the Microsoft Entra ID (Azure AD) to perform additional validation on the certificate.\n" +
				"- It can help in scenarios where intermediate certificates need to be verified.\n\n" +
				"**Considerations:**\n" +
				"- Enabling this option increases the size of each token request.\n" +
				"- Only enable if you're certain your scenario requires it.\n" +
				"- Consult Azure documentation or support if you're unsure about enabling this option.\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  client_certificate        = \"/path/to/cert.pfx\"\n" +
				"  client_certificate_password = \"certpassword\"\n" +
				"  send_certificate_chain    = true\n" +
				"}\n" +
				"```\n\n" +
				"Only enable this option if you understand its implications or if specifically instructed by Azure support.",
		},
		"username": schema.StringAttribute{
			Optional: true,
			MarkdownDescription: "Used for the 'username_password' authentication method.\n\n" +
				"The username for resource owner password credentials (ROPC) flow authentication.\n\n" +
				"**Important Security Notice:**\n" +
				"- Resource Owner Password Credentials (ROPC) is considered less secure than other authentication methods\n" +
				"- It should only be used when other, more secure methods are not possible\n" +
				"- Not recommended for production environments\n" +
				"- Does not support multi-factor authentication\n\n" +
				"Usage:\n" +
				"- Typically, this is the user's email address or User Principal Name (UPN)\n" +
				"- Ensure the user has appropriate permissions for the required operations\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  username        = \"user_name\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_USERNAME` environment variable.",
		},
		"disable_instance_discovery": schema.BoolAttribute{
			Optional: true,
			MarkdownDescription: "DisableInstanceDiscovery should be set true only by terraform runs" +
				"authenticating in disconnected clouds, or private clouds such as Azure Stack." +
				"It determines whether the credential requests Microsoft Entra instance metadata" +
				"from https://login.microsoft.com before authenticating. Setting this to true will" +
				"skip this request, making the application responsible for ensuring the configured" +
				"authority is valid and trustworthy.\n\n" +
				"Can be set using the `M365_DISABLE_INSTANCE_DISCOVERY` environment variable.",
		},
		"additionally_allowed_tenants": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			MarkdownDescription: "Specifies additional tenants for which the credential may acquire tokens." +
				"Add the wildcard value '*' to allow the credential to acquire tokens for any tenant.\n\n" +
				"Can be set using the `M365_ADDITIONALLY_ALLOWED_TENANTS` environment variable.",
		},
		"redirect_url": schema.StringAttribute{
			Optional: true,
			Description: "The redirect URL for OAuth 2.0 authentication flows that require a callback URL. " +
				"Can be set using the `M365_REDIRECT_URI` environment variable.",
			MarkdownDescription: "The redirect URL (also known as reply URL or callback URL) for OAuth 2.0 authentication flows that require a callback, such as the Authorization Code flow or interactive browser authentication.\n\n" +
				"**Important:**\n" +
				"- This URL must be registered in your Entra ID (formerly Azure AD) application settings\n" +
				"- For local development, typically use `http://localhost:port`\n" +
				"- For production, use a secure HTTPS URL\n\n" +
				"To configure in Entra ID:\n" +
				"1. Go to Azure Portal > 'Microsoft Entra ID' > 'App registrations'\n" +
				"2. Select your application\n" +
				"3. Go to 'Authentication' > 'Platform configurations'\n" +
				"4. Add or update the redirect URI\n\n" +
				"Security considerations:\n" +
				"- Use a unique path for your redirect URL to prevent potential conflicts\n" +
				"- Avoid using wildcard URLs in production environments\n" +
				"- Regularly audit and remove any unused redirect URLs\n\n" +
				"Example values:\n" +
				"- Local development: `http://localhost:8000/auth/callback`\n" +
				"- Production: `https://yourdomain.com/auth/microsoft365/callback`\n\n" +
				"Can be set using the `M365_REDIRECT_URL` environment variable.",
			Validators: []validator.String{
				validateRedirectURL(),
			},
		},
		"federated_token_file_path": schema.StringAttribute{
			Optional:    true,
			Description: "Path to a file containing a Kubernetes service account token for workload identity authentication.",
			MarkdownDescription: "Path to a file containing a Kubernetes service account token for workload identity authentication. " +
				"This field is only used with the 'workload_identity' authentication method.\n\n" +
				"In Kubernetes environments with Azure workload identity enabled, this path is typically " +
				"'/var/run/secrets/azure/tokens/azure-identity-token'. This token file is used to establish " +
				"federated identity for your workloads running in Kubernetes.\n\n" +
				"Can be set using the `AZURE_FEDERATED_TOKEN_FILE` environment variable.",
		},
		"managed_identity_id": schema.StringAttribute{
			Optional:    true,
			Description: "ID of a user-assigned managed identity to authenticate with.",
			MarkdownDescription: "ID of a user-assigned managed identity to authenticate with. This field is only used with the " +
				"'managed_identity' authentication method.\n\n" +
				"If omitted, the system-assigned managed identity will be used. If specified, it can be either:\n" +
				"- Client ID (GUID): The client ID of the user-assigned managed identity\n" +
				"- Resource ID: The full Azure resource ID of the user-assigned managed identity in the format " +
				"`/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}`\n\n" +
				"**Note:** Not all Azure hosting environments support all ID types. Some environments may have restrictions on " +
				"using certain ID formats. If you encounter errors, try using a different ID format or consult the Azure documentation " +
				"for your specific hosting environment.\n\n" +
				"Can be set using the `AZURE_CLIENT_ID` or `M365_MANAGED_IDENTITY_ID` environment variables.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^(([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})|(/subscriptions/[^/]+/resourceGroups/[^/]+/providers/Microsoft\.ManagedIdentity/userAssignedIdentities/[^/]+))$`),
					"must be either a valid GUID (client ID) or a valid Azure resource ID for a managed identity",
				),
			},
		},
		"oidc_token_file_path": schema.StringAttribute{
			Optional:    true,
			Description: "Path to a file containing an OIDC token for authentication.",
			MarkdownDescription: "Path to a file containing an OIDC token for authentication. This field is only used with the " +
				"'oidc' authentication method.\n\n" +
				"The file should contain a valid JWT assertion that will be used to authenticate the application. " +
				"This is commonly used in CI/CD pipelines or other environments that support OIDC federation with Azure AD.\n\n" +
				"Can be set using the `M365_OIDC_TOKEN_FILE_PATH` environment variable.",
		},
		"ado_service_connection_id": schema.StringAttribute{
			Optional:    true,
			Description: "Azure DevOps service connection ID for OIDC authentication.",
			MarkdownDescription: "Azure DevOps service connection ID for OIDC authentication. This field is only used with the " +
				"'oidc' authentication method when using Azure DevOps Pipelines.\n\n" +
				"Can be set using the `ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID` or `ARM_OIDC_AZURE_SERVICE_CONNECTION_ID` environment variables.",
		},
	}
}

func ClientOptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		// ----  TODO  set these with default values
		"enable_headers_inspection": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Enable inspection of HTTP headers.",
		},
		"enable_retry": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Enable automatic retries for failed requests.",
		},
		"max_retries": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Maximum number of retries for failed requests.",
		},
		"retry_delay_seconds": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Delay between retry attempts in seconds.",
		},
		"enable_redirect": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Enable automatic following of redirects.",
		},
		"max_redirects": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Maximum number of redirects to follow.",
		},
		"enable_compression": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Enable compression for HTTP requests and responses.",
		},
		// ----    set these with default values
		"custom_user_agent": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Custom User-Agent string to be sent with requests.",
		}, // TODO - set one for the provider
		"use_proxy": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable the use of a proxy for network requests.",
			MarkdownDescription: "Enables the use of a proxy server for all network requests made by the provider.\n\n" +
				"**Key points:**\n" +
				"- Default: `false`\n" +
				"- When `true`, the provider will route all HTTP requests through the specified proxy server\n" +
				"- Requires `proxy_url` to be set\n" +
				"- Useful for environments that require proxy access for internet connectivity\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  use_proxy = true\n" +
				"  proxy_url = \"http://proxy.example.com:8080\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_USE_PROXY` environment variable.",
		},
		"proxy_url": schema.StringAttribute{
			Optional:    true,
			Description: "The URL of the proxy server.",
			MarkdownDescription: "Specifies the URL of the proxy server to be used when `use_proxy` is set to `true`.\n\n" +
				"**Key points:**\n" +
				"- Must be a valid URL including the scheme (http:// or https://)\n" +
				"- Can include a port number\n" +
				"- Required when `use_proxy` is `true`\n" +
				"- Ignored if `use_proxy` is `false`\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  use_proxy = true\n" +
				"  proxy_url = \"http://proxy.example.com:8080\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_PROXY_URL` environment variable.",
			Validators: []validator.String{
				validateProxyURL(),
			},
		},
		"proxy_username": schema.StringAttribute{
			Optional:    true,
			Description: "Username for proxy authentication.",
			MarkdownDescription: "Specifies the username for authentication with the proxy server if required.\n\n" +
				"**Key points:**\n" +
				"- Optional: Only needed if your proxy server requires authentication\n" +
				"- Used in conjunction with `proxy_password`\n" +
				"- Ignored if `use_proxy` is `false` or if `proxy_url` is not set\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  use_proxy      = true\n" +
				"  proxy_url      = \"http://proxy.example.com:8080\"\n" +
				"  proxy_username = \"proxyuser\"\n" +
				"  proxy_password = \"proxypass\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_PROXY_USERNAME` environment variable.",
		},
		"proxy_password": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "Password for proxy authentication.",
			MarkdownDescription: "Specifies the password for authentication with the proxy server if required.\n\n" +
				"**Key points:**\n" +
				"- Optional: Only needed if your proxy server requires authentication\n" +
				"- Used in conjunction with `proxy_username`\n" +
				"- Treated as sensitive information and will be masked in logs\n" +
				"- Ignored if `use_proxy` is `false` or if `proxy_url` is not set\n\n" +
				"**Security note:** It's recommended to set this using an environment variable rather than in the configuration file.\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  use_proxy      = true\n" +
				"  proxy_url      = \"http://proxy.example.com:8080\"\n" +
				"  proxy_username = \"proxyuser\"\n" +
				"  proxy_password = \"proxypass\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_PROXY_PASSWORD` environment variable.",
		},
		"timeout_seconds": schema.Int64Attribute{
			Optional:    true,
			Description: "Override value for the timeout of authentication requests in seconds.",
		},
		"enable_chaos": schema.BoolAttribute{
			Optional:    true,
			Description: "Enable the chaos handler for testing purposes. When enabled, it simulates failure scenarios and random errors in API responses.",
			MarkdownDescription: "Enable the chaos handler for testing purposes. " +
				"When enabled, the chaos handler simulates specific failure scenarios " +
				"and random errors in API responses to help test the robustness and resilience " +
				"of the terraform provider against intermittent issues. This is particularly useful " +
				"for testing how the provider handles various error conditions and ensures " +
				"it can recover gracefully.\n\n" +
				"**Key points:**\n" +
				"- Default: `false`\n" +
				"- When `true`, adds a chaos handler to the middleware\n" +
				"- Injects an 'X-Chaos-Injected: true' header in affected responses\n" +
				"- Use with caution, especially in production environments\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  enable_chaos = true\n" +
				"  chaos_percentage = 20\n" +
				"}\n" +
				"```\n\n" +
				"Can also be set using the `M365_ENABLE_CHAOS` environment variable.",
		},
		"chaos_percentage": schema.Int64Attribute{
			Optional:    true,
			Description: "Percentage of requests to apply chaos testing to. Must be between 0 and 100.",
			MarkdownDescription: "Specifies the percentage of requests that should be affected by the chaos handler.\n\n" +
				"**Key points:**\n" +
				"- Value range: 0 to 100\n" +
				"- Default: 10% if not specified\n" +
				"- Determines the probability of injecting chaos into each request\n" +
				"- Higher values increase the frequency of simulated errors\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  enable_chaos = true\n" +
				"  chaos_percentage = 30  # 30% of requests will be affected\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CHAOS_PERCENTAGE` environment variable.",
		},
		"chaos_status_code": schema.Int64Attribute{
			Optional:    true,
			Description: "HTTP status code to return for chaos-affected requests. If not set, a random error status code will be used.",
			MarkdownDescription: "Specifies the HTTP status code to be returned for requests affected by the chaos handler.\n\n" +
				"**Key points:**\n" +
				"- If not set, a random error status code will be used\n" +
				"- Common error codes: 429 (Too Many Requests), 500 (Internal Server Error), 503 (Service Unavailable)\n" +
				"- Used only when `enable_chaos` is true\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  enable_chaos = true\n" +
				"  chaos_status_code = 503  # Simulate a 'Service Unavailable' error\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CHAOS_STATUS_CODE` environment variable.",
		},
		"chaos_status_message": schema.StringAttribute{
			Optional:    true,
			Description: "Custom status message to return for chaos-affected requests. If not set, a default message will be used.",
			MarkdownDescription: "Defines a custom status message to be returned for requests affected by the chaos handler.\n\n" +
				"**Key points:**\n" +
				"- If not set, a default message will be used\n" +
				"- Allows simulation of specific error messages or conditions\n" +
				"- Used only when `enable_chaos` is true\n\n" +
				"**Example usage:**\n" +
				"```hcl\n" +
				"provider \"microsoft365\" {\n" +
				"  enable_chaos = true\n" +
				"  chaos_status_message = \"Simulated server overload\"\n" +
				"}\n" +
				"```\n\n" +
				"Can be set using the `M365_CHAOS_STATUS_MESSAGE` environment variable.",
		},
	}
}

// NewMicrosoft365Provider returns a function that, when invoked, creates and returns a new instance
// of the Microsoft365 provider, which implements the terraform-plugin-framework's
// provider.Provider interface. This function is designed to accept a version string,
// which is used to track the version of the provider being created.
//
// The provider internally manages two distinct Microsoft Graph clients:
//
//  1. StableClient: A client instance configured to interact with the stable version of the
//     Microsoft Graph API.
//
//  2. BetaClient: A client instance configured to interact with the beta version of the
//     Microsoft Graph API. This client is used for operations that require access to
//     newer or experimental features that are not yet available in the stable API.
//
// The New function encapsulates these clients within the M365Provider struct, which also
// holds the provider's configuration and resources. When Terraform invokes this function,
// it ensures that the provider is correctly instantiated with all necessary clients and
// configurations, making it ready to manage Microsoft365 resources through Terraform.
func NewMicrosoft365Provider(version string, unitTestMode ...bool) func() provider.Provider {
	return func() provider.Provider {
		isUnitTestMode := false
		if len(unitTestMode) > 0 {
			isUnitTestMode = unitTestMode[0]
		}
		// Initialize with nil clients - will be set during Configure clients step.
		p := &M365Provider{
			version:      version,
			clients:      nil,
			unitTestMode: isUnitTestMode,
		}
		return p
	}
}
