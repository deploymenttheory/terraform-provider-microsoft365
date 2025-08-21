package graphBetaWindowsEnrollmentStatusPage_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	espMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_enrollment_status_page/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *espMocks.WindowsEnrollmentStatusPageMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	espMock := &espMocks.WindowsEnrollmentStatusPageMock{}
	espMock.RegisterMocks()
	return mockClient, espMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *espMocks.WindowsEnrollmentStatusPageMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	espMock := &espMocks.WindowsEnrollmentStatusPageMock{}
	espMock.RegisterErrorMocks()
	return mockClient, espMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestWindowsEnrollmentStatusPageResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "display_name", "unit-test-windows-enrollment-status-page-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "show_installation_progress", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "block_device_setup_retry_by_user", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "allow_device_reset_on_install_failure", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "install_progress_timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "custom_error_message", "Contact IT support for assistance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "selected_mobile_app_ids.#", "0"),
				),
			},
		},
	})
}

func TestWindowsEnrollmentStatusPageResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "display_name", "unit-test-windows-enrollment-status-page-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "description", "Test description for maximal enrollment status page"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "install_progress_timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "custom_error_message", "Contact IT support for assistance"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "track_install_progress_for_autopilot_only", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "selected_mobile_app_ids.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestWindowsEnrollmentStatusPageResource_WithAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigWithAssignments(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_windows_enrollment_status_page.full_with_group_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.full_with_group_assignments", "display_name", "unit-test-windows-enrollment-status-page-with-group-assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.full_with_group_assignments", "assignments.#", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.full_with_group_assignments", "selected_mobile_app_ids.#", "2"),
				),
			},
		},
	})
}


func TestWindowsEnrollmentStatusPageResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Enrollment Status Page data"),
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

func testConfigWithAssignments() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_with_group_assignments.tf")
	if err != nil {
		panic("failed to load assignments config: " + err.Error())
	}
	return unitTestConfig
}