resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "test" {
  display_name           = "Acceptance - Windows Feature Update Profile"
  feature_update_version = "Windows 11, version 23H2"

  install_feature_updates_optional                         = false
  install_latest_windows10_on_windows11_ineligible_device = false
}


