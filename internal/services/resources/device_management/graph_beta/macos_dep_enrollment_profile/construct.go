package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errUnexpectedElementType = errors.New("unexpected set element type")

// constructResource builds a depMacOSEnrollmentProfile request body from the Terraform model.
// The SDK type DepMacOSEnrollmentProfile automatically sets
// @odata.type = #microsoft.graph.depMacOSEnrollmentProfile.
func constructResource(
	ctx context.Context,
	data *MacOSDepEnrollmentProfileResourceModel,
) (graphmodels.DepMacOSEnrollmentProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDepMacOSEnrollmentProfile()

	// enrollmentProfile (base)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(
		data.RequiresUserAuthentication,
		requestBody.SetRequiresUserAuthentication,
	)
	convert.FrameworkToGraphBool(
		data.EnableAuthenticationViaCompanyPortal,
		requestBody.SetEnableAuthenticationViaCompanyPortal,
	)
	convert.FrameworkToGraphBool(
		data.RequireCompanyPortalOnSetupAssistantEnrolledDevices,
		requestBody.SetRequireCompanyPortalOnSetupAssistantEnrolledDevices,
	)

	// depEnrollmentBaseProfile (inherited)
	convert.FrameworkToGraphBool(data.IsMandatory, requestBody.SetIsMandatory)
	convert.FrameworkToGraphBool(data.SupervisedModeEnabled, requestBody.SetSupervisedModeEnabled)
	convert.FrameworkToGraphString(data.SupportDepartment, requestBody.SetSupportDepartment)
	convert.FrameworkToGraphString(data.SupportPhoneNumber, requestBody.SetSupportPhoneNumber)
	convert.FrameworkToGraphString(data.DeviceNameTemplate, requestBody.SetDeviceNameTemplate)
	convert.FrameworkToGraphBool(data.ProfileRemovalDisabled, requestBody.SetProfileRemovalDisabled)
	convert.FrameworkToGraphBool(data.ConfigurationWebURL, requestBody.SetConfigurationWebUrl)
	convert.FrameworkToGraphBool(
		data.AwaitDeviceConfigured,
		requestBody.SetWaitForDeviceConfiguredConfirmation,
	)

	// enabledSkipKeys is derived from the individual *_disabled booleans (single source
	// of truth). Privacy and Registration are intentionally excluded because Graph rejects
	// those skip-key strings even though their boolean properties work.
	requestBody.SetEnabledSkipKeys(buildEnabledSkipKeys(data))

	if err := constructAzureAdGroupIds(
		ctx,
		data.EnrollmentTimeAzureAdGroupIds,
		requestBody.SetEnrollmentTimeAzureAdGroupIds,
	); err != nil {
		return nil, fmt.Errorf("failed to set enrollment_time_azure_ad_group_ids: %w", err)
	}

	// Setup Assistant pane skip booleans (inherited)
	convert.FrameworkToGraphBool(data.LocationDisabled, requestBody.SetLocationDisabled)
	convert.FrameworkToGraphBool(data.RestoreBlocked, requestBody.SetRestoreBlocked)
	convert.FrameworkToGraphBool(data.AppleIdDisabled, requestBody.SetAppleIdDisabled)
	convert.FrameworkToGraphBool(
		data.TermsAndConditionsDisabled,
		requestBody.SetTermsAndConditionsDisabled,
	)
	convert.FrameworkToGraphBool(data.TouchIdDisabled, requestBody.SetTouchIdDisabled)
	convert.FrameworkToGraphBool(data.ApplePayDisabled, requestBody.SetApplePayDisabled)
	convert.FrameworkToGraphBool(data.SiriDisabled, requestBody.SetSiriDisabled)
	convert.FrameworkToGraphBool(data.DiagnosticsDisabled, requestBody.SetDiagnosticsDisabled)
	convert.FrameworkToGraphBool(
		data.DisplayToneSetupDisabled,
		requestBody.SetDisplayToneSetupDisabled,
	)
	convert.FrameworkToGraphBool(data.PrivacyPaneDisabled, requestBody.SetPrivacyPaneDisabled)
	convert.FrameworkToGraphBool(
		data.ScreenTimeScreenDisabled,
		requestBody.SetScreenTimeScreenDisabled,
	)

	// macOS-specific (depMacOSEnrollmentProfile)
	// Note: welcome_screen_disabled has no dedicated Graph property; it is only expressed
	// through the "Welcome" entry in enabledSkipKeys (see buildEnabledSkipKeys).
	convert.FrameworkToGraphBool(data.RegistrationDisabled, requestBody.SetRegistrationDisabled)
	convert.FrameworkToGraphBool(data.FileVaultDisabled, requestBody.SetFileVaultDisabled)
	convert.FrameworkToGraphBool(
		data.ICloudDiagnosticsDisabled,
		requestBody.SetICloudDiagnosticsDisabled,
	)
	convert.FrameworkToGraphBool(data.PassCodeDisabled, requestBody.SetPassCodeDisabled)
	convert.FrameworkToGraphBool(data.ZoomDisabled, requestBody.SetZoomDisabled)
	convert.FrameworkToGraphBool(data.ICloudStorageDisabled, requestBody.SetICloudStorageDisabled)
	convert.FrameworkToGraphBool(
		data.ChooseYourLockScreenDisabled,
		requestBody.SetChooseYourLockScreenDisabled,
	)
	convert.FrameworkToGraphBool(
		data.AccessibilityScreenDisabled,
		requestBody.SetAccessibilityScreenDisabled,
	)
	convert.FrameworkToGraphBool(
		data.AutoUnlockWithWatchDisabled,
		requestBody.SetAutoUnlockWithWatchDisabled,
	)
	convert.FrameworkToGraphBool(
		data.AutoAdvanceSetupEnabled,
		requestBody.SetAutoAdvanceSetupEnabled,
	)
	convert.FrameworkToGraphBool(
		data.RequestRequiresNetworkTether,
		requestBody.SetRequestRequiresNetworkTether,
	)
	convert.FrameworkToGraphBool(
		data.UsePlatformSSODuringSetupAssistant,
		requestBody.SetUsePlatformSSODuringSetupAssistant,
	)

	// Primary (managed local) account auto-creation
	convert.FrameworkToGraphBool(
		data.SkipPrimarySetupAccountCreation,
		requestBody.SetSkipPrimarySetupAccountCreation,
	)
	convert.FrameworkToGraphBool(
		data.SetPrimarySetupAccountAsRegularUser,
		requestBody.SetSetPrimarySetupAccountAsRegularUser,
	)
	convert.FrameworkToGraphBool(
		data.DontAutoPopulatePrimaryAccountInfo,
		requestBody.SetDontAutoPopulatePrimaryAccountInfo,
	)
	convert.FrameworkToGraphString(
		data.PrimaryAccountFullName,
		requestBody.SetPrimaryAccountFullName,
	)
	convert.FrameworkToGraphString(
		data.PrimaryAccountUserName,
		requestBody.SetPrimaryAccountUserName,
	)
	convert.FrameworkToGraphBool(data.EnableRestrictEditing, requestBody.SetEnableRestrictEditing)

	// Admin (local) account auto-creation. Use frameworkToGraphStringOrNil so removing
	// these from config clears them server-side on update.
	frameworkToGraphStringOrNil(
		data.AdminAccountUserName,
		requestBody.SetAdminAccountUserName,
	)
	frameworkToGraphStringOrNil(
		data.AdminAccountFullName,
		requestBody.SetAdminAccountFullName,
	)
	frameworkToGraphStringOrNil(
		data.AdminAccountPassword,
		requestBody.SetAdminAccountPassword,
	)
	convert.FrameworkToGraphBool(data.HideAdminAccount, requestBody.SetHideAdminAccount)

	if data.AdminAccountPasswordRotation != nil {
		rotation := graphmodels.NewDepProfileAdminAccountPasswordRotationSetting()
		convert.FrameworkToGraphInt32(
			data.AdminAccountPasswordRotation.AutoRotationPeriodInDays,
			rotation.SetAutoRotationPeriodInDays,
		)

		delay := graphmodels.NewDepProfileDelayAutoRotationSetting()
		convert.FrameworkToGraphBool(
			data.AdminAccountPasswordRotation.OnRetrievalAutoRotatePasswordEnabled,
			delay.SetOnRetrievalAutoRotatePasswordEnabled,
		)
		convert.FrameworkToGraphInt32(
			data.AdminAccountPasswordRotation.OnRetrievalDelayAutoRotatePasswordInHours,
			delay.SetOnRetrievalDelayAutoRotatePasswordInHours,
		)
		rotation.SetDepProfileDelayAutoRotationSetting(delay)

		requestBody.SetDepProfileAdminAccountPasswordRotationSetting(rotation)
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// buildEnabledSkipKeys constructs the enabledSkipKeys array from the boolean fields.
// Privacy and Registration are intentionally excluded because Microsoft Graph rejects
// those skip-key strings, even though the privacy_pane_disabled and registration_disabled
// boolean properties work correctly.
func buildEnabledSkipKeys(data *MacOSDepEnrollmentProfileResourceModel) []string {
	mappings := []struct {
		disabled types.Bool
		key      string
	}{
		// Base profile screens (Privacy excluded; Graph rejects the "Privacy" skip key)
		{data.AppleIdDisabled, "AppleID"},
		{data.ApplePayDisabled, "Payment"},
		{data.DiagnosticsDisabled, "Diagnostics"},
		{data.DisplayToneSetupDisabled, "DisplayTone"},
		{data.LocationDisabled, "Location"},
		{data.RestoreBlocked, "Restore"},
		{data.ScreenTimeScreenDisabled, "ScreenTime"},
		{data.SiriDisabled, "Siri"},
		{data.TermsAndConditionsDisabled, "TOS"},
		{data.TouchIdDisabled, "Biometric"},
		// macOS-specific screens (Registration excluded; Graph rejects the "Registration" skip key)
		{data.WelcomeScreenDisabled, "Welcome"},
		{data.AccessibilityScreenDisabled, "Accessibility"},
		{data.AutoUnlockWithWatchDisabled, "UnlockWithWatch"},
		{data.ChooseYourLockScreenDisabled, "Wallpaper"},
		{data.FileVaultDisabled, "FileVault"},
		{data.ICloudDiagnosticsDisabled, "iCloudDiagnostics"},
		{data.ICloudStorageDisabled, "iCloudStorage"},
		{data.PassCodeDisabled, "Passcode"},
		{data.ZoomDisabled, "Zoom"},
	}

	skipKeys := make([]string, 0, len(mappings))
	for _, m := range mappings {
		if m.disabled.ValueBool() {
			skipKeys = append(skipKeys, m.key)
		}
	}
	return skipKeys
}

// frameworkToGraphStringOrNil sets a Graph SDK string property, explicitly sending nil when
// the value is null or empty string so the field is cleared in the API when removed from config.
func frameworkToGraphStringOrNil(value types.String, setter func(*string)) {
	if value.IsUnknown() {
		return
	}
	if value.IsNull() || value.ValueString() == "" {
		setter(nil)
	} else {
		convert.FrameworkToGraphString(value, setter)
	}
}

// constructAzureAdGroupIds converts a Terraform set of GUID strings into []uuid.UUID.
func constructAzureAdGroupIds(ctx context.Context, set types.Set, setter func([]uuid.UUID)) error {
	if set.IsNull() || set.IsUnknown() {
		setter(nil)
		return nil
	}

	elements := set.Elements()
	result := make([]uuid.UUID, 0, len(elements))
	for i, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			return fmt.Errorf("%w at index %d: %T", errUnexpectedElementType, i, elem)
		}
		if strVal.IsNull() || strVal.IsUnknown() || strVal.ValueString() == "" {
			continue
		}
		parsed, err := uuid.Parse(strVal.ValueString())
		if err != nil {
			return fmt.Errorf("invalid GUID %q at index %d: %w", strVal.ValueString(), i, err)
		}
		result = append(result, parsed)
	}
	setter(result)
	return nil
}
