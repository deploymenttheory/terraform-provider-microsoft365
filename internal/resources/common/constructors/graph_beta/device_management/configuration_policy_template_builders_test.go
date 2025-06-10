package sharedConstructors

import (
	"testing"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestConstructChoiceSettingInstance(t *testing.T) {
	setting := ConstructChoiceSettingInstance("defId", "val", "instTmplId", "valTmplId")
	inst := setting.GetSettingInstance()
	choiceInst, ok := inst.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		t.Fatalf("Expected choice setting instance, got %T", inst)
	}
	if got := choiceInst.GetSettingDefinitionId(); got == nil || *got != "defId" {
		t.Errorf("SettingDefinitionId = %v, want 'defId'", got)
	}
	valObj := choiceInst.GetChoiceSettingValue()
	if valObj == nil || valObj.GetValue() == nil || *valObj.GetValue() != "val" {
		t.Errorf("Choice value = %v, want 'val'", valObj.GetValue())
	}
	if ref := choiceInst.GetSettingInstanceTemplateReference(); ref == nil || ref.GetSettingInstanceTemplateId() == nil || *ref.GetSettingInstanceTemplateId() != "instTmplId" {
		t.Errorf("SettingInstanceTemplateId = %v, want 'instTmplId'", ref)
	}
	if valObj.GetSettingValueTemplateReference() == nil || valObj.GetSettingValueTemplateReference().GetSettingValueTemplateId() == nil || *valObj.GetSettingValueTemplateReference().GetSettingValueTemplateId() != "valTmplId" {
		t.Errorf("SettingValueTemplateId = %v, want 'valTmplId'", valObj.GetSettingValueTemplateReference())
	}
}

func TestConstructStringSimpleSettingInstance(t *testing.T) {
	setting := ConstructStringSimpleSettingInstance("defId", "strval", "instTmplId", "valTmplId")
	inst := setting.GetSettingInstance()
	strInst, ok := inst.(models.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		t.Fatalf("Expected simple setting instance, got %T", inst)
	}
	if got := strInst.GetSettingDefinitionId(); got == nil || *got != "defId" {
		t.Errorf("SettingDefinitionId = %v, want 'defId'", got)
	}
	valObj := strInst.GetSimpleSettingValue()
	strVal, ok := valObj.(models.DeviceManagementConfigurationStringSettingValueable)
	if !ok {
		t.Fatalf("Expected string setting value, got %T", valObj)
	}
	if strVal.GetValue() == nil || *strVal.GetValue() != "strval" {
		t.Errorf("String value = %v, want 'strval'", strVal.GetValue())
	}
	if ref := strInst.GetSettingInstanceTemplateReference(); ref == nil || ref.GetSettingInstanceTemplateId() == nil || *ref.GetSettingInstanceTemplateId() != "instTmplId" {
		t.Errorf("SettingInstanceTemplateId = %v, want 'instTmplId'", ref)
	}
	if strVal.GetSettingValueTemplateReference() == nil || strVal.GetSettingValueTemplateReference().GetSettingValueTemplateId() == nil || *strVal.GetSettingValueTemplateReference().GetSettingValueTemplateId() != "valTmplId" {
		t.Errorf("SettingValueTemplateId = %v, want 'valTmplId'", strVal.GetSettingValueTemplateReference())
	}
}

func TestConstructIntSimpleSettingInstance(t *testing.T) {
	setting := ConstructIntSimpleSettingInstance("defId", 42, "instTmplId", "valTmplId")
	inst := setting.GetSettingInstance()
	intInst, ok := inst.(models.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		t.Fatalf("Expected simple setting instance, got %T", inst)
	}
	if got := intInst.GetSettingDefinitionId(); got == nil || *got != "defId" {
		t.Errorf("SettingDefinitionId = %v, want 'defId'", got)
	}
	valObj := intInst.GetSimpleSettingValue()
	intVal, ok := valObj.(models.DeviceManagementConfigurationIntegerSettingValueable)
	if !ok {
		t.Fatalf("Expected integer setting value, got %T", valObj)
	}
	if intVal.GetValue() == nil || *intVal.GetValue() != 42 {
		t.Errorf("Int value = %v, want 42", intVal.GetValue())
	}
	if ref := intInst.GetSettingInstanceTemplateReference(); ref == nil || ref.GetSettingInstanceTemplateId() == nil || *ref.GetSettingInstanceTemplateId() != "instTmplId" {
		t.Errorf("SettingInstanceTemplateId = %v, want 'instTmplId'", ref)
	}
	if intVal.GetSettingValueTemplateReference() == nil || intVal.GetSettingValueTemplateReference().GetSettingValueTemplateId() == nil || *intVal.GetSettingValueTemplateReference().GetSettingValueTemplateId() != "valTmplId" {
		t.Errorf("SettingValueTemplateId = %v, want 'valTmplId'", intVal.GetSettingValueTemplateReference())
	}
}

func TestConstructSimpleSettingCollectionInstance(t *testing.T) {
	values := []string{"a", "b"}
	setting := ConstructSimpleSettingCollectionInstance("defId", values, "instTmplId")
	inst := setting.GetSettingInstance()
	collInst, ok := inst.(models.DeviceManagementConfigurationSimpleSettingCollectionInstanceable)
	if !ok {
		t.Fatalf("Expected collection setting instance, got %T", inst)
	}
	if got := collInst.GetSettingDefinitionId(); got == nil || *got != "defId" {
		t.Errorf("SettingDefinitionId = %v, want 'defId'", got)
	}
	collVals := collInst.GetSimpleSettingCollectionValue()
	if len(collVals) != 2 {
		t.Errorf("Expected 2 values, got %d", len(collVals))
	}
	for i, v := range collVals {
		strVal, ok := v.(models.DeviceManagementConfigurationStringSettingValueable)
		if !ok {
			t.Errorf("Value %d: expected string setting value, got %T", i, v)
		}
		if strVal.GetValue() == nil || *strVal.GetValue() != values[i] {
			t.Errorf("Value %d: got %v, want %v", i, strVal.GetValue(), values[i])
		}
	}
	if ref := collInst.GetSettingInstanceTemplateReference(); ref == nil || ref.GetSettingInstanceTemplateId() == nil || *ref.GetSettingInstanceTemplateId() != "instTmplId" {
		t.Errorf("SettingInstanceTemplateId = %v, want 'instTmplId'", ref)
	}
}

func TestConstructBoolChoiceSettingInstance(t *testing.T) {
	settingTrue := ConstructBoolChoiceSettingInstance("defId", true, "instTmplId", "valTmplId")
	instTrue := settingTrue.GetSettingInstance()
	choiceInstTrue, ok := instTrue.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		t.Fatalf("Expected choice setting instance, got %T", instTrue)
	}
	valObjTrue := choiceInstTrue.GetChoiceSettingValue()
	if valObjTrue == nil || valObjTrue.GetValue() == nil || *valObjTrue.GetValue() != "defId_1" {
		t.Errorf("True value = %v, want 'defId_1'", valObjTrue.GetValue())
	}

	settingFalse := ConstructBoolChoiceSettingInstance("defId", false, "instTmplId", "valTmplId")
	instFalse := settingFalse.GetSettingInstance()
	choiceInstFalse, ok := instFalse.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		t.Fatalf("Expected choice setting instance, got %T", instFalse)
	}
	valObjFalse := choiceInstFalse.GetChoiceSettingValue()
	if valObjFalse == nil || valObjFalse.GetValue() == nil || *valObjFalse.GetValue() != "defId_0" {
		t.Errorf("False value = %v, want 'defId_0'", valObjFalse.GetValue())
	}
}
