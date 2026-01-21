package graphBetaGroupMemberAssignment

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
	ResourceName  = "microsoft365_graph_beta_groups_group_member_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupMemberAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupMemberAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupMemberAssignmentResource{}
)

func NewGroupMemberAssignmentResource() resource.Resource {
	return &GroupMemberAssignmentResource{
		ReadPermissions: []string{
			"GroupMember.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"GroupMember.ReadWrite.All",
			"Directory.ReadWrite.All",
			"Device.ReadWrite.All",               // Required for adding devices to groups
			"Application.ReadWrite.All",          // Required for adding service principals to groups
			"OrgContact.Read.All",                // Required for adding organizational contacts to groups
			"RoleManagement.ReadWrite.Directory", // Required for role-assignable groups
		},
		ResourcePath: "/groups",
	}
}

type GroupMemberAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupMemberAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupMemberAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *GroupMemberAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupMemberAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Azure AD/Entra group member assignments using the `/groups/{group-id}/members` endpoint. This resource is used to enables adding and removing users, groups, service principals, devices, and organizational contacts as members of security groups and Microsoft 365 groups.\n\n**Member Type Support by Group Type:**\n- **Security Groups**: Users, other Security groups, Devices, Service principals, and Organizational contacts\n- **Microsoft 365 Groups**: Only Users are supported\n- **Mail-enabled Security Groups**: Read-only, cannot add members\n- **Distribution Groups**: Read-only, cannot add members\n\n**Important Notes:**\n- The resource automatically validates member compatibility with the target group type\n- When adding a Group as a member, both the target and member groups must be Security groups\n- Microsoft 365 groups cannot be members of any group type\n\n**Required Permissions by Member Type:**\n- **Users**: `GroupMember.ReadWrite.All`\n- **Groups**: `GroupMember.ReadWrite.All`\n- **Devices**: `GroupMember.ReadWrite.All` + `Device.ReadWrite.All`\n- **Service Principals**: `GroupMember.ReadWrite.All` + `Application.ReadWrite.All`\n- **Organizational Contacts**: `GroupMember.ReadWrite.All` + `OrgContact.Read.All`\n- **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this group member assignment. This is a composite ID formed by combining the group ID and member ID.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the group.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"member_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the member to be added to the group. This can be a user, group, device, service principal, or organizational contact.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"member_object_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of object being added as a member. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'Group', 'Device', 'ServicePrincipal', 'OrganizationalContact'. Note: Microsoft 365 groups only support 'User' and 'Group' (where Group must be a security group), while security groups support all types.",
				Validators: []validator.String{
					stringvalidator.OneOf("User", "Group", "Device", "ServicePrincipal", "OrganizationalContact"),
				},
			},
			"member_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the member object as returned by Microsoft Graph (e.g., 'User', 'Group', 'Device', 'ServicePrincipal', 'OrganizationalContact'). Read-only.",
			},
			"member_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the member. Read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
