resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "maximal" {
  display_name                                 = "Test Maximal Windows Feature Update Profile - Unique"
  description                                  = "Maximal Windows Feature Update Profile for testing with all features"
  feature_update_version                       = "Windows 11, version 23H2"
  install_feature_updates_optional             = true
  install_latest_windows10_on_windows11_ineligible_device = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2025-04-01T00:00:00Z"
    offer_end_date_time_in_utc   = "2025-05-01T00:00:00Z"
    offer_interval_in_days       = 7
  }

  role_scope_tag_ids = ["0", "1"]
}


