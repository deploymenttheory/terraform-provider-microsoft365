package graphBetaCrossTenantAccessPolicy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaCrossTenantAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// resourceType is the Terraform resource type name from the resource package constant.
	resourceType = graphBetaCrossTenantAccessPolicy.ResourceName

	// testResource is the test helper implementation for cross-tenant access policy.
	testResource = graphBetaCrossTenantAccessPolicy.CrossTenantAccessPolicyTestResource{}
)

func loadAcceptanceConfig(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance config %s: %s", filename, err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceCrossTenantAccessPolicy_01_WithNoB2B tests that the singleton cross-tenant
// access policy can be brought under Terraform management with no allowed_cloud_endpoints,
// disabling all cross-cloud B2B collaboration.
//
// Note: CheckDestroy is intentionally nil. The crossTenantAccessPolicy is a singleton that always
// exists in the tenant and has no DELETE API endpoint. Destroy removes the resource from Terraform
// state only (restore_defaults_on_destroy defaults to false).
func TestAccResourceCrossTenantAccessPolicy_01_WithNoB2B(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
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
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting singleton cross-tenant access policy with no B2B endpoints")
				},
				Config: loadAcceptanceConfig("resource_cross_tenant_access_policy_with_no_b2b.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".with_no_b2b").ExistsInGraph(testResource),
					check.That(resourceType+".with_no_b2b").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".with_no_b2b").Key("display_name").Exists(),
					check.That(resourceType+".with_no_b2b").Key("allowed_cloud_endpoints.#").HasValue("0"),
					check.That(resourceType+".with_no_b2b").Key("restore_defaults_on_destroy").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing cross-tenant access policy")
				},
				ResourceName:      resourceType + ".with_no_b2b",
				ImportState:       true,
				ImportStateId:     "crossTenantAccessPolicy",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"restore_defaults_on_destroy",
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceCrossTenantAccessPolicy_02_WithAllowedCloudEndpoints tests configuring
// the policy with specific cloud endpoint collaborations.
func TestAccResourceCrossTenantAccessPolicy_02_WithAllowedCloudEndpoints(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
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
				PreConfig: func() {
					testlog.StepAction(resourceType, "Configuring cross-tenant access policy with allowed cloud endpoints")
				},
				Config: loadAcceptanceConfig("resource_cross_tenant_access_policy_with_endpoints.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".with_endpoints").ExistsInGraph(testResource),
					check.That(resourceType+".with_endpoints").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.#").HasValue("2"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("microsoftonline.us"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("partner.microsoftonline.cn"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing cross-tenant access policy with endpoints")
				},
				ResourceName:      resourceType + ".with_endpoints",
				ImportState:       true,
				ImportStateId:     "crossTenantAccessPolicy",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"restore_defaults_on_destroy",
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceCrossTenantAccessPolicy_03_RestoreDefaultsOnDestroy tests that the
// restore_defaults_on_destroy flag is accepted, stored correctly in state, and triggers a
// PATCH to restore service defaults when the resource is destroyed.
//
// This test intentionally uses allowed_cloud_endpoints = [] to avoid a known MS Graph API
// limitation: allowedCloudEndpoints behaves as a one-shot write per tenant session — once
// cloud endpoint values have been set and cleared, re-adding the same values returns 404.
// Since Test 02 already exercises the full endpoint lifecycle (set/verify/import), this test
// focuses solely on the restore_defaults_on_destroy behaviour.
func TestAccResourceCrossTenantAccessPolicy_03_RestoreDefaultsOnDestroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
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
				PreConfig: func() {
					testlog.StepAction(resourceType, "Configuring cross-tenant access policy with restore_defaults_on_destroy = true")
				},
				Config: loadAcceptanceConfig("resource_cross_tenant_access_policy_restore_defaults_on_destroy.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".restore_defaults").ExistsInGraph(testResource),
					check.That(resourceType+".restore_defaults").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".restore_defaults").Key("display_name").Exists(),
					check.That(resourceType+".restore_defaults").Key("restore_defaults_on_destroy").HasValue("true"),
				),
			},
		},
	})
}
