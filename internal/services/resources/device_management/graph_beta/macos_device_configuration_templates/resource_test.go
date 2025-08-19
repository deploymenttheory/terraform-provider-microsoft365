package graphBetaMacosDeviceConfigurationTemplates_test

import (
	"regexp"
	"testing"

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

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
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
				Config: testConfigCustomConfigurationMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "display_name", "unit-test-macOS-custom-configuration-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "description", "Example custom configuration template for macOS devices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload_file_name", "com.example.custom.mobileconfig"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload_name", "Custom Configuration Example"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "role_scope_tag_ids.0"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "assignments.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "assignments.*", map[string]string{
						"type": "exclusionGroupAssignmentTarget",
					}),
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
				Config: testConfigPreferenceFileMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "display_name", "unit-test-macOS-preference-file-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "description", "Configure Safari browser settings via preference file"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.file_name", "com.apple.Safari.plist"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.bundle_id", "com.apple.Safari"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.configuration_xml"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "assignments.#", "4"),
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
				Config: testConfigTrustedRootCertificateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "display_name", "unit-test-macOS-trusted-root-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "description", "Install company root certificate for secure connections"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.cert_file_name", "CompanyRootCA.cer"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.trusted_root_certificate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "assignments.#", "4"),
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
				Config: testConfigScepCertificateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "display_name", "unit-test-macOS-scep-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "description", "SCEP certificate profile for device authentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.renewal_threshold_percentage", "20"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_store", "machine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_validity_period_scale", "years"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_validity_period_value", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.subject_name_format", "custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.subject_name_format_string", "CN={{DeviceName}},O=Example Corp,C=US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.root_certificate_odata_bind", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('87654321-4321-4321-4321-210987654321')"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_size", "size2048"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_usage.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_usage.*", "digitalSignature"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_usage.*", "keyEncipherment"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.custom_subject_alternative_names.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.extended_key_usages.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.scep_server_urls.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.allow_all_apps_access", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "assignments.#", "4"),
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
				Config: testConfigPkcsCertificateMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "display_name", "unit-test-macOS-pkcs-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "description", "PKCS certificate profile for user authentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.deployment_channel", "userChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.renewal_threshold_percentage", "30"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_store", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_validity_period_scale", "months"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_validity_period_value", "12"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.subject_name_format", "commonNameIncludingEmail"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.subject_name_format_string", "CN={{UserName}},E={{EmailAddress}},O=Example Corp"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certification_authority", "ExampleCA\\ExampleCA-CA"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certification_authority_name", "ExampleCA-CA"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_template_name", "UserAuthentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.custom_subject_alternative_names.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.allow_all_apps_access", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "assignments.#", "4"),
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
				Config:      testConfigCustomConfigurationMaximal(),
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
				Config: testConfigCustomConfigurationMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "display_name", "unit-test-macOS-custom-configuration-example"),
				),
			},
			{
				Config: testConfigCustomConfigurationMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "display_name", "unit-test-macOS-custom-configuration-example"),
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
				Config: testConfigCustomConfigurationMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// macOS Device Configuration Templates Configuration Functions
func testConfigCustomConfigurationMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_custom_configuration_maximal.tf")
	if err != nil {
		panic("failed to load custom configuration maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigPreferenceFileMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_preference_file_maximal.tf")
	if err != nil {
		panic("failed to load preference file maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigTrustedRootCertificateMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_trusted_root_certificate_maximal.tf")
	if err != nil {
		panic("failed to load trusted root certificate maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigScepCertificateMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_scep_certificate_maximal.tf")
	if err != nil {
		panic("failed to load scep certificate maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigPkcsCertificateMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_pkcs_certificate_maximal.tf")
	if err != nil {
		panic("failed to load pkcs certificate maximal config: " + err.Error())
	}
	return unitTestConfig
}
