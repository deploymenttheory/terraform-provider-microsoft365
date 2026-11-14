package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// deviceRegistrationPolicyConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the PUT has propagated before accepting the read as authoritative.
//
// expected is the plan mapped through the request body that was sent, so it holds every written
// value (including the values copied from the remote policy for unconfigured attributes). The
// device registration policy has no modified timestamp, so the predicate compares the written values.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func deviceRegistrationPolicyConsistencyPredicate(expected DeviceRegistrationPolicyResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual DeviceRegistrationPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		return actual.UserDeviceQuota.Equal(expected.UserDeviceQuota) &&
			actual.MultiFactorAuthConfiguration.Equal(expected.MultiFactorAuthConfiguration) &&
			actual.LocalAdminPassword.Equal(expected.LocalAdminPassword) &&
			actual.AzureADRegistration.Equal(expected.AzureADRegistration) &&
			actual.AzureADJoin.Equal(expected.AzureADJoin)
	}
}
