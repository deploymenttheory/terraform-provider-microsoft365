package graphBetaWindowsSettingsCatalog

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing WindowsSettingsCatalog resource")
	construct.DebugPrintStruct(ctx, "Constructed WindowsSettingsCatalog Resource from model", data)

	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	construct.SetStringProperty(data.DisplayName, profile.SetName)
	construct.SetStringProperty(data.Description, profile.SetDescription)

	if len(data.RoleScopeTagIds) > 0 {
		construct.SetArrayProperty(data.RoleScopeTagIds, profile.SetRoleScopeTagIds)
	} else {
		profile.SetRoleScopeTagIds([]string{"0"})
	}

	if data.ID.IsNull() {
		platforms := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
		profile.SetPlatforms(&platforms)

		technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
		profile.SetTechnologies(&technologies)
	}

	// Handle settings
	if len(data.Settings) > 0 {
		var settings []graphmodels.DeviceManagementConfigurationSettingable

		for _, settingData := range data.Settings {
			if settingData.SettingInstance != nil {
				setting := graphmodels.NewDeviceManagementConfigurationSetting()

				// Construct setting instance based on type
				switch settingData.SettingInstance.ODataType.ValueString() {
				case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
					instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
					construct.SetStringProperty(settingData.SettingInstance.SettingDefinitionID, instance.SetSettingDefinitionId)

					// Handle choice setting value
					if settingData.SettingInstance.ChoiceSettingValue != nil {
						if settingData.SettingInstance.ChoiceSettingValue.IntValue != 0 {
							intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
							val := int32(settingData.SettingInstance.ChoiceSettingValue.IntValue)
							intValue.SetValue(&val)
							instance.SetSimpleSettingValue(intValue)
						} else if settingData.SettingInstance.ChoiceSettingValue.StringValue != "" {
							stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
							val := settingData.SettingInstance.ChoiceSettingValue.StringValue
							stringValue.SetValue(&val)
							instance.SetSimpleSettingValue(stringValue)
						}
					}
					setting.SetSettingInstance(instance)

				case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
					instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
					construct.SetStringProperty(settingData.SettingInstance.SettingDefinitionID, instance.SetSettingDefinitionId)

					// Handle choice setting value
					if settingData.SettingInstance.ChoiceSettingValue != nil {
						choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
						if settingData.SettingInstance.ChoiceSettingValue.IntValue != 0 {
							val := strconv.Itoa(int(settingData.SettingInstance.ChoiceSettingValue.IntValue))
							choiceValue.SetValue(&val)
						} else if settingData.SettingInstance.ChoiceSettingValue.StringValue != "" {
							val := settingData.SettingInstance.ChoiceSettingValue.StringValue
							choiceValue.SetValue(&val)
						}
						instance.SetChoiceSettingValue(choiceValue)
					}
					setting.SetSettingInstance(instance)

				case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
					instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
					construct.SetStringProperty(settingData.SettingInstance.SettingDefinitionID, instance.SetSettingDefinitionId)

					// Handle collection values
					if settingData.SettingInstance.ChoiceSettingValue != nil && len(settingData.SettingInstance.ChoiceSettingValue.Children) > 0 {
						var collectionValues []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

						for _, child := range settingData.SettingInstance.ChoiceSettingValue.Children {
							if child.ChoiceSettingValue.IntValue != 0 {
								intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
								val := int32(child.ChoiceSettingValue.IntValue)
								intValue.SetValue(&val)
								collectionValues = append(collectionValues, intValue)
							} else if child.ChoiceSettingValue.StringValue != "" {
								stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								val := child.ChoiceSettingValue.StringValue
								stringValue.SetValue(&val)
								collectionValues = append(collectionValues, stringValue)
							}
						}

						if len(collectionValues) > 0 {
							instance.SetSimpleSettingCollectionValue(collectionValues)
						}
					}
					setting.SetSettingInstance(instance)
				}

				settings = append(settings, setting)
			}
		}

		if len(settings) > 0 {
			profile.SetSettings(settings)
		}
	}

	tflog.Debug(ctx, "Finished constructing WindowsSettingsCatalog resource")
	return profile, nil
}
