package graphBetaWinGetApp

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignment constructs and returns a MobileAppsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *WinGetAppResourceModel) (deviceappmanagement.MobileAppsItemAssignPostRequestBodyable, error) {
	if data.Assignments == nil {
		return nil, fmt.Errorf("mobile app assignments configuration block is required")
	}

	tflog.Debug(ctx, "Starting mobile app assignment construction")

	requestBody := deviceappmanagement.NewMobileAppsItemAssignPostRequestBody()
	var assignments []graphmodels.MobileAppAssignmentable

	// Process each assignment
	for _, assignment := range data.Assignments.MobileAppAssignments {
		mobileAppAssignment := graphmodels.NewMobileAppAssignment()

		// Set Target
		target, err := constructTarget(ctx, assignment.Target)
		if err != nil {
			return nil, fmt.Errorf("error constructing target: %v", err)
		}
		mobileAppAssignment.SetTarget(target)

		// Set Intent
		if !assignment.Intent.IsNull() {
			intentValue, err := graphmodels.ParseInstallIntent(assignment.Intent.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing install intent: %v", err)
			}
			intent := intentValue.(*graphmodels.InstallIntent)
			mobileAppAssignment.SetIntent(intent)
		}

		// Set Source if specified
		if !assignment.Source.IsNull() {
			sourceValue, err := graphmodels.ParseDeviceAndAppManagementAssignmentSource(
				assignment.Source.ValueString())
			if err != nil {
				return nil, fmt.Errorf("error parsing source: %v", err)
			}
			source := sourceValue.(*graphmodels.DeviceAndAppManagementAssignmentSource)
			mobileAppAssignment.SetSource(source)
		}

		// Set SourceId if specified
		if !assignment.SourceId.IsNull() {
			sourceId := assignment.SourceId.ValueString()
			mobileAppAssignment.SetSourceId(&sourceId)
		}

		// Set Settings if present
		if assignment.Settings != nil {
			settings, err := constructSettings(ctx, assignment.Settings)
			if err != nil {
				return nil, fmt.Errorf("error constructing settings: %v", err)
			}
			mobileAppAssignment.SetSettings(settings)
		}

		assignments = append(assignments, mobileAppAssignment)
	}

	requestBody.SetMobileAppAssignments(assignments)

	if err := construct.DebugLogGraphObject(ctx, "Constructed mobile app assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructTarget constructs and returns a DeviceAndAppManagementAssignmentTargetable
func constructTarget(ctx context.Context, targetConfig sharedmodels.Target) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {

	setFilterProperties := func(target graphmodels.DeviceAndAppManagementAssignmentTargetable) error {
		if !targetConfig.DeviceAndAppManagementAssignmentFilterID.IsNull() {
			filterId := targetConfig.DeviceAndAppManagementAssignmentFilterID.ValueString()
			target.SetDeviceAndAppManagementAssignmentFilterId(&filterId)

			if !targetConfig.DeviceAndAppManagementAssignmentFilterType.IsNull() {
				filterType, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(
					targetConfig.DeviceAndAppManagementAssignmentFilterType.ValueString())
				if err != nil {
					tflog.Warn(ctx, "Failed to parse assignment filter type", map[string]interface{}{
						"error": err.Error(),
					})
				}
				filterTypeValue := filterType.(*graphmodels.DeviceAndAppManagementAssignmentFilterType)
				target.SetDeviceAndAppManagementAssignmentFilterType(filterTypeValue)
			}
		}
		return nil
	}

	if targetConfig.IsExclusionGroup.ValueBool() {
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		odataType := "#microsoft.graph.exclusionGroupAssignmentTarget"
		target.SetOdataType(&odataType)

		if !targetConfig.GroupID.IsNull() {
			groupId := targetConfig.GroupID.ValueString()
			target.SetGroupId(&groupId)
		}

		if err := setFilterProperties(target); err != nil {
			return nil, err
		}

		return target, nil
	}

	target := graphmodels.NewGroupAssignmentTarget()
	odataType := "#microsoft.graph.groupAssignmentTarget"
	target.SetOdataType(&odataType)

	if !targetConfig.GroupID.IsNull() {
		groupId := targetConfig.GroupID.ValueString()
		target.SetGroupId(&groupId)
	}

	if err := setFilterProperties(target); err != nil {
		return nil, err
	}

	return target, nil
}

// constructSettings constructs and returns a WinGetAppAssignmentSettings
func constructSettings(ctx context.Context, settingsConfig *sharedmodels.WinGetAppAssignmentSettings) (graphmodels.MobileAppAssignmentSettingsable, error) {
	settings := graphmodels.NewWinGetAppAssignmentSettings()

	// Set @odata.type for WinGet settings
	odataType := "#microsoft.graph.winGetAppAssignmentSettings"
	settings.SetOdataType(&odataType)

	// Set notifications
	if !settingsConfig.Notifications.IsNull() {
		notificationValue, err := graphmodels.ParseWinGetAppNotification(
			settingsConfig.Notifications.ValueString())
		if err != nil {
			tflog.Warn(ctx, "Failed to parse notification type", map[string]interface{}{
				"error": err.Error(),
			})
		}
		notification := notificationValue.(*graphmodels.WinGetAppNotification)
		settings.SetNotifications(notification)
	}

	// Set install time settings if present
	if settingsConfig.InstallTimeSettings != nil {
		installTimeSettings := graphmodels.NewWinGetAppInstallTimeSettings()

		if !settingsConfig.InstallTimeSettings.UseLocalTime.IsNull() {
			useLocalTime := settingsConfig.InstallTimeSettings.UseLocalTime.ValueBool()
			installTimeSettings.SetUseLocalTime(&useLocalTime)
		}

		if !settingsConfig.InstallTimeSettings.DeadlineDateTime.IsNull() {
			deadlineStr := settingsConfig.InstallTimeSettings.DeadlineDateTime.ValueString()
			if deadline, err := time.Parse(time.RFC3339, deadlineStr); err == nil {
				installTimeSettings.SetDeadlineDateTime(&deadline)
			}
		}

		settings.SetInstallTimeSettings(installTimeSettings)
	}

	// Set restart settings if present
	if settingsConfig.RestartSettings != nil {
		restartSettings := graphmodels.NewWinGetAppRestartSettings()

		if !settingsConfig.RestartSettings.GracePeriodInMinutes.IsNull() {
			gracePeriod := int32(settingsConfig.RestartSettings.GracePeriodInMinutes.ValueInt64())
			restartSettings.SetGracePeriodInMinutes(&gracePeriod)
		}

		if !settingsConfig.RestartSettings.CountdownDisplayBeforeRestartInMinutes.IsNull() {
			countdown := int32(settingsConfig.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt64())
			restartSettings.SetCountdownDisplayBeforeRestartInMinutes(&countdown)
		}

		if !settingsConfig.RestartSettings.RestartNotificationSnoozeDurationInMinutes.IsNull() {
			snooze := int32(settingsConfig.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt64())
			restartSettings.SetRestartNotificationSnoozeDurationInMinutes(&snooze)
		}

		settings.SetRestartSettings(restartSettings)
	}

	return settings, nil
}
