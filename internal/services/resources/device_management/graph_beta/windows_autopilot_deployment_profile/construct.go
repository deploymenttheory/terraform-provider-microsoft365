package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a Windows Autopilot Deployment Profile object for API requests
func constructResource(ctx context.Context, data *WindowsAutopilotDeploymentProfileResourceModel, forUpdate bool) (graphmodels.WindowsAutopilotDeploymentProfileable, error) {
	var resource graphmodels.WindowsAutopilotDeploymentProfileable

	// Determine the profile type based on device join type
	switch data.DeviceJoinType.ValueString() {
	case "microsoft_entra_hybrid_joined":
		// For hybrid domain join scenarios, use ActiveDirectoryWindowsAutopilotDeploymentProfile
		resource = graphmodels.NewActiveDirectoryWindowsAutopilotDeploymentProfile()
	case "microsoft_entra_joined":
		// For pure Azure AD/Entra joined scenarios, use AzureADWindowsAutopilotDeploymentProfile
		resource = graphmodels.NewAzureADWindowsAutopilotDeploymentProfile()
	default:
		return nil, fmt.Errorf("invalid device join type: %s", data.DeviceJoinType.ValueString())
	}

	// Basic properties
	constructors.SetStringProperty(data.DisplayName, resource.SetDisplayName)
	constructors.SetStringProperty(data.Description, resource.SetDescription)
	constructors.SetStringProperty(data.Locale, resource.SetLocale)
	constructors.SetStringProperty(data.DeviceNameTemplate, resource.SetDeviceNameTemplate)
	constructors.SetStringProperty(data.ManagementServiceAppId, resource.SetManagementServiceAppId)

	// Boolean properties
	constructors.SetBoolProperty(data.HardwareHashExtractionEnabled, resource.SetHardwareHashExtractionEnabled)
	constructors.SetBoolProperty(data.PreprovisioningAllowed, resource.SetPreprovisioningAllowed)

	// Set hybrid Azure AD join setting only for ActiveDirectory profiles
	if adProfile, ok := resource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		constructors.SetBoolProperty(data.HybridAzureADJoinSkipConnectivityCheck, adProfile.SetHybridAzureADJoinSkipConnectivityCheck)
	}

	// Device type enum
	if err := constructors.SetEnumProperty(data.DeviceType, graphmodels.ParseWindowsAutopilotDeviceType, resource.SetDeviceType); err != nil {
		return nil, fmt.Errorf("error setting device type: %v", err)
	}

	// Role scope tag IDs
	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, resource.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("error setting role scope tag IDs: %v", err)
	}

	// Set out-of-box experience setting (current)
	if data.OutOfBoxExperienceSetting != nil {
		oobe := graphmodels.NewOutOfBoxExperienceSetting()

		constructors.SetBoolProperty(data.OutOfBoxExperienceSetting.PrivacySettingsHidden, oobe.SetPrivacySettingsHidden)
		constructors.SetBoolProperty(data.OutOfBoxExperienceSetting.EulaHidden, oobe.SetEulaHidden)
		constructors.SetBoolProperty(data.OutOfBoxExperienceSetting.KeyboardSelectionPageSkipped, oobe.SetKeyboardSelectionPageSkipped)
		constructors.SetBoolProperty(data.OutOfBoxExperienceSetting.EscapeLinkHidden, oobe.SetEscapeLinkHidden)

		if err := constructors.SetEnumProperty(data.OutOfBoxExperienceSetting.UserType, graphmodels.ParseWindowsUserType, oobe.SetUserType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting user type: %v", err)
		}

		if err := constructors.SetEnumProperty(data.OutOfBoxExperienceSetting.DeviceUsageType, graphmodels.ParseWindowsDeviceUsageType, oobe.SetDeviceUsageType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting device usage type: %v", err)
		}

		resource.SetOutOfBoxExperienceSetting(oobe)
	}

	// Set enrollment status screen settings
	if data.EnrollmentStatusScreenSettings != nil {
		ess := graphmodels.NewWindowsEnrollmentStatusScreenSettings()

		constructors.SetBoolProperty(data.EnrollmentStatusScreenSettings.HideInstallationProgress, ess.SetHideInstallationProgress)
		constructors.SetBoolProperty(data.EnrollmentStatusScreenSettings.AllowDeviceUseBeforeProfileAndAppInstallComplete, ess.SetAllowDeviceUseBeforeProfileAndAppInstallComplete)
		constructors.SetBoolProperty(data.EnrollmentStatusScreenSettings.BlockDeviceSetupRetryByUser, ess.SetBlockDeviceSetupRetryByUser)
		constructors.SetBoolProperty(data.EnrollmentStatusScreenSettings.AllowLogCollectionOnInstallFailure, ess.SetAllowLogCollectionOnInstallFailure)
		constructors.SetStringProperty(data.EnrollmentStatusScreenSettings.CustomErrorMessage, ess.SetCustomErrorMessage)
		constructors.SetInt32Property(data.EnrollmentStatusScreenSettings.InstallProgressTimeoutInMinutes, ess.SetInstallProgressTimeoutInMinutes)
		constructors.SetBoolProperty(data.EnrollmentStatusScreenSettings.AllowDeviceUseOnInstallFailure, ess.SetAllowDeviceUseOnInstallFailure)

		resource.SetEnrollmentStatusScreenSettings(ess)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed Windows Autopilot Deployment Profile Resource", resource); err != nil {
		tflog.Error(ctx, "Failed to log Windows Autopilot Deployment Profile", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return resource, nil
}
