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

	if err := convert.FrameworkToGraphStringSet(
		ctx,
		data.EnabledSkipKeys,
		requestBody.SetEnabledSkipKeys,
	); err != nil {
		return nil, fmt.Errorf("failed to set enabled_skip_keys: %w", err)
	}

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

	// Admin (local) account auto-creation
	convert.FrameworkToGraphString(data.AdminAccountUserName, requestBody.SetAdminAccountUserName)
	convert.FrameworkToGraphString(data.AdminAccountFullName, requestBody.SetAdminAccountFullName)
	convert.FrameworkToGraphString(data.AdminAccountPassword, requestBody.SetAdminAccountPassword)
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
