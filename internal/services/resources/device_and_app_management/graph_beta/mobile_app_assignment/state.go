package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data MobileAppAssignmentResourceModel, assignment graphmodels.MobileAppAssignmentable) MobileAppAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(assignment.GetId())
	data.Intent = convert.GraphToFrameworkEnum(assignment.GetIntent())
	data.Source = convert.GraphToFrameworkEnum(assignment.GetSource())
	data.SourceId = convert.GraphToFrameworkString(assignment.GetSourceId())

	if target := assignment.GetTarget(); target != nil {
		data.Target = mapRemoteTargetToTerraform(target)
	}

	if settings := assignment.GetSettings(); settings != nil {
		data.Settings = mapRemoteSettingsToTerraform(settings)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	target := AssignmentTargetResourceModel{
		DeviceAndAppManagementAssignmentFilterId:   convert.GraphToFrameworkString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		DeviceAndAppManagementAssignmentFilterType: convert.GraphToFrameworkEnum(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		target.TargetType = types.StringValue("groupAssignment")
		target.GroupId = convert.GraphToFrameworkString(v.GetGroupId())
	case *graphmodels.ExclusionGroupAssignmentTarget:
		target.TargetType = types.StringValue("exclusionGroupAssignment")
		target.GroupId = convert.GraphToFrameworkString(v.GetGroupId())
	case *graphmodels.ConfigurationManagerCollectionAssignmentTarget:
		target.TargetType = types.StringValue("configurationManagerCollection")
		target.CollectionId = convert.GraphToFrameworkString(v.GetCollectionId())
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
		AndroidManagedStoreAppTrackIds: convert.GraphToFrameworkStringList(remoteSettings.GetAndroidManagedStoreAppTrackIds()),
		AutoUpdateMode:                 convert.GraphToFrameworkEnum(remoteSettings.GetAutoUpdateMode()),
	}
}

// mapIosLobSettingsToTerraform maps an iOS LOB settings to a Terraform assignment settings
func mapIosLobSettingsToTerraform(remoteSettings *graphmodels.IosLobAppAssignmentSettings) *IosLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosLobAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosStoreSettingsToTerraform maps an iOS store settings to a Terraform assignment settings
func mapIosStoreSettingsToTerraform(remoteSettings *graphmodels.IosStoreAppAssignmentSettings) *IosStoreAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosStoreAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosVppSettingsToTerraform maps an iOS VPP settings to a Terraform assignment settings
func mapIosVppSettingsToTerraform(remoteSettings *graphmodels.IosVppAppAssignmentSettings) *IosVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &IosVppAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventAutoAppUpdate:     convert.GraphToFrameworkBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceLicensing()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapMacOsLobSettingsToTerraform maps a macOS LOB settings to a Terraform assignment settings
func mapMacOsLobSettingsToTerraform(remoteSettings *graphmodels.MacOsLobAppAssignmentSettings) *MacOsLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MacOsLobAppAssignmentSettingsResourceModel{
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
	}
}

// mapMacOsVppSettingsToTerraform maps a macOS VPP settings to a Terraform assignment settings
func mapMacOsVppSettingsToTerraform(remoteSettings *graphmodels.MacOsVppAppAssignmentSettings) *MacOsVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MacOsVppAppAssignmentSettingsResourceModel{
		PreventAutoAppUpdate:     convert.GraphToFrameworkBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceLicensing()),
	}
}

// mapMicrosoftStoreSettingsToTerraform maps a Microsoft Store settings to a Terraform assignment settings
func mapMicrosoftStoreSettingsToTerraform(remoteSettings *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings) *MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel{
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWin32LobSettingsToTerraform maps a Win32 LOB settings to a Terraform assignment settings
func mapWin32LobSettingsToTerraform(remoteSettings *graphmodels.Win32LobAppAssignmentSettings) *Win32LobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	settings := &Win32LobAppAssignmentSettingsResourceModel{
		DeliveryOptimizationPriority: convert.GraphToFrameworkEnum(remoteSettings.GetDeliveryOptimizationPriority()),
		Notifications:                convert.GraphToFrameworkEnum(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		settings.InstallTimeSettings = &MobileAppInstallTimeSettingsResourceModel{
			DeadlineDateTime: convert.GraphToFrameworkTime(installSettings.GetDeadlineDateTime()),
			StartDateTime:    convert.GraphToFrameworkTime(installSettings.GetStartDateTime()),
			UseLocalTime:     convert.GraphToFrameworkBool(installSettings.GetUseLocalTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		settings.RestartSettings = &MobileAppAssignmentSettingsRestartResourceModel{
			CountdownDisplayBeforeRestart:     convert.GraphToFrameworkInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:              convert.GraphToFrameworkInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDuration: convert.GraphToFrameworkInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
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
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWindowsUniversalAppXSettingsToTerraform maps a Windows Universal AppX settings to a Terraform assignment settings
func mapWindowsUniversalAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsUniversalAppXAppAssignmentSettings) *WindowsUniversalAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &WindowsUniversalAppXAssignmentSettingsResourceModel{
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWinGetSettingsToTerraform maps a WinGet settings to a Terraform assignment settings
func mapWinGetSettingsToTerraform(remoteSettings *graphmodels.WinGetAppAssignmentSettings) *WinGetAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	winGetSettings := &WinGetAppAssignmentSettingsResourceModel{
		Notifications: convert.GraphToFrameworkEnum(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		winGetSettings.InstallTimeSettings = &WinGetAppInstallTimeSettingsResourceModel{
			UseLocalTime:     convert.GraphToFrameworkBool(installSettings.GetUseLocalTime()),
			DeadlineDateTime: convert.GraphToFrameworkTime(installSettings.GetDeadlineDateTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		winGetSettings.RestartSettings = &WinGetAppRestartSettingsResourceModel{
			CountdownDisplayBeforeRestartInMinutes:     convert.GraphToFrameworkInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:                       convert.GraphToFrameworkInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDurationInMinutes: convert.GraphToFrameworkInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return winGetSettings
}
