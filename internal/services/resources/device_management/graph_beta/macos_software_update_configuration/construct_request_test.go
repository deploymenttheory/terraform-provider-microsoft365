package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestUnitResourceMacOSSoftwareUpdateConfiguration_ResetRequest(t *testing.T) {
	windowType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"start_day": types.StringType, "end_day": types.StringType,
		"start_time": types.StringType, "end_time": types.StringType,
	}}
	window := types.ObjectValueMust(windowType.AttrTypes, map[string]attr.Value{
		"start_day": types.StringValue("monday"), "end_day": types.StringValue("friday"),
		"start_time": types.StringValue("02:00:00"), "end_time": types.StringValue("06:00:00"),
	})
	for _, tc := range []struct {
		name         string
		priority     types.String
		windows      types.List
		wantPriority string
		wantWindows  int
	}{
		{"removed", types.StringNull(), types.ListNull(windowType), "null", 0},
		{"unknown", types.StringUnknown(), types.ListUnknown(windowType), "", -1},
		{"empty", types.StringValue("low"), types.ListValueMust(windowType, []attr.Value{}), `"low"`, 0},
		{"configured", types.StringValue("high"), types.ListValueMust(windowType, []attr.Value{window}), `"high"`, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := constructResource(context.Background(), &MacOSSoftwareUpdateConfigurationResourceModel{
				Priority: tc.priority, CustomUpdateTimeWindows: tc.windows,
			})
			if err != nil {
				t.Fatal(err)
			}
			factory := store.NewBackingStoreSerializationWriterProxyFactory(jsonserialization.NewJsonSerializationWriterFactory())
			writer, err := factory.GetSerializationWriter("application/json")
			if err != nil {
				t.Fatal(err)
			}
			defer writer.Close()
			if err := writer.WriteObjectValue("", body); err != nil {
				t.Fatal(err)
			}
			content, err := writer.GetSerializedContent()
			if err != nil {
				t.Fatal(err)
			}
			var payload map[string]json.RawMessage
			if err := json.Unmarshal(content, &payload); err != nil {
				t.Fatal(err)
			}
			if got := string(payload["priority"]); got != tc.wantPriority {
				t.Fatalf("priority = %q; want %q", got, tc.wantPriority)
			}
			if tc.wantWindows == -1 {
				if _, present := payload["customUpdateTimeWindows"]; present {
					t.Fatal("unknown windows must be omitted")
				}
				return
			}
			var windows []json.RawMessage
			if err := json.Unmarshal(payload["customUpdateTimeWindows"], &windows); err != nil {
				t.Fatal(err)
			}
			if windows == nil || len(windows) != tc.wantWindows {
				t.Fatalf("customUpdateTimeWindows = %s; want an array with %d elements", payload["customUpdateTimeWindows"], tc.wantWindows)
			}
		})
	}
}
