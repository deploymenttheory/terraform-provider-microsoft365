resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "specific_date_rollout" {
  display_name                                 = "Test Specific Date Rollout Windows Feature Update Profile - Unique"
  description                                  = "Windows Feature Update Profile for testing specific date rollout scenarios"
  feature_update_version                       = "Windows 11, version 24H2"
  install_feature_updates_optional             = true
  install_latest_windows10_on_windows11_ineligible_device = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2029-08-01T00:00:00Z"
  }

  role_scope_tag_ids = ["8", "9"]
}