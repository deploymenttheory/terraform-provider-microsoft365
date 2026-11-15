package graphBetaRoleManagementPolicy_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/role_management_policy"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaRoleManagementPolicy.ResourceName

	// testResource is the test resource implementation for role management policies
	testResource = graphBetaRoleManagementPolicy.RoleManagementPolicyTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Role management policies exist for every role and cannot be created or deleted. Every test records the rules of the
// Usage Summary Reports Reader policy in its PreCheck, and a t.Cleanup PATCHes them back and verifies the restore with
// GET (R9). CheckDestroy asserts that destroy left the policy in place (R7).

// ExpirationRule: adopts the policy by role definition, changes one expiration rule, updates it and imports the policy.
func TestAccResourceRoleManagementPolicy_01_ExpirationRule(t *testing.T) {
	policyId := roleManagementPolicyPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the Usage Summary Reports Reader policy and changing one expiration rule")
				},
				Config: loadAcceptanceTestTerraform("resource_expiration_initial.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^DirectoryRole_[0-9a-fA-F-]+_[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("role_definition_id").HasValue(accRoleDefinitionId),
					check.That(resourceType+".test").Key("scope_id").HasValue("/"),
					check.That(resourceType+".test").Key("scope_type").HasValue("DirectoryRole"),
					check.That(resourceType+".test").Key("last_modified_date_time").MatchesRegex(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)),
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("1"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.is_expiration_required").HasValue("false"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating the expiration rules")
				},
				Config: loadAcceptanceTestTerraform("resource_expiration_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("2"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P60D"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.is_expiration_required").HasValue("true"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").HasValue("PT4H"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the role management policy")
				},
				ResourceName:     resourceType + ".test",
				ImportState:      true,
				ImportStateId:    policyId,
				ImportStateCheck: checkAccImportedAllRules(policyId),
			},
		},
	})
}

// Complete: manages every supported rule, so an import matches the applied state exactly.
func TestAccResourceRoleManagementPolicy_02_Complete(t *testing.T) {
	policyId := roleManagementPolicyPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Managing every supported rule of the Usage Summary Reports Reader policy")
				},
				Config: loadAcceptanceTestTerraform("resource_complete.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("3"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Eligibility.maximum_duration").HasValue("P180D"),
					check.That(resourceType+".test").Key("enablement_rules.%").HasValue("3"),
					check.That(resourceType+".test").Key("enablement_rules.Enablement_Admin_Eligibility.enabled_rules.#").HasValue("0"),
					check.That(resourceType+".test").Key("notification_rules.%").HasValue("9"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Admin_Admin_Assignment.notification_recipients.#").HasValue("0"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.is_enabled").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing the role management policy")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           policyId,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// ByPolicyID: adopts the policy by its ID; the role definition is looked up from the policy assignments.
func TestAccResourceRoleManagementPolicy_03_ByPolicyID(t *testing.T) {
	policyId := roleManagementPolicyPreCheck(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the Usage Summary Reports Reader policy by policy ID")
				},
				Config: loadAcceptanceTestTerraform("resource_by_policy_id.tf"),
				ConfigVariables: config.Variables{
					"policy_id": config.StringVariable(policyId),
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue(policyId),
					check.That(resourceType+".test").Key("role_definition_id").HasValue(accRoleDefinitionId),
					check.That(resourceType+".test").Key("scope_type").HasValue("DirectoryRole"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Eligibility.maximum_duration").HasValue("P180D"),
				),
			},
		},
	})
}
