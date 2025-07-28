package graphBetaWindowsDeviceCompliancePolicies

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
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
	ResourceName  = "graph_beta_device_management_windows_device_compliance_policy"
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
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsDeviceCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *WindowsDeviceCompliancePolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsDeviceCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsDeviceCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *WindowsDeviceCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows device compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` " +
			"endpoint. Device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security " +
			"requirements.",
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
				MarkdownDescription: "Admin provided description of the Device Configuration",
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
			// Password settings
			"password_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require a password to unlock Windows device",
			},
			"password_block_simple": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether or not to block simple password",
			},
			"password_required_to_unlock_from_idle": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require a password to unlock an idle device",
			},
			"password_minutes_of_inactivity_before_lock": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minutes of inactivity before a password is required",
			},
			"password_expiration_days": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The password expiration in days",
			},
			"password_minimum_length": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The minimum password length",
			},
			"password_minimum_character_set_count": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The number of character sets required in the password",
			},
			"password_required_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The required password type. Possible values are: deviceDefault, alphanumeric, numeric",
				Validators: []validator.String{
					stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
				},
			},
			"password_previous_password_block_count": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The number of previous passwords to prevent re-use of",
			},
			"require_healthy_device_report": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation",
			},
			// OS version settings
			"os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum Windows version",
			},
			"os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum Windows version",
			},
			"mobile_os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum Windows Phone version",
			},
			"mobile_os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum Windows Phone version",
			},
			// Device health and security settings
			"early_launch_anti_malware_driver_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled",
			},
			"bit_locker_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled",
			},
			"secure_boot_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled",
			},
			"code_integrity_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation",
			},
			"memory_integrity_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Memory Integrity (HVCI) to be reported as healthy",
			},
			"kernel_dma_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Kernel DMA Protection to be reported as healthy",
			},
			"virtualization_based_security_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Virtualization-based Security to be reported as healthy",
			},
			"firmware_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Firmware protection to be reported as healthy",
			},
			"storage_require_encryption": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require encryption on windows devices",
			},
			"active_firewall_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require active firewall on Windows devices",
			},
			"defender_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Windows Defender Antimalware on Windows devices",
			},
			"defender_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Require Windows Defender Antimalware minimum version on Windows devices",
			},
			"signature_out_of_date": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Windows Defender Antimalware Signature to be up to date on Windows devices",
			},
			"rtp_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Windows Defender Antimalware Real-Time Protection on Windows devices",
			},
			"antivirus_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require any Antivirus solution registered with Windows Security Center to be on and monitoring",
			},
			"anti_spyware_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require any AntiSpyware solution registered with Windows Security Center to be on and monitoring",
			},
			"device_threat_protection_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require that devices have enabled device threat protection",
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
			"configuration_manager_compliance_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require to consider SCCM Compliance state into consideration for Intune Compliance State",
			},
			"tpm_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Trusted Platform Module(TPM) to be present",
			},
			"valid_operating_system_build_ranges": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The valid operating system build ranges on Windows devices",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"low_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The minimum allowed OS version",
						},
						"high_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The maximum allowed OS version",
						},
					},
				},
			},
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
								validators.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
						"minimum_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The minimum OS version for the WSL distribution",
							Validators: []validator.String{
								validators.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
						"maximum_os_version": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The maximum OS version for the WSL distribution",
							Validators: []validator.String{
								validators.MutuallyExclusiveObjectAndSet("device_compliance_policy_script", "wsl_distributions"),
							},
						},
					},
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
