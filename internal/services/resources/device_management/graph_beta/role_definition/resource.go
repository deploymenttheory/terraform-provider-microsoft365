package graphBetaRoleDefinition

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_role_definition"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleDefinitionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleDefinitionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleDefinitionResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleDefinitionResource{}
)

func NewRoleDefinitionResource() resource.Resource {
	return &RoleDefinitionResource{
		ReadPermissions: []string{
			"DeviceManagementRBAC.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementRBAC.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/roleDefinitions",
	}
}

type RoleDefinitionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RoleDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *RoleDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *RoleDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RoleDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages custom role definitions in Microsoft Intune using the `/deviceManagement/roleDefinitions` endpoint. Role definitions define sets of permissions that can be assigned to administrators, enabling granular access control for device management, policy configuration, and administrative functions within Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity. This is read-only and automatically generated.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display Name of the Role definition.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Role definition.",
				Optional:            true,
			},
			"is_built_in": schema.BoolAttribute{
				MarkdownDescription: "Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.",
				Computed:            true,
			},
			"is_built_in_role_definition": schema.BoolAttribute{
				MarkdownDescription: "Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"role_permissions": schema.ListNestedAttribute{
				MarkdownDescription: "List of Role Permissions this role is allowed to perform. Not used for in-built Intune role definitions.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_resource_actions": schema.SetAttribute{
							MarkdownDescription: "Allowed actions for this role permission. This field is equivalent to 'actions' and can be used interchangeably. The API will consolidate values from both fields. Each action must start with 'Microsoft.Intune_'.",
							Optional:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^Microsoft\.Intune_`),
										"must start with 'Microsoft.Intune_'",
									),
								),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
