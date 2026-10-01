package graphBetaAuthorizationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource sends only known values. False booleans and empty sets are explicit updates.
func constructResource(ctx context.Context, data *AuthorizationPolicyResourceModel) (graphmodels.AuthorizationPolicyable, error) {
	requestBody := graphmodels.NewAuthorizationPolicy()
	convert.FrameworkToGraphBool(data.AllowedToSignUpEmailBasedSubscriptions, requestBody.SetAllowedToSignUpEmailBasedSubscriptions)
	convert.FrameworkToGraphBool(data.AllowedToUseSSPR, requestBody.SetAllowedToUseSSPR)
	convert.FrameworkToGraphBool(data.AllowEmailVerifiedUsersToJoinOrganization, requestBody.SetAllowEmailVerifiedUsersToJoinOrganization)
	if err := convert.FrameworkToGraphEnum(data.AllowInvitesFrom, graphmodels.ParseAllowInvitesFrom, requestBody.SetAllowInvitesFrom); err != nil {
		return nil, fmt.Errorf("failed to set allow_invites_from: %w", err)
	}
	convert.FrameworkToGraphBool(data.AllowUserConsentForRiskyApps, requestBody.SetAllowUserConsentForRiskyApps)
	convert.FrameworkToGraphBool(data.BlockMsolPowerShell, requestBody.SetBlockMsolPowerShell)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	if err := convert.FrameworkToGraphStringSet(ctx, data.EnabledPreviewFeatures, requestBody.SetEnabledPreviewFeatures); err != nil {
		return nil, fmt.Errorf("failed to set enabled_preview_features: %w", err)
	}
	if err := convert.FrameworkToGraphUUID(data.GuestUserRoleId, requestBody.SetGuestUserRoleId); err != nil {
		return nil, fmt.Errorf("failed to set guest_user_role_id: %w", err)
	}
	if err := convert.FrameworkToGraphStringSet(ctx, data.PermissionGrantPolicyIdsAssignedToDefaultUserRole, requestBody.SetPermissionGrantPolicyIdsAssignedToDefaultUserRole); err != nil {
		return nil, fmt.Errorf("failed to set permission_grant_policy_ids_assigned_to_default_user_role: %w", err)
	}
	if !data.DefaultUserRolePermissions.IsNull() && !data.DefaultUserRolePermissions.IsUnknown() {
		var permissions DefaultUserRolePermissionsModel
		if diags := data.DefaultUserRolePermissions.As(ctx, &permissions, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, fmt.Errorf("failed to decode default user role permissions: %v", diags)
		}
		role := graphmodels.NewDefaultUserRolePermissions()
		// This live API property is not yet exposed by the pinned SDK model.
		if !permissions.AllowedToCreateAgentIdentityBlueprints.IsNull() && !permissions.AllowedToCreateAgentIdentityBlueprints.IsUnknown() {
			role.SetAdditionalData(map[string]any{
				"allowedToCreateAgentIdentityBlueprints": permissions.AllowedToCreateAgentIdentityBlueprints.ValueBool(),
			})
		}
		convert.FrameworkToGraphBool(permissions.AllowedToCreateApps, role.SetAllowedToCreateApps)
		convert.FrameworkToGraphBool(permissions.AllowedToCreateSecurityGroups, role.SetAllowedToCreateSecurityGroups)
		convert.FrameworkToGraphBool(permissions.AllowedToCreateTenants, role.SetAllowedToCreateTenants)
		convert.FrameworkToGraphBool(permissions.AllowedToReadBitlockerKeysForOwnedDevice, role.SetAllowedToReadBitlockerKeysForOwnedDevice)
		convert.FrameworkToGraphBool(permissions.AllowedToReadOtherUsers, role.SetAllowedToReadOtherUsers)
		requestBody.SetDefaultUserRolePermissions(role)
	}
	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return requestBody, nil
}
