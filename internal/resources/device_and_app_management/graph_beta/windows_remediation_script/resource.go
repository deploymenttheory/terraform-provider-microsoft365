package graphBetaWindowsRemediationScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_windows_remediation_script"
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
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceHealthScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
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
		MarkdownDescription: "Manages a Device Health Script in Microsoft Intune. Device health scripts can detect and remediate issues on Windows devices.",
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
				MarkdownDescription: "Indicate whether PowerShell script(s) should run as 32-bit.",
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
				MarkdownDescription: "Indicate whether the script signature needs be checked.",
			},
			"detection_script_content": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The entire content of the detection PowerShell script.",
			},
			"remediation_script_content": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
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
				MarkdownDescription: "List of Scope Tag IDs for the device health script.",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")})),
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
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"assignment": AssignmentBlock(),
		},
	}
}

// AssignmentBlock returns the schema block for Windows remediation script assignments

// AssignmentBlock returns the schema block for Windows remediation script assignments
func AssignmentBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		MarkdownDescription: "List of assignment configurations for the device health script",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"all_devices": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Assign to all devices. Cannot be used with all_users or include_groups.",
				},
				"all_devices_filter_type": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Filter type for all devices assignment. Can be 'include' or 'exclude'.",
					Validators: []validator.String{
						stringvalidator.OneOf("include", "exclude"),
					},
				},
				"all_devices_filter_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Filter ID for all devices assignment.",
				},
				"all_users": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Assign to all users. Cannot be used with all_devices or include_groups.",
				},
				"all_users_filter_type": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Filter type for all users assignment. Can be 'include' or 'exclude'.",
					Validators: []validator.String{
						stringvalidator.OneOf("include", "exclude"),
					},
				},
				"all_users_filter_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Filter ID for all users assignment.",
				},
				"include_groups": schema.SetNestedAttribute{
					Optional:            true,
					MarkdownDescription: "Groups to include in the assignment. Cannot be used with all_devices or all_users.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"group_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Group ID to include.",
							},
							"include_groups_filter_type": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Filter type for include group assignment. Can be 'include' or 'exclude'.",
								Validators: []validator.String{
									stringvalidator.OneOf("include", "exclude"),
								},
							},
							"include_groups_filter_id": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Filter ID for include group assignment.",
							},
							"run_remediation_script": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(true),
								MarkdownDescription: "Whether to run the remediation script for this group assignment.",
							},
							"run_schedule": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Run schedule for this group assignment.",
								Attributes: map[string]schema.Attribute{
									"schedule_type": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Type of schedule. Can be 'daily', 'hourly', or 'once'.",
										Validators: []validator.String{
											stringvalidator.OneOf("daily", "hourly", "once"),
										},
									},
									"interval": schema.Int32Attribute{
										Optional:            true,
										Computed:            true,
										Default:             int32default.StaticInt32(1),
										MarkdownDescription: "Repeat interval for the schedule.For 'daily' the interal represents days, for 'hourly' the interval represents hours. ",
									},
									"time": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "Time of day for daily and once schedules (e.g., '14:30').",
									},
									"date": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "Date for once schedule (e.g., '2025-05-01').",
									},
									"use_utc": schema.BoolAttribute{
										Optional:            true,
										Computed:            true,
										Default:             booldefault.StaticBool(false),
										MarkdownDescription: "Whether to use UTC time.",
									},
								},
							},
						},
					},
				},

				"exclude_group_ids": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					MarkdownDescription: "Group IDs to exclude from the assignment.",
				},
			},
		},
	}
}
