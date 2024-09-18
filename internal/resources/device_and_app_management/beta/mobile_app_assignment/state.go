package graphBetaMobileAppAssignment

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *MobileAppAssignmentResourceModel, remoteResource models.MobileAppAssignmentable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))

	data.Intent = state.EnumPtrToTypeString(remoteResource.GetIntent())

	if target := remoteResource.GetTarget(); target != nil {
		mapTargetToTerraform(data, target)
	}

	if settings := remoteResource.GetSettings(); settings != nil {
		mapSettingsToTerraform(ctx, data, settings)
	}

	data.Source = state.EnumPtrToTypeString(remoteResource.GetSource())
	data.SourceID = types.StringValue(state.StringPtrToString(remoteResource.GetSourceId()))

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

// mapTargetToTerraform maps the target object from the API to the Terraform state.
func mapTargetToTerraform(data *MobileAppAssignmentResourceModel, target models.DeviceAndAppManagementAssignmentTargetable) {
	switch t := target.(type) {
	case *models.AllLicensedUsersAssignmentTarget:
		data.Target.DeviceAndAppManagementAssignmentFilterType = types.StringValue("allLicensedUsers")
	case *models.AllDevicesAssignmentTarget:
		data.Target.DeviceAndAppManagementAssignmentFilterType = types.StringValue("allDevices")
	case *models.GroupAssignmentTarget:
		data.Target.DeviceAndAppManagementAssignmentFilterType = types.StringValue("group")
		data.Target.DeviceAndAppManagementAssignmentFilterID = types.StringValue(state.StringPtrToString(t.GetGroupId()))
	}
}

// mapSettingsToTerraform maps the settings object from the API to the Terraform state.
func mapSettingsToTerraform(ctx context.Context, data *MobileAppAssignmentResourceModel, settings models.MobileAppAssignmentSettingsable) {
	winGetSettings, ok := settings.(*models.WinGetAppAssignmentSettings)
	if !ok {
		tflog.Debug(ctx, "Settings type is not WinGetAppAssignmentSettings")
		return
	}

	if notifications := winGetSettings.GetNotifications(); notifications != nil {
		data.Settings.Notifications = types.StringValue(notifications.String())
	}

	if restartSettings := winGetSettings.GetRestartSettings(); restartSettings != nil {
		data.Settings.RestartSettings.GracePeriodInMinutes = types.Int64Value(int64(*restartSettings.GetGracePeriodInMinutes()))
		data.Settings.RestartSettings.CountdownDisplayBeforeRestartInMinutes = types.Int64Value(int64(*restartSettings.GetCountdownDisplayBeforeRestartInMinutes()))
		data.Settings.RestartSettings.RestartNotificationSnoozeDurationInMinutes = types.Int64Value(int64(*restartSettings.GetRestartNotificationSnoozeDurationInMinutes()))
	}

	if installTimeSettings := winGetSettings.GetInstallTimeSettings(); installTimeSettings != nil {
		data.Settings.InstallTimeSettings.UseLocalTime = types.BoolValue(*installTimeSettings.GetUseLocalTime())

		if deadline := installTimeSettings.GetDeadlineDateTime(); deadline != nil {
			data.Settings.InstallTimeSettings.DeadlineDateTime = types.StringValue(deadline.Format(time.RFC3339))
		}
	}
}
