package graphBetaDeviceEnrollmentConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Priority = state.Int32PtrToTypeInt32(remoteResource.GetPriority())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Version = state.Int32PtrToTypeInt32(remoteResource.GetVersion())

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		data.DeviceEnrollmentConfigurationType = types.StringValue(configType.String())

		// Map configuration type-specific properties
		switch *configType {
		case graphmodels.LIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTLIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapLimitConfigToTerraform(ctx, data, remoteResource)

		case graphmodels.PLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTPLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapPlatformRestrictionsToTerraform(ctx, data, remoteResource)

		case graphmodels.SINGLEPLATFORMRESTRICTION_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapPlatformRestrictionsToTerraform(ctx, data, remoteResource)

		case graphmodels.WINDOWS10ENROLLMENTCOMPLETIONPAGECONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTWINDOWS10ENROLLMENTCOMPLETIONPAGECONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapWindows10EnrollmentCompletionPageToTerraform(ctx, data, remoteResource)

		case graphmodels.WINDOWSHELLOFORBUSINESS_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTWINDOWSHELLOFORBUSINESS_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapWindowsHelloForBusinessToTerraform(ctx, data, remoteResource)

		case graphmodels.DEVICECOMANAGEMENTAUTHORITYCONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapDeviceComanagementAuthorityToTerraform(ctx, data, remoteResource)

		case graphmodels.ENROLLMENTNOTIFICATIONSCONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapEnrollmentNotificationsToTerraform(ctx, data, remoteResource)

		default:
			tflog.Warn(ctx, fmt.Sprintf("Unhandled device enrollment configuration type: %s", *configType))
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform", map[string]interface{}{"resourceId": data.ID.ValueString()})
}

// mapLimitConfigToTerraform maps the limit configuration specific fields
func mapLimitConfigToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping limit configuration")

	if limitConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentLimitConfigurationable); ok && limitConfig != nil {
		data.DeviceEnrollmentLimit = &DeviceEnrollmentLimitModel{
			Limit: state.Int32PtrToTypeInt32(limitConfig.GetLimit()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentLimitConfigurationable")
	}
}

// mapPlatformRestrictionsToTerraform maps the platform restrictions configuration specific fields
func mapPlatformRestrictionsToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping platform restrictions configuration")

	if platformConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentPlatformRestrictionsConfigurationable); ok && platformConfig != nil {
		// Create a new platform restriction model
		platformRestriction := &NewPlatformRestrictionModel{
			Restriction: &PlatformRestrictionModel{},
		}

		// Check each platform and map the first non-nil one we find
		if androidRestriction := platformConfig.GetAndroidRestriction(); androidRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("android")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, androidRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Android platform restriction")
			return
		}

		if androidForWorkRestriction := platformConfig.GetAndroidForWorkRestriction(); androidForWorkRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("androidForWork")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, androidForWorkRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Android for Work platform restriction")
			return
		}

		if iosRestriction := platformConfig.GetIosRestriction(); iosRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("ios")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, iosRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped iOS platform restriction")
			return
		}

		if macRestriction := platformConfig.GetMacRestriction(); macRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("mac")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, macRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Mac platform restriction")
			return
		}

		if macOSRestriction := platformConfig.GetMacOSRestriction(); macOSRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("macos")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, macOSRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped macOS platform restriction")
			return
		}

		if windowsRestriction := platformConfig.GetWindowsRestriction(); windowsRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("windows")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, windowsRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Windows platform restriction")
			return
		}

		if windowsMobileRestriction := platformConfig.GetWindowsMobileRestriction(); windowsMobileRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("windowsMobile")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, windowsMobileRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Windows Mobile platform restriction")
			return
		}

		if windowsHomeSkuRestriction := platformConfig.GetWindowsHomeSkuRestriction(); windowsHomeSkuRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("windowsHomeSku")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, windowsHomeSkuRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped Windows Home SKU platform restriction")
			return
		}

		if tvosRestriction := platformConfig.GetTvosRestriction(); tvosRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("tvos")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, tvosRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped tvOS platform restriction")
			return
		}

		if visionOSRestriction := platformConfig.GetVisionOSRestriction(); visionOSRestriction != nil {
			platformRestriction.PlatformType = types.StringValue("visionos")
			mapPlatformRestrictionProperties(ctx, platformRestriction.Restriction, visionOSRestriction)
			data.NewPlatformRestriction = platformRestriction
			tflog.Debug(ctx, "Mapped VisionOS platform restriction")
			return
		}

		// If we get here, no platform restrictions were found
		tflog.Debug(ctx, "No platform restrictions found")
		data.NewPlatformRestriction = nil

	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentPlatformRestrictionsConfigurationable")
	}
}

