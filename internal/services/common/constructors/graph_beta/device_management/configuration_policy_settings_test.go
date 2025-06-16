package sharedConstructors

import (
	"context"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestConstructSettingsCatalogSettings_Simple(t *testing.T) {
	ctx := context.Background()
	jsonStr := `{"settings":[{"id":"id1","settingInstance":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance","settingDefinitionId":"defId","simpleSettingValue":{"@odata.type":"#microsoft.graph.deviceManagementConfigurationStringSettingValue","value":"foo"}}}]}`
	settings := ConstructSettingsCatalogSettings(ctx, types.StringValue(jsonStr))
	if len(settings) != 1 {
		t.Fatalf("Expected 1 setting, got %d", len(settings))
	}
	inst := settings[0].GetSettingInstance()
	if _, ok := inst.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable); !ok {
		t.Errorf("Expected simple setting instance, got %T", inst)
	}
}

func TestHandleSimpleValue_String(t *testing.T) {
	ctx := context.Background()
	val := &sharedmodels.SimpleSettingStruct{
		ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
		Value:     "bar",
	}
	res := handleSimpleValue(ctx, val)
	strVal, ok := res.(graphmodels.DeviceManagementConfigurationStringSettingValueable)
	if !ok {
		t.Fatalf("Expected string setting value, got %T", res)
	}
	if strVal.GetValue() == nil || *strVal.GetValue() != "bar" {
		t.Errorf("Value = %v, want 'bar'", strVal.GetValue())
	}
}

func TestHandleSimpleValue_Integer(t *testing.T) {
	ctx := context.Background()
	val := &sharedmodels.SimpleSettingStruct{
		ODataType: "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue",
		Value:     float64(123),
	}
	res := handleSimpleValue(ctx, val)
	intVal, ok := res.(graphmodels.DeviceManagementConfigurationIntegerSettingValueable)
	if !ok {
		t.Fatalf("Expected integer setting value, got %T", res)
	}
	if intVal.GetValue() == nil || *intVal.GetValue() != 123 {
		t.Errorf("Value = %v, want 123", intVal.GetValue())
	}
}

func TestHandleSimpleSettingCollection(t *testing.T) {
	coll := []sharedmodels.SimpleSettingCollectionStruct{{ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue", Value: "a"}, {ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue", Value: "b"}}
	res := handleSimpleSettingCollection(coll)
	if len(res) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(res))
	}
	for i, v := range res {
		if _, ok := v.(graphmodels.DeviceManagementConfigurationStringSettingValueable); !ok {
			t.Errorf("Value %d: expected string setting value, got %T", i, v)
		}
	}
}

func TestHandleChoiceSettingChildren_Empty(t *testing.T) {
	ctx := context.Background()
	res := handleChoiceSettingChildren(ctx, nil)
	if len(res) != 0 {
		t.Errorf("Expected 0 children, got %d", len(res))
	}
}

func TestHandleGroupSettingCollection_Empty(t *testing.T) {
	ctx := context.Background()
	res := handleGroupSettingCollection(ctx, nil)
	if len(res) != 0 {
		t.Errorf("Expected 0 group values, got %d", len(res))
	}
}

func TestHandleChoiceCollectionValue_Empty(t *testing.T) {
	ctx := context.Background()
	res := handleChoiceCollectionValue(ctx, nil)
	if len(res) != 0 {
		t.Errorf("Expected 0 choice values, got %d", len(res))
	}
}

func TestHandleSettingInstance_Simple(t *testing.T) {
	ctx := context.Background()
	inst := sharedmodels.SettingInstance{
		ODataType:           "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
		SettingDefinitionId: "defId",
		SimpleSettingValue: &sharedmodels.SimpleSettingStruct{
			ODataType: "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
			Value:     "foo",
		},
	}
	res := handleSettingInstance(ctx, inst)
	if res == nil {
		t.Fatal("Expected non-nil result")
	}
	if _, ok := res.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable); !ok {
		t.Errorf("Expected simple setting instance, got %T", res)
	}
}

func TestSetInstanceTemplateReference(t *testing.T) {
	inst := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	ref := &sharedmodels.SettingInstanceTemplateReference{SettingInstanceTemplateId: "tmpl"}
	setInstanceTemplateReference(inst, ref)
	if inst.GetSettingInstanceTemplateReference() == nil || *inst.GetSettingInstanceTemplateReference().GetSettingInstanceTemplateId() != "tmpl" {
		t.Errorf("TemplateId = %v, want 'tmpl'", inst.GetSettingInstanceTemplateReference())
	}
}

func TestSetValueTemplateReference_Simple(t *testing.T) {
	val := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
	ref := &sharedmodels.SettingValueTemplateReference{SettingValueTemplateId: "tmpl", UseTemplateDefault: true}
	setValueTemplateReference(val, ref)
	if val.GetSettingValueTemplateReference() == nil || *val.GetSettingValueTemplateReference().GetSettingValueTemplateId() != "tmpl" {
		t.Errorf("TemplateId = %v, want 'tmpl'", val.GetSettingValueTemplateReference())
	}
	if val.GetSettingValueTemplateReference() == nil || val.GetSettingValueTemplateReference().GetUseTemplateDefault() == nil || *val.GetSettingValueTemplateReference().GetUseTemplateDefault() != true {
		t.Errorf("UseTemplateDefault = %v, want true", val.GetSettingValueTemplateReference().GetUseTemplateDefault())
	}
}
