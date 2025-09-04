package graphBetaAppControlForBusinessBuiltInControls_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	appControlMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *appControlMocks.AppControlForBusinessBuiltInControlsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appControlMock := &appControlMocks.AppControlForBusinessBuiltInControlsMock{}
	appControlMock.RegisterMocks()
	return mockClient, appControlMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *appControlMocks.AppControlForBusinessBuiltInControlsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appControlMock := &appControlMocks.AppControlForBusinessBuiltInControlsMock{}
	appControlMock.RegisterErrorMocks()
	return mockClient, appControlMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAppControlForBusinessBuiltInControlsResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "name", "unit-test-app-control-for-business-built-in-controls-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "enable_app_control", "audit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

func TestAppControlForBusinessBuiltInControlsResource_EnableAppControl(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAuditMode(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.audit_mode"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.audit_mode", "name", "unit-test-app-control-audit-mode"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.audit_mode", "enable_app_control", "audit"),
				),
			},
			{
				Config: testConfigEnforceMode(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.enforce_mode"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.enforce_mode", "name", "unit-test-app-control-enforce-mode"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.enforce_mode", "enable_app_control", "enforce"),
				),
			},
		},
	})
}

func TestAppControlForBusinessBuiltInControlsResource_AdditionalRules(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "name", "unit-test-app-control-for-business-built-in-controls-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "additional_rules_for_trusting_apps.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "additional_rules_for_trusting_apps.*", "trust_apps_with_good_reputation"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.advanced", "additional_rules_for_trusting_apps.*", "trust_apps_from_managed_installers"),
				),
			},
		},
	})
}

func TestAppControlForBusinessBuiltInControlsResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments", "name", "unit-test-app-control-with-assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments", "assignments.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "include",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls.with_assignments", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "exclude",
					}),
				),
			},
		},
	})
}

func TestAppControlForBusinessBuiltInControlsResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid App Control for Business data"),
			},
		},
	})
}

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_acfb_built_in_controls_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_acfb_built_in_controls_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAuditMode() string {
	return `
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "audit_mode" {
  name        = "unit-test-app-control-audit-mode"
  description = "unit-test-app-control-audit-mode"
  
  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}`
}

func testConfigEnforceMode() string {
	return `
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "enforce_mode" {
  name        = "unit-test-app-control-enforce-mode"
  description = "unit-test-app-control-enforce-mode"
  
  enable_app_control = "enforce"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}`
}

func testConfigWithAssignments() string {
	return `
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "with_assignments" {
  name        = "unit-test-app-control-with-assignments"
  description = "unit-test-app-control-with-assignments"
  
  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}`
}