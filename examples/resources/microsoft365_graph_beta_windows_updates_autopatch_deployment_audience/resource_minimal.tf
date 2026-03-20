# Minimal example — creates an empty deployment audience with no members or exclusions.
# This is the singleton autopatch deployment audience resource for the tenant.

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
  timeouts = {
    create = "5m"
    read   = "5m"
    delete = "10m"
  }
}
