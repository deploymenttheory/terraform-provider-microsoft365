resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_assignments" {
  display_name = "Acceptance - Windows Quality Update Policy with Assignments"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_directory_group.test_group_include.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_directory_group.test_group_exclude.id
    }
  ]
}


