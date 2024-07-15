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
	"github.com/hashicorp/terraform-plugin-framework/types"
	khttp "github.com/microsoft/kiota-http-go"
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
	TenantID        types.String `tfsdk:"tenant_id"`
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	CertificatePath types.String `tfsdk:"certificate_path"`
	UserAssertion   types.String `tfsdk:"user_assertion"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	RedirectURL     types.String `tfsdk:"redirect_url"`
	Token           types.String `tfsdk:"token"`
	UseBeta         types.Bool   `tfsdk:"use_beta"`
	UseProxy        types.Bool   `tfsdk:"use_proxy"`
	ProxyURL        types.String `tfsdk:"proxy_url"`
	EnableChaos     types.Bool   `tfsdk:"enable_chaos"`
	AuthMethod      types.String `tfsdk:"auth_method"`
	TokenEndpoint   types.String `tfsdk:"token_endpoint"`
	ServiceRoot     types.String `tfsdk:"service_root"`
}

func (p *M365Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "M365"
	resp.Version = p.version
}

func (p *M365Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:    true,
				Description: "The tenant ID for the Azure AD application.",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The client ID for the Azure AD application.",
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The client secret for the Azure AD application. " +
					"Required for client credentials and on-behalf-of flows.",
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
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Description: "The token for the Azure AD application. " +
					"Can also be set using the `M365_API_TOKEN` environment variable.",
			},
			"use_beta": schema.BoolAttribute{
				Optional:    true,
				Description: "Use the beta version of the Microsoft Graph API.",
			},
			"use_proxy": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable the use of an HTTP proxy.",
			},
			"proxy_url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL of the HTTP proxy.",
			},
			"enable_chaos": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable chaos handler for simulating specific scenarios.",
			},
			"auth_method": schema.StringAttribute{
				Optional: true,
				Description: "The authentication method to use. " +
					"Options: 'device_code', 'client_secret', 'client_certificate', 'on_behalf_of', " +
					"'interactive_browser', 'username_password'.",
			},
			"token_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The token endpoint for the national cloud deployment.",
			},
			"service_root": schema.StringAttribute{
				Optional:    true,
				Description: "The Microsoft Graph service root endpoint for the national cloud deployment.",
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

	tenantID := data.TenantID.ValueString()
	clientID := data.ClientID.ValueString()
	useBeta := data.UseBeta.ValueBool()
	useProxy := data.UseProxy.ValueBool()
	proxyURL := data.ProxyURL.ValueString()
	enableChaos := data.EnableChaos.ValueBool()
	authMethod := data.AuthMethod.ValueString()
	tokenEndpoint := data.TokenEndpoint.ValueString()
	serviceRoot := data.ServiceRoot.ValueString()

	var cred azcore.TokenCredential
	var err error

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"M365 provider configuration error",
			"Cannot use unknown value as token",
		)
		return
	}

	if data.Token.IsNull() {
		token := os.Getenv("M365_API_TOKEN")
		if token == "" {
			resp.Diagnostics.AddError(
				"M365 provider configuration error",
				"Token cannot be an empty string",
			)
			return
		}
	}

	var transport *http.Transport
	if useProxy {
		proxyUrlParsed, err := url.Parse(proxyURL)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid proxy URL",
				"Error parsing proxy URL: "+err.Error(),
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

	if tokenEndpoint != "" {
		clientOptions.Cloud.ActiveDirectoryAuthorityHost = tokenEndpoint
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
		clientSecret := data.ClientSecret.ValueString()
		cred, err = azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "client_certificate":
		certificatePath := data.CertificatePath.ValueString()
		certFile, err := os.Open(certificatePath)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to open certificate file",
				"Error opening certificate file: "+err.Error(),
			)
			return
		}
		defer certFile.Close()

		info, err := certFile.Stat()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to stat certificate file",
				"Error stating certificate file: "+err.Error(),
			)
			return
		}

		certBytes := make([]byte, info.Size())
		_, err = certFile.Read(certBytes)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to read certificate file",
				"Error reading certificate file: "+err.Error(),
			)
			return
		}

		certs, key, err := azidentity.ParseCertificates(certBytes, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse certificates",
				"Error parsing certificates: "+err.Error(),
			)
			return
		}

		cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, certs, key, &azidentity.ClientCertificateCredentialOptions{
			ClientOptions: clientOptions,
		})
	case "on_behalf_of":
		clientSecret := data.ClientSecret.ValueString()
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
			"Error creating credentials: "+err.Error(),
		)
		return
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create authentication provider",
			"Error creating authentication provider: "+err.Error(),
		)
		return
	}

	clientOptionsGraph := msgraphgocore.GraphClientOptions{}
	middleware := msgraphgocore.GetDefaultMiddlewaresWithOptions(&clientOptionsGraph)

	if enableChaos {
		chaosHandler := khttp.NewChaosHandler()
		middleware = append(middleware, chaosHandler)
	}

	httpClient := khttp.GetDefaultClient(middleware...)
	if useProxy {
		httpClient, err = khttp.GetClientWithProxySettings(proxyURL, middleware...)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create HTTP client with proxy settings",
				"Error creating HTTP client with proxy settings: "+err.Error(),
			)
			return
		}
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider, nil, nil, httpClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create adapter",
			"Error creating adapter: "+err.Error(),
		)
		return
	}

	if serviceRoot != "" {
		adapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", serviceRoot))
	} else if useBeta {
		adapter.SetBaseUrl("https://graph.microsoft.com/beta")
	} else {
		adapter.SetBaseUrl("https://graph.microsoft.com/v1.0")
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

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
