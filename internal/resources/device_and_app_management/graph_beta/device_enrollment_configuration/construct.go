package graphBetaDeviceEnrollmentConfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// deviceEnrollmentConfigurationDispatch defines the signature for all config constructor functions
type deviceEnrollmentConfigurationDispatch func(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error)

// deviceEnrollmentConfigDispatch maps string identifiers to their enum type and constructor
var deviceEnrollmentConfigDispatch = map[string]struct {
	EnumType    graphmodels.DeviceEnrollmentConfigurationType
	Constructor deviceEnrollmentConfigurationDispatch
}{
	// Indicates that configuration is of type limit which refers to number of devices a user is allowed to enroll.
	"limit": {
		EnumType:    graphmodels.LIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructDeviceEnrollmentLimitConfig,
	},

	// Indicates that configuration is of type default limit which refers to types of devices a user is allowed to enroll by default.
	"defaultLimit": {
		EnumType:    graphmodels.DEFAULTLIMIT_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructDeviceEnrollmentLimitConfig,
	},

	// Indicates that configuration is of type default Windows Hello which refers to authentication method devices would use by default.
	"defaultWindowsHelloForBusiness": {
		EnumType:    graphmodels.DEFAULTWINDOWSHELLOFORBUSINESS_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructWindowsHelloForBusinessConfig,
	},

	// Indicates that configuration is of type Windows Hello which refers to authentication method devices would use.
	"windowsHelloForBusiness": {
		EnumType:    graphmodels.WINDOWSHELLOFORBUSINESS_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructWindowsHelloForBusinessConfig,
	},

	// Indicates that configuration is of type default platform restriction which refers to types of devices a user is allowed to enroll by default.
	"defaultPlatformRestrictions": {
		EnumType:    graphmodels.DEFAULTPLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructNewPlatformRestrictionsConfig,
	},

	// Indicates that configuration is of type platform restriction which refers to types of devices a user is allowed to enroll.
	"platformRestrictions": {
		EnumType:    graphmodels.PLATFORMRESTRICTIONS_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructNewPlatformRestrictionsConfig,
	},

	// Indicates that configuration is of type single platform restriction which refers to types of devices a user is allowed to enroll.
	"singlePlatformRestriction": {
		EnumType:    graphmodels.SINGLEPLATFORMRESTRICTION_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructNewPlatformRestrictionsConfig,
	},

	// Indicates that configuration is of type default Enrollment status page which refers to startup page displayed during OOBE in Autopilot devices by default.
	"defaultWindows10EnrollmentCompletionPageConfiguration": {
		EnumType:    graphmodels.DEFAULTWINDOWS10ENROLLMENTCOMPLETIONPAGECONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructWindows10EnrollmentCompletionPageConfig,
	},

	// Indicates that configuration is of type Enrollment status page which refers to startup page displayed during OOBE in Autopilot devices.
	"windows10EnrollmentCompletionPageConfiguration": {
		EnumType:    graphmodels.WINDOWS10ENROLLMENTCOMPLETIONPAGECONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructWindows10EnrollmentCompletionPageConfig,
	},

	// Indicates that configuration is of type Comanagement Authority which refers to policies applied to Co-Managed devices.
	"deviceComanagementAuthorityConfiguration": {
		EnumType:    graphmodels.DEVICECOMANAGEMENTAUTHORITYCONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructDeviceComanagementAuthorityConfig,
	},

	// Indicates that configuration is of type Enrollment Notification which refers to types of notification a user receives during enrollment.
	"enrollmentNotificationsConfiguration": {
		EnumType:    graphmodels.ENROLLMENTNOTIFICATIONSCONFIGURATION_DEVICEENROLLMENTCONFIGURATIONTYPE,
		Constructor: constructEnrollmentNotificationsConfig,
	},
}

