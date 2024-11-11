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
func constructResource(ctx context.Context, data *SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Settings Catalog resource")
	construct.DebugPrintStruct(ctx, "Constructed Settings Catalog Resource from model", data)

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
func constructSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	// Determine the setting type from ODataType and construct accordingly
	switch instanceConfig.ODataType.ValueString() {

	case DeviceManagementConfigurationSimpleSettingInstance:
		return buildSimpleSettingInstance(instanceConfig)

	case DeviceManagementConfigurationSimpleSettingCollectionInstance:
		return buildSimpleSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationChoiceSettingInstance:
		return buildChoiceSettingInstance(instanceConfig)

	case DeviceManagementConfigurationChoiceSettingCollectionInstance:
		return buildChoiceSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationGroupSettingInstance:
		return buildGroupSettingInstance(instanceConfig)

	case DeviceManagementConfigurationGroupSettingCollectionInstance:
		return buildGroupSettingCollectionInstance(instanceConfig)

	case DeviceManagementConfigurationSettingGroupInstance:
		return buildSettingGroupInstance(instanceConfig)

	case DeviceManagementConfigurationSettingGroupCollectionInstance:
		return buildSettingGroupCollectionInstance(instanceConfig)

	// Unsupported type
	default:
		return nil
	}
}

// simple settings

// buildSimpleSettingInstance constructs a simple setting instance from the configuration.
func buildSimpleSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

		if !simpleValue.State.IsNull() {
			stateStr := simpleValue.State.ValueString()
			parsedState, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(stateStr)
			if err == nil && parsedState != nil {
				value.SetValueState(parsedState.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
			}
		}
		instance.SetSimpleSettingValue(value)
	}

	return instance
}

// buildSimpleSettingCollectionInstance constructs a simple collection setting instance from the configuration.
func buildSimpleSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

	// Handle secret values collection with state
	if len(collectionValue.SecretValue) > 0 {
		for i, secretVal := range collectionValue.SecretValue {
			if !secretVal.IsNull() {
				value := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
				secretValueODataType := DeviceManagementConfigurationSecretSettingValue
				value.SetOdataType(&secretValueODataType)
				val := secretVal.ValueString()
				value.SetValue(&val)

				if i < len(collectionValue.State) && !collectionValue.State[i].IsNull() {
					stateStr := collectionValue.State[i].ValueString()
					parsedState, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(stateStr)
					if err == nil && parsedState != nil {
						value.SetValueState(parsedState.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
					}
				}

				values = append(values, value)
			}
		}
	}

	// Set the collection of values
	if len(values) > 0 {
		instance.SetSimpleSettingCollectionValue(values)
	}

	return instance
}

// choice settings

// buildChoiceSettingInstance constructs a choice setting instance from the configuration.
func buildChoiceSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

	children := buildChoiceSettingChildren(choiceValue.Children)
	value.SetChildren(children)

	instance.SetChoiceSettingValue(value)
	return instance
}

// buildChoiceSettingChildren recursively constructs child choice setting instances.
// While the parent is a Choice Setting Instance, its children can be of any valid setting instance type
// (Choice, Simple, Group etc). Therefore, we must route each child back through the main constructSettingInstance
// function which acts as a type router, inspecting each child's ODataType and building it with the appropriate
// constructor. This matches the Graph API's schema where a Choice Setting's children array can contain
// heterogeneous setting instance types, each requiring different value structures and processing logic.
func buildChoiceSettingChildren(childrenConfig []DeviceManagementConfigurationSettingInstanceResourceModel) []graphmodels.DeviceManagementConfigurationSettingInstanceable {
	var children []graphmodels.DeviceManagementConfigurationSettingInstanceable

	for _, childConfig := range childrenConfig {
		childInstance := constructSettingInstance(&childConfig)
		if childInstance != nil {
			children = append(children, childInstance)
		}
	}
	return children
}

// buildChoiceSettingCollectionInstance constructs a choice collection setting instance from the configuration.
func buildChoiceSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

				buildChoiceSettingCollectionInstanceChildren(value, collectionValue.Children)

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

				buildChoiceSettingCollectionInstanceChildren(value, collectionValue.Children)

				values = append(values, value)
			}
		}
	}

	// Set the collection of values
	if len(values) > 0 {
		instance.SetChoiceSettingCollectionValue(values)
	}

	return instance
}

