package graphBetaAgentInstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// agentInstanceConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the agent instance write has propagated before accepting the read as authoritative.
//
// prevModified is the last_modified_date_time value captured from state immediately before the write:
//   - For Create: pass types.StringNull() — the predicate accepts once last_modified_date_time is
//     populated, confirming the new instance is visible on the responding service node.
//   - For Update: pass state.LastModifiedDateTime — the predicate accepts once the timestamp
//     differs from the pre-update value, confirming the responding node returned the updated
//     version and not a stale cached copy.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func agentInstanceConsistencyPredicate(
	expected *AgentInstanceResourceModel,
	prevModified types.String,
) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AgentInstanceResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// last_modified_date_time is set by the API after every write.
		if actual.LastModifiedDateTime.IsNull() || actual.LastModifiedDateTime.IsUnknown() || actual.LastModifiedDateTime.ValueString() == "" {
			return false
		}

		// For update: the timestamp must have advanced past the pre-update value.
		// A stale replica returns the old object with the old last_modified_date_time.
		if !prevModified.IsNull() && !prevModified.IsUnknown() && prevModified.ValueString() != "" {
			if actual.LastModifiedDateTime.ValueString() == prevModified.ValueString() {
				return false
			}
		}

		return true
	}
}