// constructResource
func constructResource(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, forUpdate bool) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (forUpdate=%v)", ResourceName, forUpdate))

	baseConfig := graphmodels.NewDeviceEnrollmentConfiguration()

	constructors.SetStringProperty(data.DisplayName, baseConfig.SetDisplayName)
	constructors.SetStringProperty(data.Description, baseConfig.SetDescription)
	constructors.SetInt32Property(data.Priority, baseConfig.SetPriority)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, baseConfig.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	deviceEnrollmentConfigurationType := data.DeviceEnrollmentConfigurationType.ValueString()

	dispatch, ok := deviceEnrollmentConfigDispatch[deviceEnrollmentConfigurationType]
	if !ok {
		return nil, fmt.Errorf("unsupported device enrollment configuration type: %s", deviceEnrollmentConfigurationType)
	}

	requestBody, err := dispatch.Constructor(ctx, data, dispatch.EnumType)
	if err != nil {
		return nil, fmt.Errorf("failed to construct %s config: %w", deviceEnrollmentConfigurationType, err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructNewPlatformRestrictionsConfig creates a platform restrictions configuration
func constructNewPlatformRestrictionsConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing platform restrictions configuration")

	platformConfig := graphmodels.NewDeviceEnrollmentPlatformRestrictionsConfiguration()

	if data.NewPlatformRestriction != nil && !data.NewPlatformRestriction.PlatformType.IsNull() {
		restriction := graphmodels.NewDeviceEnrollmentPlatformRestriction()

		constructors.SetBoolProperty(data.NewPlatformRestriction.Restriction.PlatformBlocked, restriction.SetPlatformBlocked)
		constructors.SetBoolProperty(data.NewPlatformRestriction.Restriction.PersonalDeviceEnrollmentBlocked, restriction.SetPersonalDeviceEnrollmentBlocked)
		constructors.SetStringProperty(data.NewPlatformRestriction.Restriction.OSMinimumVersion, restriction.SetOsMinimumVersion)
		constructors.SetStringProperty(data.NewPlatformRestriction.Restriction.OSMaximumVersion, restriction.SetOsMaximumVersion)

		if err := constructors.SetStringSet(ctx, data.NewPlatformRestriction.Restriction.BlockedManufacturers, restriction.SetBlockedManufacturers); err != nil {
			return nil, fmt.Errorf("failed to set blocked manufacturers: %s", err)
		}

		if err := constructors.SetStringSet(ctx, data.NewPlatformRestriction.Restriction.BlockedSkus, restriction.SetBlockedSkus); err != nil {
			return nil, fmt.Errorf("failed to set blocked skus: %s", err)
		}

		platformType := data.NewPlatformRestriction.PlatformType.ValueString()

		switch strings.ToLower(platformType) {
		case "android":
			tflog.Debug(ctx, "Applying restrictions to Android platform")
			platformConfig.SetAndroidRestriction(restriction)
		case "androidforwork":
			tflog.Debug(ctx, "Applying restrictions to Android for Work platform")
			platformConfig.SetAndroidForWorkRestriction(restriction)
		case "ios":
			tflog.Debug(ctx, "Applying restrictions to iOS platform")
			platformConfig.SetIosRestriction(restriction)
		case "mac":
			tflog.Debug(ctx, "Applying restrictions to Mac platform")
			platformConfig.SetMacRestriction(restriction)
		case "macos":
			tflog.Debug(ctx, "Applying restrictions to macOS platform")
			platformConfig.SetMacOSRestriction(restriction)
		case "windows":
			tflog.Debug(ctx, "Applying restrictions to Windows platform")
			platformConfig.SetWindowsRestriction(restriction)
		case "windowsmobile":
			tflog.Debug(ctx, "Applying restrictions to Windows Mobile platform")
			platformConfig.SetWindowsMobileRestriction(restriction)
		case "windowshomesku":
			tflog.Debug(ctx, "Applying restrictions to Windows Home SKU platform")
			platformConfig.SetWindowsHomeSkuRestriction(restriction)
		case "tvos":
			tflog.Debug(ctx, "Applying restrictions to tvOS platform")
			platformConfig.SetTvosRestriction(restriction)
		case "visionos":
			tflog.Debug(ctx, "Applying restrictions to VisionOS platform")
			platformConfig.SetVisionOSRestriction(restriction)
		default:
			return nil, fmt.Errorf("unsupported platform type: %s", platformType)
		}
	} else {
		return nil, fmt.Errorf("platform_type is required for platform restriction configuration")
	}

	return platformConfig, nil
}

// constructDeviceEnrollmentLimitConfig creates a device enrollment limit configuration
func constructDeviceEnrollmentLimitConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing device enrollment limit configuration")

	if data.DeviceEnrollmentLimit == nil {
		return nil, fmt.Errorf("device enrollment limit block must be defined for this configuration type")
	}

	limitConfig := graphmodels.NewDeviceEnrollmentLimitConfiguration()

	constructors.SetInt32Property(data.DeviceEnrollmentLimit.Limit,
		limitConfig.SetLimit)

	return limitConfig, nil
}

