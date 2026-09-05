package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestUnitResourceNetworkPromptPolicyRule_22_EmptySchemesPatch(t *testing.T) {
	var patch map[string]any
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
		if q.Method == "PATCH" {
			decodeRequest(t, q, &patch)
		}
		var payload map[string]any
		require.NoError(t, json.Unmarshal(successFixture(t), &payload))
		payload["matchingConditions"].(map[string]any)["conversationSchemes"] = []any{}
		b, e := json.Marshal(payload)
		require.NoError(t, e)
		writeResponse(t, w, 200, b)
	})
	state := testState(t, r, testModel())
	plan := testModel()
	plan.ConversationSchemes = types.ListValueMust(conversationSchemeObjectType(), []attr.Value{})
	resp := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: state.Schema, Raw: testState(t, r, plan).Raw}}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.Equal(t, []any{}, patch["matchingConditions"].(map[string]any)["conversationSchemes"])
	require.Equal(t, "maliciousPromptDetected", patch["matchingConditions"].(map[string]any)["scanResult"])
}

func TestUnitResourceNetworkPromptPolicyRule_23_IgnoredUpdate(t *testing.T) {
	r := testClient(t, func(w http.ResponseWriter, q *http.Request) { writeResponse(t, w, 200, successFixture(t)) })
	state := testState(t, r, testModel())
	plan := testModel()
	plan.Enabled = types.BoolValue(false)
	resp := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{State: state, Plan: tfsdk.Plan{Schema: state.Schema, Raw: testState(t, r, plan).Raw}}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "Enabled")
	var got NetworkPromptPolicyRuleResourceModel
	require.False(t, resp.State.Get(context.Background(), &got).HasError())
	require.True(t, got.Enabled.ValueBool())
	require.Equal(t, testModel().ID, got.ID)
}

func TestUnitResourceNetworkPromptPolicyRule_24_JSONPathDefault(t *testing.T) {
	for _, omitted := range []bool{true, false} {
		r := testClient(t, func(w http.ResponseWriter, q *http.Request) {
			var payload map[string]any
			require.NoError(t, json.Unmarshal(successFixture(t), &payload))
			payload["matchingConditions"].(map[string]any)["conversationSchemes"].([]any)[0].(map[string]any)["jsonPath"] = ""
			b, e := json.Marshal(payload)
			require.NoError(t, e)
			writeResponse(t, w, 200, b)
		})
		model := testModel()
		fields := model.ConversationSchemes.Elements()[0].(types.Object).Attributes()
		fields["json_path"] = types.StringValue("")
		if omitted {
			fields["json_path"] = types.StringNull()
		}
		model.ConversationSchemes = types.ListValueMust(conversationSchemeObjectType(), []attr.Value{types.ObjectValueMust(conversationSchemeObjectType().AttrTypes, fields)})
		state := testState(t, r, model)
		resp := resource.ReadResponse{State: state}
		r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
		require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
		var got NetworkPromptPolicyRuleResourceModel
		require.False(t, resp.State.Get(context.Background(), &got).HasError())
		require.True(t, model.ConversationSchemes.Equal(got.ConversationSchemes))
	}
}

func TestUnitResourceNetworkPromptPolicyRule_25_ConditionalValidation(t *testing.T) {
	for _, tc := range []struct {
		typ                       string
		url, jsonPath, schemeName types.String
		invalid                   bool
	}{
		{"custom", types.StringNull(), types.StringNull(), types.StringNull(), true},
		{"custom", types.StringValue("ftp://example.com/chat"), types.StringNull(), types.StringNull(), true},
		{"custom", types.StringUnknown(), types.StringNull(), types.StringNull(), false},
		{"custom", types.StringValue("https://example.com/chat"), types.StringNull(), types.StringValue("chatGpt"), true},
		{"predefined", types.StringNull(), types.StringNull(), types.StringNull(), true},
		{"predefined", types.StringNull(), types.StringValue(""), types.StringValue("chatGpt"), true},
		{"predefined", types.StringNull(), types.StringNull(), types.StringValue("chatGpt"), false},
	} {
		r := NewNetworkPromptPolicyRuleResource().(*NetworkPromptPolicyRuleResource)
		model := testModel()
		model.ConversationSchemes = types.ListValueMust(conversationSchemeObjectType(), []attr.Value{types.ObjectValueMust(conversationSchemeObjectType().AttrTypes, map[string]attr.Value{"type": types.StringValue(tc.typ), "url": tc.url, "json_path": tc.jsonPath, "scheme_name": tc.schemeName})})
		state := testState(t, r, model)
		resp := resource.ValidateConfigResponse{}
		r.ValidateConfig(context.Background(), resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: state.Schema, Raw: state.Raw}}, &resp)
		require.Equal(t, tc.invalid, resp.Diagnostics.HasError(), "%+v: %v", tc, resp.Diagnostics)
	}
}
