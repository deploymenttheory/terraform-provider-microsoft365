package graphBetaSettingsCatalog

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateConfigurationPolicySettings maps settings from Graph  models to Terraform state
func StateConfigurationPolicySettings(ctx context.Context, data *SettingsCatalogProfileResourceModel, settingsResponse graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) error {
	tflog.Debug(ctx, "Starting to map settings from Graph  models to Terraform state")

	if settingsResponse == nil {
		tflog.Debug(ctx, "No settings response data to process")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "Settings array is empty")
		return nil
	}

	// Convert settings to our model
	deviceConfigModel := &DeviceConfigV2GraphServiceResourceModel{}
	var mappedSettings []Setting

	for _, apiSetting := range settings {
		setting, err := mapSettingToModel(ctx, apiSetting)
		if err != nil {
			tflog.Warn(ctx, "Failed to map setting", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		if setting != nil {
			mappedSettings = append(mappedSettings, *setting)
		}
	}

	deviceConfigModel.Settings = mappedSettings
	data.ConfigurationPolicy = deviceConfigModel

	tflog.Debug(ctx, fmt.Sprintf("Finished stating settings catalog resource %s with id %s", ResourceName, data.ID.ValueString()))

	return nil
}

// mapSettingToModel converts a Graph  setting to our Terraform model
func mapSettingToModel(ctx context.Context, apiSetting graphmodels.DeviceManagementConfigurationSettingable) (*Setting, error) {
	if apiSetting == nil {
		return nil, fmt.Errorf("API setting is nil")
	}

	setting := &Setting{}

	// Map setting ID (if available)
	if id := apiSetting.GetId(); id != nil {
		setting.ID = types.StringValue(*id)
	}

	// Map the setting instance
	settingInstance := apiSetting.GetSettingInstance()
	if settingInstance == nil {
		return nil, fmt.Errorf("setting instance is nil")
	}

	mappedInstance, err := mapSettingInstanceToModel(ctx, settingInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to map setting instance: %w", err)
	}

	setting.SettingInstance = *mappedInstance

	return setting, nil
}

// mapSettingInstanceToModel converts a Graph  setting instance to our model
func mapSettingInstanceToModel(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSettingInstanceable) (*SettingInstance, error) {
	if instance == nil {
		return nil, fmt.Errorf("setting instance is nil")
	}

	settingInstance := &SettingInstance{}

	// Map OData type
	if odataType := instance.GetOdataType(); odataType != nil {
		settingInstance.ODataType = types.StringValue(*odataType)
	}

	// Map setting definition ID
	if settingDefId := instance.GetSettingDefinitionId(); settingDefId != nil {
		settingInstance.SettingDefinitionId = types.StringValue(*settingDefId)
	}

	// Map instance template reference
	if instanceTemplateRef := instance.GetSettingInstanceTemplateReference(); instanceTemplateRef != nil {
		settingInstance.SettingInstanceTemplateReference = mapInstanceTemplateReference(instanceTemplateRef)
	}

	// Type-specific mapping based on the concrete type
	switch typedInstance := instance.(type) {
	case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
		if simpleValue := typedInstance.GetSimpleSettingValue(); simpleValue != nil {
			mappedSimpleValue, err := mapSimpleSettingValue(ctx, simpleValue)
			if err != nil {
				return nil, fmt.Errorf("failed to map simple setting value: %w", err)
			}
			settingInstance.SimpleSettingValue = mappedSimpleValue
		}

	case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
		if choiceValue := typedInstance.GetChoiceSettingValue(); choiceValue != nil {
			mappedChoiceValue, err := mapChoiceSettingValue(ctx, choiceValue)
			if err != nil {
				return nil, fmt.Errorf("failed to map choice setting value: %w", err)
			}
			settingInstance.ChoiceSettingValue = mappedChoiceValue
		}

	case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
		simpleCollectionValues := typedInstance.GetSimpleSettingCollectionValue()
		if len(simpleCollectionValues) > 0 {
			mappedCollection, err := mapSimpleSettingCollection(ctx, simpleCollectionValues)
			if err != nil {
				return nil, fmt.Errorf("failed to map simple setting collection: %w", err)
			}
			settingInstance.SimpleSettingCollectionValue = mappedCollection
		} else {
			// FIXED: Always initialize as empty slice
			settingInstance.SimpleSettingCollectionValue = make([]SimpleSettingCollectionStruct, 0)
		}

	case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
		choiceCollectionValues := typedInstance.GetChoiceSettingCollectionValue()
		if len(choiceCollectionValues) > 0 {
			mappedCollection, err := mapChoiceSettingCollection(ctx, choiceCollectionValues)
			if err != nil {
				return nil, fmt.Errorf("failed to map choice setting collection: %w", err)
			}
			settingInstance.ChoiceSettingCollectionValue = mappedCollection
		} else {
			// FIXED: Always initialize as empty slice
			settingInstance.ChoiceSettingCollectionValue = make([]ChoiceSettingCollectionStruct, 0)
		}

	case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
		groupCollectionValues := typedInstance.GetGroupSettingCollectionValue()
		if len(groupCollectionValues) > 0 {
			mappedCollection, err := mapGroupSettingCollection(ctx, groupCollectionValues)
			if err != nil {
				return nil, fmt.Errorf("failed to map group setting collection: %w", err)
			}
			settingInstance.GroupSettingCollectionValue = mappedCollection
		} else {
			// FIXED: Always initialize as empty slice
			settingInstance.GroupSettingCollectionValue = make([]GroupSettingCollectionStruct, 0)
		}
	}

	return settingInstance, nil
}

// mapSimpleSettingValue converts  simple setting value to our model
func mapSimpleSettingValue(ctx context.Context, value graphmodels.DeviceManagementConfigurationSimpleSettingValueable) (*SimpleSettingStruct, error) {
	if value == nil {
		return nil, fmt.Errorf("simple setting value is nil")
	}

	simpleValue := &SimpleSettingStruct{}

	if odataType := value.GetOdataType(); odataType != nil {
		simpleValue.ODataType = types.StringValue(*odataType)
	}

	if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
		simpleValue.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
	}

	switch typedValue := value.(type) {

	case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
		if intVal := typedValue.GetValue(); intVal != nil {
			simpleValue.Value = types.StringValue(strconv.Itoa(int(*intVal)))
		}

	case graphmodels.DeviceManagementConfigurationSecretSettingValueable:
		if secretVal := typedValue.GetValue(); secretVal != nil {
			simpleValue.Value = types.StringValue(*secretVal)
		}

		if valueState := typedValue.GetValueState(); valueState != nil {
			simpleValue.ValueState = types.StringValue(valueState.String())
		}

	case graphmodels.DeviceManagementConfigurationStringSettingValueable:
		if stringVal := typedValue.GetValue(); stringVal != nil {
			simpleValue.Value = types.StringValue(*stringVal)
		}

	default:
		return nil, fmt.Errorf("unsupported simple setting value type: %T", typedValue)
	}

	return simpleValue, nil
}

