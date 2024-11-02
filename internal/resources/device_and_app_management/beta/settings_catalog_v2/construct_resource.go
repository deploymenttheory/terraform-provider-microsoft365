package graphBetaSettingsCatalog

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

	platformStr := data.Platforms.ValueString()
	var platform graphmodels.DeviceManagementConfigurationPlatforms
	switch platformStr {
	case "android":
		platform = graphmodels.ANDROID_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "androidEnterprise":
		platform = graphmodels.ANDROIDENTERPRISE_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "aosp":
		platform = graphmodels.AOSP_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "iOS":
		platform = graphmodels.IOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "linux":
		platform = graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "macOS":
		platform = graphmodels.MACOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "windows10":
		platform = graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "windows10X":
		platform = graphmodels.WINDOWS10X_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	}
	profile.SetPlatforms(&platform)

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
	case DeviceManagementConfigurationSimpleSettingInstance:
		return buildSimpleSettingInstance(instanceConfig)

	case DeviceManagementConfigurationChoiceSettingInstance:
		return buildChoiceSettingInstance(instanceConfig)

	case DeviceManagementConfigurationSimpleSettingCollectionInstance:
		return buildSimpleSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationChoiceSettingCollectionInstance:
		return buildChoiceSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationSettingGroupInstance:
		return buildGroupSettingInstance(instanceConfig)

	case DeviceManagementConfigurationSettingGroupCollectionInstance:
		return buildGroupSettingCollectionInstance(instanceConfig)

	// Unsupported type
	default:
		return nil
	}
}

func buildSimpleSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	if instanceConfig.ChoiceSettingValue != nil {
		// Process simple values
		var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable
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
}

// Example builder function for ChoiceSettingInstance with recursion
func buildChoiceSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	if instanceConfig.ChoiceSettingValue != nil {
		choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
		if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
			val := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
			choiceValue.SetValue(&val)
		}

		// Recursively add children
		var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			childInstance := constructSettingInstance(child) // Pass child directly without &
			if childInstance != nil {
				childInstances = append(childInstances, childInstance)
			}
		}
		choiceValue.SetChildren(childInstances)
		instance.SetChoiceSettingValue(choiceValue)
	}

	return instance
}

func buildSimpleSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

func buildChoiceSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	if instanceConfig.ChoiceSettingValue != nil && len(instanceConfig.ChoiceSettingValue.Children) > 0 {
		var collectionValues []graphmodels.DeviceManagementConfigurationChoiceSettingValueable

		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			if child.ChoiceSettingValue != nil {
				choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
				if !child.ChoiceSettingValue.StringValue.IsNull() {
					val := child.ChoiceSettingValue.StringValue.ValueString()
					choiceValue.SetValue(&val)
				}
				collectionValues = append(collectionValues, choiceValue)
			}
		}

		if len(collectionValues) > 0 {
			instance.SetChoiceSettingCollectionValue(collectionValues)
		}
	}
	return instance
}

func buildGroupSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingInstance()
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	// Check if ChoiceSettingValue has children
	if instanceConfig.ChoiceSettingValue != nil && instanceConfig.ChoiceSettingValue.Children != nil {
		var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			// Recursively construct each child instance
			childInstance := constructSettingInstance(child)
			if childInstance != nil {
				childInstances = append(childInstances, childInstance)
			}
		}
		// Set the children for the group setting instance
		groupSettingValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
		groupSettingValue.SetChildren(childInstances)
		instance.SetGroupSettingValue(groupSettingValue)
	}

	return instance
}

func buildGroupSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	settingDefId := instanceConfig.SettingDefinitionID.ValueString() // Accessing SettingDefinitionID directly
	instance.SetSettingDefinitionId(&settingDefId)

	// Check if instanceConfig has ChoiceSettingValue and iterate through its Children if present
	if instanceConfig.ChoiceSettingValue != nil && instanceConfig.ChoiceSettingValue.Children != nil {
		var collectionValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable

		// Iterate over each child in ChoiceSettingValue.Children to construct a DeviceManagementConfigurationGroupSettingValue
		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()

			// Recursively construct the child instances for this group
			childInstance := constructSettingInstance(child)
			if childInstance != nil {
				groupValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{childInstance})
			}

			// Append the constructed group value to the collection
			collectionValues = append(collectionValues, groupValue)
		}

		// Set the constructed collection on the GroupSettingCollectionInstance
		instance.SetGroupSettingCollectionValue(collectionValues)
	}

	return instance
}

// // Constructs a setting instance based on its ODataType, properly mapping values for each supported instance type.
// func constructSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
// 	// Determine the setting type from ODataType and construct accordingly
// 	switch instanceConfig.ODataType.ValueString() {
// 	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
// 		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
// 		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
// 		instance.SetSettingDefinitionId(&settingDefId)

// 		if instanceConfig.ChoiceSettingValue != nil {
// 			var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

// 			if !instanceConfig.ChoiceSettingValue.IntValue.IsNull() {
// 				// Handle integer values
// 				intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
// 				val := instanceConfig.ChoiceSettingValue.IntValue.ValueInt32()
// 				intValue.SetValue(&val)
// 				simpleValue = intValue
// 			} else if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
// 				// Handle string values
// 				stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
// 				val := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
// 				stringValue.SetValue(&val)
// 				simpleValue = stringValue
// 			}

// 			if simpleValue != nil {
// 				instance.SetSimpleSettingValue(simpleValue)
// 			}
// 		}
// 		return instance

// 	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
// 		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
// 		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
// 		instance.SetSettingDefinitionId(&settingDefId)

// 		// Set choice setting value
// 		if instanceConfig.ChoiceSettingValue != nil {
// 			choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()

// 			if !instanceConfig.ChoiceSettingValue.IntValue.IsNull() {
// 				strValue := strconv.Itoa(int(instanceConfig.ChoiceSettingValue.IntValue.ValueInt32()))
// 				choiceValue.SetValue(&strValue)
// 			} else if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
// 				strVal := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
// 				choiceValue.SetValue(&strVal)
// 			}

// 			instance.SetChoiceSettingValue(choiceValue)
// 		}
// 		return instance

// 	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
// 		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
// 		settingDefId := instanceConfig.SettingDefinitionID.ValueString()
// 		instance.SetSettingDefinitionId(&settingDefId)

// 		// Handle collections
// 		if instanceConfig.ChoiceSettingValue != nil && len(instanceConfig.ChoiceSettingValue.Children) > 0 {
// 			var collectionValues []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

// 			for _, child := range instanceConfig.ChoiceSettingValue.Children {
// 				if child.ChoiceSettingValue != nil {
// 					var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable
// 					if !child.ChoiceSettingValue.IntValue.IsNull() {
// 						strValue := strconv.Itoa(int(child.ChoiceSettingValue.IntValue.ValueInt32()))
// 						stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
// 						stringValue.SetValue(&strValue)
// 						simpleValue = stringValue
// 					} else if !child.ChoiceSettingValue.StringValue.IsNull() {
// 						strVal := child.ChoiceSettingValue.StringValue.ValueString()
// 						stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
// 						stringValue.SetValue(&strVal)
// 						simpleValue = stringValue
// 					}
// 					if simpleValue != nil {
// 						collectionValues = append(collectionValues, simpleValue)
// 					}
// 				}
// 			}

// 			if len(collectionValues) > 0 {
// 				instance.SetSimpleSettingCollectionValue(collectionValues)
// 			}
// 		}
// 		return instance
// 	}

// 	// Unsupported type
// 	return nil
// }
