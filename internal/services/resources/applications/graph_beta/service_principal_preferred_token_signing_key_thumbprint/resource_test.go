package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_preferred_token_signing_key_thumbprint"
	thumbprintMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_preferred_token_signing_key_thumbprint/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint.ResourceName

	// testResource is the test resource implementation for the preferred token signing key thumbprint
	testResource = graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint.ServicePrincipalPreferredTokenSigningKeyThumbprintTestResource{}
)

const testServicePrincipalID = "11111111-1111-1111-1111-111111111111"

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *thumbprintMocks.ServicePrincipalPreferredTokenSigningKeyThumbprintMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	thumbprintMock := &thumbprintMocks.ServicePrincipalPreferredTokenSigningKeyThumbprintMock{}
	thumbprintMock.RegisterMocks()
	return mockClient, thumbprintMock
}

func TestUnitResourceServicePrincipalPreferredTokenSigningKeyThumbprint_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, thumbprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer thumbprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		// Destroy must clear the property via an explicit JSON null; the mock only
		// removes the stored thumbprint when the PATCH body contains the null key.
		CheckDestroy: func(s *terraform.State) error {
			if thumbprint, ok := thumbprintMocks.GetThumbprint(testServicePrincipalID); ok {
				return fmt.Errorf("preferredTokenSigningKeyThumbprint still set after destroy: %s", thumbprint)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("service_principal_id").HasValue(testServicePrincipalID),
					check.That(resourceType+".test").Key("thumbprint").HasValue("aabbccddeeff00112233445566778899aabbccdd"),
					check.That(resourceType+".test").Key("id").HasValue(testServicePrincipalID),
				),
			},
			{
				// Certificate rotation: a new thumbprint updates the service principal in place
				Config: loadUnitTestTerraform("resource_02_rotated.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("thumbprint").HasValue("ffee00112233445566778899aabbccddeeff0011"),
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

func TestUnitResourceServicePrincipalPreferredTokenSigningKeyThumbprint_02_CasingPreserved(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, thumbprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer thumbprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Graph stores the thumbprint lowercased; the configured uppercase value
				// must survive Read without producing a perpetual diff.
				Config: loadUnitTestTerraform("resource_03_uppercase.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("thumbprint").HasValue("AABBCCDDEEFF00112233445566778899AABBCCDD"),
				),
			},
			{
				Config:   loadUnitTestTerraform("resource_03_uppercase.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestUnitResourceServicePrincipalPreferredTokenSigningKeyThumbprint_03_InvalidThumbprint(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, thumbprintMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer thumbprintMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_04_invalid.tf"),
				ExpectError: regexp.MustCompile(`must be a 40-character hexadecimal SHA-1 certificate\s+thumbprint`),
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
		return rs.Primary.Attributes["service_principal_id"], nil
	}
}