// mapChoiceSettingValue converts  choice setting value to our model
func mapChoiceSettingValue(ctx context.Context, value graphmodels.DeviceManagementConfigurationChoiceSettingValueable) (*ChoiceSettingStruct, error) {
	if value == nil {
		return nil, fmt.Errorf("choice setting value is nil")
	}

	choiceValue := &ChoiceSettingStruct{}

	// Map value
	if val := value.GetValue(); val != nil {
		choiceValue.Value = types.StringValue(*val)
	}

	// Map value template reference
	if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
		choiceValue.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
	}

	// FIXED: Always initialize children, even if empty
	children := value.GetChildren()
	if len(children) > 0 {
		mappedChildren, err := mapChoiceSettingChildren(ctx, children)
		if err != nil {
			return nil, fmt.Errorf("failed to map choice setting children: %w", err)
		}
		choiceValue.Children = mappedChildren
	} else {
		// Always initialize as empty slice rather than leaving as nil
		choiceValue.Children = make([]ChoiceSettingChild, 0)
	}

	return choiceValue, nil
}

// mapSimpleSettingCollection converts  simple setting collection to our model
func mapSimpleSettingCollection(ctx context.Context, values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable) ([]SimpleSettingCollectionStruct, error) {
	var result []SimpleSettingCollectionStruct

	for _, value := range values {
		collectionItem := SimpleSettingCollectionStruct{}

		// Map OData type
		if odataType := value.GetOdataType(); odataType != nil {
			collectionItem.ODataType = types.StringValue(*odataType)
		}

		// Map value template reference
		if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
			collectionItem.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
		}

		// Map value (assuming string type for collection items)
		switch typedValue := value.(type) {
		case graphmodels.DeviceManagementConfigurationStringSettingValueable:
			if stringVal := typedValue.GetValue(); stringVal != nil {
				collectionItem.Value = types.StringValue(*stringVal)
			}
		case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
			if intVal := typedValue.GetValue(); intVal != nil {
				collectionItem.Value = types.StringValue(strconv.Itoa(int(*intVal)))
			}
		}

		result = append(result, collectionItem)
	}

	return result, nil
}

