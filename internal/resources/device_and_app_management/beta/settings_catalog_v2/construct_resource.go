package graphBetaSettingsCatalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	// Top-level OData setting type for Device Management Configuration Setting
	DeviceManagementConfigurationSetting = "#microsoft.graph.deviceManagementConfigurationSetting"
	// Derived OData types for setting instances
	DeviceManagementConfigurationChoiceSettingInstance           = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	DeviceManagementConfigurationChoiceSettingCollectionInstance = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
	DeviceManagementConfigurationSimpleSettingInstance           = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	DeviceManagementConfigurationSimpleSettingCollectionInstance = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	DeviceManagementConfigurationSettingGroupInstance            = "#microsoft.graph.deviceManagementConfigurationSettingGroupInstance"
	DeviceManagementConfigurationGroupSettingInstance            = "#microsoft.graph.deviceManagementConfigurationGroupSettingInstance"
	DeviceManagementConfigurationSettingGroupCollectionInstance  = "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance"
	DeviceManagementConfigurationGroupSettingCollectionInstance  = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
	// Derived OData types for values of setting instances
	DeviceManagementConfigurationIntegerSettingValue   = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	DeviceManagementConfigurationReferenceSettingValue = "#microsoft.graph.deviceManagementConfigurationReferenceSettingValue"
	DeviceManagementConfigurationSecretSettingValue    = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
	DeviceManagementConfigurationStringSettingValue    = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	DeviceManagementConfigurationChoiceSettingValue    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	DeviceManagementConfigurationGroupSettingValue     = "#microsoft.graph.deviceManagementConfigurationGroupSettingValue"
	// Derived OData types for values template references
	DeviceManagementConfigurationSettingInstanceTemplateReferenceODataType = "#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference"
)

// Main entry point to construct the settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Windows Settings Catalog resource")
	construct.DebugPrintStruct(ctx, "Constructed Windows Settings Catalog Resource from model", data)

	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	Name := data.Name.ValueString()
	description := data.Description.ValueString()
	profile.SetName(&Name)
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

	case DeviceManagementConfigurationChoiceSettingInstance:
		return buildChoiceSettingInstance(instanceConfig)

	case DeviceManagementConfigurationChoiceSettingCollectionInstance:
		return buildChoiceSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationGroupSettingInstance:
		return buildGroupSettingInstance(instanceConfig)

	case DeviceManagementConfigurationGroupSettingCollectionInstance:
		return buildGroupSettingCollectionInstance(instanceConfig)

	// case DeviceManagementConfigurationSettingGroupInstance:
	// 	return buildSettingGroupInstance(instanceConfig)

	// case DeviceManagementConfigurationSettingGroupCollectionInstance:
	// 	return buildSettingGroupCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationSimpleSettingInstance:
		return buildSimpleSettingInstance(instanceConfig)

	case DeviceManagementConfigurationSimpleSettingCollectionInstance:
		return buildSimpleSettingCollectionInstance(instanceConfig)

	// Unsupported type
	default:
		return nil
	}
}

// buildSimpleSettingInstance constructs a simple setting instance from the configuration.
func buildSimpleSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.SimpleSettingValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()

	odataType := DeviceManagementConfigurationSimpleSettingInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	simpleValue := instanceConfig.SimpleSettingValue

	if !simpleValue.IntValue.IsNull() {
		value := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
		intValueODataType := DeviceManagementConfigurationIntegerSettingValue
		value.SetOdataType(&intValueODataType)
		intVal := simpleValue.IntValue.ValueInt32()
		value.SetValue(&intVal)
		instance.SetSimpleSettingValue(value)
	} else if !simpleValue.StringValue.IsNull() {
		value := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
		stringValueODataType := DeviceManagementConfigurationStringSettingValue
		value.SetOdataType(&stringValueODataType)
		strVal := simpleValue.StringValue.ValueString()
		value.SetValue(&strVal)
		instance.SetSimpleSettingValue(value)
	} else if !simpleValue.SecretValue.IsNull() {
		value := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
		secretValueODataType := DeviceManagementConfigurationSecretSettingValue
		value.SetOdataType(&secretValueODataType)
		secretVal := simpleValue.SecretValue.ValueString()
		value.SetValue(&secretVal)
		instance.SetSimpleSettingValue(value)
	}

	return instance
}

