# Custom Private Access profile with home network support. Users assigned to
# this profile can turn on "Prefer Local Network" in the Global Secure Access
# client to reach devices on their directly connected local subnet.
resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "home_network" {
  name                    = "Private Access - remote workers"
  description             = "Private Access with home network support"
  traffic_forwarding_type = "private"
  priority                = 150
  state                   = "enabled"
  client_fallback_action  = "bypass"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }

  forwarding_source_matching = {
    session_type = "any"
  }

  assignment_rules = {
    membership = "enumerated"
    group_ids  = ["00000000-0000-0000-0000-000000000001"]
    platforms  = ["windows", "macOS"]
  }
}
