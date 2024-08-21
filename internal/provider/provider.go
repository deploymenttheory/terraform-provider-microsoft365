package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &M365Provider{}

// M365Provider defines the provider implementation.
type M365Provider struct {
	version string
	clients *client.GraphClients
}

// M365ProviderModel describes the provider data model.
type M365ProviderModel struct {
	TenantID                  types.String `tfsdk:"tenant_id"`
	AuthMethod                types.String `tfsdk:"auth_method"`
	ClientID                  types.String `tfsdk:"client_id"`
	ClientSecret              types.String `tfsdk:"client_secret"`
	ClientCertificate         types.String `tfsdk:"client_certificate"`
	ClientCertificatePassword types.String `tfsdk:"client_certificate_password"`
	Username                  types.String `tfsdk:"username"`
	Password                  types.String `tfsdk:"password"`
	RedirectURL               types.String `tfsdk:"redirect_url"`
	UseProxy                  types.Bool   `tfsdk:"use_proxy"`
	ProxyURL                  types.String `tfsdk:"proxy_url"`
	Cloud                     types.String `tfsdk:"cloud"`
	EnableChaos               types.Bool   `tfsdk:"enable_chaos"`
	TelemetryOptout           types.Bool   `tfsdk:"telemetry_optout"`
	DebugMode                 types.Bool   `tfsdk:"debug_mode"`
}

func (p *M365Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365"
	resp.Version = p.version
}

