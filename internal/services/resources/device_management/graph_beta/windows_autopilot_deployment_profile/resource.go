package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ResourceName  = "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsAutopilotDeploymentProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsAutopilotDeploymentProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsAutopilotDeploymentProfileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsAutopilotDeploymentProfileResource{}
)

func NewWindowsAutopilotDeploymentProfileResource() resource.Resource {
	return &WindowsAutopilotDeploymentProfileResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/windowsAutopilotDeploymentProfiles",
	}
}

type WindowsAutopilotDeploymentProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsAutopilotDeploymentProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsAutopilotDeploymentProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsAutopilotDeploymentProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsAutopilotDeploymentProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Autopilot deployment profiles using the `/deviceManagement/windowsAutopilotDeploymentProfiles` endpoint. Autopilot deployment profiles define the out-of-box experience (OOBE) settings, device naming templates, and enrollment configurations for automated Windows device provisioning and domain joining.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The profile key.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the deployment profile. Max allowed length is 200 chars. Cannot contain the following characters: ! # % ^ * ) ( - + ; ' > <",
				Validators: []validator.String{
					validate.StringLengthAtMost(200),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^!#%^*()\-+;'><]*$`),
						"display name cannot contain the following characters: ! # % ^ * ) ( - + ; ' > <",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the windows autopilotdeployment profile. Max allowed length is 1500 chars.",
				Validators: []validator.String{
					validate.StringLengthAtMost(1500),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^!#%^*()\-+;'><]*$`),
						"description cannot contain the following characters: ! # % ^ * ) ( - + ; ' > <",
					),
				},
			},
			"locale": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The locale (language) to be used when configuring the device. Possible values are: `user_select` (allows user to select language during OOBE), `os-default` (uses OS default), or specific country codes like `en-US`, `ja-JP`, `fr-FR`, etc. Default value is `os-default`.",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.OneOf("user_select", "os-default"),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-z]{2}-[A-Z]{2}$`),
							"must be a valid locale code in format 'xx-XX' (e.g., en-US, ja-JP)",
						),
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("os-default"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time of when the deployment profile was created. Read-Only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time of when the deployment profile was last modified. Read-Only.",
			},
			"device_join_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of device join to configure. Determines which Windows Autopilot deployment profile type to use. Possible values are: `microsoft_entra_joined`, `microsoft_entra_hybrid_joined`. Note: HoloLens devices must use `microsoft_entra_joined`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"microsoft_entra_joined",
						"microsoft_entra_hybrid_joined",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hardware_hash_extraction_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Select Yes to register all targeted devices to Autopilot if they are not already registered. " +
					"The next time registered devices go through the Windows Out of Box Experience (OOBE), they will go through the assigned Autopilot scenario." +
					"Please note that certain Autopilot scenarios require specific minimum builds of Windows. Please make sure your device has the required minimum build to go through the scenario." +
					"Removing this profile won't remove affected devices from Autopilot. To remove a device from Autopilot, use the Windows Autopilot Devices view.Default value is FALSE.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefaultValue(false),
				},
				Validators: []validator.Bool{
					validate.BoolCanOnlyBeFalseWhenStringEquals("device_type", "holoLens", "hardware_hash_extraction_enabled must be false when device_type is holoLens"),
				},
			},
			"device_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The template used to name the Autopilot device. This can be a custom text and can also contain either the serial number of the device, or a randomly generated number. The total length of the text generated by the template can be no more than 15 characters. For Microsoft Entra hybrid joined type of Autopilot deployment profiles, devices are named using settings specified in Domain Join configuration.",
				Validators: []validator.String{
					validate.StringLengthAtMost(15),
					validate.StringMustBeEmptyWhenStringEquals("device_join_type", "microsoft_entra_hybrid_joined", "device_name_template must not be set when 'device_join_type' is 'microsoft_entra_hybrid_joined', devices are named using settings specified in the AD Domain Join configuration"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString(""),
				},
			},
			"device_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Windows device type that this profile is applicable to. Possible values include `windowsPc`, `holoLens`, `surfaceHub2`, `surfaceHub2S`, `virtualMachine`, `unknownFutureValue`. The default is `windowsPc`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"windowsPc",
						"holoLens",
					),
				},
			},
			"preprovisioning_allowed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the user is allowed to use Windows Autopilot for pre-provisioned deployment mode during Out of Box experience (OOBE). When TRUE, indicates that Windows Autopilot for pre-provisioned deployment mode for OOBE is allowed to be used. When false, Windows Autopilot for pre-provisioned deployment mode for OOBE is not allowed. The default is FALSE.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefaultValue(false),
				},
				Validators: []validator.Bool{
					validate.BoolCanOnlyBeFalseWhenStringEquals("device_type", "holoLens", "preprovisioning_allowed must be false when device_type is holoLens"),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of role scope tags for the deployment profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"management_service_app_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Entra management service App ID which gets used during client device-based enrollment discovery.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"hybrid_azure_ad_join_skip_connectivity_check": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The Autopilot Hybrid Azure AD join flow will continue even if it does not establish domain controller connectivity during OOBE. " +
					"This should only be set to true when using `microsoft_entra_hybrid_joined` device join type, else always false.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefaultValue(false),
				},
				Validators: []validator.Bool{
					validate.BoolCanOnlyBeFalseWhenStringEquals("device_join_type", "microsoft_entra_joined", "hybrid_azure_ad_join_skip_connectivity_check can only be set to true when device_join_type is microsoft_entra_hybrid_joined"),
					validate.BoolCanOnlyBeFalseWhenStringEquals("device_type", "holoLens", "hybrid_azure_ad_join_skip_connectivity_check must be false when device_type is holoLens"),
				},
			},
			"out_of_box_experience_setting": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The Windows Autopilot Deployment Profile settings used by the device for the out-of-box experience.",
				Attributes: map[string]schema.Attribute{
					"privacy_settings_hidden": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "When TRUE, privacy settings is hidden to the end user during OOBE. When FALSE, privacy settings is shown to the end user during OOBE. Default value is FALSE.",
					},
					"eula_hidden": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, EULA is hidden to the end user during OOBE. When FALSE, EULA is shown to the end user during OOBE. Default value is FALSE.",
					},
					"user_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The type of user. Possible values are administrator and standard. Default value is administrator. Possible values are: `administrator`, `standard`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf("administrator", "standard"),
						},
					},
					"device_usage_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Entra join authentication type. Possible values are singleUser and shared. The default is singleUser. Possible values are: `singleUser`, `shared`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf("singleUser", "shared"),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"keyboard_selection_page_skipped": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "When TRUE, the keyboard selection page is hidden to the end user during OOBE if Language and Region are set. When FALSE, the keyboard selection page is skipped during OOBE.",
					},
					"escape_link_hidden": schema.BoolAttribute{
						Computed: true,
						MarkdownDescription: "When TRUE, the link that allows user to start over with a different account on company sign-in is hidden. When false, the link that allows user to start over with a different account on company sign-in is available. " +
							" This field is defaulted to TRUE for a valid api call but doesnt configure anything in the gui. This field is always required to be set to TRUE.",
					},
				},
			},
			"assignments": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The list of assignments for this deployment profile.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The type of assignment target. Possible values are: `groupAssignmentTarget`, `exclusionGroupAssignmentTarget`, `allDevicesAssignmentTarget`.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"groupAssignmentTarget",
									"exclusionGroupAssignmentTarget",
									"allDevicesAssignmentTarget",
								),
							},
						},
						"group_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The ID of the target group. Required when type is `groupAssignmentTarget` or `exclusionGroupAssignmentTarget`.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
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
