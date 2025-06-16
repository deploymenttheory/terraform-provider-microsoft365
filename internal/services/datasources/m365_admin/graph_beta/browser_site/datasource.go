package graphBetaBrowserSite

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName        = "graph_beta_m365_admin_browser_site"
	ReadDataSourceTimeout = 180
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &BrowserSiteDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &BrowserSiteDataSource{}
)

// NewBrowserSiteDataSource creates a new data source for Browser Sites
func NewBrowserSiteDataSource() datasource.DataSource {
	return &BrowserSiteDataSource{
		ReadPermissions: []string{
			"BrowserSiteLists.Read.All",
		},
	}
}

// BrowserSiteDataSource defines the data source implementation
type BrowserSiteDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name
func (d *BrowserSiteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure configures the data source with the provider client
func (d *BrowserSiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.ProviderTypeName)
}

// Schema defines the schema for the data source
func (d *BrowserSiteDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Browser Sites from Microsoft 365 Admin Centre with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `web_url`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "web_url"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"browser_site_list_assignment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The browser site list ID this browser site belongs to.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Browser Sites that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the browser site.",
						},
						"browser_site_list_assignment_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The browser site list id this browser site belongs to.",
						},
						"web_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL of the site.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
