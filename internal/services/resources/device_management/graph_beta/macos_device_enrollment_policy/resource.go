// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/mem/intune/enrollment/device-enrollment-program-enroll-macos

package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_macos_device_enrollment_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacOSDeviceEnrollmentPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacOSDeviceEnrollmentPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacOSDeviceEnrollmentPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacOSDeviceEnrollmentPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &MacOSDeviceEnrollmentPolicyResource{}

	_ resource.ResourceWithConfigValidators = &MacOSDeviceEnrollmentPolicyResource{}
)

func NewMacOSDeviceEnrollmentPolicyResource() resource.Resource {
	return &MacOSDeviceEnrollmentPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type MacOSDeviceEnrollmentPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MacOSDeviceEnrollmentPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *MacOSDeviceEnrollmentPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *MacOSDeviceEnrollmentPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances.
func (r *MacOSDeviceEnrollmentPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the full macOS Automated Device Enrollment (ADE) profile schema.
func (r *MacOSDeviceEnrollmentPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a macOS Automated Device Enrollment (ADE) profile using the `/deviceManagement/configurationPolicies` " +
			"settings catalog endpoint. This is the modern, settings-catalog-backed equivalent of the legacy `depMacOSEnrollmentProfile` API " +
			"(see `microsoft365_graph_beta_device_management_macos_dep_enrollment_profile`), and controls macOS Setup Assistant behavior for " +
			"devices enrolled via Apple Business Manager / Apple School Manager.",
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
				Computed:            true,
				MarkdownDescription: "Number of settings within the policy.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the policy is assigned to any scope.",
			},
			"platforms": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The platforms this policy applies to. Always `macOS`.",
			},
			"technologies": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The technology this policy is using. Always `enrollment`.",
			},
			"template_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The settings catalog template ID used by this policy.",
			},
			"template_family": schema.StringAttribute{
				Computed:            true,
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
			"requires_user_authentication": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				MarkdownDescription: "Whether the enrollment requires user authentication (user affinity). When `false`, the device enrolls " +
					"without an associated user (shared/kiosk device path).",
			},
			"enable_authentication_via_company_portal": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether Setup Assistant authenticates the user via Company Portal. Only applicable when " +
					"`requires_user_authentication` is `true`. Mutually exclusive with `require_company_portal_on_setup_assistant_enrolled_devices`.",
				Validators: []validator.Bool{
					customValidator.MutuallyExclusiveBool(
						"require_company_portal_on_setup_assistant_enrolled_devices",
						"enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices are mutually exclusive",
					),
				},
			},
			"require_company_portal_on_setup_assistant_enrolled_devices": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Whether Company Portal is required on Setup Assistant enrolled devices. Only applicable when " +
					"`requires_user_authentication` is `true`. Mutually exclusive with `enable_authentication_via_company_portal`.",
			},
			"await_device_configured": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether Setup Assistant waits for the local account configuration described by `admin_account` to complete before continuing.",
			},
			"admin_account": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "Local account settings created during Setup Assistant. Required when `await_device_configured` is `true`, " +
					"and must be omitted when it is `false`.",
				Attributes: map[string]schema.Attribute{
					"create_local_admin_account": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Whether Setup Assistant creates a local administrator account.",
					},
					"user_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The account (short) name for the local administrator account.",
					},
					"full_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The full name for the local administrator account.",
					},
					"hide_account": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						MarkdownDescription: "Whether to hide the local administrator account from the login window and Users & Groups.",
					},
					"password_rotation_in_days": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(0),
						MarkdownDescription: "Automatic rotation period, in days, for the local administrator account password. `0` disables rotation.",
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"create_local_primary_account": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
						MarkdownDescription: "Whether Setup Assistant also creates a separate, standard (non-admin) local primary account.",
					},
					"primary_account": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Standard local account settings. Only applicable when `create_local_primary_account` is `true`.",
						Attributes: map[string]schema.Attribute{
							"prefill_account_info": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to prefill the primary account name/full name in Setup Assistant.",
							},
							"restrict_editing": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to prevent the user from editing the prefilled primary account information. Only applicable when `prefill_account_info` is `true`.",
							},
							"full_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The full name to prefill for the primary account. Only applicable when `prefill_account_info` is `true`.",
							},
							"user_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The account (short) name to prefill for the primary account. Only applicable when `prefill_account_info` is `true`.",
							},
						},
					},
				},
			},
			"locked_enrollment_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether enrollment is locked to the authorized user/device, preventing the MDM profile from being removed before enrollment completes.",
			},
			"support_department": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The department name shown to the user on the Setup Assistant Remote Management pane.",
			},
			"support_phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The support phone number shown to the user on the Setup Assistant Remote Management pane.",
			},
			"location_services_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Location Services pane in Setup Assistant.",
			},
			"restore_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Restore from Backup pane in Setup Assistant.",
			},
			"apple_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Apple ID sign-in pane in Setup Assistant.",
			},
			"terms_and_conditions_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Terms and Conditions pane in Setup Assistant.",
			},
			"touch_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Touch ID/Face ID pane in Setup Assistant.",
			},
			"apple_pay_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Apple Pay pane in Setup Assistant.",
			},
			"siri_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Siri pane in Setup Assistant.",
			},
			"diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Diagnostics pane in Setup Assistant.",
			},
			"file_vault_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the FileVault pane in Setup Assistant.",
			},
			"icloud_diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the iCloud Analytics pane in Setup Assistant.",
			},
			"icloud_storage_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the iCloud Storage pane in Setup Assistant.",
			},
			"display_tone_setup_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Appearance (display tone) pane in Setup Assistant.",
			},
			"screen_time_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Screen Time pane in Setup Assistant.",
			},
			"privacy_pane_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Privacy pane in Setup Assistant.",
			},
			"accessibility_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Accessibility pane in Setup Assistant.",
			},
			"auto_unlock_with_watch_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Unlock with Apple Watch pane in Setup Assistant.",
			},
			"lockdown_mode_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Lockdown Mode pane in Setup Assistant.",
			},
			"software_update_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Software Update pane in Setup Assistant.",
			},
			"software_update_completed_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the post-installation Software Update Completed pane in Setup Assistant.",
			},
			"terms_of_address_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Terms of Address pane in Setup Assistant.",
			},
			"apple_intelligence_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the Apple Intelligence pane in Setup Assistant.",
			},
			"os_showcase_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the What's New (OS showcase) pane in Setup Assistant.",
			},
			"app_store_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether to hide the App Store pane in Setup Assistant.",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
