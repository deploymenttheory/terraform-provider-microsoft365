package graphBetaAdministrativeUnit_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAdministrativeUnit "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAdministrativeUnit.ResourceName

	// testResource is the test resource implementation for administrative units
	testResource = graphBetaAdministrativeUnit.AdministrativeUnitTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// AU001: User-Based Administrative Unit
func TestAccResourceAdministrativeUnit_01_AU001(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AU001 user-based administrative unit")
				},
				Config: loadAcceptanceTestTerraform("resource_au001_user_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au001_user_based").ExistsInGraph(testResource),
					check.That(resourceType+".au001_user_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au001_user_based").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-au001-user-based-[a-z0-9]{8}$`)),
					check.That(resourceType+".au001_user_based").Key("description").HasValue("Administrative unit for user-based testing"),
					check.That(resourceType+".au001_user_based").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AU001 administrative unit")
				},
				ResourceName: resourceType + ".au001_user_based",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au001_user_based"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au001_user_based")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU002: Group-Based Administrative Unit
func TestAccResourceAdministrativeUnit_02_AU002(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AU002 group-based administrative unit")
				},
				Config: loadAcceptanceTestTerraform("resource_au002_group_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au002_group_based").ExistsInGraph(testResource),
					check.That(resourceType+".au002_group_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au002_group_based").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-au002-group-based-[a-z0-9]{8}$`)),
					check.That(resourceType+".au002_group_based").Key("description").HasValue("Administrative unit for group-based testing"),
					check.That(resourceType+".au002_group_based").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AU002 administrative unit")
				},
				ResourceName: resourceType + ".au002_group_based",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au002_group_based"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au002_group_based")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU003: Mixed User and Group Administrative Unit
func TestAccResourceAdministrativeUnit_03_AU003(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AU003 mixed user and group administrative unit")
				},
				Config: loadAcceptanceTestTerraform("resource_au003_mixed.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au003_mixed").ExistsInGraph(testResource),
					check.That(resourceType+".au003_mixed").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au003_mixed").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-au003-mixed-[a-z0-9]{8}$`)),
					check.That(resourceType+".au003_mixed").Key("description").HasValue("Administrative unit for mixed user and group testing"),
					check.That(resourceType+".au003_mixed").Key("visibility").HasValue("HiddenMembership"),
					check.That(resourceType+".au003_mixed").Key("is_member_management_restricted").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AU003 administrative unit")
				},
				ResourceName: resourceType + ".au003_mixed",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au003_mixed"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au003_mixed")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

// AU004: Multi-Step Update Test
func TestAccResourceAdministrativeUnit_04_AU004_Updates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			// Step 1: Create initial administrative unit
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AU004 administrative unit - Initial state")
				},
				Config: loadAcceptanceTestTerraform("resource_au004_update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au004_update").ExistsInGraph(testResource),
					check.That(resourceType+".au004_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".au004_update").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-au004-update-[a-z0-9]{8}$`)),
					check.That(resourceType+".au004_update").Key("description").HasValue("Initial description for update testing"),
				),
			},
			// Step 2: Update description
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating AU004 - Updating description")
					testlog.WaitForConsistency("administrative unit", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_au004_update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au004_update").ExistsInGraph(testResource),
					check.That(resourceType+".au004_update").Key("description").HasValue("Updated description for update testing"),
				),
			},
			// Step 3: Change to dynamic membership
			// NOTE: This step triggers a RequiresReplace cycle on the 'visibility' attribute
			// (immutable in Graph API), causing a destroy+create. After destroy+create, both
			// the deletion of the old resource and the creation of the new resource need to
			// propagate across Entra replicas, which takes longer than a plain PATCH update.
			// We use a longer post-apply wait (30s) to allow the new resource to propagate.
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating AU004 - Changing to dynamic membership (requires destroy+recreate due to immutable visibility attribute)")
					testlog.WaitForConsistency("administrative unit", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_au004_update_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						// After a RequiresReplace (destroy+create), the new resource needs more
						// time to propagate across Entra replicas than a simple update (PATCH).
						testlog.WaitForConsistency("administrative unit", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".au004_update").ExistsInGraph(testResource),
					check.That(resourceType+".au004_update").Key("description").HasValue("Updated to dynamic membership"),
					check.That(resourceType+".au004_update").Key("membership_type").HasValue("Dynamic"),
					check.That(resourceType+".au004_update").Key("membership_rule").HasValue("(user.country -eq \"United States\")"),
					check.That(resourceType+".au004_update").Key("membership_rule_processing_state").HasValue("On"),
				),
			},
			// Step 4: Pause dynamic membership
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating AU004 - Pausing dynamic membership")
					testlog.WaitForConsistency("administrative unit", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_au004_update_step4.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".au004_update").ExistsInGraph(testResource),
					check.That(resourceType+".au004_update").Key("description").HasValue("Paused dynamic membership"),
					check.That(resourceType+".au004_update").Key("membership_rule_processing_state").HasValue("Paused"),
				),
			},
			// Final step: Import to verify state
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AU004 administrative unit")
				},
				ResourceName: resourceType + ".au004_update",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".au004_update"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".au004_update")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}
