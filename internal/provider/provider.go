package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
	clients *GraphClients
}

type GraphClients struct {
	StableClient *msgraphsdk.GraphServiceClient
	BetaClient   *msgraphbetasdk.GraphServiceClient
}

// M365ProviderModel describes the provider data model.
type M365ProviderModel struct {
	TenantID                  types.String `tfsdk:"tenant_id"`
	AuthMethod                types.String `tfsdk:"auth_method"`
	ClientID                  types.String `tfsdk:"client_id"`
	ClientSecret              types.String `tfsdk:"client_secret"`
	ClientCertificateBase64   types.String `tfsdk:"client_certificate_base64"`
	ClientCertificateFilePath types.String `tfsdk:"client_certificate_file_path"`
	ClientCertificatePassword types.String `tfsdk:"client_certificate_password"`
	Username                  types.String `tfsdk:"username"`
	Password                  types.String `tfsdk:"password"`
	RedirectURL               types.String `tfsdk:"redirect_url"`
	UseProxy                  types.Bool   `tfsdk:"use_proxy"`
	ProxyURL                  types.String `tfsdk:"proxy_url"`
	Cloud                     types.String `tfsdk:"cloud"`
	EnableChaos               types.Bool   `tfsdk:"enable_chaos"`
	TelemetryOptout           types.Bool   `tfsdk:"telemetry_optout"`
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
				MarkdownDescription: "The cloud to use for authentication and Graph / Graph Beta API requests." +
					"Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`." +
					"Can also be set using the `M365_CLOUD` environment variable.",
				Required: true,
				// Validators: []validator.String{
				// 	validateCloud(),
				// },
			},
			"auth_method": schema.StringAttribute{
				Required: true,
				Description: "The authentication method to use for the Entra ID application to authenticate the provider. " +
					"Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', " +
					"'username_password'. Can also be set using the `M365_AUTH_METHOD` environment variable.",
				// Validators: []validator.String{
				// 	validateAuthMethod(),
				// },
			},
			"tenant_id": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Description: "The M365 tenant ID for the Entra ID application. " +
					"This ID uniquely identifies your Entra ID (EID) instance. " +
					"It can be found in the Azure portal under Entra ID > Properties. " +
					"Can also be set using the `M365_TENANT_ID` environment variable.",
				// Validators: []validator.String{
				// 	validateGUID(),
				// },
			},
			"client_id": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client ID for the Entra ID application. " +
					"This ID is generated when you register an application in the Entra ID (Azure AD) " +
					"and can be found under App registrations > YourApp > Overview. " +
					"Can also be set using the `M365_CLIENT_ID` environment variable.",
				// Validators: []validator.String{
				// 	validateGUID(),
				// },
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
			"client_certificate_base64": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded PKCS#12 certificate bundle. For use when" +
					"authenticating as a Service Principal using a Client Certificate. Can also be" +
					"set using the `M365_CLIENT_CERTIFICATE_BASE64` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"client_certificate_file_path": schema.StringAttribute{
				MarkdownDescription: "The path to the Client Certificate associated with the Service" +
					"Principal for use when authenticating as a Service Principal using a Client Certificate." +
					"Can also be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"client_certificate_password": schema.StringAttribute{
				MarkdownDescription: "The password associated with the Client Certificate. For use when" +
					"authenticating as a Service Principal using a Client Certificate. Can also be set using" +
					"the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.",
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
			"redirect_url": schema.StringAttribute{
				Optional: true,
				Description: "The redirect URL for interactive browser authentication. Can also be set using" +
					"the `M365_REDIRECT_URL` environment variable.",
				Validators: []validator.String{
					validateURL(),
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
					validateURL(),
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
				Description: "Flag to indicate whether to opt out of telemetry. Default is `false`. " +
					"Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable.",
				MarkdownDescription: "Flag to indicate whether to opt out of telemetry. Default is `false`. " +
					"Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable.",
				Optional: true,
			},
		},
	}
}

