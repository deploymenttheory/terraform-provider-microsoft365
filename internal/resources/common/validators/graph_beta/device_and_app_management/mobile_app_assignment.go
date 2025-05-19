package sharedValidators

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ValidateMobileAppAssignmentSettings validates the mobile app assignment settings across all assignments
func ValidateMobileAppAssignmentSettings(ctx context.Context, appType string, config []sharedmodels.MobileAppAssignmentResourceModel, graphClient *msgraphsdk.GraphServiceClient) error {

	// Track usage of special target types
	allDevicesCount := 0
	allLicensedUsersCount := 0
	firstAllDevicesIndex := -1
	firstAllLicensedUsersIndex := -1

	// Rule 0: Validate that group IDs are only used once across all assignments
	if err := validateUniqueGroupIds(config); err != nil {
		return err
	}

	for i, assignment := range config {
		// Rule 1: Validate app-type specific requirements
		if err := validateMobileAppAssignmentType(appType, i, assignment); err != nil {
			return err
		}

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

		// Rule 5: Validate required group_id for valid assignment target types
		if err := validateRequiredGroupId(i, assignment); err != nil {
			return err
		}

		// Rule 6: Validate that assignment filters exist in the system
		if err := validateAssignmentFilterExists(ctx, i, assignment, graphClient); err != nil {
			return err
		}

		// Rule 7: Validate that group IDs exist in the system
		if err := validateGroupIdExists(ctx, i, assignment, graphClient); err != nil {
			return err
		}
	}

	return nil
}

