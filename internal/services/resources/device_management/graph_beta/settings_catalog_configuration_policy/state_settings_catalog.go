package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateConfigurationPolicySettings maps settings from Graph  models to Terraform state
func StateConfigurationPolicySettings(ctx context.Context, data *SettingsCatalogProfileResourceModel, settingsResponse graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) error {
	tflog.Debug(ctx, "Starting to map settings from Graph models to Terraform state")

	if settingsResponse == nil {
		tflog.Debug(ctx, "No settings response data to process")
		return nil
	}

	settings := settingsResponse.GetValue()
	if len(settings) == 0 {
		tflog.Debug(ctx, "Settings array is empty")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Processing %d settings from API response", len(settings)))

	// Convert settings to our model
	deviceConfigModel := &DeviceConfigV2GraphServiceResourceModel{}
	var mappedSettings []Setting
	successfulMappings := 0
	failedMappings := 0

	for i, apiSetting := range settings {
		if apiSetting == nil {
			tflog.Warn(ctx, fmt.Sprintf("Setting at index %d is nil", i))
			failedMappings++
			continue
		}

		// Log details about the setting being processed
		settingId := "unknown"
		if id := apiSetting.GetId(); id != nil {
			settingId = *id
		}

		settingDefId := "unknown"
		if instance := apiSetting.GetSettingInstance(); instance != nil {
			if defId := instance.GetSettingDefinitionId(); defId != nil {
				settingDefId = *defId
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("Mapping setting %d: ID=%s, DefinitionID=%s", i, settingId, settingDefId))

		setting, err := mapSettingToModel(ctx, apiSetting)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to map setting %d (ID: %s, DefinitionID: %s): %s", i, settingId, settingDefId, err.Error()))
			failedMappings++
			continue // This is where settings get dropped!
		}

		if setting != nil {
			mappedSettings = append(mappedSettings, *setting)
			successfulMappings++
			tflog.Debug(ctx, fmt.Sprintf("Successfully mapped setting %d (ID: %s)", i, settingId))
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Setting %d (ID: %s) mapped to nil", i, settingId))
			failedMappings++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping summary: %d successful, %d failed, %d total from API", successfulMappings, failedMappings, len(settings)))

	if failedMappings > 0 {
		tflog.Error(ctx, fmt.Sprintf("WARNING: %d settings failed to map - this will cause state inconsistency!", failedMappings))
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
		setting.ID = convert.GraphToFrameworkString(id)
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
		settingInstance.ODataType = convert.GraphToFrameworkString(odataType)
	}

	// Map setting definition ID
	if settingDefId := instance.GetSettingDefinitionId(); settingDefId != nil {
		settingInstance.SettingDefinitionId = convert.GraphToFrameworkString(settingDefId)
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
			// Always initialize as empty slice
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
			// Always initialize as empty slice
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
			// Always initialize as empty slice
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
		simpleValue.ODataType = convert.GraphToFrameworkString(odataType)
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
			simpleValue.Value = convert.GraphToFrameworkString(secretVal)
		}

		if valueState := typedValue.GetValueState(); valueState != nil {
			simpleValue.ValueState = types.StringValue(valueState.String())
		}

	case graphmodels.DeviceManagementConfigurationStringSettingValueable:
		if stringVal := typedValue.GetValue(); stringVal != nil {
			simpleValue.Value = convert.GraphToFrameworkString(stringVal)
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
		choiceValue.Value = convert.GraphToFrameworkString(val)
	}

	// Map value template reference
	if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
		choiceValue.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
	}

	// Always initialize children, even if empty
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
		if value == nil {
			continue
		}

		collectionItem := SimpleSettingCollectionStruct{}

		if odataType := value.GetOdataType(); odataType != nil {
			collectionItem.ODataType = convert.GraphToFrameworkString(odataType)
		}

		if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
			collectionItem.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
		}

		if stringVal, ok := value.(graphmodels.DeviceManagementConfigurationStringSettingValueable); ok {
			if val := stringVal.GetValue(); val != nil {
				collectionItem.Value = convert.GraphToFrameworkString(val)
			}
		} else if intVal, ok := value.(graphmodels.DeviceManagementConfigurationIntegerSettingValueable); ok {
			if val := intVal.GetValue(); val != nil {
				collectionItem.Value = types.StringValue(strconv.Itoa(int(*val)))
			}
		} else if choiceVal, ok := value.(graphmodels.DeviceManagementConfigurationChoiceSettingValueable); ok {
			if val := choiceVal.GetValue(); val != nil {
				collectionItem.Value = convert.GraphToFrameworkString(val)
			}
		} else {
			return nil, fmt.Errorf("unsupported simple setting collection value type: %T", value)
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

		if val := value.GetValue(); val != nil {
			collectionItem.Value = types.StringValue(*val)
		}

		if valueTemplateRef := value.GetSettingValueTemplateReference(); valueTemplateRef != nil {
			collectionItem.SettingValueTemplateReference = mapValueTemplateReference(valueTemplateRef)
		}

		// Always initialize children, even if empty
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

		// Always initialize children, even if empty
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

// mapChoiceSettingChildren converts Graph  choice setting children to our model
func mapChoiceSettingChildren(ctx context.Context, children []graphmodels.DeviceManagementConfigurationSettingInstanceable) ([]ChoiceSettingChild, error) {
	var result []ChoiceSettingChild

	for _, child := range children {
		if child == nil {
			continue
		}

		childItem := ChoiceSettingChild{}

		if odataType := child.GetOdataType(); odataType != nil {
			childItem.ODataType = convert.GraphToFrameworkString(odataType)
		}

		if settingDefId := child.GetSettingDefinitionId(); settingDefId != nil {
			childItem.SettingDefinitionId = convert.GraphToFrameworkString(settingDefId)
		}

		if instanceTemplateRef := child.GetSettingInstanceTemplateReference(); instanceTemplateRef != nil {
			childItem.SettingInstanceTemplateReference = mapInstanceTemplateReference(instanceTemplateRef)
		}

		// Handle different child types
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
				// Always initialize as empty slice
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
				// Always initialize as empty slice
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
				// Always initialize as empty slice
				childItem.GroupSettingCollectionValue = make([]GroupSettingCollectionStruct, 0)
			}

		default:
			return nil, fmt.Errorf("unsupported choice setting child type: %T", typedChild)
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
				// Always initialize as empty slice
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
				// Always initialize as empty slice
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
				// Always initialize as empty slice
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
				// Always initialize as empty slice
				childItem.GroupSettingCollectionValue = make([]GroupSettingCollectionStruct, 0)
			}
		}

		result = append(result, childItem)
	}

	return result, nil
}

// mapInstanceTemplateReference converts a Graph  setting instance template reference to our model
func mapInstanceTemplateReference(ref graphmodels.DeviceManagementConfigurationSettingInstanceTemplateReferenceable) *SettingInstanceTemplateReference {
	if ref == nil {
		return nil
	}

	templateRef := &SettingInstanceTemplateReference{}
	if templateId := ref.GetSettingInstanceTemplateId(); templateId != nil {
		templateRef.SettingInstanceTemplateId = convert.GraphToFrameworkString(templateId)
	}
	return templateRef
}

// mapValueTemplateReference converts a Graph  setting value template reference to our model
func mapValueTemplateReference(ref graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable) *SettingValueTemplateReference {
	if ref == nil {
		return nil
	}

	templateRef := &SettingValueTemplateReference{}
	if templateId := ref.GetSettingValueTemplateId(); templateId != nil {
		templateRef.SettingValueTemplateId = convert.GraphToFrameworkString(templateId)
	}
	if useDefault := ref.GetUseTemplateDefault(); useDefault != nil {
		templateRef.UseTemplateDefault = *useDefault
	}
	return templateRef
}
