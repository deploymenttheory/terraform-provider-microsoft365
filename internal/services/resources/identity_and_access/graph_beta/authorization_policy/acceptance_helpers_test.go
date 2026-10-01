package graphBetaAuthorizationPolicy_test

import (
	"context"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// authorizationPolicyPreCheck preserves the tenant settings for acceptance-test cleanup.
// Graph ignores null resets for risky-app consent, so an originally-null value cannot
// be safely round-tripped by tests that must explicitly configure all Boolean fields.
func authorizationPolicyPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	response, err := client.Policies().AuthorizationPolicy().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Len(t, response.GetValue(), 1)
	original := response.GetValue()[0]
	require.NotNil(t, original)
	if original.GetAllowUserConsentForRiskyApps() == nil {
		t.Skip("Graph cannot restore allowUserConsentForRiskyApps to null; use a test tenant with an explicitly configured Boolean value")
	}
	expected := authorizationPolicySnapshot(original)
	original.SetId(nil)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		_, err := client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId("authorizationPolicy").Patch(cleanupCtx, original, nil)
		if err != nil {
			t.Errorf("restore original authorization policy: %v", err)
			return
		}
		for {
			response, err := client.Policies().AuthorizationPolicy().Get(cleanupCtx, nil)
			if err == nil && response != nil && len(response.GetValue()) == 1 && reflect.DeepEqual(expected, authorizationPolicySnapshot(response.GetValue()[0])) {
				return
			}
			select {
			case <-cleanupCtx.Done():
				t.Error("authorization policy did not return to its original settings after acceptance testing")
				return
			case <-time.After(2 * time.Second):
			}
		}
	})
}

// authorizationPolicySnapshot compares all managed values while ignoring set ordering.
func authorizationPolicySnapshot(policy graphmodels.AuthorizationPolicyable) map[string]any {
	if policy == nil {
		return nil
	}
	preview := slices.Clone(policy.GetEnabledPreviewFeatures())
	grants := slices.Clone(policy.GetPermissionGrantPolicyIdsAssignedToDefaultUserRole())
	slices.Sort(preview)
	slices.Sort(grants)
	snapshot := map[string]any{
		"allowInvitesFrom":                                  policy.GetAllowInvitesFrom(),
		"allowedToSignUpEmailBasedSubscriptions":            policy.GetAllowedToSignUpEmailBasedSubscriptions(),
		"allowedToUseSSPR":                                  policy.GetAllowedToUseSSPR(),
		"allowEmailVerifiedUsersToJoinOrganization":         policy.GetAllowEmailVerifiedUsersToJoinOrganization(),
		"allowUserConsentForRiskyApps":                      policy.GetAllowUserConsentForRiskyApps(),
		"blockMsolPowerShell":                               policy.GetBlockMsolPowerShell(),
		"description":                                       policy.GetDescription(),
		"displayName":                                       policy.GetDisplayName(),
		"enabledPreviewFeatures":                            preview,
		"guestUserRoleId":                                   policy.GetGuestUserRoleId(),
		"permissionGrantPolicyIdsAssignedToDefaultUserRole": grants,
	}
	if role := policy.GetDefaultUserRolePermissions(); role != nil {
		snapshot["defaultUserRolePermissions"] = map[string]any{
			"allowedToCreateApps":                      role.GetAllowedToCreateApps(),
			"allowedToCreateAgentIdentityBlueprints":   role.GetAdditionalData()["allowedToCreateAgentIdentityBlueprints"],
			"allowedToCreateSecurityGroups":            role.GetAllowedToCreateSecurityGroups(),
			"allowedToCreateTenants":                   role.GetAllowedToCreateTenants(),
			"allowedToReadBitlockerKeysForOwnedDevice": role.GetAllowedToReadBitlockerKeysForOwnedDevice(),
			"allowedToReadOtherUsers":                  role.GetAllowedToReadOtherUsers(),
		}
	}
	return snapshot
}
