package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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

	convert.FrameworkToGraphString(data.DisplayName, resource.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, resource.SetDescription)
	convert.FrameworkToGraphString(data.Locale, resource.SetLocale)
	convert.FrameworkToGraphString(data.DeviceNameTemplate, resource.SetDeviceNameTemplate)
	convert.FrameworkToGraphString(data.ManagementServiceAppId, resource.SetManagementServiceAppId)
	convert.FrameworkToGraphBool(data.HardwareHashExtractionEnabled, resource.SetHardwareHashExtractionEnabled)
	convert.FrameworkToGraphBool(data.PreprovisioningAllowed, resource.SetPreprovisioningAllowed)

	// Set hybrid Azure AD join setting only for ActiveDirectory profiles
	if adProfile, ok := resource.(graphmodels.ActiveDirectoryWindowsAutopilotDeploymentProfileable); ok {
		convert.FrameworkToGraphBool(data.HybridAzureADJoinSkipConnectivityCheck, adProfile.SetHybridAzureADJoinSkipConnectivityCheck)
	}

	if err := convert.FrameworkToGraphEnum(data.DeviceType, graphmodels.ParseWindowsAutopilotDeviceType, resource.SetDeviceType); err != nil {
		return nil, fmt.Errorf("error setting device type: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, resource.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("error setting role scope tag IDs: %v", err)
	}

	if data.OutOfBoxExperienceSetting != nil {
		oobe := graphmodels.NewOutOfBoxExperienceSetting()

		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.PrivacySettingsHidden, oobe.SetPrivacySettingsHidden)
		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.EulaHidden, oobe.SetEulaHidden)
		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.KeyboardSelectionPageSkipped, oobe.SetKeyboardSelectionPageSkipped)
		convert.FrameworkToGraphBool(data.OutOfBoxExperienceSetting.EscapeLinkHidden, oobe.SetEscapeLinkHidden)

		if err := convert.FrameworkToGraphEnum(data.OutOfBoxExperienceSetting.UserType, graphmodels.ParseWindowsUserType, oobe.SetUserType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting user type: %v", err)
		}

		if err := convert.FrameworkToGraphEnum(data.OutOfBoxExperienceSetting.DeviceUsageType, graphmodels.ParseWindowsDeviceUsageType, oobe.SetDeviceUsageType); err != nil {
			return nil, fmt.Errorf("error setting OOBE setting device usage type: %v", err)
		}

		resource.SetOutOfBoxExperienceSetting(oobe)
	}

	if data.EnrollmentStatusScreenSettings != nil {
		ess := graphmodels.NewWindowsEnrollmentStatusScreenSettings()

		convert.FrameworkToGraphBool(data.EnrollmentStatusScreenSettings.HideInstallationProgress, ess.SetHideInstallationProgress)
		convert.FrameworkToGraphBool(data.EnrollmentStatusScreenSettings.AllowDeviceUseBeforeProfileAndAppInstallComplete, ess.SetAllowDeviceUseBeforeProfileAndAppInstallComplete)
		convert.FrameworkToGraphBool(data.EnrollmentStatusScreenSettings.BlockDeviceSetupRetryByUser, ess.SetBlockDeviceSetupRetryByUser)
		convert.FrameworkToGraphBool(data.EnrollmentStatusScreenSettings.AllowLogCollectionOnInstallFailure, ess.SetAllowLogCollectionOnInstallFailure)
		convert.FrameworkToGraphString(data.EnrollmentStatusScreenSettings.CustomErrorMessage, ess.SetCustomErrorMessage)
		convert.FrameworkToGraphInt32(data.EnrollmentStatusScreenSettings.InstallProgressTimeoutInMinutes, ess.SetInstallProgressTimeoutInMinutes)
		convert.FrameworkToGraphBool(data.EnrollmentStatusScreenSettings.AllowDeviceUseOnInstallFailure, ess.SetAllowDeviceUseOnInstallFailure)

		resource.SetEnrollmentStatusScreenSettings(ess)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed Windows Autopilot Deployment Profile Resource", resource); err != nil {
		tflog.Error(ctx, "Failed to log Windows Autopilot Deployment Profile", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return resource, nil
}
