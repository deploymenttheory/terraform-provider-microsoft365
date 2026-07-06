package graphBetaVisionOSDeviceEnrollmentPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	enrollmentPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/visionos_device_enrollment_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment initializes the mock environment for testing
func setupMockEnvironment() (*mocks.Mocks, *enrollmentPolicyMocks.VisionOSDeviceEnrollmentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	enrollmentPolicyMock := &enrollmentPolicyMocks.VisionOSDeviceEnrollmentPolicyMock{}
	enrollmentPolicyMock.RegisterMocks()
	return mockClient, enrollmentPolicyMock
}

// loadUnitTestTerraform loads a Terraform configuration file for unit testing
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// ====================================================================================
// Scenario 01: Minimal Resource
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("name").HasValue("unit-test-visionos-ade-minimal"),
					check.That(resourceType+".minimal").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".minimal").Key("await_device_configured").HasValue("true"),
					check.That(resourceType+".minimal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("platforms").HasValue("visionOS"),
					check.That(resourceType+".minimal").Key("technologies").HasValue("enrollment"),
					check.That(resourceType+".minimal").Key("dep_onboarding_settings_id").HasValue(enrollmentPolicyMocks.DepOnboardingSettingsTestID),
					check.That(resourceType+".minimal").Key("device_security_group").DoesNotExist(),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 02: Maximal Resource
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("name").HasValue("unit-test-visionos-ade-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("visionOS ADE enrollment policy exercising the full settings tree"),
					check.That(resourceType+".maximal").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".maximal").Key("await_device_configured").HasValue("true"),
					check.That(resourceType+".maximal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("support_department").HasValue("IT Support"),
					check.That(resourceType+".maximal").Key("support_phone_number").HasValue("+1-555-0100"),
					check.That(resourceType+".maximal").Key("passcode_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("apple_id_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("siri_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("tips_screen_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("location_services_disabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("touch_id_disabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 03: Minimal to Maximal Update
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_03_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_scenario_minimal_to_maximal_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-visionos-ade-update"),
					check.That(resourceType+".update_test").Key("requires_user_authentication").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_scenario_minimal_to_maximal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-visionos-ade-update-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".update_test").Key("passcode_disabled").HasValue("true"),
					check.That(resourceType+".update_test").Key("tips_screen_disabled").HasValue("true"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 04: Maximal to Minimal Update
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_04_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_scenario_maximal_to_minimal_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-visionos-ade-downgrade"),
					check.That(resourceType+".downgrade_test").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("role_scope_tag_ids.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("006_scenario_maximal_to_minimal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-visionos-ade-downgrade-minimal"),
					check.That(resourceType+".downgrade_test").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("passcode_disabled").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("description").DoesNotExist(),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 05: Device Security Group (Enrollment Time Grouping)
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_05_DeviceSecurityGroup(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_scenario_device_security_group.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".device_group").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".device_group").Key("device_security_group").HasValue("10000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 06: Device Security Group Update (clear then set)
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_06_DeviceSecurityGroupUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_scenario_device_security_group_update_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".device_group_update").Key("device_security_group").HasValue("10000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: loadUnitTestTerraform("009_scenario_device_security_group_update_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".device_group_update").Key("device_security_group").HasValue("10000000-0000-0000-0000-000000000009"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 07: Default Policy Assignment (setDefaultProfile action)
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_07_DefaultPolicyAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("011_scenario_default_policy_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".default_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".default_assignment").Key("is_default_policy_assignment").HasValue("true"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 07b: Default Policy Assignment Unset Error (cannot flip to false while the policy is
// still the DEP token's current default - rejected by validateRequest during Update, since Graph
// has no unset action and nothing else in the plan promotes a replacement)
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_07b_DefaultPolicyAssignmentUnsetError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("011_scenario_default_policy_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".default_assignment").Key("is_default_policy_assignment").HasValue("true"),
				),
			},
			{
				Config:      loadUnitTestTerraform("014_scenario_default_policy_assignment_unset_error.tf"),
				ExpectError: regexp.MustCompile(`cannot unset the default policy assignment`),
			},
		},
	})
}

// ====================================================================================
// Scenario 08: Default Policy Assignment Switch (create default, then switch to another policy)
// ====================================================================================

func TestUnitResourceVisionOSDeviceEnrollmentPolicy_08_DefaultPolicyAssignmentSwitch(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("012_scenario_default_policy_assignment_switch_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".switch_a").Key("is_default_policy_assignment").HasValue("true"),
					check.That(resourceType+".switch_b").Key("is_default_policy_assignment").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("013_scenario_default_policy_assignment_switch_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".switch_a").Key("is_default_policy_assignment").HasValue("false"),
					check.That(resourceType+".switch_b").Key("is_default_policy_assignment").HasValue("true"),
				),
			},
		},
	})
}
