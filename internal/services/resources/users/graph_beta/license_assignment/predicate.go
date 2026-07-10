package graphBetaUserLicenseAssignment

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// licenseAssignmentConsistencyPredicate returns a consistency predicate for ReadWithRetry that
// verifies the complete license assignment resource state has propagated across Microsoft Entra
// replicas before accepting the read as authoritative.
//
// Microsoft Entra uses an eventually consistent, multi-replica architecture. After a successful
// assignLicense POST (2xx), reads may be served from a replica that has not yet received the
// write. Any field in the resource could reflect pre-write stale data during this window.
// This causes terraform state issues that are due to the timing.
//
// The predicate compares the full read state against the expected state captured at write time:
//
//   - user_principal_name must be populated — a Computed field only set by a real API read,
//     so null/unknown indicates the read state has not yet been refreshed from the API.
//   - id must be non-empty — confirms the composite userId_skuId key resolved correctly.
//   - user_id must match — confirms the read is for the correct user object.
//   - sku_id must match — confirms the correct license SKU is tracked in state.
//   - the state must not be null — Read removes the resource from state when the managed SKU
//     is absent from the user's assignedLicenses, so a null state means the assignment has
//     not been observed on the responding replica yet.
//   - disabled_plans must match as a set — confirms the disabled plans visible on the
//     responding replica are exactly the ones that were written, not stale pre-write data.
//
// Retries continue until all conditions are satisfied or the context deadline is reached,
// implementing the polling pattern recommended by Microsoft for Entra eventual consistency:
// https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func licenseAssignmentConsistencyPredicate(expected *UserLicenseAssignmentResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		// Read removes the resource from state when the managed SKU is not present in the
		// user's assignedLicenses. Treat that as "write not yet confirmed" and keep polling.
		if state.Raw.IsNull() {
			return false
		}

		var actual UserLicenseAssignmentResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		// user_principal_name is Computed and only populated by a real API read.
		// Null or unknown means the state has not been refreshed from the API yet.
		if actual.UserPrincipalName.IsNull() || actual.UserPrincipalName.IsUnknown() || actual.UserPrincipalName.ValueString() == "" {
			return false
		}

		// id must be resolved to the composite userId_skuId value.
		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// user_id must match — confirms we are reading the correct user.
		if actual.UserId.ValueString() != expected.UserId.ValueString() {
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
