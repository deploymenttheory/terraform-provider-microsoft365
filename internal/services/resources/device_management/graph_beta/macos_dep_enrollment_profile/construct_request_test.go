package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"encoding/json"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestUnitResourceMacOSDepEnrollmentProfile_GroupIdsRequest(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value types.Set
		want  string
	}{
		{"omitted", types.SetNull(types.StringType), ""},
		{"unknown", types.SetUnknown(types.StringType), ""},
		{"empty", types.SetValueMust(types.StringType, []attr.Value{}), "[]"},
		{"configured", types.SetValueMust(types.StringType, []attr.Value{types.StringValue("10000000-0000-0000-0000-000000000001")}), `["10000000-0000-0000-0000-000000000001"]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, err := constructResource(context.Background(), &MacOSDepEnrollmentProfileResourceModel{EnrollmentTimeAzureAdGroupIds: tc.value})
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
			if got := string(payload["enrollmentTimeAzureAdGroupIds"]); got != tc.want {
				t.Fatalf("enrollmentTimeAzureAdGroupIds = %q; want %q", got, tc.want)
			}
			if got := string(payload["depProfileAdminAccountPasswordRotationSetting"]); got != "null" {
				t.Fatalf("removed password rotation must serialize as null; got %q", got)
			}
		})
	}
}

func TestUnitResourceMacOSDepEnrollmentProfile_DuplicateSkipKeys(t *testing.T) {
	content, err := os.ReadFile("tests/responses/enrollment_profile_duplicate_skip_keys.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		EnabledSkipKeys []string `json:"enabledSkipKeys"`
	}
	if err := json.Unmarshal(content, &fixture); err != nil {
		t.Fatal(err)
	}
	profile := graphmodels.NewDepMacOSEnrollmentProfile()
	profile.SetEnabledSkipKeys(fixture.EnabledSkipKeys)
	var data MacOSDepEnrollmentProfileResourceModel
	MapRemoteStateToTerraform(context.Background(), &data, profile, "")
	if data.EnabledSkipKeys.IsNull() || data.EnabledSkipKeys.IsUnknown() || len(data.EnabledSkipKeys.Elements()) != 27 {
		t.Fatalf("expected 27 distinct skip keys from the live PATCH response, got %s", data.EnabledSkipKeys.String())
	}
	if !data.AppearanceDisabled.ValueBool() || !data.IntelligenceDisabled.ValueBool() {
		t.Fatal("new pane controls were not recovered from the returned skip keys")
	}
}
