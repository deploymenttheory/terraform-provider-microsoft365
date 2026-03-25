package graphBetaAdministrativeUnitMembership_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAdministrativeUnitMembership "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/administrative_unit_membership"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaAdministrativeUnitMembership.ResourceName
	testResource = graphBetaAdministrativeUnitMembership.AdministrativeUnitMembershipTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// AUM001: User-Based Membership
func TestAccResourceAdministrativeUnitMembership_01_AUM001(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AUM001 user-based membership")
				},
				Config: loadAcceptanceTestTerraform("resource_aum001_user_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum001_user_based").ExistsInGraph(testResource),
					check.That(resourceType+".aum001_user_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aum001_user_based").Key("members.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AUM001 membership")
				},
				ResourceName:      resourceType + ".aum001_user_based",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// AUM002: Group-Based Membership
func TestAccResourceAdministrativeUnitMembership_02_AUM002(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AUM002 group-based membership")
				},
				Config: loadAcceptanceTestTerraform("resource_aum002_group_based.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum002_group_based").ExistsInGraph(testResource),
					check.That(resourceType+".aum002_group_based").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aum002_group_based").Key("members.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AUM002 membership")
				},
				ResourceName:      resourceType + ".aum002_group_based",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// AUM003: Mixed User and Group Membership
func TestAccResourceAdministrativeUnitMembership_03_AUM003(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AUM003 mixed membership")
				},
				Config: loadAcceptanceTestTerraform("resource_aum003_mixed.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum003_mixed").ExistsInGraph(testResource),
					check.That(resourceType+".aum003_mixed").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".aum003_mixed").Key("members.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AUM003 membership")
				},
				ResourceName:      resourceType + ".aum003_mixed",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// AUM004: Multi-Step Update Test
func TestAccResourceAdministrativeUnitMembership_04_AUM004_Updates(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AUM004 membership - Initial state")
				},
				Config: loadAcceptanceTestTerraform("resource_aum004_update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum004_update").ExistsInGraph(testResource),
					check.That(resourceType+".aum004_update").Key("members.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating AUM004 - Adding members")
					testlog.WaitForConsistency("administrative unit membership", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_aum004_update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum004_update").ExistsInGraph(testResource),
					check.That(resourceType+".aum004_update").Key("members.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating AUM004 - Removing members")
					testlog.WaitForConsistency("administrative unit membership", 20*time.Second)
					time.Sleep(20 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("resource_aum004_update_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("administrative unit membership", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".aum004_update").ExistsInGraph(testResource),
					check.That(resourceType+".aum004_update").Key("members.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AUM004 membership")
				},
				ResourceName:      resourceType + ".aum004_update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
