resource "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile" "test" {
  traffic_forwarding_type = "internet"
  state                   = "disabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }
}
