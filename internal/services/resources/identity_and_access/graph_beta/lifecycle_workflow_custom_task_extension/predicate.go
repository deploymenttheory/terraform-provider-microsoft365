package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// requiredConsecutiveFreshReads is how many reads in a row must reflect a PATCH before Update
// accepts the state. With the default 2 second retry interval this needs about 6 seconds of
// stable reads.
const requiredConsecutiveFreshReads = 4

// customTaskExtensionConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// accepts a read only once it reflects the preceding PATCH (ADR-002).
//
// Probing the tenant showed that after a PATCH returns 204, GETs alternate between the new and the
// previous custom task extension for 5 to 8 seconds before they settle (for example
// "oooNNoNNNN..." at 0.75 second intervals, where o is the old value). Accepting the first fresh
// read let the refresh that follows apply read a stale replica and plan the update again, so the
// predicate requires requiredConsecutiveFreshReads fresh reads in a row and starts counting again
// after any stale one.
//
// prevModified is last_modified_date_time captured from state before the write. A read is fresh
// when its last_modified_date_time differs from prevModified. When prevModified is null (no PATCH
// was sent) the first successful read is accepted.
func customTaskExtensionConsistencyPredicate(prevModified types.String) func(ctx context.Context, state tfsdk.State) bool {
	freshReads := 0

	return func(ctx context.Context, state tfsdk.State) bool {
		var actual LifecycleWorkflowCustomTaskExtensionResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			freshReads = 0
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			freshReads = 0
			return false
		}

		if actual.LastModifiedDateTime.IsNull() || actual.LastModifiedDateTime.IsUnknown() || actual.LastModifiedDateTime.ValueString() == "" {
			freshReads = 0
			return false
		}

		if prevModified.IsNull() || prevModified.IsUnknown() || prevModified.ValueString() == "" {
			return true
		}

		if actual.LastModifiedDateTime.ValueString() == prevModified.ValueString() {
			freshReads = 0
			return false
		}

		freshReads++
		return freshReads >= requiredConsecutiveFreshReads
	}
}
