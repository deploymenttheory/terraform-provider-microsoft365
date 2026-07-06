package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapSettingsToState parses the settings catalog tree returned by Graph back into the flat/typed
// Terraform schema, mirroring the nesting constructed in construct.go.
func mapSettingsToState(
	ctx context.Context,
	stateModel *IOSiPadOSDeviceEnrollmentPolicyResourceModel,
	settingsResponse models.DeviceManagementConfigurationSettingCollectionResponseable,
) error {
	if settingsResponse == nil {
		tflog.Debug(ctx, "Settings response is nil")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "No settings found in response")
		return nil
	}

	for _, setting := range settings {
		if setting == nil {
			continue
		}
		instance := setting.GetSettingInstance()
		if instance == nil {
			continue
		}
		defIdPtr := instance.GetSettingDefinitionId()
		if defIdPtr == nil {
			continue
		}

		switch *defIdPtr {
		case SettingDefUserAffinity:
			mapUserAffinityToState(instance, stateModel)
		case SettingDefLockedEnrollment:
			if value, _, ok := extractChoiceInstance(instance); ok {
				stateModel.LockedEnrollmentEnabled = types.BoolValue(boolFromChoiceValue(SettingDefLockedEnrollment, value))
			}
		case SettingDefDeviceNameTemplateChoices:
			stateModel.DeviceNameTemplate = mapStringChoiceToState(instance, SettingDefDeviceNameTemplateChoices, SettingDefAppleDeviceNameTemplate)
		case SettingDefActivateCellularDataChoices:
			stateModel.CellularDataActivationUrl = mapStringChoiceToState(instance, SettingDefActivateCellularDataChoices, SettingDefActivateCellularData)
		case SettingDefDepartment:
			if v, ok := extractSimpleStringInstance(instance); ok {
				stateModel.SupportDepartment = types.StringValue(v)
			}
		case SettingDefDepartmentPhone:
			if v, ok := extractSimpleStringInstance(instance); ok {
				stateModel.SupportPhoneNumber = types.StringValue(v)
			}
		default:
			mapSetupAssistantSettingToState(*defIdPtr, instance, stateModel)
		}
	}

	return nil
}

// mapUserAffinityToState parses ade_useraffinity, its ade_authenticationmethod child, and the
// ade_modernauth_awaitfinalconfiguration grandchild nested under the modern authentication choice.
func mapUserAffinityToState(instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *IOSiPadOSDeviceEnrollmentPolicyResourceModel) {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return
	}

	stateModel.RequiresUserAuthentication = types.BoolValue(boolFromChoiceValue(SettingDefUserAffinity, value))
	stateModel.EnableAuthenticationViaCompanyPortal = types.BoolValue(false)
	stateModel.RequireSetupAssistantWithModernAuthentication = types.BoolValue(false)
	stateModel.AwaitFinalConfiguration = types.BoolValue(false)

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil || *defIdPtr != SettingDefAuthenticationMethod {
			continue
		}
		authValue, authChildren, ok := extractChoiceInstance(child)
		if !ok {
			continue
		}
		switch authValue {
		case AuthenticationMethodCompanyPortal:
			stateModel.EnableAuthenticationViaCompanyPortal = types.BoolValue(true)
		case AuthenticationMethodSetupAssistantModernAuth:
			stateModel.RequireSetupAssistantWithModernAuthentication = types.BoolValue(true)

			for _, authChild := range authChildren {
				if authChild == nil {
					continue
				}
				childDefIdPtr := authChild.GetSettingDefinitionId()
				if childDefIdPtr == nil || *childDefIdPtr != SettingDefAwaitFinalConfiguration {
					continue
				}
				if awaitValue, _, ok := extractChoiceInstance(authChild); ok {
					stateModel.AwaitFinalConfiguration = types.BoolValue(boolFromChoiceValue(SettingDefAwaitFinalConfiguration, awaitValue))
				}
			}
		}
	}
}

// mapStringChoiceToState parses one of the enable/disable choice settings carrying an optional
// string child (device name template, cellular data activation URL). Returns null when the choice
// is "_0" (not configured), so state round-trips as null matching an omitted config attribute.
func mapStringChoiceToState(instance models.DeviceManagementConfigurationSettingInstanceable, choiceSettingDefinitionId, childSettingDefinitionId string) types.String {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return types.StringNull()
	}

	if !boolFromChoiceValue(choiceSettingDefinitionId, value) {
		return types.StringNull()
	}

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil || *defIdPtr != childSettingDefinitionId {
			continue
		}
		if v, ok := extractSimpleStringInstance(child); ok {
			return types.StringValue(v)
		}
	}

	return types.StringNull()
}

