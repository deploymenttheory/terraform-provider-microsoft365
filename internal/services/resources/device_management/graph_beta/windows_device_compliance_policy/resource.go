package graphBetaWindowsDeviceCompliancePolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_device_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsDeviceCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsDeviceCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsDeviceCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsDeviceCompliancePolicyResource{}
)

func NewWindowsDeviceCompliancePolicyResource() resource.Resource {
	return &WindowsDeviceCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCompliancePolicies",
	}
}

type WindowsDeviceCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsDeviceCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsDeviceCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsDeviceCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WindowsDeviceCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows device compliance policies using the `/deviceManagement/deviceCompliancePolicies` endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements. you can find out more here: 'https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-policy-create-windows'..",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The id of the driver.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the device compliance policy",
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
				MarkdownDescription: "Set of scope tag IDs for this Entity instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			// all of these fields are now deprecated.

			// "require_healthy_device_report": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation",
			// 	Default:             booldefault.StaticBool(false),
			// },
			// // Device health and security settings
			// "early_launch_anti_malware_driver_enabled": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled",
			// 	Default:             booldefault.StaticBool(false),
			// },
			// "memory_integrity_enabled": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require Memory Integrity (HVCI) to be reported as healthy",
			// 	Default:             booldefault.StaticBool(false),
			// },
			// "kernel_dma_protection_enabled": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require Kernel DMA Protection to be reported as healthy",
			// 	Default:             booldefault.StaticBool(false),
			// },
			// "virtualization_based_security_enabled": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require Virtualization-based Security to be reported as healthy",
			// 	Default:             booldefault.StaticBool(false),
			// },
			// "firmware_protection_enabled": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Require Firmware protection to be reported as healthy",
			// 	Default:             booldefault.StaticBool(false),
			// },
			"device_properties": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Device operating system version requirements and build ranges for compliance evaluation",
				Attributes: map[string]schema.Attribute{
					"os_minimum_version": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "Minimum OS version. Enter the minimum allowed version in the major.minor.build.revision number format. " +
							"To get the correct value, open a command prompt, and type ver. The ver command returns the version in the following format: " +
							"Microsoft Windows [Version 10.0.17134.1] When a device has an earlier version than the OS version you enter, it's reported as noncompliant. " +
							"A link with information on how to upgrade is shown. The end user can choose to upgrade their device. After they upgrade, they can access company resources.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.OSVersionRegex),
								"must be a valid OS version format (e.g., 10.0.22631.9999)",
							),
						},
					},
					"os_maximum_version": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "Maximum OS version:Enter the maximum allowed version, in the major.minor.build.revision number format. " +
							"To get the correct value, open a command prompt, and type ver. The ver command returns the version in the following format: " +
							"Microsoft Windows [Version 10.0.17134.1] When a device is using an OS version later than the version entered, access to organization " +
							"resources is blocked. The end user is asked to contact their IT administrator. The device can't access organization resources until " +
							"the rule is changed to allow the OS version.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.OSVersionRegex),
								"must be a valid OS version format (e.g., 10.0.22631.9999)",
							),
						},
					},
					"mobile_os_minimum_version": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "Enter the minimum allowed version, in the major.minor.build number format. When a device has an earlier version" +
							" that the OS version you enter, it's reported as noncompliant. A link with information on how to upgrade is shown. The end user can " +
							"choose to upgrade their device. After they upgrade, they can access company resources.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.OSVersionRegex),
								"must be a valid OS version format (e.g., 10.0.22631.9999)",
							),
						},
					},
					"mobile_os_maximum_version": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "Enter the maximum allowed version, in the major.minor.build number. When a device is using an OS version " +
							"later than the version entered, access to organization resources is blocked. The end user is asked to contact their IT administrator. " +
							"The device can't access organization resources until the rule is changed to allow the OS version.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.OSVersionRegex),
								"must be a valid OS version format (e.g., 10.0.22631.9999)",
							),
						},
					},
					"valid_operating_system_build_ranges": schema.SetNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The valid operating system build ranges on Windows devices",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Description for this valid operating system build range",
								},
								"low_os_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The minimum allowed OS version for this build range",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.OSVersionRegex),
											"must be a valid OS version format (e.g., 10.0.22631.9999)",
										),
									},
								},
								"high_os_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The maximum allowed OS version for this build range",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(constants.OSVersionRegex),
											"must be a valid OS version format (e.g., 10.0.22631.9999)",
										),
									},
								},
							},
						},
					},
				},
			},
			"system_security": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "System security settings for device compliance including firewall, antivirus, TPM, and encryption requirements",
				Attributes: map[string]schema.Attribute{
					"active_firewall_required": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require active firewall on Windows devices",
						Default:             booldefault.StaticBool(false),
					},
					"anti_spyware_required": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require any AntiSpyware solution registered with Windows Security Center to be on and monitoring",
						Default:             booldefault.StaticBool(false),
					},
					"antivirus_required": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require any Antivirus solution registered with Windows Security Center to be on and monitoring",
						Default:             booldefault.StaticBool(false),
					},
					"configuration_manager_compliance_required": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Require device compliance from Configuration Manager: " +
							"Not configured (default) - Intune doesn't check for any of the Configuration Manager settings for compliance. " +
							"Require - Require all settings (configuration items) in Configuration Manager to be compliant.",
						Default: booldefault.StaticBool(false),
					},
					"defender_enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware on Windows devices",
						Default:             booldefault.StaticBool(false),
					},
					"defender_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware minimum version on Windows devices",
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.OSVersionRegex),
								"must be a valid version format (e.g., 4.11.0.0)",
							),
						},
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block simple password",
						Default:             booldefault.StaticBool(false),
					},
					"password_minimum_character_set_count": schema.Int32Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require a password to unlock Windows device",
						Default:             booldefault.StaticBool(false),
					},
					"password_required_to_unlock_from_idle": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require a password to unlock an idle device",
						Default:             booldefault.StaticBool(false),
					},
					"password_required_type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The required password type. Possible values are: deviceDefault, alphanumeric, numeric",
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						Default: stringdefault.StaticString("deviceDefault"),
					},
					"rtp_enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware Real-Time Protection on Windows devices",
						Default:             booldefault.StaticBool(false),
					},
					"signature_out_of_date": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware Signature to be up to date on Windows devices",
						Default:             booldefault.StaticBool(false),
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require encryption on windows devices",
						Default:             booldefault.StaticBool(false),
					},
					"tpm_required": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require Trusted Platform Module(TPM) to be present",
						Default:             booldefault.StaticBool(false),
					},
				},
			},
			"microsoft_defender_for_endpoint": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Microsoft Defender for Endpoint device threat protection settings",
				Attributes: map[string]schema.Attribute{
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection",
						Default:             booldefault.StaticBool(false),
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Require Device Threat Protection minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet",
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers: []planmodifier.String{
							planmodifiers.DefaultValueString("unavailable"),
						},
					},
				},
			},
			// Device health settings (moved to device_health block)
			"device_health": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "Microsoft Attestation Service evaluation settings. Use these settings to confirm that a device has protective measures enabled at boot time." +
					"Learn more here 'https://learn.microsoft.com/en-us/intune/intune-service/protect/compliance-policy-create-windows?WT.mc_id=Portal-Microsoft_Intune_DeviceSettings#device-health'",
				Attributes: map[string]schema.Attribute{
					"bit_locker_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Windows BitLocker Drive Encryption encrypts all data stored on the Windows operating system volume. " +
							"BitLocker uses the Trusted Platform Module (TPM) to help protect the Windows operating system and user data. " +
							"It also helps confirm that a computer isn't tampered with, even if its left unattended, lost, or stolen. If the computer " +
							"is equipped with a compatible TPM, BitLocker uses the TPM to lock the encryption keys that protect the data. As a result, " +
							"the keys can't be accessed until the TPM verifies the state of the computer. " +
							"Not configured (default) - This setting isn't evaluated for compliance or non-compliance. " +
							"Require - The device can protect data that's stored on the drive from unauthorized access when the system is off, or hibernates.",
						Default: booldefault.StaticBool(false),
					},
					"secure_boot_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Require Secure Boot to be enabled on the device:" +
							"Not configured (default) - This setting isn't evaluated for compliance or non-compliance. " +
							"Require - The system is forced to boot to a factory trusted state. The core components that " +
							"are used to boot the machine must have correct cryptographic signatures that are trusted by " +
							"the organization that manufactured the device. The UEFI firmware verifies the signature before " +
							"it lets the machine start. If any files are tampered with, which breaks their signature, the system doesn't boot.",
						Default: booldefault.StaticBool(false),
					},
					"code_integrity_enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						MarkdownDescription: "Require code integrity: " +
							"Code integrity is a feature that validates the integrity of a driver or system file each time it's loaded into memory." +
							"Not configured (default) - This setting isn't evaluated for compliance or non-compliance." +
							"Require - Require code integrity, which detects if an unsigned driver or system file is being loaded into the kernel. " +
							"It also detects if a system file is changed by malicious software or run by a user account with administrator privileges.",
						Default: booldefault.StaticBool(false),
					},
				},
			},
			// Security settings (moved to system_security block)
			"device_compliance_policy_script": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Device compliance policy script for custom compliance. When wsl block is set, this block is computed and should not be set.",
				Attributes: map[string]schema.Attribute{
					"device_compliance_script_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The ID of the device compliance script",
					},
					"rules_content": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The base64 encoded rules content of the compliance script",
					},
				},
			},
			"custom_compliance_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether custom compliance is required",
				// does not require a default value of false.
			},
			"wsl_distributions": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Windows Subsystem for Linux distributions configuration",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"distribution": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the WSL distribution",
							Validators: []validator.String{
								validate.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
						"minimum_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The minimum OS version for the WSL distribution",
							Validators: []validator.String{
								validate.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
						"maximum_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The maximum OS version for the WSL distribution",
							Validators: []validator.String{
								validate.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
					},
				},
			},
			"scheduled_actions_for_rule": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "The list of scheduled action for this rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
										Required:            true,
										MarkdownDescription: "Number of hours to wait till the action will be enforced. Value must be between 0 and 365",
										Validators: []validator.Int32{
											int32validator.Between(0, 365),
										},
									},
									"notification_template_id": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "What notification Message template to use",
									},
									"notification_message_cc_list": schema.ListAttribute{
										ElementType:         types.StringType,
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "A list of group GUIDs to specify who to CC this notification message to",
									},
								},
							},
						},
					},
				},
			},
			"assignments": commonschemagraphbeta.ComplianceScriptAssignmentsSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
