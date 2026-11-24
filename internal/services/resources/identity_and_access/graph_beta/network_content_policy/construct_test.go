package graphBetaNetworkContentPolicy

import (
	"context"
	"encoding/json"
	"testing"

	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// encodeContentPolicyBody marshals a constructResource body the way the Config API sends it.
func encodeContentPolicyBody(t *testing.T, body map[string]any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	return payload
}

func TestUnitConstructResourceNetworkContentPolicy_CreatePayload(t *testing.T) {
	model := &NetworkContentPolicyResourceModel{
		Name:          types.StringValue("test for codex"),
		Description:   types.StringValue("sample"),
		DefaultAction: types.StringValue("allow"),
	}
	payload := encodeContentPolicyBody(t, constructResource(context.Background(), model, nil))
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

func TestUnitConstructResourceNetworkContentPolicy_UpdateWithoutPolicyRules(t *testing.T) {
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
	payload := encodeContentPolicyBody(t, constructResource(context.Background(), model, state))
	if _, exists := payload["policyRules"]; exists {
		t.Fatalf("update payload unexpectedly contains policyRules: %#v", payload)
	}
	if _, exists := payload["settings"]; exists {
		t.Fatalf("description-only update unexpectedly contains settings: %#v", payload)
	}
	if _, exists := payload["name"]; exists {
		t.Fatalf("description-only update unexpectedly contains name: %#v", payload)
	}
	if payload["description"] != "sample update" {
		t.Fatalf("description = %#v, expected sample update", payload["description"])
	}
}

func TestUnitConstructResourceNetworkContentPolicy_UpdateClearsDescription(t *testing.T) {
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
	payload := encodeContentPolicyBody(t, constructResource(context.Background(), plan, state))
	if description, exists := payload["description"]; !exists || description != "" {
		t.Fatalf("description = %#v, expected explicit empty string", description)
	}
}

func TestUnitSchemaNetworkContentPolicy_DescriptionIsOptionalAndComputedWithDefault(t *testing.T) {
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
	if !description.Computed {
		t.Fatal("description must be computed so the provider can supply its default")
	}
	if description.Default == nil {
		t.Fatal("description must have an empty-string default")
	}
}

func TestUnitMapRemoteStateToTerraformNetworkContentPolicy_MapsResponse(t *testing.T) {
	var remote contentPolicyResponse
	if err := json.Unmarshal([]byte(`{"id":"id","name":"name","description":"description","settings":{"defaultAction":"allow"}}`), &remote); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	model := &NetworkContentPolicyResourceModel{}
	MapRemoteStateToTerraform(context.Background(), model, &remote)
	if model.Description.ValueString() != "description" || model.DefaultAction.ValueString() != "allow" || model.ID.ValueString() != "id" {
		t.Fatalf("unexpected model: %#v", model)
	}
}