// mapChoiceSettingCollection converts  choice setting collection to our model
func mapChoiceSettingCollection(ctx context.Context, values []graphmodels.DeviceManagementConfigurationChoiceSettingValueable) ([]ChoiceSettingCollectionStruct, error) {
	var result []ChoiceSettingCollectionStruct

	for _, value := range values {
		collectionItem := ChoiceSettingCollectionStruct{}

		// Map value
		if val := value.GetValue(); val != nil {
			collectionItem.Value = types.StringValue(*val)
		}

		// Map value template reference
		if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
			collectionItem.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
		}

		// FIXED: Always initialize children, even if empty
		children := value.GetChildren()
		if len(children) > 0 {
			mappedChildren, err := mapChoiceSettingCollectionChildren(ctx, children)
			if err != nil {
				return nil, fmt.Errorf("failed to map choice setting collection children: %w", err)
			}
			collectionItem.Children = mappedChildren
		} else {
			// Always initialize as empty slice rather than leaving as nil
			collectionItem.Children = make([]ChoiceSettingCollectionChild, 0)
		}

		result = append(result, collectionItem)
	}

	return result, nil
}

// mapGroupSettingCollection converts  group setting collection to our model
func mapGroupSettingCollection(ctx context.Context, values []graphmodels.DeviceManagementConfigurationGroupSettingValueable) ([]GroupSettingCollectionStruct, error) {
	var result []GroupSettingCollectionStruct

	for _, value := range values {
		groupItem := GroupSettingCollectionStruct{}

		// Map value template reference
		if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
			groupItem.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
		}

		// FIXED: Always initialize children, even if empty
		children := value.GetChildren()
		if len(children) > 0 {
			mappedChildren, err := mapGroupSettingCollectionChildren(ctx, children)
			if err != nil {
				return nil, fmt.Errorf("failed to map group setting collection children: %w", err)
			}
			groupItem.Children = mappedChildren
		} else {
			// Always initialize as empty slice rather than leaving as nil
			groupItem.Children = make([]GroupSettingCollectionChild, 0)
		}

		result = append(result, groupItem)
	}

	return result, nil
}

