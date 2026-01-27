package graphBetaDeviceComplianceNotificationTemplate_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	complianceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_compliance_notification_template/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *complianceMocks.WindowsDeviceComplianceNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceMock := &complianceMocks.WindowsDeviceComplianceNotificationsMock{}
	complianceMock.RegisterMocks()
	return mockClient, complianceMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *complianceMocks.WindowsDeviceComplianceNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	complianceMock := &complianceMocks.WindowsDeviceComplianceNotificationsMock{}
	complianceMock.RegisterErrorMocks()
	return mockClient, complianceMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestUnitResourceDeviceComplianceNotificationTemplate_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "display_name", "English Compliance Notification"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.*", "1"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_02_BrandingOptions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.english"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "branding_options.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "branding_options.*", "includeCompanyLogo"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "branding_options.*", "includeCompanyName"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "branding_options.*", "includeContactInformation"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_03_LocalizedMessages(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.english"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "localized_notification_messages.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "localized_notification_messages.*", map[string]string{
						"locale":     "en-us",
						"subject":    "Immediate Action Required: Device Compliance",
						"is_default": "true",
					}),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_04_MultilingualMessages(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "display_name", "Multilingual Compliance Notification"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "localized_notification_messages.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "localized_notification_messages.*", map[string]string{
						"locale":     "en-us",
						"subject":    "Device Compliance Issue Detected",
						"is_default": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "localized_notification_messages.*", map[string]string{
						"locale":     "es-es",
						"subject":    "Problema de Cumplimiento del Dispositivo",
						"is_default": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "localized_notification_messages.*", map[string]string{
						"locale":     "fr-fr",
						"subject":    "Problème de Conformité de l'Appareil",
						"is_default": "false",
					}),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_05_AllTerraformConfigurations(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.english"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "display_name", "English Compliance Notification"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "branding_options.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "localized_notification_messages.#", "1"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "display_name", "Multilingual Compliance Notification"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "branding_options.#", "3"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "localized_notification_messages.#", "3"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_06_RoleScopeTags(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.english"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.english", "role_scope_tag_ids.*", "1"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.multilingual", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceComplianceNotificationTemplate_07_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, complianceMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer complianceMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigMinimal(),
				ExpectError: regexp.MustCompile("Invalid Windows Device Compliance Notification data"),
			},
		},
	})
}

// Configuration Functions
func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load maximal config: " + err.Error())
	}
	return unitTestConfig
}
