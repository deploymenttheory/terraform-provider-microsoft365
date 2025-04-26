resource "microsoft365_graph_beta_device_and_app_management_windows_feature_update_profile" "example" {
  display_name                                            = "Windows 11 22H2 Deployment"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = false
  role_scope_tag_ids                                      = ["8", "9"]

  # Rollout settings
  rollout_settings = {
    offer_start_date_time_in_utc = "2025-05-01T00:00:00Z"
    offer_end_date_time_in_utc   = "2025-06-30T23:59:59Z"
    offer_interval_in_days       = 7
  }

  # Optional - Timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}