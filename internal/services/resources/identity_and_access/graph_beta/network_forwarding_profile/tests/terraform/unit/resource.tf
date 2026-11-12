resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "test" {
  name                    = "unit-home-network"
  description             = "unit test"
  traffic_forwarding_type = "private"
  priority                = 150
  state                   = "disabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }

  forwarding_source_matching = {
    session_type = "any"
  }

  assignment_rules = {
    membership = "enumerated"
    group_ids  = ["a6ed41d3-45ac-4e12-8601-bdb56216e84e"]
    platforms  = ["windows", "macOS"]
  }
}
