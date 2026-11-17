resource "microsoft365_graph_beta_identity_and_access_network_device_forwarding_profile" "test" {
  name        = "unit-home-network"
  description = "unit test"
  priority    = 150
  state       = "disabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }

  assignment_rules = {
    membership = "enumerated"
    group_ids  = ["aaaaaaaa-1111-4111-8111-111111111111"]
    platforms  = ["windows", "macOS"]
  }
}
