package graphBetaAppControlForBusinessPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.AppControlForBusinessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.AppControlForBusinessPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.AppControlForBusinessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.AppControlForBusinessPolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAppControlForBusinessPolicyResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "name", "unit-test-app-control-policy-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

func TestAppControlForBusinessPolicyResource_XMLPolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "name", "unit-test-app-control-policy-minimal"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "policy_xml"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal", "name", "unit-test-app-control-policy-maximal"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal", "policy_xml"),
				),
			},
		},
	})
}

func TestAppControlForBusinessPolicyResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "name", "unit-test-app-control-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "assignments.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "include",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "exclude",
					}),
				),
			},
		},
	})
}

func TestAppControlForBusinessPolicyResource_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_app_control_for_business_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAppControlForBusinessPolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid App Control for Business policy data"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_acfb_policy_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_acfb_policy_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}
