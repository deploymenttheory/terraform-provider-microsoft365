package graphBetaGroupPolicyBooleanValue

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
	ResourceName  = "graph_beta_device_management_group_policy_boolean_value"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupPolicyBooleanValueResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupPolicyBooleanValueResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupPolicyBooleanValueResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &GroupPolicyBooleanValueResource{}
)

func NewGroupPolicyBooleanValueResource() resource.Resource {
	return &GroupPolicyBooleanValueResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type GroupPolicyBooleanValueResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *GroupPolicyBooleanValueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *GroupPolicyBooleanValueResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupPolicyBooleanValueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupPolicyBooleanValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *GroupPolicyBooleanValueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group policy presentation boolean values in Microsoft Intune using the `/deviceManagement/groupPolicyConfigurations/{groupPolicyConfigurationId}/definitionValues/{groupPolicyDefinitionValueId}/presentationValues` endpoint. This resource represents multiple boolean values for group policy presentations such as checkboxes or radio buttons within a single policy definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the group policy definition value (not individual presentation values)",
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
			"enabled": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the group policy setting is enabled or disabled",
			},
			"values": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of boolean presentation values for this group policy definition. Each presentation corresponds to a different checkbox or setting within the policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"presentation_id": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The unique identifier of the group policy presentation template. If not provided, presentations will be auto-resolved in order.",
						},
						"value": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "The boolean value for this specific presentation",
						},
					},
				},
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
