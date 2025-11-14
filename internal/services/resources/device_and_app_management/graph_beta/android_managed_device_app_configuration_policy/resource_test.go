package graphBetaAndroidManagedDeviceAppConfigurationPolicy_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	androidManagedDeviceAppConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/android_managed_device_app_configuration_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *androidManagedDeviceAppConfigurationMocks.AndroidManagedDeviceAppConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	androidManagedDeviceAppConfigurationMock := &androidManagedDeviceAppConfigurationMocks.AndroidManagedDeviceAppConfigurationMock{}
	androidManagedDeviceAppConfigurationMock.RegisterMocks()

	return mockClient, androidManagedDeviceAppConfigurationMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *androidManagedDeviceAppConfigurationMocks.AndroidManagedDeviceAppConfigurationMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	androidManagedDeviceAppConfigurationMock := &androidManagedDeviceAppConfigurationMocks.AndroidManagedDeviceAppConfigurationMock{}
	androidManagedDeviceAppConfigurationMock.RegisterErrorMocks()

	return mockClient, androidManagedDeviceAppConfigurationMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// Helper functions to load each test configuration
func testConfigMicrosoftAuthenticator() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_authenticator_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoft365Copilot() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_365_copilot_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigManagedHomeScreen() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_managed_home_screen_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftDefender() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_defender_antivirus_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftEdge() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_edge_browser_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftExcel() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_excel_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftOneDrive() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_onedrive_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftOneNote() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_onenote_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftOutlook() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_outlook_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftPowerPoint() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_powerpoint_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftTeams() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_teams_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

func testConfigMicrosoftWord() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_microsoft_word_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_Schema validates the resource schema
func TestAndroidManagedDeviceAppConfigurationPolicyResource_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "display_name", "unit-test-android-managed-device-app-configuration-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "description", "Unit test Android managed store app configuration"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "package_id", "app:com.microsoft.office.officehubrow"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "payload_json"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "app_supports_oem_config", "false"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_Minimal tests basic CRUD operations
func TestAndroidManagedDeviceAppConfigurationPolicyResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "display_name", "unit-test-android-managed-device-app-configuration-policy-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "package_id", "app:com.microsoft.office.officehubrow"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "payload_json"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal", "version"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftAuthenticator tests Microsoft Authenticator configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftAuthenticator(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftAuthenticator(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-authenticator-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "package_id", "app:com.azure.authenticator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "role_scope_tag_ids.*", "0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_authenticator_maximal", "app_supports_oem_config", "false"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_Microsoft365Copilot tests Microsoft 365 Copilot configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_Microsoft365Copilot(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoft365Copilot(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-365-copilot-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "package_id", "app:com.microsoft.office.officehubrow"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_365_copilot_maximal", "version"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_ManagedHomeScreen tests Managed Home Screen kiosk configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_ManagedHomeScreen(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigManagedHomeScreen(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-managed-home-screen-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "package_id", "app:com.microsoft.launcher.enterprise"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "profile_applicability", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.managed_home_screen_maximal", "version"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftDefender tests Microsoft Defender configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftDefender(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftDefender(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-defender-antivirus-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "package_id", "app:com.microsoft.scmx"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "profile_applicability", "default"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_defender_antivirus_maximal", "version"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftEdge tests Microsoft Edge browser configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftEdge(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftEdge(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-edge-browser-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "package_id", "app:com.microsoft.emmx"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_edge_browser_maximal", "version"),
				),
			},
		},
	})
}

// NOTE: Office 365 apps (Excel, PowerPoint, Word, OneNote, OneDrive) are now tested individually
// in their own test functions below for better isolation and failure reporting.

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftExcel tests Microsoft Excel configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftExcel(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftExcel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "package_id", "app:com.microsoft.office.excel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_excel_maximal", "permission_actions.#", "33"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftPowerPoint tests Microsoft PowerPoint configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftPowerPoint(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftPowerPoint(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "package_id", "app:com.microsoft.office.powerpoint"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_powerpoint_maximal", "permission_actions.#", "33"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftWord tests Microsoft Word configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftWord(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftWord(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "package_id", "app:com.microsoft.office.word"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_word_maximal", "permission_actions.#", "33"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneNote tests Microsoft OneNote configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneNote(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftOneNote(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "package_id", "app:com.microsoft.office.onenote"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onenote_maximal", "permission_actions.#", "33"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneDrive tests Microsoft OneDrive configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOneDrive(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftOneDrive(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "package_id", "app:com.microsoft.skydrive"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_onedrive_maximal", "permission_actions.#", "33"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOutlook tests Outlook email configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftOutlook(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftOutlook(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-outlook-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "package_id", "app:com.microsoft.office.outlook"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "profile_applicability", "androidDeviceOwner"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_outlook_maximal", "version"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftTeams tests Teams collaboration configuration
func TestAndroidManagedDeviceAppConfigurationPolicyResource_MicrosoftTeams(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMicrosoftTeams(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "display_name", "unit-test-android-managed-device-app-configuration-policy-microsoft-teams-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "package_id", "app:com.microsoft.teams"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "profile_applicability", "androidWorkProfile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "connected_apps_enabled", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "payload_json"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "permission_actions.#", "33"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "targeted_mobile_apps.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.microsoft_teams_maximal", "version"),
				),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_RequiredFields tests required field validation
func TestAndroidManagedDeviceAppConfigurationPolicyResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "test" {
  # Missing display_name
  package_id = "app:com.test"
  payload_json = jsonencode({})
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_ErrorHandling tests error scenarios
func TestAndroidManagedDeviceAppConfigurationPolicyResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "test" {
  display_name = "Test Android Managed Device App Configuration Policy"
  targeted_mobile_apps = ["00000000-0000-0000-0000-000000000000"]
  package_id = "app:com.test"
  payload_json = jsonencode({})
}
`,
				ExpectError: regexp.MustCompile(`(Invalid Android managed store app configuration data|BadRequest|failed to query Android mobile apps from Intune|no responder found)`),
			},
		},
	})
}

// TestAndroidManagedDeviceAppConfigurationPolicyResource_TargetedMobileAppsValidation tests GUID validation for targeted mobile apps
func TestAndroidManagedDeviceAppConfigurationPolicyResource_TargetedMobileAppsValidation(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, androidManagedDeviceAppConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer androidManagedDeviceAppConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test invalid GUID format
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "test" {
  display_name = "Test Android Managed Device App Configuration Policy"
  targeted_mobile_apps = ["invalid-guid", "another-invalid-guid"]
  package_id = "app:com.test"
  payload_json = jsonencode({})
}
`,
				ExpectError: regexp.MustCompile(`(Must be a valid GUID format|Invalid Attribute Value Match|Error running pre-apply plan)`),
			},
			// Test valid GUID format with real tenant app IDs
			{
				Config: `
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "test" {
  display_name = "Test Android Managed Device App Configuration Policy"
  targeted_mobile_apps = ["0e8ea6ec-63bb-436f-a89e-3adc475eb628", "9711516a-f6f8-4953-ad1f-45920ef34dda"]
  package_id = "app:com.test"
  payload_json = jsonencode({})
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.test", "targeted_mobile_apps.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.test", "targeted_mobile_apps.*", "0e8ea6ec-63bb-436f-a89e-3adc475eb628"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.test", "targeted_mobile_apps.*", "9711516a-f6f8-4953-ad1f-45920ef34dda"),
				),
			},
		},
	})
}
