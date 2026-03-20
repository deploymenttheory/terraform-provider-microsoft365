resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {
  timeouts = {
    create = "5m"
    read   = "5m"
    delete = "10m"
  }
}
