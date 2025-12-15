package graphBetaMobileAppCatalogPackage_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	MobileAppCatalogPackageMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_and_app_management/graph_beta/mobile_app_catalog_package/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const dataSourceType = "data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package"

func setupMockEnvironment() (*mocks.Mocks, *MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	packagesMock := &MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock{}
	packagesMock.RegisterMocks()
	return mockClient, packagesMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	packagesMock := &MobileAppCatalogPackageMocks.MobileAppCatalogPackagesMock{}
	packagesMock.RegisterErrorMocks()
	return mockClient, packagesMock
}

// Test 01: Get all packages - comprehensive field validation
func TestMobileAppCatalogPackageDataSource_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".all").Key("filter_type").HasValue("all"),
					check.That(dataSourceType+".all").Key("items.#").HasValue("5"),

					// ============================================
					// Item 0: 7-Zip - Complete field validation
					// ============================================

					// Base mobile app fields
					check.That(dataSourceType+".all").Key("items.0.id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(dataSourceType+".all").Key("items.0.display_name").HasValue("7-Zip (x64)"),
					check.That(dataSourceType+".all").Key("items.0.description").Exists(),
					check.That(dataSourceType+".all").Key("items.0.publisher").HasValue("Igor Pavlov"),
					check.That(dataSourceType+".all").Key("items.0.created_date_time").HasValue("0001-01-01T00:00:00Z"),
					check.That(dataSourceType+".all").Key("items.0.last_modified_date_time").HasValue("0001-01-01T00:00:00Z"),
					check.That(dataSourceType+".all").Key("items.0.is_featured").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.privacy_information_url").HasValue("https://www.7-zip.org/"),
					check.That(dataSourceType+".all").Key("items.0.information_url").HasValue("https://www.7-zip.org"),
					check.That(dataSourceType+".all").Key("items.0.developer").HasValue("Igor Pavlov"),
					check.That(dataSourceType+".all").Key("items.0.upload_state").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.publishing_state").HasValue("notPublished"),
					check.That(dataSourceType+".all").Key("items.0.is_assigned").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.role_scope_tag_ids.#").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.dependent_app_count").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.superseding_app_count").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.superseded_app_count").HasValue("0"),

					// Win32 specific fields
					check.That(dataSourceType+".all").Key("items.0.file_name").HasValue("7z2501-x64.msi"),
					check.That(dataSourceType+".all").Key("items.0.size").HasValue("1996942"),
					check.That(dataSourceType+".all").Key("items.0.install_command_line").HasValue("\"%SystemRoot%\\System32\\msiexec.exe\" /i \"7z2501-x64.msi\" /qn REBOOT=ReallySuppress"),
					check.That(dataSourceType+".all").Key("items.0.uninstall_command_line").HasValue("\"%SystemRoot%\\System32\\msiexec.exe\" /X {23170F69-40C1-2702-2501-000001000000} /qn"),
					check.That(dataSourceType+".all").Key("items.0.applicable_architectures").HasValue("none"),
					check.That(dataSourceType+".all").Key("items.0.allowed_architectures").HasValue("x64"),
					check.That(dataSourceType+".all").Key("items.0.setup_file_path").HasValue("7z2501-x64.msi"),
					check.That(dataSourceType+".all").Key("items.0.minimum_supported_windows_release").HasValue("1607"),
					check.That(dataSourceType+".all").Key("items.0.display_version").HasValue("25.01"),
					check.That(dataSourceType+".all").Key("items.0.allow_available_uninstall").HasValue("true"),
					check.That(dataSourceType+".all").Key("items.0.mobile_app_catalog_package_id").HasValue("a09730b0-93d9-bb83-a96e-c5346258734b"),

					// Rules - File System Rule 1
					check.That(dataSourceType+".all").Key("items.0.rules.#").HasValue("3"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.odata_type").HasValue("#microsoft.graph.win32LobAppFileSystemRule"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.rule_type").HasValue("detection"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.path").HasValue("%ProgramFiles%\\7-Zip"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.file_or_folder_name").HasValue("7zFM.exe"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.check_32bit_on_64system").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.operation_type").HasValue("version"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.operator").HasValue("equal"),
					check.That(dataSourceType+".all").Key("items.0.rules.0.comparison_value").HasValue("25.1.0.0"),

					// Rules - File System Rule 2
					check.That(dataSourceType+".all").Key("items.0.rules.1.odata_type").HasValue("#microsoft.graph.win32LobAppFileSystemRule"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.rule_type").HasValue("detection"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.path").HasValue("%ProgramFiles%\\7-Zip"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.file_or_folder_name").HasValue("7zFM.exe"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.check_32bit_on_64system").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.operation_type").HasValue("sizeInBytes"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.operator").HasValue("equal"),
					check.That(dataSourceType+".all").Key("items.0.rules.1.comparison_value").HasValue("998912"),

					// Rules - Registry Rule
					check.That(dataSourceType+".all").Key("items.0.rules.2.odata_type").HasValue("#microsoft.graph.win32LobAppRegistryRule"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.rule_type").HasValue("detection"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.check_32bit_on_64system").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.key_path").HasValue("HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\{23170F69-40C1-2702-2501-000001000000}"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.value_name").HasValue("DisplayVersion"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.operation_type").HasValue("string"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.operator").HasValue("equal"),
					check.That(dataSourceType+".all").Key("items.0.rules.2.comparison_value").HasValue("25.01.00.0"),

					// Install Experience
					check.That(dataSourceType+".all").Key("items.0.install_experience.run_as_account").HasValue("system"),
					check.That(dataSourceType+".all").Key("items.0.install_experience.max_run_time_in_minutes").HasValue("60"),
					check.That(dataSourceType+".all").Key("items.0.install_experience.device_restart_behavior").HasValue("basedOnReturnCode"),

					// Return Codes
					check.That(dataSourceType+".all").Key("items.0.return_codes.#").HasValue("4"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.0.return_code").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.0.type").HasValue("success"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.1.return_code").HasValue("3010"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.1.type").HasValue("softReboot"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.2.return_code").HasValue("1618"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.2.type").HasValue("retry"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.3.return_code").HasValue("1707"),
					check.That(dataSourceType+".all").Key("items.0.return_codes.3.type").HasValue("success"),

					// MSI Information
					check.That(dataSourceType+".all").Key("items.0.msi_information.product_code").HasValue("{23170F69-40C1-2702-2501-000001000000}"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.product_version").HasValue("25.01"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.upgrade_code").HasValue("{23170F69-40C1-2702-0000-000004000000}"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.requires_reboot").HasValue("false"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.package_type").HasValue("perMachine"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.product_name").HasValue("7-Zip"),
					check.That(dataSourceType+".all").Key("items.0.msi_information.publisher").HasValue("Igor Pavlov"),

					// ============================================
					// Item 1: CPU-Z - Verify non-MSI app structure
					// ============================================
					check.That(dataSourceType+".all").Key("items.1.id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(dataSourceType+".all").Key("items.1.display_name").HasValue("CPU-Z (x64)"),
					check.That(dataSourceType+".all").Key("items.1.publisher").HasValue("CPUID, Inc"),
					check.That(dataSourceType+".all").Key("items.1.file_name").HasValue("cpu-z_2.17-en.exe"),
					check.That(dataSourceType+".all").Key("items.1.size").HasValue("4740244"),
					check.That(dataSourceType+".all").Key("items.1.display_version").HasValue("2.17"),
					check.That(dataSourceType+".all").Key("items.1.install_command_line").HasValue("\"cpu-z_2.17-en.exe\" /VERYSILENT /NORESTART"),
					check.That(dataSourceType+".all").Key("items.1.uninstall_command_line").HasValue("\"%ProgramW6432%\\CPUID\\CPU-Z\\unins000.exe\" /VERYSILENT"),
					check.That(dataSourceType+".all").Key("items.1.setup_file_path").HasValue("cpu-z_2.17-en.exe"),
					check.That(dataSourceType+".all").Key("items.1.mobile_app_catalog_package_id").HasValue("a8375b62-1909-812c-ee54-044ba1b1461b"),

					// CPU-Z Rules
					check.That(dataSourceType+".all").Key("items.1.rules.#").HasValue("3"),
					check.That(dataSourceType+".all").Key("items.1.rules.0.path").HasValue("%ProgramW6432%\\CPUID\\CPU-Z"),
					check.That(dataSourceType+".all").Key("items.1.rules.0.file_or_folder_name").HasValue("cpuz.exe"),
					check.That(dataSourceType+".all").Key("items.1.rules.0.operation_type").HasValue("sizeInBytes"),
					check.That(dataSourceType+".all").Key("items.1.rules.0.comparison_value").HasValue("7281384"),

					// CPU-Z Return Codes (only 2 codes)
					check.That(dataSourceType+".all").Key("items.1.return_codes.#").HasValue("2"),
					check.That(dataSourceType+".all").Key("items.1.return_codes.0.return_code").HasValue("0"),
					check.That(dataSourceType+".all").Key("items.1.return_codes.0.type").HasValue("success"),
					check.That(dataSourceType+".all").Key("items.1.return_codes.1.return_code").HasValue("3010"),
					check.That(dataSourceType+".all").Key("items.1.return_codes.1.type").HasValue("softReboot"),

					// ============================================
					// Item 2: Adobe AIR - Verify check32BitOn64System handling
					// ============================================
					check.That(dataSourceType+".all").Key("items.2.id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(dataSourceType+".all").Key("items.2.display_name").HasValue("Adobe AIR"),
					check.That(dataSourceType+".all").Key("items.2.publisher").HasValue("HARMAN International"),
					check.That(dataSourceType+".all").Key("items.2.file_name").HasValue("AdobeAIR.exe"),
					check.That(dataSourceType+".all").Key("items.2.size").HasValue("6288034"),
					check.That(dataSourceType+".all").Key("items.2.display_version").HasValue("51.2.2.6"),
					check.That(dataSourceType+".all").Key("items.2.mobile_app_catalog_package_id").HasValue("03ca0d35-9d3b-761e-db57-2116b6f6f2ea"),

					// Adobe AIR Rules - verify check32BitOn64System = true
					check.That(dataSourceType+".all").Key("items.2.rules.#").HasValue("3"),
					check.That(dataSourceType+".all").Key("items.2.rules.0.check_32bit_on_64system").HasValue("true"),
					check.That(dataSourceType+".all").Key("items.2.rules.2.check_32bit_on_64system").HasValue("true"),

					// ============================================
					// Item 3: Dell Display Manager
					// ============================================
					check.That(dataSourceType+".all").Key("items.3.id").HasValue("00000000-0000-0000-0000-000000000003"),
					check.That(dataSourceType+".all").Key("items.3.display_name").HasValue("Dell Display Manager"),
					check.That(dataSourceType+".all").Key("items.3.publisher").HasValue("Dell, Inc."),
					check.That(dataSourceType+".all").Key("items.3.file_name").HasValue("ddmsetup2110.exe"),
					check.That(dataSourceType+".all").Key("items.3.size").HasValue("1003058"),
					check.That(dataSourceType+".all").Key("items.3.display_version").HasValue("1.56.2110"),
					check.That(dataSourceType+".all").Key("items.3.mobile_app_catalog_package_id").HasValue("6274fe37-8968-1a80-e561-5a9fceff4579"),

					// ============================================
					// Item 4: Docker Desktop - Verify hardReboot return code
					// ============================================
					check.That(dataSourceType+".all").Key("items.4.id").HasValue("00000000-0000-0000-0000-000000000004"),
					check.That(dataSourceType+".all").Key("items.4.display_name").HasValue("Docker Desktop (x64)"),
					check.That(dataSourceType+".all").Key("items.4.publisher").HasValue("Docker Inc."),
					check.That(dataSourceType+".all").Key("items.4.file_name").HasValue("Docker Desktop Installer.exe"),
					check.That(dataSourceType+".all").Key("items.4.size").HasValue("569957978"),
					check.That(dataSourceType+".all").Key("items.4.display_version").HasValue("4.49.0.208700"),
					check.That(dataSourceType+".all").Key("items.4.mobile_app_catalog_package_id").HasValue("6e8cf1b6-1a04-d641-bc1c-04a8e61bff16"),

					// Docker Return Codes - includes hardReboot
					check.That(dataSourceType+".all").Key("items.4.return_codes.#").HasValue("3"),
					check.That(dataSourceType+".all").Key("items.4.return_codes.2.return_code").HasValue("1641"),
					check.That(dataSourceType+".all").Key("items.4.return_codes.2.type").HasValue("hardReboot"),

					// Docker Rules - verify appVersion operation type
					check.That(dataSourceType+".all").Key("items.4.rules.#").HasValue("1"),
					check.That(dataSourceType+".all").Key("items.4.rules.0.operation_type").HasValue("appVersion"),
					check.That(dataSourceType+".all").Key("items.4.rules.0.operator").HasValue("greaterThanOrEqual"),
					check.That(dataSourceType+".all").Key("items.4.rules.0.comparison_value").HasValue("4.49.0.0"),
				),
			},
		},
	})
}

