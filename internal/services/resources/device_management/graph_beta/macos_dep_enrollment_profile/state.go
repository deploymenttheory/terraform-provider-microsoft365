package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the depMacOSEnrollmentProfile to Terraform state.
// depId, when non-empty, will be stated to DepOnboardingSettingsID.
// Note: admin_account_password is write-only and is never returned by Graph, so it is
// intentionally not mapped here.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *MacOSDepEnrollmentProfileResourceModel,
	profile graphmodels.DepMacOSEnrollmentProfileable,
	depId string,
) {
	if profile == nil {
		tflog.Debug(ctx, "Remote depMacOSEnrollmentProfile is nil")
		return
	}

	tflog.Debug(
		ctx,
		"Starting to map remote depMacOSEnrollmentProfile to Terraform state",
		map[string]any{
			"resourceId": convert.GraphToFrameworkString(profile.GetId()),
		},
	)

	data.ID = convert.GraphToFrameworkString(profile.GetId())
	if depId != "" {
		data.DepOnboardingSettingsID = types.StringValue(depId)
	}

	// enrollmentProfile (base)
	data.DisplayName = convert.GraphToFrameworkString(profile.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(profile.GetDescription())
	data.RequiresUserAuthentication = convert.GraphToFrameworkBool(
		profile.GetRequiresUserAuthentication(),
	)
	data.ConfigurationEndpointURL = convert.GraphToFrameworkString(
		profile.GetConfigurationEndpointUrl(),
	)
	data.EnableAuthenticationViaCompanyPortal = convert.GraphToFrameworkBool(
		profile.GetEnableAuthenticationViaCompanyPortal(),
	)
	data.RequireCompanyPortalOnSetupAssistantEnrolledDevices = convert.GraphToFrameworkBool(
		profile.GetRequireCompanyPortalOnSetupAssistantEnrolledDevices(),
	)

	// depEnrollmentBaseProfile (inherited)
	data.IsDefault = convert.GraphToFrameworkBool(profile.GetIsDefault())
	data.IsMandatory = convert.GraphToFrameworkBool(profile.GetIsMandatory())
	data.SupervisedModeEnabled = convert.GraphToFrameworkBool(profile.GetSupervisedModeEnabled())
	data.SupportDepartment = convert.GraphToFrameworkString(profile.GetSupportDepartment())
	data.SupportPhoneNumber = convert.GraphToFrameworkString(profile.GetSupportPhoneNumber())
	data.DeviceNameTemplate = convert.GraphToFrameworkString(profile.GetDeviceNameTemplate())
	data.ProfileRemovalDisabled = convert.GraphToFrameworkBool(profile.GetProfileRemovalDisabled())
	data.ConfigurationWebURL = convert.GraphToFrameworkBool(profile.GetConfigurationWebUrl())
	data.AwaitDeviceConfigured = convert.GraphToFrameworkBool(
		profile.GetWaitForDeviceConfiguredConfirmation(),
	)
	data.EnabledSkipKeys = convert.GraphToFrameworkStringSet(ctx, profile.GetEnabledSkipKeys())

	// welcome_screen_disabled has no dedicated Graph property; derive it from the
	// presence of the "Welcome" key in enabledSkipKeys.
	data.WelcomeScreenDisabled = types.BoolValue(
		slices.Contains(profile.GetEnabledSkipKeys(), "Welcome"),
	)

	groupIds := profile.GetEnrollmentTimeAzureAdGroupIds()
	groupIdStrings := make([]string, 0, len(groupIds))
	for _, g := range groupIds {
		groupIdStrings = append(groupIdStrings, g.String())
	}
	data.EnrollmentTimeAzureAdGroupIds = convert.GraphToFrameworkStringSet(ctx, groupIdStrings)

	// Setup Assistant pane skip booleans (inherited)
	data.LocationDisabled = convert.GraphToFrameworkBool(profile.GetLocationDisabled())
	data.RestoreBlocked = convert.GraphToFrameworkBool(profile.GetRestoreBlocked())
	data.AppleIdDisabled = convert.GraphToFrameworkBool(profile.GetAppleIdDisabled())
	data.TermsAndConditionsDisabled = convert.GraphToFrameworkBool(
		profile.GetTermsAndConditionsDisabled(),
	)
	data.TouchIdDisabled = convert.GraphToFrameworkBool(profile.GetTouchIdDisabled())
	data.ApplePayDisabled = convert.GraphToFrameworkBool(profile.GetApplePayDisabled())
	data.SiriDisabled = convert.GraphToFrameworkBool(profile.GetSiriDisabled())
	data.DiagnosticsDisabled = convert.GraphToFrameworkBool(profile.GetDiagnosticsDisabled())
	data.DisplayToneSetupDisabled = convert.GraphToFrameworkBool(
		profile.GetDisplayToneSetupDisabled(),
	)
	data.PrivacyPaneDisabled = convert.GraphToFrameworkBool(profile.GetPrivacyPaneDisabled())
	data.ScreenTimeScreenDisabled = convert.GraphToFrameworkBool(
		profile.GetScreenTimeScreenDisabled(),
	)

	// macOS-specific (depMacOSEnrollmentProfile)
	data.RegistrationDisabled = convert.GraphToFrameworkBool(profile.GetRegistrationDisabled())
	data.FileVaultDisabled = convert.GraphToFrameworkBool(profile.GetFileVaultDisabled())
	data.ICloudDiagnosticsDisabled = convert.GraphToFrameworkBool(
		profile.GetICloudDiagnosticsDisabled(),
	)
	data.PassCodeDisabled = convert.GraphToFrameworkBool(profile.GetPassCodeDisabled())
	data.ZoomDisabled = convert.GraphToFrameworkBool(profile.GetZoomDisabled())
	data.ICloudStorageDisabled = convert.GraphToFrameworkBool(profile.GetICloudStorageDisabled())
	data.ChooseYourLockScreenDisabled = convert.GraphToFrameworkBool(
		profile.GetChooseYourLockScreenDisabled(),
	)
	data.AccessibilityScreenDisabled = convert.GraphToFrameworkBool(
		profile.GetAccessibilityScreenDisabled(),
	)
	data.AutoUnlockWithWatchDisabled = convert.GraphToFrameworkBool(
		profile.GetAutoUnlockWithWatchDisabled(),
	)
	data.AutoAdvanceSetupEnabled = convert.GraphToFrameworkBool(
		profile.GetAutoAdvanceSetupEnabled(),
	)
	data.RequestRequiresNetworkTether = convert.GraphToFrameworkBool(
		profile.GetRequestRequiresNetworkTether(),
	)
	data.UsePlatformSSODuringSetupAssistant = convert.GraphToFrameworkBool(
		profile.GetUsePlatformSSODuringSetupAssistant(),
	)

	// Primary (managed local) account auto-creation
	data.SkipPrimarySetupAccountCreation = convert.GraphToFrameworkBool(
		profile.GetSkipPrimarySetupAccountCreation(),
	)
	data.SetPrimarySetupAccountAsRegularUser = convert.GraphToFrameworkBool(
		profile.GetSetPrimarySetupAccountAsRegularUser(),
	)
	data.DontAutoPopulatePrimaryAccountInfo = convert.GraphToFrameworkBool(
		profile.GetDontAutoPopulatePrimaryAccountInfo(),
	)
	data.PrimaryAccountFullName = convert.GraphToFrameworkString(
		profile.GetPrimaryAccountFullName(),
	)
	data.PrimaryAccountUserName = convert.GraphToFrameworkString(
		profile.GetPrimaryAccountUserName(),
	)
	data.EnableRestrictEditing = convert.GraphToFrameworkBool(profile.GetEnableRestrictEditing())

	// Admin (local) account auto-creation (password is write-only, not mapped)
	data.AdminAccountUserName = convert.GraphToFrameworkString(profile.GetAdminAccountUserName())
	data.AdminAccountFullName = convert.GraphToFrameworkString(profile.GetAdminAccountFullName())
	data.HideAdminAccount = convert.GraphToFrameworkBool(profile.GetHideAdminAccount())

	if rotation := profile.GetDepProfileAdminAccountPasswordRotationSetting(); rotation != nil {
		mapped := &AdminAccountPasswordRotationModel{
			AutoRotationPeriodInDays: convert.GraphToFrameworkInt32(
				rotation.GetAutoRotationPeriodInDays(),
			),
		}
		if delay := rotation.GetDepProfileDelayAutoRotationSetting(); delay != nil {
			mapped.OnRetrievalAutoRotatePasswordEnabled = convert.GraphToFrameworkBool(
				delay.GetOnRetrievalAutoRotatePasswordEnabled(),
			)
			mapped.OnRetrievalDelayAutoRotatePasswordInHours = convert.GraphToFrameworkInt32(
				delay.GetOnRetrievalDelayAutoRotatePasswordInHours(),
			)
		}
		data.AdminAccountPasswordRotation = mapped
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()),
	)
}
