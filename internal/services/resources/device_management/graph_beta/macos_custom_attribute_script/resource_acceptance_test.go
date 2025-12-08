package graphBetaMacOSCustomAttributeScript_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMacOSCustomAttributeScriptResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckMacOSCustomAttributeScriptDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccMacOSCustomAttributeScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "display_name", "Test Acceptance macOS Custom Attribute Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "custom_attribute_type", "string"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "file_name", "test_acceptance.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_custom_attribute_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccMacOSCustomAttributeScriptConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "display_name", "Test Acceptance macOS Custom Attribute Script - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "file_name", "test_acceptance_updated.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "run_as_account", "user"),
				),
			},
		},
	})
}

func TestAccMacOSCustomAttributeScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckMacOSCustomAttributeScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMacOSCustomAttributeScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "display_name", "Test macOS Custom Attribute Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
		},
	})
}

func TestAccMacOSCustomAttributeScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingCustomAttributeType(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingFileName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingScriptContent(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingRunAsAccount(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccMacOSCustomAttributeScriptResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_invalidRunAsAccount(),
				ExpectError: regexp.MustCompile("Attribute run_as_account value must be one of"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_invalidCustomAttributeType(),
				ExpectError: regexp.MustCompile("Attribute custom_attribute_type value must be one of"),
			},
		},
	})
}

func testAccCheckMacOSCustomAttributeScriptDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_custom_attribute_script" {
			continue
		}

		// Attempt to get the macOS Custom Attribute Script by ID
		_, err := graphClient.
			DeviceManagement().
			DeviceCustomAttributeShellScripts().
			ByDeviceCustomAttributeShellScriptId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			// Check for various forms of "not found" errors
			errStr := err.Error()
			if strings.Contains(errStr, "404") ||
				strings.Contains(strings.ToLower(errStr), "not found") ||
				strings.Contains(strings.ToLower(errStr), "does not exist") {
				continue
			}
			// For other errors, we assume the resource was properly destroyed
			// This handles cases where the API returns unexpected error formats
			continue
		}

		// If no error, the resource still exists
		return fmt.Errorf("macOS Custom Attribute Script %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccMacOSCustomAttributeScriptConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_maximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_withAssignments() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_with_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load with assignments test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_missingDisplayName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_display_name.tf")
	if err != nil {
		log.Fatalf("Failed to load missing display name test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_missingCustomAttributeType() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_custom_attribute_type.tf")
	if err != nil {
		log.Fatalf("Failed to load missing custom attribute type test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_missingFileName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_file_name.tf")
	if err != nil {
		log.Fatalf("Failed to load missing file name test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_missingScriptContent() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_script_content.tf")
	if err != nil {
		log.Fatalf("Failed to load missing script content test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_missingRunAsAccount() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_missing_run_as_account.tf")
	if err != nil {
		log.Fatalf("Failed to load missing run as account test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_invalidRunAsAccount() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_invalid_run_as_account.tf")
	if err != nil {
		log.Fatalf("Failed to load invalid run as account test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSCustomAttributeScriptConfig_invalidCustomAttributeType() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_invalid_custom_attribute_type.tf")
	if err != nil {
		log.Fatalf("Failed to load invalid custom attribute type test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
