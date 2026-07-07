package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapSettingsToState parses the settings catalog tree returned by Graph back into the flat/typed
// Terraform schema, mirroring the settings constructed in construct.go. Every visionOS ADE
// setting is a flat, childless choice or string, so no nested parsing is required.
func mapSettingsToState(
	ctx context.Context,
	stateModel *VisionOSDeviceEnrollmentPolicyResourceModel,
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
			if value, ok := extractChoiceValue(instance); ok {
				stateModel.UserAffinity = types.BoolValue(boolFromChoiceValue(SettingDefUserAffinity, value))
			}
		case SettingDefAwaitConfiguration:
			if value, ok := extractChoiceValue(instance); ok {
				stateModel.AwaitDeviceConfigured = types.BoolValue(boolFromChoiceValue(SettingDefAwaitConfiguration, value))
			}
		case SettingDefLockedEnrollment:
			if value, ok := extractChoiceValue(instance); ok {
				stateModel.LockedEnrollmentEnabled = types.BoolValue(boolFromChoiceValue(SettingDefLockedEnrollment, value))
			}
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

// mapSetupAssistantSettingToState parses one of the 14 Setup Assistant screen toggles.
func mapSetupAssistantSettingToState(settingDefinitionId string, instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *VisionOSDeviceEnrollmentPolicyResourceModel) {
	value, ok := extractChoiceValue(instance)
	if !ok {
		return
	}
	disabled := types.BoolValue(!boolFromChoiceValue(settingDefinitionId, value))

	switch settingDefinitionId {
	case SettingDefAppleId:
		stateModel.AppleIdDisabled = disabled
	case SettingDefApplePay:
		stateModel.ApplePayDisabled = disabled
	case SettingDefDiagnosticsData:
		stateModel.DiagnosticsDisabled = disabled
	case SettingDefGetStarted:
		stateModel.GetStartedScreenDisabled = disabled
	case SettingDefIntelligence:
		stateModel.AppleIntelligenceDisabled = disabled
	case SettingDefLocationServices:
		stateModel.LocationServicesDisabled = disabled
	case SettingDefPasscode:
		stateModel.PasscodeDisabled = disabled
	case SettingDefPrivacy:
		stateModel.PrivacyPaneDisabled = disabled
	case SettingDefScreenTime:
		stateModel.ScreenTimeScreenDisabled = disabled
	case SettingDefSiri:
		stateModel.SiriDisabled = disabled
	case SettingDefSoftwareUpdate:
		stateModel.SoftwareUpdateScreenDisabled = disabled
	case SettingDefTermsAndConditions:
		stateModel.TermsAndConditionsDisabled = disabled
	case SettingDefTips:
		stateModel.TipsScreenDisabled = disabled
	case SettingDefTouchFaceId:
		stateModel.TouchIdDisabled = disabled
	}
}

// extractChoiceValue returns the choice value of a choice setting instance.
func extractChoiceValue(instance models.DeviceManagementConfigurationSettingInstanceable) (string, bool) {
	choiceInstance, ok := instance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return "", false
	}
	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return "", false
	}
	if v := choiceValue.GetValue(); v != nil {
		return *v, true
	}
	return "", false
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
