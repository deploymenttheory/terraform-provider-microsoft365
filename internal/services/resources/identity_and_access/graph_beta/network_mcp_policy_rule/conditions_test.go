package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkMCPPolicyRule_22_ConditionRoundtrip(t *testing.T) {
	for _, raw := range []string{
		`{"sources":null,"destinations":null}`,
		`{"destinations":{}}`,
		`{"destinations":{"serverUrls":{"values":["https://Example.invalid/mcp/","https://example.invalid/MCP"],"matchType":"notContains"},"protocolVersions":{"values":[],"matchType":"exactMatch"},"insecureConnection":"excluded","missingPrm":"required"}}`,
		`{"destinations":{"toolMatching":{"names":{"values":["tool_a","tool_b"],"matchType":"contains"},"methods":"list,call","classifications":null},"resourceMatching":{"names":{"values":[],"matchType":"notExactMatch"},"methods":"read"},"promptMatching":{"methods":"get"}}}`,
	} {
		t.Run(raw, func(t *testing.T) {
			ctx := context.Background()
			state, err := conditionsToState(ctx, json.RawMessage(raw))
			require.NoError(t, err)
			model := testModel()
			model.MatchingConditions = state
			payload, err := ruleConditions(ctx, &model)
			require.NoError(t, err)
			if payload == nil {
				require.True(t, state.IsNull())
				return
			}
			writer := jsonserialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", payload))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			roundtrip, err := conditionsToState(ctx, encoded)
			require.NoError(t, err)
			require.True(t, state.Equal(roundtrip), "%s", encoded)
		})
	}
}

func TestUnitResourceNetworkMCPPolicyRule_23_ConditionRemoval(t *testing.T) {
	ctx := context.Background()
	state := testModel()
	var err error
	state.MatchingConditions, err = conditionsToState(
		ctx,
		json.RawMessage(
			`{"destinations":{"toolMatching":{"names":{"values":["one"],"matchType":"exactMatch"},"methods":"call"}}}`,
		),
	)
	require.NoError(t, err)
	plan := state
	plan.MatchingConditions, err = conditionsToState(
		ctx,
		json.RawMessage(`{"destinations":{"toolMatching":{"methods":"list,call"}}}`),
	)
	require.NoError(t, err)
	body, err := constructUpdateResource(ctx, &plan, &state)
	require.NoError(t, err)
	require.True(t, body.hasChanges())
	writer := jsonserialization.NewJsonSerializationWriter()
	require.NoError(t, writer.WriteObjectValue("", body))
	encoded, err := writer.GetSerializedContent()
	require.NoError(t, err)
	var patch map[string]any
	require.NoError(t, json.Unmarshal(encoded, &patch))
	dest := patch["matchingConditions"].(map[string]any)["destinations"].(map[string]any)
	tool := dest["toolMatching"].(map[string]any)
	require.Contains(t, tool, "names")
	require.Nil(t, tool["names"])
	require.Equal(t, "list,call", tool["methods"])
	require.Contains(t, dest, "serverUrls")
	require.Nil(t, dest["serverUrls"])
	plan.MatchingConditions = types.ObjectNull(conditionAttributeTypes())
	body, err = constructUpdateResource(ctx, &plan, &state)
	require.NoError(t, err)
	writer = jsonserialization.NewJsonSerializationWriter()
	require.NoError(t, writer.WriteObjectValue("", body))
	encoded, err = writer.GetSerializedContent()
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(encoded, &patch))
	require.Contains(t, patch, "matchingConditions")
	require.Nil(t, patch["matchingConditions"])
}

func TestUnitResourceNetworkMCPPolicyRule_24_UnsupportedConditions(t *testing.T) {
	for _, raw := range []string{
		`{"sources":{"users":["x"]}}`,
		`{"destinations":{"notificationMethods":"x"}}`,
		`{"destinations":{"toolMatching":{"callArguments":{"values":["x"]}}}}`,
		`{"destinations":{"serverUrls":{"values":[null],"matchType":"exactMatch"}}}`,
		`{"destinations":{"serverUrls":{"values":"x","matchType":"exactMatch"}}}`,
		`{"destinations":[]}`,
		`{"destinations":{"serverUrls":{}}}`,
		`{"destinations":{"serverUrls":{"values":[],"matchType":"future"}}}`,
		`{"destinations":{"missingPrm":"future"}}`,
	} {
		_, err := conditionsToState(context.Background(), json.RawMessage(raw))
		require.Error(t, err, raw)
	}
	model := testModel()
	model.MatchingConditions = types.ObjectUnknown(conditionAttributeTypes())
	_, err := ruleConditions(context.Background(), &model)
	require.Error(t, err)
}
