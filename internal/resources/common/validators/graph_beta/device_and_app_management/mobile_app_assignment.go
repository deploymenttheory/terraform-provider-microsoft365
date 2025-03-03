package sharedValidators

import (
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
)

// ValidateMobileAppAssignmentSettings validates the mobile app assignment settings across all assignments
func ValidateMobileAppAssignmentSettings(config []sharedmodels.MobileAppAssignmentResourceModel) error {

	// Rule 1: Validate assignment ordering
	if err := validateAssignmentOrdering(config); err != nil {
		return err
	}

	// Track usage of special target types
	allDevicesCount := 0
	allLicensedUsersCount := 0
	firstAllDevicesIndex := -1
	firstAllLicensedUsersIndex := -1

	for i, assignment := range config {
		// Rule 2: Validate install time settings based on intent
		if err := validateInstallTimeSettings(i, assignment); err != nil {
			return err
		}

		// Rule 3: Validate special target types usage
		if err := validateSpecialTargetTypes(i, assignment, &allDevicesCount, &allLicensedUsersCount, &firstAllDevicesIndex, &firstAllLicensedUsersIndex); err != nil {
			return err
		}

		// Rule 4: Validate restart settings relationships
		if err := validateRestartSettings(i, assignment); err != nil {
			return err
		}
	}

	return nil
}

// validateAssignmentOrdering ensures assignments follow the required ordering:
// 1. First tier: Sort by intent alphabetically
// 2. Second tier: Within same intent, sort by target_type alphabetically
// 3. Third tier: Within same target_type, sort by group_id alphabetically
func validateAssignmentOrdering(config []sharedmodels.MobileAppAssignmentResourceModel) error {
	if len(config) <= 1 {
		return nil // No ordering needed for 0 or 1 assignments
	}

	for i := 0; i < len(config)-1; i++ {
		current := config[i]
		next := config[i+1]

		// Get intent values, treating null as empty string for comparison
		currentIntent := ""
		nextIntent := ""
		if !current.Intent.IsNull() {
			currentIntent = current.Intent.ValueString()
		}
		if !next.Intent.IsNull() {
			nextIntent = next.Intent.ValueString()
		}

		// Compare intents (First tier)
		if currentIntent > nextIntent {
			return fmt.Errorf(
				"invalid mobile app assignment ordering between index %d and %d: intent '%s' must come before '%s'",
				i, i+1, nextIntent, currentIntent,
			)
		}

		// If intents are equal, compare target_types (Second tier)
		if currentIntent == nextIntent {
			currentTargetType := ""
			nextTargetType := ""
			if !current.Target.TargetType.IsNull() {
				currentTargetType = current.Target.TargetType.ValueString()
			}
			if !next.Target.TargetType.IsNull() {
				nextTargetType = next.Target.TargetType.ValueString()
			}

			if currentTargetType > nextTargetType {
				return fmt.Errorf(
					"invalid mobile app assignment ordering between index %d and %d: for intent '%s', target_type '%s' must come before '%s'",
					i, i+1, currentIntent, nextTargetType, currentTargetType,
				)
			}

			// If target_types are equal, compare group_ids (Third tier)
			if currentTargetType == nextTargetType {
				currentGroupID := ""
				nextGroupID := ""
				if !current.Target.GroupId.IsNull() {
					currentGroupID = current.Target.GroupId.ValueString()
				}
				if !next.Target.GroupId.IsNull() {
					nextGroupID = next.Target.GroupId.ValueString()
				}

				if currentGroupID > nextGroupID {
					return fmt.Errorf(
						"invalid mobile app assignment ordering between index %d and %d: for intent '%s' and target_type '%s', group_id '%s' must come before '%s'",
						i, i+1, currentIntent, currentTargetType, nextGroupID, currentGroupID,
					)
				}
			}
		}
	}

	return nil
}

// validateInstallTimeSettings checks if install_time_settings is set when intent is "available"
func validateInstallTimeSettings(index int, assignment sharedmodels.MobileAppAssignmentResourceModel) error {
	if !assignment.Intent.IsNull() && assignment.Intent.ValueString() == "available" {
		if assignment.Settings != nil {
			if assignment.Settings.WinGet != nil && assignment.Settings.WinGet.InstallTimeSettings != nil {
				return fmt.Errorf(
					"assignment[%d]: install_time_settings cannot be set when intent is 'available'",
					index,
				)
			}
			if assignment.Settings.Win32Lob != nil && assignment.Settings.Win32Lob.InstallTimeSettings != nil {
				return fmt.Errorf(
					"assignment[%d]: install_time_settings cannot be set when intent is 'available'",
					index,
				)
			}
			if assignment.Settings.Win32Catalog != nil && assignment.Settings.Win32Catalog.InstallTimeSettings != nil {
				return fmt.Errorf(
					"assignment[%d]: install_time_settings cannot be set when intent is 'available'",
					index,
				)
			}
		}
	}
	return nil
}

// validateSpecialTargetTypes checks if special target types are used more than once
func validateSpecialTargetTypes(index int, assignment sharedmodels.MobileAppAssignmentResourceModel,
	allDevicesCount, allLicensedUsersCount *int,
	firstAllDevicesIndex, firstAllLicensedUsersIndex *int) error {

	if !assignment.Target.TargetType.IsNull() {
		targetType := assignment.Target.TargetType.ValueString()

		switch targetType {
		case "allDevices":
			if *allDevicesCount == 0 {
				*firstAllDevicesIndex = index
			}
			*allDevicesCount++
			if *allDevicesCount > 1 {
				return fmt.Errorf(
					"assignment[%d]: target_type 'allDevices' can only be used once across all Intune app assignments. Already used in assignment[%d]",
					index, *firstAllDevicesIndex,
				)
			}
		case "allLicensedUsers":
			if *allLicensedUsersCount == 0 {
				*firstAllLicensedUsersIndex = index
			}
			*allLicensedUsersCount++
			if *allLicensedUsersCount > 1 {
				return fmt.Errorf(
					"assignment[%d]: target_type 'allLicensedUsers' can only be used once across all Intune app assignments. Already used in assignment[%d]",
					index, *firstAllLicensedUsersIndex,
				)
			}
		}
	}
	return nil
}

