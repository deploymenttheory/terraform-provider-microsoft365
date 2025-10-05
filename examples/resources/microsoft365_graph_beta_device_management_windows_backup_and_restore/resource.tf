resource "microsoft365_graph_beta_device_management_windows_backup_and_restore" "minimal" {
  state = "enabled" # "disabled" or "notConfigured"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
