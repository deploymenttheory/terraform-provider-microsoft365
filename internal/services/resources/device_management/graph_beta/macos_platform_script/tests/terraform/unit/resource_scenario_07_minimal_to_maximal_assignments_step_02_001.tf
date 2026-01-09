# Scenario 7: Step 2 - Resource with all 4 assignment types
resource "microsoft365_graph_beta_device_management_macos_platform_script" "assignment_update" {
  display_name   = "unit-test-assignment-update"
  file_name      = "test_assignment_update.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "66666666-6666-6666-6666-666666666666"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "77777777-7777-7777-7777-777777777777"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
