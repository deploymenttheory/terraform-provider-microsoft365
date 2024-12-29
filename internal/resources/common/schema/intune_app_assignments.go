package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func IntuneMobileAppAssignmentsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Configuration for Mobile App Assignment, including settings and targets for Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Key of the entity. This property is read-only.",
				Computed:            true,
			},
			"mobile_app_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the mobile app associated with this assignment.",
				Required:            true,
			},
			"mobile_app_assignments": schema.ListNestedAttribute{
				MarkdownDescription: "List of assignments for the mobile app.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Key of the assignment entity. This property is read-only.",
							Computed:            true,
						},
						"intent": schema.StringAttribute{
							MarkdownDescription: "The install intent defined by the admin. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
							Required:            true,
						},
						"source": schema.StringAttribute{
							MarkdownDescription: "The resource type which is the source for the assignment. Possible values are: direct, policySets. This property is read-only.",
							Computed:            true,
						},
						"source_id": schema.StringAttribute{
							MarkdownDescription: "The identifier of the source of the assignment. This property is read-only.",
							Computed:            true,
						},
						"target": schema.SingleNestedAttribute{
							MarkdownDescription: "The target group assignment defined by the admin.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"target_type": schema.StringAttribute{
									MarkdownDescription: "The type of target assignment. Possible values include groupAssignmentTarget, allLicensedUsersAssignmentTarget, etc.",
									Required:            true,
								},
								"group_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the target group.",
									Optional:            true,
								},
								"device_and_app_management_assignment_filter_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the filter for the target assignment.",
									Optional:            true,
								},
								"device_and_app_management_assignment_filter_type": schema.StringAttribute{
									MarkdownDescription: "The type of filter for the target assignment. Possible values are: none, include, exclude.",
									Optional:            true,
								},
								"is_exclusion_group": schema.BoolAttribute{
									MarkdownDescription: "Indicates whether this is an exclusion group.",
									Optional:            true,
								},
							},
						},
						"settings": schema.SingleNestedAttribute{
							MarkdownDescription: "The settings for target assignment defined by the admin.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"notifications": schema.StringAttribute{
									MarkdownDescription: "The notification settings for the assignment.",
									Optional:            true,
								},
								"install_time_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "Settings related to install time.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"use_local_time": schema.BoolAttribute{
											MarkdownDescription: "Whether the local device time or UTC time should be used when determining the deadline times.",
											Optional:            true,
										},
										"deadline_date_time": schema.StringAttribute{
											MarkdownDescription: "The time at which the app should be installed.",
											Optional:            true,
										},
									},
								},
								"restart_settings": schema.SingleNestedAttribute{
									MarkdownDescription: "Settings related to restarts after installation.",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"grace_period_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation.",
											Optional:            true,
										},
										"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts.",
											Optional:            true,
										},
										"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
											MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
											Optional:            true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
