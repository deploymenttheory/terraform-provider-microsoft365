package graphBetaRoleDefinitions

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
	DataSourceName = "microsoft365_graph_beta_identity_and_access_role_definitions"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &RoleDefinitionsDataSource{}
	_ datasource.DataSourceWithConfigure = &RoleDefinitionsDataSource{}
)

func NewRoleDefinitionsDataSource() datasource.DataSource {
	return &RoleDefinitionsDataSource{
		ReadPermissions: []string{
			"RoleManagement.Read.All",
		},
	}
}

type RoleDefinitionsDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (r *RoleDefinitionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *RoleDefinitionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *RoleDefinitionsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves directory role definitions from Microsoft Entra ID using the `/roleManagement/directory/roleDefinitions` endpoint. This data source is used to query built-in and custom role definitions with their permissions and scopes.",
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
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all' or 'odata'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: isBuiltIn eq true.",
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
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: displayName.",
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
				MarkdownDescription: "The list of role definitions that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the role definition.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description for the role definition.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name for the role definition.",
						},
						"is_built_in": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag indicating if the role is built in.",
						},
						"is_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag indicating if the role is enabled.",
						},
						"is_privileged": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag indicating if the role is privileged.",
						},
						"resource_scopes": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of scopes permissions granted by the role definition apply to.",
						},
						"template_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Custom template identifier that can be set when isBuiltIn is false.",
						},
						"version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates the version of the role definition.",
						},
						"role_permissions": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of permissions included in the role.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"allowed_resource_actions": schema.ListAttribute{
										Computed:            true,
										ElementType:         types.StringType,
										MarkdownDescription: "Set of tasks that can be performed on a resource.",
									},
									"condition": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Optional constraints that must be met for the permission to be effective.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
