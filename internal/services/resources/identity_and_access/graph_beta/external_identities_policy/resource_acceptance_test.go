package graphBetaExternalIdentitiesPolicy_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/external_identities_policy"
)

var (
	resourceType = policy.ResourceName
	testResource = policy.ExternalIdentitiesPolicyTestResource{}
)

// loadAcceptanceTestTerraform loads a configuration with the shared provider setup.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance config %s: %s", filename, err))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceExternalIdentitiesPolicy_01_Lifecycle exercises the singleton and restores its original setting.
// CheckDestroy is nil because destroy must leave the policy in Microsoft Graph.
func TestAccResourceExternalIdentitiesPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { externalIdentitiesPolicyPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Enabling external users leaving the tenant") },
				Config:    loadAcceptanceTestTerraform("resource_initial.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("id").
						HasValue("externalIdentityPolicy"),
					check.That(resourceType+".test").
						Key("allow_external_identities_to_leave").
						HasValue("true"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Updating the reserved data-removal setting") },
				Config:    loadAcceptanceTestTerraform("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").
						Key("allow_external_identities_to_leave").
						HasValue("true"),
					check.That(resourceType+".test").
						Key("allow_deleted_identities_data_removal").
						HasValue("true"),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "externalIdentityPolicy",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_updated.tf"), PlanOnly: true},
		},
	})
}

// externalIdentitiesPolicyPreCheck captures the original writable setting before any mutations.
func externalIdentitiesPolicyPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	original, err := client.Policies().ExternalIdentitiesPolicy().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, original)
	require.NotNil(
		t,
		original.GetAllowExternalIdentitiesToLeave(),
		"a Boolean baseline is required for reversible testing",
	)
	require.NotNil(
		t,
		original.GetAllowDeletedIdentitiesDataRemoval(),
		"a Boolean baseline is required for reversible testing",
	)
	expected := externalIdentitiesPolicySnapshot(original)
	enabled := *original.GetAllowExternalIdentitiesToLeave()
	removal := *original.GetAllowDeletedIdentitiesDataRemoval()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cleanupCancel()
		// Build a fresh SDK object so all restore fields are serialized as changes.
		body := graphmodels.NewExternalIdentitiesPolicy()
		body.SetAllowExternalIdentitiesToLeave(&enabled)
		body.SetAllowDeletedIdentitiesDataRemoval(&removal)
		_, err := client.Policies().ExternalIdentitiesPolicy().Patch(cleanupCtx, body, nil)
		if err != nil {
			t.Errorf("restore original external identities policy: %v", err)
			return
		}
		consecutive := 0
		for {
			current, err := client.Policies().ExternalIdentitiesPolicy().Get(cleanupCtx, nil)
			if err == nil &&
				reflect.DeepEqual(expected, externalIdentitiesPolicySnapshot(current)) {
				consecutive++
				if consecutive >= 3 {
					return
				}
			} else {
				consecutive = 0
			}
			select {
			case <-cleanupCtx.Done():
				t.Error("original external identities policy did not verify after cleanup")
				return
			case <-time.After(2 * time.Second):
			}
		}
	})
}

func externalIdentitiesPolicySnapshot(
	value graphmodels.ExternalIdentitiesPolicyable,
) map[string]any {
	if value == nil || value.GetAllowExternalIdentitiesToLeave() == nil {
		return nil
	}
	return map[string]any{
		"id":                                value.GetId(),
		"displayName":                       value.GetDisplayName(),
		"description":                       value.GetDescription(),
		"allowExternalIdentitiesToLeave":    value.GetAllowExternalIdentitiesToLeave(),
		"allowDeletedIdentitiesDataRemoval": value.GetAllowDeletedIdentitiesDataRemoval(),
	}
}
