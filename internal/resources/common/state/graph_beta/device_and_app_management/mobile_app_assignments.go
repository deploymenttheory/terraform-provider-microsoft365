package sharedStater

import (
	"context"
	"sort"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateMobileAppAssignment maps a remote assignment state to a Terraform assignment state
func StateMobileAppAssignment(ctx context.Context, assignments []sharedmodels.MobileAppAssignmentResourceModel, remoteAssignmentsResponse graphmodels.MobileAppAssignmentCollectionResponseable) {
	if remoteAssignmentsResponse == nil || remoteAssignmentsResponse.GetValue() == nil {
		tflog.Debug(ctx, "Remote assignments response is nil")
		return
	}

	remoteAssignments := remoteAssignmentsResponse.GetValue()
	assignments = assignments[:0]

	for _, remoteAssignment := range remoteAssignments {
		assignments = append(assignments, sharedmodels.MobileAppAssignmentResourceModel{
			Intent:   state.EnumPtrToTypeString(remoteAssignment.GetIntent()),
			Source:   state.EnumPtrToTypeString(remoteAssignment.GetSource()),
			SourceId: types.StringPointerValue(remoteAssignment.GetSourceId()),
			Target:   mapRemoteTargetToTerraform(remoteAssignment.GetTarget()),
			Settings: mapRemoteSettingsToTerraform(remoteAssignment.GetSettings()),
		})
	}

	// Sort assignments to maintain consistent order
	SortMobileAppAssignments(assignments)

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{})
}

// sortMobileAppAssignments sorts a slice of mobile app assignments
// The sort order is as follows:
// 1. First tier: Sort by intent alphabetically
// 2. Second tier: Within same intent, sort by target_type alphabetically
// 3. Third tier: Within same target_type, sort by group_id alphabetically
func SortMobileAppAssignments(assignments []sharedmodels.MobileAppAssignmentResourceModel) {
	sort.SliceStable(assignments, func(i, j int) bool {
		// First tier: Sort by intent alphabetically
		if assignments[i].Intent.ValueString() != assignments[j].Intent.ValueString() {
			return assignments[i].Intent.ValueString() < assignments[j].Intent.ValueString()
		}

		// Second tier: Within same intent, sort by target_type alphabetically
		if assignments[i].Target.TargetType.ValueString() != assignments[j].Target.TargetType.ValueString() {
			return assignments[i].Target.TargetType.ValueString() < assignments[j].Target.TargetType.ValueString()
		}

		// Third tier: Within same target_type, sort by group_id alphabetically
		if !assignments[i].Target.GroupId.IsNull() && !assignments[j].Target.GroupId.IsNull() {
			return assignments[i].Target.GroupId.ValueString() < assignments[j].Target.GroupId.ValueString()
		}

		return i < j
	})
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) sharedmodels.AssignmentTargetResourceModel {
	target := sharedmodels.AssignmentTargetResourceModel{
		DeviceAndAppManagementAssignmentFilterId:   types.StringPointerValue(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		DeviceAndAppManagementAssignmentFilterType: state.EnumPtrToTypeString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		target.TargetType = types.StringValue("groupAssignment")
		target.GroupId = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ExclusionGroupAssignmentTarget:
		target.TargetType = types.StringValue("exclusionGroupAssignment")
		target.GroupId = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ConfigurationManagerCollectionAssignmentTarget:
		target.TargetType = types.StringValue("configurationManagerCollection")
		target.CollectionId = types.StringPointerValue(v.GetCollectionId())
	case *graphmodels.AllDevicesAssignmentTarget:
		target.TargetType = types.StringValue("allDevices")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		target.TargetType = types.StringValue("allLicensedUsers")
	}

	return target
}

// mapRemoteSettingsToTerraform maps a remote assignment settings to a Terraform assignment settings
func mapRemoteSettingsToTerraform(remoteSettings graphmodels.MobileAppAssignmentSettingsable) sharedmodels.MobileAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{}
	}

	switch v := remoteSettings.(type) {
	case *graphmodels.AndroidManagedStoreAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			AndroidManagedStore: mapAndroidManagedStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosLobAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosLob: mapIosLobSettingsToTerraform(v),
		}
	case *graphmodels.IosStoreAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosStore: mapIosStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosVppAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosVpp: mapIosVppSettingsToTerraform(v),
		}
	case *graphmodels.MacOsLobAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MacOsLob: mapMacOsLobSettingsToTerraform(v),
		}
	case *graphmodels.MacOsVppAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MacOsVpp: mapMacOsVppSettingsToTerraform(v),
		}
	case *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MicrosoftStoreForBusiness: mapMicrosoftStoreSettingsToTerraform(v),
		}
	case *graphmodels.Win32LobAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			Win32Lob: mapWin32LobSettingsToTerraform(v),
		}
	case *graphmodels.WindowsAppXAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WindowsAppX: mapWindowsAppXSettingsToTerraform(v),
		}
	case *graphmodels.WindowsUniversalAppXAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WindowsUniversalAppX: mapWindowsUniversalAppXSettingsToTerraform(v),
		}
	case *graphmodels.WinGetAppAssignmentSettings:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WinGet: mapWinGetSettingsToTerraform(v),
		}
	default:
		return sharedmodels.MobileAppAssignmentSettingsResourceModel{}
	}
}

