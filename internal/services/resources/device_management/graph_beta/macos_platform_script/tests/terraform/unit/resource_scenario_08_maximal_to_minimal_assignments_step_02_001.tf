# Scenario 8: Step 2 - Resource with single assignment
resource "microsoft365_graph_beta_device_management_macos_platform_script" "assignment_downgrade" {
  display_name   = "unit-test-assignment-downgrade"
  file_name      = "test_assignment_downgrade.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
