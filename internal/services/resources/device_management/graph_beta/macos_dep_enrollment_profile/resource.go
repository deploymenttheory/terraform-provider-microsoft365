package graphBetaMacOSDepEnrollmentProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacOSDepEnrollmentProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacOSDepEnrollmentProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacOSDepEnrollmentProfileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacOSDepEnrollmentProfileResource{}

	// Enables resource-level (cross-field) configuration validation
	_ resource.ResourceWithConfigValidators = &MacOSDepEnrollmentProfileResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &MacOSDepEnrollmentProfileResource{}
)

func NewMacOSDepEnrollmentProfileResource() resource.Resource {
	return &MacOSDepEnrollmentProfileResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles",
	}
}

type MacOSDepEnrollmentProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *MacOSDepEnrollmentProfileResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *MacOSDepEnrollmentProfileResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *MacOSDepEnrollmentProfileResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *MacOSDepEnrollmentProfileResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *MacOSDepEnrollmentProfileResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a macOS Automated Device Enrollment (DEP/ADE) enrollment profile using the " +
			"`/deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles` endpoint with the " +
			"`#microsoft.graph.depMacOSEnrollmentProfile` OData type. This profile drives automated (low-touch) macOS " +
			"enrollment: skipping Setup Assistant panes, auto-creating the local admin account, and gating the desktop " +
			"until MDM configuration finishes (`await_device_configured`). Note: fully hands-off (\"zero-touch\") " +
			"provisioning is only approached when enrolling without user affinity (`requires_user_authentication = false`) " +
			"over a wired network; Apple keeps some early Setup Assistant panes (such as network and region/language) " +
			"non-skippable, so at least minimal physical interaction may still be required on a freshly wiped device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the enrollment profile. Format is `{depOnboardingSettingsId}_{profileId}`.",
			},
			"dep_onboarding_settings_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Identifier of the parent depOnboardingSetting (Apple ABM/ASM ADE token) that contains this " +
					"macOS DEP enrollment profile. If omitted, the provider resolves it from the `/deviceManagement` " +
					"endpoint's `intuneAccountId`. On tenants with multiple DEP tokens (for example, a separate Apple " +
					"Configurator token), that fallback may select the wrong token, so set this explicitly to the ABM/ADE " +
					"token id. List your tokens with `GET /deviceManagement/depOnboardingSettings` and pick the one whose " +
					"`tokenType` is the Apple ABM/ADE token (not `appleConfigurator`).",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the profile displayed in Intune.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Description of the profile. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"requires_user_authentication": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the profile requires user authentication.",
			},
			"configuration_endpoint_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Configuration endpoint url to use for enrollment. Generated by Intune.",
			},
			"enable_authentication_via_company_portal": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates to authenticate with the Company Portal instead of Apple Setup Assistant.",
				Validators: []validator.Bool{
					validate.MutuallyExclusiveBool(
						"require_company_portal_on_setup_assistant_enrolled_devices",
						"enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true",
					),
				},
			},
			"require_company_portal_on_setup_assistant_enrolled_devices": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates that the Company Portal is required on setup assistant enrolled devices.",
				Validators: []validator.Bool{
					validate.MutuallyExclusiveBool(
						"enable_authentication_via_company_portal",
						"enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true",
					),
				},
			},
			"is_default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is the default profile.",
			},
			"is_mandatory": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the profile is mandatory.",
			},
			"supervised_mode_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Supervised mode. True to enable, false otherwise.",
			},
			"support_department": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Support department information.",
			},
			"support_phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Support phone number.",
			},
			"device_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Sets a literal or name pattern for the device name.",
			},
			"profile_removal_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the profile removal option is disabled.",
			},
			"configuration_web_url": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the admin-assisted setup assistant login (web-based authentication) URL is used. Cannot be true when `use_platform_sso_during_setup_assistant` is true.",
			},
			"await_device_configured": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the device will need to wait for configured confirmation (the desktop is gated until MDM configuration finishes). Maps to `waitForDeviceConfiguredConfirmation`.",
			},
			"enabled_skip_keys": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Computed, read-only set of Setup Assistant skip keys (Apple `SkipKeys`) that the provider " +
					"sends to Graph. This is derived from the individual `*_disabled` boolean attributes; do not set it " +
					"directly. Note: `Privacy` and `Registration` are intentionally omitted from this array because the " +
					"Microsoft Graph API rejects those skip-key strings, even though the `privacy_pane_disabled` and " +
					"`registration_disabled` boolean properties work correctly.",
			},
			"enrollment_time_azure_ad_group_ids": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of enrollment-time Microsoft Entra (Azure AD) group GUIDs to be associated with the profile.",
			},
			"location_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Location service setup pane is disabled.",
			},
			"restore_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Restore setup pane is blocked.",
			},
			"apple_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Apple ID setup pane is disabled.",
			},
			"terms_and_conditions_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the 'Terms and Conditions' setup pane is disabled.",
			},
			"touch_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Touch ID setup pane is disabled.",
			},
			"apple_pay_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Apple Pay setup pane is disabled.",
			},
			"siri_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Siri setup pane is disabled.",
			},
			"diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Diagnostics setup pane is disabled.",
			},
			"display_tone_setup_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the DisplayTone setup screen is disabled.",
			},
			"privacy_pane_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Privacy screen is disabled.",
			},
			"screen_time_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Screen Time setup screen is disabled.",
			},
			"registration_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if registration is disabled.",
			},
			"welcome_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Get Started (Welcome) setup pane is disabled. macOS 15 and later.",
			},
			"file_vault_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if FileVault is disabled.",
			},
			"icloud_diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the iCloud Analytics screen is disabled.",
			},
			"pass_code_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Passcode setup pane is disabled.",
			},
			"zoom_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Zoom setup pane is disabled.",
			},
			"icloud_storage_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the iCloud Documents and Desktop screen is disabled.",
			},
			"choose_your_lock_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Choose Your Lock Screen screen is disabled.",
			},
			"accessibility_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Accessibility screen is disabled.",
			},
			"auto_unlock_with_watch_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the Unlock With Watch screen is disabled.",
			},
			"auto_advance_setup_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if Setup Assistant will automatically advance through its screens.",
			},
			"request_requires_network_tether": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the device is network-tethered to run the command.",
			},
			"use_platform_sso_during_setup_assistant": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Platform SSO is used as part of device enrollment during Setup Assistant. Cannot be true when `configuration_web_url` is true.",
			},
			"skip_primary_setup_account_creation": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will skip the user interface for primary account setup.",
			},
			"set_primary_setup_account_as_regular_user": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will set the primary account as a regular (non-admin) user.",
			},
			"dont_auto_populate_primary_account_info": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will auto-populate the primary account information.",
			},
			"primary_account_full_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The full name for the primary account.",
			},
			"primary_account_user_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The account name (short name) for the primary account.",
			},
			"enable_restrict_editing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the user will be blocked from editing the account.",
			},
			"admin_account_user_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The user name (short name) for the auto-created local admin account.",
			},
			"admin_account_full_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The full name for the auto-created local admin account.",
			},
			"admin_account_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				WriteOnly:           true,
				MarkdownDescription: "The password for the auto-created local admin account. This is a write-only, sensitive value and is never returned by the Graph API.",
			},
			"hide_admin_account": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the local admin account should be hidden.",
			},
			"admin_account_password_rotation": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for local admin account password automatic rotation (depProfileAdminAccountPasswordRotationSetting).",
				Attributes: map[string]schema.Attribute{
					"auto_rotation_period_in_days": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The number of days between automatic admin account password rotations.",
					},
					"on_retrieval_auto_rotate_password_enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Indicates whether the password is automatically rotated after retrieval.",
					},
					"on_retrieval_delay_auto_rotate_password_in_hours": schema.Int32Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The delay in hours before automatically rotating the password after retrieval.",
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