// constructWindows10EnrollmentCompletionPageConfig creates a Windows 10 enrollment completion page configuration
func constructWindows10EnrollmentCompletionPageConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing Windows 10 enrollment completion page configuration")

	if data.Windows10EnrollmentCompletionPage == nil {
		return nil, fmt.Errorf("windows10_enrollment_completion_page block must be defined for this configuration type")
	}

	completionConfig := graphmodels.NewWindows10EnrollmentCompletionPageConfiguration()

	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.AllowDeviceResetOnInstallFailure,
		completionConfig.SetAllowDeviceResetOnInstallFailure)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.AllowDeviceUseOnInstallFailure,
		completionConfig.SetAllowDeviceUseOnInstallFailure)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.AllowLogCollectionOnInstallFailure,
		completionConfig.SetAllowLogCollectionOnInstallFailure)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.AllowNonBlockingAppInstallation,
		completionConfig.SetAllowNonBlockingAppInstallation)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.BlockDeviceSetupRetryByUser,
		completionConfig.SetBlockDeviceSetupRetryByUser)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.DisableUserStatusTrackingAfterFirstUser,
		completionConfig.SetDisableUserStatusTrackingAfterFirstUser)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.InstallQualityUpdates,
		completionConfig.SetInstallQualityUpdates)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.ShowInstallationProgress,
		completionConfig.SetShowInstallationProgress)
	constructors.SetBoolProperty(data.Windows10EnrollmentCompletionPage.TrackInstallProgressForAutopilotOnly,
		completionConfig.SetTrackInstallProgressForAutopilotOnly)

	constructors.SetStringProperty(data.Windows10EnrollmentCompletionPage.CustomErrorMessage,
		completionConfig.SetCustomErrorMessage)

	constructors.SetInt32Property(data.Windows10EnrollmentCompletionPage.InstallProgressTimeoutInMinutes,
		completionConfig.SetInstallProgressTimeoutInMinutes)

	if err := constructors.SetStringSet(ctx, data.Windows10EnrollmentCompletionPage.SelectedMobileAppIds,
		completionConfig.SetSelectedMobileAppIds); err != nil {
		return nil, fmt.Errorf("failed to set selected mobile app IDs: %s", err)
	}

	return completionConfig, nil
}

// // constructSinglePlatformRestrictionConfig creates a single platform restriction configuration
// func constructSinglePlatformRestrictionConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
// 	tflog.Debug(ctx, "Constructing single platform restriction configuration")

// 	if data.NewPlatformRestriction == nil || data.NewPlatformRestriction.Restriction == nil {
// 		return nil, fmt.Errorf("new_platform_restriction block is required for singlePlatformRestriction configuration")
// 	}

// 	platformConfig := graphmodels.NewDeviceEnrollmentPlatformRestrictionConfiguration()

// 	restriction := graphmodels.NewDeviceEnrollmentPlatformRestriction()
// 	if err := setDeviceEnrollmentPlatformRestriction(ctx, data.NewPlatformRestriction.Restriction, restriction); err != nil {
// 		return nil, fmt.Errorf("failed to construct platform restriction: %s", err)
// 	}
// 	platformConfig.SetPlatformRestriction(restriction)

// 	if err := constructors.SetEnumProperty(
// 		data.NewPlatformRestriction.PlatformType,
// 		graphmodels.ParseEnrollmentRestrictionPlatformType,
// 		platformConfig.SetPlatformType,
// 	); err != nil {
// 		return nil, fmt.Errorf("failed to set platform type: %s", err)
// 	}

// 	return platformConfig, nil
// }

