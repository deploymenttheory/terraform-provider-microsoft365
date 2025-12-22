resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_008" {
  display_name                                            = "unit-test-windows-feature-update-policy-008-assignments-lifecycle"
  description                                             = "Maximal assignments lifecycle test"
  feature_update_version                                  = "Windows 11, version 24H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000003"
    }
  ]
}
