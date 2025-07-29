package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
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
	ResourceName  = "graph_beta_device_management_android_device_owner_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AndroidDeviceOwnerCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AndroidDeviceOwnerCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AndroidDeviceOwnerCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AndroidDeviceOwnerCompliancePolicyResource{}
)

func NewAndroidDeviceOwnerCompliancePolicyResource() resource.Resource {
	return &AndroidDeviceOwnerCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCompliancePolicies",
	}
}

type AndroidDeviceOwnerCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AndroidDeviceOwnerCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *AndroidDeviceOwnerCompliancePolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *AndroidDeviceOwnerCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *AndroidDeviceOwnerCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AndroidDeviceOwnerCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Android Device Owner compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` " +
			"endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security " +
			"requirements.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity. Inherited from deviceCompliancePolicy",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Admin provided name of the device configuration. Inherited from deviceCompliancePolicy",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Admin provided description of the Device Configuration. Inherited from deviceCompliancePolicy",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of Scope Tags for this Entity instance. Inherited from deviceCompliancePolicy",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			// Threat Protection settings
			"device_threat_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the policy requires devices have device threat protection enabled. When TRUE, threat protection is enabled. When FALSE, threat protection is not enabled. Default is FALSE.",
			},
			"device_threat_protection_required_security_level": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the minimum mobile threat protection risk level to that results in Intune reporting device noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet.",
				Validators: []validator.String{
					stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
				},
			},
			"advanced_threat_protection_required_security_level": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the Microsoft Defender for Endpoint (also referred to Microsoft Defender Advanced Threat Protection (MDATP)) minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet.",
				Validators: []validator.String{
					stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
				},
			},
			// Security settings
			"security_block_jailbroken_devices": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the device should not be rooted. When TRUE, if the device is detected as rooted it will be reported non-compliant. When FALSE, the device is not reported as non-compliant regardless of device rooted state. Default is FALSE.",
			},
			"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the compliance check will validate the Google Play Integrity check. When TRUE, the Google Play integrity basic check must pass to consider the device compliant. When FALSE, the Google Play integrity basic check can pass or fail and the device will be considered compliant. Default is FALSE.",
			},
			"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the compliance check will validate the Google Play Integrity check. When TRUE, the Google Play integrity device check must pass to consider the device compliant. When FALSE, the Google Play integrity device check can pass or fail and the device will be considered compliant. Default is FALSE.",
			},
			// OS version settings
			"os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum Android version required to mark the device as compliant. For example: '14'",
			},
			"os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates the maximum Android version required to mark the device as compliant. For example: '15'",
			},
			"min_android_security_patch_level": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum Android security patch level required to mark the device as compliant. Must be a valid date format (YYYY-MM-DD). Example: 2026-10-01, 2026-10-31 etc.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.DateFormatYYYYMMDDRegex), "must be a valid date in YYYY-MM-DD format"),
				},
			},
			// Password settings
			"password_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether a password is required to unlock the device. When TRUE, there must be a password set that unlocks the device for the device to be marked as compliant. When FALSE, a device is marked as compliant whether or not a password is set as required to unlock the device. Default is FALSE.",
			},
			"password_minimum_length": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum password length required to mark the device as compliant. Valid values are 4 to 16, inclusive. Valid values 4 to 16",
			},
			"password_minimum_letter_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_minimum_lower_case_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of lower case characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_minimum_non_letter_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of non-letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_minimum_numeric_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of numeric characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_minimum_symbol_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of symbol characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_minimum_upper_case_characters": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the minimum number of upper case letter characters required for device password for the device to be marked compliant. Valid values 1 to 16.",
			},
			"password_required_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the password complexity requirement for the device to be marked compliant. Possible values are: deviceDefault, required, numeric, numericComplex, alphabetic, alphanumeric, alphanumericWithSymbols, lowSecurityBiometric, customPassword.",
				Validators: []validator.String{
					stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
				},
			},
			"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the number of minutes of inactivity before a password is required.",
			},
			"password_expiration_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the number of days before the password expires. Valid values 1 to 365.",
			},
			"password_previous_password_count_to_block": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Indicates the number of previous passwords to block. Valid values 1 to 24.",
			},
			// Storage settings
			"storage_require_encryption": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether encryption on Android devices is required to mark the device as compliant.",
			},
			// Additional security settings
			"security_require_intune_app_integrity": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Intune application integrity is required to mark the device as compliant. When TRUE, Intune checks that the Intune app installed on fully managed, dedicated, or corporate-owned work profile Android Enterprise enrolled devices, is the one provided by Microsoft from the Managed Google Play store. If the check fails, the device will be reported as non-compliant. Default is FALSE.",
			},
			"require_no_pending_system_updates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the device has pending security or OS updates and sets the compliance state accordingly. When TRUE, checks if there are any pending system updates on each check in and if there are any pending security or OS version updates (System Updates), the device will be reported as non-compliant. If set to FALSE, then checks for any pending security or OS version updates (System Updates) are done without impact to device compliance state. Default is FALSE.",
			},
			"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the types of measurements and reference data used to evaluate the device SafetyNet evaluation. Evaluation is completed on the device to assess device integrity based on checks defined by Android and built into the device hardware, for example, compromised OS version or root detection. Possible values are: basic, hardwareBacked, with default value of basic.",
				Validators: []validator.String{
					stringvalidator.OneOf("basic", "hardwareBacked"),
				},
			},
			"scheduled_actions_for_rule": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The list of scheduled action for this rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Name of the scheduled action rule",
							PlanModifiers: []planmodifier.String{
								planmodifiers.DefaultValueString("unavailable"),
							},
						},
						"scheduled_action_configurations": schema.SetNestedAttribute{
							Required:            true,
							MarkdownDescription: "The list of scheduled action configurations for this compliance policy",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"action_type": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "What action to take. Possible values are: 'noAction', 'notification', 'block', 'retire', 'wipe', 'removeResourceAccessProfiles', 'pushNotification', 'remoteLock'.",
										Validators: []validator.String{
											stringvalidator.OneOf("noAction", "notification", "block", "retire", "wipe", "removeResourceAccessProfiles", "pushNotification", "remoteLock"),
										},
									},
									"grace_period_hours": schema.Int32Attribute{
										Optional:            true,
										MarkdownDescription: "Number of hours to wait till the action will be enforced",
									},
									"notification_template_id": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "What notification Message template to use",
									},
									"notification_message_cc_list": schema.ListAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										MarkdownDescription: "A list of group GUIDs to specify who to CC this notification message to",
									},
								},
							},
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.ComplianceScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
