package graphBetaAdminConsentRequestPolicy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// The tenant singleton is tested with its current writable configuration. The
// acceptance test must never enable consent workflow or replace tenant reviewers.
func TestAccResourceAdminConsentRequestPolicy_PreserveTenantConfiguration(t *testing.T) {
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	baseline, err := client.Policies().AdminConsentRequestPolicy().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, baseline)
	reviewers := make([]map[string]any, 0, len(baseline.GetReviewers()))
	for _, reviewer := range baseline.GetReviewers() {
		require.NotNil(t, reviewer.GetQuery(), "this test supports query-based reviewer fixtures only")
		require.NotNil(t, reviewer.GetQueryType())
		entry := map[string]any{"query": *reviewer.GetQuery(), "query_type": *reviewer.GetQueryType()}
		if reviewer.GetQueryRoot() != nil {
			entry["query_root"] = *reviewer.GetQueryRoot()
		}
		reviewers = append(reviewers, entry)
	}
	reviewerJSON, err := json.Marshal(reviewers)
	require.NoError(t, err)
	config := fmt.Sprintf(`resource "%s" "test" {
 is_enabled = %t
 notify_reviewers = %t
 reminders_enabled = %t
 request_duration_in_days = %d
 reviewers = %s
}`, resourceType, *baseline.GetIsEnabled(), *baseline.GetNotifyReviewers(), *baseline.GetRemindersEnabled(), *baseline.GetRequestDurationInDays(), reviewerJSON)
	address := resourceType + ".test"
	resource.Test(t, resource.TestCase{PreCheck: func() { mocks.TestAccPreCheck(t) }, ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: config, Check: resource.TestCheckResourceAttr(address, "id", "adminConsentRequestPolicy")},
		{ResourceName: address, ImportState: true, ImportStateVerify: true},
		{
			ResourceName:     address,
			ImportState:      true,
			GenerateConfig:   true,
			ImportStateKind:  resource.ImportBlockWithID,
			ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
		},
		{Config: config, PlanOnly: true},
	}})
	// Destroy must leave the service policy and its writable values intact.
	after, err := client.Policies().AdminConsentRequestPolicy().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, after)
	require.Equal(t, baseline.GetIsEnabled(), after.GetIsEnabled())
	require.Equal(t, baseline.GetNotifyReviewers(), after.GetNotifyReviewers())
	require.Equal(t, baseline.GetRemindersEnabled(), after.GetRemindersEnabled())
	require.Equal(t, baseline.GetRequestDurationInDays(), after.GetRequestDurationInDays())
	require.Equal(t, baseline.GetReviewers(), after.GetReviewers())
}
