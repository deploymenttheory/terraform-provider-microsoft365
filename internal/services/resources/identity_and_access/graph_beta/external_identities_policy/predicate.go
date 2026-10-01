package graphBetaExternalIdentitiesPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// stableConsistencyPredicate waits for consecutive reads matching the writable planned settings.
func stableConsistencyPredicate(
	expected *ExternalIdentitiesPolicyResourceModel,
) func(context.Context, tfsdk.State) bool {
	consecutive := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ExternalIdentitiesPolicyResourceModel
		if diags := state.Get(
			ctx,
			&actual,
		); diags.HasError() || actual.ID.ValueString() != singletonID ||
			!expected.AllowExternalIdentitiesToLeave.Equal(actual.AllowExternalIdentitiesToLeave) ||
			!expected.AllowDeletedIdentitiesDataRemoval.Equal(
				actual.AllowDeletedIdentitiesDataRemoval,
			) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= 3
	}
}