// mapSetupAssistantSettingToState parses one of the 32 Setup Assistant screen toggles.
func mapSetupAssistantSettingToState(settingDefinitionId string, instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *IOSiPadOSDeviceEnrollmentPolicyResourceModel) {
	value, _, ok := extractChoiceInstance(instance)
	if !ok {
		return
	}
	disabled := types.BoolValue(!boolFromChoiceValue(settingDefinitionId, value))

	switch settingDefinitionId {
	case SettingDefPasscode:
		stateModel.PasscodeDisabled = disabled
	case SettingDefLocationServices:
		stateModel.LocationServicesDisabled = disabled
	case SettingDefRestore:
		stateModel.RestoreDisabled = disabled
	case SettingDefAppleId:
		stateModel.AppleIdDisabled = disabled
	case SettingDefTermsAndConditions:
		stateModel.TermsAndConditionsDisabled = disabled
	case SettingDefTouchFaceId:
		stateModel.TouchIdDisabled = disabled
	case SettingDefApplePay:
		stateModel.ApplePayDisabled = disabled
	case SettingDefSiri:
		stateModel.SiriDisabled = disabled
	case SettingDefDiagnosticsData:
		stateModel.DiagnosticsDisabled = disabled
	case SettingDefPrivacy:
		stateModel.PrivacyPaneDisabled = disabled
	case SettingDefAndroidMigration:
		stateModel.RestoreFromAndroidDisabled = disabled
	case SettingDefIMessageFaceTime:
		stateModel.IMessageAndFaceTimeDisabled = disabled
	case SettingDefScreenTime:
		stateModel.ScreenTimeScreenDisabled = disabled
	case SettingDefSimSetup:
		stateModel.SimSetupScreenDisabled = disabled
	case SettingDefSoftwareUpdate:
		stateModel.SoftwareUpdateScreenDisabled = disabled
	case SettingDefWatchMigration:
		stateModel.WatchMigrationScreenDisabled = disabled
	case SettingDefAppearance:
		stateModel.AppearanceScreenDisabled = disabled
	case SettingDefDeviceMigration:
		stateModel.DeviceToDeviceMigrationDisabled = disabled
	case SettingDefRestoreCompleted:
		stateModel.RestoreCompletedScreenDisabled = disabled
	case SettingDefSoftwareUpdateCompleted:
		stateModel.SoftwareUpdateCompletedScreenDisabled = disabled
	case SettingDefGetStarted:
		stateModel.GetStartedScreenDisabled = disabled
	case SettingDefActionButton:
		stateModel.ActionButtonScreenDisabled = disabled
	case SettingDefSafety:
		stateModel.SafetyScreenDisabled = disabled
	case SettingDefTermsOfAddress:
		stateModel.TermsOfAddressScreenDisabled = disabled
	case SettingDefIntelligence:
		stateModel.AppleIntelligenceDisabled = disabled
	case SettingDefEnableLockdownMode:
		stateModel.LockdownModeDisabled = disabled
	case SettingDefAppStore:
		stateModel.AppStoreDisabled = disabled
	case SettingDefCameraButton:
		stateModel.CameraButtonScreenDisabled = disabled
	case SettingDefMultitasking:
		stateModel.MultitaskingScreenDisabled = disabled
	case SettingDefOSShowcase:
		stateModel.OsShowcaseScreenDisabled = disabled
	case SettingDefSafetyAndHandling:
		stateModel.SafetyAndHandlingScreenDisabled = disabled
	case SettingDefWebContentFiltering:
		stateModel.WebContentFilteringDisabled = disabled
	}
}

// extractChoiceInstance returns the choice value and children of a choice setting instance.
func extractChoiceInstance(instance models.DeviceManagementConfigurationSettingInstanceable) (string, []models.DeviceManagementConfigurationSettingInstanceable, bool) {
	choiceInstance, ok := instance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return "", nil, false
	}
	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return "", nil, false
	}
	value := ""
	if v := choiceValue.GetValue(); v != nil {
		value = *v
	}
	return value, choiceValue.GetChildren(), true
}

// extractSimpleStringInstance returns the string value of a simple setting instance.
func extractSimpleStringInstance(instance models.DeviceManagementConfigurationSettingInstanceable) (string, bool) {
	simpleInstance, ok := instance.(models.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		return "", false
	}
	simpleValue := simpleInstance.GetSimpleSettingValue()
	if simpleValue == nil {
		return "", false
	}
	if stringValue, ok := simpleValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
		if v := stringValue.GetValue(); v != nil {
			return *v, true
		}
	}
	return "", false
}

// boolFromChoiceValue returns true when a choice value carries the `_1` (enabled) suffix for the
// given setting definition ID.
func boolFromChoiceValue(settingDefinitionId string, value string) bool {
	return strings.HasSuffix(value, settingDefinitionId+"_1")
}
