package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// autopatchUpdatableAssetGroupConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the autopatch updatable asset group write has propagated before
// accepting the read as authoritative.
//
// The predicate compares the full entra_device_object_ids set from the expected plan against
// the state returned by the read. The addMembersById / removeMembersById APIs are eventually
// consistent; a stale service node may return an outdated member list that does not reflect
// the written state. Retrying until the member set matches ensures the accepted read is
// authoritative.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func autopatchUpdatableAssetGroupConsistencyPredicate(expected *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// entra_device_object_ids must match exactly — addMembersById and removeMembersById are
		// eventually consistent; a stale service node may return an incomplete or outdated set.
		if !expected.EntraDeviceObjectIds.IsNull() && !expected.EntraDeviceObjectIds.IsUnknown() {
			if !actual.EntraDeviceObjectIds.Equal(expected.EntraDeviceObjectIds) {
				return false
			}
		} else {
			// expected is null/unknown — accept only if actual is also null or empty
			if !actual.EntraDeviceObjectIds.IsNull() && !actual.EntraDeviceObjectIds.IsUnknown() && len(actual.EntraDeviceObjectIds.Elements()) > 0 {
				return false
			}
		}

		return true
	}
}
