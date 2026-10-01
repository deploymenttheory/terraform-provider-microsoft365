package graphBetaAuthorizationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// consistencyPredicate waits until a read reflects every known planned policy setting.
func consistencyPredicate(expected *AuthorizationPolicyResourceModel) func(context.Context, tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual AuthorizationPolicyResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}
		if actual.ID.ValueString() != singletonID {
			return false
		}
		for _, pair := range [][2]attr.Value{
			{expected.AllowedToSignUpEmailBasedSubscriptions, actual.AllowedToSignUpEmailBasedSubscriptions},
			{expected.AllowedToUseSSPR, actual.AllowedToUseSSPR},
			{expected.AllowEmailVerifiedUsersToJoinOrganization, actual.AllowEmailVerifiedUsersToJoinOrganization},
			{expected.AllowInvitesFrom, actual.AllowInvitesFrom},
			{expected.AllowUserConsentForRiskyApps, actual.AllowUserConsentForRiskyApps},
			{expected.BlockMsolPowerShell, actual.BlockMsolPowerShell},
			{expected.Description, actual.Description},
			{expected.DisplayName, actual.DisplayName},
			{expected.EnabledPreviewFeatures, actual.EnabledPreviewFeatures},
			{expected.GuestUserRoleId, actual.GuestUserRoleId},
			{expected.PermissionGrantPolicyIdsAssignedToDefaultUserRole, actual.PermissionGrantPolicyIdsAssignedToDefaultUserRole},
			{expected.DefaultUserRolePermissions, actual.DefaultUserRolePermissions},
		} {
			if !matchesKnownValue(pair[0], pair[1]) {
				return false
			}
		}
		return true
	}
}

func matchesKnownValue(expected, actual attr.Value) bool {
	if expected.IsNull() || expected.IsUnknown() {
		return true
	}
	if actual.IsNull() || actual.IsUnknown() {
		return false
	}
	if expectedObject, ok := expected.(types.Object); ok {
		actualObject, ok := actual.(types.Object)
		if !ok {
			return false
		}
		for name, value := range expectedObject.Attributes() {
			if !matchesKnownValue(value, actualObject.Attributes()[name]) {
				return false
			}
		}
		return true
	}
	return expected.Equal(actual)
}
