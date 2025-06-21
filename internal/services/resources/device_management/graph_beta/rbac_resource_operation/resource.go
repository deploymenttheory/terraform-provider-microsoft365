package graphBetaRBACResourceOperation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_rbac_resource_operation"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RBACResourceOperationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RBACResourceOperationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RBACResourceOperationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RBACResourceOperationResource{}
)

func NewRBACResourceOperationResource() resource.Resource {
	return &RBACResourceOperationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementRBAC.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementRBAC.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/RBACResourceOperations",
	}
}

type RBACResourceOperationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RBACResourceOperationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *RBACResourceOperationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *RBACResourceOperationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *RBACResourceOperationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *RBACResourceOperationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages resource operations in Microsoft Intune RBAC using the `/deviceManagement/RBACResourceOperations` endpoint. Resource operations define granular permissions that can be included in custom role definitions, enabling precise control over what actions administrators can perform on specific Intune resources and configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Key of the Resource Operation. Read-only, automatically generated.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource": schema.StringAttribute{
				MarkdownDescription: "Resource category to which this Operation belongs. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_name": schema.StringAttribute{
				MarkdownDescription: "Name of the Resource this operation is performed on.",
				Required:            true,
			},
			"action_name": schema.StringAttribute{
				MarkdownDescription: "Type of action this operation is going to perform. The actionName should be concise and limited to as few words as possible.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the resource operation. The description is used in mouse-over text for the operation when shown in the Azure Portal.",
				Required:            true,
			},
			"enabled_for_scope_validation": schema.BoolAttribute{
				MarkdownDescription: "Determines whether the Permission is validated for Scopes defined per Role Assignment. This property is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
