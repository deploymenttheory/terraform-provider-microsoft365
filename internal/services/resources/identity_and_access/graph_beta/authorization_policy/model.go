// REF: https://learn.microsoft.com/en-us/graph/api/resources/authorizationpolicy?view=graph-rest-beta
package graphBetaAuthorizationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthorizationPolicyResourceModel represents the tenant-wide authorization policy.
type AuthorizationPolicyResourceModel struct {
	ID                                                types.String   `tfsdk:"id"`
	AllowedToSignUpEmailBasedSubscriptions            types.Bool     `tfsdk:"allowed_to_sign_up_email_based_subscriptions"`
	AllowedToUseSSPR                                  types.Bool     `tfsdk:"allowed_to_use_sspr"`
	AllowEmailVerifiedUsersToJoinOrganization         types.Bool     `tfsdk:"allow_email_verified_users_to_join_organization"`
	AllowInvitesFrom                                  types.String   `tfsdk:"allow_invites_from"`
	AllowUserConsentForRiskyApps                      types.Bool     `tfsdk:"allow_user_consent_for_risky_apps"`
	BlockMsolPowerShell                               types.Bool     `tfsdk:"block_msol_powershell"`
	Description                                       types.String   `tfsdk:"description"`
	DisplayName                                       types.String   `tfsdk:"display_name"`
	EnabledPreviewFeatures                            types.Set      `tfsdk:"enabled_preview_features"`
	GuestUserRoleId                                   types.String   `tfsdk:"guest_user_role_id"`
	PermissionGrantPolicyIdsAssignedToDefaultUserRole types.Set      `tfsdk:"permission_grant_policy_ids_assigned_to_default_user_role"`
	DefaultUserRolePermissions                        types.Object   `tfsdk:"default_user_role_permissions"`
	Timeouts                                          timeouts.Value `tfsdk:"timeouts"`
}

// DefaultUserRolePermissionsModel represents customizable permissions for the default user role.
type DefaultUserRolePermissionsModel struct {
	AllowedToCreateAgentIdentityBlueprints   types.Bool `tfsdk:"allowed_to_create_agent_identity_blueprints"`
	AllowedToCreateApps                      types.Bool `tfsdk:"allowed_to_create_apps"`
	AllowedToCreateSecurityGroups            types.Bool `tfsdk:"allowed_to_create_security_groups"`
	AllowedToCreateTenants                   types.Bool `tfsdk:"allowed_to_create_tenants"`
	AllowedToReadBitlockerKeysForOwnedDevice types.Bool `tfsdk:"allowed_to_read_bitlocker_keys_for_owned_device"`
	AllowedToReadOtherUsers                  types.Bool `tfsdk:"allowed_to_read_other_users"`
}

func defaultUserRolePermissionsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_to_create_agent_identity_blueprints":     types.BoolType,
		"allowed_to_create_apps":                          types.BoolType,
		"allowed_to_create_security_groups":               types.BoolType,
		"allowed_to_create_tenants":                       types.BoolType,
		"allowed_to_read_bitlocker_keys_for_owned_device": types.BoolType,
		"allowed_to_read_other_users":                     types.BoolType,
	}
}
