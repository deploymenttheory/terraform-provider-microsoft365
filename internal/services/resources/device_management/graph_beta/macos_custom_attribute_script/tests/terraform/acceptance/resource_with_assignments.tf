resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "test_assignments" {
  display_name          = "Test macOS Custom Attribute Script with Assignments"
  description           = ""
  custom_attribute_type = "string"
  file_name             = "test_with_assignments.sh"
  script_content        = "#!/bin/bash\necho 'Script with assignments'\nexit 0"
  run_as_account        = "system"

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]
}