// buildChoiceSettingCollectionInstanceChildren adds children to a choice setting value if children exist in the configuration.
func buildChoiceSettingCollectionInstanceChildren(value *graphmodels.DeviceManagementConfigurationChoiceSettingValue, children []DeviceManagementConfigurationSettingInstanceResourceModel) {
	if len(children) > 0 {
		var childInstances []graphmodels.DeviceManagementConfigurationSettingInstanceable
		for _, childConfig := range children {
			childInstance := buildChoiceSettingCollectionInstance(&childConfig)
			if childInstance != nil {
				childInstances = append(childInstances, childInstance)
			}
		}
		value.SetChildren(childInstances)
	}
}

// group settings

// buildGroupSettingInstance constructs a group setting instance from the configuration.
func buildGroupSettingInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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

	var children []graphmodels.DeviceManagementConfigurationSettingInstanceable

	for _, childConfig := range instanceConfig.GroupSettingValue.Children {
		switch childConfig.ODataType.ValueString() {
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
			// Add a ChoiceSettingInstance child
			if choiceChild := buildChoiceSettingInstance(&childConfig); choiceChild != nil {
				children = append(children, choiceChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
			// Add a SimpleSettingInstance child
			if simpleChild := buildSimpleSettingInstance(&childConfig); simpleChild != nil {
				children = append(children, simpleChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
			// Add a GroupSettingCollectionInstance child
			if groupCollectionChild := buildGroupSettingCollectionInstance(&childConfig); groupCollectionChild != nil {
				children = append(children, groupCollectionChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationGroupSettingInstance":
			// Add a GroupSettingInstance child
			if groupChild := buildGroupSettingInstance(&childConfig); groupChild != nil {
				children = append(children, groupChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance":
			// Add a ChoiceSettingCollectionInstance child
			if choiceCollectionChild := buildChoiceSettingCollectionInstance(&childConfig); choiceCollectionChild != nil {
				children = append(children, choiceCollectionChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
			// Add a SimpleSettingCollectionInstance child
			if simpleCollectionChild := buildSimpleSettingCollectionInstance(&childConfig); simpleCollectionChild != nil {
				children = append(children, simpleCollectionChild)
			}
		}
	}

	value.SetChildren(children)
	instance.SetGroupSettingValue(value)
	return instance
}

// buildGroupSettingCollectionInstance constructs a group setting collection instance from the configuration.
func buildGroupSettingCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {

	if instanceConfig.GroupSettingCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
	odataType := DeviceManagementConfigurationGroupSettingCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	var groupSettingValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable

	// Process each child in GroupSettingCollectionValue's Children
	for _, childConfig := range instanceConfig.GroupSettingCollectionValue.Children {
		groupSettingValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
		groupValueODataType := DeviceManagementConfigurationGroupSettingValue
		groupSettingValue.SetOdataType(&groupValueODataType)

		var nestedChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
		switch childConfig.ODataType.ValueString() {
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
			if choiceChild := buildChoiceSettingInstance(&childConfig); choiceChild != nil {
				nestedChildren = append(nestedChildren, choiceChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
			if simpleChild := buildSimpleSettingInstance(&childConfig); simpleChild != nil {
				nestedChildren = append(nestedChildren, simpleChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationGroupSettingInstance":
			if groupChild := buildGroupSettingInstance(&childConfig); groupChild != nil {
				nestedChildren = append(nestedChildren, groupChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance":
			if choiceCollectionChild := buildChoiceSettingCollectionInstance(&childConfig); choiceCollectionChild != nil {
				nestedChildren = append(nestedChildren, choiceCollectionChild)
			}
		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
			if simpleCollectionChild := buildSimpleSettingCollectionInstance(&childConfig); simpleCollectionChild != nil {
				nestedChildren = append(nestedChildren, simpleCollectionChild)
			}
		}

		groupSettingValue.SetChildren(nestedChildren)
		groupSettingValues = append(groupSettingValues, groupSettingValue)
	}

	instance.SetGroupSettingCollectionValue(groupSettingValues)
	return instance
}

// buildSettingGroupInstance constructs a setting group instance from the configuration.
func buildSettingGroupInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.SettingGroupSettingValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSettingGroupInstance()

	odataType := DeviceManagementConfigurationSettingGroupInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	return instance
}

// buildSettingGroupCollectionInstance constructs a setting group collection instance from the configuration.
func buildSettingGroupCollectionInstance(instanceConfig *DeviceManagementConfigurationSettingInstanceResourceModel) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	if instanceConfig.SettingGroupCollectionValue == nil {
		return nil
	}

	instance := graphmodels.NewDeviceManagementConfigurationSettingGroupCollectionInstance()

	odataType := DeviceManagementConfigurationSettingGroupCollectionInstance
	instance.SetOdataType(&odataType)
	settingDefinitionID := instanceConfig.SettingDefinitionID.ValueString()
	instance.SetSettingDefinitionId(&settingDefinitionID)

	return instance
}
