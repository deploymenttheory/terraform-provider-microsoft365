package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point like ConstructSettingsCatalogProfileRequestBody
func constructResource(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Windows Settings Catalog resource")
	construct.DebugPrintStruct(ctx, "Constructed Windows Settings Catalog Resource from model", data)

	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	displayName := data.DisplayName.ValueString()
	description := data.Description.ValueString()

	profile.SetName(&displayName)
	profile.SetDescription(&description)

	platforms := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	profile.SetPlatforms(&platforms)

	technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
	profile.SetTechnologies(&technologies)

	if len(data.RoleScopeTagIds) > 0 {
		var tagIds []string
		for _, tag := range data.RoleScopeTagIds {
			tagIds = append(tagIds, tag.ValueString())
		}
		profile.SetRoleScopeTagIds(tagIds)
	} else {
		profile.SetRoleScopeTagIds([]string{"0"})
	}

	// Construct settings like ConstructSettingsCatalogSettingsRequestBody
	settings := constructSettingsCatalogSettings(ctx, data.Settings)
	profile.SetSettings(settings)

	tflog.Debug(ctx, "Finished constructing Windows Settings Catalog resource")
	return profile, nil
}

// ConstructSettingsCatalogSettingsRequestBody constructs a slice of DeviceManagementConfigurationSettingable objects from a slice of DeviceManagementConfigurationSetting objects.
func constructSettingsCatalogSettings(ctx context.Context, settingConfigs []DeviceManagementConfigurationSettingResourceModel) []graphmodels.DeviceManagementConfigurationSettingable {
	var settings []graphmodels.DeviceManagementConfigurationSettingable

	for _, settingConfig := range settingConfigs {
		setting := graphmodels.NewDeviceManagementConfigurationSetting()
		if settingConfig.SettingInstance != nil {
			settingInstance := constructSettingInstance(settingConfig.SettingInstance)
			if settingInstance != nil {
				setting.SetSettingInstance(settingInstance)
				settings = append(settings, setting)
				tflog.Debug(ctx, fmt.Sprintf("Adding setting: %s", *settingInstance.GetSettingDefinitionId()))
			}
		}
	}

	return settings
}

// Like constructSettingInstance from the shared package
func constructSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	switch instanceConfig.ODataType.ValueString() {
	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
		instance.SetSettingDefinitionId(&settingDefId)

		var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

		if instanceConfig.ChoiceSettingValue != nil {
			if !instanceConfig.ChoiceSettingValue.IntValue.IsNull() {
				intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
				val := instanceConfig.ChoiceSettingValue.IntValue.ValueInt32()
				intValue.SetValue(&val)
				simpleValue = intValue
			} else if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
				stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				val := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
				stringValue.SetValue(&val)
				simpleValue = stringValue
			}
			if simpleValue != nil {
				instance.SetSimpleSettingValue(simpleValue)
			}
		}
		return instance

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
		instance.SetSettingDefinitionId(&settingDefId)

		if instanceConfig.ChoiceSettingValue != nil {
			choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			if !instanceConfig.ChoiceSettingValue.IntValue.IsNull() {
				strValue := strconv.Itoa(int(instanceConfig.ChoiceSettingValue.IntValue.ValueInt32()))
				choiceValue.SetValue(&strValue)
			} else if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
				strVal := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
				choiceValue.SetValue(&strVal)
			}
			instance.SetChoiceSettingValue(choiceValue)
		}
		return instance

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
		instance.SetSettingDefinitionId(&settingDefId)

		if instanceConfig.ChoiceSettingValue != nil && len(instanceConfig.ChoiceSettingValue.Children) > 0 {
			var collectionValues []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

			for _, child := range instanceConfig.ChoiceSettingValue.Children {
				if child.ChoiceSettingValue != nil {
					var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable
					if !child.ChoiceSettingValue.IntValue.IsNull() {
						strValue := strconv.Itoa(int(child.ChoiceSettingValue.IntValue.ValueInt32()))
						stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
						stringValue.SetValue(&strValue)
						simpleValue = stringValue
					} else if !child.ChoiceSettingValue.StringValue.IsNull() {
						strVal := child.ChoiceSettingValue.StringValue.ValueString()
						stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
						stringValue.SetValue(&strVal)
						simpleValue = stringValue
					}
					if simpleValue != nil {
						collectionValues = append(collectionValues, simpleValue)
					}
				}
			}

			if len(collectionValues) > 0 {
				instance.SetSimpleSettingCollectionValue(collectionValues)
			}
		}
		return instance
	}

	return nil
}
