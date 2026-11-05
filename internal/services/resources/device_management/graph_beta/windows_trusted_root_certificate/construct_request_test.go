package graphBetaWindowsTrustedRootCertificate

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

//go:embed tests/test-root.cer
var requestCertificate []byte

func TestUnitCertificateScopeTagsRequest(t *testing.T) {
	ctx := context.Background()
	explicit, diags := types.SetValueFrom(ctx, types.StringType, []string{"0"})
	if diags.HasError() {
		t.Fatal(diags)
	}
	for _, tc := range []struct {
		name  string
		value types.Set
		want  []string
	}{
		{"omitted", types.SetNull(types.StringType), nil},
		{"computed", types.SetUnknown(types.StringType), nil},
		{"explicit", explicit, []string{"0"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := constructResource(ctx, &WindowsTrustedRootCertificateResourceModel{
				TrustedRootCertificate: types.StringValue(base64.StdEncoding.EncodeToString(requestCertificate)),
				DestinationStore:       types.StringValue("computerCertStoreRoot"),
				RoleScopeTagIds:        tc.value,
			})
			if err != nil {
				t.Fatal(err)
			}
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
			raw, present := payload["roleScopeTagIds"]
			if tc.want == nil {
				if present {
					t.Fatalf("omitted/computed roleScopeTagIds must be absent; got %s", raw)
				}
				return
			}
			var got []string
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("roleScopeTagIds = %v; want %v", got, tc.want)
			}
		})
	}
}
