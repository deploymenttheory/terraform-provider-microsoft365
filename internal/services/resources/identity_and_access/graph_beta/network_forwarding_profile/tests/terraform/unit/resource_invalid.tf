resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "test" {
  name                    = "unit-invalid"
  traffic_forwarding_type = "internet"
  priority                = 250
  state                   = "disabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }

  assignment_rules = {
    membership = "enumerated"
    platforms  = ["windows"]
  }
}
