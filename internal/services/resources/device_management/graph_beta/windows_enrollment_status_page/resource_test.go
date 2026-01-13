package graphBetaWindowsEnrollmentStatusPage_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsEnrollmentStatusPage "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_enrollment_status_page"
	espMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_enrollment_status_page/mocks"
	groupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *espMocks.WindowsEnrollmentStatusPageMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	espMock := &espMocks.WindowsEnrollmentStatusPageMock{}
	espMock.RegisterMocks()
	return mockClient, espMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *espMocks.WindowsEnrollmentStatusPageMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register group mocks for tests that create groups
	groupMock := &groupMocks.GroupMock{}
	groupMock.RegisterMocks()

	espMock := &espMocks.WindowsEnrollmentStatusPageMock{}
	espMock.RegisterErrorMocks()
	return mockClient, espMock
}

// Test 001: Minimal Configuration
func TestWindowsEnrollmentStatusPageResource_001_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-minimal"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("description").HasValue("Test description for minimal enrollment status page"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("show_installation_progress").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("install_quality_updates").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("allow_log_collection_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("allow_device_use_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("only_fail_selected_blocking_apps_in_technician_phase").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("install_progress_timeout_in_minutes").HasValue("120"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("custom_error_message").HasValue("Contact IT support for assistance"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("selected_mobile_app_ids.#").HasValue("0"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".minimal").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

// Test 002: Maximal Configuration
func TestWindowsEnrollmentStatusPageResource_002_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-maximal"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("description").HasValue("Test description for maximal enrollment status page"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("show_installation_progress").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_quality_updates").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_log_collection_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_use_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_fail_selected_blocking_apps_in_technician_phase").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_progress_timeout_in_minutes").HasValue("120"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("custom_error_message").HasValue("Contact IT support for assistance"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 003: Configuration with Group Assignments
func TestWindowsEnrollmentStatusPageResource_003_WithAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_with_group_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-maximal"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("description").HasValue("Test description for maximal enrollment status page"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("show_installation_progress").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_quality_updates").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_log_collection_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("allow_device_use_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("only_fail_selected_blocking_apps_in_technician_phase").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("install_progress_timeout_in_minutes").HasValue("120"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("custom_error_message").HasValue("Contact IT support for assistance"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("assignments.#").HasValue("4"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".maximal").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 004: Error and Validation Testing
func TestWindowsEnrollmentStatusPageResource_004_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid Windows Enrollment Status Page data"),
			},
		},
	})
}

// Test 005: Lifecycle Minimal to Maximal
func TestWindowsEnrollmentStatusPageResource_005_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-lifecycle"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-lifecycle"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 006: Lifecycle Maximal to Minimal
func TestWindowsEnrollmentStatusPageResource_006_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, espMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer espMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-lifecycle"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("selected_mobile_app_ids.#").HasValue("3"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("display_name").HasValue("unit-test-windows-enrollment-status-page-lifecycle"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("block_device_use_until_all_apps_and_profiles_are_installed").HasValue("true"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("allow_device_reset_on_install_failure").HasValue("false"),
					check.That(graphBetaWindowsEnrollmentStatusPage.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}
