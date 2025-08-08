resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_assignments" {
  display_name                = "Test Windows Remediation Script with Assignments"
  description                 = ""
  publisher                   = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with assignments\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with assignments\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = false
      }
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      daily_schedule = {
        interval = 1
        time     = "14:00:00"
        use_utc  = false
      }
    },
    {
      type = "allLicensedUsersAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "12:00:00"
        use_utc  = false
      }
    },
    {
      type = "allDevicesAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "18:00:00"
        use_utc  = false
      }
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    }
  ]

  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}