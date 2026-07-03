package graphBetaMacOSDeviceEnrollmentPolicy

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
	stateModel *MacOSDeviceEnrollmentPolicyResourceModel,
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
		case SettingDefAwaitConfiguration:
			mapAwaitConfigurationToState(instance, stateModel)
		case SettingDefLockedEnrollment:
			if value, _, ok := extractChoiceInstance(instance); ok {
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

// mapUserAffinityToState parses ade_macos_useraffinity and its ade_macos_authenticationmethod child.
func mapUserAffinityToState(instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *MacOSDeviceEnrollmentPolicyResourceModel) {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return
	}

	stateModel.RequiresUserAuthentication = types.BoolValue(boolFromChoiceValue(SettingDefUserAffinity, value))
	stateModel.EnableAuthenticationViaCompanyPortal = types.BoolValue(false)
	stateModel.RequireCompanyPortalOnSetupAssistantEnrolledDevices = types.BoolValue(false)

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil || *defIdPtr != SettingDefAuthenticationMethod {
			continue
		}
		authValue, _, ok := extractChoiceInstance(child)
		if !ok {
			continue
		}
		switch authValue {
		case AuthenticationMethodCompanyPortal:
			stateModel.EnableAuthenticationViaCompanyPortal = types.BoolValue(true)
		case AuthenticationMethodCompanyPortalOnSetupAssistant:
			stateModel.RequireCompanyPortalOnSetupAssistantEnrolledDevices = types.BoolValue(true)
		}
	}
}

// mapAwaitConfigurationToState parses ade_macos_awaitconfiguration and its local admin account subtree.
func mapAwaitConfigurationToState(instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *MacOSDeviceEnrollmentPolicyResourceModel) {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return
	}

	awaitConfigured := boolFromChoiceValue(SettingDefAwaitConfiguration, value)
	stateModel.AwaitDeviceConfigured = types.BoolValue(awaitConfigured)

	if !awaitConfigured {
		stateModel.AdminAccount = nil
		return
	}

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil || *defIdPtr != SettingDefCreateLocalAdmin {
			continue
		}
		stateModel.AdminAccount = mapCreateLocalAdminToState(child)
	}
}

// mapCreateLocalAdminToState parses ade_accountsettings_createlocaladmin and its children.
func mapCreateLocalAdminToState(instance models.DeviceManagementConfigurationSettingInstanceable) *AdminAccountModel {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return nil
	}

	admin := &AdminAccountModel{
		CreateLocalAdminAccount:   types.BoolValue(boolFromChoiceValue(SettingDefCreateLocalAdmin, value)),
		UserName:                  types.StringValue(""),
		FullName:                  types.StringValue(""),
		HideAccount:               types.BoolValue(false),
		PasswordRotationInDays:    types.Int64Value(0),
		CreateLocalPrimaryAccount: types.BoolValue(false),
	}

	if !admin.CreateLocalAdminAccount.ValueBool() {
		return admin
	}

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil {
			continue
		}

		switch *defIdPtr {
		case SettingDefAdminAccountName:
			if v, ok := extractSimpleStringInstance(child); ok {
				admin.UserName = types.StringValue(v)
			}
		case SettingDefAdminAccountFullName:
			if v, ok := extractSimpleStringInstance(child); ok {
				admin.FullName = types.StringValue(v)
			}
		case SettingDefHideUsersGroups:
			if v, _, ok := extractChoiceInstance(child); ok {
				admin.HideAccount = types.BoolValue(boolFromChoiceValue(SettingDefHideUsersGroups, v))
			}
		case SettingDefAdminAccountPasswordRotation:
			if v, ok := extractSimpleIntInstance(child); ok {
				admin.PasswordRotationInDays = types.Int64Value(v)
			}
		case SettingDefCreateLocalPrimary:
			primaryValue, primaryChildren, ok := extractChoiceInstance(child)
			if !ok {
				continue
			}
			createPrimary := boolFromChoiceValue(SettingDefCreateLocalPrimary, primaryValue)
			admin.CreateLocalPrimaryAccount = types.BoolValue(createPrimary)
			if createPrimary {
				for _, pc := range primaryChildren {
					if pc == nil {
						continue
					}
					pcDefId := pc.GetSettingDefinitionId()
					if pcDefId == nil || *pcDefId != SettingDefPrefillAccountInfo {
						continue
					}
					admin.PrimaryAccount = mapPrefillAccountInfoToState(pc)
				}
			}
		}
	}

	return admin
}

