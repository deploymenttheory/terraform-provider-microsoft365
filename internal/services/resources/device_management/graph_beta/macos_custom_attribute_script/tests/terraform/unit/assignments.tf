terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
  }
}

data "microsoft365_graph_beta_groups_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_assignment" {
  display_name          = "Test macOS Custom Attribute Script - Assignment"
  custom_attribute_type = "string"
  file_name             = "test_assignment.sh"
  script_content        = "#!/bin/bash\necho 'Assignment Test Value'\nexit 0"
  run_as_account        = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_graph_beta_groups_group.test_group.id
    }
  ]
}