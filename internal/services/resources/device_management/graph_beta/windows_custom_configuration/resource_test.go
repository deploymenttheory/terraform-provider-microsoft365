package graphBetaWindowsCustomConfiguration_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	windowsCustomConfigMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_custom_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *windowsCustomConfigMocks.WindowsCustomConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &windowsCustomConfigMocks.WindowsCustomConfigurationMock{}
	configMock.RegisterMocks()
	return mockClient, configMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *windowsCustomConfigMocks.WindowsCustomConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &windowsCustomConfigMocks.WindowsCustomConfigurationMock{}
	configMock.RegisterErrorMocks()
	return mockClient, configMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceWindowsCustomConfiguration_01_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_windows_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-windows-custom-configuration-example"),
					check.That(resourceType+".custom_configuration_example").Key("description").HasValue("Example Windows custom configuration profile using OMA-URI settings"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.#").HasValue("4"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.odata_type").HasValue("#microsoft.graph.omaSettingString"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.display_name").HasValue("ADMX Ingest"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.oma_uri").HasValue("./Device/Vendor/MSFT/Policy/ConfigOperations/ADMXInstall/VSCode/Policy/VSCodeADMX"),
					// Encrypted values are masked as "****" by the API; the provider must resolve
					// them back to plain text via getOmaSettingPlainTextValue.
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.value").HasValue("<?xml version=\"1.0\" encoding=\"utf-8\"?><policyDefinitions revision=\"1.1\" schemaVersion=\"1.0\"><policies /></policyDefinitions>"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.1.value").HasValue("<enabled/>\n<data id=\"ExtensionsAutoUpdate\" value=\"off\"/>"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.2.odata_type").HasValue("#microsoft.graph.omaSettingInteger"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.2.value").HasValue("30"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.3.odata_type").HasValue("#microsoft.graph.omaSettingBoolean"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.3.value").HasValue("false"),
					check.That(resourceType+".custom_configuration_example").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".custom_configuration_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_02_InvalidOmaSettingValue(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_windows_custom_configuration_invalid_integer.tf"),
				ExpectError: regexp.MustCompile("is not a valid integer"),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_02a_NonCanonicalValue(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_windows_custom_configuration_non_canonical_value.tf"),
				ExpectError: regexp.MustCompile("is not in the canonical form"),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_02b_DuplicateOmaUri(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_windows_custom_configuration_duplicate_oma_uri.tf"),
				ExpectError: regexp.MustCompile("Duplicate OMA-URI"),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_03_CreateWithError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_windows_custom_configuration_maximal.tf"),
				ExpectError: regexp.MustCompile("Error creating windows custom configuration|BadRequest"),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_04_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_windows_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-windows-custom-configuration-example"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_windows_custom_configuration_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-windows-custom-configuration-example-updated"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.#").HasValue("1"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.value").HasValue("<enabled/>\n<data id=\"ExtensionsAutoUpdate\" value=\"on\"/>"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsCustomConfiguration_05_ImportState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_windows_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
			{
				ResourceName:      resourceType + ".custom_configuration_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
