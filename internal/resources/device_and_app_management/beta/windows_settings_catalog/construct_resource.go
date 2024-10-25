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

					if settingData.SettingInstance.ChoiceSettingValue != nil {
						if !settingData.SettingInstance.ChoiceSettingValue.IntValue.IsNull() && !settingData.SettingInstance.ChoiceSettingValue.IntValue.IsUnknown() {
							intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
							construct.SetInt32Property(settingData.SettingInstance.ChoiceSettingValue.IntValue, intValue.SetValue)
							instance.SetSimpleSettingValue(intValue)
						} else if !settingData.SettingInstance.ChoiceSettingValue.StringValue.IsNull() && !settingData.SettingInstance.ChoiceSettingValue.StringValue.IsUnknown() {
							stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
							construct.SetStringProperty(settingData.SettingInstance.ChoiceSettingValue.StringValue, stringValue.SetValue)
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
						if !settingData.SettingInstance.ChoiceSettingValue.IntValue.IsNull() && !settingData.SettingInstance.ChoiceSettingValue.IntValue.IsUnknown() {
							intVal := strconv.Itoa(int(settingData.SettingInstance.ChoiceSettingValue.IntValue.ValueInt32()))
							choiceValue.SetValue(&intVal)
						} else if !settingData.SettingInstance.ChoiceSettingValue.StringValue.IsNull() && !settingData.SettingInstance.ChoiceSettingValue.StringValue.IsUnknown() {
							strVal := settingData.SettingInstance.ChoiceSettingValue.StringValue.ValueString()
							choiceValue.SetValue(&strVal)
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
							if !child.ChoiceSettingValue.IntValue.IsNull() && !child.ChoiceSettingValue.IntValue.IsUnknown() {
								// Convert Int32 value to string for setting
								intStr := strconv.Itoa(int(child.ChoiceSettingValue.IntValue.ValueInt32()))
								stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								stringValue.SetValue(&intStr)
								collectionValues = append(collectionValues, stringValue)
							} else if !child.ChoiceSettingValue.StringValue.IsNull() && !child.ChoiceSettingValue.StringValue.IsUnknown() {
								// Use String value directly if set
								strVal := child.ChoiceSettingValue.StringValue.ValueString()
								stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								stringValue.SetValue(&strVal)
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
