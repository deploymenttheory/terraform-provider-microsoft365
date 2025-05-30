package graphBetaCloudPC

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	datasourceName = "graph_beta_device_and_app_management_cloud_pc"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &CloudPCDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &CloudPCDataSource{}
)

func NewCloudPCDataSource() datasource.DataSource {
	return &CloudPCDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPCDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the datasource type name.
func (r *CloudPCDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + datasourceName
}

// Schema defines the schema for the data source
func (d *CloudPCDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PCs from Microsoft Intune with explicit filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'.",
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
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of cloud PCs that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the cloud PC.",
						},
						"aad_device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Azure AD device ID of the cloud PC.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the cloud PC.",
						},
						"image_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the OS image that's on the cloud PC.",
						},
						"managed_device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Intune managed device ID of the cloud PC.",
						},
						"managed_device_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Intune managed device name of the cloud PC.",
						},
						"provisioning_policy_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The provisioning policy ID of the cloud PC.",
						},
						"provisioning_policy_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The provisioning policy name of the cloud PC.",
						},
						"on_premises_connection_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure network connection that is applied during the provisioning of cloud PCs.",
						},
						"service_plan_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan ID of the cloud PC.",
						},
						"service_plan_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan name of the cloud PC.",
						},
						"service_plan_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan type of the cloud PC. Possible values are: `enterprise`, `business`, `unknownFutureValue`.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the cloud PC.",
						},
						"user_principal_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user principal name (UPN) of the user assigned to the cloud PC.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last modified date and time of the cloud PC.",
						},
						"status_details": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The details of the cloud PC status.",
						},
						"grace_period_end_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the grace period ends and reprovisioning/deprovisioning happens.",
						},
						"provisioning_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Specifies the type of license used when provisioning cloud PCs using this policy.",
						},
						"device_region_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the geographical region where the cloud PC is currently provisioned.",
						},
						"disk_encryption_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The disk encryption applied to the cloud PC.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

// Configure configures the data source with the provider client
func (d *CloudPCDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
