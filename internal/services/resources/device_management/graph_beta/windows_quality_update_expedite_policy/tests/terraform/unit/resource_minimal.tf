resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "minimal" {
  display_name = "Test Minimal Windows Quality Update Expedite Policy - Unique"

  expedited_update_settings = {
    quality_update_release   = "2025-04-08T00:00:00Z"
    days_until_forced_reboot = 1
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