// buildChoiceSettingInstance constructs a choice setting instance from the configuration.
func buildChoiceSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.ChoiceSettingValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()

	odataType := DeviceManagementConfigurationChoiceSettingInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	choiceValue := instanceConfig.ChoiceSettingValue
	value := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
	choiceODataType := DeviceManagementConfigurationChoiceSettingValue
	value.SetOdataType(&choiceODataType)

	if !choiceValue.StringValue.IsNull() {
		val := choiceValue.StringValue.ValueString()
		value.SetValue(&val)
	}

	// if len(choiceValue.Children) > 0 {
	// 	var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
	// 	for _, child := range choiceValue.Children {
	// 		childInstance := buildChoiceSettingInstance(&child)
	// 		if childInstance != nil {
	// 			childInstances = append(childInstances, childInstance)
	// 		}
	// 	}
	// 	value.SetChildren(childInstances)
	// }

	instance.SetChoiceSettingValue(value)
	return instance

}

// buildSimpleSettingCollectionInstance constructs a simple collection setting instance from the configuration.
func buildSimpleSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.SimpleCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()

	odataType := DeviceManagementConfigurationSimpleSettingCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	collectionValue := instanceConfig.SimpleCollectionValue
	var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

	// Handle integer values collection
	if len(collectionValue.IntValue) > 0 {
		for _, intVal := range collectionValue.IntValue {
			if !intVal.IsNull() {
				value := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
				intValueODataType := DeviceManagementConfigurationIntegerSettingValue
				value.SetOdataType(&intValueODataType)
				val := intVal.ValueInt32()
				value.SetValue(&val)
				values = append(values, value)
			}
		}
	}

	// Handle string values collection
	if len(collectionValue.StringValue) > 0 {
		for _, stringVal := range collectionValue.StringValue {
			if !stringVal.IsNull() {
				value := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				stringValueODataType := DeviceManagementConfigurationStringSettingValue
				value.SetOdataType(&stringValueODataType)
				val := stringVal.ValueString()
				value.SetValue(&val)
				values = append(values, value)
			}
		}
	}

	// Handle secret value
	if !collectionValue.SecretValue.IsNull() {
		value := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
		secretValueODataType := DeviceManagementConfigurationSecretSettingValue
		value.SetOdataType(&secretValueODataType)
		secretVal := collectionValue.SecretValue.ValueString()
		value.SetValue(&secretVal)
		values = append(values, value)
	}

	if len(values) > 0 {
		instance.SetSimpleSettingCollectionValue(values)
	}

	return instance
}

// buildChoiceSettingCollectionInstance constructs a choice collection setting instance from the configuration.
func buildChoiceSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.ChoiceCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()

	odataType := DeviceManagementConfigurationChoiceSettingCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	collectionValue := instanceConfig.ChoiceCollectionValue
	var values []graphmodels.DeviceManagementConfigurationChoiceSettingValueable

	// Handle string values collection
	if len(collectionValue.StringValue) > 0 {
		for _, stringVal := range collectionValue.StringValue {
			if !stringVal.IsNull() {
				value := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
				choiceValueODataType := DeviceManagementConfigurationChoiceSettingValue
				value.SetOdataType(&choiceValueODataType)
				val := stringVal.ValueString()
				value.SetValue(&val)

				// if len(collectionValue.Children) > 0 {
				// 	var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
				// 	for _, child := range collectionValue.Children {
				// 		if childInstance := constructSettingInstance(&child); childInstance != nil {
				// 			childInstances = append(childInstances, childInstance)
				// 		}
				// 	}
				// 	value.SetChildren(childInstances)
				// }

				values = append(values, value)
			}
		}
	}

	// Handle integer values collection
	if len(collectionValue.IntValue) > 0 {
		for _, intVal := range collectionValue.IntValue {
			if !intVal.IsNull() {
				value := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
				choiceValueODataType := DeviceManagementConfigurationChoiceSettingValue
				value.SetOdataType(&choiceValueODataType)
				val := intVal.ValueInt32()
				strVal := strconv.FormatInt(int64(val), 10)
				value.SetValue(&strVal)

				values = append(values, value)
			}
		}
	}

	if len(values) > 0 {
		instance.SetChoiceSettingCollectionValue(values)
	}

	return instance
}

