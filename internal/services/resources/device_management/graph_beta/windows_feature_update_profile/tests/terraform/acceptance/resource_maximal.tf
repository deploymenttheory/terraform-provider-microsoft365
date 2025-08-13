resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example" {
  display_name                                            = "Windows 11 24H2 Deployment x"
  description                                             = "Feature update profile for Windows 11 24H2"
  feature_update_version                                  = "Windows 11, version 24H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // rollout_settings = Make update available gradually
  rollout_settings = {
    offer_start_date_time_in_utc = "2029-05-01T00:00:00Z"
    offer_end_date_time_in_utc   = "2029-06-02T23:59:59Z"
    offer_interval_in_days       = 7
  }
}


