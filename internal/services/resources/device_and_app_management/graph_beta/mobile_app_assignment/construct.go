package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"

	// validators "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructMobileAppAssignment constructs and returns a MobileAppAssignment
func ConstructMobileAppAssignment(ctx context.Context, data MobileAppAssignmentResourceModel) (graphmodels.MobileAppAssignmentable, error) {
	tflog.Debug(ctx, "Starting mobile app assignment construction")

	assignment := graphmodels.NewMobileAppAssignment()

	// Set Intent
	if !data.Intent.IsNull() {
		intentValue, err := graphmodels.ParseInstallIntent(data.Intent.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing install intent: %v", err)
		}
		assignment.SetIntent(intentValue.(*graphmodels.InstallIntent))
	}

	// Set Target
	target, err := constructAssignmentTarget(ctx, &data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing mobile app assignment target: %v", err)
	}
	assignment.SetTarget(target)

	// Set Source
	if !data.Source.IsNull() {
		sourceValue, err := graphmodels.ParseDeviceAndAppManagementAssignmentSource(data.Source.ValueString())
		if err != nil {
			return nil, fmt.Errorf("error parsing source: %v", err)
		}
		assignment.SetSource(sourceValue.(*graphmodels.DeviceAndAppManagementAssignmentSource))
	}

	// Set SourceId
	if !data.SourceId.IsNull() {
		id := data.SourceId.ValueString()
		assignment.SetSourceId(&id)
	}

	// Set Settings
	if data.Settings != nil {
		settings, err := constructMobileAppAssignmentSettings(ctx, data.Settings)
		if err != nil {
			return nil, fmt.Errorf("error constructing settings: %v", err)
		}
		if settings != nil {
			assignment.SetSettings(settings)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed mobile app assignment", assignment); err != nil {
		tflog.Error(ctx, "Failed to log mobile app assignment", map[string]any{
			"error": err.Error(),
		})
	}

	return assignment, nil
}

// constructAssignmentTarget constructs the mobile app deployment assignment target
func constructAssignmentTarget(ctx context.Context, data *AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	if data == nil {
		return nil, fmt.Errorf("assignment target data is required")
	}

	var target graphmodels.DeviceAndAppManagementAssignmentTargetable
	targetType := data.DeviceAndAppManagementAssignmentFilterType.ValueString()

	switch data.TargetType.ValueString() {
	case "allDevices":
		target = graphmodels.NewAllDevicesAssignmentTarget()
	case "allLicensedUsers":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "androidFotaDeployment":
		androidFotaDeploymentAssignmentTarget := graphmodels.NewAndroidFotaDeploymentAssignmentTarget()
		if !data.GroupId.IsNull() {
			id := data.GroupId.ValueString()
			androidFotaDeploymentAssignmentTarget.SetGroupId(&id)
		}
		target = androidFotaDeploymentAssignmentTarget
	case "configurationManagerCollection":
		configManagerTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		if !data.CollectionId.IsNull() {
			id := data.CollectionId.ValueString()
			configManagerTarget.SetCollectionId(&id)
		}
		target = configManagerTarget
	case "exclusionGroupAssignment":
		exclusionGroupTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		if !data.GroupId.IsNull() {
			id := data.GroupId.ValueString()
			exclusionGroupTarget.SetGroupId(&id)
		}
		target = exclusionGroupTarget
	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if !data.GroupId.IsNull() {
			id := data.GroupId.ValueString()
			groupTarget.SetGroupId(&id)
		}
		target = groupTarget
	default:
		target = graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	}

	// Then set the filter properties if they exist
	if !data.DeviceAndAppManagementAssignmentFilterId.IsNull() {
		id := data.DeviceAndAppManagementAssignmentFilterId.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&id)
	}

	if !data.DeviceAndAppManagementAssignmentFilterType.IsNull() {
		filterType, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(targetType)
		if err != nil {
			return nil, fmt.Errorf("error parsing filter type: %v", err)
		}
		target.SetDeviceAndAppManagementAssignmentFilterType(filterType.(*graphmodels.DeviceAndAppManagementAssignmentFilterType))
	}

	tflog.Debug(ctx, "Finished constructing assignment target")
	return target, nil
}

