package graphBetaConditionalAccessPolicy

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

// Graph rejects null device lists, so a policy that only sets device_filter must leave them out.
func TestConstructDevicesOmitsUnsetLists(t *testing.T) {
	data := &ConditionalAccessDevices{
		IncludeDevices:      types.SetNull(types.StringType),
		ExcludeDevices:      types.SetNull(types.StringType),
		IncludeDeviceStates: types.SetUnknown(types.StringType),
		ExcludeDeviceStates: types.SetNull(types.StringType),
		DeviceFilter: &ConditionalAccessFilter{
			Mode: types.StringValue("exclude"),
			Rule: types.StringValue(`device.isCompliant -eq True`),
		},
	}

	devices, err := constructDevices(context.Background(), data)
	if err != nil {
		t.Fatalf("constructDevices returned error: %v", err)
	}
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", devices); err != nil {
		t.Fatalf("serialize devices: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized devices: %v", err)
	}

	body := string(content)
	for _, field := range []string{"includeDevices", "excludeDevices", "includeDeviceStates", "excludeDeviceStates"} {
		if strings.Contains(body, `"`+field+`"`) {
			t.Errorf("expected %s to be omitted, got %s", field, body)
		}
	}
	if !strings.Contains(body, `"deviceFilter"`) {
		t.Errorf("expected deviceFilter to be sent, got %s", body)
	}
}

func TestConstructDevicesRejectsNonStringLists(t *testing.T) {
	numbers := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)})
	cases := map[string]func(*ConditionalAccessDevices){
		"include devices":       func(data *ConditionalAccessDevices) { data.IncludeDevices = numbers },
		"exclude devices":       func(data *ConditionalAccessDevices) { data.ExcludeDevices = numbers },
		"include device states": func(data *ConditionalAccessDevices) { data.IncludeDeviceStates = numbers },
		"exclude device states": func(data *ConditionalAccessDevices) { data.ExcludeDeviceStates = numbers },
	}

	for name, apply := range cases {
		t.Run(name, func(t *testing.T) {
			data := &ConditionalAccessDevices{}
			apply(data)
			_, err := constructDevices(context.Background(), data)
			if err == nil || !strings.Contains(err.Error(), "failed to convert "+name) {
				t.Fatalf("expected failed to convert %s, got %v", name, err)
			}
		})
	}
}
