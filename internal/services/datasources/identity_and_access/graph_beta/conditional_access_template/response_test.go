package graphBetaConditionalAccessTemplate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplateResponsePreservesEvolvableValues(t *testing.T) {
	var response templateResponse
	err := json.Unmarshal([]byte(`{
		"id": "template-id",
		"scenarios": "zeroTrust, aiAgents,futureScenario",
		"details": {"grantControls": {
			"builtInControls": ["mfa", "verifiedID", "futureControl"]
		}}
	}`), &response)
	require.NoError(t, err)

	ctx := context.Background()
	state := MapRemoteStateToDataSource(ctx, &response)
	var scenarios, controls []string
	require.False(t, state.Scenarios.ElementsAs(ctx, &scenarios, false).HasError())
	require.False(t, state.Details.GrantControls.BuiltInControls.ElementsAs(ctx, &controls, false).HasError())
	require.ElementsMatch(t, []string{"zeroTrust", "aiAgents", "futureScenario"}, scenarios)
	require.Equal(t, []string{"mfa", "verifiedID", "futureControl"}, controls)
}

func TestTemplateResponseRejectsMalformedValues(t *testing.T) {
	for name, payload := range map[string]string{
		"scenarios": `{"scenarios":42}`,
		"controls":  `{"details":{"grantControls":{"builtInControls":[42]}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			var response templateResponse
			require.Error(t, json.Unmarshal([]byte(payload), &response))
		})
	}
}
