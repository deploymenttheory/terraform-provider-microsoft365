# ==============================================================================
# These locations require additional security controls like device compliance
# but are not outright blocked. NOT marked as trusted.
# ==============================================================================

# Partner/Vendor Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_partner_networks" {
  display_name = "Semi-Trusted - Partner Networks"

  ipv4_ranges = [
    "198.18.0.0/24", # Example: Partner A network
    "198.18.1.0/24", # Example: Partner B network
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

