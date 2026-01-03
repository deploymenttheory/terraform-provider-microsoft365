package graphBetaGroupOwnerAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_groups_group_owner_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupOwnerAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupOwnerAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupOwnerAssignmentResource{}
)

func NewGroupOwnerAssignmentResource() resource.Resource {
	return &GroupOwnerAssignmentResource{
		ReadPermissions: []string{
			"GroupMember.Read.All",
			"Directory.Read.All",
			"Group.Read.All",
		},
		WritePermissions: []string{
			"Group.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groups",
	}
}

type GroupOwnerAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupOwnerAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupOwnerAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *GroupOwnerAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupOwnerAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Azure AD/Entra group owner assignments using the `/groups/{group-id}/owners` endpoint. This resource enables adding and removing users or service principals as owners of security groups and Microsoft 365 groups.\n\n" +
			"**Owner Type Support by Group Type:**\n" +
			"- **Security Groups**: Users and Service principals\n" +
			"- **Microsoft 365 Groups**: Users and Service principals\n" +
			"- **Mail-enabled Security Groups**: Read-only, cannot add owners\n" +
			"- **Distribution Groups**: Read-only, cannot add owners\n\n" +
			"**Important Notes:**\n" +
			"- Owners are allowed to modify the group object\n" +
			"- The last owner (user object) of a group cannot be removed\n" +
			"- If you update group owners and created a team for the group, it can take up to 2 hours for owners to sync with Microsoft Teams\n" +
			"- If you want the owner to make changes in a team (e.g., creating a Planner plan), the owner also needs to be added as a group/team member\n\n" +
			"**Required Permissions by Owner Type:**\n" +
			"- **Users**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`\n" +
			"- **Service Principals**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`\n" +
			"- **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required\n\n" +
			"**Supported Microsoft Entra Roles:**\n" +
			"- Group owners (can modify all types of group owners)\n" +
			"- Groups Administrator (can modify all types of group owners)\n" +
			"- User Administrator (can modify user owners only)\n" +
			"- Directory Writers (can modify user owners only)\n" +
			"- Exchange Administrator (Microsoft 365 groups only)\n" +
			"- SharePoint Administrator (Microsoft 365 groups only)\n" +
			"- Teams Administrator (Microsoft 365 groups only)\n" +
			"- Yammer Administrator (Microsoft 365 groups only)\n" +
			"- Intune Administrator (security groups only)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this group owner assignment. This is a composite ID formed by combining the group ID and owner ID.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the group.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the owner to be added to the group. This can be a user or service principal.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_object_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of object being added as an owner. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'ServicePrincipal'. Both security groups and Microsoft 365 groups support both types.",
				Validators: []validator.String{
					stringvalidator.OneOf("User", "ServicePrincipal"),
				},
			},
			"owner_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the owner object as returned by Microsoft Graph (e.g., 'User', 'ServicePrincipal'). Read-only.",
			},
			"owner_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the owner. Read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
