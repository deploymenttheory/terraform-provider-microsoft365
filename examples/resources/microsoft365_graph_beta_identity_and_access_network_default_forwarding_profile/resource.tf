# Enable the built-in Internet Access profile for one group on Windows and macOS.
resource "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile" "internet" {
  traffic_forwarding_type = "internet"
  state                   = "enabled"
  client_fallback_action  = "bypass"

  assignment_rules = {
    membership = "enumerated"
    group_ids  = ["00000000-0000-0000-0000-000000000001"]
    platforms  = ["windows", "macOS"]
  }
}

# Turn on home network support for everyone assigned to the default Private
# Access profile, leaving its assignments unchanged.
resource "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile" "private" {
  traffic_forwarding_type = "private"
  state                   = "enabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }
}
