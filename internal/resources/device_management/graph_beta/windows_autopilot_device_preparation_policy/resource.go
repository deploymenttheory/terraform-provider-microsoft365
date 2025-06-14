package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_autopilot_device_preparation_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsAutopilotDevicePreparationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsAutopilotDevicePreparationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsAutopilotDevicePreparationPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsAutopilotDevicePreparationPolicyResource{}
)

func NewWindowsAutopilotDevicePreparationPolicyResource() resource.Resource {
	return &WindowsAutopilotDevicePreparationPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type WindowsAutopilotDevicePreparationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *WindowsAutopilotDevicePreparationPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsAutopilotDevicePreparationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full Windows Autopilot Device Preparation Policy schema
func (r *WindowsAutopilotDevicePreparationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Autopilot Device Preparation Policy using the `/deviceManagement/configurationPolicies` endpoint. Windows Autopilot Device Preparation is used to set up and configure new devices, getting them ready for productive use by delivering consistent configurations and enhancing the setup experience.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this policy",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				MarkdownDescription: "Optional description for the Windows Autopilot Device Preparation policy.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Entity instance.",
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{types.StringValue("0")},
					),
				),
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Creation date and time of the policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the policy",
			},
			"settings_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings with the policy. This will change over time as the resource is updated.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "Indicates if the policy is assigned to any scope",
			},
			"platforms": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The platforms this policy applies to (e.g., windows10)",
			},
			"technologies": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The technology this policy is using (e.g., enrollment)",
			},
			"template_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The template ID used by this policy",
			},
			"template_family": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The template family for this policy (e.g., enrollmentConfiguration)",
			},
			"device_security_group": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the assigned security device group that devices will be automatically added to during the Windows Autopilot Device Preparation flow. This group must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"deployment_settings": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Deployment settings for the Windows Autopilot Device Preparation policy",
				Attributes: map[string]schema.Attribute{
					"deployment_mode": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The deployment mode for the Windows Autopilot Device Preparation policy. Valid values are: 'enrollment_autopilot_dpp_deploymentmode_0' (Standard mode) or 'enrollment_autopilot_dpp_deploymentmode_1' (Enhanced mode).",
						Default:             stringdefault.StaticString("enrollment_autopilot_dpp_deploymentmode_0"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"enrollment_autopilot_dpp_deploymentmode_0", // Standard mode
								"enrollment_autopilot_dpp_deploymentmode_1", // Enhanced mode
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"deployment_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The deployment type for the Windows Autopilot Device Preparation policy. Valid values are: 'enrollment_autopilot_dpp_deploymenttype_0' (User-driven) or 'enrollment_autopilot_dpp_deploymenttype_1' (Self-deploying).",
						Default:             stringdefault.StaticString("enrollment_autopilot_dpp_deploymenttype_0"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"enrollment_autopilot_dpp_deploymenttype_0", // User-driven
								"enrollment_autopilot_dpp_deploymenttype_1", // Self-deploying
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"join_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The join type for the Windows Autopilot Device Preparation policy. Valid values are: 'enrollment_autopilot_dpp_jointype_0' (Entra ID joined) or 'enrollment_autopilot_dpp_jointype_1' (Entra ID hybrid joined).",
						Default:             stringdefault.StaticString("enrollment_autopilot_dpp_jointype_0"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"enrollment_autopilot_dpp_jointype_0", // Entra ID joined
								"enrollment_autopilot_dpp_jointype_1", // Entra ID hybrid joined
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"account_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The account type for users in the Windows Autopilot Device Preparation policy. Valid values are: 'enrollment_autopilot_dpp_accountype_0' (Standard User) or 'enrollment_autopilot_dpp_accountype_1' (Administrator).",
						Default:             stringdefault.StaticString("enrollment_autopilot_dpp_accountype_0"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"enrollment_autopilot_dpp_accountype_0", // Standard User
								"enrollment_autopilot_dpp_accountype_1", // Administrator
							),
						},
					},
				},
			},
			"oobe_settings": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Out-of-box experience settings for the Windows Autopilot Device Preparation policy",
				Attributes: map[string]schema.Attribute{
					"timeout_in_minutes": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The timeout in minutes for the Windows Autopilot Device Preparation policy. Valid range is 15-720 minutes.",
						Default:             int64default.StaticInt64(60),
						Validators: []validator.Int64{
							int64validator.Between(15, 720),
						},
					},
					"custom_error_message": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The custom error message to display if the deployment fails. Maximum length is 1000 characters.",
						Default:             stringdefault.StaticString("Contact your organization's support person for help."),
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 1000),
						},
					},
					"allow_skip": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Whether to allow users to skip setup after multiple failed attempts",
						Default:             booldefault.StaticBool(false),
					},
					"allow_diagnostics": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Whether to allow users to access diagnostics information during setup",
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"allowed_apps": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of applications that are allowed to be installed during the Windows Autopilot Device Preparation process. Maximum of 10 items.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The ID of the application.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
								),
							},
						},
						"app_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of the application. Valid values are: 'winGetApp', 'win32LobApp', 'officeSuiteApp', 'windowsUniversalAppX'.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"winGetApp",
									"win32LobApp",
									"officeSuiteApp",
									"windowsUniversalAppX",
								),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
			},
			"allowed_scripts": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of script IDs that are allowed to be executed during the Windows Autopilot Device Preparation process. Maximum of 10 items.",
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					),
				},
			},
			"assignments": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The assignment configuration for this Windows Autopilot Device Preparation policy",
				Attributes: map[string]schema.Attribute{
					"include_group_ids": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "A list of user group IDs to include in the assignment",
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
								),
							),
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
