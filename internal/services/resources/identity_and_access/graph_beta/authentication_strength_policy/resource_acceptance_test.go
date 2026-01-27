package graphBetaAuthenticationStrengthPolicy_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAuthenticationStrength "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_strength_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAuthenticationStrength.ResourceName

	// testResource is the test resource implementation for authentication strength policies
	testResource = graphBetaAuthenticationStrength.AuthenticationStrengthTestResource{}
)

func TestAccResourceAuthenticationStrengthPolicy_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			10*time.Second,
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
					testlog.StepAction(resourceType, "Creating minimal authentication strength policy")
				},
				Config: testAccConfigAuthStrengthMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication strength policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".auth_strength_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".auth_strength_minimal").Key("id").Exists(),
					check.That(resourceType+".auth_strength_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-min-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auth_strength_minimal").Key("description").HasValue("Acceptance test minimal authentication strength policy"),
					check.That(resourceType+".auth_strength_minimal").Key("allowed_combinations.#").HasValue("1"),
					check.That(resourceType+".auth_strength_minimal").Key("allowed_combinations.*").ContainsTypeSetElement("password,sms"),
					check.That(resourceType+".auth_strength_minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".auth_strength_minimal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal authentication strength policy")
				},
				ResourceName:      resourceType + ".auth_strength_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccResourceAuthenticationStrengthPolicy_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			10*time.Second,
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
					testlog.StepAction(resourceType, "Creating maximal authentication strength policy with combination configurations")
				},
				Config: testAccConfigAuthStrengthMaximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication strength policy", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".auth_strength_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".auth_strength_maximal").Key("id").Exists(),
					check.That(resourceType+".auth_strength_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-max-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auth_strength_maximal").Key("description").HasValue("Acceptance test maximal authentication strength policy with all combinations and configurations"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.#").HasValue("22"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("deviceBasedPush"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("federatedMultiFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("federatedSingleFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("fido2"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("hardwareOath,federatedSingleFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("microsoftAuthenticatorPush,federatedSingleFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("password"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("windowsHelloForBusiness"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("x509CertificateMultiFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("combination_configurations.#").HasValue("3"),
					check.That(resourceType+".auth_strength_maximal").Key("created_date_time").Exists(),
					check.That(resourceType+".auth_strength_maximal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal authentication strength policy")
				},
				ResourceName:      resourceType + ".auth_strength_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// Configuration helper functions
func testAccConfigAuthStrengthMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength minimal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigAuthStrengthMaximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_maximal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength maximal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
