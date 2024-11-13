package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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

	var technologiesStr []string
	for _, tech := range data.Technologies {
		technologiesStr = append(technologiesStr, tech.ValueString())
	}
	parsedTechnologies, _ := graphmodels.ParseDeviceManagementConfigurationTechnologies(strings.Join(technologiesStr, ","))
	profile.SetTechnologies(parsedTechnologies.(*graphmodels.DeviceManagementConfigurationTechnologies))

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

	// Create serialization writer to see the final JSON
	factory := jsonserialization.NewJsonSerializationWriterFactory()
	writer, _ := factory.GetSerializationWriter("application/json")

	// Write the profile to JSON
	_ = writer.WriteObjectValue("", profile)

	// Get the JSON bytes
	jsonBytes, _ := writer.GetSerializedContent()

	// Pretty print the JSON for debugging
	var prettyJSON map[string]interface{}
	_ = json.Unmarshal(jsonBytes, &prettyJSON)
	debugJSON, _ := json.MarshalIndent(prettyJSON, "", "  ")

	tflog.Debug(ctx, "Final JSON to be sent to API", map[string]interface{}{
		"json": string(debugJSON),
	})

	tflog.Debug(ctx, "Finished constructing Windows Settings Catalog resource")
	return profile, nil
}

func constructSettingsCatalogSettings(ctx context.Context, settingsJSON types.String) []graphmodels.DeviceManagementConfigurationSettingable {
	tflog.Debug(ctx, "Constructing settings catalog settings")

	// Parse the settings structure with settingsDetails array
	var settingsData struct {
		SettingsDetails []struct {
			ID              string `json:"id"`
			SettingInstance struct {
				ODataType           string `json:"@odata.type"`
				SettingDefinitionId string `json:"settingDefinitionId"`
				// For choice settings
				ChoiceSettingValue *struct {
					Children []struct {
						ODataType                   string `json:"@odata.type"`
						SettingDefinitionId         string `json:"settingDefinitionId"`
						GroupSettingCollectionValue []struct {
							Children []struct {
								ODataType           string `json:"@odata.type"`
								SettingDefinitionId string `json:"settingDefinitionId"`
								SimpleSettingValue  *struct {
									ODataType                     string      `json:"@odata.type"`
									Value                         string      `json:"value"`
									SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
								} `json:"simpleSettingValue,omitempty"`
								SettingInstanceTemplateReference interface{} `json:"settingInstanceTemplateReference"`
							} `json:"children"`
							SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
						} `json:"groupSettingCollectionValue,omitempty"`
						SettingInstanceTemplateReference interface{} `json:"settingInstanceTemplateReference"`
					} `json:"children"`
					Value                         string      `json:"value"`
					SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
				} `json:"choiceSettingValue,omitempty"`
				// For simple settings
				SimpleSettingValue *struct {
					ODataType                     string      `json:"@odata.type"`
					Value                         string      `json:"value"`
					SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
				} `json:"simpleSettingValue,omitempty"`
				// For simple collection settings
				SimpleSettingCollectionValue []struct {
					ODataType                     string      `json:"@odata.type"`
					Value                         string      `json:"value"`
					SettingValueTemplateReference interface{} `json:"settingValueTemplateReference"`
				} `json:"simpleSettingCollectionValue,omitempty"`
				SettingInstanceTemplateReference interface{} `json:"settingInstanceTemplateReference"`
			} `json:"settingInstance"`
		} `json:"settingsDetails"`
	}

	if err := json.Unmarshal([]byte(settingsJSON.ValueString()), &settingsData); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	// Add debug logging after unmarshaling
	tflog.Debug(ctx, "Unmarshaled settings data", map[string]interface{}{
		"data": settingsData,
	})

	settingsCollection := make([]graphmodels.DeviceManagementConfigurationSettingable, 0)

	for _, detail := range settingsData.SettingsDetails {
		baseSetting := graphmodels.NewDeviceManagementConfigurationSetting()

		switch detail.SettingInstance.ODataType {
		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
			instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			if detail.SettingInstance.ChoiceSettingValue != nil {
				choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
				choiceValue.SetValue(&detail.SettingInstance.ChoiceSettingValue.Value)

				// Handle nested children
				if len(detail.SettingInstance.ChoiceSettingValue.Children) > 0 {
					var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
					for _, child := range detail.SettingInstance.ChoiceSettingValue.Children {
						if child.ODataType == "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance" {
							groupInstance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
							groupInstance.SetOdataType(&child.ODataType)
							groupInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							// Handle group collection values
							var groupValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable
							for _, groupValue := range child.GroupSettingCollectionValue {
								groupSettingValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()

								// Handle nested children in group values
								var nestedChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
								for _, nestedChild := range groupValue.Children {
									if nestedChild.SimpleSettingValue != nil {
										simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
										simpleInstance.SetOdataType(&nestedChild.ODataType)
										simpleInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)

										simpleValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
										simpleValue.SetOdataType(&nestedChild.SimpleSettingValue.ODataType)
										simpleValue.SetValue(&nestedChild.SimpleSettingValue.Value)
										simpleInstance.SetSimpleSettingValue(simpleValue)

										nestedChildren = append(nestedChildren, simpleInstance)
									}
								}
								groupSettingValue.SetChildren(nestedChildren)
								groupValues = append(groupValues, groupSettingValue)
							}
							groupInstance.SetGroupSettingCollectionValue(groupValues)
							children = append(children, groupInstance)
						}
					}
					choiceValue.SetChildren(children)
				} else {
					choiceValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{})
				}
				instance.SetChoiceSettingValue(choiceValue)
			}

			baseSetting.SetSettingInstance(instance)

		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
			instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			if len(detail.SettingInstance.SimpleSettingCollectionValue) > 0 {
				var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
				for _, v := range detail.SettingInstance.SimpleSettingCollectionValue {
					stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
					stringValue.SetOdataType(&v.ODataType)
					stringValue.SetValue(&v.Value)
					values = append(values, stringValue)
				}
				instance.SetSimpleSettingCollectionValue(values)
			}

			baseSetting.SetSettingInstance(instance)

		case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
			instance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			if detail.SettingInstance.SimpleSettingValue != nil {
				value := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
				value.SetOdataType(&detail.SettingInstance.SimpleSettingValue.ODataType)
				value.SetValue(&detail.SettingInstance.SimpleSettingValue.Value)
				instance.SetSimpleSettingValue(value)
			}

			baseSetting.SetSettingInstance(instance)
		}

		settingsCollection = append(settingsCollection, baseSetting)
	}

	// Debug logging before returning
	tflog.Debug(ctx, "Constructed settings collection", map[string]interface{}{
		"count": len(settingsCollection),
	})

	return settingsCollection
}
