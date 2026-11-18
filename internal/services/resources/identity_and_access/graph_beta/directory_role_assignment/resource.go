package graphBetaDirectoryRoleAssignment

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_directory_role_assignment"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &DirectoryRoleAssignmentResource{}
	_ resource.ResourceWithConfigure   = &DirectoryRoleAssignmentResource{}
	_ resource.ResourceWithImportState = &DirectoryRoleAssignmentResource{}
	_ resource.ResourceWithIdentity    = &DirectoryRoleAssignmentResource{}
)

type DirectoryRoleAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func NewDirectoryRoleAssignmentResource() resource.Resource {
	return &DirectoryRoleAssignmentResource{ReadPermissions: []string{"RoleManagement.Read.Directory"}, WritePermissions: []string{"RoleManagement.ReadWrite.Directory"}}
}

func (r *DirectoryRoleAssignmentResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *DirectoryRoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *DirectoryRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if strings.TrimSpace(req.ID) == "" || strings.Contains(req.ID, "/") {
		resp.Diagnostics.AddError("Invalid import ID", "Supply the opaque Graph directory role assignment ID without an API path.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DirectoryRoleAssignmentResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{"id": identityschema.StringAttribute{RequiredForImport: true}}}
}

func (r *DirectoryRoleAssignmentResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	guid := stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a permanent active Microsoft Entra directory role assignment at tenant scope using Microsoft Graph beta `/roleManagement/directory/roleAssignments`. This grants active access; it does not create PIM eligibility or a scheduled assignment. Changing the role or principal replaces the assignment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Opaque Graph assignment ID; this is not necessarily a UUID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"role_definition_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Built-in role template ID or custom directory role definition ID. An activated directoryRole object ID is a different identifier.",
				Validators:          []validator.String{guid},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"principal_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Object ID of a user, service principal, or role-assignable group receiving the role.",
				Validators:          []validator.String{guid},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"directory_scope_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Directory scope. This resource manages tenant-wide assignments only, represented by `/`.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
