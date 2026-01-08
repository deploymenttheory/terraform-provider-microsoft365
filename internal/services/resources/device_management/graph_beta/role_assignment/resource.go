package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_role_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleAssignmentResource{}
)

func NewRoleAssignmentResource() resource.Resource {
	return &RoleAssignmentResource{
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

type RoleAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RoleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *RoleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *RoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fmt.Printf("DEBUG: ImportState called with ID: %s\n", req.ID)

	// Parse composite ID format: id/role_definition_id
	idParts := strings.Split(req.ID, "/")
	fmt.Printf("DEBUG: Parsed ID parts: %v (length: %d)\n", idParts, len(idParts))

	if len(idParts) != 2 {
		fmt.Printf("DEBUG: Invalid ID format error\n")
		resp.Diagnostics.AddError(
			"Invalid Import ID Format",
			"Import ID must be in format: id/role_definition_id",
		)
		return
	}

	fmt.Printf("DEBUG: Setting id to: %s\n", idParts[0])
	fmt.Printf("DEBUG: Setting role_definition_id to: %s\n", idParts[1])

	// Set the resource ID and role definition ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_definition_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		fmt.Printf("DEBUG: Diagnostics has errors after setting attributes\n")
		for _, diag := range resp.Diagnostics.Errors() {
			fmt.Printf("DEBUG: Error: %s - %s\n", diag.Summary(), diag.Detail())
		}
		return
	}

	fmt.Printf("DEBUG: Successfully set attributes, now calling Read to populate remaining fields\n")

	// After setting the basic attributes, call Read to populate all other attributes from the API
	readReq := resource.ReadRequest{State: resp.State}
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, readReq, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	resp.State = readResp.State

	if resp.Diagnostics.HasError() {
		fmt.Printf("DEBUG: Diagnostics has errors after Read\n")
		for _, diag := range resp.Diagnostics.Errors() {
			fmt.Printf("DEBUG: Read Error: %s - %s\n", diag.Summary(), diag.Detail())
		}
	} else {
		fmt.Printf("DEBUG: Successfully completed import with Read\n")
	}
}

func (r *RoleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages role assignments in Microsoft Intune using the `/deviceManagement/roleAssignments` endpoint. Role assignments bind role definitions to users or groups with specific scope configurations, enabling granular access control for device management and administrative functions within Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity. This is read-only and automatically generated.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name or name of the role Assignment.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"role_definition_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				MarkdownDescription: "Role definition this assignment is for. Either use a built-in role ID or reference a custom role definition." +
					"Built-in role id's are as follows: " +
					"Policy and Profile manager: 0bd113fe-6be5-400c-a28f-ae5553f9c0be" +
					"School Administrator: 2f9f4f7e-2d13-427b-adf2-361a1eef7ae8" +
					"Help Desk Operator: 9e0cc482-82df-4ab2-a24c-0c23a3f52e1e" +
					"Application Manager: c1d9fcbb-cba5-40b0-bf6b-527006585f4b" +
					"Endpoint Security Manager: c56d53a2-73d0-4502-b6bd-4a9d3dba28d5" +
					"Read Only Operator: fa1d7878-e8cb-41a1-8254-0142355c9f84" +
					"Intune Role Administrator: fb2603eb-3c87-4be3-8b5b-d58a5b4a0bc0",
			},
			"members": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The list of ids of principals that are assigned to the role. These can be user ids or group ids.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"scope_configuration": schema.ListNestedBlock{
				MarkdownDescription: "Defines the scope configuration for the role assignment. Exactly one scope configuration block is required.",
				Validators:          []validator.List{
					// Exactly one scope configuration is required
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ResourceScopes", "AllLicensedUsers", "AllDevices"),
							},
							MarkdownDescription: "The type of scope configuration. Valid values are: `ResourceScopes`, `AllLicensedUsers`, `AllDevices`.",
						},
						"resource_scopes": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "List of resource scope IDs. Required when type is `ResourceScopes`, must be empty for other types.",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.GuidRegex),
										"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
									),
								),
							},
						},
					},
				},
			},
		},
	}
}
