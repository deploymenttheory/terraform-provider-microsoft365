package graphBetaBrowserSiteList

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
	DataSourceName = "microsoft365_graph_beta_m365_admin_browser_site_list"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &BrowserSiteListDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &BrowserSiteListDataSource{}
)

// NewBrowserSiteListDataSource creates a new data source for Browser Site Lists
func NewBrowserSiteListDataSource() datasource.DataSource {
	return &BrowserSiteListDataSource{
		ReadPermissions: []string{
			"BrowserSiteLists.Read.All",
		},
	}
}

// BrowserSiteListDataSource defines the data source implementation
type BrowserSiteListDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

// Metadata returns the data source type name
func (d *BrowserSiteListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *BrowserSiteListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *BrowserSiteListDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Browser Site Lists from Microsoft 365 Admin Centre with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `published_date_time`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For date filters, use RFC3339 format (e.g., '2023-01-01T00:00:00Z').",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Browser Site Lists that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the site list.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the site list.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