// validateMobileAppAssignmentType validates app-type specific requirements for assignments
func validateMobileAppAssignmentType(appType string, index int, assignment sharedmodels.MobileAppAssignmentResourceModel) error {
	// Special handling for WindowsStoreApp type
	if appType == "WindowsStoreApp" {
		// First check if settings exists
		if assignment.Settings == nil {
			return fmt.Errorf("assignment[%d] is missing required 'settings' field for application_type '%s'", index, appType)
		}

		// Then check if win_get exists
		if assignment.Settings.WinGet == nil {
			return fmt.Errorf("assignment[%d] is missing required 'settings.win_get' field for application_type '%s'", index, appType)
		}

		// Finally check notifications
		if assignment.Settings.WinGet.Notifications.IsNull() || assignment.Settings.WinGet.Notifications.IsUnknown() {
			return fmt.Errorf("assignment[%d] is missing required 'settings.win_get.notifications' field for application_type '%s'", index, appType)
		}
	}

	// Add similar checks for other application types as needed

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

// validateRequiredGroupId checks if group_id is provided when specific target types are used
func validateRequiredGroupId(index int, assignment sharedmodels.MobileAppAssignmentResourceModel) error {
	if !assignment.Target.TargetType.IsNull() {
		targetType := assignment.Target.TargetType.ValueString()

		// List of target types that require a group_id
		requiresGroupId := map[string]bool{
			"androidFotaDeployment":    true,
			"exclusionGroupAssignment": true,
			"groupAssignment":          true,
			// Not including "configurationManagerCollection" as it uses collectionId instead
		}

		if requiresGroupId[targetType] {
			// Check if group_id exists and is not empty
			if assignment.Target.GroupId.IsNull() || assignment.Target.GroupId.ValueString() == "" {
				return fmt.Errorf(
					"assignment[%d]: target_type '%s' requires a valid group_id to be specified",
					index, targetType,
				)
			}
		}

		// Special case for configurationManagerCollection which requires collectionId
		if targetType == "configurationManagerCollection" {
			if assignment.Target.CollectionId.IsNull() || assignment.Target.CollectionId.ValueString() == "" {
				return fmt.Errorf(
					"assignment[%d]: target_type 'configurationManagerCollection' requires a valid collection_id to be specified",
					index,
				)
			}
		}
	}

	return nil
}

// validateAssignmentFilterExists checks if the specified assignment filter exists in the system
func validateAssignmentFilterExists(ctx context.Context, index int, assignment sharedmodels.MobileAppAssignmentResourceModel, client *msgraphsdk.GraphServiceClient) error {
	// Skip validation if client is nil (for backward compatibility)
	if client == nil {
		return nil
	}

	// If a filter ID is specified, verify it exists
	if !assignment.Target.DeviceAndAppManagementAssignmentFilterId.IsNull() &&
		assignment.Target.DeviceAndAppManagementAssignmentFilterId.ValueString() != "" {
		filterId := assignment.Target.DeviceAndAppManagementAssignmentFilterId.ValueString()

		_, err := client.
			DeviceManagement().
			AssignmentFilters().
			ByDeviceAndAppManagementAssignmentFilterId(filterId).
			Get(ctx, nil)

		if err != nil {
			// Check if it's a "not found" error (404) from Graph API
			if odataErr, ok := err.(*odataerrors.ODataError); ok {
				if odataErr.ResponseStatusCode == 404 {
					return fmt.Errorf(
						"assignment[%d]: specified assignment filter ID '%s' was not found in the system",
						index, filterId,
					)
				}
			}
			// For any other errors, return a generic error message
			return fmt.Errorf(
				"assignment[%d]: error validating assignment filter ID '%s': %v",
				index, filterId, err,
			)
		}
	}

	return nil
}

// validateGroupIdExists checks if the specified group ID exists in the system
func validateGroupIdExists(ctx context.Context, index int, assignment sharedmodels.MobileAppAssignmentResourceModel, client *msgraphsdk.GraphServiceClient) error {
	// Skip validation if client is nil (for backward compatibility)
	if client == nil {
		return nil
	}

	// Check if this assignment requires a group ID validation
	if !assignment.Target.TargetType.IsNull() {
		targetType := assignment.Target.TargetType.ValueString()

		requiresGroupId := map[string]bool{
			"androidFotaDeployment":    true,
			"exclusionGroupAssignment": true,
			"groupAssignment":          true,
		}

		// Only validate group IDs for target types that require them
		if requiresGroupId[targetType] && !assignment.Target.GroupId.IsNull() &&
			assignment.Target.GroupId.ValueString() != "" {
			groupId := assignment.Target.GroupId.ValueString()

			_, err := client.Groups().
				ByGroupId(groupId).
				Get(ctx, nil)

			if err != nil {
				// Check if it's a "not found" error (404) from Graph API
				if odataErr, ok := err.(*odataerrors.ODataError); ok {
					if odataErr.ResponseStatusCode == 404 {
						return fmt.Errorf(
							"assignment[%d]: specified group ID '%s' for target type '%s' was not found in the system",
							index, groupId, targetType,
						)
					}
				}
				// For any other errors, return a generic error message
				return fmt.Errorf(
					"assignment[%d]: error validating group ID '%s': %v",
					index, groupId, err,
				)
			}
		}
	}

	return nil
}

// validateUniqueGroupIds checks that group IDs are only used once across all assignments
func validateUniqueGroupIds(config []sharedmodels.MobileAppAssignmentResourceModel) error {
	usedGroupIds := make(map[string]int)

	for i, assignment := range config {
		if assignment.Target.TargetType.IsNull() {
			continue
		}

		targetType := assignment.Target.TargetType.ValueString()

		if targetType == "groupAssignment" && !assignment.Target.GroupId.IsNull() && assignment.Target.GroupId.ValueString() != "" {
			groupId := assignment.Target.GroupId.ValueString()

			if prevIndex, exists := usedGroupIds[groupId]; exists {
				return fmt.Errorf(
					"assignment[%d]: group ID '%s' is already used in assignment[%d]. A group ID can only be targeted once across all assignments",
					i, groupId, prevIndex,
				)
			}

			usedGroupIds[groupId] = i
		}
	}

	return nil
}
