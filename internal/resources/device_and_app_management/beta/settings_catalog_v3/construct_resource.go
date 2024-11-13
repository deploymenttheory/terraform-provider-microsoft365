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
									ODataType                     string                                                                     `json:"@odata.type"`
									Value                         string                                                                     `json:"value"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"simpleSettingValue,omitempty"`
								SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
							} `json:"children"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"groupSettingCollectionValue,omitempty"`
						SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					} `json:"children"`
					Value                         string                                                                     `json:"value"`
					SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				} `json:"choiceSettingValue,omitempty"`

				// For choice setting collections
				ChoiceSettingCollectionValue []struct {
					Children []struct {
						ODataType           string `json:"@odata.type"`
						SettingDefinitionId string `json:"settingDefinitionId"`
						SimpleSettingValue  *struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         interface{}                                                                `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingValue,omitempty"`
						SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					} `json:"children"`
					Value                         string                                                                     `json:"value"`
					SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				} `json:"choiceSettingCollectionValue,omitempty"`

				// For group setting collections
				GroupSettingCollectionValue []struct {
					SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
					Children                      []struct {
						ODataType           string `json:"@odata.type"`
						SettingDefinitionId string `json:"settingDefinitionId"`
						ChoiceSettingValue  *struct {
							Value                         string                                                                     `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
							Children                      []struct {
								ODataType           string `json:"@odata.type"`
								SettingDefinitionId string `json:"settingDefinitionId"`
							} `json:"children"`
						} `json:"choiceSettingValue,omitempty"`
						SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
					} `json:"children"`
				} `json:"groupSettingCollectionValue,omitempty"`

				// For simple settings
				SimpleSettingValue *struct {
					ODataType                     string                                                                     `json:"@odata.type"`
					Value                         interface{}                                                                `json:"value"`
					SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				} `json:"simpleSettingValue,omitempty"`

				// For simple collection settings
				SimpleSettingCollectionValue []struct {
					ODataType                     string                                                                     `json:"@odata.type"`
					Value                         string                                                                     `json:"value"`
					SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
				} `json:"simpleSettingCollectionValue,omitempty"`

				SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
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
				switch detail.SettingInstance.SimpleSettingValue.ODataType {
				case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
					// Handle string value
					value := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
					value.SetOdataType(&detail.SettingInstance.SimpleSettingValue.ODataType)

					// Type assertion for string
					if stringValue, ok := detail.SettingInstance.SimpleSettingValue.Value.(string); ok {
						value.SetValue(&stringValue)
					} else {
						tflog.Error(ctx, "Expected string value but got different type", map[string]interface{}{
							"value": detail.SettingInstance.SimpleSettingValue.Value,
						})
					}
					instance.SetSimpleSettingValue(value)

				case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
					// Handle integer value
					value := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
					value.SetOdataType(&detail.SettingInstance.SimpleSettingValue.ODataType)

					// Type assertion for int
					if intValue, ok := detail.SettingInstance.SimpleSettingValue.Value.(float64); ok {
						intVal := int32(intValue) // Convert float64 to int
						value.SetValue(&intVal)
					} else {
						tflog.Error(ctx, "Expected integer value but got different type", map[string]interface{}{
							"value": detail.SettingInstance.SimpleSettingValue.Value,
						})
					}
					instance.SetSimpleSettingValue(value)
				}
			}

			baseSetting.SetSettingInstance(instance)

		case "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance":
			instance := graphmodels.NewDeviceManagementConfigurationChoiceSettingCollectionInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			// Process each item in ChoiceSettingCollectionValue
			if len(detail.SettingInstance.ChoiceSettingCollectionValue) > 0 {
				var collectionValues []graphmodels.DeviceManagementConfigurationChoiceSettingValueable
				for _, choiceItem := range detail.SettingInstance.ChoiceSettingCollectionValue {
					choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
					choiceValue.SetValue(&choiceItem.Value)

					// Process children within each choice item
					var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
					for _, child := range choiceItem.Children {
						childInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
						childInstance.SetOdataType(&child.ODataType)
						childInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

						// Handle SimpleSettingValue based on type (string or integer)
						if child.SimpleSettingValue != nil {
							if stringValue, ok := child.SimpleSettingValue.Value.(string); ok {
								simpleValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								simpleValue.SetOdataType(&child.SimpleSettingValue.ODataType)
								simpleValue.SetValue(&stringValue)
								childInstance.SetSimpleSettingValue(simpleValue)
							} else if intValue, ok := child.SimpleSettingValue.Value.(float64); ok {
								intVal := int32(intValue)
								intValueSetting := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
								intValueSetting.SetOdataType(&child.SimpleSettingValue.ODataType)
								intValueSetting.SetValue(&intVal)
								childInstance.SetSimpleSettingValue(intValueSetting)
							}
						}
						children = append(children, childInstance)
					}
					choiceValue.SetChildren(children)
					collectionValues = append(collectionValues, choiceValue)
				}
				instance.SetChoiceSettingCollectionValue(collectionValues)
			}

			baseSetting.SetSettingInstance(instance)

		case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
			instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			// Process each group setting in GroupSettingCollectionValue
			if len(detail.SettingInstance.GroupSettingCollectionValue) > 0 {
				var groupValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable
				for _, groupItem := range detail.SettingInstance.GroupSettingCollectionValue {
					groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
					groupValue.SetSettingValueTemplateReference(groupItem.SettingValueTemplateReference)

					// Process children within each group item
					var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
					for _, child := range groupItem.Children {
						switch child.ODataType {
						case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
							choiceInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
							choiceInstance.SetOdataType(&child.ODataType)
							choiceInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							if child.ChoiceSettingValue != nil {
								choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
								choiceValue.SetValue(&child.ChoiceSettingValue.Value)

								// Process children within ChoiceSettingValue if present
								var choiceChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
								for _, choiceChild := range child.ChoiceSettingValue.Children {
									choiceChildInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
									choiceChildInstance.SetOdataType(&choiceChild.ODataType)
									choiceChildInstance.SetSettingDefinitionId(&choiceChild.SettingDefinitionId)
									choiceChildren = append(choiceChildren, choiceChildInstance)
								}
								choiceValue.SetChildren(choiceChildren)
								choiceInstance.SetChoiceSettingValue(choiceValue)
							}
							children = append(children, choiceInstance)
						}
					}
					groupValue.SetChildren(children)
					groupValues = append(groupValues, groupValue)
				}
				instance.SetGroupSettingCollectionValue(groupValues)
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
