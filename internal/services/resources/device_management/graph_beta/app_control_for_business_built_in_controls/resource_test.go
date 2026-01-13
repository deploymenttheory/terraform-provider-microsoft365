package graphBetaAppControlForBusinessBuiltInControls_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAppControlForBusinessBuiltInControls "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls"
	appControlMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_built_in_controls/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *appControlMocks.AppControlForBusinessBuiltInControlsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appControlMock := &appControlMocks.AppControlForBusinessBuiltInControlsMock{}
	appControlMock.RegisterMocks()
	return mockClient, appControlMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *appControlMocks.AppControlForBusinessBuiltInControlsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appControlMock := &appControlMocks.AppControlForBusinessBuiltInControlsMock{}
	appControlMock.RegisterErrorMocks()
	return mockClient, appControlMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 001: Audit Mode
func TestAppControlForBusinessBuiltInControlsResource_001_AuditMode(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_audit_mode.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("name").HasValue("unit-test-app-control-audit-mode"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".audit_mode").Key("enable_app_control").HasValue("audit"),
				),
			},
		},
	})
}

// Test 002: Enforce Mode
func TestAppControlForBusinessBuiltInControlsResource_002_EnforceMode(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_enforce_mode.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("name").HasValue("unit-test-app-control-enforce-mode"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".enforce_mode").Key("enable_app_control").HasValue("enforce"),
				),
			},
		},
	})
}

// Test 003: Minimal Configuration
func TestAppControlForBusinessBuiltInControlsResource_003_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("name").HasValue("unit-test-app-control-for-business-built-in-controls-minimal"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("enable_app_control").HasValue("audit"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("role_scope_tag_ids.#").HasValue("3"),
				),
			},
		},
	})
}

// Test 004: Maximal Configuration with Additional Rules
func TestAppControlForBusinessBuiltInControlsResource_004_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("name").HasValue("unit-test-app-control-for-business-built-in-controls-maximal"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".advanced").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 005: Lifecycle - Minimal to Maximal
func TestAppControlForBusinessBuiltInControlsResource_005_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("enable_app_control").HasValue("audit"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".lifecycle").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
				),
			},
		},
	})
}

// Test 006: Lifecycle - Maximal to Minimal (Downgrade)
func TestAppControlForBusinessBuiltInControlsResource_006_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("role_scope_tag_ids.#").HasValue("3"),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("additional_rules_for_trusting_apps.#").HasValue("2"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".downgrade").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 007: Assignments Lifecycle - Minimal to Maximal
func TestAppControlForBusinessBuiltInControlsResource_007_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_assignments_lifecycle_step_1_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_assignments_lifecycle_step_2_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_lifecycle").Key("assignments.#").HasValue("3"),
				),
			},
		},
	})
}

// Test 008: Assignments Lifecycle - Maximal to Minimal (Downgrade)
func TestAppControlForBusinessBuiltInControlsResource_008_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_assignments_downgrade_step_1_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("assignments.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_acfb_built_in_controls_assignments_downgrade_step_2_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaAppControlForBusinessBuiltInControls.ResourceName+".assignments_downgrade").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}

// Test 009: Error Handling
func TestAppControlForBusinessBuiltInControlsResource_009_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, appControlMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer appControlMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_acfb_built_in_controls_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid App Control for Business data"),
			},
		},
	})
}
