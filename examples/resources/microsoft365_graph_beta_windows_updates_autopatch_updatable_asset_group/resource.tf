# Minimal example — creates an empty updatable asset group with no device members.
# The group ID is assigned by the service and can be referenced by other resources.

resource "microsoft365_graph_beta_windows_updates_autopatch_updatable_asset_group" "example" {
  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}
