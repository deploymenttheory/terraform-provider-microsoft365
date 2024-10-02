package graphBetaMobileAppAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &MobileAppAssignmentResource{}
var _ resource.ResourceWithConfigure = &MobileAppAssignmentResource{}
var _ resource.ResourceWithImportState = &MobileAppAssignmentResource{}

func NewMobileAppAssignmentResource() resource.Resource {
	return &MobileAppAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type MobileAppAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// GetID returns the ID of a resource from the state model.
func (s *MobileAppAssignmentResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *MobileAppAssignmentResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *MobileAppAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_mobile_app_assignment"
}

// Configure sets the client for the resource.
func (r *MobileAppAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *MobileAppAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *MobileAppAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Mobile App Assignment in Microsoft Intune. Used by different app types to define the assignments" +
			"for the app. Used by winget_app, windows_web_app, windows_universal_appx, windows_microsoft_edge_app, win32_lob_app, windows_web_app, " +
			"windows_office_suite_app, managed_ios_store_app, managed_ios_lob_app, managed_android_store_app, managed_ios_lob_app, mac_web_clip, " +
			"macos_vpp_app, macos_pkg_app, macos_office_suite_app, macos_microsoft_edge_app, macOS_microsoft_defender_app, macOS_lob_app, " +
			"macOS_dmg_app, ios_vpp_app, ios_store_app, ios_lob_app, ios_ipados_web_clip, android_store_app, android_managed_webstore_app, " +
			"android_managed_store_app, android_managed_lob_app, android_for_work_app",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the mobile app assignment.",
			},
			"source_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The identifier of the source mobile app.",
			},
			"target": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The target for this assignment.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The type of target. Possible values are: allLicensedUsers, allDevices, group.",
						Validators: []validator.String{
							stringvalidator.OneOf("allLicensedUsers", "allDevices", "group"),
						},
					},
					"device_and_app_management_assignment_filter_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The ID of the filter for the target assignment.",
					},
					"device_and_app_management_assignment_filter_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The type of filter for the target assignment. Possible values are: none, include, exclude.",
						Validators: []validator.String{
							stringvalidator.OneOf("none", "include", "exclude"),
						},
					},
					"group_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The ID of the group to assign the app to. Required when type is 'group'.",
					},
				},
			},
			"intent": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The intent of the assignment. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
				Validators: []validator.String{
					stringvalidator.OneOf("available", "required", "uninstall", "availableWithoutEnrollment"),
				},
			},
			"settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The settings for this assignment.",
				Attributes: map[string]schema.Attribute{
					"notifications": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The notification setting for the assignment. Possible values are: showAll, showReboot, hideAll.",
						Validators: []validator.String{
							stringvalidator.OneOf("showAll", "showReboot", "hideAll"),
						},
					},
					"restart_settings": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The restart settings for the assignment.",
						Attributes: map[string]schema.Attribute{
							"grace_period_in_minutes": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "The grace period before a restart in minutes.",
							},
							"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "The countdown display before restart in minutes.",
							},
							"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "The snooze duration for the restart notification in minutes.",
							},
						},
					},
					"install_time_settings": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The install time settings for the assignment.",
						Attributes: map[string]schema.Attribute{
							"use_local_time": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Indicates whether to use local time for the assignment.",
							},
							"deadline_date_time": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The deadline date and time for the assignment.",
							},
						},
					},
				},
			},
			"source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The source of the assignment.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
