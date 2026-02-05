# ==============================================================================
# Make Update Available On a Specific Date
# ==============================================================================
# This example demonstrates how to schedule a Windows feature update to become
# available on a specific date. The update will not be offered to devices until
# the specified start date.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "scheduled" {
  display_name                                            = "Windows 11 25H2 - Scheduled Deployment"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
  }
}
