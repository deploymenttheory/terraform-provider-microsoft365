package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
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
	ResourceName  = "graph_beta_device_management_macos_software_update_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &MacOSSoftwareUpdateConfigurationResource{}
	_ resource.ResourceWithConfigure   = &MacOSSoftwareUpdateConfigurationResource{}
	_ resource.ResourceWithImportState = &MacOSSoftwareUpdateConfigurationResource{}
	_ resource.ResourceWithModifyPlan  = &MacOSSoftwareUpdateConfigurationResource{}
)

func NewMacOSSoftwareUpdateConfigurationResource() resource.Resource {
	return &MacOSSoftwareUpdateConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceConfigurations",
	}
}

type MacOSSoftwareUpdateConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *MacOSSoftwareUpdateConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *MacOSSoftwareUpdateConfigurationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *MacOSSoftwareUpdateConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *MacOSSoftwareUpdateConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *MacOSSoftwareUpdateConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS software update configurations using the `/deviceManagement/deviceConfigurations` endpoint. See [macOSSoftwareUpdateConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdateconfiguration?view=graph-rest-beta) for details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the macOS software update configuration.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the macOS software update configuration",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Admin provided description of the device configuration.",
				Optional:            true,
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of Scope Tags for this Entity instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"critical_update_behavior": schema.StringAttribute{
				MarkdownDescription: "Update behavior for critical updates. Possible values: `notConfigured`, `default`, `downloadOnly`, `installASAP`, `notifyOnly`, `installLater`. See [macOSSoftwareUpdateBehavior](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdatebehavior?view=graph-rest-beta).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
				},
			},
			"config_data_update_behavior": schema.StringAttribute{
				MarkdownDescription: "Update behavior for configuration data file updates. Possible values: `notConfigured`, `default`, `downloadOnly`, `installASAP`, `notifyOnly`, `installLater`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
				},
			},
			"firmware_update_behavior": schema.StringAttribute{
				MarkdownDescription: "Update behavior for firmware updates. Possible values: `notConfigured`, `default`, `downloadOnly`, `installASAP`, `notifyOnly`, `installLater`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
				},
			},
			"all_other_update_behavior": schema.StringAttribute{
				MarkdownDescription: "Update behavior for all other updates. Possible values: `notConfigured`, `default`, `downloadOnly`, `installASAP`, `notifyOnly`, `installLater`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
				},
			},
			"update_schedule_type": schema.StringAttribute{
				MarkdownDescription: "Update schedule type. Possible values: `alwaysUpdate`, `updateDuringTimeWindows`, `updateOutsideOfTimeWindows`. See [macOSSoftwareUpdateScheduleType](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdatescheduletype?view=graph-rest-beta).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("alwaysUpdate", "updateDuringTimeWindows", "updateOutsideOfTimeWindows"),
				},
			},
			"custom_update_time_windows": schema.ListNestedAttribute{
				MarkdownDescription: "Custom time windows when updates will be allowed or blocked. Maximum 20 elements. See [customUpdateTimeWindow](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customupdatetimewindow?view=graph-rest-beta).",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_day": schema.StringAttribute{
							MarkdownDescription: "Start day of the time window. Possible values: `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
							},
						},
						"end_day": schema.StringAttribute{
							MarkdownDescription: "End day of the time window. Possible values: `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
							},
						},
						"start_time": schema.StringAttribute{
							MarkdownDescription: "Start time of the window in `HH:MM:SS` format.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
									"Time must be in the format 'HH:MM:SS' (24-hour format)",
								),
							},
						},
						"end_time": schema.StringAttribute{
							MarkdownDescription: "End time of the window in `HH:MM:SS` format.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
									"Time must be in the format 'HH:MM:SS' (24-hour format)",
								),
							},
						},
					},
				},
			},
			"update_time_window_utc_offset_in_minutes": schema.Int32Attribute{
				MarkdownDescription: "Minutes indicating UTC offset for each update time window.",
				Required:            true,
			},
			"max_user_deferrals_count": schema.Int32Attribute{
				MarkdownDescription: "The maximum number of times the system allows the user to postpone an update before it's installed. Supported values: 0 - 365.",
				Optional:            true,
				Validators: []validator.Int32{
					int32validator.Between(0, 365),
				},
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "The scheduling priority for downloading and preparing the requested update. Possible values: `low`, `high`, `unknownFutureValue`. See [macOSPriority](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macospriority?view=graph-rest-beta).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("low", "high", "unknownFutureValue"),
				},
			},
			"assignments": assignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

// assignmentsSchema returns the schema for the assignments block
func assignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required:            true,
		MarkdownDescription: "Assignment configuration for the macOS software update configuration.",
		Attributes: map[string]schema.Attribute{
			"all_devices": schema.BoolAttribute{
				MarkdownDescription: "Whether to assign the configuration to all devices.",
				Required:            true,
			},
			"all_users": schema.BoolAttribute{
				MarkdownDescription: "Whether to assign the configuration to all users.",
				Required:            true,
			},
			"include_group_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of group IDs to include in the assignment.",
				Optional:            true,
			},
			"exclude_group_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of group IDs to exclude from the assignment.",
				Optional:            true,
			},
		},
	}
}
