# Scenario 8: Step 1 - Resource with all 4 assignment types
resource "microsoft365_graph_beta_device_management_macos_platform_script" "assignment_downgrade" {
  display_name   = "unit-test-assignment-downgrade"
  file_name      = "test_assignment_downgrade.sh"
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
      group_id = "88888888-8888-8888-8888-888888888888"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "99999999-9999-9999-9999-999999999999"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
