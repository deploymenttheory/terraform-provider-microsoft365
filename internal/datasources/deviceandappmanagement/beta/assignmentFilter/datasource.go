package graphBetaAssignmentFilter

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ datasource.DataSource = &AssignmentFilterDataSource{}
var _ datasource.DataSourceWithConfigure = &AssignmentFilterDataSource{}

func NewAssignmentFilterDataSource() datasource.DataSource {
	return &AssignmentFilterDataSource{}
}

type AssignmentFilterDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

func (d *AssignmentFilterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_assignment_filter"
}

func (d *AssignmentFilterDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the assignment filter.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the assignment filter.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the assignment filter.",
			},
			"platform": schema.StringAttribute{
				Computed:    true,
				Description: "The Intune device management type (platform) for the assignment filter.",
			},
			"rule": schema.StringAttribute{
				Computed:    true,
				Description: "Rule definition of the assignment filter.",
			},
			"assignment_filter_management_type": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates filter is applied to either 'devices' or 'apps' management type.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the assignment filter.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modified time of the assignment filter.",
			},
			"role_scope_tags": schema.ListAttribute{
				Computed:    true,
				Description: "Indicates role scope tags assigned for the assignment filter.",
				ElementType: types.StringType,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *AssignmentFilterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
