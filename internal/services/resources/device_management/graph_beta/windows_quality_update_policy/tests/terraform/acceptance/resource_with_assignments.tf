resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_assignments" {
  display_name = "Acceptance - Windows Quality Update Policy with Assignments"
  hotpatch_enabled = false

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


