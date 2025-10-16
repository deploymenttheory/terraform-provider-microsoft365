resource "microsoft365_graph_beta_identity_and_access_attribute_set" "test" {
  id = "Marketing"

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
