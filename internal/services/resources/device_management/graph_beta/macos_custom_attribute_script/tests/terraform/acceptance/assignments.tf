terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

data "microsoft365_graph_beta_groups_group" "test_group" {
  display_name = "Test Group"
}

resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_assignments" {
  display_name          = "Test macOS Custom Attribute Script with Assignments - Acceptance"
  custom_attribute_type = "string"
  file_name             = "test_assignment_acceptance.sh"
  script_content        = "#!/bin/bash\necho 'Assignment Acceptance Test Value'\nexit 0"
  run_as_account        = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_graph_beta_groups_group.test_group.id
    }
  ]
}