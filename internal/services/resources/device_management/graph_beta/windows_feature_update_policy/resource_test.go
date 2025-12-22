package graphBetaWindowsFeatureUpdatePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsFeatureUpdatePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_feature_update_policy"
	featureMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_feature_update_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *featureMocks.WindowsFeatureUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	featureMock := &featureMocks.WindowsFeatureUpdateProfileMock{}
	featureMock.RegisterMocks()
	return mockClient, featureMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *featureMocks.WindowsFeatureUpdateProfileMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	featureMock := &featureMocks.WindowsFeatureUpdateProfileMock{}
	featureMock.RegisterErrorMocks()
	return mockClient, featureMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestWindowsFeatureUpdatePolicyResource_001_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("display_name").HasValue("unit-test-windows-feature-update-policy-001-minimal"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("feature_update_version").HasValue("Windows 11, version 23H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("install_feature_updates_optional").HasValue("false"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("install_latest_windows10_on_windows11_ineligible_device").HasValue("false"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
				),
			},
			{
				ResourceName:      graphBetaWindowsFeatureUpdatePolicy.ResourceName + ".test_001",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_002_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("display_name").HasValue("unit-test-windows-feature-update-policy-002-maximal"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("description").HasValue("Maximal test configuration for Windows feature update policy"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("feature_update_version").HasValue("Windows 11, version 24H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("install_feature_updates_optional").HasValue("true"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("install_latest_windows10_on_windows11_ineligible_device").HasValue("true"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_002").Key("rollout_settings.offer_start_date_time_in_utc").HasValue("2025-04-01T00:00:00Z"),
				),
			},
			{
				ResourceName:      graphBetaWindowsFeatureUpdatePolicy.ResourceName + ".test_002",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_003_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-feature-update-policy-003-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("feature_update_version").HasValue("Windows 11, version 23H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("install_feature_updates_optional").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("display_name").HasValue("unit-test-windows-feature-update-policy-003-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("feature_update_version").HasValue("Windows 11, version 24H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("install_feature_updates_optional").HasValue("true"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_004_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-feature-update-policy-004-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("description").HasValue("Maximal lifecycle test configuration"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("feature_update_version").HasValue("Windows 11, version 24H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("install_feature_updates_optional").HasValue("true"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("display_name").HasValue("unit-test-windows-feature-update-policy-004-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("feature_update_version").HasValue("Windows 11, version 23H2"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("install_feature_updates_optional").HasValue("false"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_005_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_005").Key("display_name").HasValue("unit-test-windows-feature-update-policy-005-assignments-minimal"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
				),
			},
			{
				ResourceName:      graphBetaWindowsFeatureUpdatePolicy.ResourceName + ".test_005",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_006_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_006").Key("display_name").HasValue("unit-test-windows-feature-update-policy-006-assignments-maximal"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_006").Key("description").HasValue("Maximal test with multiple assignments"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_006").Key("assignments.#").HasValue("3"),
				),
			},
			{
				ResourceName:      graphBetaWindowsFeatureUpdatePolicy.ResourceName + ".test_006",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-feature-update-policy-007-assignments-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("display_name").HasValue("unit-test-windows-feature-update-policy-007-assignments-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("3"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-feature-update-policy-008-assignments-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("description").HasValue("Maximal assignments lifecycle test"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("display_name").HasValue("unit-test-windows-feature-update-policy-008-assignments-lifecycle"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("description").HasValue("Maximal assignments lifecycle test"),
					check.That(graphBetaWindowsFeatureUpdatePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

func TestWindowsFeatureUpdatePolicyResource_009_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, featureMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer featureMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Windows Feature Update Profile data"),
			},
		},
	})
}
