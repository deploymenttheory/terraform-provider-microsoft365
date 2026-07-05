package graphBetaApplicationsOnPremisesConnectorGroup_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	connectorGroupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/on_premises_connector_group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *connectorGroupMocks.OnPremisesConnectorGroupMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectorGroupMock := &connectorGroupMocks.OnPremisesConnectorGroupMock{}
	connectorGroupMock.RegisterMocks()
	return mockClient, connectorGroupMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceConnectorGroup_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigConnectorGroupMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_applications_on_premises_connector_group.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "name", "unit-test-connector-group"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "connector_group_type", "applicationProxy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "is_default", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "region", "japan"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:            "microsoft365_graph_beta_applications_on_premises_connector_group.minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceConnectorGroup_02_WithRegion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigConnectorGroupWithRegion(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_applications_on_premises_connector_group.with_region"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.with_region", "name", "unit-test-connector-group-region"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.with_region", "connector_group_type", "applicationProxy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.with_region", "is_default", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.with_region", "region", "nam"),
				),
			},
		},
	})
}

func TestUnitResourceConnectorGroup_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigConnectorGroupMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_applications_on_premises_connector_group.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "region", "japan"),
				),
			},
			{
				Config: testConfigConnectorGroupUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "name", "unit-test-connector-group-renamed"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group.minimal", "region", "eur"),
				),
			},
		},
	})
}

func testConfigConnectorGroupMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load connector group minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigConnectorGroupWithRegion() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_region.tf")
	if err != nil {
		panic("failed to load connector group region config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigConnectorGroupUpdated() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_updated.tf")
	if err != nil {
		panic("failed to load connector group updated config: " + err.Error())
	}
	return unitTestConfig
}
