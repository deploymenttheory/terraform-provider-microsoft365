package graphBetaAuthorizationPolicy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
)

// authorizationPolicyPreCheck preserves the tenant settings for acceptance-test cleanup.
// Graph ignores null resets for risky-app consent, so the acceptance configs omit it to
// keep an originally-null value restorable.
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
	expected := authorizationPolicySnapshot(original)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		restore := authorizationPolicyRestoreBody(original)
		if err := patchAuthorizationPolicyForCleanup(cleanupCtx, client, restore); err != nil {
			t.Errorf("restore original authorization policy: %v", err)
			return
		}
		primaryExpected := maps.Clone(expected)
		delete(primaryExpected, "allowInvitesFrom")
		primaryMatches, consecutive := 0, 0
		invitationsWritten := false
		var observed map[string]any
		for {
			response, err := client.Policies().AuthorizationPolicy().Get(cleanupCtx, nil)
			if err == nil && response != nil && len(response.GetValue()) == 1 {
				observed = authorizationPolicySnapshot(response.GetValue()[0])
				if !invitationsWritten {
					primaryObserved := maps.Clone(observed)
					delete(primaryObserved, "allowInvitesFrom")
					if reflect.DeepEqual(primaryExpected, primaryObserved) {
						primaryMatches++
					} else {
						primaryMatches = 0
					}
					if primaryMatches >= 3 {
						if original.GetAllowInvitesFrom() != nil {
							invitations := graphmodels.NewAuthorizationPolicy()
							invitations.SetAllowInvitesFrom(original.GetAllowInvitesFrom())
							if err := patchAuthorizationPolicyForCleanup(cleanupCtx, client, invitations); err != nil {
								t.Errorf("restore original invitation settings: %v", err)
								return
							}
						}
						invitationsWritten = true
					}
				} else if reflect.DeepEqual(expected, observed) {
					consecutive++
					if consecutive >= 3 {
						return
					}
				} else {
					consecutive = 0
				}
			} else {
				primaryMatches, consecutive = 0, 0
			}
			select {
			case <-cleanupCtx.Done():
				expectedJSON, _ := json.Marshal(expected)
				observedJSON, _ := json.Marshal(observed)
				t.Errorf("authorization policy did not return to its original settings: expected %s; last read %s; error %v", expectedJSON, observedJSON, err)
				return
			case <-time.After(2 * time.Second):
			}
		}
	})
}

// authorizationPolicyRestoreBody marks every original value for serialization.
// Hydrated SDK models track changed properties, so reusing a GET response can omit
// the unchanged values that cleanup must explicitly send back to Graph.
func authorizationPolicyRestoreBody(
	original graphmodels.AuthorizationPolicyable,
) graphmodels.AuthorizationPolicyable {
	body := graphmodels.NewAuthorizationPolicy()
	body.SetAllowInvitesFrom(original.GetAllowInvitesFrom())
	body.SetAllowedToSignUpEmailBasedSubscriptions(
		original.GetAllowedToSignUpEmailBasedSubscriptions(),
	)
	body.SetAllowedToUseSSPR(original.GetAllowedToUseSSPR())
	body.SetAllowEmailVerifiedUsersToJoinOrganization(
		original.GetAllowEmailVerifiedUsersToJoinOrganization(),
	)
	body.SetAllowUserConsentForRiskyApps(original.GetAllowUserConsentForRiskyApps())
	body.SetBlockMsolPowerShell(original.GetBlockMsolPowerShell())
	body.SetDescription(original.GetDescription())
	body.SetDisplayName(original.GetDisplayName())
	body.SetEnabledPreviewFeatures(original.GetEnabledPreviewFeatures())
	body.SetGuestUserRoleId(original.GetGuestUserRoleId())
	body.SetPermissionGrantPolicyIdsAssignedToDefaultUserRole(
		original.GetPermissionGrantPolicyIdsAssignedToDefaultUserRole(),
	)
	if role := original.GetDefaultUserRolePermissions(); role != nil {
		permissions := graphmodels.NewDefaultUserRolePermissions()
		permissions.SetAllowedToCreateApps(role.GetAllowedToCreateApps())
		permissions.SetAllowedToCreateSecurityGroups(role.GetAllowedToCreateSecurityGroups())
		permissions.SetAllowedToCreateTenants(role.GetAllowedToCreateTenants())
		permissions.SetAllowedToReadBitlockerKeysForOwnedDevice(
			role.GetAllowedToReadBitlockerKeysForOwnedDevice(),
		)
		permissions.SetAllowedToReadOtherUsers(role.GetAllowedToReadOtherUsers())
		permissions.SetAdditionalData(role.GetAdditionalData())
		body.SetDefaultUserRolePermissions(permissions)
	}
	return body
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

// checkRiskyAppConsentMatchesGraph verifies the omitted attribute mirrors Graph, including null.
func checkRiskyAppConsentMatchesGraph(address string) func(*terraform.State) error {
	return func(state *terraform.State) error {
		resource, ok := state.RootModule().Resources[address]
		if !ok {
			return fmt.Errorf("%s not found in state", address)
		}
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		response, err := client.Policies().AuthorizationPolicy().Get(ctx, nil)
		if err != nil {
			return err
		}
		if response == nil || len(response.GetValue()) != 1 {
			return fmt.Errorf("expected one authorization policy")
		}
		actual, set := resource.Primary.Attributes["allow_user_consent_for_risky_apps"]
		remote := response.GetValue()[0].GetAllowUserConsentForRiskyApps()
		if remote == nil {
			if set {
				return fmt.Errorf("expected null allow_user_consent_for_risky_apps, got %q", actual)
			}
			return nil
		}
		if want := fmt.Sprint(*remote); actual != want {
			return fmt.Errorf("expected allow_user_consent_for_risky_apps %q, got %q", want, actual)
		}
		return nil
	}
}

func patchAuthorizationPolicyForCleanup(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	body graphmodels.AuthorizationPolicyable,
) error {
	return crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
		_, err := client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId("authorizationPolicy").Patch(ctx, body, nil)
		return err
	}, crud.PreconditionRetryOptions{
		MaxRetries: 15, RetryInterval: 2 * time.Second, StatusCodes: []int{404},
		ResourceTypeName: "authorization_policy", Operation: "restore",
	})
}
