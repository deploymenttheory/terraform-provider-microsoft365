package graphBetaAppControlForBusinessBuiltInControls

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapAppControlSettingsToTerraform maps App Control settings from Graph API response to Terraform state
func MapAppControlSettingsToTerraform(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel, settingsResponse graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) error {
	// Initialize trust apps to empty set - will be updated if found in API response
	data.AdditionalRulesForTrustingApps = convert.GraphToFrameworkStringSet(ctx, []string{})

	if settingsResponse == nil {
		tflog.Debug(ctx, "No settings response data to process")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "Settings array is empty")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Processing %d settings from API response", len(settings)))

	for _, setting := range settings {
		if setting == nil {
			continue
		}

		settingInstance := setting.GetSettingInstance()
		if settingInstance == nil {
			continue
		}

		settingDefId := settingInstance.GetSettingDefinitionId()
		if settingDefId == nil {
			continue
		}

		// Look for the main App Control setting
		if *settingDefId == "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions" {
			err := extractAppControlSettings(ctx, data, settingInstance)
			if err != nil {
				return fmt.Errorf("failed to extract app control settings: %v", err)
			}
		}
	}

	return nil
}

// extractAppControlSettings extracts app control settings from the main setting instance
func extractAppControlSettings(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel, settingInstance graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	choiceInstance, ok := settingInstance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return fmt.Errorf("setting instance is not a choice setting")
	}

	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return nil
	}

	children := choiceValue.GetChildren()
	if children == nil {
		return nil
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		childDefId := child.GetSettingDefinitionId()
		if childDefId == nil {
			continue
		}

		if *childDefId == "device_vendor_msft_policy_config_applicationcontrol_built_in_controls" {
			err := extractBuiltInControlsSettings(ctx, data, child)
			if err != nil {
				return fmt.Errorf("failed to extract built-in controls settings: %v", err)
			}
		}
	}

	return nil
}

// extractBuiltInControlsSettings extracts settings from built-in controls group collection
func extractBuiltInControlsSettings(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel, child graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	groupCollection, ok := child.(graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable)
	if !ok {
		return fmt.Errorf("child setting is not a group collection")
	}

	collectionValue := groupCollection.GetGroupSettingCollectionValue()
	if len(collectionValue) == 0 {
		return nil
	}

	for _, groupValue := range collectionValue {
		if groupValue == nil {
			continue
		}

		groupChildren := groupValue.GetChildren()
		if groupChildren == nil {
			continue
		}

		for _, groupChild := range groupChildren {
			if groupChild == nil {
				continue
			}

			childDefId := groupChild.GetSettingDefinitionId()
			if childDefId == nil {
				continue
			}

			switch *childDefId {
			case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control":
				err := extractEnableAppControlSetting(ctx, data, groupChild)
				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to extract enable app control setting: %v", err))
				}
			case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps":
				err := extractTrustAppsSetting(ctx, data, groupChild)
				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to extract trust apps setting: %v", err))
				}
			}
		}
	}

	return nil
}

// extractEnableAppControlSetting extracts the enable_app_control setting value
func extractEnableAppControlSetting(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel, setting graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	choiceInstance, ok := setting.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return fmt.Errorf("setting instance is not a choice setting")
	}

	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return nil
	}

	value := choiceValue.GetValue()
	if value == nil {
		return nil
	}

	switch *value {
	case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_0":
		data.EnableAppControl = types.StringValue("audit")
	case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_1":
		data.EnableAppControl = types.StringValue("enforce")

	default:
		tflog.Warn(ctx, fmt.Sprintf("Unknown enable app control value: %s", *value))
	}

	return nil
}

// extractTrustAppsSetting extracts the trust_apps setting values
func extractTrustAppsSetting(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel, setting graphmodels.DeviceManagementConfigurationSettingInstanceable) error {
	choiceCollectionInstance, ok := setting.(graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable)
	if !ok {
		return fmt.Errorf("setting instance is not a choice collection setting")
	}

	collectionValue := choiceCollectionInstance.GetChoiceSettingCollectionValue()
	if collectionValue == nil {
		// Set empty set instead of leaving as null
		data.AdditionalRulesForTrustingApps = convert.GraphToFrameworkStringSet(ctx, []string{})
		return nil
	}

	var trustApps []string
	for _, choiceValue := range collectionValue {
		if choiceValue == nil {
			continue
		}

		value := choiceValue.GetValue()
		if value == nil {
			continue
		}

		switch *value {
		case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_0":
			trustApps = append(trustApps, "trust_apps_with_good_reputation")
		case "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_1":
			trustApps = append(trustApps, "trust_apps_from_managed_installers")
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown trust app value: %s", *value))
		}
	}

	data.AdditionalRulesForTrustingApps = convert.GraphToFrameworkStringSet(ctx, trustApps)
	return nil
}