// Test 02: Get by product ID - comprehensive single item validation
func TestMobileAppCatalogPackageDataSource_ById(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigById(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_id").Key("filter_type").HasValue("id"),
					check.That(dataSourceType+".by_id").Key("filter_value").HasValue("3a6307ef-6991-faf1-01e1-35e1557287aa"),
					check.That(dataSourceType+".by_id").Key("items.#").HasValue("1"),

					// Complete field validation for single item
					check.That(dataSourceType+".by_id").Key("items.0.id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(dataSourceType+".by_id").Key("items.0.display_name").HasValue("7-Zip (x64)"),
					check.That(dataSourceType+".by_id").Key("items.0.description").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.publisher").HasValue("Igor Pavlov"),
					check.That(dataSourceType+".by_id").Key("items.0.developer").HasValue("Igor Pavlov"),
					check.That(dataSourceType+".by_id").Key("items.0.privacy_information_url").HasValue("https://www.7-zip.org/"),
					check.That(dataSourceType+".by_id").Key("items.0.information_url").HasValue("https://www.7-zip.org"),
					check.That(dataSourceType+".by_id").Key("items.0.is_featured").HasValue("false"),
					check.That(dataSourceType+".by_id").Key("items.0.publishing_state").HasValue("notPublished"),
					check.That(dataSourceType+".by_id").Key("items.0.is_assigned").HasValue("false"),
					check.That(dataSourceType+".by_id").Key("items.0.file_name").HasValue("7z2501-x64.msi"),
					check.That(dataSourceType+".by_id").Key("items.0.size").HasValue("1996942"),
					check.That(dataSourceType+".by_id").Key("items.0.display_version").HasValue("25.01"),
					check.That(dataSourceType+".by_id").Key("items.0.allowed_architectures").HasValue("x64"),
					check.That(dataSourceType+".by_id").Key("items.0.minimum_supported_windows_release").HasValue("1607"),
					check.That(dataSourceType+".by_id").Key("items.0.allow_available_uninstall").HasValue("true"),
					check.That(dataSourceType+".by_id").Key("items.0.mobile_app_catalog_package_id").HasValue("a09730b0-93d9-bb83-a96e-c5346258734b"),

					// Install/Uninstall commands
					check.That(dataSourceType+".by_id").Key("items.0.install_command_line").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.uninstall_command_line").Exists(),
					check.That(dataSourceType+".by_id").Key("items.0.setup_file_path").HasValue("7z2501-x64.msi"),

					// Rules validation
					check.That(dataSourceType+".by_id").Key("items.0.rules.#").HasValue("3"),
					check.That(dataSourceType+".by_id").Key("items.0.rules.0.rule_type").HasValue("detection"),
					check.That(dataSourceType+".by_id").Key("items.0.rules.0.odata_type").Exists(),

					// Install experience validation
					check.That(dataSourceType+".by_id").Key("items.0.install_experience.run_as_account").HasValue("system"),
					check.That(dataSourceType+".by_id").Key("items.0.install_experience.max_run_time_in_minutes").HasValue("60"),
					check.That(dataSourceType+".by_id").Key("items.0.install_experience.device_restart_behavior").HasValue("basedOnReturnCode"),

					// Return codes validation
					check.That(dataSourceType+".by_id").Key("items.0.return_codes.#").HasValue("4"),

					// MSI information validation
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.product_code").HasValue("{23170F69-40C1-2702-2501-000001000000}"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.product_version").HasValue("25.01"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.upgrade_code").HasValue("{23170F69-40C1-2702-0000-000004000000}"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.requires_reboot").HasValue("false"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.package_type").HasValue("perMachine"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.product_name").HasValue("7-Zip"),
					check.That(dataSourceType+".by_id").Key("items.0.msi_information.publisher").HasValue("Igor Pavlov"),
				),
			},
		},
	})
}

