resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_assignments" {
  display_name               = "Test Windows Remediation Script with Assignments"
  description                = ""
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Detection script with assignments\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Remediation script with assignments\nWrite-Host 'Remediation complete'\nexit 0"
  role_scope_tag_ids         = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

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
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      run_once_schedule = {
        date    = "2030-12-31"
        time    = "23:59:59"
        use_utc = false
      }
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      hourly_schedule = {
        interval = 1
      }
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "exclude"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      daily_schedule = {
        interval = 1
        time     = "02:00:00"
        use_utc  = false
      }
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_3.id
      filter_type = "exclude"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
    }
  ]


  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}