package graphBetaMacosDeviceConfigurationTemplates_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	macosConfigMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *macosConfigMocks.MacosDeviceConfigurationTemplatesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &macosConfigMocks.MacosDeviceConfigurationTemplatesMock{}
	configMock.RegisterMocks()
	return mockClient, configMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *macosConfigMocks.MacosDeviceConfigurationTemplatesMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	configMock := &macosConfigMocks.MacosDeviceConfigurationTemplatesMock{}
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

func TestMacosDeviceConfigurationTemplatesResource_CustomConfiguration(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-macOS-custom-configuration-example"),
					check.That(resourceType+".custom_configuration_example").Key("description").HasValue("Example custom configuration template for macOS devices"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload_file_name").HasValue("com.example.custom.mobileconfig"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload_name").HasValue("Custom Configuration Example"),
					check.That(resourceType+".custom_configuration_example").Key("custom_configuration.payload").Exists(),
					check.That(resourceType+".custom_configuration_example").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".custom_configuration_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_PreferenceFile(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_preference_file_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".preference_file_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".preference_file_example").Key("display_name").HasValue("unit-test-macOS-preference-file-example"),
					check.That(resourceType+".preference_file_example").Key("description").HasValue("Configure Safari browser settings via preference file"),
					check.That(resourceType+".preference_file_example").Key("preference_file.file_name").HasValue("com.apple.Safari.plist"),
					check.That(resourceType+".preference_file_example").Key("preference_file.bundle_id").HasValue("com.apple.Safari"),
					check.That(resourceType+".preference_file_example").Key("preference_file.configuration_xml").Exists(),
					check.That(resourceType+".preference_file_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_TrustedRootCertificate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_trusted_root_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".trusted_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".trusted_cert_example").Key("display_name").HasValue("unit-test-macOS-trusted-root-certificate-example"),
					check.That(resourceType+".trusted_cert_example").Key("description").HasValue("Install company root certificate for secure connections"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.cert_file_name").HasValue("MicrosoftRootCertificateAuthority2011.cer"),
					check.That(resourceType+".trusted_cert_example").Key("trusted_certificate.trusted_root_certificate").Exists(),
					check.That(resourceType+".trusted_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_ScepCertificate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_scep_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".scep_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".scep_cert_example").Key("display_name").HasValue("unit-test-macOS-scep-certificate-example"),
					check.That(resourceType+".scep_cert_example").Key("description").HasValue("SCEP certificate profile for device authentication"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.deployment_channel").HasValue("deviceChannel"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.renewal_threshold_percentage").HasValue("20"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_store").HasValue("machine"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_validity_period_scale").HasValue("years"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.certificate_validity_period_value").HasValue("2"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.subject_name_format").HasValue("custom"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.subject_name_format_string").HasValue("CN={{DeviceName}},O=Example Corp,C=US"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.root_certificate_odata_bind").HasValue("https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('87654321-4321-4321-4321-210987654321')"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_size").HasValue("size2048"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_usage.#").HasValue("2"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_usage.*").ContainsTypeSetElement("digitalSignature"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.key_usage.*").ContainsTypeSetElement("keyEncipherment"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.custom_subject_alternative_names.#").HasValue("4"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.extended_key_usages.#").HasValue("4"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.scep_server_urls.#").HasValue("2"),
					check.That(resourceType+".scep_cert_example").Key("scep_certificate.allow_all_apps_access").HasValue("false"),
					check.That(resourceType+".scep_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_PkcsCertificate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_pkcs_certificate_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".pkcs_cert_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".pkcs_cert_example").Key("display_name").HasValue("unit-test-macOS-pkcs-certificate-example"),
					check.That(resourceType+".pkcs_cert_example").Key("description").HasValue("PKCS certificate profile for user authentication"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.deployment_channel").HasValue("userChannel"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.renewal_threshold_percentage").HasValue("30"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_store").HasValue("user"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_validity_period_scale").HasValue("months"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_validity_period_value").HasValue("12"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.subject_name_format").HasValue("commonNameIncludingEmail"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.subject_name_format_string").HasValue("CN={{UserName}},E={{EmailAddress}},O=Example Corp"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certification_authority").HasValue("ExampleCA\\ExampleCA-CA"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certification_authority_name").HasValue("ExampleCA-CA"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.certificate_template_name").HasValue("UserAuthentication"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.custom_subject_alternative_names.#").HasValue("6"),
					check.That(resourceType+".pkcs_cert_example").Key("pkcs_certificate.allow_all_apps_access").HasValue("true"),
					check.That(resourceType+".pkcs_cert_example").Key("assignments.#").HasValue("4"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_CreateWithError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_custom_configuration_maximal.tf"),
				ExpectError: regexp.MustCompile("Error creating macOS device configuration template"),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-macOS-custom-configuration-example"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").HasValue("unit-test-macOS-custom-configuration-example"),
				),
			},
		},
	})
}

func TestMacosDeviceConfigurationTemplatesResource_ImportState(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, configMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer configMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_custom_configuration_maximal.tf"),
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
