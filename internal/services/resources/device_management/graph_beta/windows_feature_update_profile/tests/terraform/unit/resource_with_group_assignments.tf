resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "group_assignments" {
  display_name           = "Test Group Assignments Windows Feature Update Profile - Unique"
  feature_update_version = "Windows 11, version 23H2"

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


