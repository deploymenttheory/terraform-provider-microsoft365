package graphBetaMobileAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAssignments constructs a list of MobileAppAssignments from the resource model.
func ConstructAssignments(ctx context.Context, assignments []MobileAppAssignmentResourceModel) []models.MobileAppAssignmentable {
	construct.DebugPrintStruct(ctx, "Constructing Mobile App assignments from model", assignments)

	var mobileAppAssignments []models.MobileAppAssignmentable

	for _, assignment := range assignments {
		mobileAppAssignment := models.NewMobileAppAssignment()

		// Set target
		target := ConstructAssignmentTarget(assignment.Target)
		mobileAppAssignment.SetTarget(target)

		// Set intent
		if !assignment.Intent.IsNull() {
			intent, err := models.ParseInstallIntent(assignment.Intent.ValueString())
			if err == nil && intent != nil {
				mobileAppAssignment.SetIntent(intent.(*models.InstallIntent))
			}
		}

		// Set settings
		settings := models.NewWinGetAppAssignmentSettings()

		// Set notifications
		if !assignment.Settings.Notifications.IsNull() {
			notifications, err := models.ParseWinGetAppNotification(assignment.Settings.Notifications.ValueString())
			if err == nil && notifications != nil {
				settings.SetNotifications(notifications.(*models.WinGetAppNotification))
			}
		}

		// Set restart settings if any fields are set
		if assignment.Settings.RestartSettings.GracePeriodInMinutes.ValueInt64() != 0 ||
			assignment.Settings.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt64() != 0 ||
			assignment.Settings.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt64() != 0 {

			restartSettings := models.NewWinGetAppRestartSettings()

			// Convert int64 to int32 and set pointers
			gracePeriod := int32(assignment.Settings.RestartSettings.GracePeriodInMinutes.ValueInt64())
			countdownBeforeRestart := int32(assignment.Settings.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt64())
			snoozeDuration := int32(assignment.Settings.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt64())

			restartSettings.SetGracePeriodInMinutes(&gracePeriod)
			restartSettings.SetCountdownDisplayBeforeRestartInMinutes(&countdownBeforeRestart)
			restartSettings.SetRestartNotificationSnoozeDurationInMinutes(&snoozeDuration)

			settings.SetRestartSettings(restartSettings)
		}

		// Set install time settings if any fields are set
		if assignment.Settings.InstallTimeSettings.UseLocalTime.ValueBool() ||
			assignment.Settings.InstallTimeSettings.DeadlineDateTime.ValueString() != "" {

			installTimeSettings := models.NewWinGetAppInstallTimeSettings()

			// Set UseLocalTime by taking a pointer of the bool
			useLocalTime := assignment.Settings.InstallTimeSettings.UseLocalTime.ValueBool()
			installTimeSettings.SetUseLocalTime(&useLocalTime)

			// Set DeadlineDateTime if it is not empty
			if deadline := assignment.Settings.InstallTimeSettings.DeadlineDateTime.ValueString(); deadline != "" {
				// Parse the string to time.Time
				parsedTime, err := time.Parse(time.RFC3339, deadline)
				if err != nil {
					// Handle parsing error
					fmt.Printf("Error parsing DeadlineDateTime: %v\n", err)
				} else {
					// Set the parsed time
					installTimeSettings.SetDeadlineDateTime(&parsedTime)
				}
			}

			settings.SetInstallTimeSettings(installTimeSettings)
		}

		mobileAppAssignment.SetSettings(settings)

		mobileAppAssignments = append(mobileAppAssignments, mobileAppAssignment)
	}

	return mobileAppAssignments
}

// ConstructAssignmentTarget constructs the assignment target based on the provided target model.
func ConstructAssignmentTarget(target AllLicensedUsersAssignmentTarget) models.DeviceAndAppManagementAssignmentTargetable {
	switch target.DeviceAndAppManagementAssignmentFilterType.ValueString() {
	case "allLicensedUsers":
		return models.NewAllLicensedUsersAssignmentTarget()
	case "allDevices":
		return models.NewAllDevicesAssignmentTarget()
	case "group":
		groupTarget := models.NewGroupAssignmentTarget()
		groupTarget.SetGroupId(target.DeviceAndAppManagementAssignmentFilterID.ValueStringPointer())

		// Set the filter ID if available
		if target.DeviceAndAppManagementAssignmentFilterID.ValueString() != "" {
			groupTarget.SetDeviceAndAppManagementAssignmentFilterId(target.DeviceAndAppManagementAssignmentFilterID.ValueStringPointer())
		}

		// Set the filter type if available
		if target.DeviceAndAppManagementAssignmentFilterType.ValueString() != "" {
			filterType, err := models.ParseDeviceAndAppManagementAssignmentFilterType(target.DeviceAndAppManagementAssignmentFilterType.ValueString())
			if err == nil && filterType != nil {
				groupTarget.SetDeviceAndAppManagementAssignmentFilterType(filterType.(*models.DeviceAndAppManagementAssignmentFilterType))
			}
		}
		return groupTarget
	default:
		return nil // This should not happen due to schema validation
	}
}
