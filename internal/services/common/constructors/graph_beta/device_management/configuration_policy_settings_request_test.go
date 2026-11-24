package sharedConstructors

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// serializeSettingsCatalogSetting writes setting with the backing-store proxy factory that the Graph request adapter
// uses, so a children slice that was set to nil is written as "children": null, as it is on the wire.
func serializeSettingsCatalogSetting(t *testing.T, setting graphmodels.DeviceManagementConfigurationSettingable) map[string]any {
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

// requireEmptyChildren fails unless value[key]["children"] is present and is an empty JSON array.
func requireEmptyChildren(t *testing.T, value any, path string) {
	t.Helper()
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s is %T, want an object", path, value)
	}
	children, present := object["children"]
	if !present {
		t.Fatalf("%s.children is missing, want []", path)
	}
	list, ok := children.([]any)
	if !ok || list == nil || len(list) != 0 {
		t.Fatalf("%s.children = %#v, want []", path, children)
	}
}

func TestConstructSettingsCatalogSettings_EmptyChildrenSerializeAsEmptyArray(t *testing.T) {
	for _, tc := range []struct {
		name     string
		settings string
		assert   func(t *testing.T, instance map[string]any)
	}{
		{
			name:     "choice_collection_without_children",
			settings: `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance","settingDefinitionId":"device_vendor_msft_defender_configuration_passiveremediation","choiceSettingCollectionValue":[{"value":"device_vendor_msft_defender_configuration_passiveremediation_1","children":[]}]}}]}`,
			assert: func(t *testing.T, instance map[string]any) {
				requireEmptyChildren(t, instance["choiceSettingCollectionValue"].([]any)[0], "choiceSettingCollectionValue[0]")
			},
		},
		{
			name:     "group_collection_without_children",
			settings: `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","settingDefinitionId":"group_definition","groupSettingCollectionValue":[{"children":[]}]}}]}`,
			assert: func(t *testing.T, instance map[string]any) {
				requireEmptyChildren(t, instance["groupSettingCollectionValue"].([]any)[0], "groupSettingCollectionValue[0]")
			},
		},
		{
			name:     "choice_without_children",
			settings: `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"device_vendor_msft_policy_config_camera_allowcamera","choiceSettingValue":{"value":"device_vendor_msft_policy_config_camera_allowcamera_1","children":[]}}}]}`,
			assert: func(t *testing.T, instance map[string]any) {
				requireEmptyChildren(t, instance["choiceSettingValue"], "choiceSettingValue")
			},
		},
		{
			name:     "nested_choice_without_children_in_group",
			settings: `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationGroupSettingCollectionInstance","settingDefinitionId":"group_definition","groupSettingCollectionValue":[{"children":[{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"group_definition_choice","choiceSettingValue":{"value":"group_definition_choice_1","children":[]}},{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance","settingDefinitionId":"group_definition_choicecollection","choiceSettingCollectionValue":[{"value":"group_definition_choicecollection_1","children":[]}]}]}]}}]}`,
			assert: func(t *testing.T, instance map[string]any) {
				children := instance["groupSettingCollectionValue"].([]any)[0].(map[string]any)["children"].([]any)
				requireEmptyChildren(t, children[0].(map[string]any)["choiceSettingValue"], "groupSettingCollectionValue[0].children[0].choiceSettingValue")
				requireEmptyChildren(t, children[1].(map[string]any)["choiceSettingCollectionValue"].([]any)[0], "groupSettingCollectionValue[0].children[1].choiceSettingCollectionValue[0]")
			},
		},
		{
			name:     "choice_with_only_unsupported_children",
			settings: `{"settings":[{"id":"0","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance","settingDefinitionId":"choice_definition","choiceSettingValue":{"value":"choice_definition_1","children":[{"@odata.type":"#microsoft.graph.unsupportedSettingInstance","settingDefinitionId":"unsupported"}]}}}]}`,
			assert: func(t *testing.T, instance map[string]any) {
				requireEmptyChildren(t, instance["choiceSettingValue"], "choiceSettingValue")
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			settings := ConstructSettingsCatalogSettings(context.Background(), types.StringValue(tc.settings))
			if len(settings) != 1 {
				t.Fatalf("got %d settings, want 1", len(settings))
			}
			payload := serializeSettingsCatalogSetting(t, settings[0])
			instance, ok := payload["settingInstance"].(map[string]any)
			if !ok {
				t.Fatalf("settingInstance missing from %#v", payload)
			}
			tc.assert(t, instance)
		})
	}
}
