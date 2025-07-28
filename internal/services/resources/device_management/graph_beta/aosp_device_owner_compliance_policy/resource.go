package graphBetaDeviceCompliancePolicies

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_device_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceCompliancePolicyResource{}
)

func NewDeviceCompliancePolicyResource() resource.Resource {
	return &DeviceCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCompliancePolicies",
	}
}

type DeviceCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *DeviceCompliancePolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *DeviceCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages device compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements. This resource supports four policy types: AOSP Device Owner, Android Device Owner, iOS, and Windows 10 compliance policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this device compliance policy",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The OData type of the compliance policy. Supported values: 'aospDeviceOwnerCompliancePolicy', 'androidDeviceOwnerCompliancePolicy', 'iosCompliancePolicy', 'windows10CompliancePolicy', 'macOSCompliancePolicy'",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"aospDeviceOwnerCompliancePolicy",
						"androidDeviceOwnerCompliancePolicy",
						"iosCompliancePolicy",
						"windows10CompliancePolicy",
						"macOSCompliancePolicy",
					),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the device compliance policy",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the device compliance policy",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of role scope tag IDs for this Entity instance",
			},
			"os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum OS version required for device compliance",
			},
			"os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum OS version allowed for device compliance",
			},
			"password_required": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether password is required on the device",
			},
			"password_required_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The type of password required",
			},
			"scheduled_actions_for_rule": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The list of scheduled action for this rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of the rule",
						},
						"scheduled_action_configurations": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The list of scheduled action configurations for this compliance policy",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"action_type": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "What action to take",
									},
									"grace_period_hours": schema.Int32Attribute{
										Optional:            true,
										MarkdownDescription: "Number of hours to wait till the action will be enforced",
									},
									"notification_template_id": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "What notification Message template to use",
									},
									"notification_message_cc_list": schema.SetAttribute{
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
			"local_actions": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The list of local actions available on this device",
			},
			"aosp_device_owner_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for AOSP Device Owner compliance",
				Attributes: map[string]schema.Attribute{
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on Android device",
					},
					"password_minimum_length": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required",
					},
				},
			},
			"android_device_owner_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for Android Device Owner compliance",
				Attributes: map[string]schema.Attribute{
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on Android device",
					},
					"password_minimum_length": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require device threat protection minimum risk level",
					},
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require mobile threat protection minimum risk level",
					},
					"password_expiration_days": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires",
					},
					"password_previous_password_count_to_block": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block",
					},
					"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require a specific SafetyNet evaluation type for Android devices",
					},
					"security_require_intune_app_integrity": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require the device to pass the Company Portal app runtime check",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection",
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require SafetyNet attestation basic integrity",
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require SafetyNet attestation certified device",
					},
				},
			},
			"ios_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for iOS compliance",
				Attributes: map[string]schema.Attribute{
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require device threat protection minimum risk level",
					},
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require mobile threat protection minimum risk level",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection",
					},
					"passcode_required_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The required passcode type",
					},
					"managed_email_profile_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to require a managed email profile",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted",
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum iOS build version",
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum iOS build version",
					},
					"passcode_minimum_character_set_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password",
					},
					"passcode_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a passcode is required",
					},
					"passcode_minutes_of_inactivity_before_screen_timeout": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out",
					},
					"passcode_expiration_days": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the passcode expires",
					},
					"passcode_previous_passcode_block_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passcodes to block",
					},
					"restricted_apps": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Require the device to not have the specified apps installed",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The display name of the app",
								},
								"app_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The application ID of the app",
								},
							},
						},
					},
				},
			},
			"macos_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for macOS compliance",
				Attributes: map[string]schema.Attribute{
					"gatekeeper_allowed_app_source": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "System and Privacy setting that determines which download locations apps can be run from on a macOS device",
					},
					"system_integrity_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require that devices have enabled system integrity protection",
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum macOS build version",
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum macOS build version",
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block simple passwords",
					},
					"password_minimum_character_set_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on Mac OS devices",
					},
					"firewall_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether the firewall should be enabled or not",
					},
					"firewall_block_all_incoming": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Corresponds to the 'Block all incoming connections' option",
					},
					"firewall_enable_stealth_mode": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Corresponds to 'Enable stealth mode.'",
					},
				},
			},
			"windows10_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for Windows 10 compliance",
				Attributes: map[string]schema.Attribute{
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require device threat protection minimum risk level",
					},
					"device_compliance_policy_script": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The PowerShell script to use for device compliance policy",
					},
					"password_required_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The required password type",
					},
					"wsl_distributions": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Windows Subsystem for Linux (WSL) distributions that should be checked for compliance",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"distribution": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The WSL distribution name",
								},
								"minimum_os_version": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The minimum OS version of the WSL distribution",
								},
								"maximum_os_version": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The maximum OS version of the WSL distribution",
								},
							},
						},
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require a password to unlock Windows device",
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block simple passwords",
					},
					"password_required_to_unlock_from_idle": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require a password to unlock an idle device",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on windows devices",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required",
					},
					"password_minimum_character_set_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password",
					},
					"active_firewall_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require Active firewall on Windows devices",
					},
					"tpm_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require Trusted Platform Module (TPM) to be present",
					},
					"antivirus_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require any antivirus solution registered with Windows Decurity Center",
					},
					"anti_spyware_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require any antispyware solution registered with Windows Decurity Center",
					},
					"defender_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware on Windows devices",
					},
					"signature_out_of_date": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware signature to be up to date on Windows devices",
					},
					"rtp_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware Real-Time Protection on Windows devices",
					},
					"defender_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware minimum version on Windows devices",
					},
					"configuration_manager_compliance_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require devices to be reported healthy by System Center Configuration Manager",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Windows 10 version",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Windows 10 version",
					},
					"mobile_os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Windows Phone version",
					},
					"mobile_os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Windows Phone version",
					},
					"secure_boot_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled",
					},
					"bit_locker_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled",
					},
					"code_integrity_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - code integrity is enabled",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
