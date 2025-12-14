# Regional Office Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_regional_offices" {
  display_name = "Trusted - Regional Offices"
  is_trusted   = true

  ipv4_ranges = [
    "192.0.2.0/24",    # Example: EMEA Office
    "198.51.100.0/24", # Example: APAC Office
    "203.0.113.0/24",  # Example: Americas Office
  ]

  ipv6_ranges = [
    "2001:db8:abcd::/48", # Example: Regional IPv6
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

