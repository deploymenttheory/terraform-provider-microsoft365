package graphBetaServicePrincipal_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipal "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal"
	servicePrincipalMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaServicePrincipal.ResourceName
	testResource = graphBetaServicePrincipal.ServicePrincipalTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() *servicePrincipalMocks.MockState {
	httpmock.Activate()
	mockState := servicePrincipalMocks.RegisterServicePrincipalMockResponders()
	return mockState
}

func TestUnitResourceServicePrincipal_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_minimal").Key("display_name").HasValue("Test Service Principal"),
					check.That(resourceType+".test_minimal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("service_principal_type").HasValue("Application"),
				),
			},
			{
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceServicePrincipal_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("app_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_maximal").Key("display_name").Exists(),
					check.That(resourceType+".test_maximal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("app_role_assignment_required").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("Maximal service principal configuration for testing"),
					check.That(resourceType+".test_maximal").Key("login_url").HasValue("https://login.example.com"),
					check.That(resourceType+".test_maximal").Key("notes").HasValue("Service principal for maximal unit testing"),
					check.That(resourceType+".test_maximal").Key("notification_email_addresses.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("preferred_single_sign_on_mode").HasValue("saml"),
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("service_principal_type").HasValue("Application"),
				),
			},
			{
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIdFunc(resourceType + ".test_maximal"),
			},
		},
	})
}

func TestUnitResourceServicePrincipal_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("app_role_assignment_required").HasValue("false"),
				),
			},
			{
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id                       = "11111111-1111-1111-1111-111111111111"
  app_role_assignment_required = true
  tags                         = ["HideApp"]
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_minimal").Key("app_role_assignment_required").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("tags.#").HasValue("1"),
				),
			},
		},
	})
}
