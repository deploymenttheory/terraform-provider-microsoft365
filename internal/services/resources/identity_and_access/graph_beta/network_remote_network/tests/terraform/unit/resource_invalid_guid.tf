resource "microsoft365_graph_beta_identity_and_access_network_remote_network" "test" {
  name                   = "unit-test-remote-network"
  region                 = "eastUS"
  forwarding_profile_ids = ["not-a-guid"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
