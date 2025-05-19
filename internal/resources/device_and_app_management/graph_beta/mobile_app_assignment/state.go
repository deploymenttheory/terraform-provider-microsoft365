package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, model MobileAppAssignmentResourceModel, assignment graphmodels.MobileAppAssignmentable) MobileAppAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return model
	}

	// Map the remote assignment to the Terraform model
	model.ID = state.StringPointerValue(assignment.GetId())
	model.Intent = state.EnumPtrToTypeString(assignment.GetIntent())
	model.Source = state.EnumPtrToTypeString(assignment.GetSource())
	model.SourceId = state.StringPointerValue(assignment.GetSourceId())

	// Map target and settings
	if target := assignment.GetTarget(); target != nil {
		model.Target = mapRemoteTargetToTerraform(target)
	}

	if settings := assignment.GetSettings(); settings != nil {
		model.Settings = mapRemoteSettingsToTerraform(settings)
	}

	tflog.Debug(ctx, "Finished mapping remote assignment to Terraform state")

	return model
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	target := AssignmentTargetResourceModel{
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

// mapRemoteSettingsToTerraform
func mapRemoteSettingsToTerraform(remoteSettings graphmodels.MobileAppAssignmentSettingsable) *MobileAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	var settings MobileAppAssignmentSettingsResourceModel

	switch v := remoteSettings.(type) {
	case *graphmodels.AndroidManagedStoreAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			AndroidManagedStore: mapAndroidManagedStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosLobAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			IosLob: mapIosLobSettingsToTerraform(v),
		}
	case *graphmodels.IosStoreAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			IosStore: mapIosStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosVppAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			IosVpp: mapIosVppSettingsToTerraform(v),
		}
	case *graphmodels.MacOsLobAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			MacOsLob: mapMacOsLobSettingsToTerraform(v),
		}
	case *graphmodels.MacOsVppAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			MacOsVpp: mapMacOsVppSettingsToTerraform(v),
		}
	case *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			MicrosoftStoreForBusiness: mapMicrosoftStoreSettingsToTerraform(v),
		}
	case *graphmodels.Win32LobAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			Win32Lob: mapWin32LobSettingsToTerraform(v),
		}
	case *graphmodels.WindowsAppXAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			WindowsAppX: mapWindowsAppXSettingsToTerraform(v),
		}
	case *graphmodels.WindowsUniversalAppXAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			WindowsUniversalAppX: mapWindowsUniversalAppXSettingsToTerraform(v),
		}
	case *graphmodels.WinGetAppAssignmentSettings:
		settings = MobileAppAssignmentSettingsResourceModel{
			WinGet: mapWinGetSettingsToTerraform(v),
		}
	default:
		return nil
	}

	return &settings
}

// mapAndroidManagedStoreSettingsToTerraform maps an Android managed store settings to a Terraform assignment settings
func mapAndroidManagedStoreSettingsToTerraform(remoteSettings *graphmodels.AndroidManagedStoreAppAssignmentSettings) *AndroidManagedStoreAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &AndroidManagedStoreAssignmentSettingsResourceModel{
		AndroidManagedStoreAppTrackIds: state.StringListToTypeList(remoteSettings.GetAndroidManagedStoreAppTrackIds()),
		AutoUpdateMode:                 state.EnumPtrToTypeString(remoteSettings.GetAutoUpdateMode()),
	}
}

