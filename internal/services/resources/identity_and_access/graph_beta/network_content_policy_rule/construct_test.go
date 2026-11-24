package graphBetaNetworkContentPolicyRule

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// encodeContentPolicyRuleBody builds the rule body and marshals it the way the Config API sends it.
func encodeContentPolicyRuleBody(t *testing.T, ctx context.Context, model *NetworkContentPolicyRuleResourceModel) map[string]any {
	t.Helper()
	body, err := constructResource(ctx, model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	return payload
}

func TestUnitConstructResourceNetworkContentPolicyRule_ObservedPortalPayload(t *testing.T) {
	ctx := context.Background()
	model := &NetworkContentPolicyRuleResourceModel{
		Name:             types.StringValue("rule name"),
		Description:      types.StringValue("rule description"),
		Action:           types.StringValue("scanPurview"),
		Priority:         types.Int64Value(101),
		Status:           types.StringValue("enabled"),
		Activities:       stringSetValue(t, ctx, "download", "upload"),
		ContentTypes:     stringSetValue(t, ctx, "text/csv", "application/pdf"),
		TextContentTypes: stringSetValue(t, ctx, "json", "plain", "html", "xml"),
		SessionTypes:     stringSetValue(t, ctx, "user", "agent"),
		Destinations: types.ListValueMust(contentPolicyRuleDestinationObjectType, []attr.Value{
			destinationValue(t, ctx, destinationTypeWebCategory, "AlcoholAndTobacco"),
			destinationValue(t, ctx, destinationTypeFQDN, "example.com", "*.example.com"),
			destinationValue(t, ctx, destinationTypeURL, "https://example.com/path"),
		}),
	}
	payload := encodeContentPolicyRuleBody(t, ctx, model)
	if payload["@odata.type"] != fileRuleODataType || payload["action"] != "scanPurview" || payload["priority"] != float64(101) {
		t.Fatalf("unexpected root payload: %#v", payload)
	}
	if payload["settings"].(map[string]any)["status"] != "enabled" {
		t.Fatalf("settings = %#v", payload["settings"])
	}
	conditions := payload["matchingConditions"].(map[string]any)
	attributes := conditions["fileAttributes"].(map[string]any)
	if attributes["activities"] != "download,upload" {
		t.Fatalf("activities = %#v", attributes["activities"])
	}
	if conditions["sources"].(map[string]any)["sessionType"] != "user,agent" {
		t.Fatalf("sources = %#v", conditions["sources"])
	}
	destinations := conditions["destinations"].([]any)
	if len(destinations) != 3 {
		t.Fatalf("destinations length = %d", len(destinations))
	}
	wantTypes := []string{filePolicyWebCategoryDestinationODataType, filePolicyFQDNDestinationODataType, filePolicyURLDestinationODataType}
	for i, want := range wantTypes {
		if destinations[i].(map[string]any)["@odata.type"] != want {
			t.Fatalf("destination %d = %#v", i, destinations[i])
		}
	}
}

func TestUnitConstructResourceNetworkContentPolicyRule_ClearedOptionalContentTypes(t *testing.T) {
	ctx := context.Background()
	model := &NetworkContentPolicyRuleResourceModel{
		Name:             types.StringValue("rule name"),
		Description:      types.StringValue("rule description"),
		Action:           types.StringValue("scanPurview"),
		Priority:         types.Int64Value(101),
		Status:           types.StringValue("enabled"),
		Activities:       stringSetValue(t, ctx, "download"),
		ContentTypes:     types.SetNull(types.StringType),
		TextContentTypes: stringSetValue(t, ctx, "plain"),
		SessionTypes:     types.SetNull(types.StringType),
		Destinations:     types.ListNull(contentPolicyRuleDestinationObjectType),
	}
	payload := encodeContentPolicyRuleBody(t, ctx, model)
	conditions := payload["matchingConditions"].(map[string]any)
	attributes := conditions["fileAttributes"].(map[string]any)
	contentTypes, ok := attributes["contentTypes"].([]any)
	if !ok || len(contentTypes) != 0 {
		t.Fatalf("contentTypes = %#v, want an explicit empty array", attributes["contentTypes"])
	}
	if textContentTypes := attributes["textContentTypes"].([]any); len(textContentTypes) != 1 || textContentTypes[0] != "plain" {
		t.Fatalf("textContentTypes = %#v", textContentTypes)
	}
	if sources := conditions["sources"].(map[string]any); len(sources) != 0 {
		t.Fatalf("sources = %#v, want an empty object", sources)
	}
	if _, exists := conditions["destinations"]; exists {
		t.Fatalf("null destinations must be omitted: %#v", conditions)
	}
}

func TestUnitMapRemoteStateToTerraformNetworkContentPolicyRule_SplitsAndSkipsNulls(t *testing.T) {
	ctx := context.Background()
	raw := `{"id":"rule","name":"n","action":"block","priority":101,"settings":{"status":"enabled"},"matchingConditions":{"fileAttributes":{"activities":"download, ,upload","contentTypes":["text/csv",null],"textContentTypes":[]},"destinations":[{"@odata.type":"` + filePolicyURLDestinationODataType + `","values":["https://example.com/a"]},null],"sources":{"sessionType":"user"}}}`
	var remote contentPolicyRuleResponse
	if err := json.Unmarshal([]byte(raw), &remote); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	var data NetworkContentPolicyRuleResourceModel
	MapRemoteStateToTerraform(ctx, &data, &remote)
	if len(data.Activities.Elements()) != 2 || len(data.ContentTypes.Elements()) != 1 || !data.TextContentTypes.IsNull() || len(data.SessionTypes.Elements()) != 1 {
		t.Fatalf("unexpected sets: %#v", data)
	}
	if len(data.Destinations.Elements()) != 1 {
		t.Fatalf("destinations = %#v", data.Destinations)
	}
	destination := data.Destinations.Elements()[0].(types.Object).Attributes()["type"].(types.String)
	if destination.ValueString() != destinationTypeURL || data.Status.ValueString() != "enabled" || data.Priority.ValueInt64() != 101 {
		t.Fatalf("unexpected mapping: %#v", data)
	}
}

func stringSetValue(t *testing.T, ctx context.Context, values ...string) types.Set {
	t.Helper()
	value, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %s", diags.Errors()[0].Detail())
	}
	return value
}

func destinationValue(t *testing.T, ctx context.Context, destinationType string, values ...string) attr.Value {
	t.Helper()
	fields := map[string]attr.Value{"type": types.StringValue(destinationType), "values": stringSetValue(t, ctx, values...)}
	return types.ObjectValueMust(contentPolicyRuleDestinationObjectType.AttrTypes, fields)
}
