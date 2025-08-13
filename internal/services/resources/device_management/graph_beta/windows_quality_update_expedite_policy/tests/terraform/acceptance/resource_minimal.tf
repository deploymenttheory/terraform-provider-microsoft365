resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test" {
  display_name = "Acceptance - Windows Quality Update Expedite Policy"

  expedited_update_settings = {
    quality_update_release   = "2025-04-08T00:00:00Z"
    days_until_forced_reboot = 1
  }
}


