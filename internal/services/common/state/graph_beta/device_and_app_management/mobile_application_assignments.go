package sharedStater

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateMobileAppAssignment maps remote assignments to a slice of assignment resource models
func StateMobileAppAssignment(ctx context.Context, assignments []sharedmodels.MobileAppAssignmentResourceModel, remoteAssignmentsResponse graphmodels.MobileAppAssignmentCollectionResponseable) []sharedmodels.MobileAppAssignmentResourceModel {
	if remoteAssignmentsResponse == nil || remoteAssignmentsResponse.GetValue() == nil || len(remoteAssignmentsResponse.GetValue()) == 0 {
		tflog.Debug(ctx, "Remote assignments response is empty")
		return []sharedmodels.MobileAppAssignmentResourceModel{}
	}

	remoteAssignments := remoteAssignmentsResponse.GetValue()

	newAssignments := make([]sharedmodels.MobileAppAssignmentResourceModel, 0, len(remoteAssignments))

	for _, remoteAssignment := range remoteAssignments {
		newAssignments = append(newAssignments, sharedmodels.MobileAppAssignmentResourceModel{
			Id:       convert.GraphToFrameworkString(remoteAssignment.GetId()),
			Intent:   convert.GraphToFrameworkEnum(remoteAssignment.GetIntent()),
			Source:   convert.GraphToFrameworkEnum(remoteAssignment.GetSource()),
			SourceId: convert.GraphToFrameworkString(remoteAssignment.GetSourceId()),
			Target:   mapRemoteTargetToTerraform(remoteAssignment.GetTarget()),
			Settings: mapRemoteSettingsToTerraform(remoteAssignment.GetSettings()),
		})
	}

	SortMobileAppAssignments(newAssignments)

	tflog.Debug(ctx, "Finished mapping remote assignments to Terraform state", map[string]interface{}{
		"assignment_count": len(newAssignments),
	})

	return newAssignments
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) sharedmodels.AssignmentTargetResourceModel {
	target := sharedmodels.AssignmentTargetResourceModel{
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
func mapRemoteSettingsToTerraform(remoteSettings graphmodels.MobileAppAssignmentSettingsable) *sharedmodels.MobileAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	var settings sharedmodels.MobileAppAssignmentSettingsResourceModel

	switch v := remoteSettings.(type) {
	case *graphmodels.AndroidManagedStoreAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			AndroidManagedStore: mapAndroidManagedStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosLobAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosLob: mapIosLobSettingsToTerraform(v),
		}
	case *graphmodels.IosStoreAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosStore: mapIosStoreSettingsToTerraform(v),
		}
	case *graphmodels.IosVppAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			IosVpp: mapIosVppSettingsToTerraform(v),
		}
	case *graphmodels.MacOsLobAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MacOsLob: mapMacOsLobSettingsToTerraform(v),
		}
	case *graphmodels.MacOsVppAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MacOsVpp: mapMacOsVppSettingsToTerraform(v),
		}
	case *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			MicrosoftStoreForBusiness: mapMicrosoftStoreSettingsToTerraform(v),
		}
	case *graphmodels.Win32LobAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			Win32Lob: mapWin32LobSettingsToTerraform(v),
		}
	case *graphmodels.WindowsAppXAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WindowsAppX: mapWindowsAppXSettingsToTerraform(v),
		}
	case *graphmodels.WindowsUniversalAppXAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WindowsUniversalAppX: mapWindowsUniversalAppXSettingsToTerraform(v),
		}
	case *graphmodels.WinGetAppAssignmentSettings:
		settings = sharedmodels.MobileAppAssignmentSettingsResourceModel{
			WinGet: mapWinGetSettingsToTerraform(v),
		}
	default:
		return nil
	}

	return &settings
}

// mapAndroidManagedStoreSettingsToTerraform maps an Android managed store settings to a Terraform assignment settings
func mapAndroidManagedStoreSettingsToTerraform(remoteSettings *graphmodels.AndroidManagedStoreAppAssignmentSettings) *sharedmodels.AndroidManagedStoreAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.AndroidManagedStoreAssignmentSettingsResourceModel{
		AndroidManagedStoreAppTrackIds: convert.GraphToFrameworkStringList(remoteSettings.GetAndroidManagedStoreAppTrackIds()),
		AutoUpdateMode:                 convert.GraphToFrameworkEnum(remoteSettings.GetAutoUpdateMode()),
	}
}

// mapIosLobSettingsToTerraform maps an iOS LOB settings to a Terraform assignment settings
func mapIosLobSettingsToTerraform(remoteSettings *graphmodels.IosLobAppAssignmentSettings) *sharedmodels.IosLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosLobAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosStoreSettingsToTerraform maps an iOS store settings to a Terraform assignment settings
func mapIosStoreSettingsToTerraform(remoteSettings *graphmodels.IosStoreAppAssignmentSettings) *sharedmodels.IosStoreAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosStoreAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapIosVppSettingsToTerraform maps an iOS VPP settings to a Terraform assignment settings
func mapIosVppSettingsToTerraform(remoteSettings *graphmodels.IosVppAppAssignmentSettings) *sharedmodels.IosVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.IosVppAppAssignmentSettingsResourceModel{
		IsRemovable:              convert.GraphToFrameworkBool(remoteSettings.GetIsRemovable()),
		PreventAutoAppUpdate:     convert.GraphToFrameworkBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceLicensing()),
		VpnConfigurationId:       convert.GraphToFrameworkString(remoteSettings.GetVpnConfigurationId()),
	}
}

