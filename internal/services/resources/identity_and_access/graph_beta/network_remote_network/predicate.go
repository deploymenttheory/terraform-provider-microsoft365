package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// requiredConsecutiveMatches guards against replica lag: one matching read can be followed by
// a stale one, so a write is trusted only after several matching reads in a row.
const requiredConsecutiveMatches = 3

// stableConsistencyPredicate waits for consecutive reads that match the planned configuration.
func stableConsistencyPredicate(expected NetworkRemoteNetworkResourceModel) func(context.Context, tfsdk.State) bool {
	consecutive := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkRemoteNetworkResourceModel
		if state.Raw.IsNull() || state.Get(ctx, &actual).HasError() ||
			!actual.ID.Equal(expected.ID) ||
			!actual.Name.Equal(expected.Name) ||
			!actual.Region.Equal(expected.Region) ||
			!actual.ForwardingProfileIDs.Equal(expected.ForwardingProfileIDs) {
			consecutive = 0
			return false
		}
		consecutive++
		return consecutive >= requiredConsecutiveMatches
	}
}
