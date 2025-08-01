package graphBetaWindowsRemediationScript_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWindowsRemediationScriptResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccWindowsRemediationScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "display_name", "Test Acceptance Windows Remediation Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test", "run_as_account", "system"),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccWindowsRemediationScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "display_name", "Test Windows Remediation Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_remediation_script.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccWindowsRemediationScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccWindowsRemediationScriptConfig_invalidRunAsAccount(),
				ExpectError: regexp.MustCompile("Attribute run_as_account value must be one of"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("M365_TENANT_ID") == "" {
		t.Skip("M365_TENANT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_ID") == "" {
		t.Skip("M365_CLIENT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_SECRET") == "" {
		t.Skip("M365_CLIENT_SECRET must be set for acceptance tests")
	}
}

func testAccWindowsRemediationScriptConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Acceptance Windows Remediation Script"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"
}
`
}

func testAccWindowsRemediationScriptConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name                = "Test Acceptance Windows Remediation Script - Updated"
  description                 = "Updated description for acceptance testing"
  publisher                   = "Terraform Provider Test Suite"
  run_as_account             = "user"
  run_as_32_bit              = true
  enforce_signature_check    = true
  detection_script_content   = <<-EOT
    # Comprehensive detection script for acceptance testing
    $computerName = $env:COMPUTERNAME
    Write-Host "Computer: $computerName"
    
    # Check for specific condition
    if (Test-Path "C:\temp\marker.txt") {
        Write-Host "Marker file found - issue detected"
        exit 1
    } else {
        Write-Host "No issues detected"
        exit 0
    }
  EOT
  
  remediation_script_content = <<-EOT
    # Comprehensive remediation script for acceptance testing
    $logPath = "C:\temp\remediation.log"
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    # Create directory if it doesn't exist
    if (!(Test-Path "C:\temp")) {
        New-Item -ItemType Directory -Path "C:\temp" -Force
    }
    
    # Log the remediation action
    Add-Content -Path $logPath -Value "$timestamp - Remediation started"
    
    # Remove the marker file
    if (Test-Path "C:\temp\marker.txt") {
        Remove-Item "C:\temp\marker.txt" -Force
        Add-Content -Path $logPath -Value "$timestamp - Marker file removed"
    }
    
    Add-Content -Path $logPath -Value "$timestamp - Remediation completed"
    Write-Host "Remediation completed successfully"
    exit 0
  EOT
  
  detection_script_parameters = [
    {
      name                                    = "CheckPath"
      description                            = "Path to check for the marker file"
      is_required                            = true
      apply_default_value_when_not_assigned = false
    }
  ]
  
  role_scope_tag_ids = ["0", "1"]
}
`
}

func testAccWindowsRemediationScriptConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_assignments" {
  display_name                = "Test Windows Remediation Script with Assignments"
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with assignments\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with assignments\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`)
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
