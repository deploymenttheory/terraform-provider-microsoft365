resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "test_24h2" {
  display_name                                            = "Acceptance - Windows 11 24H2 Feature Update Profile"
  description                                             = "Acceptance test for Windows 11 24H2 feature updates"
  feature_update_version                                  = "Windows 11, version 24H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2029-08-01T00:00:00Z"
  }
}