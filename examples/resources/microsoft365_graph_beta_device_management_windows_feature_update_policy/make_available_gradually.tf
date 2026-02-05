# ==============================================================================
# Make Update Available Gradually (Phased Rollout)
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update gradually
# over time. The update will be offered to devices in stages, starting on the
# offer_start_date and completing by the offer_end_date, with devices being
# offered the update at the specified interval.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "gradual" {
  display_name                                            = "Windows 11 25H2 - Gradual Deployment"
  description                                             = "Phased rollout of Windows 11 25H2 feature update"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
    offer_end_date_time_in_utc   = "2030-01-14T00:00:00Z"
    offer_interval_in_days       = 1
  }

  role_scope_tag_ids = ["0", "1"]
}
