package provider

import (
	"context"
	"fmt"
	"os"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &M365Provider{}

// M365Provider defines the provider implementation.
type M365Provider struct {
	version string
}

// M365ProviderModel describes the provider data model.
type M365ProviderModel struct {
	TenantID types.String `tfsdk:"tenant_id"`
	ClientID types.String `tfsdk:"client_id"`
	Token    types.String `tfsdk:"token"`
	UseBeta  types.Bool   `tfsdk:"use_beta"`
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

	var token string
	if data.Token.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"M365 provider configuration error",
			"Cannot use unknown value as token",
		)
		return
	}

	if data.Token.IsNull() {
		token = os.Getenv("M365_API_TOKEN")
	} else {
		token = data.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"M365 provider configuration error",
			"Token cannot be an empty string",
		)
		return
	}

	cred, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID: tenantID,
		ClientID: clientID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create credentials",
			"Error creating credentials: "+err.Error(),
		)
		return
	}

	var client interface{}
	if useBeta {
		client, err = msgraphbetasdk.NewGraphServiceClientWithCredentials(cred, []string{"Files.Read"})
	} else {
		client, err = msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"Files.Read"})
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Error creating client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *M365Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJMESPathCheckResource,
	}
}

func (p *M365Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewJMESPathCheckDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &M365Provider{
			version: version,
		}
	}
}
