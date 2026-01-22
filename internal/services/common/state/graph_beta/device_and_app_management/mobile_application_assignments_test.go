package sharedStater

import (
	"context"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStateMobileAppAssignment tests the StateMobileAppAssignment function
func TestStateMobileAppAssignment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		assignments []sharedmodels.MobileAppAssignmentResourceModel
		response    graphmodels.MobileAppAssignmentCollectionResponseable
		validate    func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel)
	}{
		{
			name:        "Nil response",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			response:    nil,
			validate: func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:        "Empty response",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			response: func() graphmodels.MobileAppAssignmentCollectionResponseable {
				resp := graphmodels.NewMobileAppAssignmentCollectionResponse()
				resp.SetValue([]graphmodels.MobileAppAssignmentable{})
				return resp
			}(),
			validate: func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:        "Single group assignment",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			response: func() graphmodels.MobileAppAssignmentCollectionResponseable {
				resp := graphmodels.NewMobileAppAssignmentCollectionResponse()

				assignment := graphmodels.NewMobileAppAssignment()
				assignmentId := "assignment-123"
				intent := graphmodels.REQUIRED_INSTALLINTENT
				source := graphmodels.DIRECT_DEVICEANDAPPMANAGEMENTASSIGNMENTSOURCE

				assignment.SetId(&assignmentId)
				assignment.SetIntent(&intent)
				assignment.SetSource(&source)

				target := graphmodels.NewGroupAssignmentTarget()
				groupId := "group-456"
				target.SetGroupId(&groupId)
				assignment.SetTarget(target)

				resp.SetValue([]graphmodels.MobileAppAssignmentable{assignment})
				return resp
			}(),
			validate: func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel) {
				require.Len(t, result, 1)
				assert.Equal(t, "assignment-123", result[0].Id.ValueString())
				assert.Equal(t, "required", result[0].Intent.ValueString())
				assert.Equal(t, "groupAssignment", result[0].Target.TargetType.ValueString())
				assert.Equal(t, "group-456", result[0].Target.GroupId.ValueString())
			},
		},
		{
			name:        "Multiple assignments with different target types",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			response: func() graphmodels.MobileAppAssignmentCollectionResponseable {
				resp := graphmodels.NewMobileAppAssignmentCollectionResponse()

				// Group assignment
				assignment1 := graphmodels.NewMobileAppAssignment()
				id1 := "assignment-1"
				intent1 := graphmodels.REQUIRED_INSTALLINTENT
				assignment1.SetId(&id1)
				assignment1.SetIntent(&intent1)

				target1 := graphmodels.NewGroupAssignmentTarget()
				groupId1 := "group-123"
				target1.SetGroupId(&groupId1)
				assignment1.SetTarget(target1)

				// All devices assignment
				assignment2 := graphmodels.NewMobileAppAssignment()
				id2 := "assignment-2"
				intent2 := graphmodels.AVAILABLE_INSTALLINTENT
				assignment2.SetId(&id2)
				assignment2.SetIntent(&intent2)

				target2 := graphmodels.NewAllDevicesAssignmentTarget()
				assignment2.SetTarget(target2)

				// All licensed users assignment
				assignment3 := graphmodels.NewMobileAppAssignment()
				id3 := "assignment-3"
				intent3 := graphmodels.UNINSTALL_INSTALLINTENT
				assignment3.SetId(&id3)
				assignment3.SetIntent(&intent3)

				target3 := graphmodels.NewAllLicensedUsersAssignmentTarget()
				assignment3.SetTarget(target3)

				resp.SetValue([]graphmodels.MobileAppAssignmentable{assignment1, assignment2, assignment3})
				return resp
			}(),
			validate: func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel) {
				require.Len(t, result, 3)
				// Results should be sorted by intent
				assert.Equal(t, "required", result[0].Intent.ValueString())
				assert.Equal(t, "available", result[1].Intent.ValueString())
				assert.Equal(t, "uninstall", result[2].Intent.ValueString())
			},
		},
		{
			name:        "Response with null GetValue",
			assignments: []sharedmodels.MobileAppAssignmentResourceModel{},
			response: func() graphmodels.MobileAppAssignmentCollectionResponseable {
				resp := graphmodels.NewMobileAppAssignmentCollectionResponse()
				resp.SetValue(nil)
				return resp
			}(),
			validate: func(t *testing.T, result []sharedmodels.MobileAppAssignmentResourceModel) {
				assert.Len(t, result, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StateMobileAppAssignment(ctx, tt.assignments, tt.response)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestMapRemoteTargetToTerraform tests the mapRemoteTargetToTerraform function
func TestMapRemoteTargetToTerraform(t *testing.T) {
	tests := []struct {
		name     string
		target   graphmodels.DeviceAndAppManagementAssignmentTargetable
		validate func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel)
	}{
		{
			name: "GroupAssignmentTarget",
			target: func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
				target := graphmodels.NewGroupAssignmentTarget()
				groupId := "group-123"
				target.SetGroupId(&groupId)
				return target
			}(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "groupAssignment", result.TargetType.ValueString())
				assert.Equal(t, "group-123", result.GroupId.ValueString())
			},
		},
		{
			name: "ExclusionGroupAssignmentTarget",
			target: func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
				target := graphmodels.NewExclusionGroupAssignmentTarget()
				groupId := "group-456"
				target.SetGroupId(&groupId)
				return target
			}(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "exclusionGroupAssignment", result.TargetType.ValueString())
				assert.Equal(t, "group-456", result.GroupId.ValueString())
			},
		},
		{
			name: "ConfigurationManagerCollectionAssignmentTarget",
			target: func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
				target := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
				collectionId := "collection-789"
				target.SetCollectionId(&collectionId)
				return target
			}(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "configurationManagerCollection", result.TargetType.ValueString())
				assert.Equal(t, "collection-789", result.CollectionId.ValueString())
			},
		},
		{
			name:   "AllDevicesAssignmentTarget",
			target: graphmodels.NewAllDevicesAssignmentTarget(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "allDevices", result.TargetType.ValueString())
			},
		},
		{
			name:   "AllLicensedUsersAssignmentTarget",
			target: graphmodels.NewAllLicensedUsersAssignmentTarget(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "allLicensedUsers", result.TargetType.ValueString())
			},
		},
		{
			name: "Target with filter",
			target: func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
				target := graphmodels.NewGroupAssignmentTarget()
				groupId := "group-123"
				filterId := "filter-456"
				filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE

				target.SetGroupId(&groupId)
				target.SetDeviceAndAppManagementAssignmentFilterId(&filterId)
				target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
				return target
			}(),
			validate: func(t *testing.T, result sharedmodels.AssignmentTargetResourceModel) {
				assert.Equal(t, "groupAssignment", result.TargetType.ValueString())
				assert.Equal(t, "group-123", result.GroupId.ValueString())
				assert.Equal(t, "filter-456", result.DeviceAndAppManagementAssignmentFilterId.ValueString())
				assert.Equal(t, "include", result.DeviceAndAppManagementAssignmentFilterType.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapRemoteTargetToTerraform(tt.target)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestMapRemoteSettingsToTerraform tests the mapRemoteSettingsToTerraform function
func TestMapRemoteSettingsToTerraform(t *testing.T) {
	tests := []struct {
		name     string
		settings graphmodels.MobileAppAssignmentSettingsable
		validate func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel)
	}{
		{
			name:     "Nil settings",
			settings: nil,
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				assert.Nil(t, result)
			},
		},
		{
			name: "AndroidManagedStoreAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewAndroidManagedStoreAppAssignmentSettings()
				trackIds := []string{"track-1", "track-2"}
				autoUpdateMode := graphmodels.DEFAULT_ANDROIDMANAGEDSTOREAUTOUPDATEMODE
				settings.SetAndroidManagedStoreAppTrackIds(trackIds)
				settings.SetAutoUpdateMode(&autoUpdateMode)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.AndroidManagedStore)
				assert.Equal(t, "default", result.AndroidManagedStore.AutoUpdateMode.ValueString())
			},
		},
		{
			name: "IosLobAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewIosLobAppAssignmentSettings()
				isRemovable := true
				preventBackup := false
				uninstall := true
				vpnId := "vpn-123"

				settings.SetIsRemovable(&isRemovable)
				settings.SetPreventManagedAppBackup(&preventBackup)
				settings.SetUninstallOnDeviceRemoval(&uninstall)
				settings.SetVpnConfigurationId(&vpnId)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.IosLob)
				assert.True(t, result.IosLob.IsRemovable.ValueBool())
				assert.False(t, result.IosLob.PreventManagedAppBackup.ValueBool())
				assert.True(t, result.IosLob.UninstallOnDeviceRemoval.ValueBool())
				assert.Equal(t, "vpn-123", result.IosLob.VpnConfigurationId.ValueString())
			},
		},
		{
			name: "Win32LobAppAssignmentSettings without nested settings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWin32LobAppAssignmentSettings()
				notifications := graphmodels.SHOWALL_WIN32LOBAPPNOTIFICATION
				priority := graphmodels.NOTCONFIGURED_WIN32LOBAPPDELIVERYOPTIMIZATIONPRIORITY

				settings.SetNotifications(&notifications)
				settings.SetDeliveryOptimizationPriority(&priority)
				// Explicitly set nil for install and restart settings
				settings.SetInstallTimeSettings(nil)
				settings.SetRestartSettings(nil)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.Win32Lob)
				assert.Equal(t, "showAll", result.Win32Lob.Notifications.ValueString())
				assert.Equal(t, "notConfigured", result.Win32Lob.DeliveryOptimizationPriority.ValueString())
				// Nested settings should be nil
				assert.Nil(t, result.Win32Lob.InstallTimeSettings)
				assert.Nil(t, result.Win32Lob.RestartSettings)
			},
		},
		{
			name: "Win32LobAppAssignmentSettings with install and restart settings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWin32LobAppAssignmentSettings()

				// Install time settings
				installSettings := graphmodels.NewMobileAppInstallTimeSettings()
				useLocalTime := true
				installSettings.SetUseLocalTime(&useLocalTime)
				settings.SetInstallTimeSettings(installSettings)

				// Restart settings
				restartSettings := graphmodels.NewWin32LobAppRestartSettings()
				countdown := int32(15)
				gracePeriod := int32(60)
				snooze := int32(30)
				restartSettings.SetCountdownDisplayBeforeRestartInMinutes(&countdown)
				restartSettings.SetGracePeriodInMinutes(&gracePeriod)
				restartSettings.SetRestartNotificationSnoozeDurationInMinutes(&snooze)
				settings.SetRestartSettings(restartSettings)

				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.Win32Lob)
				require.NotNil(t, result.Win32Lob.InstallTimeSettings)
				require.NotNil(t, result.Win32Lob.RestartSettings)

				assert.True(t, result.Win32Lob.InstallTimeSettings.UseLocalTime.ValueBool())
				assert.Equal(t, int32(15), result.Win32Lob.RestartSettings.CountdownDisplayBeforeRestart.ValueInt32())
				assert.Equal(t, int32(60), result.Win32Lob.RestartSettings.GracePeriodInMinutes.ValueInt32())
				assert.Equal(t, int32(30), result.Win32Lob.RestartSettings.RestartNotificationSnoozeDuration.ValueInt32())
			},
		},
		{
			name: "MacOsLobAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewMacOsLobAppAssignmentSettings()
				uninstall := true
				settings.SetUninstallOnDeviceRemoval(&uninstall)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.MacOsLob)
				assert.True(t, result.MacOsLob.UninstallOnDeviceRemoval.ValueBool())
			},
		},
		{
			name: "MacOsVppAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewMacOsVppAppAssignmentSettings()
				preventAutoUpdate := true
				preventBackup := false
				uninstall := true
				useDeviceLicensing := true

				settings.SetPreventAutoAppUpdate(&preventAutoUpdate)
				settings.SetPreventManagedAppBackup(&preventBackup)
				settings.SetUninstallOnDeviceRemoval(&uninstall)
				settings.SetUseDeviceLicensing(&useDeviceLicensing)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.MacOsVpp)
				assert.True(t, result.MacOsVpp.PreventAutoAppUpdate.ValueBool())
				assert.False(t, result.MacOsVpp.PreventManagedAppBackup.ValueBool())
				assert.True(t, result.MacOsVpp.UninstallOnDeviceRemoval.ValueBool())
				assert.True(t, result.MacOsVpp.UseDeviceLicensing.ValueBool())
			},
		},
		{
			name: "MicrosoftStoreForBusinessAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewMicrosoftStoreForBusinessAppAssignmentSettings()
				useDeviceContext := true
				settings.SetUseDeviceContext(&useDeviceContext)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.MicrosoftStoreForBusiness)
				assert.True(t, result.MicrosoftStoreForBusiness.UseDeviceContext.ValueBool())
			},
		},
		{
			name: "WindowsAppXAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWindowsAppXAppAssignmentSettings()
				useDeviceContext := true
				settings.SetUseDeviceContext(&useDeviceContext)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.WindowsAppX)
				assert.True(t, result.WindowsAppX.UseDeviceContext.ValueBool())
			},
		},
		{
			name: "WindowsUniversalAppXAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWindowsUniversalAppXAppAssignmentSettings()
				useDeviceContext := false
				settings.SetUseDeviceContext(&useDeviceContext)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.WindowsUniversalAppX)
				assert.False(t, result.WindowsUniversalAppX.UseDeviceContext.ValueBool())
			},
		},
		{
			name: "WinGetAppAssignmentSettings without nested settings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWinGetAppAssignmentSettings()
				notifications := graphmodels.SHOWALL_WINGETAPPNOTIFICATION
				settings.SetNotifications(&notifications)
				// Explicitly set nil for install and restart settings
				settings.SetInstallTimeSettings(nil)
				settings.SetRestartSettings(nil)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.WinGet)
				assert.Equal(t, "showAll", result.WinGet.Notifications.ValueString())
				// Nested settings should be nil
				assert.Nil(t, result.WinGet.InstallTimeSettings)
				assert.Nil(t, result.WinGet.RestartSettings)
			},
		},
		{
			name: "IosStoreAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewIosStoreAppAssignmentSettings()
				isRemovable := true
				settings.SetIsRemovable(&isRemovable)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.IosStore)
				assert.True(t, result.IosStore.IsRemovable.ValueBool())
			},
		},
		{
			name: "IosVppAppAssignmentSettings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewIosVppAppAssignmentSettings()
				useDeviceLicensing := true
				settings.SetUseDeviceLicensing(&useDeviceLicensing)
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.IosVpp)
				assert.True(t, result.IosVpp.UseDeviceLicensing.ValueBool())
			},
		},
		{
			name: "WinGetAppAssignmentSettings with install and restart settings",
			settings: func() graphmodels.MobileAppAssignmentSettingsable {
				settings := graphmodels.NewWinGetAppAssignmentSettings()
				notifications := graphmodels.SHOWALL_WINGETAPPNOTIFICATION
				settings.SetNotifications(&notifications)
				
				// Install time settings
				installSettings := graphmodels.NewWinGetAppInstallTimeSettings()
				useLocalTime := true
				installSettings.SetUseLocalTime(&useLocalTime)
				settings.SetInstallTimeSettings(installSettings)
				
				// Restart settings
				restartSettings := graphmodels.NewWinGetAppRestartSettings()
				countdown := int32(10)
				gracePeriod := int32(45)
				snooze := int32(20)
				restartSettings.SetCountdownDisplayBeforeRestartInMinutes(&countdown)
				restartSettings.SetGracePeriodInMinutes(&gracePeriod)
				restartSettings.SetRestartNotificationSnoozeDurationInMinutes(&snooze)
				settings.SetRestartSettings(restartSettings)
				
				return settings
			}(),
			validate: func(t *testing.T, result *sharedmodels.MobileAppAssignmentSettingsResourceModel) {
				require.NotNil(t, result)
				require.NotNil(t, result.WinGet)
				require.NotNil(t, result.WinGet.InstallTimeSettings)
				require.NotNil(t, result.WinGet.RestartSettings)
				
				assert.True(t, result.WinGet.InstallTimeSettings.UseLocalTime.ValueBool())
				assert.Equal(t, int32(10), result.WinGet.RestartSettings.CountdownDisplayBeforeRestartInMinutes.ValueInt32())
				assert.Equal(t, int32(45), result.WinGet.RestartSettings.GracePeriodInMinutes.ValueInt32())
				assert.Equal(t, int32(20), result.WinGet.RestartSettings.RestartNotificationSnoozeDurationInMinutes.ValueInt32())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapRemoteSettingsToTerraform(tt.settings)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
