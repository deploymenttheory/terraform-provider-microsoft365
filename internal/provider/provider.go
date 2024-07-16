package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &M365Provider{}

// M365Provider defines the provider implementation.
type M365Provider struct {
	version string
}

// M365ProviderModel describes the provider data model.
type M365ProviderModel struct {
	TenantID                             types.String `tfsdk:"tenant_id"`
	AuthMethod                           types.String `tfsdk:"auth_method"`
	ClientID                             types.String `tfsdk:"client_id"`
	ClientSecret                         types.String `tfsdk:"client_secret"`
	CertificatePath                      types.String `tfsdk:"certificate_path"`
	UserAssertion                        types.String `tfsdk:"user_assertion"`
	Username                             types.String `tfsdk:"username"`
	Password                             types.String `tfsdk:"password"`
	RedirectURL                          types.String `tfsdk:"redirect_url"`
	Token                                types.String `tfsdk:"token"`
	UseGraphBeta                         types.Bool   `tfsdk:"use_graph_beta"`
	UseProxy                             types.Bool   `tfsdk:"use_proxy"`
	ProxyURL                             types.String `tfsdk:"proxy_url"`
	EnableChaos                          types.Bool   `tfsdk:"enable_chaos"`
	NationalCloudDeployment              types.Bool   `tfsdk:"national_cloud_deployment"`
	NationalCloudDeploymentTokenEndpoint types.String `tfsdk:"national_cloud_deployment_token_endpoint"`
	NationalCloudDeploymentServiceRoot   types.String `tfsdk:"national_cloud_deployment_service_root"`
}

func (p *M365Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "M365"
	resp.Version = p.version
}

func (p *M365Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required: true,
				Description: "The M365 tenant ID for the Entra ID application. " +
					"This ID uniquely identifies your Entra ID (EID) instance. " +
					"It can be found in the Azure portal under Entra ID > Properties. " +
					"Can also be set using the `M365_TENANT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID(),
				},
			},
			"auth_method": schema.StringAttribute{
				Optional: true,
				Description: "The authentication method to use. " +
					"Options: 'device_code', 'client_secret', 'client_certificate', 'on_behalf_of', " +
					"'interactive_browser', 'username_password'.",
				Validators: []validator.String{
					validateAuthMethod(),
				},
			},
			"client_id": schema.StringAttribute{
				Optional: true,
				Description: "The client ID for the Entra ID application. " +
					"This ID is generated when you register an application in the Entra ID (Azure AD) " +
					"and can be found under App registrations > YourApp > Overview. " +
					"Can also be set using the `M365_CLIENT_ID` environment variable.",
				Validators: []validator.String{
					validateGUID(),
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
			"certificate_path": schema.StringAttribute{
				Optional: true,
				Description: "Path to the client certificate file. " +
					"Required for client certificate authentication.",
			},
			"user_assertion": schema.StringAttribute{
				Optional:    true,
				Description: "The user assertion for on-behalf-of authentication.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "The username for username/password authentication.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The password for username/password authentication.",
			},
			"redirect_url": schema.StringAttribute{
				Optional:    true,
				Description: "The redirect URL for interactive browser authentication.",
				Validators: []validator.String{
					validateURL(),
				},
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The token for the Azure AD application. " +
					"Can also be set using the `M365_API_TOKEN` environment variable.",
			},
			"use_graph_beta": schema.BoolAttribute{
				Optional: true,
				Description: "When set to true, the provider will use the beta version of the Microsoft Graph API " +
					"(https://graph.microsoft.com/beta). When set to false or not set, the provider will use the stable " +
					"version of the Microsoft Graph API (https://graph.microsoft.com/v1.0).",
			},
			"use_proxy": schema.BoolAttribute{
				Optional: true,
				Description: "Enables the use of an HTTP proxy for network requests. When set to true, the provider will " +
					"route all HTTP requests through the specified proxy server. This can be useful for environments that " +
					"require proxy access for internet connectivity or for monitoring and logging HTTP traffic.",
				Validators: []validator.Bool{
					validateUseProxy(),
				},
			},
			"proxy_url": schema.StringAttribute{
				Optional: true,
				Description: "Specifies the URL of the HTTP proxy server. This URL should be in a valid URL format " +
					"(e.g., 'http://proxy.example.com:8080'). When 'use_proxy' is enabled, this URL is used to configure the " +
					"HTTP client to route requests through the proxy. Ensure the proxy server is reachable and correctly " +
					"configured to handle the network traffic.",
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
					"it can recover gracefully. Use with caution in production environments.",
			},
			"national_cloud_deployment": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true if connecting to Microsoft Graph national cloud deployments. (Microsoft Cloud for US Government and Microsoft Azure and Microsoft 365 operated by 21Vianet in China.)",
			},
			"national_cloud_deployment_token_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "By default, the provider is configured to access data in the Microsoft Graph global service, using the https://graph.microsoft.com root URL to access the Microsoft Graph REST API. This field overrides this configuration to connect to Microsoft Graph national cloud deployments. Microsoft Cloud for US Government and Microsoft Azure and Microsoft 365 operated by 21Vianet in China. https://learn.microsoft.com/en-gb/graph/deployments",
				Validators: []validator.String{
					validateURL(),
					validateNationalCloudDeployment(),
				},
			},
			"national_cloud_deployment_service_root": schema.StringAttribute{
				Optional:    true,
				Description: "The Microsoft Graph service root endpoint for the national cloud deployment. Overrides the default Microsoft Graph service root endpoint (https://graph.microsoft.com/v1.0 / https://graph.microsoft.com/beta).This field overrides this configuration to connect to Microsoft Graph national cloud deployments. Microsoft Cloud for US Government and Microsoft Azure and Microsoft 365 operated by 21Vianet in China. https://learn.microsoft.com/en-gb/graph/deployments",
				Validators: []validator.String{
					validateURL(),
					validateNationalCloudDeployment(),
				},
			},
		},
	}
}

