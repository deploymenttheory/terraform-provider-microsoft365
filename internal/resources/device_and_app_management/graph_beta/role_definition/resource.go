package graphBetaRoleDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_role_definition"
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
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RoleDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *RoleDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *RoleDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RoleDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The resource `role_definition` manages a Role Definition in Microsoft 365",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune role definition",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display Name of the Role definition.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the Role definition.",
				Optional:            true,
			},
			"is_built_in": schema.BoolAttribute{
				MarkdownDescription: "Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.",
				Optional:            true,
			},
			"is_built_in_role_definition": schema.BoolAttribute{
				MarkdownDescription: "Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.",
				Optional:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				MarkdownDescription: "List of Scope Tags for this Entity instance.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"permissions": schema.ListNestedAttribute{
				MarkdownDescription: "List of Role Permissions this role is allowed to perform. These must match the actionName that is defined as part of the rolePermission.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListAttribute{
							MarkdownDescription: "Allowed actions for this role permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"resource_actions": schema.ListNestedAttribute{
							MarkdownDescription: "Resource actions for this role permission",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"allowed_resource_actions": schema.ListAttribute{
										MarkdownDescription: "Allowed Resource Actions",
										Optional:            true,
										ElementType:         types.StringType,
									},
									"not_allowed_resource_actions": schema.ListAttribute{
										MarkdownDescription: "Not Allowed Resource Actions",
										Optional:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
					},
				},
			},
			"role_permissions": schema.ListNestedAttribute{
				MarkdownDescription: "List of Role Permissions this role is allowed to perform. These must match the actionName that is defined as part of the rolePermission.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListAttribute{
							MarkdownDescription: "Allowed actions for this role permission",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"resource_actions": schema.ListNestedAttribute{
							MarkdownDescription: "Resource actions for this role permission",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"allowed_resource_actions": schema.ListAttribute{
										MarkdownDescription: "Allowed Resource Actions",
										Optional:            true,
										ElementType:         types.StringType,
									},
									"not_allowed_resource_actions": schema.ListAttribute{
										MarkdownDescription: "Not Allowed Resource Actions",
										Optional:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.RoleAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
