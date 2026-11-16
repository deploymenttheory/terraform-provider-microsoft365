package graphBetaIdentityAndAccessMobilityManagementPolicy_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphpolicies "github.com/microsoftgraph/msgraph-beta-sdk-go/policies"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	mobilityMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/mobility_management_policy/mocks"
)

// policySnapshot is the full writable state of one mobility management policy
type policySnapshot struct {
	AppliesTo                                 string
	ComplianceUrl                             string
	DiscoveryUrl                              string
	TermsOfUseUrl                             string
	IsMdmEnrollmentDuringRegistrationDisabled *bool
	IncludedGroupIDs                          []string
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// readPolicySnapshot reads the policy and its included groups with the acceptance test client
func readPolicySnapshot(ctx context.Context, policyType string) (policySnapshot, error) {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return policySnapshot{}, fmt.Errorf("create Graph client: %w", err)
	}

	var policy graphmodels.MobilityManagementPolicyable
	var groups graphmodels.GroupCollectionResponseable
	if policyType == "app" {
		policy, err = client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(mobilityMocks.IntunePolicyID).Get(ctx, nil)
		if err == nil {
			query := &graphpolicies.MobileAppManagementPoliciesItemIncludedGroupsRequestBuilderGetQueryParameters{Select: []string{"id"}}
			config := &graphpolicies.MobileAppManagementPoliciesItemIncludedGroupsRequestBuilderGetRequestConfiguration{QueryParameters: query}
			groups, err = client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(mobilityMocks.IntunePolicyID).IncludedGroups().Get(ctx, config)
		}
	} else {
		policy, err = client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(mobilityMocks.IntunePolicyID).Get(ctx, nil)
		if err == nil {
			query := &graphpolicies.MobileDeviceManagementPoliciesItemIncludedGroupsRequestBuilderGetQueryParameters{Select: []string{"id"}}
			config := &graphpolicies.MobileDeviceManagementPoliciesItemIncludedGroupsRequestBuilderGetRequestConfiguration{QueryParameters: query}
			groups, err = client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(mobilityMocks.IntunePolicyID).IncludedGroups().Get(ctx, config)
		}
	}
	if err != nil {
		return policySnapshot{}, fmt.Errorf("read %s mobility management policy: %w", policyType, err)
	}

	snapshot := policySnapshot{
		ComplianceUrl:    stringValue(policy.GetComplianceUrl()),
		DiscoveryUrl:     stringValue(policy.GetDiscoveryUrl()),
		TermsOfUseUrl:    stringValue(policy.GetTermsOfUseUrl()),
		IncludedGroupIDs: []string{},
	}
	if policy.GetAppliesTo() != nil {
		snapshot.AppliesTo = policy.GetAppliesTo().String()
	}
	if device, ok := policy.(graphmodels.MobileDeviceManagementPolicyable); ok {
		snapshot.IsMdmEnrollmentDuringRegistrationDisabled = device.GetIsMdmEnrollmentDuringRegistrationDisabled()
	}
	if groups != nil {
		for _, group := range groups.GetValue() {
			if group != nil && group.GetId() != nil {
				snapshot.IncludedGroupIDs = append(snapshot.IncludedGroupIDs, *group.GetId())
			}
		}
		slices.Sort(snapshot.IncludedGroupIDs)
	}
	return snapshot, nil
}

// snapshotPreCheck records the policy before the test and fails the test if it differs afterwards.
// It is used by tests that must leave the policy exactly as they found it.
func snapshotPreCheck(t *testing.T, policyType string) policySnapshot {
	t.Helper()
	mocks.TestAccPreCheck(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	original, err := readPolicySnapshot(ctx, policyType)
	require.NoError(t, err)
	t.Logf("original %s mobility management policy: %+v", policyType, original)

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		consecutive := 0
		for {
			current, err := readPolicySnapshot(cleanupCtx, policyType)
			if err == nil && fmt.Sprintf("%+v", current) == fmt.Sprintf("%+v", original) {
				consecutive++
				if consecutive >= 3 {
					return
				}
			} else {
				consecutive = 0
			}
			select {
			case <-cleanupCtx.Done():
				t.Errorf("%s mobility management policy did not return to its original state %+v after acceptance testing (last error: %v)", policyType, original, err)
				return
			case <-time.After(2 * time.Second):
			}
		}
	})

	return original
}

// appScopePreCheck snapshots the MAM policy and skips the test unless its scope is none.
// Changing a none scope to the test group and back affects no real user; any other scope would
// change automatic enrollment for real users while the test runs.
func appScopePreCheck(t *testing.T) {
	t.Helper()
	original := snapshotPreCheck(t, "app")
	if original.AppliesTo != "none" || len(original.IncludedGroupIDs) != 0 {
		t.Skipf("skipping: the MAM policy scope is %q with %d groups; the test only changes a policy whose scope is none", original.AppliesTo, len(original.IncludedGroupIDs))
	}

	// Safety net that runs before the snapshot comparison (cleanups run last in, first out): if a failed step
	// left the scope changed, PATCH it back to none, which also removes the test group reference.
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		current, err := readPolicySnapshot(ctx, "app")
		if err == nil && current.AppliesTo == "none" && len(current.IncludedGroupIDs) == 0 {
			return
		}
		client, err := acceptance.TestGraphClient()
		if err != nil {
			t.Errorf("create Graph client to restore the MAM scope: %v", err)
			return
		}
		body := graphmodels.NewMobileAppManagementPolicy()
		scope := graphmodels.NONE_POLICYSCOPE
		body.SetAppliesTo(&scope)
		if _, err := client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(mobilityMocks.IntunePolicyID).Patch(ctx, body, nil); err != nil {
			t.Errorf("restore the MAM scope to none: %v", err)
		}
	})
}
