package sharedmodels

import (
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// DeviceConfigV2GraphServiceModel is the root configuration model
type DeviceConfigV2GraphServiceModel struct {
	SettingsDetails []SettingDetail `json:"settingsDetails"`
}

// SettingDetail represents a single setting detail
type SettingDetail struct {
	ID              string          `json:"id"`
	SettingInstance SettingInstance `json:"settingInstance"`
}

// SettingInstance contains the core setting configuration
type SettingInstance struct {
	ODataType                        string                                                                     `json:"@odata.type"`
	SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
	SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
	SimpleSettingValue               *SimpleSettingStruct                                                       `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct                                            `json:"simpleSettingCollectionValue,omitempty"`
	ChoiceSettingValue               *ChoiceSettingStruct                                                       `json:"choiceSettingValue,omitempty"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct                                            `json:"choiceSettingCollectionValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct                                             `json:"groupSettingCollectionValue,omitempty"`
}

// SimpleSettingStruct represents a simple setting value
type SimpleSettingStruct struct {
	ODataType                     string                                                                     `json:"@odata.type"`
	SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
	Value                         interface{}                                                                `json:"value"`
	ValueState                    string                                                                     `json:"valueState,omitempty"`
}

// SimpleSettingCollectionStruct represents a collection of simple settings
type SimpleSettingCollectionStruct struct {
	ODataType                     string                                                                     `json:"@odata.type"`
	SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
	Value                         string                                                                     `json:"value"`
}

// ChoiceSettingChild represents a child element in a choice setting
type ChoiceSettingChild struct {
	ODataType                        string                                                                     `json:"@odata.type"`
	SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
	SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
	ChoiceSettingValue               *ChoiceSettingStruct                                                       `json:"choiceSettingValue,omitempty"`
	SimpleSettingValue               *SimpleSettingStruct                                                       `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct                                            `json:"simpleSettingCollectionValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct                                             `json:"groupSettingCollectionValue,omitempty"`
}

// ChoiceSettingStruct represents a choice setting
type ChoiceSettingStruct struct {
	SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
	Value                         string                                                                     `json:"value"`
	Children                      []ChoiceSettingChild                                                       `json:"children"`
}

// ChoiceSettingCollectionChild represents a child element in a choice setting collection
type ChoiceSettingCollectionChild struct {
	ODataType                        string                                                                     `json:"@odata.type"`
	SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
	SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
	SimpleSettingValue               *SimpleSettingStruct                                                       `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct                                            `json:"simpleSettingCollectionValue,omitempty"`
}

// ChoiceSettingCollectionStruct represents a collection of choice settings
type ChoiceSettingCollectionStruct struct {
	SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
	Value                         string                                                                     `json:"value"`
	Children                      []ChoiceSettingCollectionChild                                             `json:"children"`
}

// GroupSettingCollectionChild represents a child element in a group setting collection
type GroupSettingCollectionChild struct {
	ODataType                        string                                                                     `json:"@odata.type"`
	SettingDefinitionId              string                                                                     `json:"settingDefinitionId"`
	SettingInstanceTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingInstanceTemplateReference"`
	SimpleSettingValue               *SimpleSettingStruct                                                       `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct                                            `json:"simpleSettingCollectionValue,omitempty"`
	ChoiceSettingValue               *ChoiceSettingStruct                                                       `json:"choiceSettingValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct                                             `json:"groupSettingCollectionValue,omitempty"`
}

// GroupSettingCollectionStruct represents a collection of group settings
type GroupSettingCollectionStruct struct {
	SettingValueTemplateReference graphmodels.DeviceManagementConfigurationSettingValueTemplateReferenceable `json:"settingValueTemplateReference"`
	Children                      []GroupSettingCollectionChild                                              `json:"children"`
}
