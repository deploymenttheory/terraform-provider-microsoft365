package graphBetaAssignmentFilter

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
	ResourceName = "graph_beta_device_management_assignment_filter"
	ReadTimeout  = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &AssignmentFilterDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &AssignmentFilterDataSource{}
)

func NewAssignmentFilterDataSource() datasource.DataSource {
	return &AssignmentFilterDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type AssignmentFilterDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata sets the data source name
func (d *AssignmentFilterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the data source
func (d *AssignmentFilterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *AssignmentFilterDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves assignment filters from Microsoft Intune using the `/deviceManagement/assignmentFilters` endpoint. Assignment filters allow targeting of policies and apps based on device properties like OS version, manufacturer, or custom device attributes, enabling more granular deployment control beyond basic group membership.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Assignment Filters that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the assignment filter.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the assignment filter.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the assignment filter.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
