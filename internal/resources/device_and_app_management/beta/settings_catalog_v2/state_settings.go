package graphBetaSettingsCatalog

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteSettings []graphmodels.DeviceManagementConfigurationSettingable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings are nil")
		return
	}

	tflog.Debug(ctx, "Starting to map settings state to Terraform state")

	data.Settings = make([]DeviceManagementConfigurationSettingResourceModel, len(remoteSettings))
	for i, setting := range remoteSettings {
		// Base model for each setting
		settingModel := DeviceManagementConfigurationSettingResourceModel{
			ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationSetting"),
		}

		// Check for the SettingInstance within the Setting and map accordingly
		if instance := setting.GetSettingInstance(); instance != nil {
			settingInstance := &DeviceManagementConfigurationSettingInstanceResourceModel{
				SettingDefinitionID: types.StringValue(state.StringPtrToString(instance.GetSettingDefinitionId())),
				ODataType:           types.StringValue(state.StringPtrToString(instance.GetOdataType())),
			}

			// Map based on specific instance type
			switch inst := instance.(type) {
			case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingInstance)
				mapSimpleSettingInstance(ctx, inst, settingInstance)

			case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingInstance)
				mapChoiceSettingInstance(ctx, inst, settingInstance)

			case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingCollectionInstance)
				mapSimpleSettingCollectionInstance(ctx, inst, settingInstance)

			case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingCollectionInstance)
				mapChoiceSettingCollectionInstance(ctx, inst, settingInstance)

			case graphmodels.DeviceManagementConfigurationGroupSettingInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationGroupSettingInstance)
				mapGroupSettingInstance(ctx, inst, settingInstance)

			case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
				settingInstance.ODataType = types.StringValue(DeviceManagementConfigurationGroupSettingCollectionInstance)
				mapGroupSettingCollectionInstance(ctx, inst, settingInstance)

			default:
				tflog.Debug(ctx, fmt.Sprintf("Unknown setting instance type: %T", inst))
			}

			// Attach the mapped instance to the model
			settingModel.SettingInstance = settingInstance
		}

		// Add mapped setting to the data's settings slice
		data.Settings[i] = settingModel

		tflog.Debug(ctx, fmt.Sprintf("Mapped setting %d with type %s", i, settingModel.ODataType.ValueString()))
	}

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}

// Helper functions for mapping specific setting types

// mapSimpleSettingInstance maps a simple setting instance to Terraform state
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

		case graphmodels.DeviceManagementConfigurationSecretSettingValueable:
			if secretVal := v.GetValue(); secretVal != nil {
				simpleSettingValue.SecretValue = types.StringValue(*secretVal)
			}
			if valueState := v.GetValueState(); valueState != nil {
				simpleSettingValue.State = types.StringValue(valueState.String())
			}

		case graphmodels.DeviceManagementConfigurationReferenceSettingValueable:
			if refVal := v.GetValue(); refVal != nil {
				simpleSettingValue.ReferenceValue = types.StringValue(*refVal)
			}
			if note := v.GetNote(); note != nil {
				simpleSettingValue.Note = types.StringValue(*note)
			}

		case graphmodels.DeviceManagementConfigurationStringSettingValueable:
			if strVal := v.GetValue(); strVal != nil {
				simpleSettingValue.StringValue = types.StringValue(*strVal)
			}
		}

		settingInstance.SimpleSettingValue = simpleSettingValue
	}

	tflog.Debug(ctx, "Mapped simple setting instance to Terraform state")
}

func mapChoiceSettingInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if choiceValue := instance.GetChoiceSettingValue(); choiceValue != nil {
		choiceSettingValue := &ChoiceSettingValueResourceModel{
			ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"),
		}

		// Map value and children
		if value := choiceValue.GetValue(); value != nil {
			choiceSettingValue.StringValue = types.StringValue(*value)
		}

		// Map children if they exist
		if children := choiceValue.GetChildren(); len(children) > 0 {
			childrenModels := make([]DeviceManagementConfigurationSettingInstanceResourceModel, len(children))
			for i, child := range children {
				childModel := DeviceManagementConfigurationSettingInstanceResourceModel{
					ODataType:           types.StringValue(state.StringPtrToString(child.GetOdataType())),
					SettingDefinitionID: types.StringValue(state.StringPtrToString(child.GetSettingDefinitionId())),
				}
				// Ensure mapping based on child type
				switch childInst := child.(type) {
				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					mapSimpleSettingInstance(ctx, childInst, &childModel)
				}
				childrenModels[i] = childModel
			}
			choiceSettingValue.Children = childrenModels
		}

		settingInstance.ChoiceSettingValue = choiceSettingValue
	}
	tflog.Debug(ctx, "Mapped choice setting instance to Terraform state")
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

func mapGroupSettingInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationGroupSettingInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if groupValue := instance.GetGroupSettingValue(); groupValue != nil {
		groupSettingValue := &GroupSettingValueResourceModel{
			ODataType: types.StringValue(DeviceManagementConfigurationGroupSettingValue),
		}

		// Map children if they exist
		if children := groupValue.GetChildren(); len(children) > 0 {
			childrenModels := make([]DeviceManagementConfigurationSettingInstanceResourceModel, 0, len(children))

			for _, child := range children {
				childModel := DeviceManagementConfigurationSettingInstanceResourceModel{
					SettingDefinitionID: types.StringValue(state.StringPtrToString(child.GetSettingDefinitionId())),
				}

				// Determine child type and map accordingly
				switch childInst := child.(type) {
				case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
					childModel.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingInstance)
					mapChoiceSettingInstance(ctx, childInst, &childModel)

				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					childModel.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingInstance)
					mapSimpleSettingInstance(ctx, childInst, &childModel)

				case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
					childModel.ODataType = types.StringValue(DeviceManagementConfigurationChoiceSettingCollectionInstance)
					mapChoiceSettingCollectionInstance(ctx, childInst, &childModel)

				case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
					childModel.ODataType = types.StringValue(DeviceManagementConfigurationSimpleSettingCollectionInstance)
					mapSimpleSettingCollectionInstance(ctx, childInst, &childModel)
				}

				childrenModels = append(childrenModels, childModel)
			}

			groupSettingValue.Children = childrenModels
		}

		settingInstance.GroupSettingValue = groupSettingValue
	}
	tflog.Debug(ctx, "Mapped group setting instance to Terraform state")
}

func mapGroupSettingCollectionInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
	if collectionValues := instance.GetGroupSettingCollectionValue(); len(collectionValues) > 0 {
		groupSettingCollectionValue := &GroupSettingCollectionValueResourceModel{}
		childrenModels := make([]DeviceManagementConfigurationSettingInstanceResourceModel, 0)

		for _, collectionValue := range collectionValues {
			if children := collectionValue.GetChildren(); len(children) > 0 {
				for _, child := range children {
					childModel := DeviceManagementConfigurationSettingInstanceResourceModel{
						ODataType:           types.StringValue(getODataType(child)), // Helper function needed
						SettingDefinitionID: types.StringValue(state.StringPtrToString(child.GetSettingDefinitionId())),
					}

					// Map the appropriate value based on the type
					switch childInst := child.(type) {
					case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
						mapChoiceSettingInstance(ctx, childInst, &childModel)

					case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
						mapGroupSettingCollectionInstance(ctx, childInst, &childModel)

					case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
						mapChoiceSettingCollectionInstance(ctx, childInst, &childModel)

					case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
						mapSimpleSettingInstance(ctx, childInst, &childModel)

					case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
						mapSimpleSettingCollectionInstance(ctx, childInst, &childModel)
					}

					childrenModels = append(childrenModels, childModel)
				}
			}
		}

		groupSettingCollectionValue.Children = childrenModels
		settingInstance.GroupSettingCollectionValue = groupSettingCollectionValue
	}
	tflog.Debug(ctx, "Mapped group setting collection instance to Terraform state")
}

// Helper function to get OData type
func getODataType(instance interface{}) string {
	switch instance.(type) {
	case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
		return DeviceManagementConfigurationChoiceSettingInstance
	case graphmodels.DeviceManagementConfigurationGroupSettingCollectionInstanceable:
		return DeviceManagementConfigurationGroupSettingCollectionInstance
	case graphmodels.DeviceManagementConfigurationChoiceSettingCollectionInstanceable:
		return DeviceManagementConfigurationChoiceSettingCollectionInstance
	case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
		return DeviceManagementConfigurationSimpleSettingInstance
	case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
		return DeviceManagementConfigurationSimpleSettingCollectionInstance
	default:
		return ""
	}
}

// func mapSettingGroupInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSettingGroupInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
// 	settingInstance.SettingGroupSettingValue = &SettingGroupSettingValueResourceModel{
// 		ODataType: types.StringValue(DeviceManagementConfigurationSettingGroupInstance),
// 	}
// 	tflog.Debug(ctx, "Mapped setting group instance to Terraform state")
// }

// func mapSettingGroupCollectionInstance(ctx context.Context, instance graphmodels.DeviceManagementConfigurationSettingGroupCollectionInstanceable, settingInstance *DeviceManagementConfigurationSettingInstanceResourceModel) {
// 	settingInstance.SettingGroupCollectionValue = &SettingGroupCollectionValueResourceModel{
// 		ODataType: types.StringValue(DeviceManagementConfigurationSettingGroupCollectionInstance),
// 	}
// 	tflog.Debug(ctx, "Mapped setting group collection instance to Terraform state")
// }
