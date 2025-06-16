package graphBetaApplicationCategory

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
	DataSourceName = "graph_beta_device_and_app_management_application_category"
	ReadTimeout    = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &ApplicationCategoryDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &ApplicationCategoryDataSource{}
)

func NewApplicationCategoryDataSource() datasource.DataSource {
	return &ApplicationCategoryDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementApps.Read.All",
		},
	}
}

type ApplicationCategoryDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *ApplicationCategoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Schema defines the schema for the data source
func (d *ApplicationCategoryDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Application Categories from Microsoft Intune with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `last_modified_date_time`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "last_modified_date_time"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For date filters, use RFC3339 format (e.g., '2023-01-01T00:00:00Z').",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Application Categories that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the application category.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the application category.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the application category was last modified.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *ApplicationCategoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
