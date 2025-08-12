package graphBetaMacOSCustomAttributeScript_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMacOSCustomAttributeScriptResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccMacOSCustomAttributeScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "display_name", "Test Acceptance macOS Custom Attribute Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "custom_attribute_name", "AcceptanceTestAttribute"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "custom_attribute_type", "String"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "file_name", "test_acceptance.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "run_as_account", "system"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "script_content"),
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
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "custom_attribute_name", "UpdatedAttribute"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "file_name", "test_acceptance_updated.sh"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "run_as_account", "user"),
				),
			},
			// Update back to minimal configuration
			{
				Config: testAccMacOSCustomAttributeScriptConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "display_name", "Test Acceptance macOS Custom Attribute Script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test", "run_as_account", "system"),
				),
			},
		},
	})
}

func TestAccMacOSCustomAttributeScriptResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with assignments
			{
				Config: testAccMacOSCustomAttributeScriptConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "display_name", "Test macOS Custom Attribute Script with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_custom_attribute_script.test_assignments", "assignments.0.type", "groupAssignmentTarget"),
				),
			},
		},
	})
}

func TestAccMacOSCustomAttributeScriptResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingDisplayName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccMacOSCustomAttributeScriptConfig_missingCustomAttributeName(),
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

func testAccMacOSCustomAttributeScriptConfig_minimal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Acceptance macOS Custom Attribute Script"
  custom_attribute_name = "AcceptanceTestAttribute"
  custom_attribute_type = "string"
  file_name             = "test_acceptance.sh"
  script_content        = "#!/bin/bash\necho 'Acceptance test value'\nexit 0"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Acceptance macOS Custom Attribute Script - Updated"
  description           = "Updated description for acceptance testing"
  custom_attribute_name = "UpdatedAttribute"
  custom_attribute_type = "string"
  file_name             = "test_acceptance_updated.sh"
  script_content        = "#!/bin/bash\necho 'Updated acceptance test value'\ndate\nexit 0"
  run_as_account        = "user"
  role_scope_tag_ids    = ["0", "1"]
}
`
}

func testAccMacOSCustomAttributeScriptConfig_withAssignments() string {
	return fmt.Sprintf(`
data "azuread_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_assignments" {
  display_name          = "Test macOS Custom Attribute Script with Assignments"
  custom_attribute_name = "AssignmentTestAttribute"
  custom_attribute_type = "string"
  file_name             = "test_with_assignments.sh"
  script_content        = "#!/bin/bash\necho 'Script with assignments'\nexit 0"
  run_as_account        = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.azuread_group.test_group.object_id
    }
  ]
}
`)
}

func testAccMacOSCustomAttributeScriptConfig_missingDisplayName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "string"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_missingCustomAttributeName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_type = "string"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_missingCustomAttributeType() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_missingFileName() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "string"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_missingScriptContent() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "string"
  file_name             = "test.sh"
  run_as_account        = "system"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_missingRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "string"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_invalidRunAsAccount() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "string"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "invalid"
}
`
}

func testAccMacOSCustomAttributeScriptConfig_invalidCustomAttributeType() string {
	return `
resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test" {
  display_name          = "Test Script"
  custom_attribute_name = "TestAttribute"
  custom_attribute_type = "Invalid"
  file_name             = "test.sh"
  script_content        = "#!/bin/bash\necho 'test'\nexit 0"
  run_as_account        = "system"
}
`
}
