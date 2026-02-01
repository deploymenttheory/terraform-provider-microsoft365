package graphBetaApplicationOwner_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationOwner "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_owner"
	applicationOwnerMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_owner/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationOwner.ResourceName

	// testResource is the test resource implementation for application owners
	testResource = graphBetaApplicationOwner.ApplicationOwnerTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *applicationOwnerMocks.ApplicationOwnerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationOwnerMock := &applicationOwnerMocks.ApplicationOwnerMock{}
	applicationOwnerMock.RegisterMocks()
	return mockClient, applicationOwnerMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *applicationOwnerMocks.ApplicationOwnerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationOwnerMock := &applicationOwnerMocks.ApplicationOwnerMock{}
	applicationOwnerMock.RegisterErrorMocks()
	return mockClient, applicationOwnerMock
}

func TestUnitResourceApplicationOwner_01_OwnerTypeUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationOwnerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationOwnerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_owner_type_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_user").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_user").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_user").Key("owner_id").HasValue("user-11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_user").Key("owner_object_type").HasValue("User"),
					check.That(resourceType+".test_user").Key("owner_type").HasValue("User"),
					check.That(resourceType+".test_user").Key("owner_display_name").HasValue("Test User Owner"),
				),
			},
			{
				ResourceName:      resourceType + ".test_user",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_user"),
			},
		},
	})
}

func TestUnitResourceApplicationOwner_02_OwnerTypeServicePrincipal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationOwnerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationOwnerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_owner_type_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_service_principal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_service_principal").Key("application_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_service_principal").Key("owner_id").HasValue("sp-11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_service_principal").Key("owner_object_type").HasValue("ServicePrincipal"),
					check.That(resourceType+".test_service_principal").Key("owner_type").HasValue("ServicePrincipal"),
					check.That(resourceType+".test_service_principal").Key("owner_display_name").HasValue("Test Service Principal Owner"),
				),
			},
			{
				ResourceName:      resourceType + ".test_service_principal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test_service_principal"),
			},
		},
	})
}
