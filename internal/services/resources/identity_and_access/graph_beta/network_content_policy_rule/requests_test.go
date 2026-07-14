package graphBetaNetworkContentPolicyRule

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestContentPolicyRuleRequestSerializesObservedPortalPayload(t *testing.T) {
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
		Destinations: types.ListValueMust(contentPolicyRuleDestinationObjectType(), []attr.Value{
			destinationValue(t, ctx, destinationTypeWebCategory, "AlcoholAndTobacco"),
			destinationValue(t, ctx, destinationTypeFQDN, "example.com", "*.example.com"),
			destinationValue(t, ctx, destinationTypeURL, "https://example.com/path"),
		}),
	}
	body, err := constructResource(ctx, model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	requestInfo, err := newContentPolicyRuleRequestInformation(ctx, contentPolicyRuleTestRequestAdapter{}, abstractions.POST, "parent-id", "", body)
	if err != nil {
		t.Fatalf("newContentPolicyRuleRequestInformation returned error: %v", err)
	}
	if requestInfo.UrlTemplate != contentPolicyRulesURLTemplate {
		t.Fatalf("UrlTemplate = %q", requestInfo.UrlTemplate)
	}
	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload["@odata.type"] != fileRuleODataType || payload["action"] != "scanPurview" {
		t.Fatalf("unexpected root payload: %#v", payload)
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
	return types.ObjectValueMust(contentPolicyRuleDestinationObjectType().AttrTypes, map[string]attr.Value{
		"type":   types.StringValue(destinationType),
		"values": stringSetValue(t, ctx, values...),
	})
}

type contentPolicyRuleTestRequestAdapter struct{}

func (contentPolicyRuleTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (contentPolicyRuleTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (contentPolicyRuleTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (contentPolicyRuleTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}
func (contentPolicyRuleTestRequestAdapter) SetBaseUrl(string)                            {}
func (contentPolicyRuleTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (contentPolicyRuleTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}
