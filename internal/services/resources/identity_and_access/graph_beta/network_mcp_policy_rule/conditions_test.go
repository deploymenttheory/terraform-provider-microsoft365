package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// conditionsResponse returns a valid rule response carrying raw as its matchingConditions.
func conditionsResponse(raw string) *mcpPolicyRuleResponse {
	id, name, action, odataType, status := "11111111-1111-1111-1111-111111111111", "test", "allow", mcpPolicyRuleODataType, "enabled"
	priority := int64(1000)
	remote := &mcpPolicyRuleResponse{ID: &id, ODataType: &odataType, Name: &name, Action: &action, Priority: &priority, Conditions: json.RawMessage(raw)}
	remote.Settings = &struct {
		Status *string `json:"status"`
	}{Status: &status}
	return remote
}

// conditionsToState maps raw matchingConditions JSON through MapRemoteStateToTerraform.
func conditionsToState(t *testing.T, raw string) (types.Object, error) {
	t.Helper()
	model := testModel()
	err := MapRemoteStateToTerraform(context.Background(), &model, conditionsResponse(raw))
	return model.MatchingConditions, err
}

func TestUnitResourceNetworkMCPPolicyRule_22_ConditionRoundtrip(t *testing.T) {
	for _, raw := range []string{
		`{"sources":null,"destinations":null}`,
		`{"destinations":{}}`,
		`{"destinations":{"serverUrls":{"values":["https://Example.invalid/mcp/","https://example.invalid/MCP"],"matchType":"notContains"},"protocolVersions":{"values":[],"matchType":"exactMatch"},"insecureConnection":"excluded","missingPrm":"required"}}`,
		`{"destinations":{"toolMatching":{"names":{"values":["tool_a","tool_b"],"matchType":"contains"},"methods":"list,call","classifications":null},"resourceMatching":{"names":{"values":[],"matchType":"notExactMatch"},"methods":"read"},"promptMatching":{"methods":"get"}}}`,
	} {
		t.Run(raw, func(t *testing.T) {
			state, err := conditionsToState(t, raw)
			require.NoError(t, err)
			model := testModel()
			model.MatchingConditions = state
			body, err := constructResource(context.Background(), &model, nil)
			require.NoError(t, err)
			payload, ok := body["matchingConditions"]
			if !ok {
				require.True(t, state.IsNull())
				return
			}
			encoded, err := json.Marshal(payload)
			require.NoError(t, err)
			roundtrip, err := conditionsToState(t, string(encoded))
			require.NoError(t, err)
			require.True(t, state.Equal(roundtrip), "%s", encoded)
		})
	}
}

func TestUnitResourceNetworkMCPPolicyRule_23_ConditionRemoval(t *testing.T) {
	ctx := context.Background()
	state := testModel()
	var err error
	state.MatchingConditions, err = conditionsToState(t, `{"destinations":{"toolMatching":{"names":{"values":["one"],"matchType":"exactMatch"},"methods":"call"}}}`)
	require.NoError(t, err)
	plan := state
	plan.MatchingConditions, err = conditionsToState(t, `{"destinations":{"toolMatching":{"methods":"list,call"}}}`)
	require.NoError(t, err)
	body, err := constructResource(ctx, &plan, &state)
	require.NoError(t, err)
	require.Greater(t, len(body), 1)
	encoded, err := json.Marshal(body)
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
	plan.MatchingConditions = types.ObjectNull(mcpConditionObjectType.AttrTypes)
	body, err = constructResource(ctx, &plan, &state)
	require.NoError(t, err)
	encoded, err = json.Marshal(body)
	require.NoError(t, err)
	patch = nil
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
		_, err := conditionsToState(t, raw)
		require.Error(t, err, raw)
	}
	model := testModel()
	model.MatchingConditions = types.ObjectUnknown(mcpConditionObjectType.AttrTypes)
	_, err := constructResource(context.Background(), &model, nil)
	require.Error(t, err)
}
