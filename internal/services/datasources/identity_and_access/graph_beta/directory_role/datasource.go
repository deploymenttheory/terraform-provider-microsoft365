// REF: https://learn.microsoft.com/en-us/graph/api/directoryrole-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directoryrole-get?view=graph-rest-beta

package graphBetaDirectoryRole

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_directory_role"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &DirectoryRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &DirectoryRoleDataSource{}
)

func NewDirectoryRoleDataSource() datasource.DataSource {
	return &DirectoryRoleDataSource{
		ReadPermissions: []string{
			"Directory.Read.All",
			"RoleManagement.Read.Directory",
		},
	}
}

type DirectoryRoleDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *DirectoryRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *DirectoryRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *DirectoryRoleDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves activated directory roles from Microsoft Entra ID using the `/directoryRoles` endpoint. " +
			"Returns tenant-specific directoryRole object IDs required by resources such as " +
			"`microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment`. " +
			"Supports lookup by role object ID, display name, or listing all activated roles.",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The tenant-specific object ID of the activated directory role. Conflicts with `display_name` and `list_all`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("display_name"),
						path.MatchRoot("list_all"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("role_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter activated directory roles by display name (exact match, e.g. 'User Administrator'). Uses OData `$filter=displayName eq '...'`. Conflicts with `role_id` and `list_all`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("role_id"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all activated directory roles in the tenant. Conflicts with `role_id` and `display_name`.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("role_id"),
						path.MatchRoot("display_name"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of activated directory roles that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The tenant-specific object ID of the activated directory role. " +
								"e.g you can use this value as `directory_role_id` in `microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment`.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the directory role (e.g. 'User Administrator').",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the directory role.",
						},
						"role_template_id": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The well-known, cross-tenant role template ID. " +
								"This is the same identifier across all tenants for the same built-in role " +
								"(e.g. `fe930be7-5e62-47db-91af-98c3a49a38b1` for User Administrator).",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
