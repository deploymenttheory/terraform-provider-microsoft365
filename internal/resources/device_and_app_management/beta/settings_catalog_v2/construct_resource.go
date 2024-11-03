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
	// Top-level setting type for Device Management Configuration Setting
	DeviceManagementConfigurationSetting = "#microsoft.graph.deviceManagementConfigurationSetting"
	// Derived types for setting instances
	DeviceManagementConfigurationChoiceSettingInstance           = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
	DeviceManagementConfigurationChoiceSettingCollectionInstance = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
	DeviceManagementConfigurationSimpleSettingInstance           = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	DeviceManagementConfigurationSimpleSettingCollectionInstance = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
	DeviceManagementConfigurationSettingGroupInstance            = "#microsoft.graph.deviceManagementConfigurationSettingGroupInstance"
  DeviceManagementConfigurationGroupSettingInstance            = "#microsoft.graph.deviceManagementConfigurationGroupSettingInstance"
	DeviceManagementConfigurationSettingGroupCollectionInstance  = "#microsoft.graph.deviceManagementConfigurationSettingGroupCollectionInstance"
	DeviceManagementConfigurationGroupSettingCollectionInstance  = "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance"
	// Derived types from `Simple Setting Value Attributes`
	DeviceManagementConfigurationIntegerSettingValue   = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
	DeviceManagementConfigurationReferenceSettingValue = "#microsoft.graph.deviceManagementConfigurationReferenceSettingValue"
	DeviceManagementConfigurationSecretSettingValue    = "#microsoft.graph.deviceManagementConfigurationSecretSettingValue"
	DeviceManagementConfigurationStringSettingValue    = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	DeviceManagementConfigurationChoiceSettingValue    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
	// Derived types from `Setting Value Template Reference Attributes`
	DeviceManagementConfigurationSettingInstanceTemplateReferenceODataType = "#microsoft.graph.deviceManagementConfigurationSettingInstanceTemplateReference"
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
	// Set OData type explicitly
	oDataType := DeviceManagementConfigurationSimpleSettingInstance
	instance.SetOdataType(&oDataType)

	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	if instanceConfig.SimpleSettingValue != nil {
		var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

		if !instanceConfig.SimpleSettingValue.IntValue.IsNull() {
			intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
			intODataType := DeviceManagementConfigurationIntegerSettingValue
			intValue.SetOdataType(&intODataType)
			val := instanceConfig.SimpleSettingValue.IntValue.ValueInt32()
			intValue.SetValue(&val)
			simpleValue = intValue
		} else if !instanceConfig.SimpleSettingValue.StringValue.IsNull() {
			stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			strODataType := DeviceManagementConfigurationStringSettingValue
			stringValue.SetOdataType(&strODataType)
			val := instanceConfig.SimpleSettingValue.StringValue.ValueString()
			stringValue.SetValue(&val)
			simpleValue = stringValue
		}

		if simpleValue != nil {
			instance.SetSimpleSettingValue(simpleValue)
		}
	}

	// Set template reference if present
	if instanceConfig.SettingInstanceTemplateReference != nil {
		templateRef := buildTemplateReference(instanceConfig.SettingInstanceTemplateReference)
		instance.SetSettingInstanceTemplateReference(templateRef)
	}

	return instance
}

// buildChoiceSettingInstance constructs a DeviceManagementConfigurationChoiceSettingInstance from the provided configuration.
func buildChoiceSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {

	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()

	choiceInstanceODataType := DeviceManagementConfigurationChoiceSettingInstance
	instance.SetOdataType(&choiceInstanceODataType)

	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	// Check and set the choice setting value if it exists
	if instanceConfig.ChoiceSettingValue != nil {
		choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
		choiceODataType := DeviceManagementConfigurationChoiceSettingValue
		choiceValue.SetOdataType(&choiceODataType)

		// Assign the choice value string directly from instanceConfig if present
		if !instanceConfig.ChoiceSettingValue.StringValue.IsNull() {
			val := instanceConfig.ChoiceSettingValue.StringValue.ValueString()
			choiceValue.SetValue(&val)
		}

		// Recursively construct child instances if available
		if len(instanceConfig.ChoiceSettingValue.Children) > 0 {
			childInstances := make([]graphmodels.DeviceManagementConfigurationSettingInstanceable, 0)
			for _, child := range instanceConfig.ChoiceSettingValue.Children {
				childInstance := constructSettingInstance(&child)
				if childInstance != nil {
					childInstances = append(childInstances, childInstance)
				}
			}
			choiceValue.SetChildren(childInstances)
		}

		instance.SetChoiceSettingValue(choiceValue)
	}

	if instanceConfig.SettingInstanceTemplateReference != nil {
		templateRef := buildTemplateReference(instanceConfig.SettingInstanceTemplateReference)
		instance.SetSettingInstanceTemplateReference(templateRef)
	}

	return instance
}

