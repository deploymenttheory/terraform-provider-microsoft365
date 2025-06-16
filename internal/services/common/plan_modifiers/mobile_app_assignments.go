// In a file like mobileapp_assignment_plan_modifiers.go
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
)

// MobileAppAssignmentsListModifier handles the reordering of assignments by the API
func MobileAppAssignmentsListModifier() planmodifier.List {
	return &mobileAppAssignmentsModifier{}
}

// mobileAppAssignmentsModifier implements planmodifier.List
type mobileAppAssignmentsModifier struct{}

// Description returns a plain text description of the validator's behavior
func (m *mobileAppAssignmentsModifier) Description(ctx context.Context) string {
	return "Handles reordering of assignments by the Microsoft Graph API"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior
func (m *mobileAppAssignmentsModifier) MarkdownDescription(ctx context.Context) string {
	return "Handles reordering of assignments by the Microsoft Graph API"
}

// PlanModifyList is called when the provider is creating the plan for assignments
func (m *mobileAppAssignmentsModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Only apply this logic during update operations
	if req.State.Raw.IsNull() {
		return
	}

	// If plan is null, we're removing all assignments - no need to sort
	if req.Plan.Raw.IsNull() {
		return
	}

	var planAssignments, stateAssignments []sharedmodels.MobileAppAssignmentResourceModel

	diags := req.PlanValue.ElementsAs(ctx, &planAssignments, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = req.StateValue.ElementsAs(ctx, &stateAssignments, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Order doesn't matter for creation (API will reorder)
	// But for updates or reads, need to match API's ordering
	// Order doesn't matter for creation (API will reorder)
	// But for updates or reads, need to match API's ordering
	if len(planAssignments) > 0 && len(stateAssignments) > 0 {
		// Sort both the plan and state assignments the same way
		sharedstater.SortMobileAppAssignments(planAssignments)
		sharedstater.SortMobileAppAssignments(stateAssignments)

		// Create a new plan value with sorted assignments
		elemType := req.PlanValue.Type(ctx).(types.ListType).ElemType

		// Convert the sorted assignments back to a list value
		sortedAssignments := make([]attr.Value, len(planAssignments))

		for i, assignment := range planAssignments {
			// Convert each assignment to attr.Value
			assignmentMap := map[string]attr.Value{
				"id":        types.StringValue(assignment.Id.ValueString()),
				"intent":    types.StringValue(assignment.Intent.ValueString()),
				"source":    types.StringValue(assignment.Source.ValueString()),
				"source_id": assignment.SourceId,
				"target": types.ObjectValueMust(
					map[string]attr.Type{
						"target_type": types.StringType,
						"device_and_app_management_assignment_filter_id":   types.StringType,
						"device_and_app_management_assignment_filter_type": types.StringType,
						"group_id":      types.StringType,
						"collection_id": types.StringType,
					},
					map[string]attr.Value{
						"target_type": types.StringValue(assignment.Target.TargetType.ValueString()),
						"device_and_app_management_assignment_filter_id":   assignment.Target.DeviceAndAppManagementAssignmentFilterId,
						"device_and_app_management_assignment_filter_type": types.StringValue(assignment.Target.DeviceAndAppManagementAssignmentFilterType.ValueString()),
						"group_id":      assignment.Target.GroupId,
						"collection_id": assignment.Target.CollectionId,
					},
				),
			}

			// Handle settings which could be null
			if assignment.Settings != nil {
				// Create a settings object based on which type is set
				settingsMap := map[string]attr.Value{
					"android_managed_store":        types.ObjectNull(map[string]attr.Type{}),
					"ios_lob":                      types.ObjectNull(map[string]attr.Type{}),
					"ios_store":                    types.ObjectNull(map[string]attr.Type{}),
					"ios_vpp":                      types.ObjectNull(map[string]attr.Type{}),
					"macos_lob":                    types.ObjectNull(map[string]attr.Type{}),
					"macos_vpp":                    types.ObjectNull(map[string]attr.Type{}),
					"microsoft_store_for_business": types.ObjectNull(map[string]attr.Type{}),
					"win32_catalog":                types.ObjectNull(map[string]attr.Type{}),
					"win32_lob":                    types.ObjectNull(map[string]attr.Type{}),
					"win_get":                      types.ObjectNull(map[string]attr.Type{}),
					"windows_app_x":                types.ObjectNull(map[string]attr.Type{}),
					"windows_universal_app_x":      types.ObjectNull(map[string]attr.Type{}),
				}

				// Handle WinGet settings
				if assignment.Settings.WinGet != nil {
					// Initialize the winGetMap with appropriate types for nested objects
					winGetMap := map[string]attr.Value{
						"notifications": types.StringValue(assignment.Settings.WinGet.Notifications.ValueString()),
						// Initialize with empty objects that match the expected schema
						"install_time_settings": types.ObjectValueMust(
							map[string]attr.Type{
								"deadline_date_time": types.StringType,
								"use_local_time":     types.BoolType,
							},
							map[string]attr.Value{
								"deadline_date_time": types.StringNull(),
								"use_local_time":     types.BoolNull(),
							},
						),
						"restart_settings": types.ObjectValueMust(
							map[string]attr.Type{
								"countdown_display_before_restart_in_minutes":     types.Int64Type,
								"grace_period_in_minutes":                         types.Int64Type,
								"restart_notification_snooze_duration_in_minutes": types.Int64Type,
							},
							map[string]attr.Value{
								"countdown_display_before_restart_in_minutes":     types.Int64Null(),
								"grace_period_in_minutes":                         types.Int64Null(),
								"restart_notification_snooze_duration_in_minutes": types.Int64Null(),
							},
						),
					}

					// Handle install time settings
					if assignment.Settings.WinGet.InstallTimeSettings != nil {
						installTimeMap := map[string]attr.Value{
							"deadline_date_time": types.StringValue(assignment.Settings.WinGet.InstallTimeSettings.DeadlineDateTime.ValueString()),
							"use_local_time":     types.BoolValue(assignment.Settings.WinGet.InstallTimeSettings.UseLocalTime.ValueBool()),
						}

						winGetMap["install_time_settings"] = types.ObjectValueMust(
							map[string]attr.Type{
								"deadline_date_time": types.StringType,
								"use_local_time":     types.BoolType,
							},
							installTimeMap,
						)
					}

					// Handle restart settings
					if assignment.Settings.WinGet.RestartSettings != nil {
						restartMap := map[string]attr.Value{
							"countdown_display_before_restart_in_minutes":     types.Int64Value(int64(assignment.Settings.WinGet.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt32())),
							"grace_period_in_minutes":                         types.Int64Value(int64(assignment.Settings.WinGet.RestartSettings.GracePeriodInMinutes.ValueInt32())),
							"restart_notification_snooze_duration_in_minutes": types.Int64Value(int64(assignment.Settings.WinGet.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt32())),
						}

						winGetMap["restart_settings"] = types.ObjectValueMust(
							map[string]attr.Type{
								"countdown_display_before_restart_in_minutes":     types.Int64Type,
								"grace_period_in_minutes":                         types.Int64Type,
								"restart_notification_snooze_duration_in_minutes": types.Int64Type,
							},
							restartMap,
						)
					}

					// Create the win_get object with consistent schema
					settingsMap["win_get"] = types.ObjectValueMust(
						map[string]attr.Type{
							"notifications": types.StringType,
							"install_time_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
								"deadline_date_time": types.StringType,
								"use_local_time":     types.BoolType,
							}},
							"restart_settings": types.ObjectType{AttrTypes: map[string]attr.Type{
								"countdown_display_before_restart_in_minutes":     types.Int64Type,
								"grace_period_in_minutes":                         types.Int64Type,
								"restart_notification_snooze_duration_in_minutes": types.Int64Type,
							}},
						},
						winGetMap,
					)
				}

				// Handle other settings types similarly...
				// For example, for Android, iOS, etc.

				assignmentMap["settings"] = types.ObjectValueMust(
					map[string]attr.Type{
						"android_managed_store":        types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"ios_lob":                      types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"ios_store":                    types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"ios_vpp":                      types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"macos_lob":                    types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"macos_vpp":                    types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"microsoft_store_for_business": types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"win32_catalog":                types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"win32_lob":                    types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"win_get":                      types.ObjectType{AttrTypes: map[string]attr.Type{"notifications": types.StringType, "install_time_settings": types.ObjectType{AttrTypes: map[string]attr.Type{}}, "restart_settings": types.ObjectType{AttrTypes: map[string]attr.Type{}}}},
						"windows_app_x":                types.ObjectType{AttrTypes: map[string]attr.Type{}},
						"windows_universal_app_x":      types.ObjectType{AttrTypes: map[string]attr.Type{}},
					},
					settingsMap,
				)
			} else {
				// If settings is null, add an empty settings object
				assignmentMap["settings"] = types.ObjectNull(map[string]attr.Type{})
			}

			// Create the assignment object and add it to the list
			assignmentObj, diags := types.ObjectValue(elemType.(types.ObjectType).AttrTypes, assignmentMap)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}

			sortedAssignments[i] = assignmentObj
		}

		// Create the sorted list and set it as the plan value
		sortedList, diags := types.ListValue(elemType, sortedAssignments)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		resp.PlanValue = sortedList
	}
}
