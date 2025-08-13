resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_assignments" {
  display_name = "Acceptance - Windows Quality Update Expedite Policy with Assignments"
  
  expedited_update_settings = {
      quality_update_release   = "2025-04-08T00:00:00Z"
      days_until_forced_reboot = 1
    }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]
}


