package graphBetaWindowsQualityUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_management_windows_quality_update_policy"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsQualityUpdateProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsQualityUpdateProfileDataSource{}
)

// NewWindowsQualityUpdateProfileDataSource creates a new data source for Windows Driver Update Profile
func NewWindowsQualityUpdateProfileDataSource() datasource.DataSource {
	return &WindowsQualityUpdateProfileDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsQualityUpdateProfileDataSource defines the data source implementation
type WindowsQualityUpdateProfileDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

// Metadata returns the data source type name
func (d *WindowsQualityUpdateProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsQualityUpdateProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsQualityUpdateProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a Windows Driver Update Profile in Microsoft Intune.",
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
				MarkdownDescription: "List of Scope Tags for this Driver Update entity.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