// mapAndroidManagedStoreSettingsToTerraform maps an Android managed store settings to a Terraform assignment settings
func mapAndroidManagedStoreSettingsToTerraform(remoteSettings *graphmodels.AndroidManagedStoreAppAssignmentSettings) *sharedmodels.AndroidManagedStoreAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.AndroidManagedStoreAssignmentSettingsResourceModel{
		AndroidManagedStoreAppTrackIds: state.StringListToTypeList(remoteSettings.GetAndroidManagedStoreAppTrackIds()),
		AutoUpdateMode:                 state.EnumPtrToTypeString(remoteSettings.GetAutoUpdateMode()),
	}
}

// mapIosLobSettingsToTerraform maps an iOS LOB settings to a Terraform assignment settings
func mapIosLobSettingsToTerraform(remoteSettings *graphmodels.IosLobAppAssignmentSettings) *sharedmodels.IosLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosLobAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosStoreSettingsToTerraform maps an iOS store settings to a Terraform assignment settings
func mapIosStoreSettingsToTerraform(remoteSettings *graphmodels.IosStoreAppAssignmentSettings) *sharedmodels.IosStoreAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosStoreAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosVppSettingsToTerraform maps an iOS VPP settings to a Terraform assignment settings
func mapIosVppSettingsToTerraform(remoteSettings *graphmodels.IosVppAppAssignmentSettings) *sharedmodels.IosVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosVppAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventAutoAppUpdate:     state.BoolPtrToTypeBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceLicensing()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapMacOsLobSettingsToTerraform maps a macOS LOB settings to a Terraform assignment settings
func mapMacOsLobSettingsToTerraform(remoteSettings *graphmodels.MacOsLobAppAssignmentSettings) *sharedmodels.MacOsLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MacOsLobAppAssignmentSettingsResourceModel{
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
	}
}

// mapMacOsVppSettingsToTerraform maps a macOS VPP settings to a Terraform assignment settings
func mapMacOsVppSettingsToTerraform(remoteSettings *graphmodels.MacOsVppAppAssignmentSettings) *sharedmodels.MacOsVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MacOsVppAppAssignmentSettingsResourceModel{
		PreventAutoAppUpdate:     state.BoolPtrToTypeBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceLicensing()),
	}
}

// mapMicrosoftStoreSettingsToTerraform maps a Microsoft Store settings to a Terraform assignment settings
func mapMicrosoftStoreSettingsToTerraform(remoteSettings *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings) *sharedmodels.MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWin32LobSettingsToTerraform maps a Win32 LOB settings to a Terraform assignment settings
func mapWin32LobSettingsToTerraform(remoteSettings *graphmodels.Win32LobAppAssignmentSettings) *sharedmodels.Win32LobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	settings := &sharedmodels.Win32LobAppAssignmentSettingsResourceModel{
		DeliveryOptimizationPriority: state.EnumPtrToTypeString(remoteSettings.GetDeliveryOptimizationPriority()),
		Notifications:                state.EnumPtrToTypeString(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		settings.InstallTimeSettings = &sharedmodels.MobileAppInstallTimeSettingsResourceModel{
			DeadlineDateTime: state.TimeToString(installSettings.GetDeadlineDateTime()),
			StartDateTime:    state.TimeToString(installSettings.GetStartDateTime()),
			UseLocalTime:     state.BoolPtrToTypeBool(installSettings.GetUseLocalTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		settings.RestartSettings = &sharedmodels.MobileAppAssignmentSettingsRestartResourceModel{
			CountdownDisplayBeforeRestart:     state.Int32PtrToTypeInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:              state.Int32PtrToTypeInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDuration: state.Int32PtrToTypeInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return settings
}

// mapWindowsAppXSettingsToTerraform maps a Windows AppX settings to a Terraform assignment settings
func mapWindowsAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsAppXAppAssignmentSettings) *sharedmodels.WindowsAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.WindowsAppXAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWindowsUniversalAppXSettingsToTerraform maps a Windows Universal AppX settings to a Terraform assignment settings
func mapWindowsUniversalAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsUniversalAppXAppAssignmentSettings) *sharedmodels.WindowsUniversalAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.WindowsUniversalAppXAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWinGetSettingsToTerraform maps a WinGet settings to a Terraform assignment settings
func mapWinGetSettingsToTerraform(remoteSettings *graphmodels.WinGetAppAssignmentSettings) *sharedmodels.WinGetAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	winGetSettings := &sharedmodels.WinGetAppAssignmentSettingsResourceModel{
		Notifications: state.EnumPtrToTypeString(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		winGetSettings.InstallTimeSettings = &sharedmodels.WinGetAppInstallTimeSettingsResourceModel{
			UseLocalTime:     types.BoolPointerValue(installSettings.GetUseLocalTime()),
			DeadlineDateTime: state.TimeToString(installSettings.GetDeadlineDateTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		winGetSettings.RestartSettings = &sharedmodels.WinGetAppRestartSettingsResourceModel{
			CountdownDisplayBeforeRestartInMinutes:     state.Int32PtrToTypeInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:                       state.Int32PtrToTypeInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDurationInMinutes: state.Int32PtrToTypeInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return winGetSettings
}