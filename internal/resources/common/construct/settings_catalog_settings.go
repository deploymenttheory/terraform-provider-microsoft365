package construct

import (
	"context"
	"encoding/json"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructSettingsCatalogSettings constructs a collection of settings catalog settings from a JSON string
// it supports simple, choice, simpleCollection, choiceCollection, and groupCollection settings and nested
// settings within choice and group collections recursively.
func ConstructSettingsCatalogSettings(ctx context.Context, settingsJSON types.String) []graphmodels.DeviceManagementConfigurationSettingable {
	tflog.Debug(ctx, "Constructing settings catalog settings")

	var configModel sharedmodels.DeviceConfigV2GraphServiceModel
	if err := json.Unmarshal([]byte(settingsJSON.ValueString()), &configModel); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	tflog.Debug(ctx, "Unmarshaled settings data", map[string]interface{}{
		"data": configModel,
	})

	settingsCollection := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	for _, detail := range configModel.SettingsDetails {
		baseSetting := graphmodels.NewDeviceManagementConfigurationSetting()
		instance, instanceType := createBaseInstance(ctx, detail.SettingInstance.ODataType, detail.SettingInstance.SettingDefinitionId)

		if instance == nil {
			continue
		}

		switch instanceType {
		case "simple":
			simpleInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
			if detail.SettingInstance.SimpleSettingValue != nil {
				simpleInstance.SetSimpleSettingValue(handleSimpleValue(ctx, detail.SettingInstance.SimpleSettingValue))
			}
			setInstanceTemplateReference(simpleInstance, detail.SettingInstance.SettingInstanceTemplateReference)
			baseSetting.SetSettingInstance(simpleInstance)

		case "choice":
			choiceInstance := instance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
			if detail.SettingInstance.ChoiceSettingValue != nil {
				choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
				choiceValue.SetValue(&detail.SettingInstance.ChoiceSettingValue.Value)
				setValueTemplateReference(choiceValue, detail.SettingInstance.ChoiceSettingValue.SettingValueTemplateReference)

				if len(detail.SettingInstance.ChoiceSettingValue.Children) > 0 {
					choiceChildren := handleChoiceSettingChildren(ctx, detail.SettingInstance.ChoiceSettingValue.Children)
					choiceValue.SetChildren(choiceChildren)
				}

				choiceInstance.SetChoiceSettingValue(choiceValue)
			}
			setInstanceTemplateReference(choiceInstance, detail.SettingInstance.SettingInstanceTemplateReference)
			baseSetting.SetSettingInstance(choiceInstance)

		case "simpleCollection":
			collectionInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
			if len(detail.SettingInstance.SimpleSettingCollectionValue) > 0 {
				values := handleSimpleSettingCollection(detail.SettingInstance.SimpleSettingCollectionValue)
				collectionInstance.SetSimpleSettingCollectionValue(values)
			}
			setInstanceTemplateReference(collectionInstance, detail.SettingInstance.SettingInstanceTemplateReference)
			baseSetting.SetSettingInstance(collectionInstance)

		case "choiceCollection":
			collectionInstance := instance.(graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable)
			if len(detail.SettingInstance.ChoiceSettingCollectionValue) > 0 {
				values := handleChoiceCollectionValue(ctx, detail.SettingInstance.ChoiceSettingCollectionValue)
				collectionInstance.SetChoiceSettingCollectionValue(values)
			}
			setInstanceTemplateReference(collectionInstance, detail.SettingInstance.SettingInstanceTemplateReference)
			baseSetting.SetSettingInstance(collectionInstance)

		case "groupCollection":
			groupInstance := instance.(graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable)
			if len(detail.SettingInstance.GroupSettingCollectionValue) > 0 {
				values := handleGroupSettingCollection(ctx, detail.SettingInstance.GroupSettingCollectionValue)
				groupInstance.SetGroupSettingCollectionValue(values)
			}
			setInstanceTemplateReference(groupInstance, detail.SettingInstance.SettingInstanceTemplateReference)
			baseSetting.SetSettingInstance(groupInstance)
		}

		settingsCollection = append(settingsCollection, baseSetting)
	}

	tflog.Debug(ctx, "Constructed settings catalog settings", map[string]interface{}{
		"count": len(settingsCollection),
	})

	return settingsCollection
}

// createBaseInstance creates and initializes a new setting instance based on the provided OData type and setting definition ID.
// The function dynamically determines the appropriate type of the instance to create and returns it along with its type as a string.
// Supported types include:
// - Simple Setting Instance
// - Choice Setting Instance
// - Simple Setting Collection Instance
// - Choice Setting Collection Instance
// - Group Setting Collection Instance
//
// Parameters:
// - odataType: The OData type of the setting instance (e.g., "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance").
// - settingDefinitionId: The ID of the setting definition associated with this instance.
//
// Returns:
// - interface{}: The newly created setting instance, or nil if the OData type is unsupported.
// - string: A string identifier for the type of the instance (e.g., "simple", "choice").
func createBaseInstance(ctx context.Context, odataType string, settingDefinitionId string) (interface{}, string) {
	if odataType == "" {
		tflog.Error(ctx, "Invalid input: OData type is empty", map[string]interface{}{
			"odataType": odataType,
		})
		return nil, ""
	}

	switch odataType {
	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
		instance.SetOdataType(&odataType)
		instance.SetSettingDefinitionId(&settingDefinitionId)
		return instance, "simple"

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
		instance.SetOdataType(&odataType)
		instance.SetSettingDefinitionId(&settingDefinitionId)
		return instance, "choice"

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
		instance.SetOdataType(&odataType)
		instance.SetSettingDefinitionId(&settingDefinitionId)
		return instance, "simpleCollection"

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
		instance.SetOdataType(&odataType)
		instance.SetSettingDefinitionId(&settingDefinitionId)
		return instance, "choiceCollection"

	case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
		instance.SetOdataType(&odataType)
		instance.SetSettingDefinitionId(&settingDefinitionId)
		return instance, "groupCollection"
	}

	return nil, ""
}

// handleSimpleValue processes a simple setting value and returns the corresponding model instance.
// It supports various types of simple settings such as string, integer, and secret values, and dynamically
// creates the appropriate setting value instance based on the OData type.
//
// Parameters:
// - ctx: The context for logging and other operations.
// - valueStruct: A pointer to the SimpleSettingStruct containing the value to process and its associated metadata.
//
// Supported OData types:
// - "#microsoft.graph.deviceManagementConfigurationStringSettingValue" (string values)
// - "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue" (integer values)
// - "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" (secret values with value state)
//
// Returns:
//   - graphmodels.DeviceManagementConfigurationSimpleSettingValueable: The processed setting value instance,
//     or nil if the OData type is unsupported or an error occurs.
//
// Logs:
//   - Logs an error if the function is unable to process the provided valueStruct due to an unsupported type
//     or invalid value.
func handleSimpleValue(ctx context.Context, valueStruct *sharedmodels.SimpleSettingStruct) graphmodels.DeviceManagementConfigurationSimpleSettingValueable {
	if valueStruct == nil {
		return nil
	}

	var result graphmodels.DeviceManagementConfigurationSimpleSettingValueable

	switch valueStruct.ODataType {
	case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
		if strValue, ok := valueStruct.Value.(string); ok {
			stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
			stringValue.SetOdataType(&valueStruct.ODataType)
			stringValue.SetValue(&strValue)
			result = stringValue
		}

	case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
		if numValue, ok := valueStruct.Value.(float64); ok {
			intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
			intValue.SetOdataType(&valueStruct.ODataType)
			int32Value := int32(numValue)
			intValue.SetValue(&int32Value)
			result = intValue
		}

	case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue":
		if strValue, ok := valueStruct.Value.(string); ok {
			secretValue := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
			secretValue.SetOdataType(&valueStruct.ODataType)
			secretValue.SetValue(&strValue)
			if valueStruct.ValueState != "" {
				if state, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(valueStruct.ValueState); err == nil {
					secretValue.SetValueState(state.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
				}
			}
			result = secretValue
		}
	}

	if result != nil {
		setValueTemplateReference(result, valueStruct.SettingValueTemplateReference)
		return result
	}

	tflog.Error(ctx, "Failed to handle simple setting value", map[string]interface{}{
		"type":  valueStruct.ODataType,
		"value": valueStruct.Value,
	})
	return nil
}

// Helper function to handle simple setting collections
func handleSimpleSettingCollection(collectionValues []sharedmodels.SimpleSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationSimpleSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
	for _, v := range collectionValues {
		stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
		stringValue.SetOdataType(&v.ODataType)
		stringValue.SetValue(&v.Value)
		setValueTemplateReference(stringValue, v.SettingValueTemplateReference)
		values = append(values, stringValue)
	}
	return values
}

// handleChoiceSettingChildren recursively processes a list of choice setting children.
// Each child is converted into a setting instance based on its OData type and configuration,
// supporting nested recursive structures.
//
// Parameters:
// - ctx: The context for logging and operations.
// - children: A slice of sharedmodels.ChoiceSettingChild representing the choice setting children to process.
//
// Returns:
// - []graphmodels.DeviceManagementConfigurationSettingInstanceable: A slice of processed setting instances.
func handleChoiceSettingChildren(ctx context.Context, children []sharedmodels.ChoiceSettingChild) []graphmodels.DeviceManagementConfigurationSettingInstanceable {
	var result []graphmodels.DeviceManagementConfigurationSettingInstanceable

	for _, child := range children {
		instance := handleSettingInstance(ctx, sharedmodels.SettingInstance{
			ODataType:                        child.ODataType,
			SettingDefinitionId:              child.SettingDefinitionId,
			SettingInstanceTemplateReference: child.SettingInstanceTemplateReference,
			SimpleSettingValue:               child.SimpleSettingValue,
			SimpleSettingCollectionValue:     child.SimpleSettingCollectionValue,
			ChoiceSettingValue:               child.ChoiceSettingValue,
			ChoiceSettingCollectionValue:     child.ChoiceSettingCollectionValue,
			GroupSettingCollectionValue:      child.GroupSettingCollectionValue,
		})
		if instance != nil {
			result = append(result, instance)
		}
	}

	return result
}

// handleGroupSettingCollection recursively processes a list of group setting collections.
// Each group item and its children are converted into setting instances based on their OData type
// and configuration, supporting deeply nested group structures.
//
// Parameters:
// - ctx: The context for logging and operations.
// - groupValues: A slice of sharedmodels.GroupSettingCollectionStruct representing the group setting collections to process.
//
// Returns:
// - []graphmodels.DeviceManagementConfigurationGroupSettingValueable: A slice of processed group setting values.
func handleGroupSettingCollection(ctx context.Context, groupValues []sharedmodels.GroupSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationGroupSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationGroupSettingValueable

	for _, groupItem := range groupValues {
		groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
		var children []graphmodels.DeviceManagementConfigurationSettingInstanceable

		for _, child := range groupItem.Children {
			instance := handleSettingInstance(ctx, sharedmodels.SettingInstance{
				ODataType:                        child.ODataType,
				SettingDefinitionId:              child.SettingDefinitionId,
				SettingInstanceTemplateReference: child.SettingInstanceTemplateReference,
				SimpleSettingValue:               child.SimpleSettingValue,
				SimpleSettingCollectionValue:     child.SimpleSettingCollectionValue,
				ChoiceSettingValue:               child.ChoiceSettingValue,
				ChoiceSettingCollectionValue:     child.ChoiceSettingCollectionValue,
				GroupSettingCollectionValue:      child.GroupSettingCollectionValue,
			})
			if instance != nil {
				children = append(children, instance)
			}
		}

		groupValue.SetChildren(children)
		setValueTemplateReference(groupValue, groupItem.SettingValueTemplateReference)
		values = append(values, groupValue)
	}

	return values
}

// Helper function to handle choice collection values recursively
func handleChoiceCollectionValue(ctx context.Context, collectionValues []sharedmodels.ChoiceSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationChoiceSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationChoiceSettingValueable

	for _, choiceItem := range collectionValues {
		choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
		choiceValue.SetValue(&choiceItem.Value)
		setValueTemplateReference(choiceValue, choiceItem.SettingValueTemplateReference)

		var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
		for _, child := range choiceItem.Children {
			instance, instanceType := createBaseInstance(ctx, child.ODataType, child.SettingDefinitionId)
			if instance == nil {
				continue
			}

			switch instanceType {
			case "simple":
				simpleInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
				if child.SimpleSettingValue != nil {
					simpleInstance.SetSimpleSettingValue(handleSimpleValue(ctx, child.SimpleSettingValue))
				}
				setInstanceTemplateReference(simpleInstance, child.SettingInstanceTemplateReference)
				children = append(children, simpleInstance)

			case "simpleCollection":
				collectionInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
				if len(child.SimpleSettingCollectionValue) > 0 {
					simpleValues := handleSimpleSettingCollection(child.SimpleSettingCollectionValue)
					collectionInstance.SetSimpleSettingCollectionValue(simpleValues)
				}
				setInstanceTemplateReference(collectionInstance, child.SettingInstanceTemplateReference)
				children = append(children, collectionInstance)
			}
		}

		choiceValue.SetChildren(children)
		values = append(values, choiceValue)
	}

	return values
}

// handleSettingInstance processes a given setting instance by creating and configuring
// the appropriate type based on its OData type and setting definition ID. This function
// supports recursive handling of nested settings, including simple, choice, collection,
// and group settings.
//
// Parameters:
// - ctx: The context for logging and operations.
// - instance: A sharedmodels.SettingInstance object representing the setting to process.
//
// Returns:
// - graphmodels.DeviceManagementConfigurationSettingInstanceable: The configured setting instance, or nil if the input is invalid or unsupported.
func handleSettingInstance(ctx context.Context, instance sharedmodels.SettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
	baseInstance, instanceType := createBaseInstance(ctx, instance.ODataType, instance.SettingDefinitionId)
	if baseInstance == nil {
		return nil
	}

	switch instanceType {
	case "simple":
		simpleInstance := baseInstance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
		if instance.SimpleSettingValue != nil {
			simpleInstance.SetSimpleSettingValue(handleSimpleValue(ctx, instance.SimpleSettingValue))
		}
		setInstanceTemplateReference(simpleInstance, instance.SettingInstanceTemplateReference)
		return simpleInstance

	case "choice":
		choiceInstance := baseInstance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
		if instance.ChoiceSettingValue != nil {
			choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			choiceValue.SetValue(&instance.ChoiceSettingValue.Value)
			setValueTemplateReference(choiceValue, instance.ChoiceSettingValue.SettingValueTemplateReference)

			if len(instance.ChoiceSettingValue.Children) > 0 {
				children := handleChoiceSettingChildren(ctx, instance.ChoiceSettingValue.Children)
				choiceValue.SetChildren(children)
			}
			choiceInstance.SetChoiceSettingValue(choiceValue)
		}
		setInstanceTemplateReference(choiceInstance, instance.SettingInstanceTemplateReference)
		return choiceInstance

	case "simpleCollection":
		collectionInstance := baseInstance.(graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
		if len(instance.SimpleSettingCollectionValue) > 0 {
			values := handleSimpleSettingCollection(instance.SimpleSettingCollectionValue)
			collectionInstance.SetSimpleSettingCollectionValue(values)
		}
		setInstanceTemplateReference(collectionInstance, instance.SettingInstanceTemplateReference)
		return collectionInstance

	case "choiceCollection":
		collectionInstance := baseInstance.(graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable)
		if len(instance.ChoiceSettingCollectionValue) > 0 {
			values := handleChoiceCollectionValue(ctx, instance.ChoiceSettingCollectionValue)
			collectionInstance.SetChoiceSettingCollectionValue(values)
		}
		setInstanceTemplateReference(collectionInstance, instance.SettingInstanceTemplateReference)
		return collectionInstance

	case "groupCollection":
		groupInstance := baseInstance.(graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable)
		if len(instance.GroupSettingCollectionValue) > 0 {
			values := handleGroupSettingCollection(ctx, instance.GroupSettingCollectionValue)
			groupInstance.SetGroupSettingCollectionValue(values)
		}
		setInstanceTemplateReference(groupInstance, instance.SettingInstanceTemplateReference)
		return groupInstance
	}

	return nil
}

// setInstanceTemplateReference creates and assigns a SettingInstanceTemplateReference to a setting instance.
// This function adds the template reference metadata to a given setting instance if the reference is provided.
//
// Parameters:
// - instance: The setting instance implementing DeviceManagementConfigurationSettingInstanceable.
// - ref: The SettingInstanceTemplateReference containing the template ID.
//
// Example:
// Input:
//
//	{
//	    "settingInstanceTemplateReference": {
//	        "settingInstanceTemplateId": "template-id-123"
//	    }
//	}
func setInstanceTemplateReference(instance graphmodels.DeviceManagementConfigurationSettingInstanceable, ref *sharedmodels.SettingInstanceTemplateReference) {
	if ref != nil {
		templateRef := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()
		templateRef.SetSettingInstanceTemplateId(&ref.SettingInstanceTemplateId)
		instance.SetSettingInstanceTemplateReference(templateRef)
	}
}

// setValueTemplateReference creates and assigns a SettingValueTemplateReference to a value object.
// This function adds the template reference metadata to a value object if the reference is provided.
// It supports group, choice, and simple setting value types.
//
// Parameters:
// - value: The value object implementing one of the supported interfaces (Group, Choice, or Simple setting values).
// - ref: The SettingValueTemplateReference containing the template ID and `useTemplateDefault` flag.
//
// Example:
// Input:
//
//	{
//	    "value": "example-value",
//	    "settingValueTemplateReference": {
//	        "settingValueTemplateId": "template-id-456",
//	        "useTemplateDefault": true
//	    }
//	}
//
// Effect:
// The `value` will have its SettingValueTemplateReference field set to:
//
//	{
//	    "settingValueTemplateId": "template-id-456",
//	    "useTemplateDefault": true
//	}
//
// Supported Value Types:
// - GroupSettingValue
// - ChoiceSettingValue
// - SimpleSettingValue
func setValueTemplateReference(value interface{}, ref *sharedmodels.SettingValueTemplateReference) {
	if ref != nil {
		templateRef := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
		templateRef.SetSettingValueTemplateId(&ref.SettingValueTemplateId)
		templateRef.SetUseTemplateDefault(&ref.UseTemplateDefault)

		switch v := value.(type) {
		case graphmodels.DeviceManagementConfigurationGroupSettingValueable:
			v.SetSettingValueTemplateReference(templateRef)
		case graphmodels.DeviceManagementConfigurationChoiceSettingValueable:
			v.SetSettingValueTemplateReference(templateRef)
		case graphmodels.DeviceManagementConfigurationSimpleSettingValueable:
			v.SetSettingValueTemplateReference(templateRef)
		}
	}
}
