package graphBetaGroupPolicyDefinition

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kiotajson "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestListBoxRoundTrip(t *testing.T) {
	value := `[{"name":"https://example.test","value":"2"},{"name":"https://other.test","value":"1"}]`
	presentation, err := ConvertValueForType(context.Background(), "Enter the zone assignments here.", value, "#microsoft.graph.groupPolicyPresentationListBox")
	if err != nil {
		t.Fatal(err)
	}
	if got := ExtractValueFromPresentation(presentation); got != value {
		t.Fatalf("ListBox state = %q, want %q", got, value)
	}
}

func TestListBoxValidation(t *testing.T) {
	for _, value := range []string{`null`, `{}`, `[{"name":"a"}]`, `[{"name":"a","value":2}]`, `[{"name":"a","value":"b","extra":"c"}]`, `[] []`, `invalid`} {
		t.Run(value, func(t *testing.T) {
			if _, err := ConvertValueForType(context.Background(), "List", value, "#microsoft.graph.groupPolicyPresentationListBox"); err == nil {
				t.Fatalf("accepted invalid ListBox %q", value)
			}
			if err := validateValueForType(context.Background(), "List", value, "#microsoft.graph.groupPolicyPresentationListBox"); err == nil {
				t.Fatalf("validation accepted %q", value)
			}
		})
	}
}

func TestListBoxEmptyAndSpecialCharacters(t *testing.T) {
	for _, value := range []string{`[]`, `[{"name":"a,b=c","value":"line1\nline2 \\ \"quoted\""}]`, `[{"name":"","value":""}]`} {
		presentation, err := ConvertValueForType(context.Background(), "List", value, "#microsoft.graph.groupPolicyPresentationListBox")
		if err != nil {
			t.Fatal(err)
		}
		if got := ExtractValueFromPresentation(presentation); got != value {
			t.Fatalf("round trip = %q, want %q", got, value)
		}
	}
}

func TestListBoxConstructAndState(t *testing.T) {
	ctx := context.Background()
	value := `[ { "value": "2", "name": "https://example.test" } ]`
	data := GroupPolicyDefinitionResourceModel{
		Enabled: types.BoolValue(true),
		AdditionalData: map[string]any{
			"definitionTemplateID":      "definition-template",
			"definitionValueInstanceID": "definition-instance",
			"resolvedPresentations":     []ResolvedPresentation{{Label: "List", Type: "#microsoft.graph.groupPolicyPresentationListBox", TemplateID: "presentation-template", InstanceID: "presentation-instance"}},
		},
	}

	values, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType, "label": types.StringType, "value": types.StringType}}, []PresentationValue{{ID: types.StringValue("presentation-template"), Label: types.StringValue("List"), Value: types.StringValue(value)}})
	if diags.HasError() {
		t.Fatal(diags)
	}
	data.Values = values

	for _, operation := range []string{constants.TfOperationCreate, constants.TfOperationUpdate} {
		body, err := constructResource(ctx, &data, operation)
		if err != nil {
			t.Fatal(err)
		}
		writer := kiotajson.NewJsonSerializationWriter()
		if err := writer.WriteObjectValue("", body); err != nil {
			t.Fatal(err)
		}
		encoded, err := writer.GetSerializedContent()
		if err != nil {
			t.Fatal(err)
		}
		var payload map[string]any
		if err := json.Unmarshal(encoded, &payload); err != nil {
			t.Fatal(err)
		}
		field := "added"
		if operation == constants.TfOperationUpdate {
			field = "updated"
		}
		dv := payload[field].([]any)[0].(map[string]any)
		pv := dv["presentationValues"].([]any)[0].(map[string]any)
		if pv["@odata.type"] != "#microsoft.graph.groupPolicyPresentationValueList" {
			t.Fatalf("wrong ListBox discriminator: %s", encoded)
		}
		pair := pv["values"].([]any)[0].(map[string]any)
		if pair["name"] != "https://example.test" || pair["value"] != "2" {
			t.Fatalf("wrong name/value payload: %s", encoded)
		}
		if pv["presentation@odata.bind"] == nil {
			t.Fatalf("missing presentation binding: %s", encoded)
		}
		if operation == constants.TfOperationUpdate && pv["id"] != "presentation-instance" {
			t.Fatalf("missing instance ID: %s", encoded)
		}
	}
	presentation, err := ConvertValueForType(ctx, "List", value, "#microsoft.graph.groupPolicyPresentationListBox")
	if err != nil {
		t.Fatal(err)
	}
	instance := "presentation-instance"
	presentation.SetId(&instance)
	MapRemoteStateToTerraform(ctx, &data, []models.GroupPolicyPresentationValueable{presentation}, nil)
	var mapped []PresentationValue
	if diags := data.Values.ElementsAs(ctx, &mapped, false); diags.HasError() {
		t.Fatal(diags)
	}
	if len(mapped) != 1 || mapped[0].Value.ValueString() != value {
		t.Fatalf("equivalent configured JSON formatting was lost: %#v", mapped)
	}
	drift, err := ConvertValueForType(ctx, "List", `[{"name":"drift","value":"4"}]`, "#microsoft.graph.groupPolicyPresentationListBox")
	if err != nil {
		t.Fatal(err)
	}
	drift.SetId(&instance)
	MapRemoteStateToTerraform(ctx, &data, []models.GroupPolicyPresentationValueable{drift}, nil)
	if diags := data.Values.ElementsAs(ctx, &mapped, false); diags.HasError() {
		t.Fatal(diags)
	}
	if mapped[0].Value.ValueString() != `[{"name":"drift","value":"4"}]` {
		t.Fatal("remote ListBox drift was hidden")
	}
}
