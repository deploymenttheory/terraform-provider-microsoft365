package graphBetaPolicySet

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	assignmentschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
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
	ResourceName  = "graph_beta_device_management_policy_set"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &PolicySetResource{}
	_ resource.ResourceWithConfigure   = &PolicySetResource{}
	_ resource.ResourceWithImportState = &PolicySetResource{}
	_ resource.ResourceWithModifyPlan  = &PolicySetResource{}
)

func NewPolicySetResource() resource.Resource {
	return &PolicySetResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/policySets",
	}
}

type PolicySetResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *PolicySetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *PolicySetResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *PolicySetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *PolicySetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PolicySetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages policy sets in Microsoft Intune using the `/deviceAppManagement/policySets` endpoint. Policy sets allow you to group multiple device management policies and applications together for unified deployment and assignment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this policy set",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the policy set",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The optional description of the policy set",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the policy set (e.g., notAssigned, assigned)",
			},
			"error_code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The error code for the policy set (e.g., noError)",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy set was created",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the policy set was last modified",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this policy set.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"assignments": assignmentschema.InclusionGroupAndExclusionGroupAssignmentsSchema(),
			"items": schema.SetNestedAttribute{
				Optional: true,
				MarkdownDescription: "Set of policy set items ('apps', 'app configuration policies', 'app protection policies', " +
					"'device configuration profiles', 'device management configuration policies', 'device compliance policies', " +
					"'windows autopilot deployment profiles') included in this policy set.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"payload_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The ID of the policy or application being included",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidOrPrefixedGuidRegex),
									"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000 or a prefixed GUID in the format 0_00000000-0000-0000-0000-000000000000",
								),
							},
						},
						"type": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The type of policy set item. Valid values: 'app', 'app_configuration_policy', 'app_protection_policy', " +
								"'device_configuration_profile', 'device_management_configuration_policy', 'device_compliance_policy', 'windows_autopilot_deployment_profile'",
						},
						"intent": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The intent for mobile app policy set items (e.g., required, available)",
						},
						"settings": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Settings specific to the policy set item type",
							Attributes: map[string]schema.Attribute{
								"odata_type": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The OData type of the settings",
								},
								"vpn_configuration_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "VPN configuration ID for iOS store app assignment settings",
								},
								"uninstall_on_device_removal": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether to uninstall the app when the device is removed",
								},
								"is_removable": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether the app is removable by the user",
								},
								"prevent_managed_app_backup": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether to prevent managed app backup",
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
