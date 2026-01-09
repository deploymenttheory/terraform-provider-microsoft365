# Scenario 6: Step 2 - Resource with all 4 assignment types
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_maximal_assignments" {
  display_name   = "unit-test-add-maximal-assignments"
  file_name      = "test_maximal_assignments.sh"
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
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
