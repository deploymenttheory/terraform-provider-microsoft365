resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "test" {
  name                    = "unit-home-network-updated"
  traffic_forwarding_type = "private"
  priority                = 152
  state                   = "enabled"
  client_fallback_action  = "block"

  local_network_settings = {
    allow_local_subnet_bypass = false
  }

  forwarding_source_matching = {
    session_type = "agent"
  }

  assignment_rules = {
    membership = "all"
    platforms  = ["windows", "android", "iOS"]
  }
}
