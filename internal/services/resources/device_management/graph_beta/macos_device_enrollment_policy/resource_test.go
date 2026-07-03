package graphBetaMacOSDeviceEnrollmentPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	enrollmentPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_enrollment_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment initializes the mock environment for testing
func setupMockEnvironment() (*mocks.Mocks, *enrollmentPolicyMocks.MacOSDeviceEnrollmentPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	enrollmentPolicyMock := &enrollmentPolicyMocks.MacOSDeviceEnrollmentPolicyMock{}
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

func TestUnitResourceMacOSDeviceEnrollmentPolicy_01_Minimal(t *testing.T) {
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
					check.That(resourceType+".minimal").Key("name").HasValue("unit-test-macos-ade-minimal"),
					check.That(resourceType+".minimal").Key("await_device_configured").HasValue("false"),
					check.That(resourceType+".minimal").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".minimal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("platforms").HasValue("macOS"),
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

func TestUnitResourceMacOSDeviceEnrollmentPolicy_02_Maximal(t *testing.T) {
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
					check.That(resourceType+".maximal").Key("name").HasValue("unit-test-macos-ade-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("macOS ADE enrollment policy exercising the full settings tree"),
					check.That(resourceType+".maximal").Key("await_device_configured").HasValue("true"),
					check.That(resourceType+".maximal").Key("admin_account.user_name").HasValue("localadmin"),
					check.That(resourceType+".maximal").Key("admin_account.full_name").HasValue("Local Administrator"),
					check.That(resourceType+".maximal").Key("admin_account.hide_account").HasValue("true"),
					check.That(resourceType+".maximal").Key("admin_account.password_rotation_in_days").HasValue("90"),
					check.That(resourceType+".maximal").Key("admin_account.primary_account.user_name").HasValue("primaryuser"),
					check.That(resourceType+".maximal").Key("admin_account.primary_account.full_name").HasValue("Primary User"),
					check.That(resourceType+".maximal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("support_department").HasValue("IT Support"),
					check.That(resourceType+".maximal").Key("support_phone_number").HasValue("+1-555-0100"),
					check.That(resourceType+".maximal").Key("restore_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("apple_id_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("siri_disabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("file_vault_disabled").HasValue("false"),
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

func TestUnitResourceMacOSDeviceEnrollmentPolicy_03_MinimalToMaximal(t *testing.T) {
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
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-macos-ade-update"),
					check.That(resourceType+".update_test").Key("await_device_configured").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_scenario_minimal_to_maximal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("unit-test-macos-ade-update-updated"),
					check.That(resourceType+".update_test").Key("description").HasValue("Updated to maximal configuration"),
					check.That(resourceType+".update_test").Key("await_device_configured").HasValue("true"),
					check.That(resourceType+".update_test").Key("admin_account.user_name").HasValue("localadmin"),
					check.That(resourceType+".update_test").Key("admin_account.primary_account.user_name").HasValue("primaryuser"),
					check.That(resourceType+".update_test").Key("restore_disabled").HasValue("true"),
					check.That(resourceType+".update_test").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 04: Maximal to Minimal Update
// ====================================================================================

func TestUnitResourceMacOSDeviceEnrollmentPolicy_04_MaximalToMinimal(t *testing.T) {
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
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-macos-ade-downgrade"),
					check.That(resourceType+".downgrade_test").Key("await_device_configured").HasValue("true"),
					check.That(resourceType+".downgrade_test").Key("role_scope_tag_ids.#").HasValue("3"),
				),
			},
			{
				Config: loadUnitTestTerraform("006_scenario_maximal_to_minimal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("unit-test-macos-ade-downgrade-minimal"),
					check.That(resourceType+".downgrade_test").Key("await_device_configured").HasValue("false"),
					check.That(resourceType+".downgrade_test").Key("description").DoesNotExist(),
				),
			},
		},
	})
}

// ====================================================================================
// Scenario 05: Device Security Group (Enrollment Time Grouping)
// ====================================================================================

func TestUnitResourceMacOSDeviceEnrollmentPolicy_05_DeviceSecurityGroup(t *testing.T) {
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

func TestUnitResourceMacOSDeviceEnrollmentPolicy_06_DeviceSecurityGroupUpdate(t *testing.T) {
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
// Note: description length, device_security_group GUID format, and admin_account presence are
// all enforced by stringvalidator/ConfigValidator declarations in resource.go/validate.go and
// fail during ValidateResourceConfig, before any provider CRUD or API logic runs - so a dedicated
// unit test for each would only be re-testing terraform-plugin-framework's own validator
// plumbing. The one exception kept here, requireAuthenticationMethodWhenUserAuthRequired, encodes
// a non-obvious live Graph API behavior (rejection of ade_macos_authenticationmethod_0) rather
// than a generic "field must look like X" rule, so it's worth guarding against regression.
// ====================================================================================

func TestUnitResourceMacOSDeviceEnrollmentPolicy_07_Error_AuthenticationMethodRequired(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("010_scenario_error_authentication_method_required.tf"),
				ExpectError: regexp.MustCompile(`an authentication method is required`),
			},
		},
	})
}
