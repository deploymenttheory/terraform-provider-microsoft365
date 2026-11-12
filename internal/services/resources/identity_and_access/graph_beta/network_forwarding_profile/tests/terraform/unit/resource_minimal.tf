resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "test" {
  name                    = "unit-minimal"
  traffic_forwarding_type = "private"
  priority                = 150
  state                   = "disabled"
}
