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
					check.That(resourceType+".test_maximal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("app_role_assignment_required").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("Maximal service principal configuration for testing"),
					check.That(resourceType+".test_maximal").Key("login_url").HasValue("https://login.example.com"),
					check.That(resourceType+".test_maximal").Key("notes").HasValue("Service principal for maximal unit testing"),
					check.That(resourceType+".test_maximal").Key("notification_email_addresses.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("preferred_single_sign_on_mode").HasValue("saml"),
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("alternative_names.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("saml_single_sign_on_settings.relay_state").HasValue("https://example.com/relay"),
					check.That(resourceType+".test_maximal").Key("token_encryption_key_id").HasValue("cccccccc-1111-2222-3333-444444444444"),
					check.That(resourceType+".test_maximal").Key("service_principal_type").HasValue("Application"),
					// Computed properties returned by the API
					check.That(resourceType+".test_maximal").Key("app_owner_organization_id").HasValue("2cbe968c-9683-4d8a-af06-dab1bb350a04"),
					check.That(resourceType+".test_maximal").Key("created_by_app_id").HasValue("04b07795-8ddb-461a-bbee-02f9e1bf7b46"),
					// Credential collections round-trip the API wire format: base64 identifier,
					// UUID key ID and RFC3339 timestamps
					check.That(resourceType+".test_maximal").Key("key_credentials.#").HasValue("1"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.custom_key_identifier").HasValue("a8NSGsQqlkjIPN1kEpJ8QIe4AgI="),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.display_name").HasValue("CN=test-signing"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.key_id").HasValue("dddddddd-1111-2222-3333-444444444444"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.start_date_time").HasValue("2026-01-01T00:00:00Z"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.end_date_time").HasValue("2027-01-01T00:00:00Z"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.type").HasValue("AsymmetricX509Cert"),
					check.That(resourceType+".test_maximal").Key("key_credentials.0.usage").HasValue("Sign"),
					check.That(resourceType+".test_maximal").Key("password_credentials.#").HasValue("1"),
					check.That(resourceType+".test_maximal").Key("password_credentials.0.custom_key_identifier").HasValue("a8NSGsQqlkjIPN1kEpJ8QIe4AgI="),
					check.That(resourceType+".test_maximal").Key("password_credentials.0.hint").HasValue("abc"),
					check.That(resourceType+".test_maximal").Key("password_credentials.0.key_id").HasValue("eeeeeeee-1111-2222-3333-444444444444"),
					check.That(resourceType+".test_maximal").Key("password_credentials.0.end_date_time").HasValue("2027-01-01T00:00:00Z"),
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

func TestUnitResourceServicePrincipal_04_RemoveSamlSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id = "11111111-1111-1111-1111-111111111111"

  saml_single_sign_on_settings = {
    relay_state = "https://example.com/relay"
  }
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("saml_single_sign_on_settings.relay_state").HasValue("https://example.com/relay"),
				),
			},
			{
				// Removing the block must clear the property remotely (explicit JSON null)
				// and converge without an inconsistent-result error
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id = "11111111-1111-1111-1111-111111111111"
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("saml_single_sign_on_settings").DoesNotExist(),
				),
			},
		},
	})
}

func TestUnitResourceServicePrincipal_05_ClearAlternativeNames(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id            = "11111111-1111-1111-1111-111111111111"
  alternative_names = ["isExplicit=True"]
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("alternative_names.#").HasValue("1"),
				),
			},
			{
				// An explicit empty set clears previously configured values remotely
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id            = "11111111-1111-1111-1111-111111111111"
  alternative_names = []
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("alternative_names.#").HasValue("0"),
				),
			},
		},
	})
}

func TestUnitResourceServicePrincipal_06_ClearTokenEncryptionKeyId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	mockState := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mockState.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id                  = "11111111-1111-1111-1111-111111111111"
  token_encryption_key_id = "cccccccc-1111-2222-3333-444444444444"
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("token_encryption_key_id").HasValue("cccccccc-1111-2222-3333-444444444444"),
				),
			},
			{
				// Removing the attribute must clear the property remotely (explicit JSON null)
				// and converge without an inconsistent-result error
				Config: `
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id = "11111111-1111-1111-1111-111111111111"
}`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test_minimal").Key("token_encryption_key_id").DoesNotExist(),
				),
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
