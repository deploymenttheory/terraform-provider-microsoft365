package graphBetaCloudPcFrontlineServicePlan

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_windows_365_cloud_pc_frontline_service_plan"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcFrontlineServicePlanDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcFrontlineServicePlanDataSource{}
)

func NewCloudPcFrontlineServicePlanDataSource() datasource.DataSource {
	return &CloudPcFrontlineServicePlanDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcFrontlineServicePlanDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *CloudPcFrontlineServicePlanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *CloudPcFrontlineServicePlanDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *CloudPcFrontlineServicePlanDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PC frontline service plans from Microsoft Intune using the `/deviceManagement/virtualEndpoint/frontlineServicePlans` endpoint. This data source is used to query shared Cloud PC plans and their capacity utilization.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all frontline service plans, 'id' to retrieve a specific plan by its unique identifier, or 'display_name' to filter by the plan's display name.",
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the frontline service plan ID. For 'display_name', provide a substring to match against plan display names.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PC frontline service plans that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the frontline service plan.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the frontline service plan. For example, '2vCPU/8GB/128GB Front-line' or '4vCPU/16GB/256GB Front-line'.",
						},
						"total_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The total number of frontline service plans purchased by the customer.",
						},
						"used_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of service plans that have been used for the account.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
