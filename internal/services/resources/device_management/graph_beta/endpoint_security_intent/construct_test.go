package graphBetaEndpointSecurityIntent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func serializeBody(t *testing.T, body serialization.Parsable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize body: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized body: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("parse serialized body %s: %v", content, err)
	}
	return result
}

func settingsSet(values ...[3]string) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.ObjectValueMust(settingAttrTypes, map[string]attr.Value{
			"definition_id": types.StringValue(value[0]),
			"setting_type":  types.StringValue(value[1]),
			"value_json":    types.StringValue(value[2]),
		}))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: settingAttrTypes}, elements)
}

// Each setting is sent with its concrete @odata.type, because Graph's updateSettings fails with a 500 for the abstract base type.
func TestUnitConstructResourceSendsTypedSettingsDelta(t *testing.T) {
	data := &EndpointSecurityIntentResourceModel{
		DisplayName:     types.StringValue("intent"),
		Description:     types.StringValue(""),
		RoleScopeTagIds: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")}),
		Settings:        settingsSet([3]string{"def-bool", settingTypeBoolean, "true"}, [3]string{"def-string", settingTypeString, `"on"`}),
		Assignments:     types.SetNull(types.ObjectType{AttrTypes: assignmentAttrTypes}),
	}

	body, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource: %v", err)
	}

	result := serializeBody(t, body)
	delta, ok := result["settingsDelta"].([]any)
	if !ok || len(delta) != 2 {
		t.Fatalf("expected two settings in settingsDelta, got %v", result["settingsDelta"])
	}

	settingTypes := map[string]string{}
	for _, raw := range delta {
		setting := raw.(map[string]any)
		settingTypes[setting["definitionId"].(string)] = setting["@odata.type"].(string)
		if _, hasValue := setting["value"]; hasValue {
			t.Errorf("setting %v must only carry valueJson", setting)
		}
	}
	if settingTypes["def-bool"] != "#microsoft.graph.deviceManagementBooleanSettingInstance" || settingTypes["def-string"] != "#microsoft.graph.deviceManagementStringSettingInstance" {
		t.Errorf("unexpected setting types: %v", settingTypes)
	}
}

// Unset settings are omitted from the create body rather than sent as null or an empty list.
func TestUnitConstructResourceOmitsUnsetSettings(t *testing.T) {
	data := &EndpointSecurityIntentResourceModel{
		DisplayName:     types.StringValue("intent"),
		Description:     types.StringValue(""),
		RoleScopeTagIds: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")}),
		Settings:        types.SetNull(types.ObjectType{AttrTypes: settingAttrTypes}),
	}

	body, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource: %v", err)
	}

	if _, present := serializeBody(t, body)["settingsDelta"]; present {
		t.Errorf("settingsDelta must be omitted when no settings are configured")
	}
}

// An empty assignment list is still sent, because that is how assign removes every assignment.
func TestUnitConstructAssignmentsSendsEmptyList(t *testing.T) {
	data := &EndpointSecurityIntentResourceModel{Assignments: types.SetNull(types.ObjectType{AttrTypes: assignmentAttrTypes})}

	body, err := constructAssignments(context.Background(), data)
	if err != nil {
		t.Fatalf("constructAssignments: %v", err)
	}

	assignments, ok := serializeBody(t, body)["assignments"].([]any)
	if !ok || len(assignments) != 0 {
		t.Errorf("expected an empty assignments list, got %v", assignments)
	}
}

func TestUnitJsonEqual(t *testing.T) {
	cases := []struct {
		a, b  string
		equal bool
	}{
		{"true", " true ", true},
		{`{"a":1,"b":[1,2]}`, `{ "b": [1, 2], "a": 1 }`, true},
		{`"on"`, `"off"`, false},
		{"true", "not json", false},
	}
	for _, c := range cases {
		if got := jsonEqual(c.a, c.b); got != c.equal {
			t.Errorf("jsonEqual(%q, %q) = %v, want %v", c.a, c.b, got, c.equal)
		}
	}
}

func TestUnitParseImportId(t *testing.T) {
	intentId, definitionIds, err := parseImportId("44357f75-2b25-49af-a7ae-f8dbaba1659d/def-a, def-b,def-a")
	if err != nil {
		t.Fatalf("parseImportId: %v", err)
	}
	if intentId != "44357f75-2b25-49af-a7ae-f8dbaba1659d" || len(definitionIds) != 2 || definitionIds[0] != "def-a" || definitionIds[1] != "def-b" {
		t.Errorf("unexpected result %s %v", intentId, definitionIds)
	}

	if _, _, err := parseImportId("44357f75-2b25-49af-a7ae-f8dbaba1659d/"); err == nil {
		t.Errorf("expected an error for an empty settings list")
	}
	if _, _, err := parseImportId("not-a-guid"); err == nil {
		t.Errorf("expected an error for an invalid intent ID")
	}
}
