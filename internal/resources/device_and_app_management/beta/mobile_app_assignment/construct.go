package graphBetaMobileAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructResource maps the Terraform schema to the SDK model
func ConstructResource(ctx context.Context, data *MobileAppAssignmentResourceModel) (*graphmodels.MobileAppAssignment, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewMobileAppAssignment()

	// Set Intent
	intentValue, err := graphmodels.ParseInstallIntent(data.Intent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("invalid intent: %s", err)
	}
	if intentValue != nil {
		intent, ok := intentValue.(*graphmodels.InstallIntent)
		if !ok {
			return nil, fmt.Errorf("unexpected type for intent: %T", intentValue)
		}
		requestBody.SetIntent(intent)
	}

	// Set Target
	target := graphmodels.NewAllLicensedUsersAssignmentTarget()
	if !data.Target.DeviceAndAppManagementAssignmentFilterID.IsNull() {
		filterID := data.Target.DeviceAndAppManagementAssignmentFilterID.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&filterID)
	}
	if !data.Target.DeviceAndAppManagementAssignmentFilterType.IsNull() {
		filterTypeValue, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(data.Target.DeviceAndAppManagementAssignmentFilterType.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid device and app management assignment filter type: %s", err)
		}
		if filterTypeValue != nil {
			filterType, ok := filterTypeValue.(*graphmodels.DeviceAndAppManagementAssignmentFilterType)
			if !ok {
				return nil, fmt.Errorf("unexpected type for filter type: %T", filterTypeValue)
			}
			target.SetDeviceAndAppManagementAssignmentFilterType(filterType)
		}
	}
	requestBody.SetTarget(target)

	// Set Settings
	settings := graphmodels.NewWinGetAppAssignmentSettings()
	if !data.Settings.Notifications.IsNull() {
		notificationsValue, err := graphmodels.ParseWinGetAppNotification(data.Settings.Notifications.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid notifications setting: %s", err)
		}
		if notificationsValue != nil {
			notifications, ok := notificationsValue.(*graphmodels.WinGetAppNotification)
			if !ok {
				return nil, fmt.Errorf("unexpected type for notifications: %T", notificationsValue)
			}
			settings.SetNotifications(notifications)
		}
	}

	// Set Restart Settings
	restartSettings := graphmodels.NewWinGetAppRestartSettings()
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
	installTimeSettings := graphmodels.NewWinGetAppInstallTimeSettings()
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
		sourceValue, err := graphmodels.ParseDeviceAndAppManagementAssignmentSource(data.Source.ValueString())
		if err != nil {
			return nil, fmt.Errorf("invalid source: %s", err)
		}
		if sourceValue != nil {
			source, ok := sourceValue.(*graphmodels.DeviceAndAppManagementAssignmentSource)
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

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
