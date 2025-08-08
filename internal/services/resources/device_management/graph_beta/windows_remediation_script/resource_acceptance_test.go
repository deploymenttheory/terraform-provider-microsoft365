package graphBetaWindowsRemediationScript_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsRemediationScriptResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsRemediationScriptDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccWindowsRemediationScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Acceptance Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "publisher", "Terraform Provider Test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_32_bit", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "enforce_signature_check", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test", "detection_script_content"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test", "remediation_script_content"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_remediation_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccWindowsRemediationScriptConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Acceptance Windows Remediation Script - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "publisher", "Terraform Provider Test Suite"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_32_bit", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsRemediationScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsRemediationScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "display_name", "Test Windows Remediation Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.#", "5"),
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckWindowsRemediationScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsRemediationScriptConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsRemediationScriptConfig_missingPublisher(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsRemediationScriptConfig_missingRunAsAccount(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsRemediationScriptConfig_missingDetectionScript(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsRemediationScriptConfig_missingRemediationScript(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccWindowsRemediationScriptResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsRemediationScriptDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsRemediationScriptConfig_invalidRunAsAccount(),
				ExpectError: regexp.MustCompile("Attribute run_as_account value must be one of"),
			},
		},
	})
}

func testAccWindowsRemediationScriptConfig_minimal() string {
	config := mocks.LoadLocalTerraformConfig("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccWindowsRemediationScriptConfig_maximal() string {
	config := mocks.LoadLocalTerraformConfig("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccWindowsRemediationScriptConfig_withAssignments() string {
	groups := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	roleScopeTags := mocks.LoadCentralizedTerraformConfig("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	config := mocks.LoadLocalTerraformConfig("resource_with_assignments.tf")
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + config)
}

func testAccWindowsRemediationScriptConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_missingPublisher() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Script"
  run_as_account             = "system"
  detection_script_content   = "# Detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_missingRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Script"
  publisher                   = "Terraform Provider Test"
  detection_script_content   = "# Detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_missingDetectionScript() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Script"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  remediation_script_content = "# Remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_missingRemediationScript() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Script"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script\nWrite-Host 'Detection complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_invalidRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Script"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "invalid"
  detection_script_content   = "# Detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

// testAccCheckWindowsRemediationScriptDestroy verifies that Windows remediation scripts have been destroyed
func testAccCheckWindowsRemediationScriptDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_remediation_script" {
			continue
		}

		// Attempt to get the Windows remediation script by ID
		_, err := graphClient.
			DeviceManagement().
			DeviceHealthScripts().
			ByDeviceHealthScriptId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if Windows remediation script %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("Windows remediation script %s still exists", rs.Primary.ID)
	}

	return nil
}
