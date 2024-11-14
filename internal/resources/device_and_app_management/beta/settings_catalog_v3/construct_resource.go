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
						ODataType           string `json:"@odata.type"`
						SettingDefinitionId string `json:"settingDefinitionId"`

						// For SimpleSettingCollectionValue within Choice children
						SimpleSettingCollectionValue []struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         string                                                                     `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingCollectionValue,omitempty"`

						// For GroupSettingCollectionValue within Choice children
						GroupSettingCollectionValue []struct {
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
							Children                      []struct {
								SimpleSettingValue *struct {
									ODataType                     string                                                                     `json:"@odata.type"`
									Value                         interface{}                                                                `json:"value"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"simpleSettingValue,omitempty"`
								ODataType                        string                                                                     `json:"@odata.type"`
								SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
								SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
							} `json:"children"`
						} `json:"groupSettingCollectionValue,omitempty"`

						// For simple settings within choice children
						SimpleSettingValue *struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         interface{}                                                                `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingValue,omitempty"`

						// For nested choice settings within choice children
						ChoiceSettingValue *struct {
							Value                         string                                                                     `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
							Children                      []struct {
								ODataType           string `json:"@odata.type"`
								SettingDefinitionId string `json:"settingDefinitionId"`
							} `json:"children"`
						} `json:"choiceSettingValue,omitempty"`

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

						// For nested simple settings within choice setting collection
						SimpleSettingValue *struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         interface{}                                                                `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingValue,omitempty"`

						// For nested simple setting collection within choice setting collection
						SimpleSettingCollectionValue []struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         string                                                                     `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingCollectionValue,omitempty"`

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

						// For nested group setting collections within group setting collection
						GroupSettingCollectionValue []struct {
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
							Children                      []struct {
								ODataType                        string                                                                     `json:"@odata.type"`
								SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
								SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`

								// For nested simple settings within group setting collection within group setting collection
								SimpleSettingValue *struct {
									ODataType                     string                                                                     `json:"@odata.type"`
									Value                         interface{}                                                                `json:"value"`
									ValueState                    string                                                                     `json:"valueState,omitempty"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"simpleSettingValue,omitempty"`

								// For nested simple setting collections within group setting collection within group setting collection
								SimpleSettingCollectionValue []struct {
									ODataType                     string                                                                     `json:"@odata.type"`
									Value                         string                                                                     `json:"value"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"simpleSettingCollectionValue,omitempty"`

								// For nested choice settings within group setting collection within group setting collection
								ChoiceSettingValue *struct {
									Value    string `json:"value"`
									Children []struct {
										ODataType           string `json:"@odata.type"`
										SettingDefinitionId string `json:"settingDefinitionId"`
									} `json:"children"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"choiceSettingValue,omitempty"`
							} `json:"children"`
						} `json:"groupSettingCollectionValue,omitempty"`

						// For nested simple settings (string, integer, secret) within group setting collection
						SimpleSettingValue *struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         interface{}                                                                `json:"value"`
							ValueState                    string                                                                     `json:"valueState,omitempty"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingValue,omitempty"`

						// For nested choice settings within group setting collection
						ChoiceSettingValue *struct {
							Value    string `json:"value"`
							Children []struct {
								ODataType           string `json:"@odata.type"`
								SettingDefinitionId string `json:"settingDefinitionId"`
								SimpleSettingValue  *struct {
									ODataType                     string                                                                     `json:"@odata.type"`
									Value                         interface{}                                                                `json:"value"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"simpleSettingValue,omitempty"`
								ChoiceSettingValue *struct {
									Value    string `json:"value"`
									Children []struct {
										ODataType           string `json:"@odata.type"`
										SettingDefinitionId string `json:"settingDefinitionId"`
									} `json:"children"`
									SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
								} `json:"choiceSettingValue,omitempty"`
							} `json:"children"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"choiceSettingValue,omitempty"`

						// For nested simple setting collections within group setting collection
						SimpleSettingCollectionValue []struct {
							ODataType                     string                                                                     `json:"@odata.type"`
							Value                         string                                                                     `json:"value"`
							SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
						} `json:"simpleSettingCollectionValue,omitempty"`

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

				var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
				for _, child := range detail.SettingInstance.ChoiceSettingValue.Children {
					switch child.ODataType {
					case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
						// Handle SimpleSettingInstance within ChoiceSettingValue
						simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
						simpleInstance.SetOdataType(&child.ODataType)
						simpleInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

						if child.SimpleSettingValue != nil {
							switch child.SimpleSettingValue.ODataType {
							case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
								stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								stringValue.SetOdataType(&child.SimpleSettingValue.ODataType)

								if strValue, ok := child.SimpleSettingValue.Value.(string); ok {
									stringValue.SetValue(&strValue)
								}
								simpleInstance.SetSimpleSettingValue(stringValue)
							case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
								intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
								intValue.SetOdataType(&child.SimpleSettingValue.ODataType)

								if numValue, ok := child.SimpleSettingValue.Value.(float64); ok {
									int32Value := int32(numValue)
									intValue.SetValue(&int32Value)
								}
								simpleInstance.SetSimpleSettingValue(intValue)
							}
						}
						children = append(children, simpleInstance)

					case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
						// Handle nested ChoiceSettingInstance within ChoiceSettingValue
						nestedChoiceInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
						nestedChoiceInstance.SetOdataType(&child.ODataType)
						nestedChoiceInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

						if child.ChoiceSettingValue != nil {
							nestedChoiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
							nestedChoiceValue.SetValue(&child.ChoiceSettingValue.Value)

							// Process nested children within the nested ChoiceSettingValue
							var nestedChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
							for _, nestedChild := range child.ChoiceSettingValue.Children {
								switch nestedChild.ODataType {
								case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
									nestedChoiceInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
									nestedChoiceInstance.SetOdataType(&nestedChild.ODataType)
									nestedChoiceInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)
									nestedChildren = append(nestedChildren, nestedChoiceInstance)

								case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
									nestedSimpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
									nestedSimpleInstance.SetOdataType(&nestedChild.ODataType)
									nestedSimpleInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)
									nestedChildren = append(nestedChildren, nestedSimpleInstance)

								// Handle other types if necessary
								default:
									tflog.Warn(ctx, "Unhandled @odata.type for nested child", map[string]interface{}{
										"odata_type": nestedChild.ODataType,
									})
								}
							}

							nestedChoiceValue.SetChildren(nestedChildren)
							nestedChoiceInstance.SetChoiceSettingValue(nestedChoiceValue)
						}
						children = append(children, nestedChoiceInstance)

						// Handling for GroupSettingCollection within Choice
					case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
						groupCollectionInstance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
						groupCollectionInstance.SetOdataType(&child.ODataType)
						groupCollectionInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

						// Handle group collection value
						if len(child.GroupSettingCollectionValue) > 0 {
							var groupValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable
							for _, groupItem := range child.GroupSettingCollectionValue {
								groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()

								// Process children of each group item (key-value pairs)
								var groupChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
								for _, groupChild := range groupItem.Children {
									if groupChild.ODataType == "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance" {
										simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
										simpleInstance.SetOdataType(&groupChild.ODataType)
										simpleInstance.SetSettingDefinitionId(&groupChild.SettingDefinitionId)

										if groupChild.SimpleSettingValue != nil {
											switch groupChild.SimpleSettingValue.ODataType {
											case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
												stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
												stringValue.SetOdataType(&groupChild.SimpleSettingValue.ODataType)
												if strValue, ok := groupChild.SimpleSettingValue.Value.(string); ok {
													stringValue.SetValue(&strValue)
												}
												simpleInstance.SetSimpleSettingValue(stringValue)
											case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
												intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
												intValue.SetOdataType(&groupChild.SimpleSettingValue.ODataType)
												if numValue, ok := groupChild.SimpleSettingValue.Value.(float64); ok {
													int32Value := int32(numValue)
													intValue.SetValue(&int32Value)
												}
												simpleInstance.SetSimpleSettingValue(intValue)
											}
										}
										groupChildren = append(groupChildren, simpleInstance)
									}
								}

								groupValue.SetChildren(groupChildren)
								if groupItem.SettingValueTemplateReference != nil {
									groupValue.SetSettingValueTemplateReference(groupItem.SettingValueTemplateReference)
								}
								groupValues = append(groupValues, groupValue)
							}
							groupCollectionInstance.SetGroupSettingCollectionValue(groupValues)
						}

						children = append(children, groupCollectionInstance)

						// For SimpleSettingCollection within Choice
					case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
						simpleCollectionInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
						simpleCollectionInstance.SetOdataType(&child.ODataType)
						simpleCollectionInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

						if len(child.SimpleSettingCollectionValue) > 0 {
							var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
							for _, v := range child.SimpleSettingCollectionValue {
								stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
								stringValue.SetOdataType(&v.ODataType)
								stringValue.SetValue(&v.Value)
								values = append(values, stringValue)
							}
							simpleCollectionInstance.SetSimpleSettingCollectionValue(values)
						}

						children = append(children, simpleCollectionInstance)
					}
				}

				choiceValue.SetChildren(children)
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
					value := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
					value.SetOdataType(&detail.SettingInstance.SimpleSettingValue.ODataType)
					if stringValue, ok := detail.SettingInstance.SimpleSettingValue.Value.(string); ok {
						value.SetValue(&stringValue)
					} else {
						tflog.Error(ctx, "Expected string value but got different type", map[string]interface{}{
							"value": detail.SettingInstance.SimpleSettingValue.Value,
						})
					}
					instance.SetSimpleSettingValue(value)

				case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
					value := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
					value.SetOdataType(&detail.SettingInstance.SimpleSettingValue.ODataType)
					if intValue, ok := detail.SettingInstance.SimpleSettingValue.Value.(float64); ok {
						intVal := int32(intValue)
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

			// Handling for GroupSettingCollection
		case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
			instance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
			instance.SetOdataType(&detail.SettingInstance.ODataType)
			instance.SetSettingDefinitionId(&detail.SettingInstance.SettingDefinitionId)

			if len(detail.SettingInstance.GroupSettingCollectionValue) > 0 {
				var groupValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable

				for _, groupItem := range detail.SettingInstance.GroupSettingCollectionValue {
					groupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
					groupOdataType := "#microsoft.graph.deviceManagementConfigurationGroupSettingValue"
					groupValue.SetOdataType(&groupOdataType)

					if groupItem.SettingValueTemplateReference != nil {
						groupValue.SetSettingValueTemplateReference(groupItem.SettingValueTemplateReference)
					}

					var children []graphmodels.DeviceManagementConfigurationSettingInstanceable
					for _, child := range groupItem.Children {
						switch child.ODataType {

						// For nested group setting collections within group setting collection
						case "#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance":
							nestedGroupInstance := graphmodels.NewDeviceManagementConfigurationGroupSettingCollectionInstance()
							nestedGroupInstance.SetOdataType(&child.ODataType)
							nestedGroupInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							if len(child.GroupSettingCollectionValue) > 0 {
								var nestedGroupValues []graphmodels.DeviceManagementConfigurationGroupSettingValueable
								for _, nestedGroupItem := range child.GroupSettingCollectionValue {
									nestedGroupValue := graphmodels.NewDeviceManagementConfigurationGroupSettingValue()
									nestedGroupOdataType := "#microsoft.graph.deviceManagementConfigurationGroupSettingValue"
									nestedGroupValue.SetOdataType(&nestedGroupOdataType)

									var nestedChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
									for _, nestedChild := range nestedGroupItem.Children {
										switch nestedChild.ODataType {
										case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
											simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
											simpleInstance.SetOdataType(&nestedChild.ODataType)
											simpleInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)

											if nestedChild.SimpleSettingValue != nil {
												switch nestedChild.SimpleSettingValue.ODataType {

												case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
													stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
													stringValue.SetOdataType(&nestedChild.SimpleSettingValue.ODataType)
													if strValue, ok := nestedChild.SimpleSettingValue.Value.(string); ok {
														stringValue.SetValue(&strValue)
													}
													simpleInstance.SetSimpleSettingValue(stringValue)

												case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
													intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
													intValue.SetOdataType(&nestedChild.SimpleSettingValue.ODataType)
													if numValue, ok := nestedChild.SimpleSettingValue.Value.(float64); ok {
														int32Value := int32(numValue)
														intValue.SetValue(&int32Value)
													}
													simpleInstance.SetSimpleSettingValue(intValue)

												case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue":
													secretValue := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
													secretValue.SetOdataType(&nestedChild.SimpleSettingValue.ODataType)
													if strValue, ok := nestedChild.SimpleSettingValue.Value.(string); ok {
														secretValue.SetValue(&strValue)
														if nestedChild.SimpleSettingValue.ValueState != "" {
															valueState, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(nestedChild.SimpleSettingValue.ValueState)
															if err == nil {
																secretValue.SetValueState(valueState.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
															}
														}
													}
													simpleInstance.SetSimpleSettingValue(secretValue)
												}
											}
											nestedChildren = append(nestedChildren, simpleInstance)

										case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
											simpleCollectionInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
											simpleCollectionInstance.SetOdataType(&nestedChild.ODataType)
											simpleCollectionInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)

											if len(nestedChild.SimpleSettingCollectionValue) > 0 {
												var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
												for _, v := range nestedChild.SimpleSettingCollectionValue {
													stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
													stringValue.SetOdataType(&v.ODataType)
													stringValue.SetValue(&v.Value)
													values = append(values, stringValue)
												}
												simpleCollectionInstance.SetSimpleSettingCollectionValue(values)
											}
											nestedChildren = append(nestedChildren, simpleCollectionInstance)

										case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
											choiceInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
											choiceInstance.SetOdataType(&nestedChild.ODataType)
											choiceInstance.SetSettingDefinitionId(&nestedChild.SettingDefinitionId)

											if nestedChild.ChoiceSettingValue != nil {
												choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
												choiceOdataType := "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
												choiceValue.SetOdataType(&choiceOdataType)
												choiceValue.SetValue(&nestedChild.ChoiceSettingValue.Value)

												// Always include empty children array for choice settings
												choiceValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{})

												choiceInstance.SetChoiceSettingValue(choiceValue)
											}
											nestedChildren = append(nestedChildren, choiceInstance)
										}
									}
									nestedGroupValue.SetChildren(nestedChildren)
									nestedGroupValues = append(nestedGroupValues, nestedGroupValue)
								}
								nestedGroupInstance.SetGroupSettingCollectionValue(nestedGroupValues)
							}
							children = append(children, nestedGroupInstance)

							// For nested simple setting collections within group setting collection
						case "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance":
							simpleCollectionInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingCollectionInstance()
							simpleCollectionInstance.SetOdataType(&child.ODataType)
							simpleCollectionInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							// Handle simpleSettingCollectionValue
							if len(child.SimpleSettingCollectionValue) > 0 {
								var values []graphmodels.DeviceManagementConfigurationSimpleSettingValueable
								for _, valueItem := range child.SimpleSettingCollectionValue {
									stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
									stringValue.SetOdataType(&valueItem.ODataType)
									stringValue.SetValue(&valueItem.Value)
									values = append(values, stringValue)
								}
								simpleCollectionInstance.SetSimpleSettingCollectionValue(values)
							}

							children = append(children, simpleCollectionInstance)

							// For nested simple settings within group setting collection
						case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
							simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
							simpleInstance.SetOdataType(&child.ODataType)
							simpleInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							if child.SimpleSettingValue != nil {
								switch child.SimpleSettingValue.ODataType {
								case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
									stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
									stringValue.SetOdataType(&child.SimpleSettingValue.ODataType)
									if strValue, ok := child.SimpleSettingValue.Value.(string); ok {
										stringValue.SetValue(&strValue)
									}
									simpleInstance.SetSimpleSettingValue(stringValue)

								case "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue":
									intValue := graphmodels.NewDeviceManagementConfigurationIntegerSettingValue()
									intValue.SetOdataType(&child.SimpleSettingValue.ODataType)
									if numValue, ok := child.SimpleSettingValue.Value.(float64); ok {
										int32Value := int32(numValue)
										intValue.SetValue(&int32Value)
									}
									simpleInstance.SetSimpleSettingValue(intValue)

								case "#microsoft.graph.deviceManagementConfigurationSecretSettingValue":
									secretValue := graphmodels.NewDeviceManagementConfigurationSecretSettingValue()
									secretValue.SetOdataType(&child.SimpleSettingValue.ODataType)
									if strValue, ok := child.SimpleSettingValue.Value.(string); ok {
										secretValue.SetValue(&strValue)
										if child.SimpleSettingValue.ValueState != "" {
											valueState, err := graphmodels.ParseDeviceManagementConfigurationSecretSettingValueState(child.SimpleSettingValue.ValueState)
											if err == nil {
												secretValue.SetValueState(valueState.(*graphmodels.DeviceManagementConfigurationSecretSettingValueState))
											}
										}
									}
									simpleInstance.SetSimpleSettingValue(secretValue)
								}
							}
							children = append(children, simpleInstance)

							// For nested choice settings within group setting collection
						case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
							choiceInstance := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
							choiceInstance.SetOdataType(&child.ODataType)
							choiceInstance.SetSettingDefinitionId(&child.SettingDefinitionId)

							if child.ChoiceSettingValue != nil {
								choiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
								choiceValue.SetValue(&child.ChoiceSettingValue.Value)

								var choiceChildren []graphmodels.DeviceManagementConfigurationSettingInstanceable
								for _, choiceChild := range child.ChoiceSettingValue.Children {
									switch choiceChild.ODataType {
									case "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance":
										nestedChoice := graphmodels.NewDeviceManagementConfigurationChoiceSettingInstance()
										nestedChoice.SetOdataType(&choiceChild.ODataType)
										nestedChoice.SetSettingDefinitionId(&choiceChild.SettingDefinitionId)

										if choiceChild.ChoiceSettingValue != nil {
											nestedChoiceValue := graphmodels.NewDeviceManagementConfigurationChoiceSettingValue()
											nestedChoiceValue.SetValue(&choiceChild.ChoiceSettingValue.Value)
											nestedChoiceValue.SetChildren([]graphmodels.DeviceManagementConfigurationSettingInstanceable{})
											nestedChoice.SetChoiceSettingValue(nestedChoiceValue)
										}
										choiceChildren = append(choiceChildren, nestedChoice)

									case "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance":
										simpleInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
										simpleInstance.SetOdataType(&choiceChild.ODataType)
										simpleInstance.SetSettingDefinitionId(&choiceChild.SettingDefinitionId)

										if choiceChild.SimpleSettingValue != nil {
											switch choiceChild.SimpleSettingValue.ODataType {
											case "#microsoft.graph.deviceManagementConfigurationStringSettingValue":
												stringValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
												stringValue.SetOdataType(&choiceChild.SimpleSettingValue.ODataType)
												if strValue, ok := choiceChild.SimpleSettingValue.Value.(string); ok {
													stringValue.SetValue(&strValue)
												}
												simpleInstance.SetSimpleSettingValue(stringValue)
											}
										}
										choiceChildren = append(choiceChildren, simpleInstance)
									}
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
