package graphBetaReuseablePolicySettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_device_management_reuseable_policy_setting"
	ReadTimeout  = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &ReuseablePolicySettingsDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &ReuseablePolicySettingsDataSource{}
)

func NewReuseablePolicySettingsDataSource() datasource.DataSource {
	return &ReuseablePolicySettingsDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type ReuseablePolicySettingsDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the resource type name.
func (r *ReuseablePolicySettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the data source
func (d *ReuseablePolicySettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *ReuseablePolicySettingsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Reusable Policy Settings from Microsoft Intune with explicit filtering options. " +
			"Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificate " +
			"directly to an elevation rule. Like all reusable settings groups for Intune, configurations and changes made to a reusable settings " +
			"group are automatically passed to the policies that reference the group.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Reusable Policy Settings that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the reusable policy setting.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the reusable policy setting.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the reusable policy setting.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
