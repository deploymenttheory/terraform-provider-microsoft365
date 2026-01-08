package graphBetaAospDeviceOwnerCompliancePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
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
	ResourceName  = "microsoft365_graph_beta_device_management_aosp_device_owner_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AospDeviceOwnerCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AospDeviceOwnerCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AospDeviceOwnerCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &AospDeviceOwnerCompliancePolicyResource{}
)

func NewAospDeviceOwnerCompliancePolicyResource() resource.Resource {
	return &AospDeviceOwnerCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCompliancePolicies",
	}
}

type AospDeviceOwnerCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *AospDeviceOwnerCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *AospDeviceOwnerCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *AospDeviceOwnerCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AospDeviceOwnerCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages AOSP (Android Open Source Project) device owner compliance policies in Microsoft Intune using the `/deviceManagement/deviceCompliancePolicies` " +
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
				MarkdownDescription: "List of Scope Tags for this Entity instance. Inherited from deviceCompliancePolicy",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			// Password settings
			"passcode_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require a password to unlock device.",
			},
			"passcode_minimum_length": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum password length. Valid values 4 to 16",
			},
			"passcode_minutes_of_inactivity_before_lock": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minutes of inactivity before a password is required. Valid values 1 to 8640",
			},
			"passcode_required_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Type of characters in password. Possible values are: deviceDefault, required, numeric, numericComplex, alphabetic, alphanumeric, alphanumericWithSymbols, lowSecurityBiometric, customPassword.",
				Validators: []validator.String{
					stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("deviceDefault"),
				},
			},
			// OS version settings
			"os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum Android version.",
			},
			"os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum Android version.",
			},
			// Android-specific settings
			"min_android_security_patch_level": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum Android security patch level.",
			},
			// Security settings
			"security_block_jailbroken_devices": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the device should not be rooted. When TRUE, if the device is detected as rooted it will be reported non-compliant. When FALSE, the device is not reported as non-compliant regardless of device rooted state. Default is FALSE.",
			},
			"storage_require_encryption": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require encryption on Android devices.",
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
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
