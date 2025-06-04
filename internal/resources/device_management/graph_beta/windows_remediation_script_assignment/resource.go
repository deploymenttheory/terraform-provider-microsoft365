package graphBetaWindowsRemediationScriptAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_remediation_script_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceHealthScriptAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceHealthScriptAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceHealthScriptAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceHealthScriptAssignmentResource{}
)

func NewDeviceHealthScriptAssignmentResource() resource.Resource {
	return &DeviceHealthScriptAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/deviceHealthScripts",
	}
}

type DeviceHealthScriptAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceHealthScriptAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *DeviceHealthScriptAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceHealthScriptAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceHealthScriptAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *DeviceHealthScriptAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Device Health Script Assignments in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the device health script assignment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_health_script_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the device health script to attach the assignment to.",
			},
			"run_remediation_script": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicates whether remediation script should be run. Default is false.",
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"target_type": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The target group type for the script assignment. Possible values are:\n\n" +
							"- **allDevices**: Target all devices in the tenant\n" +
							"- **allLicensedUsers**: Target all licensed users in the tenant\n" +
							"- **configurationManagerCollection**: Target System Center Configuration Manager collection\n" +
							"- **exclusionGroupAssignment**: Target a specific Entra ID group for exclusion\n" +
							"- **groupAssignment**: Target a specific Entra ID group",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allDevices",
								"allLicensedUsers",
								"configurationManagerCollection",
								"exclusionGroupAssignment",
								"groupAssignment",
							),
						},
					},
					"group_id": schema.StringAttribute{
						MarkdownDescription: "The entra ID group ID for the assignment target. Required when target_type is 'groupAssignment' or 'exclusionGroupAssignment'.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
								"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
							),
						},
					},
					"collection_id": schema.StringAttribute{
						MarkdownDescription: "The SCCM group collection ID for the assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[A-Za-z]{2,8}[0-9A-Za-z]{8}$`),
								"Must be a valid SCCM collection ID format. Default collections start with 'SMS' followed by an alphanumeric ID. Custom collections start with your site code (e.g., 'MEM') followed by an alphanumeric ID.",
							),
						},
					},
					"device_and_app_management_assignment_filter_id": schema.StringAttribute{
						MarkdownDescription: "The Id of the scope filter applied to the target assignment.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
								"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
							),
						},
					},
					"device_and_app_management_assignment_filter_type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("none"),
						MarkdownDescription: "The type of scope filter for the target assignment. Defaults to 'none'. Possible values are:\n\n" +
							"- **include**: Only include devices or users matching the filter\n" +
							"- **exclude**: Exclude devices or users matching the filter\n" +
							"- **none**: No assignment filter applied",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"include",
								"exclude",
								"none",
							),
						},
					},
				},
			},
			"run_schedule": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"daily": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"interval": schema.Int32Attribute{
								Optional:            true,
								Computed:            true,
								Default:             int32default.StaticInt32(1),
								MarkdownDescription: "Number of days between runs. Default is 1. Valid values are between 1 and 30.",
							},
							"use_utc": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to use UTC time or user device's local time. Default is false (local time).",
							},
							"time": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Time of day to run the script, formatted as HH:MM:SS.",
							},
						},
					},
					"hourly": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"interval": schema.Int32Attribute{
								Optional:            true,
								Computed:            true,
								Default:             int32default.StaticInt32(1),
								MarkdownDescription: "Number of hours between runs. Default is 1. Valid values are between 1 and 23.",
							},
							"use_utc": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to use UTC time or user device's local time. Default is false (local time).",
							},
						},
					},
					"once": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"start_date_time": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Date and time to run the script, formatted as ISO 8601 timestamp.",
							},
							"use_utc": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
								MarkdownDescription: "Whether to use UTC time or user device's local time. Default is false (local time).",
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
