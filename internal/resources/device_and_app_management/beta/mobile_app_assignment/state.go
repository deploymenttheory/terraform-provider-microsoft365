package graphBetaMobileAppAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *MobileAppAssignmentResourceModel, remoteResource graphmodels.MobileAppAssignmentable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.Intent = state.EnumPtrToTypeString(remoteResource.GetIntent())
	data.SourceID = types.StringValue(state.StringPtrToString(remoteResource.GetSourceId()))
	data.Source = state.EnumPtrToTypeString(remoteResource.GetSource())

	// Map Target
	if target := remoteResource.GetTarget(); target != nil {
		data.Target = AllLicensedUsersAssignmentTargetResourceModel{
			DeviceAndAppManagementAssignmentFilterID:   types.StringValue(state.StringPtrToString(target.GetDeviceAndAppManagementAssignmentFilterId())),
			DeviceAndAppManagementAssignmentFilterType: state.EnumPtrToTypeString(target.GetDeviceAndAppManagementAssignmentFilterType()),
		}
	}

	// Map Settings
	if settings := remoteResource.GetSettings(); settings != nil {
		winGetSettings, ok := settings.(models.WinGetAppAssignmentSettingsable)
		if ok {
			data.Settings = WinGetAppAssignmentSettingsResourceModel{
				Notifications: state.EnumPtrToTypeString(winGetSettings.GetNotifications()),
			}

			// Map RestartSettings
			if restartSettings := winGetSettings.GetRestartSettings(); restartSettings != nil {
				data.Settings.RestartSettings = WinGetAppRestartSettingsResourceModel{
					GracePeriodInMinutes:                       state.Int32PtrToTypeInt64(restartSettings.GetGracePeriodInMinutes()),
					CountdownDisplayBeforeRestartInMinutes:     state.Int32PtrToTypeInt64(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
					RestartNotificationSnoozeDurationInMinutes: state.Int32PtrToTypeInt64(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
				}
			}

			// Map InstallTimeSettings
			if installTimeSettings := winGetSettings.GetInstallTimeSettings(); installTimeSettings != nil {
				data.Settings.InstallTimeSettings = WinGetAppInstallTimeSettingsResourceModel{
					UseLocalTime:     state.BoolPtrToTypeBool(installTimeSettings.GetUseLocalTime()),
					DeadlineDateTime: state.TimeToString(installTimeSettings.GetDeadlineDateTime()),
				}
			}
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
