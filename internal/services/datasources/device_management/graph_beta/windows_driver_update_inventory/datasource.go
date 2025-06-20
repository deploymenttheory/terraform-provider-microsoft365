package graphBetaWindowsDriverUpdateInventory

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "graph_beta_windows_driver_update_inventory"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsDriverUpdateInventoryDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsDriverUpdateInventoryDataSource{}
)

// NewWindowsDriverUpdateInventoryDataSource creates a new data source for Windows Driver Update Inventory
func NewWindowsDriverUpdateInventoryDataSource() datasource.DataSource {
	return &WindowsDriverUpdateInventoryDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsDriverUpdateInventoryDataSource defines the data source implementation
type WindowsDriverUpdateInventoryDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name
func (d *WindowsDriverUpdateInventoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure sets the client for the data source
func (d *WindowsDriverUpdateInventoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *WindowsDriverUpdateInventoryDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a Windows Driver Update Inventory in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the driver inventory.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the driver.",
			},
			"windows_driver_update_profile_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Windows Driver Update Profile this inventory belongs to.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
