package graphBetaWindowsUpdateCatalogItem

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
	DataSourceName = "microsoft365_graph_beta_device_management_windows_update_catalog_item"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsUpdateCatalogItemDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsUpdateCatalogItemDataSource{}
)

// NewWindowsUpdateCatalogItemDataSource creates a new data source for Windows Update Catalog Items
func NewWindowsUpdateCatalogItemDataSource() datasource.DataSource {
	return &WindowsUpdateCatalogItemDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsUpdateCatalogItemDataSource defines the data source implementation
type WindowsUpdateCatalogItemDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

// Metadata returns the data source type name
func (d *WindowsUpdateCatalogItemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsUpdateCatalogItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsUpdateCatalogItemDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Update Catalog Items from Microsoft Intune with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `release_date_time`, `end_of_support_date`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "release_date_time", "end_of_support_date"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For date filters, use RFC3339 format (e.g., '2023-01-01T00:00:00Z').",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Windows Update Catalog Items that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the catalog item.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the catalog item.",
						},
						"release_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release date time of the catalog item.",
						},
						"end_of_support_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The end of support date of the catalog item.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
