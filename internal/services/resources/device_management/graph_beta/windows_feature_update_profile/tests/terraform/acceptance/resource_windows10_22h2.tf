resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "test_win10_22h2" {
  display_name                                            = "Acceptance - Windows 10 22H2 Feature Update Profile"
  description                                             = "Acceptance test for Windows 10 22H2 feature updates"
  feature_update_version                                  = "Windows 10, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = false
  install_feature_updates_optional                        = false
}