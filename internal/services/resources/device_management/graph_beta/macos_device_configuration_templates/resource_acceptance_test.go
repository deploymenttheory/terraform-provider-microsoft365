package graphBetaMacosDeviceConfigurationTemplates_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Custom Configuration Tests
func TestAccMacosDeviceConfigurationTemplatesResource_CustomConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacosDeviceConfigurationTemplatesDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacosDeviceConfigurationTemplatesConfig_customConfiguration(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "display_name", "acc-test-macOS-custom-configuration-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "description", "Example custom configuration template for macOS devices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload_file_name", "com.example.custom.mobileconfig"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload_name", "Custom Configuration Example"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "custom_configuration.payload"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", "assignments.#", "4"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_configuration_example", ImportState: true, ImportStateVerify: true},
		},
	})
}

// Preference File Tests
func TestAccMacosDeviceConfigurationTemplatesResource_PreferenceFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacosDeviceConfigurationTemplatesDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacosDeviceConfigurationTemplatesConfig_preferenceFile(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "display_name", "acc-test-macOS-preference-file-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "description", "Configure Safari browser settings via preference file"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.file_name", "com.apple.Safari.plist"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.bundle_id", "com.apple.Safari"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "preference_file.configuration_xml"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example", "assignments.#", "4"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Trusted Root Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_TrustedRootCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacosDeviceConfigurationTemplatesDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacosDeviceConfigurationTemplatesConfig_trustedRootCertificate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "display_name", "acc-test-macOS-trusted-root-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "description", "Install company root certificate for secure connections"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.cert_file_name", "MicrosoftRootCertificateAuthority2011.cer"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "trusted_certificate.trusted_root_certificate"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example", "assignments.#", "4"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// SCEP Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_ScepCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacosDeviceConfigurationTemplatesDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacosDeviceConfigurationTemplatesConfig_scepCertificate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "display_name", "acc-test-macOS-scep-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "description", "SCEP certificate profile for device authentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.renewal_threshold_percentage", "20"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_store", "machine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_validity_period_scale", "years"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.certificate_validity_period_value", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.subject_name_format", "custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_size", "size4096"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.key_usage.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.custom_subject_alternative_names.#", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.extended_key_usages.#", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.scep_server_urls.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "scep_certificate.allow_all_apps_access", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example", "assignments.#", "4"),
				),
			},
			{
				ResourceName:            "microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scep_certificate.root_certificate_odata_bind"},
			},
		},
	})
}

// PKCS Certificate Tests
func TestAccMacosDeviceConfigurationTemplatesResource_PkcsCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacosDeviceConfigurationTemplatesDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacosDeviceConfigurationTemplatesConfig_pkcsCertificate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "display_name", "acc-test-macOS-pkcs-certificate-example"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "description", "PKCS certificate profile for user authentication"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.deployment_channel", "deviceChannel"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.renewal_threshold_percentage", "20"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_store", "machine"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_validity_period_scale", "months"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_validity_period_value", "12"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.subject_name_format", "custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.subject_name_format_string", "CN={{AAD_Device_ID}}"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certification_authority", "some-auth"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certification_authority_name", "some-name"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.certificate_template_name", "some-template-name"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.custom_subject_alternative_names.#", "6"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "pkcs_certificate.allow_all_apps_access", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", "assignments.#", "4"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example", ImportState: true, ImportStateVerify: true},
		},
	})
}

// Configuration Functions
func testAccMacosDeviceConfigurationTemplatesConfig_customConfiguration() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_custom_configuration_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load custom configuration test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccMacosDeviceConfigurationTemplatesConfig_preferenceFile() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_preference_file_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load preference file test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccMacosDeviceConfigurationTemplatesConfig_trustedRootCertificate() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_trusted_root_certificate_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load trusted root certificate test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccMacosDeviceConfigurationTemplatesConfig_scepCertificate() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_scep_certificate_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load SCEP certificate test config: %v", err)
	}
	trustedCert, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_trusted_root_certificate_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load trusted root certificate test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + trustedCert + "\n" + accTestConfig)
}

func testAccMacosDeviceConfigurationTemplatesConfig_pkcsCertificate() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	assignmentFilters, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/assignment_filter.tf")
	if err != nil {
		log.Fatalf("Failed to load assignment filters config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_pkcs_certificate_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load PKCS certificate test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + assignmentFilters + "\n" + accTestConfig)
}

func testAccCheckMacosDeviceConfigurationTemplatesDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_device_configuration_templates" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if macos device configuration template %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("macos device configuration template %s still exists", rs.Primary.ID)
	}
	return nil
}
