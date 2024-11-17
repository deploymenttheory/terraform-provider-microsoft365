package graphBetaMobileAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructResource maps the Terraform schema to the SDK model
func ConstructResource(ctx context.Context, data *MobileAppAssignmentResourceModel) (*models.MobileAppAssignment, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewMobileAppAssignment()

	// Set Intent
	intentValue, err := models.ParseInstallIntent(data.Intent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid intent: %s", err)
	}
	if intentValue != nil {
		intent, ok := intentValue.(*models.InstallIntent)
		if !ok {
			return nil, fmt.Errorf("unexpected type for intent: %T", intentValue)
		}
		requestBody.SetIntent(intent)
	}

	// Set Target
	target := models.NewAllLicensedUsersAssignmentTarget()
	if !data.Target.DeviceAndAppManagementAssignmentFilterID.IsNull() {
		filterID := data.Target.DeviceAndAppManagementAssignmentFilterID.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&filterID)
	}
	if !data.Target.DeviceAndAppManagementAssignmentFilterType.IsNull() {
		filterTypeValue, err := models.ParseDeviceAndAppManagementAssignmentFilterType(data.Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid device and app management assignment filter type: %s", err)
		}
		if filterTypeValue != nil {
			filterType, ok := filterTypeValue.(*models.DeviceAndAppManagementAssignmentFilterType)
			if !ok {
				return nil, fmt.Errorf("unexpected type for filter type: %T", filterTypeValue)
			}
			target.SetDeviceAndAppManagementAssignmentFilterType(filterType)
		}
	}
	requestBody.SetTarget(target)

	// Set Settings
	settings := models.NewWinGetAppAssignmentSettings()
	if !data.Settings.Notifications.IsNull() {
		notificationsValue, err := models.ParseWinGetAppNotification(data.Settings.Notifications.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid notifications setting: %s", err)
		}
		if notificationsValue != nil {
			notifications, ok := notificationsValue.(*models.WinGetAppNotification)
			if !ok {
				return nil, fmt.Errorf("unexpected type for notifications: %T", notificationsValue)
			}
			settings.SetNotifications(notifications)
		}
	}

	// Set Restart Settings
	restartSettings := models.NewWinGetAppRestartSettings()
	if !data.Settings.RestartSettings.GracePeriodInMinutes.IsNull() {
		gracePeriod := int32(data.Settings.RestartSettings.GracePeriodInMinutes.ValueInt64())
		restartSettings.SetGracePeriodInMinutes(&gracePeriod)
	}
	if !data.Settings.RestartSettings.CountdownDisplayBeforeRestartInMinutes.IsNull() {
		countdown := int32(data.Settings.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt64())
		restartSettings.SetCountdownDisplayBeforeRestartInMinutes(&countdown)
	}
	if !data.Settings.RestartSettings.RestartNotificationSnoozeDurationInMinutes.IsNull() {
		snoozeDuration := int32(data.Settings.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt64())
		restartSettings.SetRestartNotificationSnoozeDurationInMinutes(&snoozeDuration)
	}
	settings.SetRestartSettings(restartSettings)

	// Set Install Time Settings
	installTimeSettings := models.NewWinGetAppInstallTimeSettings()
	if !data.Settings.InstallTimeSettings.UseLocalTime.IsNull() {
		useLocalTime := data.Settings.InstallTimeSettings.UseLocalTime.ValueBool()
		installTimeSettings.SetUseLocalTime(&useLocalTime)
	}
	if !data.Settings.InstallTimeSettings.DeadlineDateTime.IsNull() {
		deadlineDateTimeStr := data.Settings.InstallTimeSettings.DeadlineDateTime.ValueString()
		deadlineDateTime, err := time.Parse(time.RFC3339, deadlineDateTimeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline date time format: %s", err)
		}
		installTimeSettings.SetDeadlineDateTime(&deadlineDateTime)
	}
	settings.SetInstallTimeSettings(installTimeSettings)

	requestBody.SetSettings(settings)

	if !data.Source.IsNull() {
		sourceValue, err := models.ParseDeviceAndAppManagementAssignmentSource(data.Source.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid source: %s", err)
		}
		if sourceValue != nil {
			source, ok := sourceValue.(*models.DeviceAndAppManagementAssignmentSource)
			if !ok {
				return nil, fmt.Errorf("unexpected type for source: %T", sourceValue)
			}
			requestBody.SetSource(source)
		}
	}

	if !data.SourceID.IsNull() {
		sourceID := data.SourceID.ValueString()
		requestBody.SetSourceId(&sourceID)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
