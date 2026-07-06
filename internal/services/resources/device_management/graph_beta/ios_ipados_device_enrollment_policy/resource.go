// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios

package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IOSiPadOSDeviceEnrollmentPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IOSiPadOSDeviceEnrollmentPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IOSiPadOSDeviceEnrollmentPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IOSiPadOSDeviceEnrollmentPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &IOSiPadOSDeviceEnrollmentPolicyResource{}

	_ resource.ResourceWithConfigValidators = &IOSiPadOSDeviceEnrollmentPolicyResource{}
)

func NewIOSiPadOSDeviceEnrollmentPolicyResource() resource.Resource {
	return &IOSiPadOSDeviceEnrollmentPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementServiceConfig.Read.All",
			"Directory.Read.All",
			"Group.Read.All",
			"GroupMember.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type IOSiPadOSDeviceEnrollmentPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// setupAssistantToggleAttribute returns the shared shape of one Setup Assistant screen toggle.
func setupAssistantToggleAttribute(description string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: description,
	}
}

// Schema defines the full iOS/iPadOS Automated Device Enrollment (ADE) profile schema.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iOS/iPadOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` " +
			"settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy `depIOSEnrollmentProfile` API, and " +
			"controls iOS/iPadOS Setup Assistant behavior for devices enrolled via Apple Business Manager / Apple School Manager.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The unique identifier for this policy.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the enrollment profile.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Creation date and time of the policy.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the policy.",
			},
			"settings_count": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Number of settings within the policy.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Indicates if the policy is assigned to any scope.",
			},
			"platforms": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The platforms this policy applies to. Always `iOS`.",
			},
			"technologies": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The technology this policy is using. Always `enrollment`.",
			},
			"template_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The settings catalog template ID used by this policy.",
			},
			"template_family": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The template family for this policy (`enrollmentConfiguration`).",
			},
			"dep_onboarding_settings_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the Apple ABM/ASM DEP onboarding token (`/deviceManagement/depOnboardingSettings`) that owns " +
					"this profile. If omitted, it is automatically resolved to the tenant's single Apple ADE/ABM (or ASM) token; if the tenant " +
					"has more than one Apple token, this must be set explicitly.",
			},
			"is_default_policy_assignment": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Whether this policy is the default iOS/iPadOS enrollment profile for its `dep_onboarding_settings_id`, " +
					"set via the dedicated `setDefaultProfile` action. Always reflects the DEP token's actual current default on " +
					"refresh, regardless of configuration.\n\n" +
					"~> **No unassign action:** Microsoft Graph does not expose an `unsetDefaultProfile`/`clearDefaultProfile` action - " +
					"`setDefaultProfile` is the only operation available. Setting this to `false` on a policy that is currently the " +
					"DEP token's default has no effect on Graph; the next refresh reports `true` again. Only setting a different " +
					"policy's `is_default_policy_assignment` to `true` changes which profile is the default. A change from `true` to " +
					"`false` while this policy is still the token's current default can therefore never converge, and the provider " +
					"rejects the update with a validation error. Promote the replacement policy first - in the same apply, give this " +
					"policy a `depends_on` for the replacement so the promotion runs first, or apply the promotion separately.",
			},
			"device_security_group": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The ID of the static Microsoft Entra security group to use for enrollment time grouping (the " +
					"\"Device group\" tab in the Intune admin center policy wizard). Devices assigned this policy become members of the " +
					"group as they enroll. This is set via the dedicated `setEnrollmentTimeDeviceMembershipTarget` action on " +
					"`/deviceManagement/configurationPolicies/{id}` (and cleared via `clearEnrollmentTimeDeviceMembershipTarget` when " +
					"removed), not via the settings catalog. The group must have the 'Intune Provisioning Client' service principal " +
					"(AppId: f1346770-5b25-470b-88bd-d5744ab7952c) set as its owner; in some tenants this service principal may appear as " +
					"'Intune Autopilot ConfidentialClient'.\n\n" +
					"~> **Known Microsoft Graph limitation:** as of this writing, `setEnrollmentTimeDeviceMembershipTarget` and " +
					"`clearEnrollmentTimeDeviceMembershipTarget` return an `Internal Server Error - 500` from the Intune backend " +
					"(`DeviceConfigV2`) when called with application permissions (client credentials) - the auth flow this provider " +
					"always uses. The identical request succeeds when made with delegated (signed-in user) permissions, e.g. from the " +
					"Intune admin center. Until Microsoft resolves this for application permissions, setting `device_security_group` " +
					"through this provider will fail on `Create` and `Update`; this is a Microsoft Graph service limitation, not a " +
					"provider defect.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"requires_user_authentication": schema.BoolAttribute{
				Required: true,
				MarkdownDescription: "Whether the enrollment requires user authentication (user affinity). When `false`, the device enrolls " +
					"without an associated user (shared/kiosk device path). When `true`, the authentication flow is selected via " +
					"`enable_authentication_via_company_portal` or `require_setup_assistant_with_modern_authentication`; when neither is " +
					"set, the legacy Setup Assistant authentication flow is used.",
			},
			"enable_authentication_via_company_portal": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether the user authenticates via the Company Portal app instead of Setup Assistant. Only applicable " +
					"when `requires_user_authentication` is `true`. Mutually exclusive with " +
					"`require_setup_assistant_with_modern_authentication`.",
				Validators: []validator.Bool{
					customValidator.MutuallyExclusiveBool(
						"require_setup_assistant_with_modern_authentication",
						"enable_authentication_via_company_portal and require_setup_assistant_with_modern_authentication are mutually exclusive",
					),
				},
			},
			"require_setup_assistant_with_modern_authentication": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether the user authenticates in Setup Assistant using modern authentication (Microsoft Entra ID). " +
					"Only applicable when `requires_user_authentication` is `true`. Mutually exclusive with " +
					"`enable_authentication_via_company_portal`.",
			},
			"await_final_configuration": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether devices are locked in Setup Assistant until all enrollment-time configuration is installed " +
					"(await final configuration). Only applicable when `require_setup_assistant_with_modern_authentication` is `true`.",
			},
			"locked_enrollment_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.",
			},
			"device_name_template": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Device name template applied to supervised devices at enrollment, e.g. `{{DEVICETYPE}}-{{SERIAL}}`. " +
					"Supports the `{{DEVICETYPE}}` and `{{SERIAL}}` substitution tokens. When omitted, device naming is not managed by " +
					"this profile.",
			},
			"cellular_data_activation_url": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "The carrier activation server URL used to activate cellular data plans on eligible devices at " +
					"enrollment. When omitted, cellular data plan activation is not configured by this profile.",
			},
			"support_department": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The department name shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 125 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 125),
				},
			},
			"support_phone_number": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The support phone number shown to the user on the Setup Assistant Remote Management pane. Must be between 1 and 50 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"passcode_disabled":                         setupAssistantToggleAttribute("Whether to hide the passcode and password lock pane in Setup Assistant. When shown, users are prompted for a passcode. Always require a passcode for unsecured devices unless access is controlled in some other way (such as through a kiosk mode configuration that restricts the device to one app). For iOS/iPadOS 7.0 and later."),
			"location_services_disabled":                setupAssistantToggleAttribute("Whether to hide the Location Services setup pane in Setup Assistant, where users can enable location services on their device. For iOS/iPadOS 7.0 and later."),
			"restore_disabled":                          setupAssistantToggleAttribute("Whether to hide the apps and data (Restore) setup pane in Setup Assistant. When shown, users setting up devices can restore or transfer data from iCloud Backup. For iOS/iPadOS 7.0 and later."),
			"apple_id_disabled":                         setupAssistantToggleAttribute("Whether to hide the Apple ID setup pane in Setup Assistant, which gives users the option to sign in with their Apple ID and use iCloud. For iOS/iPadOS 7.0 and later."),
			"terms_and_conditions_disabled":             setupAssistantToggleAttribute("Whether to hide the Apple terms and conditions pane in Setup Assistant. When shown, users are required to accept them. For iOS/iPadOS 7.0 and later."),
			"touch_id_disabled":                         setupAssistantToggleAttribute("Whether to hide the biometric (Touch ID and Face ID) setup pane in Setup Assistant, which gives users the option to set up fingerprint or facial identification on their devices. For iOS/iPadOS 8.1 and later."),
			"apple_pay_disabled":                        setupAssistantToggleAttribute("Whether to hide the Apple Pay setup pane in Setup Assistant, which gives users the option to set up Apple Pay on their devices. For iOS/iPadOS 7.0 and later."),
			"siri_disabled":                             setupAssistantToggleAttribute("Whether to hide the Siri setup pane in Setup Assistant. For iOS/iPadOS 7.0 and later."),
			"diagnostics_disabled":                      setupAssistantToggleAttribute("Whether to hide the diagnostics pane in Setup Assistant, where users can opt in to send diagnostic data to Apple. For iOS/iPadOS 7.0 and later."),
			"privacy_pane_disabled":                     setupAssistantToggleAttribute("Whether to hide the privacy setup pane in Setup Assistant. For iOS/iPadOS 11.3 and later."),
			"restore_from_android_disabled":             setupAssistantToggleAttribute("Whether to hide the Android Migration setup pane in Setup Assistant, meant for previous Android users. When shown, users can migrate data from an Android device. For iOS/iPadOS 9.0 and later."),
			"imessage_and_facetime_disabled":            setupAssistantToggleAttribute("Whether to hide the iMessage and FaceTime setup pane in Setup Assistant. For iOS/iPadOS 9.0 and later."),
			"screen_time_screen_disabled":               setupAssistantToggleAttribute("Whether to hide the Screen Time pane in Setup Assistant. For iOS/iPadOS 12.0 and later."),
			"sim_setup_screen_disabled":                 setupAssistantToggleAttribute("Whether to hide the cellular (SIM Setup) pane in Setup Assistant, where users can add a cellular plan. For iOS/iPadOS 12.0 and later."),
			"software_update_screen_disabled":           setupAssistantToggleAttribute("Whether to hide the mandatory software update screen in Setup Assistant. For iOS/iPadOS 12.0 and later."),
			"watch_migration_screen_disabled":           setupAssistantToggleAttribute("Whether to hide the Apple Watch migration pane in Setup Assistant, where users can migrate data from an Apple Watch. For iOS/iPadOS 11.0 and later."),
			"appearance_screen_disabled":                setupAssistantToggleAttribute("Whether to hide the appearance setup pane in Setup Assistant. For iOS/iPadOS 13.0 and later."),
			"device_to_device_migration_disabled":       setupAssistantToggleAttribute("Whether to hide the device-to-device migration pane in Setup Assistant. When shown, users can transfer data from an old device to their current device. The option to transfer data directly from a device isn't available for devices running iOS 13 or later."),
			"restore_completed_screen_disabled":         setupAssistantToggleAttribute("Whether to hide the Restore Completed screen shown after a backup and restore is performed during Setup Assistant."),
			"software_update_completed_screen_disabled": setupAssistantToggleAttribute("Whether to hide the screen showing all software updates that happen during Setup Assistant."),
			"get_started_screen_disabled":               setupAssistantToggleAttribute("Whether to hide the Get Started pane in Setup Assistant."),
			"action_button_screen_disabled":             setupAssistantToggleAttribute("Whether to hide the configuration pane for the action button in Setup Assistant. For iOS/iPadOS 17.0 and later."),
			"safety_screen_disabled":                    setupAssistantToggleAttribute("Whether to hide the safety (Emergency SOS) setup pane in Setup Assistant. For iOS/iPadOS 16.0 and later."),
			"terms_of_address_screen_disabled":          setupAssistantToggleAttribute("Whether to hide the terms of address pane in Setup Assistant, which gives users the option to choose how they want to be addressed throughout the system: feminine, masculine, or neutral. This Apple feature is available for select languages. For iOS/iPadOS 16.0 and later."),
			"apple_intelligence_disabled":               setupAssistantToggleAttribute("Whether to hide the Apple Intelligence setup pane in Setup Assistant, where users can configure Apple Intelligence features. For iOS/iPadOS 18.0 and later."),
			"lockdown_mode_disabled":                    setupAssistantToggleAttribute("Whether to hide the Lockdown Mode pane in Setup Assistant."),
			"app_store_disabled":                        setupAssistantToggleAttribute("Whether to hide the Apple App Store pane in Setup Assistant. For iOS/iPadOS 14.3 and later."),
			"camera_button_screen_disabled":             setupAssistantToggleAttribute("Whether to hide the camera button pane in Setup Assistant. For iOS/iPadOS 18.0 and later."),
			"multitasking_screen_disabled":              setupAssistantToggleAttribute("Whether to hide the multitasking pane in Setup Assistant. For iOS/iPadOS 26.0 and later."),
			"os_showcase_screen_disabled":               setupAssistantToggleAttribute("Whether to hide the OS showcase pane in Setup Assistant. For iOS/iPadOS 26.0 and later."),
			"safety_and_handling_screen_disabled":       setupAssistantToggleAttribute("Whether to hide the safety and handling pane in Setup Assistant. For iOS/iPadOS 18.4 and later."),
			"web_content_filtering_disabled":            setupAssistantToggleAttribute("Whether to hide the web content filtering pane in Setup Assistant. For iOS/iPadOS 18.2 and later."),
			"timeouts":                                  commonschema.ResourceTimeouts(ctx),
		},
	}
}
