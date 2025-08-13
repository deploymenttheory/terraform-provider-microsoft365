resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "exclusion_assignment" {
  display_name           = "Test Exclusion Assignment Windows Feature Update Profile - Unique"
  feature_update_version = "Windows 11, version 23H2"

  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]
}


