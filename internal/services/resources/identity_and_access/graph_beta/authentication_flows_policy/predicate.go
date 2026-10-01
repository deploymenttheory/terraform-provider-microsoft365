package graphBetaAuthenticationFlowsPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// stableConsistencyPredicate waits for consecutive reads matching the writable planned setting.
func stableConsistencyPredicate(
	expected *AuthenticationFlowsPolicyResourceModel,
) func(context.Context, tfsdk.State) bool {
	consecutive := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AuthenticationFlowsPolicyResourceModel
		if diags := state.Get(
			ctx,
			&actual,
		); diags.HasError() || actual.ID.ValueString() != singletonID ||
			!expected.SelfServiceSignUp.Equal(actual.SelfServiceSignUp) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= 3
	}
}