func (p *M365Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud": schema.StringAttribute{
				Description: "Specifies the Microsoft cloud environment for authentication and API requests. " +
					"This setting determines the endpoints used for Microsoft Graph and Graph Beta APIs. " +
					"Default is 'public'. Can be set using the `M365_CLOUD` environment variable.",
				MarkdownDescription: "Specifies the Microsoft cloud environment for authentication and API requests. " +
					"This setting determines the endpoints used for Microsoft Graph and Graph Beta APIs.\n\n" +
					"Valid values:\n" +
					"- `public`: Microsoft Azure Public Cloud (default)\n" +
					"- `dod`: US Department of Defense (DoD) Cloud\n" +
					"- `gcc`: US Government Cloud\n" +
					"- `gcchigh`: US Government High Cloud\n" +
					"- `china`: China Cloud\n" +
					"- `ex`: EagleX Cloud\n" +
					"- `rx`: Secure Cloud (RX)\n\n" +
					"Can be set using the `M365_CLOUD` environment variable.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("public", "dod", "gcc", "gcchigh", "china", "ex", "rx"),
				},
			},
			"auth_method": schema.StringAttribute{
				Required: true,
				Description: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', " +
					"'username_password'. Each method requires different credentials to be provided. " +
					"Can also be set using the `M365_AUTH_METHOD` environment variable.",
				MarkdownDescription: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options:\n" +
					"- `device_code`: Uses a device code flow for authentication.\n" +
					"- `client_secret`: Uses a client ID and secret for authentication.\n" +
					"- `client_certificate`: Uses a client certificate (.pfx) for authentication.\n" +
					"- `interactive_browser`: Opens a browser for interactive login.\n" +
					"- `username_password`: Uses username and password for authentication (not recommended for production).\n" +
					"Each method requires different credentials to be provided. Can also be set using the `M365_AUTH_METHOD` environment variable.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"device_code", "client_secret", "client_certificate", "interactive_browser", "username_password", "client_assertion",
					),
				},
			},
			"tenant_id": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Description: "The Microsoft 365 tenant ID for the Entra ID (formerly Azure AD) application. " +
					"This GUID uniquely identifies your Entra ID instance. " +
					"Can be set using the `M365_TENANT_ID` environment variable.",
				MarkdownDescription: "The Microsoft 365 tenant ID for the Entra ID (formerly Azure AD) application. " +
					"This GUID uniquely identifies your Entra ID instance.\n\n" +
					"To find your tenant ID:\n" +
					"1. Log in to the [Azure portal](https://portal.azure.com)\n" +
					"2. Navigate to 'Microsoft Entra ID' (formerly Azure Active Directory)\n" +
					"3. In the Overview page, look for 'Tenant ID'\n\n" +
					"Alternatively, you can use PowerShell:\n" +
					"```powershell\n" +
					"(Get-AzureADTenantDetails).ObjectId\n" +
					"```\n\n" +
					"Or Azure CLI:\n" +
					"```bash\n" +
					"az account show --query tenantId -o tsv\n" +
					"```\n\n" +
					"Can be set using the `M365_TENANT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID("tenant_id"),
				},
			},
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
					"Can be set using the `M365_CLIENT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID("client_id"),
				},
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client secret for the Entra ID application. Required for client credentials authentication. " +
					"This secret is generated in Entra ID and has an expiration date. " +
					"Can be set using the `M365_CLIENT_SECRET` environment variable.",
				MarkdownDescription: "The client secret for the Entra ID (formerly Azure AD) application. " +
					"This secret is required for client credentials authentication flow.\n\n" +
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
					"Can be set using the `M365_CLIENT_SECRET` environment variable.",
			},
			"client_certificate": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The path to the Client Certificate file associated with the Service " +
					"Principal for use when authenticating as a Service Principal using a Client Certificate. " +
					"Supports PKCS#12 (.pfx or .p12) file format. The file should contain the certificate, " +
					"private key with an RSA type, and optionally a password which can be defined in client_certificate_password. ",
				MarkdownDescription: "The path to the client certificate file for certificate-based authentication with Entra ID (formerly Azure AD). " +
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
					"Can be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.",
			},
			"client_certificate_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The password to decrypt the PKCS#12 (.pfx or .p12) client certificate file. " +
					"Required only if the certificate file is password-protected. " +
					"Can be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.",
				MarkdownDescription: "The password to decrypt the PKCS#12 (.pfx or .p12) client certificate file. " +
					"This is required only if the certificate file is password-protected.\n\n" +
					"Important notes:\n" +
					"- This password is used to encrypt the private key in the certificate file\n" +
					"- It's not related to any Entra ID settings, but to the certificate file itself\n" +
					"- If your PKCS#12 file was created without a password, leave this field empty\n" +
					"- Treat this password with the same level of security as the certificate itself\n\n" +
					"When creating a PKCS#12 file with OpenSSL, you'll be prompted for this password:\n" +
					"```bash\n" +
					"openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem\n" +
					"```\n\n" +
					"Can be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.",
			},
			"username": schema.StringAttribute{
				Optional: true,
				Description: "The username for resource owner password credentials (ROPC) flow. " +
					"Can be set using the `M365_USERNAME` environment variable.",
				MarkdownDescription: "The username for resource owner password credentials (ROPC) flow authentication.\n\n" +
					"**Important Security Notice:**\n" +
					"- Resource Owner Password Credentials (ROPC) is considered less secure than other authentication methods\n" +
					"- It should only be used when other, more secure methods are not possible\n" +
					"- Not recommended for production environments\n" +
					"- Does not support multi-factor authentication\n\n" +
					"Usage:\n" +
					"- Typically, this is the user's email address or User Principal Name (UPN)\n" +
					"- Ensure the user has appropriate permissions for the required operations\n\n" +
					"Can be set using the `M365_USERNAME` environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The password for resource owner password credentials (ROPC) flow. " +
					"Can be set using the `M365_PASSWORD` environment variable.",
				MarkdownDescription: "The password for resource owner password credentials (ROPC) flow authentication.\n\n" +
					"**Critical Security Warning:**\n" +
					"- Storing passwords in plain text is a significant security risk\n" +
					"- Use environment variables or secure vaults to manage this sensitive information\n" +
					"- Regularly rotate passwords and monitor for unauthorized access\n" +
					"- Consider using more secure authentication methods when possible\n\n" +
					"Can be set using the `M365_PASSWORD` environment variable.",
			},
			"redirect_url": schema.StringAttribute{
				Optional: true,
				Description: "The redirect URL for OAuth 2.0 authentication flows that require a callback URL. " +
					"Can be set using the `M365_REDIRECT_URL` environment variable.",
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
			"use_proxy": schema.BoolAttribute{
				Optional: true,
				Description: "Enables the use of an HTTP proxy for network requests. When set to true, the provider will " +
					"route all HTTP requests through the specified proxy server. This can be useful for environments that " +
					"require proxy access for internet connectivity or for monitoring and logging HTTP traffic. Can also be " +
					"set using the `M365_USE_PROXY` environment variable.",
				Validators: []validator.Bool{
					validateUseProxy(),
				},
			},
			"proxy_url": schema.StringAttribute{
				Optional: true,
				Description: "Specifies the URL of the HTTP proxy server. This URL should be in a valid URL format " +
					"(e.g., 'http://proxy.example.com:8080'). When 'use_proxy' is enabled, this URL is used to configure the " +
					"HTTP client to route requests through the proxy. Ensure the proxy server is reachable and correctly " +
					"configured to handle the network traffic. Can also be set using the `M365_PROXY_URL` environment variable.",
				MarkdownDescription: "Specifies the URL of the HTTP proxy server for routing requests when `use_proxy` is enabled.\n\n" +
					"**Format:**\n" +
					"- Must be a valid URL (e.g., `http://proxy.example.com:8080`)\n" +
					"- Supports HTTP and HTTPS protocols\n\n" +
					"**Usage:**\n" +
					"- When `use_proxy` is set to `true`, all HTTP(S) requests will be routed through this proxy\n" +
					"- Ensure the proxy server is reachable and correctly configured to handle the traffic\n\n" +
					"**Examples:**\n" +
					"- HTTP proxy: `http://proxy.example.com:8080`\n" +
					"- HTTPS proxy: `https://secure-proxy.example.com:443`\n" +
					"- Authenticated proxy: `http://username:password@proxy.example.com:8080`\n\n" +
					"**Security Considerations:**\n" +
					"- Use HTTPS for the proxy URL when possible to encrypt proxy communications\n" +
					"- If using an authenticated proxy, consider setting the URL via the environment variable to avoid exposing credentials in configuration files\n" +
					"- Ensure the proxy server is trusted and secure\n\n" +
					"Can be set using the `M365_PROXY_URL` environment variable.",
				Validators: []validator.String{
					validateProxyURL(),
				},
			},
			"enable_chaos": schema.BoolAttribute{
				Optional: true,
				Description: "Enable the chaos handler for testing purposes. " +
					"When enabled, the chaos handler can simulate specific failure scenarios " +
					"and random errors in API responses to help test the robustness and resilience " +
					"of the terraform provider against intermittent issues. This is particularly useful " +
					"for testing how the provider handles various error conditions and ensures " +
					"it can recover gracefully. Use with caution in production environments. " +
					"Can also be set using the `M365_ENABLE_CHAOS` environment variable.",
			},
			"telemetry_optout": schema.BoolAttribute{
				Optional: true,
				Description: "Flag to opt out of telemetry collection. Default is `false`. " +
					"Can be set using the `M365_TELEMETRY_OPTOUT` environment variable.",
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
			},
		},
	}
}

