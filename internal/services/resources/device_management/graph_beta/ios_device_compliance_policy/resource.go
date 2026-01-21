package graphBetaIosDeviceCompliancePolicy

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
	ResourceName  = "microsoft365_graph_beta_device_management_ios_device_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IosDeviceCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IosDeviceCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IosDeviceCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IosDeviceCompliancePolicyResource{}
)

func NewIosDeviceCompliancePolicyResource() resource.Resource {
	return &IosDeviceCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCompliancePolicies",
	}
}

type IosDeviceCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IosDeviceCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *IosDeviceCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *IosDeviceCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *IosDeviceCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages ios device compliance policies using the `/deviceManagement/deviceCompliancePolicies` endpoint. This resource is used to device compliance policies define rules and settings that devices must meet to be considered compliant with organizational security requirements.",
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
			// Passcode settings
			"passcode_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether or not to require a passcode",
			},
			"passcode_block_simple": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether or not to block simple passcodes",
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
				MarkdownDescription: "Number of days before the passcode expires. Valid values 1 to 65535",
			},
			"passcode_minimum_length": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum length of passcode. Valid values 4 to 14",
			},
			"passcode_minimum_character_set_count": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "The number of character sets required in the password",
			},
			"passcode_required_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The required passcode type. Possible values are: deviceDefault, alphanumeric, numeric",
				Validators: []validator.String{
					stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("deviceDefault"),
				},
			},
			"passcode_previous_passcode_block_count": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Number of previous passwords to block",
			},
			// OS version settings
			"os_minimum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum iOS version allowed.",
			},
			"os_maximum_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum iOS version allowed.",
			},
			"os_minimum_build_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Minimum iOS build version",
			},
			"os_maximum_build_version": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Maximum iOS build version",
			},
			// iOS security settings
			"managed_email_profile_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether or not to require a managed email profile",
			},
			"security_block_jailbroken_devices": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates the device should not be jailbroken. When TRUE, if the device is detected as jailbroken it will be reported non-compliant",
			},
			// Restricted apps
			"restricted_apps": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The application name",
						},
						"publisher": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The publisher of the application",
						},
						"app_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The application or bundle identifier of the application",
						},
						"app_store_url": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The Store URL of the application",
						},
					},
				},
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
			"advanced_threat_protection_required_security_level": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Require Microsoft Defender for Endpoint minimum risk level to report noncompliance. Possible values are: unavailable, secured, low, medium, high, notSet",
				Validators: []validator.String{
					stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("unavailable"),
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
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
