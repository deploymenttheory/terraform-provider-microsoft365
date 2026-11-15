resource "microsoft365_graph_beta_identity_and_access_network_remote_network" "test" {
  name                   = "unit-test-remote-network"
  region                 = "eastUS"
  forwarding_profile_ids = ["22222222-2222-2222-2222-222222222222"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