// validateRestartSettings validates the relationships between restart timing settings
func validateRestartSettings(index int, assignment sharedmodels.MobileAppAssignmentResourceModel) error {
	if assignment.Settings == nil {
		return nil
	}

	// Check WinGet restart settings
	if assignment.Settings.WinGet != nil && assignment.Settings.WinGet.RestartSettings != nil {
		rs := assignment.Settings.WinGet.RestartSettings
		if rs.GracePeriodInMinutes.IsNull() || rs.CountdownDisplayBeforeRestartInMinutes.IsNull() || rs.RestartNotificationSnoozeDurationInMinutes.IsNull() {
			return nil
		}

		gracePeriod := rs.GracePeriodInMinutes.ValueInt32()
		countdown := rs.CountdownDisplayBeforeRestartInMinutes.ValueInt32()
		snooze := rs.RestartNotificationSnoozeDurationInMinutes.ValueInt32()

		// Validate countdown must be less than grace period
		if countdown > gracePeriod {
			return fmt.Errorf(
				"assignment[%d]: countdown_display_before_restart_in_minutes (%d) must be less than or equal to grace_period_in_minutes (%d)",
				index, countdown, gracePeriod,
			)
		}

		// Special case: when grace period equals countdown
		if gracePeriod == countdown {
			if snooze != 1 {
				return fmt.Errorf(
					"assignment[%d]: when grace_period_in_minutes equals countdown_display_before_restart_in_minutes, restart_notification_snooze_duration_in_minutes must be 1",
					index,
				)
			}
			return nil
		}

		// Calculate maximum allowed snooze duration
		maxSnooze := (gracePeriod - countdown) / 2
		if snooze > maxSnooze {
			return fmt.Errorf(
				"assignment[%d]: restart_notification_snooze_duration_in_minutes (%d) cannot be more than half the difference between grace_period_in_minutes and countdown_display_before_restart_in_minutes (%d)",
				index, snooze, maxSnooze,
			)
		}
	}

	// Same validation for Win32Lob
	if assignment.Settings.Win32Lob != nil && assignment.Settings.Win32Lob.RestartSettings != nil {
		rs := assignment.Settings.Win32Lob.RestartSettings
		if rs.GracePeriodInMinutes.IsNull() || rs.CountdownDisplayBeforeRestart.IsNull() || rs.RestartNotificationSnoozeDuration.IsNull() {
			return nil
		}

		gracePeriod := rs.GracePeriodInMinutes.ValueInt32()
		countdown := rs.CountdownDisplayBeforeRestart.ValueInt32()
		snooze := rs.RestartNotificationSnoozeDuration.ValueInt32()

		if countdown > gracePeriod {
			return fmt.Errorf(
				"assignment[%d]: countdown_display_before_restart (%d) must be less than or equal to grace_period_in_minutes (%d)",
				index, countdown, gracePeriod,
			)
		}

		if gracePeriod == countdown {
			if snooze != 1 {
				return fmt.Errorf(
					"assignment[%d]: when grace_period_in_minutes equals countdown_display_before_restart, restart_notification_snooze_duration must be 1",
					index,
				)
			}
			return nil
		}

		maxSnooze := (gracePeriod - countdown) / 2
		if snooze > maxSnooze {
			return fmt.Errorf(
				"assignment[%d]: restart_notification_snooze_duration (%d) cannot be more than half the difference between grace_period_in_minutes and countdown_display_before_restart (%d)",
				index, snooze, maxSnooze,
			)
		}
	}

	// Same validation for Win32Catalog
	if assignment.Settings.Win32Catalog != nil && assignment.Settings.Win32Catalog.RestartSettings != nil {
		rs := assignment.Settings.Win32Catalog.RestartSettings
		if rs.GracePeriodInMinutes.IsNull() || rs.CountdownDisplayBeforeRestart.IsNull() || rs.RestartNotificationSnoozeDuration.IsNull() {
			return nil
		}

		gracePeriod := rs.GracePeriodInMinutes.ValueInt32()
		countdown := rs.CountdownDisplayBeforeRestart.ValueInt32()
		snooze := rs.RestartNotificationSnoozeDuration.ValueInt32()

		if countdown > gracePeriod {
			return fmt.Errorf(
				"assignment[%d]: countdown_display_before_restart (%d) must be less than or equal to grace_period_in_minutes (%d)",
				index, countdown, gracePeriod,
			)
		}

		if gracePeriod == countdown {
			if snooze != 1 {
				return fmt.Errorf(
					"assignment[%d]: when grace_period_in_minutes equals countdown_display_before_restart, restart_notification_snooze_duration must be 1",
					index,
				)
			}
			return nil
		}

		maxSnooze := (gracePeriod - countdown) / 2
		if snooze > maxSnooze {
			return fmt.Errorf(
				"assignment[%d]: restart_notification_snooze_duration (%d) cannot be more than half the difference between grace_period_in_minutes and countdown_display_before_restart (%d)",
				index, snooze, maxSnooze,
			)
		}
	}

	return nil
}