// Configure sets up the Microsoft365 provider with the given configuration.
// It processes the provider schema, retrieves values from the configuration or
// environment variables, sets up authentication, and initializes the Microsoft
// Graph clients.
//
// The function supports various authentication methods, proxy settings, and
// national cloud deployments. It performs the following main steps:
//  1. Extracts and validates the configuration data.
//  2. Sets up logging and context with relevant fields.
//  3. Determines cloud-specific constants and endpoints.
//  4. Configures the Entra ID client options.
//  5. Obtains credentials based on the specified authentication method.
//  6. Creates and configures the Microsoft Graph clients (stable and beta).
//
// If any errors occur during these steps, appropriate diagnostics are added
// to the response.
func (p *M365Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Microsoft365 Provider")

	var config M365ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting provider configuration", map[string]interface{}{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	data := M365ProviderModel{
		Cloud:                     types.StringValue(helpers.MultiEnvDefaultFunc([]string{"M365_CLOUD", "AZURE_CLOUD"}, config.Cloud.ValueString())),
		TenantID:                  types.StringValue(helpers.EnvDefaultFunc("M365_TENANT_ID", config.TenantID.ValueString())),
		AuthMethod:                types.StringValue(helpers.EnvDefaultFunc("M365_AUTH_METHOD", config.AuthMethod.ValueString())),
		ClientID:                  types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_ID", config.ClientID.ValueString())),
		ClientSecret:              types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_SECRET", config.ClientSecret.ValueString())),
		ClientCertificate:         types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_CERTIFICATE_FILE_PATH", config.ClientCertificate.ValueString())),
		ClientCertificatePassword: types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_CERTIFICATE_PASSWORD", config.ClientCertificatePassword.ValueString())),
		Username:                  types.StringValue(helpers.EnvDefaultFunc("M365_USERNAME", config.Username.ValueString())),
		Password:                  types.StringValue(helpers.EnvDefaultFunc("M365_PASSWORD", config.Password.ValueString())),
		RedirectURL:               types.StringValue(helpers.EnvDefaultFunc("M365_REDIRECT_URL", config.RedirectURL.ValueString())),
		UseProxy:                  types.BoolValue(helpers.EnvDefaultFuncBool("M365_USE_PROXY", config.UseProxy.ValueBool())),
		ProxyURL:                  types.StringValue(helpers.EnvDefaultFunc("M365_PROXY_URL", config.ProxyURL.ValueString())),
		EnableChaos:               types.BoolValue(helpers.EnvDefaultFuncBool("M365_ENABLE_CHAOS", config.EnableChaos.ValueBool())),
		TelemetryOptout:           types.BoolValue(helpers.EnvDefaultFuncBool("M365_TELEMETRY_OPTOUT", config.TelemetryOptout.ValueBool())),
		DebugMode:                 types.BoolValue(helpers.EnvDefaultFuncBool("M365_DEBUG_MODE", config.DebugMode.ValueBool())),
	}

	if data.DebugMode.ValueBool() {
		logDebugInfo(ctx, req, data)
	}

	ctx = tflog.SetField(ctx, "cloud", data.Cloud.ValueString())
	ctx = tflog.SetField(ctx, "auth_method", data.AuthMethod.ValueString())
	ctx = tflog.SetField(ctx, "use_proxy", data.UseProxy.ValueBool())
	ctx = tflog.SetField(ctx, "redirect_url", data.RedirectURL.ValueString())
	ctx = tflog.SetField(ctx, "proxy_url", data.ProxyURL.ValueString())
	ctx = tflog.SetField(ctx, "enable_chaos", data.EnableChaos.ValueBool())
	ctx = tflog.SetField(ctx, "telemetry_optout", data.TelemetryOptout.ValueBool())
	ctx = tflog.SetField(ctx, "debug_mode", data.DebugMode.ValueBool())

	ctx = tflog.SetField(ctx, "client_certificate_file_path", data.ClientCertificate.ValueString())
	ctx = tflog.SetField(ctx, "client_certificate_password", data.ClientCertificatePassword.ValueString())
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`(?i)client_certificate`))

	ctx = tflog.SetField(ctx, "username", data.Username.ValueString())
	ctx = tflog.SetField(ctx, "password", data.Password.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "password")

	ctx = tflog.SetField(ctx, "tenant_id", data.TenantID.ValueString())
	ctx = tflog.SetField(ctx, "client_id", data.ClientID.ValueString())
	ctx = tflog.SetField(ctx, "client_secret", data.ClientSecret.ValueString())
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "tenant_id")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_id")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_secret")

	authorityURL, apiScope, graphServiceRoot, graphBetaServiceRoot, err := setCloudConstants(data.Cloud.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Microsoft Cloud Type",
			fmt.Sprintf("An error occurred while attempting to get cloud constants for cloud type '%s'. "+
				"Please ensure the cloud type is valid. Detailed error: %s", data.Cloud.ValueString(), err.Error()),
		)
		return
	}

	ctx = tflog.SetField(ctx, "authority_url", authorityURL)
	ctx = tflog.SetField(ctx, "api_scope", apiScope)
	ctx = tflog.SetField(ctx, "graph_service_root", graphServiceRoot)
	ctx = tflog.SetField(ctx, "graph_beta_service_root", graphBetaServiceRoot)

	clientOptions, err := configureEntraIDClientOptions(
		ctx,
		data.UseProxy.ValueBool(),
		data.ProxyURL.ValueString(),
		authorityURL,
		data.TelemetryOptout.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure client options",
			fmt.Sprintf("An error occurred while attempting to configure client options. Detailed error: %s", err.Error()),
		)
		return
	}

	cred, err := obtainCredential(
		ctx,
		data,
		clientOptions,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create credentials",
			fmt.Sprintf("An error occurred while attempting to create the credentials using the provided authentication method '%s'. "+
				"This may be due to incorrect or missing credentials, misconfigured client options, or issues with the underlying authentication library. "+
				"Please verify the authentication method and credentials configuration. Detailed error: %s", data.AuthMethod.ValueString(), err.Error()),
		)
		return
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{apiScope},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create authentication provider",
			fmt.Sprintf("An error occurred while attempting to create the authentication provider using the provided credentials. "+
				"This may be due to misconfigured client options, incorrect credentials, or issues with the underlying authentication library. "+
				"Please verify your client options and credentials configuration. Detailed error: %s", err.Error()),
		)
		return
	}

	httpClient, err := configureGraphClientOptions(
		ctx,
		data.UseProxy.ValueBool(),
		data.ProxyURL.ValueString(),
		data.EnableChaos.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure Graph client options",
			fmt.Sprintf("An error occurred while attempting to configure the Microsoft Graph client options. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Stable SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Stable SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	betaAdapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Beta SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Beta SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter.SetBaseUrl(graphServiceRoot)
	betaAdapter.SetBaseUrl(graphBetaServiceRoot)

	clients := &client.GraphClients{
		StableClient: msgraphsdk.NewGraphServiceClient(stableAdapter),
		BetaClient:   msgraphbetasdk.NewGraphServiceClient(betaAdapter),
	}

	p.clients = clients

	resp.DataSourceData = clients
	resp.ResourceData = clients

	tflog.Debug(ctx, "Provider configuration completed", map[string]interface{}{
		"graph_client_set":      p.clients.StableClient != nil,
		"graph_beta_client_set": p.clients.BetaClient != nil,
		"config":                fmt.Sprintf("%+v", config),
	})
}

// New returns a function that, when invoked, creates and returns a new instance
// of the Microsoft365 provider, which implements the terraform-plugin-framework's
// provider.Provider interface. This function is designed to accept a version string,
// which is used to track the version of the provider being created.
//
// The provider internally manages two distinct Microsoft Graph clients:
// 1. StableClient: A client instance configured to interact with the stable version of the
//    Microsoft Graph API.
//
// 2. BetaClient: A client instance configured to interact with the beta version of the
//    Microsoft Graph API. This client is used for operations that require access to
//    newer or experimental features that are not yet available in the stable API.
//
// The New function encapsulates these clients within the M365Provider struct, which also
// holds the provider's configuration and resources. When Terraform invokes this function,
// it ensures that the provider is correctly instantiated with all necessary clients and
// configurations, making it ready to manage Microsoft365 resources through Terraform.

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		p := &M365Provider{
			version: version,
			clients: &client.GraphClients{},
		}
		return p
	}
}
