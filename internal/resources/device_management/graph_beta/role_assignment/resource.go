package graphBetaRoleDefinitionAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_role_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleDefinitionAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleDefinitionAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleDefinitionAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleDefinitionAssignmentResource{}
)

func NewRoleDefinitionAssignmentResource() resource.Resource {
	return &RoleDefinitionAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementRBAC.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementRBAC.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/roleAssignments",
	}
}

type RoleDefinitionAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RoleDefinitionAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *RoleDefinitionAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *RoleDefinitionAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *RoleDefinitionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource
func (r *RoleDefinitionAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages role assignments in Microsoft Intune using the `/deviceManagement/roleAssignments` endpoint. Role assignments bind role definitions to administrators and define the scope of resources they can manage, enabling delegation of administrative privileges across device management, policy configuration, and user licensing operations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Key of the Role Assignment. This is read-only and automatically generated.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"role_definition_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role definition to assign. Either this or `built_in_role_name` must be specified.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"built_in_role_name": schema.StringAttribute{
				MarkdownDescription: "The name of the built-in role to assign. Either this or `role_definition_id` must be specified.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Policy and Profile manager",
						"School Administrator",
						"Help Desk Operator",
						"Application Manager",
						"Endpoint Security Manager",
						"Read Only Operator",
						"Intune Role Administrator",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display or friendly name of the role assignment.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the role assignment.",
				Optional:            true,
			},
			"scope_members": schema.SetAttribute{
				MarkdownDescription: "Group IDs that are assigned as members of this role scope. Also known as admin_group_users_group_ids in the API.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
				},
			},
			"scope_type": schema.StringAttribute{
				MarkdownDescription: "Administrators in this role assignment can target policies, applications and remote tasks to a scope type of:" +
					"'allDevices', 'allLicensedUsers', 'allDevicesAndLicensedUsers' and 'resourceScope'. If the scope intent is for a Entra ID group then leave this empty. " +
					"Possible values are: `allDevices`, `allLicensedUsers`, `allDevicesAndLicensedUsers`, `resourceScope`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"allDevices",
						"allLicensedUsers",
						"allDevicesAndLicensedUsers",
						"resourceScope",
					),
				},
			},
			"resource_scopes": schema.SetAttribute{
				MarkdownDescription: "Administrators in this role assignment can target policies, applications and remote tasks. List of IDs of role scope member security groups from Entra ID.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
