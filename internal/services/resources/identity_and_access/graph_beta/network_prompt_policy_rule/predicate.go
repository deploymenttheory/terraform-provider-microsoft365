package graphBetaNetworkPromptPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// promptPolicyRuleConsistencyPredicate verifies that a successful write is visible on readback.
// The portal-backed endpoint can accept a PATCH while a stale replica still returns old values.
func promptPolicyRuleConsistencyPredicate(
	expected *NetworkPromptPolicyRuleResourceModel,
) func(context.Context, tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkPromptPolicyRuleResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}
		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}
		return (expected.Name.IsUnknown() || expected.Name.Equal(actual.Name)) &&
			(expected.Description.IsUnknown() || expected.Description.Equal(actual.Description)) &&
			(expected.Action.IsUnknown() || expected.Action.Equal(actual.Action)) &&
			(expected.Priority.IsUnknown() || expected.Priority.Equal(actual.Priority)) &&
			(expected.Enabled.IsUnknown() || expected.Enabled.Equal(actual.Enabled)) &&
			(expected.PromptLogging.IsUnknown() || expected.PromptLogging.Equal(actual.PromptLogging)) &&
			(expected.ScanResult.IsUnknown() || expected.ScanResult.Equal(actual.ScanResult)) &&
			(expected.ConversationSchemes.IsUnknown() || expected.ConversationSchemes.Equal(actual.ConversationSchemes))
	}
}