// mapMacOsLobSettingsToTerraform maps a macOS LOB settings to a Terraform assignment settings
func mapMacOsLobSettingsToTerraform(remoteSettings *graphmodels.MacOsLobAppAssignmentSettings) *sharedmodels.MacOsLobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MacOsLobAppAssignmentSettingsResourceModel{
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
	}
}

// mapMacOsVppSettingsToTerraform maps a macOS VPP settings to a Terraform assignment settings
func mapMacOsVppSettingsToTerraform(remoteSettings *graphmodels.MacOsVppAppAssignmentSettings) *sharedmodels.MacOsVppAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MacOsVppAppAssignmentSettingsResourceModel{
		PreventAutoAppUpdate:     convert.GraphToFrameworkBool(remoteSettings.GetPreventAutoAppUpdate()),
		PreventManagedAppBackup:  convert.GraphToFrameworkBool(remoteSettings.GetPreventManagedAppBackup()),
		UninstallOnDeviceRemoval: convert.GraphToFrameworkBool(remoteSettings.GetUninstallOnDeviceRemoval()),
		UseDeviceLicensing:       convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceLicensing()),
	}
}

// mapMicrosoftStoreSettingsToTerraform maps a Microsoft Store settings to a Terraform assignment settings
func mapMicrosoftStoreSettingsToTerraform(remoteSettings *graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings) *sharedmodels.MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel{
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWin32LobSettingsToTerraform maps a Win32 LOB settings to a Terraform assignment settings
func mapWin32LobSettingsToTerraform(remoteSettings *graphmodels.Win32LobAppAssignmentSettings) *sharedmodels.Win32LobAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	settings := &sharedmodels.Win32LobAppAssignmentSettingsResourceModel{
		DeliveryOptimizationPriority: convert.GraphToFrameworkEnum(remoteSettings.GetDeliveryOptimizationPriority()),
		Notifications:                convert.GraphToFrameworkEnum(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		settings.InstallTimeSettings = &sharedmodels.MobileAppInstallTimeSettingsResourceModel{
			DeadlineDateTime: convert.GraphToFrameworkTime(installSettings.GetDeadlineDateTime()),
			StartDateTime:    convert.GraphToFrameworkTime(installSettings.GetStartDateTime()),
			UseLocalTime:     convert.GraphToFrameworkBool(installSettings.GetUseLocalTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		settings.RestartSettings = &sharedmodels.MobileAppAssignmentSettingsRestartResourceModel{
			CountdownDisplayBeforeRestart:     convert.GraphToFrameworkInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:              convert.GraphToFrameworkInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDuration: convert.GraphToFrameworkInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
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
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWindowsUniversalAppXSettingsToTerraform maps a Windows Universal AppX settings to a Terraform assignment settings
func mapWindowsUniversalAppXSettingsToTerraform(remoteSettings *graphmodels.WindowsUniversalAppXAppAssignmentSettings) *sharedmodels.WindowsUniversalAppXAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	return &sharedmodels.WindowsUniversalAppXAssignmentSettingsResourceModel{
		UseDeviceContext: convert.GraphToFrameworkBool(remoteSettings.GetUseDeviceContext()),
	}
}

// mapWinGetSettingsToTerraform maps a WinGet settings to a Terraform assignment settings
func mapWinGetSettingsToTerraform(remoteSettings *graphmodels.WinGetAppAssignmentSettings) *sharedmodels.WinGetAppAssignmentSettingsResourceModel {
	if remoteSettings == nil {
		return nil
	}

	winGetSettings := &sharedmodels.WinGetAppAssignmentSettingsResourceModel{
		Notifications: convert.GraphToFrameworkEnum(remoteSettings.GetNotifications()),
	}

	if installSettings := remoteSettings.GetInstallTimeSettings(); installSettings != nil {
		winGetSettings.InstallTimeSettings = &sharedmodels.WinGetAppInstallTimeSettingsResourceModel{
			UseLocalTime:     convert.GraphToFrameworkBool(installSettings.GetUseLocalTime()),
			DeadlineDateTime: convert.GraphToFrameworkTime(installSettings.GetDeadlineDateTime()),
		}
	}

	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		winGetSettings.RestartSettings = &sharedmodels.WinGetAppRestartSettingsResourceModel{
			CountdownDisplayBeforeRestartInMinutes:     convert.GraphToFrameworkInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			GracePeriodInMinutes:                       convert.GraphToFrameworkInt32(restartSettings.GetGracePeriodInMinutes()),
			RestartNotificationSnoozeDurationInMinutes: convert.GraphToFrameworkInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}
	}

	return winGetSettings
}
