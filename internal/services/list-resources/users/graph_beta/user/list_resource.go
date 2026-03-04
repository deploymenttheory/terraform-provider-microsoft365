package graphBetaUsersUser

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ListResourceName = "microsoft365_graph_beta_users_user"
)

var (
	_ list.ListResource              = &UserListResource{}
	_ list.ListResourceWithConfigure = &UserListResource{}
)

func NewUserListResource() list.ListResource {
	return &UserListResource{
		ReadPermissions: []string{
			"Directory.Read.All",
			"User.Read.All",
			"User.ReadBasic.All",
		},
		ResourcePath: "/users",
	}
}

type UserListResource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the list resource type name.
func (r *UserListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ListResourceName
}

// Configure sets the client for the list resource.
func (r *UserListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ListResourceName)
	if r.client != nil {
		tflog.Debug(ctx, "Successfully configured list resource client", map[string]any{
			"list_resource": ListResourceName,
		})
	} else {
		tflog.Error(ctx, "Failed to configure list resource client - client is nil", map[string]any{
			"list_resource": ListResourceName,
		})
	}
}

// ListResourceConfigSchema defines the schema for the list resource configuration.
func (r *UserListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists users from Microsoft Entra ID using the `/users` endpoint. " +
			"This list resource is used to automatically retrieve all users across multiple pages with advanced filtering capabilities for user discovery and import. " +
			"For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.",
		Attributes: map[string]listschema.Attribute{
			"display_name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter users by display name using prefix matching. Supports the OData `startsWith` operator. " +
					"Example: `display_name_filter = \"John\"` will match \"John Smith\" and \"Johnny Doe\".",
				Optional: true,
			},
			"user_principal_name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter users by user principal name using partial matching. Supports the OData `startsWith` operator. " +
					"Example: `user_principal_name_filter = \"admin\"` will match \"admin@contoso.com\" and \"admin2@contoso.com\".",
				Optional: true,
			},
			"account_enabled_filter": listschema.BoolAttribute{
				MarkdownDescription: "Filter users by account status. Set to `true` to return only enabled accounts, " +
					"`false` for disabled accounts. Example: `account_enabled_filter = true`.",
				Optional: true,
			},
			"user_type_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter users by type. Valid values: `Member`, `Guest`. " +
					"Example: `user_type_filter = \"Member\"` returns only member users.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("Member", "Guest"),
				},
			},
			"odata_filter": listschema.StringAttribute{
				MarkdownDescription: "Advanced: Custom OData $filter query for complex filtering scenarios. " +
					"Allows direct control over the API filter expression. " +
					"Example: `odata_filter = \"accountEnabled eq true and userType eq 'Member'\"`. " +
					"When specified, this overrides individual filter parameters. " +
					"See Microsoft Graph API documentation for supported operators and syntax.",
				Optional: true,
			},
		},
	}
}
