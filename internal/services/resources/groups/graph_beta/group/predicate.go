package graphBetaGroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// groupConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the group resource write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// The predicate checks that the group ID and display_name are populated in state,
// confirming the group is visible on the responding replica.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func groupConsistencyPredicate(expected *GroupResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual GroupResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// display_name must match the written value — a stale replica returns the old version.
		if !actual.DisplayName.Equal(expected.DisplayName) {
			return false
		}

		// mail_nickname must also match.
		if !actual.MailNickname.Equal(expected.MailNickname) {
			return false
		}

		return true
	}
}
