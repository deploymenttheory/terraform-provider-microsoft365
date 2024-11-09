package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteSettings graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map settings state to Terraform state")

	if settings := remoteSettings.GetValue(); settings != nil {
		data.Settings = make([]DeviceManagementConfigurationSettingResourceModel, len(settings))
		for i, setting := range settings {
			if instance := setting.GetSettingInstance(); instance != nil {
				settingModel := DeviceManagementConfigurationSettingResourceModel{
					ODataType: types.StringValue(DeviceManagementConfigurationSetting),
				}

				settingInstance := &DeviceManagementConfigurationSettingInstanceResourceModel{
					SettingDefinitionID: types.StringValue(state.StringPtrToString(instance.GetSettingDefinitionId())),
				}

				switch inst := instance.(type) {
				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					tflog.Debug(ctx, "Mapping SimpleSettingInstance")
					settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingInstance)
					mapSimpleSettingInstance(ctx, inst, settingInstance)

				case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
					tflog.Debug(ctx, "Mapping ChoiceSettingInstance")
					settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingInstance)
					mapChoiceSettingInstance(ctx, inst, settingInstance)

				case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
					tflog.Debug(ctx, "Mapping SimpleSettingCollectionInstance")
					settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingCollectionInstance)
					mapSimpleSettingCollectionInstance(ctx, inst, settingInstance)

				case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
					tflog.Debug(ctx, "Mapping ChoiceSettingCollectionInstance")
					settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingCollectionInstance)
					mapChoiceSettingCollectionInstance(ctx, inst, settingInstance)

					// case graphmodels.DeviceManagementConfigurationSettingGroupInstanceable:
					// 	tflog.Debug(ctx, "Mapping SettingGroupInstance")
					// 	settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSettingGroupInstance)
					// 	mapGroupSettingInstance(ctx, inst, settingInstance)

					// case graphmodels.DeviceManagementConfigurationSettingGroupCollectionInstanceable:
					// 	tflog.Debug(ctx, "Mapping SettingGroupCollectionInstance")
					// 	settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSettingGroupCollectionInstance)
					// 	mapGroupSettingCollectionInstance(ctx, inst, settingInstance)
				}

				settingModel.SettingInstance = settingInstance
				data.Settings[i] = settingModel
			}
		}
	}

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}

// Helper functions for mapping specific setting types

func mapSimpleSettingInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if simpleValue := instance.GetSimpleSettingValue(); simpleValue != nil {
		simpleSettingValue := &SimpleSettingValueResourceModel{
			ODataType: types.StringValue(DeviceManagementConfigurationSimpleSettingInstance),
		}
		switch v := simpleValue.(type) {
		case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
			if intVal := v.GetValue(); intVal != nil {
				simpleSettingValue.IntValue = state.Int32PtrToTypeInt32(intVal)
			}
		case graphmodels.DeviceManagementConfigurationStringSettingValueable:
			if strVal := v.GetValue(); strVal != nil {
				simpleSettingValue.StringValue = types.StringValue(*strVal)
			}
		}
		settingInstance.SimpleSettingValue = simpleSettingValue
	}
}

func mapChoiceSettingInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if choiceValue := instance.GetChoiceSettingValue(); choiceValue != nil {
		choiceSettingValue := &ChoiceSettingValueResourceModel{
			ODataType:   types.StringValue(DeviceManagementConfigurationChoiceSettingValue),
			StringValue: types.StringValue(state.StringPtrToString(choiceValue.GetValue())),
		}
		settingInstance.ChoiceSettingValue = choiceSettingValue
	}
}

func mapSimpleSettingCollectionInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if collectionValues := instance.GetSimpleSettingCollectionValue(); len(collectionValues) > 0 {
		simpleCollectionValue := &SimpleCollectionValueResourceModel{
			ODataType: types.StringValue(DeviceManagementConfigurationSimpleSettingCollectionInstance),
		}

		for _, collectionValue := range collectionValues {
			switch v := collectionValue.(type) {
			case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
				if intVal := v.GetValue(); intVal != nil {
					simpleCollectionValue.IntValue = append(simpleCollectionValue.IntValue, state.Int32PtrToTypeInt32(intVal))
				}
			case graphmodels.DeviceManagementConfigurationStringSettingValueable:
				if strVal := v.GetValue(); strVal != nil {
					simpleCollectionValue.StringValue = append(simpleCollectionValue.StringValue, types.StringValue(*strVal))
				}
			}
		}
		settingInstance.SimpleCollectionValue = simpleCollectionValue
		tflog.Debug(ctx, "Mapped simple setting collection with integer and string values to Terraform state")
	}
}

func mapChoiceSettingCollectionInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if collectionValues := instance.GetChoiceSettingCollectionValue(); len(collectionValues) > 0 {
		choiceCollectionValue := &ChoiceCollectionValueResourceModel{
			ODataType: types.StringValue(DeviceManagementConfigurationChoiceSettingCollectionInstance),
		}

		for _, collectionValue := range collectionValues {
			if strVal := collectionValue.GetValue(); strVal != nil {
				choiceCollectionValue.StringValue = append(choiceCollectionValue.StringValue, types.StringValue(*strVal))
			}
		}
		settingInstance.ChoiceCollectionValue = choiceCollectionValue
		tflog.Debug(ctx, "Mapped choice setting collection with string values to Terraform state")
	}
}

