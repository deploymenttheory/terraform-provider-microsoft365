package graphBetaApplicationIdentifierUri_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationIdentifierUri "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_identifier_uri"
	applicationIdentifierUriMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/application_identifier_uri/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationIdentifierUri.ResourceName

	// testResource is the test resource implementation for application identifier URIs
	testResource = graphBetaApplicationIdentifierUri.ApplicationIdentifierUriTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *applicationIdentifierUriMocks.ApplicationIdentifierUriMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationIdentifierUriMock := &applicationIdentifierUriMocks.ApplicationIdentifierUriMock{}
	applicationIdentifierUriMock.RegisterMocks()
	return mockClient, applicationIdentifierUriMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *applicationIdentifierUriMocks.ApplicationIdentifierUriMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	applicationIdentifierUriMock := &applicationIdentifierUriMocks.ApplicationIdentifierUriMock{}
	applicationIdentifierUriMock.RegisterErrorMocks()
	return mockClient, applicationIdentifierUriMock
}

func TestUnitResourceApplicationIdentifierUri_01(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, applicationIdentifierUriMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer applicationIdentifierUriMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("application_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test").Key("identifier_uri").HasValue("api://11111111-1111-1111-1111-111111111111"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceType + ".test"),
			},
		},
	})
}

// testAccImportStateIdFunc returns a function that constructs the import ID
func testAccImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		appID := rs.Primary.Attributes["application_id"]
		uri := rs.Primary.Attributes["identifier_uri"]
		return fmt.Sprintf("%s/%s", appID, uri), nil
	}
}
