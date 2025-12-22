resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_005" {
  display_name                                            = "unit-test-windows-feature-update-policy-005-assignments-minimal"
  feature_update_version                                  = "Windows 11, version 23H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    }
  ]
}
