package graphBetaCrossTenantAccessDefaultSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// crossTenantAccessDefaultSettingsConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the cross-tenant access default settings write has propagated
// across Microsoft Entra replicas before accepting the read as authoritative.
//
// This is a singleton resource. The predicate compares all user-specified mutable fields from
// the expected plan against the state returned by the read. A stale replica returns the
// pre-write version, causing at least one field to differ and triggering a retry.
//
// is_service_default is Computed (set by the API) and is checked to confirm the replica
// returned the full settings object. restore_defaults_on_destroy is a Terraform-only field
// and is not compared.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func crossTenantAccessDefaultSettingsConsistencyPredicate(expected *CrossTenantAccessDefaultSettingsResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual CrossTenantAccessDefaultSettingsResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		// is_service_default is Computed — only set by the API, never from plan.
		// Its presence confirms the replica returned the full settings object.
		// When any custom settings are present in the plan, wait for is_service_default
		// to become false — the API sets it false asynchronously after any PATCH, so
		// false is the definitive signal that replication has completed across all replicas.
		if actual.IsServiceDefault.IsNull() || actual.IsServiceDefault.IsUnknown() {
			return false
		}

		hasCustomSettings := !expected.B2bCollaborationInbound.IsNull() ||
			!expected.B2bCollaborationOutbound.IsNull() ||
			!expected.B2bDirectConnectInbound.IsNull() ||
			!expected.B2bDirectConnectOutbound.IsNull() ||
			!expected.InboundTrust.IsNull() ||
			!expected.InvitationRedemptionIdentityProviderConfiguration.IsNull() ||
			!expected.TenantRestrictions.IsNull() ||
			!expected.AutomaticUserConsentSettings.IsNull()

		if hasCustomSettings && actual.IsServiceDefault.ValueBool() {
			return false
		}

		// Compare all user-specified mutable fields. Optional+Computed object fields are only
		// compared when explicitly set in the plan to avoid false failures from API defaults.
		if !expected.B2bCollaborationInbound.IsNull() && !expected.B2bCollaborationInbound.IsUnknown() {
			if !actual.B2bCollaborationInbound.Equal(expected.B2bCollaborationInbound) {
				return false
			}
		}
		if !expected.B2bCollaborationOutbound.IsNull() && !expected.B2bCollaborationOutbound.IsUnknown() {
			if !actual.B2bCollaborationOutbound.Equal(expected.B2bCollaborationOutbound) {
				return false
			}
		}
		if !expected.B2bDirectConnectInbound.IsNull() && !expected.B2bDirectConnectInbound.IsUnknown() {
			if !actual.B2bDirectConnectInbound.Equal(expected.B2bDirectConnectInbound) {
				return false
			}
		}
		if !expected.B2bDirectConnectOutbound.IsNull() && !expected.B2bDirectConnectOutbound.IsUnknown() {
			if !actual.B2bDirectConnectOutbound.Equal(expected.B2bDirectConnectOutbound) {
				return false
			}
		}
		if !expected.InboundTrust.IsNull() && !expected.InboundTrust.IsUnknown() {
			if !actual.InboundTrust.Equal(expected.InboundTrust) {
				return false
			}
		}
		if !expected.InvitationRedemptionIdentityProviderConfiguration.IsNull() && !expected.InvitationRedemptionIdentityProviderConfiguration.IsUnknown() {
			if !actual.InvitationRedemptionIdentityProviderConfiguration.Equal(expected.InvitationRedemptionIdentityProviderConfiguration) {
				return false
			}
		}
		if !expected.TenantRestrictions.IsNull() && !expected.TenantRestrictions.IsUnknown() {
			if !actual.TenantRestrictions.Equal(expected.TenantRestrictions) {
				return false
			}
		}
		if !expected.AutomaticUserConsentSettings.IsNull() && !expected.AutomaticUserConsentSettings.IsUnknown() {
			if !actual.AutomaticUserConsentSettings.Equal(expected.AutomaticUserConsentSettings) {
				return false
			}
		}

		return true
	}
}
