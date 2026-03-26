package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// crossTenantAccessPartnerSettingsConsistencyPredicate returns a consistency predicate for
// ReadWithRetry that verifies the cross-tenant access partner settings write has propagated
// across Microsoft Entra replicas before accepting the read as authoritative.
//
// The predicate compares all user-specified mutable fields from the expected plan against the
// state returned by the read. A stale replica returns the pre-write version, causing at least
// one field to differ and triggering a retry.
//
// is_service_provider and is_in_multi_tenant_organization are Computed (set by the API) and
// are not compared. hard_delete is a Terraform-only field and is not compared.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func crossTenantAccessPartnerSettingsConsistencyPredicate(expected *CrossTenantAccessPartnerSettingsResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual CrossTenantAccessPartnerSettingsResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if actual.TenantID.ValueString() != expected.TenantID.ValueString() {
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
		if !expected.AutomaticUserConsentSettings.IsNull() && !expected.AutomaticUserConsentSettings.IsUnknown() {
			if !actual.AutomaticUserConsentSettings.Equal(expected.AutomaticUserConsentSettings) {
				return false
			}
		}
		if !expected.TenantRestrictions.IsNull() && !expected.TenantRestrictions.IsUnknown() {
			if !actual.TenantRestrictions.Equal(expected.TenantRestrictions) {
				return false
			}
		}

		return true
	}
}
