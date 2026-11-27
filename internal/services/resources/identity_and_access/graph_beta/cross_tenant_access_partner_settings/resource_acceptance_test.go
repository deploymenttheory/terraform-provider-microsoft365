package graphBetaCrossTenantAccessPartnerSettings_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	graphBetaCrossTenantAccessPartnerSettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_partner_settings"
)

const testTenantID = "a22ff489-2ea9-48de-8d58-fa130b532d5d"

var (
	resourceType = graphBetaCrossTenantAccessPartnerSettings.ResourceName
	testResource = graphBetaCrossTenantAccessPartnerSettings.CrossTenantAccessPartnerSettingsTestResource{}
)

func loadAcceptanceConfig(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance config %s: %s", filename, err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceCrossTenantAccessPartnerSettings_01_Minimal tests that a partner-specific
// cross-tenant access configuration can be created with the minimal valid configuration
// (outbound B2B collaboration only).
func TestAccResourceCrossTenantAccessPartnerSettings_01_Minimal(t *testing.T) {
	r := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		CheckDestroy: testAccCheckCrossTenantAccessPartnerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Deploying minimal cross-tenant access partner settings")
				},
				Config: loadAcceptanceConfig("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("tenant_id").HasValue(testTenantID),
					check.That(r).Key("hard_delete").HasValue("true"),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal cross-tenant access partner settings with hard_delete=true")
				},
				ResourceName:      r,
				ImportState:       true,
				ImportStateId:     testTenantID + ":hard_delete=true",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceCrossTenantAccessPartnerSettings_02_Maximal tests deploying the full
// cross-tenant access partner settings configuration with all blocks populated.
func TestAccResourceCrossTenantAccessPartnerSettings_02_Maximal(t *testing.T) {
	r := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		CheckDestroy: testAccCheckCrossTenantAccessPartnerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Deploying maximal cross-tenant access partner settings with all blocks")
				},
				Config: loadAcceptanceConfig("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("tenant_id").HasValue(testTenantID),
					check.That(r).Key("hard_delete").HasValue("true"),
					check.That(r).Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_inbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_direct_connect_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_inbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(r).Key("inbound_trust.is_compliant_device_accepted").HasValue("true"),
					check.That(r).Key("inbound_trust.is_hybrid_azure_ad_joined_device_accepted").HasValue("true"),
					check.That(r).Key("tenant_restrictions.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("tenant_restrictions.applications.access_type").HasValue("blocked"),
					check.That(r).Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(r).Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal cross-tenant access partner settings with hard_delete=true")
				},
				ResourceName:      r,
				ImportState:       true,
				ImportStateId:     testTenantID + ":hard_delete=true",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestAccResourceCrossTenantAccessPartnerSettings_03_UpdateMinimalToMaximal tests the PATCH
// update path by first deploying the minimal configuration then expanding it to the full
// maximal configuration in a second step.
func TestAccResourceCrossTenantAccessPartnerSettings_03_UpdateMinimalToMaximal(t *testing.T) {
	r := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		CheckDestroy: testAccCheckCrossTenantAccessPartnerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Deploying minimal configuration")
				},
				Config: loadAcceptanceConfig("resource_03_update_minimal_to_maximal_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Expanding to maximal configuration via PATCH update")
				},
				Config: loadAcceptanceConfig("resource_03_update_minimal_to_maximal_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_inbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_direct_connect_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_inbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(r).Key("tenant_restrictions.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(r).Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
		},
	})
}

// TestAccResourceCrossTenantAccessPartnerSettings_04_UpdateMaximalToMinimal tests the PATCH
// update path in the reverse direction: deploying the full maximal configuration first, then
// reducing it to the minimal configuration in a second step.
func TestAccResourceCrossTenantAccessPartnerSettings_04_UpdateMaximalToMinimal(t *testing.T) {
	r := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		CheckDestroy: testAccCheckCrossTenantAccessPartnerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Deploying maximal configuration with all blocks")
				},
				Config: loadAcceptanceConfig("resource_04_update_maximal_to_minimal_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_inbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_direct_connect_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_inbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("b2b_direct_connect_outbound.applications.access_type").HasValue("blocked"),
					check.That(r).Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(r).Key("tenant_restrictions.users_and_groups.access_type").HasValue("blocked"),
					check.That(r).Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(r).Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Reducing to minimal configuration via PATCH update")
				},
				Config: loadAcceptanceConfig("resource_04_update_maximal_to_minimal_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("cross-tenant access partner settings", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(r).ExistsInGraph(testResource),
					check.That(r).Key("id").HasValue(testTenantID),
					check.That(r).Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(r).Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
				),
			},
		},
	})
}

// testAccCheckCrossTenantAccessPartnerSettingsDestroy checks both actual stores.
// Replica reads can briefly reappear after a 404, so require stable absence.
func testAccCheckCrossTenantAccessPartnerSettingsDestroy(s *terraform.State) error {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}
		id := rs.Primary.Attributes["tenant_id"]
		absentReads := 0
		for absentReads < 3 {
			_, activeErr := client.Policies().CrossTenantAccessPolicy().Partners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(id).Get(ctx, nil)
			activeAbsent := activeErr != nil && grapherrors.GraphError(ctx, activeErr).StatusCode == 404
			if activeErr != nil && !activeAbsent {
				return fmt.Errorf("verify active partner absence: %w", activeErr)
			}
			deletedAbsent := true
			if rs.Primary.Attributes["hard_delete"] == "true" {
				_, deletedErr := client.Policies().DeletedItems().CrossTenantPartners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(id).Get(ctx, nil)
				deletedAbsent = deletedErr != nil && grapherrors.GraphError(ctx, deletedErr).StatusCode == 404
				if deletedErr != nil && !deletedAbsent {
					return fmt.Errorf("verify permanent partner deletion: %w", deletedErr)
				}
			}
			if activeAbsent && deletedAbsent {
				absentReads++
			} else {
				absentReads = 0
			}
			if absentReads == 3 {
				break
			}
			select {
			case <-time.After(3 * time.Second):
			case <-ctx.Done():
				return fmt.Errorf("partner %s was not confirmed absent from the expected stores: %w", id, ctx.Err())
			}
		}
	}
	return nil
}
