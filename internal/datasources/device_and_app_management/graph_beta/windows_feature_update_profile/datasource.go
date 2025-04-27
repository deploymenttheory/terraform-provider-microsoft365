package graphBetaWindowsFeatureUpdateProfile

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
	DataSourceName = "graph_beta_windows_driver_update_profile"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsFeatureUpdateProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsFeatureUpdateProfileDataSource{}
)

// NewWindowsFeatureUpdateProfileDataSource creates a new data source for Windows Feature Update Profile
func NewWindowsFeatureUpdateProfileDataSource() datasource.DataSource {
	return &WindowsFeatureUpdateProfileDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsFeatureUpdateProfileDataSource defines the data source implementation
type WindowsFeatureUpdateProfileDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name
func (d *WindowsFeatureUpdateProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsFeatureUpdateProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.ProviderTypeName)
}

// Schema defines the schema for the data source
func (d *WindowsFeatureUpdateProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a Windows Feature Update Profile in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the profile.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name for the profile.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the profile which is specified by the user.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of Scope Tags for this Feature Update entity.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
