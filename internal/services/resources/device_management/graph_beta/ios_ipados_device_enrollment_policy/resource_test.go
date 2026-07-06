package graphBetaIOSiPadOSDeviceEnrollmentPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	enrollmentPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/ios_ipados_device_enrollment_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment initializes the mock environment for testing
func setupMockEnvironment() (*mocks.Mocks, *enrollmentPolicyMocks.IOSiPadOSDeviceEnrollmentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	enrollmentPolicyMock := &enrollmentPolicyMocks.IOSiPadOSDeviceEnrollmentPolicyMock{}
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

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_01_Minimal(t *testing.T) {
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
					check.That(resourceType+".minimal").Key("name").HasValue("unit-test-ios-ade-minimal"),
					check.That(resourceType+".minimal").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".minimal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("platforms").HasValue("iOS"),
					check.That(resourceType+".minimal").Key("technologies").HasValue("enrollment"),
					check.That(resourceType+".minimal").Key("dep_onboarding_settings_id").HasValue(enrollmentPolicyMocks.DepOnboardingSettingsTestID),
					check.That(resourceType+".minimal").Key("device_security_group").DoesNotExist(),
					check.That(resourceType+".minimal").Key("device_name_template").DoesNotExist(),
					check.That(resourceType+".minimal").Key("cellular_data_activation_url").DoesNotExist(),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 02: Maximal Resource
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_02_Maximal(t *testing.T) {
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
					check.That(resourceType+".maximal").Key("name").HasValue("unit-test-ios-ade-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("iOS/iPadOS ADE enrollment policy exercising the full settings tree"),
					check.That(resourceType+".maximal").Key("requires_user_authentication").HasValue("true"),
					check.That(resourceType+".maximal").Key("require_setup_assistant_with_modern_authentication").HasValue("true"),
					check.That(resourceType+".maximal").Key("await_final_configuration").HasValue("true"),
					check.That(resourceType+".maximal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("device_name_template").HasValue("{{DEVICETYPE}}-{{SERIAL}}"),
					check.That(resourceType+".maximal").Key("cellular_data_activation_url").HasValue("http://activation.carrier.net"),
					check.That(resourceType+".maximal").Key("support_department").HasValue("IT Support"),
					check.That(resourceType+".maximal").Key("support_phone_number").HasValue("+1-555-0100"),
					check.That(resourceType+".maximal").Key("passcode_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("restore_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("apple_id_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("siri_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("web_content_filtering_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("location_services_disabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("app_store_disabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 03: Minimal to Maximal Update
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_03_MinimalToMaximal(t *testing.T) {
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
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-ios-ade-update"),
					check.That(resourceType+".update_test").Key("requires_user_authentication").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_scenario_minimal_to_maximal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-ios-ade-update-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("requires_user_authentication").HasValue("true"),
					check.That(resourceType+".update_test").Key("require_setup_assistant_with_modern_authentication").HasValue("true"),
					check.That(resourceType+".update_test").Key("await_final_configuration").HasValue("true"),
					check.That(resourceType+".update_test").Key("device_name_template").HasValue("{{DEVICETYPE}}-{{SERIAL}}"),
					check.That(resourceType+".update_test").Key("passcode_disabled").HasValue("true"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 04: Maximal to Minimal Update
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_04_MaximalToMinimal(t *testing.T) {
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
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-ios-ade-downgrade"),
					check.That(resourceType+".downgrade_test").Key("requires_user_authentication").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("enable_authentication_via_company_portal").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("device_name_template").HasValue("{{DEVICETYPE}}-{{SERIAL}}"),
					check.That(resourceType+".downgrade_test").Key("role_scope_tag_ids.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("006_scenario_maximal_to_minimal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-ios-ade-downgrade-minimal"),
					check.That(resourceType+".downgrade_test").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("enable_authentication_via_company_portal").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("device_name_template").DoesNotExist(),
					check.That(resourceType+".downgrade_test").Key("description").DoesNotExist(),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 05: Device Security Group (Enrollment Time Grouping)
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_05_DeviceSecurityGroup(t *testing.T) {
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

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_06_DeviceSecurityGroupUpdate(t *testing.T) {
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
// Scenario 07: Error Cases
//
// Note: description length, device_security_group GUID format, and the mutual exclusivity of the
// authentication method booleans are all enforced by stringvalidator/attribute-validator
// declarations in resource.go and fail during ValidateResourceConfig, before any provider CRUD or
// API logic runs - so a dedicated unit test for each would only be re-testing
// terraform-plugin-framework's own validator plumbing. The one kept here,
// requireUserAuthenticationForAuthenticationOptions, encodes the non-obvious settings catalog
// shape (ade_modernauth_awaitfinalconfiguration only exists under ade_authenticationmethod_2)
// rather than a generic "field must look like X" rule, so it's worth guarding against regression.
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_07_Error_AwaitFinalConfigurationRequiresModernAuth(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("010_scenario_error_await_final_configuration_requires_modern_auth.tf"),
				ExpectError: regexp.MustCompile(`require_setup_assistant_with_modern_authentication is required`),
			},
		},
	})
}

// ====================================================================================
// Scenario 08: Default Policy Assignment (setDefaultProfile action)
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_08_DefaultPolicyAssignment(t *testing.T) {
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
// Scenario 08b: Default Policy Assignment Unset Error (cannot flip to false while the policy is
// still the DEP token's current default - rejected by validateRequest during Update, since Graph
// has no unset action and nothing else in the plan promotes a replacement)
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_08b_DefaultPolicyAssignmentUnsetError(t *testing.T) {
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
					check.That(resourceType+".default_assignment").Key("is_default_policy_assignment").HasValue("true"),
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
// Scenario 09: Default Policy Assignment Switch (create default, then switch to another policy)
// ====================================================================================

func TestUnitResourceIOSiPadOSDeviceEnrollmentPolicy_09_DefaultPolicyAssignmentSwitch(t *testing.T) {
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
