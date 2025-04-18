package graphBetaRoleScopeTag

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
	ResourceName = "graph_beta_device_and_app_management_role_scope_tag"
	ReadTimeout  = 180
)

var (
	// Basic data source interface
	_ datasource.DataSource = &RoleScopeTagDataSource{}

	// Allows the data source to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &RoleScopeTagDataSource{}
)

func NewRoleScopeTagDataSource() datasource.DataSource {
	return &RoleScopeTagDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementRBAC.Read.All",
		},
	}
}

type RoleScopeTagDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name.
func (r *RoleScopeTagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (d *RoleScopeTagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a Role Scope Tag in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier for the Role Scope Tag.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display or friendly name of the Role Scope Tag.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description of the Role Scope Tag.",
			},
			"is_built_in": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates whether this is a built-in Role Scope Tag. This property is read-only.",
			},
			"assignments": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The list of group assignments for the Intune Role Scope Tag.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *RoleScopeTagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