// mapChoiceSettingChildren converts  choice setting children to our model
func mapChoiceSettingChildren(ctx context.Context, children []graphmodels.DeviceManagementConfigurationSettingInstanceable) ([]ChoiceSettingChild, error) {
	var result []ChoiceSettingChild

	for _, child := range children {
		childItem := ChoiceSettingChild{}

		// Map basic properties
		if odataType := child.GetOdataType(); odataType != nil {
			childItem.ODataType = types.StringValue(*odataType)
		}
		if settingDefId := child.GetSettingDefinitionId(); settingDefId != nil {
			childItem.SettingDefinitionId = types.StringValue(*settingDefId)
		}

		// Map instance template reference
		if instanceTemplateRef := child.GetSettingInstanceTemplateReference(); instanceTemplateRef != nil {
			childItem.SettingInstanceTemplateReference = mapInstanceTemplateReference(instanceTemplateRef)
		}

		// Type-specific mapping
		switch typedChild := child.(type) {
		case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
			if simpleValue := typedChild.GetSimpleSettingValue(); simpleValue != nil {
				mappedSimpleValue, err := mapSimpleSettingValue(ctx, simpleValue)
				if err != nil {
					return nil, fmt.Errorf("failed to map child simple setting value: %w", err)
				}
				childItem.SimpleSettingValue = mappedSimpleValue
			}

		case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
			simpleCollectionValues := typedChild.GetSimpleSettingCollectionValue()
			if len(simpleCollectionValues) > 0 {
				mappedCollection, err := mapSimpleSettingCollection(ctx, simpleCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map child simple setting collection: %w", err)
				}
				childItem.SimpleSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.SimpleSettingCollectionValue = make([]SimpleSettingCollectionStruct, 0)
			}

		case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
			if choiceValue := typedChild.GetChoiceSettingValue(); choiceValue != nil {
				mappedChoiceValue, err := mapChoiceSettingValue(ctx, choiceValue)
				if err != nil {
					return nil, fmt.Errorf("failed to map child choice setting value: %w", err)
				}
				childItem.ChoiceSettingValue = mappedChoiceValue
			}

		case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
			choiceCollectionValues := typedChild.GetChoiceSettingCollectionValue()
			if len(choiceCollectionValues) > 0 {
				mappedCollection, err := mapChoiceSettingCollection(ctx, choiceCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map child choice setting collection: %w", err)
				}
				childItem.ChoiceSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.ChoiceSettingCollectionValue = make([]ChoiceSettingCollectionStruct, 0)
			}

		case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
			groupCollectionValues := typedChild.GetGroupSettingCollectionValue()
			if len(groupCollectionValues) > 0 {
				mappedCollection, err := mapGroupSettingCollection(ctx, groupCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map child group setting collection: %w", err)
				}
				childItem.GroupSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.GroupSettingCollectionValue = make([]GroupSettingCollectionStruct, 0)
			}
		}

		result = append(result, childItem)
	}

	return result, nil
}

// mapChoiceSettingCollectionChildren converts  choice setting collection children to our model
func mapChoiceSettingCollectionChildren(ctx context.Context, children []graphmodels.DeviceManagementConfigurationSettingInstanceable) ([]ChoiceSettingCollectionChild, error) {
	var result []ChoiceSettingCollectionChild

	for _, child := range children {
		childItem := ChoiceSettingCollectionChild{}

		// Map basic properties
		if odataType := child.GetOdataType(); odataType != nil {
			childItem.ODataType = types.StringValue(*odataType)
		}
		if settingDefId := child.GetSettingDefinitionId(); settingDefId != nil {
			childItem.SettingDefinitionId = types.StringValue(*settingDefId)
		}

		// Map instance template reference
		if instanceTemplateRef := child.GetSettingInstanceTemplateReference(); instanceTemplateRef != nil {
			childItem.SettingInstanceTemplateReference = mapInstanceTemplateReference(instanceTemplateRef)
		}

		// Type-specific mapping (choice collection children have limited types)
		switch typedChild := child.(type) {
		case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
			if simpleValue := typedChild.GetSimpleSettingValue(); simpleValue != nil {
				mappedSimpleValue, err := mapSimpleSettingValue(ctx, simpleValue)
				if err != nil {
					return nil, fmt.Errorf("failed to map choice collection child simple setting value: %w", err)
				}
				childItem.SimpleSettingValue = mappedSimpleValue
			}

		case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
			simpleCollectionValues := typedChild.GetSimpleSettingCollectionValue()
			if len(simpleCollectionValues) > 0 {
				mappedCollection, err := mapSimpleSettingCollection(ctx, simpleCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map choice collection child simple setting collection: %w", err)
				}
				childItem.SimpleSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.SimpleSettingCollectionValue = make([]SimpleSettingCollectionStruct, 0)
			}
		}

		result = append(result, childItem)
	}

	return result, nil
}

