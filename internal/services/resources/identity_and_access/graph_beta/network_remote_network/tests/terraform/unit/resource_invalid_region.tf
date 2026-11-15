resource "microsoft365_graph_beta_identity_and_access_network_remote_network" "test" {
  name                   = "unit-test-remote-network"
  region                 = "moonBase"
  forwarding_profile_ids = ["11111111-1111-1111-1111-111111111111"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
