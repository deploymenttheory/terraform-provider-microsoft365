package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
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
	ResourceName  = "graph_beta_device_management_windows_autopilot_deployment_profile"
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
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsAutopilotDeploymentProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsAutopilotDeploymentProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsAutopilotDeploymentProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *WindowsAutopilotDeploymentProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Windows Autopilot Deployment Profile in Microsoft Intune.",
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
				MarkdownDescription: "The display name of the deployment profile. Max allowed length is 200 chars.",
				Validators: []validator.String{
					validators.StringLengthAtMost(200),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the windows autopilotdeployment profile. Max allowed length is 1500 chars.",
				Validators: []validator.String{
					validators.StringLengthAtMost(1500),
				},
			},
			"language": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The language code to be used when configuring the device. E.g. en-US. The default value is os-default. Read-Only. Starting from May 2024 this property will no longer be supported and will be marked as deprecated. Use locale instead.",
			},
			"locale": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The locale (language) to be used when configuring the device. E.g. en-US. The default value is os-default.",
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
				MarkdownDescription: "The type of device join to configure. Determines which Windows Autopilot deployment profile type to use. Possible values are: `microsoft_entra_joined`, `microsoft_entra_hybrid_joined`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"microsoft_entra_joined",
						"microsoft_entra_hybrid_joined",
					),
				},
			},
			"hardware_hash_extraction_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the profile supports the extraction of hardware hash values and registration of the device into Windows Autopilot. When TRUE, indicates if hardware extraction and Windows Autopilot registration will happen on the next successful check-in. When FALSE, hardware hash extraction and Windows Autopilot registration will not happen. Default value is FALSE.",
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.BoolDefaultValue(false),
				},
			},
			"device_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The template used to name the Autopilot device. This can be a custom text and can also contain either the serial number of the device, or a randomly generated number. The total length of the text generated by the template can be no more than 15 characters.",
				Validators: []validator.String{
					validators.StringLengthAtMost(15),
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
						"surfaceHub2",
						"surfaceHub2S",
						"virtualMachine",
						"unknownFutureValue",
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
			},
			"hybrid_azure_ad_join_skip_connectivity_check": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "The Autopilot Hybrid Azure AD join flow will continue even if it does not establish domain controller connectivity during OOBE. This is only applicable for `microsoft_entra_hybrid_joined` device join type.",
			},
			"out_of_box_experience_setting": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The Windows Autopilot Deployment Profile settings used by the device for the out-of-box experience.",
				Attributes: map[string]schema.Attribute{
					"privacy_settings_hidden": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, privacy settings is hidden to the end user during OOBE. When FALSE, privacy settings is shown to the end user during OOBE. Default value is FALSE.",
					},
					"eula_hidden": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, EULA is hidden to the end user during OOBE. When FALSE, EULA is shown to the end user during OOBE. Default value is FALSE.",
					},
					"user_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The type of user. Possible values are administrator and standard. Default value is administrator. Possible values are: `administrator`, `standard`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf("administrator", "standard", "unknownFutureValue"),
						},
					},
					"device_usage_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Entra join authentication type. Possible values are singleUser and shared. The default is singleUser. Possible values are: `singleUser`, `shared`, `unknownFutureValue`.",
						Validators: []validator.String{
							stringvalidator.OneOf("singleUser", "shared", "unknownFutureValue"),
						},
					},
					"keyboard_selection_page_skipped": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, the keyboard selection page is hidden to the end user during OOBE if Language and Region are set. When FALSE, the keyboard selection page is skipped during OOBE.",
					},
					"escape_link_hidden": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE, the link that allows user to start over with a different account on company sign-in is hidden. When false, the link that allows user to start over with a different account on company sign-in is available. Default value is FALSE.",
					},
				},
			},
			"enrollment_status_screen_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The Windows Enrollment Status Screen settings for the deployment profile.",
				Attributes: map[string]schema.Attribute{
					"hide_installation_progress": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Show or hide installation progress to user.",
					},
					"allow_device_use_before_profile_and_app_install_complete": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow or block user to use device before profile and app installation complete.",
					},
					"block_device_setup_retry_by_user": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow the user to retry the setup on installation failure.",
					},
					"allow_log_collection_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow or block log collection on installation failure.",
					},
					"custom_error_message": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Set custom error message to show upon installation failure.",
					},
					"install_progress_timeout_in_minutes": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Set installation progress timeout in minutes.",
					},
					"allow_device_use_on_install_failure": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow the user to continue using the device on installation failure.",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
