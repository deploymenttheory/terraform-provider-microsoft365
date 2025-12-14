# ==============================================================================
# Restrictive locations for specific accounts (service accounts, regional users)
# These define the ONLY locations where certain accounts can sign in from.
# ==============================================================================

# Service Account Source IPs
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_service_account_sources" {
  display_name = "Allowed - Service Account Sources"
  is_trusted   = true

  ipv4_ranges = [
    "10.100.0.10/32", # Example: Build server
    "10.100.0.11/32", # Example: Automation server
    "10.100.0.12/32", # Example: Monitoring server
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

