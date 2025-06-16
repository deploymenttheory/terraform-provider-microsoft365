package sharedmodels

// DeviceConfigV2GraphServiceResourceModel is the root settings catalog model
// Officially the api only recognises 'settings' as the root key, but we use 'settings' and 'setting'
// to support both single-instance and array-based settings.
type DeviceConfigV2GraphServiceResourceModel struct {
	Settings []Setting `json:"settings,omitempty"` // For array-based settings
	Setting  *Setting  `json:"setting,omitempty"`  // For single-instance settings
}

// SettingDetail represents a single setting detail
// type SettingDetail struct {
type Setting struct {
	ID              string          `json:"id"`
	SettingInstance SettingInstance `json:"settingInstance"`
}

// SettingInstance contains the core setting configuration
type SettingInstance struct {
	ODataType                        string                            `json:"@odata.type"`
	SettingDefinitionId              string                            `json:"settingDefinitionId"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `json:"settingInstanceTemplateReference"`
	SimpleSettingValue               *SimpleSettingStruct              `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `json:"simpleSettingCollectionValue,omitempty"`
	ChoiceSettingValue               *ChoiceSettingStruct              `json:"choiceSettingValue,omitempty"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `json:"choiceSettingCollectionValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `json:"groupSettingCollectionValue,omitempty"`
}

// SimpleSettingStruct represents a simple setting value
type SimpleSettingStruct struct {
	ODataType                     string                         `json:"@odata.type"`
	SettingValueTemplateReference *SettingValueTemplateReference `json:"settingValueTemplateReference"`
	Value                         interface{}                    `json:"value"`
	ValueState                    string                         `json:"valueState,omitempty"`
}

// SimpleSettingCollectionStruct represents a collection of simple settings
type SimpleSettingCollectionStruct struct {
	ODataType                     string                         `json:"@odata.type"`
	SettingValueTemplateReference *SettingValueTemplateReference `json:"settingValueTemplateReference"`
	Value                         string                         `json:"value"`
}

// ChoiceSettingChild represents a child element in a choice setting
type ChoiceSettingChild struct {
	ODataType                        string                            `json:"@odata.type"`
	SettingDefinitionId              string                            `json:"settingDefinitionId"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `json:"settingInstanceTemplateReference"`
	ChoiceSettingValue               *ChoiceSettingStruct              `json:"choiceSettingValue,omitempty"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `json:"choiceSettingCollectionValue,omitempty"`
	SimpleSettingValue               *SimpleSettingStruct              `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `json:"simpleSettingCollectionValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `json:"groupSettingCollectionValue,omitempty"`
}

// ChoiceSettingStruct represents a choice setting
type ChoiceSettingStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `json:"settingValueTemplateReference"`
	Value                         string                         `json:"value"`
	Children                      []ChoiceSettingChild           `json:"children"`
}

// ChoiceSettingCollectionChild represents a child element in a choice setting collection
type ChoiceSettingCollectionChild struct {
	ODataType                        string                            `json:"@odata.type"`
	SettingDefinitionId              string                            `json:"settingDefinitionId"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `json:"settingInstanceTemplateReference"`
	SimpleSettingValue               *SimpleSettingStruct              `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `json:"simpleSettingCollectionValue,omitempty"`
}

// ChoiceSettingCollectionStruct represents a collection of choice settings
type ChoiceSettingCollectionStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `json:"settingValueTemplateReference"`
	Value                         string                         `json:"value"`
	Children                      []ChoiceSettingCollectionChild `json:"children"`
}

// GroupSettingCollectionChild represents a child element in a group setting collection
type GroupSettingCollectionChild struct {
	ODataType                        string                            `json:"@odata.type"`
	SettingDefinitionId              string                            `json:"settingDefinitionId"`
	SettingInstanceTemplateReference *SettingInstanceTemplateReference `json:"settingInstanceTemplateReference"`
	ChoiceSettingValue               *ChoiceSettingStruct              `json:"choiceSettingValue,omitempty"`
	ChoiceSettingCollectionValue     []ChoiceSettingCollectionStruct   `json:"choiceSettingCollectionValue,omitempty"`
	GroupSettingCollectionValue      []GroupSettingCollectionStruct    `json:"groupSettingCollectionValue,omitempty"`
	SimpleSettingValue               *SimpleSettingStruct              `json:"simpleSettingValue,omitempty"`
	SimpleSettingCollectionValue     []SimpleSettingCollectionStruct   `json:"simpleSettingCollectionValue,omitempty"`
}

// GroupSettingCollectionStruct represents a collection of group settings
type GroupSettingCollectionStruct struct {
	SettingValueTemplateReference *SettingValueTemplateReference `json:"settingValueTemplateReference"`
	Children                      []GroupSettingCollectionChild  `json:"children"`
}

// SettingInstanceTemplateReference represents the template reference at the instance level
type SettingInstanceTemplateReference struct {
	SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
}

// SettingValueTemplateReference represents the template reference at the value level
type SettingValueTemplateReference struct {
	SettingValueTemplateId string `json:"settingValueTemplateId"`
	UseTemplateDefault     bool   `json:"useTemplateDefault"`
}
