// This is a shared resource used by multiple intune app types to define the assignments for the app.
// As such, i have taken the design decision to make this a shared resource that will be called by the
// different app types that require it. Therefore it doesn't have it's own terraform CRUD operations directly

package graphBetaMobileAppAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MobileAppAssignments returns a schema attribute for device and app management resource assignments that can be reused across different app types.
func MobileAppAssignments() schema.Attribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: "The list of group assignments for this device and app management resource.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"target": schema.SingleNestedAttribute{
					Required: true,
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
					Description: "The target for this assignment.",
				},
				"intent": schema.StringAttribute{
					Required:    true,
					Description: "The intent of the assignment. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
					Validators: []validator.String{
						stringvalidator.OneOf("available", "required", "uninstall", "availableWithoutEnrollment"),
					},
				},
				"settings": schema.SingleNestedAttribute{
					Optional: true,
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
					Description: "The settings for this assignment.",
				},
			},
		},
	}
}
