package graphBetaDeviceEnrollmentConfiguration

import (
	"context"
	"fmt"
	"strings"

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

	// Map common properties
	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Priority = state.Int32PtrToTypeInt32(remoteResource.GetPriority())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Version = state.Int32PtrToTypeInt32(remoteResource.GetVersion())

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		data.DeviceEnrollmentConfigurationType = types.StringValue(string(*configType))

		// Map configuration type-specific properties
		switch *configType {
		case graphmodels.LIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTLIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapLimitConfigToTerraform(ctx, data, remoteResource)

		case graphmodels.PLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE,
			graphmodels.DEFAULTPLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapPlatformRestrictionsToTerraform(ctx, data, remoteResource)

		case graphmodels.SINGLEPLATFORMRESTRICTION_DEVICEENROLLMENTCONFIGURATIONTYPE:
			mapSinglePlatformRestrictionToTerraform(ctx, data, remoteResource)

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
		restrictions := &PlatformRestrictionModel{}

		// Map each platform restriction
		if androidRestriction := platformConfig.GetAndroidRestriction(); androidRestriction != nil {
			restrictions.AndroidRestriction = mapPlatformRestriction(ctx, androidRestriction)
		}

		if androidForWorkRestriction := platformConfig.GetAndroidForWorkRestriction(); androidForWorkRestriction != nil {
			restrictions.AndroidForWorkRestriction = mapPlatformRestriction(ctx, androidForWorkRestriction)
		}

		if iosRestriction := platformConfig.GetIosRestriction(); iosRestriction != nil {
			restrictions.IOSRestriction = mapPlatformRestriction(ctx, iosRestriction)
		}

		if macRestriction := platformConfig.GetMacRestriction(); macRestriction != nil {
			restrictions.MacRestriction = mapPlatformRestriction(ctx, macRestriction)
		}

		if macOSRestriction := platformConfig.GetMacOSRestriction(); macOSRestriction != nil {
			restrictions.MacOSRestriction = mapPlatformRestriction(ctx, macOSRestriction)
		}

		if windowsRestriction := platformConfig.GetWindowsRestriction(); windowsRestriction != nil {
			restrictions.WindowsRestriction = mapPlatformRestriction(ctx, windowsRestriction)
		}

		if windowsMobileRestriction := platformConfig.GetWindowsMobileRestriction(); windowsMobileRestriction != nil {
			restrictions.WindowsMobileRestriction = mapPlatformRestriction(ctx, windowsMobileRestriction)
		}

		if windowsHomeSkuRestriction := platformConfig.GetWindowsHomeSkuRestriction(); windowsHomeSkuRestriction != nil {
			restrictions.WindowsHomeSkuRestriction = mapPlatformRestriction(ctx, windowsHomeSkuRestriction)
		}

		if tvosRestriction := platformConfig.GetTvosRestriction(); tvosRestriction != nil {
			restrictions.TVOSRestriction = mapPlatformRestriction(ctx, tvosRestriction)
		}

		if visionOSRestriction := platformConfig.GetVisionOSRestriction(); visionOSRestriction != nil {
			restrictions.VisionOSRestriction = mapPlatformRestriction(ctx, visionOSRestriction)
		}

		data.PlatformRestriction = restrictions
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentPlatformRestrictionsConfigurationable")
	}
}

// mapPlatformRestriction maps a single platform restriction
func mapPlatformRestriction(ctx context.Context, restriction graphmodels.DeviceEnrollmentPlatformRestrictionable) *DeviceEnrollmentPlatformRestriction {
	if restriction == nil {
		return nil
	}

	return &DeviceEnrollmentPlatformRestriction{
		PlatformBlocked:                 state.BoolPtrToTypeBool(restriction.GetPlatformBlocked()),
		PersonalDeviceEnrollmentBlocked: state.BoolPtrToTypeBool(restriction.GetPersonalDeviceEnrollmentBlocked()),
		OSMinimumVersion:                types.StringPointerValue(restriction.GetOsMinimumVersion()),
		OSMaximumVersion:                types.StringPointerValue(restriction.GetOsMaximumVersion()),
		BlockedManufacturers:            state.StringSliceToSet(ctx, restriction.GetBlockedManufacturers()),
		BlockedSkus:                     state.StringSliceToSet(ctx, restriction.GetBlockedSkus()),
	}
}

// mapSinglePlatformRestrictionToTerraform maps the single platform restriction configuration
func mapSinglePlatformRestrictionToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping single platform restriction configuration")

	// For single platform restrictions, the Graph API returns a DeviceEnrollmentPlatformRestrictionsConfiguration
	if platformConfig, ok := remoteResource.(graphmodels.DeviceEnrollmentPlatformRestrictionsConfigurationable); ok && platformConfig != nil {
		// Access platformType and platformRestriction via backing store since they don't have direct getter methods

		// Try to access the platformType
		platformTypeVal, err := platformConfig.GetBackingStore().Get("platformType")
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to get platformType: %s", err))
			return
		}

		if platformTypeVal == nil {
			tflog.Warn(ctx, "Platform type is nil in single platform restriction configuration")
			return
		}

		var platformTypeStr string
		// Handle different possible types of platformType value
		switch pt := platformTypeVal.(type) {
		case *graphmodels.EnrollmentRestrictionPlatformType:
			platformTypeStr = string(*pt)
		case string:
			platformTypeStr = pt
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unexpected platform type format: %T", platformTypeVal))
			return
		}

		// Try to get the platform restriction
		restrictionVal, err := platformConfig.GetBackingStore().Get("platformRestriction")
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to get platformRestriction: %s", err))
			return
		}

		if restrictionVal == nil {
			tflog.Warn(ctx, "Platform restriction is nil in single platform restriction configuration")
			return
		}

		// Convert to the proper type
		platformRestriction, ok := restrictionVal.(graphmodels.DeviceEnrollmentPlatformRestrictionable)
		if !ok {
			tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentPlatformRestrictionable")
			return
		}

		// Map to the model
		data.NewPlatformRestriction = &NewPlatformRestrictionModel{
			PlatformType: types.StringValue(platformTypeStr),
			Restriction:  mapPlatformRestriction(ctx, platformRestriction),
		}
	} else {
		tflog.Warn(ctx, "Failed to cast to DeviceEnrollmentPlatformRestrictionsConfigurationable")
	}
}

// mapWindows10EnrollmentCompletionPageToTerraform maps the Windows 10 enrollment completion page configuration
func mapWindows10EnrollmentCompletionPageToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	tflog.Debug(ctx, "Mapping Windows 10 enrollment completion page configuration")

	if configType := remoteResource.GetDeviceEnrollmentConfigurationType(); configType != nil {
		if strings.Contains(string(*configType), "DEFAULT") {
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
