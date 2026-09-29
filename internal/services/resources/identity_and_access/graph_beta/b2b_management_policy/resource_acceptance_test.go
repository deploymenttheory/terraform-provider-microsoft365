package graphBetaIdentityAndAccessB2bManagementPolicy_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaB2bManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy"
)

var testResource = graphBetaB2bManagementPolicy.B2bManagementPolicyTestResource{}

// is_organization_default stays false so the test does not change tenant-wide B2B settings.
func TestAccResourceB2bManagementPolicy_01_Lifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating B2B management policy")
				},
				Config: testAccConfig("tests/terraform/acceptance/resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("display_name").
						MatchesRegex(regexp.MustCompile(`^acc-test-b2b-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").Key("definition.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing B2B management policy")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating B2B management policy")
				},
				Config: testAccConfig("tests/terraform/acceptance/resource_02_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("display_name").
						MatchesRegex(regexp.MustCompile(`^acc-test-b2b-updated-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").
						Key("definition.0").
						MatchesRegex(regexp.MustCompile(`AllowedDomains`)),
				),
			},
		},
	})
}

func testAccConfig(path string) string {
	accTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		panic(
			fmt.Sprintf(
				"failed to load B2B management policy acceptance config %s: %s",
				path,
				err.Error(),
			),
		)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func TestAccResourceB2bManagementPolicy_02_AllowDomains(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_03_BlockDomains(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 02_block_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_04_AnyDomain(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 03_any_domain") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_22_any_domain.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_05_AllowToBlock(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 02_block_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_06_BlockToAllow(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 02_block_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_07_RestrictedToAny(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 03_any_domain") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_22_any_domain.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_08_AnyToRestricted(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 03_any_domain") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 02_block_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_21_block_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_09_MinimalToMaximal(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 04_expanded_domains") },
				Config: testAccConfig(
					"tests/terraform/acceptance/resource_23_expanded_domains.tf",
				),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example","partner.example","sub.partner.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config: testAccConfig(
					"tests/terraform/acceptance/resource_23_expanded_domains.tf",
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_10_MaximalToMinimal(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 04_expanded_domains") },
				Config: testAccConfig(
					"tests/terraform/acceptance/resource_23_expanded_domains.tf",
				),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example","partner.example","sub.partner.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_11_OtherSettings(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 05_other_settings") },
				Config: testAccConfig(
					"tests/terraform/acceptance/resource_24_other_settings.tf",
				),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"AutoRedeemPolicy":{"AdminConsentedForUsersIntoTenantIds":[],"NoAADConsentForUsersFromTenantsIds":[]},"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example"]},"PreviewPolicy":{"Features":[]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_24_other_settings.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_12_Rename(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 01_allow_domains") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 07_renamed") },
				Config:    testAccConfig("tests/terraform/acceptance/resource_26_renamed.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").
						Key("display_name").
						MatchesRegex(regexp.MustCompile(`^acc-test-b2b-renamed-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency(
						"B2B management policy replicas before import",
						30*time.Second,
					)
					time.Sleep(30 * time.Second)
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_26_renamed.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceB2bManagementPolicy_13_FormattedJSON(t *testing.T) {
	var id string
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
				PreConfig: func() { testlog.StepAction(resourceType, "Apply inactive policy: 06_formatted_json") },
				Config: testAccConfig(
					"tests/terraform/acceptance/resource_25_formatted_json.tf",
				),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("B2B management policy", 5*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && id != current {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testAccConfig("tests/terraform/acceptance/resource_25_formatted_json.tf"),
				PlanOnly: true,
			},
		},
	})
}
