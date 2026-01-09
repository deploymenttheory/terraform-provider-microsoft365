# Scenario 7: Step 1 - Resource with single assignment
resource "microsoft365_graph_beta_device_management_macos_platform_script" "assignment_update" {
  display_name   = "unit-test-assignment-update"
  file_name      = "test_assignment_update.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "55555555-5555-5555-5555-555555555555"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
