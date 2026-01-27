package graphBetaDeviceComplianceNotificationTemplate_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceDeviceComplianceNotificationTemplate_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceComplianceNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDeviceComplianceNotificationsConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "branding_options.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "branding_options.*", "includeCompanyLogo"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "localized_notification_messages.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "localized_notification_messages.*", map[string]string{
						"locale":     "en-us",
						"subject":    "Device Compliance Required",
						"is_default": "true",
					}),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccResourceDeviceComplianceNotificationTemplate_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceComplianceNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDeviceComplianceNotificationsConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "branding_options.#", "5"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "branding_options.*", "includeCompanyLogo"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "branding_options.*", "includeCompanyName"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "branding_options.*", "includeContactInformation"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "localized_notification_messages.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "localized_notification_messages.*", map[string]string{
						"locale":     "en-us",
						"subject":    "Device Compliance Issue Detected",
						"is_default": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "localized_notification_messages.*", map[string]string{
						"locale":     "es-es",
						"subject":    "Problema de Cumplimiento del Dispositivo",
						"is_default": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", "localized_notification_messages.*", map[string]string{
						"locale":     "fr-fr",
						"subject":    "Problème de Conformité de l'Appareil",
						"is_default": "false",
					}),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_compliance_notification_template.maximal", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccResourceDeviceComplianceNotificationTemplate_05_BrandingOptions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceComplianceNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDeviceComplianceNotificationsConfig_brandingTest(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "branding_options.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "branding_options.*", "includeCompanyLogo"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "branding_options.*", "includeCompanyName"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "branding_options.*", "includeContactInformation"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccResourceDeviceComplianceNotificationTemplate_04_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsDeviceComplianceNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsDeviceComplianceNotificationsConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.minimal", "branding_options.#", "1"),
				),
			},
			{
				Config: testAccWindowsDeviceComplianceNotificationsConfig_brandingTest(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", "branding_options.#", "3"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_compliance_notification_template.branding_test", ImportState: true, ImportStateVerify: true},
		},
	})
}

// Configuration Functions
func testAccWindowsDeviceComplianceNotificationsConfig_minimal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal acceptance test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsDeviceComplianceNotificationsConfig_maximal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal acceptance test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsDeviceComplianceNotificationsConfig_brandingTest() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_branding_test.tf")
	if err != nil {
		log.Fatalf("Failed to load branding test acceptance config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccCheckWindowsDeviceComplianceNotificationsDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_device_compliance_notification_template" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			NotificationMessageTemplates().
			ByNotificationMessageTemplateId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows device compliance notification %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows device compliance notification %s still exists", rs.Primary.ID)
	}
	return nil
}
