package graphBetaSettingsCatalogConfigurationPolicyJson_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaSettingsCatalogConfigurationPolicyJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json"
	settingsCatalogPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Setup functions for mocks
func setupMockEnvironment() (*mocks.Mocks, *settingsCatalogPolicyMocks.SettingsCatalogPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &settingsCatalogPolicyMocks.SettingsCatalogPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *settingsCatalogPolicyMocks.SettingsCatalogPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &settingsCatalogPolicyMocks.SettingsCatalogPolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

// loadUnitTestTerraform loads terraform test files from tests/terraform/unit
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Camera - Simple choice setting
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_01_Camera(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_camera.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".camera").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".camera").Key("name").HasValue("Test Camera Policy"),
					check.That(resourceType+".camera").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".camera",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test 02: Task Manager - Simple choice setting
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_02_TaskManager(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_task_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".task_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".task_manager").Key("name").HasValue("Test Task Manager Policy"),
					check.That(resourceType+".task_manager").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".task_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 03: App Privacy - Simple choice setting
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_03_AppPrivacy(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_03_app_privacy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".app_privacy").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".app_privacy").Key("name").HasValue("Test App Privacy Policy"),
					check.That(resourceType+".app_privacy").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".app_privacy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 04: Cryptography - Simple choice setting
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_04_Cryptography(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_04_cryptography.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cryptography").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cryptography").Key("name").HasValue("Test Cryptography Policy"),
					check.That(resourceType+".cryptography").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".cryptography",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 05: Notifications - Simple choice setting
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_05_Notifications(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_05_notifications.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".notifications").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".notifications").Key("name").HasValue("Test Notifications Policy"),
					check.That(resourceType+".notifications").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".notifications",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 06: Attachment Manager - Multiple choice settings
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_06_AttachmentManager(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_06_attachment_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".attachment_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".attachment_manager").Key("name").HasValue("Test Attachment Manager Policy"),
					check.That(resourceType+".attachment_manager").Key("settings_count").HasValue("2"),
				),
			},
			{
				ResourceName:            resourceType + ".attachment_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 07: Credential User Interface - Multiple choice settings
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_07_CredentialUserInterface(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_07_credential_user_interface.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".credential_user_interface").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".credential_user_interface").Key("name").HasValue("Test Credential User Interface Policy"),
					check.That(resourceType+".credential_user_interface").Key("settings_count").HasValue("2"),
				),
			},
			{
				ResourceName:            resourceType + ".credential_user_interface",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 08: Remote Desktop AVD URL - Simple collection settings
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_08_RemoteDesktopAVDURL(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_08_remote_desktop_avd_url.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".remote_desktop_avd_url").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".remote_desktop_avd_url").Key("name").HasValue("Test Remote Desktop AVD URL Policy"),
					check.That(resourceType+".remote_desktop_avd_url").Key("settings_count").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".remote_desktop_avd_url",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 09: Storage Sense - Integer settings and multiple choices
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_09_StorageSense(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_09_storage_sense.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".storage_sense").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".storage_sense").Key("name").HasValue("Test Storage Sense Policy"),
					check.That(resourceType+".storage_sense").Key("settings_count").HasValue("8"),
				),
			},
			{
				ResourceName:            resourceType + ".storage_sense",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 10: Windows Connection Manager - Choice with children
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_10_WindowsConnectionManager(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_10_windows_connection_manager.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".windows_connection_manager").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".windows_connection_manager").Key("name").HasValue("Test Windows Connection Manager Policy"),
					check.That(resourceType+".windows_connection_manager").Key("settings_count").HasValue("4"),
				),
			},
			{
				ResourceName:            resourceType + ".windows_connection_manager",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 11: AutoPlay Policies - Nested choice settings
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_11_AutoPlayPolicies(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_11_autoplay_policies.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".autoplay").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".autoplay").Key("name").HasValue("Test AutoPlay Policies Policy"),
					check.That(resourceType+".autoplay").Key("settings_count").HasValue("3"),
				),
			},
			{
				ResourceName:            resourceType + ".autoplay",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 12: Defender Smartscreen - Choice with collection child
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_12_DefenderSmartscreen(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_12_defender_smartscreen.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".defender_smartscreen").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".defender_smartscreen").Key("name").HasValue("Test Defender Smartscreen Policy"),
					check.That(resourceType+".defender_smartscreen").Key("settings_count").HasValue("11"),
				),
			},
			{
				ResourceName:            resourceType + ".defender_smartscreen",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 13: Edge Extensions macOS - Multiple collections
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_13_EdgeExtensionsMacOS(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_13_edge_extensions_macos.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".edge_extensions_macos").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".edge_extensions_macos").Key("name").HasValue("Test Edge Extensions macOS Policy"),
					check.That(resourceType+".edge_extensions_macos").Key("platforms").HasValue("macOS"),
					check.That(resourceType+".edge_extensions_macos").Key("settings_count").HasValue("4"),
				),
			},
			{
				ResourceName:            resourceType + ".edge_extensions_macos",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 14: Office Configuration macOS - Nested group collections
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_14_OfficeConfigurationMacOS(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_14_office_configuration_macos.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".office_configuration_macos").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".office_configuration_macos").Key("name").HasValue("Test Office Configuration macOS Policy"),
					check.That(resourceType+".office_configuration_macos").Key("platforms").HasValue("macOS"),
					check.That(resourceType+".office_configuration_macos").Key("settings_count").HasValue("3"),
				),
			},
			{
				ResourceName:            resourceType + ".office_configuration_macos",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 15: Defender Antivirus Baseline - Complex nested group collections
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_15_DefenderAntivirusBaseline(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_15_defender_antivirus_baseline.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".defender_antivirus_baseline").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".defender_antivirus_baseline").Key("name").HasValue("Test Defender Antivirus Security Baseline Policy"),
					check.That(resourceType+".defender_antivirus_baseline").Key("settings_count").HasValue("9"),
				),
			},
			{
				ResourceName:            resourceType + ".defender_antivirus_baseline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 16: File Explorer with minimal assignments
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_16_FileExplorerMinimalAssignments(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_16_file_explorer_minimal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".file_explorer_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".file_explorer_minimal").Key("name").HasValue("Test File Explorer Policy"),
					check.That(resourceType+".file_explorer_minimal").Key("settings_count").HasValue("2"),
					check.That(resourceType+".file_explorer_minimal").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".file_explorer_minimal").Key("is_assigned").HasValue("true"),
				),
			},
			{
				ResourceName:            resourceType + ".file_explorer_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 17: Local Policies Security Options with maximal assignments
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_17_LocalPoliciesMaximalAssignments(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_17_local_policies_maximal_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".local_policies_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".local_policies_maximal").Key("name").HasValue("Test Local Policies Security Options Policy"),
					check.That(resourceType+".local_policies_maximal").Key("settings_count").HasValue("6"),
					check.That(resourceType+".local_policies_maximal").Key("assignments.#").HasValue("4"),
					check.That(resourceType+".local_policies_maximal").Key("is_assigned").HasValue("true"),
				),
			},
			{
				ResourceName:            resourceType + ".local_policies_maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings"},
			},
		},
	})
}

// Test 18: Update from minimal to maximal
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_18_UpdateMinimalToMaximal(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_18_update_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("Test Update Policy"),
					check.That(resourceType+".update_test").Key("settings_count").HasValue("2"),
					check.That(resourceType+".update_test").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_18_update_minimal_to_maximal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("Test Update Policy - Updated"),
					check.That(resourceType+".update_test").Key("settings_count").HasValue("6"),
					check.That(resourceType+".update_test").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

// Test 19: Update from maximal to minimal
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_19_UpdateMaximalToMinimal(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_19_update_maximal_to_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_reverse_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_reverse_test").Key("name").HasValue("Test Update Reverse Policy"),
					check.That(resourceType+".update_reverse_test").Key("settings_count").HasValue("6"),
					check.That(resourceType+".update_reverse_test").Key("assignments.#").HasValue("4"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_19_update_maximal_to_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update_reverse_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_reverse_test").Key("name").HasValue("Test Update Reverse Policy - Updated"),
					check.That(resourceType+".update_reverse_test").Key("settings_count").HasValue("2"),
					check.That(resourceType+".update_reverse_test").Key("assignments.#").HasValue("1"),
				),
			},
		},
	})
}
