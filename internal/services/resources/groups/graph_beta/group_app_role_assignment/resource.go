package graphBetaGroupAppRoleAssignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
	ResourceName  = "microsoft365_graph_beta_groups_group_app_role_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &GroupAppRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &GroupAppRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &GroupAppRoleAssignmentResource{}
)

func NewGroupAppRoleAssignmentResource() resource.Resource {
	return &GroupAppRoleAssignmentResource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AppRoleAssignment.ReadWrite.All",
			"Group.Read.All",
		},
		ResourcePath: "/groups",
	}
}

type GroupAppRoleAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *GroupAppRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *GroupAppRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *GroupAppRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: {group_id}/{assignment_id}
	parts := strings.Split(req.ID, "/")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID Format",
			fmt.Sprintf("Expected import ID format: {group_id}/{assignment_id}, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("target_group_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *GroupAppRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Azure AD/Entra group app role assignments using the `/groups/{group-id}/appRoleAssignments` endpoint. " +
			"This resource enables assigning app roles to security groups, allowing all direct members of the group to inherit the assigned permissions. " +
			"Security groups with dynamic memberships are supported.\n\n" +
			"**Important Notes:**\n" +
			"- All direct members of the assigned group will be considered as having the app role\n" +
			"- Additional licenses might be required to use a group to manage access to applications\n" +
			"- The resource requires three key identifiers: principal ID (group), resource ID (service principal), and app role ID\n\n" +
			"**Required Permissions:**\n" +
			"- `AppRoleAssignment.ReadWrite.All` + `Group.Read.All` (least privileged)\n" +
			"- Delegated scenarios: The signed-in user must be assigned one of the supported Microsoft Entra roles (Directory Readers, Directory Writers, Application Administrator, Cloud Application Administrator, etc.)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this app role assignment.",
			},
			"target_group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the group to which you are assigning the app role. This is the principal ID.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"resource_object_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the resource service principal that has defined the app role. This is the service principal ID of the application.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"app_role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the app role defined on the resource service principal to assign to the group. Use '00000000-0000-0000-0000-000000000000' for the default access role.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"principal_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the group (principal). Read-only.",
			},
			"resource_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the service principal (resource/application). Read-only.",
			},
			"principal_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the principal. For groups, this will always be 'Group'. Read-only.",
			},
			"creation_timestamp": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the app role assignment was created. Read-only.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
