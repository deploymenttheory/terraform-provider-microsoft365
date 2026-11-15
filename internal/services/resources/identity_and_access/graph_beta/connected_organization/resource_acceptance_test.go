package graphBetaIdentityAndAccessConnectedOrganization_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIdentityAndAccessConnectedOrganization "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/connected_organization"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessConnectedOrganization.ResourceName

	// testResource is the test resource implementation for connected organizations
	testResource = graphBetaIdentityAndAccessConnectedOrganization.ConnectedOrganizationTestResource{}

	accTestResourceName = resourceType + ".test"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// accTestCase builds a test case with the shared pre-check, providers and destroy check.
func accTestCase(t *testing.T, steps []resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
		Steps: steps,
	}
}

func importStep(description string) resource.TestStep {
	return resource.TestStep{
		PreConfig: func() {
			testlog.StepAction(resourceType, description)
		},
		ResourceName:            accTestResourceName,
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

func minimalChecks() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(accTestResourceName).ExistsInGraph(testResource),
		check.That(accTestResourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
		check.That(accTestResourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-connected-organization-[a-z0-9]{8}$`)),
		check.That(accTestResourceName).Key("description").DoesNotExist(),
		check.That(accTestResourceName).Key("state").HasValue("configured"),
		check.That(accTestResourceName).Key("identity_sources.#").HasValue("1"),
		check.That(accTestResourceName).Key("identity_sources.0.odata_type").HasValue("#microsoft.graph.domainIdentitySource"),
		check.That(accTestResourceName).Key("identity_sources.0.domain_name").HasValue("fabrikam.com"),
		check.That(accTestResourceName).Key("created_by").IsSet(),
		check.That(accTestResourceName).Key("created_date_time").IsSet(),
		check.That(accTestResourceName).Key("modified_date_time").IsSet(),
	)
}

func maximalChecks() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(accTestResourceName).ExistsInGraph(testResource),
		check.That(accTestResourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
		check.That(accTestResourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-connected-organization-maximal-[a-z0-9]{8}$`)),
		check.That(accTestResourceName).Key("description").HasValue("Connected organization managed by Terraform acceptance tests"),
		check.That(accTestResourceName).Key("state").HasValue("proposed"),
		check.That(accTestResourceName).Key("identity_sources.#").HasValue("1"),
		check.That(accTestResourceName).Key("identity_sources.0.domain_name").HasValue("fabrikam.com"),
	)
}

// Minimal: display name and one domain identity source.
func TestAccResourceConnectedOrganization_01_Minimal(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating minimal connected organization")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			Check:  minimalChecks(),
		},
		importStep("Importing minimal connected organization"),
	}))
}

// Maximal: description, state and timeouts.
func TestAccResourceConnectedOrganization_02_Maximal(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating maximal connected organization")
			},
			Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
			Check:  maximalChecks(),
		},
		importStep("Importing maximal connected organization"),
	}))
}

// Minimal to maximal: in-place update of display name, description and state.
func TestAccResourceConnectedOrganization_03_MinimalToMaximal(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating minimal connected organization")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			Check:  minimalChecks(),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Updating connected organization to maximal")
			},
			Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(accTestResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: maximalChecks(),
		},
		importStep("Importing updated connected organization"),
	}))
}

// Maximal to minimal: in-place update that clears the description and resets state.
func TestAccResourceConnectedOrganization_04_MaximalToMinimal(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating maximal connected organization")
			},
			Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
			Check:  maximalChecks(),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Updating connected organization to minimal")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(accTestResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: minimalChecks(),
		},
		importStep("Importing updated connected organization"),
	}))
}

// Microsoft Entra tenant identity source.
func TestAccResourceConnectedOrganization_05_AzureActiveDirectoryTenant(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating connected organization for a Microsoft Entra tenant")
			},
			Config: loadAcceptanceTestTerraform("resource_azure_ad_tenant.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(accTestResourceName).ExistsInGraph(testResource),
				check.That(accTestResourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-connected-organization-tenant-[a-z0-9]{8}$`)),
				check.That(accTestResourceName).Key("identity_sources.#").HasValue("1"),
				check.That(accTestResourceName).Key("identity_sources.0.odata_type").HasValue("#microsoft.graph.azureActiveDirectoryTenant"),
				check.That(accTestResourceName).Key("identity_sources.0.tenant_id").HasValue("72f988bf-86f1-41af-91ab-2d7cd011db47"),
				check.That(accTestResourceName).Key("identity_sources.0.domain_name").DoesNotExist(),
			),
		},
		importStep("Importing connected organization for a Microsoft Entra tenant"),
	}))
}

// Several domain identity sources in one connected organization.
func TestAccResourceConnectedOrganization_06_MultipleIdentitySources(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating connected organization with two domain identity sources")
			},
			Config: loadAcceptanceTestTerraform("resource_multiple_identity_sources.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(accTestResourceName).ExistsInGraph(testResource),
				check.That(accTestResourceName).Key("identity_sources.#").HasValue("2"),
			),
		},
		importStep("Importing connected organization with two domain identity sources"),
	}))
}

// Changing the identity source destroys and recreates the connected organization.
func TestAccResourceConnectedOrganization_07_IdentitySourceForcesReplace(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating minimal connected organization")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			Check:  minimalChecks(),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Replacing the identity source")
			},
			Config: loadAcceptanceTestTerraform("resource_replace_identity_source.tf"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(accTestResourceName, plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
			Check: resource.ComposeTestCheckFunc(
				check.That(accTestResourceName).ExistsInGraph(testResource),
				check.That(accTestResourceName).Key("identity_sources.#").HasValue("1"),
				check.That(accTestResourceName).Key("identity_sources.0.domain_name").HasValue("adatum.com"),
			),
		},
		importStep("Importing replaced connected organization"),
	}))
}
