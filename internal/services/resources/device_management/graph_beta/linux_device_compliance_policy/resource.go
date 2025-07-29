package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"

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
	ResourceName  = "graph_beta_device_management_linux_device_compliance_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LinuxDeviceCompliancePolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LinuxDeviceCompliancePolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LinuxDeviceCompliancePolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LinuxDeviceCompliancePolicyResource{}
)

func NewLinuxDeviceCompliancePolicyResource() resource.Resource {
	return &LinuxDeviceCompliancePolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/compliancePolicies",
	}
}

type LinuxDeviceCompliancePolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *LinuxDeviceCompliancePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *LinuxDeviceCompliancePolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *LinuxDeviceCompliancePolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *LinuxDeviceCompliancePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *LinuxDeviceCompliancePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Linux device compliance policies in Microsoft Intune using the `/deviceManagement/configurationPolicies` " +
			"endpoint. Linux device compliance policies define rules and settings that Linux devices must meet to be considered compliant with organizational " +
			"security requirements. These policies use the Settings Catalog configuration framework to provide granular control over Linux device compliance settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the Linux device compliance policy.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the Linux device compliance policy.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the Linux device compliance policy.",
			},
			"platforms": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("linux"),
				},
				MarkdownDescription: "Platform for which this policy applies. Always 'linux' for Linux device compliance policies.",
			},
			"technologies": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DefaultValueString("linuxMdm"),
				},
				MarkdownDescription: "Technology stack for the policy. Always 'linuxMdm' for Linux device compliance policies.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of Scope Tags for this Linux device compliance policy instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"settings_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings configured in this Linux device compliance policy.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates whether this Linux device compliance policy is assigned to any groups.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the Linux device compliance policy was last modified.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the Linux device compliance policy was created.",
			},

			// Individual Linux Compliance Settings
			"distribution_allowed_distros": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of allowed Linux distributions with version constraints.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Type of Linux distribution (e.g., 'ubuntu', 'rhel', 'centos', 'debian', 'suse').",
							Validators: []validator.String{
								stringvalidator.OneOf("ubuntu", "rhel", "centos", "debian", "suse"),
							},
						},
						"minimum_version": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Minimum version of the Linux distribution that is considered compliant.",
						},
						"maximum_version": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Maximum version of the Linux distribution that is considered compliant.",
						},
					},
				},
			},
			"custom_compliance_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether custom compliance rules are required for Linux devices.",
			},
			"custom_compliance_discovery_script": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Reference ID to the custom compliance discovery script for Linux devices.",
			},
			"custom_compliance_rules": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base64-encoded JSON string containing custom compliance rules for Linux devices.",
			},
			"device_encryption_required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether device encryption is required for Linux devices to be considered compliant.",
			},
			"password_policy_minimum_digits": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum number of digit characters required in the Linux device password.",
			},
			"password_policy_minimum_length": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum length required for the Linux device password.",
			},
			"password_policy_minimum_lowercase": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum number of lowercase characters required in the Linux device password.",
			},
			"password_policy_minimum_symbols": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum number of symbol characters required in the Linux device password.",
			},
			"password_policy_minimum_uppercase": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum number of uppercase characters required in the Linux device password.",
			},
			"scheduled_actions": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The list of scheduled action for this rule",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Name of the scheduled action rule. always use 'PasswordRequired' for Linux device compliance policies",
							PlanModifiers: []planmodifier.String{
								planmodifiers.DefaultValueString("PasswordRequired"),
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
