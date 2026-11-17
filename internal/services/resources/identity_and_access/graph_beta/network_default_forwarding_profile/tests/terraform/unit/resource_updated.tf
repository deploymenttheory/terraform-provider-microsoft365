resource "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile" "test" {
  traffic_forwarding_type = "internet"
  state                   = "enabled"
  client_fallback_action  = "block"

  forwarding_source_matching = {
    session_type = "agent"
  }

  assignment_rules = {
    membership = "none"
    platforms  = ["windows"]
  }
}
