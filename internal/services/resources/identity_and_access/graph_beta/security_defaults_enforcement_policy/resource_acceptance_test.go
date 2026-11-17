package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/security_defaults_enforcement_policy"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy.ResourceName

	// testResource is the test resource implementation for the security defaults enforcement policy
	testResource = graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy.SecurityDefaultsEnforcementPolicyTestResource{}
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
			testlog.StepAction(resourceType, "Importing the security defaults enforcement policy")
		},
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "identitySecurityDefaultsEnforcementPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// The security defaults enforcement policy is a tenant-wide singleton. Enabling it enforces MFA registration for
// every user in the shared test tenant and conflicts with the Conditional Access suites, so these tests only apply
// is_enabled = false, which is the tenant's current value (the PreCheck skips them otherwise). Enabling is covered by
// the unit tests. Every test records the original policy in its PreCheck and a t.Cleanup PATCHes isEnabled back and
// verifies the restore with GET (R9). CheckDestroy asserts that destroy left security defaults disabled (R7).

// Disabled: the policy is adopted with is_enabled = false, imported, and a follow-up plan is empty.
func TestAccResourceSecurityDefaultsEnforcementPolicy_01_Disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { securityDefaultsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkTenantSecurityDefaultsDisabled,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the security defaults enforcement policy with is_enabled = false")
				},
				Config: loadAcceptanceTestTerraform("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("identitySecurityDefaultsEnforcementPolicy"),
					check.That(resourceType+".test").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("display_name").MatchesRegex(regexp.MustCompile(`\S`)),
					check.That(resourceType+".test").Key("description").MatchesRegex(regexp.MustCompile(`\S`)),
					checkTenantSecurityDefaultsDisabled,
				),
			},
			acceptanceImportStep(),
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Confirming the plan is empty after apply and import")
				},
				Config:   loadAcceptanceTestTerraform("resource_disabled.tf"),
				PlanOnly: true,
			},
		},
	})
}

// ImportExisting: the existing policy is imported without a prior apply and the plan for the matching
// configuration is empty, so nothing is sent to Graph.
func TestAccResourceSecurityDefaultsEnforcementPolicy_02_ImportExisting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { securityDefaultsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkTenantSecurityDefaultsDisabled,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the existing security defaults enforcement policy")
				},
				Config:             loadAcceptanceTestTerraform("resource_disabled.tf"),
				ResourceName:       resourceType + ".test",
				ImportState:        true,
				ImportStateId:      "identitySecurityDefaultsEnforcementPolicy",
				ImportStatePersist: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Confirming the plan is empty after import")
				},
				Config:   loadAcceptanceTestTerraform("resource_disabled.tf"),
				PlanOnly: true,
			},
		},
	})
}
