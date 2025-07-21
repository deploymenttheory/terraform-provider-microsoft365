package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ValidateAssignments validates the assignments according to the following rules:
// 1. If all_devices is set, no other group assignments are allowed
// 2. If all_users is set, no other group assignments can be set
// 3. All_devices and all_users cannot be set at the same time
// 4. Exclude assignments can always be set regardless
// 5. A group can only be defined once across all include and exclude assignments
// 6. Each assignment must have exactly one schedule type defined (or none for exclusions)
// 7. group_id must be provided for groupAssignmentTarget and exclusionGroupAssignmentTarget
func ValidateAssignments(ctx context.Context, data *DeviceHealthScriptResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		// No assignments to validate
		return diags
	}

	// Extract assignments using the proper struct types
	var assignments []WindowsRemediationScriptAssignmentModel
	diags.Append(data.Assignments.ElementsAs(ctx, &assignments, false)...)
	if diags.HasError() {
		return diags
	}

	// Track group IDs to detect duplicates
	groupIDs := make(map[string]bool)

	for i, assignment := range assignments {
		// Validate target type
		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			diags.AddError(
				"Invalid Assignment Configuration",
				fmt.Sprintf("Assignment at index %d is missing a target type", i),
			)
			continue
		}

		targetType := assignment.Type.ValueString()

		// Validate group ID is provided for group targets
		if targetType == "groupAssignmentTarget" || targetType == "exclusionGroupAssignmentTarget" {
			if assignment.GroupId.IsNull() || assignment.GroupId.IsUnknown() || assignment.GroupId.ValueString() == "" {
				diags.AddError(
					"Invalid Assignment Configuration",
					fmt.Sprintf("Assignment at index %d has target type '%s' but is missing a group_id", i, targetType),
				)
			} else {
				groupID := assignment.GroupId.ValueString()

				// Validate it's not the default GUID value
				if groupID == "00000000-0000-0000-0000-000000000000" {
					diags.AddError(
						"Invalid Assignment Configuration",
						fmt.Sprintf("Assignment at index %d has target type '%s' but group_id cannot be the default value '00000000-0000-0000-0000-000000000000'", i, targetType),
					)
				} else {
					// Check for duplicate group IDs
					if _, exists := groupIDs[groupID]; exists {
						diags.AddError(
							"Duplicate Group Assignment",
							fmt.Sprintf("Group ID '%s' is assigned multiple times. Each group can only be assigned once.", groupID),
						)
					} else {
						groupIDs[groupID] = true
					}
				}
			}
		} else {
			// For allDevicesAssignmentTarget and allLicensedUsersAssignmentTarget, group_id should not be set
			if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
				diags.AddError(
					"Invalid Assignment Configuration",
					fmt.Sprintf("Assignment at index %d has target type '%s' but should not have a group_id", i, targetType),
				)
			}
		}

		// Validate filter type if filter ID is provided
		if !assignment.FilterId.IsNull() && !assignment.FilterId.IsUnknown() && assignment.FilterId.ValueString() != "" && assignment.FilterId.ValueString() != "00000000-0000-0000-0000-000000000000" {
			if assignment.FilterType.IsNull() || assignment.FilterType.IsUnknown() || assignment.FilterType.ValueString() == "" || assignment.FilterType.ValueString() == "none" {
				diags.AddError(
					"Invalid Assignment Configuration",
					fmt.Sprintf("Assignment at index %d has a filter_id but filter_type must be 'include' or 'exclude' (not 'none' or empty)", i),
				)
			} else {
				filterType := assignment.FilterType.ValueString()
				if filterType != "include" && filterType != "exclude" {
					diags.AddError(
						"Invalid Assignment Configuration",
						fmt.Sprintf("Assignment at index %d has an invalid filter_type '%s'. Must be 'include' or 'exclude' when filter_id is provided", i, filterType),
					)
				}
			}
		}

		// Validate schedule configuration - exactly one schedule type must be defined for non-exclusion targets
		scheduleCount := 0

		// Check daily schedule
		if assignment.DailySchedule != nil {
			scheduleCount++

			// Validate daily schedule has required fields
			if assignment.DailySchedule.Time.IsNull() || assignment.DailySchedule.Time.IsUnknown() || assignment.DailySchedule.Time.ValueString() == "" {
				diags.AddError(
					"Invalid Schedule Configuration",
					fmt.Sprintf("Assignment at index %d has a daily_schedule without a required 'time' field", i),
				)
			}
		}

		// Check hourly schedule
		if assignment.HourlySchedule != nil {
			scheduleCount++
		}

		// Check run once schedule
		if assignment.RunOnceSchedule != nil {
			scheduleCount++

			// Validate run once schedule has required fields
			if assignment.RunOnceSchedule.Date.IsNull() || assignment.RunOnceSchedule.Date.IsUnknown() || assignment.RunOnceSchedule.Date.ValueString() == "" {
				diags.AddError(
					"Invalid Schedule Configuration",
					fmt.Sprintf("Assignment at index %d has a run_once_schedule without a required 'date' field", i),
				)
			}

			if assignment.RunOnceSchedule.Time.IsNull() || assignment.RunOnceSchedule.Time.IsUnknown() || assignment.RunOnceSchedule.Time.ValueString() == "" {
				diags.AddError(
					"Invalid Schedule Configuration",
					fmt.Sprintf("Assignment at index %d has a run_once_schedule without a required 'time' field", i),
				)
			}
		}

		// Exclusion assignments can have no schedule (they inherit from the include assignment)
		// But include assignments must have exactly one schedule
		if targetType == "exclusionGroupAssignmentTarget" {
			// Exclusion assignments can have 0 or 1 schedule
			if scheduleCount > 1 {
				diags.AddError(
					"Multiple Schedule Configurations",
					fmt.Sprintf("Exclusion assignment at index %d has %d schedule types defined. Exclusion assignments should have at most one schedule type", i, scheduleCount),
				)
			}
		} else {
			// Include assignments must have exactly one schedule
			if scheduleCount == 0 {
				diags.AddError(
					"Missing Schedule Configuration",
					fmt.Sprintf("Assignment at index %d must have exactly one schedule type (daily_schedule, hourly_schedule, or run_once_schedule)", i),
				)
			} else if scheduleCount > 1 {
				diags.AddError(
					"Multiple Schedule Configurations",
					fmt.Sprintf("Assignment at index %d has %d schedule types defined. Only one schedule type (daily_schedule, hourly_schedule, or run_once_schedule) should be specified per assignment", i, scheduleCount),
				)
			}
		}
	}

	if diags.HasError() {
		tflog.Error(ctx, "Assignment validation failed", map[string]interface{}{
			"errors": diags.Errors(),
		})
	} else {
		tflog.Debug(ctx, "Assignment validation passed")
	}

	return diags
}