// mapPlatformRestrictionProperties maps the properties from a DeviceEnrollmentPlatformRestrictionable to a PlatformRestrictionModel
func mapPlatformRestrictionProperties(ctx context.Context, model *PlatformRestrictionModel, restriction graphmodels.DeviceEnrollmentPlatformRestrictionable) {
	model.PlatformBlocked = state.BoolPtrToTypeBool(restriction.GetPlatformBlocked())
	model.PersonalDeviceEnrollmentBlocked = state.BoolPtrToTypeBool(restriction.GetPersonalDeviceEnrollmentBlocked())
	model.OSMinimumVersion = types.StringPointerValue(restriction.GetOsMinimumVersion())
	model.OSMaximumVersion = types.StringPointerValue(restriction.GetOsMaximumVersion())
	model.BlockedManufacturers = state.StringSliceToSet(ctx, restriction.GetBlockedManufacturers())
	model.BlockedSkus = state.StringSliceToSet(ctx, restriction.GetBlockedSkus())
}

// mapWindows10EnrollmentCompletionPageToTerraform maps the Windows 10 enrollment completion page configuration
func mapWindows10EnrollmentCompletionPageToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping Windows 10 enrollment completion page configuration")

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		if *configType == graphmodels.DEFAULTWINDOWS10ENROLLMENTCOMPLETIONPAGECONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE {
			mapDefaultWindows10EnrollmentCompletionPageToTerraform(ctx, data, remoteResource)
			return
		}
	}

	if completionConfig, ok := remoteResource.(graphmodels.Windows10EnrollmentCompletionPageConfigurationable); ok && completionConfig != nil {
		data.Windows10EnrollmentCompletionPage = &Windows10EnrollmentCompletionPageModel{
			AllowDeviceResetOnInstallFailure:        state.BoolPtrToTypeBool(completionConfig.GetAllowDeviceResetOnInstallFailure()),
			AllowDeviceUseOnInstallFailure:          state.BoolPtrToTypeBool(completionConfig.GetAllowDeviceUseOnInstallFailure()),
			AllowLogCollectionOnInstallFailure:      state.BoolPtrToTypeBool(completionConfig.GetAllowLogCollectionOnInstallFailure()),
			AllowNonBlockingAppInstallation:         state.BoolPtrToTypeBool(completionConfig.GetAllowNonBlockingAppInstallation()),
			BlockDeviceSetupRetryByUser:             state.BoolPtrToTypeBool(completionConfig.GetBlockDeviceSetupRetryByUser()),
			CustomErrorMessage:                      types.StringPointerValue(completionConfig.GetCustomErrorMessage()),
			DisableUserStatusTrackingAfterFirstUser: state.BoolPtrToTypeBool(completionConfig.GetDisableUserStatusTrackingAfterFirstUser()),
			InstallProgressTimeoutInMinutes:         state.Int32PtrToTypeInt32(completionConfig.GetInstallProgressTimeoutInMinutes()),
			InstallQualityUpdates:                   state.BoolPtrToTypeBool(completionConfig.GetInstallQualityUpdates()),
			SelectedMobileAppIds:                    state.StringSliceToSet(ctx, completionConfig.GetSelectedMobileAppIds()),
			ShowInstallationProgress:                state.BoolPtrToTypeBool(completionConfig.GetShowInstallationProgress()),
			TrackInstallProgressForAutopilotOnly:    state.BoolPtrToTypeBool(completionConfig.GetTrackInstallProgressForAutopilotOnly()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to Windows10EnrollmentCompletionPageConfigurationable")
	}
}

// mapDefaultWindows10EnrollmentCompletionPageToTerraform maps the default Windows 10 enrollment completion page configuration
func mapDefaultWindows10EnrollmentCompletionPageToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping default Windows 10 enrollment completion page configuration")

	if completionConfig, ok := remoteResource.(graphmodels.Windows10EnrollmentCompletionPageConfigurationable); ok && completionConfig != nil {
		data.Windows10EnrollmentCompletionPage = &Windows10EnrollmentCompletionPageModel{
			AllowDeviceResetOnInstallFailure:        state.BoolPtrToTypeBool(completionConfig.GetAllowDeviceResetOnInstallFailure()),
			AllowLogCollectionOnInstallFailure:      state.BoolPtrToTypeBool(completionConfig.GetAllowLogCollectionOnInstallFailure()),
			CustomErrorMessage:                      types.StringPointerValue(completionConfig.GetCustomErrorMessage()),
			DisableUserStatusTrackingAfterFirstUser: state.BoolPtrToTypeBool(completionConfig.GetDisableUserStatusTrackingAfterFirstUser()),
			InstallProgressTimeoutInMinutes:         state.Int32PtrToTypeInt32(completionConfig.GetInstallProgressTimeoutInMinutes()),
			SelectedMobileAppIds:                    state.StringSliceToSet(ctx, completionConfig.GetSelectedMobileAppIds()),
			ShowInstallationProgress:                state.BoolPtrToTypeBool(completionConfig.GetShowInstallationProgress()),
			TrackInstallProgressForAutopilotOnly:    state.BoolPtrToTypeBool(completionConfig.GetTrackInstallProgressForAutopilotOnly()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DefaultWindows10EnrollmentCompletionPageConfigurationable")
	}
}

// mapWindowsHelloForBusinessToTerraform maps the Windows Hello for Business configuration
func mapWindowsHelloForBusinessToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping Windows Hello for Business configuration")

	if helloConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentWindowsHelloForBusinessConfigurationable); ok && helloConfig != nil {
		data.WindowsHelloForBusiness = &WindowsHelloForBusinessModel{
			State:                       state.EnumPtrToTypeString(helloConfig.GetState()),
			EnhancedBiometricsState:     state.EnumPtrToTypeString(helloConfig.GetEnhancedBiometricsState()),
			SecurityKeyForSignIn:        state.EnumPtrToTypeString(helloConfig.GetSecurityKeyForSignIn()),
			PinLowercaseCharactersUsage: state.EnumPtrToTypeString(helloConfig.GetPinLowercaseCharactersUsage()),
			PinUppercaseCharactersUsage: state.EnumPtrToTypeString(helloConfig.GetPinUppercaseCharactersUsage()),
			PinSpecialCharactersUsage:   state.EnumPtrToTypeString(helloConfig.GetPinSpecialCharactersUsage()),
			EnhancedSignInSecurity:      state.Int32PtrToTypeInt32(helloConfig.GetEnhancedSignInSecurity()),
			PinMinimumLength:            state.Int32PtrToTypeInt32(helloConfig.GetPinMinimumLength()),
			PinMaximumLength:            state.Int32PtrToTypeInt32(helloConfig.GetPinMaximumLength()),
			PinExpirationInDays:         state.Int32PtrToTypeInt32(helloConfig.GetPinExpirationInDays()),
			PinPreviousBlockCount:       state.Int32PtrToTypeInt32(helloConfig.GetPinPreviousBlockCount()),
			RemotePassportEnabled:       state.BoolPtrToTypeBool(helloConfig.GetRemotePassportEnabled()),
			SecurityDeviceRequired:      state.BoolPtrToTypeBool(helloConfig.GetSecurityDeviceRequired()),
			UnlockWithBiometricsEnabled: state.BoolPtrToTypeBool(helloConfig.GetUnlockWithBiometricsEnabled()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentWindowsHelloForBusinessConfigurationable")
	}
}

// mapDeviceComanagementAuthorityToTerraform maps the device co-management authority configuration
func mapDeviceComanagementAuthorityToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping device co-management authority configuration")

	if comanagementConfig, ok := remoteResource.(graphmodels.DeviceComanagementAuthorityConfigurationable); ok && comanagementConfig != nil {
		data.DeviceComanagementAuthority = &DeviceComanagementAuthorityModel{
			ConfigurationManagerAgentCommandLineArgument: types.StringPointerValue(comanagementConfig.GetConfigurationManagerAgentCommandLineArgument()),
			InstallConfigurationManagerAgent:             state.BoolPtrToTypeBool(comanagementConfig.GetInstallConfigurationManagerAgent()),
			ManagedDeviceAuthority:                       state.Int32PtrToTypeInt32(comanagementConfig.GetManagedDeviceAuthority()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceComanagementAuthorityConfigurationable")
	}
}

// mapEnrollmentNotificationsToTerraform maps the enrollment notifications configuration
func mapEnrollmentNotificationsToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping enrollment notifications configuration")

	if notificationConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentNotificationConfigurationable); ok && notificationConfig != nil {
		data.EnrollmentNotifications = &EnrollmentNotificationsModel{
			DefaultLocale:                 types.StringPointerValue(notificationConfig.GetDefaultLocale()),
			BrandingOptions:               state.EnumPtrToTypeString(notificationConfig.GetBrandingOptions()),
			PlatformType:                  state.EnumPtrToTypeString(notificationConfig.GetPlatformType()),
			TemplateType:                  state.EnumPtrToTypeString(notificationConfig.GetTemplateType()),
			NotificationMessageTemplateId: state.UUIDPtrToTypeString(notificationConfig.GetNotificationMessageTemplateId()),
			NotificationTemplates:         state.StringSliceToSet(ctx, notificationConfig.GetNotificationTemplates()),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentNotificationConfigurationable")
	}
}
