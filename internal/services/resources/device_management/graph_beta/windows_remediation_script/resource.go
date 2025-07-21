package graphBetaWindowsRemediationScript

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_remediation_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceHealthScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceHealthScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceHealthScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceHealthScriptResource{}
)

func NewDeviceHealthScriptResource() resource.Resource {
	return &DeviceHealthScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceHealthScripts",
	}
}

type DeviceHealthScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceHealthScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *DeviceHealthScriptResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceHealthScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceHealthScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ModifyPlan modifies the plan for the resource.
func (r *DeviceHealthScriptResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No modifications needed at this time
}

// Schema defines the schema for the resource.
func (r *DeviceHealthScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows remediation scripts using the `/deviceManagement/deviceHealthScripts` endpoint. Remediation scripts enable proactive detection and automatic remediation of common issues on Windows devices through PowerShell detection scripts paired with remediation scripts that execute when problems are identified.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the device health script.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device health script.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the device health script.",
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device health script publisher.",
			},
			"run_as_32_bit": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicate whether PowerShell script(s) should run as 32-bit. Default is false, which runs script in 64-bit PowerShell.",
			},
			"run_as_account": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Indicates the type of execution context. Possible values are: system, user.",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicate whether the script signature needs be checked. Default is false, which does not check the script signature.",
			},
			"detection_script_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The entire content of the detection PowerShell script.",
			},
			"remediation_script_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The entire content of the remediation PowerShell script.",
			},
			"detection_script_parameters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of ComplexType DetectionScriptParameters objects.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the param",
						},
						"description": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The description of the param",
						},
						"is_required": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether the param is required",
						},
						"apply_default_value_when_not_assigned": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether Apply DefaultValue When Not Assigned",
						},
					},
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
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the device health script.",
			},
			"is_global_script": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Determines if this is Microsoft Proprietary Script. Proprietary scripts are read-only.",
			},
			"device_health_script_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DeviceHealthScriptType for the script policy. Possible values are: deviceHealthScript, managedInstallerScript.",
				Validators: []validator.String{
					stringvalidator.OneOf("deviceHealthScript", "managedInstallerScript"),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device health script was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device health script was modified. This property is read-only.",
			},
			"highest_available_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Highest available version for a Microsoft Proprietary script.",
			},
			"assignments": AssignmentBlock(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

// AssignmentBlock returns the schema for the assignments block
func AssignmentBlock() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the Windows remediation script. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				// Target assignment fields - only one should be used at a time
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"allDevicesAssignmentTarget",
							"allLicensedUsersAssignmentTarget",
							"groupAssignmentTarget",
							"exclusionGroupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
				// Assignment filter fields
				"filter_id": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "ID of the filter to apply to the assignment.",
					Default:             stringdefault.StaticString("00000000-0000-0000-0000-000000000000"),
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
				"filter_type": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.",
					Computed:            true,
					Default:             stringdefault.StaticString("none"),
					Validators: []validator.String{
						stringvalidator.OneOf("include", "exclude", "none"),
					},
				},
				// Schedule configuration - only one should be used at a time
				"daily_schedule": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Configuration for daily schedule execution. Only one schedule type (daily_schedule, hourly_schedule, or run_once_schedule) should be specified per assignment.",
					Attributes: map[string]schema.Attribute{
						"interval": schema.Int32Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int32default.StaticInt32(1),
							MarkdownDescription: "Days between runs. Default is 1.",
						},
						"time": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Time of day in format 'HH:MM:SS'.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
									"Time must be in the format 'HH:MM:SS' (24-hour format)",
								),
							},
						},
						"use_utc": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "Whether to use UTC time. Default is false (local time).",
						},
					},
				},
				"hourly_schedule": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Configuration for hourly schedule execution. Only one schedule type (daily_schedule, hourly_schedule, or run_once_schedule) should be specified per assignment.",
					Attributes: map[string]schema.Attribute{
						"interval": schema.Int32Attribute{
							Optional:            true,
							Computed:            true,
							Default:             int32default.StaticInt32(1),
							MarkdownDescription: "Hours between runs. Default is 1.",
						},
					},
				},
				"run_once_schedule": schema.SingleNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Configuration for one-time execution. Only one schedule type (daily_schedule, hourly_schedule, or run_once_schedule) should be specified per assignment.",
					Attributes: map[string]schema.Attribute{
						"date": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Date for the one-time execution in format 'YYYY-MM-DD'.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.DateFormatYYYYMMDDRegex),
									"Date must be in the format 'YYYY-MM-DD'",
								),
							},
						},
						"time": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Time of day in format 'HH:MM:SS'.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
									"Time must be in the format 'HH:MM:SS' (24-hour format)",
								),
							},
						},
						"use_utc": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							MarkdownDescription: "Whether to use UTC time. Default is false (local time).",
						},
					},
				},
			},
		},
	}
}
