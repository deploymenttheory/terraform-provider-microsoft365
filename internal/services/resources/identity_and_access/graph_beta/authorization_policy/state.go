package graphBetaAuthorizationPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var (
	errInvalidPolicyCollection = errors.New("expected exactly one authorization policy in the response")
	errMissingSingleton        = errors.New("response is missing the authorizationPolicy singleton")
)

// mapRemoteState validates the collection returned by the beta singleton GET endpoint.
func mapRemoteState(ctx context.Context, data *AuthorizationPolicyResourceModel, response graphmodels.AuthorizationPolicyCollectionResponseable) error {
	if response == nil || len(response.GetValue()) != 1 {
		return errInvalidPolicyCollection
	}
	remote := response.GetValue()[0]
	if remote == nil || remote.GetId() == nil || *remote.GetId() != singletonID {
		return errMissingSingleton
	}
	data.ID = types.StringValue(singletonID)
	data.AllowedToSignUpEmailBasedSubscriptions = convert.GraphToFrameworkBool(remote.GetAllowedToSignUpEmailBasedSubscriptions())
	data.AllowedToUseSSPR = convert.GraphToFrameworkBool(remote.GetAllowedToUseSSPR())
	data.AllowEmailVerifiedUsersToJoinOrganization = convert.GraphToFrameworkBool(remote.GetAllowEmailVerifiedUsersToJoinOrganization())
	data.AllowInvitesFrom = convert.GraphToFrameworkEnum(remote.GetAllowInvitesFrom())
	data.AllowUserConsentForRiskyApps = convert.GraphToFrameworkBool(remote.GetAllowUserConsentForRiskyApps())
	data.BlockMsolPowerShell = convert.GraphToFrameworkBool(remote.GetBlockMsolPowerShell())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.EnabledPreviewFeatures = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remote.GetEnabledPreviewFeatures())
	data.GuestUserRoleId = convert.GraphToFrameworkUUID(remote.GetGuestUserRoleId())
	data.PermissionGrantPolicyIdsAssignedToDefaultUserRole = mapPermissionGrantPolicyIDs(ctx, remote.GetPermissionGrantPolicyIdsAssignedToDefaultUserRole(), data.PermissionGrantPolicyIdsAssignedToDefaultUserRole)
	data.DefaultUserRolePermissions = types.ObjectNull(defaultUserRolePermissionsAttributeTypes())
	if permissions := remote.GetDefaultUserRolePermissions(); permissions != nil {
		role := DefaultUserRolePermissionsModel{
			AllowedToCreateApps:                      convert.GraphToFrameworkBool(permissions.GetAllowedToCreateApps()),
			AllowedToCreateSecurityGroups:            convert.GraphToFrameworkBool(permissions.GetAllowedToCreateSecurityGroups()),
			AllowedToCreateTenants:                   convert.GraphToFrameworkBool(permissions.GetAllowedToCreateTenants()),
			AllowedToReadBitlockerKeysForOwnedDevice: convert.GraphToFrameworkBool(permissions.GetAllowedToReadBitlockerKeysForOwnedDevice()),
			AllowedToReadOtherUsers:                  convert.GraphToFrameworkBool(permissions.GetAllowedToReadOtherUsers()),
		}
		// The pinned SDK stores this newer permission in AdditionalData.
		switch value := permissions.GetAdditionalData()["allowedToCreateAgentIdentityBlueprints"].(type) {
		case *bool:
			role.AllowedToCreateAgentIdentityBlueprints = convert.GraphToFrameworkBool(value)
		case bool:
			role.AllowedToCreateAgentIdentityBlueprints = types.BoolValue(value)
		}
		value, diags := types.ObjectValueFrom(ctx, defaultUserRolePermissionsAttributeTypes(), role)
		if diags.HasError() {
			return fmt.Errorf("failed to map default user role permissions: %v", diags)
		}
		data.DefaultUserRolePermissions = value
	}
	return nil
}
