package graphBetaApplicationsOnPremisesConnectorGroupAssignment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/on_premises_connector_group_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *assignmentMocks.OnPremisesConnectorGroupAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	assignmentMock := &assignmentMocks.OnPremisesConnectorGroupAssignmentMock{}
	assignmentMock.RegisterMocks()
	return mockClient, assignmentMock
}

func TestUnitResourceConnectorGroupAssignment_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, assignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigConnectorGroupAssignmentMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group_assignment.minimal", "application_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group_assignment.minimal", "connector_group_id", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group_assignment.minimal", "connector_group_name", "Unit Test Connector Group"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_applications_on_premises_connector_group_assignment.minimal", "id", regexp.MustCompile(`^11111111-1111-1111-1111-111111111111/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa$`)),
				),
			},
			{
				ResourceName:            "microsoft365_graph_beta_applications_on_premises_connector_group_assignment.minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func testConfigConnectorGroupAssignmentMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load connector group assignment minimal config: " + err.Error())
	}
	return unitTestConfig
}
