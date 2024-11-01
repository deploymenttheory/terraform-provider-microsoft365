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
				// Always use the standard setting type at the top level
				settingModel := DeviceManagementConfigurationSettingResourceModel{
					ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationSetting"),
				}

				// Map setting instance
				settingInstance := &DeviceManagementConfigurationSettingInstance{
					SettingDefinitionID: types.StringValue(state.StringPtrToString(instance.GetSettingDefinitionId())),
				}

				// Set the correct instance OData type based on the instance type
				switch instance.(type) {
				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					settingInstance.ODataType = types.StringValue("#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance")
				case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
					settingInstance.ODataType = types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance")
				case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
					settingInstance.ODataType = types.StringValue("#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance")
				}

				// Handle different setting types
				switch v := instance.(type) {
				case graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable:
					if simpleValue := v.GetSimpleSettingValue(); simpleValue != nil {
						settingInstance.ChoiceSettingValue = &DeviceManagementConfigurationChoiceSettingValue{
							ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"),
						}

						switch sv := simpleValue.(type) {
						case graphmodels.DeviceManagementConfigurationIntegerSettingValueable:
							if intVal := sv.GetValue(); intVal != nil {
								settingInstance.ChoiceSettingValue.IntValue = state.Int32PtrToTypeInt32(intVal)
							}
						}
					}

				case graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable:
					if choiceValue := v.GetChoiceSettingValue(); choiceValue != nil {
						settingInstance.ChoiceSettingValue = &DeviceManagementConfigurationChoiceSettingValue{
							ODataType:   types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"),
							StringValue: types.StringValue(state.StringPtrToString(choiceValue.GetValue())),
						}
					}

				case graphmodels.DeviceManagementConfigurationSimpleSettingCollectionInstanceable:
					if collectionValues := v.GetSimpleSettingCollectionValue(); len(collectionValues) > 0 {
						settingInstance.ChoiceSettingValue = &DeviceManagementConfigurationChoiceSettingValue{
							ODataType: types.StringValue("#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"),
							Children:  make([]DeviceManagementConfigurationSettingInstance, len(collectionValues)),
						}

						for j, collectionValue := range collectionValues {
							switch cv := collectionValue.(type) {
							case graphmodels.DeviceManagementConfigurationStringSettingValueable:
								if strVal := cv.GetValue(); strVal != nil {
									settingInstance.ChoiceSettingValue.Children[j] = DeviceManagementConfigurationSettingInstance{
										ChoiceSettingValue: &DeviceManagementConfigurationChoiceSettingValue{
											StringValue: types.StringValue(*strVal),
										},
									}
								}
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
