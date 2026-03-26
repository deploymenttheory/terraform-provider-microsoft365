package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// autopatchDeploymentStateConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the autopatch deployment state write has propagated before accepting the read
// as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version of the object,
// causing the field to differ and triggering a retry.
//
// effective_value is Computed — set by the Windows Update service — and is checked to confirm
// the replica returned the full deployment state object.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func autopatchDeploymentStateConsistencyPredicate(expected *WindowsUpdatesAutopatchDeploymentStateResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual WindowsUpdatesAutopatchDeploymentStateResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.DeploymentId.ValueString() != expected.DeploymentId.ValueString() {
			return false
		}

		// effective_value is Computed — set by the Windows Update service, never from plan.
		// Its presence confirms the replica returned the full deployment state object.
		if actual.EffectiveValue.IsNull() || actual.EffectiveValue.IsUnknown() || actual.EffectiveValue.ValueString() == "" {
			return false
		}

		// requested_value is the sole user-specified mutable field. A stale replica returns
		// the old requested_value, which would differ from the plan.
		if !expected.RequestedValue.IsNull() && !expected.RequestedValue.IsUnknown() {
			if !actual.RequestedValue.Equal(expected.RequestedValue) {
				return false
			}
		}

		return true
	}
}