// buildGroupSettingInstance constructs a group setting instance from the configuration.
func buildGroupSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.GroupSettingValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingInstance()

	odataType := DeviceManagementConfigurationGroupSettingInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	value := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
	groupValueODataType := DeviceManagementConfigurationGroupSettingValue
	value.SetOdataType(&groupValueODataType)

	// if len(instanceConfig.GroupSettingValue.Children) > 0 {
	// 	var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
	// 	for _, child := range instanceConfig.GroupSettingValue.Children {
	// 		if childInstance := constructSettingInstance(&child); childInstance != nil {
	// 			childInstances = append(childInstances, childInstance)
	// 		}
	// 	}
	// 	value.SetChildren(childInstances)
	// }

	instance.SetGroupSettingValue(value)
	return instance
}

// buildGroupSettingCollectionInstance constructs a group collection setting instance from the configuration.
func buildGroupSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.GroupCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()

	odataType := DeviceManagementConfigurationGroupSettingCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	//collectionValue := instanceConfig.GroupCollectionValue
	var values []graphmodels.DeviceManagementConfigurationGroupSettingValueable

	// Handle children recursively
	// for _, child := range collectionValue.Children {
	// 	childInstance := constructSettingInstance(&child)
	// 	if childInstance != nil {
	// 		value := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
	// 		groupValueODataType := DeviceManagementConfigurationGroupSettingValue
	// 		value.SetOdataType(&groupValueODataType)
	// 		value.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{childInstance})
	// 		values = append(values, value)
	// 	}
	// }

	if len(values) > 0 {
		instance.SetGroupSettingCollectionValue(values)
	}

	return instance
}

// buildSettingGroupInstance constructs a setting group instance from the configuration.
func buildSettingGroupInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.GroupSettingValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSettingGroupInstance()

	odataType := DeviceManagementConfigurationSettingGroupInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	// For group settings, we use the base instance's additional data to store children
	// if len(instanceConfig.GroupSettingValue.Children) > 0 {
	// 	var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
	// 	for _, child := range instanceConfig.GroupSettingValue.Children {
	// 		if childInstance := constructSettingInstance(&child); childInstance != nil {
	// 			childInstances = append(childInstances, childInstance)
	// 		}
	// 	}

	// 	// Store children in additional data
	// 	additionalData := instance.GetAdditionalData()
	// 	additionalData["children"] = childInstances
	// 	instance.SetAdditionalData(additionalData)
	// }

	return instance
}

// buildSettingGroupCollectionInstance constructs a setting group collection instance from the configuration.
func buildSettingGroupCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.GroupCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSettingGroupCollectionInstance()

	odataType := DeviceManagementConfigurationSettingGroupCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	//collectionValue := instanceConfig.GroupCollectionValue
	// if len(collectionValue.Children) > 0 {
	// 	// Create an array of settings for the collection
	// 	var settingsArray []map[string]interface{}

	// 	for _, child := range collectionValue.Children {
	// 		if childInstance := constructSettingInstance(&child); childInstance != nil {
	// 			// Create a single group setting with its children
	// 			groupSetting := map[string]interface{}{
	// 				"@odata.type": DeviceManagementConfigurationGroupSettingValue,
	// 				"children":    []graphmodels.DeviceManagementConfigurationSettingInstanceable{childInstance},
	// 			}
	// 			settingsArray = append(settingsArray, groupSetting)
	// 		}
	// 	}

	// 	// Add the settings array to additional data
	// 	if len(settingsArray) > 0 {
	// 		additionalData := instance.GetAdditionalData()
	// 		if additionalData == nil {
	// 			additionalData = make(map[string]interface{})
	// 		}
	// 		additionalData["settingValues"] = settingsArray
	// 		instance.SetAdditionalData(additionalData)
	// 	}
	// }

	return instance
}
