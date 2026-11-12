resource "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile" "test" {
  traffic_forwarding_type = "internet"
  state                   = "disabled"
  client_fallback_action  = "bypass"
}