// Configure configures the M365Provider with the given settings. It reads
// the configuration data from the provided request, applies defaults and
// environment variable overrides as necessary, and sets up authentication
// and client options based on the configuration. If any required configuration
// is missing or invalid, it appends appropriate diagnostics to the response.
//
// The function supports various authentication methods including device code,
// client secret, client certificate, on-behalf-of, interactive browser, and
// username/password. It also handles optional proxy settings and national cloud
// deployments.
//
// Parameters:
//   - ctx: The context for the configure request.
//   - req: The configure request containing the provider configuration.
//   - resp: The configure response used to store any diagnostics and the
//     configured client.
//
// The function performs the following steps:
//  1. Extracts configuration data from the request.
//  2. Retrieves values from environment variables if not set in the configuration.
//  3. Handles token retrieval from configuration or environment.
//  4. Configures HTTP client transport for proxy if specified.
//  5. Sets up authentication using the specified method.
//  6. Creates a Microsoft Graph client with the configured authentication provider.
//
// If any errors occur during these steps, appropriate diagnostics are added
// to the response.
func (p *M365Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data M365ProviderModel

	tflog.Debug(ctx, "Configure request received")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cloud := helpers.GetEnvOrDefault(ctx, data.Cloud.ValueString(), "M365_CLOUD")
	authMethod := helpers.GetEnvOrDefault(ctx, data.AuthMethod.ValueString(), "M365_AUTH_METHOD")
	tenantID := helpers.GetEnvOrDefault(ctx, data.TenantID.ValueString(), "M365_TENANT_ID")
	clientID := helpers.GetEnvOrDefault(ctx, data.ClientID.ValueString(), "M365_CLIENT_ID")
	clientSecret := helpers.GetEnvOrDefault(ctx, data.ClientSecret.ValueString(), "M365_CLIENT_SECRET")
	clientCertificateBase64 := helpers.GetEnvOrDefault(ctx, data.ClientCertificateBase64.ValueString(), "M365_CLIENT_CERTIFICATE_BASE64")
	clientCertificateFilePath := helpers.GetEnvOrDefault(ctx, data.ClientCertificateFilePath.ValueString(), "M365_CLIENT_CERTIFICATE_FILE_PATH")
	clientCertificatePassword := helpers.GetEnvOrDefault(ctx, data.ClientCertificatePassword.ValueString(), "M365_CLIENT_CERTIFICATE_PASSWORD")
	username := helpers.GetEnvOrDefault(ctx, data.Username.ValueString(), "M365_USERNAME")
	password := helpers.GetEnvOrDefault(ctx, data.Password.ValueString(), "M365_PASSWORD")
	redirectURL := helpers.GetEnvOrDefault(ctx, data.RedirectURL.ValueString(), "M365_REDIRECT_URL")
	useProxy := helpers.GetEnvOrDefaultBool(ctx, data.UseProxy.ValueBool(), "M365_USE_PROXY")
	proxyURL := helpers.GetEnvOrDefault(ctx, data.ProxyURL.ValueString(), "M365_PROXY_URL")
	enableChaos := helpers.GetEnvOrDefaultBool(ctx, data.EnableChaos.ValueBool(), "M365_ENABLE_CHAOS")
	telemetryOptout := helpers.GetEnvOrDefaultBool(ctx, data.TelemetryOptout.ValueBool(), "M365_TELEMETRY_OPTOUT")

	// Logging to verify environment variables are being used
	tflog.Debug(ctx, "Configuration values", map[string]interface{}{
		"cloud":                     cloud,
		"authMethod":                authMethod,
		"tenantID":                  tenantID,
		"clientID":                  clientID,
		"clientSecret":              clientSecret,
		"clientCertificateBase64":   clientCertificateBase64,
		"clientCertificateFilePath": clientCertificateFilePath,
		"clientCertificatePassword": clientCertificatePassword,
		"username":                  username,
		"password":                  password,
		"redirectURL":               redirectURL,
		"useProxy":                  useProxy,
		"proxyURL":                  proxyURL,
		"enableChaos":               enableChaos,
		"telemetryOptout":           telemetryOptout,
	})

	ctx = tflog.SetField(ctx, "cloud", cloud)
	ctx = tflog.SetField(ctx, "auth_method", authMethod)
	ctx = tflog.SetField(ctx, "use_proxy", useProxy)
	ctx = tflog.SetField(ctx, "redirect_url", redirectURL)
	ctx = tflog.SetField(ctx, "proxy_url", proxyURL)
	ctx = tflog.SetField(ctx, "enable_chaos", enableChaos)

	ctx = tflog.SetField(ctx, "client_certificate_base64", clientCertificateBase64)
	ctx = tflog.SetField(ctx, "client_certificate_file_path", clientCertificateFilePath)
	ctx = tflog.SetField(ctx, "client_certificate_password", clientCertificatePassword)
	ctx = tflog.MaskAllFieldValuesRegexes(ctx, regexp.MustCompile(`(?i)client_certificate_base64`))

	ctx = tflog.SetField(ctx, "username", username)
	ctx = tflog.SetField(ctx, "password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "password")

	ctx = tflog.SetField(ctx, "tenant_id", tenantID)
	ctx = tflog.SetField(ctx, "client_id", clientID)
	ctx = tflog.SetField(ctx, "client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "tenant_id")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_id")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_secret")

	authorityURL, apiScope, graphServiceRoot, graphBetaServiceRoot, err := setCloudConstants(cloud)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Microsoft Cloud Type",
			fmt.Sprintf("An error occurred while attempting to get cloud constants for cloud type '%s'. "+
				"Please ensure the cloud type is valid. Detailed error: %s", cloud, err.Error()),
		)
		return
	}

	ctx = tflog.SetField(ctx, "authority_url", authorityURL)
	ctx = tflog.SetField(ctx, "api_scope", apiScope)
	ctx = tflog.SetField(ctx, "graph_service_root", graphServiceRoot)
	ctx = tflog.SetField(ctx, "graph_beta_service_root", graphBetaServiceRoot)

	clientOptions, err := configureEntraIDClientOptions(ctx, useProxy, proxyURL, authorityURL, telemetryOptout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure client options",
			fmt.Sprintf("An error occurred while attempting to configure client options. Detailed error: %s", err.Error()),
		)
		return
	}

	cred, err := obtainCredential(ctx, data, clientOptions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create credentials",
			fmt.Sprintf("An error occurred while attempting to create the credentials using the provided authentication method '%s'. "+
				"This may be due to incorrect or missing credentials, misconfigured client options, or issues with the underlying authentication library. "+
				"Please verify the authentication method and credentials configuration. Detailed error: %s", data.AuthMethod.ValueString(), err.Error()),
		)
		return
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{apiScope})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create authentication provider",
			fmt.Sprintf("An error occurred while attempting to create the authentication provider using the provided credentials. "+
				"This may be due to misconfigured client options, incorrect credentials, or issues with the underlying authentication library. "+
				"Please verify your client options and credentials configuration. Detailed error: %s", err.Error()),
		)
		return
	}

	httpClient, err := configureGraphClientOptions(ctx, useProxy, proxyURL, enableChaos)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure Graph client options",
			fmt.Sprintf("An error occurred while attempting to configure the Microsoft Graph client options. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider, nil, nil, httpClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Stable SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Stable SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	betaAdapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider, nil, nil, httpClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Beta SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Beta SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter.SetBaseUrl(graphServiceRoot)
	betaAdapter.SetBaseUrl(graphBetaServiceRoot)

	clients := &GraphClients{
		StableClient: msgraphsdk.NewGraphServiceClient(stableAdapter),
		BetaClient:   msgraphbetasdk.NewGraphServiceClient(betaAdapter),
	}

	p.clients = clients

	resp.DataSourceData = clients
	resp.ResourceData = clients
}

// New returns a new provider.Provider instance for the Microsoft365 provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &M365Provider{
			version: version,
			clients: &GraphClients{},
		}
	}
}
