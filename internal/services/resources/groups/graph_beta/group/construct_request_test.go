package graphBetaGroup

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestUnitGroupTypesRequest(t *testing.T) {
	ctx := context.Background()
	dynamic, diags := types.SetValueFrom(ctx, types.StringType, []string{"DynamicMembership"})
	if diags.HasError() {
		t.Fatal(diags)
	}
	for _, update := range []bool{false, true} {
		for _, tc := range []struct {
			name  string
			value types.Set
			want  []string
		}{
			{"omitted", types.SetNull(types.StringType), nil},
			{"computed", types.SetUnknown(types.StringType), nil},
			{"explicit", dynamic, []string{"DynamicMembership"}},
		} {
			t.Run(fmt.Sprintf("%s/update=%t", tc.name, update), func(t *testing.T) {
				body, err := constructResource(ctx, &GroupResourceModel{
					DisplayName: types.StringValue("test-static-security-group"),
					GroupTypes:  tc.value,
				}, update)
				if err != nil {
					t.Fatal(err)
				}
				// Match the request adapter: backing-store serialization emits
				// explicitly nil properties even when plain JSON serialization omits them.
				factory := store.NewBackingStoreSerializationWriterProxyFactory(
					jsonserialization.NewJsonSerializationWriterFactory(),
				)
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
				raw, present := payload["groupTypes"]
				if tc.want == nil {
					if present {
						t.Fatalf("omitted/computed groupTypes must be absent; got %s", raw)
					}
					return
				}
				var got []string
				if err := json.Unmarshal(raw, &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.want) {
					t.Fatalf("groupTypes = %v; want %v", got, tc.want)
				}
			})
		}
	}
}