func (p *M365Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data M365ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tenantID := getEnvOrDefault(data.TenantID.ValueString(), "M365_TENANT_ID")
	authMethod := data.AuthMethod.ValueString()
	clientID := getEnvOrDefault(data.ClientID.ValueString(), "M365_CLIENT_ID")
	clientSecret := getEnvOrDefault(data.ClientSecret.ValueString(), "M365_CLIENT_SECRET")
	useGraphBeta := data.UseGraphBeta.ValueBool()
	useProxy := data.UseProxy.ValueBool()
	proxyURL := data.ProxyURL.ValueString()
	enableChaos := data.EnableChaos.ValueBool()
	nationalCloudDeployment := data.NationalCloudDeployment.ValueBool()
	nationalCloudDeploymentTokenEndpoint := data.NationalCloudDeploymentTokenEndpoint.ValueString()
	nationalCloudDeploymentServiceRoot := data.NationalCloudDeploymentServiceRoot.ValueString()

	var cred azcore.TokenCredential
	var err error

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"M365 Provider Configuration Warning",
			"The token value is unknown in the provider configuration. "+
				"The token will be obtained from the credentials provided using the associated MS Graph authentication provider.",
		)
	}

	if data.Token.IsNull() {
		token := os.Getenv("M365_API_TOKEN")
		if token == "" {
			resp.Diagnostics.AddError(
				"M365 Provider Configuration Error",
				"The token is not set in the provider configuration and the environment variable 'M365_API_TOKEN' is empty. "+
					"A token is required for authentication. Please provide a valid token either through the provider configuration or by setting the 'M365_API_TOKEN' environment variable. "+
					"Alternatively, ensure that the credentials provided can obtain a token dynamically.",
			)
			return
		}
	}

	var transport *http.Transport
	if useProxy {
		proxyUrlParsed, err := url.Parse(proxyURL)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Proxy URL",
				fmt.Sprintf("Failed to parse the provided proxy URL '%s': %s. "+
					"Ensure the URL is correctly formatted.", proxyURL, err.Error()),
			)
			return
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrlParsed),
		}
	} else {
		transport = &http.Transport{}
	}

	authClient := &http.Client{
		Transport: transport,
	}

	clientOptions := policy.ClientOptions{
		Transport: authClient,
	}

	// Set cloud configuration for national cloud deployments
	if nationalCloudDeployment && nationalCloudDeploymentTokenEndpoint != "" {
		clientOptions.Cloud.ActiveDirectoryAuthorityHost = nationalCloudDeploymentTokenEndpoint
	}

	switch authMethod {
	case "device_code":
		cred, err = azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID: tenantID,
			ClientID: clientID,
			UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
				fmt.Println(message.Message)
				return nil
			},
			ClientOptions: clientOptions,
		})
	case "client_secret":
		cred, err = azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "client_certificate":
		certificatePath := data.CertificatePath.ValueString()
		certFile, err := os.Open(certificatePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Opening Certificate File",
				fmt.Sprintf("Failed to open the certificate file at path '%s': %s. "+
					"Ensure the file path is correct and the file is accessible.", certificatePath, err.Error()),
			)
			return
		}
		defer certFile.Close()

		info, err := certFile.Stat()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Accessing Certificate File",
				fmt.Sprintf("Failed to retrieve file information: %s. "+
					"Ensure the file exists and is accessible.", err.Error()),
			)
			return
		}

		certBytes := make([]byte, info.Size())
		_, err = certFile.Read(certBytes)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Certificate File",
				fmt.Sprintf("Failed to read the certificate file: %s. "+
					"Ensure the file is accessible and not corrupted.", err.Error()),
			)
			return
		}

		certs, key, err := azidentity.ParseCertificates(certBytes, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Parsing Certificates",
				fmt.Sprintf("Failed to parse certificates from the provided file: %s. "+
					"Ensure the file contains valid certificate data and is correctly formatted.", err.Error()),
			)
			return
		}

		cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, &azidentity.ClientCertificateCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "on_behalf_of":
		userAssertion := data.UserAssertion.ValueString()
		cred, err = azidentity.NewOnBehalfOfCredentialWithSecret(tenantID, clientID, userAssertion, clientSecret, &azidentity.OnBehalfOfCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "interactive_browser":
		redirectURL := data.RedirectURL.ValueString()
		cred, err = azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
			TenantID:      tenantID,
			ClientID:      clientID,
			RedirectURL:   redirectURL,
			ClientOptions: clientOptions,
		})
	case "username_password":
		username := data.Username.ValueString()
		password := data.Password.ValueString()
		cred, err = azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password, &azidentity.UsernamePasswordCredentialOptions{
			ClientOptions: clientOptions,
		})
	default:
		resp.Diagnostics.AddError(
			"Unsupported authentication method",
			fmt.Sprintf("The authentication method '%s' is not supported.", authMethod),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create credentials",
			fmt.Sprintf("An error occurred while attempting to create the credentials using the provided authentication method '%s'. "+
				"This may be due to incorrect or missing credentials, misconfigured client options, or issues with the underlying authentication library. "+
				"Please verify the authentication method and credentials configuration. Detailed error: %s", authMethod, err.Error()),
		)
		return
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create authentication provider",
			fmt.Sprintf("An error occurred while attempting to create the authentication provider using the provided credentials. "+
				"This may be due to misconfigured client options, incorrect credentials, or issues with the underlying authentication library. "+
				"Please verify your client options and credentials configuration. Detailed error: %s", err.Error()),
		)
		return
	}

	clientOptionsGraph := msgraphgocore.GraphClientOptions{}
	middleware := msgraphgocore.GetDefaultMiddlewaresWithOptions(&clientOptionsGraph)

	if enableChaos {
		chaosHandler := khttp.NewChaosHandler()
		middleware = append(middleware, chaosHandler)
	}

	var httpClient *http.Client
	if useProxy {
		httpClient, err = khttp.GetClientWithProxySettings(proxyURL, middleware...)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create HTTP client with proxy settings",
				fmt.Sprintf("An error occurred while attempting to create the HTTP client with the provided proxy settings. "+
					"This might be due to an invalid proxy URL, issues with the proxy server, or other network-related problems. "+
					"Please verify the proxy URL and your network connection. Detailed error: %s", err.Error()),
			)
			return
		}
	} else {
		httpClient = khttp.GetDefaultClient(middleware...)
	}

	var stableAdapter *msgraphsdk.GraphRequestAdapter
	var betaAdapter *msgraphbetasdk.GraphRequestAdapter

	if useGraphBeta {
		betaAdapter, err = msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
			authProvider, nil, nil, httpClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create Microsoft Graph Beta SDK Adapter",
				fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Beta SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
			)
			return
		}
	} else {
		stableAdapter, err = msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
			authProvider, nil, nil, httpClient)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create Microsoft Graph Stable SDK Adapter",
				fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Stable SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
			)
			return
		}
	}

	// Set the service root for national cloud deployments
	if nationalCloudDeployment && nationalCloudDeploymentServiceRoot != "" {
		if useGraphBeta {
			betaAdapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", nationalCloudDeploymentServiceRoot))
		} else {
			stableAdapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", nationalCloudDeploymentServiceRoot))
		}
	}

	var client interface{}

	if useGraphBeta {
		client = msgraphbetasdk.NewGraphServiceClient(betaAdapter)
	} else {
		client = msgraphsdk.NewGraphServiceClient(stableAdapter)
	}

	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add your resource functions here
	}
}

func (p *M365Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add your datasource functions here
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &M365Provider{
			version: version,
		}
	}
}