// mapGroupSettingCollectionChildren converts  group setting collection children to our model
func mapGroupSettingCollectionChildren(ctx context.Context, children []graphmodels.DeviceManagementConfigurationSettingInstanceable) ([]GroupSettingCollectionChild, error) {
	var result []GroupSettingCollectionChild

	for _, child := range children {
		childItem := GroupSettingCollectionChild{}

		// Map basic properties
		if odataType := child.GetOdataType(); odataType != nil {
			childItem.ODataType = types.StringValue(*odataType)
		}
		if settingDefId := child.GetSettingDefinitionId(); settingDefId != nil {
			childItem.SettingDefinitionId = types.StringValue(*settingDefId)
		}

		// Map instance template reference
		if instanceTemplateRef := child.GetSettingInstanceTemplateReference(); instanceTemplateRef != nil {
			childItem.SettingInstanceTemplateReference = mapInstanceTemplateReference(instanceTemplateRef)
		}

		// Type-specific mapping (group children can have all types)
		switch typedChild := child.(type) {
		case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
			if simpleValue := typedChild.GetSimpleSettingValue(); simpleValue != nil {
				mappedSimpleValue, err := mapSimpleSettingValue(ctx, simpleValue)
				if err != nil {
					return nil, fmt.Errorf("failed to map group child simple setting value: %w", err)
				}
				childItem.SimpleSettingValue = mappedSimpleValue
			}

		case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
			simpleCollectionValues := typedChild.GetSimpleSettingCollectionValue()
			if len(simpleCollectionValues) > 0 {
				mappedCollection, err := mapSimpleSettingCollection(ctx, simpleCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map group child simple setting collection: %w", err)
				}
				childItem.SimpleSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.SimpleSettingCollectionValue = make([]SimpleSettingCollectionStruct, 0)
			}

		case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
			if choiceValue := typedChild.GetChoiceSettingValue(); choiceValue != nil {
				mappedChoiceValue, err := mapChoiceSettingValue(ctx, choiceValue)
				if err != nil {
					return nil, fmt.Errorf("failed to map group child choice setting value: %w", err)
				}
				childItem.ChoiceSettingValue = mappedChoiceValue
			}

		case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
			choiceCollectionValues := typedChild.GetChoiceSettingCollectionValue()
			if len(choiceCollectionValues) > 0 {
				mappedCollection, err := mapChoiceSettingCollection(ctx, choiceCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map group child choice setting collection: %w", err)
				}
				childItem.ChoiceSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.ChoiceSettingCollectionValue = make([]ChoiceSettingCollectionStruct, 0)
			}

		case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
			groupCollectionValues := typedChild.GetGroupSettingCollectionValue()
			if len(groupCollectionValues) > 0 {
				mappedCollection, err := mapGroupSettingCollection(ctx, groupCollectionValues)
				if err != nil {
					return nil, fmt.Errorf("failed to map group child group setting collection: %w", err)
				}
				childItem.GroupSettingCollectionValue = mappedCollection
			} else {
				// FIXED: Always initialize as empty slice
				childItem.GroupSettingCollectionValue = make([]GroupSettingCollectionStruct, 0)
			}
		}

		result = append(result, childItem)
	}

	return result, nil
}

// mapInstanceTemplateReference converts  instance template reference to our model
func mapInstanceTemplateReference(ref graphmodels.DeviceManagementConfigurationSettingInstanceTemplateReferenceable) *SettingInstanceTemplateReference {
	if ref == nil {
		return nil
	}

	templateRef := &SettingInstanceTemplateReference{}
	if templateId := ref.GetSettingInstanceTemplateId(); templateId != nil {
		templateRef.SettingInstanceTemplateId = types.StringValue(*templateId)
	}

	return templateRef
}

// mapValueTemplateReference converts  value template reference to our model
func mapValueTemplateReference(ref graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable) *SettingValueTemplateReference {
	if ref == nil {
		return nil
	}

	templateRef := &SettingValueTemplateReference{}
	if templateId := ref.GetSettingValueTemplateId(); templateId != nil {
		templateRef.SettingValueTemplateId = types.StringValue(*templateId)
	}
	if useDefault := ref.GetUseTemplateDefault(); useDefault != nil {
		templateRef.UseTemplateDefault = *useDefault
	}

	return templateRef
}
