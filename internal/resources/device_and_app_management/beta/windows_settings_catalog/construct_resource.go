package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Windows Settings Catalog resource")
	construct.DebugPrintStruct(ctx, "Constructed Windows Settings Catalog Resource from model", data)

	// Initialize profile object
	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	// Set basic properties from data model
	displayName := data.DisplayName.ValueString()
	description := data.Description.ValueString()
	profile.SetName(&displayName)
	profile.SetDescription(&description)

	// Set platforms and technologies (static values for Windows 10 MDM)
	platforms := graphmodels.DeviceManagementConfigurationPlatforms(graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS)
	profile.SetPlatforms(&platforms)
	technologies := graphmodels.DeviceManagementConfigurationTechnologies(graphmodels.MDM_DEVICEMANAGEMENTCONFIGURATIONTECHNOLOGIES)
	profile.SetTechnologies(&technologies)

	// Handle Role Scope Tag IDs
	if len(data.RoleScopeTagIds) > 0 {
		var tagIds []string
		for _, tag := range data.RoleScopeTagIds {
			tagIds = append(tagIds, tag.ValueString())
		}
		profile.SetRoleScopeTagIds(tagIds)
	} else {
		profile.SetRoleScopeTagIds([]string{"0"})
	}

	// Construct settings and set them to profile
	settings := constructSettingsCatalogSettings(ctx, data.Settings)
	profile.SetSettings(settings)

	tflog.Debug(ctx, "Finished constructing Windows Settings Catalog resource")
	return profile, nil
}

// Constructs settings catalog settings by processing each provided setting configuration.
func constructSettingsCatalogSettings(ctx context.Context, settingConfigs []DeviceManagementConfigurationSettingResourceModel) []graphmodels.DeviceManagementConfigurationSettingable {
	var settings []graphmodels.DeviceManagementConfigurationSettingable

	for _, settingConfig := range settingConfigs {
		setting := graphmodels.NewDeviceManagementConfigurationSetting()
		if settingConfig.SettingInstance != nil {
			// Build setting instance based on the type and configuration
			settingInstance := constructSettingInstance(settingConfig.SettingInstance)
			if settingInstance != nil {
				setting.SetSettingInstance(settingInstance)
				settings = append(settings, setting)
				tflog.Debug(ctx, fmt.Sprintf("Added setting with definition ID: %s", *settingInstance.GetSettingDefinitionId()))
			}
		}
	}
	return settings
}

// Constructs a setting instance based on its ODataType, properly mapping values for each supported instance type.
func constructSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	// Determine the setting type from ODataType and construct accordingly
	switch instanceConfig.ODataType.ValueString() {
	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
		instance.SetSettingDefinitionId(&settingDefId)

		if instanceConfig.ChoiceSettingValue != nil {
			var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

			if !instanceConfig.ChoiceSettingValue.IntValue.IsNull() {
				// Handle integer values
				intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
				val := instanceConfig.ChoiceSettingValue.IntValue.ValueInt32()
				intValue.SetValue(&val)
				simpleValue = intValue
			} else if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
				// Handle string values
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

		// Set choice setting value
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

		// Handle collections
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

	// Unsupported type
	return nil
}
