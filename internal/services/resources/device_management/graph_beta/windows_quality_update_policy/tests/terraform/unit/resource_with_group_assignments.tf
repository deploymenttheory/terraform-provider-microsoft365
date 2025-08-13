resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "group_assignments" {
  display_name = "Test Group Assignments Windows Quality Update Policy - Unique"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    }
  ]
}


