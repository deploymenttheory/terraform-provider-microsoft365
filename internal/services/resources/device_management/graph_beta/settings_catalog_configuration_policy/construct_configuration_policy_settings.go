package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructSettingsCatalogSettings constructs a collection of settings catalog settings from HCL struct
// it supports simple, choice, simpleCollection, choiceCollection, and groupCollection settings and nested
// settings within choice and group collections recursively.
func ConstructSettingsCatalogSettings(ctx context.Context, configModel DeviceConfigV2GraphServiceResourceModel) []graphmodels.DeviceManagementConfigurationSettingable {
	tflog.Debug(ctx, "Constructing settings catalog settings")

	tflog.Debug(ctx, "Processing settings catalog data from HCL", map[string]interface{}{
		"settings_count": len(configModel.Settings),
	})

	constructedSettings := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	// Process all settings in the array
	for _, setting := range configModel.Settings {
		processSetting(ctx, setting, &constructedSettings)
	}

	tflog.Debug(ctx, "Constructed settings catalog settings", map[string]interface{}{
		"count": len(constructedSettings),
	})

	return constructedSettings
}

// processSetting processes an individual settings catalog setting and appends it to the settingsRequestPayload.
//
// This function is part of the constructor logic and is responsible for creating and configuring
// instances of DeviceManagementConfigurationSetting based on the provided `setting`. It dynamically
// determines the type of setting based on the `@odata.type` field and handles different setting types, including:
// - Simple Settings
// - Choice Settings
// - Simple Collection Settings
// - Choice Collection Settings
// - Group Collection Settings
//
// Why Use a Separate Function for `processSetting`?
// -------------------------------------------------
// The Microsoft Graph API returns settings in a consistent format but allows for various nested and typed
// configurations. This function abstracts the logic of initializing, configuring, and appending each setting
// instance to the request payload, ensuring modularity and readability.
//
// Why Handle Both `settings` (Array) and `setting` (Single Instance)?
// -------------------------------------------------------------------
// The API always returns settings as an array, but it supports both single and multiple settings. As a result, the constructor
// allows for the definiton a single `setting` and `settings` array.
//
// Parameters:
//   - ctx: Context for logging and cancellation.
//   - setting: A `Setting` object representing an individual setting to process.
//   - settingsRequestPayload: A pointer to a slice of `DeviceManagementConfigurationSettingable` where the processed
//     setting will be appended.
func processSetting(ctx context.Context, setting Setting, settingsRequestPayload *[]graphmodels.DeviceManagementConfigurationSettingable) {
	baseSetting := graphmodels.NewDeviceManagementConfigurationSetting()
	instance, instanceType := createBaseInstance(ctx, setting.SettingInstance.ODataType, setting.SettingInstance.SettingDefinitionId)

	if instance == nil {
		return
	}

	switch instanceType {
	case "simple":
		simpleInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
		if setting.SettingInstance.SimpleSettingValue != nil {
			simpleInstance.SetSimpleSettingValue(handleSimpleValue(ctx, setting.SettingInstance.SimpleSettingValue))
		}
		setInstanceTemplateReference(simpleInstance, setting.SettingInstance.SettingInstanceTemplateReference)
		baseSetting.SetSettingInstance(simpleInstance)

	case "choice":
		choiceInstance := instance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
		if setting.SettingInstance.ChoiceSettingValue != nil {
			choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
			convert.FrameworkToGraphString(setting.SettingInstance.ChoiceSettingValue.Value, choiceValue.SetValue)
			setValueTemplateReference(choiceValue, setting.SettingInstance.ChoiceSettingValue.SettingValueTemplateReference)

			if len(setting.SettingInstance.ChoiceSettingValue.Children) > 0 {
				choiceChildren := handleChoiceSettingChildren(ctx, setting.SettingInstance.ChoiceSettingValue.Children)
				choiceValue.SetChildren(choiceChildren)
			}

			choiceInstance.SetChoiceSettingValue(choiceValue)
		}
		setInstanceTemplateReference(choiceInstance, setting.SettingInstance.SettingInstanceTemplateReference)
		baseSetting.SetSettingInstance(choiceInstance)

	case "simpleCollection":
		collectionInstance := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
		if len(setting.SettingInstance.SimpleSettingCollectionValue) > 0 {
			values := handleSimpleSettingCollection(setting.SettingInstance.SimpleSettingCollectionValue)
			collectionInstance.SetSimpleSettingCollectionValue(values)
		}
		setInstanceTemplateReference(collectionInstance, setting.SettingInstance.SettingInstanceTemplateReference)
		baseSetting.SetSettingInstance(collectionInstance)

	case "choiceCollection":
		collectionInstance := instance.(graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable)
		if len(setting.SettingInstance.ChoiceSettingCollectionValue) > 0 {
			values := handleChoiceCollectionValue(ctx, setting.SettingInstance.ChoiceSettingCollectionValue)
			collectionInstance.SetChoiceSettingCollectionValue(values)
		}
		setInstanceTemplateReference(collectionInstance, setting.SettingInstance.SettingInstanceTemplateReference)
		baseSetting.SetSettingInstance(collectionInstance)

	case "groupCollection":
		groupInstance := instance.(graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable)
		if len(setting.SettingInstance.GroupSettingCollectionValue) > 0 {
			values := handleGroupSettingCollection(ctx, setting.SettingInstance.GroupSettingCollectionValue)
			groupInstance.SetGroupSettingCollectionValue(values)
		}
		setInstanceTemplateReference(groupInstance, setting.SettingInstance.SettingInstanceTemplateReference)
		baseSetting.SetSettingInstance(groupInstance)
	}

	*settingsRequestPayload = append(*settingsRequestPayload, baseSetting)
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
func createBaseInstance(ctx context.Context, odataType types.String, settingDefinitionId types.String) (interface{}, string) {
	// Check if odataType is null, unknown, or empty
	if odataType.IsNull() || odataType.IsUnknown() || odataType.ValueString() == "" {
		tflog.Error(ctx, "Invalid input: OData type is empty", map[string]interface{}{
			"odataType": odataType.ValueString(),
		})
		return nil, ""
	}

	// Get the actual string value for comparison
	odataTypeStr := odataType.ValueString()

	switch odataTypeStr {
	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
		// Use constructor helpers to handle types.String -> *string conversion
		convert.FrameworkToGraphString(odataType, instance.SetOdataType)
		convert.FrameworkToGraphString(settingDefinitionId, instance.SetSettingDefinitionId)
		return instance, "simple"

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
		convert.FrameworkToGraphString(odataType, instance.SetOdataType)
		convert.FrameworkToGraphString(settingDefinitionId, instance.SetSettingDefinitionId)
		return instance, "choice"

	case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
		convert.FrameworkToGraphString(odataType, instance.SetOdataType)
		convert.FrameworkToGraphString(settingDefinitionId, instance.SetSettingDefinitionId)
		return instance, "simpleCollection"

	case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
		convert.FrameworkToGraphString(odataType, instance.SetOdataType)
		convert.FrameworkToGraphString(settingDefinitionId, instance.SetSettingDefinitionId)
		return instance, "choiceCollection"

	case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
		instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
		convert.FrameworkToGraphString(odataType, instance.SetOdataType)
		convert.FrameworkToGraphString(settingDefinitionId, instance.SetSettingDefinitionId)
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
func handleSimpleValue(ctx context.Context, valueStruct *SimpleSettingStruct) graphmodels.DeviceManagementConfigurationSimpleSettingValueable {
	if valueStruct == nil {
		return nil
	}

	var result graphmodels.DeviceManagementConfigurationSimpleSettingValueable

	odataTypeStr := valueStruct.ODataType.ValueString()

	switch odataTypeStr {
	case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
		stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
		convert.FrameworkToGraphString(valueStruct.ODataType, stringValue.SetOdataType)
		convert.FrameworkToGraphString(valueStruct.Value, stringValue.SetValue)
		result = stringValue

	case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
		if intVal, err := strconv.Atoi(valueStruct.Value.ValueString()); err == nil {
			intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
			convert.FrameworkToGraphString(valueStruct.ODataType, intValue.SetOdataType)
			int32Value := int32(intVal)
			intValue.SetValue(&int32Value)
			result = intValue
		} else {
			tflog.Error(ctx, "Failed to convert string to integer", map[string]interface{}{
				"value": valueStruct.Value.ValueString(),
				"error": err.Error(),
			})
			return nil
		}

	case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue":
		secretValue := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
		convert.FrameworkToGraphString(valueStruct.ODataType, secretValue.SetOdataType)
		convert.FrameworkToGraphString(valueStruct.Value, secretValue.SetValue)

		// Handle ValueState if it's not null/unknown/empty
		if !valueStruct.ValueState.IsNull() && !valueStruct.ValueState.IsUnknown() && valueStruct.ValueState.ValueString() != "" {
			if state, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(valueStruct.ValueState.ValueString()); err == nil {
				secretValue.SetValueState(state.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
			}
		}
		result = secretValue
	}

	if result != nil {
		setValueTemplateReference(result, valueStruct.SettingValueTemplateReference)
		return result
	}

	tflog.Error(ctx, "Failed to handle simple setting value", map[string]interface{}{
		"type":  valueStruct.ODataType.ValueString(),
		"value": valueStruct.Value.ValueString(),
	})
	return nil
}

// handleSimpleSettingCollection is a helper function to handle simple setting collections
func handleSimpleSettingCollection(collectionValues []SimpleSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationSimpleSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
	for _, v := range collectionValues {
		stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
		convert.FrameworkToGraphString(v.ODataType, stringValue.SetOdataType)
		convert.FrameworkToGraphString(v.Value, stringValue.SetValue)
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
// - children: A slice of ChoiceSettingChild representing the choice setting children to process.
//
// Returns:
// - []graphmodels.DeviceManagementConfigurationSettingInstanceable: A slice of processed setting instances.
func handleChoiceSettingChildren(ctx context.Context, children []ChoiceSettingChild) []graphmodels.DeviceManagementConfigurationSettingInstanceable {
	var result []graphmodels.DeviceManagementConfigurationSettingInstanceable

	for _, child := range children {
		instance := handleSettingInstance(ctx, SettingInstance{
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
// - groupValues: A slice of GroupSettingCollectionStruct representing the group setting collections to process.
//
// Returns:
// - []graphmodels.DeviceManagementConfigurationGroupSettingValueable: A slice of processed group setting values.
func handleGroupSettingCollection(ctx context.Context, groupValues []GroupSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationGroupSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationGroupSettingValueable

	for _, groupItem := range groupValues {
		groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
		var children []graphmodels.DeviceManagementConfigurationSettingInstanceable

		for _, child := range groupItem.Children {
			instance := handleSettingInstance(ctx, SettingInstance{
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
func handleChoiceCollectionValue(ctx context.Context, collectionValues []ChoiceSettingCollectionStruct) []graphmodels.DeviceManagementConfigurationChoiceSettingValueable {
	var values []graphmodels.DeviceManagementConfigurationChoiceSettingValueable

	for _, choiceItem := range collectionValues {
		choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
		convert.FrameworkToGraphString(choiceItem.Value, choiceValue.SetValue)
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
// - instance: A SettingInstance object representing the setting to process.
//
// Returns:
// - graphmodels.DeviceManagementConfigurationSettingInstanceable: The configured setting instance, or nil if the input is invalid or unsupported.
func handleSettingInstance(ctx context.Context, instance SettingInstance) graphmodels.DeviceManagementConfigurationSettingInstanceable {
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
			convert.FrameworkToGraphString(instance.ChoiceSettingValue.Value, choiceValue.SetValue)
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
func setInstanceTemplateReference(instance graphmodels.DeviceManagementConfigurationSettingInstanceable, ref *SettingInstanceTemplateReference) {
	if ref != nil {
		templateRef := graphmodels.NewDeviceManagementConfigurationSettingInstanceTemplateReference()

		convert.FrameworkToGraphString(ref.SettingInstanceTemplateId, templateRef.SetSettingInstanceTemplateId)

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
func setValueTemplateReference(value interface{}, ref *SettingValueTemplateReference) {
	if ref == nil {
		return
	}

	templateRef := graphmodels.NewDeviceManagementConfigurationSettingValueTemplateReference()
	convert.FrameworkToGraphString(ref.SettingValueTemplateId, templateRef.SetSettingValueTemplateId)
	templateRef.SetUseTemplateDefault(&ref.UseTemplateDefault)

	// Use type assertion instead of type switch to avoid interface hierarchy issues
	if v, ok := value.(graphmodels.DeviceManagementConfigurationSimpleSettingValueable); ok {
		v.SetSettingValueTemplateReference(templateRef)
	} else if v, ok := value.(graphmodels.DeviceManagementConfigurationChoiceSettingValueable); ok {
		v.SetSettingValueTemplateReference(templateRef)
	} else if v, ok := value.(graphmodels.DeviceManagementConfigurationGroupSettingValueable); ok {
		v.SetSettingValueTemplateReference(templateRef)
	}
}
