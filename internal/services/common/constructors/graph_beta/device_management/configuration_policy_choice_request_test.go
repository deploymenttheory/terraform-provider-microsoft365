package sharedConstructors

import (
	"encoding/json"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestChoiceSettingRequestChildren(t *testing.T) {
	for _, tc := range []struct {
		name         string
		setting      models.DeviceManagementConfigurationSettingable
		wantChildren int
	}{
		{"choice", ConstructChoiceSettingInstance("definition", "definition_1", "instance-template", "value-template"), -1},
		{"boolean_true", ConstructBoolChoiceSettingInstance("definition", true, "instance-template", "value-template"), -1},
		{"boolean_false", ConstructBoolChoiceSettingInstance("definition", false, "instance-template", "value-template"), -1},
		{"nil_children", ConstructChoiceSettingWithChildren("definition", "definition_1", "instance-template", "value-template", nil), -1},
		{"empty_children", ConstructChoiceSettingWithChildren("definition", "definition_1", "instance-template", "value-template", []models.DeviceManagementConfigurationSettingInstanceable{}), 0},
		{"nested_child", ConstructChoiceSettingWithChildren("definition", "definition_1", "instance-template", "value-template", []models.DeviceManagementConfigurationSettingInstanceable{
			ConstructBoolChoiceSettingInstance("child", true, "", "").GetSettingInstance(),
		}), 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			factory := store.NewBackingStoreSerializationWriterProxyFactory(jsonserialization.NewJsonSerializationWriterFactory())
			writer, err := factory.GetSerializationWriter("application/json")
			if err != nil {
				t.Fatal(err)
			}
			defer writer.Close()
			if err := writer.WriteObjectValue("", tc.setting); err != nil {
				t.Fatal(err)
			}
			content, err := writer.GetSerializedContent()
			if err != nil {
				t.Fatal(err)
			}
			var payload struct {
				SettingInstance struct {
					ChoiceSettingValue struct {
						Children json.RawMessage `json:"children"`
					} `json:"choiceSettingValue"`
				} `json:"settingInstance"`
			}
			if err := json.Unmarshal(content, &payload); err != nil {
				t.Fatal(err)
			}
			raw := payload.SettingInstance.ChoiceSettingValue.Children
			if tc.wantChildren == -1 {
				if raw != nil {
					t.Fatalf("childless choice must omit children; got %s", raw)
				}
				return
			}
			var children []json.RawMessage
			if err := json.Unmarshal(raw, &children); err != nil {
				t.Fatal(err)
			}
			if children == nil || len(children) != tc.wantChildren {
				t.Fatalf("children = %s; want an array with %d elements", raw, tc.wantChildren)
			}
		})
	}
}
