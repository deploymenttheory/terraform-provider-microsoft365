package graphBetaRoleScopeTag

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_role_scope_tag"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &RoleScopeTagResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &RoleScopeTagResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &RoleScopeTagResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &RoleScopeTagResource{}
)

func NewRoleScopeTagResource() resource.Resource {
	return &RoleScopeTagResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementRBAC.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementRBAC.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/roleScopeTags",
	}
}

type RoleScopeTagResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *RoleScopeTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *RoleScopeTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *RoleScopeTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the role scope tags schema
func (r *RoleScopeTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages role scope tags in Microsoft Intune using the `/deviceManagement/roleScopeTags` endpoint. Role scope tags enable role-based access control by limiting administrator visibility and management capabilities to specific organizational units, devices, or resources, ensuring proper delegation of administrative responsibilities.",
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
				MarkdownDescription: "The display or friendly name of the Role Scope Tag.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString(""),
				},
				MarkdownDescription: "Description of the Role Scope Tag.",
			},
			"is_built_in": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Defines if the Role Scope Tag is built-in. This property is read-only.",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithInclusionGroupAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
