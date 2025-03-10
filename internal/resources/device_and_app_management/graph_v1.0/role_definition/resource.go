package graphRoleDefinition

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	ResourceName  = "graph_device_and_app_management_role_definition"
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
		},
		WritePermissions: []string{
			" DeviceManagementRBAC.ReadWrite.All",
		},
	}
}

type RoleDefinitionResource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *RoleDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *RoleDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphStableClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *RoleDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *RoleDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource `role_definition` manages a Role Definition in Microsoft 365",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Key of the entity. This is read-only and automatically generated.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Display Name of the Role definition.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the Role definition.",
				Optional:    true,
			},
			"is_built_in": schema.BoolAttribute{
				Description: "Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.",
				Optional:    true,
			},
			"role_permissions": schema.ListNestedAttribute{
				Description: "List of Role Permissions this role is allowed to perform. These must match the actionName that is defined as part of the rolePermission.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_resource_actions": schema.ListAttribute{
							Description: "Allowed Resource Actions",
							Optional:    true,
							ElementType: types.StringType,
						},
						"not_allowed_resource_actions": schema.ListAttribute{
							Description: "Not Allowed Resource Actions",
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
