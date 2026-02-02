package graphBetaServicePrincipalOwner_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalOwner "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_owner"
	servicePrincipalOwnerMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_owner/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalOwner.ResourceName

	// testResource is the test resource implementation for service principal owners
	testResource = graphBetaServicePrincipalOwner.ServicePrincipalOwnerTestResource{}
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *servicePrincipalOwnerMocks.ServicePrincipalOwnerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	servicePrincipalOwnerMock := &servicePrincipalOwnerMocks.ServicePrincipalOwnerMock{}
	servicePrincipalOwnerMock.RegisterMocks()
	return mockClient, servicePrincipalOwnerMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *servicePrincipalOwnerMocks.ServicePrincipalOwnerMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	servicePrincipalOwnerMock := &servicePrincipalOwnerMocks.ServicePrincipalOwnerMock{}
	servicePrincipalOwnerMock.RegisterErrorMocks()
	return mockClient, servicePrincipalOwnerMock
}

func TestUnitResourceServicePrincipalOwner_01_OwnerTypeUser(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, servicePrincipalOwnerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer servicePrincipalOwnerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_owner_type_user.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_user").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_user").Key("service_principal_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test_user").Key("owner_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_user").Key("owner_object_type").HasValue("User"),
					check.That(resourceType+".test_user").Key("owner_type").HasValue("User"),
					check.That(resourceType+".test_user").Key("owner_display_name").HasValue("Test User Owner"),
				),
			},
			{
				ResourceName:            resourceType + ".test_user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"owner_object_type"},
				ImportStateIdFunc:       testAccImportStateIdFunc(resourceType + ".test_user"),
			},
		},
	})
}

func TestUnitResourceServicePrincipalOwner_02_OwnerTypeServicePrincipal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, servicePrincipalOwnerMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer servicePrincipalOwnerMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_owner_type_service_principal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_service_principal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+/[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_service_principal").Key("service_principal_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test_service_principal").Key("owner_id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test_service_principal").Key("owner_object_type").HasValue("ServicePrincipal"),
					check.That(resourceType+".test_service_principal").Key("owner_type").HasValue("ServicePrincipal"),
					check.That(resourceType+".test_service_principal").Key("owner_display_name").HasValue("Test Service Principal Owner"),
				),
			},
			{
				ResourceName:            resourceType + ".test_service_principal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"owner_object_type"},
				ImportStateIdFunc:       testAccImportStateIdFunc(resourceType + ".test_service_principal"),
			},
		},
	})
}