// buildSimpleSettingCollectionInstance constructs a DeviceManagementConfigurationSimpleSettingCollectionInstance from the provided configuration.
func buildSimpleSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {

	instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()

	simpleCollectionODataType := DeviceManagementConfigurationSimpleSettingCollectionInstance
	instance.SetOdataType(&simpleCollectionODataType)
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	// Check if ChoiceSettingValue has child values
	if instanceConfig.ChoiceSettingValue != nil && len(instanceConfig.ChoiceSettingValue.Children) > 0 {
		var collectionValues []graphmodels.DeviceManagementConfigurationSimpleSettingValueable

		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			if child.ChoiceSettingValue != nil {
				var simpleValue graphmodels.DeviceManagementConfigurationSimpleSettingValueable

				if !child.ChoiceSettingValue.IntValue.IsNull() {
					strValue := strconv.Itoa(int(child.ChoiceSettingValue.IntValue.ValueInt32()))
					stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
					stringODataType := DeviceManagementConfigurationStringSettingValue
					stringValue.SetOdataType(&stringODataType)
					stringValue.SetValue(&strValue)
					simpleValue = stringValue
				} else if !child.ChoiceSettingValue.StringValue.IsNull() {
					strVal := child.ChoiceSettingValue.StringValue.ValueString()
					stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
					stringODataType := DeviceManagementConfigurationStringSettingValue
					stringValue.SetOdataType(&stringODataType)
					stringValue.SetValue(&strVal)
					simpleValue = stringValue
				}

				if simpleValue != nil {
					collectionValues = append(collectionValues, simpleValue)
				}
			}
		}

		// Set the collection values on the instance
		if len(collectionValues) > 0 {
			instance.SetSimpleSettingCollectionValue(collectionValues)
		}
	}

	return instance
}

// buildChoiceSettingCollectionInstance constructs a DeviceManagementConfigurationChoiceSettingCollectionInstance from the provided configuration.
func buildChoiceSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {

	instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()

	choiceCollectionODataType := DeviceManagementConfigurationChoiceSettingCollectionInstance
	instance.SetOdataType(&choiceCollectionODataType)

	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	// Check if ChoiceSettingValue has child values
	if instanceConfig.ChoiceSettingValue != nil && len(instanceConfig.ChoiceSettingValue.Children) > 0 {
		var collectionValues []graphmodels.DeviceManagementConfigurationChoiceSettingValueable

		for _, child := range instanceConfig.ChoiceSettingValue.Children {
			if child.ChoiceSettingValue != nil {
				choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()

				choiceValueODataType := DeviceManagementConfigurationChoiceSettingValue
				choiceValue.SetOdataType(&choiceValueODataType)

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

// buildGroupSettingInstance constructs a DeviceManagementConfigurationGroupSettingInstance from the provided configuration.
func buildGroupSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	// Initialize the GroupSettingInstance directly
	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingInstance()

	// Set the setting definition ID directly from instanceConfig
	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	// Check if GroupSettingValue has child instances
	if instanceConfig.GroupSettingValue != nil && len(instanceConfig.GroupSettingValue.Children) > 0 {
		var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
		for _, child := range instanceConfig.GroupSettingValue.Children {
			// Recursively construct each child instance
			childInstance := constructSettingInstance(&child)
			if childInstance != nil {
				childInstances = append(childInstances, childInstance)
			}
		}

		// Initialize the GroupSettingValue and set children
		groupSettingValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
		groupSettingValue.SetChildren(childInstances)
		instance.SetGroupSettingValue(groupSettingValue)
	}

	return instance
}

// buildGroupSettingCollectionInstance constructs a DeviceManagementConfigurationGroupSettingCollectionInstance from the provided configuration.
func buildGroupSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	groupCollectionODataType := DeviceManagementConfigurationGroupSettingCollectionInstance
	instance.SetOdataType(&groupCollectionODataType)

	settingDefId := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefId)

	if instanceConfig.GroupCollectionValue != nil && len(instanceConfig.GroupCollectionValue.Children) > 0 {
		var collectionValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable

		for _, groupChild := range instanceConfig.GroupCollectionValue.Children {
			groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
			groupValueODataType := DeviceManagementConfigurationSettingGroupInstance
			groupValue.SetOdataType(&groupValueODataType)

			// Construct each child instance within GroupSettingValue
			var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
			childInstance := constructSettingInstance(&groupChild)
			if childInstance != nil {
				childInstances = append(childInstances, childInstance)
			}
			groupValue.SetChildren(childInstances)

			// Append the constructed group value to the collection
			collectionValues = append(collectionValues, groupValue)
		}

		// Set the constructed collection on the GroupSettingCollectionInstance
		instance.SetGroupSettingCollectionValue(collectionValues)
	}

	return instance
}

// Helper function to build template references
func buildTemplateReference(templateRef *DeviceManagementConfigurationTemplateReferenceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceTemplateReferenceable {
	reference := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()

	// Set the @odata.type using the constant
	groupCollectionODataType := DeviceManagementConfigurationSettingInstanceTemplateReferenceODataType
	reference.SetOdataType(&groupCollectionODataType)

	// Set the settingInstanceTemplateId directly from templateRef
	templateId := templateRef.SettingInstanceTemplateId.ValueString()
	reference.SetSettingInstanceTemplateId(&templateId)

	return reference
}
