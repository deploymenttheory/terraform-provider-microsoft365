package graphBetaMobileAppCatalogPackage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	datasourceName = "graph_beta_device_and_app_management_mobile_app_catalog_package"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &MobileAppCatalogPackageDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MobileAppCatalogPackageDataSource{}
)

func NewMobileAppCatalogPackageDataSource() datasource.DataSource {
	return &MobileAppCatalogPackageDataSource{
		ReadPermissions: []string{
			"DeviceAppManagement.Read.All",
		},
	}
}

type MobileAppCatalogPackageDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the datasource type name.
func (r *MobileAppCatalogPackageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceName
}

// Configure sets the client for the data source
func (d *MobileAppCatalogPackageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *MobileAppCatalogPackageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves mobile app catalog packages from Microsoft Intune using the `/deviceAppManagement/MobileAppCatalogPackage` endpoint. This data source enables querying mobile app catalog packages with advanced filtering capabilities including OData queries for filtering by product name, publisher, and other properties.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `product_name`, `publisher_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "product_name", "publisher_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all' or 'odata'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: productDisplayName eq 'Microsoft Office'.",
			},
			"odata_top": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.",
			},
			"odata_skip": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $skip parameter for pagination. Only used when filter_type is 'odata'.",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.",
			},
			"odata_orderby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: productDisplayName.",
			},
			"odata_count": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "OData $count parameter to include count of total results. Only used when filter_type is 'odata'.",
			},
			"odata_search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $search parameter for full-text search. Only used when filter_type is 'odata'.",
			},
			"odata_expand": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $expand parameter to include related entities. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of mobile app catalog packages that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the mobile app catalog package.",
						},
						"product_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the product.",
						},
						"product_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the product.",
						},
						"publisher_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the publisher.",
						},
						"version_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the version.",
						},
						"branch_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the branch.",
						},
						"applicable_architectures": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The applicable architectures for the package (e.g., x64, x86, ARM64).",
						},
						"locales": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "The list of supported locales for the package.",
						},
						"package_auto_update_capable": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the package supports automatic updates.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
