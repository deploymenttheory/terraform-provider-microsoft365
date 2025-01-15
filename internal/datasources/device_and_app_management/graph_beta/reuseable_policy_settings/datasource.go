package graphBetaReuseablePolicySettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_device_and_app_management_reuseable_policy_setting"
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

func (d *ReuseablePolicySettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Reuseable Settings Policy using Settings Catalog in Microsoft Intune for Endpoint Privilege Management." +
			"Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificate" +
			"directly to an elevation rule. Like all reusable settings groups for Intune, configurations and changes made to a reusable settings" +
			"group are automatically passed to the policies that reference the group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier for this Reuseable Settings Policy",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The reusable setting display name supplied by user.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reuseable Settings Policy description",
			},
			"settings": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reuseable Settings Policy with settings catalog settings defined as a valid JSON string.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation date and time of the settings catalog policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy",
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the policy",
			},
			"referencing_configuration_policies": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of configuration policies referencing this reuseable policy",
			},
			"referencing_configuration_policy_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of configuration policies referencing this reuseable policy",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *ReuseablePolicySettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
