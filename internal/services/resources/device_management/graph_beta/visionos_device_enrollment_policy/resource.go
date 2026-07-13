// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-ios

package graphBetaVisionOSDeviceEnrollmentPolicy

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
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &VisionOSDeviceEnrollmentPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &VisionOSDeviceEnrollmentPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &VisionOSDeviceEnrollmentPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &VisionOSDeviceEnrollmentPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &VisionOSDeviceEnrollmentPolicyResource{}
)

func NewVisionOSDeviceEnrollmentPolicyResource() resource.Resource {
	return &VisionOSDeviceEnrollmentPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementServiceConfig.Read.All",
			"Directory.Read.All",
			"Group.Read.All",
			"GroupMember.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type VisionOSDeviceEnrollmentPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *VisionOSDeviceEnrollmentPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *VisionOSDeviceEnrollmentPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *VisionOSDeviceEnrollmentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances.
func (r *VisionOSDeviceEnrollmentPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
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

// Schema defines the full visionOS Automated Device Enrollment (ADE) profile schema.
func (r *VisionOSDeviceEnrollmentPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a visionOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` " +
			"settings catalog endpoint. This controls visionOS Setup Assistant behavior for Apple Vision Pro devices enrolled via " +
			"Apple Business Manager / Apple School Manager.",
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
				MarkdownDescription: "The platforms this policy applies to. Always `visionOS`.",
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
				MarkdownDescription: "Whether this policy is the default visionOS enrollment profile for its `dep_onboarding_settings_id`, " +
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
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether the enrollment requires user authentication (user affinity). visionOS Automated Device " +
					"Enrollment only supports enrollment without user affinity - Microsoft Graph rejects `true` " +
					"(`ade_useraffinitybasic_1` is not a valid option for this settings catalog template). Unlike iOS/iPadOS, this is " +
					"not user-configurable; the attribute is retained for parity with the other ADE enrollment policy resources but " +
					"must be left as `false`.",
			},
			"await_device_configured": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				MarkdownDescription: "Whether devices are locked in Setup Assistant until all enrollment-time configuration is installed " +
					"(await configuration).",
			},
			"locked_enrollment_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.",
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
			"apple_id_disabled":               setupAssistantToggleAttribute("Whether to hide the Apple ID setup pane in Setup Assistant, which gives users the option to sign in with their Apple ID and use iCloud."),
			"apple_pay_disabled":              setupAssistantToggleAttribute("Whether to hide the Apple Pay setup pane in Setup Assistant, which gives users the option to set up Apple Pay on their devices."),
			"diagnostics_disabled":            setupAssistantToggleAttribute("Whether to hide the diagnostics pane in Setup Assistant, where users can opt in to send diagnostic data to Apple."),
			"get_started_screen_disabled":     setupAssistantToggleAttribute("Whether to hide the Get Started pane in Setup Assistant."),
			"apple_intelligence_disabled":     setupAssistantToggleAttribute("Whether to hide the Apple Intelligence setup pane in Setup Assistant, where users can configure Apple Intelligence features."),
			"location_services_disabled":      setupAssistantToggleAttribute("Whether to hide the Location Services setup pane in Setup Assistant, where users can enable location services on their device."),
			"passcode_disabled":               setupAssistantToggleAttribute("Whether to hide the passcode and password lock pane in Setup Assistant. When shown, users are prompted for a passcode. Always require a passcode for unsecured devices unless access is controlled in some other way (such as through a kiosk mode configuration that restricts the device to one app)."),
			"privacy_pane_disabled":           setupAssistantToggleAttribute("Whether to hide the privacy setup pane in Setup Assistant."),
			"screen_time_screen_disabled":     setupAssistantToggleAttribute("Whether to hide the Screen Time pane in Setup Assistant."),
			"siri_disabled":                   setupAssistantToggleAttribute("Whether to hide the Siri setup pane in Setup Assistant."),
			"software_update_screen_disabled": setupAssistantToggleAttribute("Whether to hide the mandatory software update screen in Setup Assistant."),
			"terms_and_conditions_disabled":   setupAssistantToggleAttribute("Whether to hide the Apple terms and conditions pane in Setup Assistant. When shown, users are required to accept them."),
			"tips_screen_disabled":            setupAssistantToggleAttribute("Whether to hide the Tips pane in Setup Assistant."),
			"touch_id_disabled":               setupAssistantToggleAttribute("Whether to hide the biometric (Optic ID) setup pane in Setup Assistant, which gives users the option to set up biometric identification on their devices."),
			"timeouts":                        commonschema.ResourceTimeouts(ctx),
		},
	}
}
