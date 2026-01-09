# Scenario 5: Step 2 - Resource with single assignment (base resource unchanged)
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_minimal_assignment" {
  display_name   = "unit-test-add-minimal-assignment"
  file_name      = "test_minimal_assignment.sh"
  script_content = "#!/bin/bash\necho 'Test'\nexit 0"
  run_as_account = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
