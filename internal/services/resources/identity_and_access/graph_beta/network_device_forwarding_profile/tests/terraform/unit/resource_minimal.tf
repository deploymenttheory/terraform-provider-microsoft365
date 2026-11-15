resource "microsoft365_graph_beta_identity_and_access_network_device_forwarding_profile" "test" {
  name     = "unit-minimal"
  priority = 150
  state    = "disabled"
}
