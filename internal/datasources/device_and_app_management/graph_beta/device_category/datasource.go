package graphBetaDeviceCategory

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
	ResourceName = "graph_beta_device_and_app_management_device_category"
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &DeviceCategoryDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &DeviceCategoryDataSource{}
)

func NewDeviceCategoryDataSource() datasource.DataSource {
	return &DeviceCategoryDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type DeviceCategoryDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *DeviceCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (d *DeviceCategoryDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the device category.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name of the device category.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the device category.",
			},
			"role_scope_tags": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates role scope tags assigned for the device category.",
				ElementType:         types.StringType,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *DeviceCategoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
