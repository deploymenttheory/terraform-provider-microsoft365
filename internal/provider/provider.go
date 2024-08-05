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
	ClientAssertion           types.String `tfsdk:"client_assertion"`
	ClientAssertionFile       types.String `tfsdk:"client_assertion_file"`
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
				Description: "The cloud to use for authentication and Graph / Graph Beta API requests." +
					"Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`." +
					"Can also be set using the `M365_CLOUD` environment variable.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("public", "gcc", "gcchigh", "china", "dod", "ex", "rx"),
				},
			},
			"auth_method": schema.StringAttribute{
				Required: true,
				Description: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', " +
					"'username_password', 'client_assertion'. Each method requires different credentials to be provided. " +
					"Can also be set using the `M365_AUTH_METHOD` environment variable.",
				MarkdownDescription: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options:\n" +
					"- `device_code`: Uses a device code flow for authentication.\n" +
					"- `client_secret`: Uses a client ID and secret for authentication.\n" +
					"- `client_certificate`: Uses a client certificate (.pfx) for authentication.\n" +
					"- `interactive_browser`: Opens a browser for interactive login.\n" +
					"- `username_password`: Uses username and password for authentication (not recommended for production).\n" +
					"- `client_assertion`: Uses a client assertion (OIDC token) for authentication, suitable for CI/CD and server-to-server scenarios.\n\n" +
					"Each method requires different credentials to be provided. Can also be set using the `M365_AUTH_METHOD` environment variable.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"device_code",
						"client_secret",
						"client_certificate",
						"interactive_browser",
						"username_password",
						"client_assertion",
					),
				},
			},
			"tenant_id": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Description: "The M365 tenant ID for the Entra ID application. " +
					"This ID uniquely identifies your Entra ID (EID) instance. " +
					"It can be found in the Azure portal under Entra ID > Properties. " +
					"Can also be set using the `M365_TENANT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID("tenant_id"),
				},
			},
			"client_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client ID for the Entra ID application. " +
					"This ID is generated when you register an application in the Entra ID (Azure AD) " +
					"and can be found under App registrations > YourApp > Overview. " +
					"Can also be set using the `M365_CLIENT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID("client_id"),
				},
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client secret for the Entra ID application. " +
					"This secret is generated in the Entra ID (Azure AD) and is required for " +
					"authentication flows such as client credentials and on-behalf-of flows. " +
					"It can be found under App registrations > YourApp > Certificates & secrets. " +
					"Required for client credentials and on-behalf-of flows. " +
					"Can also be set using the `M365_CLIENT_SECRET` environment variable.",
			},
			"client_certificate": schema.StringAttribute{
				Description: "The path to the Client Certificate file associated with the Service " +
					"Principal for use when authenticating as a Service Principal using a Client Certificate. " +
					"Supports PKCS#12 (.pfx or .p12) file format. The file should contain the certificate, " +
					"private key, and optionally a certificate chain. " +
					"IMPORTANT: As a prerequisite, the public key certificate must be uploaded to the " +
					"Enterprise Application in Microsoft Entra ID (formerly Azure Active Directory). This can be done in the Azure Portal " +
					"under 'Enterprise Applications' > [Your App] > 'Certificates & secrets' > 'Certificates'. " +
					"Use 'client_certificate_password' if the file is encrypted. This certificate should be " +
					"associated with the application registered in Azure Entra ID. " +
					"Can also be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"client_certificate_password": schema.StringAttribute{
				Description: "The password to decrypt the PKCS#12 (.pfx or .p12) file specified in " +
					"'client_certificate_file_path'. This is required if the PKCS#12 file is password-protected. " +
					"When the certificate file is created, this password is used to encrypt the private key for " +
					"security. It's not related to any Microsoft Entra ID (formerly Azure Active Directory) settings," +
					"but rather to the certificate file itself. " +
					"If your PKCS#12 file was created without a password, this field should be left empty. " +
					"Can also be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
				Description: "The username for username/password authentication. Can also be set using the" +
					"`M365_USERNAME` environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The password for username/password authentication. Can also be set using the" +
					"`M365_PASSWORD` environment variable.",
			},
			"client_assertion": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client assertion string (OIDC token) for authentication. " +
					"This is typically a JSON Web Token (JWT) that represents the identity of the client. " +
					"It is used in the client credentials flow with client assertion. " +
					"This method is more secure than client secret as the assertion is short-lived. " +
					"Commonly used in CI/CD pipelines and server-to-server authentication scenarios. " +
					"Can also be set using the `M365_CLIENT_ASSERTION` environment variable. " +
					"If both this and `client_assertion_file` are specified, this takes precedence.",
			},
			"client_assertion_file": schema.StringAttribute{
				Optional: true,
				Description: "Path to a file containing the client assertion (OIDC token) for authentication. " +
					"This file should contain a JSON Web Token (JWT) that represents the identity of the client. " +
					"Useful when the assertion is too long to be specified directly or when it's generated externally. " +
					"The provider will read this file to obtain the assertion string. " +
					"Ensure the file permissions are set appropriately to protect the token. " +
					"Can also be set using the `M365_CLIENT_ASSERTION_FILE` environment variable. " +
					"If both this and `client_assertion` are specified, `client_assertion` takes precedence.",
			},
			"redirect_url": schema.StringAttribute{
				Optional: true,
				Description: "The redirect URL for interactive browser authentication. Can also be set using " +
					"the `M365_REDIRECT_URL` environment variable.",
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
			}, "telemetry_optout": schema.BoolAttribute{
				Optional: true,
				Description: "Flag to indicate whether to opt out of telemetry. Default is `false`. " +
					"Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable.",
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
		ClientAssertion:           types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_ASSERTION", config.ClientAssertion.ValueString())),
		ClientAssertionFile:       types.StringValue(helpers.EnvDefaultFunc("M365_CLIENT_ASSERTION_FILE", config.ClientAssertionFile.ValueString())),
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

	ctx = tflog.SetField(ctx, "client_assertion", data.ClientAssertion.ValueString())
	ctx = tflog.SetField(ctx, "client_assertion_file", data.ClientAssertionFile.ValueString())
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`(?i)client_assertion`))

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
	})
}

// New returns a new provider.Provider instance for the Microsoft365 provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		p := &M365Provider{
			version: version,
			clients: &client.GraphClients{},
		}
		return p
	}
}
