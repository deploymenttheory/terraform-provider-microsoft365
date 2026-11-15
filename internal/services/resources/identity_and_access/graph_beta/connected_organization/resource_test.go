package graphBetaIdentityAndAccessConnectedOrganization_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	connectedOrganizationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/connected_organization/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *connectedOrganizationMocks.ConnectedOrganizationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectedOrganizationMock := &connectedOrganizationMocks.ConnectedOrganizationMock{}
	connectedOrganizationMock.RegisterMocks()
	return mockClient, connectedOrganizationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *connectedOrganizationMocks.ConnectedOrganizationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectedOrganizationMock := &connectedOrganizationMocks.ConnectedOrganizationMock{}
	connectedOrganizationMock.RegisterErrorMocks()
	return mockClient, connectedOrganizationMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

var unitTestResourceName = resourceType + ".test"

// Minimal configuration: display name and one domain identity source; state defaults to configured.
func TestUnitResourceConnectedOrganization_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(unitTestResourceName).Key("display_name").HasValue("unit-test-connected-organization-minimal"),
					check.That(unitTestResourceName).Key("description").DoesNotExist(),
					check.That(unitTestResourceName).Key("state").HasValue("configured"),
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("1"),
					check.That(unitTestResourceName).Key("identity_sources.0.odata_type").HasValue("#microsoft.graph.domainIdentitySource"),
					check.That(unitTestResourceName).Key("identity_sources.0.display_name").HasValue("example.com"),
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").HasValue("example.com"),
					check.That(unitTestResourceName).Key("identity_sources.0.tenant_id").DoesNotExist(),
					check.That(unitTestResourceName).Key("created_by").HasValue("unit-test-app"),
					check.That(unitTestResourceName).Key("created_date_time").IsSet(),
					check.That(unitTestResourceName).Key("modified_date_time").IsSet(),
				),
			},
			{
				ResourceName:      unitTestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Maximal configuration: description, explicit state and timeouts.
func TestUnitResourceConnectedOrganization_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(unitTestResourceName).Key("display_name").HasValue("unit-test-connected-organization-maximal"),
					check.That(unitTestResourceName).Key("description").HasValue("Connected organization managed by Terraform unit tests"),
					check.That(unitTestResourceName).Key("state").HasValue("proposed"),
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("1"),
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").HasValue("example.com"),
				),
			},
			{
				ResourceName:            unitTestResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Minimal to maximal is an in-place PATCH because the identity sources are unchanged.
func TestUnitResourceConnectedOrganization_03_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("state").HasValue("configured"),
					check.That(unitTestResourceName).Key("description").DoesNotExist(),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(unitTestResourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("display_name").HasValue("unit-test-connected-organization-maximal"),
					check.That(unitTestResourceName).Key("description").HasValue("Connected organization managed by Terraform unit tests"),
					check.That(unitTestResourceName).Key("state").HasValue("proposed"),
				),
			},
			{
				ResourceName:            unitTestResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal to minimal clears the description and returns state to its default.
func TestUnitResourceConnectedOrganization_04_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("description").HasValue("Connected organization managed by Terraform unit tests"),
					check.That(unitTestResourceName).Key("state").HasValue("proposed"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(unitTestResourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("display_name").HasValue("unit-test-connected-organization-minimal"),
					check.That(unitTestResourceName).Key("description").DoesNotExist(),
					check.That(unitTestResourceName).Key("state").HasValue("configured"),
				),
			},
			{
				ResourceName:      unitTestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Microsoft Entra tenant identity source.
func TestUnitResourceConnectedOrganization_05_AzureActiveDirectoryTenant(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_azure_ad_tenant.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("1"),
					check.That(unitTestResourceName).Key("identity_sources.0.odata_type").HasValue("#microsoft.graph.azureActiveDirectoryTenant"),
					check.That(unitTestResourceName).Key("identity_sources.0.display_name").HasValue("Contoso"),
					check.That(unitTestResourceName).Key("identity_sources.0.tenant_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").DoesNotExist(),
				),
			},
			{
				ResourceName:      unitTestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// External domain federation identity source.
func TestUnitResourceConnectedOrganization_06_ExternalDomainFederation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_external_domain_federation.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("1"),
					check.That(unitTestResourceName).Key("identity_sources.0.odata_type").HasValue("#microsoft.graph.externalDomainFederation"),
					check.That(unitTestResourceName).Key("identity_sources.0.display_name").HasValue("Fabrikam"),
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").HasValue("fabrikam.example"),
					check.That(unitTestResourceName).Key("identity_sources.0.issuer_uri").HasValue("https://idp.fabrikam.example/saml"),
				),
			},
			{
				ResourceName:      unitTestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Several identity sources in one connected organization.
func TestUnitResourceConnectedOrganization_07_MultipleIdentitySources(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_multiple_identity_sources.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("2"),
					resource.TestCheckTypeSetElemNestedAttrs(unitTestResourceName, "identity_sources.*", map[string]string{"odata_type": "#microsoft.graph.domainIdentitySource", "domain_name": "example.net", "display_name": "example.net"}),
					resource.TestCheckTypeSetElemNestedAttrs(unitTestResourceName, "identity_sources.*", map[string]string{"odata_type": "#microsoft.graph.domainIdentitySource", "domain_name": "example.org", "display_name": "example.org"}),
				),
			},
			{
				ResourceName:      unitTestResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Changing identity sources replaces the connected organization.
func TestUnitResourceConnectedOrganization_08_IdentitySourceForcesReplace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").HasValue("example.com"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_replace_identity_source.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(unitTestResourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(unitTestResourceName).Key("identity_sources.#").HasValue("1"),
					check.That(unitTestResourceName).Key("identity_sources.0.domain_name").HasValue("example.org"),
				),
			},
		},
	})
}

// Graph errors surface through the error handler.
func TestUnitResourceConnectedOrganization_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges|403`),
			},
		},
	})
}

// An identity source with attributes that do not match its odata_type fails at plan time.
func TestUnitResourceConnectedOrganization_10_InvalidIdentitySource(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_identity_source.tf"),
				ExpectError: regexp.MustCompile(`domain_name is required when odata_type is`),
			},
		},
	})
}

// An unsupported state value fails schema validation.
func TestUnitResourceConnectedOrganization_11_InvalidState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectedOrganizationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectedOrganizationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_state.tf"),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
		},
	})
}
