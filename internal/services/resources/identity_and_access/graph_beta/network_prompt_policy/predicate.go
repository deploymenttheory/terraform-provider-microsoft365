package graphBetaNetworkPromptPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// promptPolicyConsistencyPredicate verifies that a successful write is visible on readback.
// The portal-backed endpoint can accept a PATCH while a stale replica still returns old values.
func promptPolicyConsistencyPredicate(
	expected *NetworkPromptPolicyResourceModel,
) func(context.Context, tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkPromptPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}
		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}
		return (expected.Name.IsUnknown() || expected.Name.Equal(actual.Name)) &&
			(expected.Description.IsUnknown() || expected.Description.Equal(actual.Description)) &&
			(expected.DefaultAction.IsUnknown() || expected.DefaultAction.Equal(actual.DefaultAction))
	}
}
