package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// conditionalAccessPolicyConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the conditional access policy write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// prevModified is the modified_date_time value captured from state immediately before the PATCH. The predicate
// accepts once modified_date_time differs from it, confirming the replica returned the updated version and not a
// stale cached copy from before the PATCH. Create sets no predicate: a new policy has no earlier version to confuse
// with, and ReadWithRetry already retries the 404 a lagging replica returns.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func conditionalAccessPolicyConsistencyPredicate(
	expected *ConditionalAccessPolicyResourceModel,
	prevModified types.String,
) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual ConditionalAccessPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// modified_date_time is set by Entra after every write.
		if actual.ModifiedDateTime.IsNull() || actual.ModifiedDateTime.IsUnknown() || actual.ModifiedDateTime.ValueString() == "" {
			return false
		}

		// The timestamp must have advanced past the pre-update value.
		// A stale replica returns the old object, which still carries the old modified_date_time.
		if !prevModified.IsNull() && !prevModified.IsUnknown() && prevModified.ValueString() != "" {
			if actual.ModifiedDateTime.ValueString() == prevModified.ValueString() {
				return false
			}
		}

		return true
	}
}
