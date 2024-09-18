// This is a shared resource used by multiple intune app types to define the assignments for the app.
// As such, i have taken the design decision to make this a shared resource that will be called by the
// different app types that require it. It shall be triggered by the inclusion of app assignments within
// the hcl of app resource. it shall not be a stand alone resource that can be created on its own and
// is therefore not exposed to the user.

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
	return &MobileAppAssignmentResource{}
}

type MobileAppAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
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
		Description: "Manages a Mobile App Assignment in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the mobile app assignment.",
			},
			"source_id": schema.StringAttribute{
				Optional:    true,
				Description: "The identifier of the source mobile app.",
			},
			"target": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The target for this assignment.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of target. Possible values are: allLicensedUsers, allDevices, group.",
						Validators: []validator.String{
							stringvalidator.OneOf("allLicensedUsers", "allDevices", "group"),
						},
					},
					"device_and_app_management_assignment_filter_id": schema.StringAttribute{
						Optional:    true,
						Description: "The ID of the filter for the target assignment.",
					},
					"device_and_app_management_assignment_filter_type": schema.StringAttribute{
						Optional:    true,
						Description: "The type of filter for the target assignment. Possible values are: none, include, exclude.",
						Validators: []validator.String{
							stringvalidator.OneOf("none", "include", "exclude"),
						},
					},
					"group_id": schema.StringAttribute{
						Optional:    true,
						Description: "The ID of the group to assign the app to. Required when type is 'group'.",
					},
				},
			},
			"intent": schema.StringAttribute{
				Required:    true,
				Description: "The intent of the assignment. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
				Validators: []validator.String{
					stringvalidator.OneOf("available", "required", "uninstall", "availableWithoutEnrollment"),
				},
			},
			"settings": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The settings for this assignment.",
				Attributes: map[string]schema.Attribute{
					"notifications": schema.StringAttribute{
						Optional:    true,
						Description: "The notification setting for the assignment. Possible values are: showAll, showReboot, hideAll.",
						Validators: []validator.String{
							stringvalidator.OneOf("showAll", "showReboot", "hideAll"),
						},
					},
					"restart_settings": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The restart settings for the assignment.",
						Attributes: map[string]schema.Attribute{
							"grace_period_in_minutes": schema.Int64Attribute{
								Optional:    true,
								Description: "The grace period before a restart in minutes.",
							},
							"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
								Optional:    true,
								Description: "The countdown display before restart in minutes.",
							},
							"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
								Optional:    true,
								Description: "The snooze duration for the restart notification in minutes.",
							},
						},
					},
					"install_time_settings": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "The install time settings for the assignment.",
						Attributes: map[string]schema.Attribute{
							"use_local_time": schema.BoolAttribute{
								Optional:    true,
								Description: "Indicates whether to use local time for the assignment.",
							},
							"deadline_date_time": schema.StringAttribute{
								Optional:    true,
								Description: "The deadline date and time for the assignment.",
							},
						},
					},
				},
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "The source of the assignment.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
