package graphBetaAdministrativeUnitDirectoryRoleAssignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AdministrativeUnitDirectoryRoleAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AdministrativeUnitDirectoryRoleAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AdministrativeUnitDirectoryRoleAssignmentResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &AdministrativeUnitDirectoryRoleAssignmentResource{}
)

func NewAdministrativeUnitDirectoryRoleAssignmentResource() resource.Resource {
	return &AdministrativeUnitDirectoryRoleAssignmentResource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"RoleManagement.ReadWrite.Directory",
		},
	}
}

type AdministrativeUnitDirectoryRoleAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *AdministrativeUnitDirectoryRoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AdministrativeUnitDirectoryRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource using a composite ID.
//
// Format: "administrative_unit_id/scoped_role_membership_id"
//
// Example:
//
//	terraform import microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment.example "11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222"
func (r *AdministrativeUnitDirectoryRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected format 'administrative_unit_id/scoped_role_membership_id', got: %s", req.ID),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with administrative_unit_id: %s, id: %s", ResourceName, parts[0], parts[1]))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("administrative_unit_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *AdministrativeUnitDirectoryRoleAssignmentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *AdministrativeUnitDirectoryRoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a scoped role assignment for an administrative unit in Microsoft Entra ID using the " +
			"`/administrativeUnits/{id}/scopedRoleMembers` endpoint. " +
			"Scoped role members allow directory roles (such as User Administrator or Helpdesk Administrator) to be delegated " +
			"to a user within the scope of a specific administrative unit rather than the entire tenant. " +
			"All fields are immutable after creation; any change triggers a destroy and recreate.\n\n" +
			"**Required permissions:** `RoleManagement.ReadWrite.Directory`\n\n" +
			"**Import format:** `administrative_unit_id/scoped_role_membership_id`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the scoped role membership. Assigned by the API on creation. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"administrative_unit_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the administrative unit to which this role assignment is scoped. " +
					"Changing this value forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
			},
			"directory_role_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The tenant-specific object ID of the activated directoryRole to assign within the administrative unit scope. " +
					"This is **not** the well-known roleTemplateId — it is the object ID of the directoryRole as activated in your tenant. " +
					"Use `GET /directoryRoles` to list activated roles and find the correct object ID. " +
					"Only roles that support administrative unit scoping are valid (e.g. User Administrator, Helpdesk Administrator). " +
					"Changing this value forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
			},
			"role_member_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The unique identifier of the user or service principal to assign the directory role to. " +
					"Changing this value forces a new resource to be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
			},
			"role_member_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the role member. Populated by the API after creation. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_member_user_principal_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user principal name (UPN) of the role member. Populated by the API after creation. Read-only.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