// Test 03: Get by product name
func TestMobileAppCatalogPackageDataSource_ByProductName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByProductName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_product_name").Key("filter_type").HasValue("product_name"),
					check.That(dataSourceType+".by_product_name").Key("filter_value").HasValue("7-Zip"),
					check.That(dataSourceType+".by_product_name").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`7-Zip`)),
					check.That(dataSourceType+".by_product_name").Key("items.0.publisher").HasValue("Igor Pavlov"),
					check.That(dataSourceType+".by_product_name").Key("items.0.file_name").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.mobile_app_catalog_package_id").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.install_command_line").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.uninstall_command_line").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.rules.#").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.install_experience.run_as_account").Exists(),
					check.That(dataSourceType+".by_product_name").Key("items.0.return_codes.#").Exists(),
				),
			},
		},
	})
}

// Test 04: Get by publisher name
func TestMobileAppCatalogPackageDataSource_ByPublisherName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, packagesMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer packagesMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigByPublisherName(),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".by_publisher").Key("filter_type").HasValue("publisher_name"),
					check.That(dataSourceType+".by_publisher").Key("filter_value").HasValue("Docker"),
					check.That(dataSourceType+".by_publisher").Key("items.0.display_name").MatchesRegex(regexp.MustCompile(`Docker`)),
					check.That(dataSourceType+".by_publisher").Key("items.0.publisher").HasValue("Docker Inc."),
					check.That(dataSourceType+".by_publisher").Key("items.0.file_name").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.mobile_app_catalog_package_id").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.install_command_line").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.rules.#").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.install_experience.run_as_account").Exists(),
					check.That(dataSourceType+".by_publisher").Key("items.0.return_codes.#").Exists(),
				),
			},
		},
	})
}

// Terraform config helpers
func testConfigAll() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/01_all.tf")
	if err != nil {
		panic("failed to load 01_all config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigById() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/02_by_id.tf")
	if err != nil {
		panic("failed to load 02_by_id config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByProductName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/03_by_product_name.tf")
	if err != nil {
		panic("failed to load 03_by_product_name config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigByPublisherName() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/04_by_publisher_name.tf")
	if err != nil {
		panic("failed to load 04_by_publisher_name config: " + err.Error())
	}
	return unitTestConfig
}
