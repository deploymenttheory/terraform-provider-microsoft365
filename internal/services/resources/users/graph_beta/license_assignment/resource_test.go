package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUserLicenseAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

var resourceType = graphBetaUserLicenseAssignment.ResourceName

func setupMockEnvironment() (*mocks.Mocks, *localMocks.UserLicenseAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	licenseAssignmentMock := &localMocks.UserLicenseAssignmentMock{}
	licenseAssignmentMock.RegisterMocks()
	return mockClient, licenseAssignmentMock
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load resource_minimal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load resource_maximal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMinimalToMaximal() string {
	maximalContent := testConfigMaximal()
	updatedMaximal := strings.Replace(maximalContent, "maximal", "minimal", 1)
	updatedMaximal = strings.Replace(updatedMaximal, "00000000-0000-0000-0000-000000000003", "00000000-0000-0000-0000-000000000002", 1)
	return updatedMaximal
}

func testConfigMaximalWithResourceName(resourceName string) string {
	maximalContent := testConfigMaximal()
	return strings.Replace(maximalContent, "maximal", resourceName, 1)
}

func testConfigMinimalWithResourceName(resourceName string) string {
	return fmt.Sprintf(`resource "microsoft365_graph_beta_users_user_license_assignment" "%s" {
  user_id = "00000000-0000-0000-0000-000000000003"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}`, resourceName)
}

func testConfigError() string {
	return `
resource "microsoft365_graph_beta_users_user_license_assignment" "error" {
  user_id = "invalid-user-id"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
}
`
}

func TestUnitUserLicenseAssignmentResource_Create_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("add_licenses.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("add_licenses.0.sku_id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".minimal").Key("add_licenses.0.disabled_plans.#").HasValue("0"),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Create_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000003"),
					check.That(resourceType+".maximal").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("add_licenses.0.sku_id").HasValue("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".maximal").Key("add_licenses.0.disabled_plans.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("add_licenses.1.sku_id").HasValue("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".maximal").Key("remove_licenses.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("remove_licenses.0").HasValue("88888888-8888-8888-8888-888888888888"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Update_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("add_licenses.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("add_licenses.0.sku_id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".minimal").Key("add_licenses.0.disabled_plans.#").HasValue("0"),
				),
			},
			{
				Config: testConfigMinimalToMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".minimal").Key("add_licenses.0.sku_id").HasValue("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".minimal").Key("add_licenses.0.disabled_plans.#").HasValue("2"),
					check.That(resourceType+".minimal").Key("add_licenses.1.sku_id").HasValue("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".minimal").Key("remove_licenses.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Update_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".test").Key("remove_licenses.#").HasValue("1"),
				),
			},
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("add_licenses.#").HasValue("1"),
					check.That(resourceType+".test").Key("add_licenses.0.sku_id").HasValue("33333333-3333-3333-3333-333333333333"),
				),
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Delete_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".minimal").Key("id").Exists(),
				),
			},
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					_, exists := s.RootModule().Resources[resourceType+".minimal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Delete_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".maximal").Key("id").Exists(),
				),
			},
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					_, exists := s.RootModule().Resources[resourceType+".maximal"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

func TestUnitUserLicenseAssignmentResource_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, licenseAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	licenseAssignmentMock.RegisterErrorMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("Attribute user_id Must be a valid UUID format"),
			},
		},
	})
}