// func mapGroupSettingInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSettingGroupInstanceable, settingInstance *DeviceManagementConfigurationSettingInstance) {
// 	// Initialize the GroupSettingValue for Terraform state with the correct OData type
// 	groupValue := &DeviceManagementConfigurationGroupSettingValueResourceModel{
// 		ODataType: types.StringValue(DeviceManagementConfigurationSettingGroupInstance),
// 	}

// 	// Since we don't have direct access to children, let's assume handling is based on known instance types
// 	instanceID := state.StringPtrToString(instance.GetSettingDefinitionId())
// 	tflog.Debug(ctx, "Mapping GroupSettingInstance", map[string]interface{}{"SettingDefinitionID": instanceID})

// 	// In the absence of child accessors, assume children are set in a predefined structure or need recursive construction.
// 	childInstances := retrieveChildInstances(ctx, instance) // Define this helper based on custom logic.
// 	for _, childInstance := range childInstances {
// 		childModel := DeviceManagementConfigurationSettingInstance{
// 			ODataType:           types.StringValue(state.StringPtrToString(childInstance.GetOdataType())),
// 			SettingDefinitionID: types.StringValue(state.StringPtrToString(childInstance.GetSettingDefinitionId())),
// 		}

// 		// Recursively map each specific child type
// 		switch child := childInstance.(type) {
// 		case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
// 			mapSimpleSettingInstance(ctx, child, &childModel)
// 		case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
// 			mapChoiceSettingInstance(ctx, child, &childModel)
// 		case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
// 			mapSimpleSettingCollectionInstance(ctx, child, &childModel)
// 		case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
// 			mapChoiceSettingCollectionInstance(ctx, child, &childModel)
// 		case graphmodels.DeviceManagementConfigurationSettingGroupInstanceable:
// 			mapGroupSettingInstance(ctx, child, &childModel)
// 		case graphmodels.DeviceManagementConfigurationSettingGroupCollectionInstanceable:
// 			mapGroupSettingCollectionInstance(ctx, child, &childModel)
// 		}

// 		// Append the mapped child model to the group's Children slice
// 		groupValue.Children = append(groupValue.Children, childModel)
// 	}

// 	// Set the completed GroupSettingValue on the main setting instance
// 	settingInstance.GroupSettingValue = groupValue
// 	tflog.Debug(ctx, "Mapped group setting instance to Terraform state")
// }

// func mapGroupSettingCollectionInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSettingGroupCollectionInstanceable, settingInstance *DeviceManagementConfigurationSettingInstance) {
// 	// Retrieve and process values of the group setting collection
// 	if collectionValues := instance.GetValue(); len(collectionValues) > 0 {
// 		groupCollectionValue := &DeviceManagementConfigurationGroupCollectionValueResourceModel{
// 			ODataType: types.StringValue(DeviceManagementConfigurationSettingGroupCollectionInstance),
// 		}

// 		// Iterate over each group setting within the collection
// 		for _, groupSetting := range collectionValues {
// 			groupValue := &DeviceManagementConfigurationGroupSettingValueResourceModel{
// 				ODataType: types.StringValue(DeviceManagementConfigurationSettingGroupInstance),
// 			}

// 			// Process child instances of each group setting
// 			if children := groupSetting.GetValue(); len(children) > 0 {
// 				for _, childInstance := range children {
// 					childModel := DeviceManagementConfigurationSettingInstance{
// 						ODataType:           types.StringValue(state.StringPtrToString(childInstance.GetOdataType())),
// 						SettingDefinitionID: types.StringValue(state.StringPtrToString(childInstance.GetSettingDefinitionId())),
// 					}

// 					// Map each specific child type using the correct mapping functions
// 					switch child := childInstance.(type) {
// 					case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
// 						mapSimpleSettingInstance(ctx, child, &childModel)
// 					case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
// 						mapChoiceSettingInstance(ctx, child, &childModel)
// 					case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
// 						mapSimpleSettingCollectionInstance(ctx, child, &childModel)
// 					case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
// 						mapChoiceSettingCollectionInstance(ctx, child, &childModel)
// 					case graphmodels.DeviceManagementConfigurationSettingGroupInstanceable:
// 						mapGroupSettingInstance(ctx, child, &childModel)
// 					case graphmodels.DeviceManagementConfigurationSettingGroupCollectionInstanceable:
// 						mapGroupSettingCollectionInstance(ctx, child, &childModel)
// 					}

// 					// Append each mapped child model to the group's Children slice
// 					groupValue.Children = append(groupValue.Children, childModel)
// 				}
// 			}

// 			// Append each mapped group value to the group collection's Children slice
// 			groupCollectionValue.Children = append(groupCollectionValue.Children, *groupValue)
// 		}

// 		// Set the fully constructed group collection value on the main setting instance
// 		settingInstance.GroupCollectionValue = groupCollectionValue
// 		tflog.Debug(ctx, "Mapped group setting collection instance to Terraform state")
// 	}
// }
