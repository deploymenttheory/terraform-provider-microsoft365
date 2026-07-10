package graphBetaGroupLicenseAssignment

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// groupLicenseAssignmentConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the complete group license assignment resource state has propagated across
// Microsoft Entra replicas before accepting the read as authoritative.
//
// Microsoft Entra uses an eventually consistent, multi-replica architecture, and assignLicense
// itself is asynchronous (202 Accepted): after a successful POST, reads may be served from a
// replica that has not yet received the write, or the backend licensing service may not have
// processed the change yet. Any field in the resource could reflect pre-write stale data during
// this window.
//
// The predicate compares the full read state against the expected state captured at write time:
//
//   - the state must not be null — Read removes the resource from state when the managed SKU
//     is absent from the group's assignedLicenses, so a null state means the assignment has
//     not been observed on the responding replica yet. This is the primary signal that the
//     asynchronous assignLicense processing has not completed (or has silently failed).
//   - display_name must be populated — a Computed field only set by a real API read,
//     so null/unknown indicates the read state has not yet been refreshed from the API.
//   - id must be non-empty — confirms the composite groupId_skuId key resolved correctly.
//   - group_id must match — confirms the read is for the correct group object.
//   - sku_id must match — confirms the correct license SKU is tracked in state.
//   - disabled_plans must match as a set — confirms the disabled plans visible on the
//     responding replica are exactly the ones that were written, not stale pre-write data.
//
// Retries continue until all conditions are satisfied or the context deadline is reached,
// implementing the polling pattern recommended by Microsoft for Entra eventual consistency:
// https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func groupLicenseAssignmentConsistencyPredicate(expected *GroupLicenseAssignmentResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		// Read removes the resource from state when the managed SKU is not present in the
		// group's assignedLicenses. Treat that as "write not yet confirmed" and keep polling.
		if state.Raw.IsNull() {
			return false
		}

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

		// disabled_plans must match the written value as a set, so that stale pre-write
		// data with the same cardinality but different plan IDs is not accepted.
		return uuidSetsEqual(expected.DisabledPlans, actual.DisabledPlans)
	}
}

// uuidSetsEqual compares two types.Set values of UUID strings for set equality, treating
// null/unknown as the empty set. UUIDs are compared case-insensitively: the API returns
// them in canonical lowercase form while the configured values may use any casing.
func uuidSetsEqual(a, b types.Set) bool {
	aElems := setElementsAsLowercaseStrings(a)
	bElems := setElementsAsLowercaseStrings(b)
	if len(aElems) != len(bElems) {
		return false
	}
	for elem := range aElems {
		if !bElems[elem] {
			return false
		}
	}
	return true
}

func setElementsAsLowercaseStrings(s types.Set) map[string]bool {
	result := make(map[string]bool)
	if s.IsNull() || s.IsUnknown() {
		return result
	}
	for _, elem := range s.Elements() {
		if strVal, ok := elem.(types.String); ok {
			result[strings.ToLower(strVal.ValueString())] = true
		}
	}
	return result
}
