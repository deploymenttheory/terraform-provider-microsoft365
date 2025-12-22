resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_004" {
  display_name                                            = "acc-test-windows-feature-update-policy-004-${random_string.test_suffix.result}"
  feature_update_version                                  = "Windows 11, version 23H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false
}
