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
	// Resource type name constructed from exported constants
	resourceType = constants.PROVIDER_NAME + "_" + graphBetaAuthenticationStrength.ResourceName

	// testResource is the test resource implementation for authentication strength policies
	testResource = graphBetaAuthenticationStrength.AuthenticationStrengthTestResource{}
)

func TestAccAuthenticationStrengthResource_Minimal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
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
						testlog.WaitForConsistency("authentication strength policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("display_name").HasValue("acc-test-auth-strength-min"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("description").HasValue("Acceptance test minimal authentication strength policy"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("allowed_combinations.#").HasValue("1"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("allowed_combinations.*").ContainsTypeSetElement("password,sms"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("created_date_time").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal authentication strength policy")
				},
				ResourceName:      "microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccAuthenticationStrengthResource_MFAOnly(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating MFA-only authentication strength policy")
				},
				Config: testAccConfigAuthStrengthMFAOnly(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication strength policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("display_name").HasValue("acc-test-auth-strength-mfa"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("description").HasValue("Acceptance test MFA-only authentication strength policy"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("allowed_combinations.#").HasValue("4"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("allowed_combinations.*").ContainsTypeSetElement("fido2"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("allowed_combinations.*").ContainsTypeSetElement("windowsHelloForBusiness"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("allowed_combinations.*").ContainsTypeSetElement("microsoftAuthenticatorPush,federatedSingleFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only").Key("allowed_combinations.*").ContainsTypeSetElement("x509CertificateMultiFactor"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing MFA-only authentication strength policy")
				},
				ResourceName:      "microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccAuthenticationStrengthResource_Maximal(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
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
						testlog.WaitForConsistency("authentication strength policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("id").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-authentication-strength-maximal-[0-9a-fA-F-]+$`)),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("description").HasValue("Acceptance test maximal authentication strength policy with all combinations and configurations"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.#").HasValue("22"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("deviceBasedPush"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("federatedMultiFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("federatedSingleFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("fido2"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("hardwareOath,federatedSingleFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("microsoftAuthenticatorPush,federatedSingleFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("password"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("windowsHelloForBusiness"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("x509CertificateMultiFactor"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("combination_configurations.#").HasValue("3"),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("created_date_time").Exists(),
					check.That("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal").Key("modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal authentication strength policy")
				},
				ResourceName:      "microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal",
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

func testAccConfigAuthStrengthMFAOnly() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_auth_strength_mfa_only.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication strength MFA-only config: %s", err.Error()))
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
