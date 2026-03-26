package graphBetaGroupLicenseAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// groupLicenseAssignmentConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the complete group license assignment resource state has propagated across
// Microsoft Entra replicas before accepting the read as authoritative.
//
// Microsoft Entra uses an eventually consistent, multi-replica architecture. After a successful
// assignLicense POST (2xx), reads may be served from a replica that has not yet received the
// write. Any field in the resource could reflect pre-write stale data during this window.
//
// The predicate compares the full read state against the expected state captured at write time:
//
//   - display_name must be populated — a Computed field only set by a real API read,
//     so null/unknown indicates the read state has not yet been refreshed from the API.
//   - id must be non-empty — confirms the composite groupId_skuId key resolved correctly.
//   - group_id must match — confirms the read is for the correct group object.
//   - sku_id must match — confirms the correct license SKU is tracked in state.
//   - disabled_plans count must match — when the SKU is not yet visible in assignedLicenses
//     on the responding replica, MapRemoteResourceStateToTerraform defaults disabled_plans
//     to an empty set regardless of what was written. A count mismatch is the primary signal
//     that the assignedLicenses list has not yet converged on the responding replica.
//
// Retries continue until all conditions are satisfied or the context deadline is reached,
// implementing the polling pattern recommended by Microsoft for Entra eventual consistency:
// https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func groupLicenseAssignmentConsistencyPredicate(expected *GroupLicenseAssignmentResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual GroupLicenseAssignmentResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		// display_name is Computed and only populated by a real API read.
		// Null or unknown means the state has not been refreshed from the API yet.
		if actual.DisplayName.IsNull() || actual.DisplayName.IsUnknown() || actual.DisplayName.ValueString() == "" {
			return false
		}

		// id must be resolved to the composite groupId_skuId value.
		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// group_id must match — confirms we are reading the correct group.
		if actual.GroupId.ValueString() != expected.GroupId.ValueString() {
			return false
		}

		// sku_id must match — confirms the correct license SKU is in state.
		if actual.SkuId.ValueString() != expected.SkuId.ValueString() {
			return false
		}

		// disabled_plans count must match the written value. When the SKU is not yet visible
		// in assignedLicenses on the responding replica, MapRemoteResourceStateToTerraform
		// defaults disabled_plans to an empty set regardless of what was written.
		expectedCount := len(expected.DisabledPlans.Elements())
		if actual.DisabledPlans.IsNull() || actual.DisabledPlans.IsUnknown() {
			return expectedCount == 0
		}
		return len(actual.DisabledPlans.Elements()) == expectedCount
	}
}
