package graphBetaWindowsPlatformScript_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsPlatformScriptResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccWindowsPlatformScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Acceptance Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "file_name", "test_acceptance.ps1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.test", "script_content"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_windows_platform_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccWindowsPlatformScriptConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Acceptance Windows Platform Script - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "file_name", "test_acceptance_updated.ps1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "enforce_signature_check", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_32_bit", "false"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccWindowsPlatformScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "display_name", "Test Acceptance Windows Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test", "run_as_account", "system"),
				),
			},
		},
	})
}

func TestAccWindowsPlatformScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccWindowsPlatformScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_platform_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test_assignments", "display_name", "Test Windows Platform Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_platform_script.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccWindowsPlatformScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsPlatformScriptConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsPlatformScriptConfig_missingFileName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsPlatformScriptConfig_missingScriptContent(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccWindowsPlatformScriptConfig_missingRunAsAccount(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccWindowsPlatformScriptResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsPlatformScriptConfig_invalidRunAsAccount(),
				ExpectError: regexp.MustCompile("Attribute run_as_account value must be one of"),
			},
		},
	})
}

func testAccWindowsPlatformScriptConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name    = "Test Acceptance Windows Platform Script"
  file_name       = "test_acceptance.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'Acceptance test script'\nExit 0"
  run_as_account  = "system"
}
`
}

func testAccWindowsPlatformScriptConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name              = "Test Acceptance Windows Platform Script - Updated"
  description               = "Updated description for acceptance testing"
  file_name                 = "test_acceptance_updated.ps1"
  script_content            = "# PowerShell Script\nWrite-Host 'Updated acceptance test script'\nGet-Date\nExit 0"
  run_as_account            = "user"
  role_scope_tag_ids       = ["0", "1"]
  enforce_signature_check   = true
  run_as_32_bit            = false
}
`
}

func testAccWindowsPlatformScriptConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_assignments" {
  display_name    = "Test Windows Platform Script with Assignments"
  file_name       = "test_with_assignments.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'Script with assignments'\nExit 0"
  run_as_account  = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`)
}

func testAccWindowsPlatformScriptConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  file_name       = "test.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'test'\nExit 0"
  run_as_account  = "system"
}
`
}

func testAccWindowsPlatformScriptConfig_missingFileName() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name    = "Test Script"
  script_content  = "# PowerShell Script\nWrite-Host 'test'\nExit 0"
  run_as_account  = "system"
}
`
}

func testAccWindowsPlatformScriptConfig_missingScriptContent() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.ps1"
  run_as_account  = "system"
}
`
}

func testAccWindowsPlatformScriptConfig_missingRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'test'\nExit 0"
}
`
}

func testAccWindowsPlatformScriptConfig_invalidRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.ps1"
  script_content  = "# PowerShell Script\nWrite-Host 'test'\nExit 0"
  run_as_account  = "invalid"
}
`
}
