# ==============================================================================
# IP-based locations marked as is_trusted=true will be automatically included
# in the "AllTrusted" built-in location reference in Conditional Access policies.
# Note: Country-based locations do NOT support is_trusted attribute and must be
# referenced explicitly by ID in CA policies if needed.
# ==============================================================================

# Corporate Headquarters Network
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_corporate_hq" {
  display_name = "Trusted - Corporate Headquarters"
  is_trusted   = true

  ipv4_ranges = [
    "203.0.113.0/24", # Example: HQ public IP range
    "203.0.114.0/24", # Example: HQ additional subnet
  ]

  ipv6_ranges = [
    "2001:db8:1234::/48", # Example: HQ IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

