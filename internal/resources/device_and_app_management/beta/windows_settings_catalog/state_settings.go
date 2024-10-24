package graphBetaWindowsSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *WindowsSettingsCatalogProfileResourceModel, remoteSettings graphmodels.DeviceManagementConfigurationSettingCollectionResponseable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map settings state to Terraform state")

	// Map settings from the response value
	if settings := remoteSettings.GetValue(); settings != nil {
		data.Settings = make([]DeviceManagementConfigurationSettingResourceModel, len(settings))
		for i, setting := range settings {
			if instance := setting.GetSettingInstance(); instance != nil {
				settingModel := DeviceManagementConfigurationSettingResourceModel{
					ID:        types.StringValue(state.StringPtrToString(setting.GetId())),
					ODataType: types.StringValue(state.StringPtrToString(instance.GetOdataType())),
				}

				// Map setting instance
				settingInstance := &DeviceManagementConfigurationSettingInstance{
					ODataType:           types.StringValue(state.StringPtrToString(instance.GetOdataType())),
					SettingDefinitionID: types.StringValue(state.StringPtrToString(instance.GetSettingDefinitionId())),
				}

				// Map setting instance template reference if present
				if templateRef := instance.GetSettingInstanceTemplateReference(); templateRef != nil {
					settingInstance.SettingInstanceTemplateReference = &DeviceManagementConfigurationSettingInstanceTemplateReference{
						ODataType:                 types.StringValue(state.StringPtrToString(templateRef.GetOdataType())),
						SettingInstanceTemplateID: types.StringValue(state.StringPtrToString(templateRef.GetSettingInstanceTemplateId())),
					}
				}

				// Handle different setting types
				switch v := instance.(type) {
				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					if simpleValue := v.GetSimpleSettingValue(); simpleValue != nil {
						settingInstance.ChoiceSettingValue = mapSimpleSettingValue(simpleValue)
					}

				case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
					if choiceValue := v.GetChoiceSettingValue(); choiceValue != nil {
						settingInstance.ChoiceSettingValue = &DeviceManagementConfigurationChoiceSettingValue{
							ODataType:   types.StringValue(state.StringPtrToString(choiceValue.GetOdataType())),
							StringValue: state.StringPtrToString(choiceValue.GetValue()),
						}
					}

				case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
					if collectionValues := v.GetSimpleSettingCollectionValue(); len(collectionValues) > 0 {
						settingInstance.ChoiceSettingValue = &DeviceManagementConfigurationChoiceSettingValue{
							Children: make([]DeviceManagementConfigurationSettingInstance, len(collectionValues)),
						}

						for j, collectionValue := range collectionValues {
							child := mapSimpleSettingValue(collectionValue)
							settingInstance.ChoiceSettingValue.Children[j] = DeviceManagementConfigurationSettingInstance{
								ChoiceSettingValue: child,
							}
						}
					}
				}

				settingModel.SettingInstance = settingInstance
				data.Settings[i] = settingModel
			}
		}
	}

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}

// Helper function to map simple setting values
func mapSimpleSettingValue(value graphmodels.DeviceManagementConfigurationSimpleSettingValueable) *DeviceManagementConfigurationChoiceSettingValue {
	result := &DeviceManagementConfigurationChoiceSettingValue{
		ODataType: types.StringValue(state.StringPtrToString(value.GetOdataType())),
	}

	switch v := value.(type) {
	case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
		if intVal := v.GetValue(); intVal != nil {
			result.IntValue = int32(*intVal)
		}
	case graphmodels.DeviceManagementConfigurationStringSettingValueable:
		if strVal := v.GetValue(); strVal != nil {
			result.StringValue = *strVal
		}
	}

	return result
}
