package graphBetaGroupPolicyTextValue

import (
	"context"

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
	ResourceName  = "graph_beta_device_management_group_policy_text_value"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupPolicyTextValueResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupPolicyTextValueResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupPolicyTextValueResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &GroupPolicyTextValueResource{}
)

func NewGroupPolicyTextValueResource() resource.Resource {
	return &GroupPolicyTextValueResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type GroupPolicyTextValueResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *GroupPolicyTextValueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *GroupPolicyTextValueResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupPolicyTextValueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupPolicyTextValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *GroupPolicyTextValueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group policy text values in Microsoft Intune using the" +
			"`deviceManagement/groupPolicyConfigurations('{groupPolicyConfigurationId}')/updateDefinitionValues` endpoint." +
			"This resource manages singular text values for a given group policy presentations such as text boxes " +
			"within a single group policy definition. It can also be used when integers are represented as strings. Group policy schema dependant." +
			"This resource has a hard dependency on the group policy configuration resource and it must be created before this resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the group policy presentation text value",
			},
			"group_policy_configuration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the group policy configuration that contains this presentation value",
			},
			"policy_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the group policy definition (e.g., 'Allow automatic full screen on specified sites')",
			},
			"class_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The class type of the group policy definition. Must be 'user' or 'machine'",
				Validators: []validator.String{
					stringvalidator.OneOf("user", "machine"),
				},
			},
			"category_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The category path of the group policy definition (e.g., '\\FSLogix\\Profile Containers', '\\FSLogix\\ODFC Containers'). Used to distinguish between policies with the same name in different categories",
			},
			"group_policy_definition_value_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the group policy definition value instance within the configuration (resolved automatically from policy_name and class_type)",
			},
			"presentation_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the group policy presentation template (resolved automatically from the policy definition and presentation_index)",
			},
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the group policy setting is enabled or disabled",
			},
			"value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The text value for the group policy setting",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the presentation value was created",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the presentation value was last modified",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
