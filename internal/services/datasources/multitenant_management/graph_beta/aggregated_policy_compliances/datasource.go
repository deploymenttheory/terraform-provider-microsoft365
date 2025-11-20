package graphBetaAggregatedPolicyCompliances

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
	DataSourceName = "microsoft365_graph_beta_multitenant_management_aggregated_policy_compliances"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &AggregatedPolicyCompliancesDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &AggregatedPolicyCompliancesDataSource{}
)

func NewAggregatedPolicyCompliancesDataSource() datasource.DataSource {
	return &AggregatedPolicyCompliancesDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type AggregatedPolicyCompliancesDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

// Metadata returns the datasource type name.
func (r *AggregatedPolicyCompliancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source
func (d *AggregatedPolicyCompliancesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *AggregatedPolicyCompliancesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves aggregated policy compliances from Microsoft 365 managed tenants using the `/tenantRelationships/managedTenants/aggregatedPolicyCompliances` endpoint. This data source provides an aggregate view of device compliance for managed tenants with advanced filtering capabilities.",
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
				MarkdownDescription: "The list of aggregated policy compliances that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier for the aggregated policy compliance.",
						},
						"compliance_policy_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the compliance policy.",
						},
						"compliance_policy_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the compliance policy.",
						},
						"compliance_policy_platform": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The platform for the compliance policy (e.g., 'Windows10', 'iOS', 'Android').",
						},
						"compliance_policy_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the compliance policy.",
						},
						"last_refreshed_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the compliance data was last refreshed.",
						},
						"number_of_compliant_devices": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of devices that are compliant with the policy.",
						},
						"number_of_error_devices": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of devices that encountered errors during compliance evaluation.",
						},
						"number_of_non_compliant_devices": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of devices that are not compliant with the policy.",
						},
						"policy_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the policy was last modified.",
						},
						"tenant_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the managed tenant.",
						},
						"tenant_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure Active Directory tenant identifier for the managed tenant.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}