package graphBetaWinGetApp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentStateToTerraform maps the remote mobile app assignments to the Terraform state
func MapRemoteAssignmentStateToTerraform(ctx context.Context, assignment *sharedmodels.MobileAppAssignmentResourceModel, remoteAssignmentsResponse graphmodels.MobileAppAssignmentCollectionResponseable) {
	if remoteAssignmentsResponse == nil || remoteAssignmentsResponse.GetValue() == nil {
		assignment.MobileAppAssignments = nil
		tflog.Debug(ctx, "No assignments found in remote resource")
		return
	}

	remoteAssignments := remoteAssignmentsResponse.GetValue()
	tflog.Debug(ctx, "Mapping assignments from remote resource to Terraform state", map[string]interface{}{
		"assignmentsCount": len(remoteAssignments),
	})

	mappedAssignments := make([]sharedmodels.MobileAppAssignment, len(remoteAssignments))

	for i, remoteAssignment := range remoteAssignments {
		mappedAssignment := sharedmodels.MobileAppAssignment{}

		if target := remoteAssignment.GetTarget(); target != nil {
			mappedAssignment.Target = MapTargetToTerraform(ctx, target)
		}

		if intent := remoteAssignment.GetIntent(); intent != nil {
			mappedAssignment.Intent = state.EnumPtrToTypeString(intent)
		}

		if source := remoteAssignment.GetSource(); source != nil {
			mappedAssignment.Source = state.EnumPtrToTypeString(source)
		}

		if sourceId := remoteAssignment.GetSourceId(); sourceId != nil {
			mappedAssignment.SourceId = types.StringValue(*sourceId)
		}

		if settings := remoteAssignment.GetSettings(); settings != nil {
			mappedAssignment.Settings = MapSettingsToTerraform(ctx, settings)
		}

		mappedAssignments[i] = mappedAssignment
	}

	assignment.MobileAppAssignments = mappedAssignments
	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"mappedAssignmentsCount": len(mappedAssignments),
	})
}

func MapTargetToTerraform(ctx context.Context, target graphmodels.DeviceAndAppManagementAssignmentTargetable) sharedmodels.Target {
	targetModel := sharedmodels.Target{}

	if groupTarget, ok := target.(*graphmodels.GroupAssignmentTarget); ok {
		targetModel.GroupID = types.StringValue(state.StringPtrToString(groupTarget.GetGroupId()))
		targetModel.IsExclusionGroup = types.BoolValue(false)
		if filterId := groupTarget.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
			targetModel.DeviceAndAppManagementAssignmentFilterID = types.StringValue(*filterId)
		}
		if filterType := groupTarget.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
			targetModel.DeviceAndAppManagementAssignmentFilterType = state.EnumPtrToTypeString(filterType)
		}
	} else if exclusionTarget, ok := target.(*graphmodels.ExclusionGroupAssignmentTarget); ok {
		targetModel.GroupID = types.StringValue(state.StringPtrToString(exclusionTarget.GetGroupId()))
		targetModel.IsExclusionGroup = types.BoolValue(true)
		if filterId := exclusionTarget.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil {
			targetModel.DeviceAndAppManagementAssignmentFilterID = types.StringValue(*filterId)
		}
		if filterType := exclusionTarget.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
			targetModel.DeviceAndAppManagementAssignmentFilterType = state.EnumPtrToTypeString(filterType)
		}
	}

	return targetModel
}

func MapSettingsToTerraform(ctx context.Context, settings graphmodels.MobileAppAssignmentSettingsable) *sharedmodels.WinGetAppAssignmentSettings {
	if settings == nil {
		return nil
	}

	mappedSettings := &sharedmodels.WinGetAppAssignmentSettings{}

	if winGetSettings, ok := settings.(*graphmodels.WinGetAppAssignmentSettings); ok {
		mappedSettings.Notifications = state.EnumPtrToTypeString(winGetSettings.GetNotifications())

		if installTimeSettings := winGetSettings.GetInstallTimeSettings(); installTimeSettings != nil {
			mappedSettings.InstallTimeSettings = &sharedmodels.WinGetAppInstallTimeSettings{
				UseLocalTime:     state.BoolPtrToTypeBool(installTimeSettings.GetUseLocalTime()),
				DeadlineDateTime: state.TimeToString(installTimeSettings.GetDeadlineDateTime()),
			}
		}

		if restartSettings := winGetSettings.GetRestartSettings(); restartSettings != nil {
			mappedSettings.RestartSettings = &sharedmodels.WinGetAppRestartSettings{
				GracePeriodInMinutes:                       state.Int32PtrToTypeInt64(restartSettings.GetGracePeriodInMinutes()),
				CountdownDisplayBeforeRestartInMinutes:     state.Int32PtrToTypeInt64(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
				RestartNotificationSnoozeDurationInMinutes: state.Int32PtrToTypeInt64(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
			}
		}
	}

	return mappedSettings
}
