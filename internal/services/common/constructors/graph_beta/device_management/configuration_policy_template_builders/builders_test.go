package configurationPolicyTemplateBuilders

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// serializeSetting writes setting with the backing-store proxy factory that the Graph request adapter uses, so a
// children slice that was set to nil is written as "children": null, as it is on the wire.
func serializeSetting(t *testing.T, setting models.DeviceManagementConfigurationSettingable) map[string]any {
	t.Helper()
	factory := store.NewBackingStoreSerializationWriterProxyFactory(jsonserialization.NewJsonSerializationWriterFactory())
	writer, err := factory.GetSerializationWriter("application/json")
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	if err := writer.WriteObjectValue("", setting); err != nil {
		t.Fatal(err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(content, &payload); err != nil {
		t.Fatalf("unmarshal %s: %v", content, err)
	}
	return payload
}

// valueAt follows path through nested JSON objects (string keys) and arrays (int indexes).
func valueAt(t *testing.T, value any, path ...any) any {
	t.Helper()
	for _, step := range path {
		switch key := step.(type) {
		case string:
			object, ok := value.(map[string]any)
			if !ok {
				t.Fatalf("at %v: %T is not an object", path, value)
			}
			value = object[key]
		case int:
			list, ok := value.([]any)
			if !ok || key >= len(list) {
				t.Fatalf("at %v: %#v has no element %d", path, value, key)
			}
			value = list[key]
		}
	}
	return value
}

func TestConstructBuilders_EmptyChildrenSerializeAsEmptyArray(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name      string
		construct func() (models.DeviceManagementConfigurationSettingable, error)
		paths     [][]any
	}{
		{
			name: "simple_choice",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructSimpleChoiceSetting(ctx, "choice_definition", "choice_definition_1", "", "")
			},
			paths: [][]any{{"settingInstance", "choiceSettingValue", "children"}},
		},
		{
			name: "simple_choice_with_template",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructSimpleChoiceSettingWithTemplate(ctx, "choice_definition", "choice_definition_1", "", "", false)
			},
			paths: [][]any{{"settingInstance", "choiceSettingValue", "children"}},
		},
		{
			name: "choice_with_boolean_child",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructChoiceWithBooleanSetting(ctx, "choice_definition", "choice_definition_1", "choice_definition_child", "choice_definition_child_1", "", "", "")
			},
			paths: [][]any{{"settingInstance", "choiceSettingValue", "children", 0, "choiceSettingValue", "children"}},
		},
		{
			name: "choice_collection",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructChoiceCollectionSetting(ctx, "choice_collection_definition", []string{"choice_collection_definition_1"}, "")
			},
			paths: [][]any{{"settingInstance", "choiceSettingCollectionValue", 0, "children"}},
		},
		{
			name: "choice_with_choice_children",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructChoiceWithMultipleChoiceChildren(ctx, "choice_definition", "choice_definition_1", []ChoiceChildConfig{
					{SettingDefinitionId: "choice_definition_child", SettingType: "choice", Value: "choice_definition_child_1"},
					{SettingDefinitionId: "choice_definition_collection", SettingType: "choice_collection", ChoiceValues: []string{"choice_definition_collection_1"}},
				}, "", "", false)
			},
			paths: [][]any{
				{"settingInstance", "choiceSettingValue", "children", 0, "choiceSettingValue", "children"},
				{"settingInstance", "choiceSettingValue", "children", 1, "choiceSettingCollectionValue", 0, "children"},
			},
		},
		{
			name: "group_collection",
			construct: func() (models.DeviceManagementConfigurationSettingable, error) {
				return ConstructGroupCollectionSetting(ctx, "group_definition", [][]GroupChildConfig{
					{},
					{{SettingDefinitionId: "group_definition_choice", SettingType: "choice", Value: "group_definition_choice_1"}},
				}, "")
			},
			paths: [][]any{
				{"settingInstance", "groupSettingCollectionValue", 0, "children"},
				{"settingInstance", "groupSettingCollectionValue", 1, "children", 0, "choiceSettingValue", "children"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setting, err := tc.construct()
			if err != nil {
				t.Fatal(err)
			}
			payload := serializeSetting(t, setting)
			for _, path := range tc.paths {
				children := valueAt(t, payload, path...)
				if list, ok := children.([]any); !ok || list == nil || len(list) != 0 {
					t.Errorf("%v = %#v, want []", path, children)
				}
			}
		})
	}
}
