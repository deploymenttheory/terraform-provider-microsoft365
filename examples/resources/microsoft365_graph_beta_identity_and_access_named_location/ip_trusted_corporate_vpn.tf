# Corporate VPN Endpoints
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_corporate_vpn" {
  display_name = "Trusted - Corporate VPN"
  is_trusted   = true

  ipv4_ranges = [
    "198.51.100.0/24", # Example: VPN endpoint pool
  ]

  ipv6_ranges = [
    "2001:db8:5678::/48", # Example: VPN IPv6 pool
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

