package graphBetaNetworkContentPolicy

import (
	"context"
	"encoding/json"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func TestNewContentPolicyRequestInformationSerializesCreatePayload(t *testing.T) {
	model := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue("sample"),
		DefaultAction: types.StringValue("allow"),
	}
	body, err := constructResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	requestInfo, err := newContentPolicyRequestInformation(context.Background(), contentPolicyTestRequestAdapter{}, abstractions.POST, "", body)
	if err != nil {
		t.Fatalf("newContentPolicyRequestInformation returned error: %v", err)
	}
	if requestInfo.UrlTemplate != contentPoliciesURLTemplate {
		t.Fatalf("UrlTemplate = %q, want %q", requestInfo.UrlTemplate, contentPoliciesURLTemplate)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	if payload["name"] != "test for codex" || payload["description"] != "sample" {
		t.Fatalf("unexpected identity fields: %#v", payload)
	}
	settings, ok := payload["settings"].(map[string]any)
	if !ok || settings["defaultAction"] != "allow" {
		t.Fatalf("settings = %#v, expected defaultAction allow", payload["settings"])
	}
	if _, exists := payload["policyRules"]; exists {
		t.Fatalf("create payload unexpectedly contains policyRules: %#v", payload)
	}
}

func TestNewContentPolicyRequestInformationSerializesUpdateWithoutPolicyRules(t *testing.T) {
	model := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue("sample update"),
		DefaultAction: types.StringValue("allow"),
	}
	state := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue("sample"),
		DefaultAction: types.StringValue("allow"),
	}
	body, err := constructUpdateResource(context.Background(), model, state)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	requestInfo, err := newContentPolicyRequestInformation(context.Background(), contentPolicyTestRequestAdapter{}, abstractions.PATCH, "b11b072b-93cb-4eff-ab63-6a586eccb03b", body)
	if err != nil {
		t.Fatalf("newContentPolicyRequestInformation returned error: %v", err)
	}
	if requestInfo.UrlTemplate != contentPolicyItemURLTemplate {
		t.Fatalf("UrlTemplate = %q, want %q", requestInfo.UrlTemplate, contentPolicyItemURLTemplate)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	if _, exists := payload["policyRules"]; exists {
		t.Fatalf("update payload unexpectedly contains policyRules: %#v", payload)
	}
	if _, exists := payload["settings"]; exists {
		t.Fatalf("description-only update unexpectedly contains settings: %#v", payload)
	}
	if payload["description"] != "sample update" {
		t.Fatalf("description = %#v, expected sample update", payload["description"])
	}
}

func TestConstructUpdateResourceClearsDescription(t *testing.T) {
	plan := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue(""),
		DefaultAction: types.StringValue("allow"),
	}
	state := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue("sample"),
		DefaultAction: types.StringValue("allow"),
	}
	body, err := constructUpdateResource(context.Background(), plan, state)
	if err != nil {
		t.Fatalf("constructUpdateResource returned error: %v", err)
	}
	requestInfo, err := newContentPolicyRequestInformation(context.Background(), contentPolicyTestRequestAdapter{}, abstractions.PATCH, "policy-id", body)
	if err != nil {
		t.Fatalf("newContentPolicyRequestInformation returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(requestInfo.Content, &payload); err != nil {
		t.Fatalf("failed to unmarshal request content: %v", err)
	}
	if description, exists := payload["description"]; !exists || description != "" {
		t.Fatalf("description = %#v, expected explicit empty string", description)
	}
}

func TestDescriptionSchemaIsOptionalAndNotComputed(t *testing.T) {
	resourceUnderTest := NewNetworkContentPolicyResource().(*NetworkContentPolicyResource)
	response := &frameworkresource.SchemaResponse{}
	resourceUnderTest.Schema(context.Background(), frameworkresource.SchemaRequest{}, response)

	description, ok := response.Schema.Attributes["description"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("description attribute has type %T, expected schema.StringAttribute", response.Schema.Attributes["description"])
	}
	if !description.Optional {
		t.Fatal("description must be optional")
	}
	if description.Computed {
		t.Fatal("description must not be computed")
	}
	if len(description.PlanModifiers) != 1 {
		t.Fatalf("description plan modifiers = %d, expected empty-string default modifier", len(description.PlanModifiers))
	}
}

func TestMapRemoteStateToTerraformMapsDescription(t *testing.T) {
	id, name, description, defaultAction := "id", "name", "description", "allow"
	model := &NetworkContentPolicyResourceModel{}
	MapRemoteStateToTerraform(context.Background(), model, &contentPolicyResponse{
		id: &id, name: &name, description: &description, defaultAction: &defaultAction,
	})
	if model.Description.ValueString() != description {
		t.Fatalf("description = %q, expected %q", model.Description.ValueString(), description)
	}
}

type contentPolicyTestRequestAdapter struct{}

func (contentPolicyTestRequestAdapter) Send(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) (s.Parsable, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendEnum(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendCollection(context.Context, *abstractions.RequestInformation, s.ParsableFactory, abstractions.ErrorMappings) ([]s.Parsable, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendEnumCollection(context.Context, *abstractions.RequestInformation, s.EnumFactory, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendPrimitive(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendPrimitiveCollection(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error) {
	return nil, nil
}
func (contentPolicyTestRequestAdapter) SendNoContent(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error {
	return nil
}
func (contentPolicyTestRequestAdapter) GetSerializationWriterFactory() s.SerializationWriterFactory {
	return jsonserialization.NewJsonSerializationWriterFactory()
}
func (contentPolicyTestRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}
func (contentPolicyTestRequestAdapter) SetBaseUrl(string)                            {}
func (contentPolicyTestRequestAdapter) GetBaseUrl() string {
	return "https://graph.microsoft.com/beta"
}
func (contentPolicyTestRequestAdapter) ConvertToNativeRequest(context.Context, *abstractions.RequestInformation) (any, error) {
	return nil, nil
}
