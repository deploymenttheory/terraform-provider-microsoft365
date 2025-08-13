resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "test_assignments" {
  display_name           = "Acceptance - Windows Feature Update Profile with Assignments"
  feature_update_version = "Windows 11, version 23H2"

  install_feature_updates_optional                         = false
  install_latest_windows10_on_windows11_ineligible_device = false

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