func constructMobileAppAssignmentSettings(ctx context.Context, data *MobileAppAssignmentSettingsResourceModel) (graphmodels.MobileAppAssignmentSettingsable, error) {
	if data == nil {
		return nil, nil
	}

	tflog.Debug(ctx, "Constructing mobile app assignment settings")

	// Handle Android Managed Store settings
	if data.AndroidManagedStore != nil {
		settings, err := constructAndroidManagedStoreAppAssignmentSettings(ctx, data.AndroidManagedStore)
		if err != nil {
			return nil, fmt.Errorf("error constructing Android Managed Store app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle iOS Lob App settings
	if data.IosLob != nil {
		settings, err := constructIosLobAppAssignmentSettings(data.IosLob)
		if err != nil {
			return nil, fmt.Errorf("error constructing iOS Lob app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle iOS Store App settings
	if data.IosStore != nil {
		settings, err := constructIosStoreAppAssignmentSettings(data.IosStore)
		if err != nil {
			return nil, fmt.Errorf("error constructing iOS Store app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle iOS VPP App settings
	if data.IosVpp != nil {
		settings, err := constructIosVppAppAssignmentSettings(data.IosVpp)
		if err != nil {
			return nil, fmt.Errorf("error constructing iOS VPP app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle MacOS Lob App settings
	if data.MacOsLob != nil {
		settings := graphmodels.NewMacOsLobAppAssignmentSettings()

		// Set UninstallOnDeviceRemoval
		if !data.MacOsLob.UninstallOnDeviceRemoval.IsNull() {
			settings.SetUninstallOnDeviceRemoval(data.MacOsLob.UninstallOnDeviceRemoval.ValueBoolPointer())
		}

		return settings, nil
	}

	// Handle MacOS VPP App settings
	if data.MacOsVpp != nil {
		settings, err := constructMacOsVppAppAssignmentSettings(data.MacOsVpp)
		if err != nil {
			return nil, fmt.Errorf("error constructing MacOS VPP app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle Microsoft Store for Business App settings
	if data.MicrosoftStoreForBusiness != nil {
		settings, err := constructMicrosoftStoreForBusinessAppAssignmentSettings(data.MicrosoftStoreForBusiness)
		if err != nil {
			return nil, fmt.Errorf("error constructing Microsoft Store for Business app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle Win32Catalog assignment settings
	if data.Win32Catalog != nil {
		settings, err := constructWin32CatalogAppAssignmentSettings(data.Win32Catalog)
		if err != nil {
			return nil, fmt.Errorf("error constructing Win32Catalog app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle Win32 LOB App settings
	if data.Win32Lob != nil {
		settings, err := constructWin32LobAppAssignmentSettings(data.Win32Lob)
		if err != nil {
			return nil, fmt.Errorf("error constructing Win32 LOB app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle Windows AppX App settings
	if data.WindowsAppX != nil {
		settings, err := constructWindowsAppXAssignmentSettings(data.WindowsAppX)
		if err != nil {
			return nil, fmt.Errorf("error constructing Windows AppX app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle Windows Universal AppX App settings
	if data.WindowsUniversalAppX != nil {
		settings, err := constructWindowsUniversalAppXAssignmentSettings(data.WindowsUniversalAppX)
		if err != nil {
			return nil, fmt.Errorf("error constructing Windows Universal AppX app assignment settings: %v", err)
		}
		return settings, nil
	}

	// Handle WinGet settings
	if data.WinGet != nil {
		settings, err := constructWinGetAppAssignmentSettings(data.WinGet)
		if err != nil {
			return nil, fmt.Errorf("error constructing WinGet app assignment settings: %v", err)
		}
		return settings, nil
	}

	return nil, nil
}

func constructAndroidManagedStoreAppAssignmentSettings(ctx context.Context, data *AndroidManagedStoreAssignmentSettingsResourceModel) (*graphmodels.AndroidManagedStoreAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("android Managed Store data is required")
	}

	settings := graphmodels.NewAndroidManagedStoreAppAssignmentSettings()

	// Set Android Managed Store App Track IDs
	err := convert.FrameworkToGraphStringList(ctx, data.AndroidManagedStoreAppTrackIds, settings.SetAndroidManagedStoreAppTrackIds)
	if err != nil {
		return nil, fmt.Errorf("error setting Android Managed Store App Track IDs: %v", err)
	}

	// Set Auto Update Mode
	err = convert.FrameworkToGraphEnum[*graphmodels.AndroidManagedStoreAutoUpdateMode](
		data.AutoUpdateMode,
		graphmodels.ParseAndroidManagedStoreAutoUpdateMode,
		settings.SetAutoUpdateMode,
	)
	if err != nil {
		return nil, fmt.Errorf("error setting auto update mode: %v", err)
	}

	return settings, nil
}

func constructIosLobAppAssignmentSettings(data *IosLobAppAssignmentSettingsResourceModel) (*graphmodels.IosLobAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("iOS Lob App data is required")
	}

	settings := graphmodels.NewIosLobAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.IsRemovable, settings.SetIsRemovable)
	convert.FrameworkToGraphBool(data.PreventManagedAppBackup, settings.SetPreventManagedAppBackup)
	convert.FrameworkToGraphBool(data.UninstallOnDeviceRemoval, settings.SetUninstallOnDeviceRemoval)
	convert.FrameworkToGraphString(data.VpnConfigurationId, settings.SetVpnConfigurationId)

	return settings, nil
}

func constructIosStoreAppAssignmentSettings(data *IosStoreAppAssignmentSettingsResourceModel) (*graphmodels.IosStoreAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("iOS Store App data is required")
	}

	settings := graphmodels.NewIosStoreAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.IsRemovable, settings.SetIsRemovable)
	convert.FrameworkToGraphBool(data.PreventManagedAppBackup, settings.SetPreventManagedAppBackup)
	convert.FrameworkToGraphBool(data.UninstallOnDeviceRemoval, settings.SetUninstallOnDeviceRemoval)
	convert.FrameworkToGraphString(data.VpnConfigurationId, settings.SetVpnConfigurationId)

	return settings, nil
}

func constructIosVppAppAssignmentSettings(data *IosVppAppAssignmentSettingsResourceModel) (*graphmodels.IosVppAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("iOS VPP App data is required")
	}

	settings := graphmodels.NewIosVppAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.IsRemovable, settings.SetIsRemovable)
	convert.FrameworkToGraphBool(data.PreventAutoAppUpdate, settings.SetPreventAutoAppUpdate)
	convert.FrameworkToGraphBool(data.PreventManagedAppBackup, settings.SetPreventManagedAppBackup)
	convert.FrameworkToGraphBool(data.UninstallOnDeviceRemoval, settings.SetUninstallOnDeviceRemoval)
	convert.FrameworkToGraphBool(data.UseDeviceLicensing, settings.SetUseDeviceLicensing)
	convert.FrameworkToGraphString(data.VpnConfigurationId, settings.SetVpnConfigurationId)

	return settings, nil
}

func constructMacOsVppAppAssignmentSettings(data *MacOsVppAppAssignmentSettingsResourceModel) (*graphmodels.MacOsVppAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("MacOS VPP App data is required")
	}

	settings := graphmodels.NewMacOsVppAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.PreventAutoAppUpdate, settings.SetPreventAutoAppUpdate)
	convert.FrameworkToGraphBool(data.PreventManagedAppBackup, settings.SetPreventManagedAppBackup)
	convert.FrameworkToGraphBool(data.UninstallOnDeviceRemoval, settings.SetUninstallOnDeviceRemoval)
	convert.FrameworkToGraphBool(data.UseDeviceLicensing, settings.SetUseDeviceLicensing)

	return settings, nil
}

func constructMicrosoftStoreForBusinessAppAssignmentSettings(data *MicrosoftStoreForBusinessAppAssignmentSettingsResourceModel) (*graphmodels.MicrosoftStoreForBusinessAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("microsoft Store for Business App data is required")
	}

	settings := graphmodels.NewMicrosoftStoreForBusinessAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.UseDeviceContext, settings.SetUseDeviceContext)

	return settings, nil
}

func constructWin32CatalogAppAssignmentSettings(data *Win32CatalogAppAssignmentSettingsResourceModel) (*graphmodels.Win32CatalogAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("Win32Catalog App data is required")
	}

	settings := graphmodels.NewWin32CatalogAppAssignmentSettings()

	// Set AutoUpdateSettings
	if data.AutoUpdateSettings != nil {
		autoUpdateSettings := graphmodels.NewWin32LobAppAutoUpdateSettings()

		err := convert.FrameworkToGraphEnum(data.AutoUpdateSettings.AutoUpdateSupersededAppsState, graphmodels.ParseWin32LobAutoUpdateSupersededAppsState, autoUpdateSettings.SetAutoUpdateSupersededAppsState)
		if err != nil {
			return nil, fmt.Errorf("error setting AutoUpdateSupersededAppsState: %v", err)
		}

		settings.SetAutoUpdateSettings(autoUpdateSettings)
	}

	err := convert.FrameworkToGraphEnum(data.DeliveryOptimizationPriority, graphmodels.ParseWin32LobAppDeliveryOptimizationPriority, settings.SetDeliveryOptimizationPriority)
	if err != nil {
		return nil, fmt.Errorf("error setting DeliveryOptimizationPriority: %v", err)
	}

	if data.InstallTimeSettings != nil {
		installTimeSettings := graphmodels.NewMobileAppInstallTimeSettings()

		convert.FrameworkToGraphString(data.InstallTimeSettings.DeadlineDateTime, func(value *string) {
			parsedDeadline, err := time.Parse(time.RFC3339, *value)
			if err == nil {
				installTimeSettings.SetDeadlineDateTime(&parsedDeadline)
			}
		})

		convert.FrameworkToGraphString(data.InstallTimeSettings.StartDateTime, func(value *string) {
			parsedStart, err := time.Parse(time.RFC3339, *value)
			if err == nil {
				installTimeSettings.SetStartDateTime(&parsedStart)
			}
		})

		convert.FrameworkToGraphBool(data.InstallTimeSettings.UseLocalTime, installTimeSettings.SetUseLocalTime)

		settings.SetInstallTimeSettings(installTimeSettings)
	}

	err = convert.FrameworkToGraphEnum(data.Notifications, graphmodels.ParseWin32LobAppNotification, settings.SetNotifications)
	if err != nil {
		return nil, fmt.Errorf("error setting Notifications: %v", err)
	}

	if data.RestartSettings != nil {
		restartSettings := graphmodels.NewWin32LobAppRestartSettings()

		convert.FrameworkToGraphInt32(data.RestartSettings.CountdownDisplayBeforeRestart, restartSettings.SetCountdownDisplayBeforeRestartInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.GracePeriodInMinutes, restartSettings.SetGracePeriodInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.RestartNotificationSnoozeDuration, restartSettings.SetRestartNotificationSnoozeDurationInMinutes)

		settings.SetRestartSettings(restartSettings)
	}

	return settings, nil
}

func constructWin32LobAppAssignmentSettings(data *Win32LobAppAssignmentSettingsResourceModel) (*graphmodels.Win32LobAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("Win32 LOB App data is required")
	}

	settings := graphmodels.NewWin32LobAppAssignmentSettings()

	if data.AutoUpdateSettings != nil {
		autoUpdateSettings := graphmodels.NewWin32LobAppAutoUpdateSettings()

		err := convert.FrameworkToGraphEnum(data.AutoUpdateSettings.AutoUpdateSupersededAppsState, graphmodels.ParseWin32LobAutoUpdateSupersededAppsState, autoUpdateSettings.SetAutoUpdateSupersededAppsState)
		if err != nil {
			return nil, fmt.Errorf("error setting AutoUpdateSupersededAppsState: %v", err)
		}

		settings.SetAutoUpdateSettings(autoUpdateSettings)
	}

	err := convert.FrameworkToGraphEnum(data.DeliveryOptimizationPriority, graphmodels.ParseWin32LobAppDeliveryOptimizationPriority, settings.SetDeliveryOptimizationPriority)
	if err != nil {
		return nil, fmt.Errorf("error setting DeliveryOptimizationPriority: %v", err)
	}

	if data.InstallTimeSettings != nil {
		installTimeSettings := graphmodels.NewMobileAppInstallTimeSettings()

		convert.FrameworkToGraphString(data.InstallTimeSettings.DeadlineDateTime, func(value *string) {
			parsedDeadline, err := time.Parse(time.RFC3339, *value)
			if err == nil {
				installTimeSettings.SetDeadlineDateTime(&parsedDeadline)
			}
		})

		convert.FrameworkToGraphString(data.InstallTimeSettings.StartDateTime, func(value *string) {
			parsedStart, err := time.Parse(time.RFC3339, *value)
			if err == nil {
				installTimeSettings.SetStartDateTime(&parsedStart)
			}
		})

		convert.FrameworkToGraphBool(data.InstallTimeSettings.UseLocalTime, installTimeSettings.SetUseLocalTime)

		settings.SetInstallTimeSettings(installTimeSettings)
	}

	err = convert.FrameworkToGraphEnum(data.Notifications, graphmodels.ParseWin32LobAppNotification, settings.SetNotifications)
	if err != nil {
		return nil, fmt.Errorf("error setting Notifications: %v", err)
	}

	if data.RestartSettings != nil {
		restartSettings := graphmodels.NewWin32LobAppRestartSettings()

		convert.FrameworkToGraphInt32(data.RestartSettings.CountdownDisplayBeforeRestart, restartSettings.SetCountdownDisplayBeforeRestartInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.GracePeriodInMinutes, restartSettings.SetGracePeriodInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.RestartNotificationSnoozeDuration, restartSettings.SetRestartNotificationSnoozeDurationInMinutes)

		settings.SetRestartSettings(restartSettings)
	}

	return settings, nil
}

func constructWindowsAppXAssignmentSettings(data *WindowsAppXAssignmentSettingsResourceModel) (*graphmodels.WindowsAppXAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("windows AppX App data is required")
	}

	settings := graphmodels.NewWindowsAppXAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.UseDeviceContext, settings.SetUseDeviceContext)

	return settings, nil
}

func constructWindowsUniversalAppXAssignmentSettings(data *WindowsUniversalAppXAssignmentSettingsResourceModel) (*graphmodels.WindowsUniversalAppXAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("windows Universal AppX App data is required")
	}

	settings := graphmodels.NewWindowsUniversalAppXAppAssignmentSettings()

	convert.FrameworkToGraphBool(data.UseDeviceContext, settings.SetUseDeviceContext)

	return settings, nil
}
func constructWinGetAppAssignmentSettings(data *WinGetAppAssignmentSettingsResourceModel) (*graphmodels.WinGetAppAssignmentSettings, error) {
	if data == nil {
		return nil, fmt.Errorf("winGet settings data is required")
	}

	settings := graphmodels.NewWinGetAppAssignmentSettings()

	if data.InstallTimeSettings != nil {
		installSettings := graphmodels.NewWinGetAppInstallTimeSettings()

		// Set @odata.type for install time settings
		//odataType := "microsoft.graph.winGetAppInstallTimeSettings"
		//installSettings.SetOdataType(&odataType)

		convert.FrameworkToGraphString(data.InstallTimeSettings.DeadlineDateTime, func(value *string) {
			parsedDeadline, err := time.Parse(time.RFC3339, *value)
			if err == nil {
				installSettings.SetDeadlineDateTime(&parsedDeadline)
			}
		})

		convert.FrameworkToGraphBool(data.InstallTimeSettings.UseLocalTime, installSettings.SetUseLocalTime)

		settings.SetInstallTimeSettings(installSettings)
	}

	err := convert.FrameworkToGraphEnum(data.Notifications, graphmodels.ParseWinGetAppNotification, settings.SetNotifications)
	if err != nil {
		return nil, fmt.Errorf("error setting Notifications: %v", err)
	}

	if data.RestartSettings != nil {
		restartSettings := graphmodels.NewWinGetAppRestartSettings()

		// Set @odata.type for restart settings
		//odataType := "microsoft.graph.winGetAppRestartSettings"
		//restartSettings.SetOdataType(&odataType)

		convert.FrameworkToGraphInt32(data.RestartSettings.CountdownDisplayBeforeRestartInMinutes, restartSettings.SetCountdownDisplayBeforeRestartInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.GracePeriodInMinutes, restartSettings.SetGracePeriodInMinutes)
		convert.FrameworkToGraphInt32(data.RestartSettings.RestartNotificationSnoozeDurationInMinutes, restartSettings.SetRestartNotificationSnoozeDurationInMinutes)

		settings.SetRestartSettings(restartSettings)
	}

	return settings, nil
}