// mapIosLobSettingsToTerraform maps an iOS LOB settings to a Terraform assignment settings
func mapIosLobSettingsToTerraform(remoteSettings *graphmodels.IosLobAppAssignmentSettings) *IosLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosLobAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosStoreSettingsToTerraform maps an iOS store settings to a Terraform assignment settings
func mapIosStoreSettingsToTerraform(remoteSettings *graphmodels.IosStoreAppAssignmentSettings) *IosStoreAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosStoreAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosVppSettingsToTerraform maps an iOS VPP settings to a Terraform assignment settings
func mapIosVppSettingsToTerraform(remoteSettings *graphmodels.IosVppAppAssignmentSettings) *IosVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosVppAppAssignmentSettingsResourceModel{
		IsRemovable:              state.BoolPtrToTypeBool(remoteSettings.GetIsRemovable()),
		PreventAutoAppUpdate:     state.BoolPtrToTypeBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceLicensing()),
		VpnConfigurationId:       types.StringPointerValue(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapMacOsLobSettingsToTerraform maps a macOS LOB settings to a Terraform assignment settings
func mapMacOsLobSettingsToTerraform(remoteSettings *graphmodels.MacOsLobAppAssignmentSettings) *MacOsLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MacOsLobAppAssignmentSettingsResourceModel{
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
	}
}

// mapMacOsVppSettingsToTerraform maps a macOS VPP settings to a Terraform assignment settings
func mapMacOsVppSettingsToTerraform(remoteSettings *graphmodels.MacOsVppAppAssignmentSettings) *MacOsVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MacOsVppAppAssignmentSettingsResourceModel{
		PreventAutoAppUpdate:     state.BoolPtrToTypeBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  state.BoolPtrToTypeBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: state.BoolPtrToTypeBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceLicensing()),
	}
}

// mapMicrosoftStoreSettingsToTerraform maps a Microsoft Store settings to a Terraform assignment settings
func mapMicrosoftStoreSettingsToTerraform(remoteSettings *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings) *MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWin32LobSettingsToTerraform maps a Win32 LOB settings to a Terraform assignment settings
func mapWin32LobSettingsToTerraform(remoteSettings *graphmodels.Win32LobAppAssignmentSettings) *Win32LobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	settings := &Win32LobAppAssignmentSettingsResourceModel{
		DeliveryOptimizationPriority: state.EnumPtrToTypeString(remoteSettings.GetDeliveryOptimizationPriority()),
		Notifications:                state.EnumPtrToTypeString(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		settings.InstallTimeSettings = &MobileAppInstallTimeSettingsResourceModel{
			DeadlineDateTime: state.TimeToString(installSettings.GetDeadlineDateTime()),
			StartDateTime:    state.TimeToString(installSettings.GetStartDateTime()),
			UseLocalTime:     state.BoolPtrToTypeBool(installSettings.GetUseLocalTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		settings.RestartSettings = &MobileAppAssignmentSettingsRestartResourceModel{
			CountdownDisplayBeforeRestart:     state.Int32PtrToTypeInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:              state.Int32PtrToTypeInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDuration: state.Int32PtrToTypeInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return settings
}

// mapWindowsAppXSettingsToTerraform maps a Windows AppX settings to a Terraform assignment settings
func mapWindowsAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsAppXAppAssignmentSettings) *WindowsAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &WindowsAppXAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWindowsUniversalAppXSettingsToTerraform maps a Windows Universal AppX settings to a Terraform assignment settings
func mapWindowsUniversalAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsUniversalAppXAppAssignmentSettings) *WindowsUniversalAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &WindowsUniversalAppXAssignmentSettingsResourceModel{
		UseDeviceContext: state.BoolPtrToTypeBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWinGetSettingsToTerraform maps a WinGet settings to a Terraform assignment settings
func mapWinGetSettingsToTerraform(remoteSettings *graphmodels.WinGetAppAssignmentSettings) *WinGetAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	winGetSettings := &WinGetAppAssignmentSettingsResourceModel{
		Notifications: state.EnumPtrToTypeString(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		winGetSettings.InstallTimeSettings = &WinGetAppInstallTimeSettingsResourceModel{
			UseLocalTime:     types.BoolPointerValue(installSettings.GetUseLocalTime()),
			DeadlineDateTime: state.TimeToString(installSettings.GetDeadlineDateTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		winGetSettings.RestartSettings = &WinGetAppRestartSettingsResourceModel{
			CountdownDisplayBeforeRestartInMinutes:     state.Int32PtrToTypeInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:                       state.Int32PtrToTypeInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDurationInMinutes: state.Int32PtrToTypeInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return winGetSettings
}
