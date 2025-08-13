resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "minimal" {
  display_name = "Test Minimal Windows Quality Update Policy - Unique"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


