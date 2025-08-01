package graphBetaMacOSPlatformScript_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMacOSPlatformScriptResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccMacOSPlatformScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test Acceptance macOS Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test_acceptance.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_platform_script.test", "script_content"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_macos_platform_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccMacOSPlatformScriptConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test Acceptance macOS Platform Script - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "file_name", "test_acceptance_updated.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "user"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "block_execution_notifications", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "execution_frequency", "PT1H"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "retry_count", "2"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccMacOSPlatformScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_platform_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "display_name", "Test Acceptance macOS Platform Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test", "run_as_account", "system"),
				),
			},
		},
	})
}

func TestAccMacOSPlatformScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccMacOSPlatformScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_platform_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test_assignments", "display_name", "Test macOS Platform Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_platform_script.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccMacOSPlatformScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSPlatformScriptConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSPlatformScriptConfig_missingFileName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSPlatformScriptConfig_missingScriptContent(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSPlatformScriptConfig_missingRunAsAccount(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccMacOSPlatformScriptResource_InvalidValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSPlatformScriptConfig_invalidRunAsAccount(),
				ExpectError: regexp.MustCompile("Attribute run_as_account value must be one of"),
			},
			{
				Config:      testAccMacOSPlatformScriptConfig_invalidExecutionFrequency(),
				ExpectError: regexp.MustCompile("must be a valid ISO 8601 duration"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("ARM_TENANT_ID") == "" {
		t.Skip("ARM_TENANT_ID must be set for acceptance tests")
	}
	if os.Getenv("ARM_CLIENT_ID") == "" {
		t.Skip("ARM_CLIENT_ID must be set for acceptance tests")
	}
	if os.Getenv("ARM_CLIENT_SECRET") == "" {
		t.Skip("ARM_CLIENT_SECRET must be set for acceptance tests")
	}
}

func testAccMacOSPlatformScriptConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test Acceptance macOS Platform Script"
  file_name       = "test_acceptance.sh"
  script_content  = "#!/bin/bash\necho 'Acceptance test script'\nexit 0"
  run_as_account  = "system"
}
`
}

func testAccMacOSPlatformScriptConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name                   = "Test Acceptance macOS Platform Script - Updated"
  description                    = "Updated description for acceptance testing"
  file_name                      = "test_acceptance_updated.sh"
  script_content                 = "#!/bin/bash\necho 'Updated acceptance test script'\ndate\nexit 0"
  run_as_account                 = "user"
  role_scope_tag_ids            = ["0", "1"]
  block_execution_notifications = true
  execution_frequency           = "PT1H"
  retry_count                   = 2
}
`
}

func testAccMacOSPlatformScriptConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_macos_platform_script" "test_assignments" {
  display_name    = "Test macOS Platform Script with Assignments"
  file_name       = "test_with_assignments.sh"
  script_content  = "#!/bin/bash\necho 'Script with assignments'\nexit 0"
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

func testAccMacOSPlatformScriptConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  file_name       = "test.sh"
  script_content  = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account  = "system"
}
`
}

func testAccMacOSPlatformScriptConfig_missingFileName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test Script"
  script_content  = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account  = "system"
}
`
}

func testAccMacOSPlatformScriptConfig_missingScriptContent() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.sh"
  run_as_account  = "system"
}
`
}

func testAccMacOSPlatformScriptConfig_missingRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.sh"
  script_content  = "#!/bin/bash\necho 'test'\nexit 0"
}
`
}

func testAccMacOSPlatformScriptConfig_invalidRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name    = "Test Script"
  file_name       = "test.sh"
  script_content  = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account  = "invalid"
}
`
}

func testAccMacOSPlatformScriptConfig_invalidExecutionFrequency() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_platform_script" "test" {
  display_name        = "Test Script"
  file_name           = "test.sh"
  script_content      = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account      = "system"
  execution_frequency = "invalid_duration"
}
`
}