// constructWindowsHelloForBusinessConfig creates a Windows Hello for Business configuration
func constructWindowsHelloForBusinessConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing Windows Hello for Business configuration")

	helloConfig := graphmodels.NewDeviceEnrollmentWindowsHelloForBusinessConfiguration()

	if data.WindowsHelloForBusiness == nil {
		return nil, fmt.Errorf("windows_hello_for_business block must be defined for this configuration type")
	}

	if data.WindowsHelloForBusiness == nil {
		return helloConfig, nil
	}

	err := constructors.SetEnumProperty(data.WindowsHelloForBusiness.State, graphmodels.ParseEnablement, helloConfig.SetState)
	if err != nil {
		return nil, fmt.Errorf("failed to set state: %s", err)
	}

	err = constructors.SetEnumProperty(data.WindowsHelloForBusiness.EnhancedBiometricsState, graphmodels.ParseEnablement, helloConfig.SetEnhancedBiometricsState)
	if err != nil {
		return nil, fmt.Errorf("failed to set enhanced biometrics state: %s", err)
	}

	err = constructors.SetEnumProperty(data.WindowsHelloForBusiness.SecurityKeyForSignIn, graphmodels.ParseEnablement, helloConfig.SetSecurityKeyForSignIn)
	if err != nil {
		return nil, fmt.Errorf("failed to set security key for sign in: %s", err)
	}

	err = constructors.SetEnumProperty(data.WindowsHelloForBusiness.PinLowercaseCharactersUsage, graphmodels.ParseWindowsHelloForBusinessPinUsage, helloConfig.SetPinLowercaseCharactersUsage)
	if err != nil {
		return nil, fmt.Errorf("failed to set PIN lowercase characters usage: %s", err)
	}

	err = constructors.SetEnumProperty(data.WindowsHelloForBusiness.PinUppercaseCharactersUsage, graphmodels.ParseWindowsHelloForBusinessPinUsage, helloConfig.SetPinUppercaseCharactersUsage)
	if err != nil {
		return nil, fmt.Errorf("failed to set PIN uppercase characters usage: %s", err)
	}

	err = constructors.SetEnumProperty(data.WindowsHelloForBusiness.PinSpecialCharactersUsage, graphmodels.ParseWindowsHelloForBusinessPinUsage, helloConfig.SetPinSpecialCharactersUsage)
	if err != nil {
		return nil, fmt.Errorf("failed to set PIN special characters usage: %s", err)
	}

	constructors.SetInt32Property(data.WindowsHelloForBusiness.EnhancedSignInSecurity, helloConfig.SetEnhancedSignInSecurity)
	constructors.SetInt32Property(data.WindowsHelloForBusiness.PinMinimumLength, helloConfig.SetPinMinimumLength)
	constructors.SetInt32Property(data.WindowsHelloForBusiness.PinMaximumLength, helloConfig.SetPinMaximumLength)
	constructors.SetInt32Property(data.WindowsHelloForBusiness.PinExpirationInDays, helloConfig.SetPinExpirationInDays)
	constructors.SetInt32Property(data.WindowsHelloForBusiness.PinPreviousBlockCount, helloConfig.SetPinPreviousBlockCount)

	constructors.SetBoolProperty(data.WindowsHelloForBusiness.RemotePassportEnabled, helloConfig.SetRemotePassportEnabled)
	constructors.SetBoolProperty(data.WindowsHelloForBusiness.SecurityDeviceRequired, helloConfig.SetSecurityDeviceRequired)
	constructors.SetBoolProperty(data.WindowsHelloForBusiness.UnlockWithBiometricsEnabled, helloConfig.SetUnlockWithBiometricsEnabled)

	return helloConfig, nil
}

func constructDeviceComanagementAuthorityConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing device co-management authority configuration")

	config := graphmodels.NewDeviceComanagementAuthorityConfiguration()

	if data.DeviceComanagementAuthority == nil {
		return nil, fmt.Errorf("device_comanagement_authority block must be defined for this configuration type")
	}

	constructors.SetStringProperty(data.DeviceComanagementAuthority.ConfigurationManagerAgentCommandLineArgument, config.SetConfigurationManagerAgentCommandLineArgument)
	constructors.SetBoolProperty(data.DeviceComanagementAuthority.InstallConfigurationManagerAgent, config.SetInstallConfigurationManagerAgent)
	constructors.SetInt32Property(data.DeviceComanagementAuthority.ManagedDeviceAuthority, config.SetManagedDeviceAuthority)

	return config, nil
}

func constructEnrollmentNotificationsConfig(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, configType graphmodels.DeviceEnrollmentConfigurationType) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, "Constructing enrollment notifications configuration")

	config := graphmodels.NewDeviceEnrollmentNotificationConfiguration()

	if data.DeviceComanagementAuthority == nil {
		return nil, fmt.Errorf("device_enrollment_notification block must be defined for this configuration type")
	}

	constructors.SetStringProperty(data.EnrollmentNotifications.DefaultLocale, config.SetDefaultLocale)

	if err := constructors.SetBitmaskEnumProperty(data.EnrollmentNotifications.BrandingOptions, graphmodels.ParseEnrollmentNotificationBrandingOptions, config.SetBrandingOptions); err != nil {
		return nil, fmt.Errorf("failed to set branding options: %s", err)
	}

	if err := constructors.SetEnumProperty(data.EnrollmentNotifications.PlatformType, graphmodels.ParseEnrollmentRestrictionPlatformType, config.SetPlatformType); err != nil {
		return nil, fmt.Errorf("failed to set platform type: %s", err)
	}

	if err := constructors.SetEnumProperty(data.EnrollmentNotifications.TemplateType, graphmodels.ParseEnrollmentNotificationTemplateType, config.SetTemplateType); err != nil {
		return nil, fmt.Errorf("failed to set template type: %s", err)
	}

	if err := constructors.SetUUIDProperty(data.EnrollmentNotifications.NotificationMessageTemplateId, config.SetNotificationMessageTemplateId); err != nil {
		return nil, fmt.Errorf("failed to set notification message template ID: %s", err)
	}

	if err := constructors.SetStringSet(ctx, data.EnrollmentNotifications.NotificationTemplates, config.SetNotificationTemplates); err != nil {
		return nil, fmt.Errorf("failed to set notification templates: %s", err)
	}

	return config, nil
}
