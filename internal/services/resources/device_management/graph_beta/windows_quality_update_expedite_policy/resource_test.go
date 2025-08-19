package graphBetaWindowsQualityUpdateExpeditePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	expediteMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_quality_update_expedite_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *expediteMocks.WindowsQualityUpdateExpeditePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	expediteMock := &expediteMocks.WindowsQualityUpdateExpeditePolicyMock{}
	expediteMock.RegisterMocks()
	return mockClient, expediteMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *expediteMocks.WindowsQualityUpdateExpeditePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	expediteMock := &expediteMocks.WindowsQualityUpdateExpeditePolicyMock{}
	expediteMock.RegisterErrorMocks()
	return mockClient, expediteMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsQualityUpdateExpeditePolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal", "display_name", "Test Minimal Windows Quality Update Expedite Policy - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_MinimalToMax(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal", "display_name", "Test Minimal Windows Quality Update Expedite Policy - Unique"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.maximal", "display_name", "Test Maximal Windows Quality Update Expedite Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.maximal", "description", "Maximal Windows Quality Update Expedite Policy for testing with all features"),
				),
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_GroupAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.group_assignments", "display_name", "Test Group Assignments Windows Quality Update Expedite Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.group_assignments", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.group_assignments", "assignments.0.type", "groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.group_assignments", "assignments.1.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestWindowsQualityUpdateExpeditePolicyResource_ExclusionAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, expediteMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer expediteMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigExclusionAssignment(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.exclusion_assignment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.exclusion_assignment", "display_name", "Test Exclusion Assignment Windows Quality Update Expedite Policy - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.exclusion_assignment", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy.exclusion_assignment", "assignments.0.type", "exclusionGroupAssignmentTarget"),
				),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigGroupAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_group_assignments.tf")
	if err != nil {
		panic("failed to load group assignments config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigExclusionAssignment() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_exclusion_assignment.tf")
	if err != nil {
		panic("failed to load exclusion assignment config: " + err.Error())
	}
	return unitTestConfig
}
