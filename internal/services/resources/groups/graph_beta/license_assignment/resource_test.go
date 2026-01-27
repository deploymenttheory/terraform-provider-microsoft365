package graphBetaGroupLicenseAssignment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	localMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/license_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *localMocks.GroupLicenseAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	licenseAssignmentMock := &localMocks.GroupLicenseAssignmentMock{}
	licenseAssignmentMock.RegisterMocks()
	return mockClient, licenseAssignmentMock
}

func testConfigError() string {
	return `
resource "microsoft365_graph_beta_groups_license_assignment" "error" {
  group_id = "invalid-group-id"
  sku_id = "33333333-3333-3333-3333-333333333333"
}
`
}

func TestUnitResourceGroupLicenseAssignment_01_CreateMinimal(t *testing.T) {
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
					check.That(resourceType+".minimal").Key("group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
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

func TestUnitResourceGroupLicenseAssignment_02_CreateMaximal(t *testing.T) {
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
					check.That(resourceType+".maximal").Key("group_id").HasValue("00000000-0000-0000-0000-000000000003"),
					check.That(resourceType+".maximal").Key("sku_id").HasValue("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".maximal").Key("disabled_plans.#").HasValue("2"),
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

func TestUnitResourceGroupLicenseAssignment_03_Update(t *testing.T) {
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
					check.That(resourceType+".minimal").Key("group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
				),
			},
			{
				Config: testConfigMinimalWithDisabledPlans(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("group_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".minimal").Key("sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
					check.That(resourceType+".minimal").Key("disabled_plans.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceGroupLicenseAssignment_04_DeleteMinimal(t *testing.T) {
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

func TestUnitResourceGroupLicenseAssignment_05_DeleteMaximal(t *testing.T) {
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

func TestUnitResourceGroupLicenseAssignment_06_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, licenseAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	licenseAssignmentMock.RegisterErrorMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile("must be a valid GUID"),
			},
		},
	})
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

func testConfigMinimalWithDisabledPlans() string {
	return `
resource "microsoft365_graph_beta_groups_license_assignment" "minimal" {
  group_id = "00000000-0000-0000-0000-000000000002"
  sku_id  = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE
  disabled_plans = [
    "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"
  ]
}
`
}