// mapPrefillAccountInfoToState parses ade_accountsettings_prefillaccountinfo and its children.
func mapPrefillAccountInfoToState(instance models.DeviceManagementConfigurationSettingInstanceable) *PrimaryAccountModel {
	value, children, ok := extractChoiceInstance(instance)
	if !ok {
		return nil
	}

	primary := &PrimaryAccountModel{
		PrefillAccountInfo: types.BoolValue(boolFromChoiceValue(SettingDefPrefillAccountInfo, value)),
		RestrictEditing:    types.BoolValue(false),
		FullName:           types.StringValue(""),
		UserName:           types.StringValue(""),
	}

	if !primary.PrefillAccountInfo.ValueBool() {
		return primary
	}

	for _, child := range children {
		if child == nil {
			continue
		}
		defIdPtr := child.GetSettingDefinitionId()
		if defIdPtr == nil {
			continue
		}

		switch *defIdPtr {
		case SettingDefRestrictEditing:
			if v, _, ok := extractChoiceInstance(child); ok {
				primary.RestrictEditing = types.BoolValue(boolFromChoiceValue(SettingDefRestrictEditing, v))
			}
		case SettingDefPrimaryAccountFullName:
			if v, ok := extractSimpleStringInstance(child); ok {
				primary.FullName = types.StringValue(v)
			}
		case SettingDefPrimaryAccountName:
			if v, ok := extractSimpleStringInstance(child); ok {
				primary.UserName = types.StringValue(v)
			}
		}
	}

	return primary
}

// mapSetupAssistantSettingToState parses one of the 23 Setup Assistant screen toggles.
func mapSetupAssistantSettingToState(settingDefinitionId string, instance models.DeviceManagementConfigurationSettingInstanceable, stateModel *MacOSDeviceEnrollmentPolicyResourceModel) {
	value, _, ok := extractChoiceInstance(instance)
	if !ok {
		return
	}
	disabled := types.BoolValue(!boolFromChoiceValue(settingDefinitionId, value))

	switch settingDefinitionId {
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
	case SettingDefFileVault:
		stateModel.FileVaultDisabled = disabled
	case SettingDefICloudDiagnostics:
		stateModel.ICloudDiagnosticsDisabled = disabled
	case SettingDefICloudStorage:
		stateModel.ICloudStorageDisabled = disabled
	case SettingDefAppearance:
		stateModel.DisplayToneSetupDisabled = disabled
	case SettingDefScreenTime:
		stateModel.ScreenTimeScreenDisabled = disabled
	case SettingDefPrivacy:
		stateModel.PrivacyPaneDisabled = disabled
	case SettingDefAccessibility:
		stateModel.AccessibilityScreenDisabled = disabled
	case SettingDefUnlockWithWatch:
		stateModel.AutoUnlockWithWatchDisabled = disabled
	case SettingDefEnableLockdownMode:
		stateModel.LockdownModeDisabled = disabled
	case SettingDefSoftwareUpdate:
		stateModel.SoftwareUpdateScreenDisabled = disabled
	case SettingDefSoftwareUpdateCompleted:
		stateModel.SoftwareUpdateCompletedScreenDisabled = disabled
	case SettingDefTermsOfAddress:
		stateModel.TermsOfAddressScreenDisabled = disabled
	case SettingDefIntelligence:
		stateModel.AppleIntelligenceDisabled = disabled
	case SettingDefOSShowcase:
		stateModel.OsShowcaseScreenDisabled = disabled
	case SettingDefAppStore:
		stateModel.AppStoreDisabled = disabled
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

// extractSimpleIntInstance returns the integer value of a simple setting instance.
func extractSimpleIntInstance(instance models.DeviceManagementConfigurationSettingInstanceable) (int64, bool) {
	simpleInstance, ok := instance.(models.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		return 0, false
	}
	simpleValue := simpleInstance.GetSimpleSettingValue()
	if simpleValue == nil {
		return 0, false
	}
	if intValue, ok := simpleValue.(models.DeviceManagementConfigurationIntegerSettingValueable); ok {
		if v := intValue.GetValue(); v != nil {
			return int64(*v), true
		}
	}
	return 0, false
}

// boolFromChoiceValue returns true when a choice value carries the `_1` (enabled) suffix for the
// given setting definition ID.
func boolFromChoiceValue(settingDefinitionId string, value string) bool {
	return strings.HasSuffix(value, settingDefinitionId+"_1")
}
