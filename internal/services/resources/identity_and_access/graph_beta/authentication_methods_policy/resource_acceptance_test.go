package graphBetaIdentityAndAccessAuthenticationMethodsPolicy_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaIdentityAndAccessAuthenticationMethodsPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_methods_policy"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessAuthenticationMethodsPolicy.ResourceName

	// testResource is the test resource implementation for the authentication methods policy
	testResource = graphBetaIdentityAndAccessAuthenticationMethodsPolicy.AuthenticationMethodsPolicyTestResource{}

	// externalProviders are the providers the acceptance fixtures use besides microsoft365
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// acceptanceImportStep verifies import by the fixed singleton ID.
func acceptanceImportStep() resource.TestStep {
	return resource.TestStep{
		PreConfig: func() {
			testlog.StepAction(resourceType, "Importing the authentication methods policy")
		},
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "authenticationMethodsPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// The authentication methods policy is a tenant-wide singleton. Every test records the original policy in its
// PreCheck and a t.Cleanup PATCHes the Email method configuration back and verifies with GET that the whole
// policy matches the recording (R9). The tests only change the Email method, and only its targets, using
// acc-test groups. CheckDestroy asserts that destroy left the policy in the tenant (R7).

// Minimal: nothing is configured, so the policy is adopted without any write.
func TestAccResourceAuthenticationMethodsPolicy_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { authenticationMethodsPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the authentication methods policy with an empty configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("authenticationMethodsPolicy"),
					check.That(resourceType+".test").Key("display_name").HasValue("Authentication Methods Policy"),
					check.That(resourceType+".test").Key("policy_migration_state").MatchesRegex(regexp.MustCompile(`^(preMigration|migrationInProgress|migrationComplete)$`)),
					check.That(resourceType+".test").Key("email.state").MatchesRegex(regexp.MustCompile(`^(enabled|disabled)$`)),
					check.That(resourceType+".test").Key("temporary_access_pass.state").MatchesRegex(regexp.MustCompile(`^(enabled|disabled)$`)),
					check.That(resourceType+".test").Key("sms.state").MatchesRegex(regexp.MustCompile(`^(enabled|disabled)$`)),
					check.That(resourceType+".test").Key("software_oath.state").MatchesRegex(regexp.MustCompile(`^(enabled|disabled)$`)),
					check.That(resourceType+".test").Key("hardware_oath.state").MatchesRegex(regexp.MustCompile(`^(enabled|disabled)$`)),
					checkUnmanagedMethodsUnchanged,
				),
			},
			acceptanceImportStep(),
		},
	})
}

// EmailTargets: the Email method targets an acc-test group, then switches to a second group and excludes the
// first. Other method configurations must not change.
func TestAccResourceAuthenticationMethodsPolicy_02_EmailTargets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { authenticationMethodsPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Targeting the email method at an acc-test group")
				},
				Config: loadAcceptanceTestTerraform("resource_email_include.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("email.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("email.allow_external_id_to_use_email_otp").HasValue("default"),
					check.That(resourceType+".test").Key("email.include_targets.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test", "email.include_targets.*.id", "microsoft365_graph_beta_groups_group.first", "id"),
					check.That(resourceType+".test").Key("email.exclude_targets.#").HasValue("0"),
					checkUnmanagedMethodsUnchanged,
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Switching the email method to a second acc-test group and excluding the first")
				},
				Config: loadAcceptanceTestTerraform("resource_email_exclude.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("email.include_targets.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test", "email.include_targets.*.id", "microsoft365_graph_beta_groups_group.second", "id"),
					check.That(resourceType+".test").Key("email.exclude_targets.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test", "email.exclude_targets.*.id", "microsoft365_graph_beta_groups_group.first", "id"),
					checkUnmanagedMethodsUnchanged,
				),
			},
			acceptanceImportStep(),
		},
	})
